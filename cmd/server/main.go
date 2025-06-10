package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/skip2/go-qrcode"
)

const defaultPort = "3000"

var sessions = make(map[string]string) // sessionId -> roomName
var lastRoom string

// New data structures for the employee queue system
var (
	employeeTickets      []string // Stores ticketIDs of employees in the queue
	employeeAssignments  = make(map[string]SessionDetails) // ticketID -> SessionDetails
	nextTicketID       int64 = 1
	queueMutex         sync.Mutex
	assignmentsMutex   sync.Mutex
	ticketIDMutex      sync.Mutex
)

type SessionDetails struct {
	SessionID          string `json:"sessionId"`
	RoomName           string `json:"roomName"`
	ClientURL          string `json:"clientUrl"` // URL for the client to join the session
	EmployeeRedirectURL string `json:"employeeRedirectUrl"` // URL for the employee to join the video call
}

func main() {
	// Main routes
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/empregado", serveEmployeePage)
	http.HandleFunc("/admin", serveAdminPage)
	http.HandleFunc("/medico", serveDoctorPage)
	http.HandleFunc("/cliente/", serveClientPage)
	http.HandleFunc("/debug", serveDebugPage)
	http.HandleFunc("/video", serveVideo)
	http.HandleFunc("/test", serveTestPage)
	http.HandleFunc("/comprehensive-test", serveComprehensiveTestPage)
	http.HandleFunc("/status", serveStatusDashboard)
	http.HandleFunc("/live-video-test", serveLiveVideoTest)
	http.HandleFunc("/last-room", serveLastRoom)

	// API routes
	http.HandleFunc("/api/create-session", createSession) // This will likely be used by admin/doctor flow
	http.HandleFunc("/api/token", generateToken)
	http.HandleFunc("/api/qr/", generateQRCode) // Used for client QR, admin might use a new one
	http.HandleFunc("/api/end-session", endSession)
	http.HandleFunc("/api/session-info/", getSessionInfo)

	// New Employee Queue API Endpoints
	http.HandleFunc("/api/employee/enqueue", handleEmployeeEnqueue)
	http.HandleFunc("/api/employee/queue-status", handleEmployeeQueueStatus)
	http.HandleFunc("/api/admin/employee-queue-entry-qr", handleAdminQueueEntryQR)
	http.HandleFunc("/api/admin/call-next-employee", handleAdminCallNextEmployee)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
     port = 3000
	}

	log.Printf("Server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func generateNewTicketID() string {
	ticketIDMutex.Lock()
	defer ticketIDMutex.Unlock()
	ticketID := fmt.Sprintf("ticket-%d", nextTicketID)
	nextTicketID++
	return ticketID
}

func handleEmployeeEnqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ticketID := generateNewTicketID()

	queueMutex.Lock()
	employeeTickets = append(employeeTickets, ticketID)
	position := len(employeeTickets)
	queueMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ticketId": ticketID,
		"status":   "queued",
		"position": position,
	})
}

func handleEmployeeQueueStatus(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticket")
	if ticketID == "" {
		http.Error(w, "Ticket ID required", http.StatusBadRequest)
		return
	}

	assignmentsMutex.Lock()
	assignment, assigned := employeeAssignments[ticketID]
	assignmentsMutex.Unlock()

	if assigned {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":         "assigned",
			"sessionDetails": assignment, // Send the whole assignment struct as expected by frontend
		})
		return
	}

	queueMutex.Lock()
	position := -1
	for i, t := range employeeTickets {
		if t == ticketID {
			position = i + 1
			break
		}
	}
	queueMutex.Unlock()

	if position != -1 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ticketId": ticketID,
			"status":   "queued",
			"position": position,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ticketId": ticketID,
			"status":   "unknown",
			"message":  "Ticket not found in active queue or assignments. It might have been processed or is invalid.",
		})
	}
}

func handleAdminQueueEntryQR(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	if r.Host != "" && !strings.HasPrefix(r.Host, "localhost") {
		scheme := "http"
		baseURL = fmt.Sprintf("%s://%s", scheme, r.Host)
	}

	queueEntryURL := fmt.Sprintf("%s/empregado?action=joinQueueViaQR", baseURL)

	png, err := qrcode.Encode(queueEntryURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Error generating QR code for queue entry", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func createSessionInternal(r *http.Request) (SessionDetails, error) {
	sessionId := fmt.Sprintf("session-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))
	roomName := fmt.Sprintf("room-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))

	sessions[sessionId] = roomName
	lastRoom = roomName

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	baseURL := fmt.Sprintf("http://localhost:%s", port)
	if r != nil && r.Host != "" && !strings.HasPrefix(r.Host, "localhost") {
		scheme := "http"
		baseURL = fmt.Sprintf("%s://%s", scheme, r.Host)
	}

	clientURL := fmt.Sprintf("%s/cliente/%s", baseURL, sessionId)
	employeeRedirectURL := fmt.Sprintf("%s/video?room=%s&sessionId=%s&role=employee", baseURL, roomName, sessionId)

	return SessionDetails{
		SessionID:          sessionId,
		RoomName:           roomName,
		ClientURL:          clientURL,
		EmployeeRedirectURL: employeeRedirectURL,
	}, nil
}

func handleAdminCallNextEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueMutex.Lock()
	if len(employeeTickets) == 0 {
		queueMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "queue_empty", "message": "No employees in the queue."})
		return
	}
	ticketID := employeeTickets[0]
	employeeTickets = employeeTickets[1:]
	queueMutex.Unlock()

	sessionDetails, err := createSessionInternal(r)
	if err != nil {
		http.Error(w, "Failed to create session for employee", http.StatusInternalServerError)
		return
	}

	assignmentsMutex.Lock()
	employeeAssignments[ticketID] = sessionDetails
	assignmentsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "employee_called",
		"ticketId":           ticketID,
		"sessionId":          sessionDetails.SessionID,
		"roomName":           sessionDetails.RoomName,
		"clientUrl":          sessionDetails.ClientURL,
		"adminVideoUrl":      fmt.Sprintf("%s/video?room=%s&role=admin", strings.TrimSuffix(sessionDetails.EmployeeRedirectURL, fmt.Sprintf("&sessionId=%s&role=employee", sessionDetails.SessionID)), sessionDetails.RoomName),
	})
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>QRKit Video Calling</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; }
        .nav-links { list-style: none; padding: 0; }
        .nav-links li { margin: 15px 0; }
        .nav-links a { display: block; padding: 15px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; text-align: center; }
        .nav-links a:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>QRKit Video Calling System</h1>        <ul class="nav-links">
            <li><a href="/empregado">👨‍💼 Employee Interface</a></li>
            <li><a href="/live-video-test">🎥 Live Video Test</a></li>
            <li><a href="/status">📊 Real-time Status Dashboard</a></li>
            <li><a href="/debug">🔧 Debug & Testing</a></li>
            <li><a href="/test">🧪 Complete Test Page</a></li>
            <li><a href="/comprehensive-test">🚀 Comprehensive Testing Suite</a></li>
        </ul>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func serveEmployeePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/empregado.html"))
	tmpl.Execute(w, nil)
}

func serveAdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, nil)
}

func serveClientPage(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Path[len("/cliente/"):]
	if sessionId == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/cliente.html"))
	tmpl.Execute(w, map[string]string{
		"SessionId": sessionId,
	})
}

func serveDebugPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/debug.html"))
	tmpl.Execute(w, nil)
}

func serveTestPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/test.html"))
	tmpl.Execute(w, nil)
}

func serveComprehensiveTestPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/comprehensive-test.html")
}

func serveStatusDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/status-dashboard.html")
}

func serveLiveVideoTest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/live-video-test.html")
}

func createSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionId := fmt.Sprintf("session-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))
	roomName := fmt.Sprintf("room-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))

	sessions[sessionId] = roomName
	lastRoom = roomName

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	baseURL := fmt.Sprintf("http://localhost:%s", port)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sessionId": "%s", "roomName": "%s", "clientUrl": "%s/cliente/%s"}`, sessionId, roomName, baseURL, sessionId)
}

func getSessionInfo(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Path[len("/api/session-info/"):]
	if sessionId == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	roomName, ok := sessions[sessionId]
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	if r.Host != "" && !strings.HasPrefix(r.Host, "localhost") {
		scheme := "http"
		baseURL = fmt.Sprintf("%s://%s", scheme, r.Host)
	}

	clientUrl := fmt.Sprintf("%s/cliente/%s", baseURL, sessionId)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sessionId": "%s", "roomName": "%s", "clientUrl": "%s"}`, sessionId, roomName, clientUrl)
}

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Path[len("/api/qr/"):]
	if sessionId == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	clientUrl := fmt.Sprintf("http://localhost:%s/cliente/%s", port, sessionId)

	png, err := qrcode.Encode(clientUrl, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "room and identity parameters required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	if apiKey == "" {
		apiKey = "devkey"
	}
	if apiSecret == "" {
		apiSecret = "secret"
	}

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.SetIdentity(identity).SetVideoGrant(grant)

	token, err := at.ToJWT()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

func endSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.FormValue("sessionId")
	if sessionID == "" {
		http.Error(w, "sessionId required", http.StatusBadRequest)
		return
	}

	delete(sessions, sessionID)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ended"}`)
}

func serveLastRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"room":"%s"}`, lastRoom)
}

func serveDoctorPage(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if room == "" {
		tmpl := template.Must(template.ParseFiles("templates/medico.html"))
		tmpl.Execute(w, nil)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/video.html"))
	tmpl.Execute(w, map[string]string{"Room": room})
}

func serveVideo(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	tmpl := template.Must(template.ParseFiles("templates/video.html"))
	tmpl.Execute(w, map[string]string{"Room": room})
}

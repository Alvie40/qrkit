package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/skip2/go-qrcode"
)

var sessions = make(map[string]string) // sessionId -> roomName
var lastRoom string

type QueueItem struct {
	SessionID string    `json:"sessionId"`
	RoomName  string    `json:"roomName"`
	Timestamp time.Time `json:"timestamp"`
}

var employeeQueue []QueueItem

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
	http.HandleFunc("/api/create-session", createSession)
	http.HandleFunc("/api/token", generateToken)
	http.HandleFunc("/api/qr/", generateQRCode)
	http.HandleFunc("/api/join-queue", joinQueue)
	http.HandleFunc("/api/next-session", nextSession)
	http.HandleFunc("/api/end-session", endSession)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Println("Server starting on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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
		port = "8080"
	}
	baseURL := fmt.Sprintf("http://localhost:%s", port)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sessionId": "%s", "roomName": "%s", "clientUrl": "%s/cliente/%s"}`, sessionId, roomName, baseURL, sessionId)
}

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Path[len("/api/qr/"):]
	if sessionId == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
	at.AddGrant(grant).SetIdentity(identity)

	token, err := at.ToJWT()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

func joinQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.FormValue("sessionId")
	if sessionID == "" {
		http.Error(w, "sessionId required", http.StatusBadRequest)
		return
	}
	roomName := r.FormValue("roomName")
	if roomName == "" {
		roomName = sessions[sessionID]
	}

	qi := QueueItem{SessionID: sessionID, RoomName: roomName, Timestamp: time.Now()}
	employeeQueue = append(employeeQueue, qi)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"queued"}`)
}

func nextSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(employeeQueue) == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"empty"}`)
		return
	}

	qi := employeeQueue[0]
	employeeQueue = employeeQueue[1:]

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sessionId":"%s","roomName":"%s","timestamp":"%s"}`,
		qi.SessionID, qi.RoomName, qi.Timestamp.Format(time.RFC3339))
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

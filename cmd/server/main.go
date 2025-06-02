package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/skip2/go-qrcode"
)

var lastRoom string

func main() {
	http.HandleFunc("/", serveQRPage)             // página principal
	http.HandleFunc("/generate", generateQR)      // rota que gera o QR Code
	http.HandleFunc("/video", serveVideo)         // exibe a sala da videoconferência
	http.HandleFunc("/last-room", serveLastRoom)  // novo endpoint para última sala
	http.HandleFunc("/empregado", serveEmpregado) // página do empregado
	http.HandleFunc("/medico", serveMedico)       // página do médico

	// servir arquivos estáticos se necessário
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// inicia servidor
	http.ListenAndServe(":8080", nil)
}

func serveQRPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/qrcode.html")
}

func generateQR(w http.ResponseWriter, r *http.Request) {
	// cria um nome de sala aleatório
	room := fmt.Sprintf("sala-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(99999))
	link := fmt.Sprintf("http://localhost:8080/video?room=%s", room)

	// gera QR code com o link
	png, _ := qrcode.Encode(link, qrcode.Medium, 256)
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func serveVideo(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if room != "" {
		lastRoom = room
	}
	tmpl := template.Must(template.ParseFiles("templates/video.html"))
	tmpl.Execute(w, map[string]string{"Room": room})
}

func serveLastRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"room": "%s"}`, lastRoom)
}

func serveEmpregado(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/empregado.html")
}

func serveMedico(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if room == "" {
		// Se não há sala, mostra página de espera
		http.ServeFile(w, r, "static/medico.html")
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/video.html"))
	tmpl.Execute(w, map[string]string{"Room": room})
}

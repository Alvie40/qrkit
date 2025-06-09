# QRKit – Video Call System with QR Code

A complete video call system that uses QR codes to easily connect employees and clients through WebRTC and LiveKit.

---

## 🚀 Features

* **Employee Interface:** Start video call sessions
* **QR Code Generation:** Unique QR codes per session
* **Client Interface:** Instant access via QR code or URL
* **WebRTC:** Real-time video communication
* **LiveKit Integration:** Robust WebRTC infrastructure
* **Status Dashboard (demo):** Frontend status/metrics example
* **Debug Tools:** Complete test & diagnostics suite

---

## 📋 Requirements

* Go 1.21+
* Docker & Docker Compose (optional for local/dev)
* Modern browser with WebRTC support
* Camera and microphone access

---

## 🏃‍♂️ Quick Start

### Option 1: Manual Run (Local Development)

```sh
# Compile the app
go build -o server ./cmd/server

# Run the server
./server
```

### Option 2: Docker Compose (Local Development)

```sh
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

The app connects to PostgreSQL and Redis using `DATABASE_URL` and `REDIS_URL` from docker-compose.

---

## 🚀 Deploy with Coolify (Production)

**Recommended for production!**
Coolify handles orchestration, HTTPS, and custom domains automatically.
Create a Docker Compose app in the Coolify dashboard and paste this file:

```yaml
version: "3.9"
services:
  qrkit:
    build: .
    # Do NOT expose any port here! Coolify handles reverse proxy.
    volumes:
      - ./static:/app/static
      - ./templates:/app/templates
    environment:
      - LIVEKIT_API_KEY=mykey
      - LIVEKIT_API_SECRET=mysecret
      - LIVEKIT_HOST=livekit:7880
      - DATABASE_URL=postgres://qrkit:qrkitpass@db:5432/qrkit?sslmode=disable
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
      - livekit

  livekit:
    image: livekit/livekit-server
    command: --config /livekit.yaml
    ports:
      - "7880:7880"
      - "7881:7881/udp"
      - "50000-50100:50000-50100/udp"
    volumes:
      - ./livekit-local.yaml:/livekit.yaml
    environment:
      - "LIVEKIT_KEYS=mykey:mysecret"

  db:
    image: postgres:15
    restart: always
    volumes:
      - db-data:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: qrkit
      POSTGRES_PASSWORD: qrkitpass
      POSTGRES_DB: qrkit

  redis:
    image: redis:7
    restart: always

volumes:
  db-data:
```

**In Coolify dashboard:**

* Choose `Application Type: Docker Compose`
* Set `Application Port: 8080` (the internal QRKit port)
* Coolify handles HTTPS/domain routing!

**For local testing:**
Temporarily add to the `qrkit` service:

```yaml
    ports:
      - "8080:8080"
```

Then access [http://localhost:8080](http://localhost:8080) in your browser.

---

## 🌐 Access

Once running, visit:

* **Homepage:** [http://YOUR-COOLIFY-DOMAIN](http://YOUR-COOLIFY-DOMAIN)
* **Employee Interface:** /empregado
* **Status Dashboard (demo):** /status
* **Debug Tools:** /debug
* **Comprehensive Test:** /comprehensive-test

---

## 📱 Usage

**For Employees:**

1. Go to the employee interface
2. Click "Create New Session"
3. Generate the QR code
4. Share the QR code or URL with the client
5. Wait for the client to connect
6. Start the video call

**For Clients:**

1. Scan the QR code with phone/tablet camera
2. Or open the URL provided
3. Allow camera and mic access
4. Wait for the video call to start

---

## 🏗️ Architecture

```
QRKit/
├── cmd/server/           # Main Go app
├── templates/            # HTML templates
├── static/               # Static assets (CSS, JS, HTML)
├── docker-compose.yml    # Docker config
├── livekit-local.yaml    # LiveKit config
└── README.md             # This documentation
```

**Components:**

* **Go Server:** REST API and page server
* **LiveKit:** WebRTC server
* **Frontend:** Responsive web UI
* **QR Generator:** Auto-generation of session QR codes
* **Status Dashboard:** Example frontend metrics page

---

## 🔧 API Endpoints

| Method | Path                 | Description            |
| ------ | -------------------- | ---------------------- |
| GET    | /                    | Homepage               |
| GET    | /empregado           | Employee interface     |
| GET    | /cliente/{sessionId} | Client interface       |
| POST   | /api/create-session  | Create a new session   |
| GET    | /api/token           | Generate LiveKit token |
| GET    | /api/qr/{sessionId}  | Generate QR code       |
| GET    | /status              | Status dashboard       |
| GET    | /debug               | Debug tools            |

---

## 🧪 Testing

* **HTML test pages in `static/`** allow you to manually test camera, mic, WebRTC, and LiveKit integration.
* There is no automated test suite yet, but you can add backend `go test` or frontend end-to-end tests with Playwright in the future.

---

## 🐳 Docker & Coolify

**Included services:**

* `qrkit`: Main Go application (**port 8080**, *not exposed manually*)
* `livekit`: LiveKit server (**port 7880** TCP/UDP, media UDP ports)
* `db`: PostgreSQL (persistent)
* `redis`: Cache/queues (internal)

**Best practice for production:**
Only expose LiveKit ports (WebRTC), never the `qrkit` web port. Let Coolify handle the reverse proxy!

---

## 🔐 Configuration

### **Main environment variables:**

* `LIVEKIT_API_KEY`: LiveKit API Key (default: mykey)
* `LIVEKIT_API_SECRET`: LiveKit API Secret (default: mysecret)
* `LIVEKIT_HOST`: LiveKit server host (default: livekit:7880)
* `DATABASE_URL`: postgres://qrkit:qrkitpass@db:5432/qrkit?sslmode=disable
* `REDIS_URL`: redis://redis:6379

QRKit reads these variables automatically on startup.

---

## 🚨 Troubleshooting

* **Server won't start?**
  * Check logs (`docker-compose logs qrkit`)
  * Check if port is free (`lsof -i :8080`, for local/dev)
* **Camera issues?**
  * Browser permissions
  * Use HTTPS for production (Coolify takes care of this)
* **LiveKit not connecting?**
  * Check logs and UDP port config
  * Test with `curl http://livekit:7880`
* **QR code not working?**
  * Check if the session was created, use debug tools

---

## 🔧 Development

* **Endpoints:** In `main.go`
* **Templates:** In `templates/`
* **Frontend/JS:** In `static/`
* **Containers:** `docker-compose.yml`

---

## 📊 Monitoring

Only the demo status dashboard is included.
For production metrics, consider adding Prometheus & Grafana.

---

## 🤝 Contributing

* Fork, create a branch, open PRs, manual testing encouraged!

---

## 📄 License

MIT

---

## 🆘 Support

* See this documentation
* Use the built-in debug tools
* Check logs and status dashboard

---

Built with ❤️ using Go, LiveKit, WebRTC, and Coolify.

---

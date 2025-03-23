package server

import (
	"go-websocket-chats/database"
	"go-websocket-chats/websocket"
	"log"
	"net/http"
)

type ChatServer struct {
	ConnectionManager *websocket.ConnectionManager
	serverAddress     string
}

func NewChatServer(address string) *ChatServer {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		// Continue without database
		return &ChatServer{
			ConnectionManager: websocket.NewConnectionManager(nil),
			serverAddress:     address,
		}
	}

	// Create message repository
	messageRepo := websocket.NewMessageRepository(db)

	return &ChatServer{
		ConnectionManager: websocket.NewConnectionManager(messageRepo),
		serverAddress:     address,
	}
}

func (s *ChatServer) serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

// handles WebSocket connection requests
func (s *ChatServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	websocket.HandleWebSocketConnection(s.ConnectionManager, w, r)
}

// Route setup
func (s *ChatServer) setupRoutes() {
	// Route for Homepage
	http.HandleFunc("/", s.serveHomePage)

	// Route for WebSocket connections
	http.HandleFunc("/ws", s.handleWebSocket)
	// Route for static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func (s *ChatServer) Start() error {
	go s.ConnectionManager.Run()
	s.setupRoutes()
	log.Printf("Starting chat server on %s", s.serverAddress)
	return http.ListenAndServe(s.serverAddress, nil)
}

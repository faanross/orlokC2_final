package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var WebSocketPort = 8080

var upgrader = websocket.Upgrader{
	// Allow connections from any origin for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketServer represents a simple WebSocket server
type WebSocketServer struct {
	port      int
	clients   []*websocket.Conn
	clientsMx sync.Mutex
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{
		port:    port,
		clients: make([]*websocket.Conn, 0),
	}
}

// Start begins the WebSocket server
func (s *WebSocketServer) Start() error {
	// Set up HTTP handler for the WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// Start the server
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("WebSocket server starting on %s\n", addr)

	// Start the HTTP server (this is a blocking call)
	return http.ListenAndServe(addr, nil)
}

// handleWebSocket handles WebSocket connections
func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Add client to our list
	s.clientsMx.Lock()
	s.clients = append(s.clients, conn)
	s.clientsMx.Unlock()

	fmt.Println("New WebSocket connection established")

	// Send a welcome message
	welcomeMsg := Message{
		Type:   ResponseMessage,
		Status: "connected",
	}
	conn.WriteJSON(welcomeMsg)

	// Clean up when the connection closes
	defer func() {
		conn.Close()

		// Remove client from slice
		s.clientsMx.Lock()
		for i, c := range s.clients {
			if c == conn {
				// Remove this client (order doesn't matter)
				s.clients[i] = s.clients[len(s.clients)-1]
				s.clients = s.clients[:len(s.clients)-1]
				break
			}
		}
		s.clientsMx.Unlock()
	}()

	// Message reading loop
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Connection closed or error
			break
		}

		log.Printf("Received message: %+v", msg)

		if msg.Type == CommandMessage {
			// Queue the command directly instead of calling ProcessWSCommand
			if msg.Command != "" {
				AgentCommands.QueueCommand(msg.Command)

				// Acknowledge receipt
				response := Message{
					Type:    ResponseMessage,
					Command: msg.Command,
					Output:  "Command queued for execution",
					Status:  "queued",
				}
				conn.WriteJSON(response)
			}
		}
	}
}

// Global WebSocket server instance
var GlobalWSServer *WebSocketServer

func StartWebSocketServer() {
	// Start WebSocket server in a separate goroutine
	GlobalWSServer = NewWebSocketServer(WebSocketPort)
	fmt.Printf("Starting WebSocket server on port %d...\n", WebSocketPort)
	go func() {
		err := GlobalWSServer.Start()
		if err != nil {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()

	// Give the WebSocket server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Println("WebSocket server is running. You can now connect from the web UI.")

}

// Broadcast sends a message to all connected clients
func (s *WebSocketServer) Broadcast(msg Message) {
	s.clientsMx.Lock()
	defer s.clientsMx.Unlock()

	// Send to all clients
	for i := 0; i < len(s.clients); {
		conn := s.clients[i]
		err := conn.WriteJSON(msg)

		if err != nil {
			// Remove disconnected client
			s.clients[i] = s.clients[len(s.clients)-1]
			s.clients = s.clients[:len(s.clients)-1]
		} else {
			i++
		}
	}
}

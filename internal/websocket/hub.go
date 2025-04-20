package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"fmt"
)

// Hub maintains all active client connections
type Hub struct {
	// Registered clients by user ID
	Clients map[int]*Client
	
	// Client registration channel
	Register chan *Client
	
	// Client unregistration channel
	Unregister chan *Client
	
	// Mutex for thread-safety
	mutex sync.Mutex
}

// Client represents a connected user
type Client struct {
	UserID   int
	Username string
	Send     chan []byte
	Hub      *Hub
	Conn     *Connection
}

// NewHub creates a new hub for managing clients
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client events
func (h *Hub) Run() {
	log.Println("Starting WebSocket hub")
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client.UserID] = client
			h.mutex.Unlock()
			h.broadcastUserList()
			
		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mutex.Unlock()
			h.broadcastUserList()
		}
	}
}

// broadcastUserList sends the updated user list to all clients
func (h *Hub) broadcastUserList() {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    
    // Prepare user list for sending
    var users []map[string]interface{}
    for userID, client := range h.Clients {
        users = append(users, map[string]interface{}{
            "id":       userID,
            "username": client.Username,
        })
    }

    // Create message object
    message := map[string]interface{}{
        "type":  "user_list",
        "users": users,
    }

    // Convert to JSON
    data, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error creating user list: %v", err)
        return
    }

    // Send to all clients
    h.broadcastMessage(data)
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(message []byte) {
    for _, client := range h.Clients {
		client.Send <- message
		fmt.Println(string(message))
    }
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int, message []byte) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	client, exists := h.Clients[userID]
	if !exists {
		return false
	}
	
	select {
	case client.Send <- message:
		return true
	default:
		return false
	}
}

// IsUserOnline checks if a user is currently connected
func (h *Hub) IsUserOnline(userID int) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	_, exists := h.Clients[userID]
	return exists
}
package websocket

import (
	"encoding/json"
	"forum/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// Connection wraps a websocket connection
type Connection struct {
	ws *websocket.Conn
}

// ServeWs handles websocket connections
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID int, username string) {
	// Upgrade HTTP to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	
	// Create client and connection
	conn := &Connection{ws: ws}
	client := &Client{
		UserID:   userID,
		Username: username,
		Send:     make(chan []byte, 256),
		Hub:      hub,
		Conn:     conn,
	}

	// Register with hub
	client.Hub.Register <- client

	// Start read/write processes
	go client.readPump()
	go client.writePump()
}

// readPump handles incoming messages
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.ws.Close()
	}()

	c.Conn.ws.SetReadLimit(4096)
	c.Conn.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.ws.SetPongHandler(func(string) error {
		c.Conn.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.Conn.ws.ReadMessage()
		if err != nil {
			break
		}

		// Parse the message
		var message Message
		if err := json.Unmarshal(data, &message); err != nil {
			continue
		}

		// Handle message based on type
		switch message.Type {
		case "message":
			handleChatMessage(c, message)
		case "typing":
			// Simply forward typing notification with username
			respMsg := Message{
				Type:       "typing",
				SenderID:   c.UserID,
				ReceiverID: message.ReceiverID,
				Username:   c.Username,
			}
			respData, _ := json.Marshal(respMsg)
			c.Hub.SendToUser(message.ReceiverID, respData)
		case "typing_stopped":
			// Forward typing stopped notification
			respMsg := Message{
				Type:       "typing_stopped",
				SenderID:   c.UserID,
				ReceiverID: message.ReceiverID,
			}
			respData, _ := json.Marshal(respMsg)
			c.Hub.SendToUser(message.ReceiverID, respData)
		case "get_history":
			handleHistoryRequest(c, message)
		case "get_more_history":
			handleMoreHistoryRequest(c, message)
		}
	}
}

// handleChatMessage processes chat messages
func handleChatMessage(c *Client, message Message) {
	// Set sender and store in database
	senderID := c.UserID
	receiverID := message.ReceiverID
	content := message.Content
	
	StoreMessage(senderID, receiverID, content)
	
	// Create response with timestamp
	timestamp := time.Now().Format(time.RFC3339)
	responseMsg := Message{
		Type:       "message",
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Timestamp:  timestamp,
		Username:   c.Username,
	}
	
	// Send to both sender and receiver
	respData, _ := json.Marshal(responseMsg)
	c.Hub.SendToUser(receiverID, respData)
	c.Send <- respData
}

// handleHistoryRequest gets chat history
func handleHistoryRequest(c *Client, message Message) {
	otherUserID := message.ReceiverID
	
	// Get message history between users
	messages, err := GetMessageHistory(c.UserID, otherUserID, 10)
	if err != nil {
		return
	}
	
	// Add usernames to messages
	for i := range messages {
		if messages[i].SenderID == c.UserID {
			messages[i].Username = c.Username
		} else {
			var username string
			database.Db.QueryRow("SELECT username FROM users WHERE id = ?", messages[i].SenderID).Scan(&username)
			messages[i].Username = username
		}
	}
	
	// Send history to client
	response, _ := json.Marshal(map[string]interface{}{
		"type":     "history",
		"messages": messages,
	})
	
	c.Send <- response
}

// handleMoreHistoryRequest gets older messages
func handleMoreHistoryRequest(c *Client, message Message) {
	otherUserID := message.ReceiverID
	before := message.Timestamp
	
	// Get older messages
	messages, err := GetMoreMessageHistory(c.UserID, otherUserID, before, 10)
	if err != nil {
		return
	}
	
	// Add usernames 
	for i := range messages {
		if messages[i].SenderID == c.UserID {
			messages[i].Username = c.Username
		} else {
			var username string
			database.Db.QueryRow("SELECT username FROM users WHERE id = ?", messages[i].SenderID).Scan(&username)
			messages[i].Username = username
		}
	}
	
	// Send response
	response, _ := json.Marshal(map[string]interface{}{
		"type":     "more_history",
		"messages": messages,
	})
	
	c.Send <- response
}

// writePump sends messages to the client
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
			
		case <-ticker.C:
			// Send ping to keep connection alive
			c.Conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
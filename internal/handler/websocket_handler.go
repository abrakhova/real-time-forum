package handler

import (
	"forum/internal/database"
	"forum/internal/session"
	"forum/internal/websocket"
	"net/http"
)

// WebSocket hub for managing connections
var WebSocketHub *websocket.Hub

// InitWebSocketHub creates and starts the WebSocket hub
func InitWebSocketHub() {
	WebSocketHub = websocket.NewHub()
	go WebSocketHub.Run()
}

// WebSocketHandler manages WebSocket connection requests
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user session
	userID, err := session.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch username for the authenticated user
	var username string
	err = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}

	// Serve the WebSocket connection
	websocket.ServeWs(WebSocketHub, w, r, userID, username)
}
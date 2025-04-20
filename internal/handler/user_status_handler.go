package handler

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/session"
	"forum/internal/util"
	"net/http"
)

// UserStatusHandler returns the current user's session information
func UserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	userID, err := session.GetUserIDFromSession(r)
	
	// Initialize username
	var username string
	if err == nil && userID > 0 {
		// Attempt to fetch username
		_ = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	}
	
	// Prepare response data
	data := struct {
		SessionID int    `json:"sessionID"`
		Username  string `json:"username"`
		LoggedIn  bool   `json:"loggedIn"`
	}{
		SessionID: userID,
		Username:  username,
		LoggedIn:  userID > 0,
	}
	
	// Return user status
	util.ExecuteJSON(w, data, http.StatusOK)
}
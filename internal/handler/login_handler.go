package handler

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/session"
	"forum/internal/user"
	"forum/internal/util"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ExecuteJSON(w, model.MsgData{"Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	
	// Get login credentials
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Validate inputs
	if identifier == "" || password == "" {
		util.ExecuteJSON(w, model.MsgData{"Identifier and password are required"}, http.StatusBadRequest)
		return
	}

	// Authenticate user
	userID, err := user.AuthenticateUser(identifier, password)
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Invalid identifier or password"}, http.StatusUnauthorized)
		return
	}

	// Create session
	if err := session.CreateSession(w, userID); err != nil {
		util.ExecuteJSON(w, model.MsgData{"Session creation failed"}, http.StatusInternalServerError)
		return
	}

	// Get username for response
	var username string
	_ = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)

	// Send successful login response
	util.ExecuteJSON(w, struct {
		Message   string `json:"message"`
		SessionID int    `json:"sessionID"`
		Username  string `json:"username"`
	}{
		Message:   "Login successful",
		SessionID: userID,
		Username:  username,
	}, http.StatusOK)
}
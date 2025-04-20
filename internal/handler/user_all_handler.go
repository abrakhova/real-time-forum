package handler

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/util"
	"net/http"
)

// GetAllUsersHandler returns a list of all registered users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Query all users
	rows, err := database.Db.Query("SELECT id, username FROM users")
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Failed to load users"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Struct to hold user information
	type UserInfo struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}

	var users []UserInfo

	// Scan and collect user data
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			// Skip individual errors to return as many users as possible
			continue
		}
		users = append(users, user)
	}

	// Return the user list
	util.ExecuteJSON(w, struct {
		Users []UserInfo `json:"users"`
	}{
		Users: users,
	}, http.StatusOK)
}
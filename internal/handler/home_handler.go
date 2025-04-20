package handler

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/post"
	"forum/internal/session"
	"forum/internal/util"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Get session and validate
	sessionID, err := session.GetUserIDFromSession(r)
	if err != nil || sessionID == 0 {
		util.ExecuteJSON(w, model.MsgData{"Unauthorized: Please log in to view posts"}, http.StatusUnauthorized)
		return
	}
	
	// Get username for the logged-in user
	var username string
	if sessionID > 0 {
		err = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", sessionID).Scan(&username)
		if err != nil {
			username = "" // Continue even if username fetch fails
		}
	}

	// Fetch posts
	allPosts, err := post.FetchPosts()
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Failed to load posts"}, http.StatusInternalServerError)
		return
	}
	
	// Limit to 5 most recent/popular posts
	if len(allPosts) > 5 {
		allPosts = allPosts[:5]
	}

	// Prepare response data
	data := model.Data{
		Posts:     allPosts,
		SessionID: sessionID,
		Username:  username,
	}

	// Set proper content type and return JSON response
	w.Header().Set("Content-Type", "application/json")
	util.ExecuteJSON(w, data, http.StatusOK)
}
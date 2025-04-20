package handler

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/post"
	"forum/internal/session"
	"forum/internal/util"
	"log"
	"net/http"
	"strings"
)

// CreatePostHandler handles creating a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := session.GetUserIDFromSession(r)
		if err != nil {
			util.ExecuteJSON(w, model.MsgData{"Invalid session, please log in"}, http.StatusUnauthorized)
			return
		}

		title := strings.TrimSpace(r.FormValue("title")) 
        content := strings.TrimSpace(r.FormValue("content")) 
        categories := strings.Join(r.Form["categories"], ", ")

		if title == "" || content == "" {
			util.ExecuteJSON(w, model.MsgData{"Post cannot be empty"}, http.StatusBadRequest)
			return
		}

		categoriesList := []string{
			"General",
			"Local News & Events",
			"Viking line",
			"Travel",
			"Sailing",
			"Cuisine & food",
			"Politics"}

		for _, userCategory := range r.Form["categories"] {
			catValid := false
			for _, categoryItem := range categoriesList {
				if categoryItem == userCategory {
					catValid = true
					break
				}
			}
			if !catValid {
				util.ExecuteJSON(w, model.MsgData{"Invalid category"}, http.StatusBadRequest)
				return
			}
		}

		if len(categories) == 0 {
			categories = "General"
		}

		// Insert the post into the database
		if err := post.CreatePost(userID, title, content, categories); err != nil {
			log.Println("Post creation failed:", err)
			util.ExecuteJSON(w, model.MsgData{"Post creation failed"}, http.StatusInternalServerError)
			return
		}

		id, err := post.GetPostId()
		if err != nil {
			log.Println("Database issue:", err)
			util.ExecuteJSON(w, model.MsgData{"Database issue"}, http.StatusInternalServerError)
			return
		}

		// Return JSON response with the new post ID
		util.ExecuteJSON(w, struct {
			Message string `json:"message"`
			ID      int    `json:"id"`
		}{
			Message: "Post created successfully",
			ID:      id,
		}, http.StatusOK)
	} else if r.Method == "GET" {
		sessionID, err := session.GetUserIDFromSession(r)
		if err != nil {
			sessionID = 0 // If there's an error, set sessionID to 0
		}

		var username string
		if sessionID > 0 {
			_ = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", sessionID).Scan(&username)
		}

		data := struct {
			SessionID int    `json:"sessionID"`
			Username  string `json:"username"`
		}{
			SessionID: sessionID,
			Username:  username,
		}

		// Return JSON for GET request (form data)
		util.ExecuteJSON(w, data, http.StatusOK)
	} else {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
	}
}
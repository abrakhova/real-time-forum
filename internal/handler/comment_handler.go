package handler

import (
	"forum/internal/comment"
	"forum/internal/model"
	"forum/internal/session"
	"forum/internal/util"
	"log"
	"net/http"
	"strings"
)

// CommentHandler handles adding a comment to a post
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := session.GetUserIDFromSession(r)
    if err != nil || sessionID == 0 {
        log.Println("Unauthorized access attempt to view post")
        util.ExecuteJSON(w, model.MsgData{"Unauthorized: Please log in to view posts"}, http.StatusUnauthorized)
        return
    }

	postID := r.FormValue("post_id")
	if postID == "" {
		util.ExecuteJSON(w, model.MsgData{"Post ID is missing"}, http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		util.ExecuteJSON(w, model.MsgData{"Content is missing"}, http.StatusBadRequest)
		return
	}

	if err := comment.AddComment(sessionID, postID, content); err != nil {
		log.Println("Failed to add comment:", err)
		util.ExecuteJSON(w, model.MsgData{"Failed to add the comment"}, http.StatusInternalServerError)
		return
	}

	// Send success response
	util.ExecuteJSON(w, model.MsgData{"Comment added successfully"}, http.StatusOK)
}
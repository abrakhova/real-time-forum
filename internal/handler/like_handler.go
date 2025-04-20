package handler

import (
	"forum/internal/model"
	"forum/internal/reaction"
	"forum/internal/session"
	"forum/internal/util"
	"net/http"
)

// LikeHandler handles liking or disliking a post or comment
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Validate session
	sessionID, err := session.GetUserIDFromSession(r)
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Invalid session, please log in"}, http.StatusUnauthorized)
		return
	}

	// Get request parameters
	itemID := r.FormValue("item_id")
	isComment := r.FormValue("is_comment") == "true"
	reactionType := r.FormValue("type") // "like" or "dislike"

	// Validate inputs
	if itemID == "" {
		util.ExecuteJSON(w, model.MsgData{"Item ID is missing"}, http.StatusBadRequest)
		return
	}

	// Process reaction
	if err := reaction.LikeItem(sessionID, itemID, isComment, reactionType); err != nil {
		util.ExecuteJSON(w, model.MsgData{"Failed to process reaction"}, http.StatusInternalServerError)
		return
	}

	// Send success response
	util.ExecuteJSON(w, model.MsgData{"Reaction recorded successfully"}, http.StatusOK)
}
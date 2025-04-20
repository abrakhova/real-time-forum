package handler

import (
	"forum/internal/comment"
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/post"
	"forum/internal/reaction"
	"forum/internal/session"
	"forum/internal/util"
	"net/http"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Validate session
	sessionID, err := session.GetUserIDFromSession(r)
	if err != nil || sessionID == 0 {
		util.ExecuteJSON(w, model.MsgData{"Unauthorized: Please log in to view posts"}, http.StatusUnauthorized)
		return
	}

	// Get the post ID from URL query
	postID := r.URL.Query().Get("id")
	if postID == "" || postID == "undefined" || postID == "null" {
		util.ExecuteJSON(w, model.MsgData{"Missing or invalid PostID"}, http.StatusBadRequest)
		return
	}

	// Fetch post details
	post, err := post.FetchPost(postID)
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Failed to load the post"}, http.StatusInternalServerError)
		return
	}

	// Fetch post reactions
	post.Likes, post.Dislikes, err = reaction.FetchReactionsNumber(post.ID, false)
	if err != nil {
		// Continue with zero likes/dislikes if fetch fails
		post.Likes = 0
		post.Dislikes = 0
	}

	// Fetch comments for the post
	post.Comments, err = comment.FetchCommentsForPost(post.ID)
	if err != nil {
		// Continue with empty comments if fetch fails
		post.Comments = []model.Comment{}
	}

	// Get username for the logged-in user
	var username string
	if sessionID > 0 {
		_ = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", sessionID).Scan(&username)
	}

	// Prepare and send response
	postPageData := struct {
		Post      model.Post `json:"post"`
		SessionID int        `json:"sessionID"`
		Username  string     `json:"username"`
	}{
		Post:      post,
		SessionID: sessionID,
		Username:  username,
	}

	util.ExecuteJSON(w, postPageData, http.StatusOK)
}
package handler

import (
	"database/sql"
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/session"
	"forum/internal/util"
	"net/http"
)

// FilterHandler handles filtering posts by category
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := session.GetUserIDFromSession(r)
	if err != nil || sessionID == 0 {
		util.ExecuteJSON(w, model.MsgData{"Unauthorized: Please log in to view posts"}, http.StatusUnauthorized)
		return
	}

	// Get logged-in user details
	var username string
	if sessionID > 0 {
		_ = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", sessionID).Scan(&username)
	}

	// Get filter parameters
	category := r.URL.Query().Get("category")
	userCreated := r.URL.Query().Get("user_created") == "true"
	liked := r.URL.Query().Get("liked") == "true"

	var rows *sql.Rows
	var queryErr error

	// Determine which query to run based on filter
	switch {
	case userCreated && sessionID > 0:
		rows, queryErr = database.Db.Query("SELECT id, title, category FROM posts WHERE user_id = ?", sessionID)
	case liked && sessionID > 0:
		rows, queryErr = database.Db.Query(`
			SELECT p.id, p.title, p.category 
			FROM posts p
			JOIN reactions r ON p.id = r.post_id 
			WHERE r.user_id = ? AND r.type = 'like'
		`, sessionID)
	case category != "":
		rows, queryErr = database.Db.Query("SELECT id, title, category FROM posts WHERE category LIKE ?", "%"+category+"%")
	default:
		util.ExecuteJSON(w, model.MsgData{"Invalid filter request"}, http.StatusBadRequest)
		return
	}

	if queryErr != nil {
		util.ExecuteJSON(w, model.MsgData{"Failed to fetch posts"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect filtered posts
	var posts []model.PostData
	for rows.Next() {
		var post model.PostData
		if err := rows.Scan(&post.ID, &post.Title, &post.Category); err != nil {
			util.ExecuteJSON(w, model.MsgData{"Failed to scan post"}, http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Send JSON response
	util.ExecuteJSON(w, struct {
		Category  string          `json:"category"`
		Posts     []model.PostData `json:"posts"`
		SessionID int             `json:"sessionID"`
		Username  string          `json:"username"`
	}{
		Category:  category,
		Posts:     posts,
		SessionID: sessionID,
		Username:  username,
	}, http.StatusOK)
}
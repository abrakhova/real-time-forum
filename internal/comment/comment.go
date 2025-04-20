package comment

import (
	"forum/internal/database"
	"forum/internal/model"
	"forum/internal/reaction"
)

func FetchCommentsForPost(postID int) ([]model.Comment, error) {
	commentRows, err := database.Db.Query(`
    SELECT c.id, c.user_id, c.content, u.username
    FROM comments c
    JOIN users u ON u.id = c.user_id
    WHERE c.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	var comments []model.Comment

	for commentRows.Next() {
		var comment model.Comment
		err := commentRows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Username)
		if err != nil {
			return nil, err
		}

		comment.Likes, comment.Dislikes, err = reaction.FetchReactionsNumber(comment.ID, true)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := commentRows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func AddComment(userID int, postID string, content string) error {
	query := "INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)"
	_, err := database.Db.Exec(query, userID, postID, content)
	return err
}

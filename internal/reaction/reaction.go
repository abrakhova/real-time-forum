package reaction

import (
	"database/sql"
	"fmt"
	"forum/internal/database"
)

func FetchReactionsNumber(itemID int, isComment bool) (likes, dislikes int, err error) {
	var query string
	if isComment {
		query = `
			SELECT 
				COUNT(CASE WHEN type = 'like' THEN 1 END) AS likes,
				COUNT(CASE WHEN type = 'dislike' THEN 1 END) AS dislikes
			FROM reactions 
			WHERE comment_id = ?`
	} else {
		query = `
			SELECT 
				COUNT(CASE WHEN type = 'like' THEN 1 END) AS likes,
				COUNT(CASE WHEN type = 'dislike' THEN 1 END) AS dislikes
			FROM reactions 
			WHERE post_id = ?`
	}

	err = database.Db.QueryRow(query, itemID).Scan(&likes, &dislikes)
	return
}

func LikeItem(userID int, itemID string, isComment bool, reactionType string) error {
	var table, idColumn string
	if isComment {
		table, idColumn = "reactions", "comment_id"
	} else {
		table, idColumn = "reactions", "post_id"
	}

	var existingReaction string
	query := fmt.Sprintf("SELECT type FROM %s WHERE user_id = ? AND %s = ?", table, idColumn)
	err := database.Db.QueryRow(query, userID, itemID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		insertQuery := fmt.Sprintf("INSERT INTO %s (%s, user_id, type) VALUES (?, ?, ?)", table, idColumn)
		_, err := database.Db.Exec(insertQuery, itemID, userID, reactionType)
		return err
	} else if err != nil {
		return err
	}

	if existingReaction != reactionType {
		updateQuery := fmt.Sprintf("UPDATE %s SET type = ? WHERE user_id = ? AND %s = ?", table, idColumn)
		_, err := database.Db.Exec(updateQuery, reactionType, userID, itemID)
		return err
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND %s = ?", table, idColumn)
	_, err = database.Db.Exec(deleteQuery, userID, itemID)
	return err
}

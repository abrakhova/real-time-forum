package websocket

import (
	"forum/internal/database"
	"time"
)

// Message represents a chat message
type Message struct {
	Type       string `json:"type"`
	ID         int    `json:"id,omitempty"`
	SenderID   int    `json:"sender_id,omitempty"`
	ReceiverID int    `json:"receiverID"` 
	Content    string `json:"content,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Username   string `json:"username,omitempty"`
}

// StoreMessage saves a message to the database
func StoreMessage(senderID, receiverID int, content string) error {
	timestamp := time.Now().Format(time.RFC3339)
	
	_, err := database.Db.Exec(
		"INSERT INTO private_messages (sender_id, receiver_id, content, timestamp) VALUES (?, ?, ?, ?)",
		senderID, receiverID, content, timestamp,
	)
	
	return err
}

// GetMessageHistory retrieves message history between two users
func GetMessageHistory(userID1, userID2 int, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 10
	}
	
	// Get messages between the two users
	rows, err := database.Db.Query(`
		SELECT sender_id, receiver_id, content, timestamp
		FROM private_messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY timestamp DESC
		LIMIT ?
	`, userID1, userID2, userID2, userID1, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var timestamp string
		err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content, &timestamp)
		if err != nil {
			continue
		}
		msg.Type = "message"
		msg.Timestamp = timestamp
		messages = append(messages, msg)
	}

	// Reverse order to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetMoreMessageHistory loads older messages
func GetMoreMessageHistory(userID1, userID2 int, before string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 10
	}
	
	// Parse timestamp or use current time
	var beforeTime time.Time
	if before != "" {
		var err error
		beforeTime, err = time.Parse(time.RFC3339, before)
		if err != nil {
			beforeTime = time.Now()
		}
	} else {
		beforeTime = time.Now()
	}
	
	// Get messages before the timestamp
	rows, err := database.Db.Query(`
		SELECT sender_id, receiver_id, content, timestamp
		FROM private_messages
		WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
		AND timestamp < ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, userID1, userID2, userID2, userID1, beforeTime.Format(time.RFC3339), limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var timestamp string
		err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content, &timestamp)
		if err != nil {
			continue
		}
		msg.Type = "message"
		msg.Timestamp = timestamp
		messages = append(messages, msg)
	}

	// Reverse order to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
package session

import (
	"database/sql"
	"fmt"
	"forum/internal/database"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// CreateSession generates a new session for a user
func CreateSession(w http.ResponseWriter, userID int) error {
	// Delete any existing sessions for this user
	_, err := database.Db.Exec("DELETE FROM sessions WHERE id = ?", userID)
	if err != nil {
		return err
	}

	// Generate a new session ID
	sessionID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	// Set session expiration
	expires := time.Now().Add(24 * time.Hour)
	
	// Insert new session
	_, err = database.Db.Exec(
		"INSERT INTO sessions (session_id, id, expires_at) VALUES (?, ?, ?)",
		sessionID.String(), userID, expires,
	)
	if err != nil {
		return err
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		Expires:  expires,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	
	return nil
}

// DeleteSession removes a session from the database
func DeleteSession(sessionID string) error {
	_, err := database.Db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// GetUserIDFromSession retrieves the user ID for a given session
func GetUserIDFromSession(r *http.Request) (int, error) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, fmt.Errorf("no session cookie found")
	}

	var userID int
	var expiresAt time.Time
	
	// Retrieve session details
	err = database.Db.QueryRow(
		"SELECT id, expires_at FROM sessions WHERE session_id = ?", cookie.Value,
	).Scan(&userID, &expiresAt)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("session not found")
	} else if err != nil {
		return 0, err
	}

	// Check session expiration
	if time.Now().After(expiresAt) {
		// Delete expired session
		go DeleteSession(cookie.Value)
		return 0, fmt.Errorf("session expired")
	}

	// Extend session on activity
	go extendSession(cookie.Value)
	
	return userID, nil
}

// extendSession updates the session expiration time
func extendSession(sessionID string) {
	newExpiry := time.Now().Add(24 * time.Hour)
	
	_, _ = database.Db.Exec(
		"UPDATE sessions SET expires_at = ? WHERE session_id = ?",
		newExpiry, sessionID,
	)
}

// CleanupExpiredSessions removes all expired sessions from the database
func CleanupExpiredSessions() {
	result, _ := database.Db.Exec("DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP")
	
	count, _ := result.RowsAffected()
	if count > 0 {
		fmt.Printf("Cleaned up %d expired sessions\n", count)
	}
}

// IsAuthenticated checks if a request is from an authenticated user
func IsAuthenticated(r *http.Request) bool {
	userID, err := GetUserIDFromSession(r)
	return err == nil && userID > 0
}
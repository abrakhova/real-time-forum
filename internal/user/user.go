package user

import (
	"errors"
	"forum/internal/database"

	"golang.org/x/crypto/bcrypt"
)

// EmailExists checks if an email is already registered
func EmailExists(email string) (bool, error) {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count > 0, err
}

// UsernameExists checks if a username is already taken
func UsernameExists(username string) (bool, error) {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0, err
}

// HashPassword generates a bcrypt hash for the given password
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// SaveUser adds a new user to the database
func SaveUser(username, email, hashedPassword, firstName, lastName, gender string, age int) error {
	_, err := database.Db.Exec(
		"INSERT INTO users (username, email, password, first_name, last_name, age, gender) VALUES (?, ?, ?, ?, ?, ?, ?)",
		username, email, hashedPassword, firstName, lastName, age, gender,
	)
	return err
}

// AuthenticateUser verifies user credentials
func AuthenticateUser(identifier, password string) (int, error) {
	var userID int
	var storedHash string
	
	// Try to authenticate with email
	err := database.Db.QueryRow(
		"SELECT id, password FROM users WHERE email = ?",
		identifier,
	).Scan(&userID, &storedHash)
	
	// If email lookup fails, try username
	if err != nil {
		err = database.Db.QueryRow(
			"SELECT id, password FROM users WHERE username = ?",
			identifier,
		).Scan(&userID, &storedHash)
		
		if err != nil {
			return 0, err
		}
	}
	
	// Compare passwords
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return 0, errors.New("invalid credentials")
	}
	
	return userID, nil
}
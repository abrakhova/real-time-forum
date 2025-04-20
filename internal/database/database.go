package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// initDB initializes the SQLite database and creates the necessary tables
func InitDB() {
	connectDB()
	verifyConnection()

	// Create tables
	createTables()
}

func connectDB() {
	var err error
	// Add WAL journal mode and busy_timeout to prevent most locking issues
	Db, err = sql.Open("sqlite3", "data/forum.db?_journal=WAL&_busy_timeout=5000")
	ErrorCheck("Database connection failed: ", err)
}

func verifyConnection() {
	err := Db.Ping()
	ErrorCheck("Database ping failed: ", err)
	log.Println("Database connected successfully")
}

// createTables creates the necessary tables in the database
func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password TEXT,
			first_name TEXT,
			last_name TEXT,
			age INTEGER,
			gender TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			content TEXT,
			category TEXT,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			content TEXT,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
	        comment_id INTEGER,  -- New column to store comment reactions
			type TEXT CHECK(type IN ('like', 'dislike')),
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			UNIQUE (user_id, post_id, comment_id)
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME,
		FOREIGN KEY(id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS private_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender_id INTEGER,
			receiver_id INTEGER,
			content TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(sender_id) REFERENCES users(id),
			FOREIGN KEY(receiver_id) REFERENCES users(id)
		);`,
		// Add some simple indexes to improve performance
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_id ON sessions(id);`,
	}

	for _, table := range tables {
		_, err := Db.Exec(table)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}

	log.Println("Tables created successfully")
}

func ErrorCheck(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
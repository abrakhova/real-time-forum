package main

import (
	"forum/internal/database"
	"forum/internal/handler"
	"forum/internal/session"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Initialize the database
	initializeDatabase()
	defer database.Db.Close()
	
	// Initialize the WebSocket hub
	handler.InitWebSocketHub()
	log.Println("WebSocket hub initialized")
	
	// Set up cleanup routine for sessions
	go sessionCleanupRoutine()
	
	// Set up routes and start server
	startServer()
}

// Periodically clean up expired sessions
func sessionCleanupRoutine() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			session.CleanupExpiredSessions()
		}
	}
}

func startServer() {
	// Serve static files
	serveStaticFiles()
	
	// Register auth handlers
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	
	// Register content handlers
	http.HandleFunc("/createPost", handler.CreatePostHandler)
	http.HandleFunc("/comment", handler.CommentHandler)
	http.HandleFunc("/like", handler.LikeHandler)
	http.HandleFunc("/filter", handler.FilterHandler)
	http.HandleFunc("/post", handler.ViewPostHandler)
	
	// Register user handlers
	http.HandleFunc("/user/status", handler.UserStatusHandler)
	http.HandleFunc("/user/all", handler.GetAllUsersHandler)
	
	// WebSocket endpoint
	http.HandleFunc("/ws", logRequest(handler.WebSocketHandler))
	log.Println("WebSocket endpoint registered at /ws")
	
	// Add a handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an API request
		isAPIRequest := strings.Contains(r.Header.Get("Accept"), "application/json") || 
						r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
						r.URL.Query().Get("api") == "true"
		
		log.Printf("Root path request: %s, isAPIRequest: %v", r.URL.Path, isAPIRequest)
		
		if r.URL.Path == "/" && isAPIRequest {
			// If JSON is requested, use the HomeHandler
			handler.HomeHandler(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/static/") {
			// Let the static file server handle this
			return
		} else {
			// For all other routes, serve the SPA index
			http.ServeFile(w, r, "./web/templates/index.html")
		}
	})
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeDatabase() {
	log.Println("Initializing database...")
	database.InitDB()
	session.CleanupExpiredSessions()
	log.Println("Database initialization complete")
}

func serveStaticFiles() {
	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("Static file server initialized")
}

// Middleware for logging requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	}
}
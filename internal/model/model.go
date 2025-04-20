package model

// Post represents a forum post
type Post struct {
	ID       int
	Username string
	UserID   int
	Title    string
	Content  string
	Category string
	Likes    int
	Dislikes int
	Comments []Comment
	Date     string
}

// Comment represents a comment on a post
type Comment struct {
	ID       int
	Username string
	UserID   int
	Content  string
	Likes    int
	Dislikes int
}

// PostPageData represents data for a single post page
type PostPageData struct {
	Post      Post
	SessionID int
}

// HomePageData represents summary data for posts on the home page
type HomePageData struct {
	ID       int
	Title    string
	Content  string
	Username string
	Likes    int
	Dislikes int
	Date     string
}

// Data represents the main data structure for the home page
type Data struct {
	Posts     []HomePageData
	SessionID int
	Username  string
}

// MsgData is a generic message response
type MsgData struct {
	Message string `json:"message"`
}

// PostData is a lightweight post representation
type PostData struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

// User represents a registered user
type User struct {
	ID        int
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Gender    string
}
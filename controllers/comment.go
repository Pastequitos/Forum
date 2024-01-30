package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// Comment represents a comment on a post.
type Comment struct {
	PostID   string
	UserID   string
	Content  string
	Date     string
	Username string // Assuming you have a field for the username associated with the userID
}

// AddComment handles adding comments to posts.
func AddComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("commented")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get form values
	postID := r.FormValue("post_id")
	content := r.FormValue("comment")
	cookie, err := r.Cookie("user_id")
	if err != nil {
		log.Println("Error getting user ID from cookie:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userID := cookie.Value
	date := time.Now()

	fmt.Println(postID)
	fmt.Println(content)
	fmt.Println(userID)
	fmt.Println(date)

	// Open the database connection
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert comment into the database
	_, err = db.Exec("INSERT INTO data_comments (post_id, user_id, comment, date) VALUES (?, ?, ?, ?)", postID, userID, content, date)
	if err != nil {
		log.Println("Error inserting comment into database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the post page after adding comment
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

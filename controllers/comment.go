package controllers

import (
	"database/sql"
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID := cookie.Value
	CheckSession(w, r)
	date := time.Now().Format("15h04 2 Jan 2006")

	// Open the database connection
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO data_comments (post_id, user_id, comment, date) VALUES (?, ?, ?, ?)", postID, userID, content, date)
	if err != nil {
		log.Println("Error inserting comment into database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update num_com for the corresponding post
	_, err = db.Exec("UPDATE data_post SET num_com = num_com + 1 WHERE id = ?", postID)
	if err != nil {
		log.Println("Error updating num_com:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the post page after adding comment
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

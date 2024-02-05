package controllers

import (
	"database/sql"
	"log"
	"net/http"
)

func LikeDislike(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the values from the form
	likeOrDislike := r.FormValue("likeordislike")
	postID := r.FormValue("post_id")

	// Retrieve the user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := cookie.Value
	CheckSession(w, r)
	// Open a database connection
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Check if the user has already performed an action on the post
	var existingAction string
	err = db.QueryRow("SELECT likeordislike FROM user_likeordislike WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingAction)
	switch {
	case err == sql.ErrNoRows:
		// User has not performed any action on the post
	case err != nil:
		// An unexpected error occurred
		log.Println("Error checking user action:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	default:
		// User has already performed an action on the post
		if existingAction == likeOrDislike {
			// If the user is trying to perform the same action again
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}

	// Fetch current like and dislike counts for the post from the database
	var currentLikes, currentDislikes int
	err = db.QueryRow("SELECT like, dislike FROM data_post WHERE id = ?", postID).Scan(&currentLikes, &currentDislikes)
	if err != nil {
		http.Error(w, "Failed to fetch post details from database", http.StatusInternalServerError)
		return
	}
	// Update the like and dislike counts based on the action
	switch likeOrDislike {
	case "like":
		if existingAction == "dislike" {
			currentDislikes--
		}
		currentLikes++
	case "dislike":
		if existingAction == "like" {
			currentLikes--
		}
		currentDislikes++
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// Update the database with the new like and dislike counts for the post
	_, err = db.Exec("UPDATE data_post SET like = ?, dislike = ? WHERE id = ?", currentLikes, currentDislikes, postID)
	if err != nil {
		log.Println("Error updating post details:", err)
		http.Error(w, "Failed to update post details in database", http.StatusInternalServerError)
		return
	}

	// Record the user action in the user_likeordislike table
	err = db.QueryRow("SELECT likeordislike FROM user_likeordislike WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingAction)
	switch {
	case err == sql.ErrNoRows:
		// No data exists, insert a new record
		_, err = db.Exec("INSERT INTO user_likeordislike (user_id, post_id, likeordislike) VALUES (?, ?, ?)", userID, postID, likeOrDislike)
		if err != nil {
			http.Error(w, "Failed to record user action in database", http.StatusInternalServerError)
			return
		}
	case err != nil:
		// An unexpected error occurred
		log.Println("Error checking user action:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	default:
		// Data already exists, update the existing record
		_, err = db.Exec("UPDATE user_likeordislike SET likeordislike = ? WHERE user_id = ? AND post_id = ?", likeOrDislike, userID, postID)
		if err != nil {
			http.Error(w, "Failed to update user action in database", http.StatusInternalServerError)
			return
		}
	}

	// Redirect to "/home"
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

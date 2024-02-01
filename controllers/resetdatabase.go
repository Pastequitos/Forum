package controllers

import (
	"database/sql"
	"net/http"
)

func ResetDatabase(w http.ResponseWriter, r *http.Request) {
	// Open the database connection
	db, _ := sql.Open("sqlite", "database.db")

	defer db.Close()

	// Delete all rows from data_post table
	if _, err := db.Exec("DELETE FROM data_post"); err != nil {
		return
	}

	// Delete all rows from data_comments table
	if _, err := db.Exec("DELETE FROM data_comments"); err != nil {
		return
	}

	// Delete all rows from data_comments table
	if _, err := db.Exec("DELETE FROM user_likeordislike"); err != nil {
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

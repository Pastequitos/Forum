package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "modernc.org/sqlite"
)

func Home(w http.ResponseWriter, r *http.Request) {
/* 	var errmsg []string */

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := time.Now()
		cookie, err := r.Cookie("user_id")
		if err != nil {
			log.Println("Error getting user ID from cookie:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		userID := cookie.Value

		db, err := sql.Open("sqlite", "database.db")
		if err != nil {
			log.Println("Error opening database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Query the database to get the maximum ID currently in use
		var maxID int
		row := db.QueryRow("SELECT MAX(id) FROM data_post")
		err = row.Scan(&maxID)
		if err != nil {
			log.Println("Error getting maximum post ID:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Increment the maxID to get the new post ID
		newID := maxID + 1

		insertPost := "INSERT INTO data_post (id, user_id, title, content, date, num_com) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = db.Exec(insertPost, newID, userID, title, content, date, 0)
		if err != nil {
			log.Println("Error inserting post into database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	ts := template.Must(template.ParseFiles("./ui/html/home.html"))
	err := ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

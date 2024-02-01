package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "modernc.org/sqlite"
)

func Filter(w http.ResponseWriter, r *http.Request) {
	var errmsg []string
	var posts []Post

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	category := r.FormValue("category")

	rows, err := db.Query("SELECT id, title, content, category FROM data_post WHERE category = ?", category)
	if err != nil {
		http.Error(w, "Failed to fetch posts from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category)
		if err != nil {
			// Handle error
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		// Append post to the slice
		posts = append(posts, post)
	}

	Data := struct {
		Error []string
		Posts []Post
	}{
		Error: errmsg,
		Posts: posts,
	}

	ts := template.Must(template.ParseFiles("./ui/html/home.html"))
	err = ts.Execute(w, Data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

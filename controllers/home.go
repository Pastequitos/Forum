package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "modernc.org/sqlite"
)

type Post struct {
	ID          int
	UserID      string
	Title       string
	Content     string
	Date        string
	Category    string
	NumComments int
	Like        int
	Dislike     int
	Comments    []Comment
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" && r.URL.Path != "/" {
		ErrorCode(w, r, 404, "Page not found")
		return
	}

	var errmsg []string
	var posts []Post

	CheckSession(w, r)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := time.Now().Format("15h04 2 Jan 2006")
		category := r.FormValue("category")
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}
		userID := cookie.Value
		like := 0
		dislike := 0

		db, err := sql.Open("sqlite", "database.db")
		if err != nil {
			log.Println("Error opening database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Query the database to get the maximum ID currently in use
		var maxID sql.NullInt64
		row := db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM data_post")
		err = row.Scan(&maxID)
		if err != nil {
			log.Println("Error getting maximum post ID:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if maxID is NULL
		if !maxID.Valid {
			// No posts in the database, set maxID to 0 or any other default value
			maxID.Int64 = 0
		}

		// Increment the maxID to get the new post ID
		newID := int(maxID.Int64) + 1
		insertPost := "INSERT INTO data_post (id, user_id, title, content, date, num_com, category, like, dislike) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = db.Exec(insertPost, newID, userID, title, content, date, 0, category, like, dislike)
		if err != nil {
			log.Println("Error inserting post into database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	// Fetching posts
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, user_id, title, content, date, num_com, category, like, dislike FROM data_post ORDER BY id DESC")
	if err != nil {
		log.Println("Error fetching posts from database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Date, &post.NumComments, &post.Category, &post.Like, &post.Dislike); err != nil {
			log.Println("Error scanning post row:", err)
			continue
		}

		// Fetching comments for each post
		commentRows, err := db.Query("SELECT post_id, user_id, comment, date FROM data_comments WHERE post_id = ?", post.ID)
		if err != nil {
			log.Println("Error fetching comments for post ID", post.ID, ":", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer commentRows.Close()

		for commentRows.Next() {
			var comment Comment
			if err := commentRows.Scan(&comment.PostID, &comment.UserID, &comment.Content, &comment.Date); err != nil {
				log.Println("Error scanning comment row:", err)
				continue
			}
			post.Comments = append(post.Comments, comment)
		}
		if err := commentRows.Err(); err != nil {
			log.Println("Error iterating over comment rows:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Append the post to the posts slsudo apt-get cleanice
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over post rows:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pass posts to the template
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

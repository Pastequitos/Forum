package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"main/controllers"
)

const (
	path = "database.db"
)

func main() {
	CreateDatabase(path)

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/home", controllers.Home)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/addcomment", controllers.AddComment)
	http.HandleFunc("/resetdatabase", controllers.ResetDatabase)
	http.HandleFunc("/filter", controllers.Filter)
	http.HandleFunc("/likedislike", controllers.LikeDislike)

	// Set Static file
	static := http.FileServer(http.Dir("ui"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	// Set Content-Type header for CSS files
	http.HandleFunc("/static/css/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		static.ServeHTTP(w, r)
	})

	// Set Content-Type header for image files
	http.HandleFunc("/static/media/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[1:] // Remove the leading slash
		contentType := getContentType(filePath)
		w.Header().Set("Content-Type", contentType)
		static.ServeHTTP(w, r)
	})

	log.Print("Starting server on http://localhost:3003")
	err := http.ListenAndServe(":3003", nil)
	log.Fatal(err)
}

func getContentType(filePath string) string {
	extension := filepath.Ext(filePath)
	switch extension {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream" // Default content type
	}
}

func CreateDatabase(path string) error {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS account_user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		mail TEXT,
		password TEXT
	);`)
	fmt.Println("database created")
	if err != nil {
		log.Fatal(err)
		return err

	}
	return nil
}

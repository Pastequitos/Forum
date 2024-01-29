package main

import (
	"log"
	"net/http"
	"path/filepath"

	"main/controllers"
)

func main() {
	// Display les pages
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/home", controllers.Home)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)



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

	log.Print("Starting server on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
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

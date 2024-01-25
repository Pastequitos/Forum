package main

import (
	"log"
	"net/http"

	"main/controllers"
)


func main() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/home", controllers.Home)

	// Serve static files
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Starting server on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	log.Fatal(err)
}

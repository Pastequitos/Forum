package controllers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func ErrorCode(w http.ResponseWriter, r *http.Request, status int, message string) {
	colorRed := "\033[31m" // Mise en place d'une couleur pour les erreurs
	// Mise en place de condition pour les différents types d'erreurs
	if status == 404 {
		w.WriteHeader(http.StatusNotFound)
		custTemplate, err := template.ParseFiles("./ui/html/404.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = custTemplate.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 404 : Error URL...", message) // Message qui sera afficher sur le terminal avec une précision de l'erreur
	}
	if status == 400 {
		w.WriteHeader(http.StatusNotFound)
		custTemplate, err := template.ParseFiles("./ui/html/400.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = custTemplate.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 400 : Invalid endpoint", message) // Message qui sera afficher sur le terminal avec une précision de l'erreur
	}
	if status == 500 {
		w.WriteHeader(http.StatusNotFound)
		custTemplate, err := template.ParseFiles("./ui/html/500.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = custTemplate.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 500 : Internal server error", message) // Message qui sera afficher sur le terminal avec une précision de l'erreur
	}
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./ui/html/404.html")
	if err != nil {
		ErrorCode(w, r, 404, "Template not found : 404.html")
		return
	}

	err = tmpl.Execute(w, nil)
}

func HandleServerError(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./ui/html/500.html")
	if err != nil {
		ErrorCode(w, r, 500, "Template not found : 500.html")
		return
	}

	err = tmpl.Execute(w, nil)
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./ui/html/400.html")
	if err != nil {
		ErrorCode(w, r, 400, "Template not found : 400.html")
		return
	}

	err = tmpl.Execute(w, nil)
}

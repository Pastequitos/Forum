package controllers

import (
	"log"
	"net/http"
	"text/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

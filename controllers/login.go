package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

func Login(w http.ResponseWriter, r *http.Request) {
	errmsg := "Wrong username or password"
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Open the database connection
		db, err := sql.Open("sqlite", "database.db")
		if err != nil {
			log.Println("Error opening database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Query the database for the user's details based on the username
		var dbUsername, dbPassword string
		row := db.QueryRow("SELECT username, password FROM account_user WHERE username = ?", username)
		err = row.Scan(&dbUsername, &dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				errmsg = append(errmsg, "Wrong Username")
			} else {
				log.Println("Error fetching user details:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			// Compare passwords
			if password != dbPassword {
				errmsg = append(errmsg, "Wrong password")
			} else {
				// Password is correct, redirect to home page
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
		}
	}

	// Pass error messages to the template
	data := struct {
		Error []string
	}{
		Error: errmsg,
	}

	// Render the login page with error messages
	ts := template.Must(template.ParseFiles("./ui/html/login.html"))
	err := ts.Execute(w, data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

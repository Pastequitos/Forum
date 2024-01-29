package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var errmsg []string
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get form values
		username := r.Form.Get("username")
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		// Open the database
		db, err := sql.Open("sqlite", "database.db")
		if err != nil {
			log.Println("Error opening database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		/// check username///

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM account_user WHERE username = ?", username)
		err = row.Scan(&count)
		if err != nil {
			log.Println("Error checking username:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			errmsg = append(errmsg, "Username already taken")
		}

		///// check mail //////

		var countmail int
		rowmail := db.QueryRow("SELECT COUNT(*) FROM account_user WHERE mail = ?", email)
		err = rowmail.Scan(&countmail)
		if err != nil {
			log.Println("Error checking email:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if countmail > 0 {
			errmsg = append(errmsg, "Email already taken")
		}

		if errmsg == nil {

			insertUser := "INSERT INTO account_user (username, mail, password) VALUES (?, ?, ?)"
			_, err = db.Exec(insertUser, username, email, password)
			if err != nil {
				fmt.Println(err)
				return
			}

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}

	}

	data := struct {
		Error []string
	}{
		Error: errmsg,
	}

	ts := template.Must(template.ParseFiles("./ui/html/signup.html"))
	err := ts.Execute(w, data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

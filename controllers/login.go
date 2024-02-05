package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gofrs/uuid"
	_ "github.com/gofrs/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var errmsg []string
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

				auth_token, err := uuid.NewV4()
				if err != nil {
					log.Println("Erreur lors de la génération de l'identifiant de session :", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				db.Exec("UPDATE account_user SET auth_token=? WHERE username=?", auth_token.String(), username)
				if err != nil {
					log.Println("Erreur lors de l'enregistrement de l'identifiant de session :", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Password is correct, set a cookie with the user ID
				cookie := http.Cookie{
					Name:  "user_id",
					Value: username, // Use the username as the user ID for demonstration purposes
					Path:  "/",
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    auth_token.String(), // Store some identifier for the user, like username
					HttpOnly: true,                // Cookie cannot be accessed through JavaScript
				})
				http.SetCookie(w, &cookie)

				// Redirect to the home page
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

	// Affichez la page pour écrire une discussion (write_discussion.html) avec le menu déroulant
	tmpl, err := template.ParseFiles("./ui/html/login.html")
	if err != nil {
		ErrorCode(w, r, 500, "Template not found : login.html")
		return
	}
	err = tmpl.Execute(w, data)
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	var storedSessionID string
	db, _ := sql.Open("sqlite", "database.db")
	sessionCookie, _ := r.Cookie("user_id")
	sessionidCookie, _ := r.Cookie("auth_token")
	if sessionCookie != nil {
		fmt.Println("ihavecookies")
		row := db.QueryRow("SELECT auth_token FROM account_user WHERE username=?", sessionCookie.Value)
		row.Scan(&storedSessionID)
		if storedSessionID != sessionidCookie.Value {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}
	}
}

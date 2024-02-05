package controllers

import (
	"fmt"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "user_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	cookie_token := &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	http.SetCookie(w, cookie_token)
	fmt.Println("Cookie reset")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

	http.SetCookie(w, cookie)
	fmt.Println("Cookie reset")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

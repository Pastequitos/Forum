package main

import (
	"net/http"


	controller "controller/controllers"
)

func main() {
	http.HandleFunc("/home", controller.Home())
}
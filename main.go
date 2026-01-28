package main

import (
	"net/http"
	"porfolio/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.Home)
	http.ListenAndServe(":8080", nil)
}

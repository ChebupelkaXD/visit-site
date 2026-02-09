package main

import (
	"log"
	"net/http"
	"porfolio/handlers"

	"github.com/gorilla/mux"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

func main() {
	r := mux.NewRouter()

	routes := []Route{
		{"/", handlers.Home, "GET"},
		{"/about", handlers.About, "GET"},
		{"/skills", handlers.Skills, "GET"},
		//{"/projects", handlers.Projects, "GET"},
		{"/experience", handlers.Experience, "GET"},
		{"/contacts", handlers.Contacts, "GET"},
	}

	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	for _, route := range routes {
		r.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	log.Fatal(http.ListenAndServe(":8080", r))
}

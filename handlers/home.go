package handlers

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/notfound.html",
))

func Home(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		NotFound(w, r)
		return
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := templates.ExecuteTemplate(w, "notfound.html", nil)

	if err != nil {
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

}

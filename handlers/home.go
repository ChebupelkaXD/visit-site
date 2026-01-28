package handlers

import (
	"html/template"
	"net/http"
	"time"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/home.html")

	if err != nil {
		panic(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title       string
		Greeting    string
		CurrentTime string
	}{
		Title:       "Home Page",
		Greeting:    "Hello, I'm QA!",
		CurrentTime: time.Now().Format("15:04"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

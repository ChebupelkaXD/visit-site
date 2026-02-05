package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"porfolio/data"
	"runtime"
)

var templates *template.Template

func init() {
	_, filename, _, _ := runtime.Caller(0)
	templatesDir := filepath.Join(filepath.Dir(filename), "..", "templates")

	templates = template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html")))
}

/*var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/notfound.html",
	"templates/about.html",
	"templates/skills.html",
	"templates/projects.html",
	"templates/experience.html",
	"templates/contacts.html",
))*/

func TemplateRender(w http.ResponseWriter, tmplName string) {
	err := templates.ExecuteTemplate(w, tmplName, data.Data)

	if err != nil {
		log.Printf("Template error (%s): %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "index.html")
}

func About(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "about.html")
}

func Skills(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "skills.html")
}

func Projects(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "projects.html")
}

func Experience(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "experience.html")
}

func Contacts(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "contacts.html")
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := templates.ExecuteTemplate(w, "notfound.html", nil)

	if err != nil {
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

}

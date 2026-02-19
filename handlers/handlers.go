package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"porfolio/config"
	"porfolio/data"
	"runtime"

	"gopkg.in/gomail.v2"
)

var (
	templates *template.Template
)

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

/*func Projects(w http.ResponseWriter, r *http.Request) {
	TemplateRender(w, "projects.html")
}*/

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

func ContactAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
	}

	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Email == "" || input.Message == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//log.Printf("Новый запрос на сотрудничество от %s, Email: %s, Сообщение: %s", input.Name, input.Email, input.Message)

	smtpHost := config.SMTP.Host
	smtpPort := config.SMTP.Port
	fromEmail := config.SMTP.From
	fromPassword := config.SMTP.Password
	toEmail := config.SMTP.To

	if fromEmail == "" || fromPassword == "" {
		log.Printf("Ошибка: отсутствуют SMTP credentials (FROM_EMAIL, FROM_PASSWORD)")
		http.Error(w, "Сервис временно недоступен", http.StatusInternalServerError)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", "Portfolio Site", fromEmail))
	m.SetHeader("To", toEmail)
	m.SetHeader("Reply-To", input.Email)
	m.SetHeader("Subject", fmt.Sprintf("Новое сообщение от %s", input.Name))

	body := fmt.Sprintf(`
        <h2>Новое сообщение с сайта</h2>
        <p><strong>Имя:</strong> %s</p>
        <p><strong>Email:</strong> %s</p>
        <p><strong>Сообщение:</strong></p>
        <blockquote>%s</blockquote>
        <hr>
        <small>Это автоматическое письмо с формы обратной связи</small>
    `, input.Name, input.Email, input.Message)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, fromEmail, fromPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Сообщение получено",
	})
}

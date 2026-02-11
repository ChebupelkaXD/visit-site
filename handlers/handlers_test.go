package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	t.Run("Тест рендера страниц /TemplateRender", func(t *testing.T) {
		w := httptest.NewRecorder()

		tests := []struct {
			tmplName   string
			wantStatus int
		}{
			{
				"index.html",
				200},
			{
				"notfound.html",
				200},
			{
				"about.html",
				200},
			{
				"skills.html",
				200},
			{
				"projects.html",
				200},
			{
				"experience.html",
				200},
			{
				"contacts.html",
				200},
		}

		for _, test := range tests {
			TemplateRender(w, test.tmplName)

			if w.Code != test.wantStatus {
				t.Errorf("Expected status %q, got %q", test.wantStatus, w.Code)
			}
		}
	})

	t.Run("Тест отображения страницы 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/not-found-page", nil)
		w := httptest.NewRecorder()

		NotFound(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expcted status %v, got %v", http.StatusNotFound, w.Code)
		}
	})

	t.Run("Тест контактной формы", func(t *testing.T) {
		data := `{"name": "Павел", "email": "example@mail.com", "message": "Hello World!"}`

		req, _ := http.NewRequest(http.MethodGet, "/api/contact", strings.NewReader(data))
		w := httptest.NewRecorder()

		ContactAPI(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %v", w.Code)
		}

		/*contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected JSON, got %s", contentType)
		}

		var output struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(w.Body).Decode(&output); err != nil {
			t.Errorf("Error parsing JSON")
		}*/

	})
}

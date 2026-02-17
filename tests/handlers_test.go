package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"porfolio/handlers"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitHandlersPage(t *testing.T) {
	t.Run("Тест рендера страниц TemplateRender", func(t *testing.T) {
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
			handlers.TemplateRender(w, test.tmplName)

			if w.Code != test.wantStatus {
				t.Errorf("Expected status %q, got %q", test.wantStatus, w.Code)
			}
		}
	})

	t.Run("Тест отображения страницы NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/not-found-page", nil)
		w := httptest.NewRecorder()

		handlers.NotFound(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expcted status %v, got %v", http.StatusNotFound, w.Code)
		}
	})
}

func TestUnitHandlersContactAPI(t *testing.T) {
	tests := []struct {
		desc    string
		methods []string
		target  string
		body    io.Reader
		expect  int
	}{
		{"Не разрешенный запрос", []string{http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodHead, http.MethodOptions}, "/api/contact", nil, http.StatusInternalServerError},
		{"Битый JSON", []string{http.MethodPost}, "/api/contact", strings.NewReader(`{"name": "Paul", "email": "example@mail.ru", "message": "Hello"`), http.StatusBadRequest},
		{"Не полный JSON", []string{http.MethodPost}, "/api/contact", strings.NewReader(`{"name": "", "email": "", "message": ""`), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			for _, method := range tt.methods {
				req := httptest.NewRequest(method, tt.target, tt.body)
				w := httptest.NewRecorder()

				handlers.ContactAPI(w, req)

				assert.Equal(t, w.Code, tt.expect)
			}

		})
	}
}

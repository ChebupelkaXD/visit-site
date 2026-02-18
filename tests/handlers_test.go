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

type Tests []struct {
	desc    string
	methods []string
	target  string
	body    io.Reader
	expect  int
}

func TestUnitHandlersPage(t *testing.T) {

	t.Run("Тест TemplateRender", func(t *testing.T) {
		w := httptest.NewRecorder()

		handlers.TemplateRender(w, "invisible.html")
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})

	t.Run("Тест NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/not-found-page", nil)
		w := httptest.NewRecorder()

		handlers.NotFound(w, req)

		assert.Equal(t, w.Code, http.StatusNotFound)
	})

	t.Run("Тест видимых страниц", func(t *testing.T) {
		tests := []struct {
			desc     string
			funcName func(w http.ResponseWriter, r *http.Request)
			target   string
			body     io.Reader
			expect   int
		}{
			{"Home", handlers.Home, "/", nil, http.StatusOK},
			{"About", handlers.About, "/about", nil, http.StatusOK},
			{"Skills", handlers.Skills, "/skills", nil, http.StatusOK},
			{"Experience", handlers.Experience, "/experience", nil, http.StatusOK},
			{"Contacts", handlers.Contacts, "/contacts", nil, http.StatusOK},
		}

		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tt.target, tt.body)
				w := httptest.NewRecorder()
				tt.funcName(w, req)

				assert.Equal(t, w.Code, tt.expect)
				assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
			})
		}

	})
}

func TestUnitHandlersContactAPI(t *testing.T) {
	tests := Tests{
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

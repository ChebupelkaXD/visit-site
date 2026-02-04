package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	t.Run("Главная страница", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		Home(w, req)

		want := "Im QA"
		got := w.Body.String()

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		if got != want {
			t.Errorf("Expected string: %s, got %s", want, got)
		}
	})
}

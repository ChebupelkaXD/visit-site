package tests

import (
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"porfolio/config"
	"porfolio/handlers"
	"strings"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
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

func TestIntegrationContactAPI_EmailSent(t *testing.T) {
	mockServer := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:  1025,
		HostAddress: "127.0.0.1",

		LogToStdout:       false,
		LogServerActivity: false})

	defer func() {
		if err := mockServer.Stop(); err != nil {
			t.Fatalf("Ошибка остановки SMTP сервера: %v", err)
		}
	}()

	go func() {
		if err := mockServer.Start(); err != nil {
			t.Errorf("Failed to start mock smtp: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	config.SMTP.Host = "127.0.0.1"
	config.SMTP.Port = 1025
	config.SMTP.Password = "12345"
	config.SMTP.From = "no-reply@test.local"
	config.SMTP.To = "recipient@test.local"

	body := `{"name":"Иван","email":"test@example.com","message":"Привет, хочу сотрудничать!"}`

	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handlers.ContactAPI(w, req)

	messages := mockServer.Messages()

	if len(messages) != 1 {
		t.Fatalf("ожидалось 1 письмо, получено %d", len(messages))
	}

	msg, err := mail.ReadMessage(strings.NewReader(messages[0].MsgRequest()))
	if err != nil {
		t.Fatalf("Не удалось распарсить письмо")
	}

	subject := msg.Header.Get("Subject")
	decoderSubject, err := (&mime.WordDecoder{}).DecodeHeader(subject)
	if err != nil {
		t.Fatalf("не удалось декодировать Subject: %v", err)
	}

	if !strings.Contains(decoderSubject, "Новое сообщение от") {
		t.Errorf("не найдена ожидаемая тема письма %s", messages[0].MsgRequest())
	}
}

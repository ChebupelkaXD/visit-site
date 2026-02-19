package tests

import (
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
)

func TestContactAPI_EmailSent(t *testing.T) {
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

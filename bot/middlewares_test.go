package bot

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func TestInvalidWebhookSecret(t *testing.T) {
	//: 5
	srv := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Secret: "extremely-secret", Endpoint: "/webhook-endpoint"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, srv.Config.Endpoint, nil)
	req.Header.Set("X-Gitlab-Token", "invalid-secret")

	srv.webhookSecret(testHandler).ServeHTTP(w, req)

	want := 401
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}
}

func TestWebhookSecret(t *testing.T) {
	//: 5
	srv := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Secret: "extremely-secret", Endpoint: "/webhook-endpoint"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, srv.Config.Endpoint, nil)
	req.Header.Set("X-Gitlab-Token", srv.Config.Secret)

	srv.webhookSecret(testHandler).ServeHTTP(w, req)

	want := 200
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}
}

func TestEmptyWebhookSecret(t *testing.T) {
	srv := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Secret: "", Endpoint: "/webhook-endpoint"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, srv.Config.Endpoint, nil)
	req.Header.Set("X-Gitlab-Token", srv.Config.Secret)

	srv.webhookSecret(testHandler).ServeHTTP(w, req)

	want := 200
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}
}

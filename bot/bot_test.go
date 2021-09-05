package bot

import (
	"bufio"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := loadConfig("config.yaml", "../")
	if err != nil {
		t.Errorf("config failed: %v", err)
	}
	want := "https://bot-bot.com"
	got := c.BotServer
	if got != want {
		t.Errorf("config not correctly loaded, expected %s, but got: %s", want, got)
	}
}

func TestFailedConfig(t *testing.T) {
	_, err := loadConfig("invalid", ".")
	if err == nil {
		t.Errorf("expected a failed config error, but got: %v", err)
	}
}

func TestNew(t *testing.T) {
	_, err := New("config.yaml", "../")
	if err != nil {
		t.Errorf("failed to init bot, %v", err)
	}
}

func TestPing(t *testing.T) {
	srv := Bot{
		Router: chi.NewRouter(),
		Logger: zerolog.Logger{},
		Config: &Config{Secret: "extremely-secret", Endpoint: "/webhook-endpoint"},
	}
	srv.routes(srv.Router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	srv.ServeHTTP(w, req)

	want := 200
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}
	var msg Message
	err := json.NewDecoder(w.Body).Decode(&msg)
	if err != nil {
		t.Errorf("response couldn't be decoded: %v", err)
	}

	if msg.Msg != "pong" {
		t.Errorf("expected pong response, but got: %v", err)
	}
}

func TestDecodeWebhook(t *testing.T) {
	body := strings.NewReader(`{"object_kind": "merge_request"}`)
	got, _ := decodeWebhook(bufio.NewReader(body))
	if got.ObjectKind != MergeRequest {
		t.Errorf("expected %s, but got: %v", MergeRequest, got.ObjectKind)
	}
}

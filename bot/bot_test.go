package bot

import (
	"bufio"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gitlab.com/jonny7/quetzal/policy"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := loadConfig("../examples/config.yaml")
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
	_, err := loadConfig("invalid")
	if err == nil {
		t.Errorf("expected a failed config error, but got: %v", err)
	}
}

func TestLoadPolicies(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}
	_ = b.loadPolicies("../examples/.policies.yaml")

	if len(b.Config.Policies) != 2 {
		t.Errorf("expected 2 policies, but got: %v", len(b.Config.Policies))
	}
	if b.Config.Policies[0].Name != "assign MR" {
		t.Errorf("expected name to be `%s`, but got: %s", "assign mr", b.Config.Policies[0].Name)
	}
	if b.Config.Policies[1].Actions.RemoveLabels[0] != "done" {
		t.Errorf("expected name to be `%s`, but got: %s", "done", b.Config.Policies[1].Actions.RemoveLabels[0])
	}
}

func TestNew(t *testing.T) {
	_, err := New("../examples/config.yaml", "../examples/.policies.yaml")
	if err != nil {
		t.Errorf("failed to init bot, %v", err)
	}
}

func TestPing(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: zerolog.Logger{},
		Config: &Config{Secret: "extremely-secret", Endpoint: "/webhook-endpoint"},
	}
	b.routes(b.Router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	b.ServeHTTP(w, req)

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

func TestPolicies(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}
	_ = b.loadPolicies("../examples/.policies.yaml")

	b.routes(b.Router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/policies", nil)

	b.ServeHTTP(w, req)

	var msg []policy.Policy
	err := json.NewDecoder(w.Body).Decode(&msg)
	if err != nil {
		t.Errorf("response couldn't be decoded: %v", err)
	}

	if len(msg) != 2 {
		t.Errorf("expected pong response, but got: %v", err)
	}
}

func TestReload(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint", policyPath: "../examples/.policies.yaml"},
	}

	b.routes(b.Router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reload", nil)

	b.ServeHTTP(w, req)

	want := 200
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}

	if len(b.Config.Policies) != 2 {
		t.Errorf("expected 2 policies, recevied %d", len(b.Config.Policies))
	}
}

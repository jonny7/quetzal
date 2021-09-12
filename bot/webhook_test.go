package bot

import (
	"bufio"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"io"
	"strings"
	"testing"
)

func TestDecodeWebhook(t *testing.T) {
	body := strings.NewReader(`{"object_kind": "merge_request"}`)
	got, _ := decodeWebhook(bufio.NewReader(body))
	if got.ObjectKind != MergeRequest {
		t.Errorf("expected %s, but got: %v", MergeRequest, got.ObjectKind)
	}
}

func TestDryRun(t *testing.T) {
	//: 4
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint", dryRun: true},
	}

	p := `policies:
  - name: dummy policy`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	webhook := Webhook{ObjectKind: Tag}
	got, err := webhook.handleEvent(&b)
	if got != nil && err != nil {
		t.Errorf("expected nil for got and err, but received %v, %v", got, err)
	}
}

func TestNoTriggeredPolicies(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	webhook := Webhook{ObjectKind: Comment}
	got, err := webhook.handleEvent(&b)
	if got != nil && err != nil {
		t.Errorf("expected nil for got and err, but received %v, %v", got, err)
	}
}

package bot

import (
	"bufio"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gitlab.com/jonny7/quetzal/policy"
	"io"
	"strings"
	"testing"
)

func TestDecodeWebhook(t *testing.T) {
	//: 7
	body := strings.NewReader(`{"object_kind": "merge_request"}`)
	got, _ := decodeWebhook(bufio.NewReader(body))
	if got.ObjectKind != policy.MergeRequest {
		t.Errorf("expected %s, but got: %v", policy.MergeRequest, got.ObjectKind)
	}
}

func TestDryRun(t *testing.T) {
	//: 4
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint", DryRun: true},
	}

	p := `policies:
  - name: dummy policy`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	webhook := Webhook{ObjectKind: policy.Tag}
	got, err := webhook.handleEvent(&b)
	if got != nil && err != nil {
		t.Errorf("expected nil for got and err, but received %v, %v", got, err)
	}
}

func TestNoFilteredPoliciesViaWebhook(t *testing.T) {
	//: 7
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	webhook := Webhook{ObjectKind: policy.Comment}
	got, err := webhook.handleEvent(&b)
	if got != nil && err != nil {
		t.Errorf("expected nil for got and err, but received %v, %v", got, err)
	}
}

func TestInvalidResourceValidation(t *testing.T) {
	//: 6
	webhook := Webhook{ObjectKind: policy.Deployment}
	got := webhook.ObjectKind.Validate()
	if got != nil {
		t.Errorf("expected no error as `%s` is a valid EventType", webhook.ObjectKind)
	}
}

func TestResourceValidation(t *testing.T) {
	//: 6
	invalid := policy.EventType("Invalid")
	webhook := Webhook{ObjectKind: invalid}
	got := webhook.ObjectKind.Validate()
	if got == nil {
		t.Errorf("expected an error as `%s` is a valid EventType", webhook.ObjectKind)
	}
}

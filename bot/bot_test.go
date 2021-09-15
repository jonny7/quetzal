package bot

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/jonny7/quetzal/policy"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadPolicies(t *testing.T) {
	//: 9
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	reader, _ := createReader("../examples/.policies.yaml")
	_ = b.loadPolicies(reader)

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
	_, err := New(Config{Endpoint: "/webhook/endpoint"}, "../examples/.policies.yaml")
	if err != nil {
		t.Errorf("failed to init bot, %v", err)
	}
}

func TestPing(t *testing.T) {
	//: 10
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
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

func TestPolicies(t *testing.T) {
	//: 9
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
  - name: dummy policy
    resource: issue
  - name: respond to mention
    resource: Note`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

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
		t.Errorf("expected 2 policies returned, but got: %v", err)
	}
}

func TestReload(t *testing.T) {
	//: 8
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint", PolicyPath: "../examples/.policies.yaml"},
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

func TestReloadInvalidPath(t *testing.T) {
	//: 8
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint", PolicyPath: "./invalid/.policies.yaml"},
	}

	b.routes(b.Router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reload", nil)

	b.ServeHTTP(w, req)

	want := 500
	got := w.Code

	if got != want {
		t.Errorf("expected %d, but got: %d", want, got)
	}
}

func TestFilteredEventPolicies(t *testing.T) {
	//:7
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
 - name: dummy policy
   resource: Issue Hook`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	got := b.filterPoliciesByEventType(gitlab.EventTypeIssue)
	if got[0].Name != "dummy policy" {
		t.Errorf("expected dummy policy returned")
	}
}

func TestValidatePoliciesDateProperties(t *testing.T) {
	//: 1
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
  - name: assign MR
    resource: merge_request
    conditions:
      date:
        attribute: not_a_valid_input`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	err := b.validatePolicies()
	if err == nil {
		t.Errorf("expected an error here as `not_a_valid_input` is not valid")
	}
}

func TestNoteConditionParsed(t *testing.T) {
	//: 12,7
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
 - name: show bot options
   resource: Note Hook
   conditions:
     note:
       noteType: Issue
       mentions:
         - botuser
       command: show -help`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	got := b.filterPoliciesByEventType(gitlab.EventTypeNote)
	if len(got) != 1 {
		t.Errorf("expected the 1 policy to be returned")
	}
}

func TestNoteConditionNoteTypeFilteredNil(t *testing.T) {
	//: 12,7, 13, 14
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
 - name: show bot options
   resource: Note Hook
   conditions:
     note:
       noteType: Issue
       mentions:
         - botuser
       command: show -help
 - name: some other action
   resource: Note Hook
   conditions:
     note:
       mentions:
         - botuser
       command: show -help`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	webhook := Webhook{
		eventType: gitlab.EventTypeNote,
		event: gitlab.IssueCommentEvent{
			ObjectAttributes: struct {
				ID           int            `json:"id"`
				Note         string         `json:"note"`
				NoteableType string         `json:"noteable_type"`
				AuthorID     int            `json:"author_id"`
				CreatedAt    string         `json:"created_at"`
				UpdatedAt    string         `json:"updated_at"`
				ProjectID    int            `json:"project_id"`
				Attachment   string         `json:"attachment"`
				LineCode     string         `json:"line_code"`
				CommitID     string         `json:"commit_id"`
				NoteableID   int            `json:"noteable_id"`
				System       bool           `json:"system"`
				StDiff       []*gitlab.Diff `json:"st_diff"`
				URL          string         `json:"url"`
			}{NoteableType: "Issue"},
		},
	}

	policies := b.filterPoliciesByEventType(gitlab.EventTypeNote)
	got, _ := webhook.filterAdditionalEventType(policies)

	if len(got) != 2 {
		t.Errorf("expected the 1 policy to be returned, but got: %d", len(got))
	}
}

func TestNoteConditionNoteTypeFiltered(t *testing.T) {
	//: 12,7, 13, 14
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
 - name: show bot options
   resource: Note Hook
   conditions:
     note:
       noteType: Issue
       mentions:
         - botuser
       command: show -help
 - name: some other action
   resource: Note Hook
   conditions:
     note:
       noteType: Commit
       mentions:
         - botuser
       command: show -help`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
	webhook := Webhook{
		eventType: gitlab.EventTypeNote,
		event: gitlab.CommitCommentEvent{
			ObjectAttributes: struct {
				ID           int    `json:"id"`
				Note         string `json:"note"`
				NoteableType string `json:"noteable_type"`
				AuthorID     int    `json:"author_id"`
				CreatedAt    string `json:"created_at"`
				UpdatedAt    string `json:"updated_at"`
				ProjectID    int    `json:"project_id"`
				Attachment   string `json:"attachment"`
				LineCode     string `json:"line_code"`
				CommitID     string `json:"commit_id"`
				NoteableID   int    `json:"noteable_id"`
				System       bool   `json:"system"`
				StDiff       struct {
					Diff        string `json:"diff"`
					NewPath     string `json:"new_path"`
					OldPath     string `json:"old_path"`
					AMode       string `json:"a_mode"`
					BMode       string `json:"b_mode"`
					NewFile     bool   `json:"new_file"`
					RenamedFile bool   `json:"renamed_file"`
					DeletedFile bool   `json:"deleted_file"`
				} `json:"st_diff"`
			}{NoteableType: "Commit"},
		},
	}
	policies := b.filterPoliciesByEventType(gitlab.EventTypeNote)
	got, _ := webhook.filterAdditionalEventType(policies)

	if len(got) != 1 {
		t.Errorf("expected the 1 policy to be returned, but got: %d", len(got))
	}
}

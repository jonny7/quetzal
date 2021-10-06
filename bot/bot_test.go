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
    resource: Issue Hook
  - name: respond to mention
    resource: Note Hook`
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

//func TestFilteredEventPolicies(t *testing.T) {
//	//: 7
//	b := Bot{
//		Router: chi.NewRouter(),
//		Logger: &zerolog.Logger{},
//		Config: &Config{Endpoint: "/webhook-endpoint"},
//	}
//
//	p := `policies:
//  - name: dummy policy
//    resource: Issue Hook`
//	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
//
//	webhook := policy.Webhook{EventType: gitlab.EventTypeIssue}
//	preparedPolicies := b.preparePolicies()
//	filtered := webhook.FilterEvent(preparedPolicies)
//
//	var got []policy.Policy
//	for po := range filtered {
//		got = append(got, po)
//	}
//	if got[0].Name != "dummy policy" {
//		t.Errorf("expected dummy policy returned")
//	}
//}

func TestValidatePoliciesDateProperties(t *testing.T) {
	//: 1
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
  - name: assign MR
    resource: Merge Request Hook
    conditions:
      date:
        attribute: not_a_valid_input`
	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	prep := b.preparePolicies()

	err := <-b.validatePolicies(prep)
	if err == nil {
		t.Errorf("expected an error here as `not_a_valid_input` is not valid")
	}
}

//func TestNoteConditionParsed(t *testing.T) {
//	//: 12,7
//	b := Bot{
//		Router: chi.NewRouter(),
//		Logger: &zerolog.Logger{},
//		Config: &Config{Endpoint: "/webhook-endpoint"},
//	}
//
//	p := `policies:
// - name: show bot options
//   resource: Note Hook
//   conditions:
//     note:
//       noteType: Issue
//       mentions:
//         - botuser
//       command: show -help`
//	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
//
//	webhook := policy.Webhook{EventType: gitlab.EventTypeNote}
//	preparedPolicies := b.preparePolicies()
//
//	filtered := webhook.FilterEvent(preparedPolicies)
//
//	var got []policy.Policy
//	for po := range filtered {
//		got = append(got, po)
//	}
//
//	if len(got) != 1 {
//		t.Errorf("expected 1 policy to be returned, got: %d", len(got))
//	}
//}

func TestNoteConditionNoteTypeFilteredNil(t *testing.T) {
	//: 12,7,13,14
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

	webhook := policy.Webhook{
		EventType: gitlab.EventTypeNote,
		Event: gitlab.IssueCommentEvent{
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

	preparedPolicies := b.preparePolicies()
	filtered := webhook.FilterEvent(preparedPolicies)

	var got []policy.WebhookResult
	for po := range filtered {
		got = append(got, po)
	}

	if len(got) == 1 {
		t.Errorf("expected 2 policies to be returned, but got: %d", len(got))
	}
}

func TestChannelMerge(t *testing.T) {
	b := Bot{
		Router: chi.NewRouter(),
		Logger: &zerolog.Logger{},
		Config: &Config{Endpoint: "/webhook-endpoint"},
	}

	p := `policies:
- name: show bot options
  resource: Note Hook
- name: trigger build
  resource: Note Hook
- name: trigger release
  resource: Note Hook`

	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))

	webhook := policy.Webhook{
		EventType: gitlab.EventTypeNote,
		Event: gitlab.CommitCommentEvent{
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
			}{NoteableType: "Issue"},
		},
	}

	preparedPolicies := b.preparePolicies()
	workers := make([]<-chan policy.WebhookResult, 3)
	for i := 0; i < 3; i++ {
		workers[i] = webhook.FilterEvent(preparedPolicies)
	}

	merged := mergePolicies(workers...)
	var got int
	for range merged {
		got++
	}
	if got != 3 {
		t.Errorf("expected 3 records when merged, but got: %d", got)
	}
}

//func TestNotableTypeFilter(t *testing.T) {
//	//: 12,7,13,14
//	b := Bot{
//		Router: chi.NewRouter(),
//		Logger: &zerolog.Logger{},
//		Config: &Config{Endpoint: "/webhook-endpoint"},
//	}
//
//	p := `policies:
//- name: show bot options
//  resource: Note Hook
//  conditions:
//    note:
//      noteType: Issue
//      mentions:
//        - botuser
//      command: show -help
//- name: some other action
//  resource: Note Hook
//  conditions:
//    note:
//      noteType: Commit
//      mentions:
//        - botuser
//      command: show -help`
//
//	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
//	webhook := policy.Webhook{
//		EventType: gitlab.EventTypeNote,
//		Event: gitlab.CommitCommentEvent{
//			ObjectAttributes: struct {
//				ID           int    `json:"id"`
//				Note         string `json:"note"`
//				NoteableType string `json:"noteable_type"`
//				AuthorID     int    `json:"author_id"`
//				CreatedAt    string `json:"created_at"`
//				UpdatedAt    string `json:"updated_at"`
//				ProjectID    int    `json:"project_id"`
//				Attachment   string `json:"attachment"`
//				LineCode     string `json:"line_code"`
//				CommitID     string `json:"commit_id"`
//				NoteableID   int    `json:"noteable_id"`
//				System       bool   `json:"system"`
//				StDiff       struct {
//					Diff        string `json:"diff"`
//					NewPath     string `json:"new_path"`
//					OldPath     string `json:"old_path"`
//					AMode       string `json:"a_mode"`
//					BMode       string `json:"b_mode"`
//					NewFile     bool   `json:"new_file"`
//					RenamedFile bool   `json:"renamed_file"`
//					DeletedFile bool   `json:"deleted_file"`
//				} `json:"st_diff"`
//			}{NoteableType: "Commit"},
//		},
//	}
//
//	preparedPolicies := b.preparePolicies()
//	filtered := webhook.FilterEvent(preparedPolicies)
//
//	var got []policy.Policy
//	for po := range filtered {
//		got = append(got, po)
//	}
//
//	if len(got) != 1 {
//		t.Errorf("expected 1 policy to be returned, but got: %d", len(got))
//	}
//}

//func TestFilterAdditionalType(t *testing.T) {
//	//: 12,7,13,14
//	b := Bot{
//		Router: chi.NewRouter(),
//		Logger: &zerolog.Logger{},
//		Config: &Config{Endpoint: "/webhook-endpoint"},
//	}
//
//	p := `policies:
//- name: show bot options
//  resource: Issue Hook
//  conditions:
//    state: opened
//- name: some other action
//  resource: Note Hook
//  conditions:
//    note:
//      noteType: Commit
//      mentions:
//        - botuser
//      command: show -help`
//	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
//	webhook := policy.Webhook{
//		EventType: gitlab.EventTypeNote,
//		Event: gitlab.CommitCommentEvent{
//			ObjectAttributes: struct {
//				ID           int    `json:"id"`
//				Note         string `json:"note"`
//				NoteableType string `json:"noteable_type"`
//				AuthorID     int    `json:"author_id"`
//				CreatedAt    string `json:"created_at"`
//				UpdatedAt    string `json:"updated_at"`
//				ProjectID    int    `json:"project_id"`
//				Attachment   string `json:"attachment"`
//				LineCode     string `json:"line_code"`
//				CommitID     string `json:"commit_id"`
//				NoteableID   int    `json:"noteable_id"`
//				System       bool   `json:"system"`
//				StDiff       struct {
//					Diff        string `json:"diff"`
//					NewPath     string `json:"new_path"`
//					OldPath     string `json:"old_path"`
//					AMode       string `json:"a_mode"`
//					BMode       string `json:"b_mode"`
//					NewFile     bool   `json:"new_file"`
//					RenamedFile bool   `json:"renamed_file"`
//					DeletedFile bool   `json:"deleted_file"`
//				} `json:"st_diff"`
//			}{NoteableType: "Commit"},
//		},
//	}
//
//	preparedPolicies := b.preparePolicies()
//	filtered := webhook.FilterEvent(preparedPolicies)
//
//	var got []policy.Policy
//	for po := range filtered {
//		got = append(got, po)
//	}
//
//	if len(got) != 1 {
//		t.Errorf("expected 1 policy to be returned, but got: %d", len(got))
//	}
//}
//
//func TestProcessWebhookNoConcurrencyErrors(t *testing.T) {
//	//:
//	b := Bot{
//		Router: chi.NewRouter(),
//		Logger: &zerolog.Logger{},
//		Config: &Config{Endpoint: "/webhook-endpoint"},
//	}
//
//	p := `policies:
// - name: show bot options
//   resource: Issue Hook
//   conditions:
//     state: opened
// - name: some other action
//   resource: Note Hook
//   conditions:
//     note:
//       noteType: Commit
//       mentions:
//         - botuser
//       command: show -help`
//	_ = b.loadPolicies(io.NopCloser(strings.NewReader(p)))
//
//	b.routes(b.Router)
//	w := httptest.NewRecorder()
//	payload := `{
//  "object_kind": "note",
//  "user": {
//    "id": 1,
//    "name": "Administrator",
//    "username": "root",
//    "avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40\u0026d=identicon",
//    "email": "admin@example.com"
//  },
//  "project_id": 5,
//  "project":{
//    "id": 5,
//    "name":"Gitlab Test",
//    "description":"Aut reprehenderit ut est.",
//    "web_url":"http://example.com/gitlabhq/gitlab-test",
//    "avatar_url":null,
//    "git_ssh_url":"git@example.com:gitlabhq/gitlab-test.git",
//    "git_http_url":"http://example.com/gitlabhq/gitlab-test.git",
//    "namespace":"GitlabHQ",
//    "visibility_level":20,
//    "path_with_namespace":"gitlabhq/gitlab-test",
//    "default_branch":"master",
//    "homepage":"http://example.com/gitlabhq/gitlab-test",
//    "url":"http://example.com/gitlabhq/gitlab-test.git",
//    "ssh_url":"git@example.com:gitlabhq/gitlab-test.git",
//    "http_url":"http://example.com/gitlabhq/gitlab-test.git"
//  },
//  "repository":{
//    "name": "Gitlab Test",
//    "url": "http://example.com/gitlab-org/gitlab-test.git",
//    "description": "Aut reprehenderit ut est.",
//    "homepage": "http://example.com/gitlab-org/gitlab-test"
//  },
//  "object_attributes": {
//    "id": 1243,
//    "note": "This is a commit comment. How does this work?",
//    "noteable_type": "Commit",
//    "author_id": 1,
//    "created_at": "2015-05-17 18:08:09 UTC",
//    "updated_at": "2015-05-17 18:08:09 UTC",
//    "project_id": 5,
//    "attachment":null,
//    "line_code": "bec9703f7a456cd2b4ab5fb3220ae016e3e394e3_0_1",
//    "commit_id": "cfe32cf61b73a0d5e9f13e774abde7ff789b1660",
//    "noteable_id": null,
//    "system": false,
//    "st_diff": {
//      "diff": "--- /dev/null\n+++ b/six\n@@ -0,0 +1 @@\n+Subproject commit 409f37c4f05865e4fb208c771485f211a22c4c2d\n",
//      "new_path": "six",
//      "old_path": "six",
//      "a_mode": "0",
//      "b_mode": "160000",
//      "new_file": true,
//      "renamed_file": false,
//      "deleted_file": false
//    },
//    "url": "http://example.com/gitlab-org/gitlab-test/commit/cfe32cf61b73a0d5e9f13e774abde7ff789b1660#note_1243"
//  },
//  "commit": {
//    "id": "cfe32cf61b73a0d5e9f13e774abde7ff789b1660",
//    "message": "Add submodule\n\nSigned-off-by: Example User \u003cuser@example.com.com\u003e\n",
//    "timestamp": "2014-02-27T10:06:20+02:00",
//    "url": "http://example.com/gitlab-org/gitlab-test/commit/cfe32cf61b73a0d5e9f13e774abde7ff789b1660",
//    "author": {
//      "name": "Example User",
//      "email": "user@example.com"
//    }
//  }
//}`
//	req, _ := http.NewRequest(http.MethodPost, "/webhook-endpoint", bytes.NewBuffer([]byte(payload)))
//	req.Header.Set("content-type", "application/json")
//	req.Header.Set("X-Gitlab-Event", "Note Hook")
//
//	b.ServeHTTP(w, req)
//
//	want := 200
//	got := w.Code
//
//	if got != want {
//		t.Errorf("expected %d, but got: %d", want, got)
//	}
//	var returnedPolicies []policy.Policy
//
//	err := json.NewDecoder(w.Body).Decode(&returnedPolicies)
//	if err != nil {
//		t.Errorf("response couldn't be decoded: %v", err)
//	}
//
//	if len(returnedPolicies) != 1 {
//		t.Errorf("expected 1 policy to be returned, but got: %d", len(returnedPolicies))
//	}
//}

package policy

import (
	"encoding/json"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)
	// client is the Gitlab client being tested.
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}
	return mux, server, client
}

func stubMergeEventAdaptor() MergeEventAdaptor {
	me := MergeEventAdaptor{}
	me.Project.ID = 1
	me.ObjectAttributes.IID = 234
	return me
}

func stubUpdatedMergeEventEndPoint(m MergeEventAdaptor) string {
	return fmt.Sprintf("/api/v4/projects/%d/merge_requests/%d", m.Project.ID, m.ObjectAttributes.IID)
}

// teardown closes the test HTTP server.
func teardown(server *httptest.Server) {
	server.Close()
}

func TestExecuteMethods(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	me := stubMergeEventAdaptor()
	endpoint := stubUpdatedMergeEventEndPoint(me)
	action := Action{Labels: []string{"approved"}, Mention: []string{"@jonny"}, Comment: "this has been automatically labelled"}

	// response object for MergeRequest Updates, set to the action Labels
	m := new(gitlab.MergeRequest)
	m.Labels = action.Labels

	// mock response for updateMerge Req
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&m)
		if err != nil {
			t.Errorf("failed to encode response")
		}
		return
	})

	// response object for adding a Note to MergeRequests
	n := new(gitlab.Note)
	n.Body = action.commentate()

	// mock response for add Note
	noteEndpoint := fmt.Sprintf("/api/v4/projects/%d/merge_requests/%d/notes", me.Project.ID, me.ObjectAttributes.IID)
	mux.HandleFunc(noteEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		noteErr := json.NewEncoder(w).Encode(&n)
		if noteErr != nil {
			t.Errorf("failed to encode response")
		}
		return
	})

	data := []struct {
		name     string
		updateFn gitLabUpdateFn
		expected string
		errMsg   string
	}{
		{name: "Execute Labels", updateFn: me.executeLabels, expected: endpoint, errMsg: "expected endpoint to be %s but got %s"},
		{name: "Execute Notes", updateFn: me.executeNote, expected: noteEndpoint, errMsg: "expected endpoint to be %s but got %s"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got, err := d.updateFn(action, client)
			if err != nil {
				t.Errorf("expected no error to occur in mock request")
			}
			if got != d.expected {
				t.Errorf(d.errMsg, d.expected, got)
			}
		})
	}

	t.Run("TestExecute", func(t *testing.T) {
		updatedResults := me.execute(action, client)
		if len(updatedResults) != 2 {
			t.Errorf("expected 2 updates to occur")
		}
	})
}

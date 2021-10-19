package policy

import "github.com/xanzy/go-gitlab"

// gitLabUpdateFn is allows for possible multiple action requests to
// be stacked up and executed as part of an array
type gitLabUpdateFn func(action Action, client *gitlab.Client) (string, error)

// GitLabUpdateResult reports back to the caller the series of events taken
// by the bot to update Gitlab
type GitLabUpdateResult struct {
	Action Action `json:"action"`
	// here we collect the Endpoint being called from the client to help provide
	// more info, without using reflection on a func to get the func name
	Endpoint string `json:"endpoint"`
	Error    string `json:"error"`
}

// GitLabAdaptor wraps the incoming hook so additional methods can be added
type GitLabAdaptor interface {
	Executor
	Stater
	Labeler
	Milestoner
}

// Preparer provides functionality that the GitLabAdaptor needs to determine what functionality
// is available to that type
type Preparer interface {
	updateLabels() bool
	updateState() bool
	addNote() bool
}

// Executor is how the updates to GitLab are done on a per-type basis
type Executor interface {
	prepareUpdates(action Preparer) []gitLabUpdateFn
	execute(action Action, client *gitlab.Client) []GitLabUpdateResult
}

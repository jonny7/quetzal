package policy

import "github.com/xanzy/go-gitlab"

// gitLabUpdateFn is allows for possible multiple action requests to
// be stacked up and executed as part of an array
type gitLabUpdateFn func(action Action, client *gitlab.Client) (string, error)

// GitLabUpdateResult reports back to the caller the series of events taken
// by the bot to update Gitlab
type GitLabUpdateResult struct {
	action Action
	// here we collect the endpoint being called from the client to help provide
	// more info, without using reflection on a func to get the func name
	endpoint string
	error    error
}

// GitLabAdaptor wraps the incoming hook so
// additional methods can be added
type GitLabAdaptor interface {
	Executor
}

// Executor is how the updates to GitLab are done on a per-type basis
type Executor interface {
	prepareUpdates(action Action) []gitLabUpdateFn
	execute(action Action, client *gitlab.Client) []GitLabUpdateResult
}

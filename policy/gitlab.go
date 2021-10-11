package policy

import "github.com/xanzy/go-gitlab"

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

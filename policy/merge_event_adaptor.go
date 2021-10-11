package policy

import "github.com/xanzy/go-gitlab"

// MergeEventAdaptor wraps the gitlab.MergeEvent
type MergeEventAdaptor struct {
	gitlab.MergeEvent
}

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

// prepare updates goes through the action list and determines what update requests are required.
func (m MergeEventAdaptor) prepareUpdates(action Action) []gitLabUpdateFn {
	var executables []gitLabUpdateFn
	// update status and labels
	if action.updateLabels() {
		executables = append(executables, m.executeLabels)
	}
	if action.addNote() {
		executables = append(executables, m.executeNote)
	}
	return executables
}

func (m MergeEventAdaptor) execute(action Action, client *gitlab.Client) []GitLabUpdateResult {
	updates := m.prepareUpdates(action)
	var updateResults []GitLabUpdateResult
	for _, u := range updates {
		endpoint, err := u(action, client)
		updateResults = append(updateResults, GitLabUpdateResult{
			action:   action,
			endpoint: endpoint,
			error:    err,
		})
	}
	return updateResults
}

func (m MergeEventAdaptor) executeLabels(action Action, client *gitlab.Client) (string, error) {
	opt := gitlab.UpdateMergeRequestOptions{
		AddLabels:    action.Labels,
		RemoveLabels: action.RemoveLabels,
	}
	_, resp, err := client.MergeRequests.UpdateMergeRequest(m.Project.ID, m.ObjectAttributes.IID, &opt)
	if err != nil {
		return resp.Response.Request.URL.Path, err
	}
	return resp.Response.Request.URL.Path, nil
}

func (m MergeEventAdaptor) executeNote(action Action, client *gitlab.Client) (string, error) {
	note := action.commentate()
	_, resp, err := client.Notes.CreateMergeRequestNote(m.Project.ID, m.ObjectAttributes.IID, &gitlab.CreateMergeRequestNoteOptions{Body: &note})
	if err != nil {
		return resp.Response.Request.URL.Path, err
	}
	return resp.Response.Request.URL.Path, nil
}

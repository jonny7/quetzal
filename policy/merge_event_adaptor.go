package policy

import "github.com/xanzy/go-gitlab"

// MergeEventAdaptor wraps the gitlab.MergeEvent
type MergeEventAdaptor struct {
	gitlab.MergeEvent
}

func (m MergeEventAdaptor) state() *string {
	return &m.ObjectAttributes.Action
}

// prepare updates goes through the action list and determines what update requests are required.
func (m MergeEventAdaptor) prepareUpdates(action Action) []gitLabUpdateFn {
	var executables []gitLabUpdateFn
	if action.updateLabels() {
		executables = append(executables, m.executeLabels)
	}
	if action.updateState() {
		executables = append(executables, m.executeStatus)
	}
	if action.addNote() {
		executables = append(executables, m.executeNote)
	}
	return executables
}

func (m MergeEventAdaptor) execute(action Action, client *gitlab.Client) []GitLabUpdateResult {
	updates := m.prepareUpdates(action)
	var updateResults []GitLabUpdateResult
	for _, update := range updates {
		endpoint, err := update(action, client)
		result := GitLabUpdateResult{Action: action, Endpoint: endpoint}
		if err != nil {
			result.Error = err.Error()
		}
		updateResults = append(updateResults, result)
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

func (m MergeEventAdaptor) executeStatus(action Action, client *gitlab.Client) (string, error) {
	opt := gitlab.UpdateMergeRequestOptions{
		StateEvent: &action.Status,
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

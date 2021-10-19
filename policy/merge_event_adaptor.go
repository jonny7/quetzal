package policy

import "github.com/xanzy/go-gitlab"

// MergeEventAdaptor wraps the gitlab.MergeEvent
type MergeEventAdaptor struct {
	gitlab.MergeEvent
}

func (m MergeEventAdaptor) state() []string {
	return []string{m.ObjectAttributes.Action}
}

func (m MergeEventAdaptor) labels() []string {
	var labels []string
	for _, label := range m.Labels {
		labels = append(labels, label.Name)
	}
	return sliceLower(labels)
}

func (m MergeEventAdaptor) milestone() int {
	return m.ObjectAttributes.MilestoneID
}

// prepare updates goes through the action list and determines what update requests are required.
func (m MergeEventAdaptor) prepareUpdates(action Preparer) []gitLabUpdateFn {
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
		AddLabels:    action.Labels.Labels,
		RemoveLabels: action.RemoveLabels,
	}
	_, resp, err := client.MergeRequests.UpdateMergeRequest(m.Project.ID, m.ObjectAttributes.IID, &opt)
	if err != nil {
		return resp.Response.Request.URL.Path, err
	}
	return resp.Response.Request.URL.Path, nil
}

func (m MergeEventAdaptor) executeStatus(action Action, client *gitlab.Client) (string, error) {
	if string(action.Status) == string(mergeRequestStateApproved) {
		opt := gitlab.ApproveMergeRequestOptions{SHA: &m.ObjectAttributes.LastCommit.ID}
		_, resp, err := client.MergeRequestApprovals.ApproveMergeRequest(m.Project.ID, m.ObjectAttributes.IID, &opt)
		if err != nil {
			return resp.Response.Request.URL.Path, err
		}
		return resp.Response.Request.URL.Path, nil
	}
	opt := gitlab.UpdateMergeRequestOptions{
		StateEvent: (*string)(&action.Status),
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

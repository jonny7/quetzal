package bot

import (
	"github.com/xanzy/go-gitlab"
	"gitlab.com/jonny7/quetzal/policy"
)

type Webhook struct {
	eventType gitlab.EventType
	event     interface{}
}

// filterEvent processing policies against the incoming hook and only returns policies
// that are valid for this webhook event.
func (w *Webhook) filterEvent(in <-chan policy.Policy, out chan<- policy.Policy) {
	for pol := range in {
		if pol.Resource != w.eventType {
			continue
		}
		switch ev := w.event.(type) {
		//case gitlab.BuildEvent:
		//case gitlab.DeploymentEvent:
		//case gitlab.IssueEvent:
		//case gitlab.JobEvent:
		//case gitlab.MergeEvent:
		//case gitlab.PipelineEvent:
		//case gitlab.PushEvent:
		//case gitlab.ReleaseEvent:
		//case gitlab.TagEvent:
		//case gitlab.WikiPageEvent:
		case gitlab.CommitCommentEvent:
			if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				out <- pol
			}
		case gitlab.MergeCommentEvent:
		case gitlab.IssueCommentEvent:
			if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				out <- pol
			}
		case gitlab.SnippetCommentEvent:
		default:
			out <- pol
		}
	}
	close(out)
}

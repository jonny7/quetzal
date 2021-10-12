package policy

import "fmt"

// Action struct houses how an eligible webhook
// event should be responded to
type Action struct {
	// HTTP contains allows a user to call an API endpoint
	// and use that in subsequent actions
	//HTTP *HTTP `yaml:"http,omitempty"` @todo
	// Labels identifies which labels to add to an issue
	Labels []string `yaml:"labels,omitempty"`
	// RemoveLabels defines what labels to remove
	RemoveLabels []string `yaml:"removeLabels,omitempty"`
	// Status sets the status of the issue
	Status string `yaml:"status,omitempty"`
	// Mention allows the bot to mention certain users
	Mention []string `yaml:"mention,omitempty"`
	// Comment will leave a comment on said issue
	Comment string `yaml:"comment,omitempty"`
}

// commentate builds a note body with the mentions and comment
// concatenated
func (a Action) commentate() string {
	var out string
	for _, m := range a.Mention {
		out += fmt.Sprintf("@%s ", m)
	}
	if len(a.Mention) == 0 {
		return a.Comment
	}
	return fmt.Sprintf("%s%s", out, a.Comment)
}

func (a Action) updateLabels() bool {
	if a.RemoveLabels != nil || a.Labels != nil {
		return true
	}
	return false
}

func (a Action) addNote() bool {
	if a.Mention != nil || a.Comment != "" {
		return true
	}
	return false
}

//func (a Action) buildResult()
//// returns... The type of object for update, so will need an interace maybe.
//// We know the action, but not at compile time.
//// so need an interface
//// filterEvent, we get the type, the policy checks conditions returns a bool
//// if true GitlabAdaptorType.Execute(Actioner interface, client *gitlab.Client)
///*
//We set MergeRequest's status and labels and use the client.MergeRequest.Update() to send
//Have look to look at how the Mention and comment work, it might not be possible on all types
//roll all Policy Actions into one action. Maybe there is no comment ona type
//
//So ultimately, need to validate that a Policy Resource can have a certain action taken
// */
// */
// */
// */
//*/
//
//}

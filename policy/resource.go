package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

// Resource embeds a gitlab.EventType
type Resource struct {
	EventType gitlab.EventType `yaml:"resource"`
}

// Validate of Resource ensures that the type is that of a gitlab.EventType
func (r Resource) validate() error {
	switch r.EventType {
	case gitlab.EventTypeBuild, gitlab.EventTypeDeployment, gitlab.EventTypeIssue, gitlab.EventConfidentialIssue, gitlab.EventTypeJob, gitlab.EventTypeMergeRequest, gitlab.EventTypeNote, gitlab.EventConfidentialNote, gitlab.EventTypePipeline, gitlab.EventTypePush, gitlab.EventTypeRelease, gitlab.EventTypeSystemHook, gitlab.EventTypeTagPush, gitlab.EventTypeWikiPage:
		return nil
	}
	return fmt.Errorf("`policy:resource` allowed options are: `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`. But received: %v", gitlab.EventTypeBuild,
		gitlab.EventTypeDeployment,
		gitlab.EventTypeIssue,
		gitlab.EventConfidentialIssue,
		gitlab.EventTypeJob,
		gitlab.EventTypeMergeRequest,
		gitlab.EventTypeNote,
		gitlab.EventConfidentialNote,
		gitlab.EventTypePipeline,
		gitlab.EventTypePush,
		gitlab.EventTypeRelease,
		gitlab.EventTypeSystemHook,
		gitlab.EventTypeTagPush,
		gitlab.EventTypeWikiPage, r.EventType)
}

func (r Resource) matcher(hook Webhook) bool {
	return hook.EventType == r.EventType
}

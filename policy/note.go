package policy

import (
	"fmt"
	"time"
)

// NoteType is the type of note: Commit, MergeRequest, Issue, Snippet
type NoteType string

const (
	NoteCommit       NoteType = "Commit"
	NoteMergeRequest NoteType = "MergeRequest"
	NoteIssue        NoteType = "Issue"
	NoteSnippet      NoteType = "Snippet"
)

type Note struct {
	// Type is the NoteType of the note from GitLab. If you need to narrow down
	// the type of note then use this, if left blank, then it will apply to all note types
	Type *NoteType `yaml:"noteType"`
	// Mentions looks for user's mentioned in the note
	Mentions []string `yaml:"mentions"`
	// Command is the specified string to look for if needed.
	Command string `yaml:"command"`
}

func (n Note) ConditionsMet() bool {
	time.Sleep(2 * time.Second)
	return true
}

// Validate confirms that the user provided NoteType is of an expected type
func (n NoteType) Validate() error {
	switch n {
	case NoteCommit, NoteIssue:
		return nil
	case NoteMergeRequest, NoteSnippet:
		return fmt.Errorf("MergeRequest and Snippet notes are not current supported")
	}

	return fmt.Errorf("the provided NoteType of %s is invalid", n)
}

func (n NoteType) ToString() string {
	return string(n)
}

//
//func (n *Note) ConditionMet() bool {
//	return true
//}
//
//func (n *Note) ConditionType() bool {
//	if n.Type == nil {
//		return true
//	}
//	if n.Type.ToString()
//}

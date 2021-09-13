package policy

import "fmt"

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
	Type NoteType `yaml:"noteType"`
	// Mentions looks for user's mentioned in the note
	Mentions []string `yaml:"mentions"`
	// Command is the specified string to look for if needed.
	Command string `yaml:"command"`
}

func (n NoteType) Validate() error {
	switch n {
	case NoteCommit, NoteMergeRequest, NoteIssue, NoteSnippet:
		return nil
	}
	return fmt.Errorf("the provided NoteType of %s is invalid", n)
}

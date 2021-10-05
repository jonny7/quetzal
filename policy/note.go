package policy

import (
	"fmt"
	"strings"
	"sync"
)

// NoteType is the type of note: Commit, MergeRequest, Issue, Snippet
type NoteType string

// Mentions is an array of users mentioned in a comment
type Mentions []string

// Command is a string backed type for a given command to respond to
type Command string

const (
	// NoteCommit are comments on Commits
	NoteCommit NoteType = "Commit"
	// NoteMergeRequest are comments on MergeRequests
	NoteMergeRequest NoteType = "MergeRequest"
	// NoteIssue are comments on Issues
	NoteIssue NoteType = "Issue"
	// NoteSnippet are comments on Snippets
	NoteSnippet NoteType = "Snippet"
)

// Note represents a GitLab Note, which is essentially a comment on
// a series of different scenarios and event types
type Note struct {
	// Type is the NoteType of the note from GitLab. If you need to narrow down
	// the type of note then use this, if left blank, then it will apply to all note types
	Type *NoteType `yaml:"noteType"`
	// Mentions looks for user's mentioned in the note
	Mentions Mentions `yaml:"mentions"`
	// Command is the specified string to look for if needed.
	Command Command `yaml:"command"`
}

func (m Mentions) conditionsMet(note Noter) bool {
	hookMentions := note.Mentions()
	if m == nil {
		return true
	}

	for _, person := range m {
		if existsInSlice(hookMentions, person) {
			return true
		}
	}
	return false
}

func (c Command) conditionsMet(note Noter) bool {
	if c == "" {
		return true
	}
	content, err := note.Note()
	if err != nil {
		return false
	}
	if strings.Contains(*content, string(c)) {
		return true
	}
	return false
}

// ConditionsMet confirms whether the webhook matches the policy for
// the Note type
func (n *Note) ConditionsMet(note Noter) bool {
	var wg sync.WaitGroup
	ch := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- n.Command.conditionsMet(note)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- n.Mentions.conditionsMet(note)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- n.Type.conditionsMet(note)
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	var valid = true

	for r := range ch {
		if !r {
			valid = false
			break
		}
	}
	return valid
}

func (n *NoteType) conditionsMet(note Noter) bool {
	nt, err := note.NoteType()
	if err != nil {
		return false
	}
	if n == nil {
		return true
	}
	if n.toString() == *nt {
		return true
	}
	return false
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const unsupportedError = Error("MergeRequest and Snippet notes are not current supported")

func (n *NoteType) invalidNoteType() error {
	return Error(fmt.Sprintf("the provided NoteType of %s is invalid", n.toString()))
}

// validate confirms that the user provided NoteType is of an expected type
func (n *NoteType) validate() error {
	if n == nil {
		return nil
	}
	switch *n {
	case NoteCommit, NoteIssue:
		return nil
	case NoteMergeRequest, NoteSnippet:
		return unsupportedError
	}
	return n.invalidNoteType()
}

// toString returns the NoteType as a string
func (n *NoteType) toString() string {
	if n == nil {
		return ""
	}
	return string(*n)
}

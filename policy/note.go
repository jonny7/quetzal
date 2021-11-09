package policy

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

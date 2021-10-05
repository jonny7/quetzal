package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestNoteTypeValidation(t *testing.T) {
	//: 21
	snippet := NoteSnippet
	noteIssue := NoteIssue
	invalidType := NoteType("invalid")
	data := []struct {
		name     string
		note     Note
		expected error
		errMsg   string
	}{
		{name: "Nil Note Type", note: Note{Type: nil}, expected: nil, errMsg: "expected nil as noteType can be nil"},
		{name: "Valid Note Type", note: Note{Type: &noteIssue}, expected: nil, errMsg: "expected nil as noteType is valid"},
		{name: "Unsupported Note Type", note: Note{Type: &snippet}, expected: unsupportedError, errMsg: "expected an error as type is unsupported"},
		{name: "Invalid Note Type", note: Note{Type: &invalidType}, expected: invalidType.invalidNoteType(), errMsg: "expected an error as type is invalid"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.note.Type.validate()
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestNoteTypeConditionsMet(t *testing.T) {
	//: 22
	unNotableType := MergeEventAdaptor{gitlab.MergeEvent{}}
	notableType := CommitCommentEventAdaptor{gitlab.CommitCommentEvent{}}
	notableType.ObjectAttributes.NoteableType = string(NoteIssue)
	validNoteType := NoteIssue
	validNoteType2 := NoteMergeRequest
	data := []struct {
		name     string
		note     Note
		hook     Noter
		expected bool
		errMsg   string
	}{
		{name: "Invalid Hook", note: Note{Type: &validNoteType}, hook: unNotableType, expected: false, errMsg: "MergeEvent does not have note field, so specifying a noteType is incorrect"},
		{name: "Nil Note Type", note: Note{Type: nil}, hook: notableType, expected: true, errMsg: "no NoteType is valid and conditions met should be true"},
		{name: "Matching Note Type", note: Note{Type: &validNoteType}, hook: notableType, expected: true, errMsg: "note types match so should be true"},
		{name: "Unmatched Note Type", note: Note{Type: &validNoteType2}, hook: notableType, expected: false, errMsg: "note types do not match so should be false"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.note.Type.conditionsMet(d.hook)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestMentionsConditionsMet(t *testing.T) {
	//: 22
	data := []struct {
		name     string
		note     Note
		incoming string
		expected bool
		errMsg   string
	}{
		{name: "No Policy Mentions", note: Note{Mentions: nil}, incoming: "A note with no mentions", expected: true, errMsg: "expected true as no mentions listed in policy"},
		{name: "Mentioned but not in Hook", note: Note{Mentions: []string{"@user123"}}, incoming: "A note with no mentions", expected: false, errMsg: "expected false as user123 listed in policy but not in hook"},
		{name: "Mentioned and is in Hook", note: Note{Mentions: []string{"@user123"}}, incoming: "A note which mentions @user123", expected: true, errMsg: "expected true as user123 listed in policy and appears in hook"},
		{name: "Mentioned multiple and some mentions are in Hook", note: Note{Mentions: []string{"@user123", "@jonny"}}, incoming: "A note which mentions @jonny", expected: true, errMsg: "expected true as at least one person listed in policy appears in hook"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			hook := CommitCommentEventAdaptor{gitlab.CommitCommentEvent{}}
			hook.ObjectAttributes.Note = d.incoming
			got := d.note.Mentions.conditionsMet(hook)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestCommandConditionsMet(t *testing.T) {
	//: 22
	note := CommitCommentEventAdaptor{gitlab.CommitCommentEvent{}}
	note.ObjectAttributes.Note = "some note, with no registered commands"
	note2 := note
	note2.ObjectAttributes.Note = "@bot -help please"
	noteInvalidHook := WikiEventAdaptor{gitlab.WikiPageEvent{}}
	data := []struct {
		name     string
		command  Command
		noter    Noter
		expected bool
		errMsg   string
	}{
		{name: "Empty Command", command: Command(""), noter: nil, expected: true, errMsg: "an empty policy command should be marked as conditions met being true"},
		{name: "Command in policy but not hook", command: Command("-help"), noter: note, expected: false, errMsg: "policy command not in hook, should return false"},
		{name: "Command in policy and hook", command: Command("-help"), noter: note2, expected: true, errMsg: "policy command in hook, should return true"},
		{name: "Command in policy with hook that has no notes", command: Command("-help"), noter: noteInvalidHook, expected: false, errMsg: "policy command but hook doesn't have note field, should return false"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.command.conditionsMet(d.noter)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestNoteConditionsMetIntegration(t *testing.T) {
	//: 22
	noteType := NoteMergeRequest
	p := Policy{
		Resource: Resource{
			EventType: gitlab.EventTypeNote,
		},
		Conditions: Condition{
			Note: &Note{
				Type:     &noteType,
				Mentions: Mentions{"@user123"},
				Command:  "-help",
			},
		},
	}

	hook := CommitCommentEventAdaptor{gitlab.CommitCommentEvent{}}
	hook.ObjectAttributes.Note = "@user123 -help"
	hook.ObjectAttributes.NoteableType = string(NoteMergeRequest)

	got := p.Conditions.Note.ConditionsMet(hook)
	if got != true {
		t.Errorf("expected true as policy and hook match")
	}
}

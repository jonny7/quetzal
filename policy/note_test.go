package policy

import "testing"

func TestNoteType(t *testing.T) {
	//: 12
	note := NoteIssue
	got := note.validate()
	if got != nil {
		t.Errorf("expected nil as `%s` is a valid note type, but received an error %v", note, got)
	}
}

func TestInvalidNoteType(t *testing.T) {
	//: 12
	note := NoteType("invalid")
	got := note.validate()
	if got == nil {
		t.Errorf("expected an error as `%s` is an invalid note type, but received no error %v", note, got)
	}
}

func TestUnsupportedNoteType(t *testing.T) {
	//: 12
	note := NoteMergeRequest
	got := note.validate()
	if got == nil {
		t.Errorf("expected an error as `%s` is an unsupported note type, but received no error %v", note, got)
	}
}

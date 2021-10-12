package policy

import (
	"testing"
)

func TestUpdateLabels(t *testing.T) {
	//: 4
	data := []struct {
		name     string
		action   Action
		expected bool
		errMsg   string
	}{
		{name: "No Labels or Remove Labels", action: Action{}, expected: false, errMsg: "expected false as Labels and Remove Labels are empty"},
		{name: "Labels but no Remove Labels", action: Action{Labels: []string{"added"}}, expected: true, errMsg: "expected true as Labels is not empty"},
		{name: "Remove Labels but No Labels", action: Action{RemoveLabels: []string{"removed"}}, expected: true, errMsg: "expected true as RemoveLabels is not empty"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			if d.action.updateLabels() != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	//: 5
	data := []struct {
		name     string
		action   Action
		expected bool
		errMsg   string
	}{
		{name: "No Mentions or Comments", action: Action{}, expected: false, errMsg: "expected false as Comments and Mentions are empty"},
		{name: "Mentions but no Comments", action: Action{Mention: []string{"jonny"}}, expected: true, errMsg: "expected true as Mentions are not empty"},
		{name: "Comments but no Mentions", action: Action{Comment: "I am commenting"}, expected: true, errMsg: "expected true as Comment is not empty"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			if d.action.addNote() != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestCommentCreationOnAdaptor(t *testing.T) {
	//: 6
	data := []struct {
		name     string
		action   Action
		expected string
		errMsg   string
	}{
		{name: "build comment with no mentions", action: Action{Comment: "GitLab"}, expected: "GitLab", errMsg: "expected %s but got %s"},
		{name: "build comment with 1 mention", action: Action{Comment: "GitLab", Mention: []string{"jonny"}}, expected: "@jonny GitLab", errMsg: "expected %s but got %s"},
		{name: "build comment with 2 mentions", action: Action{Comment: "GitLab", Mention: []string{"jonny", "bot"}}, expected: "@jonny @bot GitLab", errMsg: "expected %s but got %s"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			if d.action.commentate() != d.expected {
				t.Errorf(d.errMsg, d.expected, d.action.commentate())
			}
		})
	}
}

package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
	"time"
)

func TestConditionMetMilestoneNegative(t *testing.T) {
	//: 7,18
	adaptor := MergeEventAdaptor{gitlab.MergeEvent{
		ObjectKind: string(gitlab.EventTypeMergeRequest),
		ObjectAttributes: struct {
			ID                       int                 `json:"id"`
			TargetBranch             string              `json:"target_branch"`
			SourceBranch             string              `json:"source_branch"`
			SourceProjectID          int                 `json:"source_project_id"`
			AuthorID                 int                 `json:"author_id"`
			AssigneeID               int                 `json:"assignee_id"`
			AssigneeIDs              []int               `json:"assignee_ids"`
			Title                    string              `json:"title"`
			CreatedAt                string              `json:"created_at"`
			UpdatedAt                string              `json:"updated_at"`
			StCommits                []*gitlab.Commit    `json:"st_commits"`
			StDiffs                  []*gitlab.Diff      `json:"st_diffs"`
			MilestoneID              int                 `json:"milestone_id"`
			State                    string              `json:"state"`
			MergeStatus              string              `json:"merge_status"`
			TargetProjectID          int                 `json:"target_project_id"`
			IID                      int                 `json:"iid"`
			Description              string              `json:"description"`
			Position                 int                 `json:"position"`
			LockedAt                 string              `json:"locked_at"`
			UpdatedByID              int                 `json:"updated_by_id"`
			MergeError               string              `json:"merge_error"`
			MergeParams              *gitlab.MergeParams `json:"merge_params"`
			MergeWhenBuildSucceeds   bool                `json:"merge_when_build_succeeds"`
			MergeUserID              int                 `json:"merge_user_id"`
			MergeCommitSHA           string              `json:"merge_commit_sha"`
			DeletedAt                string              `json:"deleted_at"`
			ApprovalsBeforeMerge     string              `json:"approvals_before_merge"`
			RebaseCommitSHA          string              `json:"rebase_commit_sha"`
			InProgressMergeCommitSHA string              `json:"in_progress_merge_commit_sha"`
			LockVersion              int                 `json:"lock_version"`
			TimeEstimate             int                 `json:"time_estimate"`
			Source                   *gitlab.Repository  `json:"source"`
			Target                   *gitlab.Repository  `json:"target"`
			LastCommit               struct {
				ID        string     `json:"id"`
				Message   string     `json:"message"`
				Timestamp *time.Time `json:"timestamp"`
				URL       string     `json:"url"`
				Author    struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
			} `json:"last_commit"`
			WorkInProgress bool              `json:"work_in_progress"`
			URL            string            `json:"url"`
			Action         string            `json:"action"`
			OldRev         string            `json:"oldrev"`
			Assignee       *gitlab.EventUser `json:"assignee"`
		}{MilestoneID: 7},
	}}

	p := Policy{
		Conditions: Condition{
			Milestone: &Milestone{Milestone: 10},
		},
	}

	got := p.Conditions.Milestone.conditionMet(adaptor)
	if got {
		t.Errorf("expected false as milestones don't match.")
	}
}

func TestConditionsMetMilestone(t *testing.T) {
	//: 7,18
	adaptor := MergeEventAdaptor{gitlab.MergeEvent{
		ObjectAttributes: struct {
			ID                       int                 `json:"id"`
			TargetBranch             string              `json:"target_branch"`
			SourceBranch             string              `json:"source_branch"`
			SourceProjectID          int                 `json:"source_project_id"`
			AuthorID                 int                 `json:"author_id"`
			AssigneeID               int                 `json:"assignee_id"`
			AssigneeIDs              []int               `json:"assignee_ids"`
			Title                    string              `json:"title"`
			CreatedAt                string              `json:"created_at"`
			UpdatedAt                string              `json:"updated_at"`
			StCommits                []*gitlab.Commit    `json:"st_commits"`
			StDiffs                  []*gitlab.Diff      `json:"st_diffs"`
			MilestoneID              int                 `json:"milestone_id"`
			State                    string              `json:"state"`
			MergeStatus              string              `json:"merge_status"`
			TargetProjectID          int                 `json:"target_project_id"`
			IID                      int                 `json:"iid"`
			Description              string              `json:"description"`
			Position                 int                 `json:"position"`
			LockedAt                 string              `json:"locked_at"`
			UpdatedByID              int                 `json:"updated_by_id"`
			MergeError               string              `json:"merge_error"`
			MergeParams              *gitlab.MergeParams `json:"merge_params"`
			MergeWhenBuildSucceeds   bool                `json:"merge_when_build_succeeds"`
			MergeUserID              int                 `json:"merge_user_id"`
			MergeCommitSHA           string              `json:"merge_commit_sha"`
			DeletedAt                string              `json:"deleted_at"`
			ApprovalsBeforeMerge     string              `json:"approvals_before_merge"`
			RebaseCommitSHA          string              `json:"rebase_commit_sha"`
			InProgressMergeCommitSHA string              `json:"in_progress_merge_commit_sha"`
			LockVersion              int                 `json:"lock_version"`
			TimeEstimate             int                 `json:"time_estimate"`
			Source                   *gitlab.Repository  `json:"source"`
			Target                   *gitlab.Repository  `json:"target"`
			LastCommit               struct {
				ID        string     `json:"id"`
				Message   string     `json:"message"`
				Timestamp *time.Time `json:"timestamp"`
				URL       string     `json:"url"`
				Author    struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
			} `json:"last_commit"`
			WorkInProgress bool              `json:"work_in_progress"`
			URL            string            `json:"url"`
			Action         string            `json:"action"`
			OldRev         string            `json:"oldrev"`
			Assignee       *gitlab.EventUser `json:"assignee"`
		}{MilestoneID: 7},
	}}

	p := Policy{
		Conditions: Condition{
			Milestone: &Milestone{Milestone: 7},
		},
	}

	got := p.Conditions.Milestone.conditionMet(adaptor)
	if !got {
		t.Errorf("expected true as milestones match")
	}
}

func TestConditionsMetNoMilestoneInPolicy(t *testing.T) {
	//: 7,18
	adaptor := MergeEventAdaptor{gitlab.MergeEvent{
		ObjectAttributes: struct {
			ID                       int                 `json:"id"`
			TargetBranch             string              `json:"target_branch"`
			SourceBranch             string              `json:"source_branch"`
			SourceProjectID          int                 `json:"source_project_id"`
			AuthorID                 int                 `json:"author_id"`
			AssigneeID               int                 `json:"assignee_id"`
			AssigneeIDs              []int               `json:"assignee_ids"`
			Title                    string              `json:"title"`
			CreatedAt                string              `json:"created_at"`
			UpdatedAt                string              `json:"updated_at"`
			StCommits                []*gitlab.Commit    `json:"st_commits"`
			StDiffs                  []*gitlab.Diff      `json:"st_diffs"`
			MilestoneID              int                 `json:"milestone_id"`
			State                    string              `json:"state"`
			MergeStatus              string              `json:"merge_status"`
			TargetProjectID          int                 `json:"target_project_id"`
			IID                      int                 `json:"iid"`
			Description              string              `json:"description"`
			Position                 int                 `json:"position"`
			LockedAt                 string              `json:"locked_at"`
			UpdatedByID              int                 `json:"updated_by_id"`
			MergeError               string              `json:"merge_error"`
			MergeParams              *gitlab.MergeParams `json:"merge_params"`
			MergeWhenBuildSucceeds   bool                `json:"merge_when_build_succeeds"`
			MergeUserID              int                 `json:"merge_user_id"`
			MergeCommitSHA           string              `json:"merge_commit_sha"`
			DeletedAt                string              `json:"deleted_at"`
			ApprovalsBeforeMerge     string              `json:"approvals_before_merge"`
			RebaseCommitSHA          string              `json:"rebase_commit_sha"`
			InProgressMergeCommitSHA string              `json:"in_progress_merge_commit_sha"`
			LockVersion              int                 `json:"lock_version"`
			TimeEstimate             int                 `json:"time_estimate"`
			Source                   *gitlab.Repository  `json:"source"`
			Target                   *gitlab.Repository  `json:"target"`
			LastCommit               struct {
				ID        string     `json:"id"`
				Message   string     `json:"message"`
				Timestamp *time.Time `json:"timestamp"`
				URL       string     `json:"url"`
				Author    struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
			} `json:"last_commit"`
			WorkInProgress bool              `json:"work_in_progress"`
			URL            string            `json:"url"`
			Action         string            `json:"action"`
			OldRev         string            `json:"oldrev"`
			Assignee       *gitlab.EventUser `json:"assignee"`
		}{MilestoneID: 7},
	}}

	p := Policy{
		Conditions: Condition{},
	}

	got := p.Conditions.Milestone.conditionMet(adaptor)
	if !got {
		t.Errorf("expected true as policy doesn't have a milestone")
	}
}

func TestConditionsMetWebhookHasNoMilestone(t *testing.T) {
	//: 7,18
	adaptor := WikiEventAdaptor{gitlab.WikiPageEvent{
		ObjectAttributes: struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Format  string `json:"format"`
			Message string `json:"message"`
			Slug    string `json:"slug"`
			URL     string `json:"url"`
			Action  string `json:"action"`
		}{Title: "Test no State"},
	}}

	p := Policy{
		Conditions: Condition{
			Milestone: &Milestone{Milestone: 7},
		},
	}

	got := p.Conditions.Milestone.conditionMet(adaptor)
	if got {
		t.Errorf("expected false as wiki events don't have a milestone")
	}
}

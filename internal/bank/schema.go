package bank

import (
	"time"
)

// QuestionStatus constants for versioned content workflow.
const (
	StatusActive                 = "active"
	StatusApprovedActive         = "approved_active"
	StatusSeedPendingReview      = "seed_pending_review"
	StatusRetired                = "retired"
	StatusDraft                  = "draft"
	StatusCandidatePendingReview = "candidate_pending_review"
	StatusCandidateRejected      = "rejected"
)

// ApprovalAction defines the explicit action taken on a candidate proposal or question.
type ApprovalAction string

const (
	ActionApprove ApprovalAction = "approve"
	ActionReject  ApprovalAction = "reject"
	ActionRepair  ApprovalAction = "repair"
	ActionRetire  ApprovalAction = "retire"
)

// ApprovalEvent records a durable audit record of a human semantic review, approval,
// rejection, repair, or retirement action.
//
// Invariants (AGENT-CONTRACT):
// 1. Stored separately from learner evidence (attempts, mastery projections).
// 2. Requires explicit reviewer identity, timestamp, source, rule version, and action.
// 3. Learner attempts or success rates cannot produce approval events.
type ApprovalEvent struct {
	ID          string         `json:"id"`
	CandidateID string         `json:"candidate_id"`
	QuestionID  string         `json:"question_id"`
	Action      ApprovalAction `json:"action"`
	Reviewer    string         `json:"reviewer"`
	RuleVersion int            `json:"rule_version"`
	Source      string         `json:"source"`
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

func IsValidStatus(s string) bool {
	switch s {
	case StatusActive, StatusApprovedActive, StatusSeedPendingReview, StatusRetired, StatusDraft, StatusCandidatePendingReview, StatusCandidateRejected:
		return true
	default:
		return false
	}
}

// IsActiveForPractice reports whether a question's lifecycle status permits practice selection.
// Retired and draft content is strictly excluded.
func IsActiveForPractice(status string) bool {
	return status == StatusActive || status == StatusApprovedActive || status == StatusSeedPendingReview
}

// AccountsFile represents the JSON structure of curriculum/accounts.json.
type AccountsFile struct {
	SchemaVersion int           `json:"schema_version"`
	CourseScope   string        `json:"course_scope"`
	Accounts      []AccountJSON `json:"accounts"`
}

// AccountJSON represents a single account record in JSON.
type AccountJSON struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	NormalSide string `json:"normal_side"`
	ContraOf   string `json:"contra_of,omitempty"`
}

// QuestionBankFile represents the JSON structure of curriculum/seed-questions.json.
type QuestionBankFile struct {
	SchemaVersion int            `json:"schema_version"`
	Questions     []QuestionJSON `json:"questions"`
}

// QuestionJSON represents a single question template in JSON.
type QuestionJSON struct {
	ID               string             `json:"id"`
	Version          int                `json:"version"`
	FamilyID         string             `json:"family_id"`
	RuleVersion      int                `json:"rule_version"`
	Status           string             `json:"status"`
	ScenarioTemplate string             `json:"scenario_template"`
	Parameters       map[string][]int64 `json:"parameters"`
	Concepts         []string           `json:"concepts"`
	ExpectedFixture  FixtureJSON        `json:"expected_fixture"`
	Review           ReviewJSON         `json:"review"`
}

// FixtureJSON defines the expected journal postings for testing and regression.
type FixtureJSON struct {
	Postings []FixturePostingJSON `json:"postings"`
}

// FixturePostingJSON defines one posting line in the expected fixture.
type FixturePostingJSON struct {
	AccountID       string `json:"account_id"`
	Side            string `json:"side"`
	AmountParameter string `json:"amount_parameter"`
}

// ReviewJSON captures human semantic review provenance.
type ReviewJSON struct {
	Reviewer   *string `json:"reviewer"`
	ApprovedAt *string `json:"approved_at"`
	Source     string  `json:"source"`
}

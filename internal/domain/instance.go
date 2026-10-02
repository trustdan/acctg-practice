package domain

import (
	"fmt"
	"time"
)

// DrillStage identifies the sequential decision step in a progressive drill.
type DrillStage string

const (
	StageIdentifyAccount DrillStage = "identify_account" // Step 1: Which account is affected first?
	StageAccountCategory DrillStage = "account_category" // Step 2: What category is this account?
	StageDirection       DrillStage = "direction"        // Step 3: Does it increase or decrease?
	StageDebitCredit     DrillStage = "debit_credit"     // Step 4: Is this a debit or credit?
	StageCounterAccount  DrillStage = "counter_account"  // Step 5: What is the corresponding counter-account?
	StageBalancedEntry   DrillStage = "balanced_entry"   // Step 6: Full balanced journal entry assembly
	StageEquationEffect  DrillStage = "equation_effect"  // Step 7: Net effect on Assets = Liabilities + Equity
)

// ScaffoldLevel defines the fading level of pedagogical support for a question instance.
type ScaffoldLevel int

const (
	ScaffoldFull         ScaffoldLevel = 0 // Level 0: all 7 progressive stages
	ScaffoldIntermediate ScaffoldLevel = 1 // Level 1: 4 stages (Identify, Counter, Balanced, Equation)
	ScaffoldFaded        ScaffoldLevel = 2 // Level 2: 2 stages (Balanced, Equation)
)

func (s ScaffoldLevel) String() string {
	switch s {
	case ScaffoldIntermediate:
		return "Intermediate"
	case ScaffoldFaded:
		return "Faded"
	default:
		return "Full"
	}
}

// QuestionInstance represents an instantiated, reproducible question created from a bank template.
type QuestionInstance struct {
	InstanceID    string                     `json:"instance_id"`
	QuestionID    string                     `json:"question_id"`
	Version       int                        `json:"version"`
	FamilyID      string                     `json:"family_id"`
	RuleVersion   int                        `json:"rule_version"`
	PromptText    string                     `json:"prompt_text"`
	Parameters    map[string]int64           `json:"parameters"`
	RandomSeed    int64                      `json:"random_seed"`
	Entry         Entry                      `json:"entry"`
	Concepts      []string                   `json:"concepts"`
	StageAnswers  map[DrillStage]StageAnswer `json:"stage_answers"`
	ScaffoldLevel ScaffoldLevel              `json:"scaffold_level"`
}

// StageAnswer stores the derived correct answer and options snapshot for a particular drill stage.
type StageAnswer struct {
	Stage             DrillStage     `json:"stage"`
	Prompt            string         `json:"prompt"`
	CorrectOptionID   string         `json:"correct_option_id"`
	Options           []AnswerOption `json:"options"`
	CausalHint        string         `json:"causal_hint"`
	Explanation       string         `json:"explanation"`
	RelevantConceptID string         `json:"relevant_concept_id"`
}

// AnswerOption represents a single selectable choice in a multiple-choice stage prompt.
type AnswerOption struct {
	ID       string `json:"id"`
	Label    string `json:"label"`               // e.g. "a", "b", "c", "d"
	Text     string `json:"text"`                // e.g. "Cash", "Asset", "Debit"
	ErrorTag string `json:"error_tag,omitempty"` // diagnostic tag if distractor selected (e.g. "treated_advance_as_revenue")
}

// AssistanceLevel defines the level of learner assistance given on an attempt.
type AssistanceLevel string

const (
	AssistanceNone      AssistanceLevel = "none"      // First independent response
	AssistanceHinted    AssistanceLevel = "hinted"    // Learner received a Socratic hint before answering
	AssistanceRetry     AssistanceLevel = "retry"     // Learner retried after an incorrect first response
	AssistanceRevealed  AssistanceLevel = "revealed"  // Answer was shown after repeated incorrect attempts
	AssistanceReference AssistanceLevel = "reference" // Learner consulted accounting reference / cheatsheet
)

// Attempt represents a single persisted learner interaction on a drill stage.
type Attempt struct {
	AttemptID        string          `json:"attempt_id"`
	SessionID        string          `json:"session_id"`
	InstanceID       string          `json:"instance_id"`
	QuestionID       string          `json:"question_id"`
	QuestionVersion  int             `json:"question_version"`
	Stage            DrillStage      `json:"stage"`
	ConceptID        string          `json:"concept_id"`
	SelectedOptionID string          `json:"selected_option_id"`
	IsCorrect        bool            `json:"is_correct"`
	Assistance       AssistanceLevel `json:"assistance"`
	ReferenceUsed    bool            `json:"reference_used"`
	ErrorTag         string          `json:"error_tag,omitempty"`
	GradingVersion   int             `json:"grading_version"`
	AnsweredAt       time.Time       `json:"answered_at"`
}

func (a Attempt) Validate() error {
	if a.AttemptID == "" {
		return fmt.Errorf("attempt ID cannot be empty")
	}
	if a.SessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if a.InstanceID == "" {
		return fmt.Errorf("instance ID cannot be empty")
	}
	if a.QuestionID == "" {
		return fmt.Errorf("question ID cannot be empty")
	}
	if a.Stage == "" {
		return fmt.Errorf("stage cannot be empty")
	}
	if a.GradingVersion <= 0 {
		return fmt.Errorf("grading version must be positive")
	}
	if a.AnsweredAt.IsZero() {
		return fmt.Errorf("answered_at cannot be zero")
	}
	return nil
}

package exam

import (
	"errors"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

var (
	ErrHintsSuppressedInExam        = errors.New("hints are withheld during exam mode")
	ErrExplanationsSuppressedInExam = errors.New("explanations are withheld during exam mode")
	ErrReferenceSuppressedInExam    = errors.New("reference cheatsheet is withheld during exam mode")
	ErrExamAlreadyCompleted         = errors.New("exam session is already completed")
	ErrInvalidOptionID              = errors.New("invalid option ID for stage")
)

// ExamStatus represents the lifecycle state of an exam session.
type ExamStatus string

const (
	ExamStatusInProgress  ExamStatus = "in_progress"
	ExamStatusCompleted   ExamStatus = "completed"
	ExamStatusInterrupted ExamStatus = "interrupted"
	ExamStatusAbandoned   ExamStatus = "abandoned"
)

// ExamSessionRecord represents the durable storage record for an exam session.
// Persisted separately from regular practice sessions to prevent evidence contamination.
type ExamSessionRecord struct {
	ID               string     `json:"id"`
	TotalQuestions   int        `json:"total_questions"`
	TimeLimitSeconds int        `json:"time_limit_seconds"` // 0 = untimed
	ElapsedSeconds   int        `json:"elapsed_seconds"`
	Status           ExamStatus `json:"status"`
	StartedAt        time.Time  `json:"started_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Score            float64    `json:"score"` // 0.0 - 100.0
	CorrectCount     int        `json:"correct_count"`
	TotalAttempts    int        `json:"total_attempts"`
	Seed             int64      `json:"seed"`
}

// ExamAttemptRecord represents a single recorded exam interaction.
// Persisted in a dedicated table strictly separate from daily learning drill attempts.
type ExamAttemptRecord struct {
	ID               string            `json:"id"`
	ExamSessionID    string            `json:"exam_session_id"`
	InstanceID       string            `json:"instance_id"`
	QuestionID       string            `json:"question_id"`
	QuestionIndex    int               `json:"question_index"`
	Stage            domain.DrillStage `json:"stage"`
	ConceptID        string            `json:"concept_id"`
	SelectedOptionID string            `json:"selected_option_id"`
	CorrectOptionID  string            `json:"correct_option_id"`
	IsCorrect        bool              `json:"is_correct"`
	ErrorTag         string            `json:"error_tag,omitempty"`
	GradingVersion   int               `json:"grading_version"`
	AnsweredAt       time.Time         `json:"answered_at"`
}

// ExamStageResponse records the student's answer for one stage of an exam question.
type ExamStageResponse struct {
	Stage            domain.DrillStage
	ConceptID        string
	SelectedOptionID string
	CorrectOptionID  string
	IsCorrect        bool
	ErrorTag         string
	AnsweredAt       time.Time
}

// ExamQuestion tracks the runtime state of a question within an exam.
type ExamQuestion struct {
	QuestionIndex int
	Instance      *domain.QuestionInstance
	Stages        []domain.DrillStage
	CurrentStage  int
	Responses     map[domain.DrillStage]*ExamStageResponse
	IsCompleted   bool
}

// SkillSummary aggregates exam accuracy for a specific learning concept.
type SkillSummary struct {
	ConceptID   string  `json:"concept_id"`
	ConceptName string  `json:"concept_name"`
	Total       int     `json:"total"`
	Correct     int     `json:"correct"`
	Accuracy    float64 `json:"accuracy"` // percentage 0 - 100
	Status      string  `json:"status"`   // "Strong", "Needs Review", "Critical"
}

// ErrorSummary aggregates occurrences of diagnostic distractor traps.
type ErrorSummary struct {
	Tag           string `json:"tag"`
	Count         int    `json:"count"`
	Misconception string `json:"misconception"`
	Remediation   string `json:"remediation"`
}

// QuestionReviewItem holds full unmasked question details for post-exam review.
type QuestionReviewItem struct {
	QuestionIndex  int                  `json:"question_index"`
	InstanceID     string               `json:"instance_id"`
	FamilyID       string               `json:"family_id"`
	PromptText     string               `json:"prompt_text"`
	Stage          domain.DrillStage    `json:"stage"`
	StagePrompt    string               `json:"stage_prompt"`
	SelectedOption *domain.AnswerOption `json:"selected_option,omitempty"`
	CorrectOption  domain.AnswerOption  `json:"correct_option"`
	IsCorrect      bool                 `json:"is_correct"`
	ErrorTag       string               `json:"error_tag,omitempty"`
	Explanation    string               `json:"explanation"`
	Postings       []domain.Posting     `json:"postings"`
}

// ExamReport contains complete graded results, skill summaries, and error analysis.
type ExamReport struct {
	SessionID        string               `json:"session_id"`
	Status           ExamStatus           `json:"status"`
	StartedAt        time.Time            `json:"started_at"`
	CompletedAt      time.Time            `json:"completed_at"`
	Duration         time.Duration        `json:"duration"`
	TimeLimitSeconds int                  `json:"time_limit_seconds"`
	TimedOut         bool                 `json:"timed_out"`
	TotalQuestions   int                  `json:"total_questions"`
	TotalItems       int                  `json:"total_items"`
	CorrectItems     int                  `json:"correct_items"`
	Score            float64              `json:"score"` // percentage 0 - 100
	Grade            string               `json:"grade"` // A, B, C, D, F
	Skills           []SkillSummary       `json:"skills"`
	Errors           []ErrorSummary       `json:"errors"`
	Reviews          []QuestionReviewItem `json:"reviews"`
}

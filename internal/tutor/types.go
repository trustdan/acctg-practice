package tutor

import (
	"context"
	"errors"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// Standard tutor errors
var (
	ErrBudgetExceeded      = errors.New("tutor session request/token budget exceeded")
	ErrTimeout             = errors.New("tutor request timed out")
	ErrCancelled           = errors.New("tutor request cancelled")
	ErrProviderUnavailable = errors.New("tutor provider unavailable")
	ErrResponseTooLarge    = errors.New("tutor response exceeds size limit")
)

// Tutor defines the read-only contract for pedagogical accounting assistance.
//
// Invariant (AGENT-CONTRACT):
// Tutor responses are untrusted advisory prose. They NEVER directly mutate answer
// keys, attempt results, grading versions, mastery projections, or bank approval.
type Tutor interface {
	// Name returns the identifier of the tutor implementation.
	Name() string

	// Hint generates a targeted Socratic hint based on the learner's current stage and error pattern.
	Hint(ctx context.Context, req Request) (Response, error)

	// Explain provides a conceptual explanation of the underlying accounting principles.
	Explain(ctx context.Context, req Request) (Response, error)
}

// Request contains the minimum problem and pedagogical context needed for assistance.
type Request struct {
	CandidateGeneration bool
	ProblemPrompt       string
	FamilyID            string
	Stage               domain.DrillStage
	StagePrompt         string
	Options             []domain.AnswerOption
	SelectedOption      *domain.AnswerOption
	CausalHint          string
	MistakeHint         string // Reviewed option-specific hint from the immutable instance snapshot.
	Explanation         string
	ErrorTag            string
	ConceptID           string
	ReferenceRule       string
	MaxTokens           int
}

// Response contains the generated advice and provider provenance metadata.
type Response struct {
	Text           string
	Provider       string
	TokensUsed     int
	FallbackReason string
	Fallback       bool
	GeneratedAt    time.Time
}

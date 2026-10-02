package mastery

import (
	"math"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

const (
	DefaultHalfLifeDays    = 3.0
	MinHalfLifeDays        = 1.0
	MaxHalfLifeDays        = 14.0
	MinDelayedRetrievalGap = 10 * time.Minute // Minimum spacing for delayed retrieval credit
	DefaultAlpha           = 1.0
	DefaultBeta            = 1.0
	MinPracticePriority    = 0.05
	NewConceptPriority     = 1.50
)

// SessionIntensity represents learner-controlled drill pace and selection focus.
type SessionIntensity string

const (
	IntensityStandard  SessionIntensity = "standard"  // Balanced progression across new and review items
	IntensitySpaced    SessionIntensity = "spaced"    // Overdue / decayed review items prioritized
	IntensityIntensive SessionIntensity = "intensive" // Weakest concepts and new items heavily drilled
	IntensityTransfer  SessionIntensity = "transfer"  // Contrast pairs and multi-concept questions prioritized
)

// ConceptStats holds the learning evidence and mastery estimates for a single concept.
type ConceptStats struct {
	ConceptID                string               `json:"concept_id"`
	Alpha                    float64              `json:"alpha"`
	Beta                     float64              `json:"beta"`
	HalfLifeDays             float64              `json:"half_life_days"`
	IndependentAttempts      int                  `json:"independent_attempts"`
	IndependentSuccesses     int                  `json:"independent_successes"`
	AssistedAttempts         int                  `json:"assisted_attempts"`
	ReferenceUses            int                  `json:"reference_uses"`
	FirstAttemptAt           *time.Time           `json:"first_attempt_at,omitempty"`
	LastIndependentAttemptAt *time.Time           `json:"last_independent_attempt_at,omitempty"`
	LastSuccessAt            *time.Time           `json:"last_success_at,omitempty"`
	DelayedSuccesses         int                  `json:"delayed_successes"`
	HasDelayedRetrieval      bool                 `json:"has_delayed_retrieval"`
	ScaffoldLevel            domain.ScaffoldLevel `json:"scaffold_level"`
	EvidenceVersion          int                  `json:"evidence_version"`
}

func NewConceptStats(conceptID string) *ConceptStats {
	return &ConceptStats{
		ConceptID:       conceptID,
		Alpha:           DefaultAlpha,
		Beta:            DefaultBeta,
		HalfLifeDays:    DefaultHalfLifeDays,
		ScaffoldLevel:   domain.ScaffoldFull,
		EvidenceVersion: 1,
	}
}

// IsNew reports whether the concept has zero recorded attempts.
func (c *ConceptStats) IsNew() bool {
	return c.IndependentAttempts == 0 && c.AssistedAttempts == 0
}

// BaseScore computes the un-decayed probability estimate Alpha / (Alpha + Beta).
// For brand new concepts, it returns 0.0 to prevent displaying 50% prior as mastery.
func (c *ConceptStats) BaseScore() float64 {
	if c.IsNew() {
		return 0.0
	}
	total := c.Alpha + c.Beta
	if total <= 0 {
		return 0.0
	}
	return c.Alpha / total
}

// RetentionFactor calculates exponential decay 2^(-elapsed_days / half_life_days).
// Elapsed days is floored at 0 to prevent score inflation from clock rollback.
func (c *ConceptStats) RetentionFactor(now time.Time) float64 {
	if c.LastIndependentAttemptAt == nil || c.IsNew() {
		return 1.0
	}
	elapsedHours := now.Sub(*c.LastIndependentAttemptAt).Hours()
	if elapsedHours < 0 {
		elapsedHours = 0 // clock rollback protection
	}
	elapsedDays := elapsedHours / 24.0
	halfLife := c.HalfLifeDays
	if halfLife <= 0 {
		halfLife = DefaultHalfLifeDays
	}
	return math.Pow(2.0, -elapsedDays/halfLife)
}

// EffectiveScore returns the time-decayed mastery score (BaseScore * RetentionFactor).
func (c *ConceptStats) EffectiveScore(now time.Time) float64 {
	if c.IsNew() {
		return 0.0
	}
	return c.BaseScore() * c.RetentionFactor(now)
}

// PracticePriority computes the scheduling need for this concept.
// Lower effective score and longer elapsed time produce higher priority.
func (c *ConceptStats) PracticePriority(now time.Time) float64 {
	if c.IsNew() {
		return NewConceptPriority
	}
	need := 1.0 - c.EffectiveScore(now)
	if need < MinPracticePriority {
		return MinPracticePriority
	}
	return need
}

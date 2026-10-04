package mastery

import (
	"math/rand"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

// Scheduler selects next practice questions based on mastery need, decay, and anti-repeat policies.
type Scheduler struct {
	recentSettings           []string
	pendingContrast          []domain.QuestionRef
	LastSelectionWasContrast bool
	clock                    Clock
	rng                      *rand.Rand
	recentQuestions          []string // IDs of recently presented questions
	maxRecentMemory          int
	intensity                SessionIntensity
}

func NewScheduler(clock Clock, rng *rand.Rand) *Scheduler {
	if clock == nil {
		clock = RealClock{}
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	return &Scheduler{
		clock:           clock,
		rng:             rng,
		recentQuestions: make([]string, 0),
		maxRecentMemory: 3,
		intensity:       IntensityStandard,
	}
}

// SetIntensity sets the learner's chosen practice intensity mode.
func (s *Scheduler) SetIntensity(intensity SessionIntensity) {
	if intensity == "" {
		intensity = IntensityStandard
	}
	s.intensity = intensity
}

// Intensity returns the active practice intensity mode.
func (s *Scheduler) Intensity() SessionIntensity {
	if s.intensity == "" {
		return IntensityStandard
	}
	return s.intensity
}

// SelectNextQuestion picks the next question template using seeded weighted sampling.
func (s *Scheduler) SelectNextQuestion(
	availableQuestions []bank.QuestionJSON,
	projections map[string]*ConceptStats,
) *bank.QuestionJSON {
	s.LastSelectionWasContrast = false
	if len(availableQuestions) == 0 {
		return nil
	}

	// Filter out retired and inactive questions
	var eligible []bank.QuestionJSON
	for _, q := range availableQuestions {
		if bank.IsActiveForPractice(q.Status) {
			eligible = append(eligible, q)
		}
	}
	if len(eligible) == 0 {
		return nil
	}

	// At most one queued comparison; stale, changed or retired targets are skipped.
	pending := s.pendingContrast
	s.pendingContrast = nil
	for _, ref := range pending {
		for i := range eligible {
			q := &eligible[i]
			if q.ID == ref.QuestionID && q.Version == ref.QuestionVersion {
				s.LastSelectionWasContrast = true
				s.recordReviewedExposure(*q)
				return q
			}
		}
	}
	weights := make([]float64, len(eligible))
	now := s.clock.Now()
	var totalWeight float64

	for i, q := range eligible {
		weight := s.calculateQuestionWeight(q, projections, now)
		weights[i] = weight
		totalWeight += weight
	}

	if totalWeight <= 0 {
		return &eligible[0]
	}

	// Seeded weighted sampling
	target := s.rng.Float64() * totalWeight
	var running float64
	selectedIdx := 0

	for i, w := range weights {
		running += w
		if running >= target {
			selectedIdx = i
			break
		}
	}

	chosen := &eligible[selectedIdx]
	s.recordReviewedExposure(*chosen)
	return chosen
}

// RecordExposure records that a question was shown to apply anti-repeat penalties.
func (s *Scheduler) RecordExposure(questionID string) {
	s.recentQuestions = append(s.recentQuestions, questionID)
	if len(s.recentQuestions) > s.maxRecentMemory {
		s.recentQuestions = s.recentQuestions[1:]
	}
}

func (s *Scheduler) calculateQuestionWeight(
	q bank.QuestionJSON,
	projections map[string]*ConceptStats,
	now time.Time,
) float64 {
	// Base need from weakest concept
	maxNeed := MinPracticePriority
	for _, cid := range q.Concepts {
		stats, ok := projections[cid]
		var need float64
		if !ok || stats.IsNew() {
			need = NewConceptPriority
		} else {
			need = stats.PracticePriority(now)
		}
		if need > maxNeed {
			maxNeed = need
		}
	}

	weight := maxNeed

	// Intensity modulations
	switch s.intensity {
	case IntensitySpaced:
		// Prioritize concepts where retention has decayed
		for _, cid := range q.Concepts {
			if stats, ok := projections[cid]; ok && !stats.IsNew() {
				ret := stats.RetentionFactor(now)
				if ret < 0.85 {
					weight *= (2.0 - ret)
				}
			}
		}
	case IntensityIntensive:
		// Prioritize unpracticed concepts and lowest effective score
		for _, cid := range q.Concepts {
			if stats, ok := projections[cid]; !ok || stats.IsNew() {
				weight *= 2.0
			} else if stats.EffectiveScore(now) < 0.60 {
				weight *= 1.5
			}
		}
	case IntensityTransfer:
		// Prioritize contrast pairs and multi-concept questions
		if len(q.Concepts) > 1 {
			weight *= 1.5
		}
		if q.FamilyID == bank.FamilyPrepaidPurchase || q.FamilyID == bank.FamilyPrepaidConsumption ||
			q.FamilyID == bank.FamilyCustomerAdvance || q.FamilyID == bank.FamilyEarnAdvance ||
			q.FamilyID == bank.FamilyEquipmentPurchaseCash || q.FamilyID == bank.FamilyCashExpense {
			weight *= 1.8
		}
	}

	// Anti-repeat penalty
	for _, recentID := range s.recentQuestions {
		if recentID == q.ID {
			weight *= 0.20 // 80% penalty for immediate repeat
			break
		}
	}

	if policy, _, err := bank.ReviewedPedagogy(q); err == nil && policy.SettingGroup != "" {
		for _, group := range s.recentSettings {
			if group == policy.SettingGroup {
				weight *= 0.50
				break
			}
		}
	}
	if weight < MinPracticePriority {
		weight = MinPracticePriority
	}
	return weight
}

// QueueContrast stores reviewed alternatives for one subsequent comparison. It cannot chain remediation.
func (s *Scheduler) QueueContrast(inst *domain.QuestionInstance) {
	if inst == nil || inst.Pedagogy.Remediation || len(s.pendingContrast) > 0 {
		return
	}
	s.pendingContrast = append([]domain.QuestionRef(nil), inst.Pedagogy.Contrasts...)
}

func (s *Scheduler) recordReviewedExposure(q bank.QuestionJSON) {
	s.RecordExposure(q.ID)
	if policy, _, err := bank.ReviewedPedagogy(q); err == nil && policy.SettingGroup != "" {
		s.recentSettings = append(s.recentSettings, policy.SettingGroup)
		if len(s.recentSettings) > s.maxRecentMemory {
			s.recentSettings = s.recentSettings[1:]
		}
	}
}

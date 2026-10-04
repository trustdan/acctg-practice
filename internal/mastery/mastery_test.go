package mastery_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/mastery"
)

func TestNewConceptsAreLabeledNew(t *testing.T) {
	stats := mastery.NewConceptStats("cash_vs_revenue")

	if !stats.IsNew() {
		t.Fatalf("expected brand new concept to report IsNew() == true")
	}
	if stats.BaseScore() != 0.0 {
		t.Fatalf("expected BaseScore() to be 0.0 for new concept (not 50%% prior), got %f", stats.BaseScore())
	}
	now := time.Now().UTC()
	if stats.EffectiveScore(now) != 0.0 {
		t.Fatalf("expected EffectiveScore() to be 0.0 for new concept, got %f", stats.EffectiveScore(now))
	}
	if stats.PracticePriority(now) != mastery.NewConceptPriority {
		t.Fatalf("expected new concept priority to be %f, got %f", mastery.NewConceptPriority, stats.PracticePriority(now))
	}
}

func TestDelayedReviewIncreasesPriority(t *testing.T) {
	startTime := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	clock := mastery.NewMockClock(startTime)

	stats := mastery.NewConceptStats("cash_vs_revenue")
	stats.HalfLifeDays = 3.0 // 3 day half-life

	// Learner has 1 independent success at T0
	now := clock.Now()
	stats.IndependentAttempts = 1
	stats.IndependentSuccesses = 1
	stats.Alpha = 2.0 // 1 prior + 1 success
	stats.Beta = 1.0  // 1 prior
	stats.LastIndependentAttemptAt = &now

	baseScore := stats.BaseScore()
	if math.Abs(baseScore-2.0/3.0) > 1e-6 {
		t.Fatalf("expected base score ~0.667, got %f", baseScore)
	}

	eff0 := stats.EffectiveScore(clock.Now())
	prio0 := stats.PracticePriority(clock.Now())

	// Advance clock by 3 days (1 half-life)
	clock.Advance(3 * 24 * time.Hour)

	retention1 := stats.RetentionFactor(clock.Now())
	if math.Abs(retention1-0.5) > 1e-4 {
		t.Fatalf("expected retention factor to be 0.5 after 1 half-life, got %f", retention1)
	}

	eff1 := stats.EffectiveScore(clock.Now())
	prio1 := stats.PracticePriority(clock.Now())

	// Stored counts must NOT change when time advances!
	if stats.Alpha != 2.0 || stats.Beta != 1.0 || stats.IndependentAttempts != 1 {
		t.Fatalf("stored counts mutated on time advance: alpha=%f, beta=%f", stats.Alpha, stats.Beta)
	}

	// Effective score decayed: eff1 < eff0
	if eff1 >= eff0 {
		t.Errorf("expected effective score to decay after 3 days: eff0=%f, eff1=%f", eff0, eff1)
	}

	// Practice priority increased: prio1 > prio0
	if prio1 <= prio0 {
		t.Errorf("expected practice priority to INCREASE after delay: prio0=%f, prio1=%f", prio0, prio1)
	}

	// Advance clock by another 3 days (total 6 days = 2 half-lives)
	clock.Advance(3 * 24 * time.Hour)
	retention2 := stats.RetentionFactor(clock.Now())
	if math.Abs(retention2-0.25) > 1e-4 {
		t.Fatalf("expected retention factor to be 0.25 after 2 half-lives, got %f", retention2)
	}
	prio2 := stats.PracticePriority(clock.Now())
	if prio2 <= prio1 {
		t.Errorf("expected practice priority to increase further after 6 days: prio1=%f, prio2=%f", prio1, prio2)
	}

	// Clock rollback protection: Set clock to before attempt
	clock.Set(startTime.Add(-24 * time.Hour))
	rollbackRetention := stats.RetentionFactor(clock.Now())
	if rollbackRetention > 1.0 {
		t.Fatalf("retention factor inflated above 1.0 on clock rollback: %f", rollbackRetention)
	}
}

func TestHintedCorrectionsDoNotEraseFirstErrors(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	t1 := t0.Add(30 * time.Second)

	attempts := []domain.Attempt{
		// Step 1: Learner answers without aid and makes an error
		{
			AttemptID:        "att-1",
			SessionID:        "s1",
			InstanceID:       "inst-cust-adv",
			Stage:            domain.StageCounterAccount,
			ConceptID:        "cash_vs_revenue",
			SelectedOptionID: "opt_service_rev",
			IsCorrect:        false,
			Assistance:       domain.AssistanceNone, // Unassisted error
			GradingVersion:   1,
			AnsweredAt:       t0,
		},
		// Step 2: Learner retries after a Socratic hint and gets it right
		{
			AttemptID:        "att-2",
			SessionID:        "s1",
			InstanceID:       "inst-cust-adv",
			Stage:            domain.StageCounterAccount,
			ConceptID:        "cash_vs_revenue",
			SelectedOptionID: "opt_unearned_rev",
			IsCorrect:        true,
			Assistance:       domain.AssistanceRetry, // Hinted retry success
			GradingVersion:   1,
			AnsweredAt:       t1,
		},
	}

	projections := mastery.RebuildProjections(attempts)
	stats, ok := projections["cash_vs_revenue"]
	if !ok {
		t.Fatalf("missing projection for cash_vs_revenue")
	}

	// Invariant: First error must NOT be erased!
	// Beta must have incremented (+1.0)
	if stats.Beta != 2.0 {
		t.Errorf("expected Beta=2.0 (1 prior + 1 unassisted error), got %f", stats.Beta)
	}
	// Alpha must NOT have incremented (hinted retry does not inflate independent mastery)
	if stats.Alpha != 1.0 {
		t.Errorf("expected Alpha=1.0 (unassisted success only), got %f", stats.Alpha)
	}

	if stats.IndependentAttempts != 1 {
		t.Errorf("expected IndependentAttempts=1, got %d", stats.IndependentAttempts)
	}
	if stats.IndependentSuccesses != 0 {
		t.Errorf("expected IndependentSuccesses=0, got %d", stats.IndependentSuccesses)
	}
	if stats.AssistedAttempts != 1 {
		t.Errorf("expected AssistedAttempts=1, got %d", stats.AssistedAttempts)
	}

	// Base score should reflect error: 1 / (1 + 2) = 0.333
	expectedScore := 1.0 / 3.0
	if math.Abs(stats.BaseScore()-expectedScore) > 1e-4 {
		t.Errorf("expected base score %f, got %f", expectedScore, stats.BaseScore())
	}
}

func TestOneInstanceDoesNotMultiplyConceptEvidence(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)

	// Single question instance with 3 stages all testing "cash_classification"
	attempts := []domain.Attempt{
		{
			AttemptID:        "att-1",
			SessionID:        "s1",
			InstanceID:       "inst-single-instance",
			Stage:            domain.StageIdentifyAccount,
			ConceptID:        "cash_classification",
			SelectedOptionID: "opt_cash",
			IsCorrect:        true,
			Assistance:       domain.AssistanceNone,
			GradingVersion:   1,
			AnsweredAt:       t0,
		},
		{
			AttemptID:        "att-2",
			SessionID:        "s1",
			InstanceID:       "inst-single-instance",
			Stage:            domain.StageAccountCategory,
			ConceptID:        "cash_classification",
			SelectedOptionID: "opt_asset",
			IsCorrect:        true,
			Assistance:       domain.AssistanceNone,
			GradingVersion:   1,
			AnsweredAt:       t0.Add(10 * time.Second),
		},
		{
			AttemptID:        "att-3",
			SessionID:        "s1",
			InstanceID:       "inst-single-instance",
			Stage:            domain.StageDirection,
			ConceptID:        "cash_classification",
			SelectedOptionID: "opt_increase",
			IsCorrect:        true,
			Assistance:       domain.AssistanceNone,
			GradingVersion:   1,
			AnsweredAt:       t0.Add(20 * time.Second),
		},
	}

	projections := mastery.RebuildProjections(attempts)
	stats, ok := projections["cash_classification"]
	if !ok {
		t.Fatalf("missing projection for cash_classification")
	}

	// Invariant: One instance must add AT MOST 1 independent attempt and 1 success!
	if stats.IndependentAttempts != 1 {
		t.Errorf("expected 1 independent attempt (instance not multiplied), got %d", stats.IndependentAttempts)
	}
	if stats.IndependentSuccesses != 1 {
		t.Errorf("expected 1 independent success, got %d", stats.IndependentSuccesses)
	}
	if stats.Alpha != 2.0 { // 1 prior + 1 success (NOT +3!)
		t.Errorf("expected Alpha=2.0, got %f", stats.Alpha)
	}
}

func TestSchedulerSeededSelectionAndAntiRepeat(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	clock := mastery.NewMockClock(now)
	rng := rand.New(rand.NewSource(42))

	questions := []bank.QuestionJSON{
		{
			ID:       "q_cash_service",
			Status:   bank.StatusApprovedActive,
			Concepts: []string{"cash_classification", "service_revenue_classification"},
		},
		{
			ID:       "q_customer_advance",
			Status:   bank.StatusApprovedActive,
			Concepts: []string{"cash_vs_revenue", "unearned_revenue_classification"},
		},
		{
			ID:       "q_retired_example",
			Status:   bank.StatusRetired, // Retired! Must never be selected!
			Concepts: []string{"cash_vs_revenue"},
		},
	}

	// Set up projections: cash_vs_revenue is weak / new
	projections := map[string]*mastery.ConceptStats{
		"cash_classification": {
			ConceptID:                "cash_classification",
			Alpha:                    10.0,
			Beta:                     1.0,
			HalfLifeDays:             3.0,
			IndependentAttempts:      10,
			IndependentSuccesses:     10,
			LastIndependentAttemptAt: &now,
		},
		"cash_vs_revenue": mastery.NewConceptStats("cash_vs_revenue"), // Brand new
	}

	scheduler := mastery.NewScheduler(clock, rng)

	// Pick 1: Should select customer_advance because cash_vs_revenue is brand new with high priority
	selected1 := scheduler.SelectNextQuestion(questions, projections)
	if selected1 == nil {
		t.Fatalf("expected a question to be selected")
	}
	if selected1.ID == "q_retired_example" {
		t.Fatalf("retired question was selected!")
	}

	// Verify anti-repeat policy: Recording exposure reduces chance of immediate repeat
	scheduler.RecordExposure("q_customer_advance")
	scheduler.RecordExposure("q_customer_advance")

	// Next selection over multiple trials should favor q_cash_service due to anti-repeat penalty on q_customer_advance
	pickedService := false
	for i := 0; i < 20; i++ {
		pick := scheduler.SelectNextQuestion(questions, projections)
		if pick.ID == "q_cash_service" {
			pickedService = true
			break
		}
	}
	if !pickedService {
		t.Errorf("expected anti-repeat policy to allow other questions to be selected")
	}
}

func TestRetiredQuestionsStrictlyExcludedFromScheduler(t *testing.T) {
	clock := mastery.NewMockClock(time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC))
	rng := rand.New(rand.NewSource(99))
	scheduler := mastery.NewScheduler(clock, rng)

	projections := map[string]*mastery.ConceptStats{
		"cash_vs_revenue": mastery.NewConceptStats("cash_vs_revenue"),
	}

	t.Run("Return nil when bank contains only retired questions", func(t *testing.T) {
		retiredOnly := []bank.QuestionJSON{
			{ID: "ret_1", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
			{ID: "ret_2", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
			{ID: "ret_3", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
		}
		picked := scheduler.SelectNextQuestion(retiredOnly, projections)
		if picked != nil {
			t.Fatalf("expected nil when all questions are retired, got %v", picked.ID)
		}
	})

	t.Run("Never select retired questions among active questions over 100 trials", func(t *testing.T) {
		mixed := []bank.QuestionJSON{
			{ID: "act_1", Status: bank.StatusActive, Concepts: []string{"cash_vs_revenue"}},
			{ID: "ret_1", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
			{ID: "ret_2", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
			{ID: "ret_3", Status: bank.StatusRetired, Concepts: []string{"cash_vs_revenue"}},
		}
		for i := 0; i < 100; i++ {
			picked := scheduler.SelectNextQuestion(mixed, projections)
			if picked == nil {
				t.Fatalf("trial %d: expected question, got nil", i)
			}
			if picked.ID != "act_1" {
				t.Fatalf("trial %d: selected retired question %s", i, picked.ID)
			}
		}
	})
}

func TestCosmeticVariantsCannotGraduateAlone(t *testing.T) {
	// Learner solves 8 cosmetic variants within 2 minutes (< 10 min window)
	startTime := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	var attempts []domain.Attempt

	for i := 0; i < 8; i++ {
		attempts = append(attempts, domain.Attempt{
			AttemptID:        domain.NewMoney(int64(i)).FormatExact(),
			SessionID:        "sess-rapid",
			InstanceID:       domain.NewMoney(int64(i)).FormatExact(),
			QuestionID:       "q_advance",
			Stage:            domain.StageIdentifyAccount,
			ConceptID:        "unearned_revenue_classification",
			SelectedOptionID: "opt_correct",
			IsCorrect:        true,
			Assistance:       domain.AssistanceNone,
			GradingVersion:   1,
			AnsweredAt:       startTime.Add(time.Duration(i*10) * time.Second), // 10s apart
		})
	}

	projections := mastery.RebuildProjections(attempts)
	stats := projections["unearned_revenue_classification"]
	if stats == nil {
		t.Fatalf("expected stats for unearned_revenue_classification")
	}

	// Should have 8 successes, but 0 delayed successes
	if stats.IndependentSuccesses != 8 {
		t.Errorf("expected 8 independent successes, got %d", stats.IndependentSuccesses)
	}
	if stats.DelayedSuccesses != 0 {
		t.Errorf("expected 0 delayed successes from rapid cosmetic variants, got %d", stats.DelayedSuccesses)
	}
	if stats.HasDelayedRetrieval {
		t.Errorf("expected HasDelayedRetrieval to be false for rapid variants")
	}

	// Cannot graduate to ScaffoldFaded (Level 2)
	if stats.ScaffoldLevel == domain.ScaffoldFaded {
		t.Errorf("rapid cosmetic variants alone graduated concept to ScaffoldFaded!")
	}
	if stats.ScaffoldLevel != domain.ScaffoldFull {
		t.Errorf("single unreviewed setting must retain Full, got %v", stats.ScaffoldLevel)
	}

	// Half life must remain strictly DefaultHalfLifeDays (3.0)
	if stats.HalfLifeDays != mastery.DefaultHalfLifeDays {
		t.Errorf("expected fixed half life 3.0 days, got %f", stats.HalfLifeDays)
	}
}

func TestDelayedRetrievalGraduationAndBoundedHalfLife(t *testing.T) {
	startTime := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	var attempts []domain.Attempt

	// 1. Initial learning attempt at t = 0m
	attempts = append(attempts, domain.Attempt{
		AttemptID:  "att-1",
		SessionID:  "sess-1",
		InstanceID: "inst-1",
		QuestionID: "q1", SettingGroup: "context_a", PedagogyVersion: 1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       startTime,
	})

	// 2. Second attempt at t = 15m (>= 10m gap) -> First delayed success
	attempts = append(attempts, domain.Attempt{
		AttemptID:  "att-2",
		SessionID:  "sess-2",
		InstanceID: "inst-2",
		QuestionID: "q2", SettingGroup: "context_b", PedagogyVersion: 1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       startTime.Add(15 * time.Minute),
	})

	proj1 := mastery.RebuildProjections(attempts)
	s1 := proj1["cash_vs_revenue"]
	if s1.DelayedSuccesses != 1 {
		t.Fatalf("expected 1 delayed success, got %d", s1.DelayedSuccesses)
	}
	if !s1.HasDelayedRetrieval {
		t.Fatalf("expected HasDelayedRetrieval == true")
	}
	// Fixed half life with only 1 delayed success (< 2)
	if s1.HalfLifeDays != mastery.DefaultHalfLifeDays {
		t.Fatalf("expected half-life to stay fixed at 3.0 with 1 delayed success, got %f", s1.HalfLifeDays)
	}

	// 3. Third attempt at t = 45m (>= 10m gap) -> Second delayed success
	attempts = append(attempts, domain.Attempt{
		AttemptID:  "att-3",
		SessionID:  "sess-3",
		InstanceID: "inst-3",
		QuestionID: "q3", SettingGroup: "context_a", PedagogyVersion: 1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       startTime.Add(45 * time.Minute),
	})

	proj2 := mastery.RebuildProjections(attempts)
	s2 := proj2["cash_vs_revenue"]
	if s2.DelayedSuccesses != 2 {
		t.Fatalf("expected 2 delayed successes, got %d", s2.DelayedSuccesses)
	}
	// With 2 delayed successes: 3.0 * (1 + 0.5 * 1) = 4.5 days
	expectedHL := 4.5
	if math.Abs(s2.HalfLifeDays-expectedHL) > 1e-6 {
		t.Fatalf("expected bounded half-life %f, got %f", expectedHL, s2.HalfLifeDays)
	}

	// High accuracy + 3 successes + delayed retrieval -> Graduates to ScaffoldFaded
	if s2.ScaffoldLevel != domain.ScaffoldFaded {
		t.Fatalf("expected ScaffoldFaded (Level 2), got %v", s2.ScaffoldLevel)
	}
}

func TestPoorPerformanceRestoresScaffolding(t *testing.T) {
	startTime := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	var attempts []domain.Attempt

	// Build up to ScaffoldFaded
	attempts = append(attempts,
		domain.Attempt{
			AttemptID: "att-1", SessionID: "s1", InstanceID: "i1", QuestionID: "q1", SettingGroup: "context_a", PedagogyVersion: 1,
			Stage: domain.StageIdentifyAccount, ConceptID: "cash_vs_revenue",
			SelectedOptionID: "opt_1", IsCorrect: true, Assistance: domain.AssistanceNone,
			GradingVersion: 1, AnsweredAt: startTime,
		},
		domain.Attempt{
			AttemptID: "att-2", SessionID: "s2", InstanceID: "i2", QuestionID: "q2", SettingGroup: "context_b", PedagogyVersion: 1,
			Stage: domain.StageIdentifyAccount, ConceptID: "cash_vs_revenue",
			SelectedOptionID: "opt_1", IsCorrect: true, Assistance: domain.AssistanceNone,
			GradingVersion: 1, AnsweredAt: startTime.Add(15 * time.Minute),
		},
		domain.Attempt{
			AttemptID: "att-3", SessionID: "s3", InstanceID: "i3", QuestionID: "q3", SettingGroup: "context_a", PedagogyVersion: 1,
			Stage: domain.StageIdentifyAccount, ConceptID: "cash_vs_revenue",
			SelectedOptionID: "opt_1", IsCorrect: true, Assistance: domain.AssistanceNone,
			GradingVersion: 1, AnsweredAt: startTime.Add(35 * time.Minute),
		},
	)

	projBefore := mastery.RebuildProjections(attempts)
	if projBefore["cash_vs_revenue"].ScaffoldLevel != domain.ScaffoldFaded {
		t.Fatalf("expected ScaffoldFaded before error, got %v", projBefore["cash_vs_revenue"].ScaffoldLevel)
	}

	// Now learner makes an unassisted mistake on a new question
	attempts = append(attempts, domain.Attempt{
		AttemptID: "att-error", SessionID: "s4", InstanceID: "i4", QuestionID: "q4",
		Stage: domain.StageIdentifyAccount, ConceptID: "cash_vs_revenue",
		SelectedOptionID: "opt_wrong", IsCorrect: false, Assistance: domain.AssistanceNone,
		GradingVersion: 1, AnsweredAt: startTime.Add(60 * time.Minute),
	})

	projAfter := mastery.RebuildProjections(attempts)
	if projAfter["cash_vs_revenue"].ScaffoldLevel != domain.ScaffoldFull {
		t.Fatalf("expected ScaffoldFull (Level 0) after mistake, got %v", projAfter["cash_vs_revenue"].ScaffoldLevel)
	}
}

func TestReferenceUseDoesNotInflateAlpha(t *testing.T) {
	startTime := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	att := domain.Attempt{
		AttemptID:        "att-ref-1",
		SessionID:        "s-ref",
		InstanceID:       "i-ref",
		QuestionID:       "q-ref",
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt_cash",
		IsCorrect:        true,
		Assistance:       domain.AssistanceReference,
		ReferenceUsed:    true,
		GradingVersion:   1,
		AnsweredAt:       startTime,
	}

	proj := mastery.RebuildProjections([]domain.Attempt{att})
	s := proj["cash_vs_revenue"]
	if s == nil {
		t.Fatalf("expected stats")
	}

	if s.ReferenceUses != 1 {
		t.Errorf("expected 1 ReferenceUse, got %d", s.ReferenceUses)
	}
	if s.AssistedAttempts != 1 {
		t.Errorf("expected 1 AssistedAttempt, got %d", s.AssistedAttempts)
	}
	if s.IndependentSuccesses != 0 {
		t.Errorf("expected 0 IndependentSuccesses, got %d", s.IndependentSuccesses)
	}
	if s.Alpha != mastery.DefaultAlpha {
		t.Errorf("expected Alpha to remain DefaultAlpha (1.0), got %f", s.Alpha)
	}
	if s.HasDelayedRetrieval {
		t.Errorf("reference use cannot count as delayed retrieval")
	}
}

func TestComputeQuestionScaffoldLevelMinimumRule(t *testing.T) {
	q := bank.QuestionJSON{
		ID:       "q_multi",
		Concepts: []string{"concept_strong", "concept_weak"},
	}

	projections := map[string]*mastery.ConceptStats{
		"concept_strong": {
			ConceptID:           "concept_strong",
			Alpha:               5.0,
			Beta:                1.0,
			IndependentAttempts: 4,
			ScaffoldLevel:       domain.ScaffoldFaded,
		},
		"concept_weak": {
			ConceptID:           "concept_weak",
			Alpha:               1.0,
			Beta:                2.0,
			IndependentAttempts: 1,
			ScaffoldLevel:       domain.ScaffoldFull,
		},
	}

	scaffold := mastery.ComputeQuestionScaffoldLevel(q, projections, mastery.IntensityStandard)
	if scaffold != domain.ScaffoldFull {
		t.Fatalf("expected minimum scaffold level ScaffoldFull (Level 0), got %v", scaffold)
	}

	// Upgrade weak concept to Intermediate
	projections["concept_weak"].ScaffoldLevel = domain.ScaffoldIntermediate
	scaffold2 := mastery.ComputeQuestionScaffoldLevel(q, projections, mastery.IntensityStandard)
	if scaffold2 != domain.ScaffoldIntermediate {
		t.Fatalf("expected ScaffoldIntermediate (Level 1), got %v", scaffold2)
	}

	// Under IntensityIntensive, max scaffold is Intermediate even if both concepts are Faded
	projections["concept_weak"].ScaffoldLevel = domain.ScaffoldFaded
	scaffoldIntensive := mastery.ComputeQuestionScaffoldLevel(q, projections, mastery.IntensityIntensive)
	if scaffoldIntensive != domain.ScaffoldIntermediate {
		t.Fatalf("expected Intensive mode to cap scaffolding at Intermediate, got %v", scaffoldIntensive)
	}
}

func TestSchedulerIntensityModulations(t *testing.T) {
	clock := mastery.NewMockClock(time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC))
	rng := rand.New(rand.NewSource(42))
	sched := mastery.NewScheduler(clock, rng)

	if sched.Intensity() != mastery.IntensityStandard {
		t.Fatalf("expected default intensity standard, got %s", sched.Intensity())
	}

	sched.SetIntensity(mastery.IntensitySpaced)
	if sched.Intensity() != mastery.IntensitySpaced {
		t.Fatalf("expected intensity spaced, got %s", sched.Intensity())
	}

	sched.SetIntensity(mastery.IntensityTransfer)
	if sched.Intensity() != mastery.IntensityTransfer {
		t.Fatalf("expected intensity transfer, got %s", sched.Intensity())
	}
}

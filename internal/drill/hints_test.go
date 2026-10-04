package drill_test

import (
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
)

func TestMistakeHintsDependOnSubmittedOptionAndPreserveRetry(t *testing.T) {
	gen, bank, _ := setupTestDrill(t)
	q := getQuestionByID(bank, "customer_advance_basic")
	inst, err := gen.GenerateInstance(*q, 5, map[string]int64{"amount_minor_units": 5000})
	if err != nil {
		t.Fatal(err)
	}
	stage := inst.StageAnswers[domain.StageCounterAccount]
	var hints []string
	for _, option := range []string{"opt_service_rev", "opt_ar"} {
		session := drill.NewSession("mistake", inst)
		session.CurrentIndex = 4
		feedback, err := session.SubmitOption(option, time.Unix(0, 0).UTC())
		if err != nil || feedback.IsCorrect || feedback.AdvanceStage {
			t.Fatalf("wrong answer must allow one retry: %+v %v", feedback, err)
		}
		if feedback.Hint != stage.MistakeHints[option] {
			t.Fatalf("hint must come from submitted option's snapshot: %q", feedback.Hint)
		}
		repeated, err := session.RequestHint()
		if err != nil || repeated != feedback.Hint {
			t.Fatal("requesting another hint must retain the submitted misconception")
		}
		hints = append(hints, feedback.Hint)
		feedback, err = session.SubmitOption(stage.CorrectOptionID, time.Unix(1, 0).UTC())
		if err != nil || !feedback.IsCorrect || feedback.AssistanceLevel != domain.AssistanceRetry {
			t.Fatal("hinted correction must remain assisted retry evidence")
		}
		next, _ := session.CurrentStage()
		hint, _ := session.RequestHint()
		if hint != next.CausalHint {
			t.Fatal("the previous error must not carry into the next stage")
		}
	}
	if hints[0] == "" || hints[0] == hints[1] {
		t.Fatal("premature revenue and receivable confusion need different hints")
	}
}

func TestHistoricalSnapshotWithoutMistakeHintsUsesOriginalStageHint(t *testing.T) {
	gen, bank, _ := setupTestDrill(t)
	inst, err := gen.GenerateInstance(*getQuestionByID(bank, "customer_advance_basic"), 5, map[string]int64{"amount_minor_units": 5000})
	if err != nil {
		t.Fatal(err)
	}
	stage := inst.StageAnswers[domain.StageCounterAccount]
	stage.MistakeHints = nil
	stage.CausalHint = "Historical reviewed hint."
	inst.StageAnswers[stage.Stage] = stage
	session := drill.NewSession("historical", inst)
	session.CurrentIndex = 4
	feedback, err := session.SubmitOption("opt_service_rev", time.Unix(0, 0).UTC())
	if err != nil || feedback.Hint != stage.CausalHint {
		t.Fatal("historical replay must use the original hint rather than today's routing table")
	}
}

func TestDifferentQuestionVersionsHaveDistinctInstanceIDs(t *testing.T) {
	gen, bank, _ := setupTestDrill(t)
	q := *getQuestionByID(bank, "customer_advance_basic")
	params := map[string]int64{"amount_minor_units": 5000}
	original, err := gen.GenerateInstance(q, 42, params)
	if err != nil {
		t.Fatal(err)
	}
	repeated, err := gen.GenerateInstance(q, 42, params)
	if err != nil || repeated.InstanceID != original.InstanceID {
		t.Fatal("same version, seed, and parameters must remain reproducible")
	}
	q.Version++
	updated, err := gen.GenerateInstance(q, 42, params)
	if err != nil || updated.InstanceID == original.InstanceID {
		t.Fatal("new wording version must not collide with the historical instance")
	}
	if !updated.Entry.EqualNormalized(original.Entry) {
		t.Fatal("wording versions must preserve canonical accounting semantics")
	}
}

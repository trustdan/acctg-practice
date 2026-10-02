package exam

import (
	"context"
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

// Every tagged distractor must reach a specific offline hint and exam explanation,
// never the generic fallbacks.
func TestEveryDistractorTagIsRouted(t *testing.T) {
	_, _, _, gen, questions := setupTestExam(t, 1, 1, 0)
	offline := tutor.NewOfflineTutor()
	ctx := context.Background()
	fallbackHint, err := offline.Hint(ctx, tutor.Request{})
	if err != nil {
		t.Fatal(err)
	}
	_, fallbackRemediation := DistractorExplanation("unrouted_test_tag")

	checked := map[string]bool{}
	for _, q := range questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		inst, err := gen.GenerateInstance(q, 1, map[string]int64{"amount_minor_units": q.Parameters["amount_minor_units"][0]})
		if err != nil {
			t.Fatalf("%s: %v", q.ID, err)
		}
		for stage, st := range inst.StageAnswers {
			for _, opt := range st.Options {
				if opt.ErrorTag == "" || checked[opt.ErrorTag] {
					continue
				}
				checked[opt.ErrorTag] = true
				if _, rem := DistractorExplanation(opt.ErrorTag); rem == fallbackRemediation {
					t.Errorf("%s %s option %s: tag %q has no exam explanation", q.ID, stage, opt.ID, opt.ErrorTag)
				}
				hint, err := offline.Hint(ctx, tutor.Request{ErrorTag: opt.ErrorTag})
				if err != nil {
					t.Fatal(err)
				}
				if hint.Text == fallbackHint.Text {
					t.Errorf("%s %s option %s: tag %q has no offline hint", q.ID, stage, opt.ID, opt.ErrorTag)
				}
			}
		}
	}
	if len(checked) == 0 {
		t.Fatal("no tagged distractors found")
	}
}

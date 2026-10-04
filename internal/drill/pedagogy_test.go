package drill_test

import (
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"testing"
	"time"
)

func TestReviewedPedagogyContrastsAndDistractors(t *testing.T) {
	gen, seeds, cat := setupTestDrill(t)
	policy, err := bank.LoadPedagogyPolicy(curriculum.PedagogyJSON)
	if err != nil {
		t.Fatal(err)
	}
	eng := engine.NewEngine(cat)
	count := 0
	for _, meta := range policy.Questions {
		q := getQuestionByID(seeds, meta.QuestionID)
		if q == nil {
			t.Fatal("missing reviewed scenario")
		}
		for _, ref := range meta.Contrasts {
			target := getQuestionByID(seeds, ref.QuestionID)
			if target == nil {
				t.Fatal("missing partner")
			}
			for _, amount := range q.Parameters["amount_minor_units"] {
				source, err := gen.GenerateInstance(*q, 101, map[string]int64{"amount_minor_units": amount})
				if err != nil {
					t.Fatal(err)
				}
				contrast, err := gen.GenerateInstance(*target, 101, map[string]int64{"amount_minor_units": amount})
				if err != nil {
					t.Fatal(err)
				}
				if source.Entry.EqualNormalized(contrast.Entry) || !contrast.Entry.IsBalanced() {
					t.Fatal("contrast does not change a balanced event entry")
				}
				eval := eng.EvaluateEntry(engine.TransactionEvent{FamilyID: q.FamilyID, Parameters: source.Parameters}, contrast.Entry)
				if eval.IsCorrect {
					t.Fatal("balanced contrast accepted as answer")
				}
			}
		}
		if meta.DistractorProfile == "" {
			continue
		}
		count++
		for _, amount := range q.Parameters["amount_minor_units"] {
			for _, level := range []domain.ScaffoldLevel{domain.ScaffoldFull, domain.ScaffoldIntermediate, domain.ScaffoldFaded} {
				inst, err := gen.GenerateInstanceWithScaffold(*q, 101, map[string]int64{"amount_minor_units": amount}, level)
				if err != nil {
					t.Fatal(err)
				}
				st := inst.StageAnswers[domain.StageBalancedEntry]
				wrong := ""
				switch meta.DistractorProfile {
				case "service_vs_capital":
					wrong = "opt_dr_cash_cr_stock"
				case "held_cash_vs_no_entry", "seller_receipt_vs_no_entry":
					wrong = "opt_no_entry"
				case "collection_vs_advance":
					wrong = "opt_dr_cash_cr_unearned"
				case "earned_advance_vs_receivable":
					wrong = "opt_dr_ar_cr_rev"
				}
				if _, ok := bank.ReviewedAmountDistractor(meta.DistractorProfile); ok {
					wrong = "opt_wrong_amount_scope"
				}
				found := false
				seen := map[string]bool{}
				for _, opt := range st.Options {
					if seen[opt.ID] {
						t.Fatal("duplicate option identity")
					}
					seen[opt.ID] = true
					if opt.ID == wrong {
						found = true
						if opt.ErrorTag == "" || st.HintForOption(opt.ID) == "" {
							t.Fatal("new misconception missing reviewed guidance")
						}
					}
				}
				if !found || len(st.Options) > 4 {
					t.Fatalf("%s missing distractor or too many keyboard choices", q.ID)
				}
				session := drill.NewSession("reviewed", inst)
				for {
					cur, _ := session.CurrentStage()
					if cur.Stage == domain.StageBalancedEntry {
						break
					}
					session.SubmitOption(cur.CorrectOptionID, time.Unix(1, 0))
				}
				fb, err := session.SubmitOption(wrong, time.Unix(2, 0))
				if err != nil || fb.IsCorrect || fb.AdvanceStage || fb.Hint != st.HintForOption(wrong) {
					t.Fatal("new distractor not rejected with retry hint")
				}
				fb, err = session.SubmitOption(st.CorrectOptionID, time.Unix(3, 0))
				if err != nil || !fb.IsCorrect || fb.AssistanceLevel != domain.AssistanceRetry {
					t.Fatal("assisted correction changed")
				}
			}
		}
	}
	if count != 16 {
		t.Fatalf("expected sixteen reviewed distractor scenarios, got %d", count)
	}
}
func TestContrastSessionPersistsAssistanceAcrossEveryStage(t *testing.T) {
	gen, seeds, _ := setupTestDrill(t)
	q := getQuestionByID(seeds, "earn_advance_same_day")
	inst, err := gen.GenerateInstance(*q, 101, map[string]int64{"amount_minor_units": 10000})
	if err != nil {
		t.Fatal(err)
	}
	inst.Pedagogy.Remediation = true
	session := drill.NewSession("contrast", inst)
	for !session.IsCompleted {
		st, _ := session.CurrentStage()
		fb, err := session.SubmitOption(st.CorrectOptionID, time.Unix(1, 0))
		if err != nil || fb.AssistanceLevel != domain.AssistanceContrast {
			t.Fatal("guided stage became independent")
		}
	}
	for _, att := range session.Attempts {
		if att.Assistance != domain.AssistanceContrast || att.SettingGroup != inst.Pedagogy.SettingGroup || att.PedagogyVersion != 2 {
			t.Fatal("transfer snapshot lost")
		}
	}
}

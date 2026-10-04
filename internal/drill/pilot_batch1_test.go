package drill_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestPilotBatch1SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"cash_service_training":      {"cash", "service_revenue"},
		"customer_advance_training":  {"cash", "unearned_revenue"},
		"earn_advance_training":      {"unearned_revenue", "service_revenue"},
		"service_on_credit_repairs":  {"accounts_receivable", "service_revenue"},
		"collect_receivable_repairs": {"cash", "accounts_receivable"},
		"borrow_cash_future_payroll": {"cash", "notes_payable"},
	}
	pairs := [][2]string{{"cash_service_training", "customer_advance_training"}, {"customer_advance_training", "earn_advance_training"}, {"service_on_credit_repairs", "collect_receivable_repairs"}, {"cash_service_training", "borrow_cash_future_payroll"}}
	verifyPilotBatch(t, "../../docs/PILOT-BATCH1.json", expected, pairs)
}

func verifyPilotBatch(t *testing.T, packagePath string, expected map[string][2]domain.AccountID, pairs [][2]string) {
	t.Helper()
	gen, seeds, catalog := setupTestDrill(t)
	pilot, err := bank.LoadQuestionBankFile(packagePath, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if len(pilot.Questions) != len(expected) {
		t.Fatalf("review batch must contain %d complete scenarios, got %d", len(expected), len(pilot.Questions))
	}
	eng := engine.NewEngine(catalog)
	entries := map[string]domain.Entry{}
	families := map[string]string{}
	for _, q := range pilot.Questions {
		accounts, ok := expected[q.ID]
		if !ok {
			t.Fatalf("unexpected pilot ID %s", q.ID)
		}
		if review := bank.CheckSemanticWording(q.FamilyID, q.ScenarioTemplate); !review.Approved {
			t.Errorf("%s semantic wording conflict: %v", q.ID, review.Violations)
		}
		for _, old := range seeds.Questions {
			if old.ID != q.ID && strings.EqualFold(old.ScenarioTemplate, q.ScenarioTemplate) {
				t.Errorf("%s duplicates %s", q.ID, old.ID)
			}
		}
		for _, account := range []string{"Accounts Receivable", "Unearned Revenue", "Service Revenue", "Notes Payable", "Common Stock"} {
			if strings.Contains(q.ScenarioTemplate, account) {
				t.Errorf("%s scenario names an account answer", q.ID)
			}
		}
		// Preview exact reviewed overrides before activation. This copy has no
		// path into production selection or publication.
		preview := q
		preview.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			params := map[string]int64{"amount_minor_units": amount}
			canonical := domain.NewEntry(domain.Posting{AccountID: accounts[0], Side: domain.SideDebit, Amount: domain.NewMoney(amount)}, domain.Posting{AccountID: accounts[1], Side: domain.SideCredit, Amount: domain.NewMoney(amount)})
			for _, level := range []domain.ScaffoldLevel{domain.ScaffoldFull, domain.ScaffoldIntermediate, domain.ScaffoldFaded} {
				inst, err := gen.GenerateInstanceWithScaffold(preview, 101, params, level)
				if err != nil {
					t.Fatal(err)
				}
				if !inst.Entry.EqualNormalized(canonical) {
					t.Errorf("%s wrong canonical entry: %+v", q.ID, inst.Entry)
				}
				if err := inst.Entry.Validate(catalog); err != nil {
					t.Fatal(err)
				}
				if delta, err := inst.Entry.EquationEffects(catalog); err != nil || delta.DeltaAssets != delta.DeltaLiabilities.Add(delta.DeltaEquity) {
					t.Fatalf("%s equation does not reconcile: %+v %v", q.ID, delta, err)
				}
				for _, key := range drill.ScaffoldStageSequence(level) {
					stage := inst.StageAnswers[key]
					if !strings.Contains(stage.CausalHint, "?") || stage.Explanation == "" {
						t.Errorf("%s %s incomplete teaching", q.ID, key)
					}
					if strings.Contains(stage.Explanation+stage.CausalHint, "${") {
						t.Errorf("%s %s unrendered teaching", q.ID, key)
					}
					if text, ok := q.Teaching[key]; ok {
						want := strings.ReplaceAll(text.Hint, "${amount_dollars}", domain.NewMoney(amount).FormatDollars())
						if stage.CausalHint != want {
							t.Errorf("%s %s reviewed scenario hint not applied", q.ID, key)
						}
					}
				}
				session := drill.NewSession("pilot", inst)
				for !session.IsCompleted {
					stage, err := session.CurrentStage()
					if err != nil {
						t.Fatal(err)
					}
					wrong := ""
					for _, option := range stage.Options {
						if option.ID != stage.CorrectOptionID {
							wrong = option.ID
							break
						}
					}
					feedback, err := session.SubmitOption(wrong, time.Unix(int64(len(session.Attempts)), 0).UTC())
					if err != nil || feedback.IsCorrect || feedback.AdvanceStage || feedback.Hint != stage.HintForOption(wrong) {
						t.Fatalf("%s %s wrong-option feedback: %+v %v", q.ID, stage.Stage, feedback, err)
					}
					feedback, err = session.SubmitOption(stage.CorrectOptionID, time.Unix(int64(len(session.Attempts)), 0).UTC())
					if err != nil || !feedback.IsCorrect || feedback.AssistanceLevel != domain.AssistanceRetry || feedback.Explanation != stage.Explanation {
						t.Fatalf("%s %s retry feedback: %+v %v", q.ID, stage.Stage, feedback, err)
					}
				}
				if len(session.Attempts) != 2*len(drill.ScaffoldStageSequence(level)) || !session.Recap().IsBalanced {
					t.Fatalf("%s wrong attempt count or unbalanced recap", q.ID)
				}
			}
		}
		inst, err := gen.GenerateInstance(preview, 101, map[string]int64{"amount_minor_units": 5000})
		if err != nil {
			t.Fatal(err)
		}
		entries[q.ID], families[q.ID] = inst.Entry, q.FamilyID
	}
	for _, pair := range pairs {
		for _, id := range pair {
			if _, found := entries[id]; found {
				continue
			}
			q := getQuestionByID(seeds, id)
			if q == nil {
				t.Fatalf("missing reviewed contrast partner %s", id)
			}
			proc, err := eng.ProcessEvent(engine.TransactionEvent{FamilyID: q.FamilyID, Parameters: map[string]int64{"amount_minor_units": 5000}})
			if err != nil {
				t.Fatal(err)
			}
			entries[id], families[id] = proc.Entry, q.FamilyID
		}
		for i := range pair {
			id, contrast := pair[i], pair[1-i]
			wrong := entries[contrast]
			if !wrong.IsBalanced() {
				t.Fatal("contrast must be balanced to test semantics beyond balancing")
			}
			result := eng.EvaluateEntry(engine.TransactionEvent{FamilyID: families[id], Parameters: map[string]int64{"amount_minor_units": 5000}}, wrong)
			if result.IsCorrect {
				t.Errorf("%s incorrectly accepts contrast %s", id, contrast)
			}
		}
	}
}

func TestPublishedPilotMatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/PILOT-BATCH1.json")
}

func verifyPublishedPilot(t *testing.T, packagePath string) {
	t.Helper()
	_, seeds, catalog := setupTestDrill(t)
	pilot, err := bank.LoadQuestionBankFile(packagePath, catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, reviewed := range pilot.Questions {
		published := getQuestionByID(seeds, reviewed.ID)
		if published == nil || !reflect.DeepEqual(*published, reviewed) {
			t.Errorf("%s active seed does not match the reviewed package", reviewed.ID)
		}
		if reviewed.Status != bank.StatusApprovedActive {
			t.Errorf("%s reviewed pilot must be approved_active", reviewed.ID)
		}
	}
}

package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
	"reflect"
	"testing"
)

func TestQuality90OptionsRemainDeterministicAndDiagnostic(t *testing.T) {
	gen, seeds, cat := setupTestDrill(t)
	eng := engine.NewEngine(cat)
	for _, q := range seeds.Questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		for _, amount := range q.Parameters["amount_minor_units"] {
			for _, seed := range []int64{1, 42, 101} {
				inst, err := gen.GenerateInstance(q, seed, map[string]int64{"amount_minor_units": amount})
				if err != nil {
					t.Fatal(err)
				}
				again, err := gen.GenerateInstance(q, seed, map[string]int64{"amount_minor_units": amount})
				if err != nil || !reflect.DeepEqual(inst, again) {
					t.Fatalf("nonreproducible %s", q.ID)
				}
				for key, stage := range inst.StageAnswers {
					ids, texts, labels := map[string]bool{}, map[string]bool{}, map[string]bool{}
					correct := 0
					for _, opt := range stage.Options {
						if ids[opt.ID] || texts[opt.Text] || labels[opt.Label] {
							t.Fatalf("duplicate option %s/%s", q.ID, key)
						}
						ids[opt.ID] = true
						texts[opt.Text] = true
						labels[opt.Label] = true
						if opt.ID == stage.CorrectOptionID {
							correct++
							if opt.ErrorTag != "" {
								t.Fatal("correct tagged as error")
							}
						}
						if opt.ID == "opt_ap" && (q.FamilyID == bank.FamilyCashService || q.FamilyID == bank.FamilyCustomerAdvance || q.FamilyID == bank.FamilyServiceOnCredit || q.FamilyID == bank.FamilyCollectReceivable) {
							t.Fatal("generic payable filler retained")
						}
						if opt.ID == "opt_wrong_amount_scope" && (opt.ErrorTag != engine.TagWrongAmount || stage.MistakeHints[opt.ID] == "") {
							t.Fatal("amount hint/tag absent")
						}
					}
					if correct != 1 || len(stage.Options) > 4 {
						t.Fatalf("invalid choices %s/%s", q.ID, key)
					}
				}
				// A balanced numeric entry with the correct accounts can still have the wrong amount.
				wrong := domain.Entry{Postings: append([]domain.Posting(nil), inst.Entry.Postings...)}
				for i := range wrong.Postings {
					wrong.Postings[i].Amount = domain.NewMoney(amount * 2)
				}
				if !wrong.IsBalanced() || eng.EvaluateEntry(engine.TransactionEvent{FamilyID: q.FamilyID, Parameters: inst.Parameters}, wrong).IsCorrect {
					t.Fatal("balanced wrong amount accepted")
				}
			}
		}
	}
}

package drill_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// CandidateBatch represents the batch JSON structure for candidate questions.
type candidateBatch struct {
	SchemaVersion int                 `json:"schema_version"`
	Questions     []bank.QuestionJSON `json:"questions"`
}

func TestStreamM03Questions(t *testing.T) {
	// 1. Load accounts catalog
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}

	// 2. Load candidate questions from batch_m03.json
	batchPath := "../../curriculum/candidates/batch_m03.json"
	data, err := os.ReadFile(batchPath)
	if err != nil {
		t.Fatalf("failed to read batch_m03.json: %v", err)
	}

	var batch candidateBatch
	if err := json.Unmarshal(data, &batch); err != nil {
		t.Fatalf("failed to parse batch_m03.json: %v", err)
	}

	if len(batch.Questions) != 4 {
		t.Fatalf("expected 4 questions in batch_m03.json, got %d", len(batch.Questions))
	}

	// 3. Load candidate pedagogy from pedagogy_m03.json
	pedagogyPath := "../../curriculum/candidates/pedagogy_m03.json"
	pedagogyData, err := os.ReadFile(pedagogyPath)
	if err != nil {
		t.Fatalf("failed to read pedagogy_m03.json: %v", err)
	}
	pedagogyPolicy, err := bank.LoadPedagogyPolicy(pedagogyData)
	if err != nil {
		t.Fatalf("failed to validate pedagogy_m03.json: %v", err)
	}
	if len(pedagogyPolicy.Questions) != 4 {
		t.Fatalf("expected 4 pedagogy policy questions, got %d", len(pedagogyPolicy.Questions))
	}

	// 4. Setup Engine with Stream-M03 rules registered
	eng := engine.NewEngine(catalog)
	engine.RegisterM03Rules(eng)

	expectedIDs := map[string]bool{
		"q_m03_ocf_ap_increase_01":      true,
		"q_m03_ocf_prepaid_increase_02": true,
		"q_m03_ocf_unearned_revenue_03": true,
		"q_m03_ocf_quick_ratio_04":      true,
	}

	requiredStages := []domain.DrillStage{
		domain.StageIdentifyAccount,
		domain.StageAccountCategory,
		domain.StageDirection,
		domain.StageDebitCredit,
		domain.StageCounterAccount,
		domain.StageBalancedEntry,
		domain.StageEquationEffect,
	}

	for _, q := range batch.Questions {
		t.Run(q.ID, func(t *testing.T) {
			if !expectedIDs[q.ID] {
				t.Errorf("unexpected question ID: %s", q.ID)
			}

			// Invariant: Status must be candidate_pending_review
			if q.Status != bank.StatusCandidatePendingReview {
				t.Errorf("expected status %s, got %s", bank.StatusCandidatePendingReview, q.Status)
			}

			// Invariant: Triple parameter permutations (at least 3 positive amounts)
			amounts, ok := q.Parameters["amount_minor_units"]
			if !ok || len(amounts) < 3 {
				t.Fatalf("question %s must have at least 3 amount_minor_units, got %v", q.ID, amounts)
			}
			for _, a := range amounts {
				if a <= 0 {
					t.Errorf("question %s amount must be positive, got %d", q.ID, a)
				}
			}

			// Invariant: Expected fixture postings balance to the penny
			for sampleIdx, amtVal := range amounts {
				var drTotal, crTotal int64
				for _, p := range q.ExpectedFixture.Postings {
					if !catalog.Has(domain.AccountID(p.AccountID)) {
						t.Errorf("question %s fixture references unknown account %s", q.ID, p.AccountID)
					}
					if p.Side == string(domain.SideDebit) {
						drTotal += amtVal
					} else if p.Side == string(domain.SideCredit) {
						crTotal += amtVal
					}
				}
				if drTotal <= 0 || drTotal != crTotal {
					t.Errorf("sample %d for question %s unbalanced: debits=%d, credits=%d", sampleIdx, q.ID, drTotal, crTotal)
				}
			}

			// Process canonical event in engine
			testAmt := amounts[0]
			event := engine.TransactionEvent{
				FamilyID:   q.FamilyID,
				Parameters: map[string]int64{"amount_minor_units": testAmt},
			}
			proc, err := eng.ProcessEvent(event)
			if err != nil {
				t.Fatalf("ProcessEvent failed for family %s: %v", q.FamilyID, err)
			}

			if !proc.Entry.IsBalanced() {
				t.Fatalf("derived entry not balanced for %s: %+v", q.ID, proc.Entry)
			}

			// Invariant: Balance sheet equation balance holds
			if proc.Equation.DeltaAssets != proc.Equation.DeltaLiabilities+proc.Equation.DeltaEquity {
				t.Errorf("balance sheet equation not balanced for %s: dA=%d, dL=%d, dE=%d",
					q.ID, proc.Equation.DeltaAssets, proc.Equation.DeltaLiabilities, proc.Equation.DeltaEquity)
			}

			// Build scaffolding stages
			seed := int64(1337)
			money := domain.NewMoney(testAmt)
			promptText := strings.ReplaceAll(q.ScenarioTemplate, "${amount_dollars}", money.FormatDollars())

			inst := &domain.QuestionInstance{
				InstanceID:    q.ID + "-test",
				QuestionID:    q.ID,
				Version:       q.Version,
				FamilyID:      q.FamilyID,
				PromptText:    promptText,
				Parameters:    map[string]int64{"amount_minor_units": testAmt},
				RandomSeed:    seed,
				Entry:         proc.Entry,
				Concepts:      q.Concepts,
				StageAnswers:  make(map[domain.DrillStage]domain.StageAnswer),
				ScaffoldLevel: domain.ScaffoldFull,
			}

			handled := drill.BuildM03Stages(q.FamilyID, inst, money, seed, proc)
			if !handled {
				t.Fatalf("BuildM03Stages did not handle family %s", q.FamilyID)
			}

			// Validate all 7 stages
			if len(inst.StageAnswers) != 7 {
				t.Fatalf("expected 7 stage answers, got %d", len(inst.StageAnswers))
			}

			for _, st := range requiredStages {
				ans, ok := inst.StageAnswers[st]
				if !ok {
					t.Errorf("missing stage %s", st)
					continue
				}

				if ans.CorrectOptionID == "" {
					t.Errorf("stage %s has empty CorrectOptionID", st)
				}
				if len(ans.Options) < 2 {
					t.Errorf("stage %s has fewer than 2 options: %d", st, len(ans.Options))
				}

				// Find correct option
				var foundCorrect bool
				for _, opt := range ans.Options {
					if opt.ID == ans.CorrectOptionID {
						foundCorrect = true
						break
					}
				}
				if !foundCorrect {
					t.Errorf("stage %s correct option %s not found in options", st, ans.CorrectOptionID)
				}

				// Invariant: Socratic hint must end with a question mark
				if !strings.HasSuffix(strings.TrimSpace(ans.CausalHint), "?") {
					t.Errorf("stage %s CausalHint must end with '?', got: %q", st, ans.CausalHint)
				}

				// Invariant: Stage 1 explanation must NOT include 'debit' or 'credit'
				if st == domain.StageIdentifyAccount {
					lowerExp := strings.ToLower(ans.Explanation)
					if strings.Contains(lowerExp, "debit") || strings.Contains(lowerExp, "credit") {
						t.Errorf("Stage 1 explanation contains 'debit' or 'credit': %q", ans.Explanation)
					}
				}
			}
		})
	}
}

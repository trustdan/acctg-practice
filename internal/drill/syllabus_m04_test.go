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

// candidateBankFile mirrors the bank.QuestionBankFile JSON schema for isolated candidate testing.
type candidateBankFile struct {
	SchemaVersion int                 `json:"schema_version"`
	Questions     []bank.QuestionJSON `json:"questions"`
}

func TestStreamM04Questions(t *testing.T) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load canonical accounts: %v", err)
	}

	batchPath := "../../curriculum/candidates/batch_m04.json"
	data, err := os.ReadFile(batchPath)
	if err != nil {
		t.Fatalf("failed to read stream batch file %s: %v", batchPath, err)
	}

	var batch candidateBankFile
	if err := json.Unmarshal(data, &batch); err != nil {
		t.Fatalf("failed to decode stream candidate batch: %v", err)
	}

	if batch.SchemaVersion != 1 {
		t.Fatalf("expected schema_version 1, got %d", batch.SchemaVersion)
	}
	if len(batch.Questions) != 4 {
		t.Fatalf("expected exactly 4 candidate questions for Stream-M04, got %d", len(batch.Questions))
	}

	expectedIDs := []string{
		"q_m04_rev_variable_consideration_01",
		"q_m04_rev_contract_asset_ar_01",
		"q_m04_rev_right_of_return_01",
		"q_m04_rev_satisfaction_over_time_01",
	}

	for i, expectedID := range expectedIDs {
		if batch.Questions[i].ID != expectedID {
			t.Errorf("question[%d] ID mismatch: expected %q, got %q", i, expectedID, batch.Questions[i].ID)
		}
	}

	for _, q := range batch.Questions {
		t.Run(q.ID, func(t *testing.T) {
			if !strings.HasPrefix(q.ID, "q_m04_rev_") {
				t.Errorf("question ID %q does not match required namespace prefix q_m04_rev_", q.ID)
			}
			if q.Version < 1 {
				t.Errorf("question %s has invalid version: %d", q.ID, q.Version)
			}
			if q.RuleVersion < 1 {
				t.Errorf("question %s has invalid rule_version: %d", q.ID, q.RuleVersion)
			}
			if q.Status != bank.StatusCandidatePendingReview {
				t.Errorf("expected candidate_pending_review status, got %q", q.Status)
			}
			if q.Module != 4 {
				t.Errorf("expected module 4 for Stream-M04, got %d", q.Module)
			}
			if !strings.Contains(q.ScenarioTemplate, "${amount_dollars}") {
				t.Errorf("scenario template must contain ${amount_dollars}")
			}

			// Invariant: At least 3 parameter permutations with positive integer minor units
			amounts, ok := q.Parameters["amount_minor_units"]
			if !ok || len(amounts) < 3 {
				t.Fatalf("question %s must provide at least 3 amount_minor_units permutations, got %v", q.ID, amounts)
			}
			for pIdx, amt := range amounts {
				if amt <= 0 {
					t.Fatalf("question %s amount[%d] must be positive integer cents, got %d", q.ID, pIdx, amt)
				}
			}

			// Invariant: Double-Entry Balancing (Debits == Credits) for all parameter permutations
			if len(q.ExpectedFixture.Postings) < 2 {
				t.Fatalf("question %s fixture must have >= 2 postings, got %d", q.ID, len(q.ExpectedFixture.Postings))
			}
			for _, p := range q.ExpectedFixture.Postings {
				if !catalog.Has(domain.AccountID(p.AccountID)) {
					t.Fatalf("question %s references unknown account: %s", q.ID, p.AccountID)
				}
			}

			for idx, testAmt := range amounts {
				var drTotal, crTotal int64
				var deltaAssets, deltaLiab, deltaEquity int64
				for _, p := range q.ExpectedFixture.Postings {
					acc, _ := catalog.Get(domain.AccountID(p.AccountID))
					if p.Side == string(domain.SideDebit) {
						drTotal += testAmt
						switch acc.Category {
						case domain.CategoryAsset:
							deltaAssets += testAmt
						case domain.CategoryLiability:
							deltaLiab -= testAmt
						case domain.CategoryExpense:
							deltaEquity -= testAmt
						}
					} else if p.Side == string(domain.SideCredit) {
						crTotal += testAmt
						switch acc.Category {
						case domain.CategoryLiability:
							deltaLiab += testAmt
						case domain.CategoryRevenue:
							deltaEquity += testAmt
						case domain.CategoryAsset:
							deltaAssets -= testAmt
						}
					}
				}

				if drTotal != crTotal {
					t.Errorf("sample %d unbalanced entry: debits=%d, credits=%d", idx, drTotal, crTotal)
				}
				// Invariant: Derived Equation Balancing: Delta Assets == Delta Liabilities + Delta Equity
				if deltaAssets != deltaLiab+deltaEquity {
					t.Errorf("sample %d equation unbalanced: deltaA=%d, deltaL=%d, deltaE=%d", idx, deltaAssets, deltaLiab, deltaEquity)
				}
			}

			// Invariant: Review provenance recorded
			if q.Review.Reviewer == nil || *q.Review.Reviewer == "" {
				t.Errorf("question %s missing review.reviewer", q.ID)
			}
			if q.Review.ApprovedAt == nil || *q.Review.ApprovedAt == "" {
				t.Errorf("question %s missing review.approved_at", q.ID)
			}
			if q.Review.Source == "" {
				t.Errorf("question %s missing review.source", q.ID)
			}
		})
	}
}

func TestStreamM04EngineRules(t *testing.T) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}

	testCases := []struct {
		name         string
		rule         engine.EventRule
		params       map[string]int64
		wantDebits   map[domain.AccountID]int64
		wantCredits  map[domain.AccountID]int64
		badEntry     domain.Entry
		wantErrorTag string
	}{
		{
			name: "M04 Variable Consideration Rule",
			rule: engine.M04VariableConsiderationRule{},
			params: map[string]int64{
				"amount_minor_units": 240000,
			},
			wantDebits: map[domain.AccountID]int64{
				"accounts_receivable": 240000,
			},
			wantCredits: map[domain.AccountID]int64{
				"service_revenue": 240000,
			},
			badEntry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(240000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(240000)},
			),
			wantErrorTag: engine.TagCashRecordedWhenUncollected,
		},
		{
			name: "M04 Contract Asset vs AR Rule",
			rule: engine.M04ContractAssetARRule{},
			params: map[string]int64{
				"amount_minor_units": 300000,
			},
			wantDebits: map[domain.AccountID]int64{
				"accounts_receivable": 300000,
			},
			wantCredits: map[domain.AccountID]int64{
				"service_revenue": 300000,
			},
			badEntry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(300000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(300000)},
			),
			wantErrorTag: engine.TagCashRecordedWhenUncollected,
		},
		{
			name: "M04 Right of Return Advance Rule",
			rule: engine.M04RightOfReturnRule{},
			params: map[string]int64{
				"amount_minor_units": 360000,
			},
			wantDebits: map[domain.AccountID]int64{
				"cash": 360000,
			},
			wantCredits: map[domain.AccountID]int64{
				"unearned_revenue": 360000,
			},
			badEntry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(360000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(360000)},
			),
			wantErrorTag: engine.TagM04RevenueRecognizedWithUncertainReturn,
		},
		{
			name: "M04 Satisfaction Over Time Rule",
			rule: engine.M04SatisfactionOverTimeRule{},
			params: map[string]int64{
				"amount_minor_units": 400000,
			},
			wantDebits: map[domain.AccountID]int64{
				"unearned_revenue": 400000,
			},
			wantCredits: map[domain.AccountID]int64{
				"service_revenue": 400000,
			},
			badEntry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(400000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(400000)},
			),
			wantErrorTag: engine.TagM04CashDebitedOnOverTimeSatisfaction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := tc.rule.DeriveEntry(tc.params, catalog)
			if err != nil {
				t.Fatalf("DeriveEntry failed: %v", err)
			}
			if !entry.IsBalanced() {
				t.Fatalf("Entry is not balanced: %+v", entry)
			}

			// Validate postings
			for _, p := range entry.Postings {
				if p.Side == domain.SideDebit {
					expectedAmt, exists := tc.wantDebits[p.AccountID]
					if !exists || expectedAmt != p.Amount.Cents() {
						t.Errorf("unexpected debit: %s %v, want %d", p.AccountID, p.Amount, expectedAmt)
					}
				} else if p.Side == domain.SideCredit {
					expectedAmt, exists := tc.wantCredits[p.AccountID]
					if !exists || expectedAmt != p.Amount.Cents() {
						t.Errorf("unexpected credit: %s %v, want %d", p.AccountID, p.Amount, expectedAmt)
					}
				}
			}

			// Validate explanation
			expl := tc.rule.Explain(tc.params)
			if expl.Summary == "" || expl.DebitRationale == "" || expl.CreditRationale == "" || expl.EquationSummary == "" {
				t.Errorf("incomplete explanation: %+v", expl)
			}

			// Validate diagnosis of distractor
			tags, _ := tc.rule.DiagnoseError(tc.badEntry, tc.params)
			foundTag := false
			for _, tag := range tags {
				if tag == tc.wantErrorTag {
					foundTag = true
					break
				}
			}
			if !foundTag {
				t.Errorf("expected error tag %q in %v", tc.wantErrorTag, tags)
			}
		})
	}
}

func TestStreamM04StageBuilders(t *testing.T) {
	builders := []struct {
		name       string
		builderFn  func(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent)
		testAmount domain.Money
	}{
		{
			name:       "Variable Consideration Stages",
			builderFn:  drill.ExportBuildM04VariableConsiderationStages,
			testAmount: domain.NewMoney(240000),
		},
		{
			name:       "Contract Asset vs AR Stages",
			builderFn:  drill.ExportBuildM04ContractAssetARStages,
			testAmount: domain.NewMoney(300000),
		},
		{
			name:       "Right of Return Stages",
			builderFn:  drill.ExportBuildM04RightOfReturnStages,
			testAmount: domain.NewMoney(360000),
		},
		{
			name:       "Satisfaction Over Time Stages",
			builderFn:  drill.ExportBuildM04SatisfactionOverTimeStages,
			testAmount: domain.NewMoney(400000),
		},
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

	for _, tc := range builders {
		t.Run(tc.name, func(t *testing.T) {
			inst := &domain.QuestionInstance{
				StageAnswers: make(map[domain.DrillStage]domain.StageAnswer),
			}
			proc := engine.ProcessedEvent{}
			seed := int64(42)

			tc.builderFn(inst, tc.testAmount, seed, proc)

			if len(inst.StageAnswers) != 7 {
				t.Fatalf("expected 7 stages, got %d", len(inst.StageAnswers))
			}

			for _, st := range requiredStages {
				ans, ok := inst.StageAnswers[st]
				if !ok {
					t.Fatalf("missing stage answer for %s", st)
				}

				if ans.Prompt == "" {
					t.Errorf("stage %s has empty prompt", st)
				}
				if ans.CorrectOptionID == "" {
					t.Errorf("stage %s has empty correct_option_id", st)
				}
				if len(ans.Options) < 2 {
					t.Errorf("stage %s has fewer than 2 options: %d", st, len(ans.Options))
				}

				// Verify correct option is in options list
				foundCorrect := false
				for _, opt := range ans.Options {
					if opt.ID == ans.CorrectOptionID {
						foundCorrect = true
						break
					}
				}
				if !foundCorrect {
					t.Errorf("stage %s correct option %q not found in options %v", st, ans.CorrectOptionID, ans.Options)
				}

				// Invariant: Socratic hint must end in '?'
				if !strings.HasSuffix(strings.TrimSpace(ans.CausalHint), "?") {
					t.Errorf("stage %s hint %q does not end with '?'", st, ans.CausalHint)
				}

				// Invariant: Stage 1 Neutrality: Explanation MUST NOT include 'debit' or 'credit'
				if st == domain.StageIdentifyAccount {
					explLower := strings.ToLower(ans.Explanation)
					if strings.Contains(explLower, "debit") || strings.Contains(explLower, "credit") {
						t.Errorf("Stage 1 explanation contains debit/credit: %q", ans.Explanation)
					}
				}

				if ans.RelevantConceptID == "" {
					t.Errorf("stage %s has empty relevant_concept_id", st)
				}
				if ans.Explanation == "" {
					t.Errorf("stage %s has empty explanation", st)
				}
			}
		})
	}
}

func TestStreamM04PedagogyPolicy(t *testing.T) {
	pedagogyPath := "../../curriculum/candidates/pedagogy_m04.json"
	data, err := os.ReadFile(pedagogyPath)
	if err != nil {
		t.Fatalf("failed to read stream pedagogy file %s: %v", pedagogyPath, err)
	}

	policy, err := bank.LoadPedagogyPolicy(data)
	if err != nil {
		t.Fatalf("failed to validate stream pedagogy policy: %v", err)
	}

	if policy.SchemaVersion != 1 || policy.PolicyVersion < 1 {
		t.Errorf("invalid policy versions: schema=%d, policy=%d", policy.SchemaVersion, policy.PolicyVersion)
	}
	if len(policy.Questions) != 4 {
		t.Errorf("expected 4 pedagogy questions, got %d", len(policy.Questions))
	}
	for _, pq := range policy.Questions {
		if !strings.HasPrefix(pq.QuestionID, "q_m04_rev_") {
			t.Errorf("unexpected question_id in pedagogy: %s", pq.QuestionID)
		}
		if pq.SettingGroup == "" {
			t.Errorf("empty setting_group for %s", pq.QuestionID)
		}
	}
}

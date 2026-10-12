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

func buildM05TestCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	// Add required accounts for M05 and basic catalog
	_ = catalog.Add(domain.Account{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "inventory", Name: "Inventory", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "equipment", Name: "Equipment", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_payable", Name: "Accounts Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "cost_of_goods_sold", Name: "Cost of Goods Sold", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "insurance_expense", Name: "Insurance Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "bad_debt_expense", Name: "Bad Debt Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	return catalog
}

func loadM05CandidateBatch(t *testing.T, path string) *bank.QuestionBankFile {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read candidate batch file %s: %v", path, err)
	}

	var batch bank.QuestionBankFile
	if err := json.Unmarshal(data, &batch); err != nil {
		t.Fatalf("failed to unmarshal candidate batch JSON: %v", err)
	}

	if batch.SchemaVersion != 1 {
		t.Fatalf("expected schema_version 1, got %d", batch.SchemaVersion)
	}
	if len(batch.Questions) != 4 {
		t.Fatalf("expected exactly 4 questions in Stream-M05 batch, got %d", len(batch.Questions))
	}
	return &batch
}

func TestStreamM05BatchInvariants(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m05.json"
	batch := loadM05CandidateBatch(t, batchPath)
	catalog := buildM05TestCatalog()

	expectedIDs := map[string]bool{
		"q_m05_inv_fifo_costing_01":  false,
		"q_m05_inv_lifo_costing_02":  false,
		"q_m05_inv_nrv_writedown_03": false,
		"q_m05_inv_freight_in_cap_04": false,
	}

	for _, q := range batch.Questions {
		if _, ok := expectedIDs[q.ID]; !ok {
			t.Errorf("unexpected question ID: %s", q.ID)
		}
		expectedIDs[q.ID] = true

		if q.Module != 5 {
			t.Errorf("question %s: expected Module 5, got %d", q.ID, q.Module)
		}
		if q.Status != bank.StatusCandidatePendingReview {
			t.Errorf("question %s: expected status %q, got %q", q.ID, bank.StatusCandidatePendingReview, q.Status)
		}

		// Verify parameters: at least 3 distinct positive integer values
		amounts, ok := q.Parameters["amount_minor_units"]
		if !ok || len(amounts) < 3 {
			t.Fatalf("question %s: expected at least 3 amounts, got %v", q.ID, amounts)
		}
		for _, amt := range amounts {
			if amt <= 0 {
				t.Fatalf("question %s: amounts must be positive, got %d", q.ID, amt)
			}
		}

		// Verify expected fixture balancing
		if len(q.ExpectedFixture.Postings) < 2 {
			t.Fatalf("question %s: fixture must have at least 2 postings", q.ID)
		}
		for _, amt := range amounts {
			var debitSum, creditSum int64
			for _, p := range q.ExpectedFixture.Postings {
				if !catalog.Has(domain.AccountID(p.AccountID)) {
					t.Fatalf("question %s: unknown account %s in catalog", q.ID, p.AccountID)
				}
				if p.Side == string(domain.SideDebit) {
					debitSum += amt
				} else if p.Side == string(domain.SideCredit) {
					creditSum += amt
				} else {
					t.Fatalf("question %s: invalid side %s", q.ID, p.Side)
				}
			}
			if debitSum != creditSum {
				t.Fatalf("question %s: unbalanced fixture! Debits=%d, Credits=%d", q.ID, debitSum, creditSum)
			}
			if debitSum != amt {
				t.Fatalf("question %s: debit sum %d does not match transaction amount %d", q.ID, debitSum, amt)
			}
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("missing expected question ID: %s", id)
		}
	}
}

func TestStreamM05EngineRules(t *testing.T) {
	catalog := buildM05TestCatalog()

	rules := map[string]engine.EventRule{
		engine.FamilyM05InventoryFIFO:         engine.M05InventoryFIFORule{},
		engine.FamilyM05InventoryLIFO:         engine.M05InventoryLIFORule{},
		engine.FamilyM05InventoryNRVWriteDown: engine.M05InventoryNRVWriteDownRule{},
		engine.FamilyM05FreightInCapitalize:   engine.M05FreightInCapitalizeRule{},
	}

	testAmounts := []int64{100000, 250000, 500000}

	for familyID, rule := range rules {
		for _, amt := range testAmounts {
			params := map[string]int64{"amount_minor_units": amt}
			entry, err := rule.DeriveEntry(params, catalog)
			if err != nil {
				t.Fatalf("rule %s failed to derive entry: %v", familyID, err)
			}

			if !entry.IsBalanced() {
				t.Fatalf("rule %s: entry is not balanced: %+v", familyID, entry)
			}

			if err := entry.Validate(catalog); err != nil {
				t.Fatalf("rule %s: entry validation failed: %v", familyID, err)
			}

			// Validate equation balance
			effects, err := entry.EquationEffects(catalog)
			if err != nil {
				t.Fatalf("rule %s: equation effect calculation failed: %v", familyID, err)
			}
			if effects.DeltaAssets.Cents() != effects.DeltaLiabilities.Cents()+effects.DeltaEquity.Cents() {
				t.Fatalf("rule %s: fundamental equation violated! Assets=%d, Liab=%d, Equity=%d",
					familyID, effects.DeltaAssets.Cents(), effects.DeltaLiabilities.Cents(), effects.DeltaEquity.Cents())
			}

			explanation := rule.Explain(params)
			if explanation.Summary == "" || explanation.DebitRationale == "" || explanation.CreditRationale == "" || explanation.EquationSummary == "" {
				t.Fatalf("rule %s: explanation has empty fields: %+v", familyID, explanation)
			}
		}
	}
}

func TestStreamM05StageBuildersAndPedagogy(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m05.json"
	batch := loadM05CandidateBatch(t, batchPath)

	for _, q := range batch.Questions {
		for _, amt := range q.Parameters["amount_minor_units"] {
			money := domain.NewMoney(amt)
			seed := int64(1005)

			inst := &domain.QuestionInstance{
				QuestionID:   q.ID,
				StageAnswers: make(map[domain.DrillStage]domain.StageAnswer),
			}

			proc := engine.ProcessedEvent{
				Event: engine.TransactionEvent{
					FamilyID:   q.FamilyID,
					Parameters: map[string]int64{"amount_minor_units": amt},
				},
			}

			handled := drill.BuildM05Stages(q.FamilyID, inst, money, seed, proc)
			if !handled {
				t.Fatalf("unhandled family ID in BuildM05Stages: %s", q.FamilyID)
			}

			if len(inst.StageAnswers) != 7 {
				t.Fatalf("question %s: expected 7 stage answers, got %d", q.ID, len(inst.StageAnswers))
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

			for _, stage := range requiredStages {
				ans, ok := inst.StageAnswers[stage]
				if !ok {
					t.Fatalf("question %s: missing stage %s", q.ID, stage)
				}

				if ans.CorrectOptionID == "" {
					t.Fatalf("question %s stage %s: missing CorrectOptionID", q.ID, stage)
				}

				// Check that CorrectOptionID exists in Options
				foundOption := false
				for _, opt := range ans.Options {
					if opt.ID == ans.CorrectOptionID {
						foundOption = true
						break
					}
				}
				if !foundOption {
					t.Fatalf("question %s stage %s: CorrectOptionID %s not in options", q.ID, stage, ans.CorrectOptionID)
				}

				// Invariant: Socratic hint must end with '?'
				trimmedHint := strings.TrimSpace(ans.CausalHint)
				if !strings.HasSuffix(trimmedHint, "?") {
					t.Fatalf("question %s stage %s: CausalHint must end with '?', got: %q", q.ID, stage, trimmedHint)
				}

				// Invariant: Stage 1 explanation MUST NOT include "debit" or "credit"
				if stage == domain.StageIdentifyAccount {
					lowerExpl := strings.ToLower(ans.Explanation)
					if strings.Contains(lowerExpl, "debit") || strings.Contains(lowerExpl, "credit") {
						t.Fatalf("question %s Stage 1 explanation violates neutrality invariant (contains debit/credit): %s", q.ID, ans.Explanation)
					}
				}
			}
		}
	}
}

func TestStreamM05PedagogyPolicy(t *testing.T) {
	pedagogyPath := "../../curriculum/candidates/pedagogy_m05.json"
	data, err := os.ReadFile(pedagogyPath)
	if err != nil {
		t.Fatalf("failed to read pedagogy candidate file: %v", err)
	}

	var policy struct {
		SchemaVersion int `json:"schema_version"`
		Questions     []struct {
			QuestionID        string   `json:"question_id"`
			QuestionVersion   int      `json:"question_version"`
			SettingGroup      string `json:"setting_group"`
			Contrasts         []any  `json:"contrasts"`
			DistractorProfile string `json:"distractor_profile"`
		} `json:"questions"`
	}

	if err := json.Unmarshal(data, &policy); err != nil {
		t.Fatalf("failed to parse pedagogy JSON: %v", err)
	}

	if len(policy.Questions) != 4 {
		t.Fatalf("expected 4 questions in pedagogy policy, got %d", len(policy.Questions))
	}

	for _, pq := range policy.Questions {
		if pq.QuestionID == "" || pq.SettingGroup == "" || pq.DistractorProfile == "" {
			t.Fatalf("incomplete pedagogy record: %+v", pq)
		}
	}
}

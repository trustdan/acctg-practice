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

func buildM06TestCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	_ = catalog.Add(domain.Account{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "equipment", Name: "Equipment", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accumulated_depreciation", Name: "Accumulated Depreciation", Category: domain.CategoryAsset, NormalSide: domain.SideCredit, ContraOf: "equipment"})
	_ = catalog.Add(domain.Account{ID: "accounts_payable", Name: "Accounts Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "retained_earnings", Name: "Retained Earnings", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "gain_on_disposal", Name: "Gain on Disposal", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "depreciation_expense", Name: "Depreciation Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "loss_on_disposal", Name: "Loss on Disposal", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "loss_on_impairment", Name: "Loss on Impairment", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "insurance_expense", Name: "Insurance Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "cost_of_goods_sold", Name: "Cost of Goods Sold", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	return catalog
}

func loadM06CandidateBatch(t *testing.T, path string) *bank.QuestionBankFile {
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
		t.Fatalf("expected exactly 4 questions in Stream-M06 batch, got %d", len(batch.Questions))
	}
	return &batch
}

func TestStreamM06BatchInvariants(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m06.json"
	batch := loadM06CandidateBatch(t, batchPath)
	catalog := buildM06TestCatalog()

	expectedIDs := map[string]bool{
		"q_m06_ppe_straight_line_depr_01":      false,
		"q_m06_ppe_macrs_tax_timing_02":        false,
		"q_m06_ppe_asset_disposal_gain_loss_03": false,
		"q_m06_ppe_impairment_recoverability_04": false,
	}

	for _, q := range batch.Questions {
		if _, ok := expectedIDs[q.ID]; !ok {
			t.Errorf("unexpected question ID: %s", q.ID)
		}
		expectedIDs[q.ID] = true

		if q.Module != 6 {
			t.Errorf("question %s: expected Module 6, got %d", q.ID, q.Module)
		}
		if q.Status != bank.StatusCandidatePendingReview {
			t.Errorf("question %s: expected status %q, got %q", q.ID, bank.StatusCandidatePendingReview, q.Status)
		}

		// Verify scenario template has placeholder
		if !strings.Contains(q.ScenarioTemplate, "${amount_dollars}") {
			t.Errorf("question %s: scenario_template must contain ${amount_dollars}", q.ID)
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

		for idx := range amounts {
			var debitSum, creditSum int64
			for _, p := range q.ExpectedFixture.Postings {
				if !catalog.Has(domain.AccountID(p.AccountID)) {
					t.Fatalf("question %s: unknown account %s in catalog", q.ID, p.AccountID)
				}
				paramVals, ok := q.Parameters[p.AmountParameter]
				if !ok || len(paramVals) == 0 {
					t.Fatalf("question %s: fixture posting references undefined amount parameter %s", q.ID, p.AmountParameter)
				}
				pAmt := paramVals[idx%len(paramVals)]
				if p.Side == string(domain.SideDebit) {
					debitSum += pAmt
				} else if p.Side == string(domain.SideCredit) {
					creditSum += pAmt
				} else {
					t.Fatalf("question %s: invalid side %s", q.ID, p.Side)
				}
			}
			if debitSum != creditSum {
				t.Fatalf("question %s sample %d: unbalanced fixture! Debits=%d, Credits=%d", q.ID, idx, debitSum, creditSum)
			}
			if debitSum <= 0 {
				t.Fatalf("question %s sample %d: debit sum must be positive, got %d", q.ID, idx, debitSum)
			}
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("missing expected question ID: %s", id)
		}
	}
}

func TestStreamM06EngineRules(t *testing.T) {
	catalog := buildM06TestCatalog()

	rules := map[string]engine.EventRule{
		engine.FamilyM06StraightLineDepr:         engine.M06PPEStraightLineDeprRule{},
		engine.FamilyM06MACRSTaxTiming:           engine.M06PPEMACRSTaxTimingRule{},
		engine.FamilyM06AssetDisposalGainLoss:    engine.M06PPEAssetDisposalGainLossRule{},
		engine.FamilyM06ImpairmentRecoverability: engine.M06PPEImpairmentRecoverabilityRule{},
	}

	testAmounts := []int64{500000, 1000000, 2000000}

	for familyID, rule := range rules {
		for _, amt := range testAmounts {
			params := map[string]int64{
				"amount_minor_units":     amt,
				"accum_depr_minor_units": amt * 80 / 100,
				"cash_minor_units":       amt * 25 / 100,
				"gain_minor_units":       amt * 5 / 100,
			}
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

			// Validate equation balance: DeltaAssets == DeltaLiabilities + DeltaEquity
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

func TestStreamM06StageBuildersAndPedagogy(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m06.json"
	batch := loadM06CandidateBatch(t, batchPath)

	for _, q := range batch.Questions {
		for pIdx, amt := range q.Parameters["amount_minor_units"] {
			money := domain.NewMoney(amt)
			seed := int64(1006)

			inst := &domain.QuestionInstance{
				QuestionID:   q.ID,
				StageAnswers: make(map[domain.DrillStage]domain.StageAnswer),
			}

			procParams := map[string]int64{"amount_minor_units": amt}
			for pName, pVals := range q.Parameters {
				procParams[pName] = pVals[pIdx%len(pVals)]
			}

			proc := engine.ProcessedEvent{
				Event: engine.TransactionEvent{
					FamilyID:   q.FamilyID,
					Parameters: procParams,
				},
			}

			handled := drill.BuildM06Stages(q.FamilyID, inst, money, seed, proc)
			if !handled {
				t.Fatalf("unhandled family ID in BuildM06Stages: %s", q.FamilyID)
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

				// Invariant: Stage 1 explanation MUST NOT include "debit" or "credit" (case-insensitive)
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

func TestStreamM06PedagogyPolicy(t *testing.T) {
	pedagogyPath := "../../curriculum/candidates/pedagogy_m06.json"
	data, err := os.ReadFile(pedagogyPath)
	if err != nil {
		t.Fatalf("failed to read pedagogy candidate file: %v", err)
	}

	var policy struct {
		SchemaVersion int `json:"schema_version"`
		Questions     []struct {
			QuestionID        string   `json:"question_id"`
			QuestionVersion   int      `json:"question_version"`
			SettingGroup      string   `json:"setting_group"`
			Contrasts         []any    `json:"contrasts"`
			DistractorProfile string   `json:"distractor_profile"`
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

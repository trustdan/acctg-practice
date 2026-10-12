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

func buildM07TestCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	_ = catalog.Add(domain.Account{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_payable", Name: "Accounts Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "notes_payable", Name: "Notes Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "bonds_payable", Name: "Bonds Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "discount_on_bonds_payable", Name: "Discount on Bonds Payable", Category: domain.CategoryLiability, NormalSide: domain.SideDebit, ContraOf: "bonds_payable"})
	_ = catalog.Add(domain.Account{ID: "premium_on_bonds_payable", Name: "Premium on Bonds Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "contingent_liability", Name: "Contingent Liability", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "retained_earnings", Name: "Retained Earnings", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "interest_expense", Name: "Interest Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "litigation_loss", Name: "Litigation Loss", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	return catalog
}

func loadM07CandidateBatch(t *testing.T, path string) *bank.QuestionBankFile {
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
		t.Fatalf("expected exactly 4 questions in Stream-M07 batch, got %d", len(batch.Questions))
	}
	return &batch
}

func TestStreamM07BatchInvariants(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m07.json"
	batch := loadM07CandidateBatch(t, batchPath)
	catalog := buildM07TestCatalog()

	expectedIDs := map[string]bool{
		"q_m07_bnd_issuance_discount_01":       false,
		"q_m07_bnd_issuance_premium_02":        false,
		"q_m07_bnd_effective_interest_amort_03": false,
		"q_m07_bnd_contingent_liability_04":    false,
	}

	for _, q := range batch.Questions {
		if _, ok := expectedIDs[q.ID]; !ok {
			t.Errorf("unexpected question ID: %s", q.ID)
		}
		expectedIDs[q.ID] = true

		if q.Module != 7 {
			t.Errorf("question %s: expected Module 7, got %d", q.ID, q.Module)
		}
		if q.Status != bank.StatusCandidatePendingReview {
			t.Errorf("question %s: expected status %q, got %q", q.ID, bank.StatusCandidatePendingReview, q.Status)
		}

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
				paramList, ok := q.Parameters[p.AmountParameter]
				if !ok || len(paramList) <= idx {
					t.Fatalf("question %s: missing param %s at index %d", q.ID, p.AmountParameter, idx)
				}
				postingAmt := paramList[idx]
				if p.Side == string(domain.SideDebit) {
					debitSum += postingAmt
				} else if p.Side == string(domain.SideCredit) {
					creditSum += postingAmt
				} else {
					t.Fatalf("question %s: invalid side %s", q.ID, p.Side)
				}
			}
			if debitSum != creditSum {
				t.Fatalf("question %s index %d: unbalanced fixture! Debits=%d, Credits=%d", q.ID, idx, debitSum, creditSum)
			}
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("missing expected question ID: %s", id)
		}
	}
}

func TestStreamM07EngineRules(t *testing.T) {
	catalog := buildM07TestCatalog()

	rules := map[string]engine.EventRule{
		engine.FamilyM07BondIssuanceDiscount:       engine.M07BondIssuanceDiscountRule{},
		engine.FamilyM07BondIssuancePremium:        engine.M07BondIssuancePremiumRule{},
		engine.FamilyM07BondEffectiveInterestAmort: engine.M07BondEffectiveInterestAmortRule{},
		engine.FamilyM07BondContingentLiability:    engine.M07BondContingentLiabilityRule{},
	}

	testCases := map[string][]map[string]int64{
		engine.FamilyM07BondIssuanceDiscount: {
			{"amount_minor_units": 100000000, "cash_minor_units": 95000000, "discount_minor_units": 5000000},
			{"amount_minor_units": 200000000, "cash_minor_units": 190000000, "discount_minor_units": 10000000},
			{"amount_minor_units": 400000000, "cash_minor_units": 380000000, "discount_minor_units": 20000000},
		},
		engine.FamilyM07BondIssuancePremium: {
			{"amount_minor_units": 100000000, "cash_minor_units": 105000000, "premium_minor_units": 5000000},
			{"amount_minor_units": 200000000, "cash_minor_units": 210000000, "premium_minor_units": 10000000},
			{"amount_minor_units": 400000000, "cash_minor_units": 420000000, "premium_minor_units": 20000000},
		},
		engine.FamilyM07BondEffectiveInterestAmort: {
			{"amount_minor_units": 3800000, "cash_minor_units": 3000000, "discount_minor_units": 800000},
			{"amount_minor_units": 7600000, "cash_minor_units": 6000000, "discount_minor_units": 1600000},
			{"amount_minor_units": 15200000, "cash_minor_units": 12000000, "discount_minor_units": 3200000},
		},
		engine.FamilyM07BondContingentLiability: {
			{"amount_minor_units": 25000000},
			{"amount_minor_units": 50000000},
			{"amount_minor_units": 100000000},
		},
	}

	for familyID, rule := range rules {
		cases, ok := testCases[familyID]
		if !ok {
			t.Fatalf("missing test cases for family: %s", familyID)
		}
		for _, params := range cases {
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

			// Validate fundamental equation balance
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

func TestStreamM07StageBuildersAndPedagogy(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m07.json"
	batch := loadM07CandidateBatch(t, batchPath)

	for _, q := range batch.Questions {
		amounts := q.Parameters["amount_minor_units"]
		for idx, amt := range amounts {
			money := domain.NewMoney(amt)
			seed := int64(2007 + idx)

			params := map[string]int64{"amount_minor_units": amt}
			for k, vList := range q.Parameters {
				if len(vList) > idx {
					params[k] = vList[idx]
				}
			}

			inst := &domain.QuestionInstance{
				QuestionID:   q.ID,
				StageAnswers: make(map[domain.DrillStage]domain.StageAnswer),
			}

			proc := engine.ProcessedEvent{
				Event: engine.TransactionEvent{
					FamilyID:   q.FamilyID,
					Parameters: params,
				},
			}

			handled := drill.BuildM07Stages(q.FamilyID, inst, money, seed, proc)
			if !handled {
				t.Fatalf("unhandled family ID in BuildM07Stages: %s", q.FamilyID)
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

func TestStreamM07PedagogyPolicy(t *testing.T) {
	pedagogyPath := "../../curriculum/candidates/pedagogy_m07.json"
	data, err := os.ReadFile(pedagogyPath)
	if err != nil {
		t.Fatalf("failed to read pedagogy candidate file: %v", err)
	}

	var policy struct {
		SchemaVersion int `json:"schema_version"`
		Questions     []struct {
			QuestionID        string `json:"question_id"`
			QuestionVersion   int    `json:"question_version"`
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

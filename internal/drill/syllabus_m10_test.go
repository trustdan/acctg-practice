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

func buildM10TestCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	_ = catalog.Add(domain.Account{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "equipment", Name: "Equipment", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_payable", Name: "Accounts Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "notes_payable", Name: "Notes Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "retained_earnings", Name: "Retained Earnings", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	return catalog
}

func loadM10CandidateBatch(t *testing.T, path string) *bank.QuestionBankFile {
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
	if len(batch.Questions) != 3 {
		t.Fatalf("expected 3 questions in Stream-M10 candidate batch, got %d", len(batch.Questions))
	}
	return &batch
}

func TestStreamM10AccountingRulesAndEquationEffects(t *testing.T) {
	catalog := buildM10TestCatalog()

	rules := map[string]engine.EventRule{
		engine.FamilyM10ForecastingOperatingExpenseAP: engine.M10ForecastingOperatingExpenseAPRule{},
		engine.FamilyM10ForecastingCapEx:              engine.M10ForecastingCapExRule{},
		engine.FamilyM10ForecastingDebtPlug:           engine.M10ForecastingDebtPlugRule{},
	}

	testAmounts := []int64{1800000, 3600000, 7200000}

	for familyID, rule := range rules {
		for _, amt := range testAmounts {
			params := map[string]int64{
				"amount_minor_units": amt,
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

func TestStreamM10StageBuildersAndPedagogy(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m10.json"
	batch := loadM10CandidateBatch(t, batchPath)

	for _, q := range batch.Questions {
		for pIdx, amt := range q.Parameters["amount_minor_units"] {
			money := domain.NewMoney(amt)
			seed := int64(1010)

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

			handled := drill.BuildM10Stages(q.FamilyID, inst, money, seed, proc)
			if !handled {
				t.Fatalf("unhandled family ID in BuildM10Stages: %s", q.FamilyID)
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

				found := false
				for _, opt := range ans.Options {
					if opt.ID == ans.CorrectOptionID {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("question %s stage %s: correct option %s not found in options", q.ID, stage, ans.CorrectOptionID)
				}

				// Pedagogical Invariant: All CausalHints must end with a question mark ('?')
				if !strings.HasSuffix(strings.TrimSpace(ans.CausalHint), "?") {
					t.Errorf("question %s stage %s CausalHint must end with '?': %q", q.ID, stage, ans.CausalHint)
				}

				// Pedagogical Invariant: Stage 1 explanation must omit debit/credit cues
				if stage == domain.StageIdentifyAccount {
					lowerExp := strings.ToLower(ans.Explanation)
					if strings.Contains(lowerExp, "debit") || strings.Contains(lowerExp, "credit") {
						t.Errorf("question %s Stage 1 explanation leaked debit/credit cues: %q", q.ID, ans.Explanation)
					}
				}
			}
		}
	}
}

func TestStreamM10CandidateValidation(t *testing.T) {
	batchPath := "../../curriculum/candidates/batch_m10.json"
	batch := loadM10CandidateBatch(t, batchPath)

	catalog := buildM10TestCatalog()

	for _, q := range batch.Questions {
		if !bank.IsSupportedFamily(q.FamilyID) {
			t.Fatalf("question %s has unsupported family %s", q.ID, q.FamilyID)
		}
		if q.Module != 10 {
			t.Fatalf("question %s module must be 10, got %d", q.ID, q.Module)
		}

		// Verify expected fixture balances
		var totalDebit, totalCredit int64
		for _, p := range q.ExpectedFixture.Postings {
			acc, ok := catalog.Get(domain.AccountID(p.AccountID))
			if !ok {
				t.Fatalf("question %s uses unknown account %s", q.ID, p.AccountID)
			}
			_ = acc
			amt := q.Parameters[p.AmountParameter][0]
			if p.Side == string(domain.SideDebit) {
				totalDebit += amt
			} else {
				totalCredit += amt
			}
		}
		if totalDebit != totalCredit {
			t.Fatalf("question %s expected fixture does not balance: Dr %d != Cr %d", q.ID, totalDebit, totalCredit)
		}
	}
}

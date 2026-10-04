package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"strings"
	"testing"
)

func TestExpansionBatch11SemanticsTeachingAndContrasts(t *testing.T) {
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH11.json", map[string][2]domain.AccountID{
		"cash_rent_separate_insurance":       {"rent_expense", "cash"},
		"borrow_cash_collateral":             {"cash", "notes_payable"},
		"issue_shares_existing_lender":       {"cash", "common_stock"},
		"prepaid_purchase_budget_label":      {"prepaid_insurance", "cash"},
		"prepaid_consumption_recent_payment": {"insurance_expense", "prepaid_insurance"},
	}, [][2]string{
		{"cash_rent_separate_insurance", "prepaid_purchase_budget_label"},
		{"borrow_cash_collateral", "issue_shares_existing_lender"},
		{"borrow_cash_collateral", "equipment_cash_existing_note"},
		{"issue_shares_existing_lender", "borrow_cash_stockholder_note"},
		{"prepaid_purchase_budget_label", "prepaid_consumption_recent_payment"},
	})
}

func TestPublishedExpansionBatch11MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH11.json")
}

func TestExpansionBatch11PreservesUnrelatedRecordedBalances(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH11.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range pack.Questions {
		q.Status = bank.StatusApprovedActive
		for _, amount := range q.Parameters["amount_minor_units"] {
			ledger := domain.NewTAccountLedger(catalog)
			opening := map[domain.AccountID]int64{"cash": 4 * amount, "prepaid_insurance": 3 * amount, "equipment": 5 * amount, "notes_payable": 2 * amount, "common_stock": 6 * amount}
			for id, n := range opening {
				side := domain.SideDebit
				if id == "notes_payable" || id == "common_stock" {
					side = domain.SideCredit
				}
				ledger.GetOrCreate(id).Post(side, domain.NewMoney(n))
			}
			inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range inst.Entry.Postings {
				ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
			}
			want := map[domain.AccountID]int64{"cash": 4 * amount, "prepaid_insurance": 3 * amount, "equipment": 5 * amount, "notes_payable": 2 * amount, "common_stock": 6 * amount, "rent_expense": 0, "insurance_expense": 0}
			switch q.FamilyID {
			case "cash_rent":
				want["cash"] -= amount
				want["rent_expense"] = amount
			case "borrow_cash":
				want["cash"] += amount
				want["notes_payable"] += amount
			case "issue_shares":
				want["cash"] += amount
				want["common_stock"] += amount
			case "prepaid_purchase":
				want["cash"] -= amount
				want["prepaid_insurance"] += amount
			case "prepaid_consumption":
				want["prepaid_insurance"] -= amount
				want["insurance_expense"] = amount
			}
			for id, n := range want {
				if ledger.GetOrCreate(id).NetBalance.Cents() != n {
					t.Fatalf("%s %s altered unrelated balance or repeated payment", q.ID, id)
				}
			}
			if q.FamilyID == "prepaid_consumption" {
				for _, s := range inst.StageAnswers {
					text := s.Prompt + s.CausalHint + s.Explanation
					for _, h := range s.MistakeHints {
						text += h
					}
					if strings.Contains(text, "last month") {
						t.Fatal("teaching contradicts same-month purchase")
					}
				}
			}
		}
	}
}

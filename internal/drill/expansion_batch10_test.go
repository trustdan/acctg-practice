package drill_test

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"strings"
	"testing"
)

func TestExpansionBatch10SemanticsTeachingAndContrasts(t *testing.T) {
	expected := map[string][2]domain.AccountID{
		"cash_service_third_party":        {"cash", "service_revenue"},
		"service_on_credit_invoice_later": {"accounts_receivable", "service_revenue"},
		"collect_receivable_third_party":  {"cash", "accounts_receivable"},
		"customer_advance_full_price":     {"cash", "unearned_revenue"},
		"earn_advance_paperwork_later":    {"unearned_revenue", "service_revenue"},
	}
	pairs := [][2]string{
		{"cash_service_third_party", "collect_receivable_third_party"},
		{"service_on_credit_invoice_later", "cash_service_third_party"},
		{"customer_advance_full_price", "cash_service_third_party"},
		{"customer_advance_full_price", "earn_advance_paperwork_later"},
		{"earn_advance_paperwork_later", "service_on_credit_invoice_later"},
	}
	verifyPilotBatch(t, "../../docs/EXPANSION-BATCH10.json", expected, pairs)
}

func TestPublishedExpansionBatch10MatchesReviewedPackage(t *testing.T) {
	verifyPublishedPilot(t, "../../docs/EXPANSION-BATCH10.json")
}

func TestExpansionBatch10RecordedClaimAndAdvanceSettlement(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH10.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{{"service_on_credit_invoice_later", "collect_receivable_third_party"}, {"customer_advance_full_price", "earn_advance_paperwork_later"}} {
		for _, amount := range []int64{5000, 10000, 20000} {
			ledger := domain.NewTAccountLedger(catalog)
			for i, id := range pair {
				q := *getQuestionByID(pack, id)
				q.Status = bank.StatusApprovedActive
				inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
				if err != nil {
					t.Fatal(err)
				}
				for _, p := range inst.Entry.Postings {
					ledger.GetOrCreate(p.AccountID).Post(p.Side, p.Amount)
				}
				if i == 1 && (ledger.GetOrCreate("cash").NetBalance.Cents() != amount || ledger.GetOrCreate("service_revenue").NetBalance.Cents() != amount || ledger.GetOrCreate("accounts_receivable").NetBalance.Cents() != 0 || ledger.GetOrCreate("unearned_revenue").NetBalance.Cents() != 0) {
					t.Fatalf("%v duplicated earning/receipt or left settled balance", pair)
				}
			}
		}
	}
}

func TestExpansionBatch10PromptsPreservePaperworkAndPayerFacts(t *testing.T) {
	gen, _, catalog := setupTestDrill(t)
	pack, err := bank.LoadQuestionBankFile("../../docs/EXPANSION-BATCH10.json", catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range pack.Questions {
		q.Status = bank.StatusApprovedActive
		inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": 10000})
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range inst.StageAnswers {
			if q.FamilyID == "service_on_credit" && (strings.Contains(a.Prompt, "was invoiced") || strings.Contains(a.Prompt, "bills a new")) {
				t.Fatal("prompt invented invoice sent today")
			}
			if q.FamilyID == "collect_receivable" && strings.Contains(a.Prompt, "customer pays cash") {
				t.Fatal("prompt replaced third-party payer")
			}
		}
	}
}

package bank_test

import (
	"bytes"
	"encoding/json"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"testing"
)

func TestReviewedPedagogyPolicyReferencesAndStrictValidation(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	questions, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	policy, err := bank.LoadPedagogyPolicy(curriculum.PedagogyJSON)
	if err != nil {
		t.Fatal(err)
	}
	if err = policy.ValidateReferences(questions.Questions); err != nil {
		t.Fatal(err)
	}
	if len(policy.Questions) != 90 {
		t.Fatalf("expected reviewed settings for 90, got %d", len(policy.Questions))
	}
	for _, q := range questions.Questions {
		p, _, err := bank.ReviewedPedagogy(q)
		if err != nil {
			t.Fatal(err)
		}
		if bank.IsActiveForPractice(q.Status) && (p.PolicyVersion != 2 || p.SettingGroup == "") {
			t.Fatalf("missing reviewed setting: %s", q.ID)
		}
		q.Version++
		p, _, err = bank.ReviewedPedagogy(q)
		if err != nil || p.PolicyVersion != 0 {
			t.Fatal("unreviewed version inherited policy")
		}
	}
	for _, edit := range []func(*bank.PedagogyPolicy){
		func(p *bank.PedagogyPolicy) { p.Questions = append(p.Questions, p.Questions[0]) },
		func(p *bank.PedagogyPolicy) { p.Questions[0].SettingGroup = "" },
		func(p *bank.PedagogyPolicy) { p.Questions[0].DistractorProfile = "unknown" },
		func(p *bank.PedagogyPolicy) { p.Review.Reviewer = nil },
	} {
		p, _ := bank.LoadPedagogyPolicy(curriculum.PedagogyJSON)
		edit(p)
		data, _ := json.Marshal(p)
		if _, err = bank.LoadPedagogyPolicy(data); err == nil {
			t.Fatal("invalid review policy accepted")
		}
	}
	if _, err = bank.LoadPedagogyPolicy(append(curriculum.PedagogyJSON, []byte(" {}")...)); err == nil {
		t.Fatal("trailing data accepted")
	}
	q := questions.Questions[0]
	q.Status = bank.StatusRetired
	for i := range questions.Questions {
		if questions.Questions[i].ID == q.ID {
			questions.Questions[i] = q
		}
	}
	if err = policy.ValidateReferences(questions.Questions); err == nil {
		t.Fatal("inactive source accepted")
	}
}

func TestAmountProfilesRejectDifferentReviewedScenarios(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	seeds, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	for _, profile := range []string{"amount_collect_receivable_installment", "amount_prepaid_consumption_final_remaining"} {
		p, err := bank.LoadPedagogyPolicy(curriculum.PedagogyJSON)
		if err != nil {
			t.Fatal(err)
		}
		// A known profile must not be attached to a different scenario, even in the same family.
		for i := range p.Questions {
			if p.Questions[i].QuestionID == "collect_receivable_basic" {
				p.Questions[i].DistractorProfile = profile
			}
		}
		if err = p.ValidateReferences(seeds.Questions); err == nil {
			t.Fatal("amount profile accepted for unreviewed wording")
		}
	}
}

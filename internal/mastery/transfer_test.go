package mastery

import (
	"bytes"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"math/rand"
	"testing"
	"time"
)

func policyQuestions(t *testing.T) []bank.QuestionJSON {
	t.Helper()
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	return pack.Questions
}
func TestReviewedContrastIsBoundedEligibleAndCannotChain(t *testing.T) {
	qs := policyQuestions(t)
	var source bank.QuestionJSON
	for _, q := range qs {
		if q.ID == "customer_advance_same_day" {
			source = q
		}
	}
	p, _, err := bank.ReviewedPedagogy(source)
	if err != nil {
		t.Fatal(err)
	}
	inst := &domain.QuestionInstance{Pedagogy: p}
	s := NewScheduler(NewMockClock(time.Unix(1, 0)), rand.New(rand.NewSource(101)))
	s.QueueContrast(inst)
	s.QueueContrast(inst)
	q := s.SelectNextQuestion(qs, nil)
	if q.ID != "earn_advance_same_day" || !s.LastSelectionWasContrast {
		t.Fatal("did not select the reviewed partner")
	}
	inst.Pedagogy.Remediation = true
	s.QueueContrast(inst)
	s.SelectNextQuestion(qs, nil)
	if s.LastSelectionWasContrast {
		t.Fatal("guided comparison chained or queue duplicated")
	}
	inst.Pedagogy.Remediation = false
	s.QueueContrast(inst)
	for i := range qs {
		if qs[i].ID == "earn_advance_same_day" {
			qs[i].Status = bank.StatusRetired
		}
	}
	s.SelectNextQuestion(qs, nil)
	if s.LastSelectionWasContrast {
		t.Fatal("selected retired contrast")
	}
	for i := range qs {
		if qs[i].ID == "earn_advance_same_day" {
			qs[i].Status = bank.StatusApprovedActive
			qs[i].Version++
		}
	}
	s.QueueContrast(inst)
	s.SelectNextQuestion(qs, nil)
	if s.LastSelectionWasContrast {
		t.Fatal("selected unreviewed changed contrast")
	}
}
func TestCosmeticReviewedSettingGetsRepeatPenalty(t *testing.T) {
	qs := policyQuestions(t)
	var a, b bank.QuestionJSON
	for _, q := range qs {
		if q.ID == "cash_service_basic" {
			a = q
		}
		if q.ID == "cash_service_webdev" {
			b = q
		}
	}
	now := time.Unix(1, 0)
	s := NewScheduler(NewMockClock(now), rand.New(rand.NewSource(1)))
	before := s.calculateQuestionWeight(b, nil, now)
	s.recordReviewedExposure(a)
	after := s.calculateQuestionWeight(b, nil, now)
	if after >= before {
		t.Fatal("cosmetic wording in same reviewed setting got no repeat penalty")
	}
}
func transferAttempts(groups []string, assistance domain.AssistanceLevel) []domain.Attempt {
	var attempts []domain.Attempt
	for i, group := range groups {
		attempts = append(attempts, domain.Attempt{AttemptID: string(rune('a' + i)), InstanceID: string(rune('a' + i)), ConceptID: "cash_vs_revenue", SettingGroup: group, PedagogyVersion: 1, IsCorrect: true, Assistance: assistance, AnsweredAt: time.Unix(int64(i*1200), 0)})
	}
	return attempts
}
func TestTransferRequiresDifferentReviewedSettingAndUnassistedDelay(t *testing.T) {
	for _, groups := range [][]string{{"same", "same", "same", "same", "same"}, {"", "", "", "", ""}} {
		s := RebuildProjections(transferAttempts(groups, domain.AssistanceNone))["cash_vs_revenue"]
		if s.ScaffoldLevel != domain.ScaffoldFull || s.HalfLifeDays != DefaultHalfLifeDays || s.TransferDelayedSuccesses != 0 {
			t.Fatal("spaced same/unreviewed wording graduated")
		}
	}
	s := RebuildProjections(transferAttempts([]string{"a", "b", "a"}, domain.AssistanceNone))["cash_vs_revenue"]
	if s.ScaffoldLevel != domain.ScaffoldFaded || s.TransferDelayedSuccesses != 2 || s.SuccessfulSettings != 2 || s.EvidenceVersion != 2 {
		t.Fatalf("independent transfer not credited: %+v", s)
	}
	assisted := RebuildProjections(transferAttempts([]string{"a", "b", "a"}, domain.AssistanceContrast))["cash_vs_revenue"]
	if assisted.Alpha != DefaultAlpha || assisted.IndependentAttempts != 0 || assisted.SuccessfulSettings != 0 || assisted.TransferDelayedSuccesses != 0 {
		t.Fatal("contrast exposure inflated independent evidence")
	}
	repeated := RebuildProjections(transferAttempts([]string{"a", "b", "b", "b"}, domain.AssistanceNone))["cash_vs_revenue"]
	if repeated.TransferDelayedSuccesses != 1 || repeated.HalfLifeDays != DefaultHalfLifeDays {
		t.Fatal("same-setting repetition increased transfer half-life")
	}
	attempts := transferAttempts([]string{"a", "b", "a"}, domain.AssistanceNone)
	dupes := append(append([]domain.Attempt(nil), attempts...), attempts...)
	s2 := RebuildProjections(dupes)["cash_vs_revenue"]
	if s2.Alpha != s.Alpha || s2.TransferDelayedSuccesses != s.TransferDelayedSuccesses {
		t.Fatal("duplicate persisted/current attempts inflated projections")
	}
}

func TestAssistedExposurePostponesDelayedTransfer(t *testing.T) {
	atts := transferAttempts([]string{"a", "b", "b"}, domain.AssistanceNone)
	atts[1].Assistance = domain.AssistanceContrast
	atts[2].AnsweredAt = atts[1].AnsweredAt.Add(time.Minute)
	stats := RebuildProjections(atts)["cash_vs_revenue"]
	if stats.TransferDelayedSuccesses != 0 || stats.DelayedSuccesses != 0 || stats.ScaffoldLevel != domain.ScaffoldFull {
		t.Fatal("recent assisted exposure counted as delayed independent transfer")
	}
	atts[2].AnsweredAt = atts[1].AnsweredAt.Add(MinDelayedRetrievalGap)
	stats = RebuildProjections(atts)["cash_vs_revenue"]
	if stats.TransferDelayedSuccesses != 1 || stats.IndependentSuccesses != 2 || stats.AssistedAttempts != 1 || stats.ScaffoldLevel != domain.ScaffoldIntermediate {
		t.Fatal("later independent transfer not separated from guided comparison")
	}
}

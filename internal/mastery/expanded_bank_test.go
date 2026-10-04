package mastery_test

import (
	"bytes"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"testing"
	"time"
)

func TestExpandedBankSubstepsAndRapidVariantsDoNotCreateDelayedMastery(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	gen := drill.NewGenerator(cat, engine.NewEngine(cat))
	want := map[string]int{}
	var attempts []domain.Attempt
	count := 0
	start := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
	for _, q := range pack.Questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": q.Parameters["amount_minor_units"][0]})
		if err != nil {
			t.Fatal(err)
		}
		session := drill.NewSession("expanded-evidence", inst)
		seen := map[string]bool{}
		for !session.IsCompleted {
			st, err := session.CurrentStage()
			if err != nil {
				t.Fatal(err)
			}
			if _, err = session.SubmitOption(st.CorrectOptionID, start.Add(time.Duration(count)*time.Second)); err != nil {
				t.Fatal(err)
			}
			cid := st.RelevantConceptID
			if cid != "" && !seen[cid] {
				seen[cid] = true
				want[cid]++
			}
		}
		attempts = append(attempts, session.Attempts...)
		count++
	}
	if count != 90 || len(attempts) != 630 {
		t.Fatalf("bank evidence coverage %d/%d", count, len(attempts))
	}
	stats := mastery.RebuildProjections(attempts)
	for cid, n := range want {
		s := stats[cid]
		if s == nil || s.IndependentAttempts != n || s.IndependentSuccesses != n {
			t.Fatalf("substeps multiplied evidence for %s", cid)
		}
		if s.HasDelayedRetrieval || s.ScaffoldLevel == domain.ScaffoldFaded || s.HalfLifeDays != mastery.DefaultHalfLifeDays {
			t.Fatalf("rapid bank variants created delayed mastery %s", cid)
		}
	}
}

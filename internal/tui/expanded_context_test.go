package tui

import (
	"bytes"
	"context"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/tutor"
	"strings"
	"testing"
)

func TestExpandedBankTeachingAndTutorContext(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	gen := drill.NewGenerator(cat, engine.NewEngine(cat))
	offline := tutor.NewOfflineTutor()
	checked := 0
	for _, q := range pack.Questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		for _, amount := range q.Parameters["amount_minor_units"] {
			for _, level := range []domain.ScaffoldLevel{domain.ScaffoldFull, domain.ScaffoldIntermediate, domain.ScaffoldFaded} {
				inst, err := gen.GenerateInstanceWithScaffold(q, 101, map[string]int64{"amount_minor_units": amount}, level)
				if err != nil {
					t.Fatal(err)
				}
				m := &Model{CurrentInstance: inst, Session: drill.NewSession("coverage", inst)}
				for _, key := range drill.ScaffoldStageSequence(level) {
					st := inst.StageAnswers[key]
					if st.Explanation == "" || st.CausalHint == "" || strings.Contains(st.Explanation+st.CausalHint, "${") {
						t.Fatalf("%s %s incomplete teaching", q.ID, key)
					}
					for i, opt := range st.Options {
						m.SelectedOptionIndex = i
						req := m.buildTutorRequest(&st)
						if req.ProblemPrompt != inst.PromptText || req.Explanation != st.Explanation || req.SelectedOption.ID != opt.ID {
							t.Fatalf("%s incorrect tutor context", q.ID)
						}
						hint, err := offline.Hint(context.Background(), req)
						if err != nil || hint.Text == "" {
							t.Fatalf("missing hint %s", q.ID)
						}
						if st.MistakeHints[opt.ID] != "" && hint.Text != st.MistakeHints[opt.ID] {
							t.Fatalf("stored mistake hint lost %s", q.ID)
						}
						explanation, err := offline.Explain(context.Background(), req)
						if err != nil || explanation.Text != st.Explanation {
							t.Fatalf("reviewed teaching lost %s", q.ID)
						}
						if !strings.Contains(tutor.FormatUserPrompt(req, false), st.Explanation) || strings.Contains(tutor.FormatUserPrompt(req, true), "Reviewed Stage Explanation:") {
							t.Fatalf("explanation/hint context boundary broken %s", q.ID)
						}
					}
				}
				checked++
			}
		}
	}
	if checked != 810 {
		t.Fatalf("expected 810 amount/scaffold combinations, got %d", checked)
	}
}

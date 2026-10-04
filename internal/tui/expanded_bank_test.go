package tui_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/tui"
	"strings"
	"testing"
	"time"
)

func TestExpandedBankJournalAndRecapsAllAmounts(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	checked := 0
	for _, q := range m.Questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		for _, amount := range q.Parameters["amount_minor_units"] {
			inst, err := m.Generator.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": amount})
			if err != nil {
				t.Fatal(err)
			}
			m.CurrentInstance = inst
			m.Session = drill.NewSession(m.SessionID, inst)
			m.State = tui.StateDrill
			sendKey(m, "J")
			if m.JournalScenario != inst.PromptText || strings.Contains(m.JournalScenario, "${") {
				t.Fatalf("%s journal wording mismatch", q.ID)
			}
			m.JournalLines = inst.Entry.Postings
			sendKey(m, "s")
			if m.JournalFeedback == nil || !m.JournalFeedback.IsCorrect || m.JournalReconciliation == nil {
				t.Fatalf("%s journal semantics not reconciled", q.ID)
			}
			for !m.Session.IsCompleted {
				st, err := m.Session.CurrentStage()
				if err != nil {
					t.Fatal(err)
				}
				if _, err = m.Session.SubmitOption(st.CorrectOptionID, time.Unix(1, 0)); err != nil {
					t.Fatal(err)
				}
			}
			m.State = tui.StateRecap
			if !m.Session.Recap().IsBalanced || !strings.Contains(m.View(), "RECONCILIATION VERIFIED") {
				t.Fatalf("%s recap missing reconciliation", q.ID)
			}
			checked++
		}
	}
	if checked != 270 {
		t.Fatalf("expected 270 scenario/amount checks, got %d", checked)
	}
}

func TestJournalNextUsesActiveRenderedAllowedAmountsWithoutChangingDrillIndex(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	sendKey(m, "J")
	drillIndex := m.CurrentQuestionIndex
	visited := map[string]bool{}
	for i := 0; i < 90; i++ {
		sendKey(m, "n")
		q := m.Questions[m.JournalQuestionIndex]
		if !bank.IsActiveForPractice(q.Status) || strings.Contains(m.JournalScenario, "${") {
			t.Fatalf("inactive or unrendered journal scenario %s", q.ID)
		}
		amount := m.JournalEvent.Parameters["amount_minor_units"]
		allowed := false
		for _, n := range q.Parameters["amount_minor_units"] {
			if amount == n {
				allowed = true
			}
		}
		if !allowed || m.JournalCanonicalEntry.Postings[0].Amount.Cents() != amount {
			t.Fatalf("%s prompt and grading amount disagree", q.ID)
		}
		if m.CurrentQuestionIndex != drillIndex {
			t.Fatal("journal navigation changed drill position")
		}
		visited[q.ID] = true
	}
	if len(visited) != 90 {
		t.Fatalf("journal failed to traverse active bank: %d", len(visited))
	}
}

func TestExamResumeUsesSavedQuestionsAfterCurrentBankChanges(t *testing.T) {
	base, db := setupTestTUI(t)
	defer db.Close()
	cfg := tui.Config{DB: db, Catalog: base.Catalog, Questions: base.Questions, TotalQuestions: 3, Seed: 101, Clock: base.Clock, ExamMode: true, Tutor: base.Tutor, AuthStore: base.AuthStore, ModelCache: base.ModelCache}
	first, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatal(err)
	}
	original := first.ExamRunner.Questions[0].Instance
	sendSpecialKey(first, tea.KeyEnter)
	sendKey(first, "q")
	cfg.Questions = append([]bank.QuestionJSON(nil), base.Questions...)
	for i := range cfg.Questions {
		cfg.Questions[i].ScenarioTemplate = "Changed bank wording ${amount_dollars}"
		cfg.Questions[i].Version++
	}
	cfg.Questions = append(cfg.Questions, cfg.Questions[0])
	for _, auto := range []bool{true, false} {
		cfg.ResumeExam = auto
		resumed, err := tui.NewModel(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if !auto {
			sendKey(resumed, "r")
		}
		if resumed.State != tui.StateExam || resumed.ExamRunner.Questions[0].Instance.PromptText != original.PromptText || resumed.ExamRunner.Questions[0].Instance.Version != original.Version {
			t.Fatal("resume regenerated changed bank instead of original snapshot")
		}
		sendKey(resumed, "q")
	}
}

func TestExamResumeMissingSnapshotsPreservesInterruptedSession(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	record := exam.ExamSessionRecord{ID: "missing-snapshots", TotalQuestions: 2, Status: exam.ExamStatusInterrupted, StartedAt: time.Now(), Seed: 101}
	if err := db.SaveExamSession(record); err != nil {
		t.Fatal(err)
	}
	m.InterruptedExam = &record
	m.State = tui.StateExamResumePrompt
	sendKey(m, "r")
	if m.State != tui.StateExamResumePrompt || m.InterruptedExam == nil || !strings.Contains(m.View(), "interrupted session is preserved") {
		t.Fatal("failed resume silently replaced exam")
	}
	saved, err := db.GetExamSession(record.ID)
	if err != nil || saved.Status != exam.ExamStatusInterrupted {
		t.Fatal("failed resume changed persisted session")
	}
}

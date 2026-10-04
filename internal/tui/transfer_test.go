package tui

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tutor"
	"strings"
	"testing"
	"time"
)

func TestTUIReviewedContrastFollowsErrorAndPersistsAssistedEvidence(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	var qs []bank.QuestionJSON
	for _, q := range pack.Questions {
		if q.ID == "customer_advance_same_day" || q.ID == "earn_advance_same_day" {
			qs = append(qs, q)
		}
	}
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	auth, err := tutor.NewAuthStore(t.TempDir() + "/auth.json")
	if err != nil {
		t.Fatal(err)
	}
	cache, err := tutor.NewModelCache(t.TempDir() + "/cache.json")
	if err != nil {
		t.Fatal(err)
	}
	m, err := NewModel(Config{DB: db, Catalog: cat, Questions: qs, TotalQuestions: 3, Seed: 101, Clock: mastery.NewMockClock(time.Unix(1, 0)), Tutor: tutor.NewOfflineTutor(), AuthStore: auth, ModelCache: cache})
	if err != nil {
		t.Fatal(err)
	}
	var source bank.QuestionJSON
	for _, q := range qs {
		if q.ID == "customer_advance_same_day" {
			source = q
		}
	}
	inst, err := m.Generator.GenerateInstance(source, 101, map[string]int64{"amount_minor_units": 10000})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.SaveQuestionInstance(inst, m.SessionID); err != nil {
		t.Fatal(err)
	}
	m.CurrentInstance = inst
	m.Session = drill.NewSession(m.SessionID, inst)
	st, _ := m.Session.CurrentStage()
	for i, o := range st.Options {
		if o.ID != st.CorrectOptionID {
			m.SelectedOptionIndex = i
			break
		}
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	st, _ = m.Session.CurrentStage()
	m.Session.SubmitOption(st.CorrectOptionID, time.Unix(2, 0))
	for !m.Session.IsCompleted {
		st, _ = m.Session.CurrentStage()
		m.Session.SubmitOption(st.CorrectOptionID, time.Unix(3, 0))
	}
	if err = m.loadNextQuestion(); err != nil {
		t.Fatal(err)
	}
	if m.CurrentInstance.QuestionID != "earn_advance_same_day" || !m.CurrentInstance.Pedagogy.Remediation || m.CurrentInstance.Parameters["amount_minor_units"] != 10000 || m.CurrentInstance.ScaffoldLevel != domain.ScaffoldFull {
		t.Fatal("wrong, unmarked, unfaded or changed-amount contrast")
	}
	m.State = StateDrill
	m.Height = 200
	if !strings.Contains(m.View(), "Compare this event") || !strings.Contains(m.View(), "retrieval") {
		t.Fatal("learner not informed of guided comparison")
	}
	for !m.Session.IsCompleted {
		st, _ := m.Session.CurrentStage()
		for i, opt := range st.Options {
			if opt.ID == st.CorrectOptionID {
				m.SelectedOptionIndex = i
			}
		}
		m.State = StateDrill
		m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}
	for _, att := range m.Session.Attempts {
		if att.Assistance != domain.AssistanceContrast {
			t.Fatal("comparison stage independently credited")
		}
	}
	saved, err := db.GetQuestionInstance(m.CurrentInstance.InstanceID)
	if err != nil || !saved.Pedagogy.Remediation {
		t.Fatal("comparison designation not replayable")
	}
	atts, err := db.GetAttemptsForInstance(saved.InstanceID)
	if err != nil || len(atts) != 7 {
		t.Fatal("guided evidence not durable")
	}
	stats := mastery.RebuildProjections(atts)
	for _, s := range stats {
		if s.IndependentAttempts != 0 {
			t.Fatal("guided evidence increased mastery")
		}
	}
}

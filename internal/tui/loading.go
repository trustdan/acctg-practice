package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/internal/candidate"
)

type loadingTickMsg struct{}

func loadingTick() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return loadingTickMsg{} })
}

// Terminal glyphs cannot rotate an emoji; animate a rotating indicator beside it.
func (m *Model) loadingGlyph() string {
	return []string{"⏳", "⏳ /", "⏳ —", "⏳ \\"}[m.LoadingFrame%4]
}

type candidateResponseMsg struct {
	ID        int
	Candidate *candidate.CandidateQuestion
	Err       error
}

func (m *Model) cancelCandidate() {
	if m.CandidateCancel != nil {
		m.CandidateCancel()
		m.CandidateCancel = nil
	}
	m.CandidateActive = false
	m.CandidateRequestID++
}

func (m *Model) requestCandidate() (tea.Model, tea.Cmd) {
	if m.CandidateActive {
		return m, nil
	}
	m.cancelTutor()
	if m.State != StateCandidatePreview {
		m.PreviousState = m.State
		m.State = StateCandidatePreview
		m.loadCandidates()
	}
	if m.Tutor == nil || m.Tutor.Name() == "OfflineTutor" {
		m.CandidateNotice = "Connect an LLM with [t] from practice first; [g] generates local variations."
		return m, nil
	}
	family := m.CurrentInstance.FamilyID
	concept := ""
	if len(m.CurrentInstance.Concepts) > 0 {
		concept = m.CurrentInstance.Concepts[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.tutorRequestTimeout())
	m.CandidateCancel = cancel
	m.CandidateActive = true
	m.CandidateRequestID++
	id := m.CandidateRequestID
	m.LoadingFrame = 0
	m.CandidateNotice = "Sending the current transaction family and template examples to your selected LLM."
	gen := candidate.NewProviderCandidateGenerator(m.Tutor, m.Engine, m.Catalog)
	return m, tea.Batch(func() tea.Msg {
		cand, err := gen.Generate(ctx, candidate.GenerateRequest{FamilyID: family, TargetConcept: concept})
		return candidateResponseMsg{ID: id, Candidate: cand, Err: err}
	}, loadingTick())
}

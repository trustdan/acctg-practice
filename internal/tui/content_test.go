package tui_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

type contentProvider struct {
	text    string
	request tutor.Request
}

func (p *contentProvider) Name() string { return "ContentTestProvider" }
func (p *contentProvider) Hint(ctx context.Context, req tutor.Request) (tutor.Response, error) {
	return p.Explain(ctx, req)
}
func (p *contentProvider) Explain(ctx context.Context, req tutor.Request) (tutor.Response, error) {
	p.request = req
	return tutor.Response{Text: p.text, Provider: p.Name(), GeneratedAt: time.Now()}, ctx.Err()
}

func TestLLMContentSavedReviewedAndReplayed(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	seed, err := candidate.NewOfflineCandidateGenerator(m.Engine, m.Catalog).Generate(context.Background(), candidate.GenerateRequest{FamilyID: m.CurrentInstance.FamilyID, Seed: 1})
	if err != nil {
		t.Fatal(err)
	}
	teaching := map[domain.DrillStage]bank.TeachingText{domain.StageBalancedEntry: {Hint: "Which accounts change today?", Explanation: "Review the economic event before selecting the two sides."}}
	payload, _ := json.Marshal(candidate.RawProposal{FamilyID: seed.FamilyID, ScenarioTemplate: seed.ScenarioTemplate, Parameters: seed.Parameters, Concepts: seed.Concepts, Teaching: teaching})
	p := &contentProvider{text: string(payload)}
	m.Tutor = p
	before := len(m.Questions)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if !m.CandidateActive {
		t.Fatal("generation did not start")
	}
	batch := cmd().(tea.BatchMsg)
	// Advance only the animation before the provider completes.
	old := m.View()
	m.Update(batch[1]())
	if m.View() == old {
		t.Fatal("loading indicator did not animate")
	}
	m.Update(batch[0]())
	if m.CandidateActive || len(m.Questions) != before {
		t.Fatal("generation activated content or remained pending")
	}
	if !p.request.CandidateGeneration || !strings.Contains(p.request.ProblemPrompt, "Style example") {
		t.Fatal("missing structured mode/examples")
	}
	rows, err := db.ListCandidates(storage.CandidateFilter{})
	if err != nil || len(rows) < 1 {
		t.Fatal("candidate not durable", err)
	}
	if rows[0].Teaching[domain.StageBalancedEntry].Hint != teaching[domain.StageBalancedEntry].Hint {
		t.Fatal("generated hints lost in persistence")
	}
	sendKey(m, "a")
	if len(m.Questions) != before+1 {
		t.Fatalf("approval failed: %s", m.CandidateNotice)
	}
	published, err := db.GetPublishedQuestion(m.Candidates[0].ID)
	if err != nil || published.Teaching[domain.StageBalancedEntry].Hint != teaching[domain.StageBalancedEntry].Hint {
		t.Fatal("published teaching did not survive reload", err)
	}
}

func TestCanceledLLMContentCannotBeSaved(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = &contentProvider{text: `{}`}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	batch := cmd().(tea.BatchMsg)
	sendSpecialKey(m, tea.KeyEsc)
	before := len(m.Candidates)
	m.Update(batch[0]())
	if m.CandidateActive || len(m.Candidates) != before {
		t.Fatal("late canceled result was accepted")
	}
}

func TestReviewedTeachingChangesProseOnly(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	seed, _ := candidate.NewOfflineCandidateGenerator(m.Engine, m.Catalog).Generate(context.Background(), candidate.GenerateRequest{FamilyID: m.CurrentInstance.FamilyID, Seed: 1})
	seed.Teaching = map[domain.DrillStage]bank.TeachingText{domain.StageBalancedEntry: {Hint: "What changes today?", Explanation: "New reviewed explanation."}}
	q, _, err := candidate.ApproveCandidate(seed, "test reviewer", "reviewed", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	inst, err := m.Generator.GenerateInstance(*q, 1, map[string]int64{"amount_minor_units": q.Parameters["amount_minor_units"][0]})
	if err != nil {
		t.Fatal(err)
	}
	answer := inst.StageAnswers[domain.StageBalancedEntry]
	if answer.CausalHint != "What changes today?" || answer.Explanation != "New reviewed explanation." {
		t.Fatal("reviewed prose was lost")
	}
	if answer.CorrectOptionID == "" {
		t.Fatal("canonical answer lost")
	}
}

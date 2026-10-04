package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

type explanationSavedMsg struct{ Err error }

func (m *Model) captureExplanation(resp tutor.Response) {
	if m.CurrentInstance == nil || m.Session == nil {
		return
	}
	stage, err := m.Session.CurrentStage()
	if err != nil {
		return
	}
	m.PendingExplanation = &storage.SavedExplanation{
		ID: uuid.NewString(), InstanceID: m.CurrentInstance.InstanceID,
		QuestionID: m.CurrentInstance.QuestionID, QuestionVersion: m.CurrentInstance.Version,
		Stage: string(stage.Stage), Scenario: m.CurrentInstance.PromptText, StagePrompt: stage.Prompt,
		Explanation: resp.Text, Provider: resp.Provider, GeneratedAt: resp.GeneratedAt,
	}
}

func (m *Model) unsavedExplanationVisible() bool {
	return m.PendingExplanation != nil && m.ShowHint && m.TutorKind == "explain" &&
		(m.State == StateDrill || m.State == StateFeedback)
}

func explanationReadingKey(key string) bool {
	switch key {
	case "u", "d", "pgup", "pgdown", "up", "down", "j", "k", "g", "G", "home", "end":
		return true
	}
	return false
}

func (m *Model) updateExplanationSave(key string) (tea.Model, tea.Cmd) {
	if m.ExplanationSaving {
		return m, nil
	}
	switch key {
	case "y", "Y":
		if m.DB == nil {
			m.ExplanationSaveError = "Database unavailable. Retry or choose No to continue."
			return m, nil
		}
		entry := *m.PendingExplanation
		entry.SavedAt = m.Clock.Now()
		m.ExplanationSaving = true
		m.ExplanationSaveError = ""
		db := m.DB
		return m, func() tea.Msg { return explanationSavedMsg{Err: db.SaveExplanation(entry)} }
	case "n", "N":
		return m.resumeAfterExplanation()
	case "esc":
		m.ExplanationSavePrompt = false
		m.ExplanationSaveError = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) resumeAfterExplanation() (tea.Model, tea.Cmd) {
	key := m.ExplanationLeaveKey
	m.ExplanationSavePrompt = false
	m.PendingExplanation = nil
	m.ShowHint = false
	m.PageScroll = 0
	return m.Update(key)
}

func (m *Model) renderExplanationSave() string {
	text := "Would you like to save this explanation in the database?\n\n[y] Yes, save   [n] No, continue   [Esc] Keep reading"
	if m.ExplanationSaving {
		text = "Saving explanation..."
	}
	if m.ExplanationSaveError != "" {
		text += "\n\nCould not save: " + m.ExplanationSaveError
	}
	text += "\n\nSaved explanations are personal notes; they do not change grades or approved questions."
	return m.renderPageViewport(m.Styles.ExplanationBox.Render(text))
}

func (m *Model) openSavedExplanations() (tea.Model, tea.Cmd) {
	if m.DB == nil {
		m.TutorError = "Database unavailable."
		return m, nil
	}
	entries, err := m.DB.ListSavedExplanations()
	if err != nil {
		m.TutorError = err.Error()
		return m, nil
	}
	m.SavedExplanations = entries
	m.SavedExplanationIndex = 0
	m.PreviousState = m.State
	m.State = StateSavedExplanations
	return m, nil
}

func (m *Model) updateSavedExplanations(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "V":
		m.State = m.PreviousState
	case "q":
		m.closeSession()
		m.State = StateQuitting
		return m, tea.Quit
	case "right", "n":
		if len(m.SavedExplanations) > 0 {
			m.SavedExplanationIndex = (m.SavedExplanationIndex + 1) % len(m.SavedExplanations)
		}
	case "left", "p":
		if len(m.SavedExplanations) > 0 {
			m.SavedExplanationIndex = (m.SavedExplanationIndex + len(m.SavedExplanations) - 1) % len(m.SavedExplanations)
		}
	}
	return m, nil
}

func (m *Model) renderSavedExplanations(width int) string {
	parts := []string{m.Styles.CardTitle.Render("SAVED EXPLANATIONS"), "Personal AI notes — advisory, not reviewed answer keys."}
	if len(m.SavedExplanations) == 0 {
		parts = append(parts, "No saved explanations yet. Request an explanation with [e], then choose Yes when leaving it.")
	} else {
		e := m.SavedExplanations[m.SavedExplanationIndex]
		parts = append(parts, fmt.Sprintf("%d / %d | %s | %s\nSaved %s", m.SavedExplanationIndex+1, len(m.SavedExplanations), e.Provider, e.Stage, e.SavedAt.Local().Format("2006-01-02 15:04")),
			m.Styles.ScenarioBox.Width(width).Render(e.Scenario+"\n\n"+e.StagePrompt), m.Styles.ExplanationBox.Width(width).Render(RenderMarkdown(e.Explanation, width-4)))
	}
	return strings.Join(append(parts, "[n/p or Left/Right] Browse  [u/d] Scroll  [Esc/V] Return"), "\n\n")
}

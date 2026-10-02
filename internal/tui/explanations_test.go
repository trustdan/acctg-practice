package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/internal/tui"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

func finishTitleForTest(t *testing.T, m *tui.Model) {
	t.Helper()
	if m.State != tui.StateTitle {
		t.Fatalf("expected post-game title, got %v", m.State)
	}
	m.TitleFrame = 89
	m.Update(m.Init()())
}

func TestFinishedGameScoreEntryAutomaticallyStartsTitle(t *testing.T) {
	for _, cancel := range []bool{false, true} {
		m, db := setupTestTUIWithIntro(t)
		m.Intro.GameOver = true
		m.Intro.Score = 6000
		sendSpecialKey(m, tea.KeyEnter)
		if !m.Intro.InitialsEntryActive {
			t.Fatal("missing high-score entry")
		}
		if cancel {
			sendSpecialKey(m, tea.KeyEsc)
		} else {
			sendKey(m, "A")
			sendKey(m, "B")
			sendKey(m, "C")
			sendSpecialKey(m, tea.KeyEnter)
		}
		finishTitleForTest(t, m)
		if m.State != tui.StateDrill {
			t.Fatal("game-over exit did not return to practice")
		}
		db.Close()
	}
}

func generateExplanation(t *testing.T, m *tui.Model) {
	t.Helper()
	m.Tutor = &contentProvider{text: "Cash collection reduces a receivable; earning happened earlier."}
	sendKey(m, "e")
	if m.PendingExplanation == nil {
		t.Fatal("missing pending LLM explanation")
	}
}

func TestExplanationSaveDefersNavigationAndKeepsEvidenceSeparate(t *testing.T) {
	for _, key := range []string{"esc", "t", "q", "1", "e", "n", "A", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			m, db := setupTestTUI(t)
			defer db.Close()
			generateExplanation(t, m)
			beforeAttempts, _ := db.GetAllAttempts()
			beforeQuestions := len(m.Questions)
			pending := *m.PendingExplanation
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
			if key == "esc" {
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			}
			if key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}
			m.Update(msg)
			if !m.ExplanationSavePrompt || m.State != tui.StateDrill {
				t.Fatal("action ran before save choice")
			}
			if !strings.Contains(m.View(), "Would you like to save this explanation in the database?") {
				t.Fatal("save question hidden")
			}
			rows, _ := db.ListSavedExplanations()
			if len(rows) != 0 {
				t.Fatal("saved before consent")
			}
			sendKey(m, "y")
			rows, err := db.ListSavedExplanations()
			if err != nil || len(rows) != 1 || rows[0].Explanation != pending.Explanation || rows[0].QuestionID != pending.QuestionID || rows[0].Stage != pending.Stage || rows[0].Provider != pending.Provider {
				t.Fatalf("wrong saved note: %+v %v", rows, err)
			}
			if m.ExplanationSavePrompt {
				t.Fatal("save prompt did not close")
			}
			afterAttempts, _ := db.GetAllAttempts()
			if len(afterAttempts) != len(beforeAttempts) || len(m.Questions) != beforeQuestions {
				t.Fatal("saving changed grading evidence or bank")
			}
			if key == "t" && m.State != tui.StateTutorConfig {
				t.Fatal("original hotkey was lost")
			}
			if (key == "q" || key == "ctrl+c") && m.State != tui.StateQuitting {
				t.Fatal("original quit was lost")
			}
		})
	}
}

func TestExplanationDeclineCancelScrollAndMouse(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	generateExplanation(t, m)
	sendKey(m, "d")
	if m.ExplanationSavePrompt {
		t.Fatal("scrolling opened save prompt")
	}
	m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	if !m.ExplanationSavePrompt {
		t.Fatal("click did not open save prompt")
	}
	sendSpecialKey(m, tea.KeyEsc)
	if m.ExplanationSavePrompt || !m.ShowHint || m.PendingExplanation == nil {
		t.Fatal("cancel lost explanation")
	}
	sendKey(m, "t")
	sendKey(m, "n")
	rows, _ := db.ListSavedExplanations()
	if len(rows) != 0 || m.State != tui.StateTutorConfig || m.ShowHint {
		t.Fatal("decline saved content or lost navigation")
	}
}

func TestExplanationSaveFailureKeepsTextAndAllowsRetry(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	generateExplanation(t, m)
	sendSpecialKey(m, tea.KeyEsc)
	m.DB = nil
	sendKey(m, "y")
	if !m.ExplanationSavePrompt || m.PendingExplanation == nil || m.ExplanationSaveError == "" {
		t.Fatal("save failure lost text")
	}
	m.DB = db
	sendKey(m, "y")
	sendKey(m, "V")
	if m.State != tui.StateSavedExplanations || !strings.Contains(m.View(), "earning happened earlier") {
		t.Fatal("saved explanation unavailable")
	}
}

func TestOfflineExplanationDoesNotPromptForSaving(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = tutor.NewOfflineTutor()
	sendKey(m, "e")
	sendSpecialKey(m, tea.KeyEsc)
	if m.ExplanationSavePrompt || m.PendingExplanation != nil {
		t.Fatal("offline prose treated as LLM explanation")
	}
}

func TestFailedDatabaseWriteKeepsSavePrompt(t *testing.T) {
	m, db := setupTestTUI(t)
	generateExplanation(t, m)
	sendSpecialKey(m, tea.KeyEsc)
	db.Close()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd == nil || !m.ExplanationSaving {
		t.Fatal("save was not asynchronous")
	}
	m.Update(cmd())
	if !m.ExplanationSavePrompt || m.ExplanationSaving || m.ExplanationSaveError == "" || m.PendingExplanation == nil {
		t.Fatal("failed write lost explanation or navigated away")
	}
	sendKey(m, "n")
	if m.ExplanationSavePrompt || m.ShowHint {
		t.Fatal("could not decline after write failure")
	}
}

func TestFallbackExplanationDoesNotPromptForSaving(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = tutor.NewFallbackTutor(tutor.NewSimulatedProvider(tutor.SimulatedConfig{FailWith: tutor.ErrProviderUnavailable}), tutor.NewOfflineTutor(), m.TutorTimeout)
	sendKey(m, "e")
	if m.TutorResponse == nil || !m.TutorResponse.Fallback {
		t.Fatal("expected offline fallback")
	}
	sendKey(m, "t")
	if m.ExplanationSavePrompt || m.State != tui.StateTutorConfig {
		t.Fatal("fallback incorrectly asked to save an LLM explanation")
	}
}

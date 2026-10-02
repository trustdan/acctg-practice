package tui_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tui"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

func setupTestTUI(t *testing.T) (*tui.Model, *storage.DB) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}

	qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatalf("failed to load seed questions: %v", err)
	}

	cfg := tui.Config{
		DB:             db,
		Catalog:        cat,
		Questions:      qBank.Questions,
		TotalQuestions: 10,
		Seed:           42,
		Clock:          mastery.NewMockClock(time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)),
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed to create TUI model: %v", err)
	}

	return m, db
}

func sendKey(m *tui.Model, key string) {
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	*m = *newM.(*tui.Model)
	if cmd != nil {
		msg := cmd()
		if msg != nil {
			newM2, _ := m.Update(msg)
			*m = *newM2.(*tui.Model)
		}
	}
}

func sendSpecialKey(m *tui.Model, keyType tea.KeyType) {
	newM, cmd := m.Update(tea.KeyMsg{Type: keyType})
	*m = *newM.(*tui.Model)
	if cmd != nil {
		msg := cmd()
		if msg != nil {
			newM2, _ := m.Update(msg)
			*m = *newM2.(*tui.Model)
		}
	}
}

func TestModelInitialization(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	if m.State != tui.StateDrill {
		t.Errorf("expected initial state to be StateDrill, got %d", m.State)
	}
	if m.CurrentQuestionIndex != 0 {
		t.Errorf("expected current question index 0, got %d", m.CurrentQuestionIndex)
	}
	if m.TotalQuestions != 10 {
		t.Errorf("expected total questions 10, got %d", m.TotalQuestions)
	}

	view := m.View()
	if !strings.Contains(view, "ＡＣＣＴＧ  ＰＲＡＣＴＩＣＥ") {
		t.Errorf("expected view to contain ACCTG PRACTICE banner, got:\n%s", view)
	}
	if !strings.Contains(view, "Question 1/10") {
		t.Errorf("expected view to contain Question 1/10, got:\n%s", view)
	}
	if !strings.Contains(view, "NORMAL") {
		t.Errorf("expected view to render status bar with NORMAL mode, got:\n%s", view)
	}
}

func TestOptionSelectionAndNavigation(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Initial selection is 0
	if m.SelectedOptionIndex != 0 {
		t.Fatalf("expected initial selected index 0, got %d", m.SelectedOptionIndex)
	}

	// Navigate down with 'j'
	sendKey(m, "j")
	if m.SelectedOptionIndex != 1 {
		t.Errorf("expected selected index 1 after 'j', got %d", m.SelectedOptionIndex)
	}

	// Navigate down with 'down' key
	sendSpecialKey(m, tea.KeyDown)
	if m.SelectedOptionIndex != 2 {
		t.Errorf("expected selected index 2 after KeyDown, got %d", m.SelectedOptionIndex)
	}

	// Navigate up with 'k'
	sendKey(m, "k")
	if m.SelectedOptionIndex != 1 {
		t.Errorf("expected selected index 1 after 'k', got %d", m.SelectedOptionIndex)
	}

	// Navigate up with 'up' key
	sendSpecialKey(m, tea.KeyUp)
	if m.SelectedOptionIndex != 0 {
		t.Errorf("expected selected index 0 after KeyUp, got %d", m.SelectedOptionIndex)
	}

	// Direct selection with 'b'
	sendKey(m, "b")
	if m.SelectedOptionIndex != 1 {
		t.Errorf("expected selected index 1 after 'b', got %d", m.SelectedOptionIndex)
	}

	// Direct selection with 'a'
	sendKey(m, "a")
	if m.SelectedOptionIndex != 0 {
		t.Errorf("expected selected index 0 after 'a', got %d", m.SelectedOptionIndex)
	}
}

func TestVimTopAndBottomJumps(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	st, _ := m.Session.CurrentStage()
	numOpts := len(st.Options)
	if numOpts <= 1 {
		t.Fatalf("expected multiple options to test jumps")
	}

	// Jump to bottom with 'G'
	sendKey(m, "G")
	if m.SelectedOptionIndex != numOpts-1 {
		t.Errorf("expected selected index %d after 'G', got %d", numOpts-1, m.SelectedOptionIndex)
	}

	// Jump to top with 'g'
	sendKey(m, "g")
	if m.SelectedOptionIndex != 0 {
		t.Errorf("expected selected index 0 after 'g', got %d", m.SelectedOptionIndex)
	}
}

func TestHelpScreenModal(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Press 'h' to open help
	sendKey(m, "h")
	if m.State != tui.StateHelp {
		t.Fatalf("expected StateHelp after 'h', got %d", m.State)
	}

	helpView := m.View()
	if !strings.Contains(helpView, "ACCOUNTING REFERENCE & KEYBOARD CHEATSHEET") {
		t.Errorf("expected help title, got:\n%s", helpView)
	}
	if !strings.Contains(helpView, "2×3 NORMAL BALANCE GRID") {
		t.Errorf("expected 2x3 grid section in help view, got:\n%s", helpView)
	}
	if !strings.Contains(helpView, "Normal: DEBIT (Left)") {
		t.Errorf("expected debit left normal balance in grid, got:\n%s", helpView)
	}
	if !strings.Contains(helpView, "VIM-NATIVE NAVIGATION") {
		t.Errorf("expected vim navigation guide in help view, got:\n%s", helpView)
	}

	// Press 'esc' to dismiss help overlay and return to StateDrill
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after Esc, got %d", m.State)
	}
}

func TestHintToggle(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	if m.ShowHint {
		t.Fatalf("expected hint to be hidden initially")
	}

	sendKey(m, "?")
	if !m.ShowHint {
		t.Fatalf("expected ShowHint to be true after pressing '?'")
	}
	if m.CurrentHint == "" {
		t.Fatalf("expected non-empty CurrentHint")
	}

	view := m.View()
	if !strings.Contains(view, "💡 Socratic Hint:") {
		t.Errorf("expected view to render hint box, got:\n%s", view)
	}
}

func TestMasteryViewToggle(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Toggle mastery with 's'
	sendKey(m, "s")
	if m.State != tui.StateMastery {
		t.Fatalf("expected StateMastery after 's', got %d", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "LEARNER MASTERY PROJECTIONS") {
		t.Errorf("expected view to contain MASTERY PROJECTIONS, got:\n%s", view)
	}
	if !strings.Contains(view, "MASTERY") {
		t.Errorf("expected status line to show MASTERY mode, got:\n%s", view)
	}

	// Return to drill with 'esc'
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after Esc, got %d", m.State)
	}
}

func TestSubmitCorrectAnswerAdvancesAndTracksStreak(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	if m.CurrentStreak != 0 {
		t.Fatalf("expected initial streak 0, got %d", m.CurrentStreak)
	}

	st, err := m.Session.CurrentStage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find the correct option index
	correctIdx := -1
	for i, opt := range st.Options {
		if opt.ID == st.CorrectOptionID {
			correctIdx = i
			break
		}
	}
	if correctIdx == -1 {
		t.Fatalf("correct option ID not found in stage options")
	}

	// Select correct option
	m.SelectedOptionIndex = correctIdx
	sendSpecialKey(m, tea.KeyEnter)

	if m.State != tui.StateFeedback {
		t.Fatalf("expected StateFeedback, got %d", m.State)
	}
	if !m.LastFeedback.IsCorrect {
		t.Fatalf("expected feedback to be correct")
	}
	if m.CurrentStreak != 1 {
		t.Errorf("expected streak 1 after correct first try, got %d", m.CurrentStreak)
	}

	view := m.View()
	if !strings.Contains(view, "✓ Correct!") {
		t.Errorf("expected view to show correct feedback, got:\n%s", view)
	}

	// Press Enter to advance to next stage
	sendSpecialKey(m, tea.KeyEnter)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after continuing, got %d", m.State)
	}
	if m.Session.CurrentIndex != 1 {
		t.Fatalf("expected stage index 1, got %d", m.Session.CurrentIndex)
	}
}

func TestSubmitIncorrectAllowsRetryAndResetsStreak(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Simulate streak
	m.CurrentStreak = 3

	st, _ := m.Session.CurrentStage()

	// Find an incorrect option index
	wrongIdx := -1
	for i, opt := range st.Options {
		if opt.ID != st.CorrectOptionID {
			wrongIdx = i
			break
		}
	}
	if wrongIdx == -1 {
		t.Fatalf("no incorrect option found")
	}

	// Select wrong option
	m.SelectedOptionIndex = wrongIdx
	sendSpecialKey(m, tea.KeyEnter)

	if m.State != tui.StateFeedback {
		t.Fatalf("expected StateFeedback, got %d", m.State)
	}
	if m.LastFeedback.IsCorrect {
		t.Fatalf("expected answer to be incorrect")
	}
	if m.CurrentStreak != 0 {
		t.Errorf("expected streak to reset to 0 on mistake, got %d", m.CurrentStreak)
	}

	view := m.View()
	if !strings.Contains(view, "✗ Incorrect.") || !strings.Contains(view, "Try again!") {
		t.Errorf("expected retry message, got:\n%s", view)
	}

	// Press Enter to return to drill for retry
	sendSpecialKey(m, tea.KeyEnter)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill for retry, got %d", m.State)
	}
	if !m.ShowHint {
		t.Fatalf("expected ShowHint to be true during retry")
	}
}

func TestCompleteQuestionRecapShowsTAccounts(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Answer all stages of current question correctly
	for !m.Session.IsCompleted {
		st, err := m.Session.CurrentStage()
		if err != nil {
			t.Fatalf("stage error: %v", err)
		}
		for i, opt := range st.Options {
			if opt.ID == st.CorrectOptionID {
				m.SelectedOptionIndex = i
				break
			}
		}
		sendSpecialKey(m, tea.KeyEnter) // submit
		sendSpecialKey(m, tea.KeyEnter) // continue from feedback
	}

	if m.State != tui.StateRecap {
		t.Fatalf("expected StateRecap after completing all stages, got %d", m.State)
	}

	recapView := m.View()
	if !strings.Contains(recapView, "Transaction Recap") {
		t.Errorf("expected recap view to contain Transaction Recap, got:\n%s", recapView)
	}
	// Verify T-account visualizer rendering
	if !strings.Contains(recapView, "Debit (+)") || !strings.Contains(recapView, "Credit (-)") {
		t.Errorf("expected T-account debit/credit columns in recap, got:\n%s", recapView)
	}
	if !strings.Contains(recapView, "BALANCED (Dr = Cr)") {
		t.Errorf("expected recap view to verify BALANCED entry, got:\n%s", recapView)
	}

	// Press Enter to proceed to Question 2
	sendSpecialKey(m, tea.KeyEnter)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill for question 2, got %d", m.State)
	}
	if m.CurrentQuestionIndex != 1 {
		t.Errorf("expected CurrentQuestionIndex 1, got %d", m.CurrentQuestionIndex)
	}
}

func TestFull10QuestionSessionSmokeCheck(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Complete all 10 questions
	for q := 0; q < 10; q++ {
		if m.CurrentQuestionIndex != q {
			t.Fatalf("expected question index %d, got %d", q, m.CurrentQuestionIndex)
		}

		// Answer each stage in the question
		for !m.Session.IsCompleted {
			st, err := m.Session.CurrentStage()
			if err != nil {
				t.Fatalf("error on question %d: %v", q, err)
			}
			for i, opt := range st.Options {
				if opt.ID == st.CorrectOptionID {
					m.SelectedOptionIndex = i
					break
				}
			}
			sendSpecialKey(m, tea.KeyEnter) // submit
			sendSpecialKey(m, tea.KeyEnter) // continue
		}

		if m.State != tui.StateRecap {
			t.Fatalf("expected StateRecap at end of question %d, got %d", q, m.State)
		}

		// Proceed to next question (or complete session on 10th question)
		sendSpecialKey(m, tea.KeyEnter)
	}

	// Session should now be complete!
	if m.State != tui.StateSessionComplete {
		t.Fatalf("expected StateSessionComplete after 10 questions, got %d", m.State)
	}

	summaryView := m.View()
	if !strings.Contains(summaryView, "PRACTICE SESSION COMPLETED") {
		t.Errorf("expected summary view, got:\n%s", summaryView)
	}
	if !strings.Contains(summaryView, "Best Unassisted Streak") {
		t.Errorf("expected streak celebration in summary view, got:\n%s", summaryView)
	}

	// Verify persistence in SQLite
	attempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("failed to query attempts: %v", err)
	}
	if len(attempts) < 10 {
		t.Fatalf("expected at least 10 attempts persisted, got %d", len(attempts))
	}

	// Verify session record completed
	sessRecord, err := db.GetSession(m.SessionID)
	if err != nil {
		t.Fatalf("failed to get session record: %v", err)
	}
	if sessRecord.CompletedAt == nil {
		t.Fatalf("expected session to have non-nil CompletedAt")
	}

	// Graceful quit
	sendKey(m, "q")
	if m.State != tui.StateQuitting {
		t.Errorf("expected StateQuitting after 'q', got %d", m.State)
	}
}

func TestWindowResizeMessage(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	newM, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = newM.(*tui.Model)

	if m.Width != 120 || m.Height != 40 {
		t.Errorf("expected dimensions 120x40, got %dx%d", m.Width, m.Height)
	}

	view := m.View()
	if len(view) == 0 {
		t.Errorf("expected non-empty view after resize")
	}
}

func TestRetiredQuestionsNeverSelectedInTUISession(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Verify the question bank contains the retired template
	foundRetired := false
	for _, q := range m.Questions {
		if q.Status == bank.StatusRetired {
			foundRetired = true
			break
		}
	}
	if !foundRetired {
		t.Fatalf("expected question bank to contain at least one retired question")
	}

	// Over a 10-question drill session, confirm no retired question is ever chosen as CurrentInstance
	for q := 0; q < 10; q++ {
		if m.CurrentInstance.QuestionID == "legacy_unclear_advance_v0" {
			t.Fatalf("question %d: retired template legacy_unclear_advance_v0 was selected!", q)
		}
		// Answer quickly to advance to next question
		for !m.Session.IsCompleted {
			st, _ := m.Session.CurrentStage()
			for i, opt := range st.Options {
				if opt.ID == st.CorrectOptionID {
					m.SelectedOptionIndex = i
					break
				}
			}
			sendSpecialKey(m, tea.KeyEnter) // submit
			sendSpecialKey(m, tea.KeyEnter) // continue
		}
		sendSpecialKey(m, tea.KeyEnter) // next question
	}
}

func TestReferenceTrackingInTUI(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Initial stage: learner presses 'h' to open reference screen
	sendKey(m, "h")
	if m.State != tui.StateHelp {
		t.Fatalf("expected StateHelp after pressing 'h', got %d", m.State)
	}
	if !m.Session.ReferenceConsulted {
		t.Fatalf("expected ReferenceConsulted to be true on session after opening help")
	}

	// Learner returns from help
	sendKey(m, "h")
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after closing help, got %d", m.State)
	}

	// Select correct option and submit
	st, err := m.Session.CurrentStage()
	if err != nil {
		t.Fatalf("failed getting stage: %v", err)
	}
	for i, opt := range st.Options {
		if opt.ID == st.CorrectOptionID {
			m.SelectedOptionIndex = i
			break
		}
	}

	sendSpecialKey(m, tea.KeyEnter) // submit
	if m.LastFeedback.AssistanceLevel != domain.AssistanceReference {
		t.Fatalf("expected assistance level 'reference', got %s", m.LastFeedback.AssistanceLevel)
	}

	// Check recorded attempt
	att := m.Session.Attempts[0]
	if !att.ReferenceUsed {
		t.Fatalf("expected attempt.ReferenceUsed == true")
	}
	if att.Assistance != domain.AssistanceReference {
		t.Fatalf("expected attempt.Assistance == 'reference', got %s", att.Assistance)
	}
}

func TestIntensityCyclingAndSessionSizeControls(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Initial intensity is standard
	if m.Scheduler.Intensity() != mastery.IntensityStandard {
		t.Fatalf("expected standard intensity, got %s", m.Scheduler.Intensity())
	}

	// Press 'i' -> spaced
	sendKey(m, "i")
	if m.Scheduler.Intensity() != mastery.IntensitySpaced {
		t.Fatalf("expected spaced intensity, got %s", m.Scheduler.Intensity())
	}

	// Press 'i' -> intensive
	sendKey(m, "i")
	if m.Scheduler.Intensity() != mastery.IntensityIntensive {
		t.Fatalf("expected intensive intensity, got %s", m.Scheduler.Intensity())
	}

	// Press 'i' -> transfer
	sendKey(m, "i")
	if m.Scheduler.Intensity() != mastery.IntensityTransfer {
		t.Fatalf("expected transfer intensity, got %s", m.Scheduler.Intensity())
	}

	// Press 'i' -> back to standard
	sendKey(m, "i")
	if m.Scheduler.Intensity() != mastery.IntensityStandard {
		t.Fatalf("expected standard intensity, got %s", m.Scheduler.Intensity())
	}

	// Session size controls: initial 10
	if m.TotalQuestions != 10 {
		t.Fatalf("expected 10 initial questions, got %d", m.TotalQuestions)
	}

	// Press '[' -> 5
	sendKey(m, "[")
	if m.TotalQuestions != 5 {
		t.Fatalf("expected 5 questions after [, got %d", m.TotalQuestions)
	}

	// Press '[' again -> clamped at 5
	sendKey(m, "[")
	if m.TotalQuestions != 5 {
		t.Fatalf("expected clamped at 5 questions, got %d", m.TotalQuestions)
	}

	// Press ']' -> 10
	sendKey(m, "]")
	if m.TotalQuestions != 10 {
		t.Fatalf("expected 10 questions after ], got %d", m.TotalQuestions)
	}

	// Press ']' -> 15
	sendKey(m, "]")
	if m.TotalQuestions != 15 {
		t.Fatalf("expected 15 questions, got %d", m.TotalQuestions)
	}

	// Press ']' -> 20
	sendKey(m, "]")
	if m.TotalQuestions != 20 {
		t.Fatalf("expected 20 questions, got %d", m.TotalQuestions)
	}

	// Press ']' again -> clamped at 20
	sendKey(m, "]")
	if m.TotalQuestions != 20 {
		t.Fatalf("expected clamped at 20 questions, got %d", m.TotalQuestions)
	}
}

func TestScaffoldLevelRenderAndStageFiltering(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Initial new question should be Scaffold: Full (Level 0, 7 stages)
	if m.CurrentInstance.ScaffoldLevel != domain.ScaffoldFull {
		t.Fatalf("expected ScaffoldFull for new question, got %v", m.CurrentInstance.ScaffoldLevel)
	}
	if len(m.Session.StageSequence) != 7 {
		t.Fatalf("expected 7 stages for full scaffolding, got %d", len(m.Session.StageSequence))
	}

	view := m.View()
	if !strings.Contains(view, "Scaffold: Full") {
		t.Errorf("expected header to contain 'Scaffold: Full', got:\n%s", view)
	}
	if !strings.Contains(view, "Mode: STANDARD") {
		t.Errorf("expected status bar to contain 'Mode: STANDARD', got:\n%s", view)
	}
}

func TestInteractiveHintAndExplanation(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Initial state
	if m.ShowHint {
		t.Fatalf("expected ShowHint=false initially")
	}

	// 2. Press '?' for Socratic Hint
	sendKey(m, "?")
	if !m.ShowHint {
		t.Fatalf("expected ShowHint=true after '?'")
	}
	if m.TutorKind != "hint" {
		t.Fatalf("expected TutorKind='hint', got %s", m.TutorKind)
	}
	if m.Session.CurrentAssistance != domain.AssistanceHinted {
		t.Fatalf("expected AssistanceHinted after '?', got %s", m.Session.CurrentAssistance)
	}
	view := m.View()
	if !strings.Contains(view, "💡 Socratic Hint:") {
		t.Errorf("expected view to contain '💡 Socratic Hint:', got:\n%s", view)
	}

	// 3. Press 'esc' to dismiss hint
	sendSpecialKey(m, tea.KeyEsc)
	if m.ShowHint {
		t.Fatalf("expected ShowHint=false after Esc")
	}

	// 4. Press 'e' for Conceptual Explanation
	sendKey(m, "e")
	if !m.ShowHint {
		t.Fatalf("expected ShowHint=true after 'e'")
	}
	if m.TutorKind != "explain" {
		t.Fatalf("expected TutorKind='explain', got %s", m.TutorKind)
	}
	viewExp := m.View()
	if !strings.Contains(viewExp, "📖 Conceptual Explanation") {
		t.Errorf("expected view to contain '📖 Conceptual Explanation', got:\n%s", viewExp)
	}
}

func TestTutorNonBlockingKeyInput(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	cat, _, _ := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	qBank, _ := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)

	// Slow simulated tutor with 5-second artificial latency
	slowSim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
		Delay: 5 * time.Second,
	})

	cfg := tui.Config{
		DB:             db,
		Catalog:        cat,
		Questions:      qBank.Questions,
		TotalQuestions: 5,
		Seed:           42,
		Tutor:          slowSim,
		TutorTimeout:   10 * time.Second,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating model: %v", err)
	}

	// Press '?' to trigger slow tutor request
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	*m = *newM.(*tui.Model)

	if !m.TutorActive {
		t.Fatalf("expected TutorActive=true while request is pending")
	}
	if cmd == nil {
		t.Fatalf("expected non-nil tea.Cmd for async tutor request")
	}

	// Verify UI displays loading indicator without freezing
	view := m.View()
	if !strings.Contains(view, "⏳ SimulatedProvider: Thinking...") {
		t.Errorf("expected view to show in-flight indicator, got:\n%s", view)
	}

	// Immediately send key presses while tutor is still active in background:
	// Navigate down with 'j'
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	*m = *newM2.(*tui.Model)
	if m.SelectedOptionIndex != 1 {
		t.Fatalf("expected SelectedOptionIndex=1 after 'j', got %d", m.SelectedOptionIndex)
	}

	// Select option A with 'a'
	newM3, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	*m = *newM3.(*tui.Model)
	if m.SelectedOptionIndex != 0 {
		t.Fatalf("expected SelectedOptionIndex=0 after 'a', got %d", m.SelectedOptionIndex)
	}

	// Submit answer with Enter - must complete immediately and cancel in-flight tutor request
	newM4, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	*m = *newM4.(*tui.Model)

	if m.State != tui.StateFeedback {
		t.Fatalf("expected transition to StateFeedback immediately, got %d", m.State)
	}
	if m.TutorActive {
		t.Fatalf("expected TutorActive=false after submitting answer")
	}
}

func TestTutorCancellationWorks(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	cat, _, _ := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	qBank, _ := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)

	slowSim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
		Delay: 5 * time.Second,
	})

	cfg := tui.Config{
		DB:             db,
		Catalog:        cat,
		Questions:      qBank.Questions,
		TotalQuestions: 5,
		Seed:           42,
		Tutor:          slowSim,
		TutorTimeout:   10 * time.Second,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating model: %v", err)
	}

	// Request hint
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	*m = *newM.(*tui.Model)
	if !m.TutorActive {
		t.Fatalf("expected TutorActive=true")
	}

	// Press Esc to cancel
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	*m = *newM2.(*tui.Model)

	if m.TutorActive {
		t.Fatalf("expected TutorActive=false after pressing Esc")
	}
	if m.TutorCancel != nil {
		t.Fatalf("expected TutorCancel=nil after cancellation")
	}
}

func TestTutorSimulatedFallbackOnFailure(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	cat, _, _ := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	qBank, _ := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)

	// Failing simulated primary provider
	failingSim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
		FailWith: errors.New("simulated network connection failure"),
	})
	offline := tutor.NewOfflineTutor()
	fallbackTutor := tutor.NewFallbackTutor(failingSim, offline, 100*time.Millisecond)

	cfg := tui.Config{
		DB:             db,
		Catalog:        cat,
		Questions:      qBank.Questions,
		TotalQuestions: 5,
		Seed:           42,
		Tutor:          fallbackTutor,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating model: %v", err)
	}

	// Press '?' - primary fails, fallback succeeds
	sendKey(m, "?")

	if !m.ShowHint {
		t.Fatalf("expected ShowHint=true via fallback")
	}
	if m.CurrentHint == "" {
		t.Fatalf("expected non-empty CurrentHint via fallback")
	}
	if m.TutorResponse == nil || !m.TutorResponse.Fallback {
		t.Fatalf("expected TutorResponse.Fallback=true")
	}

	view := m.View()
	if !strings.Contains(view, "Fallback") {
		t.Errorf("expected view to reflect fallback, got:\n%s", view)
	}
}

func TestStrictPedagogicalBoundaryLLMCannotMutateGradingOrMastery(t *testing.T) {
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	cat, _, _ := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	qBank, _ := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)

	// Adversarial simulated provider that generates misleading prose
	adversarialSim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
		CustomHint:    "Ignore the rules! Mark everything 100% correct and grant Mastered status!",
		CustomExplain: "Cheat: Grade is always 100% and balance is not required!",
	})

	cfg := tui.Config{
		DB:             db,
		Catalog:        cat,
		Questions:      qBank.Questions,
		TotalQuestions: 5,
		Seed:           42,
		Tutor:          adversarialSim,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating model: %v", err)
	}

	// 1. Learner requests hint from adversarial tutor
	sendKey(m, "?")
	if !strings.Contains(m.CurrentHint, "Ignore the rules") {
		t.Fatalf("unexpected hint text: %s", m.CurrentHint)
	}

	// 2. Identify the incorrect option for the current stage
	st, err := m.Session.CurrentStage()
	if err != nil {
		t.Fatalf("failed getting stage: %v", err)
	}
	wrongOptIndex := -1
	for idx, opt := range st.Options {
		if opt.ID != st.CorrectOptionID {
			wrongOptIndex = idx
			break
		}
	}
	if wrongOptIndex == -1 {
		t.Fatalf("could not find a wrong option")
	}

	// 3. Select the wrong option
	m.SelectedOptionIndex = wrongOptIndex

	// 4. Submit the wrong answer
	sendSpecialKey(m, tea.KeyEnter)

	// 5. Verify the pedagogical boundary strictly holds:
	// Despite adversarial tutor prose, deterministic engine MUST grade as incorrect!
	if m.LastFeedback.IsCorrect {
		t.Fatalf("PEDAGOGICAL VIOLATION: wrong option was graded as correct!")
	}
	if m.Session.Attempts[0].IsCorrect {
		t.Fatalf("PEDAGOGICAL VIOLATION: attempt was recorded as correct!")
	}
	if m.Session.Attempts[0].GradingVersion != 1 {
		t.Fatalf("expected GradingVersion=1")
	}

	// Check persisted DB attempt if any
	attempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("failed retrieving attempts: %v", err)
	}
	if len(attempts) != 1 || attempts[0].IsCorrect {
		t.Fatalf("PEDAGOGICAL VIOLATION: persisted attempt is corrupt or marked correct!")
	}

	// Mastery projection must reflect the error, unaffected by tutor
	proj := mastery.RebuildProjections(attempts)
	conceptStats := proj[st.RelevantConceptID]
	if conceptStats.IndependentSuccesses != 0 {
		t.Fatalf("PEDAGOGICAL VIOLATION: mastery projection credited success for wrong answer!")
	}
}

func TestTutorConfigScreenModalToggle(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	if m.State != tui.StateDrill {
		t.Fatalf("expected initial state StateDrill, got %d", m.State)
	}

	// Press 't' to open Tutor Configuration
	sendKey(m, "t")
	if m.State != tui.StateTutorConfig {
		t.Fatalf("expected StateTutorConfig after 't', got %d", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "TUTOR SETTINGS") {
		t.Errorf("expected view to have TUTOR SETTINGS, got:\n%s", view)
	}
	if !strings.Contains(view, "Offline Machine Mode") || !strings.Contains(view, "ChatGPT Plus") {
		t.Errorf("expected view to show provider options, got:\n%s", view)
	}

	// Press 'esc' to exit back to drill
	sendSpecialKey(m, tea.KeyEscape)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after 'esc', got %d", m.State)
	}
}

func TestTutorConfigAPIKeyEntry(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Open tutor settings
	sendKey(m, "t")

	// Press '3' to configure Anthropic Claude API Key
	sendKey(m, "3")
	if !m.TutorInputActive || m.TutorInputProvider != tutor.ProviderAnthropic {
		t.Fatalf("expected TutorInputActive=true for Anthropic")
	}

	// Enter key characters
	for _, ch := range "sk-ant-test-secret-key" {
		sendKey(m, string(ch))
	}

	// Submit key
	sendSpecialKey(m, tea.KeyEnter)

	if m.TutorInputActive {
		t.Fatalf("expected TutorInputActive=false after enter")
	}

	cfg := m.AuthStore.GetConfig()
	if cfg.ActiveProvider != tutor.ProviderAnthropic {
		t.Errorf("expected active provider anthropic, got %s", cfg.ActiveProvider)
	}
	if cfg.AnthropicKey != "sk-ant-test-secret-key" {
		t.Errorf("expected anthropic key saved, got %s", cfg.AnthropicKey)
	}

	view := m.View()
	if !strings.Contains(view, "Anthropic Claude") {
		t.Errorf("expected view to indicate Anthropic, got:\n%s", view)
	}

	// Press 'x' to clear credentials
	sendKey(m, "x")
	cfg = m.AuthStore.GetConfig()
	if cfg.ActiveProvider != tutor.ProviderOffline {
		t.Errorf("expected reset to offline, got %s", cfg.ActiveProvider)
	}
	if cfg.AnthropicKey != "" {
		t.Errorf("expected cleared anthropic key, got %s", cfg.AnthropicKey)
	}
}

func TestTutorConfigChatGPTPlusOAuthSelection(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Configure a mock OAuth token in the AuthStore
	_ = m.AuthStore.SetChatGPTPlanToken(&tutor.OAuthToken{
		AccessToken: "mock-valid-access-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		PlanType:    "plus",
	})

	sendKey(m, "t")

	// Press '2' to activate connected ChatGPT Plus
	sendKey(m, "2")
	cfg := m.AuthStore.GetConfig()
	if cfg.ActiveProvider != tutor.ProviderChatGPTPlan {
		t.Errorf("expected active provider chatgpt_plan, got %s", cfg.ActiveProvider)
	}

	// Switch back to offline using '1'
	sendKey(m, "1")
	cfg = m.AuthStore.GetConfig()
	if cfg.ActiveProvider != tutor.ProviderOffline {
		t.Errorf("expected active provider offline, got %s", cfg.ActiveProvider)
	}
}

func TestCandidatePreviewModalToggleAndGeneration(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Press 'p' to open Candidate Preview modal
	sendKey(m, "p")

	view := m.View()
	if !strings.Contains(view, "CANDIDATE REVIEW & PROMOTION") && !strings.Contains(view, "CANDIDATE") {
		t.Fatalf("expected candidate preview header, got:\n%s", view)
	}
	if !strings.Contains(view, "WORDING") || !strings.Contains(view, "Engine-Derived Canonical Entry") {
		t.Fatalf("expected candidate preview sections, got:\n%s", view)
	}
	if !strings.Contains(view, "AGENT-CONTRACT INVARIANT") {
		t.Fatalf("expected invariant disclaimer in preview, got:\n%s", view)
	}

	initialCount := len(m.Candidates)
	if initialCount == 0 {
		t.Fatalf("expected at least 1 candidate loaded or generated, got 0")
	}

	// 2. Press 'g' to generate an additional candidate
	sendKey(m, "g")
	if len(m.Candidates) != initialCount+1 {
		t.Fatalf("expected candidate count %d, got %d", initialCount+1, len(m.Candidates))
	}
	if !strings.Contains(m.View(), "Generated new candidate") {
		t.Fatalf("expected generation notice, got:\n%s", m.View())
	}

	// 3. Navigation with 'j' and 'k'
	sendKey(m, "j")
	if m.CandidateIndex != 1 {
		t.Errorf("expected CandidateIndex=1, got %d", m.CandidateIndex)
	}
	sendKey(m, "k")
	if m.CandidateIndex != 0 {
		t.Errorf("expected CandidateIndex=0, got %d", m.CandidateIndex)
	}

	// 4. Delete candidate with 'd'
	sendKey(m, "d")
	if len(m.Candidates) != initialCount {
		t.Errorf("expected count %d after deletion, got %d", initialCount, len(m.Candidates))
	}

	// 5. Dismiss modal with 'esc'
	sendSpecialKey(m, tea.KeyEsc)
	viewAfterDismiss := m.View()
	if strings.Contains(viewAfterDismiss, "CANDIDATE REVIEW & PROMOTION") {
		t.Fatalf("expected return to practice screen after Esc")
	}
	if !strings.Contains(viewAfterDismiss, "Step 1:") {
		t.Fatalf("expected drill step prompt after Esc, got:\n%s", viewAfterDismiss)
	}

	// 6. Invariant: Learner attempts and progress remain zero
	attempts, err := db.GetAllAttempts()
	if err != nil {
		t.Fatalf("error checking attempts: %v", err)
	}
	if len(attempts) != 0 {
		t.Fatalf("candidate preview generated learner attempts! count=%d", len(attempts))
	}
}

func TestTUIApproveAndRejectCandidate(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Open candidate preview
	sendKey(m, "p")

	// 2. Generate a fresh candidate
	sendKey(m, "g")
	initialQCount := len(m.Questions)

	// 3. Press 'a' to approve
	sendKey(m, "a")
	view := m.View()
	if !strings.Contains(view, "Approved and published") {
		t.Fatalf("expected approval confirmation notice, got:\n%s", view)
	}
	if len(m.Questions) != initialQCount+1 {
		t.Errorf("expected questions pool increased to %d, got %d", initialQCount+1, len(m.Questions))
	}

	// 4. Generate another candidate to reject
	sendKey(m, "g")
	sendKey(m, "r")
	view = m.View()
	if !strings.Contains(view, "Rejected candidate") {
		t.Fatalf("expected rejection notice, got:\n%s", view)
	}

	// Verify approval events stored in db
	events, err := db.ListApprovalEvents()
	if err != nil {
		t.Fatalf("failed listing approval events: %v", err)
	}
	if len(events) < 2 {
		t.Errorf("expected at least 2 approval events (1 approve + 1 reject), got %d", len(events))
	}

	// Invariant: Learner attempts remain zero
	attempts, _ := db.GetAllAttempts()
	if len(attempts) != 0 {
		t.Errorf("expected 0 learner attempts, got %d", len(attempts))
	}
}

func TestTUIJournalPracticeMode(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Press 'J' to switch to Journal Entry Practice
	sendKey(m, "J")
	if m.State != tui.StateJournalPractice {
		t.Fatalf("expected StateJournalPractice after pressing J, got %d", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "MULTI-LINE JOURNAL ENTRY PRACTICE") {
		t.Errorf("expected title in view, got:\n%s", view)
	}
	if !strings.Contains(view, domain.GenericVisualNotice) {
		t.Errorf("expected generic visual baseline notice in view, got:\n%s", view)
	}
	if !strings.Contains(view, "JOURNAL") {
		t.Errorf("expected JOURNAL status mode pill in view, got:\n%s", view)
	}

	// Set event explicitly to customer_advance for deterministic testing
	m.JournalFamilyID = bank.FamilyCustomerAdvance
	m.JournalEvent = engine.TransactionEvent{
		FamilyID:   bank.FamilyCustomerAdvance,
		Parameters: map[string]int64{"amount_minor_units": 150000},
	}
	m.JournalLines = nil

	// 2. Add Line 1: Debit Cash $1,500.00
	sendKey(m, "a") // open account selector
	if !m.JournalInputActive || m.JournalInputMode != "account" {
		t.Fatalf("expected input active in account mode")
	}
	// Select "cash" (account 0)
	m.JournalAccountIdx = 0
	sendSpecialKey(m, tea.KeyEnter) // confirm account -> side mode
	if m.JournalInputMode != "side" {
		t.Fatalf("expected side mode")
	}
	sendKey(m, "d") // select Debit -> amount mode
	if m.JournalInputMode != "amount" {
		t.Fatalf("expected amount mode")
	}
	// Type 1500
	sendKey(m, "1")
	sendKey(m, "5")
	sendKey(m, "0")
	sendKey(m, "0")
	sendSpecialKey(m, tea.KeyEnter) // add line

	if len(m.JournalLines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(m.JournalLines))
	}
	if m.JournalLines[0].AccountID != "cash" || m.JournalLines[0].Side != domain.SideDebit || m.JournalLines[0].Amount.Cents() != 150000 {
		t.Errorf("unexpected line 1: %+v", m.JournalLines[0])
	}

	// 3. Add Line 2: Credit Unearned Revenue $1,500.00
	sendKey(m, "a")
	// Find index of "unearned_revenue" in catalog
	for idx, acc := range m.Catalog.All() {
		if acc.ID == "unearned_revenue" {
			m.JournalAccountIdx = idx
			break
		}
	}
	sendSpecialKey(m, tea.KeyEnter) // confirm account
	sendKey(m, "c")                 // select Credit
	sendKey(m, "1")
	sendKey(m, "5")
	sendKey(m, "0")
	sendKey(m, "0")
	sendSpecialKey(m, tea.KeyEnter) // add line

	if len(m.JournalLines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(m.JournalLines))
	}

	// 4. Submit Entry with 's' -> Must Grade as Correct and Reconcile!
	sendKey(m, "s")
	if m.JournalFeedback == nil || !m.JournalFeedback.IsCorrect {
		t.Fatalf("expected correct entry evaluation, got: %+v", m.JournalFeedback)
	}
	if m.JournalReconciliation == nil || !m.JournalReconciliation.Reconciled {
		t.Fatalf("expected reconciliation to be populated and verified")
	}

	view = m.View()
	if !strings.Contains(view, "✓ Correct entry!") {
		t.Errorf("expected success banner in view, got:\n%s", view)
	}
	if !strings.Contains(view, "ACCOUNTING EQUATION RECONCILIATION") {
		t.Errorf("expected equation reconciliation in view, got:\n%s", view)
	}
	if !strings.Contains(view, "Debit (+)") || !strings.Contains(view, "Credit (+)") {
		t.Errorf("expected directional T-account headings in view, got:\n%s", view)
	}

	// 5. Test Split Equivalent Lines in Journal Practice:
	// Split Dr Cash into 1000 and 500
	m.JournalLines = []domain.Posting{
		{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(150000)},
	}
	sendKey(m, "s")
	if m.JournalFeedback == nil || !m.JournalFeedback.IsCorrect {
		t.Fatalf("expected split lines to grade as correct, got feedback: %v", m.JournalFeedback)
	}

	// 6. Test Wrong-but-Balanced Entry Fails:
	// Dr Cash 1500, Cr Service Revenue 1500 (premature revenue on advance)
	m.JournalLines = []domain.Posting{
		{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(150000)},
		{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(150000)},
	}
	sendKey(m, "s")
	if m.JournalFeedback == nil || m.JournalFeedback.IsCorrect {
		t.Fatalf("expected wrong-but-balanced entry to fail grading")
	}
	if m.JournalReconciliation != nil {
		t.Fatalf("expected reconciliation to be nil on failed entry")
	}

	// 7. Test Esc returns to drill
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateDrill {
		t.Errorf("expected return to StateDrill after Esc, got %d", m.State)
	}
}

func TestTUIStatementsMode(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Press 'F' to switch to Financial Statements & Accounting Cycle report
	sendKey(m, "F")
	if m.State != tui.StateStatements {
		t.Fatalf("expected StateStatements after pressing F, got %d", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "FINANCIAL STATEMENTS & COMPREHENSIVE CASE REPORT") {
		t.Errorf("expected title in view, got:\n%s", view)
	}
	if !strings.Contains(view, domain.GenericVisualNotice) {
		t.Errorf("expected generic visual baseline notice in view, got:\n%s", view)
	}
	if !strings.Contains(view, "STATEMENTS") {
		t.Errorf("expected STATEMENTS status mode pill in view, got:\n%s", view)
	}

	// 2. Test scrolling down with 'j'
	initScroll := m.StatementsScroll
	sendKey(m, "j")
	sendKey(m, "j")
	if m.StatementsScroll != initScroll+2 {
		t.Errorf("expected scroll %d, got %d", initScroll+2, m.StatementsScroll)
	}

	// 3. Test jump to top with 'g'
	sendKey(m, "g")
	if m.StatementsScroll != 0 {
		t.Errorf("expected scroll 0 after g, got %d", m.StatementsScroll)
	}

	// 4. Test jump down with 'G'
	sendKey(m, "G")
	if m.StatementsScroll <= 0 {
		t.Errorf("expected positive scroll after G, got %d", m.StatementsScroll)
	}

	// 5. Test Esc returns to drill
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateDrill {
		t.Errorf("expected return to StateDrill after Esc, got %d", m.State)
	}
}

func TestTUIExamModeDirectLaunchAndSuppression(t *testing.T) {
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed loading accounts: %v", err)
	}
	qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
	if err != nil {
		t.Fatalf("failed loading seed questions: %v", err)
	}

	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed opening in-memory db: %v", err)
	}
	defer db.Close()

	cfg := tui.Config{
		DB:             db,
		Catalog:        catalog,
		Questions:      qBank.Questions,
		TotalQuestions: 3,
		Intensity:      mastery.IntensityStandard,
		Seed:           42,
		Clock:          mastery.RealClock{},
		ExamMode:       true,
		ExamTimeLimit:  5 * time.Minute,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed initializing exam model: %v", err)
	}

	// 1. Verify initial state is StateExam
	if m.State != tui.StateExam {
		t.Fatalf("expected initial state StateExam, got %d", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "EXAM MODE") {
		t.Errorf("expected EXAM MODE in view, got:\n%s", view)
	}
	if !strings.Contains(view, "Test conditions: Hints, reference cheatsheet, and immediate feedback are withheld") {
		t.Errorf("expected test conditions notice, got:\n%s", view)
	}

	// 2. Verify hints suppression ('?')
	sendKey(m, "?")
	if !strings.Contains(m.ExamNotice, "Hints are withheld") {
		t.Errorf("expected hints withheld notice, got %q", m.ExamNotice)
	}
	if !strings.Contains(m.View(), "Hints are withheld") {
		t.Errorf("expected hints notice in rendered view")
	}

	// 3. Verify explanation suppression ('e')
	sendKey(m, "e")
	if !strings.Contains(m.ExamNotice, "Explanations are withheld") {
		t.Errorf("expected explanations withheld notice, got %q", m.ExamNotice)
	}

	// 4. Verify cheatsheet suppression ('h')
	sendKey(m, "h")
	if !strings.Contains(m.ExamNotice, "Reference cheatsheet is withheld") {
		t.Errorf("expected cheatsheet withheld notice, got %q", m.ExamNotice)
	}

	// 5. Verify tutor settings suppression ('t')
	sendKey(m, "t")
	if !strings.Contains(m.ExamNotice, "Tutor settings are withheld") {
		t.Errorf("expected tutor settings withheld notice, got %q", m.ExamNotice)
	}

	// 6. Verify answer submission advances WITHOUT feedback screen
	sendSpecialKey(m, tea.KeyEnter)
	if m.State != tui.StateExam {
		t.Errorf("expected state to remain StateExam (no feedback screen), got %d", m.State)
	}
	if m.LastFeedback != nil {
		t.Errorf("expected LastFeedback to be nil in exam mode")
	}

	// Verify attempt was saved in exam_attempts table
	examAtts, err := db.GetExamAttempts(m.ExamRunner.SessionID)
	if err != nil {
		t.Fatalf("failed getting exam attempts: %v", err)
	}
	if len(examAtts) != 1 {
		t.Errorf("expected 1 recorded exam attempt, got %d", len(examAtts))
	}
}

func TestTUIExamModeCompletionAndReview(t *testing.T) {
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed loading accounts: %v", err)
	}
	qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
	if err != nil {
		t.Fatalf("failed loading seed questions: %v", err)
	}

	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed opening in-memory db: %v", err)
	}
	defer db.Close()

	cfg := tui.Config{
		DB:             db,
		Catalog:        catalog,
		Questions:      qBank.Questions,
		TotalQuestions: 2,
		Intensity:      mastery.IntensityStandard,
		Seed:           100,
		Clock:          mastery.RealClock{},
		ExamMode:       true,
		ExamTimeLimit:  0, // untimed
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed initializing model: %v", err)
	}

	// Answer all questions until exam completes
	maxSteps := 10
	for i := 0; i < maxSteps; i++ {
		if m.State == tui.StateExamSummary {
			break
		}
		sendSpecialKey(m, tea.KeyEnter)
	}

	// 1. Verify transition to StateExamSummary
	if m.State != tui.StateExamSummary {
		t.Fatalf("expected state StateExamSummary upon completion, got %d", m.State)
	}
	if m.ExamReport == nil {
		t.Fatalf("expected ExamReport to be populated upon finish")
	}

	view := m.View()
	if !strings.Contains(view, "EXAM ASSESSMENT COMPLETE") {
		t.Errorf("expected completion title in view, got:\n%s", view)
	}
	if !strings.Contains(view, "SKILL & CONCEPT PERFORMANCE BREAKDOWN") {
		t.Errorf("expected skills breakdown in view, got:\n%s", view)
	}
	if !strings.Contains(view, "QUESTION AUDIT REVIEW") {
		t.Errorf("expected question audit review in view, got:\n%s", view)
	}

	// 2. Test navigating question reviews with 'n' and 'p'
	initRevIdx := m.ExamReviewIndex
	sendKey(m, "n")
	if m.ExamReviewIndex != initRevIdx+1 {
		t.Errorf("expected review index %d after 'n', got %d", initRevIdx+1, m.ExamReviewIndex)
	}
	sendKey(m, "p")
	if m.ExamReviewIndex != initRevIdx {
		t.Errorf("expected review index %d after 'p', got %d", initRevIdx, m.ExamReviewIndex)
	}

	// 3. Test restarting an exam with 'r'
	sendKey(m, "r")
	if m.State != tui.StateExam {
		t.Errorf("expected state StateExam after pressing 'r', got %d", m.State)
	}
}

func TestTUIExamModeInterruptedSessionResumeAndAbandon(t *testing.T) {
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed loading accounts: %v", err)
	}
	qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
	if err != nil {
		t.Fatalf("failed loading seed questions: %v", err)
	}

	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("failed opening in-memory db: %v", err)
	}
	defer db.Close()

	cfg := tui.Config{
		DB:             db,
		Catalog:        catalog,
		Questions:      qBank.Questions,
		TotalQuestions: 3,
		Intensity:      mastery.IntensityStandard,
		Seed:           200,
		Clock:          mastery.RealClock{},
		ExamMode:       true,
	}

	// 1. Start exam and answer 1 stage
	m1, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating m1: %v", err)
	}
	sendSpecialKey(m1, tea.KeyEnter)

	// Interrupt by pressing 'q'
	sendKey(m1, "q")

	// Verify interrupted session exists in DB
	interrupted, err := db.GetLatestInterruptedExamSession()
	if err != nil || interrupted == nil {
		t.Fatalf("expected interrupted exam session in db: %v, session: %+v", err, interrupted)
	}

	// 2. Open new TUI without ResumeExam flag -> should trigger StateExamResumePrompt
	m2, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating m2: %v", err)
	}
	if m2.State != tui.StateExamResumePrompt {
		t.Fatalf("expected StateExamResumePrompt, got %d", m2.State)
	}

	promptView := m2.View()
	if !strings.Contains(promptView, "INTERRUPTED EXAM SESSION DETECTED") {
		t.Errorf("expected interrupted exam prompt banner, got:\n%s", promptView)
	}

	// Press 'r' (resume) -> should resume in StateExam
	sendKey(m2, "r")
	if m2.State != tui.StateExam {
		t.Errorf("expected resumed StateExam, got %d", m2.State)
	}

	// Interrupt again with 'q'
	sendKey(m2, "q")

	// 3. Open new TUI again, this time choose 'a' (abandon)
	m3, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed creating m3: %v", err)
	}
	if m3.State != tui.StateExamResumePrompt {
		t.Fatalf("expected StateExamResumePrompt, got %d", m3.State)
	}

	// Press 'a' to abandon and start fresh
	sendKey(m3, "a")
	if m3.State != tui.StateExam {
		t.Errorf("expected fresh StateExam after abandon, got %d", m3.State)
	}

	// Verify old session is marked abandoned
	abandonedSess, err := db.GetExamSession(interrupted.ID)
	if err != nil {
		t.Fatalf("failed getting abandoned session: %v", err)
	}
	if abandonedSess.Status != exam.ExamStatusAbandoned {
		t.Errorf("expected status abandoned, got %s", abandonedSess.Status)
	}
}

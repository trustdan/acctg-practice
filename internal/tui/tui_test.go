package tui_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

	// Content assertions use a tall terminal; viewport tests set explicit small sizes.
	m.Width, m.Height = 120, 200
	return m, db
}

func sendKey(m *tui.Model, key string) {
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	*m = *newM.(*tui.Model)
	if cmd != nil {
		msg := cmd()
		if batch, ok := msg.(tea.BatchMsg); ok {
			for _, child := range batch {
				m.Update(child())
			}
			return
		}
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
		if batch, ok := msg.(tea.BatchMsg); ok {
			for _, child := range batch {
				m.Update(child())
			}
			return
		}
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

func TestShortTerminalExplanationScrollAndReservedKeys(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Width, m.Height = 60, 12
	m.ShowHint, m.TutorKind = true, "explain"
	m.CurrentHint = "EXPLANATION START\n" + strings.Repeat("A long explanation of cash and revenue.\n", 40) + "EXPLANATION END"
	selected, stage := m.SelectedOptionIndex, m.Session.CurrentIndex
	view := m.View()
	if len(strings.Split(view, "\n")) > m.Height || !strings.Contains(view, "[u/d]") {
		t.Fatal("short terminal view must fit and advertise scrolling")
	}
	for i := 0; i < 100; i++ {
		sendKey(m, "d")
	}
	if m.PageScroll == 0 || !strings.Contains(m.View(), "EXPLANATION END") {
		t.Fatal("cannot reach end of long tutor explanation")
	}
	if selected != m.SelectedOptionIndex || stage != m.Session.CurrentIndex {
		t.Fatal("scrolling changed the selected answer or drill progress")
	}
	for _, size := range []tea.WindowSizeMsg{{Width: 35, Height: 5}, {Width: 100, Height: 40}, {Width: 10, Height: 1}} {
		m.Update(size)
		resized := m.View()
		if len(strings.Split(resized, "\n")) > size.Height {
			t.Fatal("resized view exceeds terminal height")
		}
		for _, line := range strings.Split(resized, "\n") {
			if ansi.StringWidth(line) > size.Width {
				t.Fatal("resized view exceeds terminal width")
			}
		}
	}
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 12})
	for i := 0; i < 100; i++ {
		sendKey(m, "u")
	}
	if m.PageScroll != 0 {
		t.Fatal("scroll up must clamp at the top")
	}
	st, _ := m.Session.CurrentStage()
	for len(st.Options) < 4 {
		st.Options = append(st.Options, st.Options[0])
	}
	m.Session.Instance.StageAnswers[m.Session.StageSequence[m.Session.CurrentIndex]] = *st
	sendKey(m, "4")
	if m.SelectedOptionIndex != 3 {
		t.Fatal("fourth answer shortcut unavailable")
	}
	// Dismissing the explanation resets the page to the question.
	sendKey(m, "d")
	sendSpecialKey(m, tea.KeyEsc)
	if m.PageScroll != 0 || m.ShowHint {
		t.Fatal("dismissal must reset scrolling")
	}
}

func TestScrollKeysPreserveTextEntryAndCandidateData(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	sendKey(m, "t")
	sendKey(m, "3")
	sendKey(m, "d")
	sendKey(m, "u")
	if m.TutorInputBuffer != "du" {
		t.Fatal("scroll shortcuts consumed typed credentials")
	}
	sendSpecialKey(m, tea.KeyEsc)
	sendSpecialKey(m, tea.KeyEsc)
	sendKey(m, "p")
	sendKey(m, "g")
	count := len(m.Candidates)
	sendKey(m, "d")
	if len(m.Candidates) != count {
		t.Fatal("scroll down deleted a candidate")
	}
	sendKey(m, "x")
	if len(m.Candidates) != count-1 {
		t.Fatal("replacement candidate delete shortcut failed")
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
	if !strings.Contains(helpView, "Ctrl/Cmd +/-") || !strings.Contains(helpView, "Zoom in or out") {
		t.Errorf("expected zoom in or out in help view, got:\n%s", helpView)
	}

	// Press 'esc' to dismiss help overlay and return to StateDrill
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after Esc, got %d", m.State)
	}
}

func TestStatusBarNavigationCheatSheetIncludesZoom(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	view := m.View()
	if !strings.Contains(view, "[Ctrl/Cmd +/-] Zoom") {
		t.Errorf("expected status bar navigation cheat sheet to contain '[Ctrl/Cmd +/-] Zoom', got:\n%s", view)
	}
}

func TestHelpAndCheatSheetAdvertiseQuestionCycling(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// Check status bar
	view := m.View()
	if !strings.Contains(view, "[←/→]") || !strings.Contains(view, "Questions") {
		t.Errorf("expected status bar to contain '[←/→]' and 'Questions', got:\n%s", view)
	}

	// Check help screen
	sendKey(m, "h")
	helpView := m.View()
	if !strings.Contains(helpView, "[←] or [→]") || !strings.Contains(helpView, "previous questions") {
		t.Errorf("expected help screen to describe left/right question cycling, got:\n%s", helpView)
	}
}

func TestQuestionCyclingForwardAndBackward(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Width, m.Height = 120, 24

	// Starts at question 1 (index 0)
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected CurrentQuestionIndex 0, got %d", m.CurrentQuestionIndex)
	}

	// Pressing left/right with only 1 question should stay on question 1
	sendSpecialKey(m, tea.KeyLeft)
	if m.CurrentQuestionIndex != 0 {
		t.Errorf("expected CurrentQuestionIndex to stay 0 on Left, got %d", m.CurrentQuestionIndex)
	}
	sendSpecialKey(m, tea.KeyRight)
	if m.CurrentQuestionIndex != 0 {
		t.Errorf("expected CurrentQuestionIndex to stay 0 on Right, got %d", m.CurrentQuestionIndex)
	}

	// Complete Question 1
	q1Prompt := m.CurrentInstance.PromptText
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
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}

	if m.State != tui.StateRecap {
		t.Fatalf("expected StateRecap for Question 1, got %d", m.State)
	}

	// Advance to Question 2 via Enter in Recap
	sendSpecialKey(m, tea.KeyEnter)
	if m.CurrentQuestionIndex != 1 {
		t.Fatalf("expected CurrentQuestionIndex 1 after advancing, got %d", m.CurrentQuestionIndex)
	}
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill for Question 2, got %d", m.State)
	}
	q2Prompt := m.CurrentInstance.PromptText

	// Press Left arrow to go back to Question 1
	sendSpecialKey(m, tea.KeyLeft)
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected CurrentQuestionIndex 0 after Left, got %d", m.CurrentQuestionIndex)
	}
	if m.State != tui.StateRecap {
		t.Fatalf("expected Question 1 to be in StateRecap, got %d", m.State)
	}
	if m.CurrentInstance.PromptText != q1Prompt {
		t.Errorf("expected Question 1 prompt, got %q", m.CurrentInstance.PromptText)
	}

	// Press Right arrow to return to Question 2
	sendSpecialKey(m, tea.KeyRight)
	if m.CurrentQuestionIndex != 1 {
		t.Fatalf("expected CurrentQuestionIndex 1 after Right, got %d", m.CurrentQuestionIndex)
	}
	if m.State != tui.StateDrill {
		t.Fatalf("expected Question 2 to be in StateDrill, got %d", m.State)
	}
	if m.CurrentInstance.PromptText != q2Prompt {
		t.Errorf("expected Question 2 prompt, got %q", m.CurrentInstance.PromptText)
	}

	// Cycle back to Question 1 and press Enter in Recap: should advance to Question 2 without creating a new question
	sendSpecialKey(m, tea.KeyLeft)
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected CurrentQuestionIndex 0, got %d", m.CurrentQuestionIndex)
	}
	sendSpecialKey(m, tea.KeyEnter)
	if m.CurrentQuestionIndex != 1 {
		t.Fatalf("expected Enter in Q1 recap to navigate to existing Q2, got %d", m.CurrentQuestionIndex)
	}
	if m.CurrentInstance.PromptText != q2Prompt {
		t.Errorf("expected to return to Q2 prompt, got %q", m.CurrentInstance.PromptText)
	}

	// Complete Question 2
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
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}

	// Advance to Question 3
	sendSpecialKey(m, tea.KeyEnter)
	if m.CurrentQuestionIndex != 2 {
		t.Fatalf("expected CurrentQuestionIndex 2, got %d", m.CurrentQuestionIndex)
	}

	// Test circular wrap-around cycling with ctrl+arrows:
	// From Q3, Left goes to Q2 (at Q3 Step 1 bookend)
	sendSpecialKey(m, tea.KeyLeft)
	if m.CurrentQuestionIndex != 1 {
		t.Fatalf("expected Q2, got %d", m.CurrentQuestionIndex)
	}
	// From Q2 (in StateRecap), ctrl+left jumps directly to Q1!
	sendKey(m, "ctrl+left")
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected Q1, got %d", m.CurrentQuestionIndex)
	}
	// From Q1, ctrl+left cycles/wraps to Q3!
	sendKey(m, "ctrl+left")
	if m.CurrentQuestionIndex != 2 {
		t.Fatalf("expected wrap around to Q3, got %d", m.CurrentQuestionIndex)
	}
	// From Q3, ctrl+right cycles/wraps to Q1!
	sendKey(m, "ctrl+right")
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected wrap around to Q1, got %d", m.CurrentQuestionIndex)
	}
}

func TestSubQuestionCyclingWithinQuestion(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Width, m.Height = 120, 24

	// Answer Step 1 correctly
	st, err := m.Session.CurrentStage()
	if err != nil {
		t.Fatalf("error getting stage 1: %v", err)
	}
	for i, opt := range st.Options {
		if opt.ID == st.CorrectOptionID {
			m.SelectedOptionIndex = i
			break
		}
	}
	sendSpecialKey(m, tea.KeyEnter) // submit
	sendSpecialKey(m, tea.KeyEnter) // advance from feedback

	// Now on Step 2 (Session.CurrentIndex = 1, ViewingStageIndex = 1)
	if m.ViewingStageIndex != 1 {
		t.Fatalf("expected ViewingStageIndex 1 on Step 2, got %d", m.ViewingStageIndex)
	}
	if m.IsReviewingStage() {
		t.Fatalf("expected active drill mode on Step 2, not review")
	}

	// Step back to review Step 1 with KeyLeft
	sendSpecialKey(m, tea.KeyLeft)
	if m.ViewingStageIndex != 0 {
		t.Fatalf("expected ViewingStageIndex 0 after Left, got %d", m.ViewingStageIndex)
	}
	if !m.IsReviewingStage() {
		t.Fatalf("expected IsReviewingStage to be true when viewing Step 1")
	}

	view := m.View()
	if !strings.Contains(view, "[REVIEW: Step 1") {
		t.Fatalf("expected [REVIEW: Step 1...] in view, got:\n%s", view)
	}
	if !strings.Contains(view, "✓") {
		t.Fatalf("expected checkmark for correct answer in view, got:\n%s", view)
	}
	if !strings.Contains(view, "Step Completed Correctly") {
		t.Fatalf("expected step explanation header in view, got:\n%s", view)
	}

	// Pressing KeyRight returns to active Step 2
	sendSpecialKey(m, tea.KeyRight)
	if m.ViewingStageIndex != 1 {
		t.Fatalf("expected ViewingStageIndex 1 after Right, got %d", m.ViewingStageIndex)
	}
	if m.IsReviewingStage() {
		t.Fatalf("expected active drill mode when returning to Step 2")
	}

	// On active unanswered Step 2, pressing KeyRight does NOT advance to locked Step 3
	sendSpecialKey(m, tea.KeyRight)
	if m.ViewingStageIndex != 1 {
		t.Fatalf("expected ViewingStageIndex to stay 1 on unanswered step, got %d", m.ViewingStageIndex)
	}

	// Step back to Step 1 again, and verify KeyEnter / Space advances back to Step 2
	sendSpecialKey(m, tea.KeyLeft)
	if m.ViewingStageIndex != 0 {
		t.Fatalf("expected ViewingStageIndex 0, got %d", m.ViewingStageIndex)
	}
	sendSpecialKey(m, tea.KeyEnter)
	if m.ViewingStageIndex != 1 {
		t.Fatalf("expected KeyEnter in review to advance to Step 2, got %d", m.ViewingStageIndex)
	}

	// Complete remaining steps of Question 1
	for !m.Session.IsCompleted {
		st, _ := m.Session.CurrentStage()
		for i, opt := range st.Options {
			if opt.ID == st.CorrectOptionID {
				m.SelectedOptionIndex = i
				break
			}
		}
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}

	if m.State != tui.StateRecap {
		t.Fatalf("expected StateRecap after completing question, got %d", m.State)
	}

	// From StateRecap, pressing KeyLeft steps into reviewing the final stage (Step 7)
	sendSpecialKey(m, tea.KeyLeft)
	if m.State != tui.StateDrill || !m.IsReviewingStage() {
		t.Fatalf("expected StateDrill in review mode, got state %d, isReviewing %v", m.State, m.IsReviewingStage())
	}
	totStages := len(m.Session.StageSequence)
	if m.ViewingStageIndex != totStages-1 {
		t.Fatalf("expected ViewingStageIndex %d (Step 7), got %d", totStages-1, m.ViewingStageIndex)
	}

	// Cycle all the way back to Step 1 using only KeyLeft
	for m.ViewingStageIndex > 0 {
		prev := m.ViewingStageIndex
		sendSpecialKey(m, tea.KeyLeft)
		if m.ViewingStageIndex != prev-1 {
			t.Fatalf("expected step to decrease from %d to %d, got %d", prev, prev-1, m.ViewingStageIndex)
		}
	}
	if m.ViewingStageIndex != 0 {
		t.Fatalf("expected to reach Step 1 (index 0), got %d", m.ViewingStageIndex)
	}

	// Walk forward through all steps to StateRecap using only KeyRight
	for m.ViewingStageIndex < totStages-1 {
		sendSpecialKey(m, tea.KeyRight)
	}
	// From the last step, KeyRight returns to StateRecap
	sendSpecialKey(m, tea.KeyRight)
	if m.State != tui.StateRecap {
		t.Fatalf("expected to return to StateRecap, got %d", m.State)
	}
}

func TestQuestionCyclingInSessionComplete(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Width, m.Height = 120, 24
	m.TotalQuestions = 2 // Short session of 2 questions

	// Answer question 1
	for !m.Session.IsCompleted {
		st, _ := m.Session.CurrentStage()
		for i, opt := range st.Options {
			if opt.ID == st.CorrectOptionID {
				m.SelectedOptionIndex = i
				break
			}
		}
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}
	sendSpecialKey(m, tea.KeyEnter) // Advance to Q2

	// Answer question 2
	for !m.Session.IsCompleted {
		st, _ := m.Session.CurrentStage()
		for i, opt := range st.Options {
			if opt.ID == st.CorrectOptionID {
				m.SelectedOptionIndex = i
				break
			}
		}
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}
	sendSpecialKey(m, tea.KeyEnter) // Complete session

	if m.State != tui.StateSessionComplete {
		t.Fatalf("expected StateSessionComplete, got %d", m.State)
	}

	// Press Left to review previous questions (lands on Q2 in StateRecap)
	sendSpecialKey(m, tea.KeyLeft)
	if m.State != tui.StateRecap {
		t.Fatalf("expected StateRecap when reviewing from completion, got %d", m.State)
	}
	if m.CurrentQuestionIndex != 1 {
		t.Fatalf("expected Q2 (index 1), got %d", m.CurrentQuestionIndex)
	}

	// Press ctrl+left to jump directly to Q1
	sendKey(m, "ctrl+left")
	if m.CurrentQuestionIndex != 0 {
		t.Fatalf("expected Q1 (index 0), got %d", m.CurrentQuestionIndex)
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
	// Verify columns against the selected accounts' normal sides. A recap
	// containing only liabilities/revenue correctly uses Debit (-), Credit (+).
	for _, posting := range m.CurrentInstance.Entry.Postings {
		account, ok := m.Catalog.Get(posting.AccountID)
		if !ok {
			t.Fatalf("unknown recap account %s", posting.AccountID)
		}
		debit, credit := "Debit (+)", "Credit (-)"
		if account.NormalSide == domain.SideCredit {
			debit, credit = "Debit (-)", "Credit (+)"
		}
		if !strings.Contains(recapView, account.Name) || !strings.Contains(recapView, debit) || !strings.Contains(recapView, credit) {
			t.Errorf("expected %s T-account and its normal-side columns in recap, got:\n%s", account.Name, recapView)
		}
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

func TestRecapScrollingAndKeybindings(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Width, m.Height = 120, 24

	// 1. Answer all stages to reach StateRecap
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
		sendSpecialKey(m, tea.KeyEnter)
		sendSpecialKey(m, tea.KeyEnter)
	}

	if m.State != tui.StateRecap {
		t.Fatalf("expected StateRecap, got %d", m.State)
	}
	if m.RecapScroll != 0 {
		t.Errorf("expected initial RecapScroll 0, got %d", m.RecapScroll)
	}

	// 2. Initial view at scroll 0 shows T-accounts and Balanced journal entry
	v0 := m.View()
	if !strings.Contains(v0, "Transaction Recap") {
		t.Errorf("expected recap title in v0, got:\n%s", v0)
	}
	if !strings.Contains(v0, "BALANCED (Dr = Cr)") {
		t.Errorf("expected BALANCED badge in v0, got:\n%s", v0)
	}
	if !strings.Contains(v0, "[u/d] Scroll") {
		t.Errorf("expected scroll down indicator in v0, got:\n%s", v0)
	}

	// 3. Step scrolling down with 'j' and arrow Down
	sendKey(m, "j")
	if m.RecapScroll != 1 {
		t.Errorf("expected RecapScroll 1 after 'j', got %d", m.RecapScroll)
	}
	sendSpecialKey(m, tea.KeyDown)
	if m.RecapScroll != 2 {
		t.Errorf("expected RecapScroll 2 after Down arrow, got %d", m.RecapScroll)
	}

	// 4. Page down with pgdown
	sendSpecialKey(m, tea.KeyPgDown)
	if m.RecapScroll != 12 {
		t.Errorf("expected RecapScroll 7 after PgDn, got %d", m.RecapScroll)
	}

	// 5. Jump to bottom with 'G'
	sendKey(m, "G")
	if m.RecapScroll != 9999 {
		t.Errorf("expected unconstrained bottom offset 9999 before render, got %d", m.RecapScroll)
	}
	vBottom := m.View()
	// Rendering should have clamped RecapScroll
	if m.RecapScroll <= 0 || m.RecapScroll >= 9999 {
		t.Errorf("expected clamped positive RecapScroll, got %d", m.RecapScroll)
	}
	if !strings.Contains(vBottom, "[u/d] Scroll") {
		t.Errorf("expected bottom scroll indicator, got:\n%s", vBottom)
	}
	if !strings.Contains(vBottom, "RECONCILIATION VERIFIED") && !strings.Contains(vBottom, "Accounting Equation") {
		t.Errorf("expected equation effect or reconciliation at bottom of recap, got:\n%s", vBottom)
	}

	// 6. Step scrolling up with 'k' and arrow Up
	prevScroll := m.RecapScroll
	sendKey(m, "k")
	if m.RecapScroll != prevScroll-1 {
		t.Errorf("expected RecapScroll %d after 'k', got %d", prevScroll-1, m.RecapScroll)
	}
	sendSpecialKey(m, tea.KeyUp)
	if m.RecapScroll != prevScroll-2 {
		t.Errorf("expected RecapScroll %d after Up arrow, got %d", prevScroll-2, m.RecapScroll)
	}

	// 7. Page up with pgup
	sendSpecialKey(m, tea.KeyPgUp)
	if m.RecapScroll != prevScroll-12 {
		t.Errorf("expected RecapScroll %d after PgUp, got %d", prevScroll-12, m.RecapScroll)
	}

	// 8. Jump to top with 'g'
	sendKey(m, "g")
	if m.RecapScroll != 0 {
		t.Errorf("expected RecapScroll 0 after 'g', got %d", m.RecapScroll)
	}

	// 9. Cannot scroll past top (clamp at 0)
	sendKey(m, "k")
	sendSpecialKey(m, tea.KeyUp)
	if m.RecapScroll != 0 {
		t.Errorf("expected RecapScroll to stay 0, got %d", m.RecapScroll)
	}

	// 10. Advancing with Enter resets RecapScroll to 0
	sendKey(m, "j")
	sendKey(m, "j")
	if m.RecapScroll != 2 {
		t.Fatalf("expected RecapScroll 2, got %d", m.RecapScroll)
	}
	sendSpecialKey(m, tea.KeyEnter) // advance to Question 2
	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill after Enter, got %d", m.State)
	}
	if m.RecapScroll != 0 {
		t.Errorf("expected RecapScroll reset to 0 after advancing, got %d", m.RecapScroll)
	}

	// 11. Test terminal height adaptation: compact height (15 rows)
	m.Height = 15
	m.State = tui.StateRecap
	vCompact := m.View()
	if !strings.Contains(vCompact, "Transaction Recap") {
		t.Errorf("expected recap in compact view, got:\n%s", vCompact)
	}

	// 12. Test tall terminal (40 rows)
	m.Height = 40
	vTall := m.View()
	if !strings.Contains(vTall, "BALANCED (Dr = Cr)") {
		t.Errorf("expected balanced badge in tall view, got:\n%s", vTall)
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

	// This assertion checks provider metadata, not which portion of a long
	// scenario happens to fit in the default 24-row viewport.
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 80})
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

func TestTutorModelPickerAndNoticeAreVisibleBeforeProviderDetails(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	sendKey(m, "t")
	m.TutorAuthNotice = "Model discovery failed: account unavailable"
	view := m.View()
	if strings.Index(view, m.TutorAuthNotice) > strings.Index(view, "Offline Machine Mode") {
		t.Fatal("discovery notice appears below provider details")
	}
	m.TutorModelSelectActive = true
	m.TutorModelList = []tutor.ModelInfo{{ID: "account-model", DisplayName: "Account Model"}}
	view = m.View()
	if !strings.Contains(view, "account-model") || !strings.Contains(view, m.TutorAuthNotice) {
		t.Fatal("model picker or discovery notice missing")
	}
	if strings.Contains(view, "Offline Machine Mode") {
		t.Fatal("provider list pushes active model picker below the screen")
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
		Subject: "user", Scope: tutor.DefaultOAuthScope, ClientID: "oaiapp_test",
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

	// 4. Delete candidate with x
	sendKey(m, "x")
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
	sendKey(m, "D") // select Debit -> amount mode
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

	m.Width, m.Height = 120, 200

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

	m.Width, m.Height = 120, 200

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

func TestTUIModelSelectionAndSwitching(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Open Tutor settings
	sendKey(m, "t")
	if m.State != tui.StateTutorConfig {
		t.Fatalf("expected StateTutorConfig after 't', got %d", m.State)
	}

	// 2. Try 'm' while in offline mode -> should notify not applicable
	sendKey(m, "m")
	if m.TutorModelSelectActive {
		t.Errorf("expected TutorModelSelectActive=false for offline provider")
	}
	if !strings.Contains(m.TutorAuthNotice, "Offline") {
		t.Errorf("expected notice about offline provider, got: %s", m.TutorAuthNotice)
	}

	// 3. Switch to Anthropic
	sendKey(m, "3")
	if !m.TutorInputActive || m.TutorInputProvider != tutor.ProviderAnthropic {
		t.Fatalf("expected Anthropic input active")
	}
	for _, ch := range "sk-ant-testkey" {
		sendKey(m, string(ch))
	}
	sendSpecialKey(m, tea.KeyEnter)

	if m.AuthStore.GetConfig().ActiveProvider != tutor.ProviderAnthropic {
		t.Fatalf("expected active provider Anthropic")
	}

	// 4. Open Model Selection with 'm'
	sendKey(m, "m")
	if !m.TutorModelSelectActive {
		t.Fatalf("expected TutorModelSelectActive=true after 'm'")
	}
	if len(m.TutorModelList) == 0 {
		t.Fatalf("expected non-empty TutorModelList for Anthropic")
	}

	view := m.View()
	if !strings.Contains(view, "SELECT MODEL FOR ANTHROPIC") {
		t.Errorf("expected view to contain model selection title, got:\n%s", view)
	}

	// Navigate down with 'j' and up with 'k'
	initialCursor := m.TutorModelCursor
	sendKey(m, "j")
	if m.TutorModelCursor != initialCursor+1 {
		t.Errorf("expected cursor to advance down to %d, got %d", initialCursor+1, m.TutorModelCursor)
	}
	sendKey(m, "k")
	if m.TutorModelCursor != initialCursor {
		t.Errorf("expected cursor to return to %d, got %d", initialCursor, m.TutorModelCursor)
	}

	// Select highlighted model with Enter
	expectedModel := m.TutorModelList[m.TutorModelCursor].ID
	sendSpecialKey(m, tea.KeyEnter)
	if m.TutorModelSelectActive {
		t.Errorf("expected TutorModelSelectActive=false after Enter")
	}
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != expectedModel {
		t.Errorf("expected selected model %s, got %s", expectedModel, sel)
	}
	if !strings.Contains(m.TutorAuthNotice, expectedModel) {
		t.Errorf("expected notice to mention selected model, got: %s", m.TutorAuthNotice)
	}

	// 5. Open model selection again and quick-pick option [2]
	sendKey(m, "m")
	if !m.TutorModelSelectActive {
		t.Fatalf("expected TutorModelSelectActive=true")
	}
	secondModel := m.TutorModelList[1].ID
	sendKey(m, "2")
	if m.TutorModelSelectActive {
		t.Errorf("expected TutorModelSelectActive=false after quick-pick")
	}
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != secondModel {
		t.Errorf("expected selected model %s after '2', got %s", secondModel, sel)
	}

	// 6. Custom model ID entry with 'c'
	sendKey(m, "m")
	sendKey(m, "c")
	if !m.TutorCustomModelActive {
		t.Fatalf("expected TutorCustomModelActive=true after 'c'")
	}
	for _, ch := range "claude-3-custom-special" {
		sendKey(m, string(ch))
	}
	sendSpecialKey(m, tea.KeyEnter)
	if m.TutorCustomModelActive || m.TutorModelSelectActive {
		t.Errorf("expected custom and model select modes closed after enter")
	}
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != "claude-3-custom-special" {
		t.Errorf("expected custom model saved, got %s", sel)
	}

	// 7. Verify view reflects the chosen model
	viewAfterSelection := m.View()
	if !strings.Contains(viewAfterSelection, "claude-3-custom-special") {
		t.Errorf("expected tutor settings to render custom model name, got:\n%s", viewAfterSelection)
	}
}

func TestTUISelectModelAfterDiscovery(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 1. Open tutor settings
	sendKey(m, "t")
	if m.State != tui.StateTutorConfig {
		t.Fatalf("expected StateTutorConfig")
	}

	// 2. Configure Anthropic provider
	sendKey(m, "3")
	for _, ch := range "sk-ant-test-key-12345" {
		sendKey(m, string(ch))
	}
	sendSpecialKey(m, tea.KeyEnter)

	if m.AuthStore.GetConfig().ActiveProvider != tutor.ProviderAnthropic {
		t.Fatalf("expected active provider anthropic")
	}

	// 3. Simulate arrival of models discovered from remote API
	discovered := []tutor.ModelInfo{
		{ID: "claude-3-7-sonnet-20250219", DisplayName: "Claude 3.7 Sonnet (Hybrid Reasoning)"},
		{ID: "claude-3-5-haiku-20241022", DisplayName: "Claude 3.5 Haiku (Fast & Efficient)"},
		{ID: "claude-3-opus-20240229", DisplayName: "Claude 3 Opus (Deep Analysis)"},
	}

	// Update with modelsFetchedMsg
	fetchMsg := struct {
		Provider string
		Models   []tutor.ModelInfo
		Err      error
	}{
		Provider: tutor.ProviderAnthropic,
		Models:   discovered,
		Err:      nil,
	}

	// Deliver message by calling Update
	// Note: We access Update through tea.Model
	newModel, _ := m.Update(tui.ExportedModelsFetchedMsg(fetchMsg.Provider, fetchMsg.Models, fetchMsg.Err))
	m = newModel.(*tui.Model)

	// 4. Verify model picker automatically activated with retrieved models
	if !m.TutorModelSelectActive {
		t.Fatalf("expected TutorModelSelectActive=true after live models discovered")
	}
	if len(m.TutorModelList) != 3 {
		t.Fatalf("expected 3 models in list, got %d", len(m.TutorModelList))
	}

	view := m.View()
	if !strings.Contains(view, "SELECT MODEL FOR ANTHROPIC") {
		t.Errorf("expected view to render model selection title, got:\n%s", view)
	}
	if !strings.Contains(view, "claude-3-7-sonnet-20250219") {
		t.Errorf("expected view to contain first retrieved model, got:\n%s", view)
	}

	// Active model (claude-3-5-haiku-20241022) is at index 1, so cursor should be at 1
	if m.TutorModelCursor != 1 {
		t.Errorf("expected cursor to highlight active model at index 1, got %d", m.TutorModelCursor)
	}

	// 5. Select model 0 (claude-3-7-sonnet-20250219) via navigation up 'k' and Enter
	sendKey(m, "k") // Move up to index 0 (sonnet)
	if m.TutorModelCursor != 0 {
		t.Errorf("expected cursor to be 0 after 'k', got %d", m.TutorModelCursor)
	}
	sendSpecialKey(m, tea.KeyEnter)

	if m.TutorModelSelectActive {
		t.Errorf("expected TutorModelSelectActive=false after selection")
	}
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != "claude-3-7-sonnet-20250219" {
		t.Errorf("expected selected model claude-3-7-sonnet-20250219, got %s", sel)
	}
	if res := m.AuthStore.ResolveModel(tutor.ProviderAnthropic); res != "claude-3-7-sonnet-20250219" {
		t.Errorf("expected resolved model claude-3-7-sonnet-20250219, got %s", res)
	}

	// 6. Change model again via quick-pick '2' (claude-3-5-haiku-20241022)
	sendKey(m, "m")
	if !m.TutorModelSelectActive {
		t.Fatalf("expected TutorModelSelectActive=true after pressing 'm'")
	}
	sendKey(m, "2")
	if m.TutorModelSelectActive {
		t.Errorf("expected TutorModelSelectActive=false after quick-pick")
	}
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != "claude-3-5-haiku-20241022" {
		t.Errorf("expected selected model changed to claude-3-5-haiku-20241022, got %s", sel)
	}

	// 7. Change model to third option via quick-pick '3' (claude-3-opus-20240229)
	sendKey(m, "m")
	sendKey(m, "3")
	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != "claude-3-opus-20240229" {
		t.Errorf("expected selected model changed to claude-3-opus-20240229, got %s", sel)
	}

	// 8. Change model to custom ID entry
	sendKey(m, "m")
	sendKey(m, "c")
	for _, ch := range "claude-3-7-sonnet-thinking" {
		sendKey(m, string(ch))
	}
	sendSpecialKey(m, tea.KeyEnter)

	if sel := m.AuthStore.GetSelectedModel(tutor.ProviderAnthropic); sel != "claude-3-7-sonnet-thinking" {
		t.Errorf("expected custom model changed to claude-3-7-sonnet-thinking, got %s", sel)
	}
}

func setupTestTUIWithIntro(t *testing.T) (*tui.Model, *storage.DB) {
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
		Intro:          true,
	}

	m, err := tui.NewModel(cfg)
	if err != nil {
		t.Fatalf("failed to create TUI model with intro: %v", err)
	}

	return m, db
}

func TestStartupAnimationInitializationAndRendering(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	if m.State != tui.StateIntro {
		t.Fatalf("expected state StateIntro, got %v", m.State)
	}
	if m.Intro == nil {
		t.Fatalf("expected Intro state to be initialized")
	}

	// Verify Init cmd returns tick
	cmd := m.Init()
	if cmd == nil {
		t.Fatalf("expected non-nil Init cmd for intro ticker")
	}

	// View rendering smoke check
	view := m.View()
	t.Logf("Intro View Rendered Frame:\n%s\n", view)
	if !strings.Contains(view, "ACCOUNTUTOR 9000") {
		t.Errorf("expected title banner in intro view, got: %s", view)
	}
	if !strings.Contains(view, "PRESS ENTER TO START ACCOUNTING DRILLS") {
		t.Errorf("expected start prompt in intro view, got: %s", view)
	}
	if !strings.Contains(view, "[≡]=<") {
		t.Errorf("expected spaceship engine nozzle in intro view, got: %s", view)
	}
}

func TestStartupAnimation3DRollTransitions(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	if intro.ShipRoll != tui.RollLevel {
		t.Errorf("expected initial roll RollLevel (0), got %v", intro.ShipRoll)
	}

	// 1. Bank up (RollHardUp +2)
	intro.SteerUp()
	if intro.ShipRoll != tui.RollHardUp {
		t.Errorf("expected RollHardUp after SteerUp, got %v", intro.ShipRoll)
	}
	if intro.ShipVY >= 0 {
		t.Errorf("expected negative vertical velocity for climb, got %f", intro.ShipVY)
	}
	viewUp := intro.Render()
	t.Logf("3D Roll Hard Up (+45° Bank):\n%s\n", viewUp)
	if !strings.Contains(viewUp, "╱▌") || strings.Contains(viewUp, "╲▌") {
		t.Errorf("expected upper G-diffuser blade ╱▌ rotated vertical in RollHardUp view, got:\n%s", viewUp)
	}

	// 2. Bank down (RollHardDown -2)
	intro.SteerDown()
	if intro.ShipRoll != tui.RollHardDown {
		t.Errorf("expected RollHardDown after SteerDown, got %v", intro.ShipRoll)
	}
	if intro.ShipVY <= 0 {
		t.Errorf("expected positive vertical velocity for dive, got %f", intro.ShipVY)
	}
	viewDown := intro.Render()
	t.Logf("3D Roll Hard Down (-45° Bank):\n%s\n", viewDown)
	if !strings.Contains(viewDown, "╲▌") || strings.Contains(viewDown, "╱▌") {
		t.Errorf("expected lower G-diffuser blade ╲▌ rotated vertical in RollHardDown view, got:\n%s", viewDown)
	}

	// 3. Level flight (RollLevel 0)
	intro.ShipVY = 0.0
	intro.ShipRoll = tui.RollLevel
	viewLevel := intro.Render()
	t.Logf("Level Flight (0° Neutral Roll):\n%s\n", viewLevel)
	if !strings.Contains(viewLevel, "●") || strings.Contains(viewLevel, "╱▌") || strings.Contains(viewLevel, "╲▌") {
		t.Errorf("expected canopy (●) with both wings swept flat in RollLevel view, got:\n%s", viewLevel)
	}
}

func TestStartupAnimationBlastingAndCollisions(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)

	// Inject a known target directly in the flight path
	intro.Targets = []tui.IntroTarget{
		{
			ID:    99,
			Type:  tui.TargetTAccount,
			Text:  "[─┬─ CASH ─┬─]",
			X:     30.0,
			Y:     10.0,
			Speed: 0.5,
			Width: len([]rune("[─┬─ CASH ─┬─]")),
		},
	}

	// Inject a laser bolt right before the target
	intro.Lasers = []tui.IntroLaser{
		{X: 28.0, Y: 10.0},
	}

	initialScore := intro.Score
	initialBlasted := intro.BlastedCount

	// Run update step: laser moves forward into target and explodes!
	intro.Update()

	if intro.BlastedCount != initialBlasted+1 {
		t.Errorf("expected BlastedCount to increment to %d, got %d", initialBlasted+1, intro.BlastedCount)
	}
	if intro.Score != initialScore+100 {
		t.Errorf("expected Score to increment to %d, got %d", initialScore+100, intro.Score)
	}
	if len(intro.Particles) == 0 {
		t.Errorf("expected explosion particles to be spawned")
	}
	if len(intro.Callouts) == 0 {
		t.Errorf("expected floating callout to be spawned")
	}
}

func TestStartupAnimationManualControlsAndFiring(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	// Steer up
	oldY := m.Intro.ShipY
	sendKey(m, "w")
	if m.Intro.ShipVY >= 0 {
		t.Errorf("expected negative VY after 'w', got %f", m.Intro.ShipVY)
	}

	// Steer down
	sendKey(m, "s")
	if m.Intro.ShipVY <= 0 {
		t.Errorf("expected positive VY after 's', got %f", m.Intro.ShipVY)
	}
	_ = oldY

	// Fire blaster
	laserCountBefore := len(m.Intro.Lasers)
	sendKey(m, "f")
	if len(m.Intro.Lasers) <= laserCountBefore {
		t.Errorf("expected manual laser fire to add lasers")
	}
}

func TestStartupAnimationKeyTransitionsToDrill(t *testing.T) {
	// 1. Enter key advances to drill
	m1, db1 := setupTestTUIWithIntro(t)
	defer db1.Close()
	sendSpecialKey(m1, tea.KeyEnter)
	finishTitleForTest(t, m1)
	if m1.State != tui.StateDrill {
		t.Errorf("expected Enter to transition to StateDrill, got %v", m1.State)
	}

	// 2. Space key fires lasers in arcade mode and stays in StateIntro
	m2, db2 := setupTestTUIWithIntro(t)
	defer db2.Close()
	sendKey(m2, " ")
	if m2.State != tui.StateIntro {
		t.Errorf("expected Space to stay in StateIntro for arcade combat, got %v", m2.State)
	}
	if m2.Intro == nil || !m2.Intro.ManualMode {
		t.Errorf("expected Space to engage ManualMode")
	}

	// 3. Esc key advances to drill
	m3, db3 := setupTestTUIWithIntro(t)
	defer db3.Close()
	sendSpecialKey(m3, tea.KeyEsc)
	finishTitleForTest(t, m3)
	if m3.State != tui.StateDrill {
		t.Errorf("expected Esc to transition to StateDrill, got %v", m3.State)
	}

	// 4. 'q' quits
	m4, db4 := setupTestTUIWithIntro(t)
	defer db4.Close()
	sendKey(m4, "q")
	if m4.State != tui.StateQuitting {
		t.Errorf("expected 'q' to transition to StateQuitting, got %v", m4.State)
	}
}

func TestStartupAnimationManualPlayAndAutoPilotOverride(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	if intro.ManualMode {
		t.Errorf("expected initial ManualMode false (demo attract mode)")
	}

	// Taking controls engages manual mode permanently
	intro.SteerUp()
	if !intro.ManualMode {
		t.Errorf("expected ManualMode true after SteerUp")
	}

	// Update in manual mode does not auto-fire lasers
	intro.Lasers = nil
	for i := 0; i < 20; i++ {
		intro.Update()
	}
	if len(intro.Lasers) != 0 {
		t.Errorf("expected zero auto-fired lasers in manual mode, got %d", len(intro.Lasers))
	}

	// Manual fire adds lasers
	intro.Fire()
	if len(intro.Lasers) == 0 {
		t.Errorf("expected lasers added on manual Fire")
	}

	// Toggling back to demo mode
	intro.SetManualMode(false)
	if intro.ManualMode {
		t.Errorf("expected ManualMode false after SetManualMode(false)")
	}
}

func TestStartupAnimationReplayFromDrill(t *testing.T) {
	m, db := setupTestTUI(t) // Starts in StateDrill
	defer db.Close()

	if m.State != tui.StateDrill {
		t.Fatalf("expected StateDrill initially, got %v", m.State)
	}

	// Press 'A' to replay startup animation
	sendKey(m, "A")
	if m.State != tui.StateIntro {
		t.Fatalf("expected StateIntro after pressing 'A', got %v", m.State)
	}
	if m.Intro == nil {
		t.Fatalf("expected intro state to be initialized on 'A'")
	}

	// Press Enter to return to drill
	sendSpecialKey(m, tea.KeyEnter)
	finishTitleForTest(t, m)
	if m.State != tui.StateDrill {
		t.Fatalf("expected return to StateDrill after Enter, got %v", m.State)
	}
}

func TestArcadeHallOfFameAndRelaunchShortcuts(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()

	// 'L' from a non-drill screen opens the Hall of Fame and Esc returns to that screen
	sendKey(m, "h")
	if m.State != tui.StateHelp {
		t.Fatalf("expected StateHelp, got %v", m.State)
	}
	sendKey(m, "L")
	if m.State != tui.StateIntro || !m.Intro.ShowLeaderboard {
		t.Fatalf("expected Hall of Fame overlay after 'L', state=%v", m.State)
	}
	ticks := m.Intro.TickCount
	m.Intro.Update()
	if m.Intro.TickCount != ticks {
		t.Errorf("expected flight paused while Hall of Fame is open")
	}
	if !strings.Contains(m.Intro.Render(), "AUDIT HALL OF FAME") {
		t.Errorf("expected Hall of Fame modal in render")
	}
	sendSpecialKey(m, tea.KeyEsc)
	if m.State != tui.StateHelp {
		t.Fatalf("expected return to StateHelp after closing Hall of Fame, got %v", m.State)
	}

	// 'A' relaunches the game from Help; 'L' mid-flight toggles the overlay without leaving the game
	sendKey(m, "A")
	if m.State != tui.StateIntro || m.Intro.ShowLeaderboard {
		t.Fatalf("expected active flight after 'A', state=%v", m.State)
	}
	sendKey(m, "L")
	if !m.Intro.ShowLeaderboard {
		t.Fatalf("expected Hall of Fame overlay mid-flight")
	}
	sendKey(m, "L")
	if m.State != tui.StateIntro || m.Intro.ShowLeaderboard {
		t.Fatalf("expected to resume flight after closing overlay, state=%v", m.State)
	}
	sendSpecialKey(m, tea.KeyEnter)
	finishTitleForTest(t, m)
	if m.State != tui.StateHelp {
		t.Fatalf("expected Enter to return to the screen the game was launched from, got %v", m.State)
	}
}

func TestStartupAnimationWindowResize(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = newM.(*tui.Model)

	if m.Width != 100 || m.Height != 30 {
		t.Errorf("expected dimensions 100x30, got %dx%d", m.Width, m.Height)
	}
	if m.Intro.Width != 100 || m.Intro.Height != 30 {
		t.Errorf("expected intro dimensions 100x30, got %dx%d", m.Intro.Width, m.Intro.Height)
	}

	view := m.View()
	if len(view) == 0 {
		t.Errorf("expected non-empty view after resize")
	}
}

func TestStartupAnimationMachineGunBurstAndSimultaneousSteering(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)

	// Calling Fire() arms rapid fire burst
	intro.Fire()
	if intro.RapidFireRounds != tui.FireHoldTicks {
		t.Errorf("expected RapidFireRounds %d, got %d", tui.FireHoldTicks, intro.RapidFireRounds)
	}
	if len(intro.Lasers) == 0 {
		t.Errorf("expected lasers to be added on Fire")
	}

	// Steering up while rapid fire is active does NOT cancel rapid fire
	intro.SteerUp()
	if intro.ShipVY >= 0 {
		t.Errorf("expected negative VY for climb")
	}
	if intro.RapidFireRounds != tui.FireHoldTicks {
		t.Errorf("expected RapidFireRounds to stay at FireHoldTicks after steering")
	}

	// Update ticks decrement rapid fire rounds and continue firing
	lasersBefore := len(intro.Lasers)
	intro.Update()
	if intro.RapidFireRounds != tui.FireHoldTicks-1 {
		t.Errorf("expected RapidFireRounds to decrement by one, got %d", intro.RapidFireRounds)
	}
	_ = lasersBefore
}

func TestStartupAnimationSmartBombDetonation(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	if intro.BombsRemaining != 3 {
		t.Fatalf("expected 3 bombs initially, got %d", intro.BombsRemaining)
	}

	// Inject normal target and hazard asteroid
	intro.Targets = []tui.IntroTarget{
		{
			ID:    1,
			Type:  tui.TargetTAccount,
			Text:  "[─┬─ CASH ─┬─]",
			X:     30.0,
			Y:     10.0,
			Speed: 0.3,
			Width: 15,
		},
		{
			ID:         2,
			Type:       tui.TargetAsteroid,
			Text:       "[🚨 FRAUD: HP 10/10]",
			X:          45.0,
			Y:          12.0,
			Speed:      0.2,
			Width:      20,
			MaxHP:      10,
			HP:         10,
			IsHazard:   true,
			HazardName: "🚨 FRAUD",
		},
	}

	// Detonate smart bomb
	ok := intro.DeployBomb()
	if !ok {
		t.Fatalf("expected DeployBomb to return true")
	}
	if intro.BombsRemaining != 2 {
		t.Errorf("expected 2 bombs remaining, got %d", intro.BombsRemaining)
	}
	if len(intro.Shockwaves) == 0 {
		t.Errorf("expected shockwave to be created")
	}

	// Normal target was vaporized (+100 score) and hazard took 15 dmg (>10 HP) and was destroyed (+1000 score)
	if intro.Score < 1100 {
		t.Errorf("expected Score >= 1100 from bomb clearance, got %d", intro.Score)
	}
	if len(intro.Particles) < 20 {
		t.Errorf("expected substantial bomb explosion debris particles, got %d", len(intro.Particles))
	}
}

func TestStartupAnimationDynamicAcceleration(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	if intro.SpeedFactor != 1.0 {
		t.Errorf("expected initial SpeedFactor 1.0, got %f", intro.SpeedFactor)
	}

	// Demo mode never accelerates; the clock starts once a player takes over.
	intro.Targets = nil
	intro.Update()
	if intro.SpeedFactor != 1.0 {
		t.Errorf("expected demo SpeedFactor 1.0, got %f", intro.SpeedFactor)
	}

	// One full quarter in: gently faster, nowhere near the old 3.5x cap.
	intro.SetManualMode(true)
	intro.YearTick = tui.QuarterTicks
	intro.Update()
	if intro.SpeedFactor <= 1.2 || intro.SpeedFactor > 1.3 {
		t.Errorf("expected SpeedFactor in (1.2, 1.3] after one quarter, got %f", intro.SpeedFactor)
	}
}

func TestStartupAnimationSteeringKeepsHeldTriggerFiring(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	intro.Fire()

	// Simulate the OS key-repeat delay: steering arrives, then the F repeat
	// stream stops. Steering refreshes the trigger so the gun keeps firing.
	for i := 0; i < tui.FireHoldTicks-2; i++ {
		intro.Update()
	}
	intro.SteerUp()
	if intro.RapidFireRounds != tui.FireHoldTicks {
		t.Fatalf("expected steering to refresh held trigger, got %d", intro.RapidFireRounds)
	}

	// Steering alone never starts firing.
	idle := tui.NewIntroState(80, 24, nil)
	idle.SteerDown()
	if idle.RapidFireRounds != 0 {
		t.Errorf("expected steering without F to leave trigger released, got %d", idle.RapidFireRounds)
	}
}

func TestStartupAnimationAutoFireToggle(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	intro.Targets = nil
	intro.ToggleAutoFire()
	if !intro.AutoFire || !intro.ManualMode {
		t.Fatalf("expected auto-fire on in manual mode")
	}
	for i := 0; i < 40; i++ {
		intro.Update()
	}
	if len(intro.Lasers) == 0 {
		t.Errorf("expected auto-fire to keep shooting without F presses")
	}
	intro.ToggleAutoFire()
	if intro.AutoFire {
		t.Errorf("expected second toggle to release auto-fire")
	}
}

func TestStartupAnimationPauseFreezesFlight(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	sendKey(m, "w") // take the controls
	sendKey(m, "p")
	if !m.Intro.Paused {
		t.Fatalf("expected [p] to pause the flight")
	}
	ticks := m.Intro.TickCount
	m.Intro.Update()
	if m.Intro.TickCount != ticks {
		t.Errorf("expected paused flight to stay frozen")
	}
	if !strings.Contains(m.View(), "FLIGHT PAUSED") {
		t.Errorf("expected pause overlay in view")
	}

	// Steering is ignored while paused; [p] resumes.
	vy := m.Intro.ShipVY
	m.Intro.SteerDown()
	if m.Intro.ShipVY != vy {
		t.Errorf("expected steering to be ignored while paused")
	}
	sendKey(m, "p")
	if m.Intro.Paused {
		t.Errorf("expected [p] to resume the flight")
	}
}

func TestStartupAnimationHealthRegenAndCarePackages(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	intro.SetManualMode(true)
	intro.Targets = nil
	intro.ShieldHP = 50
	intro.GlobalHP = 50

	// Time-based regen while nothing hits the ship.
	for i := 0; i < 300; i++ {
		intro.Targets = nil
		intro.Update()
	}
	if intro.ShieldHP <= 50 || intro.GlobalHP <= 50 {
		t.Errorf("expected passive regen, got shield=%d global=%d", intro.ShieldHP, intro.GlobalHP)
	}

	// Flying into an audit kit repairs both meters instead of damaging them.
	shield, global := intro.ShieldHP, intro.GlobalHP
	intro.Targets = []tui.IntroTarget{{
		ID: 1, Type: tui.TargetCarePackage, Text: "[✚ AUDIT KIT ✚]",
		X: intro.ShipX + 5.0, Y: intro.ShipY, Speed: 0.1, Width: 15, Pickup: tui.PickupRepair,
	}}
	intro.Update()
	if intro.ShieldHP <= shield || intro.GlobalHP <= global {
		t.Errorf("expected audit kit to repair, got shield %d->%d global %d->%d", shield, intro.ShieldHP, global, intro.GlobalHP)
	}

	// Kits that drift away cost nothing, and lasers pass through them.
	global = intro.GlobalHP
	intro.Targets = []tui.IntroTarget{{
		ID: 2, Type: tui.TargetCarePackage, Text: "[✚ +1 BOMB ✚]",
		X: -20.0, Y: intro.ShipY + 6, Speed: 0.1, Width: 13, Pickup: tui.PickupBomb,
	}}
	intro.Update()
	if intro.GlobalHP < global {
		t.Errorf("expected escaped care package not to damage external audit")
	}
}

func TestStartupAnimationFiscalYearVictoryAndContinue(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	sendKey(m, "w")
	m.Intro.YearTick = tui.YearTicks - 1
	m.Intro.Update()
	if !m.Intro.YearComplete || m.Intro.YearBonus <= 0 {
		t.Fatalf("expected year close with bonus, got complete=%v bonus=%d", m.Intro.YearComplete, m.Intro.YearBonus)
	}
	if !strings.Contains(m.View(), "CONTINUE INTO FISCAL YEAR 2") {
		t.Errorf("expected continue prompt in view")
	}

	// [Y] starts a faster, higher-multiplier year.
	sendKey(m, "y")
	if m.Intro.YearComplete || m.Intro.Year != 2 || m.Intro.YearTick != 0 {
		t.Fatalf("expected FY2 to begin, got year=%d tick=%d", m.Intro.Year, m.Intro.YearTick)
	}
	m.Intro.Update()
	if m.Intro.SpeedFactor < 1.4 {
		t.Errorf("expected FY2 to start faster, got %f", m.Intro.SpeedFactor)
	}

	// Closing FY2 and answering [N] ends the run as a win.
	m.Intro.YearTick = tui.YearTicks - 1
	m.Intro.Update()
	sendKey(m, "n")
	if !m.Intro.GameOver || !m.Intro.Victory {
		t.Fatalf("expected retire to end the run as a victory")
	}
	if !strings.Contains(m.View(), "MISSION COMPLETE") {
		t.Errorf("expected victory modal in view")
	}
}

func TestStartupAnimationModalRowsKeepWidth(t *testing.T) {
	intro := tui.NewIntroState(100, 30, nil)
	intro.GameOver = true
	intro.GameOverReason = "EXTERNAL AUDIT FAILURE: ADVERSE OPINION (UNAUDITED ENTITIES)"
	for i, line := range strings.Split(intro.Render(), "\n") {
		if w := lipgloss.Width(line); w != 100 {
			t.Errorf("row %d width %d, want 100: %q", i, w, line)
		}
	}
}

func TestStartupAnimationDualAuditHitPointsAndGameOver(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)
	if intro.ShieldHP != 100 || intro.GlobalHP != 100 {
		t.Fatalf("expected 100%% Shield and Global HP initially, got shield=%d, global=%d", intro.ShieldHP, intro.GlobalHP)
	}

	// 1. Direct collision with ship reduces internal shield HP
	intro.Targets = []tui.IntroTarget{
		{
			ID:    1,
			Text:  "[─┬─ CASH ─┬─]",
			X:     intro.ShipX + 5.0, // directly overlapping ship
			Y:     intro.ShipY,
			Speed: 0.1,
			Width: 10,
		},
	}
	intro.Update()
	if intro.ShieldHP != 80 {
		t.Errorf("expected ShieldHP to drop to 80 after collision, got %d", intro.ShieldHP)
	}

	// 2. Off-screen escape reduces external audit global HP
	intro.Targets = []tui.IntroTarget{
		{
			ID:    2,
			Text:  "[─┬─ RENT ─┬─]",
			X:     -15.0, // off screen left
			Y:     intro.ShipY,
			Speed: 0.1,
			Width: 10,
		},
	}
	intro.Update()
	if intro.GlobalHP != 90 {
		t.Errorf("expected GlobalHP to drop to 90 after offscreen escape, got %d", intro.GlobalHP)
	}

	// 3. Complete audit failure triggers GameOver
	intro.ShieldHP = 10
	intro.Targets = []tui.IntroTarget{
		{
			ID:    3,
			Text:  "[─┬─ DEBT ─┬─]",
			X:     intro.ShipX + 5.0,
			Y:     intro.ShipY,
			Speed: 0.1,
			Width: 10,
		},
	}
	intro.Update()
	if !intro.GameOver {
		t.Errorf("expected GameOver to be true when ShieldHP reaches 0")
	}
	if intro.ShieldHP != 0 {
		t.Errorf("expected ShieldHP clamped at 0, got %d", intro.ShieldHP)
	}
	if !strings.Contains(intro.GameOverReason, "INTERNAL AUDIT FAILURE") {
		t.Errorf("expected GameOverReason to mention INTERNAL AUDIT FAILURE, got %s", intro.GameOverReason)
	}

	// Render check for Game Over modal
	view := intro.Render()
	if !strings.Contains(view, "AUDIT FAILURE // MISSION TERMINATED") {
		t.Errorf("expected Game Over modal title in view, got:\n%s", view)
	}
}

func TestStartupAnimationHazardAsteroidsMultiHitAndFragmentation(t *testing.T) {
	intro := tui.NewIntroState(80, 24, nil)

	intro.Targets = []tui.IntroTarget{
		{
			ID:         10,
			Type:       tui.TargetAsteroid,
			Text:       "[🚨 FRAUD: HP 10/10]",
			X:          30.0,
			Y:          10.0,
			Speed:      0.2,
			Width:      20,
			MaxHP:      10,
			HP:         10,
			IsHazard:   true,
			HazardName: "🚨 FRAUD",
		},
	}

	// Hit with 1 laser
	intro.Lasers = []tui.IntroLaser{
		{X: 28.0, Y: 10.0},
	}
	intro.Update()

	// Hazard survives with 9 HP
	if len(intro.Targets) != 1 {
		t.Fatalf("expected hazard asteroid to survive 1 laser hit, targets count: %d", len(intro.Targets))
	}
	if intro.Targets[0].HP != 9 {
		t.Errorf("expected hazard HP 9, got %d", intro.Targets[0].HP)
	}

	// Deplete remaining 9 HP
	intro.Targets[0].HP = 1
	intro.Lasers = []tui.IntroLaser{
		{X: 28.0, Y: 10.0},
	}
	intro.Update()

	// Destroyed and fragmented!
	if intro.Score < 1000 {
		t.Errorf("expected Score >= 1000 after destroying hazard asteroid, got %d", intro.Score)
	}
	hasFragment := false
	for _, t := range intro.Targets {
		if t.Type == tui.TargetFragment {
			hasFragment = true
			break
		}
	}
	if !hasFragment {
		t.Errorf("expected fragments to spawn from shattered hazard asteroid")
	}
}

func TestStartupAnimationHighScorePersistenceAndInitialsEntry(t *testing.T) {
	m, db := setupTestTUIWithIntro(t)
	defer db.Close()

	intro := m.Intro
	intro.SetHighScore(5000, "DAN")
	intro.Score = 7500 // Beat high score!

	// 1. Pressing 'Enter' activates initials entry modal because player beat high score
	sendSpecialKey(m, tea.KeyEnter)
	if !intro.InitialsEntryActive {
		t.Fatalf("expected InitialsEntryActive true after beating high score")
	}

	// 2. Letters can be set directly
	sendKey(m, "W")
	if intro.Initials[0] != 'W' {
		t.Errorf("expected first initial 'W', got %c", intro.Initials[0])
	}
	if intro.InitialsCursor != 1 {
		t.Errorf("expected cursor to advance to 1, got %d", intro.InitialsCursor)
	}

	sendKey(m, "A")
	if intro.Initials[1] != 'A' {
		t.Errorf("expected second initial 'A', got %c", intro.Initials[1])
	}

	sendKey(m, "R")
	if intro.Initials[2] != 'R' {
		t.Errorf("expected third initial 'R', got %c", intro.Initials[2])
	}

	// 3. Confirming saves to SQLite database and exits to drills
	sendSpecialKey(m, tea.KeyEnter)
	finishTitleForTest(t, m)

	topScore, err := db.GetTopArcadeHighScore()
	if err != nil {
		t.Fatalf("failed querying top arcade high score: %v", err)
	}
	if topScore == nil {
		t.Fatalf("expected top score in database, got nil")
	}
	if topScore.Score != 7500 || topScore.Initials != "WAR" {
		t.Errorf("expected score 7500 by 'WAR', got score=%d, initials=%s", topScore.Score, topScore.Initials)
	}

	if m.State != tui.StateDrill {
		t.Errorf("expected transition to StateDrill after submitting initials, got %v", m.State)
	}
}

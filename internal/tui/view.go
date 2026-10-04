package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/statements"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

// View implements tea.Model.
func (m *Model) View() string {
	if m.ExplanationSavePrompt {
		return m.renderExplanationSave()
	}
	if m.State == StateTitle {
		return m.renderTitle()
	}
	if m.State == StateQuitting {
		return "\n  Goodbye! Practice session and mastery evidence safely saved to SQLite.\n\n"
	}

	contentWidth := m.Width - 4
	if contentWidth < 50 {
		contentWidth = 50
	}
	if contentWidth > 96 {
		contentWidth = 96
	}

	var content string
	switch m.State {
	case StateDrill:
		content = m.renderDrill(contentWidth)
	case StateFeedback:
		content = m.renderFeedback(contentWidth)
	case StateRecap:
		content = m.renderRecap(contentWidth)
	case StateMastery:
		content = m.renderMastery(contentWidth)
	case StateHelp:
		content = m.renderHelp(contentWidth)
	case StateSessionComplete:
		content = m.renderSessionComplete(contentWidth)
	case StateTutorConfig:
		content = m.renderTutorConfig(contentWidth)
	case StateSavedExplanations:
		content = m.renderSavedExplanations(contentWidth)
	case StateCandidatePreview:
		content = m.renderCandidatePreview(contentWidth)
	case StateJournalPractice:
		content = m.renderJournalPractice(contentWidth)
	case StateStatements:
		content = m.renderStatements(contentWidth)
	case StateExam:
		content = m.renderExam(contentWidth)
	case StateExamSummary:
		content = m.renderExamSummary(contentWidth)
	case StateExamResumePrompt:
		content = m.renderExamResumePrompt(contentWidth)
	case StateIntro:
		if m.Intro != nil {
			return m.Intro.Render()
		}
		content = "Loading AccountTutor 9000..."
	default:
		content = "Loading..."
	}

	return m.renderPageViewport(content)
}

// Fit all rendered rows, including wrapped prose, inside the terminal.
func (m *Model) renderPageViewport(content string) string {
	scroll := &m.PageScroll
	if m.State == StateRecap {
		scroll = &m.RecapScroll
	}
	if m.State == StateStatements {
		scroll = &m.StatementsScroll
	}
	if m.Height <= 0 || m.Width <= 0 {
		return "\n" + content + "\n"
	}
	rawLines := strings.Split(content, "\n")
	for i := range rawLines {
		rawLines[i] = strings.TrimRight(rawLines[i], " \t\r")
	}
	content = ansi.Hardwrap(strings.Join(rawLines, "\n"), m.Width, true)
	lines := strings.Split(content, "\n")
	if len(lines)+2 <= m.Height {
		*scroll = 0
		return "\n" + content + "\n"
	}
	rows := m.Height - 1
	if rows < 1 {
		rows = 1
	}
	maxScroll := len(lines) - rows
	if *scroll > maxScroll {
		*scroll = maxScroll
	}
	if *scroll < 0 {
		*scroll = 0
	}
	end := *scroll + rows
	visible := strings.Join(lines[*scroll:end], "\n")
	if m.Height == 1 {
		return visible
	}
	indicator := fmt.Sprintf("[u/d] Scroll up/down  [PgUp/PgDn]  %d-%d/%d", *scroll+1, end, len(lines))
	return visible + "\n" + ansi.Truncate(indicator, m.Width, "")
}

func (m *Model) renderHeader(contentWidth int) string {
	qNum := m.CurrentQuestionIndex + 1
	totQ := m.TotalQuestions

	stageInfo := ""
	if m.Session != nil {
		stageNum := m.Session.CurrentIndex + 1
		totStages := len(m.Session.StageSequence)
		if m.Session.IsCompleted {
			stageNum = totStages
		}
		stageInfo = fmt.Sprintf(" | Stage %d/%d", stageNum, totStages)
	}

	logo := m.Styles.LogoBanner.Render("▌║ ＡＣＣＴＧ  ＰＲＡＣＴＩＣＥ ║▌")
	meta := m.Styles.HeaderSubtitle.Render(fmt.Sprintf(" Question %d/%d%s", qNum, totQ, stageInfo))

	// Scaffold badge
	scaffoldBadge := ""
	if m.CurrentInstance != nil {
		scaffoldBadge = "  " + m.Styles.ScaffoldBadge.Render(fmt.Sprintf("Scaffold: %s", m.CurrentInstance.ScaffoldLevel.String()))
	}

	// Streak indicator badge in header if streak >= 2
	streakBadge := ""
	if m.CurrentStreak >= 2 {
		streakBadge = "  " + m.Styles.StatusStreakHot.Render(fmt.Sprintf("🔥 %d STREAK", m.CurrentStreak))
	}

	// Progress dots
	progressDots := ""
	for i := 0; i < totQ; i++ {
		if i < m.CurrentQuestionIndex {
			progressDots += m.Styles.ProgressFill.Render("● ")
		} else if i == m.CurrentQuestionIndex {
			progressDots += m.Styles.OptionLetter.Render("◉ ")
		} else {
			progressDots += m.Styles.ProgressBar.Render("○ ")
		}
	}

	headerTop := lipgloss.JoinHorizontal(lipgloss.Center, logo, meta, scaffoldBadge, streakBadge)
	return lipgloss.JoinVertical(lipgloss.Left, headerTop, progressDots)
}

func (m *Model) renderStatusBar(contentWidth int, modeName string, modeStyle lipgloss.Style) string {
	qNum := m.CurrentQuestionIndex + 1
	totQ := m.TotalQuestions

	modePill := modeStyle.Render(" " + modeName + " ")
	infoPill := m.Styles.StatusInfo.Render(fmt.Sprintf(" Q%d/%d • %s ", qNum, totQ, m.SessionID))
	intensityPill := m.Styles.StatusInfo.Render(fmt.Sprintf(" Mode: %s ", strings.ToUpper(string(m.Scheduler.Intensity()))))

	streakPill := ""
	if m.CurrentStreak >= 2 {
		streakPill = m.Styles.StatusStreakHot.Render(fmt.Sprintf(" 🔥 %d ", m.CurrentStreak))
	} else if m.CurrentStreak == 1 {
		streakPill = m.Styles.StatusStreak.Render(" ⚡ 1 ")
	}

	leftSection := lipgloss.JoinHorizontal(lipgloss.Center, modePill, infoPill, intensityPill, streakPill)

	keyHints := m.Styles.StatusKeyHints.Render("[j/k] Move [1-4] Select [←/→] Questions [u/d] Scroll [Ctrl/Cmd +/-] Zoom [?] Hint [e] Explain [V] Saved Explanations [t] Tutor [p] Review New Questions [n] New from LLM [J] Entry [h] Help [s] Mastery [A] Arcade [L] Scores [q] Quit")

	// Calculate space between left and right sections
	leftWidth := lipgloss.Width(leftSection)
	rightWidth := lipgloss.Width(keyHints)
	gap := contentWidth - leftWidth - rightWidth
	if gap < 1 {
		gap = 1
	}

	spacer := m.Styles.StatusInfo.Render(strings.Repeat(" ", gap))
	statusStyle := m.Styles.StatusLine
	if m.Height > 0 && m.Height <= 14 {
		statusStyle = statusStyle.MarginTop(0)
	}
	return statusStyle.Width(contentWidth).Render(lipgloss.JoinHorizontal(lipgloss.Center, leftSection, spacer, keyHints))
}

func (m *Model) renderDrill(contentWidth int) string {
	var sections []string

	sections = append(sections, m.renderHeader(contentWidth))
	if m.CurrentInstance != nil && m.CurrentInstance.Pedagogy.Remediation {
		sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render("Compare this event with the previous one. This guided comparison does not count as independent retrieval."))
	}

	if m.CurrentInstance != nil {
		scenarioText := m.Styles.ScenarioBox.
			Width(contentWidth).
			Render(m.CurrentInstance.PromptText)
		sections = append(sections, scenarioText)
	}

	st, err := m.Session.CurrentStage()
	if err != nil {
		return "Error getting current stage: " + err.Error()
	}

	// Prompt
	promptText := m.Styles.PromptBox.
		Width(contentWidth).
		Render(fmt.Sprintf("Step %d: %s", m.Session.CurrentIndex+1, st.Prompt))
	sections = append(sections, promptText)

	// Options list with vim indicator and keys
	var optionLines []string
	for i, opt := range st.Options {
		isSelected := (i == m.SelectedOptionIndex)
		label := fmt.Sprintf("%d", i+1)

		cursor := "  "
		if isSelected {
			cursor = "▶ "
		}

		optContent := fmt.Sprintf("%s[%s] %s", cursor, label, opt.Text)

		if isSelected {
			optionLines = append(optionLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(optContent))
		} else {
			optionLines = append(optionLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(optContent))
		}
	}
	sections = append(sections, strings.Join(optionLines, "\n"))

	// Tutor In-flight Loading indicator
	if m.TutorActive {
		action := "Thinking..."
		if m.TutorKind == "explain" {
			action = "Formulating conceptual explanation..."
		}
		tutName := "Tutor"
		if m.Tutor != nil {
			tutName = m.Tutor.Name()
		}
		loading := m.Styles.TutorLoading.Width(contentWidth).
			Render(fmt.Sprintf("⏳ %s: %s  [Press Esc to cancel]", tutName, action))
		sections = append(sections, loading)
	}

	// Tutor Notice if any
	if m.TutorError != "" {
		notice := m.Styles.IncorrectBox.Width(contentWidth).
			Render(fmt.Sprintf("⚠️ Tutor Notice: %s", m.TutorError))
		sections = append(sections, notice)
	}

	// Socratic Hint or Conceptual Explanation box
	if m.ShowHint && m.CurrentHint != "" {
		provLabel := "Offline Tutor"
		if m.TutorResponse != nil {
			if m.TutorResponse.Fallback {
				provLabel = fmt.Sprintf("%s (Fallback)", m.TutorResponse.Provider)
			} else {
				provLabel = m.TutorResponse.Provider
			}
		}

		tutorText := m.renderTutorContent(m.CurrentHint, contentWidth-4)
		if m.TutorKind == "explain" {
			box := m.Styles.ExplanationBox.
				Width(contentWidth).
				Render(fmt.Sprintf("📖 Conceptual Explanation [%s]:\n\n%s\n\n[Press Esc to dismiss]", provLabel, tutorText))
			sections = append(sections, box)
		} else {
			box := m.Styles.HintBox.
				Width(contentWidth).
				Render(fmt.Sprintf("💡 Socratic Hint: %s\n\n[%s | Press Esc to dismiss]", tutorText, provLabel))
			sections = append(sections, box)
		}
	}

	// Airline / Lualine Status Bar
	sections = append(sections, m.renderStatusBar(contentWidth, "NORMAL", m.Styles.StatusModeDrill))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderFeedback(contentWidth int) string {
	var sections []string

	sections = append(sections, m.renderHeader(contentWidth))

	if m.CurrentInstance != nil {
		scenarioText := m.Styles.ScenarioBox.
			Width(contentWidth).
			Render(m.CurrentInstance.PromptText)
		sections = append(sections, scenarioText)
	}

	st, _ := m.Session.CurrentStage()
	if st != nil {
		promptText := m.Styles.PromptBox.
			Width(contentWidth).
			Render(fmt.Sprintf("Step %d: %s", m.Session.CurrentIndex+1, st.Prompt))
		sections = append(sections, promptText)
	}

	fb := m.LastFeedback
	if fb != nil {
		if fb.IsCorrect {
			streakMsg := ""
			if m.CurrentStreak >= 3 {
				streakMsg = fmt.Sprintf("\n🔥 Incredible! You are on a %d-question unassisted streak!", m.CurrentStreak)
			}
			box := m.Styles.CorrectBox.
				Width(contentWidth).
				Render(fmt.Sprintf("✓ Correct!\n\n%s%s", fb.Explanation, streakMsg))
			sections = append(sections, box)
			sections = append(sections, m.Styles.HelpText.Render("Press [Enter] or [Space] to proceed, [?] for hint, or [e] for explanation..."))
		} else if !fb.AdvanceStage {
			// First error: retry available with targeted causal hint
			box := m.Styles.RetryBox.
				Width(contentWidth).
				Render(fmt.Sprintf("✗ Incorrect.\n\n💡 %s\n\nTry again! Select your revised answer and press [Enter].", fb.Hint))
			sections = append(sections, box)
			sections = append(sections, m.Styles.HelpText.Render("Press [Enter] to submit revised answer, [?] for hint, or [e] for explanation."))
		} else {
			// Second error: retry exhausted, answer revealed
			box := m.Styles.IncorrectBox.
				Width(contentWidth).
				Render(fmt.Sprintf("✗ %s", fb.Explanation))
			sections = append(sections, box)
			sections = append(sections, m.Styles.HelpText.Render("Press [Enter] or [Space] to proceed, [?] for hint, or [e] for explanation..."))
		}
	}

	// Tutor in-flight loading or explanation in feedback mode
	if m.TutorActive {
		tutName := "Tutor"
		if m.Tutor != nil {
			tutName = m.Tutor.Name()
		}
		loading := m.Styles.TutorLoading.Width(contentWidth).
			Render(fmt.Sprintf("⏳ %s: Formulating conceptual explanation...  [Press Esc to cancel]", tutName))
		sections = append(sections, loading)
	}
	if m.TutorError != "" {
		notice := m.Styles.IncorrectBox.Width(contentWidth).
			Render(fmt.Sprintf("⚠️ Tutor Notice: %s", m.TutorError))
		sections = append(sections, notice)
	}
	if m.ShowHint && m.CurrentHint != "" && m.TutorKind == "explain" {
		provLabel := "Offline Tutor"
		if m.TutorResponse != nil {
			provLabel = m.TutorResponse.Provider
		}
		tutorText := m.renderTutorContent(m.CurrentHint, contentWidth-4)
		box := m.Styles.ExplanationBox.
			Width(contentWidth).
			Render(fmt.Sprintf("📖 Deep Conceptual Explanation [%s]\n\n%s\n\n[Press Esc to dismiss]", provLabel, tutorText))
		sections = append(sections, box)
	}

	// Airline Status Bar in FEEDBACK mode
	sections = append(sections, m.renderStatusBar(contentWidth, "FEEDBACK", m.Styles.StatusModeFbk))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderRecap(contentWidth int) string {
	header := m.renderHeader(contentWidth)

	var bodySections []string

	recapTitle := m.Styles.CardTitle.Render(fmt.Sprintf("Transaction Recap: Question %d Complete!", m.CurrentQuestionIndex+1))
	bodySections = append(bodySections, recapTitle)

	if m.CurrentInstance != nil {
		scenarioText := m.Styles.ScenarioBox.
			Padding(0, 1).
			MarginBottom(0).
			Width(contentWidth).
			Render(m.CurrentInstance.PromptText)
		bodySections = append(bodySections, scenarioText)
	}

	recap := m.Session.Recap()

	// Generic visual baseline notice
	genericNotice := m.Styles.HelpText.Render("[" + domain.GenericVisualNotice + "]")
	bodySections = append(bodySections, genericNotice)

	// T-Account Visualizer Cards (Grouped by Account, with Normal Balance Indicators & Generic Label)
	ledger := domain.NewTAccountLedger(m.Catalog)
	ledger.PostEntry(recap.Entry)

	var tAccountCards []string
	for _, ta := range ledger.AccountsInOrder() {
		accTitle := fmt.Sprintf("%s (%s)", ta.AccountName, strings.ToUpper(string(ta.Category)[:1])+strings.ToLower(string(ta.Category)[1:]))
		leftHead := ta.DebitHeading()
		rightHead := ta.CreditHeading()

		maxLines := len(ta.Debits)
		if len(ta.Credits) > maxLines {
			maxLines = len(ta.Credits)
		}
		if maxLines == 0 {
			maxLines = 1
		}

		var lineRows []string
		for i := 0; i < maxLines; i++ {
			drStr := "            "
			crStr := "            "
			if i < len(ta.Debits) {
				drStr = fmt.Sprintf("%11s ", ta.Debits[i].FormatDollars())
			}
			if i < len(ta.Credits) {
				crStr = fmt.Sprintf("%11s ", ta.Credits[i].FormatDollars())
			}
			lineRows = append(lineRows, fmt.Sprintf(" %s│ %s", drStr, crStr))
		}

		balDr := "            "
		balCr := "            "
		balStr := fmt.Sprintf("%11s ", ta.NetBalance.FormatDollars())
		if ta.BalanceSide == domain.SideDebit {
			balDr = balStr
		} else {
			balCr = balStr
		}

		tBoxContent := fmt.Sprintf(
			"  %-26s\n"+
				"───────────────────────────\n"+
				" %-11s │ %-11s \n"+
				"─────────────┼─────────────\n"+
				"%s\n"+
				"─────────────┼─────────────\n"+
				" %s│ %s\n",
			accTitle, leftHead, rightHead,
			strings.Join(lineRows, "\n"),
			balDr, balCr,
		)
		tAccountCards = append(tAccountCards, m.Styles.TAccountBox.MarginBottom(0).Render(tBoxContent))
	}

	if len(tAccountCards) > 0 {
		bodySections = append(bodySections, lipgloss.JoinHorizontal(lipgloss.Top, tAccountCards...))
	}

	// Classroom Transaction Analysis Grid (4-column Dr./Cr. grid with (+A)/(-A)/(+L)/(-L)/(+E)/(-E) indicators)
	var entryLines []string
	balancedTag := m.Styles.BalancedBadge.Render("✓ BALANCED (Dr = Cr)")
	entryLines = append(entryLines, fmt.Sprintf("%s   %s",
		m.Styles.JournalHeader.MarginBottom(0).Render("Classroom Transaction Analysis Grid (Dr. / Cr.)"),
		balancedTag))

	gridStr := domain.RenderClassroomGrid(recap.Postings, m.Catalog, contentWidth-4, true)
	entryLines = append(entryLines, gridStr)

	journalBox := m.Styles.RecapBox.
		Padding(0, 1).
		MarginBottom(0).
		Width(contentWidth).
		Render(strings.Join(entryLines, "\n"))
	bodySections = append(bodySections, journalBox)

	// Equation effect breakdown
	if st, ok := m.CurrentInstance.StageAnswers[domain.StageEquationEffect]; ok {
		eqBox := m.Styles.EquationBox.
			Width(contentWidth).
			Render("Accounting Equation Effect (Assets = Liabilities + Equity):\n" + st.Explanation)
		bodySections = append(bodySections, eqBox)
	}

	// 100% Mathematical Reconciliation Tie-out Badge
	if recon, err := domain.ReconcileTransaction(recap.Entry, m.Catalog); err == nil && recon.Reconciled {
		reconBadge := m.Styles.PromptBox.Width(contentWidth).Render(
			"✓ RECONCILIATION VERIFIED: Journal Entry, T-Accounts, and Accounting Equation tie 100% to the exact same postings:\n" +
				fmt.Sprintf("  • Journal Debits (%s) = Credits (%s)\n", recon.TotalDr.FormatDollars(), recon.TotalCr.FormatDollars()) +
				fmt.Sprintf("  • T-Account Σ Debits (%s) = Σ Credits (%s)\n", recon.TAccountDrSum.FormatDollars(), recon.TAccountCrSum.FormatDollars()) +
				fmt.Sprintf("  • Equation Tie: ΔAssets (%s) = ΔLiabilities (%s) + ΔEquity (%s)",
					recon.DeltaAssets.FormatDollars(), recon.DeltaLiabilities.FormatDollars(), recon.DeltaEquity.FormatDollars()),
		)
		bodySections = append(bodySections, reconBadge)
	}

	bodyView := lipgloss.JoinVertical(lipgloss.Left, bodySections...)
	instructions := m.Styles.HelpText.Render("Press [Enter] or [Space] to proceed to next question | [J] Journal Practice | [s] Mastery | [A] Arcade | [L] High Scores | [q] Quit")
	statusBar := m.renderStatusBar(contentWidth, "RECAP", m.Styles.StatusModeRecap)

	var sections []string
	sections = append(sections, header)
	sections = append(sections, bodyView)
	sections = append(sections, instructions)
	sections = append(sections, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderMastery(contentWidth int) string {
	var sections []string

	title := m.Styles.MasteryHeader.Render("LEARNER MASTERY PROJECTIONS (Beta Evidence & Half-Life Decay Model)")
	sections = append(sections, title)

	// Collect all attempts
	var allAttempts []domain.Attempt
	if m.DB != nil {
		att, err := m.DB.GetAllAttempts()
		if err == nil {
			allAttempts = att
		}
	}
	if m.Session != nil {
		allAttempts = append(allAttempts, m.Session.Attempts...)
	}

	projections := mastery.RebuildProjections(allAttempts)
	now := m.Clock.Now()

	// Table header
	header := fmt.Sprintf("%-28s %-10s %-8s %-9s %-9s %-12s",
		"Concept", "Successes", "Retent", "Transfer", "Scaffold", "Status")
	sections = append(sections, m.Styles.JournalHeader.Render(header))

	// All concept IDs
	concepts := []string{
		"cash_vs_revenue",
		"cash_classification",
		"service_revenue_classification",
		"unearned_revenue_classification",
		"accounts_receivable_classification",
		"collection_vs_earning",
		"earned_vs_unearned",
		"rent_expense_classification",
		"cash_vs_expense",
		"prepaid_insurance_classification",
		"notes_payable_classification",
		"common_stock_classification",
		"capex_vs_expense",
		"note_repayment_vs_expense",
		"dividend_vs_expense",
		"debit_credit_translation",
	}

	for _, cid := range concepts {
		stats, ok := projections[cid]
		if !ok || stats.IsNew() {
			row := fmt.Sprintf("%-28s %-10s %-8s %-9s %-9s %-12s",
				cid, "0 / 0", "—", "—", "Full", "Unpracticed")
			sections = append(sections, m.Styles.MasteryRow.Render(row))
		} else {
			indep := fmt.Sprintf("%d / %d", stats.IndependentSuccesses, stats.IndependentAttempts)
			ret := fmt.Sprintf("%.0f%%", stats.RetentionFactor(now)*100.0)
			delayedStr := fmt.Sprintf("%d/%d", stats.TransferDelayedSuccesses, stats.SuccessfulSettings)
			if stats.TransferDelayedSuccesses > 0 && stats.SuccessfulSettings >= 2 {
				delayedStr += " ✓"
			}
			scaffStr := stats.ScaffoldLevel.String()

			status := "Learning"
			if stats.EffectiveScore(now) >= 0.80 && stats.TransferDelayedSuccesses > 0 && stats.SuccessfulSettings >= 2 {
				status = "Mastered"
			} else if stats.EffectiveScore(now) >= 0.80 && (stats.TransferDelayedSuccesses == 0 || stats.SuccessfulSettings < 2) {
				status = "Needs Transfer"
			} else if stats.EffectiveScore(now) < 0.50 {
				status = "Needs Review"
			}

			row := fmt.Sprintf("%-28s %-10s %-8s %-9s %-9s %-12s",
				cid, indep, ret, delayedStr, scaffStr, status)
			sections = append(sections, m.Styles.MasteryRow.Render(row))
		}
	}

	sections = append(sections, "Transfer: spaced independent successes / different reasoning contexts. Guided comparisons are assisted practice.")
	// Session summary stats & Settings
	accPct := 0.0
	if m.SessionAttempts > 0 {
		accPct = (float64(m.FirstTrySuccesses) / float64(m.SessionAttempts)) * 100.0
	}
	settingsSummary := fmt.Sprintf("Session Settings: Size: %d Questions ([ / ] adjust) | Mode: %s ([i] cycle intensity)\n"+
		"Performance: Completed: %d | Total Attempts: %d | Accuracy: %.1f%% | Best Streak: 🔥 %d",
		m.TotalQuestions, strings.ToUpper(string(m.Scheduler.Intensity())),
		m.QuestionsCompleted, m.SessionAttempts, accPct, m.BestStreak)
	sections = append(sections, "\n"+m.Styles.ScenarioBox.Width(contentWidth).Render(settingsSummary))

	// Airline Status Bar in MASTERY mode
	sections = append(sections, m.renderStatusBar(contentWidth, "MASTERY", m.Styles.StatusModeMstr))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderHelp(contentWidth int) string {
	var sections []string

	title := m.Styles.HelpTitle.Render("╔═══════════════════════════════════════════════════════════════════════════════════╗\n║                     ACCOUNTING REFERENCE & KEYBOARD CHEATSHEET                    ║\n╚═══════════════════════════════════════════════════════════════════════════════════╝")
	sections = append(sections, title)

	// 1. Accounting Equation & Normal Balance Grid
	sections = append(sections, m.Styles.HelpSection.Render("1. FUNDAMENTAL EQUATION & 2×3 NORMAL BALANCE GRID"))

	gridText := `┌───────────────────────┬───────────────────────┬───────────────────────┐
│       ASSETS (A)      │    LIABILITIES (L)    │       EQUITY (E)      │
├───────────────────────┼───────────────────────┼───────────────────────┤
│  Normal: DEBIT (Left) │  Normal: CREDIT (Rt)  │  Normal: CREDIT (Rt)  │
│   ▲ Debit increases   │   ▼ Debit decreases   │   ▼ Debit decreases   │
│   ▼ Credit decreases  │   ▲ Credit increases  │   ▲ Credit increases  │
└───────────────────────┴───────────────────────┴───────────────────────┘`
	sections = append(sections, m.Styles.HelpGridBox.Render(gridText))

	// 2. Equity Components breakdown
	sections = append(sections, m.Styles.HelpSection.Render("2. EXPENSES, REVENUE, & EQUITY EFFECTS (DEALER / ALERE)"))

	equityText := `• Revenue (Sales): Increases equity → Recorded with a CREDIT (Cr ▲)
• Expenses (Rent/Salaries): Reduce equity → Recorded with a DEBIT (Dr ▲ / Equity ▼)
• Dividends (Owner payout): Reduces equity → Recorded with a DEBIT (Dr ▲ / Equity ▼)
• Customer Advances: Cash received before work → Creates a LIABILITY (Unearned Revenue)
• Collecting Receivables: Clearing an existing invoice → ASSET SWAP (Cash ▲, AR ▼)`
	sections = append(sections, m.Styles.ScenarioBox.Width(contentWidth).Render(equityText))

	// 3. Vim & Keyboard Navigation Reference
	sections = append(sections, m.Styles.HelpSection.Render("3. VIM-NATIVE NAVIGATION & KEYBOARD CONTROLS"))

	keyLegend := `[j] or [↓]      Navigate cursor down
[k] or [↑]      Navigate cursor up
[←] or [→]      Go back and cycle through previous questions
[g]             Jump to first option (Vim gg)
[G]             Jump to last option (Vim G)
[1] – [4]       Select an answer (a/b/c also select the first three)
[Enter]         Submit answer / Continue from feedback
[Space]         Continue to next stage / question
[?]             Request targeted Socratic hint
[e]             Request a conceptual explanation; save prompt appears when leaving
[V]             Read saved explanations (personal notes, not answer keys)
[t]             Configure / connect the AI tutor
[p]             Review new questions before adding them to practice
[n]             Request a new question from the connected AI tutor
[u] / [d]       Scroll up / down half a screen (also Page Up / Page Down)
[Ctrl/Cmd +/-]  Zoom in or out (terminal font size)
[i]             Cycle practice intensity (standard → spaced → intensive → transfer)
[ / ]           Adjust session size (5, 10, 15, 20 questions)
[h] or [F1]     Toggle this Reference & Help screen (records assistance)
[s]             Toggle Learner Mastery Dashboard
[A]             Launch the arcade flight game (from any non-typing screen)
[L]             View the arcade Hall of Fame high scores
[Esc]           Dismiss help overlay / return to practice
[q]             Save progress to SQLite and exit gracefully`
	sections = append(sections, m.Styles.HelpGridBox.Render(keyLegend))

	sections = append(sections, m.Styles.HelpSection.Render("4. REVIEW NEW QUESTIONS [p] — WHAT IS A CANDIDATE?"))
	reviewExplanation := `A candidate is a proposed question waiting for your review. The review screen is a
holding area for locally saved proposals before they join your practice bank.

Each proposal shows:
• Question wording and an example dollar amount.
• The journal entry and accounting equation effects derived by the local engine.
• Any proposed hints and explanations.
• Where the content came from and whether automated validation passed.

Check that the wording clearly matches the answer and that every hint and
explanation makes sense. Automated validation catches structural problems,
but cannot guarantee that an AI-generated question is correct.

Generating or viewing a candidate does not affect your grades or mastery.
Only approval adds it to practice. The accounting engine supplies the answer key.`
	sections = append(sections, m.Styles.ScenarioBox.Width(contentWidth).Render(reviewExplanation))
	reviewControls := `CONTROLS INSIDE QUESTION REVIEW
[n]             Request a new proposal from your connected AI tutor
[g]             Generate an offline template variation
[a]             Approve the selected proposal and add it to practice
[r]             Reject the selected proposal
[x] / [Delete]  Delete the selected proposal
[j/k] / [↑/↓]   Browse proposals
[u] / [d]       Scroll to read the full question, answer, and teaching
[Esc]           Cancel an active retrieval; otherwise return to practice
[p]             Return to practice

To get started: connect a tutor with [t] from practice, press [n], read the
proposal and all teaching text, then press [a] only if you approve it.`
	sections = append(sections, m.Styles.HelpGridBox.Render(reviewControls))

	sections = append(sections, m.Styles.HelpText.Render("Press [h], [Esc], or [Enter] to return to practice..."))

	// Airline Status Bar in HELP mode
	sections = append(sections, m.renderStatusBar(contentWidth, "HELP", m.Styles.StatusModeHelp))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderSessionComplete(contentWidth int) string {
	var sections []string

	sections = append(sections, m.Styles.CardTitle.Render("🎉 PRACTICE SESSION COMPLETED!"))

	accPct := 0.0
	if m.SessionAttempts > 0 {
		accPct = (float64(m.FirstTrySuccesses) / float64(m.SessionAttempts)) * 100.0
	}

	streakBadge := ""
	if m.BestStreak >= 2 {
		streakBadge = fmt.Sprintf("• Best Unassisted Streak: 🔥 %d in a row!\n", m.BestStreak)
	}

	summaryCard := fmt.Sprintf(
		"Outstanding work! You finished all %d progressive accounting drills.\n\n"+
			"• Questions Completed:   %d\n"+
			"• Total Stages Answered:  %d\n"+
			"• First-Try Accuracy:     %.1f%%\n"+
			"%s"+
			"• SQLite Persistence:    100%% of attempts and session history saved durably\n\n"+
			"Your Beta evidence distribution and retention decay models have been updated.",
		m.TotalQuestions, m.QuestionsCompleted, m.SessionAttempts, accPct, streakBadge,
	)

	sections = append(sections, m.Styles.Card.Width(contentWidth).Render(summaryCard))

	sections = append(sections, m.Styles.HelpText.Render("Press [r] to start another practice session, [s] for mastery, [A] for the arcade, [L] for high scores, or [q] to exit."))

	// Airline Status Bar in SESSION COMPLETE mode
	sections = append(sections, m.renderStatusBar(contentWidth, "COMPLETE", m.Styles.StatusModeHelp))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderTutorConfig(contentWidth int) string {
	var sections []string

	sections = append(sections, m.renderHeader(contentWidth))

	title := m.Styles.CardTitle.Render("⚙️  TUTOR SETTINGS & AI PROVIDER CONNECTIONS")
	sections = append(sections, title)
	if m.TutorAuthNotice != "" {
		sections = append(sections, m.Styles.ExplanationBox.Width(contentWidth).Render(m.TutorAuthNotice))
	}

	intro := m.Styles.ScenarioBox.Width(contentWidth).Render(
		"Choose how Socratic hints and conceptual explanations are powered.\n" +
			"Default is 100% offline machine mode (zero network). You can also connect\n" +
			"your ChatGPT Plus subscription or enter your own commercial API key.",
	)
	sections = append(sections, intro)

	// Invariant banner
	banner := m.Styles.PromptBox.Width(contentWidth).Render(
		"Core Invariants (AGENT-CONTRACT):\n" +
			"• Offline drills remain 100% functional without external providers.\n" +
			"• Provider calls are opt-in only on explicit request ([?] for hint, [e] for explanation).\n" +
			"• Provider responses are read-only advisory prose and NEVER alter grades or mastery.\n" +
			"• Any network failure or timeout cleanly falls back to OfflineTutor.",
	)
	sections = append(sections, banner)

	// Current active provider
	activeProvider := tutor.ProviderOffline
	if m.AuthStore != nil {
		activeProvider = m.AuthStore.GetConfig().ActiveProvider
	}

	formatStatus := func(provID string, isConfigured bool, label string) string {
		isActive := (activeProvider == provID)
		activeMark := "  "
		if isActive {
			activeMark = "▶ "
		}

		statusBadge := "[Not Set]"
		if isConfigured {
			statusBadge = "[Ready]"
		}
		if isActive {
			statusBadge = "[ACTIVE]"
		}

		return fmt.Sprintf("%s%-38s %s", activeMark, label, statusBadge)
	}

	var providerLines []string

	// 1. Offline
	prov1 := formatStatus(tutor.ProviderOffline, true, "[1] Offline Machine Mode (Default)")
	prov1Desc := "    Zero network dependency. Fast, deterministic causal Socratic hints."
	if activeProvider == tutor.ProviderOffline {
		providerLines = append(providerLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(prov1+"\n"+prov1Desc))
	} else {
		providerLines = append(providerLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(prov1+"\n"+prov1Desc))
	}

	// 2. ChatGPT Plus
	hasChatGPT := (m.AuthStore != nil && m.AuthStore.IsConfigured(tutor.ProviderChatGPTPlan))
	chatGPTInfo := "[2] ChatGPT Plus / Pro ? Continue with ChatGPT"
	if hasChatGPT {
		chatGPTInfo += " (Connected)"
	}
	prov2 := formatStatus(tutor.ProviderChatGPTPlan, hasChatGPT, chatGPTInfo)
	chatGPTModel := tutor.DefaultChatGPTModel
	chatGPTAccount := "Not connected"
	if m.AuthStore != nil {
		chatGPTModel = m.AuthStore.ResolveModel(tutor.ProviderChatGPTPlan)
		if token := m.AuthStore.GetConfig().ChatGPTPlanToken; token != nil && hasChatGPT {
			chatGPTAccount = token.Email
			if chatGPTAccount == "" {
				chatGPTAccount = "Connected"
			}
		}
	}
	prov2Desc := fmt.Sprintf("    Model: %s • Client ID: %s\n    Uses ChatGPT Plus / Pro plan allowance (no API billing). Preview: store=false, stream=true.", chatGPTModel, chatGPTAccount)
	if activeProvider == tutor.ProviderChatGPTPlan {
		providerLines = append(providerLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(prov2+"\n"+prov2Desc))
	} else {
		providerLines = append(providerLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(prov2+"\n"+prov2Desc))
	}

	// 3. Anthropic
	hasAnthropic := (m.AuthStore != nil && m.AuthStore.IsConfigured(tutor.ProviderAnthropic))
	antInfo := "[3] Anthropic Claude (Messages API)"
	if hasAnthropic {
		antInfo += " (Key Set)"
	}
	prov3 := formatStatus(tutor.ProviderAnthropic, hasAnthropic, antInfo)
	antModel := tutor.DefaultAnthropicModel
	if m.AuthStore != nil {
		antModel = m.AuthStore.ResolveModel(tutor.ProviderAnthropic)
	}
	prov3Desc := fmt.Sprintf("    Model: %s\n    Uses ANTHROPIC_API_KEY environment variable or securely stored key.", antModel)
	if activeProvider == tutor.ProviderAnthropic {
		providerLines = append(providerLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(prov3+"\n"+prov3Desc))
	} else {
		providerLines = append(providerLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(prov3+"\n"+prov3Desc))
	}

	// 4. Google Gemini
	hasGemini := (m.AuthStore != nil && m.AuthStore.IsConfigured(tutor.ProviderGoogle))
	gemInfo := "[4] Google Gemini (Gemini API)"
	if hasGemini {
		gemInfo += " (Key Set)"
	}
	prov4 := formatStatus(tutor.ProviderGoogle, hasGemini, gemInfo)
	gemModel := tutor.DefaultGeminiModel
	if m.AuthStore != nil {
		gemModel = m.AuthStore.ResolveModel(tutor.ProviderGoogle)
	}
	prov4Desc := fmt.Sprintf("    Model: %s\n    Uses GEMINI_API_KEY environment variable or securely stored key.", gemModel)
	if activeProvider == tutor.ProviderGoogle {
		providerLines = append(providerLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(prov4+"\n"+prov4Desc))
	} else {
		providerLines = append(providerLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(prov4+"\n"+prov4Desc))
	}

	// 5. OpenAI API
	hasOpenAI := (m.AuthStore != nil && m.AuthStore.IsConfigured(tutor.ProviderOpenAI))
	oaiInfo := "[5] OpenAI API (Chat Completions)"
	if hasOpenAI {
		oaiInfo += " (Key Set)"
	}
	prov5 := formatStatus(tutor.ProviderOpenAI, hasOpenAI, oaiInfo)
	oaiModel := tutor.DefaultOpenAIAPIModel
	if m.AuthStore != nil {
		oaiModel = m.AuthStore.ResolveModel(tutor.ProviderOpenAI)
	}
	prov5Desc := fmt.Sprintf("    Model: %s\n    Uses OPENAI_API_KEY environment variable or securely stored key.", oaiModel)
	if activeProvider == tutor.ProviderOpenAI {
		providerLines = append(providerLines, m.Styles.OptionSelected.Width(contentWidth-4).Render(prov5+"\n"+prov5Desc))
	} else {
		providerLines = append(providerLines, m.Styles.OptionNormal.Width(contentWidth-4).Render(prov5+"\n"+prov5Desc))
	}

	if !m.TutorModelSelectActive && !m.TutorInputActive {
		sections = append(sections, strings.Join(providerLines, "\n"))
	} else {
		// Keep the active interaction at the top of a terminal-sized screen.
		sections = sections[:2]
		if m.TutorAuthNotice != "" {
			sections = append(sections, m.Styles.ExplanationBox.Width(contentWidth).Render(m.TutorAuthNotice))
		}
	}

	// Model Selection modal overlay if active
	if m.TutorModelSelectActive {
		activeProv := tutor.ProviderOffline
		if m.AuthStore != nil {
			activeProv = m.AuthStore.GetConfig().ActiveProvider
		}
		currentActiveModel := ""
		if m.AuthStore != nil {
			currentActiveModel = m.AuthStore.ResolveModel(activeProv)
		}

		var modelLines []string
		modelLines = append(modelLines, m.Styles.CardTitle.Render(fmt.Sprintf("🤖 SELECT MODEL FOR %s", strings.ToUpper(activeProv))))
		modelLines = append(modelLines, m.Styles.HelpText.Render("Navigate [↑/↓/j/k] • Select [Enter/1-9] • Fetch Live [r] • Custom Model [c] • Cancel [Esc]"))

		if len(m.TutorModelList) == 0 {
			modelLines = append(modelLines, m.Styles.OptionNormal.Render("  No models in catalog. Press 'r' to discover from API or 'c' to enter custom model ID."))
		} else {
			for idx, item := range m.TutorModelList {
				prefix := "  "
				if idx == m.TutorModelCursor {
					prefix = "▶ "
				}
				activeBadge := ""
				if item.ID == currentActiveModel {
					activeBadge = " [✓ ACTIVE]"
				}
				quickKey := ""
				if idx < 9 {
					quickKey = fmt.Sprintf("[%d] ", idx+1)
				}
				line := fmt.Sprintf("%s%s%-26s %s%s", prefix, quickKey, item.ID, item.DisplayName, activeBadge)
				if idx == m.TutorModelCursor {
					modelLines = append(modelLines, m.Styles.OptionSelected.Width(contentWidth-6).Render(line))
				} else {
					modelLines = append(modelLines, m.Styles.OptionNormal.Width(contentWidth-6).Render(line))
				}
			}
		}

		if m.TutorCustomModelActive {
			inputPrompt := fmt.Sprintf("✏️  Enter Custom Model ID for %s:\n> %s_\n\n[Press Enter to Save, Esc to Cancel]", strings.ToUpper(activeProv), m.TutorCustomModelBuffer)
			modelLines = append(modelLines, m.Styles.HintBox.Width(contentWidth-6).Render(inputPrompt))
		}

		modelCard := m.Styles.Card.Width(contentWidth).Render(strings.Join(modelLines, "\n"))
		sections = append(sections, modelCard)
	}

	// Key input box if active
	if m.TutorInputActive {
		masked := strings.Repeat("•", len(m.TutorInputBuffer))
		if len(m.TutorInputBuffer) > 6 {
			masked = m.TutorInputBuffer[:3] + strings.Repeat("•", len(m.TutorInputBuffer)-6) + m.TutorInputBuffer[len(m.TutorInputBuffer)-3:]
		}
		inputBox := m.Styles.HintBox.Width(contentWidth).Render(
			fmt.Sprintf("🔑 Enter %s API Key:\n> %s_\n\n[Press Enter to Save, Esc to Cancel]",
				strings.ToUpper(m.TutorInputProvider), masked),
		)
		sections = append(sections, inputBox)
	}

	// Controls footer
	footer := m.Styles.HelpText.Render("[1-5] Choose Provider  [m] Select Model  [r] Discover Models  [a] Disconnect/Change ChatGPT Account  [x] Clear  [Esc/t] Return")
	sections = append(sections, footer)

	// Airline Status Bar in TUTOR mode
	sections = append(sections, m.renderStatusBar(contentWidth, "TUTOR", m.Styles.StatusModeTutor))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderCandidatePreview(contentWidth int) string {
	var sections []string

	sections = append(sections, m.renderHeader(contentWidth))

	title := m.Styles.CardTitle.Render("💡 CANDIDATE REVIEW & PROMOTION (Stage 13)")
	sections = append(sections, title)

	if m.CandidateActive {
		sections = append(sections, m.Styles.TutorLoading.Width(contentWidth).Render(m.loadingGlyph()+" Generating new content from LLM... [Esc cancels]"))
	}
	sections = append(sections, m.Styles.HelpText.Render("Candidates are proposed questions. Review wording, answers, and hints before [a] publishes them for practice."))

	if m.CandidateNotice != "" {
		notice := m.Styles.ExplanationBox.Width(contentWidth).Render(m.CandidateNotice)
		sections = append(sections, notice)
	}

	if len(m.Candidates) == 0 {
		emptyBox := m.Styles.ScenarioBox.Width(contentWidth).Render(
			"No question candidates currently stored in candidate database.\n\n" +
				"Press [n] for a new question from your connected LLM, or [g] for a local variation.\n" +
				"Press [Esc] or [p] to return to practice.",
		)
		sections = append(sections, emptyBox)
		sections = append(sections, m.Styles.HelpText.Render("[n] New from LLM  [g] Local Variation  [p/Esc] Return to Practice"))
		sections = append(sections, m.renderStatusBar(contentWidth, "CANDIDATES", m.Styles.StatusModeHelp))
		return lipgloss.JoinVertical(lipgloss.Left, sections...)
	}

	if m.CandidateIndex < 0 {
		m.CandidateIndex = 0
	}
	if m.CandidateIndex >= len(m.Candidates) {
		m.CandidateIndex = len(m.Candidates) - 1
	}

	cand := m.Candidates[m.CandidateIndex]

	// Header metadata
	statusStr := fmt.Sprintf("Candidate %d of %d | ID: %s | Status: %s | Validation: %s",
		m.CandidateIndex+1, len(m.Candidates), cand.ID, cand.Status, cand.ValidationStatus)
	sections = append(sections, m.Styles.MasteryHeader.Render(statusStr))

	if cand.RejectionReason != "" {
		rejBox := m.Styles.IncorrectBox.Width(contentWidth).Render("⚠️ Validation Failure: " + cand.RejectionReason)
		sections = append(sections, rejBox)
	}

	// 1. Wording (Template and sample)
	sampleAmt := int64(10000)
	if amounts, ok := cand.Parameters["amount_minor_units"]; ok && len(amounts) > 0 {
		sampleAmt = amounts[0]
	}
	sampleFormatted := domain.NewMoney(sampleAmt).FormatDollars()
	instantiated := strings.ReplaceAll(cand.ScenarioTemplate, "${amount_dollars}", sampleFormatted)

	wordingContent := fmt.Sprintf(
		"WORDING (Scenario Template):\n"+
			"Template: %s\n\n"+
			"Sample:   %s",
		cand.ScenarioTemplate, instantiated,
	)
	sections = append(sections, m.Styles.ScenarioBox.Width(contentWidth).Render(wordingContent))

	// 2. Assumptions (Parameters)
	var paramStrs []string
	for pName, pVals := range cand.Parameters {
		var valStrs []string
		for _, v := range pVals {
			if strings.Contains(pName, "minor_units") {
				valStrs = append(valStrs, fmt.Sprintf("%d (%s)", v, domain.NewMoney(v).FormatDollars()))
			} else {
				valStrs = append(valStrs, fmt.Sprintf("%d", v))
			}
		}
		paramStrs = append(paramStrs, fmt.Sprintf("• %s: [%s]", pName, strings.Join(valStrs, ", ")))
	}
	assumpContent := fmt.Sprintf("ASSUMPTIONS (Parameter Space):\n%s", strings.Join(paramStrs, "\n"))
	sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(assumpContent))

	// 3. Engine-Derived Canonical Journal Entry & Equation Effects
	var entryLines []string
	entryLines = append(entryLines, m.Styles.JournalHeader.Render("Engine-Derived Canonical Entry (Local Derivation)"))
	for _, p := range cand.DerivedFixture.Postings {
		accName := p.AccountID
		if m.Catalog != nil && m.Catalog.Has(domain.AccountID(p.AccountID)) {
			acc, _ := m.Catalog.Get(domain.AccountID(p.AccountID))
			accName = acc.Name
		}
		if p.Side == string(domain.SideDebit) {
			line := fmt.Sprintf("Debit   %-32s %12s", accName, sampleFormatted)
			entryLines = append(entryLines, m.Styles.JournalDebit.Render(line))
		} else {
			line := fmt.Sprintf("  Credit  %-30s %12s", accName, sampleFormatted)
			entryLines = append(entryLines, m.Styles.JournalCredit.Render(line))
		}
	}
	isEqBalanced := (cand.DerivedEquation.DeltaAssets == cand.DerivedEquation.DeltaLiabilities.Add(cand.DerivedEquation.DeltaEquity))
	eqLine := fmt.Sprintf("Equation Delta: Assets=%s | Liabilities=%s | Equity=%s | Balanced=%t",
		cand.DerivedEquation.DeltaAssets.FormatDollars(),
		cand.DerivedEquation.DeltaLiabilities.FormatDollars(),
		cand.DerivedEquation.DeltaEquity.FormatDollars(),
		isEqBalanced)
	entryLines = append(entryLines, "\n"+eqLine)
	sections = append(sections, m.Styles.RecapBox.Width(contentWidth).Render(strings.Join(entryLines, "\n")))

	// 4. Concept Tags & Provenance
	provContent := fmt.Sprintf(
		"METADATA & PROVENANCE:\n"+
			"• Concept Tags:   [%s]\n"+
			"• Source:         %s\n"+
			"• Target Family:  %s (Rule Version %d)\n"+
			"• Generated At:   %s",
		strings.Join(cand.Concepts, ", "),
		cand.Provenance.Source,
		cand.FamilyID,
		cand.RuleVersion,
		cand.Provenance.GeneratedAt.UTC().Format(time.RFC3339),
	)
	for _, stage := range []domain.DrillStage{domain.StageIdentifyAccount, domain.StageAccountCategory, domain.StageDirection, domain.StageDebitCredit, domain.StageCounterAccount, domain.StageBalancedEntry, domain.StageEquationEffect} {
		if text, ok := cand.Teaching[stage]; ok {
			provContent += fmt.Sprintf("\n\nProposed teaching [%s] (review required):\nHint: %s\nExplanation: %s", stage, text.Hint, text.Explanation)
		}
	}
	if cand.Explanation.Summary != "" {
		provContent += fmt.Sprintf("\n• Rationale:      %s", cand.Explanation.Summary)
	}
	sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(provContent))

	// Invariant warning
	invNotice := m.Styles.HintBox.Width(contentWidth).Render(
		"🔒 AGENT-CONTRACT INVARIANT:\n" +
			"This candidate is stored outside the active bank. It CANNOT be selected for practice,\n" +
			"alter engine family rules, or affect learner mastery without human semantic review.",
	)
	sections = append(sections, invNotice)

	// Navigation bar
	footer := m.Styles.HelpText.Render("[n] New from LLM  [g] Local Variation  [a] Approve & Publish  [r] Reject  [x] Delete  [j/k or ↑/↓] Navigate  [p/Esc] Close")
	sections = append(sections, footer)

	// Status bar
	sections = append(sections, m.renderStatusBar(contentWidth, "CANDIDATES", m.Styles.StatusModeHelp))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderJournalPractice(contentWidth int) string {
	var sections []string

	sections = append(sections, m.renderHeader(contentWidth))

	title := m.Styles.CardTitle.Render("📝 MULTI-LINE JOURNAL ENTRY PRACTICE (Interactive Builder)")
	sections = append(sections, title)

	// Generic visual baseline notice
	genericNotice := m.Styles.HelpText.Render("[" + domain.GenericVisualNotice + "]")
	sections = append(sections, genericNotice)

	// Scenario box
	scenario := m.JournalScenario
	if scenario == "" {
		scenario = "Select [n] to load a business scenario."
	}
	sections = append(sections, m.Styles.ScenarioBox.Width(contentWidth).Render("Scenario: "+scenario))

	// Journal Postings Table & Classroom Grid
	var tableLines []string
	tableLines = append(tableLines, m.Styles.JournalHeader.Render("Active Classroom Transaction Grid"))

	var totalDr, totalCr domain.Money
	if len(m.JournalLines) == 0 {
		tableLines = append(tableLines, m.Styles.HelpText.Render("  (No journal lines entered yet. Press [a] to add a line)"))
	} else {
		for i, p := range m.JournalLines {
			cursor := "  "
			if i == m.JournalSelectedLine && !m.JournalInputActive {
				cursor = "▶ "
			}
			accName := string(p.AccountID)
			if m.Catalog != nil && m.Catalog.Has(p.AccountID) {
				acc, _ := m.Catalog.Get(p.AccountID)
				accName = acc.Name
			}
			effect := domain.PostingEquationEffect(p, m.Catalog)
			var lineStr string
			if p.Side == domain.SideDebit {
				lineStr = fmt.Sprintf("%sLine %d: [Dr.] %s %s   %s", cursor, i+1, accName, effect, p.Amount.FormatCommas())
				tableLines = append(tableLines, m.Styles.JournalDebit.Render(lineStr))
				totalDr = totalDr.Add(p.Amount)
			} else {
				lineStr = fmt.Sprintf("%sLine %d: [Cr.]   %s %s   %s", cursor, i+1, accName, effect, p.Amount.FormatCommas())
				tableLines = append(tableLines, m.Styles.JournalCredit.Render(lineStr))
				totalCr = totalCr.Add(p.Amount)
			}
		}

		tableLines = append(tableLines, "")
		tableLines = append(tableLines, domain.RenderClassroomGrid(m.JournalLines, m.Catalog, contentWidth-6, false))
	}

	// Balance Status
	isBalanced := (totalDr == totalCr && totalDr.IsPositive())
	var statusBadge string
	if isBalanced {
		statusBadge = m.Styles.BalancedBadge.Render("✓ BALANCED (Dr = Cr)")
	} else if totalDr.IsZero() && totalCr.IsZero() {
		statusBadge = m.Styles.HelpText.Render("[EMPTY]")
	} else {
		diff := totalDr.Sub(totalCr).Abs()
		statusBadge = m.Styles.IncorrectBox.Render(fmt.Sprintf("✗ UNBALANCED (Diff: %s)", diff.FormatDollars()))
	}

	totalsLine := fmt.Sprintf("Total Debits: %s | Total Credits: %s  %s",
		totalDr.FormatDollars(), totalCr.FormatDollars(), statusBadge)
	tableLines = append(tableLines, "\n"+totalsLine)

	sections = append(sections, m.Styles.RecapBox.Width(contentWidth).Render(strings.Join(tableLines, "\n")))

	// Input Editor Modal (when adding line)
	if m.JournalInputActive {
		var inputCard []string
		allAccounts := m.Catalog.All()
		switch m.JournalInputMode {
		case "account":
			inputCard = append(inputCard, m.Styles.CardTitle.Render("STEP 1: SELECT ACCOUNT"))
			if m.JournalAccountIdx >= 0 && m.JournalAccountIdx < len(allAccounts) {
				sel := allAccounts[m.JournalAccountIdx]
				inputCard = append(inputCard, fmt.Sprintf("  ▶ %s (%s, Normal: %s)", sel.Name, sel.Category, sel.NormalSide))
			}
			inputCard = append(inputCard, "\n[↑/k] Up  [↓/j] Down  [Enter] Confirm Account  [Esc] Cancel")
		case "side":
			inputCard = append(inputCard, m.Styles.CardTitle.Render("STEP 2: SELECT DEBIT OR CREDIT"))
			acc := allAccounts[m.JournalAccountIdx]
			sideLabel := "Debit"
			if m.JournalSide == domain.SideCredit {
				sideLabel = "Credit"
			}
			inputCard = append(inputCard, fmt.Sprintf("  Account: %s\n  Side: ▶ %s", acc.Name, strings.ToUpper(sideLabel)))
			inputCard = append(inputCard, "\n[D] Debit  [c] Credit  [Tab] Toggle  [Enter] Confirm Side  [Esc] Back")
		case "amount":
			inputCard = append(inputCard, m.Styles.CardTitle.Render("STEP 3: ENTER MONETARY AMOUNT"))
			acc := allAccounts[m.JournalAccountIdx]
			inputCard = append(inputCard, fmt.Sprintf("  Account: %s (%s)\n  Amount:  $ %s █", acc.Name, m.JournalSide, m.JournalAmountBuffer))
			inputCard = append(inputCard, "\nType dollar amount (e.g. 1500 or 1500.50). [Enter] Add Line  [Esc] Back")
		}
		sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(strings.Join(inputCard, "\n")))
	}

	// Feedback / Notice box
	if m.JournalNotice != "" && !m.JournalInputActive {
		if strings.HasPrefix(m.JournalNotice, "✓") {
			sections = append(sections, m.Styles.CorrectBox.Width(contentWidth).Render(m.JournalNotice))
		} else if strings.HasPrefix(m.JournalNotice, "✗") || strings.HasPrefix(m.JournalNotice, "⚠️") {
			sections = append(sections, m.Styles.IncorrectBox.Width(contentWidth).Render(m.JournalNotice))
		} else {
			sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(m.JournalNotice))
		}
	}

	// T-Account and Equation Reconciliation Views (when submitted and correct)
	if m.JournalReconciliation != nil && m.JournalReconciliation.Reconciled {
		// T-Accounts side-by-side
		var tCards []string
		for _, ta := range m.JournalReconciliation.Ledger.AccountsInOrder() {
			tCards = append(tCards, m.Styles.TAccountBox.Render(ta.Render()))
		}
		if len(tCards) > 0 {
			sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, tCards...))
		}

		// Equation effects
		eqSummary := fmt.Sprintf(
			"ACCOUNTING EQUATION RECONCILIATION:\n"+
				"• Δ Assets:      %12s\n"+
				"• Δ Liabilities: %12s\n"+
				"• Δ Equity:      %12s\n"+
				"• Reconciled:    %s = %s + %s (100%% T-account & Journal tie-out)",
			m.JournalReconciliation.DeltaAssets.FormatDollars(),
			m.JournalReconciliation.DeltaLiabilities.FormatDollars(),
			m.JournalReconciliation.DeltaEquity.FormatDollars(),
			m.JournalReconciliation.DeltaAssets.FormatDollars(),
			m.JournalReconciliation.DeltaLiabilities.FormatDollars(),
			m.JournalReconciliation.DeltaEquity.FormatDollars(),
		)
		sections = append(sections, m.Styles.EquationBox.Width(contentWidth).Render(eqSummary))
	}

	// Instructions footer
	footer := m.Styles.HelpText.Render("[a] Add Line  [D/c] Change Side  [x] Delete Line  [s/Enter] Submit Entry  [n] Next  [C] Clear  [Esc] Back")
	sections = append(sections, footer)

	// Airline Status Bar in JOURNAL mode
	sections = append(sections, m.renderStatusBar(contentWidth, "JOURNAL", m.Styles.StatusModeJournal))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderStatements(contentWidth int) string {
	var sections []string

	if m.StatementsReport == nil {
		c := statements.CanonicalCasePioneerConsulting()
		m.StatementsReport, _ = statements.BuildAccountingCycleReport(c, m.Catalog, m.Engine)
	}

	title := m.Styles.HelpTitle.Render("FINANCIAL STATEMENTS & COMPREHENSIVE CASE REPORT")
	sections = append(sections, title)
	sections = append(sections, m.Styles.HelpText.Render(fmt.Sprintf("[%s]", domain.GenericVisualNotice)))

	if m.StatementsReport != nil {
		fullText := m.StatementsReport.FormatReport()
		sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(fullText))
	}

	footer := m.Styles.HelpText.Render(fmt.Sprintf("[↑/k, ↓/j] Scroll (Line %d)  [g/G] Top/Bottom  [Esc/q] Return to Drill", m.StatementsScroll+1))
	sections = append(sections, footer)

	sections = append(sections, m.renderStatusBar(contentWidth, "STATEMENTS", m.Styles.StatusModeStatements))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderExam(contentWidth int) string {
	var sections []string

	if m.ExamRunner == nil {
		return "Exam session initializing..."
	}

	eq, stageAns, err := m.ExamRunner.CurrentQuestion()
	if err != nil {
		return "Exam finished or error loading question."
	}

	// 1. Top Exam Header with timer and question counter
	qNum := eq.QuestionIndex + 1
	totQ := m.ExamRunner.TotalQuestions
	stepNum := eq.CurrentStage + 1
	totSteps := len(eq.Stages)

	var timerStr string
	if m.ExamRunner.TimeLimit > 0 {
		rem := m.ExamRunner.Remaining()
		mins := int(rem.Minutes())
		secs := int(rem.Seconds()) % 60
		timerStr = fmt.Sprintf("⏳ %02d:%02d remaining", mins, secs)
	} else {
		mins := m.ExamRunner.ElapsedSeconds / 60
		secs := m.ExamRunner.ElapsedSeconds % 60
		timerStr = fmt.Sprintf("⏱ Elapsed: %02d:%02d", mins, secs)
	}

	headerBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.Styles.StatusModeExam.Render("🔒 EXAM MODE"),
		" ",
		m.Styles.StatusInfo.Render(fmt.Sprintf(" Question %d/%d • Step %d/%d ", qNum, totQ, stepNum, totSteps)),
		" ",
		m.Styles.ExamTimerPill.Render(fmt.Sprintf(" %s ", timerStr)),
	)
	sections = append(sections, headerBar)

	// 2. Strict assessment conditions reminder
	conditionsNotice := m.Styles.HelpText.Render("Test conditions: Hints, reference cheatsheet, and immediate feedback are withheld until finish.")
	sections = append(sections, conditionsNotice)

	if m.ExamNotice != "" {
		sections = append(sections, m.Styles.ExamWarningBox.Render(m.ExamNotice))
	}

	// 3. Scenario Prompt Card
	scenarioText := fmt.Sprintf("Question %d Scenario:\n\n%s", qNum, eq.Instance.PromptText)
	sections = append(sections, m.Styles.ScenarioBox.Width(contentWidth).Render(scenarioText))

	// 4. Stage Question Prompt
	promptText := fmt.Sprintf("Step %d: %s", stepNum, stageAns.Prompt)
	sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(promptText))

	// 5. Options
	var optionLines []string
	for i, opt := range stageAns.Options {
		indicator := "  "
		style := m.Styles.OptionNormal
		if i == m.SelectedOptionIndex {
			indicator = "▶ "
			style = m.Styles.OptionSelected
		}
		optLine := fmt.Sprintf("%s[%d] %s", indicator, i+1, opt.Text)
		optionLines = append(optionLines, style.Width(contentWidth).Render(optLine))
	}
	sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, optionLines...))

	// 6. Navigation Footer
	footer := m.Styles.HelpText.Render("[1-4] Select  [Enter] Submit Answer   [↑/↓, j/k] Navigate Options   [q] Save & Exit Exam")
	sections = append(sections, footer)

	// 7. Airline Status Bar
	sections = append(sections, m.renderStatusBar(contentWidth, "EXAM", m.Styles.StatusModeExam))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderExamSummary(contentWidth int) string {
	var sections []string

	if m.ExamReport == nil {
		return "Generating exam report..."
	}

	r := m.ExamReport

	// 1. Completion Header
	statusTitle := "🎉 EXAM ASSESSMENT COMPLETE"
	if r.TimedOut {
		statusTitle = "⏱ EXAM TIME EXPIRED"
	}
	sections = append(sections, m.Styles.CardTitle.Render(statusTitle))

	// 2. Score Banner
	scoreText := fmt.Sprintf(
		"• Final Score:    %d / %d items (%.1f%%)\n"+
			"• Letter Grade:   Grade %s\n"+
			"• Total Duration: %s",
		r.CorrectItems, r.TotalItems, r.Score,
		r.Grade,
		r.Duration.Round(time.Second),
	)
	if r.TimeLimitSeconds > 0 {
		scoreText += fmt.Sprintf(" (Time Limit: %s)", (time.Duration(r.TimeLimitSeconds) * time.Second).String())
	}
	sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(scoreText))

	// 3. Tabular Skill Breakdown
	var skillLines []string
	skillLines = append(skillLines, m.Styles.MasteryHeader.Render("1. SKILL & CONCEPT PERFORMANCE BREAKDOWN"))
	skillLines = append(skillLines, fmt.Sprintf("%-38s %-12s %-10s %-12s", "Skill / Concept", "Score", "Accuracy", "Status"))
	skillLines = append(skillLines, strings.Repeat("─", 74))

	for _, sk := range r.Skills {
		scoreStr := fmt.Sprintf("%d / %d", sk.Correct, sk.Total)
		accStr := fmt.Sprintf("%.1f%%", sk.Accuracy)
		badge := sk.Status
		if sk.Status == "Strong" {
			badge = m.Styles.CorrectBox.Render("Strong")
		} else if sk.Status == "Critical" {
			badge = m.Styles.IncorrectBox.Render("Critical")
		}
		skillLines = append(skillLines, fmt.Sprintf("%-38s %-12s %-10s %s", sk.ConceptName, scoreStr, accStr, badge))
	}
	sections = append(sections, m.Styles.Card.Width(contentWidth).Render(strings.Join(skillLines, "\n")))

	// 4. Diagnostic Error Pattern Analysis
	var errLines []string
	errLines = append(errLines, m.Styles.MasteryHeader.Render("2. DIAGNOSTIC ERROR MISCONCEPTIONS"))
	if len(r.Errors) == 0 {
		errLines = append(errLines, "✓ Outstanding! Zero diagnostic distractor errors recorded during this exam.")
	} else {
		for i, errItem := range r.Errors {
			errLines = append(errLines, fmt.Sprintf("[%d] %s (%d occurrences)", i+1, errItem.Misconception, errItem.Count))
			errLines = append(errLines, fmt.Sprintf("    Remediation: %s", errItem.Remediation))
		}
	}
	sections = append(sections, m.Styles.Card.Width(contentWidth).Render(strings.Join(errLines, "\n")))

	// 5. Question-by-Question Audit Review
	if len(r.Reviews) > 0 {
		if m.ExamReviewIndex < 0 {
			m.ExamReviewIndex = 0
		}
		if m.ExamReviewIndex >= len(r.Reviews) {
			m.ExamReviewIndex = len(r.Reviews) - 1
		}
		rev := r.Reviews[m.ExamReviewIndex]

		var revLines []string
		revLines = append(revLines, m.Styles.MasteryHeader.Render(fmt.Sprintf("3. QUESTION AUDIT REVIEW (Item %d of %d — use [n]/[p] to navigate)", m.ExamReviewIndex+1, len(r.Reviews))))

		statusMark := m.Styles.CorrectBox.Render("✓ CORRECT")
		if !rev.IsCorrect {
			statusMark = m.Styles.IncorrectBox.Render("✗ INCORRECT")
		}
		revLines = append(revLines, fmt.Sprintf("Question %d: %s (Stage: %s)", rev.QuestionIndex+1, statusMark, rev.Stage))
		revLines = append(revLines, fmt.Sprintf("Scenario: %s", rev.PromptText))
		revLines = append(revLines, fmt.Sprintf("Step:     %s", rev.StagePrompt))

		selText := "(none)"
		if rev.SelectedOption != nil {
			selText = fmt.Sprintf("[%s] %s", rev.SelectedOption.Label, rev.SelectedOption.Text)
		}
		revLines = append(revLines, fmt.Sprintf("Your Choice:    %s", selText))
		revLines = append(revLines, fmt.Sprintf("Correct Answer: [%s] %s", rev.CorrectOption.Label, rev.CorrectOption.Text))

		if rev.ErrorTag != "" {
			misc, _ := exam.DistractorExplanation(rev.ErrorTag)
			revLines = append(revLines, fmt.Sprintf("Misconception:  %s (%s)", rev.ErrorTag, misc))
		}
		revLines = append(revLines, fmt.Sprintf("Explanation:    %s", rev.Explanation))

		if len(rev.Postings) > 0 {
			revLines = append(revLines, "Balanced Journal Entry:")
			for _, p := range rev.Postings {
				if p.Side == domain.SideDebit {
					revLines = append(revLines, fmt.Sprintf("  Dr. %-24s %s", p.AccountID, p.Amount.FormatDollars()))
				} else {
					revLines = append(revLines, fmt.Sprintf("      Cr. %-20s %s", p.AccountID, p.Amount.FormatDollars()))
				}
			}
		}

		sections = append(sections, m.Styles.Card.Width(contentWidth).Render(strings.Join(revLines, "\n")))
	}

	// Footer instructions
	footer := m.Styles.HelpText.Render("[r] Take New Exam   [n/p, j/k] Browse Question Reviews   [s] Mastery   [q/Esc] Return to Drill")
	sections = append(sections, footer)

	sections = append(sections, m.renderStatusBar(contentWidth, "EXAM RESULTS", m.Styles.StatusModeExam))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderExamResumePrompt(contentWidth int) string {
	var sections []string
	if m.ExamNotice != "" {
		sections = append(sections, m.Styles.ExamWarningBox.Render(m.ExamNotice))
	}

	sections = append(sections, m.Styles.CardTitle.Render("⚠️  INTERRUPTED EXAM SESSION DETECTED"))

	var details string
	if m.InterruptedExam != nil {
		elapsedMin := m.InterruptedExam.ElapsedSeconds / 60
		elapsedSec := m.InterruptedExam.ElapsedSeconds % 60
		details = fmt.Sprintf(
			"An exam session was previously interrupted:\n\n"+
				"• Session ID:         %s\n"+
				"• Started At:         %s\n"+
				"• Questions:          %d total\n"+
				"• Elapsed Time:       %02d:%02d\n"+
				"• Attempts Recorded:  %d\n\n"+
				"Would you like to resume your interrupted exam, or abandon it and start fresh?",
			m.InterruptedExam.ID,
			m.InterruptedExam.StartedAt.Format("2006-01-02 15:04:05 UTC"),
			m.InterruptedExam.TotalQuestions,
			elapsedMin, elapsedSec,
			m.InterruptedExam.TotalAttempts,
		)
	} else {
		details = "An interrupted exam session was detected. Resume or start fresh?"
	}
	sections = append(sections, m.Styles.PromptBox.Width(contentWidth).Render(details))

	actions := m.Styles.Card.Width(contentWidth).Render(
		"[r / Enter]  RESUME INTERRUPTED EXAM (continue where you left off)\n" +
			"[a / n]      ABANDON INTERRUPTED EXAM (start fresh exam)\n" +
			"[q / Esc]    RETURN TO PRACTICE DRILL",
	)
	sections = append(sections, actions)

	sections = append(sections, m.renderStatusBar(contentWidth, "EXAM RESUME", m.Styles.StatusModeExam))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

package tui

import "github.com/charmbracelet/lipgloss"

// Styles holds all Lip Gloss styles used across the TUI.
type Styles struct {
	Header         lipgloss.Style
	HeaderTitle    lipgloss.Style
	HeaderSubtitle lipgloss.Style
	ProgressBar    lipgloss.Style
	ProgressFill   lipgloss.Style
	LogoBanner     lipgloss.Style

	Card        lipgloss.Style
	CardTitle   lipgloss.Style
	ScenarioBox lipgloss.Style
	PromptBox   lipgloss.Style

	OptionNormal   lipgloss.Style
	OptionSelected lipgloss.Style
	OptionLetter   lipgloss.Style
	OptionText     lipgloss.Style

	HintBox        lipgloss.Style
	ExplanationBox lipgloss.Style
	TutorLoading   lipgloss.Style
	TutorBadge     lipgloss.Style
	CorrectBox     lipgloss.Style
	IncorrectBox   lipgloss.Style
	RetryBox       lipgloss.Style

	RecapBox      lipgloss.Style
	JournalHeader lipgloss.Style
	JournalDebit  lipgloss.Style
	JournalCredit lipgloss.Style
	BalancedBadge lipgloss.Style
	EquationBox   lipgloss.Style

	// T-Account Visualization
	TAccountBox     lipgloss.Style
	TAccountTitle   lipgloss.Style
	TAccountHeader  lipgloss.Style
	TAccountDebit   lipgloss.Style
	TAccountCredit  lipgloss.Style
	TAccountDivider lipgloss.Style

	// Mastery View
	MasteryHeader lipgloss.Style
	MasteryRow    lipgloss.Style
	MasteryBadge  lipgloss.Style
	MasteryScore  lipgloss.Style
	ScaffoldBadge lipgloss.Style

	// Airline / Lualine Status Bar
	StatusLine           lipgloss.Style
	StatusModeDrill      lipgloss.Style
	StatusModeHelp       lipgloss.Style
	StatusModeMstr       lipgloss.Style
	StatusModeRecap      lipgloss.Style
	StatusModeFbk        lipgloss.Style
	StatusModeTutor      lipgloss.Style
	StatusModeJournal    lipgloss.Style
	StatusModeStatements lipgloss.Style
	StatusModeExam       lipgloss.Style
	ExamTimerPill        lipgloss.Style
	ExamWarningBox       lipgloss.Style
	StatusInfo           lipgloss.Style
	StatusStreak         lipgloss.Style
	StatusStreakHot      lipgloss.Style
	StatusKeyHints       lipgloss.Style

	// Help / Accounting Reference Screen
	HelpContainer lipgloss.Style
	HelpTitle     lipgloss.Style
	HelpSection   lipgloss.Style
	HelpGridBox   lipgloss.Style
	HelpKeyCol    lipgloss.Style
	HelpDescCol   lipgloss.Style

	// Micro-delights
	StreakCelebration lipgloss.Style
	Footer            lipgloss.Style
	KeyBadge          lipgloss.Style
	KeyDesc           lipgloss.Style
	HelpText          lipgloss.Style
	ScrollIndicator   lipgloss.Style
}

// DefaultStyles creates a vibrant, high-contrast, modern terminal aesthetic.
func DefaultStyles() Styles {
	var s Styles

	// Curated modern palette (Deep Slate, Sky Blue, Cyber Emerald, Amber Glow, Coral Red, Electric Purple)
	skyBlue := lipgloss.Color("#0284c7")
	skyCyan := lipgloss.Color("#38bdf8")
	emerald := lipgloss.Color("#10b981")
	emeraldDark := lipgloss.Color("#047857")
	emeraldBg := lipgloss.Color("#064e3b")
	amber := lipgloss.Color("#f59e0b")
	amberDark := lipgloss.Color("#b45309")
	amberBg := lipgloss.Color("#78350f")
	flameOrange := lipgloss.Color("#ff6b00")
	coralRed := lipgloss.Color("#ef4444")
	purple := lipgloss.Color("#8b5cf6")
	purpleDark := lipgloss.Color("#6d28d9")
	purpleBg := lipgloss.Color("#4c1d95")
	slateMuted := lipgloss.Color("#94a3b8")
	slateBorder := lipgloss.Color("#475569")
	slateDark := lipgloss.Color("#0f172a")
	slateCard := lipgloss.Color("#1e293b")
	white := lipgloss.Color("#f8fafc")

	s.LogoBanner = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		Background(slateDark).
		Padding(0, 1)

	s.Header = lipgloss.NewStyle().
		MarginBottom(1).
		Padding(0, 1)

	s.HeaderTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		PaddingRight(2)

	s.HeaderSubtitle = lipgloss.NewStyle().
		Foreground(slateMuted)

	s.ProgressBar = lipgloss.NewStyle().
		Foreground(slateBorder)

	s.ProgressFill = lipgloss.NewStyle().
		Foreground(skyCyan).
		Bold(true)

	s.Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(slateBorder).
		Padding(1, 2).
		MarginBottom(1)

	s.CardTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		MarginBottom(1)

	s.ScenarioBox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(skyCyan).
		PaddingLeft(2).
		MarginBottom(1).
		Foreground(white)

	s.PromptBox = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		MarginBottom(1)

	s.OptionNormal = lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(0).
		Foreground(white)

	s.OptionSelected = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(skyBlue).
		Padding(0, 1).
		MarginBottom(0)

	s.OptionLetter = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan)

	s.OptionText = lipgloss.NewStyle().
		PaddingLeft(1)

	s.HintBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(amber).
		Foreground(amber).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.ExplanationBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Foreground(white).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.TutorLoading = lipgloss.NewStyle().
		Bold(true).
		Foreground(amber).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.TutorBadge = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(purpleDark).
		Padding(0, 1)

	s.CorrectBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(emerald).
		Foreground(emerald).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.IncorrectBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(coralRed).
		Foreground(coralRed).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.RetryBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(amber).
		Foreground(amber).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	s.RecapBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(emerald).
		Padding(1, 2).
		MarginBottom(1)

	s.JournalHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		Underline(true).
		MarginBottom(1)

	s.JournalDebit = lipgloss.NewStyle().
		Foreground(white)

	s.JournalCredit = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(white)

	s.BalancedBadge = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(emeraldDark).
		Padding(0, 1)

	s.EquationBox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(slateBorder).
		PaddingTop(1).
		MarginTop(1).
		Foreground(white)

	// T-Account visualization styles
	s.TAccountBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(slateBorder).
		Padding(0, 1).
		MarginRight(1).
		MarginBottom(1)

	s.TAccountTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		Align(lipgloss.Center)

	s.TAccountHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(slateMuted)

	s.TAccountDebit = lipgloss.NewStyle().
		Foreground(emerald).
		Bold(true)

	s.TAccountCredit = lipgloss.NewStyle().
		Foreground(skyCyan).
		Bold(true)

	s.TAccountDivider = lipgloss.NewStyle().
		Foreground(slateBorder)

	// Airline / Lualine Statusline styles
	s.StatusLine = lipgloss.NewStyle().
		MarginTop(1).
		Background(slateCard)

	s.StatusModeDrill = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(skyBlue).
		Padding(0, 1)

	s.StatusModeHelp = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(emeraldDark).
		Padding(0, 1)

	s.StatusModeMstr = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(purpleDark).
		Padding(0, 1)

	s.StatusModeRecap = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(emeraldDark).
		Padding(0, 1)

	s.StatusModeFbk = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(amberDark).
		Padding(0, 1)

	s.StatusModeTutor = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(purpleDark).
		Padding(0, 1)

	s.StatusModeJournal = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(skyBlue).
		Padding(0, 1)

	s.StatusModeStatements = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(emeraldDark).
		Padding(0, 1)

	s.StatusModeExam = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(lipgloss.Color("#e11d48")).
		Padding(0, 1)

	s.ExamTimerPill = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(lipgloss.Color("#b45309")).
		Padding(0, 1)

	s.ExamWarningBox = lipgloss.NewStyle().
		Foreground(amber).
		Padding(0, 1)

	s.StatusInfo = lipgloss.NewStyle().
		Foreground(slateMuted).
		Background(slateCard).
		Padding(0, 1)

	s.StatusStreak = lipgloss.NewStyle().
		Bold(true).
		Foreground(amber).
		Background(amberBg).
		Padding(0, 1)

	s.StatusStreakHot = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(flameOrange).
		Padding(0, 1)

	s.StatusKeyHints = lipgloss.NewStyle().
		Foreground(slateMuted).
		Background(slateCard).
		Padding(0, 1)

	// Help / Reference Screen styles
	s.HelpContainer = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(skyCyan).
		Padding(1, 2).
		MarginBottom(1)

	s.HelpTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		Align(lipgloss.Center).
		MarginBottom(1)

	s.HelpSection = lipgloss.NewStyle().
		Bold(true).
		Foreground(amber).
		MarginTop(1).
		MarginBottom(1)

	s.HelpGridBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(slateBorder).
		Padding(0, 1).
		MarginBottom(1)

	s.HelpKeyCol = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		PaddingRight(2)

	s.HelpDescCol = lipgloss.NewStyle().
		Foreground(white)

	// Micro-delights & Celebration
	s.StreakCelebration = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(flameOrange).
		Padding(0, 2).
		MarginBottom(1)

	s.MasteryHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyCyan).
		Underline(true).
		MarginBottom(1)

	s.MasteryRow = lipgloss.NewStyle().
		Padding(0, 1)

	s.MasteryBadge = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	s.ScaffoldBadge = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(purpleDark).
		Padding(0, 1)

	s.MasteryScore = lipgloss.NewStyle().
		Bold(true)

	s.Footer = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(slateBorder).
		PaddingTop(1).
		MarginTop(1).
		Foreground(slateMuted)

	s.KeyBadge = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(slateBorder).
		Padding(0, 1)

	s.KeyDesc = lipgloss.NewStyle().
		Foreground(slateMuted).
		Padding(0, 1)

	s.HelpText = lipgloss.NewStyle().
		Foreground(slateMuted)

	s.ScrollIndicator = lipgloss.NewStyle().
		Bold(true).
		Foreground(amber).
		Padding(0, 1)

	_ = emeraldBg
	_ = purple
	_ = purpleBg

	return s
}

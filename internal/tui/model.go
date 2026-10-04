package tui

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/statements"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

// UIState represents the active screen or mode of the TUI.
type UIState int

const (
	StateDrill UIState = iota
	StateFeedback
	StateRecap
	StateMastery
	StateHelp
	StateSessionComplete
	StateTutorConfig
	StateCandidatePreview
	StateJournalPractice
	StateStatements
	StateExam
	StateExamSummary
	StateExamResumePrompt
	StateIntro
	StateQuitting
	StateTitle
	StateSavedExplanations
)

// StageHistoryItem records the learner's response and explanation for a completed stage.
type StageHistoryItem struct {
	StageIndex       int
	StageKey         domain.DrillStage
	SelectedOptionID string
	IsCorrect        bool
	Explanation      string
	Hint             string
	ErrorTag         string
	AssistanceLevel  domain.AssistanceLevel
}

// DrillQuestionState preserves the interactive state of a question in a practice session.
type DrillQuestionState struct {
	Index               int
	Instance            *domain.QuestionInstance
	Session             *drill.SessionState
	State               UIState
	SelectedOptionIndex int
	LastFeedback        *drill.SubmitFeedback
	CurrentHint         string
	RenderedHint        string
	RenderedHintRaw     string
	RenderedHintWidth   int
	ShowHint            bool
	TutorKind           string
	TutorResponse       *tutor.Response
	RecapScroll         int
	ViewingStageIndex   int
	StageHistory        map[int]*StageHistoryItem
}

// Model represents the Bubble Tea application state.
type Model struct {
	State         UIState
	PreviousState UIState

	// Core dependencies
	DB        *storage.DB
	Catalog   *domain.AccountCatalog
	Engine    *engine.Engine
	Generator *drill.Generator
	Scheduler *mastery.Scheduler
	Questions []bank.QuestionJSON

	// Session state
	SessionID            string
	TotalQuestions       int
	CurrentQuestionIndex int
	CurrentInstance      *domain.QuestionInstance
	Session              *drill.SessionState
	SelectedOptionIndex  int
	LastFeedback         *drill.SubmitFeedback
	CurrentHint          string
	renderedHint         string
	renderedHintRaw      string
	renderedHintWidth    int
	ShowHint             bool
	QuestionHistory      []*DrillQuestionState
	ViewingStageIndex    int
	StageHistory         map[int]*StageHistoryItem

	// Tutor integration
	Tutor                 tutor.Tutor
	TutorTimeout          time.Duration
	TutorActive           bool
	TutorCancel           context.CancelFunc
	TutorKind             string // "hint" or "explain"
	TutorResponse         *tutor.Response
	PendingExplanation    *storage.SavedExplanation
	ExplanationSavePrompt bool
	ExplanationSaving     bool
	ExplanationSaveError  string
	ExplanationLeaveKey   tea.KeyMsg
	SavedExplanations     []storage.SavedExplanation
	SavedExplanationIndex int
	TitleFrame            int
	TutorError            string
	LoadingFrame          int
	TutorRequestID        int
	CandidateActive       bool
	CandidateCancel       context.CancelFunc
	CandidateRequestID    int

	// Tutor settings state (Stage 11 & Stage 18)
	AuthStore          *tutor.AuthStore
	ModelCache         *tutor.ModelCache
	TutorInputActive   bool
	TutorInputBuffer   string
	TutorInputProvider string
	TutorAuthNotice    string
	OAuthActiveFlow    *tutor.OAuthFlow

	// Model Selection & Discovery (Stage 18)
	TutorModelSelectActive bool
	TutorModelList         []tutor.ModelInfo
	TutorModelCursor       int
	TutorCustomModelActive bool
	TutorCustomModelBuffer string
	TutorFetchingModels    bool
	TutorOAuthInputActive  bool
	TutorOAuthInputBuffer  string

	// Candidate Preview state (Stage 12)
	Candidates      []candidate.CandidateQuestion
	CandidateIndex  int
	CandidateNotice string

	// Journal Entry Practice state (Stage 14)
	JournalQuestionIndex  int
	JournalScenario       string
	JournalFamilyID       string
	JournalEvent          engine.TransactionEvent
	JournalCanonicalEntry domain.Entry
	JournalLines          []domain.Posting
	JournalSelectedLine   int
	JournalFeedback       *engine.SemanticEvaluation
	JournalReconciliation *domain.TransactionReconciliation
	JournalInputActive    bool
	JournalInputMode      string // "account", "side", "amount"
	JournalAccountIdx     int
	JournalSide           domain.Side
	JournalAmountBuffer   string
	JournalNotice         string

	// Financial Statements & Case Study state (Stage 15)
	StatementsReport *statements.AccountingCycleReport
	StatementsScroll int

	// Recap state (Stage 17)
	RecapScroll int
	PageScroll  int

	// Startup Intro Animation
	Intro       *IntroState
	IntroReturn UIState // screen to restore when leaving the arcade

	// Exam Mode state (Stage 16)
	ExamRunner      *exam.ExamRunner
	ExamReport      *exam.ExamReport
	ExamTimeLimit   time.Duration
	ExamNotice      string
	ExamReviewIndex int
	InterruptedExam *exam.ExamSessionRecord

	// Session metrics & Gamification
	SessionAttempts    int
	FirstTrySuccesses  int
	QuestionsCompleted int
	CurrentStreak      int
	BestStreak         int

	// Randomness and Clock
	RNG   *rand.Rand
	Seed  int64
	Clock mastery.Clock

	// UI layout & styling
	Styles Styles
	Width  int
	Height int
}

// Config holds configuration parameters for initializing the TUI.
type Config struct {
	DB             *storage.DB
	Catalog        *domain.AccountCatalog
	Questions      []bank.QuestionJSON
	TotalQuestions int
	Intensity      mastery.SessionIntensity
	Seed           int64
	Clock          mastery.Clock
	Tutor          tutor.Tutor
	TutorTimeout   time.Duration
	AuthStore      *tutor.AuthStore
	ModelCache     *tutor.ModelCache
	TutorModel     string
	OAuthClientID  string
	ExamMode       bool
	ExamTimeLimit  time.Duration
	ResumeExam     bool
	Intro          bool
}

// NewModel constructs and initializes a new TUI Model.
func NewModel(cfg Config) (*Model, error) {
	if cfg.Catalog == nil {
		return nil, fmt.Errorf("catalog cannot be nil")
	}
	if len(cfg.Questions) == 0 {
		return nil, fmt.Errorf("question bank cannot be empty")
	}
	if cfg.TotalQuestions <= 0 {
		cfg.TotalQuestions = 10
	}
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
	if cfg.Clock == nil {
		cfg.Clock = mastery.RealClock{}
	}
	if cfg.AuthStore == nil {
		cfg.AuthStore, _ = tutor.NewAuthStore("")
	}
	modelCache := cfg.ModelCache
	if modelCache == nil {
		cachePath, err := tutor.DefaultModelCachePath()
		if err == nil {
			modelCache, _ = tutor.NewModelCache(cachePath)
		}
	}
	if cfg.TutorModel != "" && cfg.AuthStore != nil {
		active := cfg.AuthStore.GetConfig().ActiveProvider
		if active != tutor.ProviderOffline {
			_ = cfg.AuthStore.SetSelectedModel(active, cfg.TutorModel)
		}
	}
	if cfg.OAuthClientID != "" && cfg.AuthStore != nil {
		_ = cfg.AuthStore.SetChatGPTClientID(cfg.OAuthClientID)
	}
	if cfg.TutorTimeout <= 0 {
		cfg.TutorTimeout = 3 * time.Second
	}

	if cfg.Tutor == nil {
		cfg.Tutor = tutor.BuildTutor(tutor.FactoryOptions{
			AuthStore: cfg.AuthStore,
			Timeout:   cfg.TutorTimeout,
		})
	}

	rng := rand.New(rand.NewSource(cfg.Seed))
	eng := engine.NewEngine(cfg.Catalog)
	gen := drill.NewGenerator(cfg.Catalog, eng)
	sched := mastery.NewScheduler(cfg.Clock, rng)
	if cfg.Intensity != "" {
		sched.SetIntensity(cfg.Intensity)
	}

	sessID := fmt.Sprintf("sess-%d", time.Now().Unix())

	m := &Model{
		State:          StateDrill,
		PreviousState:  StateDrill,
		DB:             cfg.DB,
		Catalog:        cfg.Catalog,
		Engine:         eng,
		Generator:      gen,
		Scheduler:      sched,
		Questions:      cfg.Questions,
		Tutor:          cfg.Tutor,
		TutorTimeout:   cfg.TutorTimeout,
		AuthStore:      cfg.AuthStore,
		ModelCache:     modelCache,
		SessionID:      sessID,
		TotalQuestions: cfg.TotalQuestions,
		RNG:            rng,
		Seed:           cfg.Seed,
		Clock:          cfg.Clock,
		Styles:         DefaultStyles(),
		Width:          80,
		Height:         24,
	}

	// Persist session if DB is provided
	if m.DB != nil {
		sessRecord := storage.SessionRecord{
			ID:        sessID,
			Mode:      "drill",
			StartedAt: m.Clock.Now(),
		}
		if err := m.DB.SaveSession(sessRecord); err != nil {
			return nil, fmt.Errorf("failed to save initial session: %w", err)
		}
	}

	// Load first question
	if err := m.loadNextQuestion(); err != nil {
		return nil, fmt.Errorf("failed to load initial question: %w", err)
	}

	m.initJournalPractice()

	// Stage 16: Exam Mode startup
	if cfg.ExamMode {
		m.ExamTimeLimit = cfg.ExamTimeLimit
		if m.DB != nil {
			interrupted, err := m.DB.GetLatestInterruptedExamSession()
			if err == nil && interrupted != nil {
				if cfg.ResumeExam {
					atts, err := m.DB.GetExamAttempts(interrupted.ID)
					if err == nil {
						snapshots, err := m.DB.GetExamQuestionSnapshots(interrupted.ID)
						var resumed *exam.ExamRunner
						if err == nil {
							resumed, err = exam.ResumeExamSnapshots(*interrupted, atts, snapshots)
						}
						if err == nil {
							m.ExamRunner = resumed
							m.State = StateExam
							return m, nil
						}
					}
				}
				m.InterruptedExam = interrupted
				m.State = StateExamResumePrompt
				return m, nil
			}
		}
		_ = m.startNewExam(cfg.ExamTimeLimit)
		m.State = StateExam
	}

	if cfg.Intro {
		m.PreviousState = m.State
		m.IntroReturn = m.State
		m.State = StateIntro
		m.initIntro()
	}

	return m, nil
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	if m.State == StateTitle {
		return titleTick()
	}
	if m.State == StateIntro {
		return introTick()
	}
	if m.State == StateExam {
		return tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })
	}
	return nil
}

// examTickMsg represents a 1-second timer tick in Exam Mode.
type examTickMsg struct{}

// tutorResponseMsg represents an asynchronous response from a tutor provider.
type tutorResponseMsg struct {
	ID   int
	Kind string
	Resp tutor.Response
	Err  error
}

// oauthCompleteMsg represents completion of OAuth browser authorization.
type oauthCompleteMsg struct {
	Token *tutor.OAuthToken
	Err   error
}

// modelsFetchedMsg represents completion of dynamic model discovery.
type modelsFetchedMsg struct {
	Provider string
	Models   []tutor.ModelInfo
	Err      error
}

// ExportedModelsFetchedMsg creates a modelsFetchedMsg for testing model discovery handlers.
func ExportedModelsFetchedMsg(provider string, models []tutor.ModelInfo, err error) tea.Msg {
	return modelsFetchedMsg{
		Provider: provider,
		Models:   models,
		Err:      err,
	}
}

// ExportedIntroTickMsg creates an introTickMsg for testing startup animation ticks.
func ExportedIntroTickMsg() tea.Msg {
	return introTickMsg{}
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case titleTickMsg:
		if m.State != StateTitle {
			return m, nil
		}
		m.TitleFrame++
		if m.TitleFrame >= titleFrames {
			return m.finishIntro()
		}
		return m, titleTick()
	case explanationSavedMsg:
		if !m.ExplanationSaving {
			return m, nil
		}
		m.ExplanationSaving = false
		if msg.Err != nil {
			m.ExplanationSaveError = msg.Err.Error()
			return m, nil
		}
		return m.resumeAfterExplanation()
	case tea.MouseMsg:
		if m.ExplanationSavePrompt {
			return m, nil
		}
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && m.unsavedExplanationVisible() {
			m.ExplanationLeaveKey = tea.KeyMsg{Type: tea.KeyEsc}
			m.ExplanationSavePrompt = true
		}
		return m, nil
	case loadingTickMsg:
		if m.TutorActive || m.CandidateActive {
			m.LoadingFrame++
			return m, loadingTick()
		}
		return m, nil
	case candidateResponseMsg:
		if msg.ID != m.CandidateRequestID || !m.CandidateActive {
			return m, nil
		}
		m.CandidateActive = false
		if m.CandidateCancel != nil {
			m.CandidateCancel()
			m.CandidateCancel = nil
		}
		m.PageScroll = 0
		if msg.Candidate != nil {
			if m.DB != nil {
				if err := m.DB.SaveCandidate(*msg.Candidate); err != nil {
					m.CandidateNotice = "Could not save candidate: " + err.Error()
					return m, nil
				}
			}
			m.Candidates = append([]candidate.CandidateQuestion{*msg.Candidate}, m.Candidates...)
			m.CandidateIndex = 0
		}
		if msg.Err != nil {
			m.CandidateNotice = "LLM generation failed: " + msg.Err.Error()
		} else {
			m.CandidateNotice = "Generated LLM candidate; saved for review. Read the wording and derived answer, then [a] approve."
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.Intro != nil {
			m.Intro.Resize(msg.Width, msg.Height)
		}
		return m, nil

	case introTickMsg:
		if m.State == StateIntro && m.Intro != nil {
			m.Intro.Update()
			return m, introTick()
		}
		return m, nil

	case examTickMsg:
		if m.State == StateExam && m.ExamRunner != nil {
			expired := m.ExamRunner.Tick(time.Second)
			if m.DB != nil {
				_ = m.DB.SaveExamSession(m.ExamRunner.ToSessionRecord())
			}
			if expired {
				m.finalizeExam()
				m.State = StateExamSummary
				return m, nil
			}
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })
		}
		return m, nil

	case tutorResponseMsg:
		if msg.ID != m.TutorRequestID || !m.TutorActive {
			return m, nil
		}
		if m.TutorCancel != nil {
			m.TutorCancel()
			m.TutorCancel = nil
		}
		m.PageScroll = 0
		m.TutorActive = false
		if msg.Err != nil {
			if !errors.Is(msg.Err, context.Canceled) {
				m.TutorError = msg.Err.Error()
			}
			return m, nil
		}
		m.TutorResponse = &msg.Resp
		m.TutorKind = msg.Kind
		m.CurrentHint = msg.Resp.Text
		m.ShowHint = true
		m.TutorError = msg.Resp.FallbackReason
		m.PendingExplanation = nil
		if msg.Kind == "explain" && msg.Resp.Text != "" && !msg.Resp.Fallback && msg.Resp.Provider != "OfflineTutor" && msg.Resp.Provider != tutor.ProviderOffline {
			m.captureExplanation(msg.Resp)
		}
		return m, nil

	case oauthCompleteMsg:
		if msg.Err != nil {
			m.TutorAuthNotice = "OAuth sign-in failed: " + msg.Err.Error()
		} else {
			m.TutorAuthNotice = "✓ Successfully connected to ChatGPT Plus! Active tutor updated."
			m.rebuildTutor()
		}
		m.OAuthActiveFlow = nil
		if msg.Err == nil {
			m.TutorFetchingModels = true
			return m, m.cmdRefreshModels(tutor.ProviderChatGPTPlan)
		}
		return m, nil

	case modelsFetchedMsg:
		m.TutorFetchingModels = false
		if msg.Err != nil {
			m.TutorAuthNotice = "Model discovery failed: " + msg.Err.Error()
		} else {
			m.TutorModelList = msg.Models
			if m.ModelCache != nil {
				_ = m.ModelCache.SetModels(msg.Provider, msg.Models)
			}
			m.TutorModelSelectActive = true
			m.TutorCustomModelActive = false
			m.TutorAuthNotice = fmt.Sprintf("✓ Discovered %d live models from %s API! Select with [1-9] or [↑/↓/Enter].", len(msg.Models), strings.ToUpper(msg.Provider))
			// Position cursor on current active model if present
			m.TutorModelCursor = 0
			if m.AuthStore != nil {
				currentModel := m.AuthStore.ResolveModel(msg.Provider)
				for idx, item := range m.TutorModelList {
					if item.ID == currentModel {
						m.TutorModelCursor = idx
						break
					}
				}
			}
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()
		if m.ExplanationSavePrompt {
			return m.updateExplanationSave(key)
		}
		if m.unsavedExplanationVisible() && !explanationReadingKey(key) {
			m.ExplanationLeaveKey = msg
			m.ExplanationSavePrompt = true
			m.ExplanationSaveError = ""
			return m, nil
		}
		if m.State == StateTitle {
			if key == "ctrl+c" || key == "q" {
				m.closeSession()
				m.State = StateQuitting
				return m, tea.Quit
			}
			return m, nil
		}
		if m.State != StateIntro && m.State != StateQuitting && !m.pageTextEntryActive() {
			if key == "u" || key == "d" || key == "pgup" || key == "pgdown" {
				step := (m.Height - 3) / 2
				if step < 1 {
					step = 1
				}
				if key == "u" || key == "pgup" {
					step = -step
				}
				scroll := &m.PageScroll
				if m.State == StateRecap {
					scroll = &m.RecapScroll
				}
				if m.State == StateStatements {
					scroll = &m.StatementsScroll
				}
				*scroll += step
				if *scroll < 0 {
					*scroll = 0
				}
				// Rendering clamps the offset to the current content and terminal size.
				_ = m.View()
				return m, nil
			}
		}
		m.PageScroll = 0

		// Global emergency exit
		if key == "ctrl+c" {
			m.cancelCandidate()
			m.cancelTutor()
			m.closeExamInterrupted()
			m.State = StateQuitting
			return m, tea.Quit
		}

		// Arcade shortcuts work from any screen that has no text entry or running timer.
		if key == "V" && m.arcadeKeysAllowed() {
			return m.openSavedExplanations()
		}
		if m.arcadeKeysAllowed() {
			switch key {
			case "A", "ctrl+a":
				return m.launchArcade(false)
			case "L":
				return m.launchArcade(true)
			}
		}

		switch m.State {
		case StateDrill:
			return m.updateDrill(key)

		case StateFeedback:
			return m.updateFeedback(key)

		case StateRecap:
			return m.updateRecap(key)

		case StateMastery:
			return m.updateMastery(key)

		case StateHelp:
			return m.updateHelp(key)

		case StateSessionComplete:
			return m.updateSessionComplete(key)

		case StateTutorConfig:
			return m.updateTutorConfig(key)

		case StateCandidatePreview:
			return m.updateCandidatePreview(key)

		case StateJournalPractice:
			return m.updateJournalPractice(key)

		case StateStatements:
			return m.updateStatements(key)

		case StateExam:
			return m.updateExam(key)

		case StateExamSummary:
			return m.updateExamSummary(key)

		case StateExamResumePrompt:
			return m.updateExamResumePrompt(key)

		case StateIntro:
			return m.updateIntro(key)
		case StateSavedExplanations:
			return m.updateSavedExplanations(key)

		case StateQuitting:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) cancelTutor() {
	m.TutorRequestID++
	if m.TutorCancel != nil {
		m.TutorCancel()
		m.TutorCancel = nil
	}
	m.TutorActive = false
}

// Only free-form text entry consumes the page scrolling letters.
func (m *Model) pageTextEntryActive() bool {
	return m.State == StateTutorConfig && (m.TutorInputActive || m.TutorOAuthInputActive || m.TutorCustomModelActive) ||
		m.State == StateJournalPractice && m.JournalInputActive && m.JournalInputMode == "amount"
}

func (m *Model) rebuildTutor() {
	if m.AuthStore == nil {
		return
	}
	m.Tutor = tutor.BuildTutor(tutor.FactoryOptions{
		AuthStore: m.AuthStore,
		Timeout:   m.TutorTimeout,
	})
}

func (m *Model) requestTutor(kind string) (tea.Model, tea.Cmd) {
	m.cancelTutor()

	st, err := m.Session.CurrentStage()
	if err != nil {
		return m, nil
	}

	// Record pedagogical assistance in session
	if kind == "hint" {
		_, _ = m.Session.RequestHint()
	} else {
		m.Session.RequestExplanation()
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.TutorTimeout)
	m.TutorCancel = cancel
	m.TutorActive = true
	m.TutorKind = kind
	m.TutorError = ""

	req := m.buildTutorRequest(st)
	id := m.TutorRequestID
	tut := m.Tutor

	cmd := func() tea.Msg {
		var resp tutor.Response
		var err error
		if kind == "hint" {
			resp, err = tut.Hint(ctx, req)
		} else {
			resp, err = tut.Explain(ctx, req)
		}
		return tutorResponseMsg{
			ID:   id,
			Kind: kind,
			Resp: resp,
			Err:  err,
		}
	}

	return m, tea.Batch(cmd, loadingTick())
}

func (m *Model) buildTutorRequest(st *domain.StageAnswer) tutor.Request {
	req := tutor.Request{
		ProblemPrompt: m.CurrentInstance.PromptText,
		FamilyID:      m.CurrentInstance.FamilyID,
		Stage:         st.Stage,
		StagePrompt:   st.Prompt,
		Options:       st.Options,
		CausalHint:    st.CausalHint,
		Explanation:   st.Explanation,
		ConceptID:     st.RelevantConceptID,
	}
	if m.SelectedOptionIndex >= 0 && m.SelectedOptionIndex < len(st.Options) {
		opt := st.Options[m.SelectedOptionIndex]
		req.SelectedOption = &opt
		req.ErrorTag = opt.ErrorTag
		req.MistakeHint = st.MistakeHints[opt.ID]
	}
	if m.LastFeedback != nil && !m.LastFeedback.IsCorrect && len(m.Session.Attempts) > 0 {
		last := m.Session.Attempts[len(m.Session.Attempts)-1]
		if last.Stage == st.Stage {
			req.SelectedOption = &m.LastFeedback.SelectedOption
			req.ErrorTag = m.LastFeedback.ErrorTag
			req.MistakeHint = st.MistakeHints[m.LastFeedback.SelectedOption.ID]
		}
	}
	if req.MistakeHint != "" {
		req.CausalHint = req.MistakeHint
	}
	return req
}

func (m *Model) updateDrill(key string) (tea.Model, tea.Cmd) {
	if m.IsReviewingStage() {
		switch key {
		case "q":
			m.cancelTutor()
			m.closeSession()
			m.State = StateQuitting
			return m, tea.Quit

		case "s":
			m.PreviousState = StateDrill
			m.State = StateMastery
			return m, nil

		case "h", "f1":
			m.PreviousState = StateDrill
			m.State = StateHelp
			return m, nil

		case "t":
			m.cancelTutor()
			m.PreviousState = StateDrill
			m.State = StateTutorConfig
			m.TutorAuthNotice = ""
			return m, nil

		case "p":
			m.cancelTutor()
			m.PreviousState = StateDrill
			m.State = StateCandidatePreview
			m.loadCandidates()
			return m, nil

		case "n":
			return m.requestCandidate()

		case "J", "ctrl+j":
			m.cancelTutor()
			m.PreviousState = StateDrill
			m.initJournalPractice()
			m.State = StateJournalPractice
			return m, nil

		case "F", "ctrl+f":
			m.cancelTutor()
			m.PreviousState = StateDrill
			if m.StatementsReport == nil {
				c := statements.CanonicalCasePioneerConsulting()
				m.StatementsReport, _ = statements.BuildAccountingCycleReport(c, m.Catalog, m.Engine)
			}
			m.StatementsScroll = 0
			m.State = StateStatements
			return m, nil

		case "ctrl+left", "shift+left", "alt+left":
			return m.cycleQuestion(-1)

		case "ctrl+right", "shift+right", "alt+right":
			return m.cycleQuestion(1)

		case "left":
			return m.cycleSubQuestion(-1)

		case "right", "enter", " ":
			return m.cycleSubQuestion(1)

		case "esc":
			if m.Session != nil && !m.Session.IsCompleted {
				m.ViewingStageIndex = m.Session.CurrentIndex
			}
			return m, nil

		default:
			return m, nil
		}
	}

	st, err := m.Session.CurrentStage()
	if err != nil {
		return m, nil
	}

	numOpts := len(st.Options)

	switch key {
	case "q":
		m.cancelTutor()
		m.closeSession()
		m.State = StateQuitting
		return m, tea.Quit

	case "s":
		m.PreviousState = StateDrill
		m.State = StateMastery
		return m, nil

	case "p":
		m.cancelTutor()
		m.PreviousState = StateDrill
		m.State = StateCandidatePreview
		m.loadCandidates()
		return m, nil
	case "n":
		return m.requestCandidate()

	case "t":
		m.cancelTutor()
		m.PreviousState = StateDrill
		m.State = StateTutorConfig
		m.TutorAuthNotice = ""
		return m, nil

	case "J", "ctrl+j":
		m.cancelTutor()
		m.PreviousState = StateDrill
		m.initJournalPractice()
		m.State = StateJournalPractice
		return m, nil

	case "F", "ctrl+f":
		m.cancelTutor()
		m.PreviousState = StateDrill
		if m.StatementsReport == nil {
			c := statements.CanonicalCasePioneerConsulting()
			m.StatementsReport, _ = statements.BuildAccountingCycleReport(c, m.Catalog, m.Engine)
		}
		m.StatementsScroll = 0
		m.State = StateStatements
		return m, nil

	case "E", "ctrl+e":
		m.cancelTutor()
		m.PreviousState = StateDrill
		if m.DB != nil {
			interrupted, err := m.DB.GetLatestInterruptedExamSession()
			if err == nil && interrupted != nil {
				m.InterruptedExam = interrupted
				m.State = StateExamResumePrompt
				return m, nil
			}
		}
		_ = m.startNewExam(m.ExamTimeLimit)
		m.State = StateExam
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })

	case "h", "f1":
		if m.Session != nil {
			m.Session.RecordReferenceUse()
		}
		m.PreviousState = StateDrill
		m.State = StateHelp
		return m, nil

	case "i":
		m.cycleIntensity()
		return m, nil

	case "[", "-":
		if m.TotalQuestions > 5 {
			m.TotalQuestions -= 5
		}
		return m, nil

	case "]", "+", "=":
		if m.TotalQuestions < 20 {
			m.TotalQuestions += 5
		}
		return m, nil

	case "?":
		return m.requestTutor("hint")

	case "e":
		return m.requestTutor("explain")

	case "esc":
		if m.TutorActive {
			m.cancelTutor()
			return m, nil
		}
		if m.ShowHint {
			m.ShowHint = false
			return m, nil
		}
		return m, nil

	// Big-picture question cycling
	case "ctrl+left", "shift+left", "alt+left":
		return m.cycleQuestion(-1)

	case "ctrl+right", "shift+right", "alt+right":
		return m.cycleQuestion(1)

	// Sub-question (stage) cycling with bookend transitions
	case "left":
		return m.cycleSubQuestion(-1)

	case "right":
		return m.cycleSubQuestion(1)

	// Vim / Arrow navigation
	case "up", "k":
		if m.SelectedOptionIndex > 0 {
			m.SelectedOptionIndex--
		} else {
			m.SelectedOptionIndex = numOpts - 1
		}
		return m, nil

	case "down", "j":
		if m.SelectedOptionIndex < numOpts-1 {
			m.SelectedOptionIndex++
		} else {
			m.SelectedOptionIndex = 0
		}
		return m, nil

	case "g":
		// Vim jump to top
		m.SelectedOptionIndex = 0
		return m, nil

	case "G":
		// Vim jump to bottom
		if numOpts > 0 {
			m.SelectedOptionIndex = numOpts - 1
		}
		return m, nil

	// Numeric choices keep lowercase u/d available for page scrolling.
	case "a", "1":
		if numOpts > 0 {
			m.SelectedOptionIndex = 0
		}
		return m, nil

	case "b", "2":
		if numOpts > 1 {
			m.SelectedOptionIndex = 1
		}
		return m, nil

	case "c", "3":
		if numOpts > 2 {
			m.SelectedOptionIndex = 2
		}
		return m, nil

	case "4":
		if numOpts > 3 {
			m.SelectedOptionIndex = 3
		}
		return m, nil

	case "enter":
		return m.submitSelectedAnswer()
	}

	return m, nil
}

func (m *Model) submitSelectedAnswer() (tea.Model, tea.Cmd) {
	m.cancelTutor()

	st, err := m.Session.CurrentStage()
	if err != nil || m.SelectedOptionIndex < 0 || m.SelectedOptionIndex >= len(st.Options) {
		return m, nil
	}

	chosenOpt := st.Options[m.SelectedOptionIndex]
	now := m.Clock.Now()

	feedback, err := m.Session.SubmitOption(chosenOpt.ID, now)
	if err != nil {
		return m, nil
	}

	m.LastFeedback = feedback
	if m.StageHistory == nil {
		m.StageHistory = make(map[int]*StageHistoryItem)
	}
	if feedback.AdvanceStage {
		stIdx := m.Session.CurrentIndex - 1
		m.StageHistory[stIdx] = &StageHistoryItem{
			StageIndex:       stIdx,
			StageKey:         st.Stage,
			SelectedOptionID: chosenOpt.ID,
			IsCorrect:        feedback.IsCorrect,
			Explanation:      feedback.Explanation,
			Hint:             feedback.Hint,
			ErrorTag:         feedback.ErrorTag,
			AssistanceLevel:  feedback.AssistanceLevel,
		}
	}
	if !feedback.IsCorrect && (st.Stage == domain.StageIdentifyAccount || st.Stage == domain.StageCounterAccount || st.Stage == domain.StageBalancedEntry || st.Stage == domain.StageEquationEffect) {
		m.Scheduler.QueueContrast(m.CurrentInstance)
	}
	m.SessionAttempts++

	if feedback.IsCorrect {
		if feedback.AssistanceLevel == domain.AssistanceNone {
			m.FirstTrySuccesses++
			m.CurrentStreak++
			if m.CurrentStreak > m.BestStreak {
				m.BestStreak = m.CurrentStreak
			}
		}
	} else {
		// Reset current streak on mistake
		m.CurrentStreak = 0
	}

	// Persist attempt to SQLite if DB is present
	if m.DB != nil && len(m.Session.Attempts) > 0 {
		latestAttempt := m.Session.Attempts[len(m.Session.Attempts)-1]
		_ = m.DB.RecordAttempt(latestAttempt)
	}

	m.State = StateFeedback
	return m, nil
}

func (m *Model) updateFeedback(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "n":
		return m.requestCandidate()
	case "q":
		m.cancelTutor()
		m.closeSession()
		m.State = StateQuitting
		return m, tea.Quit

	case "s":
		m.PreviousState = StateFeedback
		m.State = StateMastery
		return m, nil

	case "t":
		m.cancelTutor()
		m.PreviousState = StateFeedback
		m.State = StateTutorConfig
		m.TutorAuthNotice = ""
		return m, nil

	case "h", "f1":
		m.PreviousState = StateFeedback
		m.State = StateHelp
		return m, nil

	case "?":
		return m.requestTutor("hint")

	case "e":
		return m.requestTutor("explain")

	case "esc":
		if m.TutorActive {
			m.cancelTutor()
			return m, nil
		}
		if m.ShowHint {
			m.ShowHint = false
			return m, nil
		}
		return m, nil

	case "ctrl+left", "shift+left", "alt+left":
		return m.cycleQuestion(-1)

	case "ctrl+right", "shift+right", "alt+right":
		return m.cycleQuestion(1)

	case "left":
		return m.cycleSubQuestion(-1)

	case "right":
		return m.cycleSubQuestion(1)

	case "enter", " ", "a", "b", "c", "1", "2", "3", "4":
		m.cancelTutor()
		if m.LastFeedback != nil && !m.LastFeedback.AdvanceStage {
			// Retry is available! Return to drill so learner can answer again
			m.State = StateDrill
			m.ViewingStageIndex = m.Session.CurrentIndex
			m.ShowHint = true
			m.CurrentHint = m.LastFeedback.Hint
			if key != "enter" && key != " " {
				return m.updateDrill(key)
			}
			return m, nil
		}

		// Stage was completed (either correct or revealed answer)
		if m.Session.IsCompleted {
			m.RecapScroll = 0
			m.State = StateRecap
			if len(m.Session.StageSequence) > 0 {
				m.ViewingStageIndex = len(m.Session.StageSequence) - 1
			}
		} else {
			m.State = StateDrill
			m.ViewingStageIndex = m.Session.CurrentIndex
			m.SelectedOptionIndex = 0
			m.ShowHint = false
			m.LastFeedback = nil
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) updateRecap(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q":
		m.closeSession()
		m.State = StateQuitting
		return m, tea.Quit

	case "s":
		m.PreviousState = StateRecap
		m.State = StateMastery
		return m, nil

	case "t":
		m.PreviousState = StateRecap
		m.State = StateTutorConfig
		m.TutorAuthNotice = ""
		return m, nil

	case "J", "ctrl+j":
		m.PreviousState = StateRecap
		m.initJournalPractice()
		m.State = StateJournalPractice
		return m, nil

	case "h", "f1":
		m.PreviousState = StateRecap
		m.State = StateHelp
		return m, nil

	case "up", "k":
		if m.RecapScroll > 0 {
			m.RecapScroll--
		}
		return m, nil

	case "down", "j":
		m.RecapScroll++
		return m, nil

	case "pgup", "pageup":
		m.RecapScroll -= 5
		if m.RecapScroll < 0 {
			m.RecapScroll = 0
		}
		return m, nil

	case "pgdown", "pagedown":
		m.RecapScroll += 5
		return m, nil

	case "g", "home":
		m.RecapScroll = 0
		return m, nil

	case "G", "end":
		m.RecapScroll = 9999
		return m, nil

	case "ctrl+left", "shift+left", "alt+left":
		return m.cycleQuestion(-1)

	case "ctrl+right", "shift+right", "alt+right":
		return m.cycleQuestion(1)

	case "left":
		return m.cycleSubQuestion(-1)

	case "right":
		if m.CurrentQuestionIndex < len(m.QuestionHistory)-1 {
			m.restoreQuestionState(m.CurrentQuestionIndex + 1)
			return m, nil
		}
		return m.cycleQuestion(1)

	case "enter", " ":
		m.RecapScroll = 0
		m.saveCurrentQuestionState()
		if m.CurrentQuestionIndex < len(m.QuestionHistory)-1 {
			m.restoreQuestionState(m.CurrentQuestionIndex + 1)
			return m, nil
		}

		m.QuestionsCompleted++
		if m.CurrentQuestionIndex+1 >= m.TotalQuestions {
			m.closeSession()
			m.State = StateSessionComplete
			return m, nil
		}

		m.CurrentQuestionIndex++
		if err := m.loadNextQuestion(); err != nil {
			m.closeSession()
			m.State = StateSessionComplete
			return m, nil
		}

		m.State = StateDrill
		m.SelectedOptionIndex = 0
		m.ShowHint = false
		m.LastFeedback = nil
		return m, nil
	}

	return m, nil
}

func (m *Model) updateMastery(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "i":
		m.cycleIntensity()
		return m, nil

	case "[", "-":
		if m.TotalQuestions > 5 {
			m.TotalQuestions -= 5
		}
		return m, nil

	case "]", "+", "=":
		if m.TotalQuestions < 20 {
			m.TotalQuestions += 5
		}
		return m, nil

	case "t":
		m.PreviousState = StateMastery
		m.State = StateTutorConfig
		m.TutorAuthNotice = ""
		return m, nil

	case "s", "esc", "enter", " ":
		m.State = m.PreviousState
		return m, nil

	case "h", "f1":
		m.PreviousState = StateMastery
		m.State = StateHelp
		return m, nil

	case "q":
		m.closeSession()
		m.State = StateQuitting
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) cycleIntensity() {
	curr := m.Scheduler.Intensity()
	var next mastery.SessionIntensity
	switch curr {
	case mastery.IntensityStandard:
		next = mastery.IntensitySpaced
	case mastery.IntensitySpaced:
		next = mastery.IntensityIntensive
	case mastery.IntensityIntensive:
		next = mastery.IntensityTransfer
	default:
		next = mastery.IntensityStandard
	}
	m.Scheduler.SetIntensity(next)
}

func (m *Model) updateHelp(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "h", "f1", "esc", "enter", " ", "q":
		m.State = m.PreviousState
		return m, nil
	}

	return m, nil
}

func (m *Model) updateSessionComplete(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		m.State = StateQuitting
		return m, tea.Quit

	case "s":
		m.PreviousState = StateSessionComplete
		m.State = StateMastery
		return m, nil

	case "t":
		m.PreviousState = StateSessionComplete
		m.State = StateTutorConfig
		m.TutorAuthNotice = ""
		return m, nil

	case "h", "f1":
		m.PreviousState = StateSessionComplete
		m.State = StateHelp
		return m, nil

	case "left", "ctrl+left", "shift+left", "alt+left":
		if len(m.QuestionHistory) > 0 {
			m.restoreQuestionState(len(m.QuestionHistory) - 1)
			return m, nil
		}

	case "r", "enter":
		// Restart session
		m.SessionID = fmt.Sprintf("sess-%d", time.Now().Unix())
		m.CurrentQuestionIndex = 0
		m.QuestionsCompleted = 0
		m.QuestionHistory = nil
		m.SessionAttempts = 0
		m.FirstTrySuccesses = 0
		m.CurrentStreak = 0
		if m.DB != nil {
			_ = m.DB.SaveSession(storage.SessionRecord{
				ID:        m.SessionID,
				Mode:      "drill",
				StartedAt: m.Clock.Now(),
			})
		}
		_ = m.loadNextQuestion()
		m.State = StateDrill
		m.SelectedOptionIndex = 0
		m.ShowHint = false
		m.LastFeedback = nil
		return m, nil
	}

	return m, nil
}

func (m *Model) cmdRefreshModels(provider string) tea.Cmd {
	cache := m.ModelCache
	authStore := m.AuthStore
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		models, err := cache.RefreshProvider(ctx, provider, authStore, "", nil)
		return modelsFetchedMsg{
			Provider: provider,
			Models:   models,
			Err:      err,
		}
	}
}

func (m *Model) updateTutorConfig(key string) (tea.Model, tea.Cmd) {
	if m.TutorInputActive {
		switch key {
		case "enter":
			if len(m.TutorInputBuffer) > 0 && m.AuthStore != nil {
				_ = m.AuthStore.SetAPIKey(m.TutorInputProvider, m.TutorInputBuffer)
				_ = m.AuthStore.SetActiveProvider(m.TutorInputProvider)
				m.rebuildTutor()
				m.TutorAuthNotice = fmt.Sprintf("✓ %s API key saved and activated!", strings.ToUpper(m.TutorInputProvider))
			}
			m.TutorInputActive = false
			m.TutorInputBuffer = ""
			return m, nil

		case "esc":
			m.TutorInputActive = false
			m.TutorInputBuffer = ""
			return m, nil

		case "backspace":
			if len(m.TutorInputBuffer) > 0 {
				m.TutorInputBuffer = m.TutorInputBuffer[:len(m.TutorInputBuffer)-1]
			}
			return m, nil

		default:
			if len(key) == 1 {
				m.TutorInputBuffer += key
			}
			return m, nil
		}
	}

	if m.TutorOAuthInputActive {
		switch key {
		case "enter":
			if m.AuthStore != nil {
				_ = m.AuthStore.SetChatGPTClientID(m.TutorOAuthInputBuffer)
				m.rebuildTutor()
				m.TutorAuthNotice = fmt.Sprintf("✓ OpenAI OAuth Client ID saved: %s", m.TutorOAuthInputBuffer)
			}
			m.TutorOAuthInputActive = false
			m.TutorOAuthInputBuffer = ""
			return m, nil

		case "esc":
			m.TutorOAuthInputActive = false
			m.TutorOAuthInputBuffer = ""
			return m, nil

		case "backspace":
			if len(m.TutorOAuthInputBuffer) > 0 {
				m.TutorOAuthInputBuffer = m.TutorOAuthInputBuffer[:len(m.TutorOAuthInputBuffer)-1]
			}
			return m, nil

		default:
			if len(key) == 1 {
				m.TutorOAuthInputBuffer += key
			}
			return m, nil
		}
	}

	if m.TutorModelSelectActive {
		if m.TutorCustomModelActive {
			switch key {
			case "enter":
				trimmed := strings.TrimSpace(m.TutorCustomModelBuffer)
				if trimmed != "" && m.AuthStore != nil {
					activeProvider := m.AuthStore.GetConfig().ActiveProvider
					_ = m.AuthStore.SetSelectedModel(activeProvider, trimmed)
					m.rebuildTutor()
					m.TutorAuthNotice = fmt.Sprintf("✓ Active model for %s set to custom: %s", strings.ToUpper(activeProvider), trimmed)
				}
				m.TutorCustomModelActive = false
				m.TutorCustomModelBuffer = ""
				m.TutorModelSelectActive = false
				return m, nil

			case "esc":
				m.TutorCustomModelActive = false
				m.TutorCustomModelBuffer = ""
				return m, nil

			case "backspace":
				if len(m.TutorCustomModelBuffer) > 0 {
					m.TutorCustomModelBuffer = m.TutorCustomModelBuffer[:len(m.TutorCustomModelBuffer)-1]
				}
				return m, nil

			default:
				if len(key) == 1 {
					m.TutorCustomModelBuffer += key
				}
				return m, nil
			}
		}

		switch key {
		case "j", "down":
			if len(m.TutorModelList) > 0 && m.TutorModelCursor < len(m.TutorModelList)-1 {
				m.TutorModelCursor++
			}
			return m, nil

		case "k", "up":
			if m.TutorModelCursor > 0 {
				m.TutorModelCursor--
			}
			return m, nil

		case "enter":
			if len(m.TutorModelList) > 0 && m.TutorModelCursor < len(m.TutorModelList) && m.AuthStore != nil {
				activeProvider := m.AuthStore.GetConfig().ActiveProvider
				chosen := m.TutorModelList[m.TutorModelCursor].ID
				_ = m.AuthStore.SetSelectedModel(activeProvider, chosen)
				m.rebuildTutor()
				m.TutorAuthNotice = fmt.Sprintf("✓ Active model for %s set to: %s", strings.ToUpper(activeProvider), chosen)
				m.TutorModelSelectActive = false
			}
			return m, nil

		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx := int(key[0] - '1')
			if idx >= 0 && idx < len(m.TutorModelList) && m.AuthStore != nil {
				activeProvider := m.AuthStore.GetConfig().ActiveProvider
				chosen := m.TutorModelList[idx].ID
				_ = m.AuthStore.SetSelectedModel(activeProvider, chosen)
				m.rebuildTutor()
				m.TutorAuthNotice = fmt.Sprintf("✓ Active model for %s set to: %s", strings.ToUpper(activeProvider), chosen)
				m.TutorModelSelectActive = false
			}
			return m, nil

		case "r":
			if m.AuthStore != nil {
				activeProvider := m.AuthStore.GetConfig().ActiveProvider
				if activeProvider == tutor.ProviderOffline {
					m.TutorAuthNotice = "Offline provider uses deterministic rules; no remote models to discover."
					return m, nil
				}
				m.TutorFetchingModels = true
				m.TutorAuthNotice = fmt.Sprintf("Discovering live models from %s API...", strings.ToUpper(activeProvider))
				return m, m.cmdRefreshModels(activeProvider)
			}
			return m, nil

		case "c":
			m.TutorCustomModelActive = true
			m.TutorCustomModelBuffer = ""
			return m, nil

		case "esc":
			m.TutorModelSelectActive = false
			return m, nil
		}

		return m, nil
	}

	switch key {
	case "1":
		if m.AuthStore != nil {
			_ = m.AuthStore.SetActiveProvider(tutor.ProviderOffline)
			m.rebuildTutor()
			m.TutorAuthNotice = "Switched to Offline Tutor (100% local machine mode)."
		}
		return m, nil

	case "2":
		if m.OAuthActiveFlow != nil {
			return m, nil
		}
		// Continue with ChatGPT Plus (OAuth)
		if m.AuthStore != nil && m.AuthStore.IsConfigured(tutor.ProviderChatGPTPlan) {
			_ = m.AuthStore.SetActiveProvider(tutor.ProviderChatGPTPlan)
			m.rebuildTutor()
			m.TutorAuthNotice = "Activated connected ChatGPT Plus plan."
			return m, nil
		}

		chatgptTutor := tutor.NewOpenAIChatGPTPlanTutor(tutor.ChatGPTPlanConfig{
			AuthStore: m.AuthStore,
		})
		flow, authURL, err := chatgptTutor.StartOAuthFlow()
		if err != nil {
			m.TutorAuthNotice = "Failed starting OAuth flow: " + err.Error()
			return m, nil
		}
		m.OAuthActiveFlow = flow
		m.TutorAuthNotice = fmt.Sprintf("Login listener active at %s\nPlease authorize in browser: %s", flow.RedirectURI, authURL)

		return m, func() tea.Msg {
			_ = tutor.OpenSignInBrowser(authURL)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			tok, err := flow.Complete(ctx)
			return oauthCompleteMsg{Token: tok, Err: err}
		}

	case "3":
		m.TutorInputActive = true
		m.TutorInputProvider = tutor.ProviderAnthropic
		m.TutorInputBuffer = ""
		m.TutorAuthNotice = ""
		return m, nil

	case "4":
		m.TutorInputActive = true
		m.TutorInputProvider = tutor.ProviderGoogle
		m.TutorInputBuffer = ""
		m.TutorAuthNotice = ""
		return m, nil

	case "5":
		m.TutorInputActive = true
		m.TutorInputProvider = tutor.ProviderOpenAI
		m.TutorInputBuffer = ""
		m.TutorAuthNotice = ""
		return m, nil

	case "m":
		if m.AuthStore != nil {
			activeProvider := m.AuthStore.GetConfig().ActiveProvider
			if activeProvider == tutor.ProviderOffline {
				m.TutorAuthNotice = "Offline tutor uses deterministic rule engine. Choose an AI provider [2-5] first, then press 'm' to select or change its model."
				return m, nil
			}
			m.TutorModelSelectActive = true
			m.TutorCustomModelActive = false
			m.TutorModelList = m.ModelCache.GetModels(activeProvider)
			currentModel := m.AuthStore.ResolveModel(activeProvider)
			m.TutorModelCursor = 0
			for idx, item := range m.TutorModelList {
				if item.ID == currentModel {
					m.TutorModelCursor = idx
					break
				}
			}
		}
		return m, nil

	case "r":
		if m.AuthStore != nil {
			activeProvider := m.AuthStore.GetConfig().ActiveProvider
			if activeProvider == tutor.ProviderOffline {
				m.TutorAuthNotice = "Offline provider uses deterministic rules; no remote models to discover."
				return m, nil
			}
			m.TutorFetchingModels = true
			m.TutorAuthNotice = fmt.Sprintf("Discovering live models from %s API...", strings.ToUpper(activeProvider))
			return m, m.cmdRefreshModels(activeProvider)
		}
		return m, nil

	case "a":
		if m.OAuthActiveFlow != nil {
			m.OAuthActiveFlow.Close()
			m.OAuthActiveFlow = nil
		}
		if m.AuthStore != nil {
			if err := m.AuthStore.ClearCredentials(tutor.ProviderChatGPTPlan); err != nil {
				m.TutorAuthNotice = err.Error()
				return m, nil
			}
			_ = m.AuthStore.SetActiveProvider(tutor.ProviderOffline)
			m.rebuildTutor()
			m.TutorAuthNotice = "ChatGPT disconnected locally. Access can also be revoked in ChatGPT Settings. Press [2] to connect another account."
		}
		return m, nil

	case "x":
		if m.AuthStore != nil {
			active := m.AuthStore.GetConfig().ActiveProvider
			_ = m.AuthStore.ClearCredentials(active)
			_ = m.AuthStore.SetActiveProvider(tutor.ProviderOffline)
			m.rebuildTutor()
			m.TutorAuthNotice = fmt.Sprintf("Cleared credentials for %s. Reset to Offline Tutor.", active)
		}
		return m, nil

	case "esc", "t", "q":
		if m.OAuthActiveFlow != nil {
			m.OAuthActiveFlow.Close()
			m.OAuthActiveFlow = nil
		}
		m.State = m.PreviousState
		return m, nil
	}

	return m, nil
}

func (m *Model) closeSession() {
	m.cancelTutor()
	if m.DB != nil {
		_ = m.DB.CompleteSession(m.SessionID, m.Clock.Now())
	}
}

func (m *Model) loadNextQuestion() error {
	m.cancelTutor()
	m.TutorError = ""
	m.TutorResponse = nil
	m.PendingExplanation = nil
	m.ShowHint = false
	m.CurrentHint = ""
	m.RecapScroll = 0

	// Rebuild projections from all historical attempts
	var allAttempts []domain.Attempt
	if m.DB != nil {
		att, err := m.DB.GetAllAttempts()
		if err == nil {
			allAttempts = att
		}
	}
	// Also include current session attempts in case DB is absent
	if m.Session != nil {
		allAttempts = append(allAttempts, m.Session.Attempts...)
	}

	projections := mastery.RebuildProjections(allAttempts)

	chosenQ := m.Scheduler.SelectNextQuestion(m.Questions, projections)
	if chosenQ == nil {
		return fmt.Errorf("no eligible active questions available for practice session")
	}

	// Select parameter amount
	amtOptions := chosenQ.Parameters["amount_minor_units"]
	var chosenAmt int64 = 10000 // $100 default
	if len(amtOptions) > 0 {
		chosenAmt = amtOptions[m.RNG.Intn(len(amtOptions))]
	}

	if m.Scheduler.LastSelectionWasContrast && m.CurrentInstance != nil {
		priorAmount := m.CurrentInstance.Parameters["amount_minor_units"]
		for _, amount := range amtOptions {
			if amount == priorAmount {
				chosenAmt = priorAmount
				break
			}
		}
	}
	seed := m.RNG.Int63()
	params := map[string]int64{"amount_minor_units": chosenAmt}

	scaffold := mastery.ComputeQuestionScaffoldLevel(*chosenQ, projections, m.Scheduler.Intensity())
	if m.Scheduler.LastSelectionWasContrast {
		scaffold = domain.ScaffoldFull
	}
	inst, err := m.Generator.GenerateInstanceWithScaffold(*chosenQ, seed, params, scaffold)
	if err != nil {
		return fmt.Errorf("failed generating question instance: %w", err)
	}

	if m.Scheduler.LastSelectionWasContrast {
		inst.Pedagogy.Remediation = true
		inst.InstanceID += "-contrast"
	}
	// Save instance to DB
	if m.DB != nil {
		_ = m.DB.SaveQuestionInstance(inst, m.SessionID)
	}

	m.CurrentInstance = inst
	m.Session = drill.NewSession(m.SessionID, inst)
	m.SelectedOptionIndex = 0
	m.ShowHint = false
	m.LastFeedback = nil
	m.ViewingStageIndex = 0
	m.StageHistory = make(map[int]*StageHistoryItem)

	qs := &DrillQuestionState{
		Index:               m.CurrentQuestionIndex,
		Instance:            inst,
		Session:             m.Session,
		State:               StateDrill,
		SelectedOptionIndex: 0,
		ViewingStageIndex:   0,
		StageHistory:        m.StageHistory,
	}
	if m.CurrentQuestionIndex < len(m.QuestionHistory) {
		m.QuestionHistory[m.CurrentQuestionIndex] = qs
	} else {
		m.QuestionHistory = append(m.QuestionHistory, qs)
	}

	return nil
}

func (m *Model) saveCurrentQuestionState() {
	if m.QuestionHistory == nil || m.CurrentQuestionIndex < 0 || m.CurrentQuestionIndex >= len(m.QuestionHistory) {
		return
	}
	if m.State != StateDrill && m.State != StateFeedback && m.State != StateRecap {
		return
	}
	qs := m.QuestionHistory[m.CurrentQuestionIndex]
	qs.Instance = m.CurrentInstance
	qs.Session = m.Session
	qs.State = m.State
	qs.SelectedOptionIndex = m.SelectedOptionIndex
	qs.LastFeedback = m.LastFeedback
	qs.CurrentHint = m.CurrentHint
	qs.RenderedHint = m.renderedHint
	qs.RenderedHintRaw = m.renderedHintRaw
	qs.RenderedHintWidth = m.renderedHintWidth
	qs.ShowHint = m.ShowHint
	qs.TutorKind = m.TutorKind
	qs.TutorResponse = m.TutorResponse
	qs.RecapScroll = m.RecapScroll
	qs.ViewingStageIndex = m.ViewingStageIndex
	qs.StageHistory = m.StageHistory
}

func (m *Model) restoreQuestionState(index int) {
	if m.QuestionHistory == nil || index < 0 || index >= len(m.QuestionHistory) {
		return
	}
	m.cancelTutor()
	m.saveCurrentQuestionState()
	m.CurrentQuestionIndex = index
	qs := m.QuestionHistory[index]
	m.CurrentInstance = qs.Instance
	m.Session = qs.Session
	if qs.Session != nil && qs.Session.IsCompleted {
		m.State = StateRecap
		if qs.Session.StageSequence != nil && len(qs.Session.StageSequence) > 0 {
			m.ViewingStageIndex = len(qs.Session.StageSequence) - 1
		}
	} else {
		m.State = qs.State
		m.ViewingStageIndex = qs.ViewingStageIndex
	}
	m.SelectedOptionIndex = qs.SelectedOptionIndex
	m.LastFeedback = qs.LastFeedback
	m.CurrentHint = qs.CurrentHint
	m.renderedHint = qs.RenderedHint
	m.renderedHintRaw = qs.RenderedHintRaw
	m.renderedHintWidth = qs.RenderedHintWidth
	m.ShowHint = qs.ShowHint
	m.TutorKind = qs.TutorKind
	m.TutorResponse = qs.TutorResponse
	m.RecapScroll = qs.RecapScroll
	m.StageHistory = qs.StageHistory
	if m.StageHistory == nil {
		m.StageHistory = make(map[int]*StageHistoryItem)
	}
}

func (m *Model) cycleQuestion(dir int) (tea.Model, tea.Cmd) {
	if len(m.QuestionHistory) <= 1 {
		return m, nil
	}
	m.saveCurrentQuestionState()
	newIdx := m.CurrentQuestionIndex + dir
	if newIdx < 0 {
		newIdx = len(m.QuestionHistory) - 1
	} else if newIdx >= len(m.QuestionHistory) {
		newIdx = 0
	}
	m.restoreQuestionState(newIdx)
	return m, nil
}

func (m *Model) IsReviewingStage() bool {
	if m.Session == nil || m.State != StateDrill {
		return false
	}
	if m.Session.IsCompleted {
		return true
	}
	return m.ViewingStageIndex < m.Session.CurrentIndex
}

func (m *Model) getStageHistoryItem(stageIdx int) *StageHistoryItem {
	if m.StageHistory != nil {
		if item, ok := m.StageHistory[stageIdx]; ok && item != nil {
			return item
		}
	}

	if m.Session == nil || stageIdx < 0 || stageIdx >= len(m.Session.StageSequence) {
		return nil
	}

	stKey := m.Session.StageSequence[stageIdx]
	stAns, ok := m.CurrentInstance.StageAnswers[stKey]
	if !ok {
		return nil
	}

	var matchedAttempt *domain.Attempt
	for i := len(m.Session.Attempts) - 1; i >= 0; i-- {
		if m.Session.Attempts[i].Stage == stKey {
			matchedAttempt = &m.Session.Attempts[i]
			break
		}
	}

	selectedOptID := stAns.CorrectOptionID
	isCorrect := true
	errorTag := ""
	assist := domain.AssistanceNone

	if matchedAttempt != nil {
		selectedOptID = matchedAttempt.SelectedOptionID
		isCorrect = matchedAttempt.IsCorrect
		errorTag = matchedAttempt.ErrorTag
		assist = matchedAttempt.Assistance
	}

	return &StageHistoryItem{
		StageIndex:       stageIdx,
		StageKey:         stKey,
		SelectedOptionID: selectedOptID,
		IsCorrect:        isCorrect,
		Explanation:      stAns.Explanation,
		Hint:             stAns.HintForOption(selectedOptID),
		ErrorTag:         errorTag,
		AssistanceLevel:  assist,
	}
}

func (m *Model) cycleSubQuestion(dir int) (tea.Model, tea.Cmd) {
	if m.Session == nil {
		return m, nil
	}
	m.cancelTutor()

	// If in StateRecap:
	if m.State == StateRecap {
		if dir < 0 {
			// Step back into the last sub-question of this completed question
			m.State = StateDrill
			if len(m.Session.StageSequence) > 0 {
				m.ViewingStageIndex = len(m.Session.StageSequence) - 1
			}
			return m, nil
		}
		// If dir > 0 in Recap, advance to next question
		if m.CurrentQuestionIndex < len(m.QuestionHistory)-1 {
			m.restoreQuestionState(m.CurrentQuestionIndex + 1)
			return m, nil
		}
		return m.cycleQuestion(1)
	}

	// If in StateFeedback:
	if m.State == StateFeedback {
		if dir < 0 {
			activeStage := m.Session.CurrentIndex
			if m.LastFeedback != nil && m.LastFeedback.AdvanceStage {
				activeStage = m.Session.CurrentIndex - 1
			}
			if activeStage > 0 {
				m.State = StateDrill
				m.ViewingStageIndex = activeStage - 1
				return m, nil
			}
			return m.cycleQuestion(-1)
		}
		// dir > 0: advance from feedback
		if m.LastFeedback != nil && !m.LastFeedback.AdvanceStage {
			m.State = StateDrill
			m.ViewingStageIndex = m.Session.CurrentIndex
			m.ShowHint = true
			m.CurrentHint = m.LastFeedback.Hint
			return m, nil
		}
		if m.Session.IsCompleted {
			m.RecapScroll = 0
			m.State = StateRecap
			if len(m.Session.StageSequence) > 0 {
				m.ViewingStageIndex = len(m.Session.StageSequence) - 1
			}
			return m, nil
		}
		m.State = StateDrill
		m.ViewingStageIndex = m.Session.CurrentIndex
		m.SelectedOptionIndex = 0
		m.ShowHint = false
		m.LastFeedback = nil
		return m, nil
	}

	// If in StateDrill:
	totStages := len(m.Session.StageSequence)
	if dir < 0 {
		if m.ViewingStageIndex > 0 {
			m.ViewingStageIndex--
			return m, nil
		}
		// At beginning bookend (Step 1): cycle to previous question!
		return m.cycleQuestion(-1)
	}

	// dir > 0:
	if m.Session.IsCompleted {
		if m.ViewingStageIndex < totStages-1 {
			m.ViewingStageIndex++
			return m, nil
		}
		// Reached the end of stages for this completed question: return to Recap!
		m.State = StateRecap
		m.RecapScroll = 0
		return m, nil
	}

	// In-progress question:
	if m.ViewingStageIndex < m.Session.CurrentIndex {
		m.ViewingStageIndex++
		if m.ViewingStageIndex == m.Session.CurrentIndex {
			m.SelectedOptionIndex = 0
		}
		return m, nil
	}

	// Already at the active unanswered stage: future stages are locked
	return m, nil
}

func (m *Model) loadCandidates() {
	if m.DB != nil {
		cands, err := m.DB.ListCandidates(storage.CandidateFilter{})
		if err == nil && len(cands) > 0 {
			m.Candidates = cands
			m.CandidateIndex = 0
			return
		}
	}

	// If no candidates in DB, generate one on-the-fly using offline generator
	gen := candidate.NewOfflineCandidateGenerator(m.Engine, m.Catalog)
	fam := bank.FamilyCustomerAdvance
	if m.CurrentInstance != nil && m.CurrentInstance.FamilyID != "" {
		fam = m.CurrentInstance.FamilyID
	}
	cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
		FamilyID: fam,
	})
	if err == nil && cand != nil {
		if m.DB != nil {
			_ = m.DB.SaveCandidate(*cand)
		}
		m.Candidates = []candidate.CandidateQuestion{*cand}
		m.CandidateIndex = 0
	}
}

func (m *Model) updateCandidatePreview(key string) (tea.Model, tea.Cmd) {
	if m.CandidateActive {
		if key == "esc" {
			m.cancelCandidate()
			m.CandidateNotice = "LLM generation canceled."
		}
		return m, nil
	}
	switch key {
	case "n":
		return m.requestCandidate()
	case "esc", "p", "q":
		m.State = m.PreviousState
		m.CandidateNotice = ""
		return m, nil

	case "up", "k", "left":
		if len(m.Candidates) > 0 {
			if m.CandidateIndex > 0 {
				m.CandidateIndex--
			} else {
				m.CandidateIndex = len(m.Candidates) - 1
			}
		}
		return m, nil

	case "down", "j", "right":
		if len(m.Candidates) > 0 {
			if m.CandidateIndex < len(m.Candidates)-1 {
				m.CandidateIndex++
			} else {
				m.CandidateIndex = 0
			}
		}
		return m, nil

	case "g":
		// Generate a new candidate
		gen := candidate.NewOfflineCandidateGenerator(m.Engine, m.Catalog)
		targetFamily := bank.FamilyCustomerAdvance
		if m.CurrentInstance != nil && m.CurrentInstance.FamilyID != "" {
			targetFamily = m.CurrentInstance.FamilyID
		}
		cand, err := gen.Generate(context.Background(), candidate.GenerateRequest{
			FamilyID: targetFamily,
			Seed:     time.Now().UnixNano(),
		})
		if err != nil {
			m.CandidateNotice = fmt.Sprintf("Error generating candidate: %v", err)
			return m, nil
		}
		if m.DB != nil {
			_ = m.DB.SaveCandidate(*cand)
			cands, _ := m.DB.ListCandidates(storage.CandidateFilter{})
			m.Candidates = cands
		} else {
			m.Candidates = append([]candidate.CandidateQuestion{*cand}, m.Candidates...)
		}
		m.CandidateIndex = 0
		m.CandidateNotice = fmt.Sprintf("✓ Generated new candidate for %s", cand.FamilyID)
		return m, nil

	case "a":
		// Approve candidate into active bank
		if len(m.Candidates) > 0 {
			cand := m.Candidates[m.CandidateIndex]
			pub, event, err := candidate.ApproveCandidate(&cand, "tui_reviewer", "Approved in TUI candidate modal", time.Now().UTC())
			if err != nil {
				m.CandidateNotice = fmt.Sprintf("⚠️ Approval blocked: %v", err)
				return m, nil
			}
			if m.DB != nil {
				if err := m.DB.PublishQuestion(*pub, cand.ID, *event); err != nil {
					m.CandidateNotice = fmt.Sprintf("Database error publishing question: %v", err)
					return m, nil
				}
				cands, _ := m.DB.ListCandidates(storage.CandidateFilter{})
				m.Candidates = cands
			} else {
				m.Candidates[m.CandidateIndex] = cand
			}
			// Append newly approved question to active pool in this session so learner can drill on it!
			m.Questions = append(m.Questions, *pub)
			m.CandidateNotice = fmt.Sprintf("✓ Approved and published %s into active bank!", cand.ID)
		}
		return m, nil

	case "r":
		// Reject candidate
		if len(m.Candidates) > 0 {
			cand := m.Candidates[m.CandidateIndex]
			event, err := candidate.RejectCandidate(&cand, "tui_reviewer", "Rejected in TUI candidate modal", time.Now().UTC())
			if err != nil {
				m.CandidateNotice = fmt.Sprintf("⚠️ Rejection error: %v", err)
				return m, nil
			}
			if m.DB != nil {
				if err := m.DB.SaveCandidate(cand); err != nil {
					m.CandidateNotice = fmt.Sprintf("⚠️ Database error saving candidate: %v", err)
					return m, nil
				}
				if err := m.DB.RecordApprovalEvent(*event); err != nil {
					m.CandidateNotice = fmt.Sprintf("⚠️ Database error recording rejection event: %v", err)
					return m, nil
				}
				cands, _ := m.DB.ListCandidates(storage.CandidateFilter{})
				m.Candidates = cands
			} else {
				m.Candidates[m.CandidateIndex] = cand
			}
			m.CandidateNotice = fmt.Sprintf("✗ Rejected candidate %s.", cand.ID)
		}
		return m, nil

	case "x", "delete":
		// Delete current candidate
		if len(m.Candidates) > 0 {
			toDelete := m.Candidates[m.CandidateIndex]
			if m.DB != nil {
				_ = m.DB.DeleteCandidate(toDelete.ID)
				cands, _ := m.DB.ListCandidates(storage.CandidateFilter{})
				m.Candidates = cands
			} else {
				m.Candidates = append(m.Candidates[:m.CandidateIndex], m.Candidates[m.CandidateIndex+1:]...)
			}
			if m.CandidateIndex >= len(m.Candidates) && m.CandidateIndex > 0 {
				m.CandidateIndex--
			}
			m.CandidateNotice = fmt.Sprintf("Deleted candidate %s.", toDelete.ID)
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) initJournalPractice() {
	if m.CurrentInstance != nil {
		m.JournalScenario = m.CurrentInstance.PromptText
		m.JournalFamilyID = m.CurrentInstance.FamilyID
		m.JournalEvent = engine.TransactionEvent{
			FamilyID:   m.CurrentInstance.FamilyID,
			Parameters: m.CurrentInstance.Parameters,
		}
		for i, q := range m.Questions {
			if q.ID == m.CurrentInstance.QuestionID {
				m.JournalQuestionIndex = i
				break
			}
		}
	} else {
		for i, q := range m.Questions {
			if bank.IsActiveForPractice(q.Status) {
				if err := m.loadJournalTemplate(q); err == nil {
					m.JournalQuestionIndex = i
				}
				break
			}
		}
	}
	if m.Engine != nil && m.JournalEvent.FamilyID != "" {
		if processed, err := m.Engine.ProcessEvent(m.JournalEvent); err == nil {
			m.JournalCanonicalEntry = processed.Entry
		}
	}
	m.JournalLines = make([]domain.Posting, 0)
	m.JournalSelectedLine = 0
	m.JournalFeedback = nil
	m.JournalReconciliation = nil
	m.JournalInputActive = false
	m.JournalInputMode = "account"
	m.JournalAccountIdx = 0
	m.JournalSide = domain.SideDebit
	m.JournalAmountBuffer = ""
	m.JournalNotice = "Multi-line Entry Practice: Press [a] to add a line. Balance Dr & Cr, then press [s] to submit."
}

// Instantiate journal wording and grading from the same approved template and allowed amount.
func (m *Model) loadJournalTemplate(q bank.QuestionJSON) error {
	amounts := q.Parameters["amount_minor_units"]
	if len(amounts) == 0 {
		return fmt.Errorf("scenario %s has no allowed amounts", q.ID)
	}
	params := map[string]int64{"amount_minor_units": amounts[m.RNG.Intn(len(amounts))]}
	inst, err := m.Generator.GenerateInstance(q, m.RNG.Int63(), params)
	if err != nil {
		return err
	}
	m.JournalScenario = inst.PromptText
	m.JournalFamilyID = inst.FamilyID
	m.JournalEvent = engine.TransactionEvent{FamilyID: inst.FamilyID, Parameters: inst.Parameters}
	m.JournalCanonicalEntry = inst.Entry
	return nil
}

func parseAmountDollars(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	if s == "" {
		return 0, fmt.Errorf("amount is empty")
	}
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid number format")
	}
	dollarsStr := parts[0]
	if dollarsStr == "" {
		dollarsStr = "0"
	}
	var dollars int64
	if _, err := fmt.Sscanf(dollarsStr, "%d", &dollars); err != nil {
		return 0, err
	}
	if dollars < 0 {
		return 0, fmt.Errorf("negative amount")
	}
	var cents int64
	if len(parts) == 2 {
		centsStr := parts[1]
		if len(centsStr) == 1 {
			centsStr += "0"
		} else if len(centsStr) > 2 {
			centsStr = centsStr[:2]
		}
		if _, err := fmt.Sscanf(centsStr, "%d", &cents); err != nil {
			return 0, err
		}
	}
	totalCents := dollars*100 + cents
	if totalCents <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}
	return totalCents, nil
}

func (m *Model) updateJournalPractice(key string) (tea.Model, tea.Cmd) {
	allAccounts := m.Catalog.All()

	if m.JournalInputActive {
		switch m.JournalInputMode {
		case "account":
			switch key {
			case "up", "k":
				if m.JournalAccountIdx > 0 {
					m.JournalAccountIdx--
				}
				return m, nil
			case "down", "j":
				if m.JournalAccountIdx < len(allAccounts)-1 {
					m.JournalAccountIdx++
				}
				return m, nil
			case "enter":
				m.JournalInputMode = "side"
				return m, nil
			case "esc":
				m.JournalInputActive = false
				return m, nil
			}

		case "side":
			switch key {
			case "D":
				m.JournalSide = domain.SideDebit
				m.JournalInputMode = "amount"
				return m, nil
			case "c", "C":
				m.JournalSide = domain.SideCredit
				m.JournalInputMode = "amount"
				return m, nil
			case "tab":
				if m.JournalSide == domain.SideDebit {
					m.JournalSide = domain.SideCredit
				} else {
					m.JournalSide = domain.SideDebit
				}
				return m, nil
			case "enter":
				m.JournalInputMode = "amount"
				return m, nil
			case "esc":
				m.JournalInputMode = "account"
				return m, nil
			}

		case "amount":
			switch key {
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				m.JournalAmountBuffer += key
				return m, nil
			case ".":
				if !strings.Contains(m.JournalAmountBuffer, ".") {
					m.JournalAmountBuffer += "."
				}
				return m, nil
			case "backspace":
				if len(m.JournalAmountBuffer) > 0 {
					m.JournalAmountBuffer = m.JournalAmountBuffer[:len(m.JournalAmountBuffer)-1]
				}
				return m, nil
			case "enter":
				cents, err := parseAmountDollars(m.JournalAmountBuffer)
				if err != nil {
					m.JournalNotice = fmt.Sprintf("⚠️ Invalid amount: %v", err)
					return m, nil
				}
				acc := allAccounts[m.JournalAccountIdx]
				m.JournalLines = append(m.JournalLines, domain.Posting{
					AccountID: acc.ID,
					Side:      m.JournalSide,
					Amount:    domain.NewMoney(cents),
				})
				m.JournalSelectedLine = len(m.JournalLines) - 1
				m.JournalInputActive = false
				m.JournalAmountBuffer = ""
				m.JournalNotice = fmt.Sprintf("Added line: %s %s %s",
					m.JournalSide, acc.Name, domain.NewMoney(cents).FormatExact())
				return m, nil
			case "esc":
				m.JournalInputMode = "side"
				return m, nil
			}
		}
		return m, nil
	}

	// Normal controls in Journal Practice
	switch key {
	case "q", "esc":
		if m.PreviousState != 0 {
			m.State = m.PreviousState
		} else {
			m.State = StateDrill
		}
		return m, nil

	case "a":
		m.JournalInputActive = true
		m.JournalInputMode = "account"
		m.JournalAmountBuffer = ""
		m.JournalNotice = "Select account with [↑/↓], then press [Enter]."
		return m, nil

	case "up", "k":
		if m.JournalSelectedLine > 0 {
			m.JournalSelectedLine--
		}
		return m, nil

	case "down", "j":
		if m.JournalSelectedLine < len(m.JournalLines)-1 {
			m.JournalSelectedLine++
		}
		return m, nil

	case "D":
		if len(m.JournalLines) > 0 && m.JournalSelectedLine < len(m.JournalLines) {
			m.JournalLines[m.JournalSelectedLine].Side = domain.SideDebit
			m.JournalNotice = "Updated line to Debit."
		}
		return m, nil

	case "c", "C":
		if len(m.JournalLines) > 0 && m.JournalSelectedLine < len(m.JournalLines) {
			m.JournalLines[m.JournalSelectedLine].Side = domain.SideCredit
			m.JournalNotice = "Updated line to Credit."
		}
		return m, nil

	case "x", "delete":
		if len(m.JournalLines) > 0 && m.JournalSelectedLine < len(m.JournalLines) {
			m.JournalLines = append(m.JournalLines[:m.JournalSelectedLine], m.JournalLines[m.JournalSelectedLine+1:]...)
			if m.JournalSelectedLine >= len(m.JournalLines) && m.JournalSelectedLine > 0 {
				m.JournalSelectedLine--
			}
			m.JournalNotice = "Deleted line."
		}
		return m, nil

	case "s", "enter":
		if len(m.JournalLines) < 2 {
			m.JournalNotice = "Entry must have at least 2 lines (at least one Debit and one Credit)."
			return m, nil
		}
		candidateEntry := domain.NewEntry(m.JournalLines...)
		eval := m.Engine.EvaluateEntry(m.JournalEvent, candidateEntry)
		m.JournalFeedback = &eval
		if eval.IsCorrect {
			recon, err := domain.ReconcileTransaction(candidateEntry, m.Catalog)
			if err == nil {
				m.JournalReconciliation = recon
			}
			m.JournalNotice = "✓ Correct entry! All accounts, sides, and amounts match canonical event semantics."
		} else {
			m.JournalReconciliation = nil
			m.JournalNotice = "✗ " + eval.Feedback
		}
		return m, nil

	case "n":
		for offset := 1; offset <= len(m.Questions); offset++ {
			nextIdx := (m.JournalQuestionIndex + offset) % len(m.Questions)
			q := m.Questions[nextIdx]
			if !bank.IsActiveForPractice(q.Status) {
				continue
			}
			if err := m.loadJournalTemplate(q); err != nil {
				m.JournalNotice = "Could not load scenario: " + err.Error()
				return m, nil
			}
			m.JournalQuestionIndex = nextIdx
			m.JournalLines = make([]domain.Posting, 0)
			m.JournalSelectedLine = 0
			m.JournalFeedback = nil
			m.JournalReconciliation = nil
			m.JournalNotice = "Loaded next scenario. Press [a] to add a line."
			break
		}
		return m, nil

	case "clear":
		m.JournalLines = make([]domain.Posting, 0)
		m.JournalSelectedLine = 0
		m.JournalFeedback = nil
		m.JournalReconciliation = nil
		m.JournalNotice = "Cleared all lines."
		return m, nil
	}

	return m, nil
}

func (m *Model) updateStatements(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		if m.PreviousState != 0 {
			m.State = m.PreviousState
		} else {
			m.State = StateDrill
		}
		return m, nil

	case "up", "k":
		if m.StatementsScroll > 0 {
			m.StatementsScroll--
		}
		return m, nil

	case "down", "j":
		m.StatementsScroll++
		return m, nil

	case "pgup", "pageup":
		m.StatementsScroll -= 5
		if m.StatementsScroll < 0 {
			m.StatementsScroll = 0
		}
		return m, nil

	case "pgdown", "pagedown":
		m.StatementsScroll += 5
		return m, nil

	case "g", "home":
		m.StatementsScroll = 0
		return m, nil

	case "G", "end":
		m.StatementsScroll = 9999
		return m, nil
	}
	return m, nil
}

func (m *Model) startNewExam(timeLimit time.Duration) error {
	m.cancelTutor()
	m.SelectedOptionIndex = 0
	m.ExamNotice = ""
	m.ExamReport = nil
	m.ExamReviewIndex = 0

	seed := m.RNG.Int63()
	sessID := fmt.Sprintf("exam-%d", m.Clock.Now().Unix())

	runner, err := exam.NewExamRunner(exam.ExamConfig{
		SessionID:      sessID,
		Questions:      m.Questions,
		Catalog:        m.Catalog,
		Engine:         m.Engine,
		Generator:      m.Generator,
		TotalQuestions: m.TotalQuestions,
		TimeLimit:      timeLimit,
		Seed:           seed,
		Scaffold:       domain.ScaffoldFaded,
		ClockNow:       m.Clock.Now(),
	})
	if err != nil {
		return err
	}

	m.ExamRunner = runner
	m.ExamTimeLimit = timeLimit

	if m.DB != nil {
		_ = m.DB.SaveExamSession(runner.ToSessionRecord())
		for _, eq := range runner.Questions {
			_ = m.DB.SaveQuestionInstance(eq.Instance, sessID)
		}
	}
	return nil
}

func (m *Model) finalizeExam() {
	if m.ExamRunner != nil {
		report, err := m.ExamRunner.Finish(m.Clock.Now())
		if err == nil {
			m.ExamReport = report
			m.ExamReviewIndex = 0
			if m.DB != nil {
				_ = m.DB.CompleteExamSession(
					m.ExamRunner.SessionID,
					*m.ExamRunner.CompletedAt,
					report.Score,
					report.CorrectItems,
					report.TotalItems,
					m.ExamRunner.ElapsedSeconds,
				)
			}
		}
	}
}

func (m *Model) closeExamInterrupted() {
	if m.ExamRunner != nil && !m.ExamRunner.IsCompleted() {
		sessRecord := m.ExamRunner.Interrupt()
		if m.DB != nil {
			_ = m.DB.SaveExamSession(sessRecord)
		}
	}
}

func (m *Model) updateExam(key string) (tea.Model, tea.Cmd) {
	if m.ExamRunner == nil {
		m.State = StateDrill
		return m, nil
	}

	_, stageAns, err := m.ExamRunner.CurrentQuestion()
	if err != nil {
		m.finalizeExam()
		m.State = StateExamSummary
		return m, nil
	}

	numOpts := len(stageAns.Options)

	switch key {
	case "ctrl+c":
		m.closeExamInterrupted()
		m.State = StateQuitting
		return m, tea.Quit

	case "q":
		m.closeExamInterrupted()
		m.State = StateDrill
		return m, nil

	case "?":
		m.ExamNotice = "🔒 Hints are withheld during Exam Mode."
		return m, nil

	case "e":
		m.ExamNotice = "🔒 Explanations are withheld during Exam Mode."
		return m, nil

	case "h", "f1":
		m.ExamNotice = "🔒 Reference cheatsheet is withheld during Exam Mode."
		return m, nil

	case "t":
		m.ExamNotice = "🔒 Tutor settings are withheld during Exam Mode."
		return m, nil

	case "esc":
		m.ExamNotice = ""
		return m, nil

	// Vim / Arrow navigation
	case "up", "k":
		if m.SelectedOptionIndex > 0 {
			m.SelectedOptionIndex--
		} else {
			m.SelectedOptionIndex = numOpts - 1
		}
		m.ExamNotice = ""
		return m, nil

	case "down", "j":
		if m.SelectedOptionIndex < numOpts-1 {
			m.SelectedOptionIndex++
		} else {
			m.SelectedOptionIndex = 0
		}
		m.ExamNotice = ""
		return m, nil

	case "g":
		m.SelectedOptionIndex = 0
		m.ExamNotice = ""
		return m, nil

	case "G":
		if numOpts > 0 {
			m.SelectedOptionIndex = numOpts - 1
		}
		m.ExamNotice = ""
		return m, nil

	case "a", "1":
		if numOpts > 0 {
			m.SelectedOptionIndex = 0
		}
		m.ExamNotice = ""
		return m, nil

	case "b", "2":
		if numOpts > 1 {
			m.SelectedOptionIndex = 1
		}
		m.ExamNotice = ""
		return m, nil

	case "c", "3":
		if numOpts > 2 {
			m.SelectedOptionIndex = 2
		}
		m.ExamNotice = ""
		return m, nil

	case "4":
		if numOpts > 3 {
			m.SelectedOptionIndex = 3
		}
		m.ExamNotice = ""
		return m, nil

	case "enter", " ":
		if m.SelectedOptionIndex < 0 || m.SelectedOptionIndex >= len(stageAns.Options) {
			return m, nil
		}
		chosenOpt := stageAns.Options[m.SelectedOptionIndex]
		now := m.Clock.Now()

		err := m.ExamRunner.SubmitOption(chosenOpt.ID, now)
		if err != nil {
			return m, nil
		}

		// Persist exam attempt to dedicated exam_attempts table
		if m.DB != nil && len(m.ExamRunner.Attempts) > 0 {
			latest := m.ExamRunner.Attempts[len(m.ExamRunner.Attempts)-1]
			_ = m.DB.RecordExamAttempt(latest)
			_ = m.DB.SaveExamSession(m.ExamRunner.ToSessionRecord())
		}

		m.SelectedOptionIndex = 0
		m.ExamNotice = ""

		if m.ExamRunner.IsCompleted() {
			m.finalizeExam()
			m.State = StateExamSummary
			return m, nil
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) updateExamSummary(key string) (tea.Model, tea.Cmd) {
	if m.ExamReport == nil {
		m.State = StateDrill
		return m, nil
	}

	numReviews := len(m.ExamReport.Reviews)

	switch key {
	case "q", "esc":
		m.State = StateDrill
		return m, nil

	case "s":
		m.PreviousState = StateExamSummary
		m.State = StateMastery
		return m, nil

	case "r":
		_ = m.startNewExam(m.ExamTimeLimit)
		m.State = StateExam
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })

	// Navigate reviews
	case "]", "n", "down", "j":
		if m.ExamReviewIndex < numReviews-1 {
			m.ExamReviewIndex++
		} else {
			m.ExamReviewIndex = 0
		}
		return m, nil

	case "[", "p", "up", "k":
		if m.ExamReviewIndex > 0 {
			m.ExamReviewIndex--
		} else if numReviews > 0 {
			m.ExamReviewIndex = numReviews - 1
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) updateExamResumePrompt(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "r", "enter":
		if m.InterruptedExam != nil && m.DB != nil {
			atts, err := m.DB.GetExamAttempts(m.InterruptedExam.ID)
			if err == nil {
				snapshots, err := m.DB.GetExamQuestionSnapshots(m.InterruptedExam.ID)
				var resumed *exam.ExamRunner
				if err == nil {
					resumed, err = exam.ResumeExamSnapshots(*m.InterruptedExam, atts, snapshots)
				}
				if err == nil {
					m.ExamRunner = resumed
					m.State = StateExam
					m.InterruptedExam = nil
					return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })
				}
			}
		}
		m.ExamNotice = "Could not restore the original exam snapshots. The interrupted session is preserved; choose a new exam or return to practice."
		return m, nil

	case "a", "n":
		if m.InterruptedExam != nil && m.DB != nil {
			_ = m.DB.AbandonExamSession(m.InterruptedExam.ID)
		}
		m.InterruptedExam = nil
		_ = m.startNewExam(m.ExamTimeLimit)
		m.State = StateExam
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })

	case "q", "esc":
		m.State = StateDrill
		return m, nil
	}

	return m, nil
}

func (m *Model) initIntro() {
	m.Intro = NewIntroState(m.Width, m.Height, m.RNG)
	if m.DB != nil {
		top, err := m.DB.GetTopArcadeHighScore()
		if err == nil && top != nil {
			m.Intro.SetHighScore(top.Score, top.Initials)
		}
	}
}

func (m *Model) saveArcadeScore() {
	if m.Intro == nil || m.DB == nil || m.Intro.ScoreSaved {
		return
	}
	initials := m.Intro.GetInitialsString()
	if initials == "" {
		initials = "AAA"
	}
	_ = m.DB.SaveArcadeHighScore(storage.ArcadeHighScore{
		Initials:        initials,
		Score:           m.Intro.Score,
		BlastedCount:    m.Intro.BlastedCount,
		SurvivalSeconds: m.Intro.TickCount / 30,
	})
	m.Intro.ScoreSaved = true
	m.Intro.SetHighScore(m.Intro.Score, initials)
}

func (m *Model) arcadeKeysAllowed() bool {
	switch m.State {
	case StateDrill, StateFeedback, StateRecap, StateMastery, StateHelp, StateSessionComplete, StateStatements:
		return true
	}
	return false
}

// launchArcade replays the startup game from the current screen, optionally
// opening straight onto the Hall of Fame.
func (m *Model) launchArcade(showScores bool) (tea.Model, tea.Cmd) {
	m.cancelTutor()
	m.IntroReturn = m.State
	m.initIntro()
	m.State = StateIntro
	if showScores {
		m.openLeaderboard()
		m.Intro.LeaderboardOnly = true
	}
	return m, introTick()
}

func (m *Model) openLeaderboard() {
	m.Intro.Leaderboard = nil
	if m.DB != nil {
		if top, err := m.DB.ListTopArcadeHighScores(10); err == nil {
			m.Intro.Leaderboard = top
		}
	}
	m.Intro.ShowLeaderboard = true
}

func (m *Model) exitIntro() (tea.Model, tea.Cmd) {
	if m.Intro != nil && m.Intro.LeaderboardOnly {
		return m.finishIntro()
	}
	m.State = StateTitle
	m.TitleFrame = 0
	return m, titleTick()
}

func (m *Model) finishIntro() (tea.Model, tea.Cmd) {
	targetState := m.IntroReturn
	if targetState == StateIntro || targetState == 0 {
		targetState = StateDrill
	}
	m.State = targetState
	if m.State == StateExam {
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return examTickMsg{} })
	}
	return m, nil
}

func (m *Model) updateIntro(key string) (tea.Model, tea.Cmd) {
	if m.Intro == nil {
		return m, nil
	}

	// 0. Hall of Fame overlay (flight is paused underneath)
	if m.Intro.ShowLeaderboard {
		switch key {
		case "L", "l", "esc":
			if m.Intro.LeaderboardOnly {
				return m.exitIntro()
			}
			m.Intro.ShowLeaderboard = false
		case "r", "R":
			m.initIntro()
		case "enter":
			m.Intro.ShowLeaderboard = false
			if !m.Intro.LeaderboardOnly {
				return m.updateIntro("enter")
			}
			return m.exitIntro()
		case "q":
			m.State = StateQuitting
			return m, tea.Quit
		}
		return m, nil
	}

	// 1. Retro 3-Initials Entry Modal
	if m.Intro.InitialsEntryActive {
		switch key {
		case "up", "w", "k":
			m.Intro.CycleInitial(1)
			return m, nil
		case "down", "s", "j":
			m.Intro.CycleInitial(-1)
			return m, nil
		case "left", "h", "backspace":
			m.Intro.PrevInitial()
			return m, nil
		case "right", "l":
			m.Intro.NextInitial()
			return m, nil
		case "enter", "space":
			done := m.Intro.NextInitial()
			if done {
				m.saveArcadeScore()
				m.Intro.InitialsEntryActive = false
				return m.exitIntro()
			}
			return m, nil
		case "esc":
			m.saveArcadeScore()
			m.Intro.InitialsEntryActive = false
			return m.exitIntro()
		default:
			if len(key) == 1 {
				ch := rune(key[0])
				if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
					m.Intro.SetInitial(ch)
					return m, nil
				}
			}
			return m, nil
		}
	}

	// 2. Fiscal year closed: continue for a higher score, or retire as a win.
	if m.Intro.YearComplete {
		switch key {
		case "y", "Y":
			m.Intro.ContinueYear()
		case "n", "N":
			m.Intro.RetireAfterYear()
		case "L", "l":
			m.openLeaderboard()
		case "q":
			m.State = StateQuitting
			return m, tea.Quit
		}
		return m, nil
	}

	// 3. Pause overlay
	if m.Intro.Paused {
		switch key {
		case "p", "P":
			m.Intro.TogglePause()
		case "L", "l":
			m.openLeaderboard()
		case "enter", "esc":
			return m.leaveFlight()
		case "q":
			m.State = StateQuitting
			return m, tea.Quit
		}
		return m, nil
	}

	// 4. Game Over Modal Handling
	if m.Intro.GameOver {
		switch key {
		case "r", "R":
			m.initIntro()
			return m, nil
		case "L", "l":
			m.openLeaderboard()
			return m, nil
		case "enter", "esc":
			return m.leaveFlight()
		case "q":
			m.State = StateQuitting
			return m, tea.Quit
		default:
			return m, nil
		}
	}

	// 5. Active Flight & Combat
	switch key {
	case "enter", "esc":
		return m.leaveFlight()

	case "q":
		m.State = StateQuitting
		return m, tea.Quit

	case "up", "w", "k":
		m.Intro.SteerUp()
		return m, nil

	case "down", "s", "j":
		m.Intro.SteerDown()
		return m, nil

	case "b", "B":
		m.Intro.DeployBomb()
		return m, nil

	case "f", " ", "space":
		m.Intro.Fire()
		return m, nil

	case "g", "G":
		m.Intro.ToggleAutoFire()
		return m, nil

	case "p", "P":
		m.Intro.TogglePause()
		return m, nil

	case "L", "l":
		m.openLeaderboard()
		return m, nil

	default:
		return m, nil
	}
}

// leaveFlight heads to the drills, first asking for initials when the run set
// a new all-time high score.
func (m *Model) leaveFlight() (tea.Model, tea.Cmd) {
	if m.Intro.Score > m.Intro.HighScore && m.Intro.Score > 0 && !m.Intro.ScoreSaved {
		m.Intro.Paused = false
		m.Intro.InitialsEntryActive = true
		return m, nil
	}
	return m.exitIntro()
}

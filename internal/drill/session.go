package drill

import (
	"fmt"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// ScaffoldStageSequence returns the pedagogical sequence of stages based on the scaffold fading level.
// Level 0 (Full): all 7 stages
// Level 1 (Intermediate): 4 stages (Identify, Counter, Balanced, Equation)
// Level 2 (Faded): 2 stages (Balanced, Equation)
func ScaffoldStageSequence(level domain.ScaffoldLevel) []domain.DrillStage {
	switch level {
	case domain.ScaffoldIntermediate:
		return []domain.DrillStage{
			domain.StageIdentifyAccount,
			domain.StageCounterAccount,
			domain.StageBalancedEntry,
			domain.StageEquationEffect,
		}
	case domain.ScaffoldFaded:
		return []domain.DrillStage{
			domain.StageBalancedEntry,
			domain.StageEquationEffect,
		}
	default: // ScaffoldFull
		return []domain.DrillStage{
			domain.StageIdentifyAccount,
			domain.StageAccountCategory,
			domain.StageDirection,
			domain.StageDebitCredit,
			domain.StageCounterAccount,
			domain.StageBalancedEntry,
			domain.StageEquationEffect,
		}
	}
}

// SessionState tracks the runtime progress of a learner through an instantiated progressive drill.
type SessionState struct {
	SessionID          string
	Instance           *domain.QuestionInstance
	StageSequence      []domain.DrillStage
	CurrentIndex       int
	RetryAvailable     bool
	ReferenceConsulted bool
	CurrentAssistance  domain.AssistanceLevel
	Attempts           []domain.Attempt
	IsCompleted        bool
}

// SubmitFeedback represents the result returned to the UI after answering a stage.
type SubmitFeedback struct {
	IsCorrect       bool
	AdvanceStage    bool
	SelectedOption  domain.AnswerOption
	Hint            string
	Explanation     string
	ErrorTag        string
	AssistanceLevel domain.AssistanceLevel
	StageRecap      *domain.StageAnswer
}

func NewSession(sessionID string, inst *domain.QuestionInstance) *SessionState {
	candidateSeq := ScaffoldStageSequence(inst.ScaffoldLevel)

	// Filter stages present in instance
	var activeSeq []domain.DrillStage
	for _, st := range candidateSeq {
		if _, ok := inst.StageAnswers[st]; ok {
			activeSeq = append(activeSeq, st)
		}
	}

	return &SessionState{
		SessionID:          sessionID,
		Instance:           inst,
		StageSequence:      activeSeq,
		CurrentIndex:       0,
		RetryAvailable:     true,
		ReferenceConsulted: false,
		CurrentAssistance:  domain.AssistanceNone,
		Attempts:           make([]domain.Attempt, 0),
		IsCompleted:        len(activeSeq) == 0,
	}
}

// CurrentStage returns the current active stage definition.
func (s *SessionState) CurrentStage() (*domain.StageAnswer, error) {
	if s.IsCompleted || s.CurrentIndex >= len(s.StageSequence) {
		return nil, fmt.Errorf("drill session is already completed")
	}
	stageKey := s.StageSequence[s.CurrentIndex]
	stage, ok := s.Instance.StageAnswers[stageKey]
	if !ok {
		return nil, fmt.Errorf("stage %s not found in instance", stageKey)
	}
	return &stage, nil
}

// RequestHint transitions the assistance level to Hinted and returns the causal hint.
func (s *SessionState) RequestHint() (string, error) {
	stage, err := s.CurrentStage()
	if err != nil {
		return "", err
	}
	if s.CurrentAssistance == domain.AssistanceNone {
		s.CurrentAssistance = domain.AssistanceHinted
	}
	return stage.CausalHint, nil
}

// RecordReferenceUse flags that the reference cheatsheet was consulted during this stage.
func (s *SessionState) RecordReferenceUse() {
	s.ReferenceConsulted = true
	if s.CurrentAssistance == domain.AssistanceNone {
		s.CurrentAssistance = domain.AssistanceReference
	}
}

// RequestExplanation transitions the assistance level to Hinted if currently unassisted.
func (s *SessionState) RequestExplanation() {
	if s.CurrentAssistance == domain.AssistanceNone {
		s.CurrentAssistance = domain.AssistanceHinted
	}
}

// SubmitOption processes a learner answer for the current stage.
func (s *SessionState) SubmitOption(selectedOptionID string, answeredAt time.Time) (*SubmitFeedback, error) {
	stage, err := s.CurrentStage()
	if err != nil {
		return nil, err
	}

	// Find the selected option by stable ID
	var selectedOpt domain.AnswerOption
	found := false
	for _, opt := range stage.Options {
		if opt.ID == selectedOptionID {
			selectedOpt = opt
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("invalid option ID %q for stage %s", selectedOptionID, stage.Stage)
	}

	isCorrect := (selectedOptionID == stage.CorrectOptionID)
	attemptNumber := len(s.Attempts) + 1

	refUsed := s.ReferenceConsulted
	assist := s.CurrentAssistance
	if refUsed && assist == domain.AssistanceNone {
		assist = domain.AssistanceReference
	}

	attempt := domain.Attempt{
		AttemptID:        fmt.Sprintf("%s-att-%d", s.Instance.InstanceID, attemptNumber),
		SessionID:        s.SessionID,
		InstanceID:       s.Instance.InstanceID,
		QuestionID:       s.Instance.QuestionID,
		QuestionVersion:  s.Instance.Version,
		Stage:            stage.Stage,
		ConceptID:        stage.RelevantConceptID,
		SelectedOptionID: selectedOptionID,
		IsCorrect:        isCorrect,
		Assistance:       assist,
		ReferenceUsed:    refUsed,
		ErrorTag:         selectedOpt.ErrorTag,
		GradingVersion:   1,
		AnsweredAt:       answeredAt,
	}
	s.Attempts = append(s.Attempts, attempt)

	feedback := &SubmitFeedback{
		IsCorrect:       isCorrect,
		SelectedOption:  selectedOpt,
		ErrorTag:        selectedOpt.ErrorTag,
		AssistanceLevel: assist,
	}

	if isCorrect {
		feedback.AdvanceStage = true
		feedback.Explanation = stage.Explanation
		s.advanceToNextStage()
		return feedback, nil
	}

	// Incorrect answer
	if s.RetryAvailable {
		// First error: provide targeted causal hint, allow one retry
		s.RetryAvailable = false
		s.CurrentAssistance = domain.AssistanceRetry
		feedback.AdvanceStage = false
		feedback.Hint = stage.CausalHint
		if selectedOpt.ErrorTag != "" {
			feedback.Explanation = fmt.Sprintf("Incorrect. %s", stage.CausalHint)
		} else {
			feedback.Explanation = fmt.Sprintf("Incorrect. Hint: %s", stage.CausalHint)
		}
		return feedback, nil
	}

	// Second error (retry exhausted): reveal explanation and advance
	feedback.AdvanceStage = true
	feedback.Explanation = fmt.Sprintf("Incorrect on retry. The correct answer was: %s. Explanation: %s",
		stage.CorrectOptionID, stage.Explanation)
	s.CurrentAssistance = domain.AssistanceRevealed
	s.advanceToNextStage()
	return feedback, nil
}

func (s *SessionState) advanceToNextStage() {
	s.CurrentIndex++
	s.RetryAvailable = true
	s.ReferenceConsulted = false
	s.CurrentAssistance = domain.AssistanceNone
	if s.CurrentIndex >= len(s.StageSequence) {
		s.IsCompleted = true
	}
}

// TransactionRecap returns the full transaction summary and postings after completing the drill.
type TransactionRecap struct {
	Prompt     string
	Entry      domain.Entry
	IsBalanced bool
	Postings   []domain.Posting
}

func (s *SessionState) Recap() TransactionRecap {
	return TransactionRecap{
		Prompt:     s.Instance.PromptText,
		Entry:      s.Instance.Entry,
		IsBalanced: s.Instance.Entry.IsBalanced(),
		Postings:   s.Instance.Entry.Postings,
	}
}

package exam

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// ExamConfig provides parameters to initialize an ExamRunner.
type ExamConfig struct {
	SessionID      string
	Questions      []bank.QuestionJSON
	Catalog        *domain.AccountCatalog
	Engine         *engine.Engine
	Generator      *drill.Generator
	TotalQuestions int
	TimeLimit      time.Duration
	Seed           int64
	Scaffold       domain.ScaffoldLevel
	ClockNow       time.Time
}

// ExamRunner manages the execution of an assessment session adhering strictly
// to Stage 16 invariants: hints, references, and feedback are suppressed until finish.
type ExamRunner struct {
	SessionID      string
	Seed           int64
	TotalQuestions int
	TimeLimit      time.Duration
	ElapsedSeconds int
	Status         ExamStatus
	StartedAt      time.Time
	CompletedAt    *time.Time
	Questions      []*ExamQuestion
	CurrentQIndex  int
	Attempts       []ExamAttemptRecord
	TimedOut       bool
}

// NewExamRunner instantiates a deterministic exam session from approved bank questions.
func NewExamRunner(cfg ExamConfig) (*ExamRunner, error) {
	if len(cfg.Questions) == 0 {
		return nil, fmt.Errorf("cannot create exam with empty question bank")
	}
	if cfg.TotalQuestions <= 0 {
		cfg.TotalQuestions = 10
	}
	if cfg.TotalQuestions > len(cfg.Questions) {
		cfg.TotalQuestions = len(cfg.Questions)
	}
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
	if cfg.SessionID == "" {
		cfg.SessionID = fmt.Sprintf("exam-%d", time.Now().Unix())
	}
	if cfg.ClockNow.IsZero() {
		cfg.ClockNow = time.Now().UTC()
	}

	rng := rand.New(rand.NewSource(cfg.Seed))

	// Shuffle questions deterministically
	shuffled := make([]bank.QuestionJSON, len(cfg.Questions))
	copy(shuffled, cfg.Questions)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	selectedTemplates := shuffled[:cfg.TotalQuestions]
	examQuestions := make([]*ExamQuestion, 0, len(selectedTemplates))

	for i, q := range selectedTemplates {
		// Select parameter amount
		amtOptions := q.Parameters["amount_minor_units"]
		var chosenAmt int64 = 10000 // $100 default
		if len(amtOptions) > 0 {
			chosenAmt = amtOptions[rng.Intn(len(amtOptions))]
		}
		qSeed := rng.Int63()
		params := map[string]int64{"amount_minor_units": chosenAmt}

		// In Exam mode, we default to ScaffoldFaded (balanced entry + equation effect)
		// which directly tests journal entry synthesis and accounting equation impact.
		scaffold := cfg.Scaffold
		if scaffold == 0 {
			scaffold = domain.ScaffoldFaded
		}

		inst, err := cfg.Generator.GenerateInstanceWithScaffold(q, qSeed, params, scaffold)
		if err != nil {
			return nil, fmt.Errorf("failed generating exam question instance %d (%s): %w", i+1, q.ID, err)
		}

		candidateSeq := drill.ScaffoldStageSequence(scaffold)
		var activeStages []domain.DrillStage
		for _, st := range candidateSeq {
			if _, ok := inst.StageAnswers[st]; ok {
				activeStages = append(activeStages, st)
			}
		}

		examQuestions = append(examQuestions, &ExamQuestion{
			QuestionIndex: i,
			Instance:      inst,
			Stages:        activeStages,
			CurrentStage:  0,
			Responses:     make(map[domain.DrillStage]*ExamStageResponse),
			IsCompleted:   len(activeStages) == 0,
		})
	}

	return &ExamRunner{
		SessionID:      cfg.SessionID,
		Seed:           cfg.Seed,
		TotalQuestions: len(examQuestions),
		TimeLimit:      cfg.TimeLimit,
		ElapsedSeconds: 0,
		Status:         ExamStatusInProgress,
		StartedAt:      cfg.ClockNow,
		Questions:      examQuestions,
		CurrentQIndex:  0,
		Attempts:       make([]ExamAttemptRecord, 0),
	}, nil
}

// CurrentQuestion returns the active question and active stage for the student.
func (r *ExamRunner) CurrentQuestion() (*ExamQuestion, *domain.StageAnswer, error) {
	if r.IsCompleted() {
		return nil, nil, ErrExamAlreadyCompleted
	}
	if r.CurrentQIndex >= len(r.Questions) {
		return nil, nil, ErrExamAlreadyCompleted
	}

	eq := r.Questions[r.CurrentQIndex]
	if eq.CurrentStage >= len(eq.Stages) {
		return nil, nil, fmt.Errorf("question %d stages exhausted", r.CurrentQIndex+1)
	}

	stageKey := eq.Stages[eq.CurrentStage]
	stageAns, ok := eq.Instance.StageAnswers[stageKey]
	if !ok {
		return nil, nil, fmt.Errorf("stage %s not found in instance %s", stageKey, eq.Instance.InstanceID)
	}

	return eq, &stageAns, nil
}

// SubmitOption processes the student's answer for the current stage.
// Crucially, it records the choice and immediately advances WITHOUT revealing correctness,
// hints, explanations, or recap.
func (r *ExamRunner) SubmitOption(selectedOptionID string, answeredAt time.Time) error {
	if r.IsCompleted() {
		return ErrExamAlreadyCompleted
	}

	eq, stageAns, err := r.CurrentQuestion()
	if err != nil {
		return err
	}

	// Validate option exists
	var selectedOpt domain.AnswerOption
	found := false
	for _, opt := range stageAns.Options {
		if opt.ID == selectedOptionID {
			selectedOpt = opt
			found = true
			break
		}
	}
	if !found {
		return ErrInvalidOptionID
	}

	isCorrect := (selectedOptionID == stageAns.CorrectOptionID)

	resp := &ExamStageResponse{
		Stage:            stageAns.Stage,
		ConceptID:        stageAns.RelevantConceptID,
		SelectedOptionID: selectedOptionID,
		CorrectOptionID:  stageAns.CorrectOptionID,
		IsCorrect:        isCorrect,
		ErrorTag:         selectedOpt.ErrorTag,
		AnsweredAt:       answeredAt,
	}
	eq.Responses[stageAns.Stage] = resp

	attemptID := fmt.Sprintf("%s-q%d-%s-att%d", r.SessionID, eq.QuestionIndex+1, string(stageAns.Stage), len(r.Attempts)+1)
	attRecord := ExamAttemptRecord{
		ID:               attemptID,
		ExamSessionID:    r.SessionID,
		InstanceID:       eq.Instance.InstanceID,
		QuestionID:       eq.Instance.QuestionID,
		QuestionIndex:    eq.QuestionIndex,
		Stage:            stageAns.Stage,
		ConceptID:        stageAns.RelevantConceptID,
		SelectedOptionID: selectedOptionID,
		CorrectOptionID:  stageAns.CorrectOptionID,
		IsCorrect:        isCorrect,
		ErrorTag:         selectedOpt.ErrorTag,
		GradingVersion:   1,
		AnsweredAt:       answeredAt,
	}
	r.Attempts = append(r.Attempts, attRecord)

	// Advance stage or question
	eq.CurrentStage++
	if eq.CurrentStage >= len(eq.Stages) {
		eq.IsCompleted = true
		r.CurrentQIndex++
		if r.CurrentQIndex >= len(r.Questions) {
			r.Status = ExamStatusCompleted
		}
	}

	return nil
}

// RequestHint enforces the Stage 16 invariant suppressing hints during exam mode.
func (r *ExamRunner) RequestHint() error {
	return ErrHintsSuppressedInExam
}

// RequestExplanation enforces the Stage 16 invariant suppressing explanations during exam mode.
func (r *ExamRunner) RequestExplanation() error {
	return ErrExplanationsSuppressedInExam
}

// RecordReferenceUse enforces the Stage 16 invariant suppressing reference cheatsheets.
func (r *ExamRunner) RecordReferenceUse() error {
	return ErrReferenceSuppressedInExam
}

// Tick advances the timer and automatically finalizes the exam if time limit expires.
func (r *ExamRunner) Tick(delta time.Duration) bool {
	if r.IsCompleted() {
		return false
	}
	r.ElapsedSeconds += int(delta.Seconds())

	if r.TimeLimit > 0 && time.Duration(r.ElapsedSeconds)*time.Second >= r.TimeLimit {
		r.TimedOut = true
		r.Status = ExamStatusCompleted
		return true // expired
	}
	return false
}

// IsCompleted reports whether the exam is completed (all questions submitted or timed out).
func (r *ExamRunner) IsCompleted() bool {
	return r.Status == ExamStatusCompleted || r.Status == ExamStatusAbandoned
}

// Remaining returns the remaining time if timed, or 0 if untimed.
func (r *ExamRunner) Remaining() time.Duration {
	if r.TimeLimit <= 0 {
		return 0
	}
	elapsed := time.Duration(r.ElapsedSeconds) * time.Second
	rem := r.TimeLimit - elapsed
	if rem < 0 {
		return 0
	}
	return rem
}

// Interrupt marks the exam session as interrupted for explicit resumption.
func (r *ExamRunner) Interrupt() ExamSessionRecord {
	r.Status = ExamStatusInterrupted
	return r.ToSessionRecord()
}

// Finish grades all answers, sets completion time, and produces the complete ExamReport.
func (r *ExamRunner) Finish(now time.Time) (*ExamReport, error) {
	if r.Status != ExamStatusCompleted {
		r.Status = ExamStatusCompleted
	}
	r.CompletedAt = &now
	return r.GenerateReport(), nil
}

// ToSessionRecord converts runner state into a persistable ExamSessionRecord.
func (r *ExamRunner) ToSessionRecord() ExamSessionRecord {
	var compTime *time.Time
	if r.CompletedAt != nil {
		t := *r.CompletedAt
		compTime = &t
	}

	correctCount := 0
	for _, att := range r.Attempts {
		if att.IsCorrect {
			correctCount++
		}
	}
	var score float64
	if len(r.Attempts) > 0 {
		score = (float64(correctCount) / float64(len(r.Attempts))) * 100.0
	}

	var timeLimitSec int
	if r.TimeLimit > 0 {
		timeLimitSec = int(r.TimeLimit.Seconds())
	}

	return ExamSessionRecord{
		ID:               r.SessionID,
		TotalQuestions:   r.TotalQuestions,
		TimeLimitSeconds: timeLimitSec,
		ElapsedSeconds:   r.ElapsedSeconds,
		Status:           r.Status,
		StartedAt:        r.StartedAt,
		CompletedAt:      compTime,
		Score:            score,
		CorrectCount:     correctCount,
		TotalAttempts:    len(r.Attempts),
		Seed:             r.Seed,
	}
}

// GenerateReport compiles complete skill summaries, error patterns, and unmasked reviews.
func (r *ExamRunner) GenerateReport() *ExamReport {
	totalItems := 0
	correctItems := 0

	type conceptStat struct {
		name    string
		total   int
		correct int
	}
	conceptMap := make(map[string]*conceptStat)
	errorCounts := make(map[string]int)
	reviews := make([]QuestionReviewItem, 0)

	for _, eq := range r.Questions {
		for _, st := range eq.Stages {
			stageAns, ok := eq.Instance.StageAnswers[st]
			if !ok {
				continue
			}
			totalItems++

			resp := eq.Responses[st]
			isCorr := false
			var selectedOpt *domain.AnswerOption
			errTag := ""

			if resp != nil {
				isCorr = resp.IsCorrect
				errTag = resp.ErrorTag
				if isCorr {
					correctItems++
				}
				for _, opt := range stageAns.Options {
					if opt.ID == resp.SelectedOptionID {
						selectedOpt = &opt
						break
					}
				}
			}

			// Aggregate skill stats
			cid := stageAns.RelevantConceptID
			if cid == "" {
				cid = "general_accounting"
			}
			cs, ok := conceptMap[cid]
			if !ok {
				cs = &conceptStat{
					name: ConceptDisplayName(cid),
				}
				conceptMap[cid] = cs
			}
			cs.total++
			if isCorr {
				cs.correct++
			}

			// Aggregate error tags
			if errTag != "" {
				errorCounts[errTag]++
			}

			// Find correct option object
			var correctOpt domain.AnswerOption
			for _, opt := range stageAns.Options {
				if opt.ID == stageAns.CorrectOptionID {
					correctOpt = opt
					break
				}
			}

			reviews = append(reviews, QuestionReviewItem{
				QuestionIndex:  eq.QuestionIndex,
				InstanceID:     eq.Instance.InstanceID,
				FamilyID:       eq.Instance.FamilyID,
				PromptText:     eq.Instance.PromptText,
				Stage:          st,
				StagePrompt:    stageAns.Prompt,
				SelectedOption: selectedOpt,
				CorrectOption:  correctOpt,
				IsCorrect:      isCorr,
				ErrorTag:       errTag,
				Explanation:    stageAns.Explanation,
				Postings:       eq.Instance.Entry.Postings,
			})
		}
	}

	var scorePct float64
	if totalItems > 0 {
		scorePct = (float64(correctItems) / float64(totalItems)) * 100.0
	}

	// Grade calculation
	grade := "F"
	switch {
	case scorePct >= 90.0:
		grade = "A"
	case scorePct >= 80.0:
		grade = "B"
	case scorePct >= 70.0:
		grade = "C"
	case scorePct >= 60.0:
		grade = "D"
	default:
		grade = "F"
	}

	// Build skills summary list
	skills := make([]SkillSummary, 0, len(conceptMap))
	for cid, cs := range conceptMap {
		acc := 0.0
		if cs.total > 0 {
			acc = (float64(cs.correct) / float64(cs.total)) * 100.0
		}
		status := "Critical"
		if acc >= 80.0 {
			status = "Strong"
		} else if acc >= 50.0 {
			status = "Needs Review"
		}

		skills = append(skills, SkillSummary{
			ConceptID:   cid,
			ConceptName: cs.name,
			Total:       cs.total,
			Correct:     cs.correct,
			Accuracy:    acc,
			Status:      status,
		})
	}

	// Build error summaries
	errorsList := make([]ErrorSummary, 0, len(errorCounts))
	for tag, cnt := range errorCounts {
		misc, rem := DistractorExplanation(tag)
		errorsList = append(errorsList, ErrorSummary{
			Tag:           tag,
			Count:         cnt,
			Misconception: misc,
			Remediation:   rem,
		})
	}

	var compTime time.Time
	if r.CompletedAt != nil {
		compTime = *r.CompletedAt
	} else {
		compTime = r.StartedAt.Add(time.Duration(r.ElapsedSeconds) * time.Second)
	}

	duration := time.Duration(r.ElapsedSeconds) * time.Second

	return &ExamReport{
		SessionID:        r.SessionID,
		Status:           r.Status,
		StartedAt:        r.StartedAt,
		CompletedAt:      compTime,
		Duration:         duration,
		TimeLimitSeconds: int(r.TimeLimit.Seconds()),
		TimedOut:         r.TimedOut,
		TotalQuestions:   r.TotalQuestions,
		TotalItems:       totalItems,
		CorrectItems:     correctItems,
		Score:            scorePct,
		Grade:            grade,
		Skills:           skills,
		Errors:           errorsList,
		Reviews:          reviews,
	}
}

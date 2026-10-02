package exam

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// ResumeExamRunner restores an interrupted exam runner from persisted storage.
func ResumeExamRunner(
	record ExamSessionRecord,
	attempts []ExamAttemptRecord,
	allQuestions []bank.QuestionJSON,
	catalog *domain.AccountCatalog,
	eng *engine.Engine,
	gen *drill.Generator,
) (*ExamRunner, error) {
	if len(allQuestions) == 0 {
		return nil, fmt.Errorf("cannot resume exam with empty question bank")
	}

	rng := rand.New(rand.NewSource(record.Seed))

	// Recreate deterministic questions sequence
	shuffled := make([]bank.QuestionJSON, len(allQuestions))
	copy(shuffled, allQuestions)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	totQ := record.TotalQuestions
	if totQ > len(shuffled) {
		totQ = len(shuffled)
	}
	selectedTemplates := shuffled[:totQ]

	examQuestions := make([]*ExamQuestion, 0, len(selectedTemplates))
	for i, q := range selectedTemplates {
		amtOptions := q.Parameters["amount_minor_units"]
		var chosenAmt int64 = 10000
		if len(amtOptions) > 0 {
			chosenAmt = amtOptions[rng.Intn(len(amtOptions))]
		}
		qSeed := rng.Int63()
		params := map[string]int64{"amount_minor_units": chosenAmt}

		inst, err := gen.GenerateInstanceWithScaffold(q, qSeed, params, domain.ScaffoldFaded)
		if err != nil {
			return nil, fmt.Errorf("failed reconstructing question %d: %w", i+1, err)
		}

		candidateSeq := drill.ScaffoldStageSequence(domain.ScaffoldFaded)
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
			IsCompleted:   false,
		})
	}

	// Replay recorded attempts
	attemptsByInstanceAndStage := make(map[string]map[domain.DrillStage]ExamAttemptRecord)
	for _, att := range attempts {
		if _, ok := attemptsByInstanceAndStage[att.InstanceID]; !ok {
			attemptsByInstanceAndStage[att.InstanceID] = make(map[domain.DrillStage]ExamAttemptRecord)
		}
		attemptsByInstanceAndStage[att.InstanceID][att.Stage] = att
	}

	currentQIdx := 0
	for qIdx, eq := range examQuestions {
		stMap, hasAtts := attemptsByInstanceAndStage[eq.Instance.InstanceID]
		if !hasAtts {
			currentQIdx = qIdx
			break
		}

		completedStages := 0
		for _, st := range eq.Stages {
			att, ok := stMap[st]
			if ok {
				eq.Responses[st] = &ExamStageResponse{
					Stage:            st,
					ConceptID:        att.ConceptID,
					SelectedOptionID: att.SelectedOptionID,
					CorrectOptionID:  att.CorrectOptionID,
					IsCorrect:        att.IsCorrect,
					ErrorTag:         att.ErrorTag,
					AnsweredAt:       att.AnsweredAt,
				}
				completedStages++
			}
		}

		eq.CurrentStage = completedStages
		if completedStages >= len(eq.Stages) {
			eq.IsCompleted = true
			currentQIdx = qIdx + 1
		} else {
			currentQIdx = qIdx
			break
		}
	}

	var timeLimit time.Duration
	if record.TimeLimitSeconds > 0 {
		timeLimit = time.Duration(record.TimeLimitSeconds) * time.Second
	}

	status := ExamStatusInProgress
	if currentQIdx >= len(examQuestions) {
		status = ExamStatusCompleted
	}

	return &ExamRunner{
		SessionID:      record.ID,
		Seed:           record.Seed,
		TotalQuestions: len(examQuestions),
		TimeLimit:      timeLimit,
		ElapsedSeconds: record.ElapsedSeconds,
		Status:         status,
		StartedAt:      record.StartedAt,
		Questions:      examQuestions,
		CurrentQIndex:  currentQIdx,
		Attempts:       attempts,
	}, nil
}

// FormatText creates a comprehensive human-readable console report for the exam results.
func (r *ExamReport) FormatText() string {
	var b strings.Builder

	b.WriteString("==================================================================================\n")
	b.WriteString(fmt.Sprintf("ACCTG PRACTICE — EXAM ASSESSMENT REPORT\n"))
	b.WriteString(fmt.Sprintf("Session ID:   %s\n", r.SessionID))
	b.WriteString(fmt.Sprintf("Started At:   %s\n", r.StartedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("Completed At: %s\n", r.CompletedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("Duration:     %s", r.Duration.Round(time.Second)))
	if r.TimeLimitSeconds > 0 {
		b.WriteString(fmt.Sprintf(" (Time Limit: %s)", (time.Duration(r.TimeLimitSeconds) * time.Second).String()))
	}
	if r.TimedOut {
		b.WriteString(" [TIMED OUT]")
	}
	b.WriteString("\n")
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString(fmt.Sprintf("FINAL SCORE:  %d / %d items (%.1f%%) — Grade: %s\n", r.CorrectItems, r.TotalItems, r.Score, r.Grade))
	b.WriteString("==================================================================================\n\n")

	// 1. Skill breakdown table
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString("1. PERFORMANCE BREAKDOWN BY SKILL & CONCEPT\n")
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString(fmt.Sprintf("%-38s %-12s %-10s %-12s\n", "Skill / Concept", "Score", "Accuracy", "Status"))
	b.WriteString("----------------------------------------------------------------------------------\n")
	for _, sk := range r.Skills {
		scoreStr := fmt.Sprintf("%d / %d", sk.Correct, sk.Total)
		accStr := fmt.Sprintf("%.1f%%", sk.Accuracy)
		b.WriteString(fmt.Sprintf("%-38s %-12s %-10s %-12s\n", sk.ConceptName, scoreStr, accStr, sk.Status))
	}
	b.WriteString("----------------------------------------------------------------------------------\n\n")

	// 2. Error pattern analysis
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString("2. DIAGNOSTIC ERROR PATTERN ANALYSIS\n")
	b.WriteString("----------------------------------------------------------------------------------\n")
	if len(r.Errors) == 0 {
		b.WriteString("✓ Outstanding! Zero diagnostic distractor errors recorded during this exam.\n")
	} else {
		for i, errItem := range r.Errors {
			b.WriteString(fmt.Sprintf("[%d] %s (%d occurrences)\n", i+1, errItem.Misconception, errItem.Count))
			b.WriteString(fmt.Sprintf("    Remediation: %s\n\n", errItem.Remediation))
		}
	}
	b.WriteString("----------------------------------------------------------------------------------\n\n")

	// 3. Question-by-question review
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString("3. QUESTION-BY-QUESTION AUDIT REVIEW (Withheld During Exam)\n")
	b.WriteString("----------------------------------------------------------------------------------\n")
	for _, rev := range r.Reviews {
		statusMark := "✓ CORRECT"
		if !rev.IsCorrect {
			statusMark = "✗ INCORRECT"
		}
		b.WriteString(fmt.Sprintf("Question %d: [%s] (%s)\n", rev.QuestionIndex+1, statusMark, rev.Stage))
		b.WriteString(fmt.Sprintf("Scenario: %s\n", rev.PromptText))
		b.WriteString(fmt.Sprintf("Prompt:   %s\n", rev.StagePrompt))

		selText := "(no answer recorded)"
		if rev.SelectedOption != nil {
			selText = fmt.Sprintf("[%s] %s", rev.SelectedOption.Label, rev.SelectedOption.Text)
		}
		b.WriteString(fmt.Sprintf("Your Answer:    %s\n", selText))
		b.WriteString(fmt.Sprintf("Correct Answer: [%s] %s\n", rev.CorrectOption.Label, rev.CorrectOption.Text))

		if rev.ErrorTag != "" {
			misc, _ := DistractorExplanation(rev.ErrorTag)
			b.WriteString(fmt.Sprintf("Error Trap:     %s (%s)\n", rev.ErrorTag, misc))
		}
		b.WriteString(fmt.Sprintf("Explanation:    %s\n", rev.Explanation))

		if len(rev.Postings) > 0 {
			b.WriteString("Balanced Journal Entry:\n")
			for _, p := range rev.Postings {
				if p.Side == domain.SideDebit {
					b.WriteString(fmt.Sprintf("  Dr. %-24s %s\n", p.AccountID, p.Amount.FormatDollars()))
				} else {
					b.WriteString(fmt.Sprintf("      Cr. %-20s %s\n", p.AccountID, p.Amount.FormatDollars()))
				}
			}
		}
		b.WriteString("----------------------------------------------------------------------------------\n")
	}

	return b.String()
}

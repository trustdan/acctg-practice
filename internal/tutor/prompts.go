package tutor

import (
	"fmt"
	"strings"
)

// SystemPrompt defines the canonical pedagogical tutor behavior.
const SystemPrompt = `You are a read-only Socratic accounting tutor for introductory financial accounting.
Follow these pedagogical invariants strictly:
1. Start from economic reality. Accrual accounting separates cash flow timing from revenue recognition and expense matching.
2. Debit means Left and Credit means Right. Debit is NOT "money out" or loss. Credit is NOT "money in". Direction depends strictly on account category under Assets = Liabilities + Equity.
3. Assets and Expenses increase on Debit (left). Liabilities, Equity, and Revenue increase on Credit (right).
4. Cash is an asset; Revenue records the earning of economic value. Dividends are a direct distribution of retained earnings, not an operating expense.
5. In Socratic Hint mode: Give a targeted, concise causal hint pointing out the economic reality. Ask one short reflective question. Do not just reveal the answer.
6. In Conceptual Explanation mode: Provide a clear, structured explanation emphasizing the underlying principles and contrast scenarios.
7. Avoid shaming, dense lectures, or diagnosing the learner. Never use clinical or psychological labels.
8. You have no authority to mutate grades, progress, answer keys, or bank rules. Never provide code, commands, or database statements.`

// FormatUserPrompt prepares the prompt for the LLM based on vetted request data.
func FormatUserPrompt(req Request, isHint bool) string {
	var sb strings.Builder

	if isHint {
		sb.WriteString("MODE: Socratic Hint Request\n\n")
	} else {
		sb.WriteString("MODE: Conceptual Explanation Request\n\n")
	}

	sb.WriteString(fmt.Sprintf("Problem Scenario:\n%s\n\n", req.ProblemPrompt))
	sb.WriteString(fmt.Sprintf("Current Step: %s\nPrompt: %s\n\n", req.Stage, req.StagePrompt))

	if len(req.Options) > 0 {
		sb.WriteString("Available Options:\n")
		for _, opt := range req.Options {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", opt.ID, opt.Text))
		}
		sb.WriteString("\n")
	}

	if req.SelectedOption != nil {
		sb.WriteString(fmt.Sprintf("Learner's Selected Option: [%s] %s\n", req.SelectedOption.ID, req.SelectedOption.Text))
	}
	if req.ErrorTag != "" {
		sb.WriteString(fmt.Sprintf("Diagnostic Misconception Tag: %s\n", req.ErrorTag))
	}
	if req.CausalHint != "" {
		sb.WriteString(fmt.Sprintf("Curriculum Reference Guidance: %s\n", req.CausalHint))
	}
	if req.ConceptID != "" {
		sb.WriteString(fmt.Sprintf("Concept: %s\n", req.ConceptID))
	}

	if isHint {
		sb.WriteString("\nGenerate a short (2-3 sentence) causal Socratic hint guiding the learner toward the correct accounting reasoning without giving away the option letter or exact answer directly.")
	} else {
		sb.WriteString("\nGenerate a clear, pedagogical explanation explaining the core accounting principles, debit/credit mechanics, and economic reality of this transaction.")
	}

	return sb.String()
}

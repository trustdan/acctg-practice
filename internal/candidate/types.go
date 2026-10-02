package candidate

import (
	"fmt"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// CandidateStatus defines the lifecycle status of an unapproved question proposal.
type CandidateStatus string

const (
	StatusDraft                  CandidateStatus = "draft"
	StatusCandidatePendingReview CandidateStatus = "candidate_pending_review"
	StatusCandidateRejected      CandidateStatus = "rejected"
)

// ValidationStatus records the outcome of strict syntax and local semantic derivation.
type ValidationStatus string

const (
	ValidationPending ValidationStatus = "pending"
	ValidationValid   ValidationStatus = "valid"
	ValidationFailed  ValidationStatus = "failed"
)

// Provenance tracks the origin, model, prompt, and target context of a generated candidate.
type Provenance struct {
	Source        string    `json:"source"`
	Model         string    `json:"model,omitempty"`
	GeneratedAt   time.Time `json:"generated_at"`
	TargetFamily  string    `json:"target_family"`
	TargetConcept string    `json:"target_concept,omitempty"`
	PromptText    string    `json:"prompt_text,omitempty"`
}

// CandidateQuestion represents a creative question candidate outside the active question bank.
//
// Invariants (AGENT-CONTRACT):
//  1. Untrusted candidate data: Candidates are stored outside the active bank and cannot be selected for practice.
//  2. Local semantic derivation: Journal entry postings, equation deltas, and explanations are derived
//     strictly by the deterministic engine, NOT accepted blindly from model outputs.
//  3. Immutability: Generated content cannot alter family rules or learner mastery evidence.
type CandidateQuestion struct {
	Teaching         map[domain.DrillStage]bank.TeachingText `json:"teaching,omitempty"`
	ID               string                                  `json:"id"`
	FamilyID         string                                  `json:"family_id"`
	RuleVersion      int                                     `json:"rule_version"`
	Status           CandidateStatus                         `json:"status"`
	ScenarioTemplate string                                  `json:"scenario_template"`
	Parameters       map[string][]int64                      `json:"parameters"`
	Concepts         []string                                `json:"concepts"`
	DerivedFixture   bank.FixtureJSON                        `json:"derived_fixture"`
	DerivedEquation  domain.EquationDelta                    `json:"derived_equation"`
	Explanation      engine.EventExplanation                 `json:"explanation"`
	Provenance       Provenance                              `json:"provenance"`
	ValidationStatus ValidationStatus                        `json:"validation_status"`
	RejectionReason  string                                  `json:"rejection_reason,omitempty"`
}

// RawProposal captures the untrusted structured proposal parsed from JSON.
type RawProposal struct {
	Teaching         map[domain.DrillStage]bank.TeachingText `json:"teaching,omitempty"`
	ID               string                                  `json:"id,omitempty"`
	FamilyID         string                                  `json:"family_id"`
	ScenarioTemplate string                                  `json:"scenario_template"`
	Parameters       map[string][]int64                      `json:"parameters"`
	Concepts         []string                                `json:"concepts"`
	ProposedPostings []bank.FixturePostingJSON               `json:"proposed_postings,omitempty"`
}

// FormatPreview generates a clear human-readable preview of the candidate question,
// including wording, assumptions, locally derived postings, concept tags, and provenance.
func (c *CandidateQuestion) FormatPreview() string {
	var sb strings.Builder

	sb.WriteString("================================================================================\n")
	sb.WriteString(fmt.Sprintf("CANDIDATE QUESTION PREVIEW — [%s]\n", c.ID))
	sb.WriteString(fmt.Sprintf("Status: %s | Validation: %s\n", c.Status, c.ValidationStatus))
	if c.RejectionReason != "" {
		sb.WriteString(fmt.Sprintf("Rejection Reason: %s\n", c.RejectionReason))
	}
	sb.WriteString("--------------------------------------------------------------------------------\n")

	// 1. Wording (Template and Instantiated Sample)
	sb.WriteString("WORDING:\n")
	sb.WriteString(fmt.Sprintf("Template: %s\n", c.ScenarioTemplate))

	sampleAmt := int64(10000)
	if amounts, ok := c.Parameters["amount_minor_units"]; ok && len(amounts) > 0 {
		sampleAmt = amounts[0]
	}
	sampleFormatted := domain.NewMoney(sampleAmt).FormatDollars()
	instantiated := strings.ReplaceAll(c.ScenarioTemplate, "${amount_dollars}", sampleFormatted)
	sb.WriteString(fmt.Sprintf("Sample:   %s\n\n", instantiated))

	// 2. Assumptions (Parameters)
	sb.WriteString("ASSUMPTIONS (Parameters):\n")
	for param, vals := range c.Parameters {
		valStrs := make([]string, len(vals))
		for i, v := range vals {
			if strings.Contains(param, "minor_units") {
				valStrs[i] = fmt.Sprintf("%d (%s)", v, domain.NewMoney(v).FormatDollars())
			} else {
				valStrs[i] = fmt.Sprintf("%d", v)
			}
		}
		sb.WriteString(fmt.Sprintf("  - %s: [%s]\n", param, strings.Join(valStrs, ", ")))
	}
	sb.WriteString("\n")

	// 3. Derived Postings (Locally derived from accounting engine)
	sb.WriteString("DERIVED POSTINGS (Engine-Derived Canonical Entry):\n")
	if len(c.DerivedFixture.Postings) > 0 {
		for _, p := range c.DerivedFixture.Postings {
			sidePrefix := "Dr."
			if p.Side == string(domain.SideCredit) {
				sidePrefix = "Cr."
			}
			sb.WriteString(fmt.Sprintf("  %-8s %-28s %s\n", sidePrefix, p.AccountID, sampleFormatted))
		}
	} else {
		sb.WriteString("  [No valid derived postings]\n")
	}
	isEqBalanced := (c.DerivedEquation.DeltaAssets == c.DerivedEquation.DeltaLiabilities.Add(c.DerivedEquation.DeltaEquity))
	sb.WriteString(fmt.Sprintf("Equation Impact: Assets=%s, Liabilities=%s, Equity=%s (Balanced: %t)\n\n",
		c.DerivedEquation.DeltaAssets.FormatDollars(),
		c.DerivedEquation.DeltaLiabilities.FormatDollars(),
		c.DerivedEquation.DeltaEquity.FormatDollars(),
		isEqBalanced))

	// 4. Concept Tags
	sb.WriteString(fmt.Sprintf("CONCEPT TAGS: [%s]\n\n", strings.Join(c.Concepts, ", ")))

	// 5. Provenance
	sb.WriteString("PROVENANCE:\n")
	for _, stage := range []domain.DrillStage{domain.StageIdentifyAccount, domain.StageAccountCategory, domain.StageDirection, domain.StageDebitCredit, domain.StageCounterAccount, domain.StageBalancedEntry, domain.StageEquationEffect} {
		if text, ok := c.Teaching[stage]; ok {
			sb.WriteString(fmt.Sprintf("Proposed teaching [%s] (human review required):\n  Hint: %s\n  Explanation: %s\n", stage, text.Hint, text.Explanation))
		}
	}
	sb.WriteString(fmt.Sprintf("  Source:        %s\n", c.Provenance.Source))
	if c.Provenance.Model != "" {
		sb.WriteString(fmt.Sprintf("  Model:         %s\n", c.Provenance.Model))
	}
	sb.WriteString(fmt.Sprintf("  Generated At:  %s\n", c.Provenance.GeneratedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("  Target Family: %s (Rule Version %d)\n", c.FamilyID, c.RuleVersion))
	if c.Provenance.TargetConcept != "" {
		sb.WriteString(fmt.Sprintf("  Target Concept: %s\n", c.Provenance.TargetConcept))
	}
	if c.Explanation.Summary != "" {
		sb.WriteString(fmt.Sprintf("  Rationale:     %s\n", c.Explanation.Summary))
	}
	sb.WriteString("================================================================================\n")

	return sb.String()
}

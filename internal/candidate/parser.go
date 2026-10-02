package candidate

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// ParseProposal strictly parses an untrusted candidate JSON payload, validates syntax and parameters,
// and derives canonical journal postings, equation effects, and explanations LOCALLY using the deterministic engine.
//
// Invariant (AGENT-CONTRACT):
// Candidate generation never grants authority to grade or approve. Even if an external model proposes
// an entry, answers are derived locally. Malformed or contradictory proposals are rejected and remain inactive.
func ParseProposal(rawJSON string, eng *engine.Engine, catalog *domain.AccountCatalog, prov Provenance) (*CandidateQuestion, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine is required for local answer derivation")
	}
	if catalog == nil {
		return nil, fmt.Errorf("account catalog is required for validation")
	}

	cleanedJSON := extractJSONContent(rawJSON)

	var proposal RawProposal
	dec := json.NewDecoder(strings.NewReader(cleanedJSON))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&proposal); err != nil {
		// Return candidate with failed validation status and reason
		cand := &CandidateQuestion{
			ID:               generateCandidateID(proposal.FamilyID, prov.GeneratedAt),
			FamilyID:         proposal.FamilyID,
			Status:           StatusCandidateRejected,
			ValidationStatus: ValidationFailed,
			RejectionReason:  fmt.Sprintf("malformed JSON proposal: %v", err),
			Provenance:       prov,
		}
		return cand, fmt.Errorf("failed to parse candidate JSON proposal: %w", err)
	}

	candID := proposal.ID
	if strings.TrimSpace(candID) == "" {
		candID = generateCandidateID(proposal.FamilyID, prov.GeneratedAt)
	}

	cand := &CandidateQuestion{
		ID:               candID,
		FamilyID:         proposal.FamilyID,
		RuleVersion:      1,
		Status:           StatusCandidatePendingReview,
		ScenarioTemplate: strings.TrimSpace(proposal.ScenarioTemplate),
		Parameters:       proposal.Parameters,
		Concepts:         proposal.Concepts,
		Provenance:       prov,
		ValidationStatus: ValidationPending,
	}

	// 1. Validate supported family
	if !bank.IsSupportedFamily(cand.FamilyID) {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = fmt.Sprintf("unsupported event family: %q", cand.FamilyID)
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}

	// 2. Validate scenario template
	if cand.ScenarioTemplate == "" {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = "scenario_template cannot be empty"
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}
	if !strings.Contains(cand.ScenarioTemplate, "${amount_dollars}") {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = "scenario_template must include ${amount_dollars} parameter placeholder"
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}

	// 3. Validate parameters
	if len(cand.Parameters) == 0 {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = "parameters map cannot be empty"
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}
	reqParameters := bank.RequiredParameters(cand.FamilyID)
	for _, reqParam := range reqParameters {
		vals, ok := cand.Parameters[reqParam]
		if !ok || len(vals) == 0 {
			cand.Status = StatusCandidateRejected
			cand.ValidationStatus = ValidationFailed
			cand.RejectionReason = fmt.Sprintf("missing required parameter %q for family %s", reqParam, cand.FamilyID)
			return cand, fmt.Errorf("%s", cand.RejectionReason)
		}
		for i, v := range vals {
			if v <= 0 {
				cand.Status = StatusCandidateRejected
				cand.ValidationStatus = ValidationFailed
				cand.RejectionReason = fmt.Sprintf("parameter %q[%d] must be positive integer minor units, got %d", reqParam, i, v)
				return cand, fmt.Errorf("%s", cand.RejectionReason)
			}
		}
	}

	// 4. Validate concepts
	if len(cand.Concepts) == 0 {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = "concepts list cannot be empty"
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}
	for i, c := range cand.Concepts {
		if strings.TrimSpace(c) == "" {
			cand.Status = StatusCandidateRejected
			cand.ValidationStatus = ValidationFailed
			cand.RejectionReason = fmt.Sprintf("concept[%d] cannot be empty", i)
			return cand, fmt.Errorf("%s", cand.RejectionReason)
		}
	}

	// 5. Derive answers LOCALLY using the deterministic engine
	sampleAmounts := cand.Parameters["amount_minor_units"]
	sampleAmt := sampleAmounts[0]
	event := engine.TransactionEvent{
		FamilyID: cand.FamilyID,
		Parameters: map[string]int64{
			"amount_minor_units": sampleAmt,
		},
	}

	processed, err := eng.ProcessEvent(event)
	if err != nil {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = fmt.Sprintf("local engine failed to process event for family %s: %v", cand.FamilyID, err)
		return cand, fmt.Errorf("%s", cand.RejectionReason)
	}

	// Construct FixtureJSON from the engine-derived canonical entry
	var fixturePostings []bank.FixturePostingJSON
	for _, p := range processed.Entry.Postings {
		fixturePostings = append(fixturePostings, bank.FixturePostingJSON{
			AccountID:       string(p.AccountID),
			Side:            string(p.Side),
			AmountParameter: "amount_minor_units",
		})
	}
	cand.DerivedFixture = bank.FixtureJSON{Postings: fixturePostings}
	cand.DerivedEquation = processed.Equation
	cand.Explanation = processed.Explanation

	// 6. Contradictory Semantics Check
	// If the untrusted proposal provided proposed postings, evaluate them against canonical semantics.
	if len(proposal.ProposedPostings) > 0 {
		var candidatePostings []domain.Posting
		for _, pp := range proposal.ProposedPostings {
			side := domain.Side(pp.Side)
			candidatePostings = append(candidatePostings, domain.Posting{
				AccountID: domain.AccountID(pp.AccountID),
				Side:      side,
				Amount:    domain.Money(sampleAmt),
			})
		}
		candEntry := domain.Entry{Postings: candidatePostings}
		eval := eng.EvaluateEntry(event, candEntry)
		if !eval.IsCorrect {
			cand.Status = StatusCandidateRejected
			cand.ValidationStatus = ValidationFailed
			cand.RejectionReason = fmt.Sprintf("contradictory event semantics: proposal proposed postings that contradict canonical engine rules: %s", eval.Feedback)
			return cand, fmt.Errorf("%s", cand.RejectionReason)
		}
	}

	// Candidate is valid and pending human semantic review
	cand.ValidationStatus = ValidationValid
	cand.Status = StatusCandidatePendingReview

	return cand, nil
}

// extractJSONContent strips markdown code fences (```json ... ```) if an LLM wrapped its JSON response.
func extractJSONContent(input string) string {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "```") {
		// Find first newline
		nlIdx := strings.Index(trimmed, "\n")
		if nlIdx != -1 {
			trimmed = trimmed[nlIdx+1:]
		}
		// Strip trailing fence
		if endIdx := strings.LastIndex(trimmed, "```"); endIdx != -1 {
			trimmed = trimmed[:endIdx]
		}
	}
	return strings.TrimSpace(trimmed)
}

var candidateSeq uint64

func generateCandidateID(familyID string, t time.Time) string {
	if t.IsZero() {
		t = time.Now().UTC()
	}
	cleanFamily := strings.TrimSpace(familyID)
	if cleanFamily == "" {
		cleanFamily = "cand"
	}
	seq := atomic.AddUint64(&candidateSeq, 1)
	return fmt.Sprintf("cand_%s_%s_%06d_%d", cleanFamily, t.Format("20060102_150405"), (t.Nanosecond()/1000)%1000000, seq)
}

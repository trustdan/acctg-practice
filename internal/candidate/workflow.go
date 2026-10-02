package candidate

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

var approvalSeq uint64

func generateApprovalEventID(action bank.ApprovalAction, targetID string, t time.Time) string {
	if t.IsZero() {
		t = time.Now().UTC()
	}
	seq := atomic.AddUint64(&approvalSeq, 1)
	return fmt.Sprintf("appr_%s_%s_%s_%06d_%d", action, targetID, t.Format("20060102_150405"), (t.Nanosecond()/1000)%1000000, seq)
}

// ApproveCandidate evaluates an unapproved candidate against family review status and
// semantic wording rules. If valid, it returns an immutable bank.QuestionJSON ready for
// active bank publishing and an audit bank.ApprovalEvent.
//
// Invariants (AGENT-CONTRACT):
//  1. Balanced entries are necessary but insufficient: semantic wording must support the answer.
//     A balanced-but-wrong candidate (e.g. Dr Cash / Cr Revenue wording for customer_advance) is rejected.
//  2. New event families require reviewed engine rules and tests before approval.
//  3. Reviewer identity, timestamp, source, rule version, and approval action are recorded.
func ApproveCandidate(cand *CandidateQuestion, reviewer string, notes string, clock time.Time) (*bank.QuestionJSON, *bank.ApprovalEvent, error) {
	if cand == nil {
		return nil, nil, fmt.Errorf("candidate cannot be nil")
	}

	reviewer = strings.TrimSpace(reviewer)
	if reviewer == "" {
		return nil, nil, fmt.Errorf("reviewer identity is required for approval")
	}

	if cand.ValidationStatus != ValidationValid {
		return nil, nil, fmt.Errorf("cannot approve candidate %s: validation status is %q (must be %q). Rejection reason: %s",
			cand.ID, cand.ValidationStatus, ValidationValid, cand.RejectionReason)
	}

	if cand.Status == StatusCandidateRejected {
		return nil, nil, fmt.Errorf("cannot approve candidate %s: candidate is marked rejected; repair before approving", cand.ID)
	}

	// 1. Verify that the event family has reviewed engine rules and regression tests
	if !bank.IsFamilyReviewedAndTested(cand.FamilyID, cand.RuleVersion) {
		return nil, nil, fmt.Errorf("cannot approve candidate %s: family %q (rule_version %d) is not a reviewed and tested event family; new families require reviewed rule additions in the engine and passing regression tests before approval",
			cand.ID, cand.FamilyID, cand.RuleVersion)
	}

	// 2. Perform semantic wording review: catch balanced-but-wrong wording
	semResult := bank.CheckSemanticWording(cand.FamilyID, cand.ScenarioTemplate)
	if !semResult.Approved {
		return nil, nil, fmt.Errorf("cannot approve candidate %s: semantic wording review failed:\n- %s",
			cand.ID, strings.Join(semResult.Violations, "\n- "))
	}

	if clock.IsZero() {
		clock = time.Now().UTC()
	}
	approvedAtStr := clock.UTC().Format(time.RFC3339)

	// 3. Construct immutable published bank question
	pub := &bank.QuestionJSON{
		Teaching:         cand.Teaching,
		ID:               cand.ID,
		Version:          1,
		FamilyID:         cand.FamilyID,
		RuleVersion:      cand.RuleVersion,
		Status:           bank.StatusApprovedActive,
		ScenarioTemplate: cand.ScenarioTemplate,
		Parameters:       cand.Parameters,
		Concepts:         cand.Concepts,
		ExpectedFixture:  cand.DerivedFixture,
		Review: bank.ReviewJSON{
			Reviewer:   &reviewer,
			ApprovedAt: &approvedAtStr,
			Source:     cand.Provenance.Source,
		},
	}

	// 4. Construct approval event audit record
	apprID := generateApprovalEventID(bank.ActionApprove, cand.ID, clock)
	event := &bank.ApprovalEvent{
		ID:          apprID,
		CandidateID: cand.ID,
		QuestionID:  cand.ID,
		Action:      bank.ActionApprove,
		Reviewer:    reviewer,
		RuleVersion: cand.RuleVersion,
		Source:      cand.Provenance.Source,
		Notes:       notes,
		CreatedAt:   clock.UTC(),
	}

	cand.Status = bank.StatusApprovedActive

	return pub, event, nil
}

// RejectCandidate marks a candidate proposal as rejected and records an approval audit event.
func RejectCandidate(cand *CandidateQuestion, reviewer string, reason string, clock time.Time) (*bank.ApprovalEvent, error) {
	if cand == nil {
		return nil, fmt.Errorf("candidate cannot be nil")
	}

	reviewer = strings.TrimSpace(reviewer)
	if reviewer == "" {
		return nil, fmt.Errorf("reviewer identity is required for rejection")
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, fmt.Errorf("rejection reason is required")
	}

	if clock.IsZero() {
		clock = time.Now().UTC()
	}

	cand.Status = StatusCandidateRejected
	cand.RejectionReason = reason

	apprID := generateApprovalEventID(bank.ActionReject, cand.ID, clock)
	event := &bank.ApprovalEvent{
		ID:          apprID,
		CandidateID: cand.ID,
		QuestionID:  cand.ID,
		Action:      bank.ActionReject,
		Reviewer:    reviewer,
		RuleVersion: cand.RuleVersion,
		Source:      cand.Provenance.Source,
		Notes:       reason,
		CreatedAt:   clock.UTC(),
	}

	return event, nil
}

// RepairCandidate repairs a candidate proposal by updating its scenario template and/or parameters,
// re-deriving accounting answers locally using the engine, and recording a repair audit event.
func RepairCandidate(
	cand *CandidateQuestion,
	newTemplate string,
	newParams map[string][]int64,
	eng *engine.Engine,
	catalog *domain.AccountCatalog,
	reviewer string,
	notes string,
	clock time.Time,
) (*bank.ApprovalEvent, error) {
	if cand == nil {
		return nil, fmt.Errorf("candidate cannot be nil")
	}
	if eng == nil {
		return nil, fmt.Errorf("accounting engine is required for answer derivation")
	}
	if catalog == nil {
		return nil, fmt.Errorf("account catalog is required for validation")
	}

	reviewer = strings.TrimSpace(reviewer)
	if reviewer == "" {
		return nil, fmt.Errorf("reviewer identity is required for repair")
	}

	newTemplate = strings.TrimSpace(newTemplate)
	if newTemplate != "" {
		if !strings.Contains(newTemplate, "${amount_dollars}") {
			return nil, fmt.Errorf("repaired template must include '${amount_dollars}' placeholder")
		}
		cand.ScenarioTemplate = newTemplate
	}

	if newParams != nil && len(newParams) > 0 {
		reqParameters := bank.RequiredParameters(cand.FamilyID)
		for _, reqParam := range reqParameters {
			vals, ok := newParams[reqParam]
			if !ok || len(vals) == 0 {
				return nil, fmt.Errorf("missing required parameter %q for family %s", reqParam, cand.FamilyID)
			}
			for i, v := range vals {
				if v <= 0 {
					return nil, fmt.Errorf("parameter %q[%d] must be positive integer minor units, got %d", reqParam, i, v)
				}
			}
		}
		cand.Parameters = newParams
	}

	// Re-derive answers locally using engine
	sampleAmounts := cand.Parameters["amount_minor_units"]
	if len(sampleAmounts) == 0 {
		return nil, fmt.Errorf("candidate parameters missing amount_minor_units")
	}
	sampleAmt := sampleAmounts[0]
	event := engine.TransactionEvent{
		FamilyID: cand.FamilyID,
		Parameters: map[string]int64{
			"amount_minor_units": sampleAmt,
		},
	}

	processed, err := eng.ProcessEvent(event)
	if err != nil {
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = fmt.Sprintf("engine derivation failed on repaired candidate: %v", err)
		return nil, fmt.Errorf("%s", cand.RejectionReason)
	}

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

	// Check semantic wording
	semResult := bank.CheckSemanticWording(cand.FamilyID, cand.ScenarioTemplate)
	if !semResult.Approved {
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = fmt.Sprintf("repaired wording has semantic conflicts: %s", strings.Join(semResult.Violations, "; "))
	} else {
		cand.ValidationStatus = ValidationValid
		cand.RejectionReason = ""
		cand.Status = StatusCandidatePendingReview
	}

	if clock.IsZero() {
		clock = time.Now().UTC()
	}

	apprID := generateApprovalEventID(bank.ActionRepair, cand.ID, clock)
	apprEvent := &bank.ApprovalEvent{
		ID:          apprID,
		CandidateID: cand.ID,
		QuestionID:  cand.ID,
		Action:      bank.ActionRepair,
		Reviewer:    reviewer,
		RuleVersion: cand.RuleVersion,
		Source:      cand.Provenance.Source,
		Notes:       notes,
		CreatedAt:   clock.UTC(),
	}

	return apprEvent, nil
}

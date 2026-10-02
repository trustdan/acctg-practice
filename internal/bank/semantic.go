package bank

import (
	"fmt"
	"strings"
)

// SemanticReviewResult captures the outcome of evaluating scenario template wording
// against accounting event semantics.
type SemanticReviewResult struct {
	Approved   bool     `json:"approved"`
	Violations []string `json:"violations,omitempty"`
}

// CheckSemanticWording inspects the scenario wording of a proposed question against the
// economic reality and accounting rules of its event family.
//
// Invariants (AGENT-CONTRACT):
//  1. Syntactic checks and balanced entries are necessary but insufficient.
//     (e.g., Dr Cash / Cr Revenue balances, but is wrong for an unearned advance).
//  2. Variants require explicit semantic wording review to prevent contradictory or misleading text.
//  3. Timing must be unambiguous (e.g. explicitly distinguishing today vs. next month vs. on account).
func CheckSemanticWording(familyID string, scenarioTemplate string) SemanticReviewResult {
	var violations []string

	trimmed := strings.TrimSpace(scenarioTemplate)
	if trimmed == "" {
		return SemanticReviewResult{
			Approved:   false,
			Violations: []string{"scenario template cannot be empty"},
		}
	}

	if !strings.Contains(trimmed, "${amount_dollars}") {
		violations = append(violations, "scenario template must contain '${amount_dollars}' placeholder")
	}

	lower := strings.ToLower(trimmed)

	switch familyID {
	case FamilyCustomerAdvance:
		// Customer advance: Cash received in advance for FUTURE work. Liability (unearned revenue) is recorded.
		// Wording claiming services are completed today or revenue is recognized immediately contradicts this family.
		contradictoryTerms := []string{
			"completed today",
			"delivered today",
			"performed today",
			"fully delivered",
			"services completed",
			"services provided today",
			"no future obligation",
			"recognized revenue immediately",
			"recorded as revenue",
			"earned revenue today",
			"recorded service revenue",
			"earned upon receipt",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: customer_advance wording claims services delivered today or revenue earned (%q), which contradicts unearned revenue liability",
					term,
				))
			}
		}

		// Must have at least one advance / future timing indicator
		futureIndicators := []string{
			"advance", "future", "next month", "deposit", "upcoming",
			"retainer", "will provide", "will perform", "later",
			"to be delivered", "to be performed", "in advance", "unearned",
			"following month", "next quarter", "next week",
		}
		hasFutureIndicator := false
		for _, ind := range futureIndicators {
			if strings.Contains(lower, ind) {
				hasFutureIndicator = true
				break
			}
		}
		if !hasFutureIndicator {
			violations = append(violations, "semantic conflict: customer_advance wording lacks explicit advance or future delivery timing indicator (e.g., 'in advance', 'next month', 'deposit')")
		}

	case FamilyCashService:
		// Cash service: services rendered and completed TODAY for cash.
		contradictoryTerms := []string{
			"work next month",
			"advance deposit",
			"unearned",
			"future delivery",
			"services to be performed next",
			"on account",
			"billed the client",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: cash_service wording claims future delivery or credit billing (%q), contradicting immediate cash revenue",
					term,
				))
			}
		}

	case FamilyServiceOnCredit:
		// Service on credit: services rendered today, billed on account.
		contradictoryTerms := []string{
			"received cash today",
			"paid cash today",
			"cash in hand",
			"cash payment upon delivery",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: service_on_credit wording claims cash was received today (%q), contradicting on-account receivable",
					term,
				))
			}
		}

	case FamilyCollectReceivable:
		// Collect receivable: collecting cash on existing receivable; no new revenue earned today.
		contradictoryTerms := []string{
			"earned service revenue today",
			"new service performed today",
			"recognized revenue",
			"earned revenue",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: collect_receivable wording claims new revenue recognized today (%q), which produces duplicate revenue",
					term,
				))
			}
		}

	case FamilyPrepaidPurchase:
		// Prepaid purchase: cash paid for future benefit (asset). NOT an immediate expense.
		contradictoryTerms := []string{
			"incurred last month",
			"recorded as rent expense",
			"recorded as insurance expense",
			"used up immediately",
			"operating expense for the month",
			"recognized as expense immediately",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: prepaid_purchase wording claims immediate expense or past consumption (%q), contradicting asset capitalization",
					term,
				))
			}
		}

	case FamilyEquipmentPurchaseCash:
		// Equipment purchase: cash paid for long-term productive equipment (asset). NOT an operating expense.
		contradictoryTerms := []string{
			"operating expense",
			"recorded as expense",
			"supplies consumed today",
			"repair expense",
			"expense for the month",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: equipment_purchase_cash wording claims operating expense (%q), contradicting long-term equipment asset capitalization",
					term,
				))
			}
		}

	case FamilyBorrowCash:
		// Borrow cash: borrowing on note payable (liability). NOT revenue or stock.
		contradictoryTerms := []string{
			"earned revenue",
			"issued common stock",
			"customer advance",
			"service revenue",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: borrow_cash wording claims revenue or equity issuance (%q), contradicting promissory note liability",
					term,
				))
			}
		}

	case FamilyDividendCash:
		// Dividend cash: distribution to owners reducing equity. NOT an operating expense.
		contradictoryTerms := []string{
			"salary expense",
			"wage expense",
			"operating expense",
			"utility expense",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: dividend_cash wording claims operating expense (%q), contradicting equity dividend distribution",
					term,
				))
			}
		}

	case FamilyRepayNotePrincipal:
		// Repaying note principal: settles liability. NOT an expense.
		contradictoryTerms := []string{
			"operating expense",
			"interest expense only",
			"service expense",
		}
		for _, term := range contradictoryTerms {
			if isContradictionAsserted(lower, term) {
				violations = append(violations, fmt.Sprintf(
					"semantic conflict: repay_note_principal wording claims operating expense (%q), contradicting liability principal reduction",
					term,
				))
			}
		}
	}

	return SemanticReviewResult{
		Approved:   len(violations) == 0,
		Violations: violations,
	}
}

// isContradictionAsserted checks whether a contradictory term appears in the text
// without being negated (e.g. "not earned revenue", "not an operating expense", "rather than").
func isContradictionAsserted(text string, term string) bool {
	cursor := 0
	for {
		idx := strings.Index(text[cursor:], term)
		if idx == -1 {
			return false
		}
		actualIdx := cursor + idx
		prefix := text[:actualIdx]
		// Check for negation prefixes immediately preceding the term
		if strings.HasSuffix(prefix, "not ") ||
			strings.HasSuffix(prefix, "not an ") ||
			strings.HasSuffix(prefix, "not a ") ||
			strings.HasSuffix(prefix, "never ") ||
			strings.HasSuffix(prefix, "instead of ") ||
			strings.HasSuffix(prefix, "rather than ") {
			cursor = actualIdx + len(term)
			continue
		}
		return true
	}
}

package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

// Engine is the central deterministic accounting processor.
// It derives balanced entries, computes balance sheet equation deltas,
// and evaluates learner answers against canonical accounting rules.
type Engine struct {
	catalog *domain.AccountCatalog
	rules   map[string]EventRule
}

func NewEngine(catalog *domain.AccountCatalog) *Engine {
	e := &Engine{
		catalog: catalog,
		rules:   make(map[string]EventRule),
	}

	// Register canonical rules
	cashService := CashServiceRule{}
	serviceCredit := ServiceOnCreditRule{}
	custAdvance := CustomerAdvanceRule{}
	collectRec := CollectReceivableRule{}
	earnAdvance := EarnAdvanceRule{}
	cashRent := CashExpenseRule{ExpenseAccount: "rent_expense"}
	prepaidPur := PrepaidPurchaseRule{}
	prepaidCons := PrepaidConsumptionRule{}
	equipPur := EquipmentPurchaseCashRule{}
	borrowCash := BorrowCashRule{}
	repayNote := RepayNotePrincipalRule{}
	issueShares := IssueSharesRule{}
	dividendCash := DividendCashRule{}
	declareDividend := DeclareDividendRule{}
	payDividend := PayDividendPayableRule{}

	e.rules[bank.FamilyCashService] = cashService
	e.rules[bank.FamilyServiceOnCredit] = serviceCredit
	e.rules[bank.FamilyCustomerAdvance] = custAdvance
	e.rules[bank.FamilyCollectReceivable] = collectRec
	e.rules["receivable_collection"] = collectRec
	e.rules[bank.FamilyEarnAdvance] = earnAdvance
	e.rules[bank.FamilyCashRent] = cashRent
	e.rules[bank.FamilyCashExpense] = cashRent
	e.rules[bank.FamilyPrepaidPurchase] = prepaidPur
	e.rules[bank.FamilyPrepaidConsumption] = prepaidCons
	e.rules[bank.FamilyEquipmentPurchaseCash] = equipPur
	e.rules[bank.FamilyBorrowCash] = borrowCash
	e.rules["borrowing_note"] = borrowCash
	e.rules[bank.FamilyRepayNotePrincipal] = repayNote
	e.rules[bank.FamilyIssueShares] = issueShares
	e.rules["issue_stock_cash"] = issueShares
	e.rules[bank.FamilyDividendCash] = dividendCash
	e.rules[bank.FamilyDeclareDividend] = declareDividend
	e.rules[bank.FamilyPayDividendPayable] = payDividend

	return e
}

// ProcessEvent applies the rule for the event, producing the canonical balanced entry and equation delta.
func (e *Engine) ProcessEvent(event TransactionEvent) (ProcessedEvent, error) {
	rule, ok := e.rules[event.FamilyID]
	if !ok {
		return ProcessedEvent{}, fmt.Errorf("unsupported event family: %q", event.FamilyID)
	}

	entry, err := rule.DeriveEntry(event.Parameters, e.catalog)
	if err != nil {
		return ProcessedEvent{}, fmt.Errorf("failed to derive entry for %s: %w", event.FamilyID, err)
	}

	if err := entry.Validate(e.catalog); err != nil {
		return ProcessedEvent{}, fmt.Errorf("derived entry is invalid for %s: %w", event.FamilyID, err)
	}

	eqDelta, err := entry.EquationEffects(e.catalog)
	if err != nil {
		return ProcessedEvent{}, fmt.Errorf("derived entry equation error for %s: %w", event.FamilyID, err)
	}

	explanation := rule.Explain(event.Parameters)

	return ProcessedEvent{
		Event:       event,
		Entry:       entry,
		Equation:    eqDelta,
		Explanation: explanation,
	}, nil
}

// EvaluateEntry compares a candidate journal entry against the canonical entry for an event.
func (e *Engine) EvaluateEntry(event TransactionEvent, candidate domain.Entry) SemanticEvaluation {
	rule, ok := e.rules[event.FamilyID]
	if !ok {
		return SemanticEvaluation{
			IsCorrect: false,
			Feedback:  fmt.Sprintf("Unknown event family %q", event.FamilyID),
		}
	}

	canonical, err := rule.DeriveEntry(event.Parameters, e.catalog)
	if err != nil {
		return SemanticEvaluation{
			IsCorrect: false,
			Feedback:  fmt.Sprintf("Could not derive canonical entry: %v", err),
		}
	}

	// Check if candidate matches canonical postings (ignoring line order, aggregating duplicate lines)
	if canonical.EqualNormalized(candidate) {
		return SemanticEvaluation{
			IsCorrect: true,
			ErrorTags: []string{TagCorrect},
			Feedback:  "Correct entry. All accounts, sides, and amounts match canonical event semantics.",
		}
	}

	// Diagnose failure
	var errorTags []string

	// Check if unbalanced
	if !candidate.IsBalanced() {
		errorTags = append(errorTags, TagUnbalancedEntry)
	}

	normCandidate := candidate.AggregatePostings()

	// Check if reversed (debits became credits, credits became debits)
	if isReversed(canonical, candidate) || isReversed(canonical, normCandidate) {
		errorTags = append(errorTags, TagReversedSides)
	}

	// Check rule-specific diagnostic patterns (e.g. duplicate revenue, premature recognition)
	diagTags, diagFeedback := rule.DiagnoseError(normCandidate, event.Parameters)
	if len(diagTags) == 0 {
		diagTags, diagFeedback = rule.DiagnoseError(candidate, event.Parameters)
	}
	if len(diagTags) > 0 {
		errorTags = append(errorTags, diagTags...)
	}

	if len(errorTags) == 0 {
		errorTags = append(errorTags, TagWrongAccount)
	}

	feedback := diagFeedback
	if feedback == "" {
		if contains(errorTags, TagReversedSides) {
			feedback = "Your debit and credit sides are reversed. Remember: Debits go on the left, Credits on the right."
		} else if contains(errorTags, TagUnbalancedEntry) {
			feedback = fmt.Sprintf("Your entry does not balance: total debits (%s) != total credits (%s). Debits must equal credits.",
				candidate.TotalDebit().FormatDollars(), candidate.TotalCredit().FormatDollars())
		} else {
			feedback = "The selected accounts or amounts do not match the economic event."
		}
	}

	return SemanticEvaluation{
		IsCorrect: false,
		ErrorTags: errorTags,
		Feedback:  feedback,
	}
}

// entriesEqual checks if two entries contain identical postings regardless of line order.
func entriesEqual(a, b domain.Entry) bool {
	return a.EqualPostings(b)
}

func isReversed(canonical, candidate domain.Entry) bool {
	if len(canonical.Postings) != len(candidate.Postings) {
		return false
	}
	matched := make([]bool, len(candidate.Postings))
	for _, pCan := range canonical.Postings {
		found := false
		for j, pCand := range candidate.Postings {
			if !matched[j] && pCan.AccountID == pCand.AccountID && pCan.Side == pCand.Side.Opposite() && pCan.Amount == pCand.Amount {
				matched[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

package tutor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// OfflineTutor provides deterministic, zero-network pedagogical tutoring.
// It respects context cancellation and deadlines.
type OfflineTutor struct{}

// NewOfflineTutor creates a new OfflineTutor.
func NewOfflineTutor() *OfflineTutor {
	return &OfflineTutor{}
}

// Name implements Tutor.
func (o *OfflineTutor) Name() string {
	return "OfflineTutor"
}

// Hint implements Tutor.
func (o *OfflineTutor) Hint(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	hintText := o.generateSocraticHint(req)
	tokens := (len(hintText) + 3) / 4

	return Response{
		Text:        hintText,
		Provider:    "offline",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// Explain implements Tutor.
func (o *OfflineTutor) Explain(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	explainText := o.generateExplanation(req)
	tokens := (len(explainText) + 3) / 4

	return Response{
		Text:        explainText,
		Provider:    "offline",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

func (o *OfflineTutor) generateSocraticHint(req Request) string {
	// First check if there is an error tag from a submitted or inspected distractor
	if req.ErrorTag != "" {
		if tagHint := socraticHintForErrorTag(req.ErrorTag); tagHint != "" {
			return tagHint
		}
	}

	// If a vetted causal hint is provided in the request, prioritize its causal guidance
	if req.CausalHint != "" {
		return fmt.Sprintf("Consider the economic reality: %s", req.CausalHint)
	}

	// Stage-based Socratic guidance
	switch req.Stage {
	case domain.StageIdentifyAccount:
		return "Consider the economic event: which tangible resource, claim, or obligation changed first? Did cash change hands, or was a service provided?"
	case domain.StageAccountCategory:
		return "Classify the account under the accounting equation (Assets = Liabilities + Equity). Is it an economic resource owned (Asset), an obligation owed to third parties (Liability), or an owner claim / operational result (Equity)?"
	case domain.StageDirection:
		return "Did this specific account's balance increase or decrease as a result of this transaction?"
	case domain.StageDebitCredit:
		return "Remember the fundamental rule: Debit means Left and Credit means Right. Assets and Expenses increase on Debit (left); Liabilities, Equity, and Revenue increase on Credit (right)."
	case domain.StageCounterAccount:
		return "Every transaction has two sides to stay in balance. What balancing resource, claim, or obligation offsets the first account?"
	case domain.StageBalancedEntry:
		return "Check total debits and credits: every balanced journal entry must have total debits equal to total credits."
	case domain.StageEquationEffect:
		return "Reconcile the net effect on the accounting equation: does ΔAssets equal ΔLiabilities + ΔEquity?"
	default:
		return "Look at the timing of the economic exchange: separate when cash moves from when the service or benefit is earned or consumed."
	}
}

func (o *OfflineTutor) generateExplanation(req Request) string {
	var sb strings.Builder

	// If vetted explanation exists, use it as the foundation
	if req.Explanation != "" {
		sb.WriteString(req.Explanation)
		sb.WriteString("\n\n")
	}

	// Pedagogical core principles
	sb.WriteString("Core Accounting Principles:\n")
	sb.WriteString("• Fundamental Equation: Assets = Liabilities + Equity (A = L + E).\n")
	sb.WriteString("• Debit & Credit Convention: Debit means Left and Credit means Right. Debit is not 'money out' or 'loss', and Credit is not 'money in'.\n")
	sb.WriteString("• Normal Balances: Assets and Expenses increase on Debit (left). Liabilities, Equity, and Revenue increase on Credit (right).\n")
	sb.WriteString("• Accrual Principle: Revenue records earning economic value, while Cash is a separate asset. Cash flow timing does not dictate revenue recognition or expense matching.")

	// Family-specific contrast guidance
	if contrast := familyContrastExplanation(req.FamilyID); contrast != "" {
		sb.WriteString("\n\nContrast Analysis:\n")
		sb.WriteString(contrast)
	}

	return sb.String()
}

func socraticHintForErrorTag(tag string) string {
	switch tag {
	case engine.TagDuplicateRevenueOnCollection:
		return "Notice that revenue was already recognized when the service was performed on credit. Today's event is collecting cash owed. Does collecting a receivable create new revenue, or does it swap Accounts Receivable for Cash?"
	case engine.TagRevenueRecognizedPrematurely:
		return "Consider the earning process: did the company perform the service today, or did the customer merely pay an advance deposit for future work? When cash is collected before performance, what obligation does the company owe?"
	case engine.TagCashRecordedOnEarningAdvance:
		return "Check cash timing: was cash received today, or was the cash deposit collected in an earlier transaction? Today we delivered the service, reducing our unearned obligation."
	case engine.TagExpenseRecordedOnPrepaidPurchase:
		return "Consider the benefit horizon: does this payment benefit only today, or does it cover upcoming future months? Paying in advance for future benefits acquires an asset (Prepaid Expense), not an immediate period expense."
	case engine.TagExpenseRecordedOnEquipmentPurchase:
		return "Equipment provides productive economic benefit over multiple years. Purchasing long-lived operational equipment is a capital asset acquisition (asset swap), not an immediate period expense."
	case engine.TagExpenseRecordedOnLoanRepayment:
		return "Repaying loan principal extinguishes a liability (Notes Payable) that was previously borrowed. Repaying principal is not an expense—only interest represents the cost of borrowing."
	case engine.TagExpenseRecordedOnDividend:
		return "Dividends represent a direct distribution of retained earnings to owners, reducing Retained Earnings. A dividend distribution is not an operating expense and does not reduce Net Income."
	case engine.TagReversedSides:
		return "Check side placement: Debit is always the Left side, and Credit is always the Right side. Verify whether this account increases or decreases, and apply its normal side rule."
	case engine.TagUnbalancedEntry:
		return "Total debits must strictly equal total credits. Check the amounts to ensure both sides balance."
	default:
		return ""
	}
}

func familyContrastExplanation(familyID string) string {
	switch familyID {
	case "collect_receivable":
		return "Contrast: Performing a service on account recognizes Revenue and Accounts Receivable. Later, collecting cash increases Cash (Debit) and decreases Accounts Receivable (Credit). It is an asset swap and never recognizes duplicate revenue."
	case "customer_advance":
		return "Contrast: Receiving cash before performing work creates a liability (Unearned Revenue). Revenue is earned only later when the service or goods are delivered to the customer."
	case "earn_advance":
		return "Contrast: When fulfilling an advance, no new cash is received today. Instead, the liability Unearned Revenue is debited (decreased), and Service Revenue is credited (increased)."
	case "prepaid_purchase":
		return "Contrast: Paying for future insurance or rent is an asset swap (Debit Prepaid Expenses, Credit Cash). The expense is recognized month by month as the benefit is consumed."
	case "prepaid_consumption":
		return "Contrast: As prepaid coverage expires, we record Rent/Insurance Expense (Debit) and decrease the asset Prepaid Expenses (Credit). No cash moves at consumption time."
	case "equipment_purchase_cash":
		return "Contrast: Long-lived equipment is capitalized as an Asset, not expensed at purchase. It is expensed gradually over its useful life through depreciation."
	case "repay_note_principal":
		return "Contrast: Paying down loan principal reduces the liability Notes Payable (Debit) with Cash (Credit). It is a liability reduction, not an expense."
	case "dividend_cash":
		return "Contrast: Dividends are a return of capital/profits to shareholders, directly debiting Retained Earnings (Equity). Dividends are not an expense and do not appear on the Income Statement."
	default:
		return ""
	}
}

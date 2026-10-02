package domain

import (
	"fmt"
	"strings"
)

// TransactionReconciliation provides a mathematical tie-out and audit proof
// that the Journal Entry, T-Accounts, and Accounting Equation all reconcile
// to the exact same underlying postings and amounts.
type TransactionReconciliation struct {
	Entry            Entry           `json:"entry"`
	Ledger           *TAccountLedger `json:"ledger"`
	EquationDelta    EquationDelta   `json:"equation_delta"`
	TotalDr          Money           `json:"total_dr"`
	TotalCr          Money           `json:"total_cr"`
	IsBalanced       bool            `json:"is_balanced"`
	TAccountDrSum    Money           `json:"t_account_dr_sum"`
	TAccountCrSum    Money           `json:"t_account_cr_sum"`
	DeltaAssets      Money           `json:"delta_assets"`
	DeltaLiabilities Money           `json:"delta_liabilities"`
	DeltaEquity      Money           `json:"delta_equity"`
	Reconciled       bool            `json:"reconciled"`
	AuditNotes       []string        `json:"audit_notes"`
	VisualNotice     string          `json:"visual_notice"`
}

// ReconcileTransaction validates an entry, posts it to a T-account ledger, calculates its
// equation effects, and verifies that all three views reconcile to the exact same postings.
func ReconcileTransaction(entry Entry, catalog *AccountCatalog) (*TransactionReconciliation, error) {
	if catalog == nil {
		return nil, fmt.Errorf("account catalog cannot be nil for reconciliation")
	}

	r := &TransactionReconciliation{
		Entry:        entry,
		VisualNotice: GenericVisualNotice,
		AuditNotes:   make([]string, 0),
	}

	// 1. Journal Entry Validation
	if err := entry.Validate(catalog); err != nil {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Journal entry validation failed: %v", err))
		return r, err
	}

	r.TotalDr = entry.TotalDebit()
	r.TotalCr = entry.TotalCredit()
	r.IsBalanced = entry.IsBalanced()

	if !r.IsBalanced {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Journal entry does not balance: Dr %s != Cr %s",
			r.TotalDr.FormatDollars(), r.TotalCr.FormatDollars()))
		return r, fmt.Errorf("journal entry is unbalanced")
	}
	r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✓ Journal Entry balanced: Total Debits (%s) = Total Credits (%s)",
		r.TotalDr.FormatDollars(), r.TotalCr.FormatDollars()))

	// 2. T-Account Ledger Construction and Posting Tie-out
	ledger := NewTAccountLedger(catalog)
	ledger.PostEntry(entry)
	r.Ledger = ledger
	r.TAccountDrSum = ledger.TotalDebits()
	r.TAccountCrSum = ledger.TotalCredits()

	if r.TAccountDrSum != r.TotalDr {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ T-Account Dr sum (%s) != Journal Dr total (%s)",
			r.TAccountDrSum.FormatDollars(), r.TotalDr.FormatDollars()))
		return r, fmt.Errorf("T-account debit sum mismatch")
	}
	if r.TAccountCrSum != r.TotalCr {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ T-Account Cr sum (%s) != Journal Cr total (%s)",
			r.TAccountCrSum.FormatDollars(), r.TotalCr.FormatDollars()))
		return r, fmt.Errorf("T-account credit sum mismatch")
	}
	r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✓ T-Account postings tie: Σ Debits (%s) = Σ Credits (%s)",
		r.TAccountDrSum.FormatDollars(), r.TAccountCrSum.FormatDollars()))

	// 3. Equation Delta Calculation
	eq, err := entry.EquationEffects(catalog)
	if err != nil {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Equation calculation failed: %v", err))
		return r, err
	}
	r.EquationDelta = eq
	r.DeltaAssets = eq.DeltaAssets
	r.DeltaLiabilities = eq.DeltaLiabilities
	r.DeltaEquity = eq.DeltaEquity

	// 4. Verify Ledger Net Balance Changes Reconcile with Equation Delta
	var netAssets, netLiab, netEquity Money
	for _, ta := range ledger.AccountsInOrder() {
		acc, _ := catalog.Get(ta.AccountID)
		switch ta.Category {
		case CategoryAsset:
			if acc.ContraOf != "" {
				// Contra asset: credit increases contra, debit decreases
				netAssets = netAssets.Sub(ta.TotalCredits.Sub(ta.TotalDebits))
			} else {
				netAssets = netAssets.Add(ta.TotalDebits.Sub(ta.TotalCredits))
			}
		case CategoryLiability:
			netLiab = netLiab.Add(ta.TotalCredits.Sub(ta.TotalDebits))
		case CategoryEquity, CategoryRevenue:
			netEquity = netEquity.Add(ta.TotalCredits.Sub(ta.TotalDebits))
		case CategoryExpense, CategoryDividends:
			netEquity = netEquity.Sub(ta.TotalDebits.Sub(ta.TotalCredits))
		}
	}

	if netAssets != r.DeltaAssets {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Ledger asset change (%s) != Equation ΔAssets (%s)",
			netAssets.FormatDollars(), r.DeltaAssets.FormatDollars()))
		return r, fmt.Errorf("ledger net assets does not match equation delta")
	}
	if netLiab != r.DeltaLiabilities {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Ledger liability change (%s) != Equation ΔLiabilities (%s)",
			netLiab.FormatDollars(), r.DeltaLiabilities.FormatDollars()))
		return r, fmt.Errorf("ledger net liabilities does not match equation delta")
	}
	if netEquity != r.DeltaEquity {
		r.Reconciled = false
		r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✗ Ledger equity change (%s) != Equation ΔEquity (%s)",
			netEquity.FormatDollars(), r.DeltaEquity.FormatDollars()))
		return r, fmt.Errorf("ledger net equity does not match equation delta")
	}

	r.AuditNotes = append(r.AuditNotes, fmt.Sprintf("✓ Balance sheet equation ties: ΔAssets (%s) = ΔLiabilities (%s) + ΔEquity (%s)",
		r.DeltaAssets.FormatDollars(), r.DeltaLiabilities.FormatDollars(), r.DeltaEquity.FormatDollars()))
	r.AuditNotes = append(r.AuditNotes, "✓ Journal, T-Accounts, and Accounting Equation 100% reconciled to the same postings")

	r.Reconciled = true
	return r, nil
}

// RenderSummary produces a comprehensive multi-view reconciliation report.
func (r *TransactionReconciliation) RenderSummary() string {
	var b strings.Builder

	b.WriteString("================================================================================\n")
	b.WriteString("TRANSACTION POSTINGS & LEDGER RECONCILIATION REPORT\n")
	b.WriteString(fmt.Sprintf("[%s]\n", r.VisualNotice))
	b.WriteString("================================================================================\n\n")

	// Section 1: Journal Entry
	b.WriteString("1. JOURNAL ENTRY VIEW:\n")
	b.WriteString("────────────────────────────────────────────────────────────────────────────────\n")
	for _, p := range r.Entry.Postings {
		if p.Side == SideDebit {
			b.WriteString(fmt.Sprintf("  Debit   %-32s %14s\n", p.AccountID, p.Amount.FormatDollars()))
		} else {
			b.WriteString(fmt.Sprintf("    Credit  %-30s %14s\n", p.AccountID, p.Amount.FormatDollars()))
		}
	}
	b.WriteString("────────────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Total Debits: %s | Total Credits: %s (Balanced: %t)\n\n",
		r.TotalDr.FormatDollars(), r.TotalCr.FormatDollars(), r.IsBalanced))

	// Section 2: T-Accounts
	b.WriteString("2. T-ACCOUNT LEDGER CARDS VIEW:\n")
	b.WriteString("────────────────────────────────────────────────────────────────────────────────\n")
	if r.Ledger != nil {
		for _, ta := range r.Ledger.AccountsInOrder() {
			b.WriteString(ta.Render())
			b.WriteString("\n")
		}
	}

	// Section 3: Equation Reconciliation
	b.WriteString("3. ACCOUNTING EQUATION RECONCILIATION (ΔAssets = ΔLiabilities + ΔEquity):\n")
	b.WriteString("────────────────────────────────────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Δ Assets:      %12s\n", r.DeltaAssets.FormatDollars()))
	b.WriteString(fmt.Sprintf("  Δ Liabilities: %12s\n", r.DeltaLiabilities.FormatDollars()))
	b.WriteString(fmt.Sprintf("  Δ Equity:      %12s\n", r.DeltaEquity.FormatDollars()))
	b.WriteString(fmt.Sprintf("  Net Tie: %s = %s + %s\n\n",
		r.DeltaAssets.FormatDollars(), r.DeltaLiabilities.FormatDollars(), r.DeltaEquity.FormatDollars()))

	// Section 4: Audit Verification Notes
	b.WriteString("4. AUDIT PROOF NOTES:\n")
	for _, note := range r.AuditNotes {
		b.WriteString(fmt.Sprintf("  %s\n", note))
	}
	b.WriteString("================================================================================\n")

	return b.String()
}

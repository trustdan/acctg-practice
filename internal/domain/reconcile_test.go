package domain_test

import (
	"strings"
	"testing"

	"github.com/trustdan/acctg-practice/internal/domain"
)

func sampleCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	_ = catalog.Add(domain.Account{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "prepaid_rent", Name: "Prepaid Rent", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "equipment", Name: "Equipment", Category: domain.CategoryAsset, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "accumulated_depreciation", Name: "Accumulated Depreciation", Category: domain.CategoryAsset, NormalSide: domain.SideCredit, ContraOf: "equipment"})
	_ = catalog.Add(domain.Account{ID: "accounts_payable", Name: "Accounts Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "unearned_revenue", Name: "Unearned Revenue", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "notes_payable", Name: "Notes Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "retained_earnings", Name: "Retained Earnings", Category: domain.CategoryEquity, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit})
	_ = catalog.Add(domain.Account{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit})
	_ = catalog.Add(domain.Account{ID: "dividends", Name: "Dividends", Category: domain.CategoryDividends, NormalSide: domain.SideDebit})
	return catalog
}

func TestEntryNormalizationAndComparison(t *testing.T) {
	// Canonical entry: Dr Cash $1,000, Cr Service Revenue $1,000
	canonical := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
	)

	// Variant 1: Different line order (Credit first, Debit second)
	orderReversed := domain.NewEntry(
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
	)
	if !canonical.EqualPostings(orderReversed) {
		t.Errorf("EqualPostings failed for reordered lines")
	}
	if !canonical.EqualNormalized(orderReversed) {
		t.Errorf("EqualNormalized failed for reordered lines")
	}

	// Variant 2: Split debits (two Cash debits of $600 and $400)
	splitDebits := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(60000)},
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(40000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
	)
	if canonical.EqualPostings(splitDebits) {
		t.Errorf("EqualPostings should be false for split lines before aggregation")
	}
	if !canonical.EqualNormalized(splitDebits) {
		t.Errorf("EqualNormalized should be true after aggregating split debit lines")
	}

	// Variant 3: Split credits (two Service Revenue credits of $700 and $300)
	splitCredits := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(70000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(30000)},
	)
	if !canonical.EqualNormalized(splitCredits) {
		t.Errorf("EqualNormalized should be true after aggregating split credit lines")
	}

	// Variant 4: Both sides split
	bothSplit := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(40000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(60000)},
	)
	if !canonical.EqualNormalized(bothSplit) {
		t.Errorf("EqualNormalized should be true when both sides have split lines")
	}

	// Variant 5: Reversed sides (Debit Service Revenue, Credit Cash) -> MUST NOT MATCH
	reversed := domain.NewEntry(
		domain.Posting{AccountID: "service_revenue", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
	)
	if canonical.EqualNormalized(reversed) {
		t.Errorf("EqualNormalized must be false when sides are reversed")
	}

	// Variant 6: Wrong account (Unearned Revenue instead of Service Revenue) -> MUST NOT MATCH
	wrongAccount := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
	)
	if canonical.EqualNormalized(wrongAccount) {
		t.Errorf("EqualNormalized must be false when account is wrong")
	}

	// Variant 7: Wrong amount ($900 instead of $1,000) -> MUST NOT MATCH
	wrongAmount := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(90000)},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(90000)},
	)
	if canonical.EqualNormalized(wrongAmount) {
		t.Errorf("EqualNormalized must be false when amount is wrong")
	}
}

func TestTAccountRenderAndConventions(t *testing.T) {
	catalog := sampleCatalog()
	cashAcc, _ := catalog.Get("cash")
	unearnedAcc, _ := catalog.Get("unearned_revenue")

	// Asset T-Account: Debit increases (+), Credit decreases (-)
	cashT := domain.NewTAccount(cashAcc)
	cashT.Post(domain.SideDebit, domain.NewMoney(150000))
	cashT.Post(domain.SideCredit, domain.NewMoney(30000))

	if cashT.DebitHeading() != "Debit (+)" {
		t.Errorf("expected Debit (+), got %s", cashT.DebitHeading())
	}
	if cashT.CreditHeading() != "Credit (-)" {
		t.Errorf("expected Credit (-), got %s", cashT.CreditHeading())
	}
	if cashT.NetBalance.Cents() != 120000 || cashT.BalanceSide != domain.SideDebit {
		t.Errorf("expected NetBalance $1,200.00 Dr, got %s %s", cashT.NetBalance.FormatDollars(), cashT.BalanceSide)
	}

	cashRender := cashT.Render()
	if !strings.Contains(cashRender, "Cash (Asset)") {
		t.Errorf("expected title in render, got:\n%s", cashRender)
	}
	if !strings.Contains(cashRender, "Debit (+)") || !strings.Contains(cashRender, "Credit (-)") {
		t.Errorf("expected directional column headings in render, got:\n%s", cashRender)
	}
	if !strings.Contains(cashRender, "Bal: $1200.00") {
		t.Errorf("expected balance in render, got:\n%s", cashRender)
	}

	// Liability T-Account: Debit decreases (-), Credit increases (+)
	unearnedT := domain.NewTAccount(unearnedAcc)
	unearnedT.Post(domain.SideCredit, domain.NewMoney(150000))

	if unearnedT.DebitHeading() != "Debit (-)" {
		t.Errorf("expected Debit (-), got %s", unearnedT.DebitHeading())
	}
	if unearnedT.CreditHeading() != "Credit (+)" {
		t.Errorf("expected Credit (+), got %s", unearnedT.CreditHeading())
	}
	if unearnedT.NetBalance.Cents() != 150000 || unearnedT.BalanceSide != domain.SideCredit {
		t.Errorf("expected NetBalance $1,500.00 Cr, got %s %s", unearnedT.NetBalance.FormatDollars(), unearnedT.BalanceSide)
	}

	// Generic visual notice verified
	if !strings.Contains(domain.GenericVisualNotice, "Generic visual baseline") {
		t.Errorf("expected GenericVisualNotice to label generic baseline, got: %s", domain.GenericVisualNotice)
	}
}

func TestReconcileTransaction(t *testing.T) {
	catalog := sampleCatalog()

	t.Run("Customer Advance Reconciliation", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(240000)},
			domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(240000)},
		)

		recon, err := domain.ReconcileTransaction(entry, catalog)
		if err != nil {
			t.Fatalf("unexpected reconciliation error: %v", err)
		}
		if !recon.Reconciled {
			t.Fatalf("expected transaction to be reconciled")
		}
		if recon.TotalDr != recon.TotalCr || recon.TotalDr.Cents() != 240000 {
			t.Errorf("journal Dr != Cr, got Dr %d Cr %d", recon.TotalDr.Cents(), recon.TotalCr.Cents())
		}
		if recon.TAccountDrSum != recon.TAccountCrSum || recon.TAccountDrSum.Cents() != 240000 {
			t.Errorf("t-account Dr sum != Cr sum")
		}
		if recon.DeltaAssets.Cents() != 240000 || recon.DeltaLiabilities.Cents() != 240000 || recon.DeltaEquity.Cents() != 0 {
			t.Errorf("equation delta mismatch: %+v", recon.EquationDelta)
		}

		summary := recon.RenderSummary()
		if !strings.Contains(summary, "TRANSACTION POSTINGS & LEDGER RECONCILIATION REPORT") {
			t.Errorf("expected title in summary")
		}
		if !strings.Contains(summary, domain.GenericVisualNotice) {
			t.Errorf("expected generic visual baseline notice in summary")
		}
		if !strings.Contains(summary, "✓ Journal, T-Accounts, and Accounting Equation 100% reconciled") {
			t.Errorf("expected audit tie note in summary")
		}
	})

	t.Run("Split Equivalent Lines Reconcile to Exact Same Postings and Equation", func(t *testing.T) {
		// Split debits: Dr Cash 1000, Dr Cash 1400, Cr Unearned Revenue 2400
		splitEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(140000)},
			domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(240000)},
		)

		recon, err := domain.ReconcileTransaction(splitEntry, catalog)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !recon.Reconciled {
			t.Fatalf("expected split entry to reconcile completely")
		}
		if recon.DeltaAssets.Cents() != 240000 || recon.DeltaLiabilities.Cents() != 240000 || recon.DeltaEquity.Cents() != 0 {
			t.Errorf("unexpected delta for split entry: %+v", recon.EquationDelta)
		}
		// Cash T-account should have 2 debit postings summing to 240000
		cashTA := recon.Ledger.GetOrCreate("cash")
		if len(cashTA.Debits) != 2 || cashTA.TotalDebits.Cents() != 240000 {
			t.Errorf("expected 2 debits summing to $2,400 in cash T-account, got %d postings, sum %d",
				len(cashTA.Debits), cashTA.TotalDebits.Cents())
		}
	})

	t.Run("Asset Swap: Equipment purchase for cash", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "equipment", Side: domain.SideDebit, Amount: domain.NewMoney(500000)},
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(500000)},
		)
		recon, err := domain.ReconcileTransaction(entry, catalog)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if recon.DeltaAssets.Cents() != 0 || recon.DeltaLiabilities.Cents() != 0 || recon.DeltaEquity.Cents() != 0 {
			t.Errorf("asset swap must have net zero delta on assets: %+v", recon.EquationDelta)
		}
	})

	t.Run("Operating Expense: Rent paid in cash", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: domain.NewMoney(180000)},
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(180000)},
		)
		recon, err := domain.ReconcileTransaction(entry, catalog)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if recon.DeltaAssets.Cents() != -180000 || recon.DeltaLiabilities.Cents() != 0 || recon.DeltaEquity.Cents() != -180000 {
			t.Errorf("operating expense delta mismatch: %+v", recon.EquationDelta)
		}
	})

	t.Run("Unbalanced Entry Fails Reconciliation", func(t *testing.T) {
		unbalanced := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(80000)},
		)
		recon, err := domain.ReconcileTransaction(unbalanced, catalog)
		if err == nil {
			t.Fatalf("expected error for unbalanced entry")
		}
		if recon != nil && recon.Reconciled {
			t.Fatalf("unbalanced entry must not be marked reconciled")
		}
	})
}

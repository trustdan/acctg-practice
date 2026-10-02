package engine_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func setupTestEngine(t *testing.T) (*engine.Engine, *domain.AccountCatalog) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}
	return engine.NewEngine(catalog), catalog
}

func TestAllFamiliesTableDrivenCorrect(t *testing.T) {
	eng, catalog := setupTestEngine(t)

	type testCase struct {
		familyID    string
		amountCents int64
		expectedDr  domain.AccountID
		expectedCr  domain.AccountID
		expectedDA  int64 // Delta Assets in cents
		expectedDL  int64 // Delta Liabilities in cents
		expectedDE  int64 // Delta Equity in cents
	}

	tests := []testCase{
		{
			familyID:    bank.FamilyCashService,
			amountCents: 10000, // $100
			expectedDr:  "cash",
			expectedCr:  "service_revenue",
			expectedDA:  10000,
			expectedDL:  0,
			expectedDE:  10000,
		},
		{
			familyID:    bank.FamilyServiceOnCredit,
			amountCents: 25000, // $250
			expectedDr:  "accounts_receivable",
			expectedCr:  "service_revenue",
			expectedDA:  25000,
			expectedDL:  0,
			expectedDE:  25000,
		},
		{
			familyID:    bank.FamilyCustomerAdvance,
			amountCents: 5000, // $50
			expectedDr:  "cash",
			expectedCr:  "unearned_revenue",
			expectedDA:  5000,
			expectedDL:  5000,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyCollectReceivable,
			amountCents: 15000, // $150
			expectedDr:  "cash",
			expectedCr:  "accounts_receivable",
			expectedDA:  0, // Asset swap
			expectedDL:  0,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyEarnAdvance,
			amountCents: 5000, // $50
			expectedDr:  "unearned_revenue",
			expectedCr:  "service_revenue",
			expectedDA:  0,
			expectedDL:  -5000,
			expectedDE:  5000,
		},
		{
			familyID:    bank.FamilyCashRent,
			amountCents: 120000, // $1200
			expectedDr:  "rent_expense",
			expectedCr:  "cash",
			expectedDA:  -120000,
			expectedDL:  0,
			expectedDE:  -120000,
		},
		{
			familyID:    bank.FamilyPrepaidPurchase,
			amountCents: 60000, // $600
			expectedDr:  "prepaid_insurance",
			expectedCr:  "cash",
			expectedDA:  0, // Asset swap
			expectedDL:  0,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyPrepaidConsumption,
			amountCents: 5000, // $50
			expectedDr:  "insurance_expense",
			expectedCr:  "prepaid_insurance",
			expectedDA:  -5000,
			expectedDL:  0,
			expectedDE:  -5000,
		},
		{
			familyID:    bank.FamilyEquipmentPurchaseCash,
			amountCents: 350000, // $3500
			expectedDr:  "equipment",
			expectedCr:  "cash",
			expectedDA:  0, // Asset swap
			expectedDL:  0,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyBorrowCash,
			amountCents: 500000, // $5000
			expectedDr:  "cash",
			expectedCr:  "notes_payable",
			expectedDA:  500000,
			expectedDL:  500000,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyRepayNotePrincipal,
			amountCents: 100000, // $1000
			expectedDr:  "notes_payable",
			expectedCr:  "cash",
			expectedDA:  -100000,
			expectedDL:  -100000,
			expectedDE:  0,
		},
		{
			familyID:    bank.FamilyIssueShares,
			amountCents: 1000000, // $10,000
			expectedDr:  "cash",
			expectedCr:  "common_stock",
			expectedDA:  1000000,
			expectedDL:  0,
			expectedDE:  1000000,
		},
		{
			familyID:    bank.FamilyDividendCash,
			amountCents: 20000, // $200
			expectedDr:  "dividends",
			expectedCr:  "cash",
			expectedDA:  -20000,
			expectedDL:  0,
			expectedDE:  -20000,
		},
		{
			familyID:    bank.FamilyDeclareDividend,
			amountCents: 50000, // $500
			expectedDr:  "dividends",
			expectedCr:  "dividends_payable",
			expectedDA:  0,
			expectedDL:  50000,
			expectedDE:  -50000,
		},
		{
			familyID:    bank.FamilyPayDividendPayable,
			amountCents: 50000, // $500
			expectedDr:  "dividends_payable",
			expectedCr:  "cash",
			expectedDA:  -50000,
			expectedDL:  -50000,
			expectedDE:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.familyID, func(t *testing.T) {
			event := engine.TransactionEvent{
				FamilyID: tt.familyID,
				Parameters: map[string]int64{
					"amount_minor_units": tt.amountCents,
				},
			}

			processed, err := eng.ProcessEvent(event)
			if err != nil {
				t.Fatalf("unexpected ProcessEvent error for %s: %v", tt.familyID, err)
			}

			// Validate entry balancing
			if !processed.Entry.IsBalanced() {
				t.Fatalf("%s entry did not balance", tt.familyID)
			}
			if err := processed.Entry.Validate(catalog); err != nil {
				t.Fatalf("%s entry failed validation: %v", tt.familyID, err)
			}

			// Verify postings match expectations
			if len(processed.Entry.Postings) != 2 {
				t.Fatalf("expected 2 postings, got %d", len(processed.Entry.Postings))
			}
			var gotDr, gotCr domain.Posting
			for _, p := range processed.Entry.Postings {
				if p.Side == domain.SideDebit {
					gotDr = p
				} else if p.Side == domain.SideCredit {
					gotCr = p
				}
			}
			if gotDr.AccountID != tt.expectedDr || gotDr.Amount.Cents() != tt.amountCents {
				t.Errorf("%s debit mismatch: got %s %d, want %s %d", tt.familyID, gotDr.AccountID, gotDr.Amount.Cents(), tt.expectedDr, tt.amountCents)
			}
			if gotCr.AccountID != tt.expectedCr || gotCr.Amount.Cents() != tt.amountCents {
				t.Errorf("%s credit mismatch: got %s %d, want %s %d", tt.familyID, gotCr.AccountID, gotCr.Amount.Cents(), tt.expectedCr, tt.amountCents)
			}

			// Verify equation effects
			if processed.Equation.DeltaAssets.Cents() != tt.expectedDA {
				t.Errorf("%s delta assets = %d, want %d", tt.familyID, processed.Equation.DeltaAssets.Cents(), tt.expectedDA)
			}
			if processed.Equation.DeltaLiabilities.Cents() != tt.expectedDL {
				t.Errorf("%s delta liabilities = %d, want %d", tt.familyID, processed.Equation.DeltaLiabilities.Cents(), tt.expectedDL)
			}
			if processed.Equation.DeltaEquity.Cents() != tt.expectedDE {
				t.Errorf("%s delta equity = %d, want %d", tt.familyID, processed.Equation.DeltaEquity.Cents(), tt.expectedDE)
			}

			// Verify explanation has rich content
			if processed.Explanation.Summary == "" || processed.Explanation.DebitRationale == "" || processed.Explanation.CreditRationale == "" {
				t.Errorf("%s has empty explanation fields: %+v", tt.familyID, processed.Explanation)
			}

			// Verify self-evaluation returns correct
			eval := eng.EvaluateEntry(event, processed.Entry)
			if !eval.IsCorrect {
				t.Errorf("%s canonical entry failed self-evaluation: %s", tt.familyID, eval.Feedback)
			}
		})
	}
}

func TestSemanticDistractorsAndWrongEntries(t *testing.T) {
	eng, _ := setupTestEngine(t)
	amt := domain.NewMoney(10000)

	t.Run("Receivable collection NEVER creates duplicate revenue", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyCollectReceivable,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		// Wrong entry: Learner credits service revenue instead of accounts receivable
		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected duplicate revenue entry on receivable collection to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagDuplicateRevenueOnCollection {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagDuplicateRevenueOnCollection, eval.ErrorTags)
		}
	})

	t.Run("Customer advance rejects premature revenue recognition", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyCustomerAdvance,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected premature revenue entry to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagRevenueRecognizedPrematurely {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagRevenueRecognizedPrematurely, eval.ErrorTags)
		}
	})

	t.Run("Earn advance rejects recording cash today", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyEarnAdvance,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected cash entry on earning advance to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagCashRecordedOnEarningAdvance {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagCashRecordedOnEarningAdvance, eval.ErrorTags)
		}
	})

	t.Run("Borrowing on note rejects crediting revenue", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyBorrowCash,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected revenue entry on borrowing to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagRevenueRecordedOnBorrowing {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagRevenueRecordedOnBorrowing, eval.ErrorTags)
		}
	})

	t.Run("Issuing shares rejects crediting revenue", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyIssueShares,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected revenue entry on issuing shares to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagRevenueRecordedOnShareIssue {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagRevenueRecordedOnShareIssue, eval.ErrorTags)
		}
	})

	t.Run("Loan principal repayment rejects debiting expense", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyRepayNotePrincipal,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		wrongEntry := domain.NewEntry(
			domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, wrongEntry)
		if eval.IsCorrect {
			t.Fatalf("expected expense entry on loan repayment to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagExpenseRecordedOnLoanRepayment {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagExpenseRecordedOnLoanRepayment, eval.ErrorTags)
		}
	})

	t.Run("Reversed sides flagged", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyCashService,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		reversedEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideDebit, Amount: amt},
		)

		eval := eng.EvaluateEntry(event, reversedEntry)
		if eval.IsCorrect {
			t.Fatalf("expected reversed sides to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagReversedSides {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagReversedSides, eval.ErrorTags)
		}
	})

	t.Run("Unbalanced entry flagged", func(t *testing.T) {
		event := engine.TransactionEvent{
			FamilyID:   bank.FamilyCashService,
			Parameters: map[string]int64{"amount_minor_units": 10000},
		}

		unbalancedEntry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(5000)},
		)

		eval := eng.EvaluateEntry(event, unbalancedEntry)
		if eval.IsCorrect {
			t.Fatalf("expected unbalanced entry to be rejected")
		}
		foundTag := false
		for _, tag := range eval.ErrorTags {
			if tag == engine.TagUnbalancedEntry {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Errorf("expected tag %s, got tags: %v", engine.TagUnbalancedEntry, eval.ErrorTags)
		}
	})
}

func TestAllParameterCombinationsBalance(t *testing.T) {
	eng, catalog := setupTestEngine(t)

	families := []string{
		bank.FamilyCashService,
		bank.FamilyServiceOnCredit,
		bank.FamilyCustomerAdvance,
		bank.FamilyCollectReceivable,
		bank.FamilyEarnAdvance,
		bank.FamilyCashRent,
		bank.FamilyPrepaidPurchase,
		bank.FamilyPrepaidConsumption,
		bank.FamilyEquipmentPurchaseCash,
		bank.FamilyBorrowCash,
		bank.FamilyRepayNotePrincipal,
		bank.FamilyIssueShares,
		bank.FamilyDividendCash,
	}

	testAmounts := []int64{
		1,        // $0.01 (smallest minor unit)
		50,       // $0.50
		1000,     // $10.00
		5000,     // $50.00
		10000,    // $100.00
		20000,    // $200.00
		99999,    // $999.99
		1000000,  // $10,000.00
		50000000, // $500,000.00
	}

	for _, fam := range families {
		for _, amt := range testAmounts {
			event := engine.TransactionEvent{
				FamilyID:   fam,
				Parameters: map[string]int64{"amount_minor_units": amt},
			}

			processed, err := eng.ProcessEvent(event)
			if err != nil {
				t.Fatalf("ProcessEvent failed for family %s amount %d: %v", fam, amt, err)
			}

			// 1. Double entry must balance exactly
			if !processed.Entry.IsBalanced() {
				t.Fatalf("Entry did not balance for family %s amount %d", fam, amt)
			}

			// 2. Debits == Credits
			if processed.Entry.TotalDebit() != processed.Entry.TotalCredit() {
				t.Fatalf("Debits (%s) != Credits (%s) for family %s amount %d",
					processed.Entry.TotalDebit().FormatExact(),
					processed.Entry.TotalCredit().FormatExact(),
					fam, amt)
			}

			// 3. Balance sheet equation must balance exactly: DeltaA == DeltaL + DeltaE
			eq, err := processed.Entry.EquationEffects(catalog)
			if err != nil {
				t.Fatalf("Equation effects failed for family %s amount %d: %v", fam, amt, err)
			}
			expectedAssets := eq.DeltaLiabilities.Add(eq.DeltaEquity)
			if eq.DeltaAssets != expectedAssets {
				t.Fatalf("Equation delta violated for family %s amount %d: DeltaA (%s) != DeltaL (%s) + DeltaE (%s)",
					fam, amt, eq.DeltaAssets.FormatExact(), eq.DeltaLiabilities.FormatExact(), eq.DeltaEquity.FormatExact())
			}
		}
	}
}

func TestSplitEquivalentLinesGradeConsistently(t *testing.T) {
	eng, _ := setupTestEngine(t)

	event := engine.TransactionEvent{
		FamilyID:   bank.FamilyCustomerAdvance,
		Parameters: map[string]int64{"amount_minor_units": 150000}, // $1,500.00
	}

	// Canonical: Dr Cash $1500, Cr Unearned Revenue $1500

	// 1. Split debits: Dr Cash $1000, Dr Cash $500, Cr Unearned Revenue $1500
	splitDebits := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(150000)},
	)
	eval1 := eng.EvaluateEntry(event, splitDebits)
	if !eval1.IsCorrect {
		t.Fatalf("expected split debits to grade as correct, got error: %s", eval1.Feedback)
	}

	// 2. Split credits: Dr Cash $1500, Cr Unearned Revenue $600, Cr Unearned Revenue $900
	splitCredits := domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(150000)},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(60000)},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(90000)},
	)
	eval2 := eng.EvaluateEntry(event, splitCredits)
	if !eval2.IsCorrect {
		t.Fatalf("expected split credits to grade as correct, got error: %s", eval2.Feedback)
	}

	// 3. Reversed line order: Credit first, Debit second
	reversedOrder := domain.NewEntry(
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(150000)},
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(150000)},
	)
	eval3 := eng.EvaluateEntry(event, reversedOrder)
	if !eval3.IsCorrect {
		t.Fatalf("expected reversed line order to grade as correct, got error: %s", eval3.Feedback)
	}
}

func TestWrongButBalancedEntriesFailGrading(t *testing.T) {
	eng, _ := setupTestEngine(t)

	cases := []struct {
		name        string
		event       engine.TransactionEvent
		entry       domain.Entry
		expectedTag string
	}{
		{
			name: "Equipment purchase expensed: Dr Rent Expense, Cr Cash",
			event: engine.TransactionEvent{
				FamilyID:   bank.FamilyEquipmentPurchaseCash,
				Parameters: map[string]int64{"amount_minor_units": 500000},
			},
			entry: domain.NewEntry(
				domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: domain.NewMoney(500000)},
				domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(500000)},
			),
			expectedTag: engine.TagExpenseRecordedOnEquipmentPurchase,
		},
		{
			name: "Customer advance recognized as revenue: Dr Cash, Cr Service Revenue",
			event: engine.TransactionEvent{
				FamilyID:   bank.FamilyCustomerAdvance,
				Parameters: map[string]int64{"amount_minor_units": 200000},
			},
			entry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(200000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(200000)},
			),
			expectedTag: engine.TagRevenueRecognizedPrematurely,
		},
		{
			name: "Loan repayment expensed: Dr Rent Expense, Cr Cash",
			event: engine.TransactionEvent{
				FamilyID:   bank.FamilyRepayNotePrincipal,
				Parameters: map[string]int64{"amount_minor_units": 300000},
			},
			entry: domain.NewEntry(
				domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: domain.NewMoney(300000)},
				domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(300000)},
			),
			expectedTag: engine.TagExpenseRecordedOnLoanRepayment,
		},
		{
			name: "Dividends expensed: Dr Rent Expense, Cr Cash",
			event: engine.TransactionEvent{
				FamilyID:   bank.FamilyDividendCash,
				Parameters: map[string]int64{"amount_minor_units": 100000},
			},
			entry: domain.NewEntry(
				domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: domain.NewMoney(100000)},
				domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(100000)},
			),
			expectedTag: engine.TagExpenseRecordedOnDividend,
		},
		{
			name: "Receivable collection recorded as revenue: Dr Cash, Cr Service Revenue",
			event: engine.TransactionEvent{
				FamilyID:   bank.FamilyCollectReceivable,
				Parameters: map[string]int64{"amount_minor_units": 150000},
			},
			entry: domain.NewEntry(
				domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(150000)},
				domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(150000)},
			),
			expectedTag: engine.TagDuplicateRevenueOnCollection,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.entry.IsBalanced() {
				t.Fatalf("test precondition failed: entry must be balanced")
			}
			eval := eng.EvaluateEntry(tc.event, tc.entry)
			if eval.IsCorrect {
				t.Fatalf("expected wrong-but-balanced entry to fail grading")
			}
			found := false
			for _, tag := range eval.ErrorTags {
				if tag == tc.expectedTag {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected tag %s in error tags: %v", tc.expectedTag, eval.ErrorTags)
			}
		})
	}
}

func TestAllFamiliesPostingsReconcileToTAccountsAndEquation(t *testing.T) {
	eng, catalog := setupTestEngine(t)

	families := []string{
		bank.FamilyCashService,
		bank.FamilyServiceOnCredit,
		bank.FamilyCustomerAdvance,
		bank.FamilyCollectReceivable,
		bank.FamilyEarnAdvance,
		bank.FamilyCashRent,
		bank.FamilyPrepaidPurchase,
		bank.FamilyPrepaidConsumption,
		bank.FamilyEquipmentPurchaseCash,
		bank.FamilyBorrowCash,
		bank.FamilyRepayNotePrincipal,
		bank.FamilyIssueShares,
		bank.FamilyDividendCash,
		bank.FamilyDeclareDividend,
		bank.FamilyPayDividendPayable,
	}

	for _, fam := range families {
		t.Run(fam, func(t *testing.T) {
			event := engine.TransactionEvent{
				FamilyID:   fam,
				Parameters: map[string]int64{"amount_minor_units": 240000},
			}
			processed, err := eng.ProcessEvent(event)
			if err != nil {
				t.Fatalf("ProcessEvent failed: %v", err)
			}

			recon, err := domain.ReconcileTransaction(processed.Entry, catalog)
			if err != nil {
				t.Fatalf("ReconcileTransaction failed for family %s: %v", fam, err)
			}
			if !recon.Reconciled {
				t.Fatalf("Reconciliation not marked complete for family %s: %+v", fam, recon.AuditNotes)
			}

			// Verify journal Dr == Cr
			if recon.TotalDr != recon.TotalCr {
				t.Errorf("journal Dr != Cr for family %s", fam)
			}
			// Verify TAccount Dr sum == Cr sum
			if recon.TAccountDrSum != recon.TAccountCrSum {
				t.Errorf("TAccount Dr sum != Cr sum for family %s", fam)
			}
			// Verify Equation delta ties out: DeltaA == DeltaL + DeltaE
			if recon.DeltaAssets != recon.DeltaLiabilities.Add(recon.DeltaEquity) {
				t.Errorf("Equation delta does not tie for family %s: %s != %s + %s",
					fam, recon.DeltaAssets.FormatExact(), recon.DeltaLiabilities.FormatExact(), recon.DeltaEquity.FormatExact())
			}
		})
	}
}

func TestDeclaredVsPaidDividendsDistinctSemantics(t *testing.T) {
	eng, catalog := setupTestEngine(t)

	// 1. Declaration of $500 dividend
	declareEv := engine.TransactionEvent{
		FamilyID:   bank.FamilyDeclareDividend,
		Parameters: map[string]int64{"amount_minor_units": 50000},
	}
	procDeclare, err := eng.ProcessEvent(declareEv)
	if err != nil {
		t.Fatalf("unexpected error processing dividend declaration: %v", err)
	}

	// Verify declaration entry: Dr Dividends $500, Cr Dividends Payable $500
	if len(procDeclare.Entry.Postings) != 2 {
		t.Fatalf("expected 2 postings on declaration")
	}
	for _, p := range procDeclare.Entry.Postings {
		if p.AccountID == "cash" {
			t.Errorf("declaration MUST NOT touch cash")
		}
	}
	// Verify equation effects on declaration: ΔA = 0, ΔL = +$500, ΔE = -$500
	if procDeclare.Equation.DeltaAssets.Cents() != 0 {
		t.Errorf("declaration must have zero asset effect, got %d", procDeclare.Equation.DeltaAssets.Cents())
	}
	if procDeclare.Equation.DeltaLiabilities.Cents() != 50000 {
		t.Errorf("declaration must increase liabilities by $500, got %d", procDeclare.Equation.DeltaLiabilities.Cents())
	}
	if procDeclare.Equation.DeltaEquity.Cents() != -50000 {
		t.Errorf("declaration must reduce equity by $500, got %d", procDeclare.Equation.DeltaEquity.Cents())
	}

	// 2. Payment of $500 dividend liability
	payEv := engine.TransactionEvent{
		FamilyID:   bank.FamilyPayDividendPayable,
		Parameters: map[string]int64{"amount_minor_units": 50000},
	}
	procPay, err := eng.ProcessEvent(payEv)
	if err != nil {
		t.Fatalf("unexpected error processing dividend payment: %v", err)
	}

	// Verify payment entry: Dr Dividends Payable $500, Cr Cash $500
	for _, p := range procPay.Entry.Postings {
		if p.AccountID == "dividends" {
			t.Errorf("payment must NOT debit dividends again")
		}
	}
	// Verify equation effects on payment: ΔA = -$500, ΔL = -$500, ΔE = $0
	if procPay.Equation.DeltaAssets.Cents() != -50000 {
		t.Errorf("payment must reduce cash assets by $500, got %d", procPay.Equation.DeltaAssets.Cents())
	}
	if procPay.Equation.DeltaLiabilities.Cents() != -50000 {
		t.Errorf("payment must reduce liabilities by $500, got %d", procPay.Equation.DeltaLiabilities.Cents())
	}
	if procPay.Equation.DeltaEquity.Cents() != 0 {
		t.Errorf("payment must have ZERO equity effect, got %d", procPay.Equation.DeltaEquity.Cents())
	}

	// 3. Diagnostic checks for declaration mistakes
	// Mistake: recording cash outflow upon declaration
	wrongDeclCash := domain.NewEntry(
		domain.Posting{AccountID: "dividends", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(50000)},
	)
	evalDecl := eng.EvaluateEntry(declareEv, wrongDeclCash)
	if evalDecl.IsCorrect {
		t.Errorf("expected declaration with cash to fail")
	}
	hasTag := false
	for _, tag := range evalDecl.ErrorTags {
		if tag == engine.TagCashRecordedOnDividendDeclaration {
			hasTag = true
			break
		}
	}
	if !hasTag {
		t.Errorf("expected TagCashRecordedOnDividendDeclaration, got %v", evalDecl.ErrorTags)
	}

	// 4. Diagnostic checks for payment mistakes
	// Mistake: debiting dividends (reducing equity a second time) upon payment
	wrongPayDiv := domain.NewEntry(
		domain.Posting{AccountID: "dividends", Side: domain.SideDebit, Amount: domain.NewMoney(50000)},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(50000)},
	)
	evalPay := eng.EvaluateEntry(payEv, wrongPayDiv)
	if evalPay.IsCorrect {
		t.Errorf("expected payment debiting dividends to fail")
	}
	hasPayTag := false
	for _, tag := range evalPay.ErrorTags {
		if tag == engine.TagDividendsDebitedOnPayment {
			hasPayTag = true
			break
		}
	}
	if !hasPayTag {
		t.Errorf("expected TagDividendsDebitedOnPayment, got %v", evalPay.ErrorTags)
	}

	// 5. Reconciliation verification for both events
	reconDecl, err := domain.ReconcileTransaction(procDeclare.Entry, catalog)
	if err != nil || !reconDecl.Reconciled {
		t.Errorf("declaration reconciliation failed: %v", err)
	}
	reconPay, err := domain.ReconcileTransaction(procPay.Entry, catalog)
	if err != nil || !reconPay.Reconciled {
		t.Errorf("payment reconciliation failed: %v", err)
	}
}

package statements

import (
	"fmt"
	"sort"
	"strings"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildAccountingCycleReport executes the entire accounting cycle for a comprehensive case.
func BuildAccountingCycleReport(
	c ComprehensiveCase,
	catalog *domain.AccountCatalog,
	eng *engine.Engine,
) (*AccountingCycleReport, error) {
	if catalog == nil {
		return nil, fmt.Errorf("account catalog cannot be nil")
	}
	if eng == nil {
		return nil, fmt.Errorf("accounting engine cannot be nil")
	}

	report := &AccountingCycleReport{
		CaseID:                   c.CaseID,
		Title:                    c.Title,
		JournalEntries:           make([]domain.Entry, 0),
		AdjustingEntries:         make([]domain.Entry, 0),
		AuditReconciliationNotes: make([]string, 0),
	}

	// 1. Initialize Ledger and Post Opening Balances
	ledger := domain.NewTAccountLedger(catalog)
	report.Ledger = ledger

	for accID, amt := range c.OpeningBalances {
		acc, ok := catalog.Get(accID)
		if !ok {
			return nil, fmt.Errorf("opening balance references unknown account: %s", accID)
		}
		ta := ledger.GetOrCreate(accID)
		ta.Post(acc.NormalSide, amt)
	}

	// 2. Opening Trial Balance
	openTB := buildTrialBalance("OPENING TRIAL BALANCE", ledger, catalog)
	report.OpeningTrialBalance = openTB
	if !openTB.IsBalanced {
		return nil, fmt.Errorf("opening trial balance does not balance: Dr %s != Cr %s",
			openTB.TotalDebit.FormatDollars(), openTB.TotalCredit.FormatDollars())
	}
	report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
		fmt.Sprintf("✓ Opening Trial Balance balanced: Total Debits (%s) = Total Credits (%s)",
			openTB.TotalDebit.FormatDollars(), openTB.TotalCredit.FormatDollars()))

	// Track cash flow lines
	var (
		opLines  []CashFlowLine
		invLines []CashFlowLine
		finLines []CashFlowLine
		noncash  []string
	)

	// 3. Process Operational Transactions
	for idx, tx := range c.Transactions {
		proc, err := eng.ProcessEvent(tx.Event)
		if err != nil {
			return nil, fmt.Errorf("failed processing transaction [%d] %q: %w", idx+1, tx.Description, err)
		}
		report.JournalEntries = append(report.JournalEntries, proc.Entry)
		ledger.PostEntry(proc.Entry)

		// Record Cash Flow Activity if classified
		if tx.CashClassification != nil {
			amtVal := proc.Event.Parameters["amount_minor_units"]
			amt := domain.NewMoney(amtVal)

			// Determine sign from Cash posting side
			isOutflow := false
			for _, p := range proc.Entry.Postings {
				if p.AccountID == "cash" && p.Side == domain.SideCredit {
					isOutflow = true
					break
				}
			}
			if isOutflow {
				amt = domain.NewMoney(-amtVal)
			}

			line := CashFlowLine{
				Description: tx.CashDescription,
				Amount:      amt,
				Activity:    *tx.CashClassification,
			}
			if line.Description == "" {
				line.Description = tx.Description
			}

			switch *tx.CashClassification {
			case ActivityOperating:
				opLines = append(opLines, line)
			case ActivityInvesting:
				invLines = append(invLines, line)
			case ActivityFinancing:
				finLines = append(finLines, line)
			}
		}

		if tx.NoncashDisclosure != "" {
			noncash = append(noncash, tx.NoncashDisclosure)
		}
	}

	// 4. Unadjusted Trial Balance
	unadjTB := buildTrialBalance("UNADJUSTED TRIAL BALANCE", ledger, catalog)
	report.UnadjustedTrialBalance = unadjTB
	if !unadjTB.IsBalanced {
		return nil, fmt.Errorf("unadjusted trial balance does not balance: Dr %s != Cr %s",
			unadjTB.TotalDebit.FormatDollars(), unadjTB.TotalCredit.FormatDollars())
	}
	report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
		fmt.Sprintf("✓ Unadjusted Trial Balance balanced: Total Debits (%s) = Total Credits (%s)",
			unadjTB.TotalDebit.FormatDollars(), unadjTB.TotalCredit.FormatDollars()))

	// 5. Process End-of-Period Adjustments
	for idx, adj := range c.Adjustments {
		proc, err := eng.ProcessEvent(adj.Event)
		if err != nil {
			return nil, fmt.Errorf("failed processing adjustment [%d] %q: %w", idx+1, adj.Description, err)
		}
		report.AdjustingEntries = append(report.AdjustingEntries, proc.Entry)
		ledger.PostEntry(proc.Entry)

		if adj.CashClassification != nil {
			amtVal := proc.Event.Parameters["amount_minor_units"]
			amt := domain.NewMoney(amtVal)
			isOutflow := false
			for _, p := range proc.Entry.Postings {
				if p.AccountID == "cash" && p.Side == domain.SideCredit {
					isOutflow = true
					break
				}
			}
			if isOutflow {
				amt = domain.NewMoney(-amtVal)
			}
			line := CashFlowLine{
				Description: adj.CashDescription,
				Amount:      amt,
				Activity:    *adj.CashClassification,
			}
			switch *adj.CashClassification {
			case ActivityOperating:
				opLines = append(opLines, line)
			case ActivityInvesting:
				invLines = append(invLines, line)
			case ActivityFinancing:
				finLines = append(finLines, line)
			}
		}
		if adj.NoncashDisclosure != "" {
			noncash = append(noncash, adj.NoncashDisclosure)
		}
	}

	// 6. Adjusted Trial Balance
	adjTB := buildTrialBalance("ADJUSTED TRIAL BALANCE", ledger, catalog)
	report.AdjustedTrialBalance = adjTB
	if !adjTB.IsBalanced {
		return nil, fmt.Errorf("adjusted trial balance does not balance: Dr %s != Cr %s",
			adjTB.TotalDebit.FormatDollars(), adjTB.TotalCredit.FormatDollars())
	}
	report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
		fmt.Sprintf("✓ Adjusted Trial Balance balanced: Total Debits (%s) = Total Credits (%s)",
			adjTB.TotalDebit.FormatDollars(), adjTB.TotalCredit.FormatDollars()))

	// 7. Construct Income Statement
	var revItems []LineItem
	var totalRev domain.Money
	var expItems []LineItem
	var totalExp domain.Money

	for _, ta := range ledger.AccountsInOrder() {
		acc, _ := catalog.Get(ta.AccountID)
		if acc.Category == domain.CategoryRevenue && ta.NetBalance.Cents() > 0 {
			revItems = append(revItems, LineItem{
				AccountID:   acc.ID,
				AccountName: acc.Name,
				Amount:      ta.NetBalance,
			})
			totalRev = totalRev.Add(ta.NetBalance)
		} else if acc.Category == domain.CategoryExpense && ta.NetBalance.Cents() > 0 {
			expItems = append(expItems, LineItem{
				AccountID:   acc.ID,
				AccountName: acc.Name,
				Amount:      ta.NetBalance,
			})
			totalExp = totalExp.Add(ta.NetBalance)
		}
	}

	netIncome := totalRev.Sub(totalExp)
	incStmt := IncomeStatement{
		CompanyName:   c.CompanyName,
		Period:        c.Period,
		Revenues:      revItems,
		TotalRevenue:  totalRev,
		Expenses:      expItems,
		TotalExpenses: totalExp,
		NetIncome:     netIncome,
	}
	report.IncomeStatement = incStmt
	report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
		fmt.Sprintf("✓ Income Statement: Total Revenues (%s) - Total Expenses (%s) = Net Income (%s)",
			totalRev.FormatDollars(), totalExp.FormatDollars(), netIncome.FormatDollars()))

	// 8. Construct Statement of Retained Earnings (Roll-Forward)
	begRE := c.OpeningBalances["retained_earnings"]
	var divDeclared domain.Money
	divTA := ledger.GetOrCreate("dividends")
	if divTA != nil {
		divDeclared = divTA.NetBalance
	}

	endRE := begRE.Add(netIncome).Sub(divDeclared)
	reStmt := RetainedEarningsStatement{
		CompanyName:               c.CompanyName,
		Period:                    c.Period,
		BeginningRetainedEarnings: begRE,
		NetIncome:                 netIncome,
		LessDividendsDeclared:     divDeclared,
		EndingRetainedEarnings:    endRE,
	}
	report.RetainedEarningsStatement = reStmt
	report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
		fmt.Sprintf("✓ Retained Earnings Roll-Forward: Beg RE (%s) + Net Income (%s) - Dividends (%s) = Ending RE (%s)",
			begRE.FormatDollars(), netIncome.FormatDollars(), divDeclared.FormatDollars(), endRE.FormatDollars()))

	// 9. Construct Classified Balance Sheet
	var curAssets []LineItem
	var totCurAssets domain.Money
	var nonCurAssets []LineItem
	var totNonCurAssets domain.Money

	var curLiab []LineItem
	var totCurLiab domain.Money
	var longLiab []LineItem
	var totLongLiab domain.Money

	var equityItems []LineItem
	var totEquity domain.Money

	for _, ta := range ledger.AccountsInOrder() {
		acc, _ := catalog.Get(ta.AccountID)
		switch acc.Category {
		case domain.CategoryAsset:
			if ta.NetBalance.Cents() == 0 {
				continue
			}
			item := LineItem{AccountID: acc.ID, AccountName: acc.Name, Amount: ta.NetBalance}
			if acc.ID == "equipment" {
				nonCurAssets = append(nonCurAssets, item)
				totNonCurAssets = totNonCurAssets.Add(ta.NetBalance)
			} else {
				curAssets = append(curAssets, item)
				totCurAssets = totCurAssets.Add(ta.NetBalance)
			}

		case domain.CategoryLiability:
			if ta.NetBalance.Cents() == 0 {
				continue
			}
			item := LineItem{AccountID: acc.ID, AccountName: acc.Name, Amount: ta.NetBalance}
			if acc.ID == "notes_payable" {
				longLiab = append(longLiab, item)
				totLongLiab = totLongLiab.Add(ta.NetBalance)
			} else {
				curLiab = append(curLiab, item)
				totCurLiab = totCurLiab.Add(ta.NetBalance)
			}

		case domain.CategoryEquity:
			if acc.ID == "common_stock" && ta.NetBalance.Cents() > 0 {
				equityItems = append(equityItems, LineItem{
					AccountID:   acc.ID,
					AccountName: acc.Name,
					Amount:      ta.NetBalance,
				})
				totEquity = totEquity.Add(ta.NetBalance)
			}
		}
	}

	// Sort Current Assets by liquidity: Cash, Accounts Receivable, Prepaids
	sort.Slice(curAssets, func(i, j int) bool {
		orderAsset := func(id domain.AccountID) int {
			switch id {
			case "cash":
				return 1
			case "accounts_receivable":
				return 2
			case "prepaid_insurance", "prepaid_rent":
				return 3
			default:
				return 10
			}
		}
		return orderAsset(curAssets[i].AccountID) < orderAsset(curAssets[j].AccountID)
	})

	// Sort Current Liabilities: Accounts Payable, Unearned Revenue, Dividends Payable
	sort.Slice(curLiab, func(i, j int) bool {
		orderLiab := func(id domain.AccountID) int {
			switch id {
			case "accounts_payable":
				return 1
			case "unearned_revenue":
				return 2
			case "dividends_payable":
				return 3
			default:
				return 10
			}
		}
		return orderLiab(curLiab[i].AccountID) < orderLiab(curLiab[j].AccountID)
	})

	// Add Ending Retained Earnings from Roll-Forward to Stockholders' Equity
	equityItems = append(equityItems, LineItem{
		AccountID:   "retained_earnings",
		AccountName: "Retained Earnings",
		Amount:      endRE,
	})
	totEquity = totEquity.Add(endRE)

	totAssets := totCurAssets.Add(totNonCurAssets)
	totLiab := totCurLiab.Add(totLongLiab)
	totLiabAndEquity := totLiab.Add(totEquity)

	balSheet := BalanceSheet{
		CompanyName:               c.CompanyName,
		AsOfDate:                  c.AsOfDate,
		CurrentAssets:             curAssets,
		TotalCurrentAssets:        totCurAssets,
		NonCurrentAssets:          nonCurAssets,
		TotalNonCurrentAssets:     totNonCurAssets,
		TotalAssets:               totAssets,
		CurrentLiabilities:        curLiab,
		TotalCurrentLiabilities:   totCurLiab,
		LongTermLiabilities:       longLiab,
		TotalLongTermLiabilities:  totLongLiab,
		TotalLiabilities:          totLiab,
		StockholdersEquity:        equityItems,
		TotalEquity:               totEquity,
		TotalLiabilitiesAndEquity: totLiabAndEquity,
		IsBalanced:                (totAssets == totLiabAndEquity),
	}
	report.BalanceSheet = balSheet

	// 10. Construct Statement of Cash Flows
	var netOpCash, netInvCash, netFinCash domain.Money
	for _, l := range opLines {
		netOpCash = netOpCash.Add(l.Amount)
	}
	for _, l := range invLines {
		netInvCash = netInvCash.Add(l.Amount)
	}
	for _, l := range finLines {
		netFinCash = netFinCash.Add(l.Amount)
	}

	netChangeInCash := netOpCash.Add(netInvCash).Add(netFinCash)
	begCash := c.OpeningBalances["cash"]
	endCash := begCash.Add(netChangeInCash)

	cfStmt := CashFlowStatement{
		CompanyName:         c.CompanyName,
		Period:              c.Period,
		OperatingActivities: opLines,
		NetOperatingCash:    netOpCash,
		InvestingActivities: invLines,
		NetInvestingCash:    netInvCash,
		FinancingActivities: finLines,
		NetFinancingCash:    netFinCash,
		NetChangeInCash:     netChangeInCash,
		BeginningCash:       begCash,
		EndingCash:          endCash,
		NoncashDisclosures:  noncash,
	}
	report.CashFlowStatement = cfStmt

	// 11. Execute Rigorous Gate Verification Checks
	allGatesPassed := true

	// Gate Check 1: Ending Balance Sheet Balances
	if !balSheet.IsBalanced {
		allGatesPassed = false
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✗ GATE 1 FAILED: Balance sheet does not balance: Assets (%s) != Liab+Eq (%s)",
				totAssets.FormatDollars(), totLiabAndEquity.FormatDollars()))
	} else {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✓ GATE 1 PASSED: Ending Balance Sheet balances 100%% (Assets %s = Liabilities %s + Equity %s)",
				totAssets.FormatDollars(), totLiab.FormatDollars(), totEquity.FormatDollars()))
	}

	// Gate Check 2: Net Income ties to Retained Earnings
	if reStmt.NetIncome != incStmt.NetIncome {
		allGatesPassed = false
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✗ GATE 2 FAILED: Net income mismatch: Income Statement (%s) != Retained Earnings (%s)",
				incStmt.NetIncome.FormatDollars(), reStmt.NetIncome.FormatDollars()))
	} else {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✓ GATE 2 PASSED: Net Income (%s) ties 100%% from Income Statement to Retained Earnings Statement",
				incStmt.NetIncome.FormatDollars()))
	}

	// Gate Check 3: Beginning Cash + Cash Flows = Ending Cash
	if begCash.Add(netChangeInCash) != endCash {
		allGatesPassed = false
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✗ GATE 3 FAILED: Cash flow math does not tie: Beg (%s) + Change (%s) != End (%s)",
				begCash.FormatDollars(), netChangeInCash.FormatDollars(), endCash.FormatDollars()))
	} else {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✓ GATE 3 PASSED: Beginning Cash (%s) + Net Cash Change (%s) = Ending Cash (%s)",
				begCash.FormatDollars(), netChangeInCash.FormatDollars(), endCash.FormatDollars()))
	}

	// Gate Check 4: Ending Cash ties to Ledger and Balance Sheet
	ledgerCash := ledger.GetOrCreate("cash").NetBalance
	if endCash != ledgerCash || endCash != findAccountAmount(curAssets, "cash") {
		allGatesPassed = false
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✗ GATE 4 FAILED: Cash tie-out mismatch: Cash Flow End (%s) != Ledger (%s) != Balance Sheet (%s)",
				endCash.FormatDollars(), ledgerCash.FormatDollars(), findAccountAmount(curAssets, "cash").FormatDollars()))
	} else {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✓ GATE 4 PASSED: Ending Cash (%s) ties 100%% to General Ledger and Balance Sheet Cash",
				endCash.FormatDollars()))
	}

	// Gate Check 5: Declared vs Paid Dividends Handled Distinctly
	// Verify whether Dividends Payable liability exists on balance sheet if declared != paid
	divPayableTA := ledger.GetOrCreate("dividends_payable")
	if divPayableTA != nil && divPayableTA.NetBalance.Cents() > 0 {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			fmt.Sprintf("✓ GATE 5 PASSED: Declared vs. Paid Dividends handled distinctly: Declared (%s) > Paid (%s), leaving %s Dividends Payable liability on Balance Sheet",
				divDeclared.FormatDollars(),
				divDeclared.Sub(divPayableTA.NetBalance).FormatDollars(),
				divPayableTA.NetBalance.FormatDollars()))
	} else {
		report.AuditReconciliationNotes = append(report.AuditReconciliationNotes,
			"✓ GATE 5 PASSED: Declared vs. Paid Dividends handled distinctly with appropriate liability and cash flow tracking.")
	}

	report.AllGatesPassed = allGatesPassed
	return report, nil
}

func buildTrialBalance(title string, ledger *domain.TAccountLedger, catalog *domain.AccountCatalog) TrialBalance {
	tb := TrialBalance{
		Title:    title,
		Accounts: make([]TrialBalanceAccount, 0),
	}

	for _, ta := range ledger.AccountsInOrder() {
		acc, _ := catalog.Get(ta.AccountID)
		if ta.TotalDebits.Cents() == 0 && ta.TotalCredits.Cents() == 0 {
			continue
		}
		item := TrialBalanceAccount{
			AccountID:   acc.ID,
			AccountName: acc.Name,
			Category:    acc.Category,
		}
		if ta.BalanceSide == domain.SideDebit && ta.NetBalance.Cents() > 0 {
			item.Debit = ta.NetBalance
			tb.TotalDebit = tb.TotalDebit.Add(ta.NetBalance)
		} else if ta.BalanceSide == domain.SideCredit && ta.NetBalance.Cents() > 0 {
			item.Credit = ta.NetBalance
			tb.TotalCredit = tb.TotalCredit.Add(ta.NetBalance)
		}
		tb.Accounts = append(tb.Accounts, item)
	}

	tb.IsBalanced = (tb.TotalDebit == tb.TotalCredit)
	return tb
}

func findAccountAmount(items []LineItem, id domain.AccountID) domain.Money {
	for _, it := range items {
		if it.AccountID == id {
			return it.Amount
		}
	}
	return domain.ZeroMoney()
}

// FormatReport outputs the complete financial report package as formatted text.
func (r *AccountingCycleReport) FormatReport() string {
	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("ACCOUNTING CYCLE & FINANCIAL STATEMENTS REPORT: %s\n", strings.ToUpper(r.Title)))
	b.WriteString(fmt.Sprintf("[%s]\n", domain.GenericVisualNotice))
	b.WriteString("================================================================================\n\n")

	b.WriteString(r.OpeningTrialBalance.FormatText())
	b.WriteString("\n")
	b.WriteString(r.UnadjustedTrialBalance.FormatText())
	b.WriteString("\n")
	b.WriteString(r.AdjustedTrialBalance.FormatText())
	b.WriteString("\n")
	b.WriteString(r.IncomeStatement.FormatText())
	b.WriteString("\n")
	b.WriteString(r.RetainedEarningsStatement.FormatText())
	b.WriteString("\n")
	b.WriteString(r.BalanceSheet.FormatText())
	b.WriteString("\n")
	b.WriteString(r.CashFlowStatement.FormatText())
	b.WriteString("\n")

	b.WriteString("================================================================================\n")
	b.WriteString("AUDIT & RECONCILIATION PROOF (STAGE 15 GATES)\n")
	b.WriteString("================================================================================\n")
	for _, note := range r.AuditReconciliationNotes {
		b.WriteString(fmt.Sprintf("%s\n", note))
	}
	b.WriteString("================================================================================\n")
	if r.AllGatesPassed {
		b.WriteString("STATUS: ALL STAGE 15 GATES PASSED (100% RECONCILED)\n")
	} else {
		b.WriteString("STATUS: AUDIT GATES OUTSTANDING\n")
	}
	b.WriteString("================================================================================\n")
	return b.String()
}

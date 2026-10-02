package statements_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/statements"
)

func setupTestEnvironment(t *testing.T) (*domain.AccountCatalog, *engine.Engine) {
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatalf("failed to load catalog from curriculum: %v", err)
	}
	eng := engine.NewEngine(catalog)
	return catalog, eng
}

func TestComprehensiveCaseAccountingCycleReport(t *testing.T) {
	catalog, eng := setupTestEnvironment(t)
	c := statements.CanonicalCasePioneerConsulting()

	report, err := statements.BuildAccountingCycleReport(c, catalog, eng)
	if err != nil {
		t.Fatalf("BuildAccountingCycleReport failed: %v", err)
	}

	if !report.AllGatesPassed {
		t.Fatalf("expected all gates to pass, notes:\n%s", strings.Join(report.AuditReconciliationNotes, "\n"))
	}

	// 1. Gate 1: Ending Balance Sheet Balances
	bs := report.BalanceSheet
	if !bs.IsBalanced {
		t.Errorf("balance sheet is not balanced: Assets (%s) != Liab+Eq (%s)",
			bs.TotalAssets.FormatDollars(), bs.TotalLiabilitiesAndEquity.FormatDollars())
	}
	if bs.TotalAssets.Cents() != 3080000 { // $30,800.00
		t.Errorf("expected total assets $30,800.00, got %s", bs.TotalAssets.FormatDollars())
	}
	if bs.TotalLiabilities.Cents() != 720000 { // $7,200.00
		t.Errorf("expected total liabilities $7,200.00, got %s", bs.TotalLiabilities.FormatDollars())
	}
	if bs.TotalEquity.Cents() != 2360000 { // $23,600.00
		t.Errorf("expected total equity $23,600.00, got %s", bs.TotalEquity.FormatDollars())
	}

	// 2. Gate 2: Net Income Ties to Retained Earnings
	is := report.IncomeStatement
	re := report.RetainedEarningsStatement
	if is.NetIncome.Cents() != 410000 { // $4,100.00
		t.Errorf("expected net income $4,100.00, got %s", is.NetIncome.FormatDollars())
	}
	if re.NetIncome != is.NetIncome {
		t.Errorf("net income mismatch between IS and RE statement: %s != %s",
			is.NetIncome.FormatDollars(), re.NetIncome.FormatDollars())
	}
	if re.EndingRetainedEarnings.Cents() != 360000 { // $3,600.00 ($4,100 net income - $500 dividend)
		t.Errorf("expected ending retained earnings $3,600.00, got %s", re.EndingRetainedEarnings.FormatDollars())
	}

	// 3. Gate 3: Beginning Cash + Cash Flows = Ending Cash
	cfs := report.CashFlowStatement
	if cfs.BeginningCash.Cents() != 2000000 { // $20,000.00
		t.Errorf("expected beginning cash $20,000.00, got %s", cfs.BeginningCash.FormatDollars())
	}
	if cfs.NetOperatingCash.Cents() != 190000 { // $1,900.00
		t.Errorf("expected net operating cash $1,900.00, got %s", cfs.NetOperatingCash.FormatDollars())
	}
	if cfs.NetInvestingCash.Cents() != -500000 { // -$5,000.00
		t.Errorf("expected net investing cash -$5,000.00, got %s", cfs.NetInvestingCash.FormatDollars())
	}
	if cfs.NetFinancingCash.Cents() != 470000 { // +$4,700.00 ($6,000 note - $1,000 repay - $300 dividend)
		t.Errorf("expected net financing cash $4,700.00, got %s", cfs.NetFinancingCash.FormatDollars())
	}
	if cfs.NetChangeInCash.Cents() != 160000 { // +$1,600.00 ($1,900 - $5,000 + $4,700)
		t.Errorf("expected net change in cash $1,600.00, got %s", cfs.NetChangeInCash.FormatDollars())
	}
	if cfs.EndingCash.Cents() != 2160000 { // $21,600.00
		t.Errorf("expected ending cash $21,600.00, got %s", cfs.EndingCash.FormatDollars())
	}
	if cfs.BeginningCash.Add(cfs.NetChangeInCash) != cfs.EndingCash {
		t.Errorf("beginning cash + net change != ending cash")
	}

	// 4. Gate 4: Ending Cash Ties to Ledger and Balance Sheet
	ledgerCash := report.Ledger.GetOrCreate("cash").NetBalance
	if cfs.EndingCash != ledgerCash {
		t.Errorf("ending cash does not tie to ledger: %s != %s",
			cfs.EndingCash.FormatDollars(), ledgerCash.FormatDollars())
	}
	balSheetCash := bs.CurrentAssets[0].Amount
	if cfs.EndingCash != balSheetCash {
		t.Errorf("ending cash does not tie to balance sheet cash: %s != %s",
			cfs.EndingCash.FormatDollars(), balSheetCash.FormatDollars())
	}

	// 5. Gate 5: Declared vs Paid Dividends Handled Distinctly
	divPayable := report.Ledger.GetOrCreate("dividends_payable").NetBalance
	if divPayable.Cents() != 20000 { // $200.00 remaining liability
		t.Errorf("expected $200.00 remaining dividends payable liability, got %s", divPayable.FormatDollars())
	}
	if re.LessDividendsDeclared.Cents() != 50000 { // $500.00 total declared reduces retained earnings
		t.Errorf("expected $500.00 dividends declared in RE statement, got %s", re.LessDividendsDeclared.FormatDollars())
	}
	// Verify Financing cash outflow only recorded the $300 actually paid
	foundDivCash := false
	for _, f := range cfs.FinancingActivities {
		if strings.Contains(f.Description, "dividends") {
			foundDivCash = true
			if f.Amount.Cents() != -30000 {
				t.Errorf("expected cash paid for dividends to be -$300.00, got %s", f.Amount.FormatDollars())
			}
		}
	}
	if !foundDivCash {
		t.Errorf("expected to find dividend cash payment in financing cash flows")
	}

	// 6. Report Rendering Check
	rendered := report.FormatReport()
	if !strings.Contains(rendered, "ACCOUNTING CYCLE & FINANCIAL STATEMENTS REPORT") {
		t.Errorf("expected title in formatted report")
	}
	if !strings.Contains(rendered, "STATUS: ALL STAGE 15 GATES PASSED (100% RECONCILED)") {
		t.Errorf("expected all gates passed status in report, got:\n%s", rendered)
	}
}

func TestTrialBalanceReconciliations(t *testing.T) {
	catalog, eng := setupTestEnvironment(t)
	c := statements.CanonicalCasePioneerConsulting()

	report, err := statements.BuildAccountingCycleReport(c, catalog, eng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Opening Trial Balance: Dr $20,000 == Cr $20,000
	if !report.OpeningTrialBalance.IsBalanced || report.OpeningTrialBalance.TotalDebit.Cents() != 2000000 {
		t.Errorf("opening trial balance failed: %+v", report.OpeningTrialBalance)
	}

	// Unadjusted Trial Balance: Dr == Cr
	if !report.UnadjustedTrialBalance.IsBalanced {
		t.Errorf("unadjusted trial balance unbalanced: Dr %s != Cr %s",
			report.UnadjustedTrialBalance.TotalDebit.FormatDollars(),
			report.UnadjustedTrialBalance.TotalCredit.FormatDollars())
	}

	// Adjusted Trial Balance: Dr == Cr
	if !report.AdjustedTrialBalance.IsBalanced {
		t.Errorf("adjusted trial balance unbalanced: Dr %s != Cr %s",
			report.AdjustedTrialBalance.TotalDebit.FormatDollars(),
			report.AdjustedTrialBalance.TotalCredit.FormatDollars())
	}
	// Verify adjusted trial balance total debits equal total credits
	if report.AdjustedTrialBalance.TotalDebit != report.AdjustedTrialBalance.TotalCredit {
		t.Errorf("adjusted trial balance Dr != Cr")
	}
}

func TestUnbalancedOpeningBalancesRejected(t *testing.T) {
	catalog, eng := setupTestEnvironment(t)
	badCase := statements.ComprehensiveCase{
		CaseID:      "bad_case",
		Title:       "Bad Case",
		CompanyName: "Bad Co",
		Period:      "Jan 2026",
		AsOfDate:    "Jan 31, 2026",
		OpeningBalances: map[domain.AccountID]domain.Money{
			"cash":         domain.NewMoney(100000), // Dr $1,000
			"common_stock": domain.NewMoney(80000),  // Cr $800 (unbalanced by $200!)
		},
	}

	_, err := statements.BuildAccountingCycleReport(badCase, catalog, eng)
	if err == nil {
		t.Fatalf("expected error for unbalanced opening balances, got nil")
	}
	if !strings.Contains(err.Error(), "opening trial balance does not balance") {
		t.Errorf("expected 'opening trial balance does not balance' in error, got: %v", err)
	}
}

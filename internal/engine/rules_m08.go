package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// RegisterM08Rules registers the Stream-M08 rules into the Engine.
func RegisterM08Rules(e *Engine) {
	e.rules[FamilyM08LeaseCommencementROUAsset] = M08LeaseCommencementROUAssetRule{}
	e.rules[FamilyM08LeasePeriodicPayment] = M08LeasePeriodicPaymentRule{}
	e.rules[FamilyM08TreasuryStockRepurchase] = M08TreasuryStockRepurchaseRule{}
	e.rules[FamilyM08TreasuryStockReissuance] = M08TreasuryStockReissuanceRule{}
}

// 1. M08LeaseCommencementROUAssetRule: Initial lease commencement under ASC 842 at PV of lease payments.
// Debits Right-of-Use Asset (Asset increases),
// Credits Lease Liability (Liability increases).
type M08LeaseCommencementROUAssetRule struct{}

func (M08LeaseCommencementROUAssetRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "right_of_use_asset", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "lease_liability", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M08LeaseCommencementROUAssetRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded initial lease commencement under ASC 842: recognized %s Right-of-Use Asset and corresponding Lease Liability.", amt.FormatDollars()),
		DebitRationale:  "Right-of-Use Asset is an Asset that increased; under ASC 842, the lessee capitalizes the right to control and use the leased asset over the contract term by Debit (+A).",
		CreditRationale: "Lease Liability is a Liability that increased; under ASC 842, the present value of future lease payments is recognized as a financing liability by Credit (+L).",
		EquationSummary: "Assets increase (+Right-of-Use Asset); Liabilities increase (+Lease Liability); Equity is unchanged. Under ASC 842, both operating and finance leases are capitalized on the balance sheet, ending off-balance sheet operating lease treatment.",
	}
}

func (M08LeaseCommencementROUAssetRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "rent_expense" {
			return []string{TagM08LeaseExpensedAtCommencement}, "Under ASC 842, leases exceeding 12 months must be capitalized onto the balance sheet as Right-of-Use Assets at commencement rather than expensed immediately."
		}
		if p.AccountID == "common_stock" || p.AccountID == "retained_earnings" {
			return []string{TagM08ROUAssetRecordedAsEquity}, "Lease commitments represent contractual obligations owed to the lessor, not stockholders' equity."
		}
		if p.AccountID == "right_of_use_asset" && p.Side == domain.SideCredit {
			return []string{TagReversedSides}, "The Right-of-Use Asset represents a newly acquired economic resource and must be debited (+A) at commencement."
		}
	}
	return nil, ""
}

// 2. M08LeasePeriodicPaymentRule: Periodic lease payment reducing lease liability and recording interest expense.
// Debits Lease Liability (Liability decreases by principal portion),
// Debits Interest Expense (Expense increases, reducing Equity),
// Credits Cash (Asset decreases by total cash payment).
type M08LeasePeriodicPaymentRule struct{}

func (M08LeasePeriodicPaymentRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	totalAmt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	interestVal := params["interest_minor_units"]
	if interestVal <= 0 {
		interestVal = totalAmt.Cents() * 20 / 100
	}
	principalVal := params["principal_minor_units"]
	if principalVal <= 0 {
		principalVal = totalAmt.Cents() - interestVal
	}
	interestAmt := domain.NewMoney(interestVal)
	principalAmt := domain.NewMoney(principalVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "lease_liability", Side: domain.SideDebit, Amount: principalAmt},
		domain.Posting{AccountID: "interest_expense", Side: domain.SideDebit, Amount: interestAmt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: totalAmt},
	), nil
}

func (M08LeasePeriodicPaymentRule) Explain(params map[string]int64) EventExplanation {
	totalAmt, _ := getAmount(params)
	interestVal := params["interest_minor_units"]
	if interestVal <= 0 {
		interestVal = totalAmt.Cents() * 20 / 100
	}
	principalVal := params["principal_minor_units"]
	if principalVal <= 0 {
		principalVal = totalAmt.Cents() - interestVal
	}
	interestAmt := domain.NewMoney(interestVal)
	principalAmt := domain.NewMoney(principalVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Paid %s cash for periodic lease payment, reducing Lease Liability principal by %s and recognizing %s of Interest Expense.", totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars()),
		DebitRationale:  "Lease Liability is a Liability that decreased by Debit (-L) for principal repaid; Interest Expense is an Expense that increased by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Cash is an Asset that decreased by Credit (-A) for the total lease payment disbursed.",
		EquationSummary: fmt.Sprintf("Assets decrease by %s (-Cash); Liabilities decrease by %s (-Lease Liability); Equity decreases by %s (-Interest Expense). Total debits equal total credits.", totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars()),
	}
}

func (M08LeasePeriodicPaymentRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	hasInterest := false
	for _, p := range candidate.Postings {
		if p.AccountID == "interest_expense" {
			hasInterest = true
		}
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagReversedSides}, "Cash was paid to the lessor, so Cash must be credited (-A)."
		}
	}
	if !hasInterest {
		return []string{TagM08IgnoringLeaseInterest}, "Periodic lease payments amortize both the outstanding principal liability and the financing interest cost accrued over the period."
	}
	return nil, ""
}

// 3. M08TreasuryStockRepurchaseRule: Share repurchase under cost method.
// Debits Treasury Stock (Contra-Equity increases, reducing total stockholders' equity),
// Credits Cash (Asset decreases).
type M08TreasuryStockRepurchaseRule struct{}

func (M08TreasuryStockRepurchaseRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "treasury_stock", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M08TreasuryStockRepurchaseRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Repurchased company shares for %s cash under the cost method to hold in treasury.", amt.FormatDollars()),
		DebitRationale:  "Treasury Stock is a contra-equity account that increased; contra-equity accounts increase by Debit and reduce total Stockholders' Equity (-E).",
		CreditRationale: "Cash is an Asset that decreased; Assets decrease by Credit (-A).",
		EquationSummary: "Assets decrease (-Cash); Equity decreases (-Treasury Stock); Liabilities are unchanged. Under the cost method, repurchased shares are recorded at acquisition cost as a deduction from total stockholders' equity.",
	}
}

func (M08TreasuryStockRepurchaseRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "equipment" || p.AccountID == "goodwill" || p.AccountID == "patents" || p.AccountID == "digital_assets" {
			return []string{TagM08TreasuryStockRecordedAsAsset}, "Treasury stock represents repurchased company equity shares (contra-equity), never an investment asset or operating equipment."
		}
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagReversedSides}, "Cash was disbursed to repurchase shares from open market stockholders, not received."
		}
	}
	return nil, ""
}

// 4. M08TreasuryStockReissuanceRule: Reissuing treasury shares above cost.
// Debits Cash (Asset increases by total sales proceeds),
// Credits Treasury Stock (Contra-Equity decreases by original cost basis, increasing Equity),
// Credits Additional Paid-in Capital (Equity increases for the excess of proceeds over cost).
type M08TreasuryStockReissuanceRule struct{}

func (M08TreasuryStockReissuanceRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	totalAmt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	costVal := params["cost_minor_units"]
	if costVal <= 0 {
		costVal = totalAmt.Cents() * 2 / 3
	}
	apicVal := params["apic_minor_units"]
	if apicVal <= 0 {
		apicVal = totalAmt.Cents() - costVal
	}
	costAmt := domain.NewMoney(costVal)
	apicAmt := domain.NewMoney(apicVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: totalAmt},
		domain.Posting{AccountID: "treasury_stock", Side: domain.SideCredit, Amount: costAmt},
		domain.Posting{AccountID: "additional_paid_in_capital", Side: domain.SideCredit, Amount: apicAmt},
	), nil
}

func (M08TreasuryStockReissuanceRule) Explain(params map[string]int64) EventExplanation {
	totalAmt, _ := getAmount(params)
	costVal := params["cost_minor_units"]
	if costVal <= 0 {
		costVal = totalAmt.Cents() * 2 / 3
	}
	apicVal := params["apic_minor_units"]
	if apicVal <= 0 {
		apicVal = totalAmt.Cents() - costVal
	}
	costAmt := domain.NewMoney(costVal)
	apicAmt := domain.NewMoney(apicVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Reissued treasury shares for %s cash, removing %s cost basis of treasury stock and crediting %s to Additional Paid-in Capital.", totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increased by Debit (+A) for the total cash proceeds received from reissuing shares.",
		CreditRationale: "Treasury Stock is a contra-equity account that decreased by Credit at historical cost (+E); Additional Paid-in Capital is an Equity account credited (+E) for the excess of proceeds over cost.",
		EquationSummary: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (Treasury Stock contra-equity decreases by %s and Additional Paid-in Capital increases by %s); Liabilities are unchanged. Under US GAAP, transactions in a company's own stock cannot generate income statement gains or losses.", totalAmt.FormatDollars(), totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars()),
	}
}

func (M08TreasuryStockReissuanceRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "gain_on_disposal" || p.AccountID == "service_revenue" {
			return []string{TagM08GainRecordedOnTreasuryReissuance}, "Under US GAAP, a corporation can never recognize an income statement gain or revenue from transactions in its own equity shares; the excess over cost is credited to Additional Paid-in Capital."
		}
		if p.AccountID == "retained_earnings" && p.Side == domain.SideDebit {
			return []string{TagM08RetainedEarningsDebitedOnReissuance}, "When treasury shares are reissued above cost, the excess increases Additional Paid-in Capital, not Retained Earnings."
		}
		if p.AccountID == "cash" && p.Side == domain.SideCredit {
			return []string{TagReversedSides}, "Cash was received from reissuing treasury shares to investors and must be debited (+A)."
		}
	}
	return nil, ""
}

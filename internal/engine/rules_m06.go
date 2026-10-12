package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M06PPEStraightLineDeprRule: Straight-line depreciation adjusting entry ((Cost - Salvage) / Useful Life).
// Debits Depreciation Expense (Expense increases, reducing Equity),
// Credits Accumulated Depreciation (Contra-Asset increases, reducing net Assets).
type M06PPEStraightLineDeprRule struct{}

func (M06PPEStraightLineDeprRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "depreciation_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accumulated_depreciation", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M06PPEStraightLineDeprRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded annual straight-line depreciation adjusting entry of %s.", amt.FormatDollars()),
		DebitRationale:  "Depreciation Expense is an operating Expense that increased; Expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Accumulated Depreciation is a Contra-Asset account that increased by Credit, reducing the net carrying book value of Equipment (-A).",
		EquationSummary: "Assets decrease (-Accumulated Depreciation); Equity decreases (-Depreciation Expense); Liabilities are unchanged. Under the straight-line method, annual depreciation expense equals (Historical Cost - Estimated Salvage Value) divided by Useful Life.",
	}
}

func (M06PPEStraightLineDeprRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "equipment" && p.Side == domain.SideCredit {
			return []string{TagM06DirectAssetCredit}, "Under US GAAP, depreciation is credited to the contra-asset Accumulated Depreciation rather than directly reducing the Equipment cost account."
		}
		if p.AccountID == "rent_expense" || p.AccountID == "insurance_expense" {
			return []string{TagWrongAccount}, "Periodic equipment cost allocation must be recorded as Depreciation Expense, not general operating rent or insurance."
		}
	}
	return nil, ""
}

// 2. M06PPEMACRSTaxTimingRule: Accelerated MACRS / double-declining timing difference vs. straight-line book depreciation.
// For financial reporting, the company records straight-line book depreciation.
// Debits Depreciation Expense (Expense increases, reducing Equity),
// Credits Accumulated Depreciation (Contra-Asset increases, reducing net Assets).
type M06PPEMACRSTaxTimingRule struct{}

func (M06PPEMACRSTaxTimingRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "depreciation_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accumulated_depreciation", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M06PPEMACRSTaxTimingRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded %s of straight-line book depreciation expense for financial reporting under US GAAP.", amt.FormatDollars()),
		DebitRationale:  "Depreciation Expense is an operating Expense that increased; Expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Accumulated Depreciation is a Contra-Asset that increased by Credit, reducing net Equipment carrying value (-A).",
		EquationSummary: "Assets decrease (-Accumulated Depreciation); Equity decreases (-Depreciation Expense); Liabilities are unchanged. While tax returns use accelerated MACRS for rapid tax deductions, financial statements use straight-line depreciation to reflect economic consumption, creating a temporary timing difference.",
	}
}

func (M06PPEMACRSTaxTimingRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "equipment" && p.Side == domain.SideCredit {
			return []string{TagM06DirectAssetCredit}, "Depreciation should be credited to Accumulated Depreciation, maintaining historical equipment cost on the ledger."
		}
	}
	return nil, ""
}

// 3. M06PPEAssetDisposalGainLossRule: Derecognition of PPE upon sale for cash, removing Equipment and Accumulated Depreciation, recognizing Gain on Disposal.
// Debits Cash (Asset increases),
// Debits Accumulated Depreciation (Contra-Asset removed),
// Credits Equipment (Asset removed at historical cost),
// Credits Gain on Disposal (Other Income increases, increasing Equity).
type M06PPEAssetDisposalGainLossRule struct{}

func (M06PPEAssetDisposalGainLossRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}

	accumVal, ok := params["accum_depr_minor_units"]
	if !ok || accumVal <= 0 {
		accumVal = amt.Cents() * 80 / 100
	}
	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() * 25 / 100
	}
	gainVal, ok := params["gain_minor_units"]
	if !ok || gainVal <= 0 {
		gainVal = cashVal - (amt.Cents() - accumVal)
	}

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(cashVal)},
		domain.Posting{AccountID: "accumulated_depreciation", Side: domain.SideDebit, Amount: domain.NewMoney(accumVal)},
		domain.Posting{AccountID: "equipment", Side: domain.SideCredit, Amount: amt},
		domain.Posting{AccountID: "gain_on_disposal", Side: domain.SideCredit, Amount: domain.NewMoney(gainVal)},
	), nil
}

func (M06PPEAssetDisposalGainLossRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	accumVal := params["accum_depr_minor_units"]
	if accumVal <= 0 {
		accumVal = amt.Cents() * 80 / 100
	}
	cashVal := params["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() * 25 / 100
	}
	gainVal := params["gain_minor_units"]
	if gainVal <= 0 {
		gainVal = cashVal - (amt.Cents() - accumVal)
	}

	cashAmt := domain.NewMoney(cashVal)
	accumAmt := domain.NewMoney(accumVal)
	gainAmt := domain.NewMoney(gainVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Derecognized equipment upon sale: received %s cash, eliminated %s accumulated depreciation and %s original cost, recognizing %s Gain on Disposal.", cashAmt.FormatDollars(), accumAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Cash increases by %s (+A); Accumulated Depreciation is debited for %s (+A net) to eliminate the contra-asset balance of the disposed equipment.", cashAmt.FormatDollars(), accumAmt.FormatDollars()),
		CreditRationale: fmt.Sprintf("Equipment is credited for %s (-A) to remove original cost; Gain on Disposal increases by %s (+R/+E) for cash proceeds exceeding book value.", amt.FormatDollars(), gainAmt.FormatDollars()),
		EquationSummary: fmt.Sprintf("Net Assets increase by %s (+Cash %s, +Accumulated Depreciation removal %s, -Equipment %s); Equity increases by %s (+Gain on Disposal); Liabilities are unchanged.", gainAmt.FormatDollars(), cashAmt.FormatDollars(), accumAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars()),
	}
}

func (M06PPEAssetDisposalGainLossRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" {
			return []string{TagM06GainTreatedAsRevenue}, "Proceeds from disposing of capital equipment generate a non-operating Gain on Disposal, not core operating Service Revenue."
		}
	}
	return nil, ""
}

// 4. M06PPEImpairmentRecoverabilityRule: ASC 360 two-step impairment test (undiscounted cash flows < carrying value triggers write-down to fair value via Loss on Impairment).
// Debits Loss on Impairment (Expense increases, reducing Equity),
// Credits Accumulated Depreciation (Contra-Asset increases, reducing net Assets).
type M06PPEImpairmentRecoverabilityRule struct{}

func (M06PPEImpairmentRecoverabilityRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "loss_on_impairment", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "accumulated_depreciation", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M06PPEImpairmentRecoverabilityRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded %s equipment impairment write-down under ASC 360 following recoverability test failure.", amt.FormatDollars()),
		DebitRationale:  "Loss on Impairment is recognized in earnings as an operating expense/loss; Expenses increase by Debit, reducing Stockholders' Equity (-E).",
		CreditRationale: "Accumulated Depreciation is credited to reduce the net carrying value of equipment to its fair value (-A).",
		EquationSummary: "Assets decrease (-Accumulated Depreciation); Equity decreases (-Loss on Impairment); Liabilities are unchanged. Under ASC 360, when undiscounted future cash flows fall below carrying value (Step 1), an impairment loss is measured as Carrying Value minus Fair Value (Step 2); write-downs cannot be reversed under US GAAP.",
	}
}

func (M06PPEImpairmentRecoverabilityRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" {
			return []string{TagCashRecordedWhenUncollected}, "Impairment is a non-cash valuation write-down; no cash is exchanged."
		}
	}
	return nil, ""
}

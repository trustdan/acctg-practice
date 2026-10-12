package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M07BondIssuanceDiscountRule: Issuing bonds at a discount when market rate > stated coupon rate.
// Debits Cash (Asset increases by net proceeds),
// Debits Discount on Bonds Payable (Contra-Liability increases, normal debit balance),
// Credits Bonds Payable (Liability increases by face value).
type M07BondIssuanceDiscountRule struct{}

func (M07BondIssuanceDiscountRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	discountVal, ok := params["discount_minor_units"]
	if !ok || discountVal <= 0 {
		discountVal = amt.Cents() * 5 / 100
	}
	discount := domain.NewMoney(discountVal)

	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cash := domain.NewMoney(cashVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: cash},
		domain.Posting{AccountID: "discount_on_bonds_payable", Side: domain.SideDebit, Amount: discount},
		domain.Posting{AccountID: "bonds_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M07BondIssuanceDiscountRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	discountVal := params["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 5 / 100
	}
	discount := domain.NewMoney(discountVal)
	cash := domain.NewMoney(amt.Cents() - discountVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Issued %s face value bonds at a discount, receiving %s cash and recording %s Discount on Bonds Payable.", amt.FormatDollars(), cash.FormatDollars(), discount.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Cash increases by %s (+A); Discount on Bonds Payable is a Contra-Liability that increases by Debit (%s) to offset bonds to net carrying value.", cash.FormatDollars(), discount.FormatDollars()),
		CreditRationale: fmt.Sprintf("Bonds Payable is a Liability that increases by %s (+L) for the contractual maturity obligation.", amt.FormatDollars()),
		EquationSummary: fmt.Sprintf("Assets increase by %s (+Cash); Net Liabilities increase by %s (+Bonds Payable %s minus Discount %s); Equity is unchanged.", cash.FormatDollars(), cash.FormatDollars(), amt.FormatDollars(), discount.FormatDollars()),
	}
}

func (M07BondIssuanceDiscountRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	hasDiscount := false
	for _, p := range candidate.Postings {
		if p.AccountID == "discount_on_bonds_payable" {
			hasDiscount = true
		}
		if p.AccountID == "service_revenue" || p.AccountID == "interest_expense" {
			return []string{TagM07DiscountRecordedAsLoss}, "A bond discount is not an immediate financing loss or expense; it is a contra-liability amortized over the bond term."
		}
	}
	if !hasDiscount {
		return []string{TagM07DiscountRecordedAsLoss}, "Issuing bonds below par requires recording a debit to Discount on Bonds Payable (contra-liability)."
	}
	return nil, ""
}

// 2. M07BondIssuancePremiumRule: Issuing bonds at a premium when market rate < stated coupon rate.
// Debits Cash (Asset increases by gross proceeds),
// Credits Premium on Bonds Payable (Adjunct Liability increases, normal credit balance),
// Credits Bonds Payable (Liability increases by face value).
type M07BondIssuancePremiumRule struct{}

func (M07BondIssuancePremiumRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	premiumVal, ok := params["premium_minor_units"]
	if !ok || premiumVal <= 0 {
		premiumVal = amt.Cents() * 5 / 100
	}
	premium := domain.NewMoney(premiumVal)

	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() + premiumVal
	}
	cash := domain.NewMoney(cashVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: cash},
		domain.Posting{AccountID: "premium_on_bonds_payable", Side: domain.SideCredit, Amount: premium},
		domain.Posting{AccountID: "bonds_payable", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M07BondIssuancePremiumRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	premiumVal := params["premium_minor_units"]
	if premiumVal <= 0 {
		premiumVal = amt.Cents() * 5 / 100
	}
	premium := domain.NewMoney(premiumVal)
	cash := domain.NewMoney(amt.Cents() + premiumVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Issued %s face value bonds at a premium, receiving %s cash and recording %s Premium on Bonds Payable.", amt.FormatDollars(), cash.FormatDollars(), premium.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Cash increases by %s (+A) for the total proceeds received from bondholders.", cash.FormatDollars()),
		CreditRationale: fmt.Sprintf("Bonds Payable increases by %s (+L) for the face maturity obligation; Premium on Bonds Payable is an Adjunct Liability that increases by %s (+L) to reflect the carrying value above par.", amt.FormatDollars(), premium.FormatDollars()),
		EquationSummary: fmt.Sprintf("Assets increase by %s (+Cash); Total Liabilities increase by %s (+Bonds Payable %s plus +Premium %s); Stockholders' Equity is completely unaffected by debt issuance.", cash.FormatDollars(), cash.FormatDollars(), amt.FormatDollars(), premium.FormatDollars()),
	}
}

func (M07BondIssuancePremiumRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" || p.AccountID == "gain_on_disposal" {
			return []string{TagM07PremiumRecordedAsGain}, "Bond premium is an adjunct liability added to debt on the balance sheet, not an immediate gain or revenue."
		}
	}
	return nil, ""
}

// 3. M07BondEffectiveInterestAmortRule: Semiannual interest payment and discount amortization under effective interest method.
// Debits Interest Expense (Carrying Value * market rate / 2),
// Credits Cash (Face Value * coupon rate / 2),
// Credits Discount on Bonds Payable (difference amortizing discount toward par).
type M07BondEffectiveInterestAmortRule struct{}

func (M07BondEffectiveInterestAmortRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params) // Interest Expense
	if err != nil {
		return domain.Entry{}, err
	}
	discountVal, ok := params["discount_minor_units"]
	if !ok || discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	discount := domain.NewMoney(discountVal)

	cashVal, ok := params["cash_minor_units"]
	if !ok || cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cash := domain.NewMoney(cashVal)

	return domain.NewEntry(
		domain.Posting{AccountID: "interest_expense", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: cash},
		domain.Posting{AccountID: "discount_on_bonds_payable", Side: domain.SideCredit, Amount: discount},
	), nil
}

func (M07BondEffectiveInterestAmortRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	discountVal := params["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	discount := domain.NewMoney(discountVal)
	cash := domain.NewMoney(amt.Cents() - discountVal)

	return EventExplanation{
		Summary:         fmt.Sprintf("Recorded semiannual bond interest under the effective interest method: paid %s cash coupon, amortized %s discount, and recognized %s Interest Expense.", cash.FormatDollars(), discount.FormatDollars(), amt.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Interest Expense is an Expense that increases by Debit (%s), reflecting the true economic borrowing cost (carrying value times market interest rate).", amt.FormatDollars()),
		CreditRationale: fmt.Sprintf("Cash decreases by Credit (%s) for the contractual coupon paid; Discount on Bonds Payable (contra-liability) is credited (%s) as the discount amortizes toward zero.", cash.FormatDollars(), discount.FormatDollars()),
		EquationSummary: fmt.Sprintf("Assets decrease by %s (-Cash); Net Liabilities increase by %s (+carrying value as contra-liability discount shrinks); Equity decreases by %s (-Interest Expense). Equation balances: -%s = +%s + (-%s).", cash.FormatDollars(), discount.FormatDollars(), amt.FormatDollars(), cash.FormatDollars(), discount.FormatDollars(), amt.FormatDollars()),
	}
}

func (M07BondEffectiveInterestAmortRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	hasDiscount := false
	for _, p := range candidate.Postings {
		if p.AccountID == "discount_on_bonds_payable" {
			hasDiscount = true
		}
	}
	if !hasDiscount {
		return []string{TagM07IgnoringBondAmortization}, "Under the effective interest method, interest expense exceeds the cash coupon paid by the amount of discount amortization, which must be credited to Discount on Bonds Payable."
	}
	return nil, ""
}

// 4. M07BondContingentLiabilityRule: ASC 450 contingent liability accrual (probable and reasonably estimable litigation loss).
// Debits Litigation Loss (Expense / Loss increases, reducing Equity),
// Credits Contingent Liability (Liability increases).
type M07BondContingentLiabilityRule struct{}

func (M07BondContingentLiabilityRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "litigation_loss", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "contingent_liability", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M07BondContingentLiabilityRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Accrued a %s contingent litigation loss and liability under ASC 450 because the loss is probable and reasonably estimable.", amt.FormatDollars()),
		DebitRationale:  fmt.Sprintf("Litigation Loss is an Expense/Loss that increases by Debit (%s), reducing current period earnings and Stockholders' Equity (-E).", amt.FormatDollars()),
		CreditRationale: fmt.Sprintf("Contingent Liability is a Liability that increases by Credit (%s), establishing the estimated legal claim obligation.", amt.FormatDollars()),
		EquationSummary: fmt.Sprintf("Assets are unchanged ($0); Liabilities increase by %s (+Contingent Liability); Equity decreases by %s (-Litigation Loss). Fundamental equation balances identically: 0 = +%s + (-%s).", amt.FormatDollars(), amt.FormatDollars(), amt.FormatDollars(), amt.FormatDollars()),
	}
}

func (M07BondContingentLiabilityRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "accounts_payable" {
			return []string{TagM07ContingentRecordedAsPayable}, "A pending litigation contingency is a legal loss contingency recorded in a dedicated Contingent Liability account, not routine supplier Accounts Payable."
		}
		if p.AccountID == "cash" && p.Side == domain.SideCredit {
			return []string{TagCashRecordedWhenUncollected}, "The lawsuit has not yet been paid or settled; accrual under ASC 450 establishes a liability prior to cash payout."
		}
	}
	return nil, ""
}

package engine

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// 1. M04VariableConsiderationRule: Fixed fee billed on account while uncertain bonus is constrained under ASC 606 Step 3.
type M04VariableConsiderationRule struct{}

func (M04VariableConsiderationRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M04VariableConsiderationRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Completed services and billed client %s on account, properly constraining uncertain variable consideration under ASC 606 Step 3.", amt.FormatDollars()),
		DebitRationale:  "Accounts Receivable is an Asset that increases because the firm holds an unconditional right to the fixed fee; Assets increase by Debit.",
		CreditRationale: "Service Revenue increases Equity for the unconstrained transaction price earned; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Accounts Receivable); Equity increases (+Service Revenue).",
	}
}

func (M04VariableConsiderationRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagCashRecordedWhenUncollected}, "The services were billed on credit terms; cash has not yet been collected from the client."
		}
	}
	return nil, ""
}

// 2. M04ContractAssetARRule: Unconditional right to consideration (Accounts Receivable) vs conditional right (Contract Asset) under ASC 606-10-45.
type M04ContractAssetARRule struct{}

func (M04ContractAssetARRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "accounts_receivable", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M04ContractAssetARRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Satisfied performance milestone with unconditional legal right to payment under ASC 606-10-45; billed %s on account.", amt.FormatDollars()),
		DebitRationale:  "Accounts Receivable is an Asset that increases because the right to payment is unconditional (dependent only on passage of time); Assets increase by Debit.",
		CreditRationale: "Service Revenue increases Equity as the performance obligation has been fulfilled; Revenue increases by Credit.",
		EquationSummary: "Assets increase (+Accounts Receivable); Equity increases (+Service Revenue).",
	}
}

func (M04ContractAssetARRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagCashRecordedWhenUncollected}, "Payment is on net 30 credit terms; cash has not yet been received."
		}
	}
	return nil, ""
}

// 3. M04RightOfReturnRule: Cash advance subject to right of return / refund liability under ASC 606-10-55-23 recorded as Unearned Revenue.
type M04RightOfReturnRule struct{}

func (M04RightOfReturnRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M04RightOfReturnRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Received %s cash advance subject to full money-back return/refund terms prior to performance under ASC 606; recorded as contract liability (Unearned Revenue).", amt.FormatDollars()),
		DebitRationale:  "Cash is an Asset that increases upon collecting funds; Assets increase by Debit.",
		CreditRationale: "Unearned Revenue is a Liability that increases for the obligation to either deliver satisfactory service or refund customer consideration; Liabilities increase by Credit.",
		EquationSummary: "Assets increase (+Cash); Liabilities increase (+Unearned Revenue).",
	}
}

func (M04RightOfReturnRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "service_revenue" && p.Side == domain.SideCredit {
			return []string{TagM04RevenueRecognizedWithUncertainReturn, TagRevenueRecognizedPrematurely}, "Revenue cannot be recognized before services are performed and while full refund rights remain pending."
		}
	}
	return nil, ""
}

// 4. M04SatisfactionOverTimeRule: Proportional satisfaction of performance obligation over time under ASC 606 Step 5, reducing contract liability.
type M04SatisfactionOverTimeRule struct{}

func (M04SatisfactionOverTimeRule) DeriveEntry(params map[string]int64, catalog *domain.AccountCatalog) (domain.Entry, error) {
	amt, err := getAmount(params)
	if err != nil {
		return domain.Entry{}, err
	}
	return domain.NewEntry(
		domain.Posting{AccountID: "unearned_revenue", Side: domain.SideDebit, Amount: amt},
		domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: amt},
	), nil
}

func (M04SatisfactionOverTimeRule) Explain(params map[string]int64) EventExplanation {
	amt, _ := getAmount(params)
	return EventExplanation{
		Summary:         fmt.Sprintf("Delivered %s of contracted services over time under ASC 606 Step 5; relieved contract liability (Unearned Revenue) and recognized revenue.", amt.FormatDollars()),
		DebitRationale:  "Unearned Revenue is a Liability that decreases as performance obligations are fulfilled over time; Liabilities decrease by Debit.",
		CreditRationale: "Service Revenue increases Equity as revenue is earned; Revenue increases by Credit.",
		EquationSummary: "Liabilities decrease (-Unearned Revenue); Equity increases (+Service Revenue).",
	}
}

func (M04SatisfactionOverTimeRule) DiagnoseError(candidate domain.Entry, params map[string]int64) ([]string, string) {
	for _, p := range candidate.Postings {
		if p.AccountID == "cash" && p.Side == domain.SideDebit {
			return []string{TagM04CashDebitedOnOverTimeSatisfaction, TagCashRecordedOnEarningAdvance}, "No new cash was received; the client previously paid in advance, so the existing contract liability must be relieved."
		}
	}
	return nil, ""
}

package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM07Stages dispatches to the appropriate modular stage builder for Stream-M07.
func BuildM07Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM07BondIssuanceDiscount:
		BuildM07BondIssuanceDiscountStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM07BondIssuancePremium:
		BuildM07BondIssuancePremiumStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM07BondEffectiveInterestAmort:
		BuildM07BondEffectiveInterestAmortStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM07BondContingentLiability:
		BuildM07BondContingentLiabilityStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM07BondIssuanceDiscountStages builds 7 scaffolding stages for issuing bonds at a discount.
func BuildM07BondIssuanceDiscountStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	discountVal := proc.Event.Parameters["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 5 / 100
	}
	cashVal := proc.Event.Parameters["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cashMoney := domain.NewMoney(cashVal)
	discountMoney := domain.NewMoney(discountVal)
	faceMoney := amt

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_discount_bonds", Text: "Discount on Bonds Payable"},
		{ID: "opt_loss_issuance", Text: "Loss on Bond Issuance", ErrorTag: engine.TagM07DiscountRecordedAsLoss},
		{ID: "opt_interest_expense", Text: "Interest Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_premium_bonds", Text: "Premium on Bonds Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 7011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When bonds are issued below par value because the market interest rate exceeds the coupon rate, which contra-liability account captures the difference between face value and cash received?",
		CorrectOptionID:   "opt_discount_bonds",
		Options:           stage1Opts,
		CausalHint:        "Which contra-liability account offsets the face value of the debt on the balance sheet to reflect the lower initial borrowing proceeds?",
		Explanation:       "Discount on Bonds Payable is a contra-liability account that records the difference when bonds are issued for less than their maturity face value.",
		RelevantConceptID: "discount_on_bonds_payable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_contra_liability", Text: "Liability (Contra-Liability)"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 7012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Discount on Bonds Payable?",
		CorrectOptionID:   "opt_contra_liability",
		Options:           stage2Opts,
		CausalHint:        "Does this valuation account directly offset a long-term debt liability, or does it represent an independent resource or periodic expense?",
		Explanation:       "Discount on Bonds Payable is a Contra-Liability account that carries a normal debit balance and is subtracted directly from Bonds Payable on the balance sheet.",
		RelevantConceptID: "contra_liability",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 7013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the balance in Discount on Bonds Payable increase or decrease when the discounted bond is originally issued?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When new bonds are sold below par value, does the cumulative discount attached to the debt portfolio grow larger or smaller?",
		Explanation:       "Discount on Bonds Payable increases as the initial discount is established upon issuance.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 7014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger increases a Contra-Liability account with a normal debit balance?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Contra-liabilities have normal balances opposite to liabilities. Do they increase on the left or the right side?",
		Explanation:       "Contra-liability accounts increase on the left side (Debit), offsetting the credit balance of the primary liability.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_bonds_payable", Text: "Bonds Payable (Liability)"},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 7015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liability account records the full par value principal repayment obligation due to bondholders at maturity?",
		CorrectOptionID:   "opt_bonds_payable",
		Options:           stage5Opts,
		CausalHint:        "Which long-term debt account must always be credited for the full contractual face value owed upon maturity?",
		Explanation:       "Bonds Payable is credited for the contractual face amount owed at maturity.",
		RelevantConceptID: "bonds_payable_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_balanced_discount", Text: fmt.Sprintf("Debit Cash %s / Debit Discount on Bonds Payable %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), discountMoney.FormatDollars(), faceMoney.FormatDollars())},
		{ID: "opt_no_discount", Text: fmt.Sprintf("Debit Cash %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), cashMoney.FormatDollars()), ErrorTag: engine.TagM07DiscountRecordedAsLoss},
		{ID: "opt_loss_discount", Text: fmt.Sprintf("Debit Cash %s / Debit Loss on Bond Issuance %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), discountMoney.FormatDollars(), faceMoney.FormatDollars()), ErrorTag: engine.TagM07DiscountRecordedAsLoss},
	}, seed, 7016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record issuing the bonds at a discount?",
		CorrectOptionID:   "opt_balanced_discount",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry records the cash proceeds, recognizes the contra-liability discount, and records the full face value debt obligation?",
		Explanation:       fmt.Sprintf("Debit Cash %s for proceeds received, Debit Discount on Bonds Payable %s for the difference below par, and Credit Bonds Payable %s for the full face value. The discount is amortized over the life of the bond as additional interest expense.", cashMoney.FormatDollars(), discountMoney.FormatDollars(), faceMoney.FormatDollars()),
		RelevantConceptID: "bond_discount_issuance",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing bonds at a discount affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Bonds Payable %s, -Discount %s); Equity unchanged.", cashMoney.FormatDollars(), cashMoney.FormatDollars(), faceMoney.FormatDollars(), discountMoney.FormatDollars())},
			{ID: "opt_assets_up_equity_down", Text: fmt.Sprintf("Assets increase by %s; Equity decreases by %s.", cashMoney.FormatDollars(), discountMoney.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Liabilities remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 7017),
		CausalHint:        "Cash was received, long-term debt increased, and a contra-liability reduced gross debt to carrying value. What is the net impact on Assets, Liabilities, and Equity?",
		Explanation:       fmt.Sprintf("Assets increase by the cash proceeds (+Cash %s); Net Liabilities increase by the net carrying value (+Bonds Payable %s minus Discount %s = +%s); Equity is completely unaffected upon debt issuance.", cashMoney.FormatDollars(), faceMoney.FormatDollars(), discountMoney.FormatDollars(), cashMoney.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM07BondIssuancePremiumStages builds 7 scaffolding stages for issuing bonds at a premium.
func BuildM07BondIssuancePremiumStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	premiumVal := proc.Event.Parameters["premium_minor_units"]
	if premiumVal <= 0 {
		premiumVal = amt.Cents() * 5 / 100
	}
	cashVal := proc.Event.Parameters["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() + premiumVal
	}
	cashMoney := domain.NewMoney(cashVal)
	premiumMoney := domain.NewMoney(premiumVal)
	faceMoney := amt

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_premium_bonds", Text: "Premium on Bonds Payable"},
		{ID: "opt_gain_issuance", Text: "Gain on Bond Issuance", ErrorTag: engine.TagM07PremiumRecordedAsGain},
		{ID: "opt_interest_revenue", Text: "Interest Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_discount_bonds", Text: "Discount on Bonds Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 7021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When bonds are issued above par value because the stated coupon rate exceeds the market interest rate, which adjunct liability account captures the excess cash proceeds above face value?",
		CorrectOptionID:   "opt_premium_bonds",
		Options:           stage1Opts,
		CausalHint:        "Which adjunct liability account attaches directly to Bonds Payable to reflect the premium proceeds paid by eager investors?",
		Explanation:       "Premium on Bonds Payable is an adjunct liability account that records the excess proceeds received over the maturity face value.",
		RelevantConceptID: "premium_on_bonds_payable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_adjunct_liability", Text: "Liability (Adjunct Liability)"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_revenue", Text: "Revenue"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 7022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Premium on Bonds Payable?",
		CorrectOptionID:   "opt_adjunct_liability",
		Options:           stage2Opts,
		CausalHint:        "Does this account add to the carrying value of a long-term liability, or does it represent immediate corporate income?",
		Explanation:       "Premium on Bonds Payable is an adjunct Liability account with a normal credit balance that is added to Bonds Payable on the balance sheet.",
		RelevantConceptID: "adjunct_liability",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 7023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Premium on Bonds Payable increase or decrease when the bond is initially issued at a premium?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When new debt is issued above par value, does the recorded bond premium balance increase or decrease?",
		Explanation:       "Premium on Bonds Payable increases as the premium is initially recognized at issuance.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 7024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger increases an adjunct Liability account?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Adjunct liability accounts follow the normal balance of the debt they attach to. Do liabilities increase on the debit or credit side?",
		Explanation:       "Adjunct liabilities increase on the right side (Credit), augmenting the credit balance of Bonds Payable.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 7025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is increased for the total cash proceeds received from the bond issuance?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Which liquid asset account receives the upfront proceeds paid by bond investors upon purchasing the debt issue?",
		Explanation:       "Cash is debited for the total funds received from bondholders.",
		RelevantConceptID: "cash_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_balanced_premium", Text: fmt.Sprintf("Debit Cash %s / Credit Premium on Bonds Payable %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), premiumMoney.FormatDollars(), faceMoney.FormatDollars())},
		{ID: "opt_gain_premium", Text: fmt.Sprintf("Debit Cash %s / Credit Gain on Bond Issuance %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), premiumMoney.FormatDollars(), faceMoney.FormatDollars()), ErrorTag: engine.TagM07PremiumRecordedAsGain},
		{ID: "opt_no_premium", Text: fmt.Sprintf("Debit Cash %s / Credit Bonds Payable %s", cashMoney.FormatDollars(), cashMoney.FormatDollars()), ErrorTag: engine.TagWrongAccount},
	}, seed, 7026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record issuing the bonds at a premium?",
		CorrectOptionID:   "opt_balanced_premium",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry records the full cash proceeds received, the face value of the debt, and the adjunct liability premium?",
		Explanation:       fmt.Sprintf("Debit Cash %s for total proceeds received, Credit Premium on Bonds Payable %s for the amount above par, and Credit Bonds Payable %s for the contractual face value.", cashMoney.FormatDollars(), premiumMoney.FormatDollars(), faceMoney.FormatDollars()),
		RelevantConceptID: "bond_premium_issuance",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing bonds at a premium affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Bonds Payable %s, +Premium %s); Equity unchanged.", cashMoney.FormatDollars(), cashMoney.FormatDollars(), faceMoney.FormatDollars(), premiumMoney.FormatDollars())},
			{ID: "opt_assets_up_equity_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", cashMoney.FormatDollars(), premiumMoney.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Liabilities remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 7027),
		CausalHint:        "Total cash received increases Assets, while gross debt plus premium increases total carrying liabilities. Is Equity affected upon issuing debt?",
		Explanation:       fmt.Sprintf("Assets increase by total cash received (+Cash %s); Liabilities increase by the full carrying value (+Bonds Payable %s plus +Premium %s = +%s); Equity is completely unchanged by debt financing.", cashMoney.FormatDollars(), faceMoney.FormatDollars(), premiumMoney.FormatDollars(), cashMoney.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM07BondEffectiveInterestAmortStages builds 7 scaffolding stages for effective interest method discount amortization.
func BuildM07BondEffectiveInterestAmortStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	discountVal := proc.Event.Parameters["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	cashVal := proc.Event.Parameters["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	interestMoney := amt
	cashMoney := domain.NewMoney(cashVal)
	discountMoney := domain.NewMoney(discountVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_interest_expense", Text: "Interest Expense"},
		{ID: "opt_discount_bonds", Text: "Discount on Bonds Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_bonds_payable", Text: "Bonds Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 7031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under the effective interest method, which expense account captures the true economic borrowing cost calculated as carrying value multiplied by the market interest rate?",
		CorrectOptionID:   "opt_interest_expense",
		Options:           stage1Opts,
		CausalHint:        "Which periodic income statement account reflects the total economic cost of servicing the bond debt for the period?",
		Explanation:       "Interest Expense represents the periodic borrowing cost, computed under the effective interest method as the carrying value times the market interest rate.",
		RelevantConceptID: "interest_expense_recognition",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 7032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Interest Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the periodic cost of debt financing represent an expired operational expense on the income statement, or a balance sheet obligation?",
		Explanation:       "Interest Expense is an operating/financing Expense account on the income statement that reduces Net Income.",
		RelevantConceptID: "interest_expense_recognition",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 7033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Interest Expense increase or decrease when the semiannual interest payment and amortization are recorded?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "As financing costs accrue over the passage of a semiannual interest period, does total expense accumulate higher or lower?",
		Explanation:       "Interest Expense increases as the financing cost for the six-month period is recognized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 7034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger increases an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses reduce stockholders' equity and carry normal debit balances. Do expenses increase on the left or the right side?",
		Explanation:       "Expenses increase by Debit on the left side, reducing Stockholders' Equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_discount_bonds", Text: "Discount on Bonds Payable (Contra-Liability)"},
		{ID: "opt_bonds_payable", Text: "Bonds Payable (Liability)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 7035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which contra-liability account is credited to gradually reduce its balance and amortize it toward zero at bond maturity?",
		CorrectOptionID:   "opt_discount_bonds",
		Options:           stage5Opts,
		CausalHint:        "Which contra-liability account holds the unamortized discount and must be credited so that the bond carrying value accretes toward par?",
		Explanation:       "Discount on Bonds Payable is credited to amortize a portion of the discount, bringing the net carrying value closer to face value.",
		RelevantConceptID: "bond_discount_amortization",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_balanced_amort", Text: fmt.Sprintf("Debit Interest Expense %s / Credit Cash %s / Credit Discount on Bonds Payable %s", interestMoney.FormatDollars(), cashMoney.FormatDollars(), discountMoney.FormatDollars())},
		{ID: "opt_ignoring_amort", Text: fmt.Sprintf("Debit Interest Expense %s / Credit Cash %s", cashMoney.FormatDollars(), cashMoney.FormatDollars()), ErrorTag: engine.TagM07IgnoringBondAmortization},
		{ID: "opt_reversed_amort", Text: fmt.Sprintf("Debit Interest Expense %s / Debit Discount on Bonds Payable %s / Credit Cash %s", interestMoney.FormatDollars(), discountMoney.FormatDollars(), cashMoney.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 7036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the semiannual interest payment and discount amortization?",
		CorrectOptionID:   "opt_balanced_amort",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the effective interest expense (carrying value times market rate), credits cash for the contractual coupon, and credits the discount for the difference?",
		Explanation:       fmt.Sprintf("Debit Interest Expense %s for effective interest (Carrying Value * market rate / 2), Credit Cash %s for contractual coupon (Face * coupon rate / 2), and Credit Discount on Bonds Payable %s for the difference. The amortization increases bond carrying value toward face value.", interestMoney.FormatDollars(), cashMoney.FormatDollars(), discountMoney.FormatDollars()),
		RelevantConceptID: "effective_interest_method",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording the semiannual interest payment and discount amortization affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_liab_up_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_liab_up_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Liabilities increase by %s (+Carrying Value via -Discount); Equity decreases by %s (-Interest Expense).", cashMoney.FormatDollars(), discountMoney.FormatDollars(), interestMoney.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", interestMoney.FormatDollars(), interestMoney.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 7037),
		CausalHint:        "Cash left the firm, the contra-liability discount decreased (increasing net debt carrying value), and interest was expensed. What is the net impact on the equation?",
		Explanation:       fmt.Sprintf("Assets decrease by cash paid (-Cash %s); Net Liabilities increase by the discount amortized (+carrying value %s); Equity decreases by total interest expense (-Interest Expense %s). Equation check: -%s = +%s + (-%s), which balances perfectly.", cashMoney.FormatDollars(), discountMoney.FormatDollars(), interestMoney.FormatDollars(), cashMoney.FormatDollars(), discountMoney.FormatDollars(), interestMoney.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM07BondContingentLiabilityStages builds 7 scaffolding stages for ASC 450 contingent liability accrual.
func BuildM07BondContingentLiabilityStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_litigation_loss", Text: "Litigation Loss"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagM07ContingentRecordedAsPayable},
		{ID: "opt_contingent_liability", Text: "Contingent Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance", ErrorTag: engine.TagWrongAccount},
	}, seed, 7041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under ASC 450, when an adverse litigation judgment against the company is both probable and reasonably estimable, which expense/loss account captures the estimated loss?",
		CorrectOptionID:   "opt_litigation_loss",
		Options:           stage1Opts,
		CausalHint:        "Under ASC 450, when an adverse legal outcome is probable and estimable, which income statement loss account recognizes the financial impact?",
		Explanation:       "Litigation Loss is the expense/loss account recognized in earnings when a loss contingency meets the ASC 450 accrual threshold.",
		RelevantConceptID: "loss_contingency",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense (Loss)"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 7042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Litigation Loss?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an accrued legal loss represent an immediate reduction in current period earnings, or a future economic obligation?",
		Explanation:       "Litigation Loss is an Expense/Loss account reported on the income statement that directly reduces pre-tax operating income and net earnings.",
		RelevantConceptID: "loss_contingency",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 7043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Litigation Loss increase or decrease when the loss contingency is formally accrued?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When an estimated loss contingency is newly accrued, does the cumulative total of losses for the period increase or decrease?",
		Explanation:       "Litigation Loss increases as the newly estimated contingency is recognized in the current financial period.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 7044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Expense or Loss account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Losses reduce stockholders' equity and carry normal debit balances. On which side of the ledger do losses increase?",
		Explanation:       "Losses and expenses increase on the left side (Debit), reducing Stockholders' Equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_contingent_liability", Text: "Contingent Liability (Liability)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagM07ContingentRecordedAsPayable},
	}, seed, 7045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liability account is credited to record the accrued legal obligation under ASC 450?",
		CorrectOptionID:   "opt_contingent_liability",
		Options:           stage5Opts,
		CausalHint:        "Which liability account reflects the estimated claim obligation arising from the probable loss contingency?",
		Explanation:       "Contingent Liability is credited to establish the estimated balance sheet obligation resulting from the pending lawsuit.",
		RelevantConceptID: "contingent_liability_accrual",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_loss_cr_liability", Text: fmt.Sprintf("Debit Litigation Loss %s / Credit Contingent Liability %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_loss_cr_cash", Text: fmt.Sprintf("Debit Litigation Loss %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_dr_loss_cr_ap", Text: fmt.Sprintf("Debit Litigation Loss %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM07ContingentRecordedAsPayable},
	}, seed, 7046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to accrue the contingent litigation loss under ASC 450?",
		CorrectOptionID:   "opt_dr_loss_cr_liability",
		Options:           stage6Opts,
		CausalHint:        "Which balanced entry records the estimated lawsuit loss in earnings and establishes the contingent liability on the balance sheet?",
		Explanation:       fmt.Sprintf("Debit Litigation Loss %s and Credit Contingent Liability %s. Under ASC 450, accrual is mandatory when the loss is both probable and reasonably estimable. If merely reasonably possible, footnote disclosure alone is required.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "asc_450",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does accruing a contingent litigation loss affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_liab_up_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_liab_up_eq_down", Text: fmt.Sprintf("Liabilities increase by %s (+Contingent Liability); Equity decreases by %s (-Litigation Loss); Assets unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 7047),
		CausalHint:        "A liability obligation was recorded and an expense was recognized in earnings with no cash paid yet. How does this affect Assets, Liabilities, and Equity?",
		Explanation:       fmt.Sprintf("Assets are unchanged ($0); Liabilities increase (+Contingent Liability %s); Equity decreases (-Litigation Loss %s). Check: 0 = +%s + (-%s), which preserves the fundamental equation.", amt.FormatDollars(), amt.FormatDollars(), amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

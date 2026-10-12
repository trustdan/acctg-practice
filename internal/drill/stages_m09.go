package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM09Stages dispatches to the appropriate modular stage builder for Stream-M09.
func BuildM09Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM09CookieJarReserve:
		BuildM09CookieJarReserveStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM09WarrantyAccrual:
		BuildM09WarrantyAccrualStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM09CookieJarReserveStages builds 7 scaffolding stages for cookie jar reserve buffering.
func BuildM09CookieJarReserveStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_bad_debt_expense", Text: "Bad Debt Expense"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable", ErrorTag: engine.TagDirectWriteOffExpense},
		{ID: "opt_allowance_doubtful", Text: "Allowance for Doubtful Accounts", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 9011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When management decides to intentionally increase estimated uncollectible allowances above historical loss rates during a profitable period to build a reserve buffer, which operating expense account records this estimated loss?",
		CorrectOptionID:   "opt_bad_debt_expense",
		Options:           stage1Opts,
		CausalHint:        "Which income statement expense account records the cost of expected customer credit defaults?",
		Explanation:       "Bad Debt Expense is the operating expense account used to recognize estimated uncollectible customer accounts under the allowance method.",
		RelevantConceptID: "cookie_jar_reserves",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 9012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Bad Debt Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an expense account reflect an operational resource outflow, a future obligation, or equity?",
		Explanation:       "Bad Debt Expense is an operating Expense that reduces net income on the income statement.",
		RelevantConceptID: "cookie_jar_reserves",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 9013)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage3Opts,
		CausalHint:        "Expenses increase with debits and have normal debit balances. Which side increases an expense?",
		Explanation:       "Expenses increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 9014)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does recording this discretionary reserve addition increase or decrease Bad Debt Expense?",
		CorrectOptionID:   "opt_increase",
		Options:           stage4Opts,
		CausalHint:        "When a company accrues additional bad debt expense for the period, does the expense account balance increase or decrease?",
		Explanation:       "Bad Debt Expense increases as additional uncollectible costs are accrued.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_allowance_doubtful", Text: "Allowance for Doubtful Accounts (Contra-Asset)"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagDirectWriteOffExpense},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_retained_earnings", Text: "Retained Earnings (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 9015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which balance sheet account is credited to establish the reserve buffer for future uncollectible accounts?",
		CorrectOptionID:   "opt_allowance_doubtful",
		Options:           stage5Opts,
		CausalHint:        "Under the GAAP allowance method, which contra-asset account absorbs estimated uncollectible accounts without reducing customer ledgers directly?",
		Explanation:       "Allowance for Doubtful Accounts is credited as a contra-asset against Accounts Receivable, establishing the reserve on the balance sheet.",
		RelevantConceptID: "cookie_jar_reserves",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_bde_cr_allowance", Text: fmt.Sprintf("Debit Bad Debt Expense %s / Credit Allowance for Doubtful Accounts %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_allowance_cr_bde", Text: fmt.Sprintf("Debit Allowance for Doubtful Accounts %s / Credit Bad Debt Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_bde_cr_ar", Text: fmt.Sprintf("Debit Bad Debt Expense %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDirectWriteOffExpense},
	}, seed, 9016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry recording the discretionary cookie jar reserve addition?",
		CorrectOptionID:   "opt_dr_bde_cr_allowance",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the operating expense debit with the contra-asset credit?",
		Explanation:       fmt.Sprintf("Debit Bad Debt Expense %s and Credit Allowance for Doubtful Accounts %s. This establishes the reserve while reducing current-period reported earnings.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cookie_jar_reserves",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_eq_cookie_jar", Text: fmt.Sprintf("Assets decrease by %s (-Allowance); Equity decreases by %s (-Bad Debt Exp); Liabilities unchanged. Discretionary over-accruals create 'cookie jar' reserves that suppress current earnings to artificially inflate future lean quarters, degrading earnings quality.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_eq_swap", Text: fmt.Sprintf("Asset swap: Cash increases and receivables decrease by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_eq_liab_up", Text: fmt.Sprintf("Liabilities increase by %s and equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 9017)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this discretionary accrual affect the accounting equation, and what are its implications for financial reporting quality?",
		CorrectOptionID:   "opt_eq_cookie_jar",
		Options:           stage7Opts,
		CausalHint:        "The allowance reduces net receivables on the asset side, and bad debt expense reduces retained earnings. What does this discretionary pattern imply for future earnings?",
		Explanation:       fmt.Sprintf("Assets decrease (-Allowance for Doubtful Accounts %s) and Equity decreases (-Bad Debt Expense %s); Liabilities are unchanged. Intentionally over-accruing uncollectible reserves creates a 'cookie jar' buffer that management can later reverse into earnings during difficult periods, masking true operational volatility.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cookie_jar_reserves",
	}
}

// BuildM09WarrantyAccrualStages builds 7 scaffolding stages for product warranty accruals under the matching principle.
func BuildM09WarrantyAccrualStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_warranty_expense", Text: "Warranty Expense"},
		{ID: "opt_warranty_liability", Text: "Warranty Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cogs", Text: "Cost of Goods Sold", ErrorTag: engine.TagWrongAccount},
	}, seed, 9021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under the accrual matching principle (ASC 460), when a company sells products with warranty coverage, which operating expense account records the estimated future defect costs in the period of sale?",
		CorrectOptionID:   "opt_warranty_expense",
		Options:           stage1Opts,
		CausalHint:        "Which operating expense account matches estimated future warranty repair costs against the revenues earned from product sales?",
		Explanation:       "Warranty Expense is the operating expense account that matches estimated repair and replacement obligations against current sales.",
		RelevantConceptID: "warranty_accrual",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 9022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Warranty Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an account recording the cost of servicing customer warranty claims represent an expense, asset, or liability?",
		Explanation:       "Warranty Expense is an operating Expense that reduces net operating income.",
		RelevantConceptID: "warranty_accrual",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 9023)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage3Opts,
		CausalHint:        "All expense accounts increase on the debit side. Which side increases an expense?",
		Explanation:       "Expenses increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 9024)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does estimating future warranty claims increase or decrease Warranty Expense in the period of sale?",
		CorrectOptionID:   "opt_increase",
		Options:           stage4Opts,
		CausalHint:        "When a company accrues estimated warranty costs for newly sold products, does the expense balance increase or decrease?",
		Explanation:       "Warranty Expense increases as estimated warranty costs are recognized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_warranty_liability", Text: "Warranty Liability (Liability)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagM09WarrantyExpensedWhenPaid},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 9025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liability account is credited to recognize the company's estimated obligation to repair or replace defective goods sold?",
		CorrectOptionID:   "opt_warranty_liability",
		Options:           stage5Opts,
		CausalHint:        "Which accrued liability account represents the estimated future claims owed to customers under product warranties?",
		Explanation:       "Warranty Liability is an accrued liability credited to recognize the estimated future cost of satisfying customer warranty claims.",
		RelevantConceptID: "warranty_accrual",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_warr_exp_cr_liab", Text: fmt.Sprintf("Debit Warranty Expense %s / Credit Warranty Liability %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_warr_liab_cr_exp", Text: fmt.Sprintf("Debit Warranty Liability %s / Credit Warranty Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_warr_exp_cr_cash", Text: fmt.Sprintf("Debit Warranty Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM09WarrantyExpensedWhenPaid},
	}, seed, 9026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry recording the product warranty accrual in the period of sale?",
		CorrectOptionID:   "opt_dr_warr_exp_cr_liab",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the operating expense debit with the accrued warranty liability credit?",
		Explanation:       fmt.Sprintf("Debit Warranty Expense %s and Credit Warranty Liability %s. This complies with GAAP matching by recognizing the expense when the product is sold rather than when claims are paid.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "warranty_accrual",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_eq_warranty", Text: fmt.Sprintf("Liabilities increase by %s (+Warranty Liability); Equity decreases by %s (-Warranty Expense); Assets unchanged. Accruing estimated warranties in the period of sale adheres to the matching principle, preventing cash-basis distortion of net operating profit.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_eq_assets_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash) and Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_eq_swap", Text: fmt.Sprintf("Asset swap: Cash increases and inventory decreases by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 9027)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this entry affect the balance sheet equation, and why is accruing estimated warranties essential for reporting quality?",
		CorrectOptionID:   "opt_eq_warranty",
		Options:           stage7Opts,
		CausalHint:        "A liability obligation is recognized and operating expense reduces equity. Why is recognizing estimated warranties superior to waiting for cash repair claims?",
		Explanation:       fmt.Sprintf("Liabilities increase (+Warranty Liability %s) and Equity decreases (-Warranty Expense %s); Assets are unchanged. Accrual accounting and the matching principle prevent firms from delaying recognition of known product defect costs, ensuring reported earnings accurately reflect product profitability.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "warranty_accrual",
	}
}

package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM08Stages dispatches to the appropriate modular stage builder for Stream-M08.
func BuildM08Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM08LeaseCommencementROUAsset:
		BuildM08LeaseCommencementStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM08LeasePeriodicPayment:
		BuildM08LeasePeriodicPaymentStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM08TreasuryStockRepurchase:
		BuildM08TreasuryStockRepurchaseStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM08TreasuryStockReissuance:
		BuildM08TreasuryStockReissuanceStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM08LeaseCommencementStages builds 7 scaffolding stages for initial lease commencement under ASC 842.
func BuildM08LeaseCommencementStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_rou_asset", Text: "Right-of-Use Asset"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagM08LeaseExpensedAtCommencement},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagM08ROUAssetRecordedAsEquity},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagWrongAccount},
	}, seed, 8011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under ASC 842, which asset account capitalizes the lessee's right to control and operate the leased facility over the contract term?",
		CorrectOptionID:   "opt_rou_asset",
		Options:           stage1Opts,
		CausalHint:        "Under ASC 842, which asset account capitalizes the lessee's right to control and operate the leased facility over the contract term?",
		Explanation:       "Right-of-Use Asset is the non-current asset account that capitalizes the lessee's economic right to control and utilize the leased physical property across the contract term.",
		RelevantConceptID: "rou_asset_recognition",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 8012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Right-of-Use Asset?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does the capitalized right to use leased facilities represent an economic resource on the balance sheet, or an expired period cost?",
		Explanation:       "Right-of-Use Asset is an Asset account on the balance sheet representing the capitalized economic resource of leased property.",
		RelevantConceptID: "rou_asset_recognition",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 8013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Right-of-Use Asset increase or decrease at lease commencement?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "At lease commencement, does the initial recognition of the right-of-use asset add to or reduce total balance sheet resources?",
		Explanation:       "The Right-of-Use Asset increases from zero as the lease commitment begins and the asset is capitalized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 8014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets carry normal debit balances on the left side of the ledger. On which side is an asset increase recorded?",
		Explanation:       "Assets increase by Debit on the left side of the ledger.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_lease_liability", Text: "Lease Liability (Liability)"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_retained_earnings", Text: "Retained Earnings (Equity)", ErrorTag: engine.TagM08ROUAssetRecordedAsEquity},
	}, seed, 8015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which balance sheet liability account reflects the present value of future lease payments owed to the lessor?",
		CorrectOptionID:   "opt_lease_liability",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet liability account reflects the present value of future lease payments owed to the lessor?",
		Explanation:       "Lease Liability is credited to record the contractual present value obligation to make future lease payments.",
		RelevantConceptID: "lease_liability_asc842",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_rou_cr_ll", Text: fmt.Sprintf("Debit Right-of-Use Asset %s / Credit Lease Liability %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM08LeaseExpensedAtCommencement},
		{ID: "opt_dr_ll_cr_rou", Text: fmt.Sprintf("Debit Lease Liability %s / Credit Right-of-Use Asset %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 8016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record initial lease commencement under ASC 842?",
		CorrectOptionID:   "opt_dr_rou_cr_ll",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry capitalizes the right-of-use asset and recognizes the corresponding long-term lease liability under ASC 842?",
		Explanation:       fmt.Sprintf("Debit Right-of-Use Asset %s and Credit Lease Liability %s. Under ASC 842, all leases exceeding 12 months are recognized on the balance sheet at the present value of lease payments.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "lease_commencement",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording initial lease commencement affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Right-of-Use Asset); Liabilities increase by %s (+Lease Liability); Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 8017),
		CausalHint:        "When a lease is capitalized by recording both an asset and a liability for equal present value amounts, what is the net effect on the accounting equation?",
		Explanation:       fmt.Sprintf("Assets increase (+Right-of-Use Asset %s) and Liabilities increase (+Lease Liability %s); Equity is unchanged. ASC 842 eliminated off-balance sheet operating lease financing.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM08LeasePeriodicPaymentStages builds 7 scaffolding stages for periodic lease payment under ASC 842.
func BuildM08LeasePeriodicPaymentStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	totalAmt := amt
	interestVal := inst.Parameters["interest_minor_units"]
	if interestVal <= 0 {
		interestVal = totalAmt.Cents() * 20 / 100
	}
	principalVal := inst.Parameters["principal_minor_units"]
	if principalVal <= 0 {
		principalVal = totalAmt.Cents() - interestVal
	}
	interestAmt := domain.NewMoney(interestVal)
	principalAmt := domain.NewMoney(principalVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_lease_liability", Text: "Lease Liability"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
	}, seed, 8021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When a periodic lease payment is disbursed, which balance sheet liability is reduced for the principal portion of the payment?",
		CorrectOptionID:   "opt_lease_liability",
		Options:           stage1Opts,
		CausalHint:        "When a periodic lease payment is disbursed, which balance sheet liability is reduced for the principal portion of the payment?",
		Explanation:       "Lease Liability is the financing obligation account that reflects the remaining principal balance owed to the lessor, which is reduced as payments are made.",
		RelevantConceptID: "lease_liability_amortization",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 8022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Lease Liability?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does the lease obligation owed to the lessor represent a legal liability, an owned asset, or stockholders' equity?",
		Explanation:       "Lease Liability is a Liability account on the balance sheet representing financing obligations owed under lease contracts.",
		RelevantConceptID: "lease_liability_amortization",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 8023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the principal Lease Liability increase or decrease when making this lease payment?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "When cash is paid to reduce the outstanding obligation owed to the lessor, does the lease liability increase or decrease?",
		Explanation:       "Lease Liability decreases because the tenant has paid down a portion of the contractual principal debt.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 8024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to reduce a Liability account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities maintain normal credit balances on the right side. On which side is a liability reduction entered?",
		Explanation:       "Liabilities decrease by Debit on the left side of the ledger.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_interest_and_cash", Text: "Interest Expense (Debit) and Cash (Credit)"},
		{ID: "opt_cash_only", Text: "Cash only (omitting interest expense)", ErrorTag: engine.TagM08IgnoringLeaseInterest},
		{ID: "opt_rou_asset_only", Text: "Right-of-Use Asset only", ErrorTag: engine.TagWrongAccount},
	}, seed, 8025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Besides reducing the lease liability, which accounts reflect the financing charge accrued and the total cash disbursed?",
		CorrectOptionID:   "opt_interest_and_cash",
		Options:           stage5Opts,
		CausalHint:        "Besides reducing the lease liability, which accounts reflect the period financing cost incurred and the cash disbursed to the lessor?",
		Explanation:       "Interest Expense is debited for the financing cost accrued on the lease liability, and Cash is credited for the total payment disbursed.",
		RelevantConceptID: "interest_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ll_dr_int_cr_cash", Text: fmt.Sprintf("Debit Lease Liability %s, Debit Interest Expense %s, and Credit Cash %s", principalAmt.FormatDollars(), interestAmt.FormatDollars(), totalAmt.FormatDollars())},
		{ID: "opt_dr_ll_cr_cash_only", Text: fmt.Sprintf("Debit Lease Liability %s and Credit Cash %s (omitting interest)", totalAmt.FormatDollars(), totalAmt.FormatDollars()), ErrorTag: engine.TagM08IgnoringLeaseInterest},
		{ID: "opt_reversed", Text: fmt.Sprintf("Debit Cash %s, Credit Lease Liability %s, and Credit Interest Expense %s", totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 8026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for the periodic lease payment?",
		CorrectOptionID:   "opt_dr_ll_dr_int_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry splits the periodic lease payment between principal liability reduction and financing interest expense while crediting cash?",
		Explanation:       fmt.Sprintf("Debit Lease Liability %s (principal reduction), Debit Interest Expense %s (financing charge), and Credit Cash %s (total payment). Under ASC 842, periodic payments amortize both the liability balance and interest cost.", principalAmt.FormatDollars(), interestAmt.FormatDollars(), totalAmt.FormatDollars()),
		RelevantConceptID: "periodic_lease_payment",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does the periodic lease payment affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_liab_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_liab_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Liabilities decrease by %s (-Lease Liability); Equity decreases by %s (-Interest Expense).", totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars())},
			{ID: "opt_assets_down_liab_down_only", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s; Equity unchanged.", totalAmt.FormatDollars(), totalAmt.FormatDollars()), ErrorTag: engine.TagM08IgnoringLeaseInterest},
			{ID: "opt_no_effect", Text: "Total Assets and Liabilities remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 8027),
		CausalHint:        "Cash was disbursed, lease principal was reduced, and interest expense was recognized. What is the net impact across Assets, Liabilities, and Equity?",
		Explanation:       fmt.Sprintf("Assets decrease (-Cash %s); Liabilities decrease (-Lease Liability %s); Equity decreases (-Interest Expense %s). The equation balances: -%s = -%s + (-%s).", totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars(), totalAmt.FormatDollars(), principalAmt.FormatDollars(), interestAmt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM08TreasuryStockRepurchaseStages builds 7 scaffolding stages for share repurchase under the cost method.
func BuildM08TreasuryStockRepurchaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_treasury_stock", Text: "Treasury Stock"},
		{ID: "opt_equipment", Text: "Equipment", ErrorTag: engine.TagM08TreasuryStockRecordedAsAsset},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_investments", Text: "Investments in Marketable Securities", ErrorTag: engine.TagM08TreasuryStockRecordedAsAsset},
	}, seed, 8031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When a company buys back its own shares to hold in treasury under the cost method, which contra-equity account is debited?",
		CorrectOptionID:   "opt_treasury_stock",
		Options:           stage1Opts,
		CausalHint:        "When a company buys back its own shares to hold in treasury under the cost method, which contra-equity account is debited?",
		Explanation:       "Treasury Stock is the contra-equity account that records the company's buyback of its own outstanding shares from the open market, reducing total stockholders' equity.",
		RelevantConceptID: "treasury_stock_cost_method",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_contra_equity", Text: "Contra-Equity (Equity)"},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagM08TreasuryStockRecordedAsAsset},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_expense", Text: "Expense"},
	}, seed, 8032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Treasury Stock?",
		CorrectOptionID:   "opt_contra_equity",
		Options:           stage2Opts,
		CausalHint:        "Does repurchased company stock represent an asset owned by the business, or a deduction from total stockholders' equity?",
		Explanation:       "Treasury Stock is a Contra-Equity account that directly offsets and reduces total Stockholders' Equity on the balance sheet.",
		RelevantConceptID: "contra_equity",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 8033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the balance in the Treasury Stock account increase or decrease when purchasing shares?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When additional company shares are bought back and held in treasury, does the balance in the Treasury Stock account increase or decrease?",
		Explanation:       "Treasury Stock increases in balance, which exerts a larger negative deduction against total stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 8034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase a Contra-Equity account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Contra-equity accounts have the opposite normal balance of equity. On which side does a contra-equity account increase?",
		Explanation:       "Contra-equity accounts increase by Debit on the left side, reducing overall stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagWrongAccount},
	}, seed, 8035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is reduced when paying cash to repurchase shares from investors?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Which liquid asset account is disbursed and reduced to settle the share repurchase transaction with selling investors?",
		Explanation:       "Cash is credited because corporate funds were paid out to purchase the shares on the open market.",
		RelevantConceptID: "cash_outflow_financing",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ts_cr_cash", Text: fmt.Sprintf("Debit Treasury Stock %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_ts", Text: fmt.Sprintf("Debit Cash %s / Credit Treasury Stock %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_inv_cr_cash", Text: fmt.Sprintf("Debit Investment in Stock %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM08TreasuryStockRecordedAsAsset},
	}, seed, 8036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the share repurchase under the cost method?",
		CorrectOptionID:   "opt_dr_ts_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry records the share buyback under the cost method by debiting the contra-equity account and crediting cash?",
		Explanation:       fmt.Sprintf("Debit Treasury Stock %s and Credit Cash %s. Under the cost method, repurchased shares are recorded at acquisition cost as a contra-equity account.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "treasury_stock_cost_method",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does repurchasing company shares for cash affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Treasury Stock); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 8037),
		CausalHint:        "When cash is disbursed to buy back shares into treasury, how are total Assets, Liabilities, and Equity affected?",
		Explanation:       fmt.Sprintf("Assets decrease (-Cash %s) and Equity decreases (-Treasury Stock %s); Liabilities are unchanged. Total stockholders' equity is reduced by the cost of the shares repurchased.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM08TreasuryStockReissuanceStages builds 7 scaffolding stages for reissuing treasury shares above cost.
func BuildM08TreasuryStockReissuanceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	totalAmt := amt
	costVal := inst.Parameters["cost_minor_units"]
	if costVal <= 0 {
		costVal = totalAmt.Cents() * 2 / 3
	}
	apicVal := inst.Parameters["apic_minor_units"]
	if apicVal <= 0 {
		apicVal = totalAmt.Cents() - costVal
	}
	costAmt := domain.NewMoney(costVal)
	apicAmt := domain.NewMoney(apicVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_gain_on_disposal", Text: "Gain on Reissuance of Stock", ErrorTag: engine.TagM08GainRecordedOnTreasuryReissuance},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 8041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When treasury shares are reissued to market investors for cash, which asset account increases from the proceeds received?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "When treasury shares are reissued to outside investors for cash, which asset account increases from the incoming proceeds?",
		Explanation:       "Cash is the liquid asset account that receives the incoming funds from selling the treasury shares to new investors.",
		RelevantConceptID: "treasury_stock_reissuance",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 8042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "What category of balance sheet account is Cash when receiving proceeds from a stock reissuance?",
		Explanation:       "Cash is a Current Asset on the balance sheet representing available liquid funds.",
		RelevantConceptID: "treasury_stock_reissuance",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 8043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon receiving investor proceeds for the reissued shares?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When investors pay funds to purchase reissued treasury shares, does the company's cash balance increase or decrease?",
		Explanation:       "Cash increases because new financial capital has been received by the corporation.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 8044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "On which side of the ledger is an increase to an Asset account recorded?",
		Explanation:       "Assets increase by Debit on the left side of the ledger.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ts_and_apic", Text: "Treasury Stock (at cost) and Additional Paid-in Capital (for excess)"},
		{ID: "opt_ts_and_gain", Text: "Treasury Stock and Gain on Sale of Stock", ErrorTag: engine.TagM08GainRecordedOnTreasuryReissuance},
		{ID: "opt_common_and_re", Text: "Common Stock and Retained Earnings", ErrorTag: engine.TagWrongAccount},
	}, seed, 8045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which equity accounts are credited when reissuing treasury shares above original cost?",
		CorrectOptionID:   "opt_ts_and_apic",
		Options:           stage5Opts,
		CausalHint:        "When treasury shares are reissued above cost, which equity accounts are credited for the cost basis and the premium received?",
		Explanation:       "Treasury Stock is credited to remove the historical cost of shares reissued, and Additional Paid-in Capital is credited for proceeds received above cost.",
		RelevantConceptID: "apic_treasury_stock",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_ts_apic", Text: fmt.Sprintf("Debit Cash %s, Credit Treasury Stock %s (cost), and Credit Additional Paid-in Capital %s (excess)", totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars())},
		{ID: "opt_dr_cash_cr_ts_gain", Text: fmt.Sprintf("Debit Cash %s, Credit Treasury Stock %s, and Credit Gain on Reissuance %s", totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars()), ErrorTag: engine.TagM08GainRecordedOnTreasuryReissuance},
		{ID: "opt_dr_cash_cr_ts_only", Text: fmt.Sprintf("Debit Cash %s / Credit Treasury Stock %s", totalAmt.FormatDollars(), totalAmt.FormatDollars()), ErrorTag: engine.TagWrongAccount},
	}, seed, 8046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record reissuing treasury shares above cost?",
		CorrectOptionID:   "opt_dr_cash_cr_ts_apic",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry records total cash proceeds, removes treasury stock at original cost, and credits additional paid-in capital for the excess?",
		Explanation:       fmt.Sprintf("Debit Cash %s, Credit Treasury Stock %s (cost), and Credit Additional Paid-in Capital %s (excess). Under US GAAP, no gain or loss is recognized on transactions in a company's own stock.", totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars()),
		RelevantConceptID: "no_gain_or_loss_on_treasury",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does reissuing treasury shares above cost affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+APIC %s, -Treasury Stock contra %s); Liabilities unchanged.", totalAmt.FormatDollars(), totalAmt.FormatDollars(), apicAmt.FormatDollars(), costAmt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", totalAmt.FormatDollars(), totalAmt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Equity remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 8047),
		CausalHint:        "Cash proceeds were received, contra-equity was removed at cost, and paid-in capital was created for the excess. What is the net impact on the balance sheet equation?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Equity increases (+%s total: Treasury Stock contra-equity decreases by %s and APIC increases by %s); Liabilities are unchanged. Total equation balances: +%s = 0 + (+%s).", totalAmt.FormatDollars(), totalAmt.FormatDollars(), costAmt.FormatDollars(), apicAmt.FormatDollars(), totalAmt.FormatDollars(), totalAmt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

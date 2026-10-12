package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM06Stages dispatches to the appropriate modular stage builder for Stream-M06.
func BuildM06Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM06StraightLineDepr:
		BuildM06StraightLineDeprStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM06MACRSTaxTiming:
		BuildM06MACRSTaxTimingStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM06AssetDisposalGainLoss:
		BuildM06AssetDisposalGainLossStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM06ImpairmentRecoverability:
		BuildM06ImpairmentRecoverabilityStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM06StraightLineDeprStages builds 7 scaffolding stages for straight-line depreciation adjusting entry.
func BuildM06StraightLineDeprStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_depreciation_expense", Text: "Depreciation Expense"},
		{ID: "opt_equipment", Text: "Equipment", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation", ErrorTag: engine.TagM06AccumulatedDeprExpensed},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
	}, seed, 6011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "At fiscal year-end, which expense account is used to record the periodic allocation of equipment cost over its useful life?",
		CorrectOptionID:   "opt_depreciation_expense",
		Options:           stage1Opts,
		CausalHint:        "When equipment is utilized during an accounting period, which operational expense account captures the consumed service potential of the asset?",
		Explanation:       "Depreciation Expense is the operating expense account that records the periodic cost allocation of physical productive assets over their estimated useful life.",
		RelevantConceptID: "straight_line_depreciation",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 6012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Depreciation Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the periodic expired cost of using productive machinery represent an operating expense on the income statement, or an unexpired asset balance?",
		Explanation:       "Depreciation Expense is an operating Expense account on the income statement that directly reduces operating income and Net Income.",
		RelevantConceptID: "straight_line_depreciation",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 6013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Depreciation Expense increase or decrease when recording this year-end adjusting entry?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "As another accounting period expires and equipment utility is consumed, does total recorded depreciation expense accumulate upward or downward?",
		Explanation:       "Depreciation Expense increases as additional asset cost expires and is recognized in current-period earnings.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 6014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses carry normal debit balances and reduce stockholders' equity. Do expenses increase on the left or the right side?",
		Explanation:       "Expenses increase by Debit on the left side, reflecting a reduction in stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation (Contra-Asset)"},
		{ID: "opt_equipment", Text: "Equipment (Asset)", ErrorTag: engine.TagM06DirectAssetCredit},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 6015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which contra-asset balance sheet account is credited to accumulate the total cost expired to date without altering original equipment cost?",
		CorrectOptionID:   "opt_accumulated_depreciation",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet contra-account maintains the cumulative depreciation recognized against equipment across its operational life?",
		Explanation:       "Accumulated Depreciation is a contra-asset account credited to offset Equipment, allowing financial statement users to see both historical cost and net carrying value.",
		RelevantConceptID: "accumulated_depreciation",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_exp_cr_accum", Text: fmt.Sprintf("Debit Depreciation Expense %s / Credit Accumulated Depreciation %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_equip", Text: fmt.Sprintf("Debit Depreciation Expense %s / Credit Equipment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM06DirectAssetCredit},
		{ID: "opt_dr_accum_cr_exp", Text: fmt.Sprintf("Debit Accumulated Depreciation %s / Credit Depreciation Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 6016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced adjusting journal entry to record annual straight-line depreciation?",
		CorrectOptionID:   "opt_dr_exp_cr_accum",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the periodic operating expense with the contra-asset that reduces equipment carrying value?",
		Explanation:       fmt.Sprintf("Debit Depreciation Expense %s and Credit Accumulated Depreciation %s. Under the straight-line method, annual depreciation is computed as (Historical Cost - Salvage Value) / Useful Life.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "straight_line_depreciation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording straight-line depreciation affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Accumulated Depreciation); Equity decreases by %s (-Depreciation Expense); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 6017),
		CausalHint:        "How do an increase in a contra-asset and an increase in an operating expense impact total assets and total equity?",
		Explanation:       fmt.Sprintf("Assets decrease (-Accumulated Depreciation %s) and Equity decreases (-Depreciation Expense %s); Liabilities are unchanged. Accumulated depreciation reduces net PPE on the balance sheet.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM06MACRSTaxTimingStages builds 7 scaffolding stages for straight-line book depreciation vs accelerated tax MACRS timing differences.
func BuildM06MACRSTaxTimingStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_depreciation_expense", Text: "Depreciation Expense"},
		{ID: "opt_equipment", Text: "Equipment", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation", ErrorTag: engine.TagM06AccumulatedDeprExpensed},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagWrongAccount},
	}, seed, 6021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "For financial reporting purposes under US GAAP, which account records the annual straight-line cost allocation of equipment when accelerated MACRS is used on tax returns?",
		CorrectOptionID:   "opt_depreciation_expense",
		Options:           stage1Opts,
		CausalHint:        "For GAAP financial statements, which operating expense account records the economic consumption of productive equipment regardless of accelerated tax method rules?",
		Explanation:       "Depreciation Expense is the GAAP financial reporting expense account that reflects economic asset expiration under the matching principle.",
		RelevantConceptID: "macrs_vs_straight_line",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 6022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Depreciation Expense in financial reporting?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does book depreciation represent an income statement expense reducing current earnings, or a balance sheet resource?",
		Explanation:       "Depreciation Expense is an operating Expense on the income statement that reduces book pre-tax income.",
		RelevantConceptID: "macrs_vs_straight_line",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 6023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Depreciation Expense increase or decrease when recognizing book depreciation for the year?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing periodic asset expiration increase or decrease the total operating expenses reported on the income statement?",
		Explanation:       "Depreciation Expense increases as additional asset cost is recognized during the reporting period.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 6024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger increases an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expense accounts carry normal debit balances that reduce equity. Which side records an increase?",
		Explanation:       "Expenses increase by Debit on the left side of the ledger, reducing stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation (Contra-Asset)"},
		{ID: "opt_equipment", Text: "Equipment (Asset)", ErrorTag: engine.TagM06DirectAssetCredit},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 6025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which contra-asset account is credited on the company's financial books to offset equipment to its net carrying value?",
		CorrectOptionID:   "opt_accumulated_depreciation",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet contra-asset account accumulates the periodic straight-line depreciation recognized on the company's financial books?",
		Explanation:       "Accumulated Depreciation is credited to offset the Equipment cost account, presenting net book value on the balance sheet.",
		RelevantConceptID: "depreciation_timing_difference",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_exp_cr_accum", Text: fmt.Sprintf("Debit Depreciation Expense %s / Credit Accumulated Depreciation %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_equip", Text: fmt.Sprintf("Debit Depreciation Expense %s / Credit Equipment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM06DirectAssetCredit},
		{ID: "opt_dr_accum_cr_exp", Text: fmt.Sprintf("Debit Accumulated Depreciation %s / Credit Depreciation Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 6026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record GAAP book depreciation expense?",
		CorrectOptionID:   "opt_dr_exp_cr_accum",
		Options:           stage6Opts,
		CausalHint:        "Which entry correctly records straight-line depreciation expense for financial reporting and credits the contra-asset account?",
		Explanation:       fmt.Sprintf("Debit Depreciation Expense %s and Credit Accumulated Depreciation %s. GAAP requires straight-line depreciation when it best matches economic consumption; accelerated MACRS deductions on tax filings create temporary timing differences.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "depreciation_timing_difference",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording straight-line book depreciation affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Accumulated Depreciation); Equity decreases by %s (-Depreciation Expense); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 6027),
		CausalHint:        "When Accumulated Depreciation increases and Depreciation Expense is recognized, what is the net impact on the balance sheet equation?",
		Explanation:       fmt.Sprintf("Assets decrease (-Accumulated Depreciation %s) and Equity decreases (-Depreciation Expense %s); Liabilities are unchanged. Book depreciation reduces book net income and net asset carrying value.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM06AssetDisposalGainLossStages builds 7 scaffolding stages for derecognition of equipment upon sale with gain.
func BuildM06AssetDisposalGainLossStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	accumVal := proc.Event.Parameters["accum_depr_minor_units"]
	if accumVal <= 0 {
		accumVal = amt.Cents() * 80 / 100
	}
	cashVal := proc.Event.Parameters["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() * 25 / 100
	}
	gainVal := proc.Event.Parameters["gain_minor_units"]
	if gainVal <= 0 {
		gainVal = cashVal - (amt.Cents() - accumVal)
	}

	cashAmt := domain.NewMoney(cashVal)
	accumAmt := domain.NewMoney(accumVal)
	gainAmt := domain.NewMoney(gainVal)
	bookValAmt := domain.NewMoney(amt.Cents() - accumVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_gain_on_disposal", Text: "Gain on Disposal"},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagM06GainTreatedAsRevenue},
		{ID: "opt_depreciation_expense", Text: "Depreciation Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_retained_earnings", Text: "Retained Earnings", ErrorTag: engine.TagWrongAccount},
	}, seed, 6031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When equipment is sold for cash proceeds exceeding its net book value (historical cost minus accumulated depreciation), which account records the excess gain recognized in earnings?",
		CorrectOptionID:   "opt_gain_on_disposal",
		Options:           stage1Opts,
		CausalHint:        "When a capital asset is sold for more than its remaining book carrying value, which account captures this peripheral transaction profit in net income?",
		Explanation:       "Gain on Disposal is the account that captures the excess of net disposal proceeds over the carrying book value of the retired asset.",
		RelevantConceptID: "gain_on_disposal",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_revenue", Text: "Revenue / Other Income"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_expense", Text: "Expense"},
	}, seed, 6032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Gain on Disposal?",
		CorrectOptionID:   "opt_revenue",
		Options:           stage2Opts,
		CausalHint:        "Does a gain resulting from selling productive equipment represent other income that increases stockholders' equity, or an asset?",
		Explanation:       "Gain on Disposal is classified under Other Income / Revenues on the income statement, increasing Net Income and Stockholders' Equity.",
		RelevantConceptID: "gain_on_disposal",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 6033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Gain on Disposal increase or decrease when recognizing the profit from the equipment sale?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When a profitable asset sale transaction occurs, does total gain recognized on the income statement increase or decrease?",
		Explanation:       "Gain on Disposal increases to reflect the net gain realized upon derecognizing the asset.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 6034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase a Gain or Revenue account?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Revenues and gains increase stockholders' equity and maintain normal credit balances. Do they increase by debit on the left or credit on the right?",
		Explanation:       "Revenues and gains increase by Credit on the right side of the ledger, reflecting an addition to stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_equipment", Text: "Equipment (Asset)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation (Contra-Asset)", ErrorTag: engine.TagM06DisposalAccumDeprRetained},
	}, seed, 6035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account must be credited for its original cost to completely remove the sold equipment from the balance sheet?",
		CorrectOptionID:   "opt_equipment",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet asset account records the original acquisition cost of the sold equipment and must now be completely eliminated?",
		Explanation:       "Equipment is credited for its full original historical cost to eliminate the derecognized asset from the company's ledger.",
		RelevantConceptID: "asset_derecognition",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_accum_cr_equip_gain", Text: fmt.Sprintf("Debit Cash %s, Debit Accumulated Depreciation %s / Credit Equipment %s, Credit Gain on Disposal %s", cashAmt.FormatDollars(), accumAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars())},
		{ID: "opt_dr_cash_cr_equip_gain", Text: fmt.Sprintf("Debit Cash %s / Credit Equipment %s, Credit Gain on Disposal %s", cashAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars()), ErrorTag: engine.TagM06DisposalAccumDeprRetained},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", cashAmt.FormatDollars(), cashAmt.FormatDollars()), ErrorTag: engine.TagM06GainTreatedAsRevenue},
	}, seed, 6036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the cash receipt, derecognition of equipment and accumulated depreciation, and gain on disposal?",
		CorrectOptionID:   "opt_dr_cash_accum_cr_equip_gain",
		Options:           stage6Opts,
		CausalHint:        "Which compound entry brings in cash proceeds, eliminates accumulated depreciation, zeroes out original equipment cost, and recognizes the gain?",
		Explanation:       fmt.Sprintf("Debit Cash %s, Debit Accumulated Depreciation %s, Credit Equipment %s, and Credit Gain on Disposal %s. Equipment carrying value was %s (%s cost - %s depr); selling for %s cash yields a %s gain.", cashAmt.FormatDollars(), accumAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars(), bookValAmt.FormatDollars(), amt.FormatDollars(), accumAmt.FormatDollars(), cashAmt.FormatDollars(), gainAmt.FormatDollars()),
		RelevantConceptID: "asset_disposal",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does derecognizing the equipment upon sale for a gain affect the fundamental accounting equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Net Assets increase by %s (+Cash %s, +Accumulated Depreciation removal %s, -Equipment %s); Equity increases by %s (+Gain on Disposal); Liabilities unchanged.", gainAmt.FormatDollars(), cashAmt.FormatDollars(), accumAmt.FormatDollars(), amt.FormatDollars(), gainAmt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", cashAmt.FormatDollars(), cashAmt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 6037),
		CausalHint:        "Net book value of the asset was surrendered in exchange for cash. How does this net gain affect total assets and total equity?",
		Explanation:       fmt.Sprintf("Net Assets increase by the gain amount of %s (+Cash %s minus Net Book Value %s); Equity increases by %s (+Gain on Disposal); Liabilities are unchanged.", gainAmt.FormatDollars(), cashAmt.FormatDollars(), bookValAmt.FormatDollars(), gainAmt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM06ImpairmentRecoverabilityStages builds 7 scaffolding stages for ASC 360 two-step impairment testing.
func BuildM06ImpairmentRecoverabilityStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_loss_on_impairment", Text: "Loss on Impairment"},
		{ID: "opt_depreciation_expense", Text: "Depreciation Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cost_of_goods_sold", Text: "Cost of Goods Sold", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation", ErrorTag: engine.TagM06AccumulatedDeprExpensed},
	}, seed, 6041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under ASC 360, when equipment fails the Step 1 recoverability test and is written down to fair value, which income statement account records the write-down?",
		CorrectOptionID:   "opt_loss_on_impairment",
		Options:           stage1Opts,
		CausalHint:        "Under ASC 360, when carrying value cannot be recovered through future undiscounted cash flows, which loss account recognizes the sudden write-down to fair value in operating income?",
		Explanation:       "Loss on Impairment is the income statement account that reflects the sudden permanent reduction of an asset's carrying value down to its fair value.",
		RelevantConceptID: "loss_on_impairment",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 6042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Loss on Impairment?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an impairment write-down represent an operating loss/expense reducing current period net earnings, or an asset?",
		Explanation:       "Loss on Impairment is an Expense/Loss account presented within continuing operations on the income statement that reduces Net Income.",
		RelevantConceptID: "asc_360_impairment",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 6043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Loss on Impairment increase or decrease when recording this asset write-down?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recording a new impairment loss write-down increase or decrease total losses and expenses on the income statement?",
		Explanation:       "Loss on Impairment increases as the asset write-down is recognized in current-period earnings.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 6044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase a Loss or Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Losses and expenses reduce stockholders' equity and carry normal debit balances. Do they increase on the left or right side?",
		Explanation:       "Losses increase by Debit on the left side, reducing Stockholders' Equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_accumulated_depreciation", Text: "Accumulated Depreciation (Contra-Asset)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 6045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which contra-asset account is credited to write down the net carrying value of the equipment to its fair value?",
		CorrectOptionID:   "opt_accumulated_depreciation",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet contra-account is increased to write down the net carrying amount of the impaired equipment?",
		Explanation:       "Accumulated Depreciation is credited to write down the carrying value of the equipment to fair value while preserving historical cost.",
		RelevantConceptID: "fair_value_measurement",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_loss_cr_accum", Text: fmt.Sprintf("Debit Loss on Impairment %s / Credit Accumulated Depreciation %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_accum_cr_loss", Text: fmt.Sprintf("Debit Accumulated Depreciation %s / Credit Loss on Impairment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_loss_cr_cash", Text: fmt.Sprintf("Debit Loss on Impairment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 6046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the ASC 360 equipment impairment write-down?",
		CorrectOptionID:   "opt_dr_loss_cr_accum",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry charges the impairment write-down to earnings and increases accumulated depreciation to reduce carrying value?",
		Explanation:       fmt.Sprintf("Debit Loss on Impairment %s and Credit Accumulated Depreciation %s. Under ASC 360, when undiscounted future cash flows fall below carrying value, the impairment loss equals Carrying Value minus Fair Value.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "asc_360_impairment",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording the ASC 360 impairment write-down affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Accumulated Depreciation); Equity decreases by %s (-Loss on Impairment); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 6047),
		CausalHint:        "How do crediting Accumulated Depreciation and debiting Loss on Impairment impact total assets and total equity?",
		Explanation:       fmt.Sprintf("Assets decrease (-Accumulated Depreciation %s) and Equity decreases (-Loss on Impairment %s); Liabilities are unchanged. Under US GAAP, impairment write-downs establish a new cost basis and cannot be reversed even if fair value recovers subsequently.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

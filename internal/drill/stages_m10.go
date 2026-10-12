package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM10Stages dispatches to the appropriate modular stage builder for Stream-M10.
func BuildM10Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM10ForecastingOperatingExpenseAP:
		BuildM10ForecastingOperatingExpenseAPStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM10ForecastingCapEx:
		BuildM10ForecastingCapExStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM10ForecastingDebtPlug:
		BuildM10ForecastingDebtPlugStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM10ForecastingOperatingExpenseAPStages builds 7 scaffolding stages for pro-forma opex linked to working capital (AP).
func BuildM10ForecastingOperatingExpenseAPStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_rent_expense", Text: "Rent Expense"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagM10ForecastingAPOmitted},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 1011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "In constructing a pro-forma financial statement forecast (Module 10), which operating expense account models the periodic facility costs required for ongoing operations?",
		CorrectOptionID:   "opt_rent_expense",
		Options:           stage1Opts,
		CausalHint:        "Which income statement expense account captures projected periodic facility costs in the operating budget?",
		Explanation:       "Rent Expense is an operating expense account used to model projected ongoing facility and occupancy costs in the pro-forma income statement.",
		RelevantConceptID: "forecasting_operating_working_capital",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Rent Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an account reflecting operational resource consumption represent an expense, asset, or liability?",
		Explanation:       "Rent Expense is an operating Expense that reduces net income on the pro-forma income statement.",
		RelevantConceptID: "forecasting_operating_working_capital",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 1013)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage3Opts,
		CausalHint:        "Expenses have normal debit balances and increase on the debit side. Which side increases an expense?",
		Explanation:       "Expenses increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1014)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does projecting additional operating rent increase or decrease Rent Expense in the pro-forma forecast?",
		CorrectOptionID:   "opt_increase",
		Options:           stage4Opts,
		CausalHint:        "When modeling future operational periods, does the projected expense balance increase or decrease?",
		Explanation:       "Rent Expense increases as operating expenses are projected for future periods.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagM10ForecastingAPOmitted},
		{ID: "opt_retained_earnings", Text: "Retained Earnings (Equity)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 1015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Assuming trade credit payment terms based on historical days-payable-outstanding (DPO), which liability account is credited to model the projected obligation to vendors?",
		CorrectOptionID:   "opt_accounts_payable",
		Options:           stage5Opts,
		CausalHint:        "When operating expenses are incurred on credit rather than paid immediately in cash, which current trade liability increases?",
		Explanation:       "Accounts Payable is credited to reflect the projected trade liabilities for operating expenses incurred on account.",
		RelevantConceptID: "forecasting_operating_working_capital",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_rent_cr_ap", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_ap_cr_rent", Text: fmt.Sprintf("Debit Accounts Payable %s / Credit Rent Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM10ForecastingAPOmitted},
	}, seed, 1016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry modeling projected operating expenses incurred on account in pro-forma statements?",
		CorrectOptionID:   "opt_dr_rent_cr_ap",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the operating expense debit with the trade liability credit?",
		Explanation:       fmt.Sprintf("Debit Rent Expense %s and Credit Accounts Payable %s. This models the income statement expense and the corresponding working capital trade payable.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_operating_working_capital",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_eq_fore_ap", Text: fmt.Sprintf("Liabilities increase by %s (+Accounts Payable); Equity decreases by %s (-Rent Expense); Assets unchanged. In integrated financial forecasting, linking operating expenses to working capital drivers (like DPO) ensures that the pro-forma income statement seamlessly feeds projected balance sheet trade liabilities and indirect operating cash flows.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_eq_assets_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash) and Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_eq_swap", Text: fmt.Sprintf("Asset swap: Cash increases and receivables decrease by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 1017)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this projected entry affect the balance sheet equation, and why is modeling trade payables critical for forecast accuracy?",
		CorrectOptionID:   "opt_eq_fore_ap",
		Options:           stage7Opts,
		CausalHint:        "Liabilities increase while operating expense reduces equity. Why is linking expenses to accounts payable important for pro-forma cash flow forecasts?",
		Explanation:       fmt.Sprintf("Liabilities increase (+Accounts Payable %s) and Equity decreases (-Rent Expense %s); Assets are unchanged. Modeling trade liabilities via operational drivers (such as days-payable-outstanding) enables the financial model to accurately forecast cash disbursements and operating cash flow timing.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_operating_working_capital",
	}
}

// BuildM10ForecastingCapExStages builds 7 scaffolding stages for pro-forma CapEx equipment purchases.
func BuildM10ForecastingCapExStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_equipment", Text: "Equipment"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagM10CapExExpensedImmediately},
		{ID: "opt_depreciation_expense", Text: "Depreciation Expense", ErrorTag: engine.TagM10CapExExpensedImmediately},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagWrongAccount},
	}, seed, 1021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When forecasting long-term balance sheet capacity expansion to support projected sales growth (Module 10), which property, plant, and equipment asset account captures capital expenditures (CapEx)?",
		CorrectOptionID:   "opt_equipment",
		Options:           stage1Opts,
		CausalHint:        "Which physical asset account captures long-term productive machinery and equipment acquired to expand operating capacity?",
		Explanation:       "Equipment is a non-current asset account for capital expenditures that reflects newly acquired productive capacity on the pro-forma balance sheet.",
		RelevantConceptID: "forecasting_capex",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagM10CapExExpensedImmediately},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Equipment?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does long-lived productive machinery represent an economic asset resource, an operating expense, or a liability?",
		Explanation:       "Equipment is a non-current physical Asset that provides future productive economic benefits over multiple periods.",
		RelevantConceptID: "forecasting_capex",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 1023)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage3Opts,
		CausalHint:        "Assets have normal debit balances and increase on the debit side. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1024)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does acquiring additional machinery increase or decrease the Equipment balance on the projected balance sheet?",
		CorrectOptionID:   "opt_increase",
		Options:           stage4Opts,
		CausalHint:        "When new physical equipment is purchased, does total productive equipment capacity increase or decrease?",
		Explanation:       "Equipment increases as new productive assets are capitalized into property, plant, and equipment.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_rent_expense", Text: "Rent Expense (Expense)", ErrorTag: engine.TagM10CapExExpensedImmediately},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 1025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "When the projected equipment purchase is funded from corporate cash reserves, which liquid asset account is credited?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Which asset account is credited to reflect the cash outflow disbursed for capital equipment acquisitions?",
		Explanation:       "Cash is credited for the cash outflow used to acquire the productive machinery.",
		RelevantConceptID: "forecasting_capex",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_equip_cr_cash", Text: fmt.Sprintf("Debit Equipment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_equip", Text: fmt.Sprintf("Debit Cash %s / Credit Equipment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM10CapExExpensedImmediately},
	}, seed, 1026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry modeling projected capital expenditures (CapEx) for equipment in pro-forma statements?",
		CorrectOptionID:   "opt_dr_equip_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the Equipment asset debit with the Cash credit?",
		Explanation:       fmt.Sprintf("Debit Equipment %s and Credit Cash %s. This capitalizes the long-term productive machinery on the balance sheet rather than immediately expensing the cash outflow.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_capex",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_eq_capex", Text: fmt.Sprintf("Asset swap: Long-term assets increase (+Equipment %s) and liquid assets decrease (-Cash %s); Liabilities and Equity unchanged. On the pro-forma cash flow statement, this capital expenditure (CapEx) appears as an investing cash outflow, directly linking balance sheet expansion to the cash budget.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_eq_liab_up", Text: fmt.Sprintf("Liabilities increase by %s and equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_eq_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s and equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 1027)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this projected CapEx entry affect the balance sheet equation, and how does it connect to the cash budget?",
		CorrectOptionID:   "opt_eq_capex",
		Options:           stage7Opts,
		CausalHint:        "Equipment increases while Cash decreases by an equal amount. How does this asset exchange flow through the projected statement of cash flows?",
		Explanation:       fmt.Sprintf("Asset swap: Equipment increases (+Equipment %s) and Cash decreases (-Cash %s); Liabilities and Equity are unchanged. Projecting capital expenditures (CapEx) ensures productive capacity keeps pace with forecasted sales growth and flows directly into investing cash outflows on the pro-forma cash flow statement.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_capex",
	}
}

// BuildM10ForecastingDebtPlugStages builds 7 scaffolding stages for pro-forma External Financing Needed (EFN) debt plugs.
func BuildM10ForecastingDebtPlugStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagM10EFNTreatedAsEquityWithoutShares},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 1031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "In pro-forma balance sheet modeling (Module 10), when projected asset expansion exceeds projected spontaneous financing, creating an External Financing Needed (EFN) shortfall funded via long-term bank notes, which asset account receives the borrowed proceeds?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "Which liquid asset account receives the borrowed debt proceeds to fund projected asset expansion?",
		Explanation:       "Cash is the liquid asset account that reflects the financing proceeds obtained to fund the projected asset growth.",
		RelevantConceptID: "forecasting_debt_plug",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 1032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does cash and bank deposits represent an economic resource owned by the company?",
		Explanation:       "Cash is a Current Asset representing liquid purchasing power available to the firm.",
		RelevantConceptID: "forecasting_debt_plug",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 1033)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage3Opts,
		CausalHint:        "Assets have normal debit balances and increase on the debit side. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1034)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does receiving borrowed loan proceeds increase or decrease Cash on the projected balance sheet?",
		CorrectOptionID:   "opt_increase",
		Options:           stage4Opts,
		CausalHint:        "When funds are borrowed to plug external financing needs, does the cash balance increase or decrease?",
		Explanation:       "Cash increases as borrowed loan funds are deposited into corporate bank accounts.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagM10EFNTreatedAsEquityWithoutShares},
		{ID: "opt_retained_earnings", Text: "Retained Earnings (Equity)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagWrongAccount},
	}, seed, 1035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which long-term debt liability account is credited to model the External Financing Needed (EFN) borrowing plug?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage5Opts,
		CausalHint:        "Which contractual debt liability account records the formal loan obligation owed to lenders?",
		Explanation:       "Notes Payable is credited to record the formal loan obligation incurred to balance projected assets with projected liabilities and equity.",
		RelevantConceptID: "forecasting_debt_plug",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_notes_cr_cash", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cash_cr_stock", Text: fmt.Sprintf("Debit Cash %s / Credit Common Stock %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM10EFNTreatedAsEquityWithoutShares},
	}, seed, 1036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry modeling the External Financing Needed (EFN) debt plug in pro-forma statements?",
		CorrectOptionID:   "opt_dr_cash_cr_notes",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the Cash asset debit with the Notes Payable liability credit?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Notes Payable %s. This models the cash inflow and the corresponding debt obligation required to fund projected operations.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_debt_plug",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_eq_efn", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Notes Payable); Equity unchanged. In pro-forma financial statement forecasting, the External Financing Needed (EFN) debt plug reconciles projected asset requirements with projected funding, linking the projected balance sheet to financing cash inflows.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_eq_swap", Text: fmt.Sprintf("Asset swap: Cash increases and receivables decrease by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_eq_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s and equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 1037)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this projected debt plug affect the balance sheet equation, and how does it reconcile the financial model?",
		CorrectOptionID:   "opt_eq_efn",
		Options:           stage7Opts,
		CausalHint:        "Cash increases on the asset side and notes payable increases on the liability side. How does this plug balance the pro-forma balance sheet?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Liabilities increase (+Notes Payable %s); Equity is unchanged. In pro-forma financial statement forecasting, the External Financing Needed (EFN) debt plug reconciles projected asset requirements with projected funding, linking the projected balance sheet to financing cash inflows.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "forecasting_debt_plug",
	}
}

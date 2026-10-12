package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM03OCFAccountsPayableIncreaseStages builds the 7 scaffolding stages for operating A/P increase.
func BuildM03OCFAccountsPayableIncreaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ap", Text: "Accounts Payable"},
		{ID: "opt_prepaid", Text: "Prepaid Insurance", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_unearned", Text: "Unearned Revenue", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 3011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which liability account tracks trade obligations owed to suppliers for operating services purchased on credit?",
		CorrectOptionID:   "opt_ap",
		Options:           stage1Opts,
		CausalHint:        "When operating goods or services are received without immediate cash payment, which account records the debt owed to suppliers?",
		Explanation:       "Accounts Payable represents obligations to pay suppliers for operating goods or services previously received.",
		RelevantConceptID: "accounts_payable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability (Current Liability)"},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 3012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Payable?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does an obligation to settle unpaid bills with vendors represent an Asset, a Liability, or an Equity claim?",
		Explanation:       "Accounts Payable is a Current Liability representing short-term claims by trade suppliers.",
		RelevantConceptID: "accounts_payable_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 3013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Payable increase or decrease when operating expenses are incurred on account?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does purchasing services without paying cash right away add to or reduce the total amount owed to trade suppliers?",
		Explanation:       "Accounts Payable increases as new supplier invoices are incurred on credit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 3014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to a Liability (Accounts Payable) recorded?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have a normal credit balance; which side records an increase?",
		Explanation:       "Liabilities increase by Credit (right side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_rent_expense", Text: "Rent Expense (Expense)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 3015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account records the cost of the operating facility used during the period?",
		CorrectOptionID:   "opt_rent_expense",
		Options:           stage5Opts,
		CausalHint:        "The company utilized leased operating space this period. Which account reflects this expired operating cost?",
		Explanation:       "Rent Expense is debited because the operating service was consumed during this period.",
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_rent_cr_ap", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_dr_ap_cr_rent", Text: fmt.Sprintf("Debit Accounts Payable %s / Credit Rent Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 3016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for incurring operating rent on account?",
		CorrectOptionID:   "opt_dr_rent_cr_ap",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the operating expense incurred while reflecting the unpaid supplier liability?",
		Explanation:       fmt.Sprintf("Debit Rent Expense %s to recognize the period cost and Credit Accounts Payable %s to record the liability to the landlord.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_add_ap", Text: fmt.Sprintf("Liabilities increase (+A/P %s); Equity decreases (-Rent Expense %s). On the indirect Statement of Cash Flows, the increase in Accounts Payable is ADDED to Net Income.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_deduct_ap", Text: fmt.Sprintf("Liabilities increase (+A/P %s); Equity decreases (-Rent Expense %s). On the indirect Statement of Cash Flows, the increase in Accounts Payable is DEDUCTED from Net Income.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM03OCFPayableWrongDirection},
		{ID: "opt_cash_down", Text: fmt.Sprintf("Assets decrease (-Cash %s); Equity decreases (-Rent Expense %s). Operating cash flows decrease immediately.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 3017)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this transaction affect the balance sheet equation, and how is the increase in Accounts Payable treated on the indirect Statement of Cash Flows?",
		CorrectOptionID:   "opt_add_ap",
		Options:           stage7Opts,
		CausalHint:        "Net Income was reduced by this rent expense, but was any cash paid out this period? To reconcile Net Income to actual cash flow, should the change in payables be added back or subtracted?",
		Explanation:       fmt.Sprintf("Liabilities increase (+Accounts Payable %s) and Equity decreases (-Rent Expense %s). Because Net Income was reduced by the expense but no cash was disbursed, the increase in Accounts Payable must be ADDED back to Net Income in operating cash flows.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "statement_of_cash_flows_indirect",
	}
}

// BuildM03OCFPrepaidIncreaseStages builds the 7 scaffolding stages for prepaid asset increase.
func BuildM03OCFPrepaidIncreaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_prepaid", Text: "Prepaid Insurance"},
		{ID: "opt_expense", Text: "Insurance Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_unearned", Text: "Unearned Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 3021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records operating benefits paid for in cash today that will be consumed in future periods?",
		CorrectOptionID:   "opt_prepaid",
		Options:           stage1Opts,
		CausalHint:        "When cash is spent today for operating coverage that extends across future periods, which asset account stores this future benefit?",
		Explanation:       "Prepaid Insurance is an asset representing future insurance protection paid for in advance.",
		RelevantConceptID: "prepaid_expenses_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset (Current Asset)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_liability", Text: "Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 3022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Prepaid Insurance?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does the right to receive future insurance protection without additional payments qualify as an economic resource owned by the company?",
		Explanation:       "Prepaid Insurance is a Current Asset because it provides future economic benefits within the operating cycle.",
		RelevantConceptID: "prepaid_expenses_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 3023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Prepaid Insurance increase or decrease when paying policy premiums in advance?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does acquiring an additional period of advance coverage add to the company's prepaid economic resources?",
		Explanation:       "Prepaid Insurance increases as new coverage is purchased before consumption.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 3024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Prepaid Insurance) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have a normal debit balance; which side records an increase?",
		Explanation:       "Assets increase by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_ins_expense", Text: "Insurance Expense (Expense)", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_ap", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 3025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account decreases when cash is disbursed to the insurance underwriter?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid asset was paid out to acquire the insurance contract?",
		Explanation:       "Cash is credited because funds were paid out today.",
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_prepaid_cr_cash", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_cash", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_dr_cash_cr_prepaid", Text: fmt.Sprintf("Debit Cash %s / Credit Prepaid Insurance %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 3026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying insurance premiums in advance?",
		CorrectOptionID:   "opt_dr_prepaid_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry capitalizes the future benefit as an asset rather than prematurely expensing it, while reflecting the cash outflow?",
		Explanation:       fmt.Sprintf("Debit Prepaid Insurance %s to record the unexpired resource and Credit Cash %s for the cash disbursement.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_deduct_prepaid", Text: fmt.Sprintf("Total Assets are unchanged (+Prepaid Insurance %s, -Cash %s). On the indirect Statement of Cash Flows, the increase in Prepaid Insurance is DEDUCTED from Net Income.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_add_prepaid", Text: fmt.Sprintf("Total Assets are unchanged (+Prepaid Insurance %s, -Cash %s); on the indirect Statement of Cash Flows, the increase in Prepaid Insurance is ADDED to Net Income.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM03OCFPrepaidWrongDirection},
		{ID: "opt_asset_down", Text: fmt.Sprintf("Assets decrease by %s and Equity decreases by %s; operating cash flow is not adjusted.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
	}, seed, 3027)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this payment affect total assets, and how is the increase in Prepaid Insurance treated on the indirect Statement of Cash Flows?",
		CorrectOptionID:   "opt_deduct_prepaid",
		Options:           stage7Opts,
		CausalHint:        "Cash was spent, but was any expense deducted in Net Income this period? To reflect this cash outflow in operating activities under the indirect method, do you add or deduct the increase in prepaids?",
		Explanation:       fmt.Sprintf("Assets change composition (+Prepaid Insurance %s, -Cash %s) with no immediate net income effect. Because cash was spent but no expense was deducted on the Income Statement, the increase in Prepaid Insurance must be DEDUCTED from Net Income in the operating section.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "statement_of_cash_flows_indirect",
	}
}

// BuildM03OCFUnearnedRevenueIncreaseStages builds the 7 scaffolding stages for customer advance retainer increase.
func BuildM03OCFUnearnedRevenueIncreaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned", Text: "Unearned Revenue"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 3031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records customer cash received in advance of performing operating services?",
		CorrectOptionID:   "opt_unearned",
		Options:           stage1Opts,
		CausalHint:        "When cash is received from a client before the job is done, which account records the performance obligation owed to that client?",
		Explanation:       "Unearned Revenue represents the liability obligation to perform services or return the customer's funds.",
		RelevantConceptID: "unearned_revenue_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability (Current Liability)"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 3032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Unearned Revenue?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does an unperformed obligation to deliver services to a customer represent Revenue, an Asset, or a Liability?",
		Explanation:       "Unearned Revenue is a Current Liability representing a performance obligation owed to a customer.",
		RelevantConceptID: "unearned_revenue_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 3033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Unearned Revenue increase or decrease when receiving a client retainer in advance?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does receiving customer cash for future work add to the company's remaining performance obligations?",
		Explanation:       "Unearned Revenue increases when advance client payments are accepted.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 3034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to a Liability (Unearned Revenue) recorded?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have a normal credit balance; which side records an increase?",
		Explanation:       "Liabilities increase by Credit (right side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_ar", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 3035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liquid asset account increases when the customer's advance payment is received?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What tangible asset came into the company's bank account when the customer paid?",
		Explanation:       "Cash is debited because liquid funds were received today.",
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_dr_unearned_cr_cash", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 3036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for collecting a customer retainer in advance?",
		CorrectOptionID:   "opt_dr_cash_cr_unearned",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the cash inflow while deferring revenue until the performance obligation is met?",
		Explanation:       fmt.Sprintf("Debit Cash %s to record the receipt of money and Credit Unearned Revenue %s to record the performance liability.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "operating_cash_flow_adjustments",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_add_unearned", Text: fmt.Sprintf("Assets increase (+Cash %s); Liabilities increase (+Unearned Revenue %s). On the indirect Statement of Cash Flows, the increase in Unearned Revenue is ADDED to Net Income.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_deduct_unearned", Text: fmt.Sprintf("Assets increase (+Cash %s); Liabilities increase (+Unearned Revenue %s). On the indirect Statement of Cash Flows, the increase in Unearned Revenue is DEDUCTED from Net Income.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM03OCFUnearnedRevenueWrongDirection},
		{ID: "opt_rev_earned", Text: fmt.Sprintf("Assets increase (+Cash %s); Equity increases (+Revenue %s). Operating cash flows require no adjustment.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecognizedPrematurely},
	}, seed, 3037)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does this transaction affect the balance sheet equation, and how is the change in Unearned Revenue treated on the indirect Statement of Cash Flows?",
		CorrectOptionID:   "opt_add_unearned",
		Options:           stage7Opts,
		CausalHint:        "Operating cash was collected today, but was any revenue included in Net Income yet? Under the indirect method, should the increase in this operating liability be added or subtracted to reach operating cash flow?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Liabilities increase (+Unearned Revenue %s). Because cash was collected from customers but not yet included in Net Income, the increase in Unearned Revenue must be ADDED to Net Income on the indirect Statement of Cash Flows.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "statement_of_cash_flows_indirect",
	}
}

// BuildM03OCFQuickRatioStages builds the 7 scaffolding stages for settling payables in cash and Quick Ratio analysis.
func BuildM03OCFQuickRatioStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ap", Text: "Accounts Payable"},
		{ID: "opt_prepaid", Text: "Prepaid Insurance", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 3041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which liability account is reduced when cash is disbursed to satisfy outstanding supplier balances?",
		CorrectOptionID:   "opt_ap",
		Options:           stage1Opts,
		CausalHint:        "Which trade liability account records obligations owed to suppliers for past credit purchases?",
		Explanation:       "Accounts Payable is the liability account representing debts owed to trade suppliers.",
		RelevantConceptID: "accounts_payable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability (Current Liability)"},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 3042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Payable?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Is an obligation owed to trade suppliers an Asset, a Liability, or an Equity account?",
		Explanation:       "Accounts Payable is a Current Liability and represents a short-term obligation included in the denominator of both Current and Quick ratios.",
		RelevantConceptID: "quick_ratio_liquidity",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase", ErrorTag: engine.TagReversedSides},
	}, seed, 3043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Payable increase or decrease when the company pays supplier invoices?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "Does sending payments to vendors reduce the total debt owed to them?",
		Explanation:       "Accounts Payable decreases as the company satisfies its vendor obligations.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 3044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is a decrease to a Liability (Accounts Payable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities increase on the credit side; which side records a reduction?",
		Explanation:       "Liabilities decrease by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_ar", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_prepaid", Text: "Prepaid Insurance (Asset)", ErrorTag: engine.TagWrongAccount},
	}, seed, 3045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which quick asset account decreases when supplier bills are paid?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid asset is disbursed from the bank account to pay suppliers?",
		Explanation:       "Cash is credited because liquid funds were spent to pay the suppliers.",
		RelevantConceptID: "quick_ratio_liquidity",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ap_cr_cash", Text: fmt.Sprintf("Debit Accounts Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_ap", Text: fmt.Sprintf("Debit Cash %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
	}, seed, 3046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying Accounts Payable in cash?",
		CorrectOptionID:   "opt_dr_ap_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry reduces the supplier liability while recording the outflow of liquid cash?",
		Explanation:       fmt.Sprintf("Debit Accounts Payable %s to satisfy the obligation and Credit Cash %s to record the disbursement.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "quick_ratio_liquidity",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_qr_increase", Text: "The Quick Ratio increases (subtracting equal dollars from numerator and denominator when ratio > 1.0 increases the ratio). In contrast, buying inventory with cash reduces the Quick Ratio because Inventory is excluded from Quick Assets."},
		{ID: "opt_qr_decrease", Text: "The Quick Ratio decreases because liquid cash was spent, just as when inventory is purchased.", ErrorTag: engine.TagM03QuickRatioWrongDirection},
		{ID: "opt_qr_no_change", Text: "The Quick Ratio is unchanged because Cash and Accounts Payable offset, and inventory purchases also have no effect on the Quick Ratio.", ErrorTag: engine.TagM03InventoryTreatedAsQuickAsset},
	}, seed, 3047)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "Assuming the initial Quick Ratio exceeds 1.0 (e.g., 1.5), how does paying Accounts Payable in cash affect the Quick Ratio, and how does this contrast with purchasing inventory with cash?",
		CorrectOptionID:   "opt_qr_increase",
		Options:           stage7Opts,
		CausalHint:        "Quick Ratio = (Cash + Marketable Securities + Receivables) / Current Liabilities. When ratio > 1.0, what happens mathematically when you subtract equal amounts from numerator and denominator, and why is inventory excluded from quick assets?",
		Explanation:       fmt.Sprintf("Paying Accounts Payable subtracts equal dollars from Quick Assets (-Cash %s) and Current Liabilities (-A/P %s). For any ratio > 1.0 (e.g. 150/100 = 1.5 -> 130/80 = 1.625), subtracting equal amounts INCREASES the ratio. In contrast, purchasing inventory with cash converts Cash (a quick asset) into Inventory (a non-quick asset), reducing Quick Assets without changing Current Liabilities, which causes the Quick Ratio to decrease.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "quick_ratio_liquidity",
	}
}

// BuildM03Stages provides a dispatch helper for Stream-M03 event families.
func BuildM03Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM03OCFAccountsPayableIncrease:
		BuildM03OCFAccountsPayableIncreaseStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM03OCFPrepaidIncrease:
		BuildM03OCFPrepaidIncreaseStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM03OCFUnearnedRevenueIncrease:
		BuildM03OCFUnearnedRevenueIncreaseStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM03OCFQuickRatio:
		BuildM03OCFQuickRatioStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

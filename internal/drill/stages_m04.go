package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// buildM04VariableConsiderationStages builds the 7 drill stages for ASC 606 Variable Consideration (Constraint on Variable Consideration).
func buildM04VariableConsiderationStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_unearned", Text: "Unearned Revenue", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 4011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When the firm finishes services and bills the client on account for the fixed fee while an uncertain incentive bonus is constrained under ASC 606, which primary account is adjusted for the uncollected claim?",
		CorrectOptionID:   "opt_ar",
		Options:           stage1Opts,
		CausalHint:        "Has the customer already transferred liquid funds today, or does the firm hold an unpaid claim on account for completed work?",
		Explanation:       "Accounts Receivable records the unpaid claim for services performed where the customer has been billed on account.",
		RelevantConceptID: "variable_consideration_constraint",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_liability", Text: "Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 4012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What classification of account is Accounts Receivable?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does an unpaid customer invoice represent an economic resource with future cash inflows, or an obligation owed to someone else?",
		Explanation:       "Accounts Receivable is an Asset representing a legally enforceable claim to receive cash in the future.",
		RelevantConceptID: "account_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 4013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Receivable increase or decrease when the firm bills the customer for completed services?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "As newly billed amounts are added to client balances, does total money owed to the firm increase or decrease?",
		Explanation:       "Accounts Receivable increases because a new customer claim has been established.",
		RelevantConceptID: "direction_of_change",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 4014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account such as Accounts Receivable?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have normal debit balances on the left side of the ledger. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_revenue", Text: "Service Revenue"},
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Liability)", ErrorTag: engine.TagRevenueDeferredWhenEarned},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 4015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Under the ASC 606 variable consideration constraint, which earned account is credited for the fixed transaction price?",
		CorrectOptionID:   "opt_service_revenue",
		Options:           stage5Opts,
		CausalHint:        "Because the performance obligation was completed and the fixed fee is not subject to significant reversal, what account records earnings for the period?",
		Explanation:       "Service Revenue is credited because the performance obligation has been fulfilled and the fixed transaction price is unconstrained.",
		RelevantConceptID: "service_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 4016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the billing of the fixed fee under ASC 606?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the increase in the asset claim with the increase in earned revenue for the period?",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s increases the asset; Credit Service Revenue %s increases revenue and equity. The uncertain bonus remains excluded until meeting the constraint criteria.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "balanced_entry_mechanics",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Accounts Receivable); Equity increases by %s (+Service Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_no_effect", Text: "Total Assets and Equity remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 4017)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does billing the fixed fee on account affect the balance sheet equation?",
		CorrectOptionID:   "opt_assets_up_eq_up",
		Options:           stage7Opts,
		CausalHint:        "An asset increased by the billed amount and earned revenue increased net income. How does this impact total assets and stockholders' equity?",
		Explanation:       fmt.Sprintf("Assets increase (+Accounts Receivable %s) and Equity increases (+Service Revenue %s via net income). Liabilities remain unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// buildM04ContractAssetARStages builds the 7 drill stages for ASC 606 Contract Asset vs. Accounts Receivable (Unconditional Right vs. Conditional Right).
func buildM04ContractAssetARStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_unearned", Text: "Unearned Revenue", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 4021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under ASC 606-10-45, when an entity's right to payment is unconditional (dependent only upon the passage of time), which account is recognized for the billed claim?",
		CorrectOptionID:   "opt_ar",
		Options:           stage1Opts,
		CausalHint:        "When performance obligations are fully certified and payment is subject only to the passage of time, which standard asset account records unconditional client claims?",
		Explanation:       "Accounts Receivable is the required presentation under ASC 606 when an entity possesses an unconditional right to consideration, unlike a contract asset which depends on future performance.",
		RelevantConceptID: "contract_asset_vs_receivable",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 4022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Receivable under ASC 606 presentation rules?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does an unconditional legal right to collect payment from a customer constitute an owned economic resource or an obligation?",
		Explanation:       "Accounts Receivable is a Current Asset representing an unconditional right to consideration.",
		RelevantConceptID: "account_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 4023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Receivable increase or decrease when the firm satisfies its milestone and establishes an unconditional billing?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When a new unconditional claim is established on net 30 credit terms, does the total receivable balance increase or decrease?",
		Explanation:       "Accounts Receivable increases as new unconditional claims are established.",
		RelevantConceptID: "direction_of_change",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 4024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase the asset Accounts Receivable?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets carry normal debit balances. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_revenue", Text: "Service Revenue"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 4025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "When an unconditional receivable is recorded upon fulfilling performance obligations under ASC 606, which revenue account is credited?",
		CorrectOptionID:   "opt_service_revenue",
		Options:           stage5Opts,
		CausalHint:        "What account reflects the value of completed customer services earned during the operating period?",
		Explanation:       "Service Revenue is credited to recognize earnings upon the satisfaction of performance obligations.",
		RelevantConceptID: "service_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 4026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the unconditional billing under ASC 606?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry increases the unconditional asset claim and increases earned revenue by identical amounts?",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s increases assets; Credit Service Revenue %s increases revenue and equity. Under ASC 606-10-45, this unconditional claim must be presented as a receivable rather than a contract asset.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "balanced_entry_mechanics",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Accounts Receivable); Equity increases by %s (+Service Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_no_effect", Text: "Total Assets and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 4027)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does recording the unconditional billing affect the fundamental balance sheet equation?",
		CorrectOptionID:   "opt_assets_up_eq_up",
		Options:           stage7Opts,
		CausalHint:        "Receivables increased total assets, while service revenue flowed into retained earnings. How does this shift total assets and equity?",
		Explanation:       fmt.Sprintf("Assets increase by %s (+Accounts Receivable) and Equity increases by %s (+Service Revenue via retained earnings), preserving equation balance.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// buildM04RightOfReturnStages builds the 7 drill stages for ASC 606 Right of Return / Refund Obligation (Unearned Revenue).
func buildM04RightOfReturnStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagM04RevenueRecognizedWithUncertainReturn},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 4031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When a customer pays cash in advance for an advisory package that includes a 30-day full money-back satisfaction guarantee prior to service delivery, which liquid asset account is adjusted for the inflow?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What account directly tracks the company's checking account balance and immediately available legal funds?",
		Explanation:       "Cash records the immediate inflow of liquid currency received from the customer.",
		RelevantConceptID: "cash_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 4032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What classification of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does money held in the bank represent an owned economic resource or an obligation owed to an outside party?",
		Explanation:       "Cash is a Current Asset representing the most liquid store of economic purchasing power.",
		RelevantConceptID: "account_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease", ErrorTag: engine.TagReversedSides},
	}, seed, 4033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon receiving the customer's upfront payment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When a customer transfers money into the corporate checking account, does total cash increase or decrease?",
		Explanation:       "Cash increases as new funds are deposited into the company's account.",
		RelevantConceptID: "direction_of_change",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 4034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account such as Cash?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets reside on the left side of the fundamental accounting equation. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Liability)"},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagM04RevenueRecognizedWithUncertainReturn},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 4035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Under ASC 606, because services have not yet been performed and the payment is subject to full return/refund terms, which contract liability account is credited?",
		CorrectOptionID:   "opt_unearned_revenue",
		Options:           stage5Opts,
		CausalHint:        "Because the company owes future performance or a complete refund to the customer, what liability account captures unperformed obligations?",
		Explanation:       "Unearned Revenue is credited because the company has a contract liability to provide service or return the consideration.",
		RelevantConceptID: "unearned_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM04RevenueRecognizedWithUncertainReturn},
		{ID: "opt_dr_unearned_cr_cash", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 4036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the customer advance subject to right-of-return terms?",
		CorrectOptionID:   "opt_dr_cash_cr_unearned",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the debit for the liquid asset received with the credit for the contract liability obligation?",
		Explanation:       fmt.Sprintf("Debit Cash %s increases assets; Credit Unearned Revenue %s increases contract liabilities. Revenue cannot be recognized before performance obligations are satisfied and while return uncertainty remains open.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "balanced_entry_mechanics",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_no_effect", Text: "Total Assets and Liabilities remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 4037)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does receiving this upfront advance subject to refund terms affect the accounting equation?",
		CorrectOptionID:   "opt_assets_up_liab_up",
		Options:           stage7Opts,
		CausalHint:        "Cash was collected, creating an obligation to perform or refund. What happens to total assets and total liabilities?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Liabilities increase (+Unearned Revenue %s). Equity is completely unaffected until performance occurs.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// buildM04SatisfactionOverTimeStages builds the 7 drill stages for ASC 606 Step 5 (Performance Obligation Satisfied Over Time).
func buildM04SatisfactionOverTimeStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagM04CashDebitedOnOverTimeSatisfaction},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 4041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under ASC 606 Step 5, when an entity delivers consulting services over time to fulfill a performance obligation previously paid in advance, which existing liability account is adjusted first?",
		CorrectOptionID:   "opt_unearned_revenue",
		Options:           stage1Opts,
		CausalHint:        "What existing liability account was established when the client previously deposited funds before work began?",
		Explanation:       "Unearned Revenue represents the contract liability that must be reduced as the firm satisfies its performance obligation over time.",
		RelevantConceptID: "contract_liability_satisfaction",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_equity", Text: "Equity", ErrorTag: engine.TagWrongAccount},
	}, seed, 4042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Unearned Revenue?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does receiving customer payments in advance create an obligation to perform future services, or an asset?",
		Explanation:       "Unearned Revenue is a Liability representing a performance obligation owed to a customer.",
		RelevantConceptID: "account_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase", ErrorTag: engine.TagReversedSides},
	}, seed, 4043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Unearned Revenue increase or decrease when the firm satisfies its performance obligation over time?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "As promised services are progressively performed and delivered, does the firm's remaining obligation to the client grow or shrink?",
		Explanation:       "Unearned Revenue decreases because the obligation has been partially or fully satisfied through performance.",
		RelevantConceptID: "direction_of_change",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 4044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to decrease a Liability account such as Unearned Revenue?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have normal credit balances. Which side reduces a liability?",
		Explanation:       "Liabilities decrease on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_revenue", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagM04CashDebitedOnOverTimeSatisfaction},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagWrongAccount},
	}, seed, 4045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "When the contract liability is relieved upon delivering services over time, which account is credited to recognize earnings?",
		CorrectOptionID:   "opt_service_revenue",
		Options:           stage5Opts,
		CausalHint:        "What account captures the value of customer services performed and earned during the reporting period?",
		Explanation:       "Service Revenue is credited to record revenue earned as performance obligations are satisfied.",
		RelevantConceptID: "service_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_unearned_cr_rev", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM04CashDebitedOnOverTimeSatisfaction},
		{ID: "opt_dr_rev_cr_unearned", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 4046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record satisfying the performance obligation over time?",
		CorrectOptionID:   "opt_dr_unearned_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Which entry pairs the debit reducing the contract liability with the credit recognizing earned revenue?",
		Explanation:       fmt.Sprintf("Debit Unearned Revenue %s decreases liabilities; Credit Service Revenue %s increases revenue and equity. Under ASC 606 Step 5, fulfilling the obligation over time converts the contract liability into earned revenue.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "balanced_entry_mechanics",
	}

	stage7Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liab_down_eq_up", Text: fmt.Sprintf("Liabilities decrease by %s (-Unearned Revenue); Equity increases by %s (+Service Revenue); Assets unchanged.", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		{ID: "opt_no_effect", Text: "Total Liabilities and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 4047)

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:             domain.StageEquationEffect,
		Prompt:            "How does satisfying this performance obligation over time affect the balance sheet equation?",
		CorrectOptionID:   "opt_liab_down_eq_up",
		Options:           stage7Opts,
		CausalHint:        "A contract liability was fulfilled without exchanging cash, while revenue increased equity. How do total liabilities and equity change?",
		Explanation:       fmt.Sprintf("Liabilities decrease (-Unearned Revenue %s) and Equity increases (+Service Revenue %s via net income). Assets remain completely unaffected.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// Exported stage builder aliases for test package access.
var (
	ExportBuildM04VariableConsiderationStages = buildM04VariableConsiderationStages
	ExportBuildM04ContractAssetARStages       = buildM04ContractAssetARStages
	ExportBuildM04RightOfReturnStages         = buildM04RightOfReturnStages
	ExportBuildM04SatisfactionOverTimeStages  = buildM04SatisfactionOverTimeStages
)

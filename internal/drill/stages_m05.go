package drill

import (
	"fmt"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// BuildM05Stages dispatches to the appropriate modular stage builder for Stream-M05.
func BuildM05Stages(familyID string, inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) bool {
	switch familyID {
	case engine.FamilyM05InventoryFIFO:
		BuildM05InventoryFIFOStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM05InventoryLIFO:
		BuildM05InventoryLIFOStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM05InventoryNRVWriteDown:
		BuildM05InventoryNRVStages(inst, amt, seed, proc)
		return true
	case engine.FamilyM05FreightInCapitalize:
		BuildM05FreightInCapStages(inst, amt, seed, proc)
		return true
	default:
		return false
	}
}

// BuildM05InventoryFIFOStages builds 7 scaffolding stages for perpetual FIFO inventory costing.
func BuildM05InventoryFIFOStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cogs", Text: "Cost of Goods Sold"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagM05InventoryRecordedAtSelling},
		{ID: "opt_inventory", Text: "Inventory", ErrorTag: engine.TagWrongAccount},
	}, seed, 5011)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under a perpetual inventory system using FIFO, which expense account records the cost of inventory sold and delivered to customers?",
		CorrectOptionID:   "opt_cogs",
		Options:           stage1Opts,
		CausalHint:        "When merchandise is transferred to buyers, which primary expense captures the historical acquisition cost of the sold goods?",
		Explanation:       "Cost of Goods Sold is the expense account that records the historical acquisition cost of inventory sold to customers to generate sales revenue.",
		RelevantConceptID: "cogs_recognition",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 5012)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cost of Goods Sold?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the cost of merchandise sold reduce net earnings for the accounting period, or represent an unexpired resource?",
		Explanation:       "Cost of Goods Sold is an operating Expense account on the income statement that directly reduces Gross Profit and Net Income.",
		RelevantConceptID: "cogs_recognition",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 5013)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cost of Goods Sold increase or decrease as merchandise is delivered to the customer?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When a new sales transaction occurs, does total expired cost for the period accumulate upward or downward?",
		Explanation:       "Cost of Goods Sold increases because additional product cost has expired during the accounting period.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 5014)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses reduce stockholders' equity and carry normal debit balances. Do they increase on the left or the right side?",
		Explanation:       "Expenses increase by Debit on the left side, reflecting a decrease in stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_inventory", Text: "Inventory (Asset)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 5015)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which balance sheet asset account is reduced to reflect the departure of merchandise from the warehouse?",
		CorrectOptionID:   "opt_inventory",
		Options:           stage5Opts,
		CausalHint:        "Which current asset stored the goods at acquisition cost prior to sale and must now be reduced?",
		Explanation:       "Inventory is credited to remove the historical cost of items that have been shipped and sold.",
		RelevantConceptID: "inventory_costing_fifo",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cogs_cr_inv", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Inventory %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_inv_cr_cogs", Text: fmt.Sprintf("Debit Inventory %s / Credit Cost of Goods Sold %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cogs_cr_cash", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 5016)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the FIFO cost of goods sold?",
		CorrectOptionID:   "opt_dr_cogs_cr_inv",
		Options:           stage6Opts,
		CausalHint:        "Which entry matches the expired product cost to expense while crediting the asset account holding the merchandise?",
		Explanation:       fmt.Sprintf("Debit Cost of Goods Sold %s and Credit Inventory %s. Under FIFO, the oldest cost layers are expensed first, leaving newer costs in ending inventory.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "inventory_costing_fifo",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording FIFO cost of goods sold affect the fundamental balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Inventory); Equity decreases by %s (-Cost of Goods Sold); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 5017),
		CausalHint:        "Inventory decreased on the balance sheet and Cost of Goods Sold reduced net income. How does this affect Assets and Equity?",
		Explanation:       fmt.Sprintf("Assets decrease (-Inventory %s) and Equity decreases (-Cost of Goods Sold %s); Liabilities are unchanged. In periods of rising prices, FIFO yields lower COGS and higher ending inventory compared to LIFO.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM05InventoryLIFOStages builds 7 scaffolding stages for perpetual LIFO inventory costing.
func BuildM05InventoryLIFOStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cogs", Text: "Cost of Goods Sold"},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagM05InventoryRecordedAtSelling},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable", ErrorTag: engine.TagWrongAccount},
	}, seed, 5021)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under a perpetual inventory system using LIFO, which expense account records the most recent inventory acquisition costs matched against revenue?",
		CorrectOptionID:   "opt_cogs",
		Options:           stage1Opts,
		CausalHint:        "When merchandise is sold under LIFO, which expense account captures the most recent batch costs matched against sales revenue?",
		Explanation:       "Cost of Goods Sold records the expired cost of inventory sold under LIFO, matching the latest purchase costs against current revenue.",
		RelevantConceptID: "cogs_recognition",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 5022)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cost of Goods Sold?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the cost of inventory sold represent an operating expense on the income statement, or a future economic resource on the balance sheet?",
		Explanation:       "Cost of Goods Sold is an Expense account that reduces net operating income on the income statement.",
		RelevantConceptID: "cogs_recognition",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 5023)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cost of Goods Sold increase or decrease when recording this LIFO sales transaction?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing cost of goods sold for a new sale transaction increase or decrease total expenses for the period?",
		Explanation:       "Cost of Goods Sold increases as the cost of the delivered inventory is recognized in earnings.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 5024)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expense accounts increase on the normal balance side that reduces equity. Do expenses increase with a debit on the left or a credit on the right?",
		Explanation:       "Expenses increase on the left side (Debit), reducing Stockholders' Equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_inventory", Text: "Inventory (Asset)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 5025)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is credited to remove the sold merchandise from the company's books?",
		CorrectOptionID:   "opt_inventory",
		Options:           stage5Opts,
		CausalHint:        "Which asset account is reduced when merchandise physically leaves the warehouse under a perpetual inventory tracking system?",
		Explanation:       "Inventory is credited to remove the carrying value of the merchandise sold.",
		RelevantConceptID: "inventory_costing_lifo",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cogs_cr_inv", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Inventory %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_inv_cr_cogs", Text: fmt.Sprintf("Debit Inventory %s / Credit Cost of Goods Sold %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cogs_cr_cash", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 5026)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the LIFO cost of goods sold?",
		CorrectOptionID:   "opt_dr_cogs_cr_inv",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry records the recent cost layers in Cost of Goods Sold and removes the sold units from the asset account?",
		Explanation:       fmt.Sprintf("Debit Cost of Goods Sold %s and Credit Inventory %s. During inflation, LIFO matches recent higher costs with current revenue, reducing pre-tax income and taxes.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "inventory_costing_lifo",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording LIFO cost of goods sold affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Inventory); Equity decreases by %s (-Cost of Goods Sold); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets, Liabilities, and Equity remain completely unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 5027),
		CausalHint:        "What is the net effect on the fundamental balance sheet equation when inventory is credited and cost of goods sold is debited?",
		Explanation:       fmt.Sprintf("Assets decrease (-Inventory %s) and Equity decreases (-Cost of Goods Sold %s); Liabilities are unchanged. Note: US GAAP allows LIFO under the LIFO Conformity Rule, but IFRS (IAS 2) strictly prohibits LIFO.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM05InventoryNRVStages builds 7 scaffolding stages for Lower of Cost or Net Realizable Value write-downs (ASC 330).
func BuildM05InventoryNRVStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cogs", Text: "Cost of Goods Sold"},
		{ID: "opt_bad_debt_expense", Text: "Bad Debt Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 5031)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under US GAAP (ASC 330), when inventory net realizable value falls below cost due to damage or obsolescence, which primary expense/cost account absorbs the write-down?",
		CorrectOptionID:   "opt_cogs",
		Options:           stage1Opts,
		CausalHint:        "Under US GAAP (ASC 330), when inventory utility drops below its recorded cost, which primary cost/expense account commonly absorbs this valuation write-down?",
		Explanation:       "Cost of Goods Sold (or a dedicated write-down loss account) absorbs the reduction in inventory value when net realizable value falls below historical cost.",
		RelevantConceptID: "lower_of_cost_or_nrv",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 5032)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cost of Goods Sold when absorbing an inventory write-down?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an inventory valuation write-down represent an expense/loss that reduces current period earnings, or a balance sheet liability?",
		Explanation:       "Cost of Goods Sold is an Expense account that reduces current period gross margin and net income.",
		RelevantConceptID: "lower_of_cost_or_nrv",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 5033)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cost of Goods Sold increase or decrease when recognizing this lower of cost or NRV impairment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing a write-down for impaired or obsolete inventory increase or decrease the total expense reported on the income statement?",
		Explanation:       "Cost of Goods Sold increases to absorb the valuation write-down loss.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 5034)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expense and loss accounts reduce stockholders' equity. On which side is an increase to an expense account entered?",
		Explanation:       "Expenses increase by Debit on the left side, reducing equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_inventory", Text: "Inventory (Asset)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 5035)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is directly credited to reduce the inventory carrying value to net realizable value?",
		CorrectOptionID:   "opt_inventory",
		Options:           stage5Opts,
		CausalHint:        "Which balance sheet asset account must be directly reduced to establish the new lower cost basis for impaired merchandise?",
		Explanation:       "Inventory is credited to reduce its carrying value down to its estimated net realizable value.",
		RelevantConceptID: "lower_of_cost_or_nrv",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cogs_cr_inv", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Inventory %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_inv_cr_cogs", Text: fmt.Sprintf("Debit Inventory %s / Credit Cost of Goods Sold %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_loss_cr_ap", Text: fmt.Sprintf("Debit Cost of Goods Sold %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 5036)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the lower of cost or NRV inventory write-down?",
		CorrectOptionID:   "opt_dr_cogs_cr_inv",
		Options:           stage6Opts,
		CausalHint:        "Which balanced journal entry recognizes the inventory write-down in earnings and reduces the carrying value of inventory?",
		Explanation:       fmt.Sprintf("Debit Cost of Goods Sold %s and Credit Inventory %s. Under US GAAP (ASC 330), writing down inventory establishes a new, lower cost basis.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "lower_of_cost_or_nrv",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording the lower of cost or NRV write-down affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Inventory); Equity decreases by %s (-Cost of Goods Sold); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: "Total Assets and Equity remain unchanged.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 5037),
		CausalHint:        "How does writing down inventory to its net realizable value affect the balance sheet equation (Assets, Liabilities, and Equity)?",
		Explanation:       fmt.Sprintf("Assets decrease (-Inventory %s) and Equity decreases (-Cost of Goods Sold %s); Liabilities are unchanged. Under US GAAP, write-downs cannot be reversed even if market prices recover subsequently.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

// BuildM05FreightInCapStages builds 7 scaffolding stages for Freight-In inventory capitalization (FOB Shipping Point).
func BuildM05FreightInCapStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if inst.StageAnswers == nil {
		inst.StageAnswers = make(map[domain.DrillStage]domain.StageAnswer)
	}

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_inventory", Text: "Inventory"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagM05FreightInExpensed},
	}, seed, 5041)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under US GAAP, when a purchaser pays freight charges to transport goods FOB shipping point to its warehouse, which asset account capitalizes these shipping costs?",
		CorrectOptionID:   "opt_inventory",
		Options:           stage1Opts,
		CausalHint:        "Under US GAAP, all costs necessary to get inventory ready and in place for resale are capitalized. Which asset account directly absorbs these inbound freight costs?",
		Explanation:       "Inventory is an asset account that capitalizes all acquisition and inbound transportation costs necessary to prepare goods for resale.",
		RelevantConceptID: "freight_in_capitalization",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 5042)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Inventory?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does capitalized merchandise inventory represent an economic resource owned by the company with future economic value, or an expired expense?",
		Explanation:       "Inventory is a Current Asset on the balance sheet representing goods held for sale in the normal course of business.",
		RelevantConceptID: "freight_in_capitalization",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 5043)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Inventory increase or decrease when capitalizing inbound transportation freight charges?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does adding freight-in shipping charges to the recorded value of inventory increase or decrease the asset account balance?",
		Explanation:       "Inventory increases because freight-in is an essential product cost capitalized into the inventory asset base.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 5044)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side of the ledger is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets maintain normal debit balances on the left side of the ledger. On which side do assets increase?",
		Explanation:       "Assets increase by Debit on the left side of the ledger.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagWrongAccount},
	}, seed, 5045)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is reduced when paying cash for the freight-in transportation service?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Which asset account is reduced when the company disburses currency to pay the transportation carrier for shipping services?",
		Explanation:       "Cash is credited because currency was paid to the shipping carrier.",
		RelevantConceptID: "freight_in_capitalization",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_inv_cr_cash", Text: fmt.Sprintf("Debit Inventory %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_inv", Text: fmt.Sprintf("Debit Cash %s / Credit Inventory %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_freight_cr_cash", Text: fmt.Sprintf("Debit Delivery Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagM05FreightInExpensed},
	}, seed, 5046)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the freight-in inventory capitalization?",
		CorrectOptionID:   "opt_dr_inv_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which journal entry capitalizes the freight-in costs into inventory and records the immediate cash payment?",
		Explanation:       fmt.Sprintf("Debit Inventory %s and Credit Cash %s. Inbound shipping costs (freight-in) attach to the product as an asset; outbound shipping to customers (freight-out) is a selling expense.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "freight_in_capitalization",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does capitalizing freight-in paid in cash affect the balance sheet equation?",
		CorrectOptionID: "opt_asset_swap_no_net_effect",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap_no_net_effect", Text: fmt.Sprintf("Assets increase by %s (+Inventory) and decrease by %s (-Cash); Total Assets, Liabilities, and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 5047),
		CausalHint:        "When cash is exchanged to increase the capitalized cost basis of inventory, what is the net effect on total assets, liabilities, and equity?",
		Explanation:       fmt.Sprintf("Total Assets remain unchanged due to an asset exchange (+Inventory %s, -Cash %s); Liabilities and Equity are unaffected. Incurring capitalized inventory costs does not hit the income statement until the goods are sold.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

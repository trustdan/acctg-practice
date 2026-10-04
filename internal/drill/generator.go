package drill

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// Generator builds reproducible, instantiated progressive drills from question templates.
type Generator struct {
	catalog *domain.AccountCatalog
	engine  *engine.Engine
}

func NewGenerator(catalog *domain.AccountCatalog, eng *engine.Engine) *Generator {
	return &Generator{
		catalog: catalog,
		engine:  eng,
	}
}

// GenerateInstanceWithScaffold instantiates a question template with an explicit scaffold fading level.
func (g *Generator) GenerateInstanceWithScaffold(
	q bank.QuestionJSON,
	seed int64,
	paramValues map[string]int64,
	scaffold domain.ScaffoldLevel,
) (*domain.QuestionInstance, error) {
	inst, err := g.GenerateInstance(q, seed, paramValues)
	if err != nil {
		return nil, err
	}
	inst.ScaffoldLevel = scaffold
	return inst, nil
}

// GenerateInstance instantiates a question template with a deterministic seed and parameter value.
func (g *Generator) GenerateInstance(q bank.QuestionJSON, seed int64, paramValues map[string]int64) (*domain.QuestionInstance, error) {
	amtVal, ok := paramValues["amount_minor_units"]
	if !ok || amtVal <= 0 {
		return nil, fmt.Errorf("missing or invalid amount_minor_units parameter")
	}
	money := domain.NewMoney(amtVal)

	// Replace scenario template placeholders
	promptText := q.ScenarioTemplate
	promptText = strings.ReplaceAll(promptText, "${amount_dollars}", money.FormatDollars())
	promptText = strings.ReplaceAll(promptText, "${amount_exact}", money.FormatExact())

	// Process canonical event in engine
	event := engine.TransactionEvent{
		FamilyID:   q.FamilyID,
		Parameters: paramValues,
	}
	processed, err := g.engine.ProcessEvent(event)
	if err != nil {
		return nil, fmt.Errorf("failed to process event in engine: %w", err)
	}

	instance := &domain.QuestionInstance{
		InstanceID:    fmt.Sprintf("%s-v%d-%d-%d", q.ID, q.Version, seed, amtVal),
		QuestionID:    q.ID,
		Version:       q.Version,
		FamilyID:      q.FamilyID,
		RuleVersion:   q.RuleVersion,
		PromptText:    promptText,
		Parameters:    paramValues,
		RandomSeed:    seed,
		Entry:         processed.Entry,
		Concepts:      q.Concepts,
		StageAnswers:  make(map[domain.DrillStage]domain.StageAnswer),
		ScaffoldLevel: domain.ScaffoldFull,
	}

	// Build progressive stages
	switch q.FamilyID {
	case bank.FamilyCustomerAdvance:
		buildCustomerAdvanceStages(instance, money, seed, processed)
	case bank.FamilyCashService:
		buildCashServiceStages(instance, money, seed, processed)
	case bank.FamilyServiceOnCredit:
		buildServiceOnCreditStages(instance, money, seed, processed)
	case bank.FamilyCollectReceivable:
		buildCollectReceivableStages(instance, money, seed, processed)
	case bank.FamilyEarnAdvance:
		buildEarnAdvanceStages(instance, money, seed, processed)
	case bank.FamilyCashRent, bank.FamilyCashExpense:
		buildCashRentStages(instance, money, seed, processed)
	case bank.FamilyBorrowCash:
		buildBorrowCashStages(instance, money, seed, processed)
	case bank.FamilyIssueShares:
		buildIssueSharesStages(instance, money, seed, processed)
	case bank.FamilyPrepaidPurchase:
		buildPrepaidPurchaseStages(instance, money, seed, processed)
	case bank.FamilyPrepaidConsumption:
		buildPrepaidConsumptionStages(instance, money, seed, processed)
	case bank.FamilyEquipmentPurchaseCash:
		buildEquipmentPurchaseCashStages(instance, money, seed, processed)
	case bank.FamilyRepayNotePrincipal:
		buildRepayNotePrincipalStages(instance, money, seed, processed)
	case bank.FamilyDividendCash:
		buildDividendCashStages(instance, money, seed, processed)
	default:
		buildGenericStages(instance, money, seed, processed)
	}

	if bank.RequiresReviewProvenance(q.Status) {
		for stage, text := range q.Teaching {
			if answer, ok := instance.StageAnswers[stage]; ok {
				answer.CausalHint = strings.ReplaceAll(text.Hint, "${amount_dollars}", money.FormatDollars())
				answer.Explanation = strings.ReplaceAll(text.Explanation, "${amount_dollars}", money.FormatDollars())
				instance.StageAnswers[stage] = answer
			}
		}
	}
	policy, profile, err := bank.ReviewedPedagogy(q)
	if err != nil {
		return nil, fmt.Errorf("invalid pedagogy policy: %w", err)
	}
	instance.Pedagogy = policy
	if policy.PolicyVersion > 0 {
		instance.InstanceID += fmt.Sprintf("-p%d", policy.PolicyVersion)
	}
	applyReviewedDistractors(instance, money, seed, profile)
	snapshotMistakeHints(instance)
	return instance, nil
}

// Helper to shuffle options deterministically using a seeded PRNG.
func shuffleOptions(options []domain.AnswerOption, seed int64, salt int64) []domain.AnswerOption {
	r := rand.New(rand.NewSource(seed + salt))
	shuffled := make([]domain.AnswerOption, len(options))
	copy(shuffled, options)

	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	labels := []string{"a", "b", "c", "d", "e"}
	for i := range shuffled {
		if i < len(labels) {
			shuffled[i].Label = labels[i]
		}
	}
	return shuffled
}

func buildCustomerAdvanceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Primary Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 101)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "What did the company receive today, and which account records it?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "Set aside whether the work is done. What came into the business today?",
		Explanation:       "The customer paid cash today, so Cash is affected. Service Revenue is tempting, but receiving money is not the same as earning it; the work happens next month.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Account Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 102)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash something the company owns and can use, something it owes, or the owners' stake?",
		Explanation:       "Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).",
		RelevantConceptID: "cash_classification",
	}

	// Stage 3: Direction of change
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 103)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon receiving this customer payment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Compare the company's cash before and after today's payment. Is it higher or lower?",
		Explanation:       "The company holds more money after the payment, so Cash increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 104)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Cash) recorded in double-entry bookkeeping?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on the same side as its normal balance. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account Identification (Key pedagogical decision!)
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_rev", Text: "Unearned Revenue (Liability)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_ar", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 105)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What counter-account balances this receipt for services still to be performed?",
		CorrectOptionID:   "opt_unearned_rev",
		Options:           stage5Opts,
		CausalHint:        "Has the company done the work yet? If not, what does it now owe the customer?",
		Explanation:       "The company has been paid but has not done the work, so it owes the customer the service (or a refund). That obligation is Unearned Revenue, a Liability. Service Revenue would record earning before any work is done. Accounts Receivable would mean the customer still owes money, but the customer has already paid.",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry Assembly
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecognizedPrematurely},
		{ID: "opt_dr_unearned_cr_cash", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 106)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this customer advance?",
		CorrectOptionID:   "opt_dr_cash_cr_unearned",
		Options:           stage6Opts,
		CausalHint:        "The customer paid for work the company will do next month. What did the company gain today, and what does it now owe?",
		Explanation:       fmt.Sprintf("Debit Cash %s (asset up) and Credit Unearned Revenue %s (liability up). Crediting Service Revenue instead would report revenue for work not yet done.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation and Transaction Recap
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecognizedPrematurely},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 107),
		CausalHint:        "Place each account in the entry under Assets, Liabilities, or Equity. Has the company earned anything yet?",
		Explanation:       fmt.Sprintf("Assets increase by %s (Cash) and Liabilities increase by %s (Unearned Revenue). Equity is unchanged because nothing has been earned yet; it increases next month, when the work is done.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildCashServiceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Primary Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: engine.TagRevenueDeferredWhenEarned},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 201)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account reflects the immediate payment received today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What came into the business today, and in what form?",
		Explanation:       "The customer paid cash today, so Cash is affected. Accounts Receivable would apply only if the customer still owed the money, and Unearned Revenue only if the payment came before the work.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Account Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 202)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash something the company owns and can use, something it owes, or the owners' stake?",
		Explanation:       "Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).",
		RelevantConceptID: "cash_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 203)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon receiving this payment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Compare the company's cash before and after today's payment. Is it higher or lower?",
		Explanation:       "The company holds more money after the payment, so Cash increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 204)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in an Asset (Cash) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on the same side as its normal balance. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: engine.TagRevenueDeferredWhenEarned},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 205)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account records what the company did for the customer in exchange for the cash?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "Is the work finished, or does the company still owe it to the customer?",
		Explanation:       "The work was completed today, so the company has earned revenue: Service Revenue. Unearned Revenue would mean the work is still owed, but it is already done.",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rev_cr_cash", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueDeferredWhenEarned},
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagWrongAccount},
	}, seed, 206)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this cash service transaction?",
		CorrectOptionID:   "opt_dr_cash_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "The customer paid today for work finished today. What did the company receive, and what did it earn?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Service Revenue %s. Payment and work happened on the same day, so there is no receivable and no unearned revenue.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Service Revenue); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueDeferredWhenEarned},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 207),
		CausalHint:        "Place each account in the entry under Assets, Liabilities, or Equity. Did the company earn anything, and does anyone owe anyone afterward?",
		Explanation:       fmt.Sprintf("Assets increase by %s (Cash). Equity increases by %s because revenue increases equity. Liabilities are unchanged because nothing is owed afterward.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildServiceOnCreditStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: engine.TagRevenueDeferredWhenEarned},
	}, seed, 301)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "The customer owes payment for services completed today and has not paid yet. Which account records this claim?",
		CorrectOptionID:   "opt_ar",
		Options:           stage1Opts,
		CausalHint:        "Did any money arrive today? If not, what does the company hold instead?",
		Explanation:       "The company has the right to collect from the customer later, which is Accounts Receivable. Cash is tempting, but no money arrived today.",
		RelevantConceptID: "accounts_receivable_classification",
	}

	// Stage 2: Account Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 302)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Receivable?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is a right to collect money later something the company owns, something it owes, or the owners' stake?",
		Explanation:       "Accounts Receivable is an Asset: the right to collect cash later. It is not Revenue. Revenue records the earning; the receivable records the amount still to be collected.",
		RelevantConceptID: "accounts_receivable_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 303)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Receivable increase or decrease when completed work creates a new unpaid customer claim?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "For the work completed today, does the invoice create a new amount to collect or settle an existing one?",
		Explanation:       "The customer now owes the company money, so Accounts Receivable increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 304)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Accounts Receivable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on the same side as its normal balance. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: engine.TagRevenueDeferredWhenEarned},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 305)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "The work was completed today but has not been paid for. Which account balances the entry?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "Under accrual accounting, does revenue wait for the cash, or for the work?",
		Explanation:       "Revenue is recorded when the work is done, not when cash arrives, so the balancing account is Service Revenue. Cash would record a payment that has not happened.",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 306)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this service on credit?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "The work is done but unpaid. What does the company now hold, and what did it earn?",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s and Credit Service Revenue %s. Debiting Cash would record money that has not been collected yet.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this service on credit affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Accounts Receivable); Equity increases by %s (+Service Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueDeferredWhenEarned},
		}, seed, 307),
		CausalHint:        "Place each account in the entry under Assets, Liabilities, or Equity. Did the company earn anything, and did any cash move?",
		Explanation:       fmt.Sprintf("Assets increase by %s (Accounts Receivable). Equity increases by %s because revenue increases equity. Liabilities are unchanged. Cash stays the same until the customer pays.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildCollectReceivableStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagDuplicateRevenueOnCollection},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 401)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Cash arrives to settle the customer's invoice. Which account reflects the cash received today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What arrived in the business today?",
		Explanation:       "The customer paid money today, so Cash is affected. Service Revenue is tempting, but this payment is for work already recorded as revenue last month.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 402)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash something the company owns and can use, something it owes, or the owners' stake?",
		Explanation:       "Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).",
		RelevantConceptID: "cash_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 403)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon collecting this customer payment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Compare the company's cash before and after today's payment. Is it higher or lower?",
		Explanation:       "The company holds more money after the payment, so Cash increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 404)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Cash) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on the same side as its normal balance. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagDuplicateRevenueOnCollection},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 405)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?",
		CorrectOptionID:   "opt_ar",
		Options:           stage5Opts,
		CausalHint:        "Revenue was recorded last month, when the work was done. What has the company been holding since then that this payment settles?",
		Explanation:       "The payment settles the amount the customer owed, so the balancing account is Accounts Receivable, which decreases. Crediting Service Revenue again would count last month's work twice.",
		RelevantConceptID: "collection_vs_earning",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_ar", Text: fmt.Sprintf("Debit Cash %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDuplicateRevenueOnCollection},
		{ID: "opt_dr_ar_cr_cash", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 406)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for collecting this receivable?",
		CorrectOptionID:   "opt_dr_cash_cr_ar",
		Options:           stage6Opts,
		CausalHint:        "The work was recorded as revenue last month. What did the company receive today, and what did that payment settle?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Accounts Receivable %s. Revenue is not touched; it was recorded when the work was done.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does collecting an existing receivable affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Cash increases (+%s) and Accounts Receivable decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDuplicateRevenueOnCollection},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagWrongAccount},
		}, seed, 407),
		CausalHint:        "Place each account in the entry under Assets, Liabilities, or Equity. Did total assets change, or did one asset turn into another?",
		Explanation:       fmt.Sprintf("Cash increases by %s and Accounts Receivable decreases by %s: one asset turned into another. Total assets, liabilities, and equity are unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "collection_vs_earning",
	}
}

func buildEarnAdvanceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedOnEarningAdvance},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 501)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "No new cash changes hands when this prepaid work is completed. Which account recorded the obligation before completion?",
		CorrectOptionID:   "opt_unearned_rev",
		Options:           stage1Opts,
		CausalHint:        "When the customer paid earlier, what did the company owe in return?",
		Explanation:       "The earlier payment created an obligation to do the work, recorded as Unearned Revenue. Completing the prepaid work settles it. Cash is not involved in this completion entry because the receipt was already recorded.",
		RelevantConceptID: "earned_vs_unearned",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 502)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Unearned Revenue?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?",
		Explanation:       "Unearned Revenue is a Liability: the company owes the customer work, or a refund. Despite its name, it is not a Revenue account.",
		RelevantConceptID: "unearned_revenue_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 503)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Unearned Revenue increase or decrease as the promised service is delivered?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "Before completion, the company owed the customer this work. After finishing it, how much does it still owe?",
		Explanation:       "The obligation has been fulfilled, so Unearned Revenue decreases.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 504)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is a decrease to a Liability (Unearned Revenue) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Which side is a liability's normal balance, and does a decrease go on that side or the opposite one?",
		Explanation:       "Liabilities have a normal credit balance, so a decrease is recorded as a Debit (left side). This debit does not mean money went out; no new cash moves when the prepaid work is completed.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedOnEarningAdvance},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagAdvanceConfusedWithReceivable},
	}, seed, 505)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "The prepaid work is completed today. Which account balances the entry?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "After completing the prepaid work today, has the company earned anything new?",
		Explanation:       "Completing the work earns revenue, so the balancing account is Service Revenue. Cash would record a payment, but the receipt was recorded earlier.",
		RelevantConceptID: "earned_vs_unearned",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_unearned_cr_rev", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedOnEarningAdvance},
		{ID: "opt_dr_rev_cr_unearned", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 506)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for earning this advance?",
		CorrectOptionID:   "opt_dr_unearned_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "The customer paid earlier, and the prepaid work is now complete. What does the company no longer owe, and what has it now earned?",
		Explanation:       fmt.Sprintf("Debit Unearned Revenue %s and Credit Service Revenue %s. No new cash is recorded at completion because the receipt was recorded earlier.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does earning a customer advance affect the accounting equation?",
		CorrectOptionID: "opt_liab_down_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_liab_down_eq_up", Text: fmt.Sprintf("Liabilities decrease by %s (-Unearned Revenue); Equity increases by %s (+Service Revenue); Total Assets unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedOnEarningAdvance},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedOnEarningAdvance},
		}, seed, 507),
		CausalHint:        "Place each account in the entry under Assets, Liabilities, or Equity. Does any new cash move when the prepaid work is completed?",
		Explanation:       fmt.Sprintf("Liabilities decrease by %s (Unearned Revenue). Equity increases by %s because revenue increases equity. Total assets are unchanged in this completion entry because the cash receipt was recorded earlier.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "earned_vs_unearned",
	}
}

func buildCashRentStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_rent_expense", Text: "Rent Expense"},
		{ID: "opt_prepaid_rent", Text: "Prepaid Rent", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 601)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the cost of using the premises this month?",
		CorrectOptionID:   "opt_rent_expense",
		Options:           stage1Opts,
		CausalHint:        "Does today's payment buy future use of the premises or pay for the current month's use?",
		Explanation:       "Rent Expense records the current month's occupancy cost. Prepaid Rent would represent future use, and Accounts Payable would represent an amount still owed; this rent was paid today.",
		RelevantConceptID: "rent_expense_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
	}, seed, 602)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Rent Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does this account track a future resource, an obligation, or a cost incurred in operating the business?",
		Explanation:       "Rent Expense is an Expense: a cost of using the premises this month. It reduces equity through income; it is not an asset providing future use.",
		RelevantConceptID: "rent_expense_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 603)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the expense account balance increase or decrease when recording this rent payment?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing this month's rent add to the costs accumulated this period or remove a previously recorded cost?",
		Explanation:       "Rent Expense increases as this month's cost is recorded. Cash decreases, but that does not make the expense account decrease.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 604)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in an Expense (Rent Expense) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an expense's normal balance?",
		Explanation:       "Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 605)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when the company pays for this month's use of the premises?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Did the company hand over money today or only promise to pay later?",
		Explanation:       "Cash decreases because the company paid today, so Cash is credited. Recording a payable instead would leave the already-paid amount outstanding.",
		RelevantConceptID: "cash_vs_expense",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rent", Text: fmt.Sprintf("Debit Cash %s / Credit Rent Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_rent_cr_ap", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_dr_prepaid_cr_cash", Text: fmt.Sprintf("Debit Prepaid Rent %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagWrongAccount},
	}, seed, 606)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying cash rent?",
		CorrectOptionID:   "opt_dr_rent_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "The premises were used this month and paid for today. What cost was incurred, and what resource was given up?",
		Explanation:       fmt.Sprintf("Debit Rent Expense %s and Credit Cash %s. The cost belongs to this month, so treating it as a prepayment would defer a cost already incurred; a payable would ignore today's payment.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does paying cash rent affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Rent Expense); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_net_change", Text: "No change; expenses do not affect balance sheet.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 607),
		CausalHint:        "Did the payment leave the company with a new future resource, or was this month's benefit already used?",
		Explanation:       fmt.Sprintf("Assets decrease by %s (Cash) and equity decreases by %s through Rent Expense. No debt was settled: the rent was not previously owed. The expense increases even though equity decreases.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_expense",
	}
}

func buildBorrowCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagRevenueRecordedOnBorrowing},
	}, seed, 701)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "What did the company receive today, and which account records it?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What came into the business from the lender today?",
		Explanation:       "Cash records the money received. Notes Payable records the separate repayment obligation; Service Revenue would imply earning from work, which borrowing does not create.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 702)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash a resource the company owns, an amount it owes, or an ownership claim?",
		Explanation:       "Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 703)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon borrowing money?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Compare the money held before and after today's receipt. Is there more or less?",
		Explanation:       "Cash increases because money arrived today. A future repayment obligation does not cancel today's receipt.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 704)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Cash) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagRevenueRecordedOnBorrowing},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagWrongAccount},
	}, seed, 705)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account balances the receipt from the lender?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage5Opts,
		CausalHint:        "Did the lender buy ownership, pay for work, or expect the money back?",
		Explanation:       "Notes Payable records the signed promise to repay, so it is credited as the liability increases. Service Revenue would claim earning; Common Stock would claim an ownership investment. Neither happened.",
		RelevantConceptID: "notes_payable_classification",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecordedOnBorrowing},
		{ID: "opt_dr_notes_cr_cash", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 706)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for borrowing cash on a note?",
		CorrectOptionID:   "opt_dr_cash_cr_notes",
		Options:           stage6Opts,
		CausalHint:        "The company received money under a signed loan agreement. What did it gain, and what must it do later?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Notes Payable %s. Cash increased and a debt was created. Crediting Service Revenue instead would report borrowing as income.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does borrowing cash affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Notes Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecordedOnBorrowing},
			{ID: "opt_no_net_change", Text: "No net change in total assets.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 707),
		CausalHint:        "Does receiving a loan create earning, an ownership contribution, or a repayment obligation?",
		Explanation:       fmt.Sprintf("Assets increase by %s (Cash) and liabilities increase by %s (Notes Payable). Equity is unchanged because the cash was borrowed, not earned or contributed by owners.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildIssueSharesStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagRevenueRecordedOnShareIssue},
	}, seed, 801)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "What did the company receive from the investors today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What resource did investors hand over in exchange for their shares?",
		Explanation:       "Cash records the money received. Common Stock records the ownership contribution on the other side of the entry; Service Revenue would imply work for a customer, which did not happen.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 802)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash a resource the company owns, an amount it owes, or an ownership claim?",
		Explanation:       "Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 803)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon issuing stock for cash?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Compare the money held before and after today's receipt. Is there more or less?",
		Explanation:       "Cash increases because money arrived today. A future repayment obligation does not cancel today's receipt.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 4: Debit or Credit
	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 804)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Cash) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagRevenueRecordedOnShareIssue},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 805)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account balances the receipt from investors who bought shares?",
		CorrectOptionID:   "opt_common_stock",
		Options:           stage5Opts,
		CausalHint:        "Did the investors buy ownership, lend money to be repaid, or pay for work?",
		Explanation:       "Common Stock records capital contributed in exchange for shares. Service Revenue would treat investment as earning, and Notes Payable would treat the investors as lenders.",
		RelevantConceptID: "common_stock_classification",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_stock", Text: fmt.Sprintf("Debit Cash %s / Credit Common Stock %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRevenueRecordedOnShareIssue},
		{ID: "opt_dr_stock_cr_cash", Text: fmt.Sprintf("Debit Common Stock %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 806)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for issuing common stock for cash?",
		CorrectOptionID:   "opt_dr_cash_cr_stock",
		Options:           stage6Opts,
		CausalHint:        "The company received money in exchange for common shares. What resource came in, and what claim did the investors receive?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Common Stock %s. Owners contributed capital; no services were earned and no loan was created.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing common shares affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Common Stock); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_net_change", Text: "No net change in total assets.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 807),
		CausalHint:        "What did the investors receive for their money: an ownership stake or a promise of repayment?",
		Explanation:       fmt.Sprintf("Assets increase by %s (Cash) and equity increases by %s (Common Stock). Liabilities are unchanged. Equity rose through contributed capital, not through revenue.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildPrepaidPurchaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_cash", Text: "Cash"},
	}, seed, 901)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the insurance coverage acquired today?",
		CorrectOptionID:   "opt_prepaid_insurance",
		Options:           stage1Opts,
		CausalHint:        "Has the purchased coverage been used yet, or can the company still use it in future months?",
		Explanation:       "Prepaid Insurance records the coverage still available. Insurance Expense would mean coverage was used, but none has expired today. Cash records the payment, not the remaining coverage.",
		RelevantConceptID: "prepaid_insurance_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 902)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Prepaid Insurance?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is unused coverage a future benefit the company controls or a cost already used up?",
		Explanation:       "Prepaid Insurance is an Asset: the right to future coverage. It becomes an expense as coverage is used; paying today does not mean the entire benefit was consumed today.",
		RelevantConceptID: "prepaid_insurance_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 903)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Prepaid Insurance account balance increase or decrease upon purchasing the policy?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does buying the policy add to the coverage still available or use up coverage already purchased?",
		Explanation:       "Prepaid Insurance increases because the company acquired unused coverage. Cash decreases, but a different asset increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 904)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Prepaid Insurance) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 905)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when the policy is paid for today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Did the company hand over money today or only promise to pay later?",
		Explanation:       "Cash decreases because the company paid today, so Cash is credited. Recording a payable instead would leave the already-paid amount outstanding.",
		RelevantConceptID: "prepaid_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_prepaid_cr_cash", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_dr_cash_cr_prepaid", Text: fmt.Sprintf("Debit Cash %s / Credit Prepaid Insurance %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_prepaid_cr_ap", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 906)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for purchasing this prepaid insurance policy?",
		CorrectOptionID:   "opt_dr_prepaid_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "The policy was paid for today and no coverage has expired. What benefit was acquired, and what resource was given up?",
		Explanation:       fmt.Sprintf("Debit Prepaid Insurance %s and Credit Cash %s. Unused coverage is an asset. Debiting Insurance Expense would charge future coverage to today's income.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does purchasing prepaid insurance affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Prepaid Insurance increases (+%s) and Cash decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Insurance Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 907),
		CausalHint:        "Did paying for unused coverage use up a benefit or exchange money for a benefit still available?",
		Explanation:       fmt.Sprintf("Prepaid Insurance increases by %s and Cash decreases by %s. Total assets, liabilities, and equity are unchanged. Equity would fall when coverage is consumed, not when this unused policy is purchased.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "prepaid_vs_expense",
	}
}

func buildPrepaidConsumptionStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_insurance_expense", Text: "Insurance Expense"},
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedOnPrepaidExpiration},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
	}, seed, 1001)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the cost of insurance coverage used this month?",
		CorrectOptionID:   "opt_insurance_expense",
		Options:           stage1Opts,
		CausalHint:        "Was coverage acquired today, or was part of previously purchased coverage used up?",
		Explanation:       "Insurance Expense records the coverage used this month. Prepaid Insurance tracks the coverage remaining; Cash does not change because the policy was paid for earlier.",
		RelevantConceptID: "insurance_expense_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 1002)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Insurance Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the account describe protection still available or protection already used during the month?",
		Explanation:       "Insurance Expense is an Expense: the cost of protection already used. Prepaid Insurance is the Asset for protection still available; these accounts describe different parts of the same policy.",
		RelevantConceptID: "insurance_expense_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1003)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Insurance Expense increase or decrease as coverage expires during the month?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing the coverage used this month add to this period's costs or remove a previously recorded cost?",
		Explanation:       "Insurance Expense increases as consumed coverage is recorded. The prepaid asset decreases; the cost account increases.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1004)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in an Expense (Insurance Expense) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an expense's normal balance?",
		Explanation:       "Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedOnPrepaidExpiration},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 1005)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when already-paid insurance coverage is used?",
		CorrectOptionID:   "opt_prepaid_insurance",
		Options:           stage5Opts,
		CausalHint:        "What balance has been holding the coverage that is now used up?",
		Explanation:       "Prepaid Insurance decreases as its coverage is consumed, so it is credited. Cash would record another payment, but no money leaves today.",
		RelevantConceptID: "prepaid_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_expense_cr_prepaid", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Prepaid Insurance %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedOnPrepaidExpiration},
		{ID: "opt_dr_prepaid_cr_expense", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Insurance Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_no_entry", Text: "No entry; the full policy remains prepaid.", ErrorTag: engine.TagPrepaidNotExpensedOnConsumption},
	}, seed, 1006)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for recording this expired insurance?",
		CorrectOptionID:   "opt_dr_expense_cr_prepaid",
		Options:           stage6Opts,
		CausalHint:        "Coverage bought earlier has now been used, with no payment today. What cost arose, and what future benefit remains smaller?",
		Explanation:       fmt.Sprintf("Debit Insurance Expense %s and Credit Prepaid Insurance %s. Crediting Cash would record a second payment. No entry would leave used-up coverage reported as a future resource.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recognizing expired insurance affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Prepaid Insurance); Equity decreases by %s (-Insurance Expense); Cash is unaffected.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_cash", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Insurance Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedOnPrepaidExpiration},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1007),
		CausalHint:        "Is the consumed coverage still a resource available for the future, and did money move today?",
		Explanation:       fmt.Sprintf("Assets decrease by %s (Prepaid Insurance) and equity decreases by %s through Insurance Expense. Cash and liabilities are unchanged. Keeping all balances unchanged would fail to recognize the coverage consumed.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "prepaid_vs_expense",
	}
}

func buildEquipmentPurchaseCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_equipment", Text: "Equipment"},
		{ID: "opt_expense_misconception", Text: "An expense for the cost of the machinery", ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
	}, seed, 1101)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the newly acquired machinery?",
		CorrectOptionID:   "opt_equipment",
		Options:           stage1Opts,
		CausalHint:        "Will the purchased machinery provide use only today or over future years?",
		Explanation:       "Equipment records the machinery the company now owns. An immediate expense would treat the full multi-year resource as already used; Cash records how it was paid for.",
		RelevantConceptID: "equipment_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1102)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Equipment?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is machinery that can be used in future years a resource still owned or a cost already consumed?",
		Explanation:       "Equipment is an Asset: a productive resource available over future years. Paying cash does not make the entire purchase an immediate operating expense.",
		RelevantConceptID: "equipment_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1103)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Equipment account balance increase or decrease upon acquisition?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "After this purchase, does the company own more machinery or less?",
		Explanation:       "Equipment increases because new machinery was acquired. Cash decreases, but the Equipment account tracks machinery, not money.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1104)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in an Asset (Equipment) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "An account increases on its normal-balance side. Which side is an asset's normal balance?",
		Explanation:       "Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagPayableRecordedForCashPayment},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
	}, seed, 1105)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when the machinery is paid for today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Did the company hand over money today or only promise to pay later?",
		Explanation:       "Cash decreases because the company paid today, so Cash is credited. Recording a payable instead would leave the already-paid amount outstanding.",
		RelevantConceptID: "capex_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_equip_cr_cash", Text: fmt.Sprintf("Debit Equipment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit an expense for the machinery %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_dr_cash_cr_equip", Text: fmt.Sprintf("Debit Cash %s / Credit Equipment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1106)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this cash equipment purchase?",
		CorrectOptionID:   "opt_dr_equip_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Machinery with years of future use was paid for today. What resource was acquired, and what resource was given up?",
		Explanation:       fmt.Sprintf("Debit Equipment %s and Credit Cash %s. Both are assets. An expense debit would treat the full future benefit as already used; a payable would ignore the cash payment.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does purchasing equipment for cash affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Equipment increases (+%s) and Cash decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1107),
		CausalHint:        "Did this cash purchase leave the company with another resource, or was the entire benefit used today?",
		Explanation:       fmt.Sprintf("Equipment increases by %s and Cash decreases by %s. Total assets, liabilities, and equity are unchanged. Buying on credit would increase a liability, but this machinery was paid for today.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "capex_vs_expense",
	}
}

func buildRepayNotePrincipalStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable"},
		{ID: "opt_expense_misconception", Text: "An expense for the loan payment", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 1201)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records what the company owed the bank before today's principal payment?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage1Opts,
		CausalHint:        "What was recorded when the company originally received the loan?",
		Explanation:       "Notes Payable tracks the principal owed to the bank. Repayment settles that debt; an expense would imply a new cost of operating or borrowing. This transaction includes principal only, not interest.",
		RelevantConceptID: "notes_payable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1202)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Notes Payable?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does a signed promise to repay the bank describe a resource owned or an obligation owed?",
		Explanation:       "Notes Payable is a Liability: principal the company must repay. An expense records a cost incurred; the loan principal is previously borrowed money, not a new expense.",
		RelevantConceptID: "notes_payable_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 1203)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Notes Payable increase or decrease when principal is repaid?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "After repaying part of the principal, does the company owe the bank more or less?",
		Explanation:       "Notes Payable decreases as principal is settled. Borrowing would increase the debt, but today's event is repayment.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1204)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is a decrease in a Liability (Notes Payable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Which side is a liability's normal balance, and does a decrease use that side or the opposite side?",
		Explanation:       "Liabilities have a normal credit balance, so a decrease is a Debit (left side). The debit reduces debt; it does not mean cash increased.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
	}, seed, 1205)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when the company pays principal today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Did the company hand over money today or only promise to pay later?",
		Explanation:       "Cash decreases because the company paid today, so Cash is credited. Recording a payable instead would leave the already-paid amount outstanding.",
		RelevantConceptID: "principal_repayment_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_notes_cr_cash", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit an expense for the loan payment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1206)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for repaying note principal with cash?",
		CorrectOptionID:   "opt_dr_notes_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Previously borrowed principal is repaid today, with interest excluded. What obligation is settled, and what resource leaves?",
		Explanation:       fmt.Sprintf("Debit Notes Payable %s and Credit Cash %s. Debt and cash both decrease. Debiting an expense instead would reduce income for a principal payment that is not a cost.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does repaying loan principal affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_liab_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Liabilities decrease by %s (-Notes Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1207),
		CausalHint:        "Does repayment of previously borrowed principal incur a new cost or settle an existing obligation?",
		Explanation:       fmt.Sprintf("Assets decrease by %s (Cash) and liabilities decrease by %s (Notes Payable). Equity is unchanged because principal repayment is not an expense. Both sides of the equation fall equally.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "principal_repayment_vs_expense",
	}
}

func buildDividendCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dividends", Text: "Dividends"},
		{ID: "opt_expense_misconception", Text: "An expense for the payment to stockholders", ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash"},
	}, seed, 1301)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records cash distributions paid directly to stockholders?",
		CorrectOptionID:   "opt_dividends",
		Options:           stage1Opts,
		CausalHint:        "Are the stockholders being paid for work, or receiving a distribution because they own shares?",
		Explanation:       "Dividends records the distribution declared and paid to owners today. An expense would treat owners as suppliers of a business service; this distribution is not a cost of earning revenue.",
		RelevantConceptID: "dividends_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dividends_cat", Text: "Dividends"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
	}, seed, 1302)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Dividends?",
		CorrectOptionID:   "opt_dividends_cat",
		Options:           stage2Opts,
		CausalHint:        "Does this account track operating costs, a resource owned, or distributions to owners?",
		Explanation:       "Dividends is the Dividends category, which reduces equity. It is not an Expense and does not reduce net income; the account is closed to Retained Earnings at period end.",
		RelevantConceptID: "dividends_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1303)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Dividends account balance increase or decrease when recording this dividend distribution?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does declaring this distribution add to the dividends accumulated this period or reverse a previous distribution?",
		Explanation:       "Dividends increases as a new distribution is recorded, even though equity decreases. The distribution account and total equity move in opposite directions.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1304)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in Dividends recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Which side is the Dividends account's normal balance, and where does an increase belong?",
		Explanation:       "Dividends has a normal debit balance, so an increase is a Debit (left side). Crediting Dividends would reduce recorded distributions; an equity reduction does not mean the Dividends account decreases.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_dividends_payable", Text: "Dividends Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 1305)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when this newly declared distribution is paid today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Did the company hand over money today or only promise to pay later?",
		Explanation:       "Cash decreases because the company paid today, so Cash is credited. Recording a payable instead would leave the already-paid amount outstanding.",
		RelevantConceptID: "dividend_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_div_cr_cash", Text: fmt.Sprintf("Debit Dividends %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit an expense for the dividend %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_dr_cash_cr_div", Text: fmt.Sprintf("Debit Cash %s / Credit Dividends %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1306)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying cash dividends?",
		CorrectOptionID:   "opt_dr_div_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "The dividend was declared and paid today, with no earlier declaration. What distribution was recorded, and what resource left?",
		Explanation:       fmt.Sprintf("Debit Dividends %s and Credit Cash %s. It is an owner distribution, not an expense. Debiting Dividends Payable would settle an earlier declaration, but no such payable existed.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does paying cash dividends affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Dividends); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_net_change", Text: "No change in equity; dividends are an asset.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1307),
		CausalHint:        "Does distributing cash to owners change operating income or the owners' stake directly?",
		Explanation:       fmt.Sprintf("Assets decrease by %s (Cash) and equity decreases by %s through Dividends. Liabilities are unchanged because declaration and payment happened today. Net income is unchanged: dividends are not an expense.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "dividend_vs_expense",
	}
}

func buildGenericStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	if len(proc.Entry.Postings) < 2 {
		return
	}
	drPosting := proc.Entry.Postings[0]
	crPosting := proc.Entry.Postings[1]

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:           domain.StageIdentifyAccount,
		Prompt:          fmt.Sprintf("Which account is debited in this %s transaction?", inst.FamilyID),
		CorrectOptionID: fmt.Sprintf("opt_%s", drPosting.AccountID),
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: fmt.Sprintf("opt_%s", drPosting.AccountID), Text: string(drPosting.AccountID)},
			{ID: fmt.Sprintf("opt_%s", crPosting.AccountID), Text: string(crPosting.AccountID)},
			{ID: "opt_cash", Text: "cash"},
		}, seed, 501),
		CausalHint:        proc.Explanation.DebitRationale,
		Explanation:       proc.Explanation.DebitRationale,
		RelevantConceptID: "debit_credit_translation",
	}
}

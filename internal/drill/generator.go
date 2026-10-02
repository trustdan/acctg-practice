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
		InstanceID:    fmt.Sprintf("%s-%d-%d", q.ID, seed, amtVal),
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

	if q.Status == bank.StatusApprovedActive {
		for stage, text := range q.Teaching {
			if answer, ok := instance.StageAnswers[stage]; ok {
				answer.CausalHint = strings.ReplaceAll(text.Hint, "${amount_dollars}", money.FormatDollars())
				answer.Explanation = strings.ReplaceAll(text.Explanation, "${amount_dollars}", money.FormatDollars())
				instance.StageAnswers[stage] = answer
			}
		}
	}
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
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: "treated_advance_as_revenue"},
		{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: "confused_advance_with_receivable"},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: "confused_with_payable"},
	}, seed, 101)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account is directly affected by the cash payment received from the customer today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "The company physically received money today. What asset account tracks cash inflows?",
		Explanation:       "Cash is received immediately, so the Cash account is affected.",
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
		CausalHint:        "Cash represents an economic resource owned and controlled by the company.",
		Explanation:       "Cash is a current Asset representing available financial resources.",
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
		CausalHint:        "Money came into the business bank account from the customer.",
		Explanation:       "Receiving cash increases the Cash account balance.",
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
		CausalHint:        "Assets have a normal debit balance. Debits increase assets; credits decrease them.",
		Explanation:       "Assets increase on the Debit (left) side.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account Identification (Key pedagogical decision!)
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_rev", Text: "Unearned Revenue (Liability)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue / Equity)", ErrorTag: "revenue_recognized_prematurely"},
		{ID: "opt_ar", Text: "Accounts Receivable (Asset)", ErrorTag: "confused_advance_with_receivable"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: "confused_with_owner_investment"},
	}, seed, 105)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What counter-account balances this entry for services to be performed next month?",
		CorrectOptionID:   "opt_unearned_rev",
		Options:           stage5Opts,
		CausalHint:        "The service has not been performed yet. Revenue cannot be recognized before it is earned. Receiving cash before performing work creates an obligation (liability) to the customer.",
		Explanation:       "Because the service has not been provided, the company has an obligation to perform in the future, recorded as Unearned Revenue (a liability).",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry Assembly
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "revenue_recognized_prematurely"},
		{ID: "opt_dr_unearned_cr_cash", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "service_on_credit"},
	}, seed, 106)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this customer advance?",
		CorrectOptionID:   "opt_dr_cash_cr_unearned",
		Options:           stage6Opts,
		CausalHint:        "You need a Debit to Cash (asset up) and a Credit to Unearned Revenue (liability up) for equal amounts.",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Unearned Revenue %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation and Transaction Recap
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_advance_as_equity"},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: "confused_with_asset_swap"},
		}, seed, 107),
		CausalHint:        "Cash is an Asset. Unearned Revenue is a Liability. Has any Equity/Revenue changed?",
		Explanation:       fmt.Sprintf("Assets increase by %s and Liabilities increase by %s. Both sides of the equation remain in balance.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildCashServiceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Primary Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
		{ID: "opt_ap", Text: "Accounts Payable"},
	}, seed, 201)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account reflects the immediate payment received today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "The business received currency/funds today.",
		Explanation:       "Cash is received immediately today, increasing the Cash account.",
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
		CausalHint:        "Cash represents an economic resource owned and controlled by the company.",
		Explanation:       "Cash is a current Asset.",
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
		CausalHint:        "Money came into the business from the customer.",
		Explanation:       "Cash increases with the receipt of money.",
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
		CausalHint:        "Assets have a normal debit balance; increases go on the left side.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue", ErrorTag: "treated_service_as_unearned"},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: "confused_with_borrowing"},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: "confused_with_payable"},
	}, seed, 205)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Because services were performed today and cash was received, what counter-account earns this inflow?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "The work is completed today. When work is performed, what account records earnings?",
		Explanation:       "Work completed today with immediate cash payment is recorded as Service Revenue.",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rev_cr_cash", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_service_as_unearned"},
	}, seed, 206)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this cash service transaction?",
		CorrectOptionID:   "opt_dr_cash_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Debit the asset that increased (Cash) and credit the revenue account that increased (Service Revenue).",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Service Revenue %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Service Revenue); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Revenue); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_service_as_unearned"},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: "confused_with_asset_swap"},
		}, seed, 207),
		CausalHint:        "Cash increases total assets. Does earning revenue increase owner's equity?",
		Explanation:       fmt.Sprintf("Assets increase by %s and Equity increases by %s through earned revenue.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildServiceOnCreditStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_before_collection"},
		{ID: "opt_ap", Text: "Accounts Payable"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
	}, seed, 301)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "The customer was invoiced for services completed today and has not paid yet. Which account records this claim?",
		CorrectOptionID:   "opt_ar",
		Options:           stage1Opts,
		CausalHint:        "The customer owes the company for work completed on credit.",
		Explanation:       "Accounts Receivable is debited because the company holds a claim to collect cash in the future.",
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
		CausalHint:        "Accounts Receivable is an economic resource (a legal claim) owned and controlled by the company.",
		Explanation:       "Accounts Receivable is an Asset.",
		RelevantConceptID: "accounts_receivable_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 303)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Receivable increase or decrease when the company bills a new customer?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "A new claim to collect future money was created.",
		Explanation:       "Accounts Receivable increases when new services are billed on account.",
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
		CausalHint:        "Assets increase on the debit (left) side.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_before_collection"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
		{ID: "opt_ap", Text: "Accounts Payable"},
	}, seed, 305)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "The work has been completed today. What account earns this inflow under accrual accounting?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "Under accrual accounting, revenue is recognized when performance is satisfied, regardless of when cash is collected.",
		Explanation:       "Service Revenue is credited because the service was performed today.",
		RelevantConceptID: "cash_vs_revenue",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "cash_recorded_before_collection"},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
	}, seed, 306)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this service on credit?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Debit Accounts Receivable (asset up) and Credit Service Revenue (equity up).",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s and Credit Service Revenue %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this service on credit affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Accounts Receivable); Equity increases by %s (+Service Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: "confused_with_asset_swap"},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_revenue_as_liability"},
		}, seed, 307),
		CausalHint:        "Accounts Receivable increases assets; Service Revenue increases equity.",
		Explanation:       fmt.Sprintf("Assets increase by %s (+Accounts Receivable) and Equity increases by %s (+Service Revenue).", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildCollectReceivableStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: "duplicate_revenue_on_collection"},
		{ID: "opt_ap", Text: "Accounts Payable"},
	}, seed, 401)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "The customer pays cash to settle their invoice. Which account reflects the cash received today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "The business received cash funds today.",
		Explanation:       "Cash is received, so Cash is debited.",
		RelevantConceptID: "cash_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 402)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Cash is an asset.",
		Explanation:       "Cash is an Asset.",
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
		CausalHint:        "Funds came in.",
		Explanation:       "Cash increases.",
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
		CausalHint:        "Assets increase on the debit side.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: "duplicate_revenue_on_collection"},
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
		{ID: "opt_ap", Text: "Accounts Payable"},
	}, seed, 405)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "The customer is paying an invoice where revenue was already recognized last month. What account must be credited?",
		CorrectOptionID:   "opt_ar",
		Options:           stage5Opts,
		CausalHint:        "Revenue was already recognized in the prior period. Crediting revenue again would double-count sales! Which asset was holding the customer's promise to pay?",
		Explanation:       "Accounts Receivable is credited to clear the existing claim. No new revenue is recorded.",
		RelevantConceptID: "collection_vs_earning",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_ar", Text: fmt.Sprintf("Debit Cash %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "duplicate_revenue_on_collection"},
		{ID: "opt_dr_ar_cr_cash", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
	}, seed, 406)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for collecting this receivable?",
		CorrectOptionID:   "opt_dr_cash_cr_ar",
		Options:           stage6Opts,
		CausalHint:        "Debit Cash (cash up) and Credit Accounts Receivable (receivable cleared).",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Accounts Receivable %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does collecting an existing receivable affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Cash increases (+%s) and Accounts Receivable decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "duplicate_revenue_on_collection"},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars())},
		}, seed, 407),
		CausalHint:        "One asset (Cash) went up, and another asset (Accounts Receivable) went down by the exact same amount.",
		Explanation:       "This is an asset exchange. Total assets are unchanged, and no new equity/revenue is recognized.",
		RelevantConceptID: "collection_vs_earning",
	}
}

func buildEarnAdvanceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_rev", Text: "Unearned Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_on_earning_advance"},
		{ID: "opt_ar", Text: "Accounts Receivable"},
	}, seed, 501)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Services paid for last month are now completed today. Which liability account is fulfilled and debited?",
		CorrectOptionID:   "opt_unearned_rev",
		Options:           stage1Opts,
		CausalHint:        "No cash changed hands today. The liability representing future work is being discharged.",
		Explanation:       "Unearned Revenue is fulfilled and reduced by debiting it.",
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
		CausalHint:        "Unearned Revenue represents an obligation to deliver goods or services to the customer.",
		Explanation:       "Unearned Revenue is a Liability.",
		RelevantConceptID: "unearned_revenue_classification",
	}

	// Stage 3: Direction
	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 503)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the liability Unearned Revenue increase or decrease as the promised service is delivered?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "The obligation to do future work has been fulfilled, reducing the outstanding liability.",
		Explanation:       "Unearned Revenue decreases when performance is satisfied.",
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
		CausalHint:        "Liabilities have a normal credit balance. Decreases go on the opposite side (debit).",
		Explanation:       "Liabilities decrease by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_on_earning_advance"},
		{ID: "opt_ar", Text: "Accounts Receivable"},
	}, seed, 505)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Now that the service is performed, what counter-account recognizes the earned revenue?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "Performance is complete. Revenue is recognized when earned.",
		Explanation:       "Service Revenue is credited because the earnings process is complete.",
		RelevantConceptID: "earned_vs_unearned",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_unearned_cr_rev", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "cash_recorded_on_earning_advance"},
		{ID: "opt_dr_rev_cr_unearned", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
	}, seed, 506)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for earning this advance?",
		CorrectOptionID:   "opt_dr_unearned_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Debit Unearned Revenue (liability down) and Credit Service Revenue (revenue up).",
		Explanation:       fmt.Sprintf("Debit Unearned Revenue %s and Credit Service Revenue %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does earning a customer advance affect the accounting equation?",
		CorrectOptionID: "opt_liab_down_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_liab_down_eq_up", Text: fmt.Sprintf("Liabilities decrease by %s (-Unearned Revenue); Equity increases by %s (+Service Revenue); Total Assets unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "cash_recorded_on_earning_advance"},
			{ID: "opt_no_net_change", Text: "No effect on any balance; memo entry only."},
		}, seed, 507),
		CausalHint:        "Liabilities decreased because work was done; Equity increased because revenue was earned. Total assets did not change.",
		Explanation:       fmt.Sprintf("Liabilities decrease by %s and Equity increases by %s. Total assets are unaffected.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "earned_vs_unearned",
	}
}

func buildCashRentStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_rent_expense", Text: "Rent Expense"},
		{ID: "opt_prepaid_rent", Text: "Prepaid Rent"},
		{ID: "opt_ap", Text: "Accounts Payable"},
	}, seed, 601)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the rent expense incurred for the current month?",
		CorrectOptionID:   "opt_rent_expense",
		Options:           stage1Opts,
		CausalHint:        "The rent is for the current month, meaning it is an operating expense of the current period.",
		Explanation:       "Rent Expense is debited for current-period occupancy cost.",
		RelevantConceptID: "rent_expense_classification",
	}

	// Stage 2: Category
	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense (Equity reduction)"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
	}, seed, 602)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Rent Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Expenses reflect costs of doing business that reduce owner equity.",
		Explanation:       "Rent Expense is an Expense account.",
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
		CausalHint:        "Incurring an expense increases the total accumulated expenses for the period.",
		Explanation:       "Expense balances increase when additional costs are incurred.",
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
		CausalHint:        "Expenses have normal debit balances because they reduce equity.",
		Explanation:       "Expenses increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_ap", Text: "Accounts Payable", ErrorTag: "confused_with_payable"},
		{ID: "opt_service_rev", Text: "Service Revenue"},
	}, seed, 605)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Because rent was paid immediately in cash, what account is credited?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Cash was paid out of the company's account.",
		Explanation:       "Cash is credited because funds were disbursed.",
		RelevantConceptID: "cash_vs_expense",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_rent_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rent", Text: fmt.Sprintf("Debit Cash %s / Credit Rent Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
		{ID: "opt_dr_rent_cr_ap", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "confused_with_payable"},
	}, seed, 606)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying cash rent?",
		CorrectOptionID:   "opt_dr_rent_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit Rent Expense (expense up) and Credit Cash (asset down).",
		Explanation:       fmt.Sprintf("Debit Rent Expense %s and Credit Cash %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does paying cash rent affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Rent Expense); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_no_net_change", Text: "No change; expenses do not affect balance sheet."},
		}, seed, 607),
		CausalHint:        "Cash decreased assets; rent expense decreased equity.",
		Explanation:       fmt.Sprintf("Assets decrease by %s and Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_expense",
	}
}

func buildBorrowCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_notes_payable", Text: "Notes Payable"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: "treated_borrowing_as_revenue"},
	}, seed, 701)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the cash borrowed and deposited today?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "Funds were borrowed and deposited into the bank account.",
		Explanation:       "Cash is received and debited.",
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
		CausalHint:        "Cash is an asset.",
		Explanation:       "Cash is an Asset.",
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
		CausalHint:        "Money came into the company.",
		Explanation:       "Cash increases.",
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
		CausalHint:        "Assets increase by Debit.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: "treated_borrowing_as_revenue"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
	}, seed, 705)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What counter-account records the formal obligation to repay the lender?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage5Opts,
		CausalHint:        "Borrowing money creates a liability (debt to be repaid), not earned revenue.",
		Explanation:       "Notes Payable is credited for the promissory note liability.",
		RelevantConceptID: "notes_payable_classification",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_borrowing_as_revenue"},
		{ID: "opt_dr_notes_cr_cash", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
	}, seed, 706)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for borrowing cash on a note?",
		CorrectOptionID:   "opt_dr_cash_cr_notes",
		Options:           stage6Opts,
		CausalHint:        "Debit Cash (asset up) and Credit Notes Payable (liability up).",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Notes Payable %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does borrowing cash affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Notes Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Revenue); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_borrowing_as_revenue"},
			{ID: "opt_no_net_change", Text: "No net change in total assets."},
		}, seed, 707),
		CausalHint:        "Cash increased assets; debt increased liabilities. Has any equity changed?",
		Explanation:       fmt.Sprintf("Assets increase by %s and Liabilities increase by %s. Equity is unaffected.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildIssueSharesStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	// Stage 1: Identify Account
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_common_stock", Text: "Common Stock"},
		{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: "treated_equity_as_revenue"},
	}, seed, 801)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the cash received from the owner investment?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "Investors paid cash into the company.",
		Explanation:       "Cash is received and debited.",
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
		CausalHint:        "Cash is an asset.",
		Explanation:       "Cash is an Asset.",
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
		CausalHint:        "The company receives funds from investors.",
		Explanation:       "Cash increases.",
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
		CausalHint:        "Assets increase by Debit.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 5: Counter-Account
	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: "treated_equity_as_revenue"},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)"},
	}, seed, 805)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What counter-account records the equity interest issued to the owners?",
		CorrectOptionID:   "opt_common_stock",
		Options:           stage5Opts,
		CausalHint:        "Owner investment provides contributed capital (Common Stock), not earned revenue from business operations.",
		Explanation:       "Common Stock is credited for contributed capital.",
		RelevantConceptID: "common_stock_classification",
	}

	// Stage 6: Balanced Entry
	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_stock", Text: fmt.Sprintf("Debit Cash %s / Credit Common Stock %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "treated_equity_as_revenue"},
		{ID: "opt_dr_stock_cr_cash", Text: fmt.Sprintf("Debit Common Stock %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "reversed_sides"},
	}, seed, 806)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for issuing common stock for cash?",
		CorrectOptionID:   "opt_dr_cash_cr_stock",
		Options:           stage6Opts,
		CausalHint:        "Debit Cash (asset up) and Credit Common Stock (equity up).",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Common Stock %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	// Stage 7: Equation Effect
	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing common shares affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Equity increases by %s (+Common Stock); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s; Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_no_net_change", Text: "No net change in total assets."},
		}, seed, 807),
		CausalHint:        "Cash increases assets; Common Stock increases equity.",
		Explanation:       fmt.Sprintf("Assets increase by %s and Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "cash_vs_revenue",
	}
}

func buildPrepaidPurchaseStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: "confused_with_payable"},
		{ID: "opt_cash", Text: "Cash"},
	}, seed, 901)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the advance payment for 12 months of insurance coverage?",
		CorrectOptionID:   "opt_prepaid_insurance",
		Options:           stage1Opts,
		CausalHint:        "The insurance policy covers 12 future months. An advance payment for future economic benefits is an asset, not an immediate operating expense.",
		Explanation:       "Prepaid Insurance is an asset account debited when paying for future coverage in advance.",
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
		CausalHint:        "Prepaid Insurance represents an economic resource (the legal right to future protection) owned and controlled by the company.",
		Explanation:       "Prepaid Insurance is an Asset.",
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
		CausalHint:        "A new 12-month policy was purchased and added to the company's prepaid resources.",
		Explanation:       "Prepaid Insurance increases upon purchasing the policy.",
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
		CausalHint:        "Assets have a normal debit balance; increases go on the debit (left) side.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_insurance_expense", Text: "Insurance Expense", ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: "confused_with_payable"},
		{ID: "opt_service_revenue", Text: "Service Revenue"},
	}, seed, 905)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Because the policy was paid immediately in cash, what counter-account is credited?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Cash was disbursed from the company's bank account today.",
		Explanation:       "Cash is an asset that decreased and is credited.",
		RelevantConceptID: "prepaid_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_prepaid_cr_cash", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
		{ID: "opt_dr_cash_cr_prepaid", Text: fmt.Sprintf("Debit Cash %s / Credit Prepaid Insurance %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_prepaid_cr_ap", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "confused_with_payable"},
	}, seed, 906)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for purchasing this prepaid insurance policy?",
		CorrectOptionID:   "opt_dr_prepaid_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit the asset that increased (Prepaid Insurance) and credit the asset disbursed (Cash).",
		Explanation:       fmt.Sprintf("Debit Prepaid Insurance %s and Credit Cash %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does purchasing prepaid insurance affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Prepaid Insurance increases (+%s) and Cash decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Insurance Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnPrepaidPurchase},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars())},
		}, seed, 907),
		CausalHint:        "One asset increased (Prepaid Insurance) and another asset decreased (Cash) by the same amount. No expense has expired yet!",
		Explanation:       fmt.Sprintf("This is an asset exchange: Prepaid Insurance increases by %s, Cash decreases by %s. Total Assets and Equity are unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "prepaid_vs_expense",
	}
}

func buildPrepaidConsumptionStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_insurance_expense", Text: "Insurance Expense"},
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_on_prepaid_expiration"},
		{ID: "opt_rent_expense", Text: "Rent Expense"},
	}, seed, 1001)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "One month of insurance coverage has expired. Which expense account recognizes the cost incurred for the month?",
		CorrectOptionID:   "opt_insurance_expense",
		Options:           stage1Opts,
		CausalHint:        "The economic benefit of one month of coverage has been consumed. Expenses reflect consumed assets.",
		Explanation:       "Insurance Expense is debited for the cost of insurance coverage used during the period.",
		RelevantConceptID: "insurance_expense_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense (Equity reduction)"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue"},
	}, seed, 1002)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Insurance Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Insurance Expense represents expired asset costs, which reduce owner equity.",
		Explanation:       "Insurance Expense is an Expense account.",
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
		CausalHint:        "Incurring an expense increases the total expenses accumulated for the period.",
		Explanation:       "Expenses increase as costs are recognized.",
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
		CausalHint:        "Expenses have normal debit balances; increases are recorded on the debit (left) side.",
		Explanation:       "Expenses increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_prepaid_insurance", Text: "Prepaid Insurance"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: "cash_recorded_on_prepaid_expiration"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
		{ID: "opt_service_revenue", Text: "Service Revenue"},
	}, seed, 1005)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account holding the unexpired coverage must be reduced as coverage expires?",
		CorrectOptionID:   "opt_prepaid_insurance",
		Options:           stage5Opts,
		CausalHint:        "No cash is paid when insurance expires; the policy was already paid for in advance. The prepaid asset balance is reduced as coverage expires.",
		Explanation:       "Prepaid Insurance is credited to reduce the remaining asset balance.",
		RelevantConceptID: "prepaid_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_expense_cr_prepaid", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Prepaid Insurance %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Insurance Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "cash_recorded_on_prepaid_expiration"},
		{ID: "opt_dr_prepaid_cr_expense", Text: fmt.Sprintf("Debit Prepaid Insurance %s / Credit Insurance Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1006)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for recording this expired insurance?",
		CorrectOptionID:   "opt_dr_expense_cr_prepaid",
		Options:           stage6Opts,
		CausalHint:        "Debit Insurance Expense (expense up) and Credit Prepaid Insurance (asset down).",
		Explanation:       fmt.Sprintf("Debit Insurance Expense %s and Credit Prepaid Insurance %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recognizing expired insurance affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Prepaid Insurance); Equity decreases by %s (-Insurance Expense); Cash is unaffected.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_cash", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Insurance Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: "cash_recorded_on_prepaid_expiration"},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only.", ErrorTag: "confused_with_asset_swap"},
		}, seed, 1007),
		CausalHint:        "Assets decrease because Prepaid Insurance expired. Equity decreases because an expense occurred. Cash was unaffected today.",
		Explanation:       fmt.Sprintf("Assets decrease by %s (-Prepaid Insurance) and Equity decreases by %s (-Insurance Expense). Zero cash moved today.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "prepaid_vs_expense",
	}
}

func buildEquipmentPurchaseCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_equipment", Text: "Equipment"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
	}, seed, 1101)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the productive machinery and equipment acquired today?",
		CorrectOptionID:   "opt_equipment",
		Options:           stage1Opts,
		CausalHint:        "Equipment has multi-year productive utility. Long-term productive resources are capitalized as assets, not expensed immediately.",
		Explanation:       "Equipment is an asset account debited for acquired physical machinery.",
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
		CausalHint:        "Equipment is a tangible, long-term productive economic resource owned and controlled by the company.",
		Explanation:       "Equipment is a non-current Asset.",
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
		CausalHint:        "Purchasing new machinery increases the total equipment owned.",
		Explanation:       "Equipment increases when new assets are purchased.",
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
		CausalHint:        "Assets have a normal debit balance; increases are recorded on the debit (left) side.",
		Explanation:       "Assets increase by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
		{ID: "opt_common_stock", Text: "Common Stock"},
	}, seed, 1105)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Because the equipment was purchased for cash, what account is credited?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Cash was disbursed from the company bank account.",
		Explanation:       "Cash is credited because financial assets were spent.",
		RelevantConceptID: "capex_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_equip_cr_cash", Text: fmt.Sprintf("Debit Equipment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
		{ID: "opt_dr_cash_cr_equip", Text: fmt.Sprintf("Debit Cash %s / Credit Equipment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1106)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for this cash equipment purchase?",
		CorrectOptionID:   "opt_dr_equip_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit Equipment (asset up) and Credit Cash (asset down).",
		Explanation:       fmt.Sprintf("Debit Equipment %s and Credit Cash %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does purchasing equipment for cash affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset exchange: Equipment increases (+%s) and Cash decreases (-%s); Total Assets and Equity unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnEquipmentPurchase},
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars())},
		}, seed, 1107),
		CausalHint:        "Buying equipment is a capital expenditure (asset swap), NOT an operating expense. Total assets and equity do not change.",
		Explanation:       fmt.Sprintf("This is an asset exchange: Equipment increases by %s and Cash decreases by %s. Total assets and equity remain unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "capex_vs_expense",
	}
}

func buildRepayNotePrincipalStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_service_revenue", Text: "Service Revenue"},
	}, seed, 1201)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which liability account is reduced when repaying principal on the bank promissory note?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage1Opts,
		CausalHint:        "Repaying loan principal settles a previously recorded debt obligation (liability).",
		Explanation:       "Notes Payable is debited to reduce the outstanding loan balance.",
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
		CausalHint:        "Notes Payable represents a formal written debt obligation to a creditor.",
		Explanation:       "Notes Payable is a Liability.",
		RelevantConceptID: "notes_payable_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 1203)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the liability Notes Payable increase or decrease upon principal repayment?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "Paying off debt reduces the remaining obligation owed to the bank.",
		Explanation:       "Notes Payable decreases when principal is paid.",
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
		CausalHint:        "Liabilities have a normal credit balance; decreasing a liability requires a Debit.",
		Explanation:       "Liabilities decrease by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
		{ID: "opt_service_revenue", Text: "Service Revenue"},
	}, seed, 1205)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What account reflects the cash disbursed to settle the loan principal?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Funds flowed out of the business bank account.",
		Explanation:       "Cash is credited because money was disbursed.",
		RelevantConceptID: "principal_repayment_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_notes_cr_cash", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1206)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for repaying note principal with cash?",
		CorrectOptionID:   "opt_dr_notes_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit Notes Payable (liability down) and Credit Cash (asset down).",
		Explanation:       fmt.Sprintf("Debit Notes Payable %s and Credit Cash %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does repaying loan principal affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_liab_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Liabilities decrease by %s (-Notes Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Expense).", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
			{ID: "opt_no_net_change", Text: "No net change in total assets; asset swap only."},
		}, seed, 1207),
		CausalHint:        "Repaying loan principal settles debt (liabilities down) with money (assets down). Principal reduction is NOT an expense!",
		Explanation:       fmt.Sprintf("Assets decrease by %s (-Cash) and Liabilities decrease by %s (-Notes Payable). Equity is unaffected.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "principal_repayment_vs_expense",
	}
}

func buildDividendCashStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dividends", Text: "Dividends"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_service_revenue", Text: "Service Revenue"},
		{ID: "opt_cash", Text: "Cash"},
	}, seed, 1301)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records cash distributions paid directly to stockholders?",
		CorrectOptionID:   "opt_dividends",
		Options:           stage1Opts,
		CausalHint:        "Dividends represent a direct distribution of corporate earnings to stockholders, not an operating expense of doing business.",
		Explanation:       "Dividends is debited to track distributions of profits to shareholders.",
		RelevantConceptID: "dividends_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dividends_cat", Text: "Dividends (Equity reduction)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_asset", Text: "Asset"},
	}, seed, 1302)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Dividends?",
		CorrectOptionID:   "opt_dividends_cat",
		Options:           stage2Opts,
		CausalHint:        "Dividends distribute accumulated earnings directly to owners, reducing Equity.",
		Explanation:       "Dividends is a distribution account that directly reduces Equity.",
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
		CausalHint:        "The accumulated total of dividends distributed during the period increases.",
		Explanation:       "Dividends increases as new distributions are made.",
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
		CausalHint:        "Dividends reduce equity, so they have a normal debit balance and increase on the debit side.",
		Explanation:       "Dividends increases by Debit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
	}, seed, 1305)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Because dividends were paid in cash, what account is credited?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "Cash was disbursed to stockholders.",
		Explanation:       "Cash is credited because cash was distributed.",
		RelevantConceptID: "dividend_vs_expense",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_div_cr_cash", Text: fmt.Sprintf("Debit Dividends %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit Rent Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnDividend},
		{ID: "opt_dr_cash_cr_div", Text: fmt.Sprintf("Debit Cash %s / Credit Dividends %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1306)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying cash dividends?",
		CorrectOptionID:   "opt_dr_div_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit Dividends (equity reduction) and Credit Cash (asset down).",
		Explanation:       fmt.Sprintf("Debit Dividends %s and Credit Cash %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "debit_credit_translation",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does paying cash dividends affect the accounting equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Dividends); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_no_net_change", Text: "No change in equity; dividends are an asset.", ErrorTag: "confused_with_asset_swap"},
		}, seed, 1307),
		CausalHint:        "Dividends reduce assets (-Cash) and reduce equity (-Dividends). Dividends are NOT an expense on the income statement.",
		Explanation:       fmt.Sprintf("Assets decrease by %s (-Cash) and Equity decreases by %s (-Dividends). Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars()),
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

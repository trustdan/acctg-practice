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

	// Copy and populate secondary parameters aligned by amount index
	resolvedParams := make(map[string]int64, len(paramValues))
	for k, v := range paramValues {
		resolvedParams[k] = v
	}

	amtIdx := -1
	if amounts, ok := q.Parameters["amount_minor_units"]; ok {
		for i, a := range amounts {
			if a == amtVal {
				amtIdx = i
				break
			}
		}
	}

	// Replace scenario template placeholders
	promptText := q.ScenarioTemplate
	promptText = strings.ReplaceAll(promptText, "${amount_dollars}", money.FormatDollars())
	promptText = strings.ReplaceAll(promptText, "${amount_exact}", money.FormatExact())

	if amtIdx >= 0 {
		for paramKey, paramVals := range q.Parameters {
			if paramKey != "amount_minor_units" && amtIdx < len(paramVals) {
				pVal := paramVals[amtIdx]
				if _, exists := resolvedParams[paramKey]; !exists {
					resolvedParams[paramKey] = pVal
				}
				pMoney := domain.NewMoney(pVal)
				baseName := strings.TrimSuffix(paramKey, "_minor_units")
				promptText = strings.ReplaceAll(promptText, fmt.Sprintf("${%s_dollars}", baseName), pMoney.FormatDollars())
				promptText = strings.ReplaceAll(promptText, fmt.Sprintf("${%s_exact}", baseName), pMoney.FormatExact())
			}
		}
	}

	// Process canonical event in engine
	event := engine.TransactionEvent{
		FamilyID:   q.FamilyID,
		Parameters: resolvedParams,
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
		Parameters:    resolvedParams,
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
	case bank.FamilyBadDebtAllowance:
		buildBadDebtAllowanceStages(instance, money, seed, processed)
	case bank.FamilyIssueBondsPar:
		buildIssueBondsParStages(instance, money, seed, processed)
	case bank.FamilyAcquireGoodwill:
		buildAcquireGoodwillStages(instance, money, seed, processed)
	case bank.FamilyCreditSalesCashFlow:
		buildCreditSalesCashFlowStages(instance, money, seed, processed)
	case bank.FamilyWriteOffReceivable:
		buildWriteOffReceivableStages(instance, money, seed, processed)
	case bank.FamilyIssueBondsDiscount:
		buildIssueBondsDiscountStages(instance, money, seed, processed)
	case bank.FamilyCapitalizeSoftwareFeasibility:
		buildCapitalizeSoftwareFeasibilityStages(instance, money, seed, processed)
	case bank.FamilyCurrentRatioAnalysis:
		buildCurrentRatioAnalysisStages(instance, money, seed, processed)
	case bank.FamilyRepurchaseTreasuryStock:
		buildRepurchaseTreasuryStockStages(instance, money, seed, processed)
	case bank.FamilyPayBondInterest:
		buildPayBondInterestStages(instance, money, seed, processed)
	case bank.FamilyRepo105Analysis:
		buildRepo105AnalysisStages(instance, money, seed, processed)
	case bank.FamilyAppleMultiElementRevenue:
		buildAppleMultiElementRevenueStages(instance, money, seed, processed)
	case bank.FamilyGoodwillImpairment:
		buildGoodwillImpairmentStages(instance, money, seed, processed)
	case bank.FamilyDigitalAssetAccounting:
		buildDigitalAssetAccountingStages(instance, money, seed, processed)
	case bank.FamilyReportingPolicyChoice:
		buildReportingPolicyChoiceStages(instance, money, seed, processed)
	case bank.FamilyFinancialForecasting:
		buildFinancialForecastingStages(instance, money, seed, processed)
	case bank.FamilyM03OCFAccountsPayableIncrease:
		BuildM03OCFAccountsPayableIncreaseStages(instance, money, seed, processed)
	case bank.FamilyM03OCFPrepaidIncrease:
		BuildM03OCFPrepaidIncreaseStages(instance, money, seed, processed)
	case bank.FamilyM03OCFUnearnedRevenueIncrease:
		BuildM03OCFUnearnedRevenueIncreaseStages(instance, money, seed, processed)
	case bank.FamilyM03OCFQuickRatio:
		BuildM03OCFQuickRatioStages(instance, money, seed, processed)
	case bank.FamilyM04VariableConsideration:
		buildM04VariableConsiderationStages(instance, money, seed, processed)
	case bank.FamilyM04ContractAssetAR:
		buildM04ContractAssetARStages(instance, money, seed, processed)
	case bank.FamilyM04RightOfReturn:
		buildM04RightOfReturnStages(instance, money, seed, processed)
	case bank.FamilyM04SatisfactionOverTime:
		buildM04SatisfactionOverTimeStages(instance, money, seed, processed)
	case bank.FamilyM05InventoryFIFO:
		BuildM05InventoryFIFOStages(instance, money, seed, processed)
	case bank.FamilyM05InventoryLIFO:
		BuildM05InventoryLIFOStages(instance, money, seed, processed)
	case bank.FamilyM05InventoryNRVWriteDown:
		BuildM05InventoryNRVStages(instance, money, seed, processed)
	case bank.FamilyM05FreightInCapitalize:
		BuildM05FreightInCapStages(instance, money, seed, processed)
	case bank.FamilyM06StraightLineDepr, bank.FamilyM06MACRSTaxTiming, bank.FamilyM06AssetDisposalGainLoss, bank.FamilyM06ImpairmentRecoverability:
		BuildM06Stages(q.FamilyID, instance, money, seed, processed)
	case bank.FamilyM07BondIssuanceDiscount, bank.FamilyM07BondIssuancePremium, bank.FamilyM07BondEffectiveInterestAmort, bank.FamilyM07BondContingentLiability:
		BuildM07Stages(q.FamilyID, instance, money, seed, processed)
	case bank.FamilyM08LeaseCommencementROUAsset, bank.FamilyM08LeasePeriodicPayment, bank.FamilyM08TreasuryStockRepurchase, bank.FamilyM08TreasuryStockReissuance:
		BuildM08Stages(q.FamilyID, instance, money, seed, processed)
	case bank.FamilyM09CookieJarReserve, bank.FamilyM09WarrantyAccrual:
		BuildM09Stages(q.FamilyID, instance, money, seed, processed)
	case bank.FamilyM10ForecastingOperatingExpenseAP, bank.FamilyM10ForecastingCapEx, bank.FamilyM10ForecastingDebtPlug:
		BuildM10Stages(q.FamilyID, instance, money, seed, processed)
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

func buildBadDebtAllowanceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_bad_debt_expense", Text: "Bad Debt Expense"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable", ErrorTag: engine.TagDirectWriteOffExpense},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagWrongAccount},
	}, seed, 1401)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the estimated uncollectible customer receivables for the period?",
		CorrectOptionID:   "opt_bad_debt_expense",
		Options:           stage1Opts,
		CausalHint:        "Under the matching principle, which expense account captures the cost of credit customers failing to pay?",
		Explanation:       "Bad Debt Expense records the estimated uncollectible accounts in the period the sales occurred. Offsetting Accounts Receivable directly would use the direct write-off method, which is not GAAP.",
		RelevantConceptID: "bad_debt_allowance",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1402)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Bad Debt Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does this account track a business operating cost, a resource, or an obligation?",
		Explanation:       "Bad Debt Expense is an Expense that reduces net income and Stockholders' Equity.",
		RelevantConceptID: "bad_debt_allowance",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1403)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Bad Debt Expense increase or decrease when recording this year-end estimate?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recording an expense add to the expenses accumulated this period?",
		Explanation:       "Bad Debt Expense increases as additional estimated uncollectible costs are recorded for the period.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1404)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase in an Expense recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses have a normal debit balance. Which side records an increase?",
		Explanation:       "Expenses increase on the Debit (left) side, which in turn reduces Stockholders' Equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_allowance", Text: "Allowance for Doubtful Accounts (Contra-Asset)"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable (Asset)", ErrorTag: engine.TagDirectWriteOffExpense},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagWrongAccount},
	}, seed, 1405)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account balances Bad Debt Expense under the allowance method?",
		CorrectOptionID:   "opt_allowance",
		Options:           stage5Opts,
		CausalHint:        "We do not yet know which specific customers will not pay. Which contra-asset reserve account is credited?",
		Explanation:       "Allowance for Doubtful Accounts is credited. It is a contra-asset that offsets Accounts Receivable on the balance sheet without removing specific customer balances before they default.",
		RelevantConceptID: "bad_debt_allowance",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_exp_cr_allow", Text: fmt.Sprintf("Debit Bad Debt Expense %s / Credit Allowance for Doubtful Accounts %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_ar", Text: fmt.Sprintf("Debit Bad Debt Expense %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDirectWriteOffExpense},
		{ID: "opt_dr_allow_cr_exp", Text: fmt.Sprintf("Debit Allowance for Doubtful Accounts %s / Credit Bad Debt Expense %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1406)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced adjusting entry under the allowance method?",
		CorrectOptionID:   "opt_dr_exp_cr_allow",
		Options:           stage6Opts,
		CausalHint:        "Under the allowance method, debit the expense and credit the contra-asset reserve.",
		Explanation:       fmt.Sprintf("Debit Bad Debt Expense %s and Credit Allowance for Doubtful Accounts %s. Direct write-off to Accounts Receivable is rejected by GAAP matching.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "bad_debt_allowance",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording estimated uncollectible accounts affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Allowance for Doubtful Accounts); Equity decreases by %s (-Bad Debt Expense); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_net_change", Text: "No net change in assets; expenses only affect net income.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1407),
		CausalHint:        "The allowance is a contra-asset (credit decreases total assets), and expense reduces equity.",
		Explanation:       fmt.Sprintf("Net assets decrease by %s through the contra-asset Allowance for Doubtful Accounts, and Equity decreases by %s via Bad Debt Expense.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "bad_debt_allowance",
	}
}

func buildIssueBondsParStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash"},
		{ID: "opt_bonds_payable", Text: "Bonds Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagStockIssuedInsteadOfBond},
	}, seed, 1501)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "What did the company receive today from bondholders, and which account records it?",
		CorrectOptionID:   "opt_cash",
		Options:           stage1Opts,
		CausalHint:        "What liquid resource entered the business upon issuing the bonds?",
		Explanation:       "Cash records the money received from issuing bonds to investors.",
		RelevantConceptID: "cash_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1502)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Cash?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is cash a resource the company owns and controls?",
		Explanation:       "Cash is an Asset because the company owns and can use it.",
		RelevantConceptID: "cash_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1503)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Cash increase or decrease upon issuing bonds for cash?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Did money flow into the company from bond investors?",
		Explanation:       "Cash increases because the company received funds from the bond issue.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1504)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Cash) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have a normal debit balance. Which side records an increase?",
		Explanation:       "Assets increase on the Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_bonds_payable", Text: "Bonds Payable (Liability)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagStockIssuedInsteadOfBond},
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)", ErrorTag: engine.TagRevenueRecordedOnBorrowing},
	}, seed, 1505)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account balances the cash received from bond investors?",
		CorrectOptionID:   "opt_bonds_payable",
		Options:           stage5Opts,
		CausalHint:        "Did investors purchase debt notes with repayment obligations or equity ownership shares?",
		Explanation:       "Bonds Payable is a long-term liability recording the principal obligation owed to bondholders. Common Stock represents equity ownership, not debt.",
		RelevantConceptID: "bonds_payable_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_bonds", Text: fmt.Sprintf("Debit Cash %s / Credit Bonds Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_stock", Text: fmt.Sprintf("Debit Cash %s / Credit Common Stock %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagStockIssuedInsteadOfBond},
		{ID: "opt_dr_bonds_cr_cash", Text: fmt.Sprintf("Debit Bonds Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1506)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for issuing bonds at par?",
		CorrectOptionID:   "opt_dr_cash_cr_bonds",
		Options:           stage6Opts,
		CausalHint:        "Debit the asset received and credit the long-term liability obligation incurred.",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Bonds Payable %s. Bonds are debt financing, not equity.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "bonds_payable_classification",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing bonds at par affect the accounting equation?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Bonds Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_liab_up_eq_down", Text: fmt.Sprintf("Liabilities increase by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1507),
		CausalHint:        "Cash increases assets, and Bonds Payable increases liabilities. No equity changes occur upon borrowing.",
		Explanation:       fmt.Sprintf("Assets increase by %s (+Cash) and liabilities increase by %s (+Bonds Payable). Equity is unaffected by borrowing principal.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "bonds_payable_classification",
	}
}

func buildAcquireGoodwillStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_goodwill", Text: "Goodwill"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagGoodwillExpensed},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 1601)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which asset account records the excess purchase price paid over the fair value of net identifiable assets acquired?",
		CorrectOptionID:   "opt_goodwill",
		Options:           stage1Opts,
		CausalHint:        "What intangible asset account captures the premium paid for an established business's reputation, workforce, and synergies?",
		Explanation:       "Goodwill is capitalized as an intangible asset representing the premium paid over the fair value of identifiable net assets acquired.",
		RelevantConceptID: "goodwill_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagGoodwillExpensed},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_liability", Text: "Liability"},
	}, seed, 1602)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Goodwill?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is Goodwill an economic resource expected to provide future benefits, or an immediate period expense?",
		Explanation:       "Goodwill is an indefinite-lived intangible Asset, not an immediate operating expense.",
		RelevantConceptID: "goodwill_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1603)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Goodwill increase or decrease upon completing this business acquisition?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Was a new intangible asset acquired today?",
		Explanation:       "Goodwill increases because the company acquired a new business unit.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1604)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Goodwill) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Which side is an asset's normal balance?",
		Explanation:       "Assets increase by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
		{ID: "opt_accounts_payable", Text: "Accounts Payable (Liability)"},
	}, seed, 1605)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which other account changes when cash is paid to complete the acquisition?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What resource did the buyer transfer to the seller to complete the purchase?",
		Explanation:       "Cash decreases because the purchase was paid in cash, so Cash is credited.",
		RelevantConceptID: "goodwill_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_goodwill_cr_cash", Text: fmt.Sprintf("Debit Goodwill %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_cash", Text: fmt.Sprintf("Debit an expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagGoodwillExpensed},
		{ID: "opt_dr_cash_cr_goodwill", Text: fmt.Sprintf("Debit Cash %s / Credit Goodwill %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1606)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for acquiring Goodwill paid in cash?",
		CorrectOptionID:   "opt_dr_goodwill_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Debit the new intangible asset acquired and credit the cash resource paid.",
		Explanation:       fmt.Sprintf("Debit Goodwill %s and Credit Cash %s. Goodwill is capitalized, not expensed.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "goodwill_classification",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does acquiring goodwill for cash affect the balance sheet equation?",
		CorrectOptionID: "opt_asset_swap",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Net Assets change is $0 (+Goodwill %s, -Cash %s); Liabilities and Equity are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1607),
		CausalHint:        "One asset (Goodwill) increased while another asset (Cash) decreased by the same amount.",
		Explanation:       fmt.Sprintf("This is an asset swap: Goodwill increases by %s and Cash decreases by %s. Net Assets, Liabilities, and Equity remain completely unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "goodwill_classification",
	}
}

func buildCreditSalesCashFlowStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_service_rev", Text: "Service Revenue"},
	}, seed, 1701)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which account records the customer claim when services are completed on credit?",
		CorrectOptionID:   "opt_ar",
		Options:           stage1Opts,
		CausalHint:        "The customer was billed but has not paid. Which account tracks uncollected customer receivables?",
		Explanation:       "Accounts Receivable records the claim to future cash. No cash was received today.",
		RelevantConceptID: "accounts_receivable_classification",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1702)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Receivable?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is a claim to future cash an economic resource owned by the company?",
		Explanation:       "Accounts Receivable is an Asset representing future economic benefits.",
		RelevantConceptID: "accounts_receivable_classification",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1703)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Receivable increase or decrease when performing services on account?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does completing unpaid work add to the customer claims owed to the company?",
		Explanation:       "Accounts Receivable increases as a new customer invoice is issued.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1704)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Accounts Receivable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have a normal debit balance. Which side records an increase?",
		Explanation:       "Assets increase by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_rev", Text: "Service Revenue (Revenue)"},
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Liability)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 1705)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account balances Accounts Receivable when service performance is complete today?",
		CorrectOptionID:   "opt_service_rev",
		Options:           stage5Opts,
		CausalHint:        "The service work is done today. Which account records earning from completed work?",
		Explanation:       "Service Revenue is credited because the performance obligation was satisfied today.",
		RelevantConceptID: "service_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1706)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for services performed on credit?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Debit the receivable claim and credit earned revenue.",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s and Credit Service Revenue %s.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "service_revenue_classification",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this affect the equation, and how is it treated on the indirect Statement of Cash Flows?",
		CorrectOptionID: "opt_cf_deduct",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_cf_deduct", Text: fmt.Sprintf("Assets increase (+A/R %s); Equity increases (+Revenue %s). In the indirect SCF, +A/R is deducted from Net Income because no cash was collected.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_cf_add", Text: fmt.Sprintf("Operating cash flow increases by %s because revenue was earned.", amt.FormatDollars()), ErrorTag: engine.TagCashFlowOperatingWrongDirection},
			{ID: "opt_no_effect", Text: "No effect on the Statement of Cash Flows because no cash moved.", ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1707),
		CausalHint:        "Net Income includes this revenue, but no cash came in. In the indirect method operating section, do you add or subtract the increase in receivables?",
		Explanation:       fmt.Sprintf("Assets increase by %s (+A/R) and Equity increases by %s (+Revenue). On the Statement of Cash Flows (indirect method), because Net Income includes this %s of revenue but no cash was received, the increase in Accounts Receivable must be DEDUCTED from Net Income to calculate Operating Cash Flow.", amt.FormatDollars(), amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "statement_of_cash_flows_indirect",
	}
}

func buildWriteOffReceivableStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_allowance", Text: "Allowance for Doubtful Accounts"},
		{ID: "opt_bad_debt_expense", Text: "Bad Debt Expense", ErrorTag: engine.TagWriteOffExpensedAgain},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 1801)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which contra-asset account is debited when writing off a specific uncollectible customer balance under the allowance method?",
		CorrectOptionID:   "opt_allowance",
		Options:           stage1Opts,
		CausalHint:        "Bad Debt Expense was already recognized when estimating uncollectibles at period end. Which reserve account is debited to absorb this specific customer default?",
		Explanation:       "Allowance for Doubtful Accounts is the contra-asset reserve that absorbs specific customer defaults under the allowance method, matching bad debt expense in advance.",
		RelevantConceptID: "allowance_write_off",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset (Contra-Asset)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagWriteOffExpensedAgain},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1802)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Allowance for Doubtful Accounts?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Is the allowance an immediate income statement expense, or a balance sheet reserve that offsets Accounts Receivable?",
		Explanation:       "Allowance for Doubtful Accounts is a contra-asset account offsetting Accounts Receivable on the balance sheet.",
		RelevantConceptID: "allowance_write_off",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 1803)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Allowance for Doubtful Accounts reserve increase or decrease when absorbing this write-off?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "Does writing off an uncollectible customer account reduce the balance of the previously accumulated allowance reserve?",
		Explanation:       "Allowance for Doubtful Accounts decreases because the reserve is being used to write off a specific default.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1804)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is a decrease to a Contra-Asset (Allowance for Doubtful Accounts) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Contra-assets have a normal credit balance. Which side records a decrease to a contra-asset?",
		Explanation:       "Contra-assets have a normal credit balance, so a decrease is recorded with a Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ar", Text: "Accounts Receivable (Asset)"},
		{ID: "opt_bad_debt_expense", Text: "Bad Debt Expense (Expense)", ErrorTag: engine.TagWriteOffExpensedAgain},
		{ID: "opt_cash", Text: "Cash (Asset)"},
	}, seed, 1805)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account is credited to remove the specific customer's uncollectible balance from the books?",
		CorrectOptionID:   "opt_ar",
		Options:           stage5Opts,
		CausalHint:        "The customer will never pay. Which account tracked the specific uncollected invoice that must now be removed?",
		Explanation:       "Accounts Receivable is credited to eliminate the uncollectible customer account from the company's ledger.",
		RelevantConceptID: "allowance_write_off",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_allowance_cr_ar", Text: fmt.Sprintf("Debit Allowance for Doubtful Accounts %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_expense_cr_ar", Text: fmt.Sprintf("Debit Bad Debt Expense %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagWriteOffExpensedAgain},
		{ID: "opt_dr_ar_cr_allowance", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Allowance for Doubtful Accounts %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 1806)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to write off the specific uncollectible account?",
		CorrectOptionID:   "opt_dr_allowance_cr_ar",
		Options:           stage6Opts,
		CausalHint:        "Which entry reduces the contra-asset reserve and eliminates the uncollectible customer account?",
		Explanation:       fmt.Sprintf("Debit Allowance for Doubtful Accounts %s and Credit Accounts Receivable %s. Do not debit Bad Debt Expense again.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "allowance_write_off",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does writing off an uncollectible account affect the accounting equation and Net Realizable A/R?",
		CorrectOptionID: "opt_net_ar_unchanged",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_net_ar_unchanged", Text: fmt.Sprintf("Total Assets, Liabilities, and Equity change by $0. Gross A/R decreases by %s and the contra-asset decreases by %s, leaving Net Realizable A/R and Net Income completely unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagWriteOffExpensedAgain},
			{ID: "opt_liab_down_assets_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1807),
		CausalHint:        "Gross A/R went down, but the offsetting allowance also went down by the same amount. Did Net Realizable Value change?",
		Explanation:       fmt.Sprintf("This entry has a $0 net effect on total assets: Gross A/R decreases by %s, and Allowance for Doubtful Accounts decreases by %s. Net Realizable Accounts Receivable and Net Income are completely unchanged upon write-off.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "allowance_write_off",
	}
}

func buildIssueBondsDiscountStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	discountAmt := domain.NewMoney(amt.Cents() * 5 / 100)
	cashAmt := amt.Sub(discountAmt)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_discount", Text: "Discount on Bonds Payable"},
		{ID: "opt_interest_expense", Text: "Interest Expense", ErrorTag: engine.TagDiscountRecordedAsLoss},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagStockIssuedInsteadOfBond},
	}, seed, 1901)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which contra-liability account records the difference when bonds are issued for cash proceeds below face value?",
		CorrectOptionID:   "opt_discount",
		Options:           stage1Opts,
		CausalHint:        "Does issuing bonds at a discount represent an immediate loss, or a contra-liability account that offsets Bonds Payable to carrying value?",
		Explanation:       "Discount on Bonds Payable is a contra-liability account recording the discount below par value.",
		RelevantConceptID: "bond_discount_issuance",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability (Contra-Liability)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagDiscountRecordedAsLoss},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 1902)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Discount on Bonds Payable?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Is this account an immediate expense, or a valuation contra-account paired with the Bonds Payable liability?",
		Explanation:       "Discount on Bonds Payable is a Contra-Liability that offsets Bonds Payable on the balance sheet.",
		RelevantConceptID: "bond_discount_issuance",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 1903)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Discount on Bonds Payable increase or decrease upon bond issuance at a discount?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "A new bond contract is issued below par. Does the accumulated bond discount increase or decrease?",
		Explanation:       "Discount on Bonds Payable increases as the new discount is recorded upon issuance.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 1904)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to a Contra-Liability (Discount on Bonds Payable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have a normal credit balance. Which side is the normal balance of a contra-liability?",
		Explanation:       "Contra-liabilities have a normal Debit balance, the opposite of the normal credit balance of liabilities.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_bonds_payable", Text: "Bonds Payable (Liability)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)", ErrorTag: engine.TagStockIssuedInsteadOfBond},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)"},
	}, seed, 1905)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liability account is credited for the full face maturity value of the issued bonds?",
		CorrectOptionID:   "opt_bonds_payable",
		Options:           stage5Opts,
		CausalHint:        "Which long-term liability account records the principal maturity repayment obligation owed to bondholders?",
		Explanation:       "Bonds Payable is credited for the full par/face value of the debt obligation owed at maturity.",
		RelevantConceptID: "bond_discount_issuance",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_disc_cr_bonds", Text: fmt.Sprintf("Debit Cash %s, Debit Discount on Bonds Payable %s / Credit Bonds Payable %s", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_exp_cr_bonds", Text: fmt.Sprintf("Debit Cash %s, Debit Interest Expense %s / Credit Bonds Payable %s", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDiscountRecordedAsLoss},
		{ID: "opt_dr_cash_cr_bonds", Text: fmt.Sprintf("Debit Cash %s / Credit Bonds Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
	}, seed, 1906)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for issuing bonds at a discount?",
		CorrectOptionID:   "opt_dr_cash_disc_cr_bonds",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the cash proceeds, the contra-liability discount, and the full face value debt obligation?",
		Explanation:       fmt.Sprintf("Debit Cash %s, Debit Discount on Bonds Payable %s, and Credit Bonds Payable %s. The discount is amortized over time, not expensed at issuance.", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "bond_discount_issuance",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does issuing bonds at a discount affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_up_net_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_net_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Net Liabilities increase by %s (+Bonds Payable %s, -Discount %s); Equity is unchanged.", cashAmt.FormatDollars(), cashAmt.FormatDollars(), amt.FormatDollars(), discountAmt.FormatDollars())},
			{ID: "opt_assets_up_eq_down", Text: fmt.Sprintf("Assets increase by %s; Equity decreases by %s; Liabilities increase by %s.", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagDiscountRecordedAsLoss},
			{ID: "opt_assets_up_liab_up_par", Text: fmt.Sprintf("Assets increase by %s; Liabilities increase by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 1907),
		CausalHint:        "Cash received increases Assets. Bonds Payable increases liabilities, but Discount offsets it. Is Equity affected upon borrowing?",
		Explanation:       fmt.Sprintf("Assets increase by %s (+Cash); Net Liabilities increase by the net carrying value of %s (+Bonds Payable %s minus -Discount %s); Equity is completely unaffected by debt issuance proceeds.", cashAmt.FormatDollars(), cashAmt.FormatDollars(), amt.FormatDollars(), discountAmt.FormatDollars()),
		RelevantConceptID: "bond_discount_issuance",
	}
}

func buildCapitalizeSoftwareFeasibilityStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cap_software", Text: "Capitalized Software Costs"},
		{ID: "opt_rent_expense", Text: "Operating Expense / R&D Expense", ErrorTag: engine.TagExpensingFeasibleSoftware},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2001)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under US GAAP (ASC 985-20), which intangible asset account capitalizes software development coding costs incurred after technological feasibility?",
		CorrectOptionID:   "opt_cap_software",
		Options:           stage1Opts,
		CausalHint:        "Costs incurred AFTER technological feasibility is established are capitalized as an asset. Which account captures capitalized software costs?",
		Explanation:       "Capitalized Software Costs is an intangible asset account capturing software development costs incurred between technological feasibility and general commercial release.",
		RelevantConceptID: "software_capitalization",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset (Intangible Asset)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpensingFeasibleSoftware},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2002)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Capitalized Software Costs?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does capitalized software represent an economic resource providing future revenue, or an immediate period expense?",
		Explanation:       "Capitalized Software Costs is an intangible Asset, not an immediate operating expense.",
		RelevantConceptID: "software_capitalization",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2003)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Capitalized Software Costs increase or decrease as development expenditures are capitalized?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does the capitalized software asset balance increase as coding and testing costs are capitalized?",
		Explanation:       "Capitalized Software Costs increases as qualifying development costs are capitalized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 2004)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is an increase to an Asset (Capitalized Software Costs) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have a normal debit balance. Which side increases an asset?",
		Explanation:       "Assets increase by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)"},
	}, seed, 2005)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which account is credited when development costs are paid in cash?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid resource did the company pay out for developer compensation?",
		Explanation:       "Cash is credited because cash was paid to fund the development expenditures.",
		RelevantConceptID: "software_capitalization",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cap_cr_cash", Text: fmt.Sprintf("Debit Capitalized Software Costs %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_cash", Text: fmt.Sprintf("Debit an expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpensingFeasibleSoftware},
		{ID: "opt_dr_cash_cr_cap", Text: fmt.Sprintf("Debit Cash %s / Credit Capitalized Software Costs %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2006)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record software coding costs incurred post-feasibility?",
		CorrectOptionID:   "opt_dr_cap_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry capitalizes the software development asset and records the cash paid for development?",
		Explanation:       fmt.Sprintf("Debit Capitalized Software Costs %s and Credit Cash %s. Under US GAAP, feasibility mandates capitalization rather than immediate R&D expensing.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "software_capitalization",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does capitalizing software development costs affect the accounting equation?",
		CorrectOptionID: "opt_asset_swap_soft",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap_soft", Text: fmt.Sprintf("Asset swap: Capitalized Software Costs increases (+A %s) and Cash decreases (-A %s); Total Assets, Liabilities, and Equity remain unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpensingFeasibleSoftware},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2007),
		CausalHint:        "One asset increased while another asset decreased by the exact same amount. How does this asset exchange affect the accounting equation?",
		Explanation:       fmt.Sprintf("This is an asset swap: Capitalized Software Costs increases by %s and Cash decreases by %s. Net Assets, Liabilities, and Equity remain completely unchanged.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "software_capitalization",
	}
}

func buildCurrentRatioAnalysisStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_ap", Text: "Accounts Payable"},
		{ID: "opt_rent_expense", Text: "Operating Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue"},
	}, seed, 2101)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which liability account decreases when a retail company pays cash to satisfy outstanding supplier invoices?",
		CorrectOptionID:   "opt_ap",
		Options:           stage1Opts,
		CausalHint:        "Trade debts owed to suppliers for inventory or supplies are recorded in which current liability account?",
		Explanation:       "Accounts Payable records amounts owed to trade suppliers; paying supplier invoices settles this liability.",
		RelevantConceptID: "current_ratio_liquidity",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability (Current Liability)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2102)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Payable?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Is an obligation owed to trade suppliers an Asset, a Liability, or an Equity account?",
		Explanation:       "Accounts Payable is a Current Liability representing short-term claims by trade creditors.",
		RelevantConceptID: "current_ratio_liquidity",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_decrease", Text: "Decrease"},
		{ID: "opt_increase", Text: "Increase"},
	}, seed, 2103)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Accounts Payable increase or decrease when the company pays supplier balances in cash?",
		CorrectOptionID:   "opt_decrease",
		Options:           stage3Opts,
		CausalHint:        "Does paying supplier bills reduce the outstanding trade liability owed to creditors?",
		Explanation:       "Accounts Payable decreases as the debt to suppliers is satisfied.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)"},
	}, seed, 2104)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "How is a decrease to a Liability (Accounts Payable) recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities increase by Credit. Which side records a decrease?",
		Explanation:       "Liabilities decrease by Debit (left side).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable (Asset)"},
		{ID: "opt_common_stock", Text: "Common Stock (Equity)"},
	}, seed, 2105)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which current asset account decreases when supplier bills are paid?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid payment medium was transferred to suppliers to satisfy the invoice?",
		Explanation:       "Cash is credited because cash was disbursed to pay the suppliers.",
		RelevantConceptID: "current_ratio_liquidity",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ap_cr_cash", Text: fmt.Sprintf("Debit Accounts Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_exp_cr_cash", Text: fmt.Sprintf("Debit an expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_dr_cash_cr_ap", Text: fmt.Sprintf("Debit Cash %s / Credit Accounts Payable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2106)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying Accounts Payable in cash?",
		CorrectOptionID:   "opt_dr_ap_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the reduction to Accounts Payable alongside the cash outflow?",
		Explanation:       fmt.Sprintf("Debit Accounts Payable %s and Credit Cash %s. Settle the liability without recording a second expense.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "current_ratio_liquidity",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does paying Accounts Payable affect the equation, Working Capital, and the Current Ratio (assuming Current Ratio > 1.0)?",
		CorrectOptionID: "opt_cr_increases",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_cr_increases", Text: fmt.Sprintf("Assets decrease by %s (-Cash) and Liabilities decrease by %s (-Accounts Payable). Working Capital is unchanged, and the Current Ratio increases.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_cr_decreases", Text: fmt.Sprintf("Assets decrease by %s; the Current Ratio decreases because liquid cash was spent.", amt.FormatDollars()), ErrorTag: engine.TagCurrentRatioWrongDirection},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		}, seed, 2107),
		CausalHint:        "If a company has a Current Ratio greater than 1.0, what happens to the ratio when equal amounts are subtracted from both Current Assets and Current Liabilities?",
		Explanation:       fmt.Sprintf("Assets decrease (-Cash %s) and Liabilities decrease (-A/P %s). Working Capital (CA - CL) is unchanged. For any firm with Current Ratio > 1.0, subtracting equal amounts from numerator and denominator increases the Current Ratio.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "current_ratio_liquidity",
	}
}

func buildRepurchaseTreasuryStockStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_treasury_stock", Text: "Treasury Stock"},
		{ID: "opt_investment", Text: "Marketable Securities / Investment Asset", ErrorTag: engine.TagTreasuryStockRecordedAsAsset},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
	}, seed, 2201)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When a corporation buys back its own shares to hold in treasury, which contra-equity account is debited?",
		CorrectOptionID:   "opt_treasury_stock",
		Options:           stage1Opts,
		CausalHint:        "Does buying back company shares create an outside asset investment, or a contra-equity balance reducing stockholders' equity?",
		Explanation:       "Treasury Stock records the cost of the corporation's own shares repurchased from stockholders to hold in treasury.",
		RelevantConceptID: "treasury_stock_accounting",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_equity", Text: "Equity (Contra-Equity)"},
		{ID: "opt_asset", Text: "Asset", ErrorTag: engine.TagTreasuryStockRecordedAsAsset},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagExpenseRecordedOnLoanRepayment},
		{ID: "opt_liability", Text: "Liability"},
	}, seed, 2202)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Treasury Stock?",
		CorrectOptionID:   "opt_equity",
		Options:           stage2Opts,
		CausalHint:        "Can a corporation own an asset representing an ownership claim on itself, or does treasury stock offset total stockholders' equity?",
		Explanation:       "Treasury Stock is a contra-equity account that offsets stockholders' equity on the balance sheet.",
		RelevantConceptID: "treasury_stock_accounting",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2203)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Treasury Stock account balance increase or decrease when shares are repurchased?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does buying back more shares increase the accumulated cost of shares held in treasury?",
		Explanation:       "Treasury Stock increases because additional treasury shares are acquired at cost.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2204)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "To increase a contra-equity account with a normal debit balance, which side is recorded?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Treasury Stock has a normal debit balance (opposite to regular equity). To increase it, do you record on the left or right?",
		Explanation:       "Treasury Stock increases on the left side (Debit), which has the effect of reducing total stockholders' equity.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_retained_earnings", Text: "Retained Earnings (Equity)", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_notes_payable", Text: "Notes Payable (Liability)", ErrorTag: engine.TagWrongAccount},
	}, seed, 2205)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "What asset account is credited for the cash disbursed to repurchase the shares?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid asset was paid out to buy the shares back from market investors?",
		Explanation:       "Cash is credited because cash was disbursed to purchase the treasury shares.",
		RelevantConceptID: "cash_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ts_cr_cash", Text: fmt.Sprintf("Debit Treasury Stock %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_ts", Text: fmt.Sprintf("Debit Cash %s / Credit Treasury Stock %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_asset_cr_cash", Text: fmt.Sprintf("Debit Investment in Stock (Asset) %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagTreasuryStockRecordedAsAsset},
	}, seed, 2206)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for repurchasing common stock into treasury?",
		CorrectOptionID:   "opt_dr_ts_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the acquisition of treasury stock and the cash outflow?",
		Explanation:       fmt.Sprintf("Debit Treasury Stock %s and Credit Cash %s. Treasury shares are recorded at cost as contra-equity, not as an asset.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "treasury_stock_accounting",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does repurchasing treasury stock affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Treasury Stock); Liabilities are unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset swap: Cash decreases by %s and Investments increase by %s; Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagTreasuryStockRecordedAsAsset},
			{ID: "opt_liab_up", Text: fmt.Sprintf("Liabilities increase by %s; Assets decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2207),
		CausalHint:        "Cash paid reduces Assets. What happens to total Stockholders' Equity when Treasury Stock (a contra-equity account) increases?",
		Explanation:       fmt.Sprintf("Assets decrease (-Cash %s) and Equity decreases (-Treasury Stock %s); Liabilities are unaffected. Treasury stock reduces total stockholders' equity.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

func buildPayBondInterestStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	discountVal := inst.Parameters["discount_minor_units"]
	if discountVal <= 0 {
		discountVal = amt.Cents() * 20 / 100
	}
	cashVal := inst.Parameters["cash_minor_units"]
	if cashVal <= 0 {
		cashVal = amt.Cents() - discountVal
	}
	cashAmt := domain.NewMoney(cashVal)
	discountAmt := domain.NewMoney(discountVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_interest_expense", Text: "Interest Expense"},
		{ID: "opt_bonds_payable", Text: "Bonds Payable", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagWrongAccount},
	}, seed, 2301)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Which expense account records the total cost of borrowing, including cash interest and discount amortization?",
		CorrectOptionID:   "opt_interest_expense",
		Options:           stage1Opts,
		CausalHint:        "Does paying semi-annual bond interest reduce the principal liability owed at maturity, or record the periodic cost of borrowing as an expense?",
		Explanation:       "Interest Expense records the total cost of borrowed funds for the period, combining cash coupon paid and discount amortized.",
		RelevantConceptID: "bond_interest_expense",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense"},
		{ID: "opt_liability", Text: "Liability", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2302)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Interest Expense?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does the periodic cost of debt borrowing reduce equity as an expense, or represent the principal liability itself?",
		Explanation:       "Interest Expense is an expense account that reduces net income and retained earnings.",
		RelevantConceptID: "bond_interest_expense",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2303)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Interest Expense increase or decrease for the semi-annual period?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing financing costs for the period increase accumulated interest expense?",
		Explanation:       "Interest Expense increases as borrowing costs for the period are recognized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2304)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Expense account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses reduce equity and have normal debit balances. Do they increase on the left or right side?",
		Explanation:       "Expenses increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash_and_discount", Text: "Cash and Discount on Bonds Payable"},
		{ID: "opt_cash_only", Text: "Cash only (omitting discount)", ErrorTag: engine.TagIgnoringBondAmortization},
		{ID: "opt_bonds_payable", Text: "Bonds Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2305)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which accounts are credited when paying the coupon and amortizing the bond discount?",
		CorrectOptionID:   "opt_cash_and_discount",
		Options:           stage5Opts,
		CausalHint:        "What liquid asset is paid to bondholders, and which contra-liability account is amortized toward par?",
		Explanation:       "Cash is credited for the cash coupon paid, and Discount on Bonds Payable is credited to amortize the contra-liability toward face value.",
		RelevantConceptID: "discount_amortization",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_int_cr_cash_disc", Text: fmt.Sprintf("Debit Interest Expense %s, Credit Cash %s, and Credit Discount on Bonds Payable %s", amt.FormatDollars(), cashAmt.FormatDollars(), discountAmt.FormatDollars())},
		{ID: "opt_dr_int_cr_cash_only", Text: fmt.Sprintf("Debit Interest Expense %s and Credit Cash %s (omitting amortization)", cashAmt.FormatDollars(), cashAmt.FormatDollars()), ErrorTag: engine.TagIgnoringBondAmortization},
		{ID: "opt_reversed", Text: fmt.Sprintf("Debit Cash %s, Debit Discount %s, Credit Interest Expense %s", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2306)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for paying the bond coupon and amortizing the discount?",
		CorrectOptionID:   "opt_dr_int_cr_cash_disc",
		Options:           stage6Opts,
		CausalHint:        "Which entry records total interest expense, cash paid, and discount amortization?",
		Explanation:       fmt.Sprintf("Debit Interest Expense %s, Credit Cash %s, and Credit Discount on Bonds Payable %s. Total interest expense exceeds the cash coupon by the amount of discount amortized.", amt.FormatDollars(), cashAmt.FormatDollars(), discountAmt.FormatDollars()),
		RelevantConceptID: "bond_interest_expense",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does the bond interest payment and discount amortization affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_liab_up_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_liab_up_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Net Liabilities increase by %s (+Amortization); Equity decreases by %s (-Interest Expense).", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_eq_down_cash_only", Text: fmt.Sprintf("Assets decrease by %s (-Cash); Equity decreases by %s (-Expense); Liabilities are unchanged.", cashAmt.FormatDollars(), cashAmt.FormatDollars()), ErrorTag: engine.TagIgnoringBondAmortization},
			{ID: "opt_no_liab_effect", Text: fmt.Sprintf("Assets decrease by %s; Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2307),
		CausalHint:        "Cash paid reduces Assets. Interest expense reduces Equity. Crediting the contra-liability discount increases carrying liabilities. Do both sides balance?",
		Explanation:       fmt.Sprintf("Assets decrease by %s (-Cash); Net Liabilities increase by %s (+carrying value as discount shrinks); Equity decreases by %s (-Interest Expense). Equation balances: -%s = +%s + (-%s).", cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars(), cashAmt.FormatDollars(), discountAmt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

func buildRepo105AnalysisStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_notes_payable", Text: "Notes Payable / Short-Term Financing Liability"},
		{ID: "opt_revenue", Text: "Sales Revenue", ErrorTag: engine.TagRepoTreatedAsSale},
		{ID: "opt_trading_securities", Text: "Trading Securities (Asset Sale)", ErrorTag: engine.TagRepoTreatedAsSale},
	}, seed, 2401)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "Under US GAAP economic substance, which liability account records the obligation to repay cash borrowed under a repurchase agreement?",
		CorrectOptionID:   "opt_notes_payable",
		Options:           stage1Opts,
		CausalHint:        "Under economic reality, does pledging collateral for short-term cash represent an asset sale or a collateralized financing loan (Notes Payable)?",
		Explanation:       "Notes Payable records the contractual obligation to repay the borrowed cash, reflecting economic substance as collateralized financing.",
		RelevantConceptID: "repurchase_agreements",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagRepoTreatedAsSale},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2402)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Notes Payable in a repurchase agreement?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does an obligation to repurchase securities and repay cash represent a liability to creditors?",
		Explanation:       "Notes Payable is a Liability representing the obligation to repay the counterparty.",
		RelevantConceptID: "repurchase_agreements",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2403)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does the Notes Payable debt obligation increase or decrease upon receiving borrowing proceeds?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does borrowing cash under a repo contract increase outstanding debt liabilities?",
		Explanation:       "Notes Payable increases because a new borrowing obligation has been incurred.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2404)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase a Liability account with a normal credit balance?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have normal credit balances. Do they increase on the left or right side?",
		Explanation:       "Liabilities increase on the right side (Credit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_revenue", Text: "Sales Revenue", ErrorTag: engine.TagRepoTreatedAsSale},
		{ID: "opt_securities", Text: "Investment Securities", ErrorTag: engine.TagRepoTreatedAsSale},
	}, seed, 2405)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which asset account is debited for the cash received from the counterparty?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid asset was received from the counterparty upon pledging the collateral?",
		Explanation:       "Cash is debited because liquid cash funds were received.",
		RelevantConceptID: "cash_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_notes", Text: fmt.Sprintf("Debit Cash %s / Credit Notes Payable %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_securities", Text: fmt.Sprintf("Debit Cash %s / Credit Investment Securities (Sale) %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagRepoTreatedAsSale},
		{ID: "opt_reversed", Text: fmt.Sprintf("Debit Notes Payable %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2406)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry for recording the repurchase agreement under collateralized financing?",
		CorrectOptionID:   "opt_dr_cash_cr_notes",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the cash borrowing and the secured debt liability under collateralized financing?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Notes Payable %s. In Repo 105, Lehman aggressively booked this as a sale to temporarily remove assets and hide leverage.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "repurchase_agreements",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does transparently recording this borrowing affect the balance sheet and financial leverage?",
		CorrectOptionID: "opt_assets_up_liab_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Notes Payable); Equity is unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_asset_swap_hidden_leverage", Text: fmt.Sprintf("Asset swap: Cash increases and Securities decrease by %s; Liabilities and leverage appear reduced.", amt.FormatDollars()), ErrorTag: engine.TagRepoTreatedAsSale},
			{ID: "opt_equity_up", Text: fmt.Sprintf("Assets increase by %s and Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2407),
		CausalHint:        "Assets increase (+Cash) and Liabilities increase (+Notes Payable). How does transparently recognizing this borrowing affect the financial leverage ratio (Assets/Equity)?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Liabilities increase (+Notes Payable %s); Equity is unchanged. Total Assets expand, maintaining true leverage, unlike Lehman's off-balance-sheet sale treatment which cosmetically depressed reported leverage before quarter-end.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "leverage_ratio_analysis",
	}
}

func buildAppleMultiElementRevenueStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	deliveredVal := inst.Parameters["delivered_revenue_minor_units"]
	if deliveredVal <= 0 {
		deliveredVal = amt.Cents() * 80 / 100
	}
	unearnedVal := inst.Parameters["unearned_revenue_minor_units"]
	if unearnedVal <= 0 {
		unearnedVal = amt.Cents() - deliveredVal
	}
	deliveredAmt := domain.NewMoney(deliveredVal)
	unearnedAmt := domain.NewMoney(unearnedVal)

	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Deferred Revenue)"},
		{ID: "opt_service_revenue", Text: "Service Revenue (100% Immediate)", ErrorTag: engine.TagPrematureRevenueRecognition},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2501)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When selling bundled hardware with future multi-year software updates, which liability account records the deferred portion for undelivered software?",
		CorrectOptionID:   "opt_unearned_revenue",
		Options:           stage1Opts,
		CausalHint:        "When a portion of the bundle price covers future performance obligations (software updates over 2 years), which liability account records that unearned portion?",
		Explanation:       "Unearned Revenue records the consideration allocated to undelivered performance obligations until future software updates are delivered over time.",
		RelevantConceptID: "unearned_revenue_liability",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagPrematureRevenueRecognition},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_asset", Text: "Asset"},
	}, seed, 2502)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Unearned Revenue?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does receiving customer cash in advance for undelivered future software updates create an obligation (liability) to deliver?",
		Explanation:       "Unearned Revenue is a Liability representing the obligation to perform future services or provide software updates.",
		RelevantConceptID: "unearned_revenue_liability",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2503)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Unearned Revenue increase or decrease when the bundled sales cash is collected?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When customer cash is received for undelivered future software updates, does the Unearned Revenue liability increase or decrease?",
		Explanation:       "Unearned Revenue increases because the company assumed a new obligation to deliver updates.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2504)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase a Liability account with a normal credit balance?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have normal credit balances. Do they increase on the left or right side?",
		Explanation:       "Liabilities increase on the right side (Credit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash_and_revenue", Text: "Cash (Asset) and Service Revenue (Earned)"},
		{ID: "opt_cash_only", Text: "Cash only (crediting all to Service Revenue)", ErrorTag: engine.TagPrematureRevenueRecognition},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2505)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which accounts record the cash inflow and the point-of-sale revenue for delivered hardware?",
		CorrectOptionID:   "opt_cash_and_revenue",
		Options:           stage5Opts,
		CausalHint:        "Which revenue account records the earned portion allocated to the delivered hardware device?",
		Explanation:       "Cash is debited for the full proceeds, and Service Revenue is credited for the portion allocated to the delivered hardware.",
		RelevantConceptID: "point_of_sale_revenue",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_rev_unearned", Text: fmt.Sprintf("Debit Cash %s, Credit Service Revenue %s, Credit Unearned Revenue %s", amt.FormatDollars(), deliveredAmt.FormatDollars(), unearnedAmt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev_100", Text: fmt.Sprintf("Debit Cash %s and Credit Service Revenue %s (100%% Immediate)", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPrematureRevenueRecognition},
		{ID: "opt_reversed", Text: fmt.Sprintf("Debit Service Revenue %s, Debit Unearned Revenue %s, Credit Cash %s", deliveredAmt.FormatDollars(), unearnedAmt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2506)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry allocating consideration between delivered hardware and deferred software updates?",
		CorrectOptionID:   "opt_dr_cash_cr_rev_unearned",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the cash proceeds, point-of-sale revenue for delivered hardware, and deferred revenue for software obligations?",
		Explanation:       fmt.Sprintf("Debit Cash %s, Credit Service Revenue %s, and Credit Unearned Revenue %s. Consideration is allocated between delivered and undelivered performance obligations.", amt.FormatDollars(), deliveredAmt.FormatDollars(), unearnedAmt.FormatDollars()),
		RelevantConceptID: "multi_element_revenue_recognition",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this multi-element revenue allocation affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_up_liab_up_eq_up",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up_eq_up", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Rev); Equity increases by %s (+Service Rev).", amt.FormatDollars(), unearnedAmt.FormatDollars(), deliveredAmt.FormatDollars())},
			{ID: "opt_assets_up_eq_up_100", Text: fmt.Sprintf("Assets increase by %s (+Cash) and Equity increases by %s (+Revenue 100%%); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagPrematureRevenueRecognition},
			{ID: "opt_no_liab", Text: fmt.Sprintf("Assets increase by %s; Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2507),
		CausalHint:        "Cash received increases Assets. How does splitting the transaction between earned revenue and unearned liability affect Liabilities and Equity?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s), Liabilities increase (+Unearned Revenue %s), and Equity increases (+Service Revenue %s). Equation balances: +%s = +%s + %s.", amt.FormatDollars(), unearnedAmt.FormatDollars(), deliveredAmt.FormatDollars(), amt.FormatDollars(), unearnedAmt.FormatDollars(), deliveredAmt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

func buildGoodwillImpairmentStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_loss_on_impairment", Text: "Loss on Impairment"},
		{ID: "opt_rent_expense", Text: "Rent Expense", ErrorTag: engine.TagWrongAccount},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 2601)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When testing goodwill for impairment at year-end and reporting unit carrying value exceeds fair value, which expense/loss account records the asset write-down?",
		CorrectOptionID:   "opt_loss_on_impairment",
		Options:           stage1Opts,
		CausalHint:        "When a reporting unit's fair value drops below its carrying value, which loss account absorbs this non-cash asset impairment write-down?",
		Explanation:       "Loss on Impairment is recognized in operating earnings to record the drop in goodwill value below carrying value.",
		RelevantConceptID: "goodwill_impairment",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_expense", Text: "Expense (Loss)"},
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2602)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Loss on Impairment?",
		CorrectOptionID:   "opt_expense",
		Options:           stage2Opts,
		CausalHint:        "Does an impairment loss reduce current period earnings and equity, or represent a future resource or debt?",
		Explanation:       "Loss on Impairment is an Expense/Loss account that reduces Net Income and Stockholders' Equity on the income statement.",
		RelevantConceptID: "goodwill_impairment",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2603)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Loss on Impairment increase or decrease when recognizing this asset write-down?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "Does recognizing a new impairment loss increase or decrease the accumulated loss recorded for the period?",
		Explanation:       "Loss on Impairment increases as the asset write-down is recognized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2604)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Expense/Loss account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Expenses and losses reduce equity and have normal debit balances. Do they increase on the left or right side?",
		Explanation:       "Expenses and losses have normal debit balances, increasing on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_goodwill", Text: "Goodwill (Asset)"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount},
	}, seed, 2605)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which intangible asset account is credited to write down the unrecoverable carrying value?",
		CorrectOptionID:   "opt_goodwill",
		Options:           stage5Opts,
		CausalHint:        "Which recorded intangible asset became impaired and must be written down on the balance sheet?",
		Explanation:       "Goodwill is credited to reduce its balance sheet carrying value down to implied recoverable value.",
		RelevantConceptID: "goodwill_impairment",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_loss_cr_gw", Text: fmt.Sprintf("Debit Loss on Impairment %s / Credit Goodwill %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_gw_cr_loss", Text: fmt.Sprintf("Debit Goodwill %s / Credit Loss on Impairment %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_loss_cr_cash", Text: fmt.Sprintf("Debit Loss on Impairment %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 2606)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the goodwill impairment loss?",
		CorrectOptionID:   "opt_dr_loss_cr_gw",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the impairment loss in earnings and writes down the goodwill asset?",
		Explanation:       fmt.Sprintf("Debit Loss on Impairment %s and Credit Goodwill %s. Goodwill impairment is a non-cash write-down that directly reduces asset carrying value.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "goodwill_impairment",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does recording goodwill impairment affect the balance sheet equation?",
		CorrectOptionID: "opt_assets_down_eq_down",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s (-Goodwill); Equity decreases by %s (-Loss on Impairment); Liabilities unchanged.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_down_liab_down", Text: fmt.Sprintf("Assets decrease by %s; Liabilities decrease by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_no_effect", Text: fmt.Sprintf("Total Assets and Equity remain completely unchanged."), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2607),
		CausalHint:        "Goodwill decreased and an impairment loss reduced net income. How does this write-down affect total Assets and Equity?",
		Explanation:       fmt.Sprintf("Assets decrease (-Goodwill %s) and Equity decreases (-Loss on Impairment %s); Liabilities are unchanged. Under US GAAP (ASC 350), goodwill impairment cannot be reversed even if fair value subsequently recovers.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "accounting_equation_invariants",
	}
}

func buildDigitalAssetAccountingStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_digital_assets", Text: "Digital Assets (Cryptocurrency)"},
		{ID: "opt_cash_equivalents", Text: "Cash Equivalents", ErrorTag: engine.TagCryptoRecordedAsCashEquivalent},
		{ID: "opt_rent_expense", Text: "Operating Expense", ErrorTag: engine.TagCryptoRecordedAsCashEquivalent},
	}, seed, 2701)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When a corporation purchases cryptocurrency (such as Tesla purchasing Bitcoin) for cash, which account records this indefinite-lived intangible asset?",
		CorrectOptionID:   "opt_digital_assets",
		Options:           stage1Opts,
		CausalHint:        "Under US GAAP (ASC 350), are cryptocurrencies treated as legal cash equivalents, or capitalized as indefinite-lived intangible assets?",
		Explanation:       "Digital Assets records cryptocurrency holdings, classified as indefinite-lived intangible assets under US GAAP.",
		RelevantConceptID: "digital_asset_accounting",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset (Intangible Asset)"},
		{ID: "opt_expense", Text: "Expense", ErrorTag: engine.TagCryptoRecordedAsCashEquivalent},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2702)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Digital Assets?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does purchasing cryptocurrency represent an immediate expense, or a productive intangible resource owned by the company?",
		Explanation:       "Digital Assets is an Intangible Asset with an indefinite useful life.",
		RelevantConceptID: "digital_asset_accounting",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2703)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Digital Assets increase or decrease when the corporation acquires cryptocurrency?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When the company buys and takes custody of cryptocurrency, does its digital asset balance increase or decrease?",
		Explanation:       "Digital Assets increases as new digital holdings are acquired.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2704)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have normal debit balances. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagCryptoRecordedAsCashEquivalent},
		{ID: "opt_accounts_payable", Text: "Accounts Payable"},
	}, seed, 2705)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liquid asset account is credited for the funds paid to purchase the cryptocurrency?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid payment medium was transferred from bank reserves to purchase the digital tokens?",
		Explanation:       "Cash is credited because cash was disbursed to execute the purchase.",
		RelevantConceptID: "cash_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_digital_cr_cash", Text: fmt.Sprintf("Debit Digital Assets %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_digital", Text: fmt.Sprintf("Debit Cash %s / Credit Digital Assets %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_exp_cr_cash", Text: fmt.Sprintf("Debit Operating Expense %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCryptoRecordedAsCashEquivalent},
	}, seed, 2706)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry to record the cash purchase of cryptocurrency?",
		CorrectOptionID:   "opt_dr_digital_cr_cash",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the acquisition of the intangible digital asset and the cash outflow?",
		Explanation:       fmt.Sprintf("Debit Digital Assets %s and Credit Cash %s. The purchase is an asset exchange, not an immediate operating expense.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "digital_asset_accounting",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this acquisition affect the accounting equation, and how are subsequent market fluctuations treated under US GAAP?",
		CorrectOptionID: "opt_asset_swap_crypto",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_asset_swap_crypto", Text: fmt.Sprintf("Asset swap: Digital Assets increases by %s (+A), Cash decreases by %s (-A); Liabilities and Equity unchanged. Subsequent price declines require impairment losses, but price increases cannot be recognized as gains until sale.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s and Equity increases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_assets_down_eq_down", Text: fmt.Sprintf("Assets decrease by %s and Equity decreases by %s.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2707),
		CausalHint:        "One asset increased while another asset decreased by the exact same amount. How does this asset exchange affect the accounting equation?",
		Explanation:       fmt.Sprintf("Asset swap: Digital Assets increases (+A %s) and Cash decreases (-A %s); Total Assets, Liabilities, and Equity remain unchanged. Under US GAAP (ASC 350), digital assets are tested for impairment whenever market prices drop below carrying value, but price appreciation cannot be recognized as a gain until realized upon sale.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "digital_asset_accounting",
	}
}

func buildReportingPolicyChoiceStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue (Deferred Revenue)"},
		{ID: "opt_service_revenue", Text: "Service Revenue (Aggressive Immediate Recognition)", ErrorTag: engine.TagAggressiveRevenueAcceleration},
		{ID: "opt_accounts_payable", Text: "Accounts Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2801)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "When receiving customer cash for service contracts with pending performance obligations, which liability account is credited under a conservative revenue recognition policy?",
		CorrectOptionID:   "opt_unearned_revenue",
		Options:           stage1Opts,
		CausalHint:        "Under a conservative financial reporting policy, when cash is collected before all performance obligations are complete, which liability account defers the revenue?",
		Explanation:       "Unearned Revenue records customer consideration received before performance obligations are fully satisfied, adhering to conservative accrual principles.",
		RelevantConceptID: "reporting_strategy",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_revenue", Text: "Revenue", ErrorTag: engine.TagAggressiveRevenueAcceleration},
		{ID: "opt_equity", Text: "Equity"},
		{ID: "opt_asset", Text: "Asset"},
	}, seed, 2802)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Unearned Revenue?",
		CorrectOptionID:   "opt_liability",
		Options:           stage2Opts,
		CausalHint:        "Does receiving customer money in advance create an obligation to perform future services?",
		Explanation:       "Unearned Revenue is a Liability representing the obligation to deliver future services or goods.",
		RelevantConceptID: "reporting_strategy",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2803)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does Unearned Revenue increase or decrease when the advance payment is received?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When a new advance customer receipt arrives with pending work, does the unearned obligation increase or decrease?",
		Explanation:       "Unearned Revenue increases as the new customer advance obligation is recognized.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_credit", Text: "Credit (Right side)"},
		{ID: "opt_debit", Text: "Debit (Left side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2804)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase a Liability account?",
		CorrectOptionID:   "opt_credit",
		Options:           stage4Opts,
		CausalHint:        "Liabilities have normal credit balances. Which side increases a liability?",
		Explanation:       "Liabilities increase on the right side (Credit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_cash", Text: "Cash (Asset)"},
		{ID: "opt_service_revenue", Text: "Service Revenue", ErrorTag: engine.TagAggressiveRevenueAcceleration},
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable"},
	}, seed, 2805)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which liquid asset account is debited for the cash received from the customer?",
		CorrectOptionID:   "opt_cash",
		Options:           stage5Opts,
		CausalHint:        "What liquid payment medium was received from the customer today?",
		Explanation:       "Cash is debited because cash proceeds were received.",
		RelevantConceptID: "cash_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_cash_cr_unearned", Text: fmt.Sprintf("Debit Cash %s / Credit Unearned Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s (Aggressive Acceleration)", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagAggressiveRevenueAcceleration},
		{ID: "opt_dr_unearned_cr_cash", Text: fmt.Sprintf("Debit Unearned Revenue %s / Credit Cash %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
	}, seed, 2806)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry recording the receipt under a conservative reporting policy?",
		CorrectOptionID:   "opt_dr_cash_cr_unearned",
		Options:           stage6Opts,
		CausalHint:        "Which entry records the cash receipt and defers revenue recognition until performance obligations are satisfied?",
		Explanation:       fmt.Sprintf("Debit Cash %s and Credit Unearned Revenue %s. Deferring revenue until earned preserves earnings quality and prevents premature revenue recognition.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "reporting_strategy",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this conservative policy choice affect the balance sheet equation compared to an aggressive policy?",
		CorrectOptionID: "opt_assets_up_liab_up_conservative",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_liab_up_conservative", Text: fmt.Sprintf("Assets increase by %s (+Cash); Liabilities increase by %s (+Unearned Rev); Equity unchanged. Conservative policy defers earnings, whereas aggressive policy would prematurely inflate Equity.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_assets_up_eq_up", Text: fmt.Sprintf("Assets increase by %s and Equity increases by %s immediately.", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagAggressiveRevenueAcceleration},
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset swap: Cash increases and receivables decrease by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
		}, seed, 2807),
		CausalHint:        "Cash increased and an unearned liability was recognized. How does this conservative choice affect Liabilities vs Equity?",
		Explanation:       fmt.Sprintf("Assets increase (+Cash %s) and Liabilities increase (+Unearned Revenue %s); Equity is unchanged. A conservative reporting strategy defers revenue until delivery, buffering future earnings, whereas aggressive acceleration immediately inflates current-period net income.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "reporting_strategy",
	}
}

func buildFinancialForecastingStages(inst *domain.QuestionInstance, amt domain.Money, seed int64, proc engine.ProcessedEvent) {
	stage1Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_accounts_receivable", Text: "Accounts Receivable"},
		{ID: "opt_cash", Text: "Cash", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount},
	}, seed, 2901)

	inst.StageAnswers[domain.StageIdentifyAccount] = domain.StageAnswer{
		Stage:             domain.StageIdentifyAccount,
		Prompt:            "In pro-forma financial statement forecasting, which asset account models projected uncollected credit sales owed by customers?",
		CorrectOptionID:   "opt_accounts_receivable",
		Options:           stage1Opts,
		CausalHint:        "When forecasting credit sales that customers will pay on invoice terms, which current asset account tracks the uncollected customer claims?",
		Explanation:       "Accounts Receivable is the asset that tracks projected uncollected customer sales in pro-forma financial statements.",
		RelevantConceptID: "financial_forecasting",
	}

	stage2Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_asset", Text: "Asset"},
		{ID: "opt_revenue", Text: "Revenue"},
		{ID: "opt_liability", Text: "Liability"},
		{ID: "opt_equity", Text: "Equity"},
	}, seed, 2902)

	inst.StageAnswers[domain.StageAccountCategory] = domain.StageAnswer{
		Stage:             domain.StageAccountCategory,
		Prompt:            "What category of account is Accounts Receivable?",
		CorrectOptionID:   "opt_asset",
		Options:           stage2Opts,
		CausalHint:        "Does a claim against customers for future payment represent an asset, liability, or equity?",
		Explanation:       "Accounts Receivable is a Current Asset representing customer obligations owed to the business.",
		RelevantConceptID: "financial_forecasting",
	}

	stage3Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_increase", Text: "Increase"},
		{ID: "opt_decrease", Text: "Decrease"},
	}, seed, 2903)

	inst.StageAnswers[domain.StageDirection] = domain.StageAnswer{
		Stage:             domain.StageDirection,
		Prompt:            "Does projected Accounts Receivable increase or decrease when modeling future credit sales?",
		CorrectOptionID:   "opt_increase",
		Options:           stage3Opts,
		CausalHint:        "When projecting new quarterly credit sales, does the uncollected customer claim balance increase or decrease?",
		Explanation:       "Accounts Receivable increases as new credit sales are projected.",
		RelevantConceptID: "debit_credit_translation",
	}

	stage4Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_debit", Text: "Debit (Left side)"},
		{ID: "opt_credit", Text: "Credit (Right side)", ErrorTag: engine.TagReversedSides},
	}, seed, 2904)

	inst.StageAnswers[domain.StageDebitCredit] = domain.StageAnswer{
		Stage:             domain.StageDebitCredit,
		Prompt:            "Which side is recorded to increase an Asset account?",
		CorrectOptionID:   "opt_debit",
		Options:           stage4Opts,
		CausalHint:        "Assets have normal debit balances. Which side increases an asset?",
		Explanation:       "Assets increase on the left side (Debit).",
		RelevantConceptID: "debit_credit_translation",
	}

	stage5Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_service_revenue", Text: "Service Revenue (Revenue)"},
		{ID: "opt_cash", Text: "Cash (Asset)", ErrorTag: engine.TagCashRecordedWhenUncollected},
		{ID: "opt_unearned_revenue", Text: "Unearned Revenue"},
	}, seed, 2905)

	inst.StageAnswers[domain.StageCounterAccount] = domain.StageAnswer{
		Stage:             domain.StageCounterAccount,
		Prompt:            "Which revenue account is credited for the projected sales in the pro-forma income statement?",
		CorrectOptionID:   "opt_service_revenue",
		Options:           stage5Opts,
		CausalHint:        "Which income statement account captures projected earned sales from customer deliveries?",
		Explanation:       "Service Revenue is credited for the projected earnings generated in the forecast period.",
		RelevantConceptID: "service_revenue_classification",
	}

	stage6Opts := shuffleOptions([]domain.AnswerOption{
		{ID: "opt_dr_ar_cr_rev", Text: fmt.Sprintf("Debit Accounts Receivable %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars())},
		{ID: "opt_dr_rev_cr_ar", Text: fmt.Sprintf("Debit Service Revenue %s / Credit Accounts Receivable %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagReversedSides},
		{ID: "opt_dr_cash_cr_rev", Text: fmt.Sprintf("Debit Cash %s / Credit Service Revenue %s", amt.FormatDollars(), amt.FormatDollars()), ErrorTag: engine.TagCashRecordedWhenUncollected},
	}, seed, 2906)

	inst.StageAnswers[domain.StageBalancedEntry] = domain.StageAnswer{
		Stage:             domain.StageBalancedEntry,
		Prompt:            "What is the complete balanced journal entry modeling projected credit sales in the pro-forma forecast?",
		CorrectOptionID:   "opt_dr_ar_cr_rev",
		Options:           stage6Opts,
		CausalHint:        "Which entry projects the uncollected receivable claim and the earned sales in the pro-forma model?",
		Explanation:       fmt.Sprintf("Debit Accounts Receivable %s and Credit Service Revenue %s. Projected credit sales expand the receivable asset and increase projected operating earnings.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "financial_forecasting",
	}

	inst.StageAnswers[domain.StageEquationEffect] = domain.StageAnswer{
		Stage:           domain.StageEquationEffect,
		Prompt:          "How does this entry affect the projected balance sheet equation, and why do disaggregated drivers provide superior forecasting accuracy?",
		CorrectOptionID: "opt_assets_up_eq_up_forecasting",
		Options: shuffleOptions([]domain.AnswerOption{
			{ID: "opt_assets_up_eq_up_forecasting", Text: fmt.Sprintf("Assets increase by %s (+A/R); Equity increases by %s (+Service Rev); Liabilities unchanged. Disaggregated unit drivers (volume × price) isolate operational bottlenecks, unlike crude top-line percentages.", amt.FormatDollars(), amt.FormatDollars())},
			{ID: "opt_asset_swap", Text: fmt.Sprintf("Asset swap: Cash increases and receivables decrease by %s.", amt.FormatDollars()), ErrorTag: engine.TagEquationEffectMissed},
			{ID: "opt_aggregated_growth", Text: fmt.Sprintf("Revenue increases by %s based on simple top-line percentage growth without operational drivers.", amt.FormatDollars()), ErrorTag: engine.TagAggregatedForecastAssumption},
		}, seed, 2907),
		CausalHint:        "Receivables increased and revenue increased equity. Why is modeling bottom-up unit volume times price more reliable than simple percentage growth?",
		Explanation:       fmt.Sprintf("Assets increase (+Accounts Receivable %s) and Equity increases (+Service Revenue %s); Liabilities are unchanged. In financial statement forecasting, bottom-up disaggregated unit drivers (volume × price, store count) provide structural predictive accuracy and enable working capital tie-outs, whereas crude top-line percentages obscure operational constraints.", amt.FormatDollars(), amt.FormatDollars()),
		RelevantConceptID: "financial_forecasting",
	}
}



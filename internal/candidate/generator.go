package candidate

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

// GenerateRequest specifies parameters for candidate generation.
type GenerateRequest struct {
	FamilyID      string // Optional: if empty, chosen from weak concepts or random supported family
	TargetConcept string // Optional: target weak concept
	Seed          int64  // Deterministic PRNG seed (0 for current time)
}

// CandidateGenerator defines the interface for creating structured question proposals.
type CandidateGenerator interface {
	Name() string
	Generate(ctx context.Context, req GenerateRequest) (*CandidateQuestion, error)
}

// ConceptToFamilyMap maps accounting concepts to the supported families that test them.
var ConceptToFamilyMap = map[string][]string{
	"cash_vs_revenue":                    {bank.FamilyCustomerAdvance, bank.FamilyCashService, bank.FamilyCollectReceivable, bank.FamilyEarnAdvance},
	"unearned_revenue_classification":    {bank.FamilyCustomerAdvance, bank.FamilyEarnAdvance},
	"accounts_receivable_classification": {bank.FamilyServiceOnCredit, bank.FamilyCollectReceivable},
	"cash_classification":                {bank.FamilyCashService, bank.FamilyBorrowCash, bank.FamilyDividendCash},
	"service_revenue_classification":     {bank.FamilyCashService, bank.FamilyServiceOnCredit},
	"prepaid_purchase":                   {bank.FamilyPrepaidPurchase},
	"prepaid_consumption":                {bank.FamilyPrepaidConsumption},
	"prepaid_expenses":                   {bank.FamilyPrepaidPurchase, bank.FamilyPrepaidConsumption},
	"capital_vs_expense":                 {bank.FamilyEquipmentPurchaseCash, bank.FamilyCashExpense, bank.FamilyCashRent},
	"debt_vs_equity":                     {bank.FamilyBorrowCash, bank.FamilyIssueShares},
	"note_payable_principal":             {bank.FamilyBorrowCash, bank.FamilyRepayNotePrincipal},
	"dividends_vs_expense":               {bank.FamilyDividendCash, bank.FamilyCashExpense},
	"debit_credit_translation":           {bank.FamilyCashService, bank.FamilyCustomerAdvance, bank.FamilyServiceOnCredit},
}

// SelectWeakestConcept examines learner mastery projections and returns the concept needing the most practice,
// along with the best matching supported family.
func SelectWeakestConcept(projections map[string]*mastery.ConceptStats, now time.Time) (string, string) {
	if len(projections) == 0 {
		return "cash_vs_revenue", bank.FamilyCustomerAdvance
	}

	type conceptNeed struct {
		id       string
		priority float64
	}
	var needs []conceptNeed
	for cid, stats := range projections {
		needs = append(needs, conceptNeed{
			id:       cid,
			priority: stats.PracticePriority(now),
		})
	}

	sort.Slice(needs, func(i, j int) bool {
		return needs[i].priority > needs[j].priority
	})

	for _, n := range needs {
		families, ok := ConceptToFamilyMap[n.id]
		if ok && len(families) > 0 {
			return n.id, families[0]
		}
	}

	return "cash_vs_revenue", bank.FamilyCustomerAdvance
}

// -----------------------------------------------------------------------------
// Offline Candidate Generator
// -----------------------------------------------------------------------------

// OfflineCandidateGenerator deterministically creates creative candidate questions
// without any network or external provider dependencies.
type OfflineCandidateGenerator struct {
	eng     *engine.Engine
	catalog *domain.AccountCatalog
}

func NewOfflineCandidateGenerator(eng *engine.Engine, catalog *domain.AccountCatalog) *OfflineCandidateGenerator {
	return &OfflineCandidateGenerator{
		eng:     eng,
		catalog: catalog,
	}
}

func (o *OfflineCandidateGenerator) Name() string {
	return "offline_generator"
}

type templateVariation struct {
	scenarioTemplate string
	parameters       map[string][]int64
	concepts         []string
}

var offlineVariations = map[string][]templateVariation{
	bank.FamilyCashService: {
		{
			scenarioTemplate: "A client hires an architectural drafting firm and pays ${amount_dollars} in cash upon receipt of the final blueprints today. No prior bill was issued.",
			parameters:       map[string][]int64{"amount_minor_units": {45000, 90000, 150000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "service_revenue_classification", "debit_credit_translation"},
		},
		{
			scenarioTemplate: "A cybersecurity consultancy performs an on-site server penetration test today and receives ${amount_dollars} cash upon completion. No invoice was previously sent.",
			parameters:       map[string][]int64{"amount_minor_units": {35000, 75000, 120000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "service_revenue_classification", "debit_credit_translation"},
		},
		{
			scenarioTemplate: "A specialized mobile veterinary clinic visits an equestrian ranch today and receives ${amount_dollars} cash immediately upon completing examinations. No credit was extended.",
			parameters:       map[string][]int64{"amount_minor_units": {25000, 50000, 85000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "service_revenue_classification", "debit_credit_translation"},
		},
	},
	bank.FamilyCustomerAdvance: {
		{
			scenarioTemplate: "A corporate client pays ${amount_dollars} cash today as an upfront retainer for cloud migration services scheduled to commence next quarter. No migration work has started.",
			parameters:       map[string][]int64{"amount_minor_units": {60000, 120000, 240000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "unearned_revenue_classification", "liability_increase_credit"},
		},
		{
			scenarioTemplate: "A catering company receives ${amount_dollars} cash today for a wedding banquet that will take place in two months. No food or services have been delivered.",
			parameters:       map[string][]int64{"amount_minor_units": {50000, 100000, 180000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "unearned_revenue_classification", "liability_increase_credit"},
		},
		{
			scenarioTemplate: "A commercial photography studio receives ${amount_dollars} cash today from a retail fashion brand to reserve a studio shoot scheduled for next month. No photos have been taken.",
			parameters:       map[string][]int64{"amount_minor_units": {30000, 60000, 120000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "unearned_revenue_classification", "liability_increase_credit"},
		},
	},
	bank.FamilyServiceOnCredit: {
		{
			scenarioTemplate: "A digital marketing agency delivers a comprehensive SEO campaign today and issues an invoice for ${amount_dollars} payable in 30 days. No cash is collected today.",
			parameters:       map[string][]int64{"amount_minor_units": {40000, 80000, 160000}},
			concepts:         []string{"cash_vs_revenue", "accounts_receivable_classification", "service_revenue_classification", "asset_increase_debit"},
		},
		{
			scenarioTemplate: "A legal firm provides 15 hours of patent advisory services today and bills the corporate client ${amount_dollars} on account, payable net 30. No cash was received.",
			parameters:       map[string][]int64{"amount_minor_units": {75000, 150000, 300000}},
			concepts:         []string{"cash_vs_revenue", "accounts_receivable_classification", "service_revenue_classification", "asset_increase_debit"},
		},
	},
	bank.FamilyCollectReceivable: {
		{
			scenarioTemplate: "A design agency collects ${amount_dollars} cash today from a client in full payment of an invoice billed and recognized as revenue last month. No new service was performed today.",
			parameters:       map[string][]int64{"amount_minor_units": {40000, 80000, 160000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "accounts_receivable_classification", "asset_swap"},
		},
		{
			scenarioTemplate: "A freight carrier receives ${amount_dollars} cash today from a shipper settling an account receivable recorded 30 days ago. The revenue was recognized at the time of delivery last month.",
			parameters:       map[string][]int64{"amount_minor_units": {55000, 110000, 220000}},
			concepts:         []string{"cash_vs_revenue", "cash_classification", "accounts_receivable_classification", "asset_swap"},
		},
	},
	bank.FamilyEarnAdvance: {
		{
			scenarioTemplate: "A web development agency finishes and launches an e-commerce platform today for a client who had paid ${amount_dollars} in advance last month. No cash is received today.",
			parameters:       map[string][]int64{"amount_minor_units": {60000, 120000, 240000}},
			concepts:         []string{"cash_vs_revenue", "unearned_revenue_classification", "service_revenue_classification", "liability_decrease_debit"},
		},
		{
			scenarioTemplate: "A private flight academy delivers 25 hours of ground instruction today to a trainee pilot who previously paid ${amount_dollars} cash in advance two weeks ago. Zero cash changes hands today.",
			parameters:       map[string][]int64{"amount_minor_units": {35000, 70000, 140000}},
			concepts:         []string{"cash_vs_revenue", "unearned_revenue_classification", "service_revenue_classification", "liability_decrease_debit"},
		},
	},
	bank.FamilyCashRent: {
		{
			scenarioTemplate: "A medical clinic pays ${amount_dollars} cash today for its clinic office lease covering the current calendar month. No prior liability existed.",
			parameters:       map[string][]int64{"amount_minor_units": {80000, 160000, 320000}},
			concepts:         []string{"cash_classification", "expense_classification", "equity_reduction_debit", "asset_decrease_credit"},
		},
	},
	bank.FamilyBorrowCash: {
		{
			scenarioTemplate: "The company signs a 3-year commercial promissory note and receives ${amount_dollars} cash today from First Community Bank. This is borrowed debt, not earned revenue.",
			parameters:       map[string][]int64{"amount_minor_units": {1000000, 2500000, 5000000}},
			concepts:         []string{"debt_vs_equity", "cash_classification", "liability_increase_credit", "balance_sheet_expansion"},
		},
	},
	bank.FamilyIssueShares: {
		{
			scenarioTemplate: "The corporation issues 2,000 shares of common stock to outside venture investors in exchange for ${amount_dollars} cash today. This represents equity financing, not revenue.",
			parameters:       map[string][]int64{"amount_minor_units": {2000000, 5000000, 10000000}},
			concepts:         []string{"debt_vs_equity", "cash_classification", "equity_increase_credit", "balance_sheet_expansion"},
		},
	},
	bank.FamilyCashExpense: {
		{
			scenarioTemplate: "A graphic design studio pays ${amount_dollars} cash today to an independent IT contractor for on-site workstation diagnostics completed this morning.",
			parameters:       map[string][]int64{"amount_minor_units": {15000, 30000, 60000}},
			concepts:         []string{"expense_classification", "cash_classification", "equity_reduction_debit", "asset_decrease_credit"},
		},
	},
	bank.FamilyPrepaidPurchase: {
		{
			scenarioTemplate: "A trucking enterprise pays ${amount_dollars} cash today for an annual commercial vehicle liability policy covering the next 12 future months. Zero insurance has yet expired.",
			parameters:       map[string][]int64{"amount_minor_units": {120000, 240000, 480000}},
			concepts:         []string{"prepaid_expenses", "cash_classification", "asset_swap", "accrual_timing"},
		},
	},
	bank.FamilyPrepaidConsumption: {
		{
			scenarioTemplate: "At month-end, the company records the expiration of one month of insurance coverage from an annual policy previously prepaid for ${amount_dollars}. No cash changes hands today.",
			parameters:       map[string][]int64{"amount_minor_units": {10000, 20000, 40000}},
			concepts:         []string{"prepaid_expenses", "expense_classification", "asset_decrease_credit", "accrual_matching"},
		},
	},
	bank.FamilyEquipmentPurchaseCash: {
		{
			scenarioTemplate: "A regional roastery purchases an industrial coffee bean roaster today for ${amount_dollars} in cash. The equipment has an estimated 8-year useful life and is capitalized as a noncurrent asset.",
			parameters:       map[string][]int64{"amount_minor_units": {850000, 1500000, 3000000}},
			concepts:         []string{"capital_vs_expense", "cash_classification", "asset_swap", "balance_sheet_invariance"},
		},
	},
	bank.FamilyRepayNotePrincipal: {
		{
			scenarioTemplate: "The company pays ${amount_dollars} cash today to a commercial lender to retire the principal of an outstanding promissory note that matured today. (Ignore interest).",
			parameters:       map[string][]int64{"amount_minor_units": {500000, 1000000, 2000000}},
			concepts:         []string{"note_payable_principal", "cash_classification", "liability_decrease_debit", "balance_sheet_contraction"},
		},
	},
	bank.FamilyDividendCash: {
		{
			scenarioTemplate: "The board of directors declares and pays ${amount_dollars} in cash dividends to shareholders today. This is a direct distribution of retained earnings, not an operating expense.",
			parameters:       map[string][]int64{"amount_minor_units": {100000, 250000, 500000}},
			concepts:         []string{"dividends_vs_expense", "cash_classification", "equity_reduction_debit", "asset_decrease_credit"},
		},
	},
}

func (o *OfflineCandidateGenerator) Generate(ctx context.Context, req GenerateRequest) (*CandidateQuestion, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	familyID := req.FamilyID
	if familyID == "" && req.TargetConcept != "" {
		if families, ok := ConceptToFamilyMap[req.TargetConcept]; ok && len(families) > 0 {
			familyID = families[0]
		}
	}
	if familyID == "" {
		// Pick customer_advance by default as the core contrast family
		familyID = bank.FamilyCustomerAdvance
	}

	variations, ok := offlineVariations[familyID]
	if !ok || len(variations) == 0 {
		return nil, fmt.Errorf("no creative variations defined for supported family %q", familyID)
	}

	seed := req.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))
	chosenIdx := rng.Intn(len(variations))
	chosen := variations[chosenIdx]

	prov := Provenance{
		Source:        "offline_generator",
		GeneratedAt:   time.Now().UTC(),
		TargetFamily:  familyID,
		TargetConcept: req.TargetConcept,
	}

	proposal := RawProposal{
		ID:               generateCandidateID(familyID, prov.GeneratedAt),
		FamilyID:         familyID,
		ScenarioTemplate: chosen.scenarioTemplate,
		Parameters:       chosen.parameters,
		Concepts:         chosen.concepts,
	}

	payload, err := json.Marshal(proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proposal: %w", err)
	}

	return ParseProposal(string(payload), o.eng, o.catalog, prov)
}

// -----------------------------------------------------------------------------
// Provider Candidate Generator (LLM / External Tutor Provider)
// -----------------------------------------------------------------------------

// ProviderCandidateGenerator prompts an external AI tutor provider or simulated provider
// to produce a structured JSON question proposal, strictly parses the output, and derives
// answers locally using the engine.
type ProviderCandidateGenerator struct {
	tut     tutor.Tutor
	eng     *engine.Engine
	catalog *domain.AccountCatalog
}

func NewProviderCandidateGenerator(tut tutor.Tutor, eng *engine.Engine, catalog *domain.AccountCatalog) *ProviderCandidateGenerator {
	return &ProviderCandidateGenerator{
		tut:     tut,
		eng:     eng,
		catalog: catalog,
	}
}

func (p *ProviderCandidateGenerator) Name() string {
	if p.tut != nil {
		return p.tut.Name()
	}
	return "provider_generator"
}

func (p *ProviderCandidateGenerator) Generate(ctx context.Context, req GenerateRequest) (*CandidateQuestion, error) {
	if p.tut == nil {
		return nil, fmt.Errorf("tutor provider is not configured")
	}

	familyID := req.FamilyID
	if familyID == "" && req.TargetConcept != "" {
		if families, ok := ConceptToFamilyMap[req.TargetConcept]; ok && len(families) > 0 {
			familyID = families[0]
		}
	}
	if familyID == "" {
		familyID = bank.FamilyCustomerAdvance
	}

	prompt := FormatCandidateGenerationPrompt(familyID, req.TargetConcept)

	tutorReq := tutor.Request{
		CandidateGeneration: true,
		ProblemPrompt:       prompt,
		FamilyID:            familyID,
		ConceptID:           req.TargetConcept,
		MaxTokens:           1500,
	}

	resp, err := p.tut.Explain(ctx, tutorReq)
	if err == nil && resp.Fallback {
		err = fmt.Errorf("connected provider unavailable: %s", resp.FallbackReason)
	}
	if err != nil {
		prov := Provenance{
			Source:        p.tut.Name(),
			GeneratedAt:   time.Now().UTC(),
			TargetFamily:  familyID,
			TargetConcept: req.TargetConcept,
			PromptText:    prompt,
		}
		cand := &CandidateQuestion{
			ID:               generateCandidateID(familyID, prov.GeneratedAt),
			FamilyID:         familyID,
			Status:           StatusCandidateRejected,
			ValidationStatus: ValidationFailed,
			RejectionReason:  fmt.Sprintf("tutor provider call failed: %v", err),
			Provenance:       prov,
		}
		return cand, err
	}

	prov := Provenance{
		Source:        p.tut.Name(),
		Model:         resp.Provider,
		GeneratedAt:   resp.GeneratedAt,
		TargetFamily:  familyID,
		TargetConcept: req.TargetConcept,
		PromptText:    prompt,
	}

	// Strictly parse untrusted LLM output and derive answers locally via engine
	cand, parseErr := ParseProposal(resp.Text, p.eng, p.catalog, prov)
	if cand != nil {
		// Provider-selected IDs must never overwrite an existing local candidate.
		cand.ID = generateCandidateID(familyID, time.Now().UTC())
	}
	if parseErr == nil && cand.FamilyID != familyID {
		cand.Status = StatusCandidateRejected
		cand.ValidationStatus = ValidationFailed
		cand.RejectionReason = "provider returned a different transaction family"
		parseErr = fmt.Errorf("%s", cand.RejectionReason)
	}
	if parseErr != nil {
		// Even if parse fails, cand holds the failed validation state and reason
		return cand, parseErr
	}

	return cand, nil
}

// FormatCandidateGenerationPrompt crafts a prompt instructing an LLM to produce valid candidate JSON.
func FormatCandidateGenerationPrompt(familyID string, targetConcept string) string {
	var sb strings.Builder
	sb.WriteString("You are a financial accounting curriculum generator proposing a new practice problem candidate.\n")
	sb.WriteString(fmt.Sprintf("Target Family: %s\n", familyID))
	if targetConcept != "" {
		sb.WriteString(fmt.Sprintf("Target Concept: %s\n", targetConcept))
	}
	sb.WriteString("\nGenerate a single new introductory accounting practice question candidate in strict JSON format.\n")
	sb.WriteString("INVARIANTS:\n")
	sb.WriteString("1. Respond ONLY with a valid JSON object matching the schema below. Do NOT include conversational prose.\n")
	sb.WriteString("2. The scenario must be clear and unambiguous, with economic timing explicitly stated (e.g. today vs. next month vs. last month).\n")
	sb.WriteString("3. Must use the exact placeholder ${amount_dollars} in the scenario_template.\n")
	sb.WriteString("4. Parameters must include 'amount_minor_units' as a list of positive integers (e.g. 5000 = $50.00).\n")
	sb.WriteString("5. JSON SCHEMA:\n")
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"family_id\": %q,\n", familyID))
	if examples := offlineVariations[familyID]; len(examples) > 0 {
		sb.WriteString(fmt.Sprintf("  \"scenario_template\": %q,\n", examples[0].scenarioTemplate))
	}
	sb.WriteString("  \"parameters\": {\n")
	sb.WriteString("    \"amount_minor_units\": [5000, 10000, 20000]\n")
	sb.WriteString("  },\n")
	sb.WriteString("  \"concepts\": [\"cash_vs_revenue\", \"unearned_revenue_classification\"]\n")
	sb.WriteString("}\n")
	// Examples use the requested family, avoiding misleading cross-family wording.
	for i, example := range offlineVariations[familyID] {
		if i == 2 {
			break
		}
		payload, _ := json.Marshal(RawProposal{FamilyID: familyID, ScenarioTemplate: example.scenarioTemplate, Parameters: example.parameters, Concepts: example.concepts})
		sb.WriteString("\nStyle example (write a fresh variation): " + string(payload) + "\n")
	}
	sb.WriteString("Use simple dollar amounts and preserve this family's economic meaning. Answers are derived locally. Add a teaching object mapping balanced_entry and counter_account to objects with hint and explanation strings. Hints should ask one short causal question without revealing the answer; explanations should explain cash versus earning and debit=left, credit=right. Do not include option letters. All teaching prose requires human review.\n")
	return sb.String()
}

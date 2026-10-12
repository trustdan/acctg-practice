package bank

import (
	"github.com/trustdan/acctg-practice/internal/domain"
	"time"
)

// QuestionStatus constants for versioned content workflow.
const (
	StatusActive                 = "active"
	StatusApprovedActive         = "approved_active"
	StatusSeedPendingReview      = "seed_pending_review"
	StatusRetired                = "retired"
	StatusDraft                  = "draft"
	StatusCandidatePendingReview = "candidate_pending_review"
	StatusCandidateRejected      = "rejected"
)

// ApprovalAction defines the explicit action taken on a candidate proposal or question.
type ApprovalAction string

const (
	ActionApprove ApprovalAction = "approve"
	ActionReject  ApprovalAction = "reject"
	ActionRepair  ApprovalAction = "repair"
	ActionRetire  ApprovalAction = "retire"
)

// ApprovalEvent records a durable audit record of a human semantic review, approval,
// rejection, repair, or retirement action.
//
// Invariants (AGENT-CONTRACT):
// 1. Stored separately from learner evidence (attempts, mastery projections).
// 2. Requires explicit reviewer identity, timestamp, source, rule version, and action.
// 3. Learner attempts or success rates cannot produce approval events.
type ApprovalEvent struct {
	ID          string         `json:"id"`
	CandidateID string         `json:"candidate_id"`
	QuestionID  string         `json:"question_id"`
	Action      ApprovalAction `json:"action"`
	Reviewer    string         `json:"reviewer"`
	RuleVersion int            `json:"rule_version"`
	Source      string         `json:"source"`
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

func IsValidStatus(s string) bool {
	switch s {
	case StatusActive, StatusApprovedActive, StatusSeedPendingReview, StatusRetired, StatusDraft, StatusCandidatePendingReview, StatusCandidateRejected:
		return true
	default:
		return false
	}
}

// IsActiveForPractice reports whether a question's lifecycle status permits practice selection.
// Retired and draft content is strictly excluded.
func IsActiveForPractice(status string) bool {
	return status == StatusActive || status == StatusApprovedActive || status == StatusSeedPendingReview
}

// RequiresReviewProvenance reports whether a status must carry reviewer and approval
// provenance. Reviewed teaching overrides apply exactly to these statuses.
func RequiresReviewProvenance(status string) bool {
	return status == StatusActive || status == StatusApprovedActive
}

// AccountsFile represents the JSON structure of curriculum/accounts.json.
type AccountsFile struct {
	SchemaVersion int           `json:"schema_version"`
	CourseScope   string        `json:"course_scope"`
	Accounts      []AccountJSON `json:"accounts"`
}

// AccountJSON represents a single account record in JSON.
type AccountJSON struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	NormalSide string `json:"normal_side"`
	ContraOf   string `json:"contra_of,omitempty"`
}

// QuestionBankFile represents the JSON structure of curriculum/seed-questions.json.
type QuestionBankFile struct {
	SchemaVersion int            `json:"schema_version"`
	Questions     []QuestionJSON `json:"questions"`
}

// QuestionJSON represents a single question template in JSON.
type QuestionJSON struct {
	Teaching         map[domain.DrillStage]TeachingText `json:"teaching,omitempty"`
	ID               string                             `json:"id"`
	Version          int                                `json:"version"`
	FamilyID         string                             `json:"family_id"`
	RuleVersion      int                                `json:"rule_version"`
	Status           string                             `json:"status"`
	ScenarioTemplate string                             `json:"scenario_template"`
	Parameters       map[string][]int64                 `json:"parameters"`
	Concepts         []string                           `json:"concepts"`
	Module           int                                `json:"module,omitempty"`
	Modules          []int                              `json:"modules,omitempty"`
	ExpectedFixture  FixtureJSON                        `json:"expected_fixture"`
	Review           ReviewJSON                         `json:"review"`
}

// CourseModule describes a course syllabus module for Foster School Financial Reporting & Analysis.
type CourseModule struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AllCourseModules returns the 10 canonical modules from the UW Foster syllabus.
func AllCourseModules() []CourseModule {
	return []CourseModule{
		{Number: 1, Title: "Module 1: Intro Transactions", Description: "Recording economic transactions & cash vs accrual"},
		{Number: 2, Title: "Module 2: Double-Entry Mechanics", Description: "Debits/Credits, T-accounts, 4 financial statements"},
		{Number: 3, Title: "Module 3: Statement Analysis & Cash Flows", Description: "Component percentages, ratios, direct/indirect SCF"},
		{Number: 4, Title: "Module 4: Financial Statements Project", Description: "Comprehensive accounting cycle & statement creation"},
		{Number: 5, Title: "Module 5: Revenue & Receivables", Description: "Revenue recognition, allowance method, Apple case"},
		{Number: 6, Title: "Module 6: Equity & Intangibles", Description: "Stock issuance, dividends, goodwill, Groupon pro-formas"},
		{Number: 7, Title: "Module 7: Long-Term Liabilities", Description: "Bonds payable, discount/premium, Disney 100-yr bonds"},
		{Number: 8, Title: "Module 8: Leverage & Crypto", Description: "Lehman Repo 105, monitoring, Tesla Bitcoin crypto"},
		{Number: 9, Title: "Module 9: Reporting Strategy & R&D", Description: "Aggressive vs conservative, Microsoft R&D capitalization"},
		{Number: 10, Title: "Module 10: Financial Forecasting", Description: "Aggregated vs disaggregated forecasting assumptions"},
	}
}

// QuestionModules returns the module numbers associated with a question.
// If explicit modules are set, they are returned; otherwise, modules are inferred
// from the question's family ID and concept tags.
func QuestionModules(q QuestionJSON) []int {
	if len(q.Modules) > 0 {
		return q.Modules
	}
	if q.Module > 0 {
		return []int{q.Module}
	}
	// Inferred mapping from family ID
	switch q.FamilyID {
	case FamilyCashService, FamilyCashRent, FamilyCashExpense,
		FamilyPrepaidPurchase, FamilyPrepaidConsumption, FamilyEquipmentPurchaseCash:
		return []int{1, 2}
	case FamilyCustomerAdvance, FamilyServiceOnCredit, FamilyCollectReceivable, FamilyEarnAdvance,
		FamilyBadDebtAllowance, FamilyWriteOffReceivable, FamilyAppleMultiElementRevenue:
		return []int{1, 2, 5}
	case FamilyIssueShares, FamilyDividendCash, FamilyDeclareDividend, FamilyPayDividendPayable,
		FamilyAcquireGoodwill, FamilyRepurchaseTreasuryStock, FamilyGoodwillImpairment:
		return []int{1, 2, 6}
	case FamilyBorrowCash, FamilyRepayNotePrincipal,
		FamilyIssueBondsPar, FamilyIssueBondsDiscount, FamilyIssueBondsPremium, FamilyPayBondInterest:
		return []int{1, 2, 7}
	case FamilyCreditSalesCashFlow:
		return []int{1, 2, 3}
	case FamilyCurrentRatioAnalysis:
		return []int{1, 2, 3, 8}
	case FamilyRepo105Analysis, FamilyDigitalAssetAccounting:
		return []int{1, 2, 8}
	case FamilyCapitalizeSoftwareFeasibility, FamilyReportingPolicyChoice,
		FamilyM09CookieJarReserve, FamilyM09WarrantyAccrual:
		return []int{1, 2, 9}
	case FamilyFinancialForecasting,
		FamilyM10ForecastingOperatingExpenseAP, FamilyM10ForecastingCapEx, FamilyM10ForecastingDebtPlug:
		return []int{1, 2, 10}
	case FamilyM03OCFAccountsPayableIncrease, FamilyM03OCFPrepaidIncrease, FamilyM03OCFUnearnedRevenueIncrease:
		return []int{1, 2, 3}
	case FamilyM03OCFQuickRatio:
		return []int{1, 2, 3, 8}
	case FamilyM04VariableConsideration, FamilyM04ContractAssetAR, FamilyM04RightOfReturn, FamilyM04SatisfactionOverTime:
		return []int{1, 2, 4}
	case FamilyM05InventoryFIFO, FamilyM05InventoryLIFO, FamilyM05InventoryNRVWriteDown, FamilyM05FreightInCapitalize:
		return []int{1, 2, 5}
	case FamilyM06StraightLineDepr, FamilyM06MACRSTaxTiming, FamilyM06AssetDisposalGainLoss, FamilyM06ImpairmentRecoverability:
		return []int{1, 2, 6}
	case FamilyM07BondIssuanceDiscount, FamilyM07BondIssuancePremium, FamilyM07BondEffectiveInterestAmort, FamilyM07BondContingentLiability:
		return []int{1, 2, 7}
	case FamilyM08LeaseCommencementROUAsset, FamilyM08LeasePeriodicPayment, FamilyM08TreasuryStockRepurchase, FamilyM08TreasuryStockReissuance:
		return []int{1, 2, 8}
	default:
		return []int{1, 2}
	}
}

// MatchesModules reports whether the question matches any of the selected modules.
// If selected is empty, all questions match.
func MatchesModules(q QuestionJSON, selected map[int]bool) bool {
	if len(selected) == 0 {
		return true
	}
	for _, m := range QuestionModules(q) {
		if selected[m] {
			return true
		}
	}
	return false
}

// FilterQuestionsByModules filters a slice of questions by selected module numbers.
func FilterQuestionsByModules(questions []QuestionJSON, selectedModules []int) []QuestionJSON {
	if len(selectedModules) == 0 {
		return questions
	}
	selMap := make(map[int]bool, len(selectedModules))
	for _, m := range selectedModules {
		selMap[m] = true
	}
	var filtered []QuestionJSON
	for _, q := range questions {
		if MatchesModules(q, selMap) {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

// TeachingText is reviewed wording only; it cannot change stage answers.
type TeachingText struct {
	Hint        string `json:"hint"`
	Explanation string `json:"explanation"`
}

// FixtureJSON defines the expected journal postings for testing and regression.
type FixtureJSON struct {
	Postings []FixturePostingJSON `json:"postings"`
}

// FixturePostingJSON defines one posting line in the expected fixture.
type FixturePostingJSON struct {
	AccountID       string `json:"account_id"`
	Side            string `json:"side"`
	AmountParameter string `json:"amount_parameter"`
}

// ReviewJSON captures human semantic review provenance.
type ReviewJSON struct {
	Reviewer   *string `json:"reviewer"`
	ApprovedAt *string `json:"approved_at"`
	Source     string  `json:"source"`
}

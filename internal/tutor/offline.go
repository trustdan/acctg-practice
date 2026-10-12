package tutor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// OfflineTutor provides deterministic, zero-network pedagogical tutoring.
// It respects context cancellation and deadlines.
type OfflineTutor struct{}

// NewOfflineTutor creates a new OfflineTutor.
func NewOfflineTutor() *OfflineTutor {
	return &OfflineTutor{}
}

// Name implements Tutor.
func (o *OfflineTutor) Name() string {
	return "OfflineTutor"
}

// Hint implements Tutor.
func (o *OfflineTutor) Hint(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	hintText := o.generateSocraticHint(req)
	tokens := (len(hintText) + 3) / 4

	return Response{
		Text:        hintText,
		Provider:    "offline",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// Explain implements Tutor.
func (o *OfflineTutor) Explain(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	explainText := o.generateExplanation(req)
	tokens := (len(explainText) + 3) / 4

	return Response{
		Text:        explainText,
		Provider:    "offline",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

func (o *OfflineTutor) generateSocraticHint(req Request) string {
	if req.MistakeHint != "" {
		return req.MistakeHint
	}
	// First check if there is an error tag from a submitted or inspected distractor
	if req.ErrorTag != "" {
		if tagHint := socraticHintForErrorTag(req.ErrorTag); tagHint != "" {
			return tagHint
		}
	}

	// If a vetted causal hint is provided in the request, prioritize its causal guidance
	if req.CausalHint != "" {
		return fmt.Sprintf("Consider the economic reality: %s", req.CausalHint)
	}

	// Stage-based Socratic guidance
	switch req.Stage {
	case domain.StageIdentifyAccount:
		return "Consider the economic event: which tangible resource, claim, or obligation changed first? Did cash change hands, or was a service provided?"
	case domain.StageAccountCategory:
		return "Classify the account under the accounting equation (Assets = Liabilities + Equity). Is it an economic resource owned (Asset), an obligation owed to third parties (Liability), or an owner claim / operational result (Equity)?"
	case domain.StageDirection:
		return "Did this specific account's balance increase or decrease as a result of this transaction?"
	case domain.StageDebitCredit:
		return "Remember the fundamental rule: Debit means Left and Credit means Right. Assets and Expenses increase on Debit (left); Liabilities, Equity, and Revenue increase on Credit (right)."
	case domain.StageCounterAccount:
		return "Every transaction has two sides to stay in balance. What balancing resource, claim, or obligation offsets the first account?"
	case domain.StageBalancedEntry:
		return "Check total debits and credits: every balanced journal entry must have total debits equal to total credits."
	case domain.StageEquationEffect:
		return "Reconcile the net effect on the accounting equation: does ΔAssets equal ΔLiabilities + ΔEquity?"
	default:
		return "Look at the timing of the economic exchange: separate when cash moves from when the service or benefit is earned or consumed."
	}
}

func (o *OfflineTutor) generateExplanation(req Request) string {
	var sb strings.Builder

	// If vetted explanation exists, use it as the foundation
	if req.Explanation != "" {
		return req.Explanation
	}

	// Pedagogical core principles
	sb.WriteString("Core Accounting Principles:\n")
	sb.WriteString("• Fundamental Equation: Assets = Liabilities + Equity (A = L + E).\n")
	sb.WriteString("• Debit & Credit Convention: Debit means Left and Credit means Right. Debit is not 'money out' or 'loss', and Credit is not 'money in'.\n")
	sb.WriteString("• Normal Balances: Assets and Expenses increase on Debit (left). Liabilities, Equity, and Revenue increase on Credit (right).\n")
	sb.WriteString("• Accrual Principle: Revenue records earning economic value, while Cash is a separate asset. Cash flow timing does not dictate revenue recognition or expense matching.")

	// Family-specific contrast guidance
	if contrast := familyContrastExplanation(req.FamilyID); contrast != "" {
		sb.WriteString("\n\nContrast Analysis:\n")
		sb.WriteString(contrast)
	}

	return sb.String()
}

func socraticHintForErrorTag(tag string) string {
	switch tag {
	case engine.TagDuplicateRevenueOnCollection:
		return "Notice that revenue was already recognized when the service was performed on credit. Today's event is collecting cash owed. Does collecting a receivable create new revenue, or does it swap Accounts Receivable for Cash?"
	case engine.TagRevenueRecognizedPrematurely:
		return "Consider the earning process: did the company perform the service today, or did the customer merely pay an advance deposit for future work? When cash is collected before performance, what obligation does the company owe?"
	case engine.TagCashRecordedOnEarningAdvance:
		return "Check cash timing: does completing the prepaid work bring a new payment, or was its receipt recorded in an earlier transaction? Completing the service reduces the unearned obligation."
	case engine.TagExpenseRecordedOnPrepaidPurchase:
		return "Consider the benefit horizon: does this payment cover only today, or upcoming months? Paying in advance for future coverage acquires an asset (Prepaid Insurance), not an expense of the current period."
	case engine.TagExpenseRecordedOnEquipmentPurchase:
		return "Equipment provides productive economic benefit over multiple years. Purchasing long-lived operational equipment is a capital asset acquisition (asset swap), not an immediate period expense."
	case engine.TagExpenseRecordedOnLoanRepayment:
		return "Repaying loan principal extinguishes a liability (Notes Payable) that was previously borrowed. Repaying principal is not an expense—only interest represents the cost of borrowing."
	case engine.TagExpenseRecordedOnDividend:
		return "Dividends distribute earnings to stockholders. They are recorded in the Dividends account, which reduces equity. A dividend is not an expense and does not reduce Net Income."
	case engine.TagRevenueRecordedOnBorrowing:
		return "Where did this cash come from? A lender expects it back. Does money the company must repay mean it earned something, or that it owes something?"
	case engine.TagRevenueRecordedOnShareIssue:
		return "Who provided this cash, and what did they receive in return? Did the company perform work for it, or did investors buy an ownership stake?"
	case engine.TagPrepaidNotExpensedOnConsumption:
		return "Part of the coverage the company paid for has now been used up. Should the full amount still be reported as a future benefit?"
	case engine.TagCashRecordedOnDividendDeclaration:
		return "Did any cash leave the company on the declaration date, or did the company only commit to pay later?"
	case engine.TagDividendsDebitedOnPayment:
		return "Was this dividend already declared and recorded earlier? If so, what does today's payment settle?"
	case engine.TagCashRecordedWhenUncollected:
		return "Did the customer hand over any money today? If not, what does the company hold instead of cash?"
	case engine.TagCashRecordedOnPrepaidExpiration:
		return "Did any cash leave the company today, or was the coverage paid for when the policy was bought?"
	case engine.TagRevenueDeferredWhenEarned:
		return "Has the company already done the work? If the work is finished, does the company still owe the customer anything?"
	case engine.TagAdvanceConfusedWithReceivable:
		return "Who owes whom? The customer has already paid. Does the customer owe the company money, or does the company owe the customer future work?"
	case engine.TagPayableRecordedForCashPayment:
		return "Did the company pay today or promise to pay later? If cash already left, is anything still owed?"
	case engine.TagEquationEffectMissed:
		return "Is this only an exchange of one asset for another, or did a liability or equity account change too? Place each account in the entry under Assets, Liabilities, or Equity."
	case engine.TagWrongAmount:
		return "Which amount belongs to today's event: the stated payment or portion, or a larger total that includes other transactions?"
	case engine.TagWrongAccount:
		return "Re-read the scenario: what changed hands today, and who owes whom afterward? Choose the account that describes that change."
	case engine.TagReversedSides:
		return "Check side placement: Debit is always the Left side, and Credit is always the Right side. Verify whether this account increases or decreases, and apply its normal side rule."
	case engine.TagUnbalancedEntry:
		return "Total debits must strictly equal total credits. Check the amounts to ensure both sides balance."
	case engine.TagDirectWriteOffExpense:
		return "Do we know yet which specific customers will default? Under the allowance method, which contra-asset account establishes the credit reserve under GAAP matching?"
	case engine.TagWriteOffExpensedAgain:
		return "Was Bad Debt Expense already estimated and recorded earlier? Does writing off a specific customer account reduce net income a second time or reduce the established allowance reserve?"
	case engine.TagStockIssuedInsteadOfBond:
		return "Did investors buy ownership shares in the company, or did they lend money expecting repayment of principal and interest? Choose the liability account that records long-term debt obligations."
	case engine.TagGoodwillExpensed:
		return "Does the acquired business provide ongoing future economic benefits? Under GAAP, should the premium paid over net identifiable assets be capitalized as an intangible asset or expensed today?"
	case engine.TagCashFlowOperatingWrongDirection:
		return "Net income includes this credit revenue, but did the company receive cash? On the indirect Statement of Cash Flows, should uncollected revenue be added or deducted from Net Income?"
	case engine.TagDiscountRecordedAsLoss:
		return "Does issuing bonds below par represent an immediate operating loss, or does the discount represent additional borrowing cost amortized over the bond term as interest expense?"
	case engine.TagExpensingFeasibleSoftware:
		return "Under US GAAP (ASC 985-20), once technological feasibility is established, coding and testing costs provide demonstrable future economic benefit. Should these costs be expensed immediately or capitalized as an intangible asset?"
	case engine.TagCurrentRatioWrongDirection:
		return "Both Current Assets (Cash) and Current Liabilities (Accounts Payable) decreased by the same dollar amount. If the initial Current Ratio was greater than 1.0, what mathematical effect does subtracting equal amounts from numerator and denominator have?"
	case engine.TagTreasuryStockRecordedAsAsset:
		return "Can a corporation hold an asset representing an ownership claim on itself? Under US GAAP, should repurchased shares be recorded as contra-equity reducing stockholders' equity?"
	case engine.TagIgnoringBondAmortization:
		return "Does issuing bonds at a discount mean total borrowing cost equals only the cash coupon paid, or must discount amortization be added to interest expense?"
	case engine.TagRepoTreatedAsSale:
		return "Does transferring liquid collateral with an obligation to repurchase it days later represent a true sale, or a collateralized financing liability (Notes Payable)?"
	case engine.TagPrematureRevenueRecognition:
		return "Has the company fulfilled all performance obligations today, or must consideration for future deliverables (such as multi-period software updates) be deferred in Unearned Revenue?"
	case engine.TagGoodwillAmortized:
		return "Under US GAAP (ASC 350), does goodwill have a finite life subject to periodic amortization, or is it an indefinite-lived asset tested annually for impairment?"
	case engine.TagCryptoRecordedAsCashEquivalent:
		return "Under US GAAP, are cryptocurrencies classified as cash or cash equivalents, or as indefinite-lived intangible assets subject to impairment testing?"
	case engine.TagAggressiveRevenueAcceleration:
		return "Under conservative financial reporting policy, can revenue be recognized immediately if ongoing performance obligations remain to be delivered?"
	case engine.TagAggregatedForecastAssumption:
		return "In financial forecasting, do crude top-line percentage growth extrapolations provide the same structural accuracy as disaggregated operational drivers (unit volume times price)?"
	case engine.TagM03OCFPayableWrongDirection:
		return "When an operating expense was incurred on account, did it reduce Net Income without any cash leaving the firm? Should this non-cash payable increase be added or deducted to reconcile to cash flow?"
	case engine.TagM03OCFPrepaidWrongDirection:
		return "If cash was paid out for insurance coverage but not yet expensed on the income statement, must this asset increase be deducted from Net Income to calculate operating cash flow?"
	case engine.TagM03OCFUnearnedRevenueWrongDirection:
		return "If customers paid cash in advance that was not yet recognized in Net Income, must this liability increase be added to Net Income to reflect operating cash inflow?"
	case engine.TagM03QuickRatioWrongDirection:
		return "When both Quick Assets (Cash) and Current Liabilities (A/P) decrease by the identical dollar amount, and the ratio was greater than 1.0, does the ratio increase or decrease?"
	case engine.TagM03InventoryTreatedAsQuickAsset:
		return "Does the acid-test (Quick) ratio include less liquid operating assets like inventory, or only cash, marketable securities, and accounts receivable?"
	case engine.TagM04RevenueRecognizedWithUncertainReturn:
		return "Under ASC 606, if the customer holds a full right of return and significant reversal is possible, can revenue be recognized immediately, or must a contract liability (Unearned Revenue) be recorded?"
	case engine.TagM04CashDebitedOnOverTimeSatisfaction:
		return "Was cash collected again today, or is the company fulfilling an advance obligation that was already collected and recorded in Unearned Revenue?"
	case engine.TagM05InventoryRecordedAtSelling:
		return "When merchandise is sold under perpetual inventory, is the inventory reduction recorded at cost or at retail selling price?"
	case engine.TagM05FreightInExpensed:
		return "Under GAAP, are shipping costs necessary to acquire and receive merchandise inventory expensed immediately, or capitalized as part of the Inventory asset cost?"
	case engine.TagM06AccumulatedDeprExpensed:
		return "Is Accumulated Depreciation an operating expense on the income statement, or a contra-asset account offsetting the historical cost of equipment on the balance sheet?"
	case engine.TagM06DirectAssetCredit:
		return "Under US GAAP, does periodic depreciation directly reduce the original equipment cost account, or is it accumulated in a separate contra-asset account?"
	case engine.TagM06GainTreatedAsRevenue:
		return "Does selling an operational machine generate core recurring customer revenue, or a non-operating peripheral gain on disposal?"
	case engine.TagM06DisposalAccumDeprRetained:
		return "When an asset is retired and removed from the books, must the cumulative depreciation associated with that specific asset also be removed?"
	case engine.TagM07DiscountRecordedAsLoss:
		return "Does a discount on bonds payable represent an immediate period loss, or an adjustment to borrowing cost amortized over the life of the bond?"
	case engine.TagM07PremiumRecordedAsGain:
		return "Does receiving a bond premium count as immediate upfront gain, or does it reduce interest expense across the life of the debt?"
	case engine.TagM07IgnoringBondAmortization:
		return "Under the effective interest method, does interest expense equal only the cash coupon payment, or must the difference amortize the bond discount or premium?"
	case engine.TagM07ContingentRecordedAsPayable:
		return "Does an accrued litigation loss represent trade accounts payable to routine suppliers, or a specific contingent litigation liability?"
	case engine.TagM08LeaseExpensedAtCommencement:
		return "Under ASC 842, can a long-term lease commitment be expensed immediately at signing, or must a Right-of-Use Asset and Lease Liability be capitalized?"
	case engine.TagM08ROUAssetRecordedAsEquity:
		return "Does a lease agreement represent an equity ownership interest, or a contractual financing obligation owed to a lessor?"
	case engine.TagM08IgnoringLeaseInterest:
		return "Does a lease payment solely reduce the principal obligation, or does it cover both financing interest expense and liability reduction?"
	case engine.TagM08TreasuryStockRecordedAsAsset:
		return "Can a company hold an investment asset in its own stock, or is repurchased stock classified as a contra-equity account?"
	case engine.TagM08GainRecordedOnTreasuryReissuance:
		return "Can transactions in a corporation's own capital stock ever generate income statement gains or losses under US GAAP?"
	case engine.TagM09DiscretionaryReserveOveraccrual:
		return "Can management intentionally overstate uncollectible reserves in profitable quarters to buffer future earnings without compromising reporting quality?"
	case engine.TagM09WarrantyExpensedWhenPaid:
		return "Under accrual accounting and the matching principle, must product warranty costs be estimated and accrued in the period of sale rather than delayed until cash repairs are performed?"
	case engine.TagM10ForecastingAPOmitted:
		return "When forecasting operating expenses on trade credit, does the entry credit Cash or the trade liability Accounts Payable?"
	case engine.TagM10CapExExpensedImmediately:
		return "Do multi-period capital expenditures for productive machinery belong on the income statement as immediate expense, or on the balance sheet as capitalized Equipment?"
	case engine.TagM10EFNTreatedAsEquityWithoutShares:
		return "When pro-forma asset expansion requires external debt borrowing, does the financing plug credit Notes Payable or Common Stock?"
	default:
		return ""
	}
}

func familyContrastExplanation(familyID string) string {
	switch familyID {
	case "collect_receivable":
		return "Contrast: Performing a service on account recognizes Revenue and Accounts Receivable. Later, collecting cash increases Cash (Debit) and decreases Accounts Receivable (Credit). It is an asset swap and never recognizes duplicate revenue."
	case "customer_advance":
		return "Contrast: Receiving cash before performing work creates a liability (Unearned Revenue). Revenue is earned only later when the service or goods are delivered to the customer."
	case "earn_advance":
		return "Contrast: When fulfilling an advance, no new cash is received in the completion transaction. Instead, the liability Unearned Revenue is debited (decreased), and Service Revenue is credited (increased)."
	case "prepaid_purchase":
		return "Contrast: Paying now for future coverage exchanges one asset for another (Debit Prepaid Insurance, Credit Cash). Insurance Expense is recognized month by month as the coverage expires."
	case "prepaid_consumption":
		return "Contrast: As prepaid coverage expires, record Insurance Expense (Debit) and reduce the asset Prepaid Insurance (Credit). No cash moves when coverage expires."
	case "equipment_purchase_cash":
		return "Contrast: Long-lived equipment is capitalized as an Asset, not expensed at purchase. It is expensed gradually over its useful life through depreciation."
	case "repay_note_principal":
		return "Contrast: Paying down loan principal reduces the liability Notes Payable (Debit) with Cash (Credit). It is a liability reduction, not an expense."
	case "dividend_cash":
		return "Contrast: Dividends distribute earnings to stockholders. Debit the Dividends account, which reduces equity and is closed to Retained Earnings at period end. Dividends are not an expense and do not appear on the Income Statement."
	case "repurchase_treasury_stock":
		return "Contrast: Repurchasing company common stock reduces Cash (Credit) and increases Treasury Stock (Debit), a contra-equity account reducing stockholders' equity. It is not an outside asset investment."
	case "pay_bond_interest":
		return "Contrast: When bonds are issued at a discount, total Interest Expense (Debit) exceeds the cash coupon paid (Credit Cash). The difference is credited to Discount on Bonds Payable to amortize the contra-liability toward par."
	case "repo105_analysis":
		return "Contrast: Under US GAAP economic substance, short-term repurchase agreements are collateralized loans (Debit Cash, Credit Notes Payable), maintaining leverage transparency rather than being treated as asset sales."
	case "apple_multielement_revenue":
		return "Contrast: In bundled arrangements with future software updates, only the delivered hardware value is recognized as point-of-sale revenue; the undelivered software obligations are credited to Unearned Revenue and recognized over time."
	case "goodwill_impairment":
		return "Contrast: Goodwill is not amortized. When a reporting unit's fair value falls below carrying value, an impairment loss is debited (-Equity) and Goodwill is credited (-Assets). Impairment cannot be reversed under US GAAP."
	case "digital_asset_accounting":
		return "Contrast: Under US GAAP (Tesla Bitcoin), digital assets are classified as indefinite-lived intangibles (Debit Digital Assets, Credit Cash). Impairment losses are recognized when prices decline, but holding gains cannot be recognized until sold."
	case "reporting_policy_choice":
		return "Contrast: A conservative reporting strategy defers revenue to Unearned Revenue until all performance obligations are fulfilled, preserving earnings quality, whereas an aggressive policy prematurely accelerates revenue recognition."
	case "financial_forecasting":
		return "Contrast: Pro-forma forecasting built from disaggregated operational drivers (unit volume, selling price per unit) provides actionable visibility into capacity constraints and working capital, unlike crude aggregated growth percentages."
	case "m09_qual_cookie_jar_reserve":
		return "Contrast: Management discretion in setting accounting estimates should reflect bona fide economic expectations; intentionally over-accruing bad debt reserves creates a cookie jar reserve that obscures true earnings quality."
	case "m09_qual_warranty_accrual":
		return "Contrast: Under the matching principle (ASC 460), estimated future product warranty claims are matched against sales revenue in the current period by debiting Warranty Expense and crediting Warranty Liability."
	case "m10_fore_operating_expense_ap":
		return "Contrast: In pro-forma forecasting, linking operating expenses to trade payables ensures working capital metrics (such as DPO) realistically reflect the timing of operating cash outflows."
	case "m10_fore_capex_equipment_cash":
		return "Contrast: Capital expenditures (CapEx) purchase long-term productive assets (Debit Equipment, Credit Cash), expanding balance sheet capacity and flowing into investing cash outflows rather than operating expenses."
	case "m10_fore_debt_plug_cash_notes":
		return "Contrast: External Financing Needed (EFN) reconciles projected asset growth with projected spontaneous funding; plugging EFN with debt issuance (Debit Cash, Credit Notes Payable) connects the balance sheet directly to financing cash flows."
	default:
		return ""
	}
}

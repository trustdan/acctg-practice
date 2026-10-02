package statements

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// CanonicalCasePioneerConsulting provides the standard comprehensive case study for Stage 15.
// Features: opening trial balance, 8 operational events, 2 adjusting entries,
// distinct dividend declaration and payment, full statements, and cash flow reconciliation.
func CanonicalCasePioneerConsulting() ComprehensiveCase {
	actOperating := ActivityOperating
	actInvesting := ActivityInvesting
	actFinancing := ActivityFinancing

	return ComprehensiveCase{
		CaseID:      "case_pioneer_m1",
		Title:       "Pioneer Consulting Services — Month 1 Comprehensive Case",
		Description: "Comprehensive month 1 accounting cycle: initial capitalization, operating transactions, adjustments, financial statements, and cash flow reconciliation.",
		CompanyName: "Pioneer Consulting Services, Inc.",
		Period:      "For the Month Ended January 31, 2026",
		AsOfDate:    "January 31, 2026",
		OpeningBalances: map[domain.AccountID]domain.Money{
			"cash":              domain.NewMoney(2000000), // $20,000.00
			"common_stock":      domain.NewMoney(2000000), // $20,000.00
			"retained_earnings": domain.ZeroMoney(),
		},
		Transactions: []CaseTransaction{
			{
				Date:        "Jan 02",
				Description: "Prepaid 12-month commercial office insurance in cash",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyPrepaidPurchase,
					Parameters: map[string]int64{"amount_minor_units": 240000}, // $2,400.00
				},
				CashClassification: &actOperating,
				CashDescription:    "Cash paid for prepaid insurance",
			},
			{
				Date:        "Jan 05",
				Description: "Purchased computer and office equipment for cash",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyEquipmentPurchaseCash,
					Parameters: map[string]int64{"amount_minor_units": 500000}, // $5,000.00
				},
				CashClassification: &actInvesting,
				CashDescription:    "Cash paid for office equipment acquisition",
			},
			{
				Date:        "Jan 08",
				Description: "Received cash advance from client for 3-month consulting retainer",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyCustomerAdvance,
					Parameters: map[string]int64{"amount_minor_units": 300000}, // $3,000.00
				},
				CashClassification: &actOperating,
				CashDescription:    "Cash collected in advance from customers (unearned)",
			},
			{
				Date:        "Jan 12",
				Description: "Provided consulting advisory services on credit to client",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyServiceOnCredit,
					Parameters: map[string]int64{"amount_minor_units": 450000}, // $4,500.00
				},
				// Noncash revenue transaction
			},
			{
				Date:        "Jan 17",
				Description: "Collected cash on account from client billed on Jan 12",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyCollectReceivable,
					Parameters: map[string]int64{"amount_minor_units": 250000}, // $2,500.00
				},
				CashClassification: &actOperating,
				CashDescription:    "Collections from accounts receivable",
			},
			{
				Date:        "Jan 20",
				Description: "Borrowed cash from regional bank on a 2-year note payable",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyBorrowCash,
					Parameters: map[string]int64{"amount_minor_units": 600000}, // $6,000.00
				},
				CashClassification: &actFinancing,
				CashDescription:    "Cash proceeds from issuance of notes payable",
			},
			{
				Date:        "Jan 22",
				Description: "Paid monthly office rent expense in cash",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyCashRent,
					Parameters: map[string]int64{"amount_minor_units": 120000}, // $1,200.00
				},
				CashClassification: &actOperating,
				CashDescription:    "Cash paid for office rent",
			},
			{
				Date:        "Jan 25",
				Description: "Repaid note principal to bank in cash",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyRepayNotePrincipal,
					Parameters: map[string]int64{"amount_minor_units": 100000}, // $1,000.00
				},
				CashClassification: &actFinancing,
				CashDescription:    "Principal repayment on notes payable",
			},
			{
				Date:        "Jan 27",
				Description: "Board of directors declared $500 dividend payable to shareholders",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyDeclareDividend,
					Parameters: map[string]int64{"amount_minor_units": 50000}, // $500.00
				},
				// Declaration: noncash! Incurs liability Dividends Payable and reduces Retained Earnings.
			},
			{
				Date:        "Jan 28",
				Description: "Paid partial cash distribution of $300 towards declared dividend obligation",
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyPayDividendPayable,
					Parameters: map[string]int64{"amount_minor_units": 30000}, // $300.00
				},
				CashClassification: &actFinancing,
				CashDescription:    "Cash paid for dividends",
			},
		},
		Adjustments: []CaseTransaction{
			{
				Date:         "Jan 31",
				Description:  "Adjusting Entry: Expired 1 month of prepaid insurance ($2,400 / 12 months)",
				IsAdjustment: true,
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyPrepaidConsumption,
					Parameters: map[string]int64{"amount_minor_units": 20000}, // $200.00
				},
			},
			{
				Date:         "Jan 31",
				Description:  "Adjusting Entry: Earned 1 month of customer advance retainer ($3,000 / 3 months)",
				IsAdjustment: true,
				Event: engine.TransactionEvent{
					FamilyID:   bank.FamilyEarnAdvance,
					Parameters: map[string]int64{"amount_minor_units": 100000}, // $1,000.00
				},
			},
		},
	}
}

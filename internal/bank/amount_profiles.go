package bank

// Amount profiles are reviewed for these exact scenarios; unspecified totals remain symbolic.
type AmountDistractor struct {
	FamilyID string
	Debit    string
	Credit   string
	Scope    string
	Hint     string
}

var amountDistractors = map[string]AmountDistractor{
	"amount_collect_receivable_installment":      {"collect_receivable", "Cash", "Accounts Receivable", "the full original invoice, including the amount still unpaid", "Did today's payment settle the whole invoice, or only the stated installment?"},
	"amount_customer_advance_additional_receipt": {"customer_advance", "Cash", "Unearned Revenue", "the entire course price, including the earlier recorded receipt", "Which receipt is new today, and which receipt is already in the books?"},
	"amount_earn_advance_completed_session":      {"earn_advance", "Unearned Revenue", "Service Revenue", "the entire prepaid course, including sessions still owed", "Did the company complete every prepaid session, or only the separately priced session stated today?"},
	"amount_earn_advance_final_remaining":        {"earn_advance", "Unearned Revenue", "Service Revenue", "the original full course price, including earlier recorded earning", "What obligation remained immediately before the final session after earlier earning was recorded?"},
	"amount_prepaid_consumption_partial_policy":  {"prepaid_consumption", "Insurance Expense", "Prepaid Insurance", "the full policy premium, including unused future coverage", "Has all purchased protection expired, or does some coverage remain available?"},
	"amount_prepaid_consumption_final_remaining": {"prepaid_consumption", "Insurance Expense", "Prepaid Insurance", "the original full premium, including earlier recorded consumption", "What prepaid balance remained before this entry after earlier months were expensed?"},
	"amount_prepaid_purchase_discounted_policy":  {"prepaid_purchase", "Prepaid Insurance", "Cash", "the higher advertised premium before the agreed price reduction", "Was the advertised price ever owed or paid, or was the reduced agreed premium the entire purchase cost?"},
	"amount_repay_principal_early_partial":       {"repay_note_principal", "Notes Payable", "Cash", "the entire principal balance before payment, including principal still owed", "Did the bank receive the entire outstanding principal, or only today's accepted partial payment?"},
	"amount_dividend_cash_total_distribution":    {"dividend_cash", "Dividends", "Cash", "the stated total multiplied by the number of recipients", "Is the stated amount for each stockholder, or the total paid by the corporation to everyone together?"},
	"amount_equipment_cash_additional_machine":   {"equipment_purchase_cash", "Equipment", "Cash", "the combined cost of the new and previously recorded machines", "Which machine was purchased today, and which cost is already recorded?"},
}

func ReviewedAmountDistractor(profile string) (AmountDistractor, bool) {
	p, ok := amountDistractors[profile]
	return p, ok
}

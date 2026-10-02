package domain_test

import (
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestMoney(t *testing.T) {
	m1 := domain.NewMoney(5000)
	if m1.Cents() != 5000 {
		t.Fatalf("expected 5000 cents, got %d", m1.Cents())
	}
	if m1.Dollars() != 50.0 {
		t.Fatalf("expected 50.0 dollars, got %f", m1.Dollars())
	}
	if m1.FormatDollars() != "$50" {
		t.Fatalf("expected $50, got %s", m1.FormatDollars())
	}
	if m1.FormatExact() != "$50.00" {
		t.Fatalf("expected $50.00, got %s", m1.FormatExact())
	}

	m2 := domain.NewMoney(1250)
	if m2.FormatDollars() != "$12.50" {
		t.Fatalf("expected $12.50, got %s", m2.FormatDollars())
	}

	neg := domain.NewMoney(-5000)
	if neg.FormatDollars() != "-$50" {
		t.Fatalf("expected -$50, got %s", neg.FormatDollars())
	}
	if neg.FormatExact() != "-$50.00" {
		t.Fatalf("expected -$50.00, got %s", neg.FormatExact())
	}
	if neg.Abs() != m1 {
		t.Fatalf("expected abs to be 5000, got %d", neg.Abs())
	}

	sum := m1.Add(m2)
	if sum.Cents() != 6250 {
		t.Fatalf("expected 6250, got %d", sum.Cents())
	}

	diff := m1.Sub(m2)
	if diff.Cents() != 3750 {
		t.Fatalf("expected 3750, got %d", diff.Cents())
	}

	fromDC, err := domain.FromDollarsAndCents(100, 50)
	if err != nil || fromDC.Cents() != 10050 {
		t.Fatalf("unexpected FromDollarsAndCents result: %v, err=%v", fromDC, err)
	}

	_, err = domain.FromDollarsAndCents(100, 150)
	if err == nil {
		t.Fatalf("expected error for cents > 99")
	}
}

func TestCategoryAndNormalSide(t *testing.T) {
	tests := []struct {
		cat      domain.Category
		normal   domain.Side
		expected bool
	}{
		{domain.CategoryAsset, domain.SideDebit, true},
		{domain.CategoryLiability, domain.SideCredit, true},
		{domain.CategoryEquity, domain.SideCredit, true},
		{domain.CategoryRevenue, domain.SideCredit, true},
		{domain.CategoryExpense, domain.SideDebit, true},
		{domain.CategoryDividends, domain.SideDebit, true},
		{domain.Category("crypto"), domain.SideUnknown, false},
	}

	for _, tt := range tests {
		if tt.cat.IsValid() != tt.expected {
			t.Errorf("category %s IsValid expected %v, got %v", tt.cat, tt.expected, tt.cat.IsValid())
		}
		if tt.cat.NormalSide() != tt.normal {
			t.Errorf("category %s NormalSide expected %s, got %s", tt.cat, tt.normal, tt.cat.NormalSide())
		}
	}
}

func TestSideForDirection(t *testing.T) {
	tests := []struct {
		normal    domain.Side
		direction domain.Direction
		expected  domain.Side
	}{
		{domain.SideDebit, domain.DirectionIncrease, domain.SideDebit},
		{domain.SideDebit, domain.DirectionDecrease, domain.SideCredit},
		{domain.SideCredit, domain.DirectionIncrease, domain.SideCredit},
		{domain.SideCredit, domain.DirectionDecrease, domain.SideDebit},
	}

	for _, tt := range tests {
		got := domain.SideForDirection(tt.normal, tt.direction)
		if got != tt.expected {
			t.Errorf("SideForDirection(%s, %s) = %s; want %s", tt.normal, tt.direction, got, tt.expected)
		}
		revDir := domain.DirectionForSide(tt.normal, got)
		if revDir != tt.direction {
			t.Errorf("DirectionForSide(%s, %s) = %s; want %s", tt.normal, got, revDir, tt.direction)
		}
	}
}

func TestAccountCatalog(t *testing.T) {
	catalog := domain.NewAccountCatalog()

	acc1 := domain.Account{
		ID:         "cash",
		Name:       "Cash",
		Category:   domain.CategoryAsset,
		NormalSide: domain.SideDebit,
	}
	if err := catalog.Add(acc1); err != nil {
		t.Fatalf("unexpected error adding account: %v", err)
	}

	// Duplicate add must fail
	if err := catalog.Add(acc1); err == nil {
		t.Fatalf("expected error adding duplicate account")
	}

	// Mismatched normal side must fail
	badSide := domain.Account{
		ID:         "bad_cash",
		Name:       "Bad Cash",
		Category:   domain.CategoryAsset,
		NormalSide: domain.SideCredit, // Asset normal side is debit!
	}
	if err := catalog.Add(badSide); err == nil {
		t.Fatalf("expected error adding account with invalid normal side")
	}

	// Empty name must fail
	badName := domain.Account{
		ID:         "unnamed",
		Name:       "",
		Category:   domain.CategoryAsset,
		NormalSide: domain.SideDebit,
	}
	if err := catalog.Add(badName); err == nil {
		t.Fatalf("expected error adding account with empty name")
	}

	got, found := catalog.Get("cash")
	if !found || got.Name != "Cash" {
		t.Fatalf("expected to find Cash account, got %+v", got)
	}
}

func setupTestCatalog() *domain.AccountCatalog {
	catalog := domain.NewAccountCatalog()
	accounts := []domain.Account{
		{ID: "cash", Name: "Cash", Category: domain.CategoryAsset, NormalSide: domain.SideDebit},
		{ID: "accounts_receivable", Name: "Accounts Receivable", Category: domain.CategoryAsset, NormalSide: domain.SideDebit},
		{ID: "equipment", Name: "Equipment", Category: domain.CategoryAsset, NormalSide: domain.SideDebit},
		{ID: "unearned_revenue", Name: "Unearned Revenue", Category: domain.CategoryLiability, NormalSide: domain.SideCredit},
		{ID: "notes_payable", Name: "Notes Payable", Category: domain.CategoryLiability, NormalSide: domain.SideCredit},
		{ID: "common_stock", Name: "Common Stock", Category: domain.CategoryEquity, NormalSide: domain.SideCredit},
		{ID: "service_revenue", Name: "Service Revenue", Category: domain.CategoryRevenue, NormalSide: domain.SideCredit},
		{ID: "rent_expense", Name: "Rent Expense", Category: domain.CategoryExpense, NormalSide: domain.SideDebit},
		{ID: "dividends", Name: "Dividends", Category: domain.CategoryDividends, NormalSide: domain.SideDebit},
	}
	for _, a := range accounts {
		_ = catalog.Add(a)
	}
	return catalog
}

func TestEntryBalancingAndEquationEffects(t *testing.T) {
	catalog := setupTestCatalog()

	t.Run("Balanced cash service: Dr Cash 100, Cr Service Revenue 100", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(10000)},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(10000)},
		)
		if !entry.IsBalanced() {
			t.Fatalf("expected entry to be balanced")
		}
		delta, err := entry.EquationEffects(catalog)
		if err != nil {
			t.Fatalf("unexpected equation error: %v", err)
		}
		if delta.DeltaAssets.Cents() != 10000 || delta.DeltaLiabilities.Cents() != 0 || delta.DeltaEquity.Cents() != 10000 {
			t.Fatalf("unexpected delta: %+v", delta)
		}
	})

	t.Run("Customer advance: Dr Cash 50, Cr Unearned Revenue 50", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(5000)},
			domain.Posting{AccountID: "unearned_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(5000)},
		)
		delta, err := entry.EquationEffects(catalog)
		if err != nil {
			t.Fatalf("unexpected equation error: %v", err)
		}
		if delta.DeltaAssets.Cents() != 5000 || delta.DeltaLiabilities.Cents() != 5000 || delta.DeltaEquity.Cents() != 0 {
			t.Fatalf("unexpected delta: %+v", delta)
		}
	})

	t.Run("Receivable collection: Dr Cash 75, Cr Accounts Receivable 75", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(7500)},
			domain.Posting{AccountID: "accounts_receivable", Side: domain.SideCredit, Amount: domain.NewMoney(7500)},
		)
		delta, err := entry.EquationEffects(catalog)
		if err != nil {
			t.Fatalf("unexpected equation error: %v", err)
		}
		// Asset swap: net asset change is 0!
		if delta.DeltaAssets.Cents() != 0 || delta.DeltaLiabilities.Cents() != 0 || delta.DeltaEquity.Cents() != 0 {
			t.Fatalf("unexpected delta: %+v", delta)
		}
	})

	t.Run("Cash rent expense: Dr Rent Expense 200, Cr Cash 200", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "rent_expense", Side: domain.SideDebit, Amount: domain.NewMoney(20000)},
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(20000)},
		)
		delta, err := entry.EquationEffects(catalog)
		if err != nil {
			t.Fatalf("unexpected equation error: %v", err)
		}
		if delta.DeltaAssets.Cents() != -20000 || delta.DeltaLiabilities.Cents() != 0 || delta.DeltaEquity.Cents() != -20000 {
			t.Fatalf("unexpected delta: %+v", delta)
		}
	})

	t.Run("Cash dividends: Dr Dividends 30, Cr Cash 30", func(t *testing.T) {
		entry := domain.NewEntry(
			domain.Posting{AccountID: "dividends", Side: domain.SideDebit, Amount: domain.NewMoney(3000)},
			domain.Posting{AccountID: "cash", Side: domain.SideCredit, Amount: domain.NewMoney(3000)},
		)
		delta, err := entry.EquationEffects(catalog)
		if err != nil {
			t.Fatalf("unexpected equation error: %v", err)
		}
		if delta.DeltaAssets.Cents() != -3000 || delta.DeltaLiabilities.Cents() != 0 || delta.DeltaEquity.Cents() != -3000 {
			t.Fatalf("unexpected delta: %+v", delta)
		}
	})

	t.Run("Unbalanced entry rejected", func(t *testing.T) {
		unbalanced := domain.NewEntry(
			domain.Posting{AccountID: "cash", Side: domain.SideDebit, Amount: domain.NewMoney(10000)},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(5000)},
		)
		if unbalanced.IsBalanced() {
			t.Fatalf("expected entry to NOT be balanced")
		}
		if err := unbalanced.Validate(catalog); err == nil {
			t.Fatalf("expected validation error for unbalanced entry")
		}
	})

	t.Run("Unknown account rejected", func(t *testing.T) {
		badAcc := domain.NewEntry(
			domain.Posting{AccountID: "unknown_account", Side: domain.SideDebit, Amount: domain.NewMoney(10000)},
			domain.Posting{AccountID: "service_revenue", Side: domain.SideCredit, Amount: domain.NewMoney(10000)},
		)
		if err := badAcc.Validate(catalog); err == nil {
			t.Fatalf("expected validation error for unknown account")
		}
	})
}

func TestAttemptValidation(t *testing.T) {
	valid := domain.Attempt{
		AttemptID:        "att-1",
		SessionID:        "sess-1",
		InstanceID:       "inst-1",
		QuestionID:       "q-1",
		QuestionVersion:  1,
		Stage:            domain.StageIdentifyAccount,
		ConceptID:        "cash_vs_revenue",
		SelectedOptionID: "opt-1",
		IsCorrect:        true,
		Assistance:       domain.AssistanceNone,
		GradingVersion:   1,
		AnsweredAt:       time.Now().UTC(),
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	invalid := valid
	invalid.AttemptID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for empty AttemptID")
	}

	invalid = valid
	invalid.GradingVersion = 0
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for non-positive GradingVersion")
	}
}

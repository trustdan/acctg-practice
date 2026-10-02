package bank_test

import (
	"strings"
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

func TestCanonicalAccountsValidation(t *testing.T) {
	catalog, file, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("expected canonical curriculum/accounts.json to validate, got error: %v", err)
	}
	if file.SchemaVersion != 1 {
		t.Fatalf("expected schema version 1, got %d", file.SchemaVersion)
	}
	if catalog.Count() != 14 {
		t.Fatalf("expected 14 accounts in canonical catalog, got %d", catalog.Count())
	}

	// Verify key accounts exist
	required := []domain.AccountID{
		"cash",
		"accounts_receivable",
		"prepaid_insurance",
		"equipment",
		"accounts_payable",
		"unearned_revenue",
		"notes_payable",
		"dividends_payable",
		"common_stock",
		"retained_earnings",
		"service_revenue",
		"rent_expense",
		"insurance_expense",
		"dividends",
	}
	for _, id := range required {
		acc, ok := catalog.Get(id)
		if !ok {
			t.Errorf("expected account %s to exist in canonical catalog", id)
		}
		if acc.Name == "" {
			t.Errorf("expected account %s to have non-empty name", id)
		}
	}
}

func TestCanonicalSeedQuestionsValidation(t *testing.T) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}

	qBank, err := bank.LoadQuestionBankFile("../../curriculum/seed-questions.json", catalog)
	if err != nil {
		t.Fatalf("expected canonical curriculum/seed-questions.json to validate, got error: %v", err)
	}
	if qBank.SchemaVersion != 1 {
		t.Fatalf("expected schema version 1, got %d", qBank.SchemaVersion)
	}
	if len(qBank.Questions) != 19 {
		t.Fatalf("expected 19 canonical question templates (18 active, 1 retired), got %d", len(qBank.Questions))
	}

	activeCount := 0
	retiredCount := 0
	for _, q := range qBank.Questions {
		if q.Status == bank.StatusActive || q.Status == bank.StatusApprovedActive {
			activeCount++
			if q.Review.Reviewer == nil || *q.Review.Reviewer == "" {
				t.Errorf("active question %s missing reviewer", q.ID)
			}
			if q.Review.ApprovedAt == nil || *q.Review.ApprovedAt == "" {
				t.Errorf("active question %s missing approved_at", q.ID)
			}
		} else if q.Status == bank.StatusRetired {
			retiredCount++
		}
	}

	if activeCount != 18 {
		t.Errorf("expected 18 active templates, got %d", activeCount)
	}
	if retiredCount != 1 {
		t.Errorf("expected 1 retired template, got %d", retiredCount)
	}
}

func TestAccountsValidationRejections(t *testing.T) {
	t.Run("Reject unknown JSON field", func(t *testing.T) {
		jsonStr := `{
			"schema_version": 1,
			"course_scope": "provisional",
			"unexpected_field": "surprise",
			"accounts": [
				{"id": "cash", "name": "Cash", "category": "asset", "normal_side": "debit"}
			]
		}`
		_, _, err := bank.LoadAccounts(strings.NewReader(jsonStr))
		if err == nil {
			t.Fatalf("expected error for unknown JSON field")
		}
	})

	t.Run("Reject bad account ID with uppercase or spaces", func(t *testing.T) {
		jsonStr := `{
			"schema_version": 1,
			"course_scope": "provisional",
			"accounts": [
				{"id": "Bad Cash Account", "name": "Cash", "category": "asset", "normal_side": "debit"}
			]
		}`
		_, _, err := bank.LoadAccounts(strings.NewReader(jsonStr))
		if err == nil {
			t.Fatalf("expected error for bad account ID")
		}
	})

	t.Run("Reject duplicate account ID", func(t *testing.T) {
		jsonStr := `{
			"schema_version": 1,
			"course_scope": "provisional",
			"accounts": [
				{"id": "cash", "name": "Cash", "category": "asset", "normal_side": "debit"},
				{"id": "cash", "name": "Cash Duplicate", "category": "asset", "normal_side": "debit"}
			]
		}`
		_, _, err := bank.LoadAccounts(strings.NewReader(jsonStr))
		if err == nil {
			t.Fatalf("expected error for duplicate account ID")
		}
	})

	t.Run("Reject invalid normal side for category", func(t *testing.T) {
		jsonStr := `{
			"schema_version": 1,
			"course_scope": "provisional",
			"accounts": [
				{"id": "cash", "name": "Cash", "category": "asset", "normal_side": "credit"}
			]
		}`
		_, _, err := bank.LoadAccounts(strings.NewReader(jsonStr))
		if err == nil {
			t.Fatalf("expected error for mismatched normal side")
		}
	})

	t.Run("Reject unsupported schema version", func(t *testing.T) {
		jsonStr := `{
			"schema_version": 2,
			"course_scope": "provisional",
			"accounts": [
				{"id": "cash", "name": "Cash", "category": "asset", "normal_side": "debit"}
			]
		}`
		_, _, err := bank.LoadAccounts(strings.NewReader(jsonStr))
		if err == nil {
			t.Fatalf("expected error for schema version 2")
		}
	})
}

func TestQuestionBankValidationRejections(t *testing.T) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}

	validQuestionJSON := `{
		"id": "cash_service_test",
		"version": 1,
		"family_id": "cash_service",
		"rule_version": 1,
		"status": "seed_pending_review",
		"scenario_template": "The company completes services today and receives ${amount_dollars} cash.",
		"parameters": {
			"amount_minor_units": [5000, 10000]
		},
		"concepts": ["cash_vs_revenue", "debit_credit_translation"],
		"expected_fixture": {
			"postings": [
				{"account_id": "cash", "side": "debit", "amount_parameter": "amount_minor_units"},
				{"account_id": "service_revenue", "side": "credit", "amount_parameter": "amount_minor_units"}
			]
		},
		"review": {
			"reviewer": null,
			"approved_at": null,
			"source": "unit test"
		}
	}`

	wrapBank := func(q string) string {
		return `{"schema_version": 1, "questions": [` + q + `]}`
	}

	t.Run("Reject unknown JSON field in question", func(t *testing.T) {
		badJSON := strings.Replace(validQuestionJSON, `"source": "unit test"`, `"source": "unit test", "unknown_extra": 123`, 1)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(badJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for unknown field in question JSON")
		}
	})

	t.Run("Reject unsupported family ID", func(t *testing.T) {
		badJSON := strings.Replace(validQuestionJSON, `"family_id": "cash_service"`, `"family_id": "unsupported_derivatives"`, 1)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(badJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for unsupported family ID")
		}
	})

	t.Run("Reject duplicate question ID", func(t *testing.T) {
		bankJSON := `{"schema_version": 1, "questions": [` + validQuestionJSON + `, ` + validQuestionJSON + `]}`
		_, err := bank.LoadQuestionBank(strings.NewReader(bankJSON), catalog)
		if err == nil {
			t.Fatalf("expected rejection for duplicate question ID")
		}
	})

	t.Run("Reject unbalanced fixture (debits != credits)", func(t *testing.T) {
		unbalancedJSON := strings.Replace(
			validQuestionJSON,
			`"side": "credit"`,
			`"side": "debit"`, // Two debits, zero credits!
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(unbalancedJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for unbalanced fixture")
		}
	})

	t.Run("Reject fixture referencing unknown account", func(t *testing.T) {
		unknownAccJSON := strings.Replace(
			validQuestionJSON,
			`"account_id": "service_revenue"`,
			`"account_id": "nonexistent_cryptocurrency_holding"`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(unknownAccJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for fixture referencing unknown account")
		}
	})

	t.Run("Reject invalid/negative parameter amount", func(t *testing.T) {
		negativeAmtJSON := strings.Replace(
			validQuestionJSON,
			`"amount_minor_units": [5000, 10000]`,
			`"amount_minor_units": [5000, -100]`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(negativeAmtJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for negative parameter amount")
		}
	})

	t.Run("Reject missing required parameter", func(t *testing.T) {
		missingParamJSON := strings.Replace(
			validQuestionJSON,
			`"amount_minor_units": [5000, 10000]`,
			`"other_param": [5000]`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(missingParamJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for missing required parameter")
		}
	})

	t.Run("Reject invalid status", func(t *testing.T) {
		badStatusJSON := strings.Replace(
			validQuestionJSON,
			`"status": "seed_pending_review"`,
			`"status": "invalid_status_value"`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(badStatusJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for invalid question status")
		}
	})

	t.Run("Reject active question without reviewer", func(t *testing.T) {
		activeNoReviewerJSON := strings.Replace(
			validQuestionJSON,
			`"status": "seed_pending_review"`,
			`"status": "active"`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(activeNoReviewerJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for active question missing reviewer")
		}
	})

	t.Run("Reject active question without approved_at timestamp", func(t *testing.T) {
		activeNoApprovedAtJSON := strings.Replace(
			validQuestionJSON,
			`"reviewer": null`,
			`"reviewer": "prof_dan"`,
			1,
		)
		activeNoApprovedAtJSON = strings.Replace(
			activeNoApprovedAtJSON,
			`"status": "seed_pending_review"`,
			`"status": "active"`,
			1,
		)
		_, err := bank.LoadQuestionBank(strings.NewReader(wrapBank(activeNoApprovedAtJSON)), catalog)
		if err == nil {
			t.Fatalf("expected rejection for active question missing approved_at")
		}
	})
}

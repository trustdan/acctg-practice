package bank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/trustdan/acctg-practice/internal/domain"
)

var validIDRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

// LoadAccounts parses and strictly validates an accounts JSON reader into a domain.AccountCatalog.
func LoadAccounts(r io.Reader) (*domain.AccountCatalog, *AccountsFile, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var file AccountsFile
	if err := dec.Decode(&file); err != nil {
		return nil, nil, fmt.Errorf("failed to decode accounts JSON (check for syntax or unknown fields): %w", err)
	}

	if file.SchemaVersion != 1 {
		return nil, nil, fmt.Errorf("unsupported accounts schema_version: %d (expected 1)", file.SchemaVersion)
	}
	if strings.TrimSpace(file.CourseScope) == "" {
		return nil, nil, fmt.Errorf("course_scope cannot be empty")
	}
	if len(file.Accounts) == 0 {
		return nil, nil, fmt.Errorf("accounts array cannot be empty")
	}

	catalog := domain.NewAccountCatalog()
	seenIDs := make(map[string]struct{})

	for i, accJSON := range file.Accounts {
		id := strings.TrimSpace(accJSON.ID)
		if id == "" {
			return nil, nil, fmt.Errorf("account[%d] ID cannot be empty", i)
		}
		if !validIDRegex.MatchString(id) {
			return nil, nil, fmt.Errorf("account[%d] ID %q must be lowercase alphanumeric snake_case", i, id)
		}
		if _, exists := seenIDs[id]; exists {
			return nil, nil, fmt.Errorf("duplicate account ID: %q", id)
		}
		seenIDs[id] = struct{}{}

		name := strings.TrimSpace(accJSON.Name)
		if name == "" {
			return nil, nil, fmt.Errorf("account %q name cannot be empty", id)
		}

		cat := domain.Category(accJSON.Category)
		if !cat.IsValid() {
			return nil, nil, fmt.Errorf("account %q has invalid category: %q", id, accJSON.Category)
		}

		side := domain.Side(accJSON.NormalSide)
		if !side.IsValid() {
			return nil, nil, fmt.Errorf("account %q has invalid normal_side: %q", id, accJSON.NormalSide)
		}

		account := domain.Account{
			ID:         domain.AccountID(id),
			Name:       name,
			Category:   cat,
			NormalSide: side,
			ContraOf:   domain.AccountID(strings.TrimSpace(accJSON.ContraOf)),
		}

		if err := catalog.Add(account); err != nil {
			return nil, nil, fmt.Errorf("account %q validation error: %w", id, err)
		}
	}

	// Validate contra_of references
	for _, acc := range catalog.All() {
		if acc.ContraOf != "" {
			if !catalog.Has(acc.ContraOf) {
				return nil, nil, fmt.Errorf("contra account %q references unknown account %q", acc.ID, acc.ContraOf)
			}
		}
	}

	return catalog, &file, nil
}

// LoadAccountsFile opens and validates an accounts JSON file on disk.
func LoadAccountsFile(path string) (*domain.AccountCatalog, *AccountsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read accounts file %s: %w", path, err)
	}
	return LoadAccounts(bytes.NewReader(data))
}

// LoadQuestionBank parses and strictly validates a question bank JSON reader against an AccountCatalog.
func LoadQuestionBank(r io.Reader, catalog *domain.AccountCatalog) (*QuestionBankFile, error) {
	if catalog == nil {
		return nil, fmt.Errorf("account catalog must be provided to validate question bank")
	}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var bank QuestionBankFile
	if err := dec.Decode(&bank); err != nil {
		return nil, fmt.Errorf("failed to decode question bank JSON (check for syntax or unknown fields): %w", err)
	}

	if bank.SchemaVersion != 1 {
		return nil, fmt.Errorf("unsupported question bank schema_version: %d (expected 1)", bank.SchemaVersion)
	}
	if len(bank.Questions) == 0 {
		return nil, fmt.Errorf("question bank cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, q := range bank.Questions {
		id := strings.TrimSpace(q.ID)
		if id == "" {
			return nil, fmt.Errorf("question[%d] ID cannot be empty", i)
		}
		if !validIDRegex.MatchString(id) {
			return nil, fmt.Errorf("question[%d] ID %q must be lowercase alphanumeric snake_case", i, id)
		}
		if _, exists := seenIDs[id]; exists {
			return nil, fmt.Errorf("duplicate question ID: %q", id)
		}
		seenIDs[id] = struct{}{}

		if q.Version < 1 {
			return nil, fmt.Errorf("question %q version must be >= 1, got %d", id, q.Version)
		}
		if q.RuleVersion < 1 {
			return nil, fmt.Errorf("question %q rule_version must be >= 1, got %d", id, q.RuleVersion)
		}

		if !IsValidStatus(q.Status) {
			return nil, fmt.Errorf("question %q has invalid status: %q", id, q.Status)
		}

		if !IsSupportedFamily(q.FamilyID) {
			return nil, fmt.Errorf("question %q has unsupported family_id: %q", id, q.FamilyID)
		}

		if strings.TrimSpace(q.ScenarioTemplate) == "" {
			return nil, fmt.Errorf("question %q scenario_template cannot be empty", id)
		}

		// Validate parameters
		if len(q.Parameters) == 0 {
			return nil, fmt.Errorf("question %q parameters cannot be empty", id)
		}
		reqParameters := RequiredParameters(q.FamilyID)
		for _, reqParam := range reqParameters {
			vals, ok := q.Parameters[reqParam]
			if !ok || len(vals) == 0 {
				return nil, fmt.Errorf("question %q missing required parameter %q", id, reqParam)
			}
			for pIdx, v := range vals {
				if v <= 0 {
					return nil, fmt.Errorf("question %q parameter %q[%d] must be positive integer minor units, got %d", id, reqParam, pIdx, v)
				}
			}
		}

		// Validate concepts
		if len(q.Concepts) == 0 {
			return nil, fmt.Errorf("question %q concepts cannot be empty", id)
		}
		for cIdx, concept := range q.Concepts {
			if strings.TrimSpace(concept) == "" {
				return nil, fmt.Errorf("question %q concept[%d] cannot be empty", id, cIdx)
			}
		}

		// Validate expected fixture postings
		if len(q.ExpectedFixture.Postings) < 2 {
			return nil, fmt.Errorf("question %q fixture must contain at least 2 postings for double entry, got %d",
				id, len(q.ExpectedFixture.Postings))
		}

		for pIdx, p := range q.ExpectedFixture.Postings {
			accID := domain.AccountID(p.AccountID)
			if !catalog.Has(accID) {
				return nil, fmt.Errorf("question %q fixture posting[%d] references unknown account: %q", id, pIdx, p.AccountID)
			}
			side := domain.Side(p.Side)
			if !side.IsValid() {
				return nil, fmt.Errorf("question %q fixture posting[%d] has invalid side: %q", id, pIdx, p.Side)
			}
			paramVals, ok := q.Parameters[p.AmountParameter]
			if !ok || len(paramVals) == 0 {
				return nil, fmt.Errorf("question %q fixture posting[%d] references undefined amount_parameter: %q", id, pIdx, p.AmountParameter)
			}
		}

		// Validate fixture balancing for each sample parameter value
		sampleAmounts := q.Parameters["amount_minor_units"]
		for _, amt := range sampleAmounts {
			var totalDebit, totalCredit int64
			for _, p := range q.ExpectedFixture.Postings {
				if p.Side == string(domain.SideDebit) {
					totalDebit += amt
				} else if p.Side == string(domain.SideCredit) {
					totalCredit += amt
				}
			}
			if totalDebit <= 0 || totalDebit != totalCredit {
				return nil, fmt.Errorf("question %q expected fixture is unbalanced for amount %d: debits=%d, credits=%d",
					id, amt, totalDebit, totalCredit)
			}
		}

		// Validate review provenance
		if strings.TrimSpace(q.Review.Source) == "" {
			return nil, fmt.Errorf("question %q review.source cannot be empty", id)
		}
		if RequiresReviewProvenance(q.Status) {
			if q.Review.Reviewer == nil || strings.TrimSpace(*q.Review.Reviewer) == "" {
				return nil, fmt.Errorf("active question %q requires review.reviewer", id)
			}
			if q.Review.ApprovedAt == nil || strings.TrimSpace(*q.Review.ApprovedAt) == "" {
				return nil, fmt.Errorf("active question %q requires review.approved_at", id)
			}
		}
	}

	return &bank, nil
}

// LoadQuestionBankFile opens and strictly validates a question bank JSON file on disk.
func LoadQuestionBankFile(path string, catalog *domain.AccountCatalog) (*QuestionBankFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read question bank file %s: %w", path, err)
	}
	return LoadQuestionBank(bytes.NewReader(data), catalog)
}

package bank_test

import (
	"testing"

	"github.com/trustdan/acctg-practice/internal/bank"
)

func TestAllCourseModules(t *testing.T) {
	modules := bank.AllCourseModules()
	if len(modules) != 10 {
		t.Fatalf("expected 10 course modules, got %d", len(modules))
	}
	for i, m := range modules {
		if m.Number != i+1 {
			t.Errorf("module[%d] has Number %d, expected %d", i, m.Number, i+1)
		}
		if m.Title == "" || m.Description == "" {
			t.Errorf("module[%d] has empty title or description", i)
		}
	}
}

func TestQuestionModules(t *testing.T) {
	// Explicit modules
	q1 := bank.QuestionJSON{
		ID:       "q1",
		FamilyID: bank.FamilyCashService,
		Modules:  []int{1, 4},
	}
	m1 := bank.QuestionModules(q1)
	if len(m1) != 2 || m1[0] != 1 || m1[1] != 4 {
		t.Errorf("expected [1, 4], got %v", m1)
	}

	// Explicit single module
	q2 := bank.QuestionJSON{
		ID:       "q2",
		FamilyID: bank.FamilyCashService,
		Module:   5,
	}
	m2 := bank.QuestionModules(q2)
	if len(m2) != 1 || m2[0] != 5 {
		t.Errorf("expected [5], got %v", m2)
	}

	// Inferred from customer_advance -> [1, 2, 5]
	q3 := bank.QuestionJSON{
		ID:       "q3",
		FamilyID: bank.FamilyCustomerAdvance,
	}
	m3 := bank.QuestionModules(q3)
	if !bank.MatchesModules(q3, map[int]bool{5: true}) {
		t.Errorf("expected customer_advance to match module 5, got %v", m3)
	}
	if bank.MatchesModules(q3, map[int]bool{8: true}) {
		t.Errorf("did not expect customer_advance to match module 8")
	}

	// Inferred from borrow_cash -> [1, 2, 7]
	q4 := bank.QuestionJSON{
		ID:       "q4",
		FamilyID: bank.FamilyBorrowCash,
	}
	if !bank.MatchesModules(q4, map[int]bool{7: true}) {
		t.Errorf("expected borrow_cash to match module 7")
	}

	// Filter
	all := []bank.QuestionJSON{q1, q2, q3, q4}
	filtered := bank.FilterQuestionsByModules(all, []int{7})
	if len(filtered) != 1 || filtered[0].ID != "q4" {
		t.Errorf("expected only q4 in filtered list for module 7, got %v", filtered)
	}

	// Filter with empty modules returns all
	allBack := bank.FilterQuestionsByModules(all, nil)
	if len(allBack) != len(all) {
		t.Errorf("expected all questions when selectedModules is nil")
	}
}

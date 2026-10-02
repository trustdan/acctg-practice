package bank

import (
	"testing"
)

func TestSemanticWordingReviewCustomerAdvance(t *testing.T) {
	// Valid advance wording
	valid := "Acme Consulting received a ${amount_dollars} cash retainer in advance from a client for web design work to be performed next month."
	res := CheckSemanticWording(FamilyCustomerAdvance, valid)
	if !res.Approved || len(res.Violations) > 0 {
		t.Fatalf("expected valid advance to pass semantic review, got violations: %v", res.Violations)
	}

	// Balanced-but-wrong: describes services completed today, but labeled as customer_advance
	contradictory := "Client paid ${amount_dollars} in cash today for website consulting services that were fully delivered and completed today."
	res = CheckSemanticWording(FamilyCustomerAdvance, contradictory)
	if res.Approved {
		t.Fatalf("expected contradictory wording (completed today) to fail semantic review")
	}
	foundConflict := false
	for _, v := range res.Violations {
		if len(v) > 0 {
			foundConflict = true
		}
	}
	if !foundConflict {
		t.Fatalf("expected semantic conflict violation, got: %v", res.Violations)
	}

	// Missing future timing indicator
	noTiming := "Acme received ${amount_dollars} from a client."
	res = CheckSemanticWording(FamilyCustomerAdvance, noTiming)
	if res.Approved {
		t.Fatalf("expected wording without advance/future indicator to fail semantic review")
	}

	// Missing placeholder
	missingPlaceholder := "Acme received $1,000 cash advance for future services."
	res = CheckSemanticWording(FamilyCustomerAdvance, missingPlaceholder)
	if res.Approved {
		t.Fatalf("expected missing placeholder to fail semantic review")
	}
}

func TestSemanticWordingReviewOtherFamilies(t *testing.T) {
	// Prepaid purchase: claiming immediate expense fails
	contradictoryPrepaid := "Acme paid ${amount_dollars} cash for rent expense incurred last month and used up immediately."
	res := CheckSemanticWording(FamilyPrepaidPurchase, contradictoryPrepaid)
	if res.Approved {
		t.Fatalf("expected contradictory prepaid wording to fail")
	}

	// Equipment purchase: claiming operating expense fails
	contradictoryEquip := "Acme purchased a delivery truck for ${amount_dollars} cash and recorded it as an operating expense for the month."
	res = CheckSemanticWording(FamilyEquipmentPurchaseCash, contradictoryEquip)
	if res.Approved {
		t.Fatalf("expected contradictory equipment wording to fail")
	}

	// Borrowing cash: claiming earned revenue fails
	contradictoryBorrow := "Acme borrowed ${amount_dollars} from the bank and recorded it as earned revenue."
	res = CheckSemanticWording(FamilyBorrowCash, contradictoryBorrow)
	if res.Approved {
		t.Fatalf("expected contradictory borrowing wording to fail")
	}

	// Dividend cash: claiming wage expense fails
	contradictoryDividend := "Acme paid ${amount_dollars} to its owners and recorded it as salary expense."
	res = CheckSemanticWording(FamilyDividendCash, contradictoryDividend)
	if res.Approved {
		t.Fatalf("expected contradictory dividend wording to fail")
	}

	// Receivable collection: claiming new revenue earned fails
	contradictoryCollect := "Acme collected ${amount_dollars} from a customer on account and recognized revenue today."
	res = CheckSemanticWording(FamilyCollectReceivable, contradictoryCollect)
	if res.Approved {
		t.Fatalf("expected contradictory collection wording to fail")
	}
}

func TestFamilyReviewedAndTested(t *testing.T) {
	// All initial families are tested at rule_version 1
	if !IsFamilyReviewedAndTested(FamilyCustomerAdvance, 1) {
		t.Errorf("expected customer_advance v1 to be reviewed and tested")
	}
	if !IsFamilyReviewedAndTested(FamilyPrepaidPurchase, 1) {
		t.Errorf("expected prepaid_purchase v1 to be reviewed and tested")
	}

	// New / unknown family is NOT reviewed and tested
	if IsFamilyReviewedAndTested("crypto_arbitrage", 1) {
		t.Errorf("crypto_arbitrage should not be reviewed and tested")
	}

	// Unreviewed rule version is NOT reviewed and tested
	if IsFamilyReviewedAndTested(FamilyCustomerAdvance, 2) {
		t.Errorf("customer_advance v2 should not be reviewed and tested yet")
	}
}

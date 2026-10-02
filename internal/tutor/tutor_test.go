package tutor

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestOfflineTutorHintAndExplain(t *testing.T) {
	tutor := NewOfflineTutor()

	req := Request{
		ProblemPrompt: "On January 15, Acme Inc. pays $1,200 cash for 12 months of future insurance coverage.",
		FamilyID:      "prepaid_purchase",
		Stage:         domain.StageIdentifyAccount,
		StagePrompt:   "Which primary account is affected first?",
		CausalHint:    "Consider whether paying cash in advance acquires a future right/asset or creates an immediate expense.",
		Explanation:   "Prepaid Insurance is an Asset that increases via Debit. Cash decreases via Credit.",
		ErrorTag:      engine.TagExpenseRecordedOnPrepaidPurchase,
		ConceptID:     "prepaid_expenses",
	}

	ctx := context.Background()

	// Hint
	hintResp, err := tutor.Hint(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error from Hint: %v", err)
	}
	if hintResp.Text == "" {
		t.Errorf("expected non-empty hint text")
	}
	if hintResp.Provider != "offline" {
		t.Errorf("expected provider 'offline', got %q", hintResp.Provider)
	}
	if hintResp.Fallback {
		t.Errorf("expected Fallback=false for OfflineTutor")
	}
	// Check that specific distractor hint was incorporated
	if !strings.Contains(hintResp.Text, "Prepaid Expense") {
		t.Errorf("expected hint to mention Prepaid Expense, got: %s", hintResp.Text)
	}

	// Explain
	explainResp, err := tutor.Explain(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error from Explain: %v", err)
	}
	if explainResp.Text == "" {
		t.Errorf("expected non-empty explanation text")
	}
	if !strings.Contains(explainResp.Text, "Assets = Liabilities + Equity") {
		t.Errorf("expected explanation to mention fundamental equation, got: %s", explainResp.Text)
	}
	if !strings.Contains(explainResp.Text, "Debit means Left and Credit means Right") {
		t.Errorf("expected explanation to define debit/credit correctly, got: %s", explainResp.Text)
	}
	if !strings.Contains(explainResp.Text, "Prepaid Expenses") {
		t.Errorf("expected contrast analysis for prepaid_purchase, got: %s", explainResp.Text)
	}
}

func TestOfflineTutorHonorsContextCancellation(t *testing.T) {
	tutor := NewOfflineTutor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req := Request{
		ProblemPrompt: "Test prompt",
		Stage:         domain.StageIdentifyAccount,
	}

	_, err := tutor.Hint(ctx, req)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}

	_, err = tutor.Explain(ctx, req)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

func TestSimulatedProviderSuccess(t *testing.T) {
	budget := NewBudget(5, 1000)
	sim := NewSimulatedProvider(SimulatedConfig{
		Delay:      10 * time.Millisecond,
		Budget:     budget,
		CustomHint: "Custom simulated hint for testing.",
	})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Test prompt"}

	resp, err := sim.Hint(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "Custom simulated hint for testing." {
		t.Errorf("unexpected text: %s", resp.Text)
	}
	if resp.Provider != "simulated" {
		t.Errorf("expected provider 'simulated', got %q", resp.Provider)
	}
	if budget.RemainingRequests() != 4 {
		t.Errorf("expected 4 remaining requests, got %d", budget.RemainingRequests())
	}
}

func TestSimulatedProviderCancellationDuringDelay(t *testing.T) {
	sim := NewSimulatedProvider(SimulatedConfig{
		Delay: 200 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 20ms
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := sim.Hint(ctx, Request{ProblemPrompt: "Test"})
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
	if elapsed >= 180*time.Millisecond {
		t.Fatalf("did not cancel early, elapsed: %v", elapsed)
	}
}

func TestSimulatedProviderTimeout(t *testing.T) {
	sim := NewSimulatedProvider(SimulatedConfig{
		Delay: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := sim.Hint(ctx, Request{ProblemPrompt: "Test"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestSimulatedProviderErrorInjection(t *testing.T) {
	injectedErr := errors.New("simulated network outage")
	sim := NewSimulatedProvider(SimulatedConfig{
		FailWith: injectedErr,
	})

	_, err := sim.Hint(context.Background(), Request{ProblemPrompt: "Test"})
	if !errors.Is(err, injectedErr) {
		t.Fatalf("expected injected error, got: %v", err)
	}
}

func TestSimulatedProviderBudgetExceeded(t *testing.T) {
	budget := NewBudget(1, 1000)
	sim := NewSimulatedProvider(SimulatedConfig{
		Budget:     budget,
		CustomHint: "First response",
	})

	// First request succeeds
	_, err := sim.Hint(context.Background(), Request{ProblemPrompt: "Test 1"})
	if err != nil {
		t.Fatalf("first request should succeed, got: %v", err)
	}

	// Second request exceeds budget
	_, err = sim.Hint(context.Background(), Request{ProblemPrompt: "Test 2"})
	if !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("expected ErrBudgetExceeded, got: %v", err)
	}
}

func TestSimulatedProviderResponseSizeLimit(t *testing.T) {
	sim := NewSimulatedProvider(SimulatedConfig{
		MaxResponseBytes: 15,
		CustomHint:       "This response is way too long for the 15-byte limit.",
	})

	_, err := sim.Hint(context.Background(), Request{ProblemPrompt: "Test"})
	if !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got: %v", err)
	}
}

func TestFallbackTutorSuccessPath(t *testing.T) {
	primary := NewSimulatedProvider(SimulatedConfig{
		CustomHint: "Primary provider response",
	})
	fallback := NewOfflineTutor()
	wrapper := NewFallbackTutor(primary, fallback, 1*time.Second)

	resp, err := wrapper.Hint(context.Background(), Request{ProblemPrompt: "Test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "Primary provider response" {
		t.Errorf("expected primary text, got: %s", resp.Text)
	}
	if resp.Fallback {
		t.Errorf("expected Fallback=false when primary succeeds")
	}
}

func TestFallbackTutorGracefulFallbackOnTimeout(t *testing.T) {
	// Primary takes 200ms, but timeout is set to 20ms
	primary := NewSimulatedProvider(SimulatedConfig{
		Delay: 200 * time.Millisecond,
	})
	fallback := NewOfflineTutor()
	wrapper := NewFallbackTutor(primary, fallback, 20*time.Millisecond)

	req := Request{
		ProblemPrompt: "Acme collects $500 cash from customer.",
		Stage:         domain.StageIdentifyAccount,
		CausalHint:    "Cash increases on debit.",
	}

	resp, err := wrapper.Hint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error during fallback: %v", err)
	}
	if !resp.Fallback {
		t.Errorf("expected Fallback=true when primary times out")
	}
	if resp.Provider != "offline" {
		t.Errorf("expected provider 'offline' from fallback, got %q", resp.Provider)
	}
	if resp.Text == "" {
		t.Errorf("expected non-empty fallback text")
	}
}

func TestFallbackTutorGracefulFallbackOnError(t *testing.T) {
	primary := NewSimulatedProvider(SimulatedConfig{
		FailWith: errors.New("connection reset by peer"),
	})
	fallback := NewOfflineTutor()
	wrapper := NewFallbackTutor(primary, fallback, 1*time.Second)

	req := Request{
		ProblemPrompt: "Acme pays rent.",
		Stage:         domain.StageIdentifyAccount,
		CausalHint:    "Rent Expense is recorded.",
	}

	resp, err := wrapper.Hint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error during fallback: %v", err)
	}
	if !resp.Fallback {
		t.Errorf("expected Fallback=true on primary error")
	}
	if resp.Provider != "offline" {
		t.Errorf("expected provider 'offline' from fallback, got %q", resp.Provider)
	}
}

func TestFallbackTutorRespectsParentContextCancellation(t *testing.T) {
	primary := NewSimulatedProvider(SimulatedConfig{
		Delay: 500 * time.Millisecond,
	})
	fallback := NewOfflineTutor()
	wrapper := NewFallbackTutor(primary, fallback, 1*time.Second)

	parentCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel parent immediately

	_, err := wrapper.Hint(parentCtx, Request{ProblemPrompt: "Test"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled when parent is cancelled, got: %v", err)
	}
}

func TestPedagogicalInvariants(t *testing.T) {
	offline := NewOfflineTutor()
	ctx := context.Background()

	// Check all major families for pedagogical safety
	families := []string{
		"collect_receivable",
		"customer_advance",
		"earn_advance",
		"prepaid_purchase",
		"prepaid_consumption",
		"equipment_purchase_cash",
		"repay_note_principal",
		"dividend_cash",
	}

	for _, fam := range families {
		req := Request{
			ProblemPrompt: "Scenario test for family " + fam,
			FamilyID:      fam,
			Stage:         domain.StageIdentifyAccount,
			StagePrompt:   "Select account",
		}

		resp, err := offline.Explain(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", fam, err)
		}

		text := strings.ToLower(resp.Text)

		// 1. Debit cannot be treated as "money out"
		if strings.Contains(text, "debit is money out") || strings.Contains(text, "debit means money out") {
			t.Errorf("pedagogical violation in %s: debit described as money out", fam)
		}

		// 2. Must emphasize debit=left, credit=right
		if !strings.Contains(text, "debit means left") || !strings.Contains(text, "credit means right") {
			t.Errorf("pedagogical violation in %s: missing debit=left / credit=right convention", fam)
		}

		// 3. Must not shame learner or use medical/diagnostic terms
		forbiddenTerms := []string{"lazy", "dumb", "stupid", "diagnosis", "pathology", "disease", "deficit", "disorder"}
		for _, term := range forbiddenTerms {
			if strings.Contains(text, term) {
				t.Errorf("pedagogical violation in %s: found forbidden term %q", fam, term)
			}
		}
	}
}

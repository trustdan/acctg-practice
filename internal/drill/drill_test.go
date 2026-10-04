package drill_test

import (
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func setupTestDrill(t *testing.T) (*drill.Generator, *bank.QuestionBankFile, *domain.AccountCatalog) {
	catalog, _, err := bank.LoadAccountsFile("../../curriculum/accounts.json")
	if err != nil {
		t.Fatalf("failed to load accounts: %v", err)
	}
	eng := engine.NewEngine(catalog)
	gen := drill.NewGenerator(catalog, eng)

	qBank, err := bank.LoadQuestionBankFile("../../curriculum/seed-questions.json", catalog)
	if err != nil {
		t.Fatalf("failed to load questions: %v", err)
	}

	return gen, qBank, catalog
}

func getQuestionByID(qBank *bank.QuestionBankFile, id string) *bank.QuestionJSON {
	for i := range qBank.Questions {
		if qBank.Questions[i].ID == id {
			return &qBank.Questions[i]
		}
	}
	return nil
}

func TestCustomerAdvanceDrillFullRun(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)

	q := getQuestionByID(qBank, "customer_advance_basic")
	if q == nil {
		t.Fatalf("expected customer_advance_basic in question bank")
	}

	seed := int64(12345)
	params := map[string]int64{"amount_minor_units": 10000} // $100

	inst, err := gen.GenerateInstance(*q, seed, params)
	if err != nil {
		t.Fatalf("unexpected GenerateInstance error: %v", err)
	}

	// Verify scenario template expansion
	if !strings.Contains(inst.PromptText, "$100") {
		t.Fatalf("expected prompt text to contain $100, got: %s", inst.PromptText)
	}

	session := drill.NewSession("sess-1", inst)
	now := time.Now().UTC()

	// Stage 1: Identify Account -> Select Cash
	st1, err := session.CurrentStage()
	if err != nil || st1.Stage != domain.StageIdentifyAccount {
		t.Fatalf("expected StageIdentifyAccount, got %v", st1)
	}
	fb1, err := session.SubmitOption("opt_cash", now)
	if err != nil || !fb1.IsCorrect || !fb1.AdvanceStage {
		t.Fatalf("unexpected Stage 1 response: fb=%+v, err=%v", fb1, err)
	}
	if fb1.AssistanceLevel != domain.AssistanceNone {
		t.Errorf("expected AssistanceNone for unassisted first try, got %s", fb1.AssistanceLevel)
	}

	// Stage 2: Account Category -> Select Asset
	st2, err := session.CurrentStage()
	if err != nil || st2.Stage != domain.StageAccountCategory {
		t.Fatalf("expected StageAccountCategory, got %v", st2)
	}
	fb2, err := session.SubmitOption("opt_asset", now)
	if err != nil || !fb2.IsCorrect || !fb2.AdvanceStage {
		t.Fatalf("unexpected Stage 2 response: fb=%+v, err=%v", fb2, err)
	}

	// Stage 3: Direction -> Select Increase
	st3, err := session.CurrentStage()
	if err != nil || st3.Stage != domain.StageDirection {
		t.Fatalf("expected StageDirection, got %v", st3)
	}
	fb3, err := session.SubmitOption("opt_increase", now)
	if err != nil || !fb3.IsCorrect || !fb3.AdvanceStage {
		t.Fatalf("unexpected Stage 3 response: fb=%+v, err=%v", fb3, err)
	}

	// Stage 4: Debit or Credit -> Select Debit
	st4, err := session.CurrentStage()
	if err != nil || st4.Stage != domain.StageDebitCredit {
		t.Fatalf("expected StageDebitCredit, got %v", st4)
	}
	fb4, err := session.SubmitOption("opt_debit", now)
	if err != nil || !fb4.IsCorrect || !fb4.AdvanceStage {
		t.Fatalf("unexpected Stage 4 response: fb=%+v, err=%v", fb4, err)
	}

	// Stage 5: Counter-Account -> TEST TARGETED HINT FOR REVENUE CONFUSION!
	st5, err := session.CurrentStage()
	if err != nil || st5.Stage != domain.StageCounterAccount {
		t.Fatalf("expected StageCounterAccount, got %v", st5)
	}

	// Learner makes the classic mistake: selects Service Revenue instead of Unearned Revenue
	fb5Wrong, err := session.SubmitOption("opt_service_rev", now)
	if err != nil {
		t.Fatalf("unexpected error submitting distractor: %v", err)
	}
	if fb5Wrong.IsCorrect {
		t.Fatalf("expected Service Revenue on customer advance to be evaluated as incorrect")
	}
	if fb5Wrong.AdvanceStage {
		t.Fatalf("expected AdvanceStage to be false on first error to allow retry")
	}
	if fb5Wrong.ErrorTag != "revenue_recognized_prematurely" {
		t.Errorf("expected ErrorTag revenue_recognized_prematurely, got %s", fb5Wrong.ErrorTag)
	}
	if fb5Wrong.Hint == "" {
		t.Fatalf("expected targeted Socratic hint, got empty string")
	}
	if !strings.Contains(fb5Wrong.Hint, "owe the customer") || !strings.Contains(fb5Wrong.Hint, "?") {
		t.Errorf("expected hint to ask about the unfinished obligation, got: %s", fb5Wrong.Hint)
	}

	// Now learner retries with Unearned Revenue
	fb5Retry, err := session.SubmitOption("opt_unearned_rev", now)
	if err != nil || !fb5Retry.IsCorrect || !fb5Retry.AdvanceStage {
		t.Fatalf("unexpected Stage 5 retry response: fb=%+v, err=%v", fb5Retry, err)
	}
	if fb5Retry.AssistanceLevel != domain.AssistanceRetry {
		t.Errorf("expected AssistanceRetry for retry attempt, got %s", fb5Retry.AssistanceLevel)
	}

	// Stage 6: Balanced Entry Assembly -> Select Dr Cash / Cr Unearned Revenue
	st6, err := session.CurrentStage()
	if err != nil || st6.Stage != domain.StageBalancedEntry {
		t.Fatalf("expected StageBalancedEntry, got %v", st6)
	}
	fb6, err := session.SubmitOption("opt_dr_cash_cr_unearned", now)
	if err != nil || !fb6.IsCorrect || !fb6.AdvanceStage {
		t.Fatalf("unexpected Stage 6 response: fb=%+v, err=%v", fb6, err)
	}

	// Stage 7: Equation Effect -> Select Assets up, Liab up
	st7, err := session.CurrentStage()
	if err != nil || st7.Stage != domain.StageEquationEffect {
		t.Fatalf("expected StageEquationEffect, got %v", st7)
	}
	fb7, err := session.SubmitOption("opt_assets_up_liab_up", now)
	if err != nil || !fb7.IsCorrect || !fb7.AdvanceStage {
		t.Fatalf("unexpected Stage 7 response: fb=%+v, err=%v", fb7, err)
	}

	// Session is now complete
	if !session.IsCompleted {
		t.Fatalf("expected session to be completed")
	}

	// Verify transaction recap
	recap := session.Recap()
	if !recap.IsBalanced {
		t.Fatalf("expected recap entry to be balanced")
	}
	if len(recap.Postings) != 2 {
		t.Fatalf("expected 2 postings in recap, got %d", len(recap.Postings))
	}
	if recap.Postings[0].AccountID != "cash" || recap.Postings[0].Side != domain.SideDebit {
		t.Errorf("expected posting[0] to be Dr Cash, got %+v", recap.Postings[0])
	}
	if recap.Postings[1].AccountID != "unearned_revenue" || recap.Postings[1].Side != domain.SideCredit {
		t.Errorf("expected posting[1] to be Cr Unearned Revenue, got %+v", recap.Postings[1])
	}

	// Verify persisted attempts in session
	if len(session.Attempts) != 8 { // 7 stages + 1 retry on stage 5 = 8 attempts
		t.Fatalf("expected 8 recorded attempts, got %d", len(session.Attempts))
	}
}

func TestStableOptionIDsNotLabels(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "customer_advance_basic")
	params := map[string]int64{"amount_minor_units": 10000}

	// Create sessions with different seeds
	inst1, _ := gen.GenerateInstance(*q, 111, params)
	inst2, _ := gen.GenerateInstance(*q, 999, params)

	sess1 := drill.NewSession("s1", inst1)
	sess2 := drill.NewSession("s2", inst2)

	// Both must accept "opt_cash" as correct, regardless of what label ("a", "b", "c", "d") it was given
	now := time.Now().UTC()
	fb1, err := sess1.SubmitOption("opt_cash", now)
	if err != nil || !fb1.IsCorrect {
		t.Fatalf("sess1 failed to recognize opt_cash: %v", err)
	}

	fb2, err := sess2.SubmitOption("opt_cash", now)
	if err != nil || !fb2.IsCorrect {
		t.Fatalf("sess2 failed to recognize opt_cash: %v", err)
	}

	// Submitting a label ("a") must fail because label is display-only, not stable ID
	_, err = sess1.SubmitOption("a", now)
	if err == nil {
		t.Fatalf("expected error submitting display label 'a' as option ID")
	}
}

func TestSeedReproducibility(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "customer_advance_basic")
	params := map[string]int64{"amount_minor_units": 10000}

	// Same seed produces identical option order across all stages
	instA, _ := gen.GenerateInstance(*q, 42, params)
	instB, _ := gen.GenerateInstance(*q, 42, params)

	if instA.PromptText != instB.PromptText {
		t.Fatalf("prompt text mismatch with same seed")
	}

	for stageKey, stageA := range instA.StageAnswers {
		stageB, ok := instB.StageAnswers[stageKey]
		if !ok {
			t.Fatalf("missing stage %s in instB", stageKey)
		}
		if len(stageA.Options) != len(stageB.Options) {
			t.Fatalf("option count mismatch for stage %s", stageKey)
		}
		for i := range stageA.Options {
			if stageA.Options[i].ID != stageB.Options[i].ID {
				t.Fatalf("option order mismatch at index %d for stage %s: %s vs %s",
					i, stageKey, stageA.Options[i].ID, stageB.Options[i].ID)
			}
			if stageA.Options[i].Label != stageB.Options[i].Label {
				t.Fatalf("option label mismatch at index %d for stage %s", i, stageKey)
			}
		}
	}

	// Different seed produces different shuffle order
	instC, _ := gen.GenerateInstance(*q, 999999, params)
	st5A := instA.StageAnswers[domain.StageCounterAccount]
	st5C := instC.StageAnswers[domain.StageCounterAccount]

	differentFound := false
	for i := range st5A.Options {
		if st5A.Options[i].ID != st5C.Options[i].ID {
			differentFound = true
			break
		}
	}
	if !differentFound {
		t.Fatalf("expected different seeds to shuffle options differently")
	}
}

func TestExplicitHintRequest(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "customer_advance_basic")
	params := map[string]int64{"amount_minor_units": 10000}

	inst, _ := gen.GenerateInstance(*q, 42, params)
	sess := drill.NewSession("s-hint", inst)

	hint, err := sess.RequestHint()
	if err != nil {
		t.Fatalf("unexpected RequestHint error: %v", err)
	}
	if hint == "" {
		t.Fatalf("expected non-empty hint")
	}
	if sess.CurrentAssistance != domain.AssistanceHinted {
		t.Errorf("expected CurrentAssistance to be AssistanceHinted, got %s", sess.CurrentAssistance)
	}

	// Learner answers correctly after hint
	fb, err := sess.SubmitOption("opt_cash", time.Now().UTC())
	if err != nil || !fb.IsCorrect {
		t.Fatalf("unexpected submit result: %v", err)
	}
	if fb.AssistanceLevel != domain.AssistanceHinted {
		t.Errorf("expected recorded assistance to be AssistanceHinted, got %s", fb.AssistanceLevel)
	}
}

func TestExhaustedRetriesAdvancesWithRevealedAnswer(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "customer_advance_basic")
	params := map[string]int64{"amount_minor_units": 10000}

	inst, _ := gen.GenerateInstance(*q, 42, params)
	sess := drill.NewSession("s-exhaust", inst)
	now := time.Now().UTC()

	// Error 1: allows retry
	fb1, _ := sess.SubmitOption("opt_notes_payable", now)
	if fb1.IsCorrect || fb1.AdvanceStage {
		t.Fatalf("expected first error to pause for retry")
	}

	// Error 2: retry exhausted -> advances stage with revealed answer
	fb2, _ := sess.SubmitOption("opt_ar", now)
	if fb2.IsCorrect {
		t.Fatalf("expected second attempt to be incorrect")
	}
	if !fb2.AdvanceStage {
		t.Fatalf("expected AdvanceStage to be true when retries are exhausted")
	}
	if !strings.Contains(fb2.Explanation, "was: Cash.") {
		t.Errorf("expected revealed answer to show the correct option text, got: %s", fb2.Explanation)
	}

	// Verify next stage is now active
	st, _ := sess.CurrentStage()
	if st.Stage != domain.StageAccountCategory {
		t.Errorf("expected StageAccountCategory after revealed answer, got %s", st.Stage)
	}
}

func TestAllActiveQuestionsGenerateAllSevenStages(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)

	requiredStages := []domain.DrillStage{
		domain.StageIdentifyAccount,
		domain.StageAccountCategory,
		domain.StageDirection,
		domain.StageDebitCredit,
		domain.StageCounterAccount,
		domain.StageBalancedEntry,
		domain.StageEquationEffect,
	}

	for _, q := range qBank.Questions {
		if q.Status != bank.StatusActive && q.Status != bank.StatusApprovedActive {
			continue // skip retired or draft
		}

		amt := q.Parameters["amount_minor_units"][0]
		params := map[string]int64{"amount_minor_units": amt}
		inst, err := gen.GenerateInstance(q, 42, params)
		if err != nil {
			t.Fatalf("failed to generate instance for %s: %v", q.ID, err)
		}

		if len(inst.StageAnswers) < 7 {
			t.Fatalf("question %s only generated %d stages, expected 7", q.ID, len(inst.StageAnswers))
		}

		for _, st := range requiredStages {
			ans, ok := inst.StageAnswers[st]
			if !ok {
				t.Errorf("question %s missing stage %s", q.ID, st)
				continue
			}
			if ans.Prompt == "" {
				t.Errorf("question %s stage %s has empty prompt", q.ID, st)
			}
			if ans.CorrectOptionID == "" {
				t.Errorf("question %s stage %s has empty CorrectOptionID", q.ID, st)
			}
			if len(ans.Options) < 2 {
				t.Errorf("question %s stage %s has fewer than 2 options", q.ID, st)
			}
			if ans.CausalHint == "" {
				t.Errorf("question %s stage %s has empty CausalHint", q.ID, st)
			}
			// Verify CorrectOptionID exists among options
			foundCorrect := false
			for _, opt := range ans.Options {
				if opt.ID == ans.CorrectOptionID {
					foundCorrect = true
					break
				}
			}
			if !foundCorrect {
				t.Errorf("question %s stage %s: CorrectOptionID %q not found in options", q.ID, st, ans.CorrectOptionID)
			}
		}
	}
}

func TestPrepaidContrastPairDrill(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)

	// 1. Prepaid Purchase: asset swap (Cash down, Prepaid Insurance up; no expense)
	qPur := getQuestionByID(qBank, "prepaid_insurance_retail")
	if qPur == nil {
		t.Fatalf("expected prepaid_insurance_retail in question bank")
	}
	instPur, err := gen.GenerateInstance(*qPur, 42, map[string]int64{"amount_minor_units": 120000})
	if err != nil {
		t.Fatalf("failed to generate prepaid purchase: %v", err)
	}
	st1Pur := instPur.StageAnswers[domain.StageIdentifyAccount]
	if st1Pur.CorrectOptionID != "opt_prepaid_insurance" {
		t.Errorf("expected opt_prepaid_insurance for purchase, got %s", st1Pur.CorrectOptionID)
	}
	st7Pur := instPur.StageAnswers[domain.StageEquationEffect]
	if st7Pur.CorrectOptionID != "opt_asset_swap" {
		t.Errorf("expected opt_asset_swap for purchase recap, got %s", st7Pur.CorrectOptionID)
	}

	// 2. Prepaid Consumption: expense recognition (Prepaid Insurance down, Insurance Expense up; zero cash)
	qCons := getQuestionByID(qBank, "prepaid_consumption_retail")
	if qCons == nil {
		t.Fatalf("expected prepaid_consumption_retail in question bank")
	}
	instCons, err := gen.GenerateInstance(*qCons, 42, map[string]int64{"amount_minor_units": 10000})
	if err != nil {
		t.Fatalf("failed to generate prepaid consumption: %v", err)
	}
	st1Cons := instCons.StageAnswers[domain.StageIdentifyAccount]
	if st1Cons.CorrectOptionID != "opt_insurance_expense" {
		t.Errorf("expected opt_insurance_expense for consumption, got %s", st1Cons.CorrectOptionID)
	}
	st7Cons := instCons.StageAnswers[domain.StageEquationEffect]
	if st7Cons.CorrectOptionID != "opt_assets_down_eq_down" {
		t.Errorf("expected opt_assets_down_eq_down for consumption recap, got %s", st7Cons.CorrectOptionID)
	}
}

func TestScaffoldFadingLevels(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "prepaid_insurance_retail")
	if q == nil {
		t.Fatalf("expected prepaid_insurance_retail in question bank")
	}

	// 1. Full scaffolding (Level 0) -> 7 stages
	instFull, err := gen.GenerateInstanceWithScaffold(*q, 100, map[string]int64{"amount_minor_units": 120000}, domain.ScaffoldFull)
	if err != nil {
		t.Fatalf("failed generating full instance: %v", err)
	}
	sessFull := drill.NewSession("sess-full", instFull)
	if len(sessFull.StageSequence) != 7 {
		t.Fatalf("expected 7 stages for ScaffoldFull, got %d", len(sessFull.StageSequence))
	}
	expectedFull := []domain.DrillStage{
		domain.StageIdentifyAccount,
		domain.StageAccountCategory,
		domain.StageDirection,
		domain.StageDebitCredit,
		domain.StageCounterAccount,
		domain.StageBalancedEntry,
		domain.StageEquationEffect,
	}
	for i, st := range expectedFull {
		if sessFull.StageSequence[i] != st {
			t.Errorf("stage %d: expected %s, got %s", i, st, sessFull.StageSequence[i])
		}
	}

	// 2. Intermediate scaffolding (Level 1) -> 4 stages
	instInterm, err := gen.GenerateInstanceWithScaffold(*q, 101, map[string]int64{"amount_minor_units": 120000}, domain.ScaffoldIntermediate)
	if err != nil {
		t.Fatalf("failed generating intermediate instance: %v", err)
	}
	sessInterm := drill.NewSession("sess-interm", instInterm)
	if len(sessInterm.StageSequence) != 4 {
		t.Fatalf("expected 4 stages for ScaffoldIntermediate, got %d", len(sessInterm.StageSequence))
	}
	expectedInterm := []domain.DrillStage{
		domain.StageIdentifyAccount,
		domain.StageCounterAccount,
		domain.StageBalancedEntry,
		domain.StageEquationEffect,
	}
	for i, st := range expectedInterm {
		if sessInterm.StageSequence[i] != st {
			t.Errorf("intermediate stage %d: expected %s, got %s", i, st, sessInterm.StageSequence[i])
		}
	}

	// 3. Faded scaffolding (Level 2) -> 2 stages
	instFaded, err := gen.GenerateInstanceWithScaffold(*q, 102, map[string]int64{"amount_minor_units": 120000}, domain.ScaffoldFaded)
	if err != nil {
		t.Fatalf("failed generating faded instance: %v", err)
	}
	sessFaded := drill.NewSession("sess-faded", instFaded)
	if len(sessFaded.StageSequence) != 2 {
		t.Fatalf("expected 2 stages for ScaffoldFaded, got %d", len(sessFaded.StageSequence))
	}
	expectedFaded := []domain.DrillStage{
		domain.StageBalancedEntry,
		domain.StageEquationEffect,
	}
	for i, st := range expectedFaded {
		if sessFaded.StageSequence[i] != st {
			t.Errorf("faded stage %d: expected %s, got %s", i, st, sessFaded.StageSequence[i])
		}
	}
}

func TestReferenceConsultationTracking(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "prepaid_insurance_retail")
	if q == nil {
		t.Fatalf("expected question")
	}

	inst, err := gen.GenerateInstance(*q, 200, map[string]int64{"amount_minor_units": 120000})
	if err != nil {
		t.Fatalf("failed generating instance: %v", err)
	}
	sess := drill.NewSession("sess-ref", inst)

	// Stage 1: Learner opens cheatsheet reference before answering
	sess.RecordReferenceUse()
	st1, err := sess.CurrentStage()
	if err != nil {
		t.Fatalf("failed getting current stage: %v", err)
	}

	fb1, err := sess.SubmitOption(st1.CorrectOptionID, time.Now())
	if err != nil {
		t.Fatalf("failed submitting option: %v", err)
	}
	if !fb1.IsCorrect {
		t.Fatalf("expected correct")
	}
	if fb1.AssistanceLevel != domain.AssistanceReference {
		t.Fatalf("expected feedback assistance to be %s, got %s", domain.AssistanceReference, fb1.AssistanceLevel)
	}

	// Check persisted attempt on session
	if len(sess.Attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(sess.Attempts))
	}
	att1 := sess.Attempts[0]
	if !att1.ReferenceUsed {
		t.Fatalf("expected ReferenceUsed to be true")
	}
	if att1.Assistance != domain.AssistanceReference {
		t.Fatalf("expected Assistance to be %s, got %s", domain.AssistanceReference, att1.Assistance)
	}

	// Stage 2: Learner answers without reference
	st2, err := sess.CurrentStage()
	if err != nil {
		t.Fatalf("failed getting stage 2: %v", err)
	}
	fb2, err := sess.SubmitOption(st2.CorrectOptionID, time.Now())
	if err != nil {
		t.Fatalf("failed submitting option: %v", err)
	}
	if fb2.AssistanceLevel != domain.AssistanceNone {
		t.Fatalf("expected assistance None for unassisted stage 2, got %s", fb2.AssistanceLevel)
	}
	att2 := sess.Attempts[1]
	if att2.ReferenceUsed {
		t.Fatalf("expected ReferenceUsed to be false for stage 2")
	}
	if att2.Assistance != domain.AssistanceNone {
		t.Fatalf("expected Assistance to be None, got %s", att2.Assistance)
	}
}

func TestTeachingOverridesFollowReviewProvenance(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	base := *getQuestionByID(qBank, "cash_service_basic")
	base.Teaching = map[domain.DrillStage]bank.TeachingText{
		domain.StageCounterAccount: {Hint: "Reviewed hint ${amount_dollars}", Explanation: "Reviewed explanation."},
	}
	params := map[string]int64{"amount_minor_units": 5000}

	cases := []struct {
		status string
		want   bool
	}{
		{bank.StatusActive, true},
		{bank.StatusApprovedActive, true},
		{bank.StatusSeedPendingReview, false},
	}
	for _, tc := range cases {
		q := base
		q.Status = tc.status
		inst, err := gen.GenerateInstance(q, 1, params)
		if err != nil {
			t.Fatalf("%s: %v", tc.status, err)
		}
		got := inst.StageAnswers[domain.StageCounterAccount]
		applied := got.CausalHint == "Reviewed hint $50" && got.Explanation == "Reviewed explanation."
		if applied != tc.want {
			t.Errorf("%s: override applied=%v, want %v (hint %q)", tc.status, applied, tc.want, got.CausalHint)
		}
		if got.CorrectOptionID != "opt_service_rev" {
			t.Errorf("%s: canonical answer changed to %q", tc.status, got.CorrectOptionID)
		}
	}
}

func TestRetryFailureShowsCorrectOptionText(t *testing.T) {
	gen, qBank, _ := setupTestDrill(t)
	q := getQuestionByID(qBank, "customer_advance_basic")
	inst, err := gen.GenerateInstance(*q, 7, map[string]int64{"amount_minor_units": 5000})
	if err != nil {
		t.Fatal(err)
	}
	session := drill.NewSession("sess-retry", inst)
	now := time.Now()
	if _, err := session.SubmitOption("opt_service_rev", now); err != nil {
		t.Fatal(err)
	}
	fb, err := session.SubmitOption("opt_service_rev", now)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fb.Explanation, "The correct answer was: Cash.") || strings.Contains(fb.Explanation, "opt_") {
		t.Fatalf("retry feedback should show option text, got %q", fb.Explanation)
	}
}

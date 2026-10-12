# Stream-M04 Accounting & Pedagogical Review

**Reviewer:** Worker Agent 2 (Stream-M04)  
**Date:** 2026-10-12  
**Target Syllabus Module:** Module 4 (Revenue Recognition / ASC 606)  
**Assigned Scope:** 4 candidate questions (`q_m04_rev_*`) covering Variable Consideration Constraint, Contract Asset vs. Accounts Receivable presentation, Right of Return Accounting / Refund Liability, and Performance Obligation Satisfaction Over Time.

---

## 1. Accounting Standards & Syllabus Alignment

### A. Variable Consideration & The Constraint (ASC 606 Step 3)
- **Question:** `q_m04_rev_variable_consideration_01`
- **Technical Standard:** ASC 606-10-32-11 through 32-14 states that an entity shall include in the transaction price some or all of an amount of variable consideration only to the extent that it is probable that a significant reversal in the amount of cumulative revenue recognized will not occur when the uncertainty associated with the variable consideration is subsequently resolved.
- **Scenario Design:** The scenario features a fixed consulting fee ($1,200 / $2,400 / $4,800) combined with an uncertain incentive bonus of $2,000 conditioned on client-side operating cost savings outside the firm's direct influence.
- **Accounting Mechanics:** The variable bonus fails the constraint threshold and must be excluded from the transaction price. Revenue is recognized strictly for the unconstrained fixed fee on account:
  - Debit `accounts_receivable`
  - Credit `service_revenue`
- **Misconception Tag:** `TagM04VariableConsiderationUnconstrained` captures learners who aggressively recognize variable bonuses before meeting the ASC 606 constraint.

### B. Contract Asset vs. Accounts Receivable Presentation (ASC 606 Balance Sheet Presentation)
- **Question:** `q_m04_rev_contract_asset_ar_01`
- **Technical Standard:** ASC 606-10-45-3 through 45-5 draws a fundamental distinction:
  - An **account receivable** is an entity's right to consideration that is *unconditional*—that is, only the passage of time is required before payment of that consideration is due.
  - A **contract asset** is an entity's right to consideration in exchange for goods or services that the entity has transferred to a customer when that right is conditioned on something other than the passage of time (such as future performance of other deliverables).
- **Scenario Design:** The firm completes a distinct engineering milestone that has received client sign-off. Under the contractual terms, the right to bill is unconditional on net 30 payment terms; only the passage of time is required before payment is due.
- **Accounting Mechanics:** Because the right is unconditional, the debit must be presented as an Account Receivable rather than a Contract Asset:
  - Debit `accounts_receivable`
  - Credit `service_revenue`
- **Misconception Tag:** `TagM04UnconditionalRightMisclassifiedAsContractAsset` diagnoses learners who fail to recognize that unconditional claims must be classified as accounts receivable.

### C. Right of Return Accounting & Refund Liabilities (ASC 606-10-55-23)
- **Question:** `q_m04_rev_right_of_return_01`
- **Technical Standard:** Under ASC 606-10-55-23, when an entity receives consideration from a customer and expects to refund some or all of the consideration, it recognizes a refund liability. If an entity receives upfront payment for services not yet performed subject to a 100% money-back guarantee where returns cannot be reliably estimated, no revenue can be recognized; the entire cash receipt must be deferred as a contract liability (Unearned Revenue / Refund Obligation).
- **Scenario Design:** A customer pays cash in advance for an executive seminar series with a 30-day full money-back guarantee, prior to any performance.
- **Accounting Mechanics:**
  - Debit `cash`
  - Credit `unearned_revenue`
- **Misconception Tag:** `TagM04RevenueRecognizedWithUncertainReturn` catches premature revenue recognition when significant refund/return rights remain unexpired and performance has not occurred.

### D. Performance Obligation Satisfied Over Time (ASC 606 Step 5)
- **Question:** `q_m04_rev_satisfaction_over_time_01`
- **Technical Standard:** Under ASC 606-10-25-27, an entity transfers control of a good or service over time and satisfies a performance obligation over time if the customer simultaneously receives and consumes the benefits provided by the entity's performance. Fulfilling the obligation relieves the contract liability (Unearned Revenue) created upon cash collection.
- **Scenario Design:** The client had previously paid in advance (held in Unearned Revenue). Over the current month, ongoing compliance services are rendered and consumed over time, satisfying the performance obligation proportionally.
- **Accounting Mechanics:**
  - Debit `unearned_revenue`
  - Credit `service_revenue`
- **Misconception Tag:** `TagM04CashDebitedOnOverTimeSatisfaction` catches learners who mistakenly debit Cash again upon service delivery, failing to understand contract liability relief.

---

## 2. Invariant & Pedagogical Verification Matrix

| Invariant | Requirement | Status | Verification Detail |
|---|---|:---:|---|
| **Integer Minor Units** | All financial values stored in integer cents | **PASS** | Amounts encoded as 120,000 to 800,000 cents ($1,200 to $8,000). Zero float rounding. |
| **Double-Entry Balance** | $\sum \text{Debits} == \sum \text{Credits}$ for all amounts | **PASS** | Validated across all 3 parameter permutations for all 4 templates in unit test. |
| **Equation Invariance** | $\Delta \text{Assets} == \Delta \text{Liabilities} + \Delta \text{Equity}$ | **PASS** | Derived statement effects tie out identically for all 4 questions. |
| **Socratic Hint Invariant** | Hints guide causally and end with `?` | **PASS** | All 28 stage hints across 4 templates end with `?`. No direct answer disclosure. |
| **Stage 1 Neutrality** | Stage 1 explanation omits "debit" & "credit" | **PASS** | Automated check confirmed 0 occurrences of "debit" or "credit" in Stage 1 explanations. |
| **Parameter Variety** | $\ge 3$ distinct amount parameters per question | **PASS** | Each template provides exactly 3 distinct amount choices for randomized replay. |
| **Account Catalog** | All postings reference existing canonical catalog | **PASS** | Postings use `accounts_receivable`, `service_revenue`, `cash`, `unearned_revenue`. |

---

## 3. Reviewer Sign-Off

All 4 question templates in `curriculum/candidates/batch_m04.json`, their pedagogy policy entries in `curriculum/candidates/pedagogy_m04.json`, their engine rules in `internal/engine/rules_m04.go`, and their stage builders in `internal/drill/stages_m04.go` have been verified under automated test harness `internal/drill/syllabus_m04_test.go`.

The candidate package is isolated and ready for ingestion by the Integrator Agent during the Wave 1 reducer phase.

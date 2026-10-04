# Content audit

Current status (October 3, 2026): the focused quality work at 90 is complete under pedagogy policy 2. [Quality closeout](QUALITY-REVIEW90.md) records ten amount profiles, contextual account choices and the remaining external validation gates. Earlier sections retain the findings and next steps at their original checkpoints.

Step 1 of [CONTENT-EXPANSION-PLAN.md](CONTENT-EXPANSION-PLAN.md). Findings only: no content, status, or code was changed. Recommendations need human semantic approval before any wording is versioned.

Sources inspected: curriculum/seed-questions.json, curriculum/accounts.json, internal/drill/generator.go, internal/drill/session.go, internal/tutor/offline.go, internal/exam/errors.go, internal/engine/event.go. No syllabus or slides available.

Teaching (prompts, options, hints, explanations) is defined per **family** in generator.go, with optional reviewed per-scenario teaching overrides. Findings below describe the audited baseline; later implementation records supersede resolved findings.

Stage numbering: S1 identify account, S2 category, S3 direction, S4 debit/credit, S5 counter-account, S6 balanced entry, S7 equation effect.

## Group 1 — revenue-side families (9 scenarios, 5 families) — 2026-10-02

### Scenario matrix

| Scenario | Family | Canonical entry | Timing cues in wording | Difficulty | Contrast partner | Provenance |
|---|---|---|---|---|---|---|
| cash_service_basic | cash_service | Dr Cash / Cr Service Revenue | Performed today, paid today, "no invoice previously issued" | Basic | customer_advance_basic (performance timing) | course_staff, 2026-10-01T12:00Z, unverified |
| cash_service_webdev | cash_service | Dr Cash / Cr Service Revenue | Delivered today, paid on delivery, no prior invoices **or advances** | Basic | none in same context | same |
| customer_advance_basic | customer_advance | Dr Cash / Cr Unearned Revenue | Paid today, performed next month, none performed | Basic | cash_service_basic; earn_advance_basic | same |
| customer_advance_logistics | customer_advance | Dr Cash / Cr Unearned Revenue | Paid today, transport next month, none occurred | Basic | earn_advance_logistics (same business) | same |
| service_on_credit_basic | service_on_credit | Dr A/R / Cr Service Revenue | Completed today, invoiced, not paid | Basic | collect_receivable_basic | same |
| service_credit_consulting | service_on_credit | Dr A/R / Cr Service Revenue | Completed today, invoiced net 30, no cash today | Basic | collect_receivable_consulting (same business) | same |
| collect_receivable_basic | collect_receivable | Dr Cash / Cr A/R | Paid today; revenue and receivable recorded last month | Basic | service_on_credit_basic | same |
| collect_receivable_consulting | collect_receivable | Dr Cash / Cr A/R | Paid today; revenue recognized last month | Basic | service_credit_consulting | same |
| earn_advance_basic | earn_advance | Dr Unearned Revenue / Cr Service Revenue | Performed today; paid last month; no cash today | Basic | customer_advance_basic | same |
| earn_advance_logistics | earn_advance | Dr Unearned Revenue / Cr Service Revenue | Delivered today; paid last month; no cash today | Basic | customer_advance_logistics | same |

Misconceptions targeted: premature revenue on advance (customer_advance), duplicate revenue on collection (collect_receivable), cash on earning an advance (earn_advance), cash-before-collection (service_on_credit), advance-vs-earned (cash_service). All 10 are single-event, two-line, basic difficulty; within a family only business setting and amount vary.

Provenance: all 19 records (including retired) carry the identical reviewer `course_staff` and timestamp `2026-10-01T12:00:00Z`. Per the plan, this does not establish that a semantic review occurred. Treat as unverified until the user confirms.

### Scenario wording findings

| ID | Finding | Severity | Suggested direction |
|---|---|---|---|
| W1 | cash_service_basic rules out a prior invoice but not a prior advance. The webdev variant states both. | Medium | Add "and no advance payment was received earlier." |
| W2 | earn_advance_basic and earn_advance_logistics name "Unearned Revenue" in the scenario. That gives away the S1 answer and, through S1's prompt, S2-S4 as well. | Medium | Describe the prior event economically ("the company recorded an obligation to perform"), or keep it deliberately as a scaffold and record that decision. Reviewer's call. |
| W3 | Every scenario is basic. No family includes a partial-amount, decrease-led, or mixed-timing cue. | Low (expansion gap) | Address in pilot batches, not by editing these records. |

### Family teaching findings (apply to both scenarios in each family)

**Answer-revealing hints.** The hint shown after a first error states or nearly states the answer, so the "assisted retry" tests reading rather than reasoning:

| Family | Stages whose hint reveals the answer |
|---|---|
| customer_advance | S5 ("creates an obligation (liability)"; only one liability option), S6 (names the full entry), S7 |
| cash_service | S4 ("increases go on the left side"), S6 (names the entry) |
| service_on_credit | S4, S6, S7 |
| collect_receivable | S2 ("Cash is an asset."), S4, S6, S7 |
| earn_advance | S4, S6, S7; also the S1 **prompt** says "Which liability account is fulfilled and debited?", which leaks S2, S3, and S4 |

The S4 and S6 hints reveal the answer in all five families.

**Explanations that only repeat the answer.** S2-S4 and S6 explanations in every family restate the answer ("Assets increase by Debit." / "Debit Cash $X and Credit ... $X.") and do not say why the tempting alternative fails. collect_receivable S3 is just "Cash increases."

**Wording that conflicts with the course's core distinction.**

| ID | Location | Finding | Severity |
|---|---|---|---|
| T1 | service_on_credit S5 prompt | "What account earns this inflow under accrual accounting?" No cash flows in a credit sale, so "inflow" contradicts the scenario and blurs cash versus revenue. | High |
| T2 | cash_service S5 prompt | "what counter-account earns this inflow?" frames revenue as earning cash. | High |
| T3 | service_on_credit S1 and earn_advance S1 explanations | Say "is debited" before the learner reaches the S4 debit/credit decision. | Low |
| T4 | Category labels | S2 offers "Revenue" and "Equity" as separate categories (matching accounts.json `revenue`). S5 option text says "Service Revenue (Revenue / Equity)", and S7 says "Equity increases (+Service Revenue)". These are defensible but inconsistent; confirm the course convention. | Low |

**Distractors and error tags.**

| ID | Finding | Severity |
|---|---|---|
| D1 | **Fixed 2026-10-02 (tag normalization).** Eleven option tags are not canonical engine tags, so no offline hint or exam explanation routes from them: `treated_advance_as_revenue`, `confused_advance_with_receivable`, `confused_with_payable`, `confused_with_owner_investment`, `service_on_credit`, `treated_service_as_unearned`, `confused_with_borrowing`, `cash_recorded_before_collection`, `treated_revenue_as_liability`, `confused_with_asset_swap`, `treated_advance_as_equity`. | Medium |
| D2 | **Fixed 2026-10-02.** customer_advance S7: "Assets up / Equity up (+Revenue)" is tagged `treated_advance_as_equity`; the misconception is premature revenue (`revenue_recognized_prematurely`). | Medium |
| D3 | earn_advance S7: "No effect on any balance; memo entry only" is untagged and implausible. A stronger distractor is "Assets up / Liabilities down" (records new cash). | Low |
| D4 | collect_receivable S7: "Assets up / Liabilities up" is untagged. It is the customer-advance confusion and should be tagged as such. | Low |
| D5 | Several S1/S5 distractors have no tag (cash_service S1 A/R, Unearned, A/P; service_on_credit S1/S5 A/P, Unearned; collect_receivable S5 Unearned, A/P; earn_advance S1/S5 A/R). | Low |
| D6 | cash_service S6 omits the most instructive contrasts: Dr A/R / Cr Revenue (credit-sale confusion) and Dr Cash / Cr A/R (collection confusion). | Low |

### Delivery findings (plan step 2; not content edits)

| ID | Location | Finding | Severity |
|---|---|---|---|
| X1 | internal/drill/session.go:209, shown by internal/tui/view.go:336 | After a failed retry the learner sees the raw option ID: "The correct answer was: opt_unearned_rev." It should show the option text. **Fixed 2026-10-02.** | High |
| X2 | internal/drill/session.go:198-203 | The first-error hint is the stage's single `CausalHint` whichever distractor was chosen. The `ErrorTag` only changes the "Hint:" prefix. Misconception-specific hints exist only in the tutor path (`[t]`/offline tutor), and only for the 9 canonical tags. | Medium |
| X3 | internal/drill/generator.go:111 | Seed `teaching` overrides apply only to `approved_active`; all shipped seeds are `active`, so per-scenario teaching cannot be delivered without a status change. | Medium (known) |
| X4 | internal/tutor/offline.go:107-111 | Every offline explanation appends the same four-bullet "Core Accounting Principles" lecture, which the teaching standard asks to avoid. | Low |

Seen in passing, for group 2: the offline contrast text says dividends are recorded by "directly debiting Retained Earnings" (offline.go:164), and the dividend error hint says the dividend reduces Retained Earnings, but the canonical entry is Dr Dividends / Cr Cash. The prepaid contrasts say "Prepaid Expenses" where the canonical account is Prepaid Insurance.

### Reviewer decisions (user, 2026-10-02)

1. **W2 — No.** Scenarios must not name the account the learner is asked to identify. Describe prior events economically (for example, "recorded an obligation to perform the work"). Revise both earn_advance scenarios as new versions.
2. **T4 — Revenue is its own category.** This matches accounts.json (`revenue`). Teaching must label Service Revenue "Revenue", not "Revenue / Equity". The equation stage may still say that revenue increases equity, because that describes the effect, not the category. Expenses and dividends follow the same pattern: each is its own category, with its effect on equity stated separately.
3. **Provenance — Unknown.** Treat the existing `course_staff` / `2026-10-01T12:00:00Z` records as unverified. Do not alter or delete them, since that would rewrite history. When the user reviews a revised scenario in the terminal, record the new version with the user as reviewer and the actual review date. Unrevised records keep their original provenance and remain flagged here as unverified.
4. **X3 — Resolved in code.** Reviewed teaching overrides now apply to every status that requires review provenance (`active`, `approved_active`) through `bank.RequiresReviewProvenance`, which the validator and generator share. `seed_pending_review` is practice-eligible without a reviewer, so it never gets override wording. One rule now decides both whether a record must be reviewed and whether its teaching is used. Adding or changing `teaching` on a seed must bump `version` so historical instances keep their snapshots.

## Group 2 — remaining families (8 scenarios, 8 families) — 2026-10-02

### Scenario matrix

| Scenario | Family | Canonical entry | Timing cues in wording | Contrast partner | Provenance |
|---|---|---|---|---|---|
| cash_rent_basic | cash_rent | Dr Rent Expense / Cr Cash | Paid today for the current month | prepaid_insurance_retail (current vs future benefit); dividend_cash_retail | course_staff, unverified |
| borrow_cash_basic | borrow_cash | Dr Cash / Cr Notes Payable | Borrowed today, note due in two years, ignore interest | issue_shares_basic, cash_service_basic (source of cash); repay_principal_logistics | same |
| issue_shares_basic | issue_shares | Dr Cash / Cr Common Stock | Issued today for cash from investors | borrow_cash_basic (debt vs equity) | same |
| prepaid_insurance_retail | prepaid_purchase | Dr Prepaid Insurance / Cr Cash | Paid today for 12 months, none expired yet | prepaid_consumption_retail (same business); cash_rent_basic | same |
| prepaid_consumption_retail | prepaid_consumption | Dr Insurance Expense / Cr Prepaid Insurance | Month-end expiry, no cash today | prepaid_insurance_retail | same |
| equipment_cash_rental | equipment_purchase_cash | Dr Equipment / Cr Cash | Paid today, 5-year life | cash_rent_basic (capitalize vs expense) | same |
| repay_principal_logistics | repay_note_principal | Dr Notes Payable / Cr Cash | Paid today, principal only, ignore interest | borrow_cash_basic | same |
| dividend_cash_retail | dividend_cash | Dr Dividends / Cr Cash | "Pays cash dividends today", with no declaration timing | cash_rent_basic (distribution vs expense) | same |

All eight scenarios are basic, one per family. Only the prepaid pair shares a business context. The prepaid amounts are consistent: the consumption amounts ($100/$200/$300) are one-twelfth of the purchase amounts ($1,200/$2,400/$3,600).

### Scenario wording findings

| ID | Finding | Severity | Suggested direction |
|---|---|---|---|
| W4 | dividend_cash_retail does not say whether the dividend was declared earlier. If it was, the correct entry is Dr Dividends Payable / Cr Cash (the engine supports `pay_dividend_payable` and tags `dividends_debited_on_payment`). The wording supports two answers. | **High** | New version: "declares and pays ... today; no dividend was previously declared." |
| W5 | cash_rent_basic does not rule out rent previously recorded as owed. | Low | Add "No rent was previously recorded as owed." |
| W6 | The retail boutique appears in the prepaid pair and the dividend scenario ("retail corporation"), and logistics appears in group 1 and repay_principal. That is acceptable, but cross-family reuse of a business does not count as transfer breadth. | Info | Track in the expansion matrix. |

### Family teaching findings

**Prompts that name the answer** (the scenario-level rule from decision 1 applies to stage prompts too):

| Family | Prompt | Leaks |
|---|---|---|
| cash_rent S1 | "Which account records the **rent expense** incurred...?" | S1 answer and S2 category |
| prepaid_consumption S1 | "Which **expense** account recognizes...?" | S2 category |
| repay_note_principal S1 | "Which **liability** account is **reduced**...?" | S2 and S3 |
| borrow_cash S5 | "...the formal **obligation to repay** the lender?" | Points straight to Notes Payable (only liability option) |

**Answer-revealing hints:** borrow_cash and issue_shares S2 ("Cash is an asset."), S4 (hint identical to the explanation), and S6. In every group-2 family, the S4 and S6 hints state the answer and the S7 hint states the effect. prepaid_purchase S1 hint says "is an asset, not an immediate operating expense", which nearly gives the answer.

**Explanations that only repeat the answer:** same pattern as group 1 (S2-S4, S6). Several S1 explanations say "is debited" before S4 (rent, prepaid purchase, prepaid consumption, equipment, repay, dividends, borrow, shares).

**Category labels (decision 2):** S2 options read "Expense (Equity reduction)" and "Dividends (Equity reduction)", while group 1 shows "Revenue" alone. Normalize to bare category names ("Expense", "Dividends", "Revenue"), and leave the equity effect to the hint, explanation, and S7.

**Distractors and error tags:**

| ID | Finding | Severity |
|---|---|---|
| D7 | **Rent Expense** is used as the "expensed it" distractor for equipment (S1, S5, S6), loan principal (S1, S5, S6), and dividends (S1, S5, S6). A learner would not confuse buying machinery, repaying a loan, or paying dividends with rent, so the misconception being tested (treating it as *an* expense) is lost. The catalog has no generic, interest, or equipment expense account. Options: add catalog accounts (e.g. Interest Expense, Repairs/Maintenance Expense) through reviewed rules, or phrase the distractor as the misconception ("an expense account for the cost"). | **High** |
| D8 | **Fixed 2026-10-02.** Non-canonical tags in group 2: `treated_borrowing_as_revenue` (4 options) and `treated_equity_as_revenue` (3), duplicating unused engine tags `revenue_recorded_on_borrowing` and `revenue_recorded_on_share_issue`; `cash_recorded_on_prepaid_expiration` (4), with no engine equivalent; `confused_with_payable` (5); `confused_with_asset_swap` (2). | Medium |
| D9 | **Partly fixed 2026-10-02:** borrowing and share-issue tags are now attached; all five now have offline hints. `prepaid_not_expensed_on_consumption` still needs a distractor (D11), and the two dividend-declaration tags need families not in the bank. Originally: engine tags never attached to any drill option: `revenue_recorded_on_borrowing`, `revenue_recorded_on_share_issue`, `prepaid_not_expensed_on_consumption`, `cash_recorded_on_dividend_declaration`, `dividends_debited_on_payment`. Their exam explanations therefore never fire from drills. | Medium |
| D10 | Implausible or untagged filler: Service Revenue in cash_rent S5, prepaid S5, and repay S1/S5; "No change; expenses do not affect balance sheet" (rent S7); "No change in equity; dividends are an asset" (dividend S7); "No net change...asset swap only" (borrow, shares, repay S7). | Low |
| D11 | Useful contrasts that are missing: rent S6 lacks Dr Prepaid Rent / Cr Cash (prepaid confusion); prepaid_consumption lacks a "no entry needed" or "not expensed" option for `prepaid_not_expensed_on_consumption`; equipment S7 "Assets up / Liabilities up" (bought on account) is untagged; borrow S5 Common Stock is untagged (owner-investment confusion). | Low |
| D12 | cash_rent S1 offers "Prepaid Rent", which is not in curriculum/accounts.json. That is acceptable as a distractor but inconsistent with the canonical catalog; adding `prepaid_rent` would also enable a rent-prepaid contrast family. | Low |

### Tutor and exam terminology

| ID | Location | Finding | Severity |
|---|---|---|---|
| T5 | internal/tutor/offline.go:164 | Dividend contrast: "directly debiting Retained Earnings (Equity)". The canonical entry debits **Dividends**, which reduces equity and is closed to Retained Earnings at period end. **Fixed 2026-10-02.** | **High** |
| T6 | internal/tutor/offline.go:137 | Dividend error hint: "reducing Retained Earnings". Should name the Dividends account first. **Fixed 2026-10-02.** | Medium |
| T7 | internal/tutor/offline.go:131, 156, 158 | "Prepaid Expense(s)" and "Rent/Insurance Expense" where the canonical accounts are Prepaid Insurance and Insurance Expense. **Fixed 2026-10-02.** | Medium |
| T8 | internal/tutor/offline.go:122-145 | **Fixed 2026-10-02.** The offline tutor has no hints for the borrowing, share-issue, prepaid-consumption, or dividend-payable tags (D9). | Medium |

### Reviewer decisions (user delegated to implementer, 2026-10-02)

5. **D7 — option (b).** The syllabus is unavailable, so interest is not confirmed in scope; the catalog is unchanged. **Fixed 2026-10-02:** the Rent Expense distractors for equipment, loan principal, and dividends now state the misconception ("An expense for the cost of the machinery" / "Debit an expense for the machinery $X / Credit Cash $X", and likewise for the loan payment and the payment to stockholders), under option ID `opt_expense_misconception`. The meaningless "credit an expense" S5 options were replaced by Notes Payable (equipment) and Dividends Payable (dividends), both tagged `wrong_account`, or dropped (repay).
6. **W4 — revise.** Draft wording is in "Scenario wording batch 1" below, pending the user's approval of the exact text.

## Audit summary — all 18 active scenarios

| Severity | Content (wording, teaching, distractors) | Delivery / code |
|---|---|---|
| High | T1, T2 ("earns this inflow"); W4 (dividend timing); D7 (Rent Expense distractors); S4/S6 hints reveal answers in all 13 families | X1 (raw option ID); T5 (dividend contrast) |
| Medium | W1, W2; prompts naming answers (earn_advance S1, rent S1, consumption S1, repay S1, borrow S5); D1, D2, D8 | X2 (hint ignores distractor); D9, T6-T8 (tag routing and tutor wording); X3 (resolved) |
| Low | W3, W5; T3, T4 and category-label normalization; D3-D6, D10-D12 | X4 (repeated lecture) |

Proposed fix order:
1. ~~X1, then T5-T7 (tutor wording)~~ — done 2026-10-02; new tutor strings are listed in HANDOFF for review.
2. ~~Tag normalization (D1, D8, D9; adds T8 hints)~~ — done 2026-10-02, see "Tag normalization" below. map the non-canonical tags to engine constants, adding engine tags where a real misconception has none, so every tagged distractor routes to a hint and an exam explanation.
3. Teaching rewrite per family (hints, explanations, prompts, category labels, distractors), reviewed in the terminal in groups of 5-10 scenarios. Scenario wording changes (W1, W2, W4, W5) ship as new versions with real reviewer provenance.
4. X2 misconception-routed first-error hints, after tags are canonical.

## Tag normalization — 2026-10-02

Every drill distractor tag is now an engine constant, and `internal/exam/tag_routing_test.go` fails if any tagged option in an active seed lacks a specific offline hint or exam explanation.

New engine tags: `cash_recorded_when_uncollected` and `cash_recorded_on_prepaid_expiration` (the engine already emitted both as literals), `revenue_deferred_when_earned`, `advance_confused_with_receivable`, `payable_recorded_for_cash_payment`, `equation_effect_missed`.

| Old tag | Now |
|---|---|
| treated_advance_as_revenue, treated_advance_as_equity | revenue_recognized_prematurely |
| confused_advance_with_receivable, service_on_credit | advance_confused_with_receivable |
| treated_service_as_unearned, treated_revenue_as_liability | revenue_deferred_when_earned |
| cash_recorded_before_collection | cash_recorded_when_uncollected |
| treated_borrowing_as_revenue | revenue_recorded_on_borrowing |
| treated_equity_as_revenue | revenue_recorded_on_share_issue |
| confused_with_asset_swap | equation_effect_missed |
| confused_with_payable | payable_recorded_for_cash_payment (rent, prepaid); wrong_account (cash_service, customer_advance, where the customer paid) |
| confused_with_owner_investment, confused_with_borrowing | wrong_account |

Historical attempts keep the tag they were recorded with. Old tag strings still render through the exam's generic fallback.

## Scenario wording batch 1 — approved and published 2026-10-03

Five revised scenario templates, approved by the user without edits on 2026-10-03 and published in curriculum/seed-questions.json as `version: 2` of the same IDs (reviewer `trustdan`, approved_at `2026-10-03T19:02:01Z`). This resolves W1, W2, W4, and W5. Canonical entries, families, amounts, and concepts are unchanged. Version 1 instances and attempts keep their snapshotted prompts; the v1 `course_staff` review records remain in git history (the bank holds one version per ID).

| ID | Finding | Proposed v2 wording |
|---|---|---|
| cash_service_basic | W1 | The company completes services today and receives ${amount_dollars} cash today. No invoice was issued and no advance payment was received earlier. |
| earn_advance_basic | W2 | The company completes ${amount_dollars} of services today. The customer paid for this work last month, and the company recorded its obligation to perform the work at that time. No cash changes hands today. |
| earn_advance_logistics | W2 | The freight logistics carrier safely delivers cargo today, completing a shipping contract for which the shipper paid ${amount_dollars} last month. The carrier recorded its obligation to deliver the cargo at that time. No cash changes hands today. |
| cash_rent_basic | W5 | The company pays ${amount_dollars} cash today for rent for the current month. No rent was previously recorded as owed or paid in advance. |
| dividend_cash_retail | W4 | A retail corporation declares and pays ${amount_dollars} of cash dividends to its stockholders today. No dividend was declared earlier. |

## Teaching rewrite batch 1: revenue families — reviewed and applied by delegated review (2026-10-03)

Covers cash_service, customer_advance, service_on_credit, collect_receivable, and earn_advance (generator.go; shared by both scenarios in each family). Correct answers, option IDs, and grading are unchanged. Instances already stored keep their snapshotted text.

Rules applied:
- A hint is a question that points to the economic fact or rule needed, not to the answer. Hints are shown after a first wrong answer, whichever wrong option was picked.
- Hints must make sense on their own. The scaffold can skip stages, so the S1/S5 hints cannot rely on S2-S4, and the S6/S7 hints cannot rely on any earlier stage.
- An explanation states the answer, says why, and says why the most tempting wrong option fails. It is shown after a correct answer and after the answer is revealed.
- S1 text never says "debited" or "credited" (T3). Category labels are bare ("Revenue"), and the effect on equity is stated separately (decision 2). Prompts never name the answer (T1, T2, earn_advance S1).
- Debit means left, not "money in" or "money out".

`${amount}` is the instance amount.

### Shared stage text (Cash in every family; adapted for other assets)

| Stage | Field | Proposed |
|---|---|---|
| S2 Cash category | Hint | Is cash something the company owns and can use, something it owes, or the owners' stake? |
| S2 Cash category | Explanation | Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance). |
| S3 Cash increase | Hint | Compare the company's cash before and after today's payment. Is it higher or lower? |
| S3 Cash increase | Explanation | The company holds more money after the payment, so Cash increases. |
| S4 Asset increase | Hint | An account increases on the same side as its normal balance. Which side is an asset's normal balance? |
| S4 Asset increase | Explanation | Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category. |
| S7 (all) | Hint opening | Place each account in the entry under Assets, Liabilities, or Equity. (Each family then adds one question; see below.) |

### customer_advance (Dr Cash / Cr Unearned Revenue)

| Stage | Field | Current | Proposed |
|---|---|---|---|
| S1 | Prompt | Which account is directly affected by the cash payment received from the customer today? | What did the company receive today, and which account records it? |
| S1 | Hint | The company physically received money today. What asset account tracks cash inflows? | Set aside whether the work is done. What came into the business today? |
| S1 | Explanation | Cash is received immediately, so the Cash account is affected. | The customer paid cash today, so Cash is affected. Service Revenue is tempting, but receiving money is not the same as earning it; the work happens next month. |
| S2-S4 | | | Shared text above. |
| S5 | Option | Service Revenue (Revenue / Equity) | Service Revenue (Revenue) |
| S5 | Hint | The service has not been performed yet. Revenue cannot be recognized before it is earned. Receiving cash before performing work creates an obligation (liability) to the customer. | Has the company done the work yet? If not, what does it now owe the customer? |
| S5 | Explanation | Because the service has not been provided, the company has an obligation to perform in the future, recorded as Unearned Revenue (a liability). | The company has been paid but has not done the work, so it owes the customer the service (or a refund). That obligation is Unearned Revenue, a Liability. Service Revenue would record earning before any work is done. Accounts Receivable would mean the customer still owes money, but the customer has already paid. |
| S6 | Hint | You need a Debit to Cash (asset up) and a Credit to Unearned Revenue (liability up) for equal amounts. | The customer paid for work the company will do next month. What did the company gain today, and what does it now owe? |
| S6 | Explanation | Debit Cash ${amount} and Credit Unearned Revenue ${amount}. | Debit Cash ${amount} (asset up) and Credit Unearned Revenue ${amount} (liability up). Crediting Service Revenue instead would report revenue for work not yet done. |
| S7 | Hint | Cash is an Asset. Unearned Revenue is a Liability. Has any Equity/Revenue changed? | Shared opening, then: Has the company earned anything yet? |
| S7 | Explanation | Assets increase by ${amount} and Liabilities increase by ${amount}. Both sides of the equation remain in balance. | Assets increase by ${amount} (Cash) and Liabilities increase by ${amount} (Unearned Revenue). Equity is unchanged because nothing has been earned yet; it increases next month, when the work is done. |

### cash_service (Dr Cash / Cr Service Revenue)

| Stage | Field | Current | Proposed |
|---|---|---|---|
| S1 | Options | A/R, Unearned Revenue, A/P untagged | A/R `wrong_account`; Unearned Revenue `revenue_deferred_when_earned`; A/P `wrong_account` (D5) |
| S1 | Hint | The business received currency/funds today. | What came into the business today, and in what form? |
| S1 | Explanation | Cash is received immediately today, increasing the Cash account. | The customer paid cash today, so Cash is affected. Accounts Receivable would apply only if the customer still owed the money, and Unearned Revenue only if the payment came before the work. |
| S2-S4 | | | Shared text above. |
| S5 | Prompt (T2) | Because services were performed today and cash was received, what counter-account earns this inflow? | Which account records what the company did for the customer in exchange for the cash? |
| S5 | Hint | The work is completed today. When work is performed, what account records earnings? | Is the work finished, or does the company still owe it to the customer? |
| S5 | Explanation | Work completed today with immediate cash payment is recorded as Service Revenue. | The work was completed today, so the company has earned revenue: Service Revenue. Unearned Revenue would mean the work is still owed, but it is already done. |
| S6 | Options (D6) | 3 options | Add "Debit Accounts Receivable ${amount} / Credit Service Revenue ${amount}", tagged `wrong_account` (records a credit sale when the customer paid). |
| S6 | Hint | Debit the asset that increased (Cash) and credit the revenue account that increased (Service Revenue). | The customer paid today for work finished today. What did the company receive, and what did it earn? |
| S6 | Explanation | Debit Cash ${amount} and Credit Service Revenue ${amount}. | Debit Cash ${amount} and Credit Service Revenue ${amount}. Payment and work happened on the same day, so there is no receivable and no unearned revenue. |
| S7 | Hint | Cash increases total assets. Does earning revenue increase owner's equity? | Shared opening, then: Did the company earn anything, and does anyone owe anyone afterward? |
| S7 | Explanation | Assets increase by ${amount} and Equity increases by ${amount} through earned revenue. | Assets increase by ${amount} (Cash). Equity increases by ${amount} because revenue increases equity. Liabilities are unchanged because nothing is owed afterward. |

### service_on_credit (Dr Accounts Receivable / Cr Service Revenue)

| Stage | Field | Current | Proposed |
|---|---|---|---|
| S1 | Options | A/P, Unearned Revenue untagged | A/P `wrong_account`; Unearned Revenue `revenue_deferred_when_earned` (D5) |
| S1 | Hint | The customer owes the company for work completed on credit. | Did any money arrive today? If not, what does the company hold instead? |
| S1 | Explanation (T3) | Accounts Receivable is debited because the company holds a claim to collect cash in the future. | The company has the right to collect from the customer later, which is Accounts Receivable. Cash is tempting, but no money arrived today. |
| S2 | Hint | Accounts Receivable is an economic resource (a legal claim) owned and controlled by the company. | Is a right to collect money later something the company owns, something it owes, or the owners' stake? |
| S2 | Explanation | Accounts Receivable is an Asset. | Accounts Receivable is an Asset: the right to collect cash later. It is not Revenue. Revenue records the earning; the receivable records the amount still to be collected. |
| S3 | Hint | A new claim to collect future money was created. | For the work completed today, does the invoice create a new amount to collect or settle an existing one? |
| S3 | Explanation | Accounts Receivable increases when new services are billed on account. | The customer now owes the company money, so Accounts Receivable increases. |
| S4 | | | Shared asset text above. |
| S5 | Prompt (T1) | The work has been completed today. What account earns this inflow under accrual accounting? | The work was completed today but has not been paid for. Which account balances the entry? |
| S5 | Options | Unearned Revenue, A/P untagged | Unearned Revenue `revenue_deferred_when_earned`; A/P `wrong_account` (D5) |
| S5 | Hint | Under accrual accounting, revenue is recognized when performance is satisfied, regardless of when cash is collected. | Under accrual accounting, does revenue wait for the cash, or for the work? |
| S5 | Explanation | Service Revenue is credited because the service was performed today. | Revenue is recorded when the work is done, not when cash arrives, so the balancing account is Service Revenue. Cash would record a payment that has not happened. |
| S6 | Hint | Debit Accounts Receivable (asset up) and Credit Service Revenue (equity up). | The work is done but unpaid. What does the company now hold, and what did it earn? |
| S6 | Explanation | Debit Accounts Receivable ${amount} and Credit Service Revenue ${amount}. | Debit Accounts Receivable ${amount} and Credit Service Revenue ${amount}. Debiting Cash would record money that has not been collected yet. |
| S7 | Hint | Accounts Receivable increases assets; Service Revenue increases equity. | Shared opening, then: Did the company earn anything, and did any cash move? |
| S7 | Explanation | Assets increase by ${amount} (+Accounts Receivable) and Equity increases by ${amount} (+Service Revenue). | Assets increase by ${amount} (Accounts Receivable). Equity increases by ${amount} because revenue increases equity. Liabilities are unchanged. Cash stays the same until the customer pays. |

### collect_receivable (Dr Cash / Cr Accounts Receivable)

| Stage | Field | Current | Proposed |
|---|---|---|---|
| S1 | Options | A/P untagged | A/P `wrong_account` (D5) |
| S1 | Hint | The business received cash funds today. | What arrived in the business today? |
| S1 | Explanation (T3) | Cash is received, so Cash is debited. | The customer paid money today, so Cash is affected. Service Revenue is tempting, but this payment is for work already recorded as revenue last month. |
| S2 | Options | Asset, Liability, Equity | Add "Revenue", matching the other families. |
| S2-S4 | | | Shared text above. (Replaces the S2 hint "Cash is an asset.") |
| S5 | Options | Unearned Revenue, A/P untagged | Both `wrong_account` (D5) |
| S5 | Hint | Revenue was already recognized in the prior period. Crediting revenue again would double-count sales! Which asset was holding the customer's promise to pay? | Revenue was recorded last month, when the work was done. What has the company been holding since then that this payment settles? |
| S5 | Explanation | Accounts Receivable is credited to clear the existing claim. No new revenue is recorded. | The payment settles the amount the customer owed, so the balancing account is Accounts Receivable, which decreases. Crediting Service Revenue again would count last month's work twice. |
| S6 | Hint | Debit Cash (cash up) and Credit Accounts Receivable (receivable cleared). | The work was recorded as revenue last month. What did the company receive today, and what did that payment settle? |
| S6 | Explanation | Debit Cash ${amount} and Credit Accounts Receivable ${amount}. | Debit Cash ${amount} and Credit Accounts Receivable ${amount}. Revenue is not touched; it was recorded when the work was done. |
| S7 | Options (D4) | "Assets increase / Liabilities increase" untagged | Tagged `wrong_account` (treats the payment as a customer advance) |
| S7 | Hint | One asset (Cash) went up, and another asset (Accounts Receivable) went down by the exact same amount. | Shared opening, then: Did total assets change, or did one asset turn into another? |
| S7 | Explanation | This is an asset exchange. Total assets are unchanged, and no new equity/revenue is recognized. | Cash increases by ${amount} and Accounts Receivable decreases by ${amount}: one asset turned into another. Total assets, liabilities, and equity are unchanged. |

### earn_advance (Dr Unearned Revenue / Cr Service Revenue)

| Stage | Field | Current | Proposed |
|---|---|---|---|
| S1 | Prompt | Services paid for last month are now completed today. Which liability account is fulfilled and debited? | No cash changes hands today. Which account recorded what the company owed this customer before today's work? |
| S1 | Options | A/R untagged | A/R `advance_confused_with_receivable` (D5) |
| S1 | Hint | No cash changed hands today. The liability representing future work is being discharged. | When the customer paid last month, what did the company owe in return? |
| S1 | Explanation (T3) | Unearned Revenue is fulfilled and reduced by debiting it. | Last month's payment created an obligation to do the work, recorded as Unearned Revenue. Today's work settles it. Cash is not involved, because the money arrived last month. |
| S2 | Hint | Unearned Revenue represents an obligation to deliver goods or services to the customer. | Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake? |
| S2 | Explanation | Unearned Revenue is a Liability. | Unearned Revenue is a Liability: the company owes the customer work, or a refund. Despite its name, it is not a Revenue account. |
| S3 | Prompt | Does the liability Unearned Revenue increase or decrease as the promised service is delivered? | Does Unearned Revenue increase or decrease as the promised service is delivered? (Drops "the liability", the S2 answer.) |
| S3 | Hint | The obligation to do future work has been fulfilled, reducing the outstanding liability. | Before today, the company owed the customer this work. After finishing it, how much does it still owe? |
| S3 | Explanation | Unearned Revenue decreases when performance is satisfied. | The obligation has been fulfilled, so Unearned Revenue decreases. |
| S4 | Hint | Liabilities have a normal credit balance. Decreases go on the opposite side (debit). | Which side is a liability's normal balance, and does a decrease go on that side or the opposite one? |
| S4 | Explanation | Liabilities decrease by Debit. | Liabilities have a normal credit balance, so a decrease is recorded as a Debit (left side). This debit does not mean money went out; no cash moved today. |
| S5 | Options | A/R untagged | A/R `advance_confused_with_receivable` (D5) |
| S5 | Hint | Performance is complete. Revenue is recognized when earned. | After completing the prepaid work today, has the company earned anything new? |
| S5 | Explanation | Service Revenue is credited because the earnings process is complete. | Completing the work earns revenue, so the balancing account is Service Revenue. Cash would record a payment, but the customer paid last month. |
| S6 | Hint | Debit Unearned Revenue (liability down) and Credit Service Revenue (revenue up). | The customer paid last month, and the work is done today. What does the company no longer owe, and what has it now earned? |
| S6 | Explanation | Debit Unearned Revenue ${amount} and Credit Service Revenue ${amount}. | Debit Unearned Revenue ${amount} and Credit Service Revenue ${amount}. No cash is recorded today, because the cash was recorded when it arrived last month. |
| S7 | Options (D3) | "No effect on any balance; memo entry only." (untagged) | Replace with "Assets increase by ${amount} (+Cash); Liabilities increase by ${amount} (+Unearned Revenue); Equity unchanged.", tagged `cash_recorded_on_earning_advance` (records the advance again). |
| S7 | Hint | Liabilities decreased because work was done; Equity increased because revenue was earned. Total assets did not change. | Shared opening, then: Did any cash move today? |
| S7 | Explanation | Liabilities decrease by ${amount} and Equity increases by ${amount}. Total assets are unaffected. | Liabilities decrease by ${amount} (Unearned Revenue). Equity increases by ${amount} because revenue increases equity. Total assets are unchanged, because the cash arrived last month. |

### Not in this batch

- X2 (choosing the hint based on which wrong option was picked) is a separate code change, planned for after this batch.
- `wrong_account` is a weak label for three distractors: A/R on a cash sale, Unearned Revenue on a collection, and liabilities rising on a collection. Each describes a real mistake. They could get their own labels later; for now, the generic hint ("who owes whom afterward?") fits all three.

### Delegated review and implementation - 2026-10-03

The user explicitly delegated review to Codex ("I don't have the energy to review it, can you?"). Codex reviewed economic timing, canonical postings, stage independence, plausible distractors, and terminology, then applied this batch. This is a delegated assistant review, not a claim that the user read or personally approved every string.

Review refinements: the credit-sale direction hint now avoids assuming a prior zero balance; earn_advance S5 now asks "The prepaid work is completed today. Which account balances the entry?" instead of naming revenue in the prompt. Both core classification distinctions and the generic wrong_account tags are retained. The implausible memo option is replaced with opt_assets_up_liab_up, preserving the identity of historical snapshotted options.

Affected active templates were incremented by one version, with delegated reviewer provenance and unchanged scenario wording, amounts, concept tags, fixtures, and rule versions. There are ten affected scenarios (two per revenue family); earlier audit references to nine were a counting error. Saved instances retain their stage snapshots. This batch adds no scenarios.

## Teaching rewrite batch 2 and X2 completed - 2026-10-03

The user instructed Codex to proceed after delegating review. The remaining eight-family teaching package is recorded in [TEACHING-BATCH2-REVIEW.md](TEACHING-BATCH2-REVIEW.md). All seven stages now have question hints and causal explanations; early prompts no longer name the account/category/side being asked, and category options use bare category names. Shared hints and explanations remain family-specific where economic timing differs. Current rent now includes the prepaid misconception; insurance consumption includes the no-entry misconception with `prepaid_not_expensed_on_consumption`. Missing account/equation tags were attached without changing canonical accounts or rules.

X2 now routes first-error feedback by the actual submitted option. Reviewed `mistake_hints` are snapshotted on each stage at generation time; saved instances without them retain their original stage hint. Repeated hint requests preserve the submitted error until the stage advances. Offline help uses the snapshot when present; tutor requests no longer carry a previous stage's error forward. Unknown or generic tags fall back to the reviewed stage question.

All eighteen active templates receive a new version for X2; the eight remaining families also include the teaching rewrite in that version. Provenance explicitly identifies delegated assistant review, not personal user review. Rule versions and canonical fixtures remain unchanged. Tests cover all thirteen families at all allowed amounts with faded scaffolding, selected-option routing, one assisted retry, correct-option tags, historical fallback, restart/duplicate-save immutability, and tutor context. Full `go test ./...` and `go vet ./...` passed before the final additional regression checks; affected-package checks passed afterward.

Windows terminal checks: equipment full scaffold showed its targeted expense-misconception hint, retry, causal explanation, and scrollable offline explanation. An isolated in-memory harness completed prepaid consumption with faded scaffold, rejecting no-entry, accepting the assisted correction and equation effect, and showing the balanced recap. Remaining low-severity distractor wording and X4 repeated offline explanation lecture are deferred; no new family/account is needed for the pilot.

## Integration follow-up at 90 - 2026-10-03

See [CONTENT-INTEGRATION90.md](CONTENT-INTEGRATION90.md) for exhaustive journal/recap, exam/snapshot, tutoring/context and evidence checks. X4 resolved for bank-backed explanations: OfflineTutor now returns the reviewed stage explanation without appending the repeated general lecture. Missing-teaching requests retain generic fallback. Reviewed explanations also reach optional AI explanation prompts, without being included in hint requests. No active template/version was rewritten.

Journal next now renders allowed amounts and excludes inactive records; exams filter eligibility and resume original saved snapshots. Remaining priorities are explicit reviewed contrast scheduling, distinct-setting retrieval evidence, missing source/role/no-entry/amount distractors and weak payable fillers. The review is complete; those quality items remain open rather than being inferred from 90 scenarios or a green test suite.


## Reviewed transfer follow-up at 90 - 2026-10-03

[Transfer review](TRANSFER-REVIEW90.md) records policy 1: all 90 reviewed setting bindings, 13 matched pairs, bounded assisted after-error selection and six priority distractor profiles. Evidence version 2 requires two successful reviewed groups and delayed group-changing retrieval for either scaffold reduction; assisted exposures postpone delay and guided comparisons cannot inflate independent evidence. Migration 9 preserves snapshots/history with empty historical transfer metadata. The earlier scheduler/evidence limitations described above are resolved by this follow-up.

Bank remains 90 active/one retired/13 families. Reviewed contrast and distinct-setting eligibility are implemented; remaining weak fillers and amount-interpretation distractors need focused review before further breadth. Course evidence, 180 and native-platform/live-provider release checks remain open. See HANDOFF for actual final verification results.

# Pilot batch 1: review and activation

Six scenarios at $50/$100/$200, version 1, rule version 1. Reviewed by Codex under the user's delegated review and instruction to proceed on 2026-10-03; metadata explicitly records assistant review, not personal user review. No syllabus/slides are available, so all scenarios stay within existing families and account metadata.

The complete reviewed package is [PILOT-BATCH1.json](PILOT-BATCH1.json). All seven stages, answer choices, explanations, and option-specific hints are visible in [PILOT-BATCH1-PREVIEW.md](PILOT-BATCH1-PREVIEW.md). These files are review artifacts; runtime loads the copies in curriculum/seed-questions.json. Regression tests require those copies to match the approved package.

| ID | Meaningful fact / misconception | Canonical debit / credit | Contrast partner | Concepts |
|---|---|---|---|---|
| cash_service_training | Finished workshop, paid at its end, no earlier payment/invoice or future work | Cash / Service Revenue | customer_advance_training; earn_advance_training | Cash versus earning; asset classification; debit/credit |
| customer_advance_training | Workshop reservations paid before any instruction, entire promise still outstanding | Cash / Unearned Revenue | cash_service_training; earn_advance_training | Cash versus earning; obligation; debit/credit |
| earn_advance_training | Entire workshop completed after earlier payment and recorded obligation; no new cash | Unearned Revenue / Service Revenue | customer_advance_training; cash_service_training | Earned versus unearned; liability decrease; noncash earning |
| service_on_credit_repairs | Work complete and item released, invoice payable in 30 days, no deposit/payment | Accounts Receivable / Service Revenue | collect_receivable_repairs | Receivable classification; earning before collection |
| collect_receivable_repairs | Full settlement of prior repair invoice, earning and unpaid amount recorded earlier, no new work | Cash / Accounts Receivable | service_on_credit_repairs | Collection versus earning; asset exchange; duplicate revenue |
| borrow_cash_future_payroll | Loan received today, earmarked for later payroll, no spending today | Cash / Notes Payable | cash_service_training; issue_shares_basic | Borrowing versus earning/investment; intent versus today's event |

The training trio changes performance/payment timing while holding the setting stable. The repair pair changes whether work or settlement happens today. The borrowing question adds an intended future use cue that must not alter today's entry. These are matched reasoning contrasts, not six isolated business renamings. Amount substitutions and stages are counted separately from scenario breadth.

Teaching review: every scenario inherits complete reviewed category/direction/side guidance and supplies its own identity, counter-account, entry, and equation teaching. No scenario names the account answer. Hints ask about economic facts; explanations explain why the tempting timing/source alternative fails. Wrong choices use reviewed family distractors and tag routing. No new account or family is introduced. Each possible parameter amount balances and reconciles; contrast entries balance but are rejected for the opposite event.

Before activation, the six-scenario regression passed all allowed amounts and full/intermediate/faded sessions, including first errors and assisted corrections. Representative Windows terminal previews used an isolated harness with an in-memory database: a full customer advance rejected Service Revenue at identity and accepted Cash; a faded earning-an-advance rejected a second cash receipt and accepted the obligation-to-revenue entry. The preview artifact records every other stage and choice. Current-family full/faded terminal and offline-help checks were completed before drafting this batch.

Activation: six approved_active templates were copied into the seed bank with actual delegated-review timestamp 2026-10-03T19:55:34Z. Bank now contains 24 active scenarios, one retired, 13 active families, 72 scenario/amount combinations, and 168 full-scaffold scenario-stage combinations. Historical snapshots remain unchanged. Explicit contrast links are review metadata here; automatic matched-pair scheduling remains a later integration decision.

Final full-suite checks and fresh-binary reconciliation are recorded in HANDOFF.md. No release or commit is implied by activation into the local source bank. Next batch: six prepaid/expense and remaining-gap scenarios toward the 30-scenario pilot.

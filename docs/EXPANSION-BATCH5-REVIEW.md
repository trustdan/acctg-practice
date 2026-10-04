# Expansion batch 5: prior transactions and partial amounts

Reviewer **Codex (review delegated by trustdan)**; approved at **2026-10-03T23:08:04Z** under continuing delegation recorded in HANDOFF.md. This records assistant review rather than personal reading of each string. Complete package: [EXPANSION-BATCH5.json](EXPANSION-BATCH5.json). Seven-stage teaching/options/tags/stored hints: [EXPANSION-BATCH5-PREVIEW.md](EXPANSION-BATCH5-PREVIEW.md).

| Scenario | Reasoning cue | Canonical entry |
|---|---|---|
| cash_service_separate_job | Today's payment belongs to a new completed job; an older invoice stays unpaid. | Dr Cash / Cr Service Revenue |
| service_on_credit_new_client | Another client's earlier receipt does not pay the new client's completed job. | Dr Accounts Receivable / Cr Service Revenue |
| collect_receivable_installment | Partial collection preserves the unpaid remainder of an already-recorded invoice. | Dr Cash / Cr Accounts Receivable |
| customer_advance_additional_receipt | Final additional payment arrives before instruction; the earlier advance stays recorded. | Dr Cash / Cr Unearned Revenue |
| earn_advance_completed_session | One separately priced prepaid session is completed; other prepaid sessions remain owed. | Dr Unearned Revenue / Cr Service Revenue |

No syllabus/slides were located in the repository inventory. All five use existing accounts/families/fixtures/rule version 1 and $50/$100/$200, version 1. The stated amount means today's receipt or completed session, not the larger total. No allocation calculation, discounts/write-offs, new rules, or multi-event entry is introduced.

Reviewed all inherited/custom prompts, hints, explanations, options/tags, and stored first-error hints. Each scenario customizes identity/counter-account/entry/equation teaching; completed-session earning also customizes direction to preserve the remaining obligation. Draft review refined credit to a new client because its shared direction prompt names a new customer, and advance to a final additional receipt because already-paid hints would be ambiguous for an unpaid course remainder. Earlier drafts stayed inactive. Existing active versions/family teaching remain unchanged.

Contrast regressions reject collection versus new earning, cash versus credit earning, advance versus earning, credit earning versus releasing an advance, and collection versus advance even though opposite entries balance. The cash-service old-invoice misconception is covered by engine contrast and equation distractor; this batch adds no new family option sets. Generic filler options/unnarrowed wrong-account hints remain previously deferred limitations. Shared completed-work prompts refer to today's completed session; custom teaching preserves later sessions.

Draft go test ./internal/drill -run TestExpansionBatch5 -count=1 passed, including after refinements and ledger checks. Shared checks cover strict loading, supported wording, duplicate/answer-leak screening, all amounts/scaffolds, canonical entries/effects, rendering, wrong-answer/retry sessions, and balanced opposite-event rejection. Publication equality verifies active copies. Ledger checks preserve old receivables/unrelated earlier cash, unpaid invoice and unearned-session remainders, and an earlier smaller advance plus today's increment. Opening amounts are illustrative test inputs, not extra scenario parameters.

Manual Windows PTY checks before activation used in-memory SQLite, explicit OfflineTutor, isolated auth/cache, and an ignored harness:

- Faded separate cash job rejected receivable debit and asset-swap effect, accepted assisted corrections, and displayed old-invoice-versus-new-earning feedback.
- Faded completed session rejected new cash, accepted Unearned Revenue/Service Revenue and equation effect, and rendered balanced noncash recap.
- Full additional advance completed first four stages, rejected revenue at counter-account, showed still-owed-work hint, accepted liability retry, and rendered incremental-obligation feedback plus scrollable offline help.

All quit cleanly. Credit/collection were reviewed in complete rendered artifacts and automated checks, not individually walked in terminal. No personal learner history/live providers used. Harness source removed; compiled scratch remains ignored.

Initial full-suite run exposed existing TestTUIApproveAndRejectCandidate failure: offline catering's “in two months” did not match the advance future indicators. Added “in advance” to that pending template in internal/candidate/generator.go without loosening gates. New candidate_test.go regression checks 30 deterministic seeds for passing wording while staying pending. Regression passed; TUI approval/rejection passed five runs. Stored/active generated candidates untouched.

Final go test ./... and go vet ./... passed; gofmt -l . clean; git diff --check passed. Root build succeeded with existing nonfatal module stat-cache permission warning. Fresh executable reconciled **45/45** journal/T-account/equation effects and listed approved additions. CLI uses a fixed reconciliation amount; tests cover allowed amounts.

Gate passed: **45 active, one retired, 13 families**, 135 scenario/amount combinations and 315 full-scaffold scenario-stage combinations; 27 expansion scenarios have teaching overrides. Five revenue families and borrowing have four each; other seven families have three each. Contrast links are documented/regression-tested, not automatically scheduled. Next checkpoint 50; 45 more scenarios to 90. No commit, release-version change, or publication.

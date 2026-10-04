# Expansion batch 6 review

Approved under continuing delegated review on 2026-10-03 by Codex (review delegated by trustdan), recorded at 2026-10-03T23:19:42Z. This records assistant review, not independent personal human review. The reviewed package is [EXPANSION-BATCH6.json](EXPANSION-BATCH6.json); [EXPANSION-BATCH6-PREVIEW.md](EXPANSION-BATCH6-PREVIEW.md) captures all seven stages, options, explanations, error tags, and stored mistake hints. Exact approved copies are active in the seed bank.

| Scenario | Recognition and contrast reviewed |
| --- | --- |
| equipment_cash_existing_note | Delivered equipment bought with existing cash: debit Equipment, credit Cash. An unrelated old note remains unchanged; neither borrowing nor principal settlement occurs today. |
| prepaid_purchase_machine_cover | Future insurance on already owned machinery: debit Prepaid Insurance, credit Cash. The machine is the object insured, not the resource acquired today. |
| prepaid_consumption_no_claim | Previously recorded protection expires with time: debit Insurance Expense, credit Prepaid Insurance. Consumption does not require an accident, claim, or new cash payment. Only the used portion is recognized. |
| repay_principal_final_payment | Remaining principal is paid in full: debit Notes Payable, credit Cash. No interest, fees, expense, or new loan; the note reaches zero. |
| dividend_cash_prior_earnings | Distribution declared and paid today from previously recorded profits: debit Dividends, credit Cash. Prior profits do not create revenue today; no earlier declaration/payable, wages, or debt settlement is implied. |

Each addition uses version 1, existing rule version 1, and $50/$100/$200 transaction amounts. Scenario-specific teaching covers the recognition distinction and associated mistakes; the principal scenario also customizes direction teaching. Shared family stages were reviewed alongside overrides. No new account, family, or rule was introduced. Repository inventory contained no actual syllabus/slides; further subjects remain outside the supported scope.

## Verification

The draft was checked before activation. Expansion regression tests exercise every allowed amount and full/intermediate/faded scaffold, canonical entries, equation effects, balanced recaps, rendered teaching, wrong answers and assisted retries, and rejection of balanced entries belonging to contrasted events. Publication tests compare active scenarios with the reviewed package exactly. Additional ledger checks open an existing note and confirm equipment purchase leaves it unchanged, while final principal payment clears it, for all three amounts.

Manual Windows terminal checks before activation covered full equipment purchase with a wrong Notes Payable counter-account, targeted hint, corrected Cash answer, old-debt explanation, and scrollable offline help; faded no-claim consumption with a wrong no-entry answer, retry and balanced recap; and faded final repayment with a wrong expense answer, retry and balanced recap. The equipment walkthrough stopped after counter-account/help. Insurance purchase and dividends were reviewed in the complete rendered artifact and automated checks, not individual terminal walkthroughs. Checks used an in-memory database, explicit OfflineTutor, and isolated ignored auth/cache; no learner history or provider calls. Scratch harness source was removed afterward.

After activation, go test ./... and go vet ./... passed; gofmt -l . was clean and git diff --check passed with existing CRLF conversion warnings. Root acctg.exe rebuilt successfully with the existing nonfatal module stat-cache permission warning. The fresh executable reconciled 50/50 journal, T-account, and equation results and listed approved additions in the active bank. CLI reconciliation uses a fixed amount; regressions separately cover allowed scenario amounts.

The bank now has 50 active scenarios, one retired, and 13 families. See [CONTENT-CHECKPOINT50.md](CONTENT-CHECKPOINT50.md) for cumulative coverage and repetition review. Existing versions, earlier work, and historical snapshots were preserved. No commit, release version change, or external publication was performed.

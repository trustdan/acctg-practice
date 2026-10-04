# Quality closeout at 90

Reviewed 2026-10-04T02:17:25Z (October 3, 2026, America/Chicago) by Codex under trustdan's continuing delegated review. This records assistant review, not personal human reading. The bank retains 90 active scenarios, one retired record and 13 families. No syllabus or slides are available; no new subject or accounting rule is introduced.

## Amount interpretation

Pedagogy policy 2 adds ten exact-scenario profiles. Correct accounts, keys, canonical amounts and statement effects remain unchanged. Each new balanced-entry choice uses the correct debit/credit accounts with the wrong amount scope, the existing `wrong_amount` tag and a stored causal hint. Where four choices already exist, the side-reversal choice is replaced; cash, timing, payable and no-entry misconceptions remain available. Three-choice stages gain a fourth option.

| Scenario | Selectable amount misconception |
|---|---|
| collect_receivable_installment | Full invoice instead of today's installment |
| customer_advance_additional_receipt | Full course price including the earlier receipt |
| earn_advance_completed_session | All sessions including work still owed |
| earn_advance_final_remaining | Original advance including earning already recorded |
| prepaid_consumption_partial_policy | Full premium including unused protection |
| prepaid_consumption_final_remaining | Original premium including earlier consumption |
| prepaid_purchase_discounted_policy | Advertised price instead of agreed paid cost |
| repay_principal_early_partial | Entire principal instead of accepted partial payment |
| dividend_cash_total_distribution | Total multiplied by recipient count |
| equipment_cash_additional_machine | Combined cost including the previously recorded machine |

Original totals and recipient counts are not specified in these published scenarios. The new choices therefore name the wrong scope explicitly rather than inventing numerical facts. They test amount interpretation, not multi-parameter arithmetic. Numeric journal grading separately rejects balanced entries at the wrong amount. Profiles validate against both the exact scenario ID and family.

## Contextual account choices

All 32 cash-service, credit-service, collection and advance scenarios replace generic Accounts Payable fillers in account stages. Cash receipt is contrasted with ownership or a receivable; credit service contrasts the unpaid resource with earning and ownership; collection contrasts incoming cash with the old claim and lender debt; advances distinguish the incoming resource from loan financing. Real payable misconceptions remain in current cash expenses and purchases. Counter-account hints now explicitly distinguish service payment from loans, customer advances from ownership, and old-invoice collection from future-work payment. General category, direction and side choices remain straightforward alternatives; no artificial diagnostic tags are invented.

[Generated preview](QUALITY-PREVIEW90.md) contains the scenario facts, relevant stage prompts, all choices, keys, tags, stored hints and explanations at $100/seed 101. All changed account stages and ten amount stages were read, including the stockholder-service capital option retained from policy 1. Review caught generic loan/ownership hints and replaced them with the causal distinctions above. Other allowed amounts and all scaffold levels are covered by automated regressions.

## Versioning and evidence

Policy 2 changes newly generated instance IDs (`-p2`) and saved option/hint snapshots. Existing question versions, prior saved snapshots and grades remain unchanged. Setting groups and 13 matched pairs retain their reviewed definitions; evidence projection version 2 and migration 9 remain current. Changing only the pedagogy version creates no new setting or independent evidence. No database migration, learner-history rewrite, provider request, commit or release publication is part of this pass.

## Gate

The focused prerequisites for further breadth are complete: integration/replay coverage, focused offline explanations, guided matched contrasts, distinct-setting delayed evidence, six priority misconception repairs, amount interpretation and generic payable cleanup. This is a content/behavioral quality gate, not empirical validation of learning effectiveness. Further expansion can proceed in reviewed five-scenario batches within the supported families when requested. The 180 target is still unstarted; actual course evidence and native-platform/live-provider release checks remain separate open gates.

## Verification

- `go test ./...` and `go vet ./...` passed after final policy approval. `gofmt -l .` produced no output; `git diff --check` passed with existing CRLF warnings.
- Profile regressions cover all sixteen profiled scenarios at their three allowed amounts and all three scaffold levels: option identity, targeted hints, wrong-answer rejection and assisted retry. Bank validation rejects profiles bound to the wrong scenario/family.
- New all-bank checks cover 90 scenarios x three amounts x three seeds (810 generated instances), deterministic option snapshots, unique IDs/text/labels, at most four choices, correct keys without error tags, removal of generic payable fillers and rejection of balanced wrong-amount entries.
- New SQLite checks save all 90 current-policy snapshots, attempt duplicate writes with changed wording, reopen and compare all persisted instance fields including choices/order/hints and pedagogy. Concepts live on attempts rather than the instance table. Existing historical migration/replay checks pass; no personal learner database was opened.
- Existing exhaustive journal/recap, full/intermediate/faded tutor context, exam snapshot restart/resume, transfer/assistance/delay and evidence-deduplication checks passed against policy 2.
- Rebuilt root `acctg.exe`; fresh `--reconcile-all --db :memory:` verified 90/90 active questions. CLI reconciliation uses its fixed $1,500 amount; allowed-amount checks are separate regressions. Build emitted the known nonfatal module stat-cache access warning.

The generated preview is artifact inspection, not a new manual interactive terminal session. No live-provider or native macOS/Linux check was performed. Existing representative Windows PTY checks remain recorded in the integration and transfer reviews.

# Question and teaching expansion plan

Current status (October 3, 2026): the focused quality work at 90 is complete under pedagogy policy 2. [Quality closeout](QUALITY-REVIEW90.md) records ten amount profiles, contextual account choices and the remaining external validation gates. Earlier sections retain the findings and next steps at their original checkpoints.

Proposed roadmap, 2026-10-02. This does not change or approve active content.

## Baseline

At the 2026-10-02 planning baseline, the shipped bank contained 19 scenarios: 18 active and one retired, across 13 active transaction families. Each active scenario had three allowed amounts (54 scenario/amount combinations) and seven full-scaffold stages (126 scenario-stage combinations). Amount substitutions and substeps are not independent scenarios. Current progress is recorded below.

At that baseline, the generator defined 92 hint/explanation pairs, including generic fallback stages, and the offline tutor provided nine error-tag hint cases and eight family contrast cases. The original seeds had no custom teaching. The status-dependent override issue has since been resolved: reviewed `active` and `approved_active` templates both receive overrides. New pilot scenarios include per-scenario teaching; option-specific hints are immutable instance snapshots.

Counts exclude local published candidates, saved AI explanations, and learner history. No syllabus or slides are available locally; preserve supported semantics and verify course scope before adding topics.

## Targets

| Family | Planning baseline | 5x milestone | 10x milestone |
|---|---:|---:|---:|
| Cash services | 2 | 8 | 16 |
| Customer advances | 2 | 8 | 16 |
| Services on credit | 2 | 8 | 16 |
| Receivable collection | 2 | 8 | 16 |
| Earning an advance | 2 | 8 | 16 |
| Current rent expense | 1 | 6 | 12 |
| Borrowing | 1 | 6 | 12 |
| Share issuance | 1 | 6 | 12 |
| Prepaid purchase | 1 | 7 | 14 |
| Prepaid consumption | 1 | 7 | 14 |
| Equipment purchase | 1 | 6 | 12 |
| Principal repayment | 1 | 6 | 12 |
| Dividends | 1 | 6 | 12 |
| **Total** | **18** | **90** | **180** |

First milestone: add 72 scenarios. Second: add another 90. Full scaffolding yields 630 or 1,260 scenario-stage combinations. Every scenario needs reviewed teaching for all applicable stages; shared wording is acceptable where effective.

## Implementation order

1. **Audit the current 18.** Build a matrix of family, stage, misconception, timing, difficulty, contrast partner, and provenance. Flag answer-revealing hints, explanations that only repeat the answer, weak distractors, repetition, and terminology inconsistent with canonical postings. In particular, check dividend account wording. Existing reviewer metadata does not establish that a fresh review occurred.
2. **Improve teaching delivery.** Resolve the status-dependent teaching override behavior through the approved content workflow. Version changed wording and preserve history. Use existing schema where possible; add misconception routing or contrast metadata only when the audit establishes a need. Keep offline teaching complete and optional tutor context accurate.
3. **Pilot at 30 scenarios, in two batches of six additions.** Improve the existing 18 in review groups of 5-10 and draft the first six additions emphasizing cash versus revenue. Complete their quality gate before drafting the next six, emphasizing prepaid versus expense and remaining pilot gaps. Supply each with stage hints, explanations, plausible distractors, and an identified contrast partner. Review in the terminal before scaling; human semantic approval precedes activation.
4. **Expand toward 90 in batches of 5-10 scenarios.** Treat 50, 70, and 90 as cumulative checkpoints, not batch sizes. Fill the target matrix. Vary reasoning cues: prior recognition, performance/payment timing, future benefit/current consumption, source of cash, decreases, and noncash events. Cosmetic business renaming alone does not fulfill breadth goals. Keep wording within supported event/account rules.
5. **Verify integration and teaching coverage.** Check full/faded scaffolding, journal practice, equation recaps, exams, offline tutor, and optional AI context. Inspect whether explicit matched-contrast scheduling needs implementation. Ensure repeated variants and assisted substeps do not inflate mastery.
6. **Expand to 180 after the 90 gate, still in batches of 5-10.** Add distinct transfer contexts and fill coverage gaps. If current families become repetitive, confirm course scope before adding accruals, payables, depreciation, interest, inventory, or multi-account families. Each new family requires reviewed rules, metadata, grading, and regression tests; these topics are outside the first milestone commitment.

## Batch workflow

Use five scenarios by default; increase to at most ten only when review remains manageable. Count each scenario with its complete teaching package as one item. Amount substitutions and individual stages do not count toward the batch size. The two six-scenario pilot batches follow the same workflow.

1. Choose the batch's coverage gaps, misconceptions, and meaningful contrast partners before drafting.
2. Draft the complete package for every scenario: wording, canonical answer, all applicable stage hints/explanations, distractors, tags, and provenance.
3. Review every scenario against the teaching standard and acceptance gates below. Compare it with existing content and its batch peers for repetition, ambiguity, answer leakage, and semantic consistency. Automated checks do not replace human semantic approval.
4. Correct flagged items and rerun affected checks. Keep unresolved drafts outside the active bank; reduce the batch size if review quality becomes difficult to sustain.
5. Record scenario IDs/versions, actual checks, review decisions, unresolved issues, and updated coverage in the handoff. Publish only explicitly approved content. Start the next batch after the current batch's review is complete and issues are resolved or explicitly deferred outside active content.

Milestone counts are targets, not reasons to accept weaker content. Each batch must pass its own quality gate before expansion continues.

## Teaching standard

- Questions explicitly distinguish today, previous recognition/payment, and future performance or benefit.
- Hints ask one short causal question about the current decision without prematurely naming the answer. Preserve one assisted retry.
- Explanations connect economic reality to account, category, direction, side, and equation effect as relevant. Explain why the tempting alternative fails; avoid repeated general lectures.
- Contrasts change a meaningful fact and explain the resulting entry difference. Link separately reviewed scenarios where supported.
- Distractors represent recognizable misconceptions. Balanced but semantically incorrect entries must be rejected.
- Review wording, answer, hints, explanations, distractors, tags, and contrasts together. Record real provenance; never invent approvals.

## Acceptance gates

- Report approved active scenarios separately from drafts, retired records, amount substitutions, and substeps.
- Each scenario has complete stage teaching, explicit timing, supported semantics, correct tags, and human semantic approval.
- All allowed amounts derive balanced canonical entries and reconciled effects; teaching agrees with those entries.
- Verify misconception distractors and contrast pairs beyond balance checks.
- Retired questions stay excluded and historical attempts retain original versions and snapshots.
- When implementation changes, run applicable content/behavior regressions and required Go formatting, tests, and vet checks. Record actual results.
- Manually inspect representative terminal examples with full and faded scaffolding and offline help. Publish approved batches only.

## Progress - 2026-10-03

Existing eighteen scenarios have completed the teaching rewrite and X2 option-specific hint routing under delegated review. Both six-scenario pilot batches are reviewed and active; see [PILOT-BATCH1-REVIEW.md](PILOT-BATCH1-REVIEW.md) and [PILOT-BATCH2-REVIEW.md](PILOT-BATCH2-REVIEW.md). Current local source bank: **30 active, one retired, 13 families**, 90 scenario/amount combinations and 210 full-scaffold scenario-stage combinations. Cash services, advances, credit services, collections, and earning advances each have three scenarios. Rent, borrowing, shares, prepaid purchase, prepaid consumption, principal repayment, and dividends each have two; equipment has one.

The 30-scenario pilot content milestone is reached. Next task: expand toward 90 in reviewed batches of 5-10 (default five), beginning with remaining coverage gaps such as equipment. Include complete teaching and meaningful fact changes; cosmetic business renaming alone does not count. Use 50/70/90 as cumulative checkpoints. The targets remain 90 and 180; 30 does not complete either milestone. Explicit contrast scheduling and later integration decisions remain documented separately. PLAN.md stage history and outstanding release tasks remain unchanged.

## Expansion batch 3 completed - 2026-10-03

[Equipment and financing batch review](EXPANSION-BATCH3-REVIEW.md) records five additions, complete teaching, supported semantics, contrasted facts, delegated provenance, and terminal checks. Current bank: **35 active, one retired, 13 families**, 105 scenario/amount combinations and 245 full-scaffold scenario-stage combinations. Equipment rises from one to three; current rent, borrowing, and share issuance rise from two to three. The five revenue/timing families still each have three. Prepaid purchase, prepaid consumption, principal repayment, and dividends remain at two each. Seventeen expansion scenarios now have teaching overrides.

Next: another reviewed five-scenario batch prioritizing the four families still at two, with meaningful recognition/consumption, principal/distribution, and cash-source contrasts. Keep additional subjects outside scope until course evidence supports them. The next cumulative checkpoint is 50; the 90 milestone still needs 55 scenarios. Scheduling and other previously deferred integration items remain unchanged.

## Expansion batch 4 completed - 2026-10-03

[Coverage timing and payment-purpose review](EXPANSION-BATCH4-REVIEW.md) records five complete additions: a future renewal while current coverage continues, partial prepaid consumption, principal repayment on an equipment loan, distributions to owners who also work for the company, and a loan from a stockholder. Current bank: **40 active, one retired, 13 families**, 120 scenario/amount combinations and 280 full-scaffold scenario-stage combinations. All families have three scenarios except borrowing, which has four. Twenty-two expansion scenarios have teaching overrides.

Next: a reviewed five-scenario batch of cash/revenue timing and prior-recognition contrasts, extending the five revenue families toward their eight-scenario targets. Keep payment timing, performance timing, previously recorded amounts, and today's event explicit; use meaningful transfer cues rather than business renaming. The next checkpoint remains 50, followed by 70 and 90. Fifty additional scenarios remain to reach 90. Previously deferred scheduling/integration and course-evidence boundaries are unchanged.

## Expansion batch 5 completed - 2026-10-03

[Prior transactions and partial amounts review](EXPANSION-BATCH5-REVIEW.md) records five additions: a new paid job alongside an old invoice, a new client's unpaid job alongside another client's receipt, partial invoice collection, additional advance receipt, and earning one prepaid session with others remaining owed. Current bank: **45 active, one retired, 13 families**, 135 scenario/amount combinations and 315 full-scaffold scenario-stage combinations. Five revenue families and borrowing have four scenarios each; seven other families have three each. Twenty-seven expansion scenarios have teaching overrides.

Next: reviewed five-scenario equipment/prepaid/principal/distribution batch to reach **50**, then cumulative coverage/repetition review at that checkpoint. Keep amount meanings, timing, and prior recognition explicit within supported scope. Forty-five more scenarios remain to reach 90; 70/90 remain later checkpoints. Previously deferred integration items are unchanged.

## Expansion batch 6 and 50-scenario checkpoint completed - 2026-10-03

[Batch 6 review](EXPANSION-BATCH6-REVIEW.md) records equipment purchased alongside existing debt, future insurance on owned machinery, expired protection without a claim, final principal repayment, and distributions from prior earnings. Current bank: **50 active, one retired, 13 families**, 150 scenario/amount combinations and 350 full-scaffold scenario-stage combinations. Thirty-two expansion scenarios have teaching overrides.

[The cumulative checkpoint review](CONTENT-CHECKPOINT50.md) records family targets, remaining coverage, semantic repetition, meaningful transfer cues, and older amount sets retained without version changes. Rent and shares have three scenarios each; all other families have four. Next: a reviewed five-scenario batch prioritizing rent and share issuance with distinct reasoning cues. Forty further scenarios remain to reach 90; cumulative checkpoints at 70/90 and previously deferred integration items remain outstanding.

## Expansion batch 7 completed - 2026-10-03

[Rent and ownership-role review](EXPANSION-BATCH7-REVIEW.md) records five additions: current occupancy before opening, shareholder-landlord rent, customer investment, company issuance alongside a private share sale, and a shareholder loan following an earlier investment. Current bank: **55 active, one retired, 13 families**, 165 scenario/amount combinations and 385 full-scaffold scenario-stage combinations. Thirty-seven expansion scenarios have teaching overrides. Rent, shares, and borrowing have five each; other families have four.

Next: reviewed five-scenario revenue timing/prior-recognition batch toward eight per revenue family. Thirty-five additions remain to 90; cumulative checkpoints at 70/90 and previously deferred integration items remain outstanding. No new subjects without course evidence.

## Expansion batch 8 completed - 2026-10-03

[Booking, payment deadlines, and same-day sequence review](EXPANSION-BATCH8-REVIEW.md) records five additions: unpaid booking followed by paid completion, immediately due but unpaid repair, overdue full collection, morning advance, and evening earning of that recorded advance. Shared advance/earning teaching now distinguishes the transaction from the calendar date; saved history retains its snapshots. Current bank: **60 active, one retired, 13 families**, 180 scenario/amount combinations and 420 full-scaffold scenario-stage combinations; 42 expansion scenarios with overrides. Five revenue families plus rent/borrowing/shares have five each; other five families have four.

Next: reviewed five-scenario prepaid purchase/consumption, equipment, principal, and distribution batch toward 65. Thirty additions remain to 90; cumulative checkpoints 70/90 and previously deferred integration items remain outstanding. New subjects await course evidence.

## Expansion batch 9 completed - 2026-10-03

[Purchase cost, consumption, and payment review](EXPANSION-BATCH9-REVIEW.md) records five additions: reduced agreed insurance premium, coverage consumed while closed, secondhand equipment with future use, early partial principal payment, and a total owner distribution spread among recipients. Current bank: **65 active, one retired, 13 families**, with **five active scenarios per family**; 195 scenario/amount combinations, 455 full-scaffold scenario-stage combinations, 47 expansion scenarios with overrides.

Next: reviewed five-scenario revenue timing/recognition batch to reach 70, then the cumulative coverage and repetition review at that checkpoint. Twenty-five additions remain to 90. The target allocation remains eight per revenue family, seven per prepaid family, and six per remaining family; current equal counts do not change that allocation. Previously deferred integration items remain outstanding; new subjects await course evidence.

## Expansion batch 10 and 70-scenario checkpoint completed - 2026-10-03

[Batch 10 review](EXPANSION-BATCH10-REVIEW.md) records third-party service payment/collection, completed unpaid work before invoice paperwork, full prepayment versus performance, and completed prepaid instruction before routine certificate paperwork. Shared prompts now preserve invoice and payer facts; historical snapshots remain unchanged. Current bank: **70 active, one retired, 13 families**, 210 scenario/amount combinations and 490 full-scaffold scenario-stage combinations; 52 expansion scenarios with teaching overrides. Five revenue families have six each; remaining eight have five.

[Cumulative coverage/repetition review](CONTENT-CHECKPOINT70.md) records actual wording, semantic overlaps, distractor limits, and the remaining twenty-scenario allocation. Next: five additions prioritizing rent, borrowing, shares, prepaid purchase, and prepaid consumption with distinct cues. Twenty remain to 90; next cumulative checkpoint 90. Existing supported scope, deferred integration, and course-evidence boundaries remain unchanged.

## Expansion batch 11 completed - 2026-10-03

[Rent, security, and recorded-benefit review](EXPANSION-BATCH11-REVIEW.md) records current rent alongside separate future insurance, a loan secured by retained equipment, new shares issued to an existing lender, future insurance despite a cash-budget expense label, and same-month policy purchase/consumption. Current bank: **75 active, one retired, 13 families**, 225 scenario/amount combinations, 525 full-scaffold scenario-stage combinations; 57 expansion scenarios with teaching overrides. Revenue families, rent, borrowing, shares, and both prepaid families have six each; equipment/principal/dividends have five. Rent/borrowing/shares reach their planned allocation of six.

Next: a reviewed five-scenario batch covering equipment, principal, dividends, prepaid purchase, and prepaid consumption. Fifteen additions remain to 90; remaining allocation is two per revenue family, one per prepaid family, and one each for equipment/principal/dividends. Next cumulative checkpoint 90. Preserve course-evidence boundaries and review novelty before filling counts; deferred integration remains separate.

## Expansion batch 12 completed - 2026-10-03

[Prior costs, counterparty roles, and final coverage review](EXPANSION-BATCH12-REVIEW.md) records an additional machine's new cost, bank-landlord principal repayment, a service-provider stockholder's separate ownership distribution, a brokerage's own future policy, and the final remaining portion of insurance. Current bank: **80 active, one retired, 13 families**, 240 scenario/amount combinations and 560 full-scaffold scenario-stage combinations; 62 expansion scenarios with teaching overrides. Both prepaid families have seven; every other family has six. All eight non-revenue families meet their planned allocations.

Next: two reviewed five-scenario batches, one addition per revenue family in each, toward 85 then 90. Ten remain to 90, all in the five revenue families. Compare against prior wording and recognition cues before approval. Cumulative review at 90 remains outstanding; deferred integration and course-evidence boundaries unchanged.

## Expansion batch 13 completed - 2026-10-03

[Company perspective and payment-application review](EXPANSION-BATCH13-REVIEW.md) records a stockholder acting as customer, customer-held withdrawn cash, old-invoice collection alongside an unpaid booking, internal setup before promised performance, and earning the final remaining course advance. Current bank: **85 active, one retired, 13 families**, 255 scenario/amount combinations and 595 full-scaffold scenario-stage combinations; 67 expansion scenarios with teaching overrides. Revenue and prepaid families have seven each; remaining six families have six.

Next: final five additions, one per revenue family, followed by cumulative actual-wording/coverage/repetition review at 90. Five remain to 90. Counts do not establish new rules or complete deferred integration; preserve course-evidence boundaries and reject purely cosmetic additions.

## Expansion batch 14 and 90-scenario checkpoint completed - 2026-10-03

[Batch 14 review](EXPANSION-BATCH14-REVIEW.md) records company-held currency before deposit, an erroneous cash-sale label, early old-invoice collection, customer-versus-seller books, and prepaid completion alongside a separate unpaid invoice. Complete rendered teaching/options/stored hints reviewed under continuing delegated authorization. The mixed-balance earning case prompted a scoped identify-account hint correction; original snapshots remain preserved.

Current bank: **90 active, one retired, 13 families**, 270 allowed-amount combinations, 630 full-scaffold scenario-stage combinations, 72 additions with teaching overrides. Five revenue families have eight each, both prepaid families seven, other six families six. Every 90 allocation is met. [Cumulative checkpoint](CONTENT-CHECKPOINT90.md) records actual wording, semantic repetition, meaningful transfer comparisons and remaining distractor/teaching limits. Tests, vet, formatting, representative full/faded offline terminal checks, exact publication equality and 90/90 executable reconciliations passed.

The 90 content breadth/review milestone is reached. Next: step 5 integration and teaching-coverage review (journal/exam/recap/offline/AI context, repetition and evidence behavior, matched-contrast scheduling decision, weak distractors and X4). The full integration gate and 180 target remain incomplete. Begin further five-scenario expansion after that review; no new subjects without course evidence. Native-platform/live-provider release checks remain separate.

## Integration and teaching coverage reviewed at 90 - 2026-10-03

[Integration review](CONTENT-INTEGRATION90.md) verifies 270 journal/recap amount combinations, 810 full/intermediate/faded teaching/context combinations, all 90 exam scenarios and 180 review items, SQLite snapshot restart, changed-bank TUI resume, and actual-bank evidence deduplication. Fixed journal next's raw placeholders/fixed amount/retired selection/drill-index coupling, exam eligibility and snapshot resume, missing reviewed explanation in AI context, and X4 repeated lecture for bank-backed offline help. Final tests/vet/formatting and 90/90 executable reconciliation passed; representative offline/journal/exam terminal checks recorded precisely.

Step 5 review is complete; the broader transfer-quality gate remains open. Next: explicitly reviewed pair metadata with bounded after-error contrast scheduling, versioned distinct-setting retrieval evidence, and prioritized missing/weak distractors. Current scheduler still penalizes only exact IDs; time-separated same wording can satisfy faded retrieval criteria. Do not treat the passing integration checks as verified cross-setting transfer or begin expansion toward 180 until these quality issues are handled. Bank remains 90 active/one retired; no new content, history rewrite or external publication.


## Reviewed transfer follow-up at 90 - 2026-10-03

[Transfer review](TRANSFER-REVIEW90.md) records policy 1: all 90 reviewed setting bindings, 13 matched pairs, bounded assisted after-error selection and six priority distractor profiles. Evidence version 2 requires two successful reviewed groups and delayed group-changing retrieval for either scaffold reduction; assisted exposures postpone delay and guided comparisons cannot inflate independent evidence. Migration 9 preserves snapshots/history with empty historical transfer metadata. The earlier scheduler/evidence limitations described above are resolved by this follow-up.

Bank remains 90 active/one retired/13 families. Reviewed contrast and distinct-setting eligibility are implemented; remaining weak fillers and amount-interpretation distractors need focused review before further breadth. Course evidence, 180 and native-platform/live-provider release checks remain open. See HANDOFF for actual final verification results.

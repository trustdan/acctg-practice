# Question and teaching expansion plan

Proposed roadmap, 2026-10-02. This does not change or approve active content.

## Baseline

The shipped bank contains 19 scenarios: 18 active and one retired, across 13 active transaction families. Each active scenario has three allowed amounts (54 scenario/amount combinations) and seven full-scaffold stages (126 scenario-stage combinations). Amount substitutions and substeps are not independent scenarios.

The generator defines 92 hint/explanation pairs, including generic fallback stages. The offline tutor provides nine error-tag hint cases and eight family contrast cases, plus shared guidance. None of the shipped scenarios has custom `teaching` overrides. The generator currently applies overrides only to `approved_active` questions; shipped scenarios use `active`.

Counts exclude local published candidates, saved AI explanations, and learner history. No syllabus or slides are available locally; preserve supported semantics and verify course scope before adding topics.

## Targets

| Family | Current | 5x milestone | 10x milestone |
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

Next task: audit the existing content in groups of 5-10, then prepare the first six-scenario pilot batch for review. Prepare the second six only after the first batch's quality gate. PLAN.md stage history and outstanding release tasks remain unchanged.

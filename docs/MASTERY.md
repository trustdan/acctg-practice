# Evidence, forgetting, and selection

## Proposed v0 model

This is a transparent scheduling heuristic, not a validated memory model. Maintain concept evidence independently from question-family/variant exposure. Inject a UTC clock for tests. Compute decay at read time; do not periodically overwrite stored evidence.

For each concept store beta-model counts alpha/beta (initially 1/1), half_life_days (default 3), last_independent_attempt_at, independent_attempt_count, assisted_attempt_count, and evidence version. Base score = alpha/(alpha+beta). If practiced, retention factor = 2^(-elapsed_days/half_life_days), elapsed floored at zero. Display practice priority from score × retention, alongside evidence count and last practice date. Never show 50% prior as demonstrated mastery; unseen concepts are 'new'.

On a graded first response without aid: success adds 1 to alpha; error adds 1 to beta. Count a concept at most once per instance, using its designated evidence stage. Hints before submission remove that response from independent evidence; preserve it as assisted evidence. An unassisted first error stays an error even if a hinted retry succeeds. Hinted success/revealed answers do not improve base score or reset last independent practice. General tutorials never update counts.

Start with a fixed, configurable half-life; do not invent personalized decay from sparse data. Stage 09 may cautiously adapt it from delayed independent retrieval, with explicit versioning and bounded values. Unassisted performance at trivial repeated wording is not sufficient for graduation.

## Selection

Keep a floor on every eligible question's weight. Concept need = max(0.05, 1 − effective score); new concepts use a configurable introduction weight. Family weight combines its weakest eligible concepts, verified current-course priority, and a recent-exposure penalty. Select using seeded weighted sampling. Reserve about 20% for mixed review; ensure new/weak concepts are not starved. Initially use transparent additive factors rather than multiplying weakness by forgetting, which could zero out a needed concept.

Avoid immediate exact-variant repeats except deliberate contrast/remediation. Track family exposure and wording exposure separately from concept strength. Cosmetic amount changes do not count as transfer. Require several independent successes across at least two reviewed settings and a delayed review before reducing scaffolding; defaults are configurable and pedagogical, not statistical guarantees.

## Acceptance checks

A clock advance lowers effective score without changing stored counts. Assisted retries do not erase errors. Clock rollback does not inflate scores. New concepts are labeled new. Replay reconstructs the same projection. One transaction cannot add five successes to the same concept. A repeatedly memorized variant does not dominate selection. Export and reset require explicit user action; bank upgrades preserve prior attempts.


## Evidence version 2: reviewed transfer

See [transfer review](TRANSFER-REVIEW90.md) for policy 1 and the explicit 90-template setting inventory. Both intermediate and faded guidance require independent success in two reviewed reasoning groups and a delayed success changing the group. Delay is at least ten minutes since the last independent response and last exposure, including assisted practice. Intermediate additionally requires score 0.60/two successes; faded requires score 0.80/three. Errors restore full guidance. These thresholds are scheduling heuristics.

Guided matched comparisons retain contrast assistance through every stage and do not earn independent or transfer credit. Half-life stays three days until two delayed transfers support bounded extension. Cosmetic rewrites and amount changes share reviewed groups; recent groups receive a selection penalty. Duplicate attempt IDs are ignored during replay.

Migration 9 persists policy/setting snapshots. Earlier grades and snapshots remain unchanged; empty historical metadata receives no inferred transfer credit. Rebuilt projections may therefore restore full guidance or shorter half-life despite old success counts. Changed/unreviewed templates require explicit future review before earning setting credit.

## Teaching policy 2 at 90

Policy 2 changes distractor choices and targeted hints while retaining the reviewed groups and comparison pairs from policy 1. Policy version changes do not create a distinct successful setting. Evidence projection remains version 2, with migration 9 preserving instance and attempt metadata. See [quality closeout](QUALITY-REVIEW90.md); the scheduler and mastery thresholds are unchanged by this distractor pass.

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

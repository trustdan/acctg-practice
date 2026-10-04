# Integration and teaching coverage review at 90

Current status (October 3, 2026): the focused quality work at 90 is complete under pedagogy policy 2. [Quality closeout](QUALITY-REVIEW90.md) records ten amount profiles, contextual account choices and the remaining external validation gates. Earlier sections retain the findings and next steps at their original checkpoints.

Reviewed October 3, 2026 (America/Chicago). Scope: the 90 active seed scenarios, one excluded retired record, thirteen supported families. No syllabus/slides available; no new subject, bank version, or scenario was added. This is the expansion-plan step 5 review, with concrete integration fixes and explicit remaining quality work. It does not approve expansion to 180 or claim live-provider/native-platform validation.

## Verified coverage

| Surface | Actual check | Result |
| --- | --- | --- |
| Journal practice | Every active scenario at all three allowed amounts: rendered prompt, canonical submission, semantic grading and reconciliation. | 270/270 scenario/amount combinations pass. |
| Journal next | Traverse all 90 active records, including a raw bank containing the retired record. Check rendered amounts against each template and grading amount; preserve drill position. | All 90 reachable; retired excluded; no placeholder or hard-coded amount. |
| Recap | Complete all seven progressive stages, render recap and verify reconciliation for each allowed amount. | 270/270 recaps pass. |
| Teaching and tutor context | Full/intermediate/faded scaffolds, every allowed amount, every selected option. Check teaching presence, substitution, exact scenario/stage/option context, stored mistake hints and focused offline explanation. | 810 amount/scaffold combinations, 3,510 active stages pass. |
| Optional AI prompt | Explanation mode includes the reviewed stage explanation; hint mode excludes that answer prose. All expanded scenarios tested locally. | Request formatting/context pass; no live provider calls. |
| Exams | Construct from raw seed bank including retired content, complete faded stages, review explanations/postings after finish. | 90 active questions, 180 reviewed items; no retired selection. |
| Exam replay | Save all 90 snapshots to a temporary SQLite file, close/reopen, preserve insertion order, prompt, stages/options/hints and entry; restore response position. Exercise both TUI resume paths against changed/appended current bank. | Original snapshots and question position preserved. Missing snapshots preserve interrupted session and show error. |
| Evidence | Run all 630 full-scaffold responses across actual 90 scenarios within a short simulated window; compare one contribution per instance/concept. Rerun existing error/retry, reference, decay, delayed-retrieval, scheduler and exam-separation tests. | No substep multiplication or rapid-variant delayed retrieval/faded graduation; exam attempts remain separate. |

Exhaustive checks establish structural and behavioral coverage, not an independent semantic rereading of every historical teaching sentence. Existing delegated batch reviews and checkpoint 90 remain the semantic provenance. CLI reconciliation uses a fixed amount separately from these allowed-amount checks.

## Integration fixes

Journal next previously displayed raw scenario templates containing amount placeholders while grading a fixed $1,500 transaction. It also advanced CurrentQuestionIndex, which belongs to the progressive drill, and could visit retired records. Journal navigation now owns its position, skips inactive records, and instantiates wording and entry together using an allowed amount. Entering from a drill retains the current instantiated event. Canonical submission continues to use engine semantics; journal-only navigation does not add mastery evidence.

Exam construction now filters active content itself before sampling and caps size at the eligible count. Resume previously regenerated a shuffle from the current bank, which can change after expansion or wording/version edits. Both TUI resume paths now load the already saved question snapshots in insertion order and replay answers against them. Count/missing-stage/missing-instance checks reject incomplete snapshots; unknown attempt instances are rejected. A failed interactive resume keeps the interrupted session and displays a notice rather than silently starting a new exam. No schema migration or personal data rewrite was needed. The legacy reconstruction API remains for callers supplying an unchanged bank; application resume uses snapshots.

AI explanation formatting previously dropped Request.Explanation even though the TUI supplied reviewed scenario-specific teaching. Explanation requests now include that reference; hint requests exclude it. This improves local context fidelity without changing provider adapters, network behavior, budgets, grades, answer keys or learner history.

X4 is resolved for normal bank-backed offline explanations: return the reviewed stage explanation directly instead of appending the same four-principle lecture and generic family text each time. Active teaching remains complete at every stage/scaffold. Requests without reviewed teaching retain the existing generic reference fallback. This does not claim rewritten generic fallback hints/contrasts or live model output quality.

## Repetition, transfer and scheduling decision

Inspection confirms exact-question anti-repeat memory is three recent IDs with a weight penalty, not a ban. Other wordings from the same family are not separately penalized. Larger families receive more total sampling opportunities when individual needs are equal. Transfer intensity raises selected family weights; it does not enforce any of the documented matched pairs. No scheduler or evidence formula was silently retuned tonight.

Rapid cosmetic variants cannot produce delayed retrieval or faded scaffolding, as verified with the actual expanded bank. They can still increase independent counts and reach intermediate scaffolding. After sufficient time gaps, repeated instances of one wording can satisfy the current faded gate: distinct reviewed settings are not enforced. Existing designated first evidence per instance/concept and assistance separation hold; they do not establish cross-setting transfer. This is a concrete gap against the stronger goal in MASTERY.md, not an assertion that the current projection proves transfer.

Decision: implement explicit reviewed contrast metadata and a bounded after-error contrast queue next, rather than infer matching from family names or inflate concept evidence for showing a partner. Keep the current weighted scheduler until that reviewed metadata and behavior are available. Paired remediation must remain assisted/exposure evidence; later independent retrieval must remain separate. Then define reviewed setting groups and a versioned transfer-evidence requirement before claiming that repeated wording cannot graduate across time. Preserve attempt snapshots and document any projection version change. Neither schema nor mastery history was changed by this review.

## Distractor priorities

1. Full-entry source/role coverage: stockholder service payment versus new capital, early/old invoice receipt versus advance, and completed prepaid job versus creating a new receivable. Existing explanations and engine contrast rejection are correct, but not all alternatives are selectable.
2. No-entry misconceptions for waiting on a bank deposit or relying on a customer's books. Journal free entry covers accounting semantics; the multiple-choice stages do not yet offer those choices.
3. Portion-versus-original-total amounts for final prepaid earning/consumption, added equipment and total dividends. Current allowed-amount checks are not tests of learner amount interpretation.
4. Replace weak payable fillers and add diagnostic tags where they represent a real misconception. Generic classification/direction choices need not each receive an invented diagnostic label.

Published wording/options should change through versioned delegated review. No existing active record was rewritten or retired solely to hide repetition or make this gate appear complete.

## Checks and stopping point

Final go test ./... and go vet ./... passed, gofmt -l . clean, git diff --check passed with existing CRLF warnings. Workspace GOCACHE=dist/go-cache used. Root acctg.exe rebuilt successfully with the existing nonfatal module stat-cache permission warning. Fresh --reconcile-all --db :memory: verified 90/90.

Manual Windows PTY: opened bank-backed offline help on company-held currency and saw only the focused reviewed explanation; entered journal practice and advanced to the mistaken cash-sale scenario with a rendered allowed $200 amount. Separate exam run displayed $50 and withheld explanation/hint/feedback while advancing the first entry answer; interrupted and returned to practice, then quit. Both harnesses used in-memory DB, isolated ignored auth/cache and explicit OfflineTutor; no personal history/provider calls. Journal submission, complete exam review and changed-bank/restart replay are automated checks, not claimed manual full sessions. Temporary harness source removed.

Files: internal/tui/model.go and view.go; internal/exam/runner.go and report.go; internal/storage/exam.go; internal/tutor/offline.go, prompts.go and tutor_test.go. Added expanded-bank/context tests under tui, exam, mastery and tutor, plus storage/exam_snapshot_test.go. Updated expansion plan, audit, checkpoint and handoff. Prior uncommitted work retained; no commit, release-version change or external publication.

The integration/teaching review is complete and the discovered delivery/replay defects are fixed. The broader transfer-quality gate remains open for reviewed contrast scheduling, distinct-setting evidence and distractor improvements. Next session should address those at 90 before further breadth expansion. The 180 target and native-platform/live-provider release checks remain outstanding. This is a verified stopping point for tonight.


## Reviewed transfer follow-up at 90 - 2026-10-03

[Transfer review](TRANSFER-REVIEW90.md) records policy 1: all 90 reviewed setting bindings, 13 matched pairs, bounded assisted after-error selection and six priority distractor profiles. Evidence version 2 requires two successful reviewed groups and delayed group-changing retrieval for either scaffold reduction; assisted exposures postpone delay and guided comparisons cannot inflate independent evidence. Migration 9 preserves snapshots/history with empty historical transfer metadata. The earlier scheduler/evidence limitations described above are resolved by this follow-up.

Bank remains 90 active/one retired/13 families. Reviewed contrast and distinct-setting eligibility are implemented; remaining weak fillers and amount-interpretation distractors need focused review before further breadth. Course evidence, 180 and native-platform/live-provider release checks remain open. See HANDOFF for actual final verification results.

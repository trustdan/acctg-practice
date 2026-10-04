# Bank format and creative validation

The initial JSON format uses schema_version=1, stable IDs, version numbers, family IDs, status, scenario templates, finite parameter values, concept tags, and reviewer provenance. See curriculum/seed-questions.json. The embedded bank has a strict validator and 90 active reviewed scenarios plus one retired record across 13 supported families. Answer keys are derived by the canonical event rule. Expected postings in seed fixtures are regression expectations, never an alternate model-maintained grading path.

## Admission pipeline

Creative output → strict parsing → candidate storage → supported family/parameter check → derived accounting checks → semantic wording review → approval event → active versioned bank. Never activate upon generation, repeated use, or high learner scores. Existing approved templates can generate finite parameter substitutions offline without reviewing each amount.

Reject unsupported accounts, unknown fields, invalid amounts, ambiguous timing, unsupported event families, incorrect concept mappings, or conflicts between wording and event semantics. Balance alone cannot detect a semantically wrong key. The program cannot prove unrestricted natural language means the asserted event; therefore new wording requires human review. New families require accounting review plus event rules and regression tests before approval. Model self-review can assist but cannot substitute for this gate.

Parameter changes are cosmetic only if they preserve semantic preconditions. Structural setting changes need review. Contrast sets use separate event IDs with explicit timing. Do not let the model freely change account names, earnings timing, returnability, taxes, interest, or financing terms inside a template intended to preserve the same answer.

Version published records immutably. Retire flawed variants and exclude future selection; preserve historical instances and annotate any affected grades with a correction event. Never silently rewrite evidence. Review can be a simple local CLI workflow before a creative-mode UI exists.

## Reviewed teaching policy

`curriculum/pedagogy.json` schema 1 / policy 2 binds exact template IDs and versions to reviewed setting groups, symmetric matched contrasts and named distractor profiles. Validation rejects unknown fields, invalid provenance, unavailable references and incompatible profiles. Active seed references require coverage; unreviewed local template versions receive no inferred seed-policy metadata.

Generated IDs include the policy version. Saved stages preserve option IDs, order, diagnostic tags and causal hints. Policy 2 retains policy 1's groups and pairs, adds ten amount-scope profiles and replaces generic payable fillers in four revenue families. It does not change canonical entries or template versions. See [quality review](QUALITY-REVIEW90.md), [preview](QUALITY-PREVIEW90.md) and [transfer policy](TRANSFER-REVIEW90.md).

Migration 9 stores pedagogy metadata on immutable instances and setting/policy fields on attempts. Historical metadata stays empty; original grading and replay remain available. Exam resume loads the original stored instances rather than generating today's bank again. Assistant reviews in the current seed expansion explicitly record the user's delegated review; they must not be described as the user's personal reading.

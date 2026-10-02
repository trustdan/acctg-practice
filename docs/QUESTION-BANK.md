# Bank format and creative validation

The initial JSON format uses schema_version=1, stable IDs, version numbers, family IDs, status, scenario templates, finite parameter values, concept tags, and reviewer provenance. See curriculum/seed-questions.json. These are proposed domain fixtures; Stage 02 must implement a schema and strict validator before runtime use. Answer keys are derived by the canonical event rule. Expected postings in seed fixtures are regression expectations, never an alternate model-maintained grading path.

## Admission pipeline

Creative output → strict parsing → candidate storage → supported family/parameter check → derived accounting checks → semantic wording review → approval event → active versioned bank. Never activate upon generation, repeated use, or high learner scores. Existing approved templates can generate finite parameter substitutions offline without reviewing each amount.

Reject unsupported accounts, unknown fields, invalid amounts, ambiguous timing, unsupported event families, incorrect concept mappings, or conflicts between wording and event semantics. Balance alone cannot detect a semantically wrong key. The program cannot prove unrestricted natural language means the asserted event; therefore new wording requires human review. New families require accounting review plus event rules and regression tests before approval. Model self-review can assist but cannot substitute for this gate.

Parameter changes are cosmetic only if they preserve semantic preconditions. Structural setting changes need review. Contrast sets use separate event IDs with explicit timing. Do not let the model freely change account names, earnings timing, returnability, taxes, interest, or financing terms inside a template intended to preserve the same answer.

Version published records immutably. Retire flawed variants and exclude future selection; preserve historical instances and annotate any affected grades with a correction event. Never silently rewrite evidence. Review can be a simple local CLI workflow before a creative-mode UI exists.

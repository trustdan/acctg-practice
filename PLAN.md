# Formal staged implementation plan

All stages below are pending. Gates are exit criteria, not claims of completion. The first-week priority is Stages 00–07; cut later scope before sacrificing usable offline practice. The one-week target is a prioritization constraint, not a guaranteed estimate.

## Stage 00 — Repository orientation

Read existing instructions and inventory syllabus/materials, project code, target OS, and Go availability. Merge this package rather than overwriting existing instructions. Locate actual course paths. Record unknowns and baseline commands in HANDOFF.md.

Gate: repository state preserved; course files located or absence recorded; agent instruction entry explicitly loaded.

## Stage 01 — Course grounding

Read syllabus and available early slides/examples. Populate COURSE-MAP.md with source locations. Verify vocabulary, progression, grid, assessment rules, and GAAP/IFRS. Identify narrow initial learning objectives. Missing visual must not block generic debit/credit practice.

Gate: each claimed course requirement has evidence; generic assumptions are labeled; initial topic inclusion/exclusion list exists.

## Stage 02 — Domain and bank contract

Define integer money, categories, normal sides, directions, postings, family semantics, instance snapshots, stage answers, and strict versioned JSON schema. Implement bank validation. Formalize seed fixture fields; keep semantic family rules in reviewed code and wording in versioned data.

Gate: bad account IDs, duplicate IDs, invalid parameters, unknown fields, unbalanced fixtures, and unsupported families are rejected. Canonical examples validate.

## Stage 03 — Deterministic accounting engine

Implement cash services, service on credit, customer advance, receivable collection, earning an advance, cash expense, prepaid purchase/consumption, equipment purchase, borrowing/principal repayment, and share issue as course scope permits. Exclude unverified advanced rules. Derive entry and equation effects from one event model.

Gate: table-driven correct/wrong semantic examples; all generated parameter combinations balance; cash collection never creates duplicate revenue. No network/UI dependency.

## Stage 04 — Progressive drill vertical slice

Choose a reviewed fixture and produce ordered stage prompts, stable option IDs, seeded order, answer grading, one retry, and local hints/explanations. Model assistance as explicit evidence state. Include a transaction recap showing both postings.

Gate: a customer-advance drill works from event to balanced entry; targeted hint appears for revenue confusion; answer option labels do not serve as identity; same seed reproduces the instance.

## Stage 05 — Durable history

Select SQLite driver and pin compatible dependencies. Add migrations, sessions, instances, attempts, assistance metadata, and idempotent writes. Use platform data path; ignore personal data in git. Provide backup/export and migration failure messages.

Gate: restart restores progress; replay uses original question versions and option order; duplicate writes cannot inflate history; failed migration preserves the prior database.

## Stage 06 — Mastery and scheduler v0

Implement MASTERY.md with clock/random injection. Track concept evidence and question exposure separately. Add configurable decay, seeded weighted selection, review mix, and anti-repeat policy. Make stats transparent.

Gate: delayed review increases priority; hinted corrections do not erase first errors; one instance does not multiply concept evidence; new concepts are labeled new. Rebuild projections from attempts.

## Stage 07 — Offline Bubble Tea release

Build keyboard UI, help, feedback, summary, mastery view, and graceful exit. Pin maintained library versions after official-doc verification. Support terminal resize, plain/color-safe feedback, narrow screens, and nonblocking update commands. No provider initialization requirement.

Gate: complete a 10-question session offline, quit/restart, and retain attempts. j/k/arrows/Enter and a–d work; text input shortcuts do not collide. Run gofmt, go test ./..., go vet ./... and a manual terminal smoke check. This is the first useful release.

## Stage 08 — Bank breadth and contrast review

Expand to 12–20 reviewed templates if course scope warrants. Add explicit already-recorded receivable contrasts, owner/lender cash sources, decreases, and varied business settings. Add review metadata and retirement mechanism. Do not introduce advanced topics just to meet a count.

Gate: each template has semantic review provenance, concept mapping, rule version, and regression expectation; timing is unambiguous; retired content is not selected.

## Stage 09 — Scaffolding and delayed transfer

Add per-skill scaffold fading, reference-use tracking, and delayed retrieval requirements. Keep half-life fixed unless enough delayed evidence supports a documented bounded adjustment. Provide learner-facing controls for session size and intensity.

Gate: repeated cosmetic variants cannot alone graduate a concept; poor performance restores scaffolding; changes are versioned and replayable.

## Stage 10 — Read-only tutor abstraction

Implement OfflineTutor and fake asynchronous provider with request context, timeout, cancellation, size limit, and fallback. Add ? and e without blocking key input. No automatic AI call on every answer.

Gate: simulated slow/error provider cannot freeze UI; cancellation works; provider output has no path to grades, bank approval, or mastery.

## Stage 11 — Real provider integration (ChatGPT Plus and API keys)

Implement network tutor providers supporting two user-selected connection paths, while keeping offline mode the default:

### Path A: ChatGPT Plus Plan (`OpenAIChatGPTPlanTutor`)
For learners wanting to use their existing ChatGPT Plus/subscription plan allowance without an API key or client secret, using OpenAI's official Sign in with ChatGPT flow for local/open-source personal applications.
- Browser-based login initiated from TUI ("Continue with ChatGPT").
- Ephemeral local HTTP loopback listener for OAuth callback.
- Dynamic Client Registration (DCR), PKCE code challenge, state, and nonce validation.
- Secure local storage of OAuth credentials in user data directory; automatic token refresh.
- Explicit verification of granted ChatGPT-plan permission before initiating inference.
- Model discovery and direct streamed requests to the OpenAI Responses API from Go.
- Enforce subscription-route preview limitations: `store: false`, `stream: true`, and parameter restrictions.
- References:
  - [Registration and sign-in](https://developers.openai.com/siwc/token-sharing-open-source/sign-in?utm_source=chatgpt.com)
  - [Official integration walkthrough](https://developers.openai.com/cookbook/articles/sign-in-with-chatgpt?utm_source=chatgpt.com)
  - [Preview limitations](https://developers.openai.com/siwc/token-sharing-open-source/preview-limitations?utm_source=chatgpt.com)
  - [Codex app-server token consumption](https://developers.openai.com/siwc/token-sharing-open-source/codex-app-server?utm_source=chatgpt.com)

### Path B: User API Keys (`APIKeyTutor` adapters)
For learners preferring their own commercial API key rather than a ChatGPT Plus subscription:
- Supported adapters: **Anthropic** (`ANTHROPIC_API_KEY`), **Google Gemini** (`GEMINI_API_KEY`), and **OpenAI API** (`OPENAI_API_KEY`).
- Keys are supplied either via standard environment variables or entered interactively in the TUI and saved securely to the local OS user-data directory (never committed to git or stored in plain tracked files).
- Explicit per-session request budgeting (`max_requests_per_session`), timeout boundaries, and cancellation support.

### Common Invariants
- Explicit opt-in: The app never silently switches between subscription and paid API billing, and never makes background network calls without user invocation.
- Strict offline fallback: Any network or authentication error falls back cleanly to `OfflineTutor`.
- Machine mode drills remain 100% offline. No credentials or tokens are scraped from external tools or browser cookies.

Gate: TUI allows selecting between ChatGPT Plus sign-in, API key entry (Anthropic/Google/OpenAI), or offline mode; credentials stored securely outside git; streamed and non-streamed responses conform to tutor contract; subscription preview restrictions honored; offline fallback works cleanly on failure; mock tests cover error cases and cancellation.

## Stage 12 — Creative candidate generation

Generate structured candidates for supported families and weak concepts. Parse strictly, derive answers locally, and store candidates outside active bank. Expose preview with wording, assumptions, derived postings, concept tags, and provenance. Store failed proposals with reasons only if useful.

Gate: malformed or contradictory output stays inactive; generated content cannot change a family rule; no generation updates learner evidence.

## Stage 13 — Review and promotion workflow

Implement explicit approve/reject/repair actions. Require semantic wording review for variants; new families require reviewed rule additions and tests. Record reviewer, timestamp, source, rule version, and approval event. Publish immutable versions and retire flawed ones.

Gate: balanced-but-wrong advance/revenue candidate is caught by review; candidates never auto-promote based on use or score; historical attempts survive updates.

## Stage 14 — Journal entry and T-account practice

Add structured multi-line journal-entry entry and consistent T-account/equation rendering. Compare entries independent of line order, aggregate duplicate lines where appropriate, and reject wrong account/side/amount. Match professor visuals only after inspection.

Gate: journal/T-account/equation views reconcile to the same postings; split equivalent lines grade consistently; wrong but balanced entries fail. Generic visual is labeled if classroom visual remains unavailable.

## Stage 15 — Multi-event and statement construction

Add opening balances, transaction batches, adjustments, retained-earnings roll-forward, balance sheet, income statement, and course-required cash flow method. Use verified event-specific cash classifications and noncash disclosure metadata. Start with one small comprehensive case.

Gate: ending balance sheet balances; net income ties to retained earnings; beginning cash + cash flows = ending cash; ending cash ties to ledger; declared versus paid dividends handled distinctly.

## Stage 16 — Exam mode

Reuse engine and approved bank with hints/reference/feedback suppressed until finish. Optional timer; summaries by skill and error pattern. Separate learning versus exam evidence; never expose a hidden answer through tutor shortcuts.

Gate: answers withheld throughout session; results reproducible; history preserved; interrupted session resumes or closes explicitly.

## Stage 17 — Quarter maintenance and release hardening

Document build/run/backup/provider configuration and new-family review procedures. Test target OS and terminal, database upgrade, bank upgrade, and no-provider startup. Add new course units only from evidence. Review accessibility, response latency, and private-data handling.

Gate: clean setup and offline drill from a fresh clone; upgrade preserves history; dependency locks and instruction files agree; unresolved limitations documented.

## Suggested first-week allocation

Day 1: repository/course grounding and contracts. Days 2–3: engine and progressive slice. Days 3–4: persistence and scheduler. Days 4–5: offline TUI. Days 6–7: actual practice, contrast bank review, and fixes. Defer providers and creative UI if offline release needs attention. Adjust this allocation to observed progress rather than rushing gates.

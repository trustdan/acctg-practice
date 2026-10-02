# Architecture and contracts

## Dependency direction

TUI calls engine, storage, and tutor interfaces. Engine depends on account rules, event families, scheduler, and plain domain types. Storage persists evidence and content approval. Tutor consumes read-only snapshots and produces prose or candidate data. Neither tutor nor UI can bypass the grading service.

Suggested eventual packages: cmd/acctg; internal/domain; internal/bank; internal/engine; internal/mastery; internal/storage; internal/tui; internal/tutor. These are proposed boundaries, not directories that must all exist before the first drill works.

## Domain contracts

Account: stable ID, display name, category, normal side, optional contra-of ID. Direction maps through account metadata; contra accounts cannot be inferred from broad category alone. Revenue and expense are temporary accounts affecting equity, not cash accounts. Dividends reduce equity and are not an expense.

Event family: stable semantic ID, rule version, parameters, derived journal postings, stage builders, statement effects, concept mappings, and reviewed variants. Machine generation substitutes only parameters whose meanings remain invariant. Changing 'next month' to 'last month' is a new semantic event, not a cosmetic variation.

Question instance: stable instance ID, family/version, variant/version, parameters, random seed, option order, and derived answer snapshot. Do not use amount changes as evidence of new conceptual transfer.

Attempt: session/instance/stage IDs, concept IDs, submitted option or structured response, deterministic result, first/assisted/retry/revealed status, UTC timestamp, grade version, and optional error tags derived from known distractors. Error tags are hypotheses about the chosen answer, not diagnoses.

## SQLite outline

schema_migrations; sessions; question_instances; attempts; mastery_projections; candidate_questions; approval_events. Bank files are versioned in git; the database references content versions and stores sufficient snapshots for replay. Store approval provenance separately from learner evidence. Transactions make instance/attempt writes atomic; duplicate event IDs are idempotent. Backup before migration; never silently reset a corrupt database. Use platform user-data directories by default and a --data-dir override, not a tracked repository database.

## TUI states

Selection → prompt → answer → correct feedback or hint → optional retry → explanation → next stage → transaction recap → next question. Separate mastery/help screens and a cancellable tutor request. Display one decision at a time, question context, selected option, and visible controls. a–d answers, j/k or arrows navigate, Enter submits, ? hints, e explains, s mastery, q quits, Ctrl-C safely exits. Disable conflicting answer shortcuts during text entry. Exam mode changes feedback policy rather than accounting rules.

## Tutor boundary

Conceptual Go contract: Hint(context, TutorRequest), Explain(context, TutorRequest), GenerateCandidates(context, GenerationRequest). Requests contain a vetted problem snapshot, stage, learner response, relevant rule, and minimal aggregate error context. Responses contain prose or strictly parsed candidate fields; no executable commands or progress edits.

Implement `OfflineTutor` first (Stage 10). Machine-mode drills remain 100% offline with zero provider dependency.

In Stage 11, the app provides a pluggable provider interface supporting both subscription-based and key-based connections:

### Option 1: ChatGPT Plus Plan (`OpenAIChatGPTPlanTutor`)
Official "Sign in with ChatGPT" flow for personal and open-source applications (no API key required):
1. User selects "Continue with ChatGPT" in the TUI.
2. App performs Dynamic Client Registration (DCR) and generates PKCE code verifier, challenge, state, and nonce.
3. App launches default browser to OpenAI sign-in requesting ChatGPT-plan usage permission.
4. Ephemeral local HTTP loopback listener receives the authorization code callback, validates state and nonce, exchanges the code, and confirms the granted plan permission.
5. Tokens are securely stored in the local OS user-data directory (never in git) and refreshed as needed.
6. The provider invokes the OpenAI Responses API directly from Go using streaming (`stream: true`) and `store: false`, respecting subscription-route preview parameter constraints.

### Option 2: Commercial API Keys (`APIKeyTutor` adapters)
For learners who prefer to supply their own API keys instead of a ChatGPT Plus subscription:
- **Anthropic**: `AnthropicTutor` using `ANTHROPIC_API_KEY` (Messages API).
- **Google**: `GoogleGeminiTutor` using `GEMINI_API_KEY` (Gemini API).
- **OpenAI**: `OpenAIAPITutor` using `OPENAI_API_KEY` (Responses or Chat Completions API).
- Keys are loaded from standard environment variables or configured interactively in the TUI and stored securely in the local user-data path. Keys are never saved in git or plain repository files.

### Safety and Boundaries
- Clean fallback: If authentication fails, requests time out, or the network is unreachable, the system gracefully falls back to `OfflineTutor`.
- No silent billing: The app will never silently switch between subscription allowance and paid API keys.
- Request bounds: Set per-session request budgets (`max_requests_per_session`), timeouts, cancellation contexts, and response length caps.
- No credential scraping: Never extract tokens or session cookies from Codex or browser profiles. No hidden background calls during keyboard drills.

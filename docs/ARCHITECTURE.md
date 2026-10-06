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

### Option 3: Local LM Studio server (`LMStudioTutor`, Stage 29)
For learners who run a model on their own machine (no key, account or billing):
- Talks to LM Studio's OpenAI-compatible server (default `http://localhost:1234/v1`): `GET /models` for discovery and `POST /chat/completions` for hints and explanations. `OpenAIAPITutor` and `LMStudioTutor` share one chat-completions helper (`callChatCompletions`).
- Base URL resolution: `LMSTUDIO_BASE_URL`, then the stored `lmstudio_url` (set with `--lmstudio-url` or `[l]` in Tutor Settings), then the default. Only loopback hosts are accepted unless the learner opts in with `--lmstudio-allow-remote`/`LMSTUDIO_ALLOW_REMOTE=1`, so learner context stays on the machine by default.
- No Authorization header is sent unless `LM_API_TOKEN` is set (LM Studio authentication is off by default). The token is never persisted.
- `<think>…</think>` blocks are stripped (including an unterminated block when output is truncated) and `reasoning_content` is never decoded. `max_tokens` defaults to 4096 (`--lmstudio-max-tokens`) so reasoning models are not cut off mid-thought.
- With no model selected, the first non-embedding model reported by `/models` is used. Construction and startup make no network call. Tutor Settings `[6]` performs one asynchronous `/models` check; `--test-llm` probes the server.
- `BuildTutor` gives LM Studio a provider timeout of at least 90 seconds for cold model loads and CPU inference. The TUI request context outlasts `FallbackTutor.Timeout()` so a slow primary falls back offline rather than surfacing an error. The app never launches or manages the LM Studio process.

### Streamed replies (`StreamingTutor`, Stage 30a)
- Optional interface beside `Tutor`: `HintStream`/`ExplainStream(ctx, req, onUpdate)`. `onUpdate` receives a `StreamUpdate` (visible text so far, a `Thinking` flag, provider) synchronously and never after the call returns. On a mid-stream error the returned `Response` carries the partial text alongside the error.
- Implemented by `OpenAIAPITutor` and `LMStudioTutor` through the shared `streamChatCompletions` SSE reader (`stream: true`, `stream_options.include_usage`, `data:` lines, `[DONE]`). Visible text is recomputed from the accumulated raw text, so `<think>` blocks and tags split across chunks never show; `reasoning_content` deltas only set `Thinking`. The 1 MiB response cap applies to the stream (`ErrResponseTooLarge`). Budget usage is recorded once per completed or aborted request. Anthropic, Gemini and ChatGPT plan streaming are Stage 30b; those providers keep the non-streamed path.
- `FallbackTutor` implements `StreamingTutor` and `CanStream()`. Instead of a total deadline, a stream gets a first-chunk deadline (the provider timeout) and then `StreamIdleTimeout` (30s) between chunks. An error before any visible text falls back to `OfflineTutor`; an error after it returns the partial text with `Response.Incomplete` set; caller cancellation returns the caller's error.
- TUI delivery: the request runs in one `tea.Cmd`; the provider callback keeps only the newest update in a one-slot channel, and a self-reissuing `tea.Cmd` reads it every 100ms, so the event loop never blocks and markdown re-rendering is throttled. Esc cancels the context (closing the response body) and keeps any streamed text visible as incomplete. Incomplete replies are never offered for saving.

### Safety and Boundaries
- Clean fallback: If authentication fails, requests time out, or the network is unreachable, the system gracefully falls back to `OfflineTutor`.
- No silent billing: The app will never silently switch between subscription allowance and paid API keys.
- Request bounds: Set per-session request budgets (`max_requests_per_session`), timeouts, cancellation contexts, and response length caps.
- No credential scraping: Never extract tokens or session cookies from Codex or browser profiles. No hidden background calls during keyboard drills.

# Formal staged implementation plan

This document preserves the staged design and historical implementation reports. Gates require recorded verification; historical Complete/PASSED labels do not close missing course-source, native-platform or live-provider checks. The first-week priority is Stages 00–07; cut later scope before sacrificing usable offline practice. The one-week target is a prioritization constraint, not a guaranteed estimate.

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

## Stage 17 — TUI viewport ergonomics and recap scrolling

Implement smooth keyboard-driven vertical scrolling in the transaction recap view (`StateRecap`) and audit all multi-component screens for standard terminal dimensions (e.g., 80x24 rows).
- Support scrolling in `StateRecap` with `j`/`k`, `↑`/`↓`, `PageUp`/`PageDown`, `g` (jump to top), and `G` (jump to bottom).
- Maintain an explicit scroll offset clamped between `0` and `max(0, totalLines - viewportHeight)`.
- Render a clear, subtle visual scroll indicator (e.g. `▲ Scroll: [Line X/Y] (j/k or ↑/↓) ▼`) when content exceeds the terminal height.
- Reset scroll offset to `0` upon loading a new question or re-entering recap.
- Audit tall modal/view components (Financial Statements `StateStatements`, Journal Entry Practice `StateJournalPractice`, Exam Summary `StateExamSummary`, and Candidate Review `StateCandidateReview`) for consistent scroll boundaries and zero visual clipping.

Gate: Transaction recap can be scrolled smoothly via keyboard without clipping on standard terminal heights (>= 15 rows); scroll offset is strictly clamped; indicators reflect position accurately; advancement (`Enter`/`Space`) and navigation remain responsive.

## Stage 18 — Provider authentication repair and dynamic model discovery (Complete)

Resolve external tutor authentication failures and enable real-time model discovery across providers rather than relying on frozen snapshots:
- **OpenAI OAuth Login Repair & Diagnostics**:
  - Diagnosed and resolved the `auth.openai.com` authorization rejection caused by unverified/unregistered client IDs (`client_id=acctg-practice-client`).
  - Supported user-configured or organization OAuth client registration parameters via `OPENAI_OAUTH_CLIENT_ID` env var, `--oauth-client-id` CLI flag, and TUI Tutor Settings modal (`o`).
  - Surfaced rich, actionable error messages and diagnostics to the learner via the local callback server HTML page and TUI notice if OAuth configuration is incomplete or rejected.
  - Ensured zero disruption to offline practice and seamless fallback to `APIKeyTutor` (`OPENAI_API_KEY`) or `OfflineTutor`.
- **Dynamic Multi-Provider Model Discovery & Caching**:
  - Implemented dynamic model catalog discovery for all supported AI providers:
    - **OpenAI**: Query `GET /v1/models` (filtering for chat and reasoning models like `gpt-4o`, `gpt-4o-mini`, `o1`, `o3-mini`).
    - **Anthropic**: Query Anthropic Models API `GET /v1/models` (listing active Claude models like `claude-3-7-sonnet`, `claude-3-5-sonnet`, `claude-3-5-haiku`).
    - **Google Gemini**: Query Gemini API `GET /v1beta/models` (filtering for models supporting `generateContent`).
  - Cached discovered model catalogs locally in user data directory (`models_cache.json`) with timestamps and instant, non-blocking startup.
  - Maintained vetted offline fallback defaults if network or discovery calls fail.
- **Dynamic Model Selection & Switching**:
  - Implemented an interactive model selector within the TUI Tutor Settings modal (`t` -> `m`), allowing learners to browse discovered models, select with Enter or quick-keys `[1-9]`, dynamically change models, or input arbitrary custom model IDs (`c`).
  - Persisted user-selected models per provider in `AuthConfig` (`SelectedModels map[string]string`) and passed dynamically into all subsequent tutor hint/explanation requests.
  - Exposed CLI options for model listing and selection (`--tutor-model=<name>`, `--list-models`, `--fetch-models` / `--refresh-models`).

Gate: PASSED. TUI Tutor Settings modal allows browsing and selecting live models fetched from OpenAI, Anthropic, and Google APIs; discovered models are cached locally; OpenAI OAuth authorization is repaired with explicit registration setup, actionable diagnostics, and graceful fallback; mock tests verify API parsing, caching, and error resilience without external network dependencies.
## Stage 19 — Comprehensive documentation, hyperlinked Table of Contents, and Mac/Wine guide (Complete)

Transform `README.md` from an initial planning starter into an authoritative, publication-quality user hub:
- **Introductory Overview**: Concise, compelling blurb explaining the keyboard-driven, offline-first accounting tutor (**AccountTutor 9000**).
- **Prominent Table of Contents**: Directly following the intro blurb, provide a clean, hyperlinked Table of Contents with working markdown anchor links to all major document sections.
- **Prominent macOS / Wine Execution Guide**:
  - First-class, detailed instructions for running the pre-built Windows binary on macOS using Wine (`brew install --cask wine-stable`, running `wine ./acctg.exe`, terminal font/cursor rendering tips, keyboard mapping notes).
  - Alternative native macOS execution via Go (`brew install go`, `go build -o acctg ./cmd/acctg`, `./acctg`).
  - Cross-compilation instructions for Apple Silicon and Intel Macs (`GOOS=darwin GOARCH=arm64 go build -o acctg-mac ./cmd/acctg`).
- **Feature Tour & Keyboard Cheatsheet**: Exhaustive documentation of all practice modes (Progressive Drill, Journal Entry Practice, Financial Statements, Exam Mode, Tutor Configuration, Candidate Review).
- **CLI Reference & Curriculum**: Complete reference of CLI flags, embedded curriculum accounts, and accounting invariants.

Gate: PASSED. `README.md` fully rewritten with zero placeholder text; all internal markdown links navigate to valid anchors; macOS Wine instructions tested and prominently featured; complete CLI and keyboard cheat sheets match current binary capabilities.

## Stage 20 — Contributor workflow, PR/issue instructions, and release hardening (Complete)

Establish structured collaboration guidelines and finalize quarter maintenance:
- **Pull Request Guidelines & Template**:
  - Add `.github/PULL_REQUEST_TEMPLATE.md` and `CONTRIBUTING.md`.
  - Enforce automated quality gates: `gofmt -s -w .`, `go vet ./...`, `go test ./...`.
  - Guidelines for contributing new accounting families and fixtures (requiring reviewed rules, table-driven semantic tests, distractor tagging, and balanced postings).
  - Invariant enforcement: strictly prevent secrets, API keys, or personal learner attempt history from entering git commits.
- **Issue Filing Guidelines & Templates**:
  - Create issue templates for Accounting Error / Reconciliation Disputes, Feature Requests / Curriculum Scope, and Terminal / Cross-Platform Bugs.
- **Quarter Maintenance, Multi-Platform Packaging, and Release Hardening**:
  - Validate clean setup and offline practice from a fresh repository clone.
  - Native pre-compiled macOS release packaging (`acctg-mac-arm64` for Apple Silicon M1-M4, `acctg-mac-amd64` for Intel Macs, and universal binary).
  - Double-clickable [`Launch-Tutor.command`](Launch-Tutor.command) launcher script allowing Mac classmates to launch without typing terminal commands.
  - Document macOS Gatekeeper quarantine removal (`xattr -d com.apple.quarantine`) and Finder Right-Click override.
  - Retain and maintain Wine execution (`wine ./acctg.exe`) as a tested secondary fallback for users who only receive the Windows binary.
  - Multi-platform smoke test matrix (Windows PowerShell, macOS native ARM64 & Terminal via Wine, Linux).
  - SQLite database migration and upgrade regression verification.
  - Final binary packaging and release notes.

Gate: PASSED. Contribution and issue filing templates in place; PR template enforces test and accounting invariant checklists; clean build from fresh clone verified; zero unresolved lint or vet issues; release packaging complete.

## Stage 21 - Transaction Analysis Grid (implemented; source verification open)

The implemented four-column layout has automated formatting checks. No instructor slides are currently available to substantiate a course-specific match. Historical intended layout:
- Inspect and match the 4-column classroom transaction analysis grid:
  - Column 1: `Dr.` / `Cr.` side indicator.
  - Column 2: Account Name with Equation Effect indicator: `Cash (+A)`, `Loan (+L)`, `Common Stock (+E)`, `Rent Expense (-E)`, `Accounts Payable (-L)`.
  - Column 3: Debit amount (e.g. `$80,000`).
  - Column 4: Credit amount (e.g. `$80,000`).
  - Clean rectangular grid / box borders around all cells.
- Calculate equation effect per posting based on normal balance and direction ($\Delta A$, $\Delta L$, $\Delta E$).
- Update `internal/domain`: add posting equation effect calculations and grid rendering.
- Update `internal/tui/view.go`: upgrade `renderRecap` and `renderJournalPractice` to use the 4-column classroom grid.
- Update `GenericVisualNotice` to reflect verified classroom grid format.

Gate: formatting implementation verified; instructor-slide match remains unverified. Transaction recaps and journal practice render 4-column bordered grids; equation indicators accurately reflect account category and direction; debits and credits align to dedicated columns; unit tests verify formatting and equation tags.

## Stage 22 — Unencumbered Arcade Intro Flight & Combat Override (Complete)

Enhance the startup spaceship arcade animation (`StateIntro`) with persistent manual piloting:
- Auto-pilot and auto-firing remain active by default as demo/attract mode.
- Transition immediately to persistent manual control as soon as user presses a game control key (`W`, `S`, `↑`, `↓`, `F`, `Space`).
- Disable auto-movement and auto-steering once user takes control: ship flies responsive to user input with aerodynamic friction.
- Disable auto-shooting once user takes control: blaster cannons only fire when explicitly triggered by the user (`F` or `Space`).
- Allow the player to play unencumbered indefinitely until they decide to transition to accounting drills (`Enter` or `Esc`).
- Display clear HUD status badges: `PILOT: MANUAL [PLAY]` vs `PILOT: AUTO-PILOT [DEMO]`.

Gate: PASSED. Touching controls permanently disengages auto-pilot and auto-cannons for the session; player flies and shoots with zero automated interference; drill begins cleanly on Enter/Esc; unit tests cover manual transition and state persistence.

## Stage 23 — Complete Linux Distribution & Desktop Integration (Complete)

Provide first-class Linux support matching macOS and Windows:
- Standalone native binaries for both `linux/amd64` and `linux/arm64` (Raspberry Pi, ARM cloud, Asahi Linux, Chromebooks).
- Create `launch-tutor.sh` double-clickable launcher with terminal emulator auto-detection (`gnome-terminal`, `konsole`, `xfce4-terminal`, `xterm`, etc.) and Wine fallback.
- Create standard XDG desktop entry `accounttutor.desktop`.
- Update release packaging scripts (`build_releases.sh`, `build_releases.ps1`) to bundle Linux distributions.
- Expand `README.md` with comprehensive Linux guide (installation, execution, terminal settings, shortcuts, and troubleshooting).

Gate: PASSED. Clean cross-compilation for `linux/amd64` and `linux/arm64`; `launch-tutor.sh` detects environment and architecture; desktop entry conforms to XDG standards; documentation provides complete Linux onboarding.

## Stage 24 — LLM Linkages Verification & OAuth Diagnostics (Complete)

Test, diagnose, and harden external AI tutor connections:
- Update default OpenAI OAuth client ID to `acctg-practice`.
- Provide actionable in-app and browser diagnostics for OpenAI `invalid_client` ("This app is unavailable"): use automatic open-source dynamic registration (supersedes the original developer-registration diagnosis) and offer immediate 1-click fallback to API keys (`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `GEMINI_API_KEY`) or 100% offline mode.
- Add `--test-llm` / `--test-providers` CLI command to verify all configured LLM provider connections and report diagnostic status.
- Table-driven unit tests for all provider adapters and OAuth error scenarios without external network dependencies.

Gate: PASSED. `DefaultDCRClientID` matches registered identifier `acctg-practice`; `--test-llm` reports status across all 5 providers; OAuth errors display actionable resolution guidance; mock tests verify resilience.

## Stage 25 — Startup Arcade Combat Enhancements (Rapid Machine Gun & Smart Bombs) (Complete)

Enhance the startup spaceship arcade game (`StateIntro`) with responsive continuous firing mechanics and smart bomb capabilities:
- **Continuous Machine Gun Fire & Simultaneous Steering**:
  - Support holding down `F` (or `Space`) for rapid machine-gun blaster fire.
  - Implement a rapid-fire burst counter / cooldown in `IntroState` so that pressing or holding `F` initiates sustained blaster bursts across game ticks.
  - Ensure flight steering inputs (`W`, `S`, `↑`, `↓`, `k`, `j`) and laser firing operate seamlessly without cancelling or blocking each other in Bubble Tea's event stream.
- **Smart Bomb Deployment (`B` key)**:
  - Separate bomb deployment from normal laser blasters (reassign `B` / `b` to smart bomb).
  - Track player bomb ordnance stock (e.g. 3 tactical bombs per flight).
  - Trigger expanding ASCII shockwave ring (`IntroShockwave`) spanning the viewport with radial debris and screen flash.
  - Deal massive damage to all active threats on screen (vaporizing standard accounting entities and heavily damaging heavy obstacles).
  - Add combat callouts (`BOMB DETONATED!`, `EMERGENCY AUDIT!`, `TOTAL CLEARANCE!`) and HUD bomb counter display (`BOMBS: 💣💣💣 [B]`).

Gate: PASSED. Holding `F` produces continuous blaster fire; steering during fire does not stall movement or lasers; pressing `B` consumes ordnance, renders an expanding shockwave, damages entities, and awards points; unit tests verify rapid fire cooldown and bomb detonation.

## Stage 26 — Progressive Acceleration & Dual Audit Hit Points (Internal & External Audit HP) (Complete)

Transform the startup animation into a high-stakes, escalating arcade challenge with survival stakes and dual health systems:
- **Progressive Flight Acceleration**:
  - Start at a measured, accessible cruise velocity (`SpeedMultiplier = 1.0`).
  - Progressively accelerate entity speeds, starfield parallax drift, and spawn density over time based on elapsed ticks and player score.
- **Dual Audit Health System**:
  - **Internal Audit Shields (`ShieldHP`, 100%)**: Decreases when incoming accounting entities or obstacles collide directly with the spaceship's fuselage/wings (failed to dodge or blast threats in time).
  - **Global / External Audit Integrity (`GlobalHP`, 100%)**: Decreases whenever an accounting item escapes off the left edge of the screen into the wild without being audited/blasted.
  - Color-coded dual health bars on the arcade HUD:
    `INT AUDIT [SHIELD]: [████████░░] 80%  │  EXT AUDIT [GLOBAL]: [██████████] 100%`
- **Audit Failure & Remediation Flow**:
  - If `ShieldHP <= 0`, trigger "INTERNAL AUDIT FAILURE — DEFICIENT CONTROLS!"
  - If `GlobalHP <= 0`, trigger "EXTERNAL AUDIT FAILURE — ADVERSE OPINION!"
  - Visual explosion sequence with Game Over summary, final score, and prompt to retry flight (`R`) or begin accounting drills (`Enter`/`Esc`).

Gate: PASSED. Target speed and starfield parallax accelerate over time; ship collision reduces internal shield HP; missed entities reduce external audit global HP; reaching 0 HP in either meter triggers audit failure game over screen; unit tests verify damage and acceleration math.

## Stage 27 — Heavy Accounting Obstacles (Fraud & Insider Trading Asteroids) (Complete)

Introduce large, multi-hit accounting hazard obstacles requiring concentrated fire or tactical bombs:
- **Hazard Asteroids Catalog**:
  - Multi-line, high-visibility ASCII obstacles representing severe accounting violations:
    - `[🚨 FRAUD: HP 10/10]` (10 HP)
    - `[⚠️ INSIDER TRADING: HP 12/12]` (12 HP)
    - `[💣 MATERIAL WEAKNESS: HP 8/8]` (8 HP)
    - `[💸 PONZI SCHEME: HP 14/14]` (14 HP)
- **Multi-Hit Durability & Visual Feedback**:
  - Require multiple laser hits or a direct smart bomb detonation to neutralize.
  - Display remaining hazard HP or flashing damage indicator on hit with spark particle emissions.
- **Fragmentation & High-Value Rewards**:
  - When destroyed, break into 2–3 smaller sub-targets (e.g. `[SHREDDED EVIDENCE]`, `[SEC PENALTY]`, `[RESTITUTION]`) that drift and can be blasted for bonus score.
  - Large explosion effect and floating callouts (`SEC INJUNCTION! +1,000`, `CRIMINAL REFERRAL!`).
  - Severe penalty to Internal Audit Shield HP on direct collision, and major hit to External Audit Global HP if allowed to escape off-screen.

Gate: PASSED. Heavy obstacles spawn at higher difficulty intervals; require multiple laser hits or smart bomb to destroy; fragment into sub-targets upon destruction; deal amplified damage on collision or escape; unit tests verify multi-hit damage and fragmentation.

## Stage 28 — Persistent Arcade High Scores & 3-Initials Hall of Fame (SQLite Integration) (Complete)

Preserve arcade performance across sessions in the local SQLite database with retro arcade initials entry:
- **SQLite Database Migration (Migration 6)**:
  - Add `arcade_high_scores` table (`id`, `initials`, `score`, `blasted_count`, `survival_seconds`, `created_at`).
  - Pure-Go implementation with index on `score DESC`.
  - Methods in `internal/storage`: `GetTopArcadeHighScore()`, `SaveArcadeHighScore()`, `ListTopArcadeHighScores()`.
- **Conditional 3-Initials Arcade Entry**:
  - When game concludes (game over or exit to drills), check if the current score beats the all-time high score in the database.
  - If new high score achieved: prompt with retro arcade 3-character initials selector (`[ _ ] [ _ ] [ _ ]`) using `↑`/`↓` or `A-Z` to cycle characters, and `Enter`/`Space` to advance and commit.
  - If score does NOT beat the high score: do NOT prompt for initials; display standard score summary and preserve existing high score.
- **HUD & Leaderboard Integration**:
  - Display all-time high score on HUD: `HIGH: [DAN] 12,450`.
  - Flashing indicator when player surpasses the high score during live flight (`*** NEW ALL-TIME HIGH! ***`).
  - CLI flag `--high-scores` to view the top 10 arcade pilots from the terminal.

Gate: PASSED. Migration 6 applies cleanly on fresh and upgraded databases; high scores persist in SQLite; player is only prompted for initials when setting a new all-time high score; top score renders on arcade HUD; unit tests verify storage operations and initials entry state machine.

## Stage 29 — Local-first LM Studio tutor provider (Complete; live check open)

Add an optional local tutor provider backed by LM Studio's OpenAI-compatible server (default `http://localhost:1234/v1`; `GET /v1/models`, `POST /v1/chat/completions`). Offline mode remains the default and startup never contacts the server. Verify current [LM Studio OpenAI-compatibility docs](https://lmstudio.ai/docs/developer/openai-compat) before implementation.

- **Shared chat-completions client**: Extract the request/decode body of `OpenAIAPITutor.call` in `internal/tutor/apikey.go` into one helper used by both the OpenAI and new `LMStudioTutor` adapters. Add `ProviderLMStudio = "lmstudio"`.
- **No key required**: Send a placeholder bearer token, or `LM_API_TOKEN` when the learner has enabled LM Studio token authentication.
- **Local-model output hygiene**: Strip `<think>…</think>` blocks and ignore `reasoning_content` before `CleanLaTeXMath`. Make `max_tokens` configurable with a higher local default so reasoning models are not truncated mid-thought.
- **Configuration and privacy**: `LMStudioURL` in `AuthConfig`, overridable by `LMSTUDIO_BASE_URL` and `--lmstudio-url`. Accept loopback hosts only (`localhost`, `127.0.0.1`, `::1`) unless the learner explicitly opts in to a remote host, so local-first tutoring keeps learner context on the machine. `IsConfigured` checks only that a URL is set; no network probe.
- **Factory and timeouts**: `BuildTutor` wraps `LMStudioTutor` in `FallbackTutor` with a longer provider default (about 90s) to cover cold model loads and CPU inference. Existing cancellation and offline fallback are unchanged.
- **Model discovery**: `DiscoverLMStudioModels` lists loaded models from `/v1/models` without the OpenAI name filter, hides embedding models, defaults to the first loaded model and reports "no model loaded" clearly. Existing `[m]`/`[r]`/`[c]` selector keys work unchanged.
- **TUI**: Tutor Settings gains `[6] LM Studio (local)` showing URL, active model and a no-key badge, plus URL editing via the existing input buffer. Selecting it performs one asynchronous `/v1/models` check and shows either a connected/model count notice or actionable guidance (start the server in LM Studio's Developer tab or with `lms server start`). The TUI does not launch or manage the LM Studio process.
- **Diagnostics**: Include `lmstudio` in `--test-llm`, `--list-models` and `--fetch-models`, with server-not-running and no-model-loaded messages instead of auth guidance.
- **Documentation and in-app help** (shipped with this stage, not deferred):
  - README `## AI tutor`: an LM Studio setup walkthrough (install, download and load a model, start the server, press `t` then `6`), `LMSTUDIO_BASE_URL`/`LM_API_TOKEN`/`--lmstudio-url`, the loopback-only default and remote opt-in, troubleshooting (server not running, no model loaded, slow first reply from a cold load, reasoning-model output) and a note that local tutoring needs no account or billing. Update the CLI flag list and keyboard reference.
  - Help screen (`renderHelp` in `internal/tui/view.go`): list the LM Studio option beside `[t]` and mention local tutoring in the getting-started text. Tutor Settings modal shows short inline setup guidance when the server is unreachable.
  - OVERVIEW provider options and settled table, `config/example.toml` (`provider = "lmstudio"`, base URL, timeout, max tokens), RELEASE-NOTES unreleased entry, ARCHITECTURE provider list if it enumerates providers, and HANDOFF.
- **Unchanged invariants**: Responses remain untrusted advisory prose with no path to grades, answer keys, bank approval or mastery; chat-only tutoring earns no evidence.

Gate: `httptest` coverage for success, think-tag stripping, no required Authorization header, connection-refused fallback, timeout, cancellation, model-list parsing (including none loaded) and non-loopback rejection; TUI update test for `[6]`; help-screen test asserts the LM Studio option is advertised; startup path makes no network call; gofmt, go vet ./... and go test ./... pass; all documentation above updated and its local links resolve. One manual run against a real LM Studio server with a loaded model is recorded as performed or remains an open live-provider check.

Gate status (2026-10-05): PASSED for automated coverage. httptest, TUI update and help-screen tests pass; gofmt, go vet ./... and go test ./... are clean; documentation is updated and its local links resolve. One deviation from the plan: no placeholder bearer token is sent. The Authorization header is sent only when `LM_API_TOKEN` is set. The manual run against a real LM Studio server remains an open live-provider check.

## Stage 30 — Streamed tutor replies (Planned; after Stage 29)

Show hint and explanation text incrementally as it is generated. This mainly benefits slower local models but applies to all network providers. Verify each provider's current streaming documentation before implementation.

- **Optional streaming contract**: Add a `StreamingTutor` interface (e.g. `HintStream`/`ExplainStream` emitting text deltas and a final `Response`) alongside the existing `Tutor` interface. Providers without streaming keep working through the non-streamed path.
- **Adapters**: Server-sent-event parsing for OpenAI-format chat completions (`stream: true`, `data:` lines, `[DONE]`; shared by OpenAI and LM Studio), Anthropic Messages streaming and Gemini `streamGenerateContent?alt=sse`. The ChatGPT plan adapter already requests `stream: true`; surface its deltas instead of buffering.
- **Nonblocking TUI delivery**: Deliver deltas through a self-reissuing `tea.Cmd` reading from a channel, never blocking the event loop. Throttle re-rendering (glamour markdown, `CleanLaTeXMath` on accumulated text) to a fixed interval rather than every token. Show a "thinking…" indicator while a `<think>` block is open; never display its contents.
- **Timeouts, limits and cancellation**: Replace total-request timeout with first-token and idle timeouts for streams. Enforce `ErrResponseTooLarge` on accumulated size during the stream. Esc/cancel closes the response body promptly.
- **Failure semantics**: Error before the first delta falls back to `OfflineTutor` as today. Error mid-stream keeps the partial text visibly marked incomplete and does not save it as an explanation; only completed responses are persisted. Budget usage is recorded once per completed or aborted request.
- **Documentation and in-app help** (shipped with this stage, not deferred):
  - README `## AI tutor`: what streaming looks like, which providers stream, how to cancel a reply in progress, what an "incomplete" reply means and that it is not saved, and any new flag or setting (e.g. turning streaming off). Update the keyboard reference if cancel/stop keys change.
  - Help screen and status bar (`renderHelp`, `renderStatusBar` in `internal/tui/view.go`): advertise the cancel key while a reply is streaming and explain the "thinking…" and incomplete markers.
  - ARCHITECTURE tutor contract (optional `StreamingTutor`, delivery path, timeout semantics), `config/example.toml`, RELEASE-NOTES unreleased entry and HANDOFF.
- **Unchanged invariants**: Streaming changes presentation only; no path to grades, answer keys, bank approval or mastery.

Gate: fake SSE server tests for chunk parsing per provider, split/partial lines, `[DONE]`, think-block suppression across chunk boundaries, first-token/idle timeouts, mid-stream error, size limit and cancellation; TUI update tests showing key input stays responsive during a slow stream and partial replies are not saved; help-screen/status-bar test asserts the streaming cancel key is advertised; gofmt, go vet ./... and go test ./... pass; all documentation above updated and its local links resolve. Manual streamed run against at least one real provider is recorded or remains an open live-provider check.

## Implementation roadmap allocation

- Stages 00–07: Core offline machine drill, accounting engine, SQLite persistence, and Bubble Tea TUI (Complete).
- Stages 08–11: Bank breadth, adaptive scaffolding, read-only tutor abstraction, and provider integrations (Complete).
- Stages 12–16: Candidate generation/review, journal entry & T-accounts, financial statements, and exam mode (Complete).
- Stage 17: TUI viewport ergonomics and recap scrolling (Complete).
- Stage 18: Provider authentication repair and dynamic multi-provider model discovery (Complete).
- Stage 19: Comprehensive documentation, hyperlinked Table of Contents, and prominent Mac/Wine guide (Complete).
- Stage 20: Contributor workflow, PR/issue instructions, and release hardening (Complete).
- Stage 21: Transaction Analysis Grid (implemented; course-source verification open).
- Stage 22: Unencumbered Arcade Intro Flight & Combat Override (Complete).
- Stage 23: Complete Linux Distribution & Desktop Integration (Complete).
- Stage 24: LLM Linkages Verification & OAuth Diagnostics (Complete).
- Stage 25: Startup Arcade Combat Enhancements — Machine Gun Autofire & Smart Bomb Detonation (Complete).
- Stage 26: Dynamic Flight Acceleration & Dual Audit Hit Points (Internal Audit Shields & External Audit Integrity) (Complete).
- Stage 27: Heavy Accounting Hazards — Multi-Hit Fraud & Insider Trading Asteroids with Fragmentation Debris (Complete).
- Stage 28: Persistent Arcade High Scores & Old-School 3-Initials Hall of Fame (SQLite Schema Migration & TUI Entry) (Complete).
- Stage 29: Local-first LM Studio tutor provider (Complete; live LM Studio run open).
- Stage 30: Streamed tutor replies across providers (Planned; after Stage 29).

## Current quality gate at 90 - October 3, 2026

Bank breadth is 90 active scenarios, one retired record, 13 existing families. Integration and teaching delivery, immutable exam replay, focused offline explanations, guided matched contrasts, distinct-setting delayed retrieval, priority misconception choices, amount-scope choices and generic payable cleanup are covered by the checks recorded in [QUALITY-REVIEW90.md](docs/QUALITY-REVIEW90.md). Current pedagogy policy is 2; evidence projection is 2; schema migration is 9.

This closes the focused prerequisites before additional content breadth. The next content work, when requested, is a reviewed five-scenario batch within existing supported scope, following [CONTENT-EXPANSION-PLAN.md](docs/CONTENT-EXPANSION-PLAN.md). The 180 target, actual syllabus/slides, native macOS/Linux smoke checks, live-provider verification and empirical learning effectiveness remain open. Preserve existing history, template versions and delegated review provenance.

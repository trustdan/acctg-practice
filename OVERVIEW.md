# Situational overview

## Goal and learner context

Build a local, keyboard-driven deliberate practice environment for Financial Accounting (**AccountTutor 9000**). The immediate goal is useful catch-up practice within one week; the architecture should support ongoing learning without requiring a full accounting simulator first.

The learner prefers Go and model-agnostic tooling, trivial dollar amounts, progressive questions, Socratic hints, and persistent mastery with time-based forgetting. The hardest distinction is cash movement versus income, alongside debit/credit direction. The course currently uses journal entries, T-accounts, a spreadsheet, and a reported two-row, three-column debit/credit visual. The precise visual is unverified. Exercises begin with a fact pattern and ask for accounts, categories, direction, and debit/credit.

## Product decisions

Use Go + Bubble Tea for the TUI and SQLite for local progress. Choose maintained dependency versions during implementation and pin them. The offline engine owns accounting rules, answer keys, grading, selection, history, and mastery. An optional AI tutor explains, offers Socratic hints, and proposes question candidates. Machine mode must remain fully functional without a provider.

The reusable master bank contains semantic transaction families, vetted scenario variants, parameter bounds, concept tags, and provenance. Machine mode instantiates known families with deterministic randomization; creative mode proposes expansions. Candidate generation never grants authority to grade. Exam mode is a later extension with feedback withheld until the end.

AI creates. Rules validate. Machine drills. History adapts.

## Learning flow

Economic event → account identity → category → increase/decrease → debit/credit → paired account → balanced entry → statement effects. Present one decision at a time. Fade scaffolding only when unassisted performance supports it, and restore scaffolding when performance declines.

Distinguish cash received for work earned now, customer advances for future work, revenue earned on credit, collection of a previously recorded receivable, borrowing, and owner investment. These all separate receipt of cash from earning revenue. Contrast problems must explicitly state whether revenue was previously recognized.

## Why the separation matters

The model should not be the answer key. Balanced debits and credits alone do not prove a transaction is correct: Dr Cash / Cr Revenue can balance while being wrong for a customer advance. Semantic event rules and reviewed wording must support the answer. Successful learner attempts do not validate generated content.

Likewise, a percentage displayed as mastery is a scheduling estimate, not a clinical measure or validated probability of retention. Track concept evidence separately from exposure to particular question wording. Hinted retries are useful learning but are not independent evidence of mastery.

## What is settled and what is open

| Settled preference | Requires repository verification |
|---|---|
| Go, TUI, SQLite, offline-first | Go/Bubble Tea/SQLite dependency versions and target OS |
| Machine and creative modes; offline fallback | Course vocabulary and actual Week 1–2 scope |
| Optional tutor: ChatGPT Plus (OAuth), API keys (Anthropic, Google, OpenAI) or local LM Studio (loopback only by default) | Exact professor grid and statement formats |
| Dynamic live model selection across providers (not frozen snapshots) | Provider model catalog API deprecations/rate limits |
| Canonical reusable question bank | Exam rules and permitted aids |
| Socratic progressive drills; contrast scenarios | Half-life tuning from observed performance |
| Time-decayed concept scheduling | |
| Viewport vertical scrolling on recap and tall views | Terminal dimension edge cases (< 15 rows) |

Provider options for the optional AI tutor accommodate both subscription and direct-key workflows:
1. **ChatGPT Plus / Subscription**: OpenAI's official "Sign in with ChatGPT" flow for personal/open-source apps using PKCE, local loopback callback, registered client ID configuration, and streamed Responses API (`store: false`, `stream: true`).
2. **User API Keys**: Direct adapters for Anthropic Claude (`ANTHROPIC_API_KEY`), Google Gemini (`GEMINI_API_KEY`), or OpenAI API (`OPENAI_API_KEY`), loaded via environment variables or interactive TUI entry.
3. **Local model (LM Studio)**: An OpenAI-compatible adapter for a model running in LM Studio on the learner's own machine. No key, account or billing; only loopback hosts are accepted unless the learner explicitly opts in to a remote server. OpenAI API and LM Studio replies stream in as they are written; a reply stopped early is shown as incomplete and never saved.
4. **Dynamic Model Discovery**: Live queries to provider model catalog endpoints (`/v1/models` and `/v1beta/models`) with local disk caching and offline fallbacks, allowing learners to select current model releases over time.

All credentials and tokens reside securely in local user-data storage and are never committed to git. Machine-mode drills remain strictly offline, and the application will never silently switch between subscription allowance and paid API billing.

## Initial release and later growth

The application delivers offline progressive drills, reviewed seed families, immediate feedback, attempts saved durably, concept scheduling, and a readable mastery screen. Subsequent extensions added provider integration, creative candidate generation/review, full journal-entry input, T-account rendering, statement construction, and exam mode. 

The application now includes scrolling recaps, provider selection/model discovery, contributor templates, platform launchers and the startup arcade with persistent scores. Code and automated verification do not establish native-platform or live-provider release validation; those checks remain open.

The embedded bank currently has 90 active reviewed scenarios, one retired record and 13 supported families. The integration pass fixed journal amount/navigation consistency and immutable exam resume, and made bank-backed offline help use focused reviewed explanations. Pedagogy policy 2 adds reviewed setting groups, 13 guided comparison pairs, six priority misconception profiles, ten amount-scope profiles and contextual account choices. Evidence version 2 requires independent success across settings plus delayed transfer before reducing guidance. Existing history remains replayable; historical attempts without reviewed metadata receive no inferred transfer credit.

The focused quality prerequisites at 90 are complete. Further breadth remains a separately reviewed sequence toward 180; no additional scenarios are part of this pass. See [quality closeout](docs/QUALITY-REVIEW90.md), [expansion plan](docs/CONTENT-EXPANSION-PLAN.md) and [handoff](docs/HANDOFF.md). Actual syllabus/slides, course-specific visual verification and empirical learning effectiveness remain unresolved.

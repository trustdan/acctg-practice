# Situational overview

## Goal and learner context

Build a local, keyboard-driven financial-accounting practice environment for UW Foster ACCTG 502. The immediate goal is useful catch-up practice within one week; the architecture should support the quarter without requiring a full accounting simulator first.

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
| Optional tutor: ChatGPT Plus (OAuth) or API keys (Anthropic, Google, OpenAI) | Exact professor grid and statement formats |
| Canonical reusable question bank | Exam rules and permitted aids |
| Socratic progressive drills; contrast scenarios | Half-life tuning from observed performance |
| Time-decayed concept scheduling | |

Provider options for the optional AI tutor in Stage 11 accommodate both subscription and direct-key workflows:
1. **ChatGPT Plus / Subscription**: OpenAI's official "Sign in with ChatGPT" flow for personal/open-source apps (token sharing) using DCR, PKCE, local loopback callback, and streamed Responses API (`store: false`, `stream: true`). No API key or client secret required.
2. **User API Keys**: Direct adapters for Anthropic Claude (`ANTHROPIC_API_KEY`), Google Gemini (`GEMINI_API_KEY`), or OpenAI API (`OPENAI_API_KEY`), loaded via environment variables or interactive TUI entry.

All credentials and tokens reside securely in local user-data storage and are never committed to git. Machine-mode drills remain strictly offline, and the application will never silently switch between subscription allowance and paid API billing.

## Initial release and later growth

The first release delivers offline progressive drills, reviewed seed families, immediate feedback, attempts saved durably, concept scheduling, and a readable mastery screen. AI integration, creative candidate review, full journal-entry input, T-account rendering, statement construction, and exam mode follow that release. A learner can begin practicing before the whole roadmap is complete.

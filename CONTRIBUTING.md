# Contributing to AccountTutor 9000

Thank you for your interest in contributing to **AccountTutor 9000 (`accountutor-9000` / `acctg`)**! This project is an open-source, keyboard-driven terminal application engineered for students and professionals learning **Financial Accounting** to build rapid accounting reflexes through deliberate practice.

Our guiding architectural motto:
> **"AI creates. Rules validate. Machine drills. History adapts."**

Before contributing, please review this guide and our canonical governance document, [docs/AGENT-CONTRACT.md](file:///c:/Users/Dan/acctg-practice/docs/AGENT-CONTRACT.md).

---

## Table of Contents

- [Code of Conduct & Pedagogy](#code-of-conduct--pedagogy)
- [Core Invariants](#core-invariants)
- [Development Setup](#development-setup)
- [Repository Tour](#repository-tour)
- [How to Contribute a New Transaction Family](#how-to-contribute-a-new-transaction-family)
- [How to Contribute Question Templates](#how-to-contribute-question-templates)
- [Pull Request & Review Process](#pull-request--review-process)
- [Hygiene & Security](#hygiene--security)

---

## Code of Conduct & Pedagogy

We are dedicated to building a supportive, precision-oriented learning tool for accounting students. When contributing code, question templates, hints, or explanations:

1. **Precision Grounded in Economic Reality**: Explanations must start from what actually happened economically (e.g. cash was received today for work to be performed next month; cash paid today for 12 months of future insurance benefits).
2. **Pedagogical Invariants**:
   - **Debit = Left, Credit = Right**. Never define debit as "money in" or "good", or credit as "money out" or "bad". Directional impact depends entirely on account classification ($A = L + E$).
   - **Cash is an Asset**: Revenue is an earned equity inflow, not cash. Cash increases with a debit; revenue increases with a credit.
   - **Accrual Timing Decoupled**: Receipts and payments must be explicitly separated from revenue recognition and expense incurrence.
3. **Encouraging, Socratic Tone**: Explanations must never shame the learner, condescend, or diagnose cognitive retention as a medical or psychological fact. Hints should guide causal deduction.

---

## Core Invariants

All contributions must preserve these non-negotiable architectural invariants:

1. **100% Offline Default**: The core machine practice engine, drill generator, journal practice, financial statements, exam mode, and local storage have **zero network dependencies**. No cloud login, API key, or provider is required at startup.
2. **Integer Minor Units**: All monetary amounts are stored and calculated as 64-bit integer minor units (cents). Floating-point arithmetic (`float32`, `float64`) for money is strictly prohibited.
3. **Strict Entry Balancing**: Every journal entry must mathematically balance ($\sum \text{Debits} == \sum \text{Credits}$). Every transaction must reconcile across Journal $\equiv$ T-Accounts $\equiv$ Accounting Equation ($\Delta A = \Delta L + \Delta E$).
4. **Deterministic Accounting Engine**: Answer keys, grading, and trial balances are derived solely by the deterministic accounting engine (`internal/engine`). AI models and LLMs are treated as **untrusted, read-only advisory tools** and can **never** mutate answer keys, grading logic, or learner mastery.
5. **Durable, Replayable History**: Attempts are persisted immutably with question versions, seed, parameters, option order, assistance level, and grading version. Bank updates must never invalidate or rewrite historical learner attempts.
6. **Isolated Assessment Evidence**: Exam attempts (`exam_attempts`) are strictly segregated from daily drill learning attempts (`attempts`). Mastery projections and forgetting curves are unpolluted by exam sessions.

---

## Development Setup

### Prerequisites
- **Go 1.24+** (tested on Go 1.24 and Go 1.27 on Windows, macOS, and Linux).
- Optional: **Wine** on macOS/Linux if testing the Windows `.exe` fallback.

### Building and Running
```bash
# Clone the repository
git clone https://github.com/trustdan/acctg-practice.git
cd acctg-practice

# Run directly from source
go run ./cmd/acctg

# Build standalone executable
go build -o acctg ./cmd/acctg

# Run unit tests across all packages
go test -v ./...

# Run code formatter
gofmt -s -w .

# Run static analysis
go vet ./...
```

---

## Repository Tour

```
acctg-practice/
├── cmd/acctg/                 # Application entry point & CLI flag dispatcher
├── curriculum/                # Canonical course curriculum data
│   ├── accounts.json          # 14-account standard chart of accounts
│   └── seed-questions.json    # 19 reviewed question templates (18 active, 1 retired)
├── internal/
│   ├── bank/                  # Bank validation, family catalog, semantic wording checks
│   ├── candidate/             # Creative candidate generation, preview, & review workflow
│   ├── domain/                # Integer money, accounts, journal entries, T-accounts, reconciliation
│   ├── drill/                 # Progressive drill state machine, adaptive scaffolding, hints
│   ├── engine/                # Deterministic accounting rules, entry evaluation, distractor tags
│   ├── exam/                  # Exam mode runner, timer, report generator, error pattern analysis
│   ├── mastery/               # Bayesian Beta evidence, forgetting curves, delayed retrieval scheduler
│   ├── statements/            # Financial statements (BS, IS, RE, SCF), accounting cycle, audit proofs
│   ├── storage/               # Pure-Go SQLite persistence, versioned migrations (v1-v5), backups
│   ├── tui/                   # Interactive Bubble Tea terminal UI, viewports, modal overlays
│   └── tutor/                 # Read-only tutor adapters (Offline, ChatGPT OAuth, Claude, Gemini, OpenAI)
├── docs/                      # Authoritative specifications and handoff logs
│   ├── AGENT-CONTRACT.md      # Canonical shared contract and non-negotiable rules
│   ├── ARCHITECTURE.md        # Technical architecture documentation
│   └── HANDOFF.md             # Chronological development handoff and verification log
├── scripts/                   # Build and release automation scripts
└── Launch-Tutor.command       # Double-clickable macOS Finder launcher
```

---

## How to Contribute a New Transaction Family

Adding a new transaction family (e.g. `bad_debt_allowance`, `bond_issuance`, `inventory_purchase`) expands the deterministic accounting engine and curriculum bank. Follow these 6 steps:

### 1. Verify / Add Accounts
Check `curriculum/accounts.json`. If a new account is needed, add it with:
- `id`: unique snake_case string (e.g. `allowance_for_doubtful_accounts`)
- `name`: standard textbook name (e.g. "Allowance for Doubtful Accounts")
- `category`: one of `asset`, `liability`, `equity`, `revenue`, `expense`, `dividend`
- `normal_side`: `debit` or `credit`
- `contra`: boolean (e.g. `true` for contra-assets like accumulated depreciation or allowance)

### 2. Register Family in Bank
In `internal/bank/family.go`:
- Add the family constant (e.g. `FamilyBadDebtAllowance = "bad_debt_allowance"`).
- Register the family in `SupportedFamilies()` with its required and optional parameters.
- Define its reviewed rule version in `ReviewedFamilyRuleVersion()`.

### 3. Implement Deterministic Rule in Engine
In `internal/engine/rules.go` (and `internal/engine/event.go`):
- Implement the rule struct satisfying the `Rule` interface.
- Return the canonical debit/credit postings using integer minor units.
- Derive the balance sheet equation delta ($\Delta A = \Delta L + \Delta E$).
- Implement distractor tagging: identify common student conceptual mistakes (e.g. crediting Accounts Receivable directly under the allowance method) and assign a diagnostic error tag.

### 4. Write Table-Driven Engine Tests
In `internal/engine/engine_test.go`:
- Add table-driven tests verifying:
  - Canonical parameter combinations produce strictly balanced entries.
  - Balance sheet equation delta matches the postings.
  - Common wrong entries are flagged with the expected distractor error tags.

### 5. Add Progressive Drill Support
In `internal/drill/generator.go`:
- Ensure the family generates all required progressive drill stages:
  - Account identification
  - Account classification
  - Increase / decrease direction
  - Debit / credit assignment
  - Paired balancing account
  - Balanced journal entry summary
  - Balance sheet equation delta recap

### 6. Add Reviewed Seed Question
In `curriculum/seed-questions.json`:
- Create a clear, pedagogically grounded scenario template with unambiguous timing (explicit "today", "last month", "next month", or "12 future months").
- Set `reviewer` and `approved_at` timestamp.
- Verify using:
  ```bash
  go test -v ./internal/bank/...
  ./acctg.exe --reconcile-all --db :memory:
  ```

---

## How to Contribute Question Templates

You can propose question templates either through pull requests to `curriculum/seed-questions.json` or by generating and reviewing them via the CLI/TUI:

1. **Generate Candidate**:
   ```bash
   ./acctg.exe --generate-candidate --candidate-family=customer_advance
   ```
2. **Preview & Validate**:
   ```bash
   ./acctg.exe --preview-candidate=<id>
   ```
3. **Approve & Promote**:
   ```bash
   ./acctg.exe --approve-candidate=<id> --reviewer="YourName" --notes="Reviewed for Week 2"
   ```

All question candidates must pass semantic wording checks (`internal/bank/semantic.go`) to ensure balanced-but-conceptually-wrong scenarios (such as claiming revenue is earned immediately on an unearned advance) are caught and rejected.

---

## Pull Request & Review Process

1. **Fork & Branch**: Create a feature branch (e.g. `feature/bad-debt-family` or `fix/recap-scroll-indicator`).
2. **Run Quality Gates**:
   ```bash
   gofmt -s -w .
   go vet ./...
   go test -v ./...
   ```
   All tests must pass with 0 failures and 0 external network requests.
3. **Submit PR**: Open a pull request using our [Pull Request Template](file:///c:/Users/Dan/acctg-practice/.github/PULL_REQUEST_TEMPLATE.md). Complete the Invariant & Quality Gate Checklist.
4. **Code Review**: A maintainer will review the PR for accounting soundness, test coverage, and invariant compliance.

---

## Hygiene & Security

- **NEVER Commit Secrets**: Never commit API keys (`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `GEMINI_API_KEY`), OAuth tokens, or client secrets. All tutor authentication is stored locally in user configuration directories (`0600` permissions) outside the git working tree.
- **NEVER Commit Learner Progress**: Never commit personal SQLite database files (`*.db`, `*.db-wal`, `*.db-shm`) or attempt history. Check `git status` before committing.

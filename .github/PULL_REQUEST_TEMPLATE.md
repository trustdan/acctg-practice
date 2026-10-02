## Pull Request Summary

<!-- Provide a concise summary of the changes, the problem solved, and the relevant issue or stage. -->

### Changes Proposed
- 

### Related Issue / Milestone
- Resolves #<!-- issue number --> or implements Stage <!-- stage number -->

---

### Invariant & Quality Gate Checklist

Before submitting this pull request, please verify that your changes adhere to all non-negotiable invariants established in [docs/AGENT-CONTRACT.md](file:///c:/Users/Dan/acctg-practice/docs/AGENT-CONTRACT.md) and [CONTRIBUTING.md](file:///c:/Users/Dan/acctg-practice/CONTRIBUTING.md):

#### 1. Accounting Invariants
- [ ] **Integer Minor Units**: All monetary calculations, ledger entries, and statements use integer minor units (cents). Floating-point currency arithmetic is strictly prohibited.
- [ ] **Entry Balancing**: All journal entries mathematically balance ($\sum \text{Debits} == \sum \text{Credits}$).
- [ ] **Normal Balances**: Debit/Credit directions strictly conform to account classifications ($A = L + E$). Debit increases Assets, Expenses, and Dividends; Credit increases Liabilities, Stockholders' Equity, and Revenues.
- [ ] **Equation Delta Reconciles**: Transaction postings tie out directly to balance sheet equation deltas ($\Delta A = \Delta L + \Delta E$).
- [ ] **Accrual Timing Decoupled from Cash**: Cash receipts are strictly separated from revenue recognition (e.g. customer advance vs. service revenue); cash disbursements are strictly separated from expense incurrence (e.g. prepaid asset vs. operating expense).
- [ ] **Dividend Distinction**: Declared dividends (reducing equity and creating `Dividends Payable` liability with zero cash flow) are handled distinctly from paid dividends (reducing liability and cash with zero equity effect).
- [ ] **Distractor Error Tagging**: New or modified question options diagnose underlying student misconceptions with canonical error tags (e.g. `TagDuplicateRevenueOnCollection`, `TagExpenseRecordedOnEquipmentPurchase`).

#### 2. Deterministic Engine & AI Boundary
- [ ] **Deterministic Rules**: Answer keys, grading, trial balances, and mastery progression are derived strictly by the offline engine (`internal/engine`).
- [ ] **Read-Only AI Output**: Tutor and LLM responses are treated as untrusted, advisory prose. No AI generation path has permission to mutate grading, answer keys, or learner mastery.
- [ ] **Offline Machine Mode**: The default application starts up and functions 100% offline with zero network calls, zero provider requirements, and zero cloud dependencies.

#### 3. Code Quality & Automated Checks
- [ ] Code formatted with `gofmt -s -w .` (zero diff).
- [ ] Static analysis passes with `go vet ./...` (0 errors or warnings).
- [ ] Full test suite passes: `go test -v ./...` (100% passing across all packages).
- [ ] Zero external network calls in unit tests (all external endpoints mocked locally).

#### 4. Privacy, Security & Hygiene
- [ ] **Zero Secrets**: No API keys, OAuth client secrets, or private tokens committed to git.
- [ ] **Zero Learner Data**: No personal attempt history, SQLite database files (`*.db`, `*.db-wal`, `*.db-shm`), or progress files included in the PR.
- [ ] Documentation updated (`README.md`, `PLAN.md`, `docs/HANDOFF.md`) if CLI flags, keyboard controls, or conventions were added or modified.

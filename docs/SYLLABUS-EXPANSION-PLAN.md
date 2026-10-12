# Syllabus Expansion & Interleaved Breadth Plan (Modules 3–10)

## Executive Summary

This plan defines the architectural foundation and content pipeline to expand **AccountTutor 9000** from its current baseline (Modules 1, 2, and 4: basic double-entry mechanics, T-accounts, and direct financial statements) across the full University of Washington **Financial Reporting and Analysis (Foster School of Business)** syllabus (Modules 3–10).

Rather than writing all Module 5 questions before touching Module 6 or 7, this plan executes a **breadth-first, interleaved strategy**:
1. **Build out architectural foundations first** (contra-accounts, bond liabilities, intangibles, analytical ratio evaluations, and indirect cash flow models).
2. **Deploy question batches of 4**, drawing across different modules in each batch (e.g., one from Module 5, one from Module 7, one from Module 6, one from Module 3), so learners gain immediate practice across the full syllabus early on, even when each topic has an initial single question.

---

## Architecture Phase: Building the Engine Foundations

Before questions can be admitted to the bank, the offline deterministic engine must be updated to support the necessary accounts, contra-directions, and evaluation rules.

### 1. Account Catalog & Contra-Account Semantics (`curriculum/accounts.json`, `internal/domain`)
Currently, the catalog supports 14 standard accounts without explicit contra-account handling. The catalog must be extended with:
- **Contra-Assets:**
  - `allowance_for_doubtful_accounts` (Category: `asset`, Normal side: `credit`, `is_contra: true`)
- **Contra-Liabilities & Adjuncts:**
  - `discount_on_bonds_payable` (Category: `liability`, Normal side: `debit`, `is_contra: true`)
  - `premium_on_bonds_payable` (Category: `liability`, Normal side: `credit`, `is_adjunct: true`)
- **Long-Term Liabilities & Equity:**
  - `bonds_payable` (Category: `liability`, Normal side: `credit`)
  - `treasury_stock` (Category: `equity`, Normal side: `debit`, `is_contra: true`)
  - `additional_paid_in_capital` (Category: `equity`, Normal side: `credit`)
- **Intangibles & Operating Expenses/Revenues:**
  - `goodwill` (Category: `asset`, Normal side: `debit`)
  - `patents` (Category: `asset`, Normal side: `debit`)
  - `capitalized_software_costs` (Category: `asset`, Normal side: `debit`)
  - `bad_debt_expense` (Category: `expense`, Normal side: `debit`)
  - `interest_expense` (Category: `expense`, Normal side: `debit`)
  - `gain_loss_on_bond_retirement` (Category: `gain_loss`, Normal side: `credit`/`debit`)

### 2. Deterministic Rules Engine Expansion (`internal/bank/family.go`, `internal/engine/rules.go`)
Register and implement deterministic rules with integer minor-unit math and distractor tags:

| Module | Semantic Family ID | Debit Account(s) | Credit Account(s) | Key Pedagogical Misconception Tagged |
|---|---|---|---|---|
| **Mod 5** | `bad_debt_allowance` | Bad Debt Expense | Allowance for Doubtful Accts | `TagDirectWriteOffExpense` (crediting A/R directly) |
| **Mod 5** | `write_off_receivable` | Allowance for Doubtful Accts | Accounts Receivable | `TagWriteOffExpensedAgain` (debiting Bad Debt Expense) |
| **Mod 7** | `issue_bonds_par` | Cash | Bonds Payable | `TagStockIssuedInsteadOfBond` |
| **Mod 7** | `issue_bonds_discount` | Cash, Discount on Bonds | Bonds Payable | `TagDiscountRecordedAsLoss` |
| **Mod 7** | `issue_bonds_premium` | Cash | Bonds Payable, Premium on Bonds | `TagPremiumRecordedAsGain` |
| **Mod 7** | `pay_bond_interest` | Interest Expense | Cash, [Discount / Premium] | `TagIgnoringBondAmortization` |
| **Mod 6** | `acquire_goodwill` | Tangible Assets, Goodwill | Liabilities Assumed, Cash | `TagGoodwillExpensed` |
| **Mod 6** | `repurchase_treasury_stock` | Treasury Stock | Cash | `TagTreasuryStockRecordedAsAsset` |
| **Mod 9** | `capitalize_software_feasibility`| Capitalized Software Costs | Cash / Accounts Payable | `TagExpensingFeasibleSoftware` |

### 3. Indirect Cash Flow Operating Adjustments (`internal/statements`)
Support the Module 3 Topic 2 indirect cash flow rules:
- Starting point: **Net Income**.
- Add non-cash expenses: $+$ Depreciation / Bad Debt.
- Working capital adjustments:
  - Increases in current operating assets $\rightarrow$ subtract.
  - Decreases in current operating assets $\rightarrow$ add.
  - Increases in current operating liabilities $\rightarrow$ add.
  - Decreases in current operating liabilities $\rightarrow$ subtract.

### 4. Ratio & Analytical Evaluation Framework
Support calculation-based and conceptual multiple-choice drills for:
- **Module 3:** Common-size statements (% of revenue / % of assets), Current Ratio ($\text{Current Assets} / \text{Current Liabilities}$).
- **Module 8:** Financial Leverage Ratio ($\text{Total Assets} / \text{Total Equity}$ or $\text{Total Debt} / \text{Total Equity}$, Lehman Repo 105 off-balance-sheet impact).
- **Module 6 & 9:** Pro-forma / Non-GAAP adjustments (Groupon) and Aggressive vs. Conservative revenue/expense timing (Apple, Microsoft).

---

## Interleaved Batch Strategy: Batches of 4 Across Diverse Modules

To achieve early syllabus coverage, question generation is structured into **interleaved batches of 4**. Each batch takes exactly one scenario from 4 different modules.

### Batch 1: Core Extended Mechanics (First Broad Slice)
- **Question 1.1 (Module 5):** Year-end adjusting entry for estimated uncollectible accounts using the allowance method ($5,000 bad debt expense).
- **Question 1.2 (Module 7):** Disney-style long-term bond issuance at face value (par) ($100,000 cash for 10-year notes).
- **Question 1.3 (Module 6):** Facebook/WhatsApp style business acquisition recording acquired identifiable assets, liabilities, and residual Goodwill.
- **Question 1.4 (Module 3):** Statement of Cash Flows: indirect operating adjustment for an increase in Accounts Receivable (deduct from Net Income).

### Batch 2: Contra-Accounts & Analytical Ratios
- **Question 2.1 (Module 5):** Specific customer account write-off under the allowance method ($1,200), demonstrating that Net Accounts Receivable remains unchanged.
- **Question 2.2 (Module 7):** Bond issuance at a discount (market rate > coupon rate) with Discount on Bonds Payable.
- **Question 2.3 (Module 9):** Microsoft-style software development cost treatment: capitalizing post-technological feasibility coding costs vs. expensing pre-feasibility research.
- **Question 2.4 (Module 3 / 8):** Calculating Current Ratio and assessing liquidity for a retailer (Costco vs. Nordstrom contrast).

### Batch 3: Equity, Interest & Off-Balance Sheet Mechanics
- **Question 3.1 (Module 6):** Repurchase of company common stock into Treasury Stock (contra-equity, cash outflow).
- **Question 3.2 (Module 7):** Semi-annual bond coupon payment and discount amortization.
- **Question 3.3 (Module 8):** Lehman Brothers Repo 105 case analysis: classifying short-term repurchase agreements as sales vs. collateralized borrowings to manage leverage.
- **Question 3.4 (Module 5):** Apple case: Point-of-sale hardware revenue vs. deferred subscription revenue for multi-element software updates.

### Batch 4: Advanced Policy, Intangibles & Forecasting
- **Question 4.1 (Module 6):** Goodwill impairment testing: recognizing impairment loss when fair value drops below carrying value.
- **Question 4.2 (Module 8):** Cryptocurrency accounting (Tesla Bitcoin): classifying digital assets as indefinite-lived intangibles and recording impairment vs. fair value gains.
- **Question 4.3 (Module 9):** Identifying aggressive vs. conservative financial reporting choices (revenue recognition acceleration vs. reserve buffering).
- **Question 4.4 (Module 10):** Financial statement forecasting: distinguishing aggregated top-line forecasts from disaggregated unit-level assumptions.

---

## Acceptance Gates per Batch

Each batch of 4 must satisfy the following checks before being approved into the active bank:
1. **Mathematical & Double-Entry Integrity:** Every journal entry balances to the penny ($\sum \text{Debits} == \sum \text{Credits}$).
2. **Deterministic Grading:** Answers are derived by engine rules, not by language model prompting.
3. **Socratic Distractor Coverage:** Every question includes targeted causal hints for anticipated student misconceptions (e.g., expensing a write-off twice or treating treasury stock as an asset).
4. **Replay & Test Suite Pass:** `./acctg.exe --reconcile-all --db :memory:` and `go test ./...` pass with 100% compliance.
5. **Review Provenance:** All 4 questions contain an explicit reviewer, approved timestamp, and immutable schema version.

---

## Progress Tracker

| Batch | Module Scope | Status | Notes |
|:---:|---|:---:|---|
| **Architecture** | Mod 3, 5, 6, 7, 8, 9, 10 | **Completed** | Added 11 new accounts (25 total), contra-semantics, and engine rules |
| **Batch 1** | Mod 5 (Bad Debt), Mod 7 (Bonds Par), Mod 6 (Goodwill), Mod 3 (Indirect OCF) | **Completed** | First interleaved cross-module slice (4 active questions, 94 total) |
| **Batch 2** | Mod 5 (Write-Off), Mod 7 (Bond Discount), Mod 9 (Software R&D), Mod 3/8 (Current Ratio) | **Completed** | Second cross-module slice (4 active questions, 98 total) |
| **Batch 3** | Mod 6 (Treasury Stock), Mod 7 (Bond Interest), Mod 8 (Lehman Repo), Mod 5 (Apple Rev) | **Completed** | Third cross-module slice (4 active questions, 102 total) |
| **Batch 4** | Mod 6 (Goodwill Impairment), Mod 8 (Crypto), Mod 9 (Aggressive/Conservative), Mod 10 (Forecasting) | **Completed** | Fourth cross-module slice (4 active questions, 106 total); full 10-module interleaved syllabus coverage achieved |

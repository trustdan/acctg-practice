# Stream-M08 Accounting & Pedagogical Review

**Reviewer:** Worker Agent 3 (Stream-M08: Module 8 Leases ASC 842 & Stockholders' Equity / Treasury Stock)  
**Date:** 2026-10-12  
**Candidate Batch:** `curriculum/candidates/batch_m08.json`  
**Candidate Policy:** `curriculum/candidates/pedagogy_m08.json`  
**Harness Tests:** `internal/drill/syllabus_m08_test.go`  

---

## 1. Syllabus & Topic Scope

Stream-M08 targets critical curriculum areas from the University of Washington Foster School of Business *Financial Reporting and Analysis* course for **Module 8: Leverage, Off-Balance Sheet Financing & Capital Structure** (with secondary ties to **Module 2: Double-Entry Mechanics & Accrual Cycle**).

| Question ID | Scenario & Accounting Focus | US GAAP / Standard Citation | Key Concepts |
|---|---|---|---|
| `q_m08_lse_commencement_rou_asset_01` | Initial Lease Commencement: Capitalization of Right-of-Use Asset and Lease Liability at present value of future lease payments. | ASC 842-20-25-1 / ASC 842-20-30-1 | `lease_commencement`, `rou_asset_recognition`, `lease_liability_asc842`, `present_value_measurement` |
| `q_m08_lse_periodic_lease_payment_02` | Periodic Lease Payment: Cash disbursement split between principal lease liability reduction and financing interest expense. | ASC 842-20-35-1 / ASC 842-20-35-4 | `periodic_lease_payment`, `lease_liability_amortization`, `interest_expense`, `asc_842_subsequent_measurement` |
| `q_m08_lse_treasury_stock_repurchase_03` | Treasury Stock Repurchase: Share buyback under the cost method recorded as a contra-equity deduction. | ASC 505-30-30-1 (Treasury Stock Cost Method) | `treasury_stock_cost_method`, `share_repurchases`, `contra_equity`, `cash_outflow_financing` |
| `q_m08_lse_treasury_stock_reissuance_04` | Treasury Stock Reissuance Above Cost: Reissue proceeds in excess of cost credited to Additional Paid-in Capital (no P&L gain). | ASC 505-30-25-2 / ASC 505-30-30-2 | `treasury_stock_reissuance`, `apic_treasury_stock`, `no_gain_or_loss_on_treasury`, `stockholders_equity` |

---

## 2. Technical Accounting Analysis

### A. ASC 842 Lease Capitalization Framework
- **Historical vs. Current Standard:** Prior to ASC 842, operating leases were kept off the balance sheet, permitting companies to incur substantial debt-like obligations without reflecting leverage on the balance sheet.
- **ASC 842 Rule:** Under ASC 842, all leases with terms exceeding 12 months must be capitalized onto the balance sheet at commencement:
  - **Right-of-Use (ROU) Asset:** Represents the lessee's right to control the specified asset over the lease term. Measured at the present value of lease payments plus initial direct costs.
  - **Lease Liability:** Represents the obligation to make future lease payments, discounted using the rate implicit in the lease or the lessee's incremental borrowing rate.
  - **Commencement Entry:**
    $$\text{Debit: Right-of-Use Asset} \quad (\text{+Asset})$$
    $$\text{Credit: Lease Liability} \quad (\text{+Liability})$$
  - **Balance Sheet Effect:** Expands both Assets and Liabilities equally; Equity is completely unaffected at commencement.

### B. Subsequent Measurement & Periodic Lease Payments
- When the lessee makes periodic contractual lease payments:
  - Cash is credited for the contractual disbursement.
  - The lease liability is amortized using the effective interest method.
  - The payment covers two distinct economic components:
    1. **Financing Interest Cost:** Accrued at the discount rate on the outstanding liability balance $\rightarrow$ Debited to `interest_expense` (reduces Net Income and Stockholders' Equity).
    2. **Principal Reduction:** The remainder of the payment reduces the carrying liability $\rightarrow$ Debited to `lease_liability` (reduces Liabilities).
  - **Balance Sheet Effect:**
    $$\Delta\text{Assets} (-\text{Cash}) = \Delta\text{Liabilities} (-\text{Lease Liability}) + \Delta\text{Equity} (-\text{Interest Expense})$$

### C. Treasury Stock Repurchase (Cost Method)
- When a corporation acquires its own outstanding common shares from the open market:
  - Under the **Cost Method** (the prevailing practice under US GAAP ASC 505-30), the entire acquisition price is debited to `treasury_stock`.
  - **Contra-Equity Nature:** Treasury Stock is not an asset (a corporation cannot own an equity investment in itself). It carries a normal debit balance and is presented on the balance sheet as a direct reduction of total stockholders' equity.
  - **Entry:**
    $$\text{Debit: Treasury Stock} \quad (\text{Contra-Equity } \rightarrow -\text{Equity})$$
    $$\text{Credit: Cash} \quad (-\text{Asset})$$
  - **Cash Flow Classification:** Classified under Financing Activities as a cash outflow.

### D. Treasury Stock Reissuance Above Cost
- When previously repurchased treasury shares are reissued to investors at a price exceeding historical buyback cost:
  - Cash is debited for total proceeds received (+Asset).
  - Treasury Stock is credited for its historical cost basis (removing contra-equity $\rightarrow$ +Equity).
  - **Strict No-Income Rule:** Under US GAAP (ASC 505-30-25-2), a company can never recognize an income statement gain or loss on transactions involving its own equity shares.
  - The excess of reissue proceeds over acquisition cost is credited to `additional_paid_in_capital` (+Equity).
  - **Entry:**
    $$\text{Debit: Cash} \quad (+\text{Asset, Total Proceeds})$$
    $$\text{Credit: Treasury Stock} \quad (+\text{Equity, Cost Basis})$$
    $$\text{Credit: Additional Paid-in Capital} \quad (+\text{Equity, Excess over Cost})$$
  - **Balance Sheet Effect:** $\Delta\text{Assets} = +\text{Cash}$, $\Delta\text{Equity} = +\text{Proceeds}$, $\Delta\text{Liabilities} = 0$.

---

## 3. Diagnostic Misconceptions & Error Tags

Stream-M08 defines the following diagnostic tags in `internal/engine/event_m08.go`:

1. `TagM08LeaseExpensedAtCommencement` (`"m08_lease_expensed_at_commencement"`): Learner debits Rent Expense instead of capitalizing the Right-of-Use Asset at lease commencement.
2. `TagM08ROUAssetRecordedAsEquity` (`"m08_rou_asset_recorded_as_equity"`): Learner confuses lease rights/obligations with equity capital contributions.
3. `TagM08IgnoringLeaseInterest` (`"m08_ignoring_lease_interest"`): Learner attributes the full periodic cash payment to liability principal, failing to accrue periodic financing interest expense.
4. `TagM08TreasuryStockRecordedAsAsset` (`"m08_treasury_stock_recorded_as_asset"`): Learner categorizes repurchased treasury shares as marketable securities or investment assets instead of contra-equity.
5. `TagM08GainRecordedOnTreasuryReissuance` (`"m08_gain_recorded_on_treasury_reissuance"`): Learner records an income statement gain on treasury share reissuance rather than crediting Additional Paid-in Capital.
6. `TagM08RetainedEarningsDebitedOnReissuance` (`"m08_retained_earnings_debited_on_reissuance"`): Learner misroutes reissuance excess to Retained Earnings instead of APIC.

---

## 4. Verification & Invariants Checklist

- [x] **Integer Minor Units:** All parameter amounts encoded as integer cents (`amount_minor_units`, `interest_minor_units`, `principal_minor_units`, `cost_minor_units`, `apic_minor_units`).
- [x] **Double-Entry Balancing:** Every transaction has $\sum \text{Debits} == \sum \text{Credits}$ to the penny across all permutations.
- [x] **Equation Invariant:** Derived balance sheet equation effects strictly satisfy $\Delta \text{Assets} == \Delta \text{Liabilities} + \Delta \text{Equity}$.
- [x] **Socratic Hint Invariant:** Every `CausalHint` across all 7 stages ends with a question mark (`?`) and prompts reflective causal reasoning without revealing answers.
- [x] **Stage 1 Neutrality:** Explanations for Stage 1 (`IdentifyAccount`) strictly omit `"debit"` and `"credit"`.
- [x] **Triple Parameter Permutations:** Each question defines 3 distinct parameter sets to ensure replay variety.
- [x] **Isolated Unit Tests:** Standalone test suite in `internal/drill/syllabus_m08_test.go` executed and passed 100% cleanly (`PASS: TestStreamM08BatchInvariants`, `TestStreamM08EngineRules`, `TestStreamM08StageBuildersAndPedagogy`, `TestStreamM08PedagogyPolicy`).
- [x] **Zero Contention:** No shared files or hardcoded test assertions were modified.

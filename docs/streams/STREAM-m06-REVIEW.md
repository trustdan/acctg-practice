# Stream-M06 Accounting & Pedagogical Review

**Reviewer:** Worker Agent 1 (Stream-M06: Module 6 PPE, Depreciation & Impairments)  
**Date:** 2026-10-12  
**Candidate Batch:** `curriculum/candidates/batch_m06.json`  
**Candidate Policy:** `curriculum/candidates/pedagogy_m06.json`  
**Harness Tests:** `internal/drill/syllabus_m06_test.go`  

---

## 1. Syllabus & Topic Scope

Stream-M06 targets core curriculum requirements from the University of Washington Foster School of Business *Financial Reporting and Analysis* course for **Module 6: Non-Current Assets, Depreciation & Impairments** (with fundamental ties to **Module 1: Accounting Equation** and **Module 2: Accrual Adjusting Entries**).

| Question ID | Scenario & Accounting Focus | US GAAP / Standard Citation | Key Concepts |
|---|---|---|---|
| `q_m06_ppe_straight_line_depr_01` | Straight-Line Depreciation: periodic cost allocation of physical asset service potential; contra-asset presentation. | ASC 360-10-35-4 / US GAAP | `straight_line_depreciation`, `accumulated_depreciation`, `contra_asset`, `adjusting_entries` |
| `q_m06_ppe_macrs_tax_timing_02` | Book vs. Tax Depreciation: straight-line GAAP book expense vs. accelerated MACRS tax deductions; temporary timing difference. | ASC 740 / IRC § 168 / ASC 360 | `macrs_vs_straight_line`, `depreciation_timing_difference`, `book_tax_differences`, `contra_asset` |
| `q_m06_ppe_asset_disposal_gain_loss_03` | Equipment Derecognition & Sale: derecognizing original cost and accumulated depreciation, recognizing Gain on Disposal for proceeds above book value. | ASC 360-10-40 | `asset_disposal`, `gain_on_disposal`, `asset_derecognition`, `book_value_calculation` |
| `q_m06_ppe_impairment_recoverability_04` | ASC 360 Two-Step Impairment Test: Step 1 recoverability failure (undiscounted cash flows < carrying value) triggering Step 2 write-down to fair value; permanent under US GAAP. | ASC 360-10-35-17 through 21 | `asc_360_impairment`, `recoverability_test`, `fair_value_measurement`, `loss_on_impairment` |

---

## 2. Technical Accounting Analysis

### A. Straight-Line Depreciation & Contra-Asset Mechanics
- **Cost Allocation Principle:** Depreciation in financial accounting is an allocation of historical cost over useful life in proportion to economic benefits received, *not* a valuation technique measuring market price fluctuations.
- **Formula:**
  $$\text{Annual Depreciation Expense} = \frac{\text{Historical Cost} - \text{Estimated Salvage Value}}{\text{Estimated Useful Life (Years)}}$$
- **Contra-Asset Presentation:** Rather than crediting `equipment` directly, US GAAP mandates crediting `accumulated_depreciation`. This preserves the historical cost on the balance sheet while allowing users to observe cumulative cost consumption:
  $$\text{Net Carrying Value (Book Value)} = \text{Equipment Historical Cost} - \text{Accumulated Depreciation}$$

### B. Book vs. Tax Depreciation Timing Differences (Straight-Line vs. MACRS)
- **Financial Reporting (Book):** Straight-line depreciation is used to match periodic cost against revenues generated over useful life (matching principle).
- **Tax Compliance (MACRS):** The Modified Accelerated Cost Recovery System (MACRS) under IRC § 168 allows accelerated deductions (e.g., 200% double-declining balance switching to straight-line, ignoring salvage value).
- **Temporary Timing Difference:** In early years of asset life, MACRS tax depreciation exceeds straight-line book depreciation, reducing current taxable income and creating a Deferred Tax Liability (DTL). In later years, book depreciation exceeds tax depreciation, reversing the difference. Financial statements record book depreciation expense to report economic reality to shareholders.

### C. PPE Derecognition & Disposal Accounting
- **Derecognition Standard (ASC 360-10-40):** Upon sale or retirement of a capital asset, the entity must:
  1. Record cash or other consideration received (Debit `cash`).
  2. Eliminate the entire accumulated depreciation balance associated with the asset to date (Debit `accumulated_depreciation`).
  3. Remove the historical cost from the asset account (Credit `equipment`).
  4. Recognize any difference as a gain or loss in operating/other income:
     $$\text{Gain (Loss) on Disposal} = \text{Proceeds} - (\text{Cost} - \text{Accumulated Depreciation})$$
- When proceeds exceed book value, a `gain_on_disposal` is credited; when proceeds are below book value, a `loss_on_disposal` is debited.

### D. ASC 360 Two-Step Impairment Testing
- **Triggering Events:** Long-lived assets held and used are tested for impairment whenever events or changes in circumstances indicate that carrying value may not be recoverable (e.g., significant adverse change in business climate, physical damage, technological obsolescence).
- **Step 1: Recoverability Test:**
  - Compare carrying value to the **sum of undiscounted expected future cash flows** directly associated with the asset:
    - If $\sum \text{Undiscounted Cash Flows} \ge \text{Carrying Value} \implies$ **No Impairment** (Asset is recoverable, even if fair value is lower).
    - If $\sum \text{Undiscounted Cash Flows} < \text{Carrying Value} \implies$ **Impairment Triggered** (Proceed to Step 2).
- **Step 2: Measurement of Loss:**
  $$\text{Impairment Loss} = \text{Carrying Value} - \text{Fair Value}$$
- **GAAP vs. IFRS Divergence:**
  - Under **US GAAP (ASC 360)**, once an impairment loss is recognized on an asset held and used, the reduced carrying amount establishes a new cost basis. **Subsequent restoration/reversal of impairment losses is strictly prohibited.**
  - Under **IFRS (IAS 36)**, impairment testing uses a single-step discounted recoverable amount, and impairment losses (except on goodwill) can be reversed if market conditions improve.

---

## 3. Diagnostic Misconceptions & Error Tags

Stream-M06 defines the following diagnostic tags in `internal/engine/event_m06.go`:

1. `TagM06DirectAssetCredit` (`"m06_direct_asset_credit"`): Student directly credits Equipment instead of using the contra-asset Accumulated Depreciation.
2. `TagM06AccumulatedDeprExpensed` (`"m06_accumulated_depreciation_expensed"`): Student treats Accumulated Depreciation as an income statement expense rather than a balance sheet contra-asset.
3. `TagM06TaxDepreciationRecordedOnBooks` (`"m06_tax_depreciation_recorded_on_books"`): Student attempts to record tax MACRS accelerated depreciation on GAAP financial statements.
4. `TagM06GainTreatedAsRevenue` (`"m06_gain_treated_as_operating_revenue"`): Student records cash proceeds from capital asset sales as core operating Service Revenue rather than Gain on Disposal.
5. `TagM06DisposalAccumDeprRetained` (`"m06_disposal_accum_depr_retained"`): Student omits the removal of Accumulated Depreciation when derecognizing sold equipment.
6. `TagM06ImpairmentWithoutRecoverability` (`"m06_impairment_without_recoverability_test"`): Student calculates impairment without verifying Step 1 undiscounted cash flow recoverability failure.
7. `TagM06ImpairmentReversedUnderGAAP` (`"m06_impairment_reversed_under_gaap"`): Student assumes US GAAP allows reversing an impairment write-down when fair value subsequently recovers.

---

## 4. Verification & Invariants Checklist

- [x] **Integer Minor Units:** All financial values encoded in integer cents (`amount_minor_units`: 600000, 800000, 5000000, 4000000, etc.).
- [x] **Double-Entry Balancing:** Every transaction has $\sum \text{Debits} == \sum \text{Credits}$ to the penny across all parameter permutations.
- [x] **Compound Entry Balancing:** Disposal question (`q_m06_ppe_asset_disposal_gain_loss_03`) balances 4 simultaneous postings (Cash Dr + AccumDepr Dr == Equipment Cr + Gain Cr).
- [x] **Equation Invariant:** Derived balance sheet equation effects strictly satisfy $\Delta \text{Assets} == \Delta \text{Liabilities} + \Delta \text{Equity}$.
- [x] **Socratic Hint Invariant:** Every `CausalHint` across all 7 stages ends in a question mark (`?`) and prompts causal understanding rather than revealing answers.
- [x] **Stage 1 Neutrality:** Explanations for Stage 1 (`IdentifyAccount`) strictly omit `"debit"` and `"credit"`.
- [x] **Triple Parameter Permutations:** Each question defines at least 3 distinct monetary amounts for replay variety.
- [x] **Isolated Unit Tests:** Standalone test suite in `internal/drill/syllabus_m06_test.go` executed and passed 100% cleanly (`PASS: TestStreamM06BatchInvariants`, `TestStreamM06EngineRules`, `TestStreamM06StageBuildersAndPedagogy`, `TestStreamM06PedagogyPolicy`).
- [x] **Zero Contention:** No shared files, accounts catalog, or central test assertions were modified.

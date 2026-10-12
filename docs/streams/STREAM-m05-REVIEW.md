# Stream-M05 Accounting & Pedagogical Review

**Reviewer:** Worker Agent 3 (Stream-M05: Module 5 Inventory & COGS)  
**Date:** 2026-10-12  
**Candidate Batch:** `curriculum/candidates/batch_m05.json`  
**Candidate Policy:** `curriculum/candidates/pedagogy_m05.json`  
**Harness Tests:** `internal/drill/syllabus_m05_test.go`  

---

## 1. Syllabus & Topic Scope

Stream-M05 targets core curriculum requirements from the University of Washington Foster School of Business *Financial Reporting and Analysis* course for **Module 5: Inventory & Cost of Goods Sold** (with secondary ties to **Module 2: Double-Entry Mechanics & Accrual Cycle**).

| Question ID | Scenario & Accounting Focus | US GAAP / Standard Citation | Key Concepts |
|---|---|---|---|
| `q_m05_inv_fifo_costing_01` | Perpetual FIFO Costing: oldest costs expensed to COGS upon delivery; inventory reduced. | ASC 330 / US GAAP | `inventory_costing_fifo`, `cogs_recognition`, `perpetual_inventory_system`, `matching_principle` |
| `q_m05_inv_lifo_costing_02` | Perpetual LIFO Costing: newest costs matched against current sales revenue during inflation; tax shield & conformity. | IRC § 472 / ASC 330 vs. IFRS IAS 2 | `inventory_costing_lifo`, `cogs_recognition`, `inflation_tax_shield`, `lifo_conformity_rule` |
| `q_m05_inv_nrv_writedown_03` | Lower of Cost or Net Realizable Value (NRV): write-down of damaged/obsolete inventory directly to COGS/loss; irreversible. | ASC 330-10-35 | `lower_of_cost_or_nrv`, `inventory_writedown`, `asc_330`, `conservatism_principle` |
| `q_m05_inv_freight_in_cap_04` | Freight-In Capitalization: inbound transportation costs paid in cash under FOB shipping point capitalized to Inventory. | ASC 330-10-30-1 | `freight_in_capitalization`, `fob_shipping_point`, `inventory_cost_basis`, `asset_exchange` |

---

## 2. Technical Accounting Analysis

### A. FIFO vs. LIFO Perpetual Costing Dynamics
- **FIFO (First-In, First-Out):** Assumes units purchased first are sold first. In an inflationary environment with rising inventory costs:
  - COGS reflects older, lower acquisition costs $\rightarrow$ **Higher Gross Profit and Net Income**.
  - Ending Inventory reflects recent, higher costs $\rightarrow$ **Higher Balance Sheet Asset Valuation**.
  - Cash outflows for income taxes increase due to higher taxable income.
- **LIFO (Last-In, First-Out):** Assumes units purchased most recently are sold first:
  - COGS reflects newest, higher acquisition costs $\rightarrow$ **Lower Gross Profit and Net Income**.
  - Ending Inventory reflects older, stale costs $\rightarrow$ **Lower Balance Sheet Asset Valuation**.
  - Tax savings in periods of inflation (LIFO Tax Shield).
  - **LIFO Conformity Rule (IRC § 472):** If LIFO is used on the federal tax return, it *must* also be used in financial reporting.
  - **GAAP vs. IFRS Divergence:** LIFO is permitted under US GAAP, but strictly prohibited under IFRS (IAS 2).

### B. Lower of Cost or Net Realizable Value (ASC 330)
- Under US GAAP (ASC 330), inventory measured using FIFO or average cost is evaluated at each reporting date at the lower of historical cost or **Net Realizable Value (NRV)**:
  $$\text{NRV} = \text{Estimated Selling Price} - \text{Predictable Costs of Completion, Disposal, and Transportation}$$
- When NRV $<$ Cost, an immediate write-down is recognized:
  - Debit: Cost of Goods Sold (or Loss on Inventory Write-down)
  - Credit: Inventory
- **Non-Reversal Rule:** Under US GAAP, once written down, the new cost basis cannot be increased if market prices subsequently recover (unlike IFRS, which allows reversals up to original cost).

### C. Freight-In Capitalization vs. Freight-Out Expensing
- **Freight-In (Inbound Transportation):** Incurred by the purchaser under FOB Shipping Point terms to bring goods to the warehouse ready for sale. Under ASC 330, all expenditures directly necessary to bring the asset to its intended location and condition are **product costs** capitalized into the `inventory` asset.
- **Freight-Out (Outbound Delivery):** Incurred by the seller to transport goods to customers under FOB Destination terms. This is a **period operating expense** (selling expense), never capitalized into inventory.

---

## 3. Diagnostic Misconceptions & Error Tags

Stream-M05 defines the following diagnostic tags in `internal/engine/event_m05.go`:

1. `TagM05LifoUsedInIFRS` (`"m05_lifo_used_in_ifrs"`): Learner assumes LIFO is accepted internationally under IFRS.
2. `TagM05FreightInExpensed` (`"m05_freight_in_expensed"`): Learner debits Delivery Expense or Operating Expense instead of capitalizing into Inventory.
3. `TagM05NRVWriteDownReversed` (`"m05_nrv_writedown_reversed"`): Learner believes US GAAP permits upward recovery reversals after a write-down.
4. `TagM05InventoryRecordedAtSelling` (`"m05_inventory_recorded_at_selling_price"`): Learner records inventory removal at retail sales price rather than historical cost basis.
5. `TagM05COGSConfusedWithOperating` (`"m05_cogs_confused_with_operating_expense"`): Learner routes cost of merchandise to general operating expense accounts (e.g., Rent Expense).
6. `TagM05InventoryCreditedOnPurchase` (`"m05_inventory_credited_on_purchase"`): Learner inverts the debit/credit side on inventory addition.

---

## 4. Verification & Invariants Checklist

- [x] **Integer Minor Units:** All parameters expressed in integer cents (`amount_minor_units`: 1200000, 1500000, 500000, 80000, etc.).
- [x] **Double-Entry Balancing:** Every transaction has $\sum \text{Debits} == \sum \text{Credits}$ to the penny.
- [x] **Equation Invariant:** Derived balance sheet equation effects strictly satisfy $\Delta \text{Assets} == \Delta \text{Liabilities} + \Delta \text{Equity}$.
- [x] **Socratic Hint Invariant:** Every `CausalHint` across all 7 stages ends in a question mark (`?`) and prompts reflective thinking rather than giving away the target answer.
- [x] **Stage 1 Neutrality:** Explanations for Stage 1 (`IdentifyAccount`) omit `"debit"` and `"credit"`.
- [x] **Triple Parameter Permutations:** Each question defines 3 distinct monetary amounts to ensure replay variety.
- [x] **Isolated Unit Tests:** Standalone test suite in `internal/drill/syllabus_m05_test.go` executed and passed 100% cleanly (`PASS: TestStreamM05BatchInvariants`, `TestStreamM05EngineRules`, `TestStreamM05StageBuildersAndPedagogy`, `TestStreamM05PedagogyPolicy`).
- [x] **Zero Contention:** No shared files or hardcoded test assertions were modified.

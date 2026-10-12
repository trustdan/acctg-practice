# Stream-M07 Accounting & Pedagogical Review

**Reviewer:** Worker Agent 2 (Stream-M07: Module 7 Liabilities & Debt Financing)  
**Date:** 2026-10-12  
**Candidate Batch:** `curriculum/candidates/batch_m07.json`  
**Candidate Policy:** `curriculum/candidates/pedagogy_m07.json`  
**Harness Tests:** `internal/drill/syllabus_m07_test.go`  

---

## 1. Syllabus & Topic Scope

Stream-M07 implements core long-term debt and liability curriculum requirements from the University of Washington Foster School of Business *Financial Reporting and Analysis* course for **Module 7: Liabilities, Bond Pricing & Effective Interest Amortization** (with secondary ties to **Module 2: Accrual Cycle & Double-Entry Mechanics**).

| Question ID | Scenario & Accounting Focus | US GAAP / Standard Citation | Key Concepts |
|---|---|---|---|
| `q_m07_bnd_issuance_discount_01` | Issuing bonds at a discount when market yield exceeds coupon rate; contra-liability mechanics. | ASC 470 / ASC 835-30 | `bond_discount_issuance`, `contra_liability`, `discount_on_bonds_payable_classification`, `bonds_payable_classification` |
| `q_m07_bnd_issuance_premium_02` | Issuing bonds at a premium when coupon rate exceeds market yield; adjunct liability mechanics. | ASC 470 / ASC 835-30 | `bond_premium_issuance`, `adjunct_liability`, `premium_on_bonds_payable_classification`, `bonds_payable_classification` |
| `q_m07_bnd_effective_interest_amort_03` | Semiannual coupon payment and discount amortization under effective interest method; compound 3-posting entry. | ASC 835-30-35 | `effective_interest_method`, `bond_discount_amortization`, `interest_expense_recognition`, `carrying_value_accretion` |
| `q_m07_bnd_contingent_liability_04` | Accrual of probable and reasonably estimable litigation loss contingency vs. disclosure-only threshold. | ASC 450-20 | `contingent_liability_accrual`, `asc_450`, `probable_and_estimable`, `loss_contingency` |

---

## 2. Technical Accounting Analysis

### A. Bond Pricing & Valuation Mechanics: Par, Discount, and Premium
- **Present Value Determination:** A bond's market price is determined by discounting two separate cash flow streams to present value at the prevailing **market rate of interest (yield)** at the date of issuance:
  1. The lump-sum principal face value paid at maturity: $\text{PV} = \text{Face Value} \times (1 + r)^{-n}$.
  2. The periodic ordinary annuity of cash coupon interest payments: $\text{PV} = \text{Coupon} \times \left[\frac{1 - (1+r)^{-n}}{r}\right]$.
- **Issuance at a Discount (Market Rate > Stated Coupon Rate):**
  - When investors can earn higher returns elsewhere in the market, they will only purchase the bond if its upfront price is reduced below par value.
  - The issuer debits `cash` for net proceeds received and credits `bonds_payable` for the contractual maturity face value.
  - The shortfall is debited to `discount_on_bonds_payable`, a **contra-liability account** (normal debit balance). On the balance sheet, it is presented as a direct deduction from `bonds_payable` to yield the initial net carrying value:
    $$\text{Carrying Value} = \text{Bonds Payable (Face)} - \text{Discount on Bonds Payable}$$
  - **Pedagogical Note:** Discount is *not* an immediate borrowing loss or expense; it represents additional interest expense that will be recognized over the life of the debt as the discount is amortized.
- **Issuance at a Premium (Market Rate < Stated Coupon Rate):**
  - When the bond's stated coupon rate offers a return higher than prevailing market yields, investors compete to buy the bond and pay more than par value.
  - The issuer debits `cash` for the total proceeds received and credits `bonds_payable` for the maturity face value.
  - The excess proceeds are credited to `premium_on_bonds_payable`, an **adjunct liability account** (normal credit balance). On the balance sheet, it is presented as an addition to `bonds_payable`:
    $$\text{Carrying Value} = \text{Bonds Payable (Face)} + \text{Premium on Bonds Payable}$$

### B. The Effective Interest Method (ASC 835-30) vs. Straight-Line Amortization
- **US GAAP Requirement:** Under ASC 835-30-35-2, the **effective interest method** is mandatory for amortizing bond discounts and premiums unless the straight-line method results in financial statements that are not materially different.
- **Economic Principles of the Effective Interest Method:**
  1. **Constant Rate of Return:** Produces a periodic interest expense equal to a constant percentage of the carrying value of the debt (the historical market interest rate at issuance).
  2. **Interest Expense Calculation:**
     $$\text{Periodic Interest Expense} = \text{Beginning Carrying Value} \times \text{Semiannual Market Rate}$$
  3. **Cash Coupon Payment:**
     $$\text{Cash Paid} = \text{Face Value} \times \text{Semiannual Stated Rate}$$
  4. **Discount Amortization:**
     $$\text{Discount Amortization} = \text{Interest Expense} - \text{Cash Paid}$$
  - Because interest expense exceeds cash paid, the difference is credited to `discount_on_bonds_payable`. This reduces the contra-liability balance, causing the net carrying value of the bond to **accrete** upward toward par value each period until reaching exactly par value at maturity.

### C. ASC 450 Contingent Liabilities: Accrual vs. Disclosure Thresholds
Under ASC 450 (Contingencies), the accounting treatment of a loss contingency depends on the likelihood of the future confirming event and whether the amount can be quantified:

```
┌─────────────────────────────────┬──────────────────────────────────┬─────────────────────────────────┐
│ Likelihood Threshold            │ Reasonably Estimable Amount      │ Cannot Be Reasonably Estimated  │
├─────────────────────────────────┼──────────────────────────────────┼─────────────────────────────────┤
│ Probable (likely to occur)      │ ACCRUE Loss & Liability (Entry)  │ Disclose in Footnotes Only      │
│ Reasonably Possible             │ Disclose in Footnotes Only       │ Disclose in Footnotes Only      │
│ Remote (slight chance)          │ No Accrual / No Disclosure       │ No Accrual / No Disclosure      │
└─────────────────────────────────┴──────────────────────────────────┴─────────────────────────────────┘
```

- When an adverse legal judgment is both **probable** and **reasonably estimable**, US GAAP strictly requires:
  - Debit: `litigation_loss` (reduces net income and Retained Earnings)
  - Credit: `contingent_liability` (records the balance sheet claim)
- If either condition is not met (e.g., probable but unknown amount, or reasonably possible), accrual is prohibited and footnote disclosure is required.

---

## 3. Diagnostic Misconceptions & Error Tags

Stream-M07 defines the following targeted diagnostic error tags in `internal/engine/event_m07.go`:

1. `TagM07DiscountRecordedAsLoss` (`"m07_discount_recorded_as_loss"`): Learner immediately expenses the bond discount as an upfront loss or financing expense rather than capitalizing it into a contra-liability account amortized over the bond term.
2. `TagM07PremiumRecordedAsGain` (`"m07_premium_recorded_as_gain"`): Learner credits the bond premium to revenue or immediate gain on issuance rather than establishing an adjunct liability.
3. `TagM07StraightLineUsedInsteadOfEffective` (`"m07_straight_line_used_instead_of_effective"`): Learner divides total discount evenly by periods instead of applying the constant effective yield required under ASC 835-30.
4. `TagM07IgnoringBondAmortization` (`"m07_ignoring_bond_amortization"`): Learner equates periodic interest expense with the cash coupon payment, failing to amortize the bond discount.
5. `TagM07ContingencyDisclosedWithoutAccrual` (`"m07_contingency_disclosed_without_accrual"`): Learner assumes all legal contingencies are footnote-only disclosures, failing to accrue a probable and estimable claim.
6. `TagM07ContingentRecordedAsPayable` (`"m07_contingent_recorded_as_trade_payable"`): Learner routes legal liability accruals to routine trade `accounts_payable`.
7. `TagM07RemoteContingencyAccrued` (`"m07_remote_contingency_accrued"`): Learner prematurely books an accrual for a claim that is merely possible or remote.

---

## 4. Verification & Invariants Checklist

- [x] **Integer Minor Units:** All parameters expressed in integer cents (`amount_minor_units`: 100000000, 3800000, 25000000, etc.).
- [x] **Double-Entry Balancing:** Every transaction has $\sum \text{Debits} == \sum \text{Credits}$ to the penny across all parameter permutations.
- [x] **Equation Invariant:** Derived balance sheet equation effects strictly satisfy $\Delta \text{Assets} == \Delta \text{Liabilities} + \Delta \text{Equity}$ for all 4 rule families.
- [x] **Socratic Hint Invariant:** Every `CausalHint` across all 28 stages ends in a question mark (`?`) and prompts reflective causal reasoning without revealing account names.
- [x] **Stage 1 Neutrality:** Explanations for Stage 1 (`IdentifyAccount`) omit `"debit"` and `"credit"`.
- [x] **Triple Parameter Permutations:** Each question defines at least 3 distinct monetary amounts for replay variety.
- [x] **Isolated Unit Tests:** Standalone test suite in `internal/drill/syllabus_m07_test.go` executed and passed 100% cleanly (`PASS: TestStreamM07BatchInvariants`, `TestStreamM07EngineRules`, `TestStreamM07StageBuildersAndPedagogy`, `TestStreamM07PedagogyPolicy`).
- [x] **Zero Contention:** No central files (`curriculum/seed-questions.json`, `curriculum/pedagogy.json`, `curriculum/accounts.json`, `internal/drill/generator.go`, or test counters) were modified.

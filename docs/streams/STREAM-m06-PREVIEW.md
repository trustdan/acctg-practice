# Stream-M06 Rendered Stage Review & Pedagogical Preview

This document contains the rendered 7-stage scaffolding preview for all candidate questions produced in **Stream-M06** (Module 6: Property, Plant & Equipment, Depreciation, and Impairments) under Wave 2 of the Multi-Agent Question Expansion Plan.

---

## Question 1: `q_m06_ppe_straight_line_depr_01`

- **Title / Topic:** Straight-Line Depreciation Adjusting Entry
- **Family ID:** `m06_ppe_straight_line_depr`
- **Module Scope:** Module 6 (Property, Plant & Equipment), Module 1 (Accounting Equation), Module 2 (Accrual Adjustments)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[600000, 1200000, 2400000]` ($6,000, $12,000, $24,000)
- **Concepts:** `["straight_line_depreciation", "accumulated_depreciation", "contra_asset", "adjusting_entries"]`
- **Scenario Template:**
  > "At fiscal year-end, the company records annual straight-line depreciation of ${amount_dollars} on factory equipment with an estimated useful life of 5 years and salvage value. Record the adjusting journal entry to recognize annual depreciation expense."

### Rendered Scaffolding Stages (Example: $6,000.00)

#### Stage 1: Identify Account
- **Prompt:** At fiscal year-end, which expense account is used to record the periodic allocation of equipment cost over its useful life?
- **Options:**
  - `opt_depreciation_expense`: Depreciation Expense *(Correct)*
  - `opt_equipment`: Equipment *(Tag: `wrong_account`)*
  - `opt_accumulated_depreciation`: Accumulated Depreciation *(Tag: `m06_accumulated_depreciation_expensed`)*
  - `opt_rent_expense`: Rent Expense *(Tag: `wrong_account`)*
- **Causal Hint:** When equipment is utilized during an accounting period, which operational expense account captures the consumed service potential of the asset?
- **Explanation:** Depreciation Expense is the operating expense account that records the periodic cost allocation of physical productive assets over their estimated useful life.

#### Stage 2: Account Category
- **Prompt:** What category of account is Depreciation Expense?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does the periodic expired cost of using productive machinery represent an operating expense on the income statement, or an unexpired asset balance?
- **Explanation:** Depreciation Expense is an operating Expense account on the income statement that directly reduces operating income and Net Income.

#### Stage 3: Direction
- **Prompt:** Does Depreciation Expense increase or decrease when recording this year-end adjusting entry?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** As another accounting period expires and equipment utility is consumed, does total recorded depreciation expense accumulate upward or downward?
- **Explanation:** Depreciation Expense increases as additional asset cost expires and is recognized in current-period earnings.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expenses carry normal debit balances and reduce stockholders' equity. Do expenses increase on the left or the right side?
- **Explanation:** Expenses increase by Debit on the left side, reflecting a reduction in stockholders' equity.

#### Stage 5: Counter Account
- **Prompt:** Which contra-asset balance sheet account is credited to accumulate the total cost expired to date without altering original equipment cost?
- **Options:**
  - `opt_accumulated_depreciation`: Accumulated Depreciation (Contra-Asset) *(Correct)*
  - `opt_equipment`: Equipment (Asset) *(Tag: `m06_direct_asset_credit`)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
- **Causal Hint:** Which balance sheet contra-account maintains the cumulative depreciation recognized against equipment across its operational life?
- **Explanation:** Accumulated Depreciation is a contra-asset account credited to offset Equipment, allowing financial statement users to see both historical cost and net carrying value.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced adjusting journal entry to record annual straight-line depreciation?
- **Options:**
  - `opt_dr_exp_cr_accum`: Debit Depreciation Expense $6,000.00 / Credit Accumulated Depreciation $6,000.00 *(Correct)*
  - `opt_dr_exp_cr_equip`: Debit Depreciation Expense $6,000.00 / Credit Equipment $6,000.00 *(Tag: `m06_direct_asset_credit`)*
  - `opt_dr_accum_cr_exp`: Debit Accumulated Depreciation $6,000.00 / Credit Depreciation Expense $6,000.00 *(Tag: `reversed_sides`)*
- **Causal Hint:** Which entry pairs the periodic operating expense with the contra-asset that reduces equipment carrying value?
- **Explanation:** Debit Depreciation Expense $6,000.00 and Credit Accumulated Depreciation $6,000.00. Under the straight-line method, annual depreciation is computed as (Historical Cost - Salvage Value) / Useful Life.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording straight-line depreciation affect the fundamental balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $6,000.00 (-Accumulated Depreciation); Equity decreases by $6,000.00 (-Depreciation Expense); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $6,000.00; Liabilities decrease by $6,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** How do an increase in a contra-asset and an increase in an operating expense impact total assets and total equity?
- **Explanation:** Assets decrease (-Accumulated Depreciation $6,000.00) and Equity decreases (-Depreciation Expense $6,000.00); Liabilities are unchanged. Accumulated depreciation reduces net PPE on the balance sheet.

---

## Question 2: `q_m06_ppe_macrs_tax_timing_02`

- **Title / Topic:** Accelerated MACRS vs Straight-Line Book Timing Difference
- **Family ID:** `m06_ppe_macrs_tax_timing`
- **Module Scope:** Module 6 (Property, Plant & Equipment), Module 1 (Accounting Equation), Module 2 (Accrual Adjustments)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[800000, 1600000, 3200000]` ($8,000, $16,000, $32,000)
- **Concepts:** `["macrs_vs_straight_line", "depreciation_timing_difference", "book_tax_differences", "contra_asset"]`
- **Scenario Template:**
  > "For financial statement reporting under US GAAP, the company uses straight-line depreciation to reflect the steady economic consumption of its operating machinery, recording ${amount_dollars} of book depreciation expense for the year, while using accelerated MACRS on its tax return. Record the company's financial accounting adjusting entry for book depreciation."

### Rendered Scaffolding Stages (Example: $8,000.00)

#### Stage 1: Identify Account
- **Prompt:** For financial reporting purposes under US GAAP, which account records the annual straight-line cost allocation of equipment when accelerated MACRS is used on tax returns?
- **Options:**
  - `opt_depreciation_expense`: Depreciation Expense *(Correct)*
  - `opt_equipment`: Equipment *(Tag: `wrong_account`)*
  - `opt_accumulated_depreciation`: Accumulated Depreciation *(Tag: `m06_accumulated_depreciation_expensed`)*
  - `opt_insurance_expense`: Insurance Expense *(Tag: `wrong_account`)*
- **Causal Hint:** For GAAP financial statements, which operating expense account records the economic consumption of productive equipment regardless of accelerated tax method rules?
- **Explanation:** Depreciation Expense is the GAAP financial reporting expense account that reflects economic asset expiration under the matching principle.

#### Stage 2: Account Category
- **Prompt:** What category of account is Depreciation Expense in financial reporting?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does book depreciation represent an income statement expense reducing current earnings, or a balance sheet resource?
- **Explanation:** Depreciation Expense is an operating Expense on the income statement that reduces book pre-tax income.

#### Stage 3: Direction
- **Prompt:** Does Depreciation Expense increase or decrease when recognizing book depreciation for the year?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** Does recognizing periodic asset expiration increase or decrease the total operating expenses reported on the income statement?
- **Explanation:** Depreciation Expense increases as additional asset cost is recognized during the reporting period.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger increases an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expense accounts carry normal debit balances that reduce equity. Which side records an increase?
- **Explanation:** Expenses increase by Debit on the left side of the ledger, reducing stockholders' equity.

#### Stage 5: Counter Account
- **Prompt:** Which contra-asset account is credited on the company's financial books to offset equipment to its net carrying value?
- **Options:**
  - `opt_accumulated_depreciation`: Accumulated Depreciation (Contra-Asset) *(Correct)*
  - `opt_equipment`: Equipment (Asset) *(Tag: `m06_direct_asset_credit`)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
- **Causal Hint:** Which balance sheet contra-asset account accumulates the periodic straight-line depreciation recognized on the company's financial books?
- **Explanation:** Accumulated Depreciation is credited to offset the Equipment cost account, presenting net book value on the balance sheet.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record GAAP book depreciation expense?
- **Options:**
  - `opt_dr_exp_cr_accum`: Debit Depreciation Expense $8,000.00 / Credit Accumulated Depreciation $8,000.00 *(Correct)*
  - `opt_dr_exp_cr_equip`: Debit Depreciation Expense $8,000.00 / Credit Equipment $8,000.00 *(Tag: `m06_direct_asset_credit`)*
  - `opt_dr_accum_cr_exp`: Debit Accumulated Depreciation $8,000.00 / Credit Depreciation Expense $8,000.00 *(Tag: `reversed_sides`)*
- **Causal Hint:** Which entry correctly records straight-line depreciation expense for financial reporting and credits the contra-asset account?
- **Explanation:** Debit Depreciation Expense $8,000.00 and Credit Accumulated Depreciation $8,000.00. GAAP requires straight-line depreciation when it best matches economic consumption; accelerated MACRS deductions on tax filings create temporary timing differences.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording straight-line book depreciation affect the balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $8,000.00 (-Accumulated Depreciation); Equity decreases by $8,000.00 (-Depreciation Expense); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $8,000.00; Liabilities decrease by $8,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** When Accumulated Depreciation increases and Depreciation Expense is recognized, what is the net impact on the balance sheet equation?
- **Explanation:** Assets decrease (-Accumulated Depreciation $8,000.00) and Equity decreases (-Depreciation Expense $8,000.00); Liabilities are unchanged. Book depreciation reduces book net income and net asset carrying value.

---

## Question 3: `q_m06_ppe_asset_disposal_gain_loss_03`

- **Title / Topic:** Equipment Derecognition & Gain on Disposal
- **Family ID:** `m06_ppe_asset_disposal_gain_loss`
- **Module Scope:** Module 6 (Property, Plant & Equipment), Module 1 (Accounting Equation), Module 2 (Accrual Mechanics)
- **Status:** `candidate_pending_review`
- **Parameters:**
  - `amount_minor_units` (Equipment Cost): `[5000000, 10000000, 20000000]` ($50,000, $100,000, $200,000)
  - `accum_depr_minor_units`: `[4000000, 8000000, 16000000]` ($40,000, $80,000, $160,000)
  - `cash_minor_units`: `[1250000, 2500000, 5000000]` ($12,500, $25,000, $50,000)
  - `gain_minor_units`: `[250000, 500000, 1000000]` ($2,500, $5,000, $10,000)
- **Concepts:** `["asset_disposal", "gain_on_disposal", "asset_derecognition", "book_value_calculation"]`
- **Scenario Template:**
  > "The company retires and sells factory equipment with an original cost of ${amount_dollars} and accumulated depreciation of 80% of original cost. The equipment is sold to an external buyer for cash proceeds equal to 25% of original cost, recognizing a gain on disposal. Record the derecognition of the asset and the sale transaction."

### Rendered Scaffolding Stages (Example: $50,000.00 Cost, $40,000.00 Accum Depr, $12,500.00 Cash, $2,500.00 Gain)

#### Stage 1: Identify Account
- **Prompt:** When equipment is sold for cash proceeds exceeding its net book value (historical cost minus accumulated depreciation), which account records the excess gain recognized in earnings?
- **Options:**
  - `opt_gain_on_disposal`: Gain on Disposal *(Correct)*
  - `opt_service_revenue`: Service Revenue *(Tag: `m06_gain_treated_as_operating_revenue`)*
  - `opt_depreciation_expense`: Depreciation Expense *(Tag: `wrong_account`)*
  - `opt_retained_earnings`: Retained Earnings *(Tag: `wrong_account`)*
- **Causal Hint:** When a capital asset is sold for more than its remaining book carrying value, which account captures this peripheral transaction profit in net income?
- **Explanation:** Gain on Disposal is the account that captures the excess of net disposal proceeds over the carrying book value of the retired asset.

#### Stage 2: Account Category
- **Prompt:** What category of account is Gain on Disposal?
- **Options:**
  - `opt_revenue`: Revenue / Other Income *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_expense`: Expense
- **Causal Hint:** Does a gain resulting from selling productive equipment represent other income that increases stockholders' equity, or an asset?
- **Explanation:** Gain on Disposal is classified under Other Income / Revenues on the income statement, increasing Net Income and Stockholders' Equity.

#### Stage 3: Direction
- **Prompt:** Does Gain on Disposal increase or decrease when recognizing the profit from the equipment sale?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** When a profitable asset sale transaction occurs, does total gain recognized on the income statement increase or decrease?
- **Explanation:** Gain on Disposal increases to reflect the net gain realized upon derecognizing the asset.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase a Gain or Revenue account?
- **Options:**
  - `opt_credit`: Credit (Right side) *(Correct)*
  - `opt_debit`: Debit (Left side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Revenues and gains increase stockholders' equity and maintain normal credit balances. Do they increase by debit on the left or credit on the right?
- **Explanation:** Revenues and gains increase by Credit on the right side of the ledger, reflecting an addition to stockholders' equity.

#### Stage 5: Counter Account
- **Prompt:** Which asset account must be credited for its original cost to completely remove the sold equipment from the balance sheet?
- **Options:**
  - `opt_equipment`: Equipment (Asset) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_accumulated_depreciation`: Accumulated Depreciation (Contra-Asset) *(Tag: `m06_disposal_accum_depr_retained`)*
- **Causal Hint:** Which balance sheet asset account records the original acquisition cost of the sold equipment and must now be completely eliminated?
- **Explanation:** Equipment is credited for its full original historical cost to eliminate the derecognized asset from the company's ledger.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the cash receipt, derecognition of equipment and accumulated depreciation, and gain on disposal?
- **Options:**
  - `opt_dr_cash_accum_cr_equip_gain`: Debit Cash $12,500.00, Debit Accumulated Depreciation $40,000.00 / Credit Equipment $50,000.00, Credit Gain on Disposal $2,500.00 *(Correct)*
  - `opt_dr_cash_cr_equip_gain`: Debit Cash $12,500.00 / Credit Equipment $50,000.00, Credit Gain on Disposal $2,500.00 *(Tag: `m06_disposal_accum_depr_retained`)*
  - `opt_dr_cash_cr_rev`: Debit Cash $12,500.00 / Credit Service Revenue $12,500.00 *(Tag: `m06_gain_treated_as_operating_revenue`)*
- **Causal Hint:** Which compound entry brings in cash proceeds, eliminates accumulated depreciation, zeroes out original equipment cost, and recognizes the gain?
- **Explanation:** Debit Cash $12,500.00, Debit Accumulated Depreciation $40,000.00, Credit Equipment $50,000.00, and Credit Gain on Disposal $2,500.00. Equipment carrying value was $10,000.00 ($50,000.00 cost - $40,000.00 depr); selling for $12,500.00 cash yields a $2,500.00 gain.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does derecognizing the equipment upon sale for a gain affect the fundamental accounting equation?
- **Options:**
  - `opt_assets_up_eq_up`: Net Assets increase by $2,500.00 (+Cash $12,500.00, +Accumulated Depreciation removal $40,000.00, -Equipment $50,000.00); Equity increases by $2,500.00 (+Gain on Disposal); Liabilities unchanged. *(Correct)*
  - `opt_assets_up_liab_up`: Assets increase by $12,500.00; Liabilities increase by $12,500.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain completely unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** Net book value of the asset was surrendered in exchange for cash. How does this net gain affect total assets and total equity?
- **Explanation:** Net Assets increase by the gain amount of $2,500.00 (+Cash $12,500.00 minus Net Book Value $10,000.00); Equity increases by $2,500.00 (+Gain on Disposal); Liabilities are unchanged.

---

## Question 4: `q_m06_ppe_impairment_recoverability_04`

- **Title / Topic:** ASC 360 Two-Step Long-Lived Asset Impairment Test
- **Family ID:** `m06_ppe_impairment_recoverability`
- **Module Scope:** Module 6 (Property, Plant & Equipment), Module 1 (Accounting Equation), Module 2 (Accrual Adjustments)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[4000000, 8000000, 12000000]` ($40,000, $80,000, $120,000)
- **Concepts:** `["asc_360_impairment", "recoverability_test", "fair_value_measurement", "loss_on_impairment"]`
- **Scenario Template:**
  > "At year-end, the company performs an impairment review of production equipment under ASC 360 following technological changes. The equipment carrying value of $100,000 exceeds undiscounted expected future cash flows of $75,000 (triggering impairment under Step 1 recoverability). The fair value is measured at $60,000, requiring an impairment loss write-down of ${amount_dollars} (Step 2 measurement). Record the impairment loss adjusting entry."

### Rendered Scaffolding Stages (Example: $40,000.00)

#### Stage 1: Identify Account
- **Prompt:** Under ASC 360, when equipment fails the Step 1 recoverability test and is written down to fair value, which income statement account records the write-down?
- **Options:**
  - `opt_loss_on_impairment`: Loss on Impairment *(Correct)*
  - `opt_depreciation_expense`: Depreciation Expense *(Tag: `wrong_account`)*
  - `opt_cost_of_goods_sold`: Cost of Goods Sold *(Tag: `wrong_account`)*
  - `opt_accumulated_depreciation`: Accumulated Depreciation *(Tag: `m06_accumulated_depreciation_expensed`)*
- **Causal Hint:** Under ASC 360, when carrying value cannot be recovered through future undiscounted cash flows, which loss account recognizes the sudden write-down to fair value in operating income?
- **Explanation:** Loss on Impairment is the income statement account that reflects the sudden permanent reduction of an asset's carrying value down to its fair value.

#### Stage 2: Account Category
- **Prompt:** What category of account is Loss on Impairment?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does an impairment write-down represent an operating loss/expense reducing current period net earnings, or an asset?
- **Explanation:** Loss on Impairment is an Expense/Loss account presented within continuing operations on the income statement that reduces Net Income.

#### Stage 3: Direction
- **Prompt:** Does Loss on Impairment increase or decrease when recording this asset write-down?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** Does recording a new impairment loss write-down increase or decrease total losses and expenses on the income statement?
- **Explanation:** Loss on Impairment increases as the asset write-down is recognized in current-period earnings.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase a Loss or Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Losses and expenses reduce stockholders' equity and carry normal debit balances. Do they increase on the left or right side?
- **Explanation:** Losses increase by Debit on the left side, reducing Stockholders' Equity.

#### Stage 5: Counter Account
- **Prompt:** Which contra-asset account is credited to write down the net carrying value of the equipment to its fair value?
- **Options:**
  - `opt_accumulated_depreciation`: Accumulated Depreciation (Contra-Asset) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_accounts_payable`: Accounts Payable (Liability) *(Tag: `wrong_account`)*
- **Causal Hint:** Which balance sheet contra-account is increased to write down the net carrying amount of the impaired equipment?
- **Explanation:** Accumulated Depreciation is credited to write down the carrying value of the equipment to fair value while preserving historical cost.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the ASC 360 equipment impairment write-down?
- **Options:**
  - `opt_dr_loss_cr_accum`: Debit Loss on Impairment $40,000.00 / Credit Accumulated Depreciation $40,000.00 *(Correct)*
  - `opt_dr_accum_cr_loss`: Debit Accumulated Depreciation $40,000.00 / Credit Loss on Impairment $40,000.00 *(Tag: `reversed_sides`)*
  - `opt_dr_loss_cr_cash`: Debit Loss on Impairment $40,000.00 / Credit Cash $40,000.00 *(Tag: `cash_recorded_when_uncollected`)*
- **Causal Hint:** Which journal entry charges the impairment write-down to earnings and increases accumulated depreciation to reduce carrying value?
- **Explanation:** Debit Loss on Impairment $40,000.00 and Credit Accumulated Depreciation $40,000.00. Under ASC 360, when undiscounted future cash flows fall below carrying value, the impairment loss equals Carrying Value minus Fair Value.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording the ASC 360 impairment write-down affect the balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $40,000.00 (-Accumulated Depreciation); Equity decreases by $40,000.00 (-Loss on Impairment); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $40,000.00; Liabilities decrease by $40,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** How do crediting Accumulated Depreciation and debiting Loss on Impairment impact total assets and total equity?
- **Explanation:** Assets decrease (-Accumulated Depreciation $40,000.00) and Equity decreases (-Loss on Impairment $40,000.00); Liabilities are unchanged. Under US GAAP, impairment write-downs establish a new cost basis and cannot be reversed even if fair value recovers subsequently.

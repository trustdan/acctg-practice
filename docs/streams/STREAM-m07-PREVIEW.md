# Stream-M07 Rendered Stage Review & Pedagogical Preview

This document contains the rendered 7-stage scaffolding preview for all candidate questions produced in **Stream-M07** (Module 7: Liabilities, Bond Pricing & Effective Interest Amortization) under Wave 2 of the Multi-Agent Question Expansion Plan.

---

## Question 1: `q_m07_bnd_issuance_discount_01`

- **Title / Topic:** Issuing Bonds at a Discount (Market Rate > Stated Coupon Rate)
- **Family ID:** `m07_bnd_issuance_discount`
- **Module Scope:** Module 7 (Liabilities & Debt Financing), Module 2 (Double-Entry Mechanics & Accrual Cycle)
- **Status:** `candidate_pending_review`
- **Parameters:**
  - `amount_minor_units`: `[100000000, 200000000, 400000000]` ($1,000,000, $2,000,000, $4,000,000 face value)
  - `cash_minor_units`: `[95000000, 190000000, 380000000]` ($950,000, $1,900,000, $3,800,000 cash proceeds at 95%)
  - `discount_minor_units`: `[5000000, 10000000, 20000000]` ($50,000, $100,000, $200,000 discount at 5%)
- **Concepts:** `["bond_discount_issuance", "contra_liability", "discount_on_bonds_payable_classification", "bonds_payable_classification"]`
- **Scenario Template:**
  > "To fund long-term factory modernization, the company issues 10-year bonds with a contractual face value of ${amount_dollars}. Because the market interest rate exceeds the stated coupon rate, investors require a discount, and the company receives cash proceeds of 95% of face value with a 5% discount below par."

### Rendered Scaffolding Stages (Example: $1,000,000 Face Value / $950,000 Proceeds / $50,000 Discount)

#### Stage 1: Identify Account
- **Prompt:** When bonds are issued below par value because the market interest rate exceeds the coupon rate, which contra-liability account captures the difference between face value and cash received?
- **Options:**
  - `opt_discount_bonds`: Discount on Bonds Payable *(Correct)*
  - `opt_loss_issuance`: Loss on Bond Issuance *(Tag: `m07_discount_recorded_as_loss`)*
  - `opt_interest_expense`: Interest Expense *(Tag: `wrong_account`)*
  - `opt_premium_bonds`: Premium on Bonds Payable *(Tag: `wrong_account`)*
- **Causal Hint:** Which contra-liability account offsets the face value of the debt on the balance sheet to reflect the lower initial borrowing proceeds?
- **Explanation:** Discount on Bonds Payable is a contra-liability account that records the difference when bonds are issued for less than their maturity face value.

#### Stage 2: Account Category
- **Prompt:** What category of account is Discount on Bonds Payable?
- **Options:**
  - `opt_contra_liability`: Liability (Contra-Liability) *(Correct)*
  - `opt_asset`: Asset
  - `opt_expense`: Expense
  - `opt_equity`: Equity
- **Causal Hint:** Does this valuation account directly offset a long-term debt liability, or does it represent an independent resource or periodic expense?
- **Explanation:** Discount on Bonds Payable is a Contra-Liability account that carries a normal debit balance and is subtracted directly from Bonds Payable on the balance sheet.

#### Stage 3: Direction
- **Prompt:** Does the balance in Discount on Bonds Payable increase or decrease when the discounted bond is originally issued?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** When new bonds are sold below par value, does the cumulative discount attached to the debt portfolio grow larger or smaller?
- **Explanation:** Discount on Bonds Payable increases as the initial discount is established upon issuance.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger increases a Contra-Liability account with a normal debit balance?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Contra-liabilities have normal balances opposite to liabilities. Do they increase on the left or the right side?
- **Explanation:** Contra-liability accounts increase on the left side (Debit), offsetting the credit balance of the primary liability.

#### Stage 5: Counter Account
- **Prompt:** Which liability account records the full par value principal repayment obligation due to bondholders at maturity?
- **Options:**
  - `opt_bonds_payable`: Bonds Payable (Liability) *(Correct)*
  - `opt_notes_payable`: Notes Payable (Liability) *(Tag: `wrong_account`)*
  - `opt_common_stock`: Common Stock (Equity) *(Tag: `wrong_account`)*
- **Causal Hint:** Which long-term debt account must always be credited for the full contractual face value owed upon maturity?
- **Explanation:** Bonds Payable is credited for the contractual face amount owed at maturity.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record issuing the bonds at a discount?
- **Options:**
  - `opt_balanced_discount`: Debit Cash $950,000.00 / Debit Discount on Bonds Payable $50,000.00 / Credit Bonds Payable $1,000,000.00 *(Correct)*
  - `opt_no_discount`: Debit Cash $950,000.00 / Credit Bonds Payable $950,000.00 *(Tag: `m07_discount_recorded_as_loss`)*
  - `opt_loss_discount`: Debit Cash $950,000.00 / Debit Loss on Bond Issuance $50,000.00 / Credit Bonds Payable $1,000,000.00 *(Tag: `m07_discount_recorded_as_loss`)*
- **Causal Hint:** Which journal entry records the cash proceeds, recognizes the contra-liability discount, and records the full face value debt obligation?
- **Explanation:** Debit Cash $950,000.00 for proceeds received, Debit Discount on Bonds Payable $50,000.00 for the difference below par, and Credit Bonds Payable $1,000,000.00 for the full face value. The discount is amortized over the life of the bond as additional interest expense.

#### Stage 7: Equation Effect
- **Prompt:** How does issuing bonds at a discount affect the fundamental balance sheet equation?
- **Options:**
  - `opt_assets_up_liab_up`: Assets increase by $950,000.00 (+Cash); Liabilities increase by $950,000.00 (+Bonds Payable $1,000,000.00, -Discount $50,000.00); Equity unchanged. *(Correct)*
  - `opt_assets_up_equity_down`: Assets increase by $950,000.00; Equity decreases by $50,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets and Liabilities remain completely unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** Cash was received, long-term debt increased, and a contra-liability reduced gross debt to carrying value. What is the net impact on Assets, Liabilities, and Equity?
- **Explanation:** Assets increase by the cash proceeds (+Cash $950,000.00); Net Liabilities increase by the net carrying value (+Bonds Payable $1,000,000.00 minus Discount $50,000.00 = +$950,000.00); Equity is completely unaffected upon debt issuance.

---

## Question 2: `q_m07_bnd_issuance_premium_02`

- **Title / Topic:** Issuing Bonds at a Premium (Market Rate < Stated Coupon Rate)
- **Family ID:** `m07_bnd_issuance_premium`
- **Module Scope:** Module 7 (Liabilities & Debt Financing), Module 2 (Double-Entry Mechanics & Accrual Cycle)
- **Status:** `candidate_pending_review`
- **Parameters:**
  - `amount_minor_units`: `[100000000, 200000000, 400000000]` ($1,000,000, $2,000,000, $4,000,000 face value)
  - `cash_minor_units`: `[105000000, 210000000, 420000000]` ($1,050,000, $2,100,000, $4,200,000 cash proceeds at 105%)
  - `premium_minor_units`: `[5000000, 10000000, 20000000]` ($50,000, $100,000, $200,000 premium at 5%)
- **Concepts:** `["bond_premium_issuance", "adjunct_liability", "premium_on_bonds_payable_classification", "bonds_payable_classification"]`
- **Scenario Template:**
  > "The company issues 10-year bonds with a contractual face value of ${amount_dollars}. Because the stated coupon rate exceeds the prevailing market yield, investors pay a premium above par, yielding cash proceeds equal to 105% of face value with a 5% premium."

### Rendered Scaffolding Stages (Example: $1,000,000 Face Value / $1,050,000 Proceeds / $50,000 Premium)

#### Stage 1: Identify Account
- **Prompt:** When bonds are issued above par value because the stated coupon rate exceeds the market interest rate, which adjunct liability account captures the excess cash proceeds above face value?
- **Options:**
  - `opt_premium_bonds`: Premium on Bonds Payable *(Correct)*
  - `opt_gain_issuance`: Gain on Bond Issuance *(Tag: `m07_premium_recorded_as_gain`)*
  - `opt_interest_revenue`: Interest Revenue *(Tag: `wrong_account`)*
  - `opt_discount_bonds`: Discount on Bonds Payable *(Tag: `wrong_account`)*
- **Causal Hint:** Which adjunct liability account attaches directly to Bonds Payable to reflect the premium proceeds paid by eager investors?
- **Explanation:** Premium on Bonds Payable is an adjunct liability account that records the excess proceeds received over the maturity face value.

#### Stage 2: Account Category
- **Prompt:** What category of account is Premium on Bonds Payable?
- **Options:**
  - `opt_adjunct_liability`: Liability (Adjunct Liability) *(Correct)*
  - `opt_asset`: Asset
  - `opt_revenue`: Revenue
  - `opt_equity`: Equity
- **Causal Hint:** Does this account add to the carrying value of a long-term liability, or does it represent immediate corporate income?
- **Explanation:** Premium on Bonds Payable is an adjunct Liability account with a normal credit balance that is added to Bonds Payable on the balance sheet.

#### Stage 3: Direction
- **Prompt:** Does Premium on Bonds Payable increase or decrease when the bond is initially issued at a premium?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** When new debt is issued above par value, does the recorded bond premium balance increase or decrease?
- **Explanation:** Premium on Bonds Payable increases as the premium is initially recognized at issuance.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger increases an adjunct Liability account?
- **Options:**
  - `opt_credit`: Credit (Right side) *(Correct)*
  - `opt_debit`: Debit (Left side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Adjunct liability accounts follow the normal balance of the debt they attach to. Do liabilities increase on the debit or credit side?
- **Explanation:** Adjunct liabilities increase on the right side (Credit), augmenting the credit balance of Bonds Payable.

#### Stage 5: Counter Account
- **Prompt:** Which asset account is increased for the total cash proceeds received from the bond issuance?
- **Options:**
  - `opt_cash`: Cash (Asset) *(Correct)*
  - `opt_accounts_receivable`: Accounts Receivable (Asset) *(Tag: `wrong_account`)*
  - `opt_common_stock`: Common Stock (Equity) *(Tag: `wrong_account`)*
- **Causal Hint:** Which liquid asset account receives the upfront proceeds paid by bond investors upon purchasing the debt issue?
- **Explanation:** Cash is debited for the total funds received from bondholders.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record issuing the bonds at a premium?
- **Options:**
  - `opt_balanced_premium`: Debit Cash $1,050,000.00 / Credit Premium on Bonds Payable $50,000.00 / Credit Bonds Payable $1,000,000.00 *(Correct)*
  - `opt_gain_premium`: Debit Cash $1,050,000.00 / Credit Gain on Bond Issuance $50,000.00 / Credit Bonds Payable $1,000,000.00 *(Tag: `m07_premium_recorded_as_gain`)*
  - `opt_no_premium`: Debit Cash $1,050,000.00 / Credit Bonds Payable $1,050,000.00 *(Tag: `wrong_account`)*
- **Causal Hint:** Which journal entry records the full cash proceeds received, the face value of the debt, and the adjunct liability premium?
- **Explanation:** Debit Cash $1,050,000.00 for total proceeds received, Credit Premium on Bonds Payable $50,000.00 for the amount above par, and Credit Bonds Payable $1,000,000.00 for the contractual face value.

#### Stage 7: Equation Effect
- **Prompt:** How does issuing bonds at a premium affect the fundamental balance sheet equation?
- **Options:**
  - `opt_assets_up_liab_up`: Assets increase by $1,050,000.00 (+Cash); Liabilities increase by $1,050,000.00 (+Bonds Payable $1,000,000.00, +Premium $50,000.00); Equity unchanged. *(Correct)*
  - `opt_assets_up_equity_up`: Assets increase by $1,050,000.00; Equity increases by $50,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets and Liabilities remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** Total cash received increases Assets, while gross debt plus premium increases total carrying liabilities. Is Equity affected upon issuing debt?
- **Explanation:** Assets increase by total cash received (+Cash $1,050,000.00); Liabilities increase by the full carrying value (+Bonds Payable $1,000,000.00 plus +Premium $50,000.00 = +$1,050,000.00); Equity is completely unchanged by debt financing.

---

## Question 3: `q_m07_bnd_effective_interest_amort_03`

- **Title / Topic:** Semiannual Bond Coupon Payment & Discount Amortization (Effective Interest Method)
- **Family ID:** `m07_bnd_effective_interest_amort`
- **Module Scope:** Module 7 (Liabilities & Debt Financing), Module 2 (Double-Entry Mechanics & Accrual Cycle)
- **Status:** `candidate_pending_review`
- **Parameters:**
  - `amount_minor_units`: `[3800000, 7600000, 15200000]` ($38,000, $76,000, $152,000 Interest Expense)
  - `cash_minor_units`: `[3000000, 6000000, 12000000]` ($30,000, $60,000, $120,000 Cash coupon)
  - `discount_minor_units`: `[800000, 1600000, 3200000]` ($8,000, $16,000, $32,000 Discount amortized)
- **Concepts:** `["effective_interest_method", "bond_discount_amortization", "interest_expense_recognition", "carrying_value_accretion"]`
- **Scenario Template:**
  > "At the semiannual interest payment date, the company applies the effective interest method under ASC 835-30 for bonds issued at a discount. Total economic interest expense of ${amount_dollars} is recognized (carrying value multiplied by the semiannual market rate), while contractual cash coupon interest of 6% annual (3% semiannual) is paid to bondholders. The difference amortizes the bond discount."

### Rendered Scaffolding Stages (Example: $38,000 Interest Expense / $30,000 Cash / $8,000 Amortization)

#### Stage 1: Identify Account
- **Prompt:** Under the effective interest method, which expense account captures the true economic borrowing cost calculated as carrying value multiplied by the market interest rate?
- **Options:**
  - `opt_interest_expense`: Interest Expense *(Correct)*
  - `opt_discount_bonds`: Discount on Bonds Payable *(Tag: `wrong_account`)*
  - `opt_cash`: Cash *(Tag: `wrong_account`)*
  - `opt_bonds_payable`: Bonds Payable *(Tag: `wrong_account`)*
- **Causal Hint:** Which periodic income statement account reflects the total economic cost of servicing the bond debt for the period?
- **Explanation:** Interest Expense represents the periodic borrowing cost, computed under the effective interest method as the carrying value times the market interest rate.

#### Stage 2: Account Category
- **Prompt:** What category of account is Interest Expense?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_liability`: Liability
  - `opt_asset`: Asset
  - `opt_equity`: Equity
- **Causal Hint:** Does the periodic cost of debt financing represent an expired operational expense on the income statement, or a balance sheet obligation?
- **Explanation:** Interest Expense is an operating/financing Expense account on the income statement that reduces Net Income.

#### Stage 3: Direction
- **Prompt:** Does Interest Expense increase or decrease when the semiannual interest payment and amortization are recorded?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** As financing costs accrue over the passage of a semiannual interest period, does total expense accumulate higher or lower?
- **Explanation:** Interest Expense increases as the financing cost for the six-month period is recognized.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger increases an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expenses reduce stockholders' equity and carry normal debit balances. Do expenses increase on the left or the right side?
- **Explanation:** Expenses increase by Debit on the left side, reducing Stockholders' Equity.

#### Stage 5: Counter Account
- **Prompt:** Which contra-liability account is credited to gradually reduce its balance and amortize it toward zero at bond maturity?
- **Options:**
  - `opt_discount_bonds`: Discount on Bonds Payable (Contra-Liability) *(Correct)*
  - `opt_bonds_payable`: Bonds Payable (Liability) *(Tag: `wrong_account`)*
  - `opt_notes_payable`: Notes Payable (Liability) *(Tag: `wrong_account`)*
- **Causal Hint:** Which contra-liability account holds the unamortized discount and must be credited so that the bond carrying value accretes toward par?
- **Explanation:** Discount on Bonds Payable is credited to amortize a portion of the discount, bringing the net carrying value closer to face value.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the semiannual interest payment and discount amortization?
- **Options:**
  - `opt_balanced_amort`: Debit Interest Expense $38,000.00 / Credit Cash $30,000.00 / Credit Discount on Bonds Payable $8,000.00 *(Correct)*
  - `opt_ignoring_amort`: Debit Interest Expense $30,000.00 / Credit Cash $30,000.00 *(Tag: `m07_ignoring_bond_amortization`)*
  - `opt_reversed_amort`: Debit Interest Expense $38,000.00 / Debit Discount on Bonds Payable $8,000.00 / Credit Cash $30,000.00 *(Tag: `reversed_sides`)*
- **Causal Hint:** Which entry records the effective interest expense (carrying value times market rate), credits cash for the contractual coupon, and credits the discount for the difference?
- **Explanation:** Debit Interest Expense $38,000.00 for effective interest (Carrying Value * market rate / 2), Credit Cash $30,000.00 for contractual coupon (Face * coupon rate / 2), and Credit Discount on Bonds Payable $8,000.00 for the difference. The amortization increases bond carrying value toward face value.

#### Stage 7: Equation Effect
- **Prompt:** How does recording the semiannual interest payment and discount amortization affect the balance sheet equation?
- **Options:**
  - `opt_assets_down_liab_up_eq_down`: Assets decrease by $30,000.00 (-Cash); Liabilities increase by $8,000.00 (+Carrying Value via -Discount); Equity decreases by $38,000.00 (-Interest Expense). *(Correct)*
  - `opt_assets_down_eq_down`: Assets decrease by $38,000.00; Equity decreases by $38,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** Cash left the firm, the contra-liability discount decreased (increasing net debt carrying value), and interest was expensed. What is the net impact on the equation?
- **Explanation:** Assets decrease by cash paid (-Cash $30,000.00); Net Liabilities increase by the discount amortized (+carrying value $8,000.00); Equity decreases by total interest expense (-Interest Expense $38,000.00). Equation check: -$30,000.00 = +$8,000.00 + (-$38,000.00), which balances perfectly.

---

## Question 4: `q_m07_bnd_contingent_liability_04`

- **Title / Topic:** ASC 450 Contingent Liability Accrual (Probable & Reasonably Estimable Loss)
- **Family ID:** `m07_bnd_contingent_liability`
- **Module Scope:** Module 7 (Liabilities & Debt Financing), Module 2 (Double-Entry Mechanics & Accrual Cycle)
- **Status:** `candidate_pending_review`
- **Parameters:**
  - `amount_minor_units`: `[25000000, 50000000, 100000000]` ($250,000, $500,000, $1,000,000)
- **Concepts:** `["contingent_liability_accrual", "asc_450", "probable_and_estimable", "loss_contingency"]`
- **Scenario Template:**
  > "At fiscal year-end, the company's legal counsel confirms that an adverse judgment in a pending intellectual property infringement lawsuit is probable, and the ultimate liability loss is reasonably estimable at ${amount_dollars}. Under ASC 450, the company accrues the loss and corresponding liability."

### Rendered Scaffolding Stages (Example: $250,000 Accrued Loss)

#### Stage 1: Identify Account
- **Prompt:** Under ASC 450, when an adverse litigation judgment against the company is both probable and reasonably estimable, which expense/loss account captures the estimated loss?
- **Options:**
  - `opt_litigation_loss`: Litigation Loss *(Correct)*
  - `opt_accounts_payable`: Accounts Payable *(Tag: `m07_contingent_recorded_as_trade_payable`)*
  - `opt_contingent_liability`: Contingent Liability *(Tag: `wrong_account`)*
  - `opt_prepaid_insurance`: Prepaid Insurance *(Tag: `wrong_account`)*
- **Causal Hint:** Under ASC 450, when an adverse legal outcome is probable and estimable, which income statement loss account recognizes the financial impact?
- **Explanation:** Litigation Loss is the expense/loss account recognized in earnings when a loss contingency meets the ASC 450 accrual threshold.

#### Stage 2: Account Category
- **Prompt:** What category of account is Litigation Loss?
- **Options:**
  - `opt_expense`: Expense (Loss) *(Correct)*
  - `opt_liability`: Liability
  - `opt_asset`: Asset
  - `opt_equity`: Equity
- **Causal Hint:** Does an accrued legal loss represent an immediate reduction in current period earnings, or a future economic obligation?
- **Explanation:** Litigation Loss is an Expense/Loss account reported on the income statement that directly reduces pre-tax operating income and net earnings.

#### Stage 3: Direction
- **Prompt:** Does Litigation Loss increase or decrease when the loss contingency is formally accrued?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** When an estimated loss contingency is newly accrued, does the cumulative total of losses for the period increase or decrease?
- **Explanation:** Litigation Loss increases as the newly estimated contingency is recognized in the current financial period.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Expense or Loss account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Losses reduce stockholders' equity and carry normal debit balances. On which side of the ledger do losses increase?
- **Explanation:** Losses and expenses increase on the left side (Debit), reducing Stockholders' Equity.

#### Stage 5: Counter Account
- **Prompt:** Which liability account is credited to record the accrued legal obligation under ASC 450?
- **Options:**
  - `opt_contingent_liability`: Contingent Liability (Liability) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_accounts_payable`: Accounts Payable (Liability) *(Tag: `m07_contingent_recorded_as_trade_payable`)*
- **Causal Hint:** Which liability account reflects the estimated claim obligation arising from the probable loss contingency?
- **Explanation:** Contingent Liability is credited to establish the estimated balance sheet obligation resulting from the pending lawsuit.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to accrue the contingent litigation loss under ASC 450?
- **Options:**
  - `opt_dr_loss_cr_liability`: Debit Litigation Loss $250,000.00 / Credit Contingent Liability $250,000.00 *(Correct)*
  - `opt_dr_loss_cr_cash`: Debit Litigation Loss $250,000.00 / Credit Cash $250,000.00 *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_dr_loss_cr_ap`: Debit Litigation Loss $250,000.00 / Credit Accounts Payable $250,000.00 *(Tag: `m07_contingent_recorded_as_trade_payable`)*
- **Causal Hint:** Which balanced entry records the estimated lawsuit loss in earnings and establishes the contingent liability on the balance sheet?
- **Explanation:** Debit Litigation Loss $250,000.00 and Credit Contingent Liability $250,000.00. Under ASC 450, accrual is mandatory when the loss is both probable and reasonably estimable. If merely reasonably possible, footnote disclosure alone is required.

#### Stage 7: Equation Effect
- **Prompt:** How does accruing a contingent litigation loss affect the fundamental balance sheet equation?
- **Options:**
  - `opt_liab_up_eq_down`: Liabilities increase by $250,000.00 (+Contingent Liability); Equity decreases by $250,000.00 (-Litigation Loss); Assets unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $250,000.00; Liabilities decrease by $250,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain completely unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** A liability obligation was recorded and an expense was recognized in earnings with no cash paid yet. How does this affect Assets, Liabilities, and Equity?
- **Explanation:** Assets are unchanged ($0); Liabilities increase (+Contingent Liability $250,000.00); Equity decreases (-Litigation Loss $250,000.00). Check: 0 = +$250,000.00 + (-$250,000.00), which preserves the fundamental equation.

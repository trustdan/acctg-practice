# Stream-M05 Rendered Stage Review & Pedagogical Preview

This document contains the rendered 7-stage scaffolding preview for all candidate questions produced in **Stream-M05** (Module 5: Inventory & Cost of Goods Sold) under Wave 1 of the Multi-Agent Question Expansion Plan.

---

## Question 1: `q_m05_inv_fifo_costing_01`

- **Title / Topic:** Perpetual FIFO Inventory Costing & COGS Recognition
- **Family ID:** `m05_inv_fifo_costing`
- **Module Scope:** Module 5 (Inventory & Cost of Goods Sold), Module 2 (Double-Entry Mechanics)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[1200000, 2400000, 3600000]` ($12,000, $24,000, $36,000)
- **Concepts:** `["inventory_costing_fifo", "cogs_recognition", "perpetual_inventory_system", "matching_principle"]`
- **Scenario Template:**
  > "Under a perpetual inventory system with rising acquisition costs, the company sells merchandise to a customer and removes inventory using FIFO (First-In, First-Out). The cost assigned to goods sold from the earliest purchase batches is ${amount_dollars}."

### Rendered Scaffolding Stages (Example: $12,000.00)

#### Stage 1: Identify Account
- **Prompt:** Under a perpetual inventory system using FIFO, which expense account records the cost of inventory sold and delivered to customers?
- **Options:**
  - `opt_cogs`: Cost of Goods Sold *(Correct)*
  - `opt_rent_expense`: Rent Expense *(Tag: `wrong_account`)*
  - `opt_service_revenue`: Service Revenue *(Tag: `m05_inventory_recorded_at_selling_price`)*
  - `opt_inventory`: Inventory *(Tag: `wrong_account`)*
- **Causal Hint:** When merchandise is transferred to buyers, which primary expense captures the historical acquisition cost of the sold goods?
- **Explanation:** Cost of Goods Sold is the expense account that records the historical acquisition cost of inventory sold to customers to generate sales revenue.

#### Stage 2: Account Category
- **Prompt:** What category of account is Cost of Goods Sold?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does the cost of merchandise sold reduce net earnings for the accounting period, or represent an unexpired resource?
- **Explanation:** Cost of Goods Sold is an operating Expense account on the income statement that directly reduces Gross Profit and Net Income.

#### Stage 3: Direction
- **Prompt:** Does Cost of Goods Sold increase or decrease as merchandise is delivered to the customer?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** When a new sales transaction occurs, does total expired cost for the period accumulate upward or downward?
- **Explanation:** Cost of Goods Sold increases because additional product cost has expired during the accounting period.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expenses reduce stockholders' equity and carry normal debit balances. Do they increase on the left or the right side?
- **Explanation:** Expenses increase by Debit on the left side, reflecting a decrease in stockholders' equity.

#### Stage 5: Counter Account
- **Prompt:** Which balance sheet asset account is reduced to reflect the departure of merchandise from the warehouse?
- **Options:**
  - `opt_inventory`: Inventory (Asset) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_accounts_payable`: Accounts Payable (Liability) *(Tag: `payable_recorded_for_cash_payment`)*
- **Causal Hint:** Which current asset stored the goods at acquisition cost prior to sale and must now be reduced?
- **Explanation:** Inventory is credited to remove the historical cost of items that have been shipped and sold.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the FIFO cost of goods sold?
- **Options:**
  - `opt_dr_cogs_cr_inv`: Debit Cost of Goods Sold $12,000.00 / Credit Inventory $12,000.00 *(Correct)*
  - `opt_dr_inv_cr_cogs`: Debit Inventory $12,000.00 / Credit Cost of Goods Sold $12,000.00 *(Tag: `reversed_sides`)*
  - `opt_dr_cogs_cr_cash`: Debit Cost of Goods Sold $12,000.00 / Credit Cash $12,000.00 *(Tag: `cash_recorded_when_uncollected`)*
- **Causal Hint:** Which entry matches the expired product cost to expense while crediting the asset account holding the merchandise?
- **Explanation:** Debit Cost of Goods Sold $12,000.00 and Credit Inventory $12,000.00. Under FIFO, the oldest cost layers are expensed first, leaving newer costs in ending inventory.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording FIFO cost of goods sold affect the fundamental balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $12,000.00 (-Inventory); Equity decreases by $12,000.00 (-Cost of Goods Sold); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $12,000.00; Liabilities decrease by $12,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** Inventory decreased on the balance sheet and Cost of Goods Sold reduced net income. How does this affect Assets and Equity?
- **Explanation:** Assets decrease (-Inventory $12,000.00) and Equity decreases (-Cost of Goods Sold $12,000.00); Liabilities are unchanged. In periods of rising prices, FIFO yields lower COGS and higher ending inventory compared to LIFO.

---

## Question 2: `q_m05_inv_lifo_costing_02`

- **Title / Topic:** Perpetual LIFO Inventory Costing & Tax Conformity
- **Family ID:** `m05_inv_lifo_costing`
- **Module Scope:** Module 5 (Inventory & Cost of Goods Sold), Module 2 (Double-Entry Mechanics)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[1500000, 3000000, 4500000]` ($15,000, $30,000, $45,000)
- **Concepts:** `["inventory_costing_lifo", "cogs_recognition", "inflation_tax_shield", "lifo_conformity_rule"]`
- **Scenario Template:**
  > "Under a perpetual inventory system in an inflationary environment, the company delivers inventory to a client and applies LIFO (Last-In, First-Out) costing, transferring the most recent purchase cost layer of ${amount_dollars} to cost of goods sold."

### Rendered Scaffolding Stages (Example: $15,000.00)

#### Stage 1: Identify Account
- **Prompt:** Under a perpetual inventory system using LIFO, which expense account records the most recent inventory acquisition costs matched against revenue?
- **Options:**
  - `opt_cogs`: Cost of Goods Sold *(Correct)*
  - `opt_insurance_expense`: Insurance Expense *(Tag: `wrong_account`)*
  - `opt_service_revenue`: Service Revenue *(Tag: `m05_inventory_recorded_at_selling_price`)*
  - `opt_accounts_receivable`: Accounts Receivable *(Tag: `wrong_account`)*
- **Causal Hint:** When merchandise is sold under LIFO, which expense account captures the most recent batch costs matched against sales revenue?
- **Explanation:** Cost of Goods Sold records the expired cost of inventory sold under LIFO, matching the latest purchase costs against current revenue.

#### Stage 2: Account Category
- **Prompt:** What category of account is Cost of Goods Sold?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does the cost of inventory sold represent an operating expense on the income statement, or a future economic resource on the balance sheet?
- **Explanation:** Cost of Goods Sold is an Expense account that reduces net operating income on the income statement.

#### Stage 3: Direction
- **Prompt:** Does Cost of Goods Sold increase or decrease when recording this LIFO sales transaction?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** Does recognizing cost of goods sold for a new sale transaction increase or decrease total expenses for the period?
- **Explanation:** Cost of Goods Sold increases as the cost of the delivered inventory is recognized in earnings.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expense accounts increase on the normal balance side that reduces equity. Do expenses increase with a debit on the left or a credit on the right?
- **Explanation:** Expenses increase on the left side (Debit), reducing Stockholders' Equity.

#### Stage 5: Counter Account
- **Prompt:** Which asset account is credited to remove the sold merchandise from the company's books?
- **Options:**
  - `opt_inventory`: Inventory (Asset) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_common_stock`: Common Stock (Equity) *(Tag: `wrong_account`)*
- **Causal Hint:** Which asset account is reduced when merchandise physically leaves the warehouse under a perpetual inventory tracking system?
- **Explanation:** Inventory is credited to remove the carrying value of the merchandise sold.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the LIFO cost of goods sold?
- **Options:**
  - `opt_dr_cogs_cr_inv`: Debit Cost of Goods Sold $15,000.00 / Credit Inventory $15,000.00 *(Correct)*
  - `opt_dr_inv_cr_cogs`: Debit Inventory $15,000.00 / Credit Cost of Goods Sold $15,000.00 *(Tag: `reversed_sides`)*
  - `opt_dr_cogs_cr_cash`: Debit Cost of Goods Sold $15,000.00 / Credit Cash $15,000.00 *(Tag: `cash_recorded_when_uncollected`)*
- **Causal Hint:** Which journal entry records the recent cost layers in Cost of Goods Sold and removes the sold units from the asset account?
- **Explanation:** Debit Cost of Goods Sold $15,000.00 and Credit Inventory $15,000.00. During inflation, LIFO matches recent higher costs with current revenue, reducing pre-tax income and taxes.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording LIFO cost of goods sold affect the balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $15,000.00 (-Inventory); Equity decreases by $15,000.00 (-Cost of Goods Sold); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $15,000.00; Liabilities decrease by $15,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets, Liabilities, and Equity remain completely unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** What is the net effect on the fundamental balance sheet equation when inventory is credited and cost of goods sold is debited?
- **Explanation:** Assets decrease (-Inventory $15,000.00) and Equity decreases (-Cost of Goods Sold $15,000.00); Liabilities are unchanged. Note: US GAAP allows LIFO under the LIFO Conformity Rule, but IFRS (IAS 2) strictly prohibits LIFO.

---

## Question 3: `q_m05_inv_nrv_writedown_03`

- **Title / Topic:** Lower of Cost or Net Realizable Value (NRV) Inventory Write-Down (ASC 330)
- **Family ID:** `m05_inv_nrv_writedown`
- **Module Scope:** Module 5 (Inventory & Cost of Goods Sold), Module 2 (Double-Entry Mechanics)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[500000, 1000000, 2000000]` ($5,000, $10,000, $20,000)
- **Concepts:** `["lower_of_cost_or_nrv", "inventory_writedown", "asc_330", "conservatism_principle"]`
- **Scenario Template:**
  > "At financial statement date, the company performs valuation testing under ASC 330. Due to obsolescence, ending inventory has a net realizable value below carrying cost, requiring an immediate write-down of ${amount_dollars}."

### Rendered Scaffolding Stages (Example: $5,000.00)

#### Stage 1: Identify Account
- **Prompt:** Under US GAAP (ASC 330), when inventory net realizable value falls below cost due to damage or obsolescence, which primary expense/cost account absorbs the write-down?
- **Options:**
  - `opt_cogs`: Cost of Goods Sold *(Correct)*
  - `opt_bad_debt_expense`: Bad Debt Expense *(Tag: `wrong_account`)*
  - `opt_rent_expense`: Rent Expense *(Tag: `wrong_account`)*
  - `opt_accounts_payable`: Accounts Payable *(Tag: `payable_recorded_for_cash_payment`)*
- **Causal Hint:** Under US GAAP (ASC 330), when inventory utility drops below its recorded cost, which primary cost/expense account commonly absorbs this valuation write-down?
- **Explanation:** Cost of Goods Sold (or a dedicated write-down loss account) absorbs the reduction in inventory value when net realizable value falls below historical cost.

#### Stage 2: Account Category
- **Prompt:** What category of account is Cost of Goods Sold when absorbing an inventory write-down?
- **Options:**
  - `opt_expense`: Expense *(Correct)*
  - `opt_asset`: Asset
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does an inventory valuation write-down represent an expense/loss that reduces current period earnings, or a balance sheet liability?
- **Explanation:** Cost of Goods Sold is an Expense account that reduces current period gross margin and net income.

#### Stage 3: Direction
- **Prompt:** Does Cost of Goods Sold increase or decrease when recognizing this lower of cost or NRV impairment?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** Does recognizing a write-down for impaired or obsolete inventory increase or decrease the total expense reported on the income statement?
- **Explanation:** Cost of Goods Sold increases to absorb the valuation write-down loss.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Expense account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Expense and loss accounts reduce stockholders' equity. On which side is an increase to an expense account entered?
- **Explanation:** Expenses increase by Debit on the left side, reducing equity.

#### Stage 5: Counter Account
- **Prompt:** Which asset account is directly credited to reduce the inventory carrying value to net realizable value?
- **Options:**
  - `opt_inventory`: Inventory (Asset) *(Correct)*
  - `opt_cash`: Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
  - `opt_accounts_payable`: Accounts Payable (Liability) *(Tag: `payable_recorded_for_cash_payment`)*
- **Causal Hint:** Which balance sheet asset account must be directly reduced to establish the new lower cost basis for impaired merchandise?
- **Explanation:** Inventory is credited to reduce its carrying value down to its estimated net realizable value.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the lower of cost or NRV inventory write-down?
- **Options:**
  - `opt_dr_cogs_cr_inv`: Debit Cost of Goods Sold $5,000.00 / Credit Inventory $5,000.00 *(Correct)*
  - `opt_dr_inv_cr_cogs`: Debit Inventory $5,000.00 / Credit Cost of Goods Sold $5,000.00 *(Tag: `reversed_sides`)*
  - `opt_dr_loss_cr_ap`: Debit Cost of Goods Sold $5,000.00 / Credit Accounts Payable $5,000.00 *(Tag: `payable_recorded_for_cash_payment`)*
- **Causal Hint:** Which balanced journal entry recognizes the inventory write-down in earnings and reduces the carrying value of inventory?
- **Explanation:** Debit Cost of Goods Sold $5,000.00 and Credit Inventory $5,000.00. Under US GAAP (ASC 330), writing down inventory establishes a new, lower cost basis.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does recording the lower of cost or NRV write-down affect the balance sheet equation?
- **Options:**
  - `opt_assets_down_eq_down`: Assets decrease by $5,000.00 (-Inventory); Equity decreases by $5,000.00 (-Cost of Goods Sold); Liabilities unchanged. *(Correct)*
  - `opt_assets_down_liab_down`: Assets decrease by $5,000.00; Liabilities decrease by $5,000.00. *(Tag: `equation_effect_missed`)*
  - `opt_no_effect`: Total Assets and Equity remain unchanged. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** How does writing down inventory to its net realizable value affect the balance sheet equation (Assets, Liabilities, and Equity)?
- **Explanation:** Assets decrease (-Inventory $5,000.00) and Equity decreases (-Cost of Goods Sold $5,000.00); Liabilities are unchanged. Under US GAAP, write-downs cannot be reversed even if market prices recover subsequently.

---

## Question 4: `q_m05_inv_freight_in_cap_04`

- **Title / Topic:** Freight-In Inventory Capitalization (FOB Shipping Point)
- **Family ID:** `m05_inv_freight_in_cap`
- **Module Scope:** Module 5 (Inventory & Cost of Goods Sold), Module 2 (Double-Entry Mechanics)
- **Status:** `candidate_pending_review`
- **Parameters:** `amount_minor_units`: `[80000, 160000, 320000]` ($800, $1,600, $3,200)
- **Concepts:** `["freight_in_capitalization", "fob_shipping_point", "inventory_cost_basis", "asset_exchange"]`
- **Scenario Template:**
  > "The company purchases merchandise shipped FOB shipping point and pays ${amount_dollars} cash for freight-in shipping costs required to bring the inventory to its warehouse ready for resale."

### Rendered Scaffolding Stages (Example: $800.00)

#### Stage 1: Identify Account
- **Prompt:** Under US GAAP, when a purchaser pays freight charges to transport goods FOB shipping point to its warehouse, which asset account capitalizes these shipping costs?
- **Options:**
  - `opt_inventory`: Inventory *(Correct)*
  - `opt_rent_expense`: Rent Expense *(Tag: `wrong_account`)*
  - `opt_service_revenue`: Service Revenue *(Tag: `wrong_account`)*
  - `opt_insurance_expense`: Insurance Expense *(Tag: `m05_freight_in_expensed`)*
- **Causal Hint:** Under US GAAP, all costs necessary to get inventory ready and in place for resale are capitalized. Which asset account directly absorbs these inbound freight costs?
- **Explanation:** Inventory is an asset account that capitalizes all acquisition and inbound transportation costs necessary to prepare goods for resale.

#### Stage 2: Account Category
- **Prompt:** What category of account is Inventory?
- **Options:**
  - `opt_asset`: Asset *(Correct)*
  - `opt_expense`: Expense
  - `opt_liability`: Liability
  - `opt_equity`: Equity
- **Causal Hint:** Does capitalized merchandise inventory represent an economic resource owned by the company with future economic value, or an expired expense?
- **Explanation:** Inventory is a Current Asset on the balance sheet representing goods held for sale in the normal course of business.

#### Stage 3: Direction
- **Prompt:** Does Inventory increase or decrease when capitalizing inbound transportation freight charges?
- **Options:**
  - `opt_increase`: Increase *(Correct)*
  - `opt_decrease`: Decrease
- **Causal Hint:** Does adding freight-in shipping charges to the recorded value of inventory increase or decrease the asset account balance?
- **Explanation:** Inventory increases because freight-in is an essential product cost capitalized into the inventory asset base.

#### Stage 4: Debit or Credit
- **Prompt:** Which side of the ledger is recorded to increase an Asset account?
- **Options:**
  - `opt_debit`: Debit (Left side) *(Correct)*
  - `opt_credit`: Credit (Right side) *(Tag: `reversed_sides`)*
- **Causal Hint:** Assets maintain normal debit balances on the left side of the ledger. On which side do assets increase?
- **Explanation:** Assets increase by Debit on the left side of the ledger.

#### Stage 5: Counter Account
- **Prompt:** Which asset account is reduced when paying cash for the freight-in transportation service?
- **Options:**
  - `opt_cash`: Cash (Asset) *(Correct)*
  - `opt_accounts_payable`: Accounts Payable (Liability) *(Tag: `payable_recorded_for_cash_payment`)*
  - `opt_service_revenue`: Service Revenue (Revenue) *(Tag: `wrong_account`)*
- **Causal Hint:** Which asset account is reduced when the company disburses currency to pay the transportation carrier for shipping services?
- **Explanation:** Cash is credited because currency was paid to the shipping carrier.

#### Stage 6: Balanced Entry
- **Prompt:** What is the complete balanced journal entry to record the freight-in inventory capitalization?
- **Options:**
  - `opt_dr_inv_cr_cash`: Debit Inventory $800.00 / Credit Cash $800.00 *(Correct)*
  - `opt_dr_cash_cr_inv`: Debit Cash $800.00 / Credit Inventory $800.00 *(Tag: `reversed_sides`)*
  - `opt_dr_freight_cr_cash`: Debit Delivery Expense $800.00 / Credit Cash $800.00 *(Tag: `m05_freight_in_expensed`)*
- **Causal Hint:** Which journal entry capitalizes the freight-in costs into inventory and records the immediate cash payment?
- **Explanation:** Debit Inventory $800.00 and Credit Cash $800.00. Inbound shipping costs (freight-in) attach to the product as an asset; outbound shipping to customers (freight-out) is a selling expense.

#### Stage 7: Balance Sheet Equation Effect
- **Prompt:** How does capitalizing freight-in paid in cash affect the balance sheet equation?
- **Options:**
  - `opt_asset_swap_no_net_effect`: Assets increase by $800.00 (+Inventory) and decrease by $800.00 (-Cash); Total Assets, Liabilities, and Equity unchanged. *(Correct)*
  - `opt_assets_down_eq_down`: Assets decrease by $800.00; Equity decreases by $800.00. *(Tag: `equation_effect_missed`)*
  - `opt_assets_up_liab_up`: Assets increase by $800.00; Liabilities increase by $800.00. *(Tag: `equation_effect_missed`)*
- **Causal Hint:** When cash is exchanged to increase the capitalized cost basis of inventory, what is the net effect on total assets, liabilities, and equity?
- **Explanation:** Total Assets remain unchanged due to an asset exchange (+Inventory $800.00, -Cash $800.00); Liabilities and Equity are unaffected. Incurring capitalized inventory costs does not hit the income statement until the goods are sold.

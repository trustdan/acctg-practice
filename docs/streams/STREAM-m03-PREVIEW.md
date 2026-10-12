# Stream-M03 Question Preview: Indirect Cash Flow & Liquidity Ratios

**Stream ID:** Stream-M03  
**Module:** Module 3: Statement Analysis & Cash Flows (UW Foster School of Business)  
**Secondary Module:** Module 8: Statement Analysis & Liquidity  
**Status:** Candidate Pending Review (`candidate_pending_review`)  
**Scope:** 4 New Questions (`q_m03_ocf_*`)  

---

## 1. Question `q_m03_ocf_ap_increase_01`: Accounts Payable Increase & Indirect OCF

### Scenario
> "During the year, the company incurs $10,000 of operating rent on account from property suppliers, increasing Accounts Payable. Under the indirect method for the Statement of Cash Flows, how is this operating payable change reconciled to Net Income?"

### Parameter Permutations
- `$5,000` (`500000` minor units)
- `$10,000` (`1000000` minor units)
- `$25,000` (`2500000` minor units)

### Canonical Balanced Journal Entry
- **Debit:** Rent Expense `$10,000.00`
- **Credit:** Accounts Payable `$10,000.00`

### 7-Stage Scaffolding Walkthrough

1. **Stage 1 (Identify Account):**
   - **Prompt:** Which liability account tracks trade obligations owed to suppliers for operating services purchased on credit?
   - **Options:**
     - [x] Accounts Payable
     - [ ] Prepaid Insurance *(Tag: `wrong_account`)*
     - [ ] Service Revenue *(Tag: `wrong_account`)*
     - [ ] Unearned Revenue *(Tag: `advance_confused_with_receivable`)*
   - **Causal Hint:** When operating goods or services are received without immediate cash payment, which account records the debt owed to suppliers?
   - **Explanation:** Accounts Payable represents obligations to pay suppliers for operating goods or services previously received.

2. **Stage 2 (Account Category):**
   - **Prompt:** What category of account is Accounts Payable?
   - **Options:**
     - [x] Liability (Current Liability)
     - [ ] Asset *(Tag: `wrong_account`)*
     - [ ] Expense *(Tag: `expense_recorded_on_loan_repayment`)*
     - [ ] Equity *(Tag: `wrong_account`)*
   - **Causal Hint:** Does an obligation to settle unpaid bills with vendors represent an Asset, a Liability, or an Equity claim?
   - **Explanation:** Accounts Payable is a Current Liability representing short-term claims by trade suppliers.

3. **Stage 3 (Direction):**
   - **Prompt:** Does Accounts Payable increase or decrease when operating expenses are incurred on account?
   - **Options:**
     - [x] Increase
     - [ ] Decrease *(Tag: `reversed_sides`)*
   - **Causal Hint:** Does purchasing services without paying cash right away add to or reduce the total amount owed to trade suppliers?
   - **Explanation:** Accounts Payable increases as new supplier invoices are incurred on credit.

4. **Stage 4 (Debit/Credit Side):**
   - **Prompt:** How is an increase to a Liability (Accounts Payable) recorded?
   - **Options:**
     - [x] Credit (Right side)
     - [ ] Debit (Left side) *(Tag: `reversed_sides`)*
   - **Causal Hint:** Liabilities have a normal credit balance; which side records an increase?
   - **Explanation:** Liabilities increase by Credit (right side).

5. **Stage 5 (Counter Account):**
   - **Prompt:** Which account records the cost of the operating facility used during the period?
   - **Options:**
     - [x] Rent Expense (Expense)
     - [ ] Cash (Asset) *(Tag: `cash_recorded_when_uncollected`)*
     - [ ] Common Stock (Equity) *(Tag: `wrong_account`)*
   - **Causal Hint:** The company utilized leased operating space this period. Which account reflects this expired operating cost?
   - **Explanation:** Rent Expense is debited because the operating service was consumed during this period.

6. **Stage 6 (Balanced Entry):**
   - **Prompt:** What is the complete balanced journal entry for incurring operating rent on account?
   - **Options:**
     - [x] Debit Rent Expense $10,000.00 / Credit Accounts Payable $10,000.00
     - [ ] Debit Rent Expense $10,000.00 / Credit Cash $10,000.00 *(Tag: `payable_recorded_for_cash_payment`)*
     - [ ] Debit Accounts Payable $10,000.00 / Credit Rent Expense $10,000.00 *(Tag: `reversed_sides`)*
   - **Causal Hint:** Which entry records the operating expense incurred while reflecting the unpaid supplier liability?
   - **Explanation:** Debit Rent Expense $10,000.00 to recognize the period cost and Credit Accounts Payable $10,000.00 to record the liability to the landlord.

7. **Stage 7 (Equation Effect & Indirect Cash Flow):**
   - **Prompt:** How does this transaction affect the balance sheet equation, and how is the increase in Accounts Payable treated on the indirect Statement of Cash Flows?
   - **Options:**
     - [x] Liabilities increase (+A/P $10,000.00); Equity decreases (-Rent Expense $10,000.00). On the indirect Statement of Cash Flows, the increase in Accounts Payable is ADDED to Net Income.
     - [ ] Liabilities increase (+A/P $10,000.00); Equity decreases (-Rent Expense $10,000.00). On the indirect Statement of Cash Flows, the increase in Accounts Payable is DEDUCTED from Net Income. *(Tag: `m03_ocf_payable_wrong_direction`)*
     - [ ] Assets decrease (-Cash $10,000.00); Equity decreases (-Rent Expense $10,000.00). Operating cash flows decrease immediately. *(Tag: `equation_effect_missed`)*
   - **Causal Hint:** Net Income was reduced by this rent expense, but was any cash paid out this period? To reconcile Net Income to actual cash flow, should the change in payables be added back or subtracted?
   - **Explanation:** Liabilities increase (+Accounts Payable $10,000.00) and Equity decreases (-Rent Expense $10,000.00). Because Net Income was reduced by the expense but no cash was disbursed, the increase in Accounts Payable must be ADDED back to Net Income in operating cash flows.

---

## 2. Question `q_m03_ocf_prepaid_increase_02`: Prepaid Asset Increase & Indirect OCF

### Scenario
> "The company pays $6,000 cash in advance for insurance coverage spanning the upcoming operating year. Under the indirect method for the Statement of Cash Flows, how is this increase in Prepaid Insurance treated in Operating Activities?"

### Parameter Permutations
- `$3,000` (`300000` minor units)
- `$6,000` (`600000` minor units)
- `$12,000` (`1200000` minor units)

### Canonical Balanced Journal Entry
- **Debit:** Prepaid Insurance `$6,000.00`
- **Credit:** Cash `$6,000.00`

### 7-Stage Scaffolding Walkthrough

1. **Stage 1 (Identify Account):**
   - **Prompt:** Which account records operating benefits paid for in cash today that will be consumed in future periods?
   - **Options:**
     - [x] Prepaid Insurance
     - [ ] Insurance Expense *(Tag: `expense_recorded_on_prepaid_purchase`)*
     - [ ] Accounts Payable *(Tag: `wrong_account`)*
     - [ ] Unearned Revenue *(Tag: `wrong_account`)*
   - **Causal Hint:** When cash is spent today for operating coverage that extends across future periods, which asset account stores this future benefit?
   - **Explanation:** Prepaid Insurance is an asset representing future insurance protection paid for in advance.

2. **Stage 2 (Account Category):**
   - **Prompt:** What category of account is Prepaid Insurance?
   - **Options:**
     - [x] Asset (Current Asset)
     - [ ] Expense *(Tag: `expense_recorded_on_prepaid_purchase`)*
     - [ ] Liability *(Tag: `wrong_account`)*
     - [ ] Equity *(Tag: `wrong_account`)*
   - **Causal Hint:** Does the right to receive future insurance protection without additional payments qualify as an economic resource owned by the company?
   - **Explanation:** Prepaid Insurance is a Current Asset because it provides future economic benefits within the operating cycle.

3. **Stage 3 (Direction):**
   - **Prompt:** Does Prepaid Insurance increase or decrease when paying policy premiums in advance?
   - **Options:**
     - [x] Increase
     - [ ] Decrease *(Tag: `reversed_sides`)*
   - **Causal Hint:** Does acquiring an additional period of advance coverage add to the company's prepaid economic resources?
   - **Explanation:** Prepaid Insurance increases as new coverage is purchased before consumption.

4. **Stage 4 (Debit/Credit Side):**
   - **Prompt:** How is an increase to an Asset (Prepaid Insurance) recorded?
   - **Options:**
     - [x] Debit (Left side)
     - [ ] Credit (Right side) *(Tag: `reversed_sides`)*
   - **Causal Hint:** Assets have a normal debit balance; which side records an increase?
   - **Explanation:** Assets increase by Debit (left side).

5. **Stage 5 (Counter Account):**
   - **Prompt:** Which account decreases when cash is disbursed to the insurance underwriter?
   - **Options:**
     - [x] Cash (Asset)
     - [ ] Insurance Expense (Expense) *(Tag: `expense_recorded_on_prepaid_purchase`)*
     - [ ] Accounts Payable (Liability) *(Tag: `wrong_account`)*
   - **Causal Hint:** What liquid asset was paid out to acquire the insurance contract?
   - **Explanation:** Cash is credited because funds were paid out today.

6. **Stage 6 (Balanced Entry):**
   - **Prompt:** What is the complete balanced journal entry for paying insurance premiums in advance?
   - **Options:**
     - [x] Debit Prepaid Insurance $6,000.00 / Credit Cash $6,000.00
     - [ ] Debit Insurance Expense $6,000.00 / Credit Cash $6,000.00 *(Tag: `expense_recorded_on_prepaid_purchase`)*
     - [ ] Debit Cash $6,000.00 / Credit Prepaid Insurance $6,000.00 *(Tag: `reversed_sides`)*
   - **Causal Hint:** Which entry capitalizes the future benefit as an asset rather than prematurely expensing it, while reflecting the cash outflow?
   - **Explanation:** Debit Prepaid Insurance $6,000.00 to record the unexpired resource and Credit Cash $6,000.00 for the cash disbursement.

7. **Stage 7 (Equation Effect & Indirect Cash Flow):**
   - **Prompt:** How does this payment affect total assets, and how is the increase in Prepaid Insurance treated on the indirect Statement of Cash Flows?
   - **Options:**
     - [x] Total Assets are unchanged (+Prepaid Insurance $6,000.00, -Cash $6,000.00). On the indirect Statement of Cash Flows, the increase in Prepaid Insurance is DEDUCTED from Net Income.
     - [ ] Total Assets are unchanged (+Prepaid Insurance $6,000.00, -Cash $6,000.00); on the indirect Statement of Cash Flows, the increase in Prepaid Insurance is ADDED to Net Income. *(Tag: `m03_ocf_prepaid_wrong_direction`)*
     - [ ] Assets decrease by $6,000.00 and Equity decreases by $6,000.00; operating cash flow is not adjusted. *(Tag: `expense_recorded_on_prepaid_purchase`)*
   - **Causal Hint:** Cash was spent, but was any expense deducted in Net Income this period? To reflect this cash outflow in operating activities under the indirect method, do you add or deduct the increase in prepaids?
   - **Explanation:** Assets change composition (+Prepaid Insurance $6,000.00, -Cash $6,000.00) with no immediate net income effect. Because cash was spent but no expense was deducted on the Income Statement, the increase in Prepaid Insurance must be DEDUCTED from Net Income in the operating section.

---

## 3. Question `q_m03_ocf_unearned_revenue_03`: Unearned Revenue Increase & Indirect OCF

### Scenario
> "The company receives $8,000 cash in advance from a client for operating consulting services to be delivered next quarter. Under the indirect method for the Statement of Cash Flows, how is this increase in Unearned Revenue treated?"

### Parameter Permutations
- `$4,000` (`400000` minor units)
- `$8,000` (`800000` minor units)
- `$16,000` (`1600000` minor units)

### Canonical Balanced Journal Entry
- **Debit:** Cash `$8,000.00`
- **Credit:** Unearned Revenue `$8,000.00`

### 7-Stage Scaffolding Walkthrough

1. **Stage 1 (Identify Account):**
   - **Prompt:** Which account records customer cash received in advance of performing operating services?
   - **Options:**
     - [x] Unearned Revenue
     - [ ] Service Revenue *(Tag: `revenue_recognized_prematurely`)*
     - [ ] Accounts Receivable *(Tag: `advance_confused_with_receivable`)*
     - [ ] Accounts Payable *(Tag: `wrong_account`)*
   - **Causal Hint:** When cash is received from a client before the job is done, which account records the performance obligation owed to that client?
   - **Explanation:** Unearned Revenue represents the liability obligation to perform services or return the customer's funds.

2. **Stage 2 (Account Category):**
   - **Prompt:** What category of account is Unearned Revenue?
   - **Options:**
     - [x] Liability (Current Liability)
     - [ ] Revenue *(Tag: `revenue_recognized_prematurely`)*
     - [ ] Asset *(Tag: `wrong_account`)*
     - [ ] Equity *(Tag: `wrong_account`)*
   - **Causal Hint:** Does an unperformed obligation to deliver services to a customer represent Revenue, an Asset, or a Liability?
   - **Explanation:** Unearned Revenue is a Current Liability representing a performance obligation owed to a customer.

3. **Stage 3 (Direction):**
   - **Prompt:** Does Unearned Revenue increase or decrease when receiving a client retainer in advance?
   - **Options:**
     - [x] Increase
     - [ ] Decrease *(Tag: `reversed_sides`)*
   - **Causal Hint:** Does receiving customer cash for future work add to the company's remaining performance obligations?
   - **Explanation:** Unearned Revenue increases when advance client payments are accepted.

4. **Stage 4 (Debit/Credit Side):**
   - **Prompt:** How is an increase to a Liability (Unearned Revenue) recorded?
   - **Options:**
     - [x] Credit (Right side)
     - [ ] Debit (Left side) *(Tag: `reversed_sides`)*
   - **Causal Hint:** Liabilities have a normal credit balance; which side records an increase?
   - **Explanation:** Liabilities increase by Credit (right side).

5. **Stage 5 (Counter Account):**
   - **Prompt:** Which liquid asset account increases when the customer's advance payment is received?
   - **Options:**
     - [x] Cash (Asset)
     - [ ] Service Revenue (Revenue) *(Tag: `revenue_recognized_prematurely`)*
     - [ ] Accounts Receivable (Asset) *(Tag: `advance_confused_with_receivable`)*
   - **Causal Hint:** What tangible asset came into the company's bank account when the customer paid?
   - **Explanation:** Cash is debited because liquid funds were received today.

6. **Stage 6 (Balanced Entry):**
   - **Prompt:** What is the complete balanced journal entry for collecting a customer retainer in advance?
   - **Options:**
     - [x] Debit Cash $8,000.00 / Credit Unearned Revenue $8,000.00
     - [ ] Debit Cash $8,000.00 / Credit Service Revenue $8,000.00 *(Tag: `revenue_recognized_prematurely`)*
     - [ ] Debit Unearned Revenue $8,000.00 / Credit Cash $8,000.00 *(Tag: `reversed_sides`)*
   - **Causal Hint:** Which entry records the cash inflow while deferring revenue until the performance obligation is met?
   - **Explanation:** Debit Cash $8,000.00 to record the receipt of money and Credit Unearned Revenue $8,000.00 to record the performance liability.

7. **Stage 7 (Equation Effect & Indirect Cash Flow):**
   - **Prompt:** How does this transaction affect the balance sheet equation, and how is the change in Unearned Revenue treated on the indirect Statement of Cash Flows?
   - **Options:**
     - [x] Assets increase (+Cash $8,000.00); Liabilities increase (+Unearned Revenue $8,000.00). On the indirect Statement of Cash Flows, the increase in Unearned Revenue is ADDED to Net Income.
     - [ ] Assets increase (+Cash $8,000.00); Liabilities increase (+Unearned Revenue $8,000.00). On the indirect Statement of Cash Flows, the increase in Unearned Revenue is DEDUCTED from Net Income. *(Tag: `m03_ocf_unearned_revenue_wrong_direction`)*
     - [ ] Assets increase (+Cash $8,000.00); Equity increases (+Revenue $8,000.00). Operating cash flows require no adjustment. *(Tag: `revenue_recognized_prematurely`)*
   - **Causal Hint:** Operating cash was collected today, but was any revenue included in Net Income yet? Under the indirect method, should the increase in this operating liability be added or subtracted to reach operating cash flow?
   - **Explanation:** Assets increase (+Cash $8,000.00) and Liabilities increase (+Unearned Revenue $8,000.00). Because cash was collected from customers but not yet included in Net Income, the increase in Unearned Revenue must be ADDED to Net Income on the indirect Statement of Cash Flows.

---

## 4. Question `q_m03_ocf_quick_ratio_04`: Quick Ratio vs Inventory / Prepaids

### Scenario
> "A company currently possesses a Quick Ratio (Acid-Test Ratio) of 1.5. The company pays $20,000 cash to settle supplier balances in Accounts Payable. How does this payment affect the Quick Ratio, and why does this contrast with purchasing inventory?"

### Parameter Permutations
- `$10,000` (`1000000` minor units)
- `$20,000` (`2000000` minor units)
- `$40,000` (`4000000` minor units)

### Canonical Balanced Journal Entry
- **Debit:** Accounts Payable `$20,000.00`
- **Credit:** Cash `$20,000.00`

### 7-Stage Scaffolding Walkthrough

1. **Stage 1 (Identify Account):**
   - **Prompt:** Which liability account is reduced when cash is disbursed to satisfy outstanding supplier balances?
   - **Options:**
     - [x] Accounts Payable
     - [ ] Prepaid Insurance *(Tag: `wrong_account`)*
     - [ ] Notes Payable *(Tag: `wrong_account`)*
     - [ ] Service Revenue *(Tag: `wrong_account`)*
   - **Causal Hint:** Which trade liability account records obligations owed to suppliers for past credit purchases?
   - **Explanation:** Accounts Payable is the liability account representing debts owed to trade suppliers.

2. **Stage 2 (Account Category):**
   - **Prompt:** What category of account is Accounts Payable?
   - **Options:**
     - [x] Liability (Current Liability)
     - [ ] Asset *(Tag: `wrong_account`)*
     - [ ] Expense *(Tag: `expense_recorded_on_loan_repayment`)*
     - [ ] Equity *(Tag: `wrong_account`)*
   - **Causal Hint:** Is an obligation owed to trade suppliers an Asset, a Liability, or an Equity account?
   - **Explanation:** Accounts Payable is a Current Liability and represents a short-term obligation included in the denominator of both Current and Quick ratios.

3. **Stage 3 (Direction):**
   - **Prompt:** Does Accounts Payable increase or decrease when the company pays supplier invoices?
   - **Options:**
     - [x] Decrease
     - [ ] Increase *(Tag: `reversed_sides`)*
   - **Causal Hint:** Does sending payments to vendors reduce the total debt owed to them?
   - **Explanation:** Accounts Payable decreases as the company satisfies its vendor obligations.

4. **Stage 4 (Debit/Credit Side):**
   - **Prompt:** How is a decrease to a Liability (Accounts Payable) recorded?
   - **Options:**
     - [x] Debit (Left side)
     - [ ] Credit (Right side) *(Tag: `reversed_sides`)*
   - **Causal Hint:** Liabilities increase on the credit side; which side records a reduction?
   - **Explanation:** Liabilities decrease by Debit (left side).

5. **Stage 5 (Counter Account):**
   - **Prompt:** Which quick asset account decreases when supplier bills are paid?
   - **Options:**
     - [x] Cash (Asset)
     - [ ] Accounts Receivable (Asset) *(Tag: `wrong_account`)*
     - [ ] Prepaid Insurance (Asset) *(Tag: `wrong_account`)*
   - **Causal Hint:** What liquid asset is disbursed from the bank account to pay suppliers?
   - **Explanation:** Cash is credited because liquid funds were spent to pay the suppliers.

6. **Stage 6 (Balanced Entry):**
   - **Prompt:** What is the complete balanced journal entry for paying Accounts Payable in cash?
   - **Options:**
     - [x] Debit Accounts Payable $20,000.00 / Credit Cash $20,000.00
     - [ ] Debit Cash $20,000.00 / Credit Accounts Payable $20,000.00 *(Tag: `reversed_sides`)*
     - [ ] Debit Rent Expense $20,000.00 / Credit Cash $20,000.00 *(Tag: `expense_recorded_on_loan_repayment`)*
   - **Causal Hint:** Which entry reduces the supplier liability while recording the outflow of liquid cash?
   - **Explanation:** Debit Accounts Payable $20,000.00 to satisfy the obligation and Credit Cash $20,000.00 to record the disbursement.

7. **Stage 7 (Equation Effect & Quick Ratio Mechanics):**
   - **Prompt:** Assuming the initial Quick Ratio exceeds 1.0 (e.g., 1.5), how does paying Accounts Payable in cash affect the Quick Ratio, and how does this contrast with purchasing inventory with cash?
   - **Options:**
     - [x] The Quick Ratio increases (subtracting equal dollars from numerator and denominator when ratio > 1.0 increases the ratio). In contrast, buying inventory with cash reduces the Quick Ratio because Inventory is excluded from Quick Assets.
     - [ ] The Quick Ratio decreases because liquid cash was spent, just as when inventory is purchased. *(Tag: `m03_quick_ratio_wrong_direction`)*
     - [ ] The Quick Ratio is unchanged because Cash and Accounts Payable offset, and inventory purchases also have no effect on the Quick Ratio. *(Tag: `m03_inventory_treated_as_quick_asset`)*
   - **Causal Hint:** Quick Ratio = (Cash + Marketable Securities + Receivables) / Current Liabilities. When ratio > 1.0, what happens mathematically when you subtract equal amounts from numerator and denominator, and why is inventory excluded from quick assets?
   - **Explanation:** Paying Accounts Payable subtracts equal dollars from Quick Assets (-Cash $20,000.00) and Current Liabilities (-A/P $20,000.00). For any ratio > 1.0 (e.g. 150/100 = 1.5 -> 130/80 = 1.625), subtracting equal amounts INCREASES the ratio. In contrast, purchasing inventory with cash converts Cash (a quick asset) into Inventory (a non-quick asset), reducing Quick Assets without changing Current Liabilities, which causes the Quick Ratio to decrease.

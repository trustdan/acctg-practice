# Stream-M04 Rendered Preview: ASC 606 Revenue Recognition

Stream ID: `Stream-M04`  
Module Scope: Module 4 (Revenue Recognition / ASC 606)  
Target Topics: Variable Consideration Constraint, Contract Asset vs. Accounts Receivable, Right of Return Accounting, Over-Time Performance Obligation Satisfaction.  
Sample Evaluation Amount: $2,400.00 (240,000 cents), Seed: 4001.

---

## 1. q_m04_rev_variable_consideration_01

### Scenario
A consulting firm completes technical advisory services under a contract with a fixed fee of $2,400.00 and an incentive bonus of $2,000 contingent on client cost savings. Under the ASC 606 variable consideration constraint, the bonus is excluded from the transaction price because meeting the cost savings target is volatile and outside the firm's influence. The firm bills the client on account for the fixed fee of $2,400.00 with net 30 payment terms. No cash has been collected. Record the billing of the fixed fee under ASC 606.

### Stages

#### Stage 1: identify_account
- **Prompt:** When the firm finishes services and bills the client on account for the fixed fee while an uncertain incentive bonus is constrained under ASC 606, which primary account is adjusted for the uncollected claim?
- **Causal Hint:** Has the customer already transferred liquid funds today, or does the firm hold an unpaid claim on account for completed work?
- **Explanation:** Accounts Receivable records the unpaid claim for services performed where the customer has been billed on account.
- **Options:**
  - `Accounts Receivable` *(id=opt_ar; correct=true)*
  - `Cash` *(id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)*
  - `Unearned Revenue` *(id=opt_unearned; correct=false; tag=advance_confused_with_receivable)*

#### Stage 2: account_category
- **Prompt:** What classification of account is Accounts Receivable?
- **Causal Hint:** Does an unpaid customer invoice represent an economic resource with future cash inflows, or an obligation owed to someone else?
- **Explanation:** Accounts Receivable is an Asset representing a legally enforceable claim to receive cash in the future.
- **Options:**
  - `Asset` *(id=opt_asset; correct=true)*
  - `Revenue` *(id=opt_revenue; correct=false; tag=wrong_account)*
  - `Liability` *(id=opt_liability; correct=false; tag=wrong_account)*
  - `Equity` *(id=opt_equity; correct=false; tag=wrong_account)*

#### Stage 3: direction
- **Prompt:** Does Accounts Receivable increase or decrease when the firm bills the customer for completed services?
- **Causal Hint:** As newly billed amounts are added to client balances, does total money owed to the firm increase or decrease?
- **Explanation:** Accounts Receivable increases because a new customer claim has been established.
- **Options:**
  - `Increase` *(id=opt_increase; correct=true)*
  - `Decrease` *(id=opt_decrease; correct=false; tag=reversed_sides)*

#### Stage 4: debit_credit
- **Prompt:** Which side is recorded to increase an Asset account such as Accounts Receivable?
- **Causal Hint:** Assets have normal debit balances on the left side of the ledger. Which side increases an asset?
- **Explanation:** Assets increase on the left side (Debit).
- **Options:**
  - `Debit (Left side)` *(id=opt_debit; correct=true)*
  - `Credit (Right side)` *(id=opt_credit; correct=false; tag=reversed_sides)*

#### Stage 5: counter_account
- **Prompt:** Under the ASC 606 variable consideration constraint, which earned account is credited for the fixed transaction price?
- **Causal Hint:** Because the performance obligation was completed and the fixed fee is not subject to significant reversal, what account records earnings for the period?
- **Explanation:** Service Revenue is credited because the performance obligation has been fulfilled and the fixed transaction price is unconstrained.
- **Options:**
  - `Service Revenue (Revenue / Equity)` *(id=opt_service_revenue; correct=true)*
  - `Unearned Revenue (Liability)` *(id=opt_unearned_revenue; correct=false; tag=revenue_deferred_when_earned)*
  - `Cash (Asset)` *(id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)*

#### Stage 6: balanced_entry
- **Prompt:** What is the complete balanced journal entry to record the billing of the fixed fee under ASC 606?
- **Causal Hint:** Which entry pairs the increase in the asset claim with the increase in earned revenue for the period?
- **Explanation:** Debit Accounts Receivable $2,400.00 increases the asset; Credit Service Revenue $2,400.00 increases revenue and equity. The uncertain bonus remains excluded until meeting the constraint criteria.
- **Options:**
  - `Debit Accounts Receivable $2,400.00 / Credit Service Revenue $2,400.00` *(id=opt_dr_ar_cr_rev; correct=true)*
  - `Debit Cash $2,400.00 / Credit Service Revenue $2,400.00` *(id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)*
  - `Debit Service Revenue $2,400.00 / Credit Accounts Receivable $2,400.00` *(id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)*

#### Stage 7: equation_effect
- **Prompt:** How does billing the fixed fee on account affect the balance sheet equation?
- **Causal Hint:** An asset increased by the billed amount and earned revenue increased net income. How does this impact total assets and stockholders' equity?
- **Explanation:** Assets increase (+Accounts Receivable $2,400.00) and Equity increases (+Service Revenue $2,400.00 via net income). Liabilities remain unchanged.
- **Options:**
  - `Assets increase by $2,400.00 (+Accounts Receivable); Equity increases by $2,400.00 (+Service Revenue); Liabilities unchanged.` *(id=opt_assets_up_eq_up; correct=true)*
  - `Assets increase by $2,400.00; Liabilities increase by $2,400.00.` *(id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)*
  - `Total Assets and Equity remain completely unchanged.` *(id=opt_no_effect; correct=false; tag=equation_effect_missed)*

---

## 2. q_m04_rev_contract_asset_ar_01

### Scenario
Under ASC 606-10-45, an entity distinguishes between a contract asset (conditional right) and an account receivable (unconditional right). A systems engineering firm successfully completes a milestone that has received full client sign-off. The contract stipulates that invoice billing is now unconditional, subject only to the passage of time before payment is due on net 30 credit terms. The firm issues an invoice on account for $3,000.00. Record the establishment of this unconditional claim under ASC 606.

### Stages

#### Stage 1: identify_account
- **Prompt:** Under ASC 606-10-45, when an entity's right to payment is unconditional (dependent only upon the passage of time), which account is recognized for the billed claim?
- **Causal Hint:** When performance obligations are fully certified and payment is subject only to the passage of time, which standard asset account records unconditional client claims?
- **Explanation:** Accounts Receivable is the required presentation under ASC 606 when an entity possesses an unconditional right to consideration, unlike a contract asset which depends on future performance.
- **Options:**
  - `Accounts Receivable` *(id=opt_ar; correct=true)*
  - `Unearned Revenue` *(id=opt_unearned; correct=false; tag=advance_confused_with_receivable)*
  - `Cash` *(id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)*

#### Stage 2: account_category
- **Prompt:** What category of account is Accounts Receivable under ASC 606 presentation rules?
- **Causal Hint:** Does an unconditional legal right to collect payment from a customer constitute an owned economic resource or an obligation?
- **Explanation:** Accounts Receivable is a Current Asset representing an unconditional right to consideration.
- **Options:**
  - `Asset` *(id=opt_asset; correct=true)*
  - `Liability` *(id=opt_liability; correct=false; tag=wrong_account)*
  - `Equity` *(id=opt_equity; correct=false; tag=wrong_account)*
  - `Revenue` *(id=opt_revenue; correct=false; tag=wrong_account)*

#### Stage 3: direction
- **Prompt:** Does Accounts Receivable increase or decrease when the firm satisfies its milestone and establishes an unconditional billing?
- **Causal Hint:** When a new unconditional claim is established on net 30 credit terms, does the total receivable balance increase or decrease?
- **Explanation:** Accounts Receivable increases as new unconditional claims are established.
- **Options:**
  - `Increase` *(id=opt_increase; correct=true)*
  - `Decrease` *(id=opt_decrease; correct=false; tag=reversed_sides)*

#### Stage 4: debit_credit
- **Prompt:** Which side is recorded to increase the asset Accounts Receivable?
- **Causal Hint:** Assets carry normal debit balances. Which side increases an asset?
- **Explanation:** Assets increase on the left side (Debit).
- **Options:**
  - `Debit (Left side)` *(id=opt_debit; correct=true)*
  - `Credit (Right side)` *(id=opt_credit; correct=false; tag=reversed_sides)*

#### Stage 5: counter_account
- **Prompt:** When an unconditional receivable is recorded upon fulfilling performance obligations under ASC 606, which revenue account is credited?
- **Causal Hint:** What account reflects the value of completed customer services earned during the operating period?
- **Explanation:** Service Revenue is credited to recognize earnings upon the satisfaction of performance obligations.
- **Options:**
  - `Service Revenue (Revenue / Equity)` *(id=opt_service_revenue; correct=true)*
  - `Accounts Payable (Liability)` *(id=opt_accounts_payable; correct=false; tag=wrong_account)*
  - `Cash (Asset)` *(id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)*

#### Stage 6: balanced_entry
- **Prompt:** What is the complete balanced journal entry to record the unconditional billing under ASC 606?
- **Causal Hint:** Which journal entry increases the unconditional asset claim and increases earned revenue by identical amounts?
- **Explanation:** Debit Accounts Receivable $3,000.00 increases assets; Credit Service Revenue $3,000.00 increases revenue and equity. Under ASC 606-10-45, this unconditional claim must be presented as a receivable rather than a contract asset.
- **Options:**
  - `Debit Accounts Receivable $3,000.00 / Credit Service Revenue $3,000.00` *(id=opt_dr_ar_cr_rev; correct=true)*
  - `Debit Service Revenue $3,000.00 / Credit Accounts Receivable $3,000.00` *(id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)*
  - `Debit Cash $3,000.00 / Credit Service Revenue $3,000.00` *(id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)*

#### Stage 7: equation_effect
- **Prompt:** How does recording the unconditional billing affect the fundamental balance sheet equation?
- **Causal Hint:** Receivables increased total assets, while service revenue flowed into retained earnings. How does this shift total assets and equity?
- **Explanation:** Assets increase by $3,000.00 (+Accounts Receivable) and Equity increases by $3,000.00 (+Service Revenue via retained earnings), preserving equation balance.
- **Options:**
  - `Assets increase by $3,000.00 (+Accounts Receivable); Equity increases by $3,000.00 (+Service Revenue); Liabilities unchanged.` *(id=opt_assets_up_eq_up; correct=true)*
  - `Assets increase by $3,000.00; Liabilities increase by $3,000.00.` *(id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)*
  - `Total Assets and Equity remain unchanged.` *(id=opt_no_effect; correct=false; tag=equation_effect_missed)*

---

## 3. q_m04_rev_right_of_return_01

### Scenario
Under ASC 606-10-55-23, when a contract includes a right of return and returns cannot be reliably estimated prior to performance, consideration received is recorded as a contract liability (Unearned Revenue) rather than revenue. A digital education provider receives $3,600.00 cash today for an annual executive seminar series starting next month, accompanied by an unconditional 30-day money-back guarantee. No services have been performed. Record today's cash receipt under ASC 606.

### Stages

#### Stage 1: identify_account
- **Prompt:** When a customer pays cash in advance for an advisory package that includes a 30-day full money-back satisfaction guarantee prior to service delivery, which liquid asset account is adjusted for the inflow?
- **Causal Hint:** What account directly tracks the company's checking account balance and immediately available legal funds?
- **Explanation:** Cash records the immediate inflow of liquid currency received from the customer.
- **Options:**
  - `Cash` *(id=opt_cash; correct=true)*
  - `Service Revenue` *(id=opt_service_revenue; correct=false; tag=m04_revenue_recognized_with_uncertain_return)*
  - `Accounts Receivable` *(id=opt_accounts_receivable; correct=false; tag=advance_confused_with_receivable)*

#### Stage 2: account_category
- **Prompt:** What classification of account is Cash?
- **Causal Hint:** Does money held in the bank represent an owned economic resource or an obligation owed to an outside party?
- **Explanation:** Cash is a Current Asset representing the most liquid store of economic purchasing power.
- **Options:**
  - `Asset` *(id=opt_asset; correct=true)*
  - `Liability` *(id=opt_liability; correct=false; tag=wrong_account)*
  - `Revenue` *(id=opt_revenue; correct=false; tag=wrong_account)*
  - `Equity` *(id=opt_equity; correct=false; tag=wrong_account)*

#### Stage 3: direction
- **Prompt:** Does Cash increase or decrease upon receiving the customer's upfront payment?
- **Causal Hint:** When a customer transfers money into the corporate checking account, does total cash increase or decrease?
- **Explanation:** Cash increases as new funds are deposited into the company's account.
- **Options:**
  - `Increase` *(id=opt_increase; correct=true)*
  - `Decrease` *(id=opt_decrease; correct=false; tag=reversed_sides)*

#### Stage 4: debit_credit
- **Prompt:** Which side is recorded to increase an Asset account such as Cash?
- **Causal Hint:** Assets reside on the left side of the fundamental accounting equation. Which side increases an asset?
- **Explanation:** Assets increase on the left side (Debit).
- **Options:**
  - `Debit (Left side)` *(id=opt_debit; correct=true)*
  - `Credit (Right side)` *(id=opt_credit; correct=false; tag=reversed_sides)*

#### Stage 5: counter_account
- **Prompt:** Under ASC 606, because services have not yet been performed and the payment is subject to full return/refund terms, which contract liability account is credited?
- **Causal Hint:** Because the company owes future performance or a complete refund to the customer, what liability account captures unperformed obligations?
- **Explanation:** Unearned Revenue is credited because the company has a contract liability to provide service or return the consideration.
- **Options:**
  - `Unearned Revenue (Liability)` *(id=opt_unearned_revenue; correct=true)*
  - `Service Revenue (Revenue)` *(id=opt_service_revenue; correct=false; tag=m04_revenue_recognized_with_uncertain_return)*
  - `Accounts Payable (Liability)` *(id=opt_accounts_payable; correct=false; tag=wrong_account)*

#### Stage 6: balanced_entry
- **Prompt:** What is the complete balanced journal entry to record the customer advance subject to right-of-return terms?
- **Causal Hint:** Which entry pairs the debit for the liquid asset received with the credit for the contract liability obligation?
- **Explanation:** Debit Cash $3,600.00 increases assets; Credit Unearned Revenue $3,600.00 increases contract liabilities. Revenue cannot be recognized before performance obligations are satisfied and while return uncertainty remains open.
- **Options:**
  - `Debit Cash $3,600.00 / Credit Unearned Revenue $3,600.00` *(id=opt_dr_cash_cr_unearned; correct=true)*
  - `Debit Cash $3,600.00 / Credit Service Revenue $3,600.00` *(id=opt_dr_cash_cr_rev; correct=false; tag=m04_revenue_recognized_with_uncertain_return)*
  - `Debit Unearned Revenue $3,600.00 / Credit Cash $3,600.00` *(id=opt_dr_unearned_cr_cash; correct=false; tag=reversed_sides)*

#### Stage 7: equation_effect
- **Prompt:** How does receiving this upfront advance subject to refund terms affect the accounting equation?
- **Causal Hint:** Cash was collected, creating an obligation to perform or refund. What happens to total assets and total liabilities?
- **Explanation:** Assets increase (+Cash $3,600.00) and Liabilities increase (+Unearned Revenue $3,600.00). Equity is completely unaffected until performance occurs.
- **Options:**
  - `Assets increase by $3,600.00 (+Cash); Liabilities increase by $3,600.00 (+Unearned Revenue); Equity unchanged.` *(id=opt_assets_up_liab_up; correct=true)*
  - `Assets increase by $3,600.00; Equity increases by $3,600.00.` *(id=opt_assets_up_eq_up; correct=false; tag=equation_effect_missed)*
  - `Total Assets and Liabilities remain unchanged.` *(id=opt_no_effect; correct=false; tag=equation_effect_missed)*

---

## 4. q_m04_rev_satisfaction_over_time_01

### Scenario
Under ASC 606 Step 5, revenue is recognized over time as the customer simultaneously receives and consumes benefits from performance. A corporate client previously paid in advance for six months of ongoing compliance advisory, recorded in Unearned Revenue. During the current month, the firm delivers ongoing compliance services over time, satisfying $4,000.00 of its performance obligation. No cash changes hands today. Record the revenue earned and contract liability relief for the month under ASC 606.

### Stages

#### Stage 1: identify_account
- **Prompt:** Under ASC 606 Step 5, when an entity delivers consulting services over time to fulfill a performance obligation previously paid in advance, which existing liability account is adjusted first?
- **Causal Hint:** What existing liability account was established when the client previously deposited funds before work began?
- **Explanation:** Unearned Revenue represents the contract liability that must be reduced as the firm satisfies its performance obligation over time.
- **Options:**
  - `Unearned Revenue` *(id=opt_unearned_revenue; correct=true)*
  - `Cash` *(id=opt_cash; correct=false; tag=m04_cash_debited_on_over_time_satisfaction)*
  - `Accounts Payable` *(id=opt_accounts_payable; correct=false; tag=wrong_account)*

#### Stage 2: account_category
- **Prompt:** What category of account is Unearned Revenue?
- **Causal Hint:** Does receiving customer payments in advance create an obligation to perform future services, or an asset?
- **Explanation:** Unearned Revenue is a Liability representing a performance obligation owed to a customer.
- **Options:**
  - `Liability` *(id=opt_liability; correct=true)*
  - `Revenue` *(id=opt_revenue; correct=false; tag=wrong_account)*
  - `Asset` *(id=opt_asset; correct=false; tag=wrong_account)*
  - `Equity` *(id=opt_equity; correct=false; tag=wrong_account)*

#### Stage 3: direction
- **Prompt:** Does Unearned Revenue increase or decrease when the firm satisfies its performance obligation over time?
- **Causal Hint:** As promised services are progressively performed and delivered, does the firm's remaining obligation to the client grow or shrink?
- **Explanation:** Unearned Revenue decreases because the obligation has been partially or fully satisfied through performance.
- **Options:**
  - `Decrease` *(id=opt_decrease; correct=true)*
  - `Increase` *(id=opt_increase; correct=false; tag=reversed_sides)*

#### Stage 4: debit_credit
- **Prompt:** Which side is recorded to decrease a Liability account such as Unearned Revenue?
- **Causal Hint:** Liabilities have normal credit balances. Which side reduces a liability?
- **Explanation:** Liabilities decrease on the left side (Debit).
- **Options:**
  - `Debit (Left side)` *(id=opt_debit; correct=true)*
  - `Credit (Right side)` *(id=opt_credit; correct=false; tag=reversed_sides)*

#### Stage 5: counter_account
- **Prompt:** When the contract liability is relieved upon delivering services over time, which account is credited to recognize earnings?
- **Causal Hint:** What account captures the value of customer services performed and earned during the reporting period?
- **Explanation:** Service Revenue is credited to record revenue earned as performance obligations are satisfied.
- **Options:**
  - `Service Revenue (Revenue / Equity)` *(id=opt_service_revenue; correct=true)*
  - `Cash (Asset)` *(id=opt_cash; correct=false; tag=m04_cash_debited_on_over_time_satisfaction)*
  - `Accounts Receivable (Asset)` *(id=opt_accounts_receivable; correct=false; tag=wrong_account)*

#### Stage 6: balanced_entry
- **Prompt:** What is the complete balanced journal entry to record satisfying the performance obligation over time?
- **Causal Hint:** Which entry pairs the debit reducing the contract liability with the credit recognizing earned revenue?
- **Explanation:** Debit Unearned Revenue $4,000.00 decreases liabilities; Credit Service Revenue $4,000.00 increases revenue and equity. Under ASC 606 Step 5, fulfilling the obligation over time converts the contract liability into earned revenue.
- **Options:**
  - `Debit Unearned Revenue $4,000.00 / Credit Service Revenue $4,000.00` *(id=opt_dr_unearned_cr_rev; correct=true)*
  - `Debit Cash $4,000.00 / Credit Service Revenue $4,000.00` *(id=opt_dr_cash_cr_rev; correct=false; tag=m04_cash_debited_on_over_time_satisfaction)*
  - `Debit Service Revenue $4,000.00 / Credit Unearned Revenue $4,000.00` *(id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)*

#### Stage 7: equation_effect
- **Prompt:** How does satisfying this performance obligation over time affect the balance sheet equation?
- **Causal Hint:** A contract liability was fulfilled without exchanging cash, while revenue increased equity. How do total liabilities and equity change?
- **Explanation:** Liabilities decrease (-Unearned Revenue $4,000.00) and Equity increases (+Service Revenue $4,000.00 via net income). Assets remain completely unaffected.
- **Options:**
  - `Liabilities decrease by $4,000.00 (-Unearned Revenue); Equity increases by $4,000.00 (+Service Revenue); Assets unchanged.` *(id=opt_liab_down_eq_up; correct=true)*
  - `Assets increase by $4,000.00; Equity increases by $4,000.00.` *(id=opt_assets_up_eq_up; correct=false; tag=equation_effect_missed)*
  - `Total Liabilities and Equity remain unchanged.` *(id=opt_no_effect; correct=false; tag=equation_effect_missed)*

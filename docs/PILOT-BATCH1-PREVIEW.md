# Pilot batch 1: complete teaching preview

Reviewed snapshots for six version-1 scenarios, shown at $50 with seed 100. Amounts allowed: $50/$100/$200. Stage 2-4 teaching is shared with the reviewed family; stage 1/5/6/7 text has per-scenario overrides. This artifact includes canonical choices, misconception tags, and snapshotted option-specific hints. Provenance and contrast partners: [PILOT-BATCH1-REVIEW.md](PILOT-BATCH1-REVIEW.md). Approved package: [PILOT-BATCH1.json](PILOT-BATCH1.json).

## cash_service_training v1

A training company finishes a one-day workshop today. Participants pay $50 cash at the end of the workshop today. No participant was invoiced or paid before today, and no work remains for these payments.

### identify_account

Which account reflects the immediate payment received today?

Hint: What did the participants hand over at the end of today's workshop?

Explanation: Cash records the participants' payment today. Accounts Receivable would mean payment is still owed; Unearned Revenue would mean the workshop is still owed. Neither remains outstanding.

- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)
- opt_cash: Cash (tag: ; hint: )
- opt_ar: Accounts Receivable (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- opt_asset: Asset (tag: ; hint: )
- opt_revenue: Revenue (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )

### direction

Does Cash increase or decrease upon receiving this payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase in an Asset (Cash) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

Which account records what the company did for the customer in exchange for the cash?

Hint: Was the workshop completed before the participants paid, or is the company still committed to provide it later?

Explanation: Service Revenue records the completed workshop. The cash was earned today; Unearned Revenue would imply that this workshop still needs to be provided.

- opt_ap: Accounts Payable (tag: wrong_account; hint: Was the workshop completed before the participants paid, or is the company still committed to provide it later?)
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: Was the workshop completed before the participants paid, or is the company still committed to provide it later?)
- opt_service_rev: Service Revenue (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: The workshop is finished and its participants paid today. What resource came in, and what was earned?

Explanation: Debit Cash $50 and Credit Service Revenue $50. There is no unpaid invoice or future workshop obligation for this payment.

- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: wrong_account; hint: The workshop is finished and its participants paid today. What resource came in, and what was earned?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_rev_cr_cash: Debit Service Revenue $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_unearned: Debit Cash $50 / Credit Unearned Revenue $50 (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: For this completed, paid workshop, does anyone still owe money or work?

Explanation: Assets increase by $50 (Cash) and equity increases by $50 through Service Revenue. No liability remains for this workshop.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: For this completed, paid workshop, does anyone still owe money or work?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity is unchanged. (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Service Revenue); Liabilities are unchanged. (tag: ; hint: )

## customer_advance_training v1

A training company receives $50 cash today from participants reserving places at a workshop next month. None of the instruction has been provided, and the company must still deliver the full workshop. No payment or related obligation was recorded earlier.

### identify_account

What did the company receive today, and which account records it?

Hint: What did the participants hand over today, before attending any instruction?

Explanation: Cash records the payment received today. Receiving money does not mean the workshop has been earned; all instruction still lies ahead.

- opt_ap: Accounts Payable (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)
- opt_service_rev: Service Revenue (tag: revenue_recognized_prematurely; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)
- opt_ar: Accounts Receivable (tag: advance_confused_with_receivable; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)
- opt_cash: Cash (tag: ; hint: )

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- opt_revenue: Revenue (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )

### direction

Does Cash increase or decrease upon receiving this customer payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- opt_decrease: Decrease (tag: ; hint: )
- opt_increase: Increase (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Cash) recorded in double-entry bookkeeping?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

What counter-account balances this entry for services to be performed next month?

Hint: After accepting these reservations, what must the company still provide to the participants?

Explanation: Unearned Revenue records the company's obligation to deliver the workshop. Service Revenue would count instruction that has not happened; Accounts Receivable would claim the participants have not paid.

- opt_ar: Accounts Receivable (Asset) (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)
- opt_unearned_rev: Unearned Revenue (Liability) (tag: ; hint: )
- opt_common_stock: Common Stock (Equity) (tag: wrong_account; hint: After accepting these reservations, what must the company still provide to the participants?)
- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: The participants paid today but will attend next month. What resource arrived, and what promise remains?

Explanation: Debit Cash $50 and Credit Unearned Revenue $50. The company has more cash and owes the workshop; it has not yet earned this payment.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)
- opt_dr_cash_cr_unearned: Debit Cash $50 / Credit Unearned Revenue $50 (tag: ; hint: )
- opt_dr_unearned_cr_cash: Debit Unearned Revenue $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Has reserving a place earned the company's payment, or does the company still owe all the instruction?

Explanation: Assets increase by $50 (Cash) and liabilities increase by $50 (Unearned Revenue). Equity is unchanged because no instruction has been delivered.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Has reserving a place earned the company's payment, or does the company still owe all the instruction?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity is unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue); Liabilities are unchanged. (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)

## earn_advance_training v1

A training company finishes a one-day workshop today. Participants paid $50 last month, and the company recorded its obligation to provide this workshop then. Today's completed instruction satisfies the whole obligation for those payments. No cash changes hands today.

### identify_account

No cash changes hands today. Which account recorded what the company owed this customer before today's work?

Hint: Before today's workshop, what did the company owe the participants who had already paid?

Explanation: Unearned Revenue recorded the workshop still owed after last month's payment. Completing today's instruction settles that obligation; Cash was already recorded last month.

- opt_ar: Accounts Receivable (tag: advance_confused_with_receivable; hint: The customer paid earlier. Was the company waiting for money or still owing the promised work?)
- opt_cash: Cash (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)
- opt_unearned_rev: Unearned Revenue (tag: ; hint: )

### account_category

What category of account is Unearned Revenue?

Hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?

Explanation: Unearned Revenue is a Liability: the company owes the customer work, or a refund. Despite its name, it is not a Revenue account.

- opt_asset: Asset (tag: ; hint: )
- opt_revenue: Revenue (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )

### direction

Does Unearned Revenue increase or decrease as the promised service is delivered?

Hint: Before today, the company owed the customer this work. After finishing it, how much does it still owe?

Explanation: The obligation has been fulfilled, so Unearned Revenue decreases.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is a decrease to a Liability (Unearned Revenue) recorded?

Hint: Which side is a liability's normal balance, and does a decrease go on that side or the opposite one?

Explanation: Liabilities have a normal credit balance, so a decrease is recorded as a Debit (left side). This debit does not mean money went out; no cash moved today.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

The prepaid work is completed today. Which account balances the entry?

Hint: Has completing the promised workshop earned anything even though no new payment arrived?

Explanation: Service Revenue records the workshop earned today. Cash would repeat last month's receipt, and Accounts Receivable would imply the already-paid participants still owe money.

- opt_ar: Accounts Receivable (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)
- opt_service_rev: Service Revenue (tag: ; hint: )
- opt_cash: Cash (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: Today's workshop fulfills the whole promise paid for last month. What obligation ends, and what is now earned?

Explanation: Debit Unearned Revenue $50 and Credit Service Revenue $50. The workshop obligation is settled and earning is recorded; there is no second cash receipt.

- opt_dr_unearned_cr_rev: Debit Unearned Revenue $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)
- opt_dr_rev_cr_unearned: Debit Service Revenue $50 / Credit Unearned Revenue $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Can the company earn the prepaid workshop fee today without receiving another payment?

Explanation: Liabilities decrease by $50 (Unearned Revenue) and equity increases by $50 through Service Revenue. Assets are unchanged because the cash arrived last month.

- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity unchanged. (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)
- opt_liab_down_eq_up: Liabilities decrease by $50 (-Unearned Revenue); Equity increases by $50 (+Service Revenue); Total Assets unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50; Equity increases by $50. (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)

## service_on_credit_repairs v1

A repair company finishes a customer's repair today and releases the repaired item. It invoices the customer for $50, payable in 30 days. The customer has made no payment or deposit, no work remains, and the company had not previously recorded this repair.

### identify_account

The customer was invoiced for services completed today and has not paid yet. Which account records this claim?

Hint: Did money arrive when the repaired item was released, or is the invoice still unpaid?

Explanation: Accounts Receivable records the right to collect for the completed repair. Cash would record a payment that has not happened; Unearned Revenue would describe unfinished work, but the repair is complete.

- opt_cash: Cash (tag: cash_recorded_when_uncollected; hint: Did the customer pay today, or does the company still have a right to collect later?)
- opt_ar: Accounts Receivable (tag: ; hint: )
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: Did money arrive when the repaired item was released, or is the invoice still unpaid?)

### account_category

What category of account is Accounts Receivable?

Hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?

Explanation: Accounts Receivable is an Asset: the right to collect cash later. It is not Revenue. Revenue records the earning; the receivable records the amount still to be collected.

- opt_equity: Equity (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_revenue: Revenue (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )

### direction

Does Accounts Receivable increase or decrease when the company bills a new customer?

Hint: For the work completed today, does the invoice create a new amount to collect or settle an existing one?

Explanation: The customer now owes the company money, so Accounts Receivable increases.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Accounts Receivable) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

Hint: Is earning based on finishing the repair or on waiting 30 days for payment?

Explanation: Service Revenue records the completed repair today. An unpaid invoice delays cash collection, not earning; Cash is not the counter-account because no money arrived.

- opt_service_rev: Service Revenue (tag: ; hint: )
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_cash: Cash (tag: cash_recorded_when_uncollected; hint: Did any payment arrive today, or does the customer still owe money for completed work?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: Is earning based on finishing the repair or on waiting 30 days for payment?)

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: The repair is complete and the invoice is unpaid. What right does the company hold, and what has it earned?

Explanation: Debit Accounts Receivable $50 and Credit Service Revenue $50. The customer owes the billed amount; no cash or customer advance is recorded.

- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_rev_cr_ar: Debit Service Revenue $50 / Credit Accounts Receivable $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: cash_recorded_when_uncollected; hint: Did any payment arrive today, or does the customer still owe money for completed work?)

### equation_effect

How does this service on credit affect the accounting equation?

Hint: Does an unpaid invoice for completed work represent a resource, even before money arrives?

Explanation: Assets increase by $50 (Accounts Receivable) and equity increases by $50 through Service Revenue. Cash and liabilities are unchanged.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does an unpaid invoice for completed work represent a resource, even before money arrives?)
- opt_assets_up_eq_up: Assets increase by $50 (+Accounts Receivable); Equity increases by $50 (+Service Revenue); Liabilities unchanged. (tag: ; hint: )
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50; Equity unchanged. (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)

## collect_receivable_repairs v1

A repair company receives $50 cash today to settle a customer's unpaid repair invoice in full. The repair was finished last month, when the company recorded both the earning and the customer's unpaid amount. No repair or other work is performed for this payment today.

### identify_account

The customer pays cash to settle their invoice. Which account reflects the cash received today?

Hint: What did the customer hand over today to settle the old repair invoice?

Explanation: Cash records the payment received today. Service Revenue would count the repair a second time: the earning was already recorded last month.

- opt_ap: Accounts Payable (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)
- opt_cash: Cash (tag: ; hint: )
- opt_service_rev: Service Revenue (tag: duplicate_revenue_on_collection; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- opt_asset: Asset (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_revenue: Revenue (tag: ; hint: )

### direction

Does Cash increase or decrease upon collecting this customer payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

The customer is paying an invoice where revenue was already recognized last month. What account must be credited?

Hint: What previously recorded unpaid amount disappears when the customer settles this invoice?

Explanation: Accounts Receivable decreases because the customer's unpaid amount is settled. Service Revenue was recorded last month; crediting it again would count the same repair twice.

- opt_ar: Accounts Receivable (tag: ; hint: )
- opt_service_rev: Service Revenue (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: What previously recorded unpaid amount disappears when the customer settles this invoice?)
- opt_unearned_rev: Unearned Revenue (tag: wrong_account; hint: What previously recorded unpaid amount disappears when the customer settles this invoice?)

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: Last month's repair and earning were already recorded. What resource arrived today, and what unpaid amount was settled?

Explanation: Debit Cash $50 and Credit Accounts Receivable $50. The payment replaces the right to collect with cash; it does not create new earning.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_dr_ar_cr_cash: Debit Accounts Receivable $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_ar: Debit Cash $50 / Credit Accounts Receivable $50 (tag: ; hint: )

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Did today's settlement create an additional resource, or turn the existing claim into cash?

Explanation: Cash increases by $50 and Accounts Receivable decreases by $50. Total assets, liabilities, and equity are unchanged; the repair's earning was recorded last month.

- opt_asset_swap: Asset exchange: Cash increases (+$50) and Accounts Receivable decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue). (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: wrong_account; hint: Did today's settlement create an additional resource, or turn the existing claim into cash?)

## borrow_cash_future_payroll v1

A company receives $50 cash today from a bank and signs a note requiring repayment in two years. It plans to use the money to pay employees next month. No employee payment or other spending occurs today, and no customer work or share issuance is involved. Ignore interest.

### identify_account

What did the company receive today, and which account records it?

Hint: What came into the business today, before any employee payment was made?

Explanation: Cash records the bank's payment into the company today. Planning to use the money for payroll later does not mean cash has been spent or a payroll cost has been recorded today.

- opt_cash: Cash (tag: ; hint: )
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)
- opt_service_rev: Service Revenue (tag: revenue_recorded_on_borrowing; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)

### account_category

What category of account is Cash?

Hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

Explanation: Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.

- opt_liability: Liability (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )

### direction

Does Cash increase or decrease upon borrowing money?

Hint: Compare the money held before and after today's receipt. Is there more or less?

Explanation: Cash increases because money arrived today. A future repayment obligation does not cancel today's receipt.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

Which account balances the receipt from the lender?

Hint: Does the intended future payroll use change what the bank expects back?

Explanation: Notes Payable records the bank loan. Future spending plans do not turn borrowing into Service Revenue or an ownership contribution.

- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_notes_payable: Notes Payable (Liability) (tag: ; hint: )
- opt_common_stock: Common Stock (Equity) (tag: wrong_account; hint: Does the intended future payroll use change what the bank expects back?)

### balanced_entry

What is the complete balanced journal entry for borrowing cash on a note?

Hint: Money was borrowed today and nothing was spent. What resource arrived, and what must the company repay?

Explanation: Debit Cash $50 and Credit Notes Payable $50. Borrowing creates cash and debt. Planned future payroll is a separate event, and no expense is recorded in this entry.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: ; hint: )
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)

### equation_effect

How does borrowing cash affect the accounting equation?

Hint: Does planning to spend the loan later reduce today's cash or create income today?

Explanation: Assets increase by $50 (Cash) and liabilities increase by $50 (Notes Payable). Equity is unchanged; neither earning nor spending occurred today.

- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue); Liabilities unchanged. (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_no_net_change: No net change in total assets. (tag: equation_effect_missed; hint: Does planning to spend the loan later reduce today's cash or create income today?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Notes Payable); Equity is unchanged. (tag: ; hint: )


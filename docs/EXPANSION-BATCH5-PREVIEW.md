# Expansion batch 5: complete teaching preview

Version-1 snapshots at $50, seed 100. Every allowed amount and scaffold is regression-tested. Shared stages and overrides appear with every option, tag, and stored mistake hint. Review: [EXPANSION-BATCH5-REVIEW.md](EXPANSION-BATCH5-REVIEW.md).

## cash_service_separate_job v1

A consulting company completes a new, separate job for a returning client today and receives $50 cash for that new job today. This new job had no earlier payment, invoice, or accounting entry, and no work remains for it. The client also owes an older invoice whose work and earning were recorded last month, but none of today's payment applies to that older invoice. Record only the new job and its payment.

### identify_account

Which account reflects the immediate payment received today?

Hint: What resource arrived today for the separate job that was just completed?

Explanation: Cash records today's payment for the new job. The older invoice remains unpaid. An existing customer balance does not make every later receipt a collection of that old balance.

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

Hint: Did today's payment settle the old invoice or pay for newly completed work?

Explanation: Service Revenue records the new job earned today. Accounts Receivable would reduce the older unpaid invoice, but the payment is explicitly for the separate new job. Unearned Revenue would imply that this new work remains to be done.

- opt_ap: Accounts Payable (tag: wrong_account; hint: Did today's payment settle the old invoice or pay for newly completed work?)
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: Did today's payment settle the old invoice or pay for newly completed work?)
- opt_service_rev: Service Revenue (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: Only the new job was completed and paid for today. What resource arrived and what new earning occurred?

Explanation: Debit Cash $50 and Credit Service Revenue $50. Do not credit Accounts Receivable: the older invoice is not being paid. No liability remains for the new completed job.

- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: wrong_account; hint: Only the new job was completed and paid for today. What resource arrived and what new earning occurred?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_rev_cr_cash: Debit Service Revenue $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_unearned: Debit Cash $50 / Credit Unearned Revenue $50 (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Did the new job create new earning or merely convert the old invoice into cash?

Explanation: Assets increase by $50 (Cash) and equity increases by $50 through Service Revenue. Liabilities and the older receivable are unchanged. The receipt belongs to newly earned work, not an old claim.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Did the new job create new earning or merely convert the old invoice into cash?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity is unchanged. (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Service Revenue); Liabilities are unchanged. (tag: ; hint: )

## service_on_credit_new_client v1

A consulting company completes a job for a new client today and invoices $50, due next month. A different client paid for their own completed job last month, but that earlier payment does not cover this new client's job. No cash or deposit has been received for the new job, no work remains for it, and it has not previously been recorded. Record only today's new completed job.

### identify_account

The customer was invoiced for services completed today and has not paid yet. Which account records this claim?

Hint: Does the other client's earlier payment cover this new invoice or does the new client still owe for today's separate job?

Explanation: Accounts Receivable records the unpaid claim for the new job. The earlier payment settled a different client's job. Cash would record money not received for this job, and Unearned Revenue would describe work still owed rather than completed.

- opt_cash: Cash (tag: cash_recorded_when_uncollected; hint: Did the customer pay today, or does the company still have a right to collect later?)
- opt_ar: Accounts Receivable (tag: ; hint: )
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: Does the other client's earlier payment cover this new invoice or does the new client still owe for today's separate job?)

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

Hint: Has the new job been earned even though its payment will arrive next month?

Explanation: Service Revenue records the new job completed today. The other client's payment is unrelated; it neither defers this earning nor supplies cash for the new invoice.

- opt_service_rev: Service Revenue (tag: ; hint: )
- opt_unearned_rev: Unearned Revenue (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)
- opt_cash: Cash (tag: cash_recorded_when_uncollected; hint: Did any payment arrive today, or does the customer still owe money for completed work?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: Has the new job been earned even though its payment will arrive next month?)

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: The separate new job is finished and remains unpaid. What claim and earning arose today?

Explanation: Debit Accounts Receivable $50 and Credit Service Revenue $50. No cash is recorded for the new job. An earlier payment for different work is not an advance available to release here.

- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_rev_cr_ar: Debit Service Revenue $50 / Credit Accounts Receivable $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: cash_recorded_when_uncollected; hint: Did any payment arrive today, or does the customer still owe money for completed work?)

### equation_effect

How does this service on credit affect the accounting equation?

Hint: Does another client's earlier payment remove the value of today's new unpaid claim?

Explanation: Assets increase by $50 (Accounts Receivable) and equity increases by $50 through Service Revenue. Cash and liabilities are unchanged. The new job is earned today despite collection being due next month.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does another client's earlier payment remove the value of today's new unpaid claim?)
- opt_assets_up_eq_up: Assets increase by $50 (+Accounts Receivable); Equity increases by $50 (+Service Revenue); Liabilities unchanged. (tag: ; hint: )
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50; Equity unchanged. (tag: revenue_deferred_when_earned; hint: After completing today's work, does the company still owe that work to the customer?)

## collect_receivable_installment v1

A consulting company receives $50 cash today as a partial payment of a larger invoice. It completed all the invoiced work last month and recorded both the full earning and the unpaid amount then. The stated amount is only today's installment, not the full invoice total. The rest remains owed, with no discount or write-off. No new work is performed today.

### identify_account

The customer pays cash to settle their invoice. Which account reflects the cash received today?

Hint: What resource arrived today even though part of the old invoice remains unpaid?

Explanation: Cash records the installment received today. The work was already earned and recorded last month, so receiving only part of the payment does not create new revenue.

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

Hint: What previously recorded claim becomes smaller by the installment just received?

Explanation: Accounts Receivable decreases only by the cash received. The unpaid remainder stays a claim against the client. Service Revenue was already recorded in full and must not be counted again.

- opt_ar: Accounts Receivable (tag: ; hint: )
- opt_service_rev: Service Revenue (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_ap: Accounts Payable (tag: wrong_account; hint: What previously recorded claim becomes smaller by the installment just received?)
- opt_unearned_rev: Unearned Revenue (tag: wrong_account; hint: What previously recorded claim becomes smaller by the installment just received?)

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: Record only today's installment on the already-recorded invoice. What resource arrived and what part of the claim was settled?

Explanation: Debit Cash $50 and Credit Accounts Receivable $50. Use the installment amount rather than the larger invoice total. Revenue remains unchanged, and the unpaid balance stays receivable.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_dr_ar_cr_cash: Debit Accounts Receivable $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_ar: Debit Cash $50 / Credit Accounts Receivable $50 (tag: ; hint: )

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Does collecting part of the old claim add to total resources or replace that part with cash?

Explanation: Cash increases by $50 and Accounts Receivable decreases by $50. Total assets, liabilities, and equity are unchanged. The remaining claim is preserved; no discount, write-off, or new earning is recorded.

- opt_asset_swap: Asset exchange: Cash increases (+$50) and Accounts Receivable decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue). (tag: duplicate_revenue_on_collection; hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?)
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: wrong_account; hint: Does collecting part of the old claim add to total resources or replace that part with cash?)

## customer_advance_additional_receipt v1

A training company receives $50 cash today as the final payment for a course beginning next month. The customer made a smaller first payment last month, and the company recorded that receipt and its obligation to provide instruction then. None of the instruction has been provided. After today the course is paid for in full. The stated amount is only the additional cash received today, not the total course price. No receivable was recorded for the final payment and today's receipt has not previously been recorded.

### identify_account

What did the company receive today, and which account records it?

Hint: Did additional money arrive today or is this simply a repeat of last month's receipt?

Explanation: Cash records the additional payment actually received today. Last month's receipt remains recorded and must not be recorded a second time. The stated amount is only this new receipt.

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

Hint: Has any instruction been delivered for this additional payment, or is the course still entirely ahead?

Explanation: Unearned Revenue increases for the additional advance. An existing obligation does not prevent another advance from increasing it. Service Revenue would count instruction not delivered; no recorded receivable is being settled.

- opt_ar: Accounts Receivable (Asset) (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)
- opt_unearned_rev: Unearned Revenue (Liability) (tag: ; hint: )
- opt_common_stock: Common Stock (Equity) (tag: wrong_account; hint: Has any instruction been delivered for this additional payment, or is the course still entirely ahead?)
- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: Only the final payment arrives today and the instruction is still ahead. What resource and additional obligation arise?

Explanation: Debit Cash $50 and Credit Unearned Revenue $50. Record only the new receipt, not the full course price or last month's payment. The earlier liability remains in place and increases by today's advance.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)
- opt_dr_cash_cr_unearned: Debit Cash $50 / Credit Unearned Revenue $50 (tag: ; hint: )
- opt_dr_unearned_cr_cash: Debit Unearned Revenue $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_ar_cr_rev: Debit Accounts Receivable $50 / Credit Service Revenue $50 (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does a previously recorded advance make today's additional payment earned before any instruction occurs?

Explanation: Assets increase by $50 (Cash) and liabilities increase by $50 (Unearned Revenue). Equity is unchanged. The earlier advance stays recorded; today adds only the final payment to cash and the prepaid obligation.

- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does a previously recorded advance make today's additional payment earned before any instruction occurs?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity is unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue); Liabilities are unchanged. (tag: revenue_recognized_prematurely; hint: Has the company performed the promised work yet, or does it still owe the customer that work?)

## earn_advance_completed_session v1

A training company completes one separately priced session today from a course paid for in full last month. Last month it recorded the payment and its obligation to provide all the sessions. The completed session has a stated price of $50; that amount was included in the recorded advance and has not previously been earned or recognized. Other prepaid sessions are still owed. No cash changes hands today. Record only this completed session.

### identify_account

No cash changes hands today. Which account recorded what the company owed this customer before today's work?

Hint: Which earlier balance recorded the obligation for the session that was paid for before today?

Explanation: Unearned Revenue held the prepaid obligation. Completing this session releases only its stated portion. Other prepaid sessions remain obligations; Cash was recorded when the full course payment arrived last month.

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

Hint: Does finishing one prepaid session reduce the work still owed even though other sessions remain?

Explanation: Unearned Revenue decreases by the completed session's price. The remaining obligation does not need to reach zero before a completed portion can be earned.

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

Hint: Has the completed session been earned even though later prepaid sessions still remain?

Explanation: Service Revenue records the session delivered today. Cash would repeat last month's receipt, and Accounts Receivable would imply this already-paid session still needs collection. Remaining sessions do not postpone earning from this completed session.

- opt_ar: Accounts Receivable (tag: advance_confused_with_receivable; hint: The customer has already paid. Is there still money to collect for this work?)
- opt_service_rev: Service Revenue (tag: ; hint: )
- opt_cash: Cash (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: Only this completed session is earned today. What portion of the prepaid obligation ends and what earning replaces it?

Explanation: Debit Unearned Revenue $50 and Credit Service Revenue $50. Use the completed session's price rather than the full course payment. The other sessions remain unearned and no second cash receipt is recorded.

- opt_dr_unearned_cr_rev: Debit Unearned Revenue $50 / Credit Service Revenue $50 (tag: ; hint: )
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)
- opt_dr_rev_cr_unearned: Debit Service Revenue $50 / Credit Unearned Revenue $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Can one completed prepaid session reduce an obligation without changing cash or settling every session?

Explanation: Liabilities decrease by $50 (Unearned Revenue) and equity increases by $50 through Service Revenue. Assets are unchanged. Only the completed session is released; the obligation for the remaining prepaid sessions stays recorded.

- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Unearned Revenue); Equity unchanged. (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)
- opt_liab_down_eq_up: Liabilities decrease by $50 (-Unearned Revenue); Equity increases by $50 (+Service Revenue); Total Assets unchanged. (tag: ; hint: )
- opt_assets_up_eq_up: Assets increase by $50; Equity increases by $50. (tag: cash_recorded_on_earning_advance; hint: Did money arrive today, or was the payment recorded before today's work was completed?)

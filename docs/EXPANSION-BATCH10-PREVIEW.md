# Expansion batch 10 rendered review

Draft preview at $100, seed 101. All seven stages include custom/inherited teaching and stored option hints.

## cash_service_third_party

A tutoring company completes a lesson today and receives $100 cash today from the customer's parent on the customer's behalf. This is full payment for the completed lesson, not a gift, investment, or loan. No earlier payment or accounting entry exists and no instruction remains owed. Record only today's completed lesson and receipt.

### identify_account

Which account reflects the immediate payment received today?

Hint: What resource arrived even though someone else paid for the customer?

Explanation: Cash records the actual receipt. A parent paying for the customer does not change the resource received.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived even though someone else paid for the customer?
- b: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable (id=opt_ar; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- a: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- b: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- c: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- d: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

### direction

Does Cash increase or decrease upon receiving this payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?

### debit_credit

How is an increase in an Asset (Cash) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

### counter_account

Which account records what the company did for the customer in exchange for the cash?

Hint: What did the company provide in exchange for this payment?

Explanation: Service Revenue records the completed lesson. The payer is settling the service price, not contributing capital or lending money.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What did the company provide in exchange for this payment?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: What did the company provide in exchange for this payment?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: What did the company provide in exchange for this payment?

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: Was the lesson completed and paid for, regardless of who handed over the money?

Explanation: Debit Cash $100 and Credit Service Revenue $100. No old receivable or advance existed; the third-party payer does not create a financing event.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
  Stored hint: Was the lesson completed and paid for, regardless of who handed over the money?
- b: Debit Service Revenue $100 / Credit Cash $100 (id=opt_dr_rev_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: Was the lesson completed and paid for, regardless of who handed over the money?

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does changing the payer remove either the receipt or the completed service?

Explanation: Assets increase by $100 through Cash and equity increases through Service Revenue. Liabilities are unchanged.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Service Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Does changing the payer remove either the receipt or the completed service?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does changing the payer remove either the receipt or the completed service?

## service_on_credit_invoice_later

A repair company completes a repair for a new customer today for an agreed price of $100. The customer owes the full enforceable price now but has not paid. The office will send the invoice next week; sending it is only paperwork and no further repair or customer acceptance is required. No deposit or earlier accounting entry exists. Record the completed unpaid repair today.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

Hint: Does delaying paperwork remove the enforceable unpaid claim for completed work?

Explanation: Accounts Receivable records the existing claim. The company need not wait for next week's paperwork to record completed, unpaid work.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Does delaying paperwork remove the enforceable unpaid claim for completed work?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Does delaying paperwork remove the enforceable unpaid claim for completed work?

### account_category

What category of account is Accounts Receivable?

Hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?

Explanation: Accounts Receivable is an Asset: the right to collect cash later. It is not Revenue. Revenue records the earning; the receivable records the amount still to be collected.

- a: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?
- b: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?
- c: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?
- d: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is a right to collect money later something the company owns, something it owes, or the owners' stake?

### direction

Does Accounts Receivable increase or decrease when completed work creates a new unpaid customer claim?

Hint: Has today's completed repair created a new unpaid claim?

Explanation: Accounts Receivable increases by the agreed price. No earlier claim was recorded and no cash has arrived.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Has today's completed repair created a new unpaid claim?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Has today's completed repair created a new unpaid claim?

### debit_credit

How is an increase to an Asset (Accounts Receivable) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

Hint: Is the repair finished even though the invoice will be sent later?

Explanation: Service Revenue records the repair earned today. The delayed invoice does not postpone completed work.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Is the repair finished even though the invoice will be sent later?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Is the repair finished even though the invoice will be sent later?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: What claim and earning exist now, before the office sends the bill?

Explanation: Debit Accounts Receivable $100 and Credit Service Revenue $100. No entry would omit completed earning; a Cash debit would invent payment.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?
- b: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=true; tag=)
  Stored hint: What claim and earning exist now, before the office sends the bill?
- c: Debit Service Revenue $100 / Credit Accounts Receivable $100 (id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does this service on credit affect the accounting equation?

Hint: Can a completed service create a resource before billing paperwork is sent?

Explanation: Assets increase by $100 through Accounts Receivable and equity increases through Service Revenue. Cash and liabilities are unchanged.

- a: Assets increase by $100 (+Accounts Receivable); Equity increases by $100 (+Service Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Can a completed service create a resource before billing paperwork is sent?
- b: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Can a completed service create a resource before billing paperwork is sent?

## collect_receivable_third_party

A repair company receives $100 cash today from a customer's parent, paying on the customer's behalf to settle the full recorded invoice. The repair was completed last month and its full earning and unpaid customer balance were recorded then. This receipt equals that invoice balance; no gift, loan, investment, fee, discount, or new service is involved. Record only today's collection.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

Hint: What resource arrives when another person pays the customer's old invoice?

Explanation: Cash records the receipt. Who pays does not move last month's earning into today.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrives when another person pays the customer's old invoice?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- a: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- b: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- c: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- d: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

### direction

Does Cash increase or decrease upon collecting this customer payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

Hint: Does the payment settle the recorded customer claim or buy a new service?

Explanation: Accounts Receivable decreases because the parent settles the existing invoice on the customer's behalf. Revenue was already recorded.

- a: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=wrong_account)
  Stored hint: Does the payment settle the recorded customer claim or buy a new service?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Does the payment settle the recorded customer claim or buy a new service?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Does the payment settle the recorded customer claim or buy a new service?

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: What old claim ends when its full amount arrives from the customer's parent?

Explanation: Debit Cash $100 and Credit Accounts Receivable $100. The payer's identity creates neither new revenue nor a company loan or investment.

- a: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: What old claim ends when its full amount arrives from the customer's parent?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Has another service been earned, or has the old claim become cash?

Explanation: Cash increases and Accounts Receivable decreases by $100. Total assets, liabilities, and equity are unchanged.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue). (id=opt_assets_up_eq_up; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=wrong_account)
  Stored hint: Has another service been earned, or has the old claim become cash?
- c: Asset exchange: Cash increases (+$100) and Accounts Receivable decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Has another service been earned, or has the old claim become cash?

## customer_advance_full_price

A tutoring company receives $100 cash in advance today, the full agreed price of a lesson scheduled for next week. The receipt says paid in full, but no instruction or separate booking service has been delivered and the company still owes the entire lesson. No earlier payment or related accounting entry exists. Record only today's advance receipt.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource arrives before any instruction is supplied?

Explanation: Cash records the full receipt today. The payment receipt establishes payment, not completed instruction.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrives before any instruction is supplied?
- d: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?

### account_category

What category of account is Cash?

Hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

Explanation: Cash is an Asset: a resource the company owns and can use. It is not Revenue. Revenue records earning, and cash can arrive without any earning (a loan, a customer advance).

- a: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- b: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- c: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?
- d: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is cash something the company owns and can use, something it owes, or the owners' stake?

### direction

Does Cash increase or decrease upon receiving this customer payment?

Hint: Compare the company's cash before and after today's payment. Is it higher or lower?

Explanation: The company holds more money after the payment, so Cash increases.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Compare the company's cash before and after today's payment. Is it higher or lower?

### debit_credit

How is an increase to an Asset (Cash) recorded in double-entry bookkeeping?

Hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is recorded as a Debit (left side). Debit only means left; whether a debit increases an account depends on the account's category.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on the same side as its normal balance. Which side is an asset's normal balance?

### counter_account

What counter-account balances this receipt for services still to be performed?

Hint: Does paid in full mean the company has finished the promised lesson?

Explanation: Unearned Revenue records the entire lesson still owed. Full payment removes the customer's unpaid price, not the company's performance obligation.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Does paid in full mean the company has finished the promised lesson?
- c: Accounts Receivable (Asset) (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: Does paid in full mean the company has finished the promised lesson?

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: What does the company receive and still owe when the customer prepays the whole price?

Explanation: Debit Cash $100 and Credit Unearned Revenue $100. Full prepayment is not completed earning; no customer claim remains unpaid.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- b: Debit Unearned Revenue $100 / Credit Cash $100 (id=opt_dr_unearned_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=true; tag=)
  Stored hint: What does the company receive and still owe when the customer prepays the whole price?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does paying the whole price remove work that has not begun?

Explanation: Assets and liabilities increase by $100. Equity is unchanged because the lesson remains wholly unperformed.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: Does paying the whole price remove work that has not begun?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does paying the whole price remove work that has not begun?

## earn_advance_paperwork_later

A tutoring company completes the entire prepaid lesson today. The customer paid $100 cash last week and the company already recorded that receipt and its obligation to provide this lesson. The customer's attendance certificate will be printed next week; it is only paperwork, not a separate promised service or condition of completion. No instruction remains, no earning for this lesson was previously recorded, and no cash changes hands today. Record only today's completion.

### identify_account

No new cash changes hands when this prepaid work is completed. Which account recorded the obligation before completion?

Hint: Which recorded obligation has the completed instruction now satisfied?

Explanation: Unearned Revenue records the earlier obligation. Printing routine paperwork later does not leave any promised instruction owed.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Unearned Revenue (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Which recorded obligation has the completed instruction now satisfied?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer paid earlier. Was the company waiting for money or still owing the promised work?

### account_category

What category of account is Unearned Revenue?

Hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?

Explanation: Unearned Revenue is a Liability: the company owes the customer work, or a refund. Despite its name, it is not a Revenue account.

- a: Liability (id=opt_liability; correct=true; tag=)
  Stored hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?
- b: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?
- c: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?
- d: Asset (id=opt_asset; correct=false; tag=)
  Stored hint: Is an obligation to do work for a customer something the company owns, something it owes, or the owners' stake?

### direction

Does Unearned Revenue increase or decrease as the promised service is delivered?

Hint: Does any promised service remain after this entire lesson is complete?

Explanation: Unearned Revenue decreases by the full lesson price. This lesson's obligation reaches zero; routine paperwork does not preserve it.

- a: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: Does any promised service remain after this entire lesson is complete?
- b: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: Does any promised service remain after this entire lesson is complete?

### debit_credit

How is a decrease to a Liability (Unearned Revenue) recorded?

Hint: Which side is a liability's normal balance, and does a decrease go on that side or the opposite one?

Explanation: Liabilities have a normal credit balance, so a decrease is recorded as a Debit (left side). This debit does not mean money went out; no new cash moves when the prepaid work is completed.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: Which side is a liability's normal balance, and does a decrease go on that side or the opposite one?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: Which side is a liability's normal balance, and does a decrease go on that side or the opposite one?

### counter_account

The prepaid work is completed today. Which account balances the entry?

Hint: When was the promised lesson delivered: today or when paperwork is printed?

Explanation: Service Revenue records today's completed instruction. Cash was already recorded last week; a new receivable would charge again for prepaid work.

- a: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: When was the promised lesson delivered: today or when paperwork is printed?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: What recorded obligation ends when prepaid instruction is finished?

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Do not repeat last week's Cash receipt or delay earning for routine paperwork.

- a: Debit Service Revenue $100 / Credit Unearned Revenue $100 (id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 (id=opt_dr_unearned_cr_rev; correct=true; tag=)
  Stored hint: What recorded obligation ends when prepaid instruction is finished?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Does printing a certificate later require another cash receipt today?

Explanation: Liabilities decrease and equity increases by $100; assets are unchanged. The earlier receipt and today's earning are separate events.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Assets increase by $100; Equity increases by $100. (id=opt_assets_up_eq_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Liabilities decrease by $100 (-Unearned Revenue); Equity increases by $100 (+Service Revenue); Total Assets unchanged. (id=opt_liab_down_eq_up; correct=true; tag=)
  Stored hint: Does printing a certificate later require another cash receipt today?

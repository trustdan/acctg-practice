# Expansion batch 8 rendered review

Draft preview at $100, seed 101. All seven stages include inherited teaching and stored option hints.

## cash_service_unpaid_booking

A tutoring company completes a lesson today and receives $100 cash for it today. The customer reserved the appointment last week, but that booking involved no payment, invoice, work, or accounting entry. Today's lesson is fully completed and no work remains for this payment. Record only the completed lesson and today's receipt.

### identify_account

Which account reflects the immediate payment received today?

Hint: Did last week's unpaid appointment booking supply any money before today?

Explanation: Cash records today's receipt. The earlier booking created neither an advance receipt nor a recorded customer balance.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Did last week's unpaid appointment booking supply any money before today?
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

Hint: Was this lesson earned last week when booked or today when delivered?

Explanation: Service Revenue records the lesson completed today. An appointment reservation alone did not earn the lesson price. Accounts Receivable would imply a previously recorded unpaid claim; Unearned Revenue would imply work still owed.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Was this lesson earned last week when booked or today when delivered?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Was this lesson earned last week when booked or today when delivered?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Was this lesson earned last week when booked or today when delivered?

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: The appointment was only booked earlier; the lesson is now complete and paid for. What receipt and earning occur today?

Explanation: Debit Cash $100 and Credit Service Revenue $100. Do not release an advance or collect an old receivable: neither existed. The earlier booking does not postpone or duplicate today's earning.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
  Stored hint: The appointment was only booked earlier; the lesson is now complete and paid for. What receipt and earning occur today?
- b: Debit Service Revenue $100 / Credit Cash $100 (id=opt_dr_rev_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: The appointment was only booked earlier; the lesson is now complete and paid for. What receipt and earning occur today?

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Did the unpaid booking create resources or earning that today's receipt merely replaces?

Explanation: Assets increase by $100 (Cash) and equity increases by $100 through Service Revenue. Liabilities are unchanged. Both receipt and earning occur today; the unpaid booking created neither balance.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Service Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Did the unpaid booking create resources or earning that today's receipt merely replaces?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Did the unpaid booking create resources or earning that today's receipt merely replaces?

## service_on_credit_due_today

A repair company completes a repair for a new customer today and invoices the customer for $100, due immediately today. The customer has not paid, and the company retains an enforceable claim for the full price. No deposit was received, no work remains, and this repair has not previously been recorded. Record the completed repair while its invoice is still unpaid.

### identify_account

The customer was invoiced for services completed today and has not paid yet. Which account records this claim?

Hint: Does an invoice being due today mean money has actually arrived?

Explanation: Accounts Receivable records the unpaid claim for the completed repair. An immediate due date does not create Cash before the customer pays. No deposit or unearned obligation exists.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Does an invoice being due today mean money has actually arrived?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Does an invoice being due today mean money has actually arrived?

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

Does Accounts Receivable increase or decrease when the company bills a new customer?

Hint: While payment remains outstanding, has the company gained a claim for this new completed repair?

Explanation: Accounts Receivable increases by the newly invoiced price. Being due immediately does not reduce an unpaid claim or turn it into cash.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: While payment remains outstanding, has the company gained a claim for this new completed repair?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: While payment remains outstanding, has the company gained a claim for this new completed repair?

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

Hint: Has the repair been delivered even though the immediately due payment is still missing?

Explanation: Service Revenue records the repair completed today. Earning follows delivery here, not whether the customer met the invoice's due date. No further repair work is owed.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Has the repair been delivered even though the immediately due payment is still missing?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Has the repair been delivered even though the immediately due payment is still missing?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: The repair is finished, but the due-today invoice remains unpaid. What claim and earning arose?

Explanation: Debit Accounts Receivable $100 and Credit Service Revenue $100. A Cash debit would invent a payment. Waiting for payment would omit an earned claim; a customer-advance credit would imply unfinished work.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?
- b: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=true; tag=)
  Stored hint: The repair is finished, but the due-today invoice remains unpaid. What claim and earning arose?
- c: Debit Service Revenue $100 / Credit Accounts Receivable $100 (id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does this service on credit affect the accounting equation?

Hint: Can an unpaid claim be a resource even when its payment deadline is today?

Explanation: Assets increase by $100 (Accounts Receivable) and equity increases by $100 through Service Revenue. Cash and liabilities are unchanged. Payment timing does not erase the earned claim.

- a: Assets increase by $100 (+Accounts Receivable); Equity increases by $100 (+Service Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Can an unpaid claim be a resource even when its payment deadline is today?
- b: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Can an unpaid claim be a resource even when its payment deadline is today?

## collect_receivable_overdue

A repair company receives $100 cash today to settle an overdue invoice in full. The repair was completed last month, and the company recorded the full earning and unpaid customer balance then. The amount received equals the recorded invoice balance. There is no interest, late fee, discount, or write-off. No new repair or other work is performed today. Record only collection of the overdue invoice.

### identify_account

The customer pays cash to settle their invoice. Which account reflects the cash received today?

Hint: What resource arrived today for work already completed and recorded last month?

Explanation: Cash records today's overdue payment. Receiving it late does not move the already recorded earning into today.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived today for work already completed and recorded last month?
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

The customer is paying an invoice where revenue was already recognized last month. What account must be credited?

Hint: Which previously recorded claim ends when the customer finally pays the full invoice?

Explanation: Accounts Receivable decreases to settle the old invoice. Service Revenue was already recorded last month; lateness alone creates no new revenue or advance obligation.

- a: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=wrong_account)
  Stored hint: Which previously recorded claim ends when the customer finally pays the full invoice?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Which previously recorded claim ends when the customer finally pays the full invoice?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Which previously recorded claim ends when the customer finally pays the full invoice?

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: The old invoice is paid in full with no fee or adjustment. What resource arrives and what existing claim ends?

Explanation: Debit Cash $100 and Credit Accounts Receivable $100. Revenue is not recorded again. No expense, interest, or write-off is part of this full collection.

- a: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: The old invoice is paid in full with no fee or adjustment. What resource arrives and what existing claim ends?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Does paying after the deadline create another earning or exchange the old claim for cash?

Explanation: Cash increases by $100 and Accounts Receivable decreases by $100. Total assets, liabilities, and equity are unchanged. The delay introduces no fee or loss in this scenario.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue). (id=opt_assets_up_eq_up; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=wrong_account)
  Stored hint: Does paying after the deadline create another earning or exchange the old claim for cash?
- c: Asset exchange: Cash increases (+$100) and Accounts Receivable decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does paying after the deadline create another earning or exchange the old claim for cash?

## customer_advance_same_day

At 9 a.m. today, a tutoring company receives $100 cash in advance for a lesson scheduled for 6 p.m. today. No instruction or separate booking service has been provided by 9 a.m.; the company still owes the entire lesson. This is the full lesson price and no payment, invoice, or related obligation was recorded earlier. Record only the morning receipt before the lesson begins.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource has arrived by 9 a.m., before any lesson is delivered?

Explanation: Cash records the morning receipt. The lesson occurring later on the same day does not mean it has already been delivered.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource has arrived by 9 a.m., before any lesson is delivered?
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

Hint: At the morning recording time, has the company supplied the lesson or does it still owe it?

Explanation: Unearned Revenue records the full lesson still owed at 9 a.m. Service Revenue would recognize instruction before delivery. A short wait until evening does not remove the obligation.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: At the morning recording time, has the company supplied the lesson or does it still owe it?
- c: Accounts Receivable (Asset) (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: At the morning recording time, has the company supplied the lesson or does it still owe it?

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: Record only the morning receipt, with the lesson still ahead. What resource and obligation arise then?

Explanation: Debit Cash $100 and Credit Unearned Revenue $100. The lesson's same-day schedule does not earn it at 9 a.m. Completion that evening will be a separate event with no second cash receipt.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- b: Debit Unearned Revenue $100 / Credit Cash $100 (id=opt_dr_unearned_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=true; tag=)
  Stored hint: Record only the morning receipt, with the lesson still ahead. What resource and obligation arise then?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does promising delivery later today remove the work still owed at 9 a.m.?

Explanation: Assets increase by $100 (Cash) and liabilities increase by $100 (Unearned Revenue). Equity is unchanged at the morning receipt. Same-day performance remains future performance until the lesson is delivered.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: Does promising delivery later today remove the work still owed at 9 a.m.?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does promising delivery later today remove the work still owed at 9 a.m.?

## earn_advance_same_day

At 7 p.m. today, a tutoring company finishes the entire lesson for which a customer paid $100 cash at 9 a.m. today. The company already recorded the morning receipt and its obligation to provide this lesson. The full lesson price was included in that recorded advance and has not previously been recognized as earned. No work remains and no cash changes hands at 7 p.m. Record only completion of the lesson this evening.

### identify_account

No new cash changes hands when this prepaid work is completed. Which account recorded the obligation before completion?

Hint: Which recorded morning balance held the obligation that this evening's lesson now satisfies?

Explanation: Unearned Revenue records the obligation created by the morning advance. The morning receipt is already recorded and must not be repeated this evening.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Unearned Revenue (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Which recorded morning balance held the obligation that this evening's lesson now satisfies?
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

Hint: Once the entire prepaid lesson is delivered, is any instruction still owed for that payment?

Explanation: Unearned Revenue decreases by the full lesson price and this lesson's obligation reaches zero. Completing the lesson releases a recorded liability even when payment and delivery share a calendar date.

- a: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: Once the entire prepaid lesson is delivered, is any instruction still owed for that payment?
- b: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: Once the entire prepaid lesson is delivered, is any instruction still owed for that payment?

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

Hint: What has been earned by finishing the lesson, without receiving another payment?

Explanation: Service Revenue records the lesson delivered this evening. Cash was recorded in the morning; Accounts Receivable would invent an unpaid claim for a lesson already paid in full.

- a: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: What has been earned by finishing the lesson, without receiving another payment?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: The morning cash is already recorded and the prepaid lesson is now complete. What obligation ends and what earning replaces it?

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Do not record Cash again. No entry would leave the completed lesson unearned; an unpaid-claim debit would charge again for work already paid for.

- a: Debit Service Revenue $100 / Credit Unearned Revenue $100 (id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 (id=opt_dr_unearned_cr_rev; correct=true; tag=)
  Stored hint: The morning cash is already recorded and the prepaid lesson is now complete. What obligation ends and what earning replaces it?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Does completing the recorded prepaid lesson move any additional cash this evening?

Explanation: Liabilities decrease by $100 (Unearned Revenue) and equity increases by $100 through Service Revenue. Assets are unchanged at 7 p.m. The earlier cash receipt and this later earning remain separate events.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Assets increase by $100; Equity increases by $100. (id=opt_assets_up_eq_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Liabilities decrease by $100 (-Unearned Revenue); Equity increases by $100 (+Service Revenue); Total Assets unchanged. (id=opt_liab_down_eq_up; correct=true; tag=)
  Stored hint: Does completing the recorded prepaid lesson move any additional cash this evening?

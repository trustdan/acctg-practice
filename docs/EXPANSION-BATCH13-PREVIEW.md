# Expansion batch 13 rendered review

Draft preview at $100, seed 101. All seven stages include custom/inherited teaching and stored option hints.

## cash_service_stockholder_customer

A corporation completes a repair today and receives $100 cash today from a customer who also owns shares in the corporation. The full payment is solely the agreed price of this completed repair, not a new investment or loan. No shares are issued, no prior invoice, advance, or entry exists for this repair, and no work remains. Record only today's completed service and receipt from the corporation's perspective.

### identify_account

Which account reflects the immediate payment received today?

Hint: What resource arrived from the stockholder acting as a customer today?

Explanation: Cash records today's receipt for the repair. Holding shares does not change the resource received.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived from the stockholder acting as a customer today?
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

Hint: Did this payment buy new ownership or pay for completed customer work?

Explanation: Service Revenue records the completed repair. The payer's existing share ownership creates no new Common Stock or loan in this transaction.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Did this payment buy new ownership or pay for completed customer work?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Did this payment buy new ownership or pay for completed customer work?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Did this payment buy new ownership or pay for completed customer work?

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: What receipt and earning occur when an existing stockholder pays for a completed repair?

Explanation: Debit Cash $100 and Credit Service Revenue $100. Do not record share capital, borrowing, an old collection, or a future-work obligation.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
  Stored hint: What receipt and earning occur when an existing stockholder pays for a completed repair?
- b: Debit Service Revenue $100 / Credit Cash $100 (id=opt_dr_rev_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: What receipt and earning occur when an existing stockholder pays for a completed repair?

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does the customer's ownership role remove today's completed earning?

Explanation: Assets increase through Cash and equity increases through Service Revenue by $100. Liabilities and existing Common Stock are unchanged.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Service Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Does the customer's ownership role remove today's completed earning?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does the customer's ownership role remove today's completed earning?

## service_on_credit_customer_withdrawal

A repair company completes a repair for a new customer today and invoices $100. Before visiting, the customer withdrew that amount from their own bank account and still holds it, but has given no money to the repair company. The company has an enforceable claim for the full completed repair; no deposit, prior entry, or remaining work exists. Record only the company's completed unpaid service, not the customer's personal withdrawal.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

Hint: Has the company received the customer's money or only gained a right to collect it?

Explanation: Accounts Receivable records the company's unpaid claim. Cash held by the customer is not cash held by the company.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Has the company received the customer's money or only gained a right to collect it?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Has the company received the customer's money or only gained a right to collect it?

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

Hint: Does today's unpaid repair create or settle the company's claim?

Explanation: Accounts Receivable increases by the repair price. A withdrawal in the customer's own account does not pay the company.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does today's unpaid repair create or settle the company's claim?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does today's unpaid repair create or settle the company's claim?

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

Hint: Was the company's repair earned even though the customer still holds the money?

Explanation: Service Revenue records the completed repair. The customer's personal cash position does not postpone this earning.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Was the company's repair earned even though the customer still holds the money?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Was the company's repair earned even though the customer still holds the money?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: Whose cash would a company Cash debit describe before payment is handed over?

Explanation: Debit Accounts Receivable $100 and Credit Service Revenue $100. A Cash debit would record the customer's money as the company's without a receipt.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?
- b: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=true; tag=)
  Stored hint: Whose cash would a company Cash debit describe before payment is handed over?
- c: Debit Service Revenue $100 / Credit Accounts Receivable $100 (id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does this service on credit affect the accounting equation?

Hint: What resource belongs to the company while the customer still holds the withdrawn cash?

Explanation: Assets increase through Accounts Receivable and equity increases through Service Revenue by $100. Company Cash and liabilities are unchanged.

- a: Assets increase by $100 (+Accounts Receivable); Equity increases by $100 (+Service Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: What resource belongs to the company while the customer still holds the withdrawn cash?
- b: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: What resource belongs to the company while the customer still holds the withdrawn cash?

## collect_receivable_future_booking

A consulting company receives $100 cash today solely to settle a recorded invoice for work completed last month. The company already recorded that earning and the unpaid customer amount then; the receipt equals the full invoice balance. The customer also books a separate appointment for next month, but makes no payment for it and no new work is performed today. No fee, discount, or write-off applies. Record only collection of the old invoice, not an advance for the unpaid future booking.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

Hint: What resource arrived for the old invoice, apart from the unpaid new booking?

Explanation: Cash records the receipt applied to last month's recorded invoice. No portion pays for the future appointment.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived for the old invoice, apart from the unpaid new booking?
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

Hint: Which recorded claim does the explicitly applied payment settle?

Explanation: Accounts Receivable decreases for the old invoice. Revenue was already recorded; the separate unpaid appointment creates no advance receipt.

- a: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=wrong_account)
  Stored hint: Which recorded claim does the explicitly applied payment settle?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Which recorded claim does the explicitly applied payment settle?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Which recorded claim does the explicitly applied payment settle?

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: What old balance ends when all of this payment applies to its invoice?

Explanation: Debit Cash $100 and Credit Accounts Receivable $100. Do not repeat revenue or credit Unearned Revenue for an appointment receiving none of this money.

- a: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: What old balance ends when all of this payment applies to its invoice?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Does merely booking future work change how the old invoice payment is applied?

Explanation: Cash increases and Accounts Receivable decreases by $100. Total assets, liabilities, and equity are unchanged. The unpaid booking is separate.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue). (id=opt_assets_up_eq_up; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=wrong_account)
  Stored hint: Does merely booking future work change how the old invoice payment is applied?
- c: Asset exchange: Cash increases (+$100) and Accounts Receivable decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does merely booking future work change how the old invoice payment is applied?

## customer_advance_internal_setup

A tutoring company receives $100 cash in advance today for an entire lesson next month. Staff have arranged chairs and tested their own equipment, but no instruction or other promised customer service has been delivered. Setup is internal preparation, not a separately sold or promised service; the company still owes the full lesson. This is the full price and no prior payment or related entry exists. Record only the receipt, with no separate setup cost or payment in this event.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource arrived before the promised instruction was delivered?

Explanation: Cash records the receipt. Internal preparation does not change the money received into completed instruction.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived before the promised instruction was delivered?
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

Hint: Has internal setup supplied the promised lesson or does the company still owe it?

Explanation: Unearned Revenue records the full lesson still owed. Internal readiness is not a separate promised service earning any part of this price.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Has internal setup supplied the promised lesson or does the company still owe it?
- c: Accounts Receivable (Asset) (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: Has internal setup supplied the promised lesson or does the company still owe it?

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: What money and obligation arise before any promised customer service is delivered?

Explanation: Debit Cash $100 and Credit Unearned Revenue $100. Do not recognize Service Revenue merely because staff prepared to teach; no promised performance is complete.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- b: Debit Unearned Revenue $100 / Credit Cash $100 (id=opt_dr_unearned_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=true; tag=)
  Stored hint: What money and obligation arise before any promised customer service is delivered?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: Does arranging the room satisfy the customer's right to the full lesson?

Explanation: Assets and liabilities increase by $100; equity is unchanged. Internal preparation does not release the lesson obligation.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: Does arranging the room satisfy the customer's right to the full lesson?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does arranging the room satisfy the customer's right to the full lesson?

## earn_advance_final_remaining

A training company completes the final separately priced session of a prepaid course today. The customer paid for the entire course earlier, and the company recorded the receipt and its obligation then. All earlier sessions and their earning were already recorded. Today's final session has a price of $100, equal to the entire remaining recorded advance immediately before completion, not the original full course price. No work remains for this course and no cash changes hands today. Record only earning the last remaining portion.

### identify_account

No new cash changes hands when this prepaid work is completed. Which account recorded the obligation before completion?

Hint: Which recorded obligation remains after earlier sessions were already earned?

Explanation: Unearned Revenue records the remaining final-session obligation. Earlier receipt and earning are not posted again.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Unearned Revenue (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Which recorded obligation remains after earlier sessions were already earned?
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

Hint: After the final session is delivered, does any recorded course obligation remain?

Explanation: Unearned Revenue decreases by the final session's price to zero. Previously earned sessions were already removed from that balance.

- a: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: After the final session is delivered, does any recorded course obligation remain?
- b: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: After the final session is delivered, does any recorded course obligation remain?

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

Hint: What is earned by today's final session without another receipt?

Explanation: Service Revenue records only today's completed session. Cash was recorded earlier, and earlier sessions were already recognized as earned.

- a: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: What is earned by today's final session without another receipt?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: How much of the recorded prepaid course remains to be earned today?

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Use the remaining final-session price, not the original full course price; do not repeat receipt or earlier earning.

- a: Debit Service Revenue $100 / Credit Unearned Revenue $100 (id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 (id=opt_dr_unearned_cr_rev; correct=true; tag=)
  Stored hint: How much of the recorded prepaid course remains to be earned today?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Does finishing the course bring more cash or release its last recorded obligation?

Explanation: Liabilities decrease and equity increases by $100; assets are unchanged. The course obligation reaches zero and previously recorded earning stays intact.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Assets increase by $100; Equity increases by $100. (id=opt_assets_up_eq_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Liabilities decrease by $100 (-Unearned Revenue); Equity increases by $100 (+Service Revenue); Total Assets unchanged. (id=opt_liab_down_eq_up; correct=true; tag=)
  Stored hint: Does finishing the course bring more cash or release its last recorded obligation?

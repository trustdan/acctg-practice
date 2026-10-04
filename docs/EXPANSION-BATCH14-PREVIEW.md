# Expansion batch 14 rendered review

Draft preview at $100, seed 101. All seven stages include custom/inherited teaching, options, tags and stored option hints.

## cash_service_undeposited_notes

A repair company completes a repair today and receives $100 in physical currency from the customer today. The notes are in the company's locked cash box; the company will deposit them in its bank account tomorrow. This is full payment for the completed repair. No earlier payment, invoice, or accounting entry exists, and no work remains. Record today's completed repair and receipt before the bank deposit.

### identify_account

Which account reflects the immediate payment received today?

Hint: Does money already held in the company's cash box belong to it before a bank deposit?

Explanation: Cash includes currency already held by the company. Depositing these same notes tomorrow does not postpone today's receipt.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Does money already held in the company's cash box belong to it before a bank deposit?
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

Hint: Does the company hold more currency after the customer hands over the notes?

Explanation: Cash increases by the currency received today, even while the notes remain in the cash box.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does the company hold more currency after the customer hands over the notes?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does the company hold more currency after the customer hands over the notes?

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

Hint: Was the promised repair finished when the company received the notes?

Explanation: Service Revenue records the completed repair. A later bank deposit does not make the completed service unearned.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Was the promised repair finished when the company received the notes?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Was the promised repair finished when the company received the notes?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Was the promised repair finished when the company received the notes?

### balanced_entry

What is the complete balanced journal entry for this cash service transaction?

Hint: Have both the company's money and its completed earning changed before the deposit?

Explanation: Debit Cash $100 and Credit Service Revenue $100. The company already holds the payment; no customer claim or future performance remains.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
  Stored hint: Have both the company's money and its completed earning changed before the deposit?
- b: Debit Service Revenue $100 / Credit Cash $100 (id=opt_dr_rev_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: Have both the company's money and its completed earning changed before the deposit?

### equation_effect

How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: What resource and earning exist today while the bank balance has not yet changed?

Explanation: Assets increase through Cash and equity increases through Service Revenue by $100. Physical currency is a company asset even before deposit; liabilities are unchanged.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Service Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: What resource and earning exist today while the bank balance has not yet changed?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: What resource and earning exist today while the bank balance has not yet changed?

## service_on_credit_cash_sale_label

A repair company completes a repair for a new customer today for an agreed price of $100. The office mistakenly marks its job ticket cash sale, but the customer has made no payment and the company has an enforceable claim for the entire price. The company invoices the customer today. No deposit, prior accounting entry, or remaining work exists. Record the completed unpaid repair using the actual facts, regardless of the ticket label.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

Hint: What does the company actually hold: customer money or an unpaid right to collect?

Explanation: Accounts Receivable records the enforceable unpaid claim. The mistaken cash-sale label cannot create money the company did not receive.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What does the company actually hold: customer money or an unpaid right to collect?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: What does the company actually hold: customer money or an unpaid right to collect?

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

Hint: Does completing unpaid work create or settle the company's claim?

Explanation: Accounts Receivable increases by the agreed repair price. No cash was received despite the job-ticket label.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does completing unpaid work create or settle the company's claim?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does completing unpaid work create or settle the company's claim?

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

Hint: Did the repair get completed even though the payment label is mistaken?

Explanation: Service Revenue records the completed repair. Earning depends on performance here, while the unpaid price remains a receivable.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Did the repair get completed even though the payment label is mistaken?
- b: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: Did the repair get completed even though the payment label is mistaken?
- c: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash (id=opt_cash; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?

### balanced_entry

What is the complete balanced journal entry for this service on credit?

Hint: Which entry follows the unpaid claim and completed work rather than the ticket label?

Explanation: Debit Accounts Receivable $100 and Credit Service Revenue $100. Debiting Cash would invent a receipt; deferring earning would ignore the completed repair.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_when_uncollected)
  Stored hint: Did any payment arrive today, or does the customer still owe money for completed work?
- b: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=true; tag=)
  Stored hint: Which entry follows the unpaid claim and completed work rather than the ticket label?
- c: Debit Service Revenue $100 / Credit Accounts Receivable $100 (id=opt_dr_rev_cr_ar; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does this service on credit affect the accounting equation?

Hint: What resource did the company gain without receiving any money?

Explanation: Assets increase through Accounts Receivable and equity increases through Service Revenue by $100. Cash and liabilities are unchanged regardless of the mistaken label.

- a: Assets increase by $100 (+Accounts Receivable); Equity increases by $100 (+Service Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: What resource did the company gain without receiving any money?
- b: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: What resource did the company gain without receiving any money?

## collect_receivable_before_due

A consulting company receives $100 cash today to settle a recorded invoice in full, two weeks before its due date. All the invoiced work was completed last month, when the company recorded the full earning and unpaid customer amount. The receipt equals the invoice balance; no discount, fee, interest, or write-off applies. No new work is performed today. Record only collection of this existing invoice, even though payment arrived before the deadline.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

Hint: What resource arrived when the customer settled the existing invoice early?

Explanation: Cash records the money received today. Early payment changes the collection date, not the already completed work.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived when the customer settled the existing invoice early?
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

Hint: Did this payment settle an existing claim or pay for work still owed?

Explanation: Accounts Receivable decreases as the recorded invoice is settled. An early due-date payment is not an advance for future work when the work was already completed and recorded.

- a: Unearned Revenue (id=opt_unearned_rev; correct=false; tag=wrong_account)
  Stored hint: Did this payment settle an existing claim or pay for work still owed?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: Did this payment settle an existing claim or pay for work still owed?
- d: Accounts Receivable (id=opt_ar; correct=true; tag=)
  Stored hint: Did this payment settle an existing claim or pay for work still owed?

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

Hint: Which recorded customer balance disappears when the early payment arrives?

Explanation: Debit Cash $100 and Credit Accounts Receivable $100. The old earning stays recorded once; neither new revenue nor a future-service obligation arises.

- a: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: Which recorded customer balance disappears when the early payment arrives?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does collecting an existing receivable affect the accounting equation?

Hint: Does replacing an unpaid claim with money add a new total asset or new earning?

Explanation: Cash rises and Accounts Receivable falls by $100, leaving total assets, liabilities and equity unchanged. Paying before the due date does not create a second earning.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue). (id=opt_assets_up_eq_up; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=wrong_account)
  Stored hint: Does replacing an unpaid claim with money add a new total asset or new earning?
- c: Asset exchange: Cash increases (+$100) and Accounts Receivable decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does replacing an unpaid claim with money add a new total asset or new earning?

## customer_advance_customer_books

A tutoring company receives $100 cash today from a business customer in advance for a lesson next month. The customer has entered the payment in its own books, but that does not record any entry for the tutoring company. No instruction or other promised service has been provided, and the tutoring company still owes the full lesson. This is the full price, with no earlier receipt or related entry in the tutoring company's books. Record only the tutoring company's receipt and remaining promise.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource entered the tutoring company, regardless of the customer's own bookkeeping?

Explanation: Cash records the tutoring company's receipt. The customer's separate books do not record the seller's entry.

- a: Accounts Payable (id=opt_ap; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource entered the tutoring company, regardless of the customer's own bookkeeping?
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

Hint: What does the tutoring company still owe after receiving payment for next month?

Explanation: Unearned Revenue records the tutoring company's obligation to provide the lesson. The customer's own entry does not deliver any instruction or earn the seller's revenue.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: What does the tutoring company still owe after receiving payment for next month?
- c: Accounts Receivable (Asset) (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: What does the tutoring company still owe after receiving payment for next month?

### balanced_entry

What is the complete balanced journal entry for this customer advance?

Hint: Which two changes belong in the seller's books when money arrives before its lesson?

Explanation: Debit Cash $100 and Credit Unearned Revenue $100. Record the tutoring company's money and promise, using its own perspective rather than the customer's books.

- a: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?
- b: Debit Unearned Revenue $100 / Credit Cash $100 (id=opt_dr_unearned_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=true; tag=)
  Stored hint: Which two changes belong in the seller's books when money arrives before its lesson?
- d: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?

### equation_effect

Transaction recap: How does this transaction affect the accounting equation (Assets = Liabilities + Equity)?

Hint: What resource and remaining obligation belong to the tutoring company today?

Explanation: Assets and liabilities each increase by $100. Equity is unchanged because no lesson has been delivered. The customer's separate accounting does not change the seller's performance.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: What resource and remaining obligation belong to the tutoring company today?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities are unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: What resource and remaining obligation belong to the tutoring company today?

## earn_advance_separate_unpaid_job

A consulting company completes an entire prepaid job today. The client paid $100 earlier, and the company already recorded that receipt and its obligation to perform this job. The full amount remains recorded as owed service immediately before completion; none of this job's earning was previously recorded and no work remains afterward. The client also owes a separate, already recorded invoice for a different completed job, which stays unpaid today. No cash changes hands today. Record only completion of the prepaid job; leave the unrelated invoice unchanged.

### identify_account

No new cash changes hands when this prepaid work is completed. Which account recorded the obligation before completion?

Hint: Which recorded promise was satisfied today while the separate invoice stayed unpaid?

Explanation: Unearned Revenue tracks the prepaid job's remaining obligation. Completing this job removes that obligation; the unrelated receivable remains unchanged.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Unearned Revenue (id=opt_unearned_rev; correct=true; tag=)
  Stored hint: Which recorded promise was satisfied today while the separate invoice stayed unpaid?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: For this prepaid work, was the company waiting for money or still owing the promised service?

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

Hint: Does finishing the whole prepaid job leave more or less service still owed for its payment?

Explanation: Unearned Revenue decreases by the full prepaid job price. The separate unpaid invoice does not leave a performance obligation for this completed prepaid job.

- a: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: Does finishing the whole prepaid job leave more or less service still owed for its payment?
- b: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: Does finishing the whole prepaid job leave more or less service still owed for its payment?

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

Hint: What has the company earned by finishing this job, independent of the other unpaid invoice?

Explanation: Service Revenue records today's completed prepaid job. Neither its earlier cash receipt nor the other job's previously recorded earning is repeated.

- a: Service Revenue (id=opt_service_rev; correct=true; tag=)
  Stored hint: What has the company earned by finishing this job, independent of the other unpaid invoice?
- b: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Accounts Receivable (id=opt_ar; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?

### balanced_entry

What is the complete balanced journal entry for earning this advance?

Hint: Which entry removes the completed job's promise without collecting the separate invoice?

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Cash and Accounts Receivable stay unchanged; the separate invoice remains collectible.

- a: Debit Service Revenue $100 / Credit Unearned Revenue $100 (id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 (id=opt_dr_unearned_cr_rev; correct=true; tag=)
  Stored hint: Which entry removes the completed job's promise without collecting the separate invoice?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?

### equation_effect

How does earning a customer advance affect the accounting equation?

Hint: Can completing a prepaid job change obligations and earning while an unrelated claim stays the same?

Explanation: Liabilities decrease and equity increases by $100; total assets are unchanged. Cash and the separate Accounts Receivable balance do not change.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Unearned Revenue); Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Assets increase by $100; Equity increases by $100. (id=opt_assets_up_eq_up; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- c: Liabilities decrease by $100 (-Unearned Revenue); Equity increases by $100 (+Service Revenue); Total Assets unchanged. (id=opt_liab_down_eq_up; correct=true; tag=)
  Stored hint: Can completing a prepaid job change obligations and earning while an unrelated claim stays the same?

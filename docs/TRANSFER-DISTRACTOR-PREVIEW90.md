# Reviewed distractor policy v1 preview

Question wording/answers are the existing template v1. Pedagogy policy v1 versions revised options and setting/contrast metadata. $100, seed 101; full stage/options/teaching/hints below.

## cash_service_stockholder_customer (service_vs_capital)

A corporation completes a repair today and receives $100 cash today from a customer who also owns shares in the corporation. The full payment is solely the agreed price of this completed repair, not a new investment or loan. No shares are issued, no prior invoice, advance, or entry exists for this repair, and no work remains. Record only today's completed service and receipt from the corporation's perspective.

Setting: cash_service.counterparty_role

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

- a: Common Stock (id=opt_common_stock; correct=false; tag=wrong_account)
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

- a: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: Debit Cash $100 / Credit Common Stock $100 (id=opt_dr_cash_cr_stock; correct=false; tag=wrong_account)
  Stored hint: Does this payment buy new ownership or settle the price of completed work?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: What receipt and earning occur when an existing stockholder pays for a completed repair?
- d: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
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

## collect_receivable_future_booking (collection_vs_advance)

A consulting company receives $100 cash today solely to settle a recorded invoice for work completed last month. The company already recorded that earning and the unpaid customer amount then; the receipt equals the full invoice balance. The customer also books a separate appointment for next month, but makes no payment for it and no new work is performed today. No fee, discount, or write-off applies. Record only collection of the old invoice, not an advance for the unpaid future booking.

Setting: collect_receivable.separate_event

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

- a: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=wrong_account)
  Stored hint: What work does this payment cover: the already recorded invoice or an undelivered service?
- d: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: What old balance ends when all of this payment applies to its invoice?

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

## cash_service_undeposited_notes (held_cash_vs_no_entry)

A repair company completes a repair today and receives $100 in physical currency from the customer today. The notes are in the company's locked cash box; the company will deposit them in its bank account tomorrow. This is full payment for the completed repair. No earlier payment, invoice, or accounting entry exists, and no work remains. Record today's completed repair and receipt before the bank deposit.

Setting: cash_service.entity_boundary

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

- a: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=revenue_deferred_when_earned)
  Stored hint: After completing today's work, does the company still owe that work to the customer?
- b: No entry until the cash is deposited in the bank (id=opt_no_entry; correct=false; tag=wrong_account)
  Stored hint: Does the company already control this currency before it reaches the bank?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=true; tag=)
  Stored hint: Have both the company's money and its completed earning changed before the deposit?
- d: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=wrong_account)
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

## collect_receivable_before_due (collection_vs_advance)

A consulting company receives $100 cash today to settle a recorded invoice in full, two weeks before its due date. All the invoiced work was completed last month, when the company recorded the full earning and unpaid customer amount. The receipt equals the invoice balance; no discount, fee, interest, or write-off applies. No new work is performed today. Record only collection of this existing invoice, even though payment arrived before the deadline.

Setting: collect_receivable.paperwork_deadline

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

- a: Debit Accounts Receivable $100 / Credit Cash $100 (id=opt_dr_ar_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=duplicate_revenue_on_collection)
  Stored hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=false; tag=wrong_account)
  Stored hint: What work does this payment cover: the already recorded invoice or an undelivered service?
- d: Debit Cash $100 / Credit Accounts Receivable $100 (id=opt_dr_cash_cr_ar; correct=true; tag=)
  Stored hint: Which recorded customer balance disappears when the early payment arrives?

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

## customer_advance_customer_books (seller_receipt_vs_no_entry)

A tutoring company receives $100 cash today from a business customer in advance for a lesson next month. The customer has entered the payment in its own books, but that does not record any entry for the tutoring company. No instruction or other promised service has been provided, and the tutoring company still owes the full lesson. This is the full price, with no earlier receipt or related entry in the tutoring company's books. Record only the tutoring company's receipt and remaining promise.

Setting: customer_advance.entity_boundary

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

- a: Debit Cash $100 / Credit Unearned Revenue $100 (id=opt_dr_cash_cr_unearned; correct=true; tag=)
  Stored hint: Which two changes belong in the seller's books when money arrives before its lesson?
- b: No seller entry because the customer recorded the payment (id=opt_no_entry; correct=false; tag=wrong_account)
  Stored hint: Does the customer's bookkeeping record the seller's money or finish the seller's promised lesson?
- c: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recognized_prematurely)
  Stored hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- d: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: The customer has already paid. Is there still money to collect for this work?

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

## earn_advance_separate_unpaid_job (earned_advance_vs_receivable)

A consulting company completes an entire prepaid job today. The client paid $100 earlier, and the company already recorded that receipt and its obligation to perform this job. The full amount remains recorded as owed service immediately before completion; none of this job's earning was previously recorded and no work remains afterward. The client also owes a separate, already recorded invoice for a different completed job, which stays unpaid today. No cash changes hands today. Record only completion of the prepaid job; leave the unrelated invoice unchanged.

Setting: earn_advance.separate_event

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

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=cash_recorded_on_earning_advance)
  Stored hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 (id=opt_dr_unearned_cr_rev; correct=true; tag=)
  Stored hint: Which entry removes the completed job's promise without collecting the separate invoice?
- c: Debit Accounts Receivable $100 / Credit Service Revenue $100 (id=opt_dr_ar_cr_rev; correct=false; tag=advance_confused_with_receivable)
  Stored hint: For this completed prepaid job, is its price still owed by the customer?
- d: Debit Service Revenue $100 / Credit Unearned Revenue $100 (id=opt_dr_rev_cr_unearned; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

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

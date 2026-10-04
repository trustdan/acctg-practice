# Quality preview at 90: policy 2

Actual generated options, keys, diagnostic tags and stored hints. Seed 101; $100. Other allowed amounts and scaffolds are checked by regression tests. Unspecified original totals remain symbolic; no amount is invented.

## cash_service_basic

The company completes services today and receives $100 cash today. No invoice was issued and no advance payment was received earlier.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: The customer paid cash today, so Cash is affected. Accounts Receivable would apply only if the customer still owed the money, and Unearned Revenue only if the payment came before the work.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: The work was completed today, so the company has earned revenue: Service Revenue. Unearned Revenue would mean the work is still owed, but it is already done.

## customer_advance_basic

A customer pays $100 cash today for services the company will perform next month. None of the service has been performed.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: The customer paid cash today, so Cash is affected. Service Revenue is tempting, but receiving money is not the same as earning it; the work happens next month.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: The company has been paid but has not done the work, so it owes the customer the service (or a refund). That obligation is Unearned Revenue, a Liability. Service Revenue would record earning before any work is done. Accounts Receivable would mean the customer still owes money, but the customer has already paid.

## service_on_credit_basic

The company completes $100 of services today and invoices the customer. The customer has not paid.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: The company has the right to collect from the customer later, which is Accounts Receivable. Cash is tempting, but no money arrived today.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Revenue is recorded when the work is done, not when cash arrives, so the balancing account is Service Revenue. Cash would record a payment that has not happened.

## collect_receivable_basic

A customer pays $100 cash today to settle an invoice. The company already recorded the service revenue and receivable last month.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: The customer paid money today, so Cash is affected. Service Revenue is tempting, but this payment is for work already recorded as revenue last month.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: The payment settles the amount the customer owed, so the balancing account is Accounts Receivable, which decreases. Crediting Service Revenue again would count last month's work twice.

## cash_service_webdev

A web development studio delivers a finished client website today and receives $100 cash upon delivery today. No prior invoices or advances existed.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: The customer paid cash today, so Cash is affected. Accounts Receivable would apply only if the customer still owed the money, and Unearned Revenue only if the payment came before the work.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: The work was completed today, so the company has earned revenue: Service Revenue. Unearned Revenue would mean the work is still owed, but it is already done.

## service_credit_consulting

A cybersecurity consulting firm completes an audit today and sends an invoice for $100 to the client, payable in 30 days. No cash was received today.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: The company has the right to collect from the customer later, which is Accounts Receivable. Cash is tempting, but no money arrived today.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Revenue is recorded when the work is done, not when cash arrives, so the balancing account is Service Revenue. Cash would record a payment that has not happened.

## collect_receivable_consulting

The consulting firm receives $100 cash today from a client settling an invoice for audit services that were completed and recognized as revenue last month.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: The customer paid money today, so Cash is affected. Service Revenue is tempting, but this payment is for work already recorded as revenue last month.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: The payment settles the amount the customer owed, so the balancing account is Accounts Receivable, which decreases. Crediting Service Revenue again would count last month's work twice.

## customer_advance_logistics

A freight logistics carrier receives $100 cash today from a shipper for cargo transport scheduled to take place next month. None of the shipping service has occurred yet.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: The customer paid cash today, so Cash is affected. Service Revenue is tempting, but receiving money is not the same as earning it; the work happens next month.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: The company has been paid but has not done the work, so it owes the customer the service (or a refund). That obligation is Unearned Revenue, a Liability. Service Revenue would record earning before any work is done. Accounts Receivable would mean the customer still owes money, but the customer has already paid.

## cash_service_training

A training company finishes a one-day workshop today. Participants pay $100 cash at the end of the workshop today. No participant was invoiced or paid before today, and no work remains for these payments.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash records the participants' payment today. Accounts Receivable would mean payment is still owed; Unearned Revenue would mean the workshop is still owed. Neither remains outstanding.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the completed workshop. The cash was earned today; Unearned Revenue would imply that this workshop still needs to be provided.

## customer_advance_training

A training company receives $100 cash today from participants reserving places at a workshop next month. None of the instruction has been provided, and the company must still deliver the full workshop. No payment or related obligation was recorded earlier.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the payment received today. Receiving money does not mean the workshop has been earned; all instruction still lies ahead.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue records the company's obligation to deliver the workshop. Service Revenue would count instruction that has not happened; Accounts Receivable would claim the participants have not paid.

## service_on_credit_repairs

A repair company finishes a customer's repair today and releases the repaired item. It invoices the customer for $100, payable in 30 days. The customer has made no payment or deposit, no work remains, and the company had not previously recorded this repair.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the right to collect for the completed repair. Cash would record a payment that has not happened; Unearned Revenue would describe unfinished work, but the repair is complete.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the completed repair today. An unpaid invoice delays cash collection, not earning; Cash is not the counter-account because no money arrived.

## collect_receivable_repairs

A repair company receives $100 cash today to settle a customer's unpaid repair invoice in full. The repair was finished last month, when the company recorded both the earning and the customer's unpaid amount. No repair or other work is performed for this payment today.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records the payment received today. Service Revenue would count the repair a second time: the earning was already recorded last month.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases because the customer's unpaid amount is settled. Service Revenue was recorded last month; crediting it again would count the same repair twice.

## prepaid_consumption_partial_policy

At month-end, a repair corporation records $100 of insurance coverage used during this month from a policy paid for and recorded as a prepaid asset earlier. The stated amount is only this month's consumed portion, not the full policy price. Additional unused coverage remains for later months. No cash is paid today, and this month's consumption has not previously been recorded.

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

- a: No entry; the full policy remains prepaid. [ID opt_no_entry; tag prepaid_not_expensed_on_consumption]
  Hint: Can coverage already used this month still be reported as protection available for the future?
- b: Debit Insurance Expense / Credit Prepaid Insurance, both for the full policy premium, including unused future coverage [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Has all purchased protection expired, or does some coverage remain available?
- c: Debit Insurance Expense $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag cash_recorded_on_prepaid_expiration]
  Hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- d: Debit Insurance Expense $100 / Credit Prepaid Insurance $100 [ID opt_dr_expense_cr_prepaid; tag ]

Key: opt_dr_expense_cr_prepaid.

Explanation: Debit Insurance Expense $100 and Credit Prepaid Insurance $100. This records this month's portion rather than the whole policy. No entry would leave used coverage in the asset, while a Cash credit would record another payment.

## cash_service_separate_job

A consulting company completes a new, separate job for a returning client today and receives $100 cash for that new job today. This new job had no earlier payment, invoice, or accounting entry, and no work remains for it. The client also owes an older invoice whose work and earning were recorded last month, but none of today's payment applies to that older invoice. Record only the new job and its payment.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash records today's payment for the new job. The older invoice remains unpaid. An existing customer balance does not make every later receipt a collection of that old balance.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the new job earned today. Accounts Receivable would reduce the older unpaid invoice, but the payment is explicitly for the separate new job. Unearned Revenue would imply that this new work remains to be done.

## service_on_credit_new_client

A consulting company completes a job for a new client today and invoices $100, due next month. A different client paid for their own completed job last month, but that earlier payment does not cover this new client's job. No cash or deposit has been received for the new job, no work remains for it, and it has not previously been recorded. Record only today's new completed job.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the unpaid claim for the new job. The earlier payment settled a different client's job. Cash would record money not received for this job, and Unearned Revenue would describe work still owed rather than completed.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the new job completed today. The other client's payment is unrelated; it neither defers this earning nor supplies cash for the new invoice.

## collect_receivable_installment

A consulting company receives $100 cash today as a partial payment of a larger invoice. It completed all the invoiced work last month and recorded both the full earning and the unpaid amount then. The stated amount is only today's installment, not the full invoice total. The rest remains owed, with no discount or write-off. No new work is performed today.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records the installment received today. The work was already earned and recorded last month, so receiving only part of the payment does not create new revenue.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases only by the cash received. The unpaid remainder stays a claim against the client. Service Revenue was already recorded in full and must not be counted again.

### balanced_entry

What is the complete balanced journal entry for collecting this receivable?

- a: Debit Accounts Receivable $100 / Credit Cash $100 [ID opt_dr_ar_cr_cash; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit Cash $100 / Credit Service Revenue $100 [ID opt_dr_cash_cr_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Debit Cash / Credit Accounts Receivable, both for the full original invoice, including the amount still unpaid [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Did today's payment settle the whole invoice, or only the stated installment?
- d: Debit Cash $100 / Credit Accounts Receivable $100 [ID opt_dr_cash_cr_ar; tag ]

Key: opt_dr_cash_cr_ar.

Explanation: Debit Cash $100 and Credit Accounts Receivable $100. Use the installment amount rather than the larger invoice total. Revenue remains unchanged, and the unpaid balance stays receivable.

## customer_advance_additional_receipt

A training company receives $100 cash today as the final payment for a course beginning next month. The customer made a smaller first payment last month, and the company recorded that receipt and its obligation to provide instruction then. None of the instruction has been provided. After today the course is paid for in full. The stated amount is only the additional cash received today, not the total course price. No receivable was recorded for the final payment and today's receipt has not previously been recorded.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the additional payment actually received today. Last month's receipt remains recorded and must not be recorded a second time. The stated amount is only this new receipt.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue increases for the additional advance. An existing obligation does not prevent another advance from increasing it. Service Revenue would count instruction not delivered; no recorded receivable is being settled.

### balanced_entry

What is the complete balanced journal entry for this customer advance?

- a: Debit Cash $100 / Credit Unearned Revenue $100 [ID opt_dr_cash_cr_unearned; tag ]
- b: Debit Cash / Credit Unearned Revenue, both for the entire course price, including the earlier recorded receipt [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Which receipt is new today, and which receipt is already in the books?
- c: Debit Cash $100 / Credit Service Revenue $100 [ID opt_dr_cash_cr_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- d: Debit Accounts Receivable $100 / Credit Service Revenue $100 [ID opt_dr_ar_cr_rev; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?

Key: opt_dr_cash_cr_unearned.

Explanation: Debit Cash $100 and Credit Unearned Revenue $100. Record only the new receipt, not the full course price or last month's payment. The earlier liability remains in place and increases by today's advance.

## earn_advance_completed_session

A training company completes one separately priced session today from a course paid for in full last month. Last month it recorded the payment and its obligation to provide all the sessions. The completed session has a stated price of $100; that amount was included in the recorded advance and has not previously been earned or recognized. Other prepaid sessions are still owed. No cash changes hands today. Record only this completed session.

### balanced_entry

What is the complete balanced journal entry for earning this advance?

- a: Debit Cash $100 / Credit Service Revenue $100 [ID opt_dr_cash_cr_rev; tag cash_recorded_on_earning_advance]
  Hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 [ID opt_dr_unearned_cr_rev; tag ]
- c: Debit Unearned Revenue / Credit Service Revenue, both for the entire prepaid course, including sessions still owed [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Did the company complete every prepaid session, or only the separately priced session stated today?
- d: Debit Service Revenue $100 / Credit Unearned Revenue $100 [ID opt_dr_rev_cr_unearned; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

Key: opt_dr_unearned_cr_rev.

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Use the completed session's price rather than the full course payment. The other sessions remain unearned and no second cash receipt is recorded.

## cash_service_unpaid_booking

A tutoring company completes a lesson today and receives $100 cash for it today. The customer reserved the appointment last week, but that booking involved no payment, invoice, work, or accounting entry. Today's lesson is fully completed and no work remains for this payment. Record only the completed lesson and today's receipt.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash records today's receipt. The earlier booking created neither an advance receipt nor a recorded customer balance.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the lesson completed today. An appointment reservation alone did not earn the lesson price. Accounts Receivable would imply a previously recorded unpaid claim; Unearned Revenue would imply work still owed.

## service_on_credit_due_today

A repair company completes a repair for a new customer today and invoices the customer for $100, due immediately today. The customer has not paid, and the company retains an enforceable claim for the full price. No deposit was received, no work remains, and this repair has not previously been recorded. Record the completed repair while its invoice is still unpaid.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the unpaid claim for the completed repair. An immediate due date does not create Cash before the customer pays. No deposit or unearned obligation exists.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the repair completed today. Earning follows delivery here, not whether the customer met the invoice's due date. No further repair work is owed.

## collect_receivable_overdue

A repair company receives $100 cash today to settle an overdue invoice in full. The repair was completed last month, and the company recorded the full earning and unpaid customer balance then. The amount received equals the recorded invoice balance. There is no interest, late fee, discount, or write-off. No new repair or other work is performed today. Record only collection of the overdue invoice.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records today's overdue payment. Receiving it late does not move the already recorded earning into today.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases to settle the old invoice. Service Revenue was already recorded last month; lateness alone creates no new revenue or advance obligation.

## customer_advance_same_day

At 9 a.m. today, a tutoring company receives $100 cash in advance for a lesson scheduled for 6 p.m. today. No instruction or separate booking service has been provided by 9 a.m.; the company still owes the entire lesson. This is the full lesson price and no payment, invoice, or related obligation was recorded earlier. Record only the morning receipt before the lesson begins.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the morning receipt. The lesson occurring later on the same day does not mean it has already been delivered.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue records the full lesson still owed at 9 a.m. Service Revenue would recognize instruction before delivery. A short wait until evening does not remove the obligation.

## prepaid_purchase_discounted_policy

A company pays $100 cash today for an insurance policy covering next month only. The stated amount is the entire agreed premium after an upfront price reduction; the insurer's higher advertised price was never recorded or owed. No coverage has begun or been used today, no amount was previously recorded as owed, and there are no fees or separate rebates. Record only the purchase at the price actually paid.

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

- a: Debit Prepaid Insurance $100 / Credit Cash $100 [ID opt_dr_prepaid_cr_cash; tag ]
- b: Debit Insurance Expense $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag expense_recorded_on_prepaid_purchase]
  Hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Debit Prepaid Insurance / Credit Cash, both for the higher advertised premium before the agreed price reduction [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Was the advertised price ever owed or paid, or was the reduced agreed premium the entire purchase cost?
- d: Debit Prepaid Insurance $100 / Credit Accounts Payable $100 [ID opt_dr_prepaid_cr_ap; tag payable_recorded_for_cash_payment]
  Hint: Was the purchase paid for today, or is there still an unpaid amount owed?

Key: opt_dr_prepaid_cr_cash.

Explanation: Debit Prepaid Insurance $100 and Credit Cash $100. Record the actual premium, not an unaccepted advertised price or separate discount income. None of the purchased coverage is an expense yet.

## repay_principal_early_partial

A company pays $100 cash today to its bank as a voluntary early repayment of part of the principal on a previously recorded promissory note. The bank accepts the payment before its scheduled due date. The principal balance before payment is larger than today's payment, so some principal remains owed. This payment contains no interest, fee, or penalty, and no new borrowing occurs; ignore interest. Record only today's principal payment.

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

- a: Debit Cash $100 / Credit Notes Payable $100 [ID opt_dr_cash_cr_notes; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit an expense for the loan payment $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag expense_recorded_on_loan_repayment]
  Hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- c: Debit Notes Payable / Credit Cash, both for the entire principal balance before payment, including principal still owed [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Did the bank receive the entire outstanding principal, or only today's accepted partial payment?
- d: Debit Notes Payable $100 / Credit Cash $100 [ID opt_dr_notes_cr_cash; tag ]

Key: opt_dr_notes_cr_cash.

Explanation: Debit Notes Payable $100 and Credit Cash $100. Record the actual payment now, not at the later due date. Do not debit an expense or clear the full note; unpaid principal remains.

## dividend_cash_total_distribution

A corporation's board declares and pays a total cash dividend of $100 today, divided among several stockholders according to their share ownership. The stated amount is the corporation's entire distribution to all recipients together, not an amount for each stockholder. No dividend was declared or recorded as owed earlier. No shares are repurchased, no loan is repaid, and no payment is for work. Record only this total declaration and payment.

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

- a: Debit Cash $100 / Credit Dividends $100 [ID opt_dr_cash_cr_div; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- b: Debit an expense for the dividend $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag expense_recorded_on_dividend]
  Hint: Were the stockholders paid for providing a service, or because they own shares?
- c: Debit Dividends / Credit Cash, both for the stated total multiplied by the number of recipients [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Is the stated amount for each stockholder, or the total paid by the corporation to everyone together?
- d: Debit Dividends $100 / Credit Cash $100 [ID opt_dr_div_cr_cash; tag ]

Key: opt_dr_div_cr_cash.

Explanation: Debit Dividends $100 and Credit Cash $100. The sum paid to all stockholders is the stated amount. Do not multiply it by the number of owners or debit Common Stock: no shares were repurchased.

## cash_service_third_party

A tutoring company completes a lesson today and receives $100 cash today from the customer's parent on the customer's behalf. This is full payment for the completed lesson, not a gift, investment, or loan. No earlier payment or accounting entry exists and no instruction remains owed. Record only today's completed lesson and receipt.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash records the actual receipt. A parent paying for the customer does not change the resource received.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the completed lesson. The payer is settling the service price, not contributing capital or lending money.

## service_on_credit_invoice_later

A repair company completes a repair for a new customer today for an agreed price of $100. The customer owes the full enforceable price now but has not paid. The office will send the invoice next week; sending it is only paperwork and no further repair or customer acceptance is required. No deposit or earlier accounting entry exists. Record the completed unpaid repair today.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the existing claim. The company need not wait for next week's paperwork to record completed, unpaid work.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the repair earned today. The delayed invoice does not postpone completed work.

## collect_receivable_third_party

A repair company receives $100 cash today from a customer's parent, paying on the customer's behalf to settle the full recorded invoice. The repair was completed last month and its full earning and unpaid customer balance were recorded then. This receipt equals that invoice balance; no gift, loan, investment, fee, discount, or new service is involved. Record only today's collection.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records the receipt. Who pays does not move last month's earning into today.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases because the parent settles the existing invoice on the customer's behalf. Revenue was already recorded.

## customer_advance_full_price

A tutoring company receives $100 cash in advance today, the full agreed price of a lesson scheduled for next week. The receipt says paid in full, but no instruction or separate booking service has been delivered and the company still owes the entire lesson. No earlier payment or related accounting entry exists. Record only today's advance receipt.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the full receipt today. The payment receipt establishes payment, not completed instruction.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue records the entire lesson still owed. Full payment removes the customer's unpaid price, not the company's performance obligation.

## equipment_cash_additional_machine

A company pays $100 cash today to buy and take delivery and ownership of an additional machine for several future years of business use. It already owns another machine whose cost was recorded earlier; that older machine stays owned and in use. Today's stated amount is only the new machine's full price, not the combined cost of both machines. No trade-in, repair, disposal, earlier payable, or loan transaction occurs. Record only the new purchase, with no depreciation or other costs.

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

- a: Debit Equipment $100 / Credit Cash $100 [ID opt_dr_equip_cr_cash; tag ]
- b: Debit Cash $100 / Credit Equipment $100 [ID opt_dr_cash_cr_equip; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Equipment / Credit Cash, both for the combined cost of the new and previously recorded machines [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: Which machine was purchased today, and which cost is already recorded?
- d: Debit an expense for the machinery $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag expense_recorded_on_equipment_purchase]
  Hint: Was the machinery's entire benefit used today, or can it provide use over future years?

Key: opt_dr_equip_cr_cash.

Explanation: Debit Equipment $100 and Credit Cash $100. Leave the older machine's recorded cost intact; do not repost the combined cost or expense the new multi-year resource.

## prepaid_consumption_final_remaining

At month-end, a company records $100 of insurance protection used during the final month of a multi-month policy. The policy was paid for and recorded as a prepaid asset earlier, and all earlier months' consumption was already recorded. The stated amount equals the entire remaining prepaid balance immediately before today's entry, not the original full premium. No future coverage remains, no cash is paid today, and this final portion has not previously been expensed. Record only the last remaining portion.

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

- a: No entry; the full policy remains prepaid. [ID opt_no_entry; tag prepaid_not_expensed_on_consumption]
  Hint: Can coverage already used this month still be reported as protection available for the future?
- b: Debit Insurance Expense / Credit Prepaid Insurance, both for the original full premium, including earlier recorded consumption [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: What prepaid balance remained before this entry after earlier months were expensed?
- c: Debit Insurance Expense $100 / Credit Cash $100 [ID opt_dr_expense_cr_cash; tag cash_recorded_on_prepaid_expiration]
  Hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- d: Debit Insurance Expense $100 / Credit Prepaid Insurance $100 [ID opt_dr_expense_cr_prepaid; tag ]

Key: opt_dr_expense_cr_prepaid.

Explanation: Debit Insurance Expense $100 and Credit Prepaid Insurance $100. Use only the remaining balance, not the historical full premium, and do not repeat payment.

## cash_service_stockholder_customer

A corporation completes a repair today and receives $100 cash today from a customer who also owns shares in the corporation. The full payment is solely the agreed price of this completed repair, not a new investment or loan. No shares are issued, no prior invoice, advance, or entry exists for this repair, and no work remains. Record only today's completed service and receipt from the corporation's perspective.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash records today's receipt for the repair. Holding shares does not change the resource received.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Does this payment buy new ownership or settle the price of completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the completed repair. The payer's existing share ownership creates no new Common Stock or loan in this transaction.

## service_on_credit_customer_withdrawal

A repair company completes a repair for a new customer today and invoices $100. Before visiting, the customer withdrew that amount from their own bank account and still holds it, but has given no money to the repair company. The company has an enforceable claim for the full completed repair; no deposit, prior entry, or remaining work exists. Record only the company's completed unpaid service, not the customer's personal withdrawal.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the company's unpaid claim. Cash held by the customer is not cash held by the company.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the completed repair. The customer's personal cash position does not postpone this earning.

## collect_receivable_future_booking

A consulting company receives $100 cash today solely to settle a recorded invoice for work completed last month. The company already recorded that earning and the unpaid customer amount then; the receipt equals the full invoice balance. The customer also books a separate appointment for next month, but makes no payment for it and no new work is performed today. No fee, discount, or write-off applies. Record only collection of the old invoice, not an advance for the unpaid future booking.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records the receipt applied to last month's recorded invoice. No portion pays for the future appointment.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases for the old invoice. Revenue was already recorded; the separate unpaid appointment creates no advance receipt.

## customer_advance_internal_setup

A tutoring company receives $100 cash in advance today for an entire lesson next month. Staff have arranged chairs and tested their own equipment, but no instruction or other promised customer service has been delivered. Setup is internal preparation, not a separately sold or promised service; the company still owes the full lesson. This is the full price and no prior payment or related entry exists. Record only the receipt, with no separate setup cost or payment in this event.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the receipt. Internal preparation does not change the money received into completed instruction.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue records the full lesson still owed. Internal readiness is not a separate promised service earning any part of this price.

## earn_advance_final_remaining

A training company completes the final separately priced session of a prepaid course today. The customer paid for the entire course earlier, and the company recorded the receipt and its obligation then. All earlier sessions and their earning were already recorded. Today's final session has a price of $100, equal to the entire remaining recorded advance immediately before completion, not the original full course price. No work remains for this course and no cash changes hands today. Record only earning the last remaining portion.

### balanced_entry

What is the complete balanced journal entry for earning this advance?

- a: Debit Cash $100 / Credit Service Revenue $100 [ID opt_dr_cash_cr_rev; tag cash_recorded_on_earning_advance]
  Hint: Does completing the prepaid work bring another payment, or was its receipt already recorded?
- b: Debit Unearned Revenue $100 / Credit Service Revenue $100 [ID opt_dr_unearned_cr_rev; tag ]
- c: Debit Unearned Revenue / Credit Service Revenue, both for the original full course price, including earlier recorded earning [ID opt_wrong_amount_scope; tag wrong_amount]
  Hint: What obligation remained immediately before the final session after earlier earning was recorded?
- d: Debit Service Revenue $100 / Credit Unearned Revenue $100 [ID opt_dr_rev_cr_unearned; tag reversed_sides]
  Hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

Key: opt_dr_unearned_cr_rev.

Explanation: Debit Unearned Revenue $100 and Credit Service Revenue $100. Use the remaining final-session price, not the original full course price; do not repeat receipt or earlier earning.

## cash_service_undeposited_notes

A repair company completes a repair today and receives $100 in physical currency from the customer today. The notes are in the company's locked cash box; the company will deposit them in its bank account tomorrow. This is full payment for the completed repair. No earlier payment, invoice, or accounting entry exists, and no work remains. Record today's completed repair and receipt before the bank deposit.

### identify_account

Which account reflects the immediate payment received today?

- a: Cash [ID opt_cash; tag ]
- b: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Which resource arrived today, rather than which ownership account might explain its source?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- d: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: What arrived today, and is the company still waiting to collect that money?

Key: opt_cash.

Explanation: Cash includes currency already held by the company. Depositing these same notes tomorrow does not postpone today's receipt.

### counter_account

Which account records what the company did for the customer in exchange for the cash?

- a: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Was this completed work paid for today, or is its price still owed?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does this receipt pay for completed work, or create money the company must repay to a lender?

Key: opt_service_rev.

Explanation: Service Revenue records the completed repair. A later bank deposit does not make the completed service unearned.

## service_on_credit_cash_sale_label

A repair company completes a repair for a new customer today for an agreed price of $100. The office mistakenly marks its job ticket cash sale, but the customer has made no payment and the company has an enforceable claim for the entire price. The company invoices the customer today. No deposit, prior accounting entry, or remaining work exists. Record the completed unpaid repair using the actual facts, regardless of the ticket label.

### identify_account

The customer owes payment for services completed today and has not paid yet. Which account records this claim?

- a: Service Revenue [ID opt_service_rev; tag wrong_account]
  Hint: Which resource represents the unpaid customer's obligation, rather than the earning itself?
- b: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did the customer pay today, or does the company still have a right to collect later?
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable records the enforceable unpaid claim. The mistaken cash-sale label cannot create money the company did not receive.

### counter_account

The work was completed today but has not been paid for. Which account balances the entry?

- a: Common Stock [ID opt_common_stock; tag wrong_account]
  Hint: Did the customer buy ownership, or receive completed work?
- b: Service Revenue [ID opt_service_rev; tag ]
- c: Unearned Revenue [ID opt_unearned_rev; tag revenue_deferred_when_earned]
  Hint: After completing today's work, does the company still owe that work to the customer?
- d: Cash [ID opt_cash; tag cash_recorded_when_uncollected]
  Hint: Did any payment arrive today, or does the customer still owe money for completed work?

Key: opt_service_rev.

Explanation: Service Revenue records the completed repair. Earning depends on performance here, while the unpaid price remains a receivable.

## collect_receivable_before_due

A consulting company receives $100 cash today to settle a recorded invoice in full, two weeks before its due date. All the invoiced work was completed last month, when the company recorded the full earning and unpaid customer amount. The receipt equals the invoice balance; no discount, fee, interest, or write-off applies. No new work is performed today. Record only collection of this existing invoice, even though payment arrived before the deadline.

### identify_account

Cash arrives to settle the customer's invoice. Which account reflects the cash received today?

- a: Cash [ID opt_cash; tag ]
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Accounts Receivable [ID opt_ar; tag wrong_account]
  Hint: Which resource arrived today, rather than the claim that the payment settles?

Key: opt_cash.

Explanation: Cash records the money received today. Early payment changes the collection date, not the already completed work.

### counter_account

Payment settles an invoice whose revenue was already recognized earlier. What account must be credited?

- a: Unearned Revenue [ID opt_unearned_rev; tag wrong_account]
  Hint: Does this payment settle previously recorded work, or pay for a separate service still owed?
- b: Service Revenue [ID opt_service_rev; tag duplicate_revenue_on_collection]
  Hint: Was this work already recorded as revenue when it was completed, or is new work being earned today?
- c: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: Does the customer owe the company, or is the company repaying a lender?
- d: Accounts Receivable [ID opt_ar; tag ]

Key: opt_ar.

Explanation: Accounts Receivable decreases as the recorded invoice is settled. An early due-date payment is not an advance for future work when the work was already completed and recorded.

## customer_advance_customer_books

A tutoring company receives $100 cash today from a business customer in advance for a lesson next month. The customer has entered the payment in its own books, but that does not record any entry for the tutoring company. No instruction or other promised service has been provided, and the tutoring company still owes the full lesson. This is the full price, with no earlier receipt or related entry in the tutoring company's books. Record only the tutoring company's receipt and remaining promise.

### identify_account

What did the company receive today, and which account records it?

- a: Notes Payable [ID opt_notes_payable; tag wrong_account]
  Hint: What resource arrived from the customer, and is this a loan or payment for future work?
- b: Service Revenue [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash [ID opt_cash; tag ]
- d: Accounts Receivable [ID opt_ar; tag advance_confused_with_receivable]
  Hint: Set aside how the payment was earned or financed. What resource actually arrived today?

Key: opt_cash.

Explanation: Cash records the tutoring company's receipt. The customer's separate books do not record the seller's entry.

### counter_account

What counter-account balances this receipt for services still to be performed?

- a: Service Revenue (Revenue) [ID opt_service_rev; tag revenue_recognized_prematurely]
  Hint: Has the company performed the promised work yet, or does it still owe the customer that work?
- b: Unearned Revenue (Liability) [ID opt_unearned_rev; tag ]
- c: Accounts Receivable (Asset) [ID opt_ar; tag advance_confused_with_receivable]
  Hint: The customer has already paid. Is there still money to collect for this work?
- d: Common Stock (Equity) [ID opt_common_stock; tag wrong_account]
  Hint: Did this customer buy new ownership, or pay for work the company still owes?

Key: opt_unearned_rev.

Explanation: Unearned Revenue records the tutoring company's obligation to provide the lesson. The customer's own entry does not deliver any instruction or earn the seller's revenue.

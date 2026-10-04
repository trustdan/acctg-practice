# Expansion batch 9 rendered review

Draft preview at $100, seed 101. All seven stages include inherited teaching and stored option hints.

## prepaid_purchase_discounted_policy

A company pays $100 cash today for an insurance policy covering next month only. The stated amount is the entire agreed premium after an upfront price reduction; the insurer's higher advertised price was never recorded or owed. No coverage has begun or been used today, no amount was previously recorded as owed, and there are no fees or separate rebates. Record only the purchase at the price actually paid.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Does the lower agreed price change whether the purchased protection is still available for the future?

Explanation: Prepaid Insurance records the unused protection at its actual purchase cost. A price reduction does not make future coverage an expense today or create a separate earning from the insurer's advertised price.

- a: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Does the lower agreed price change whether the purchased protection is still available for the future?
- d: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Does the lower agreed price change whether the purchased protection is still available for the future?

### account_category

What category of account is Prepaid Insurance?

Hint: Is unused coverage a future benefit the company controls or a cost already used up?

Explanation: Prepaid Insurance is an Asset: the right to future coverage. It becomes an expense as coverage is used; paying today does not mean the entire benefit was consumed today.

- a: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is unused coverage a future benefit the company controls or a cost already used up?
- b: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is unused coverage a future benefit the company controls or a cost already used up?
- c: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is unused coverage a future benefit the company controls or a cost already used up?
- d: Expense (id=opt_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?

### direction

Does the Prepaid Insurance account balance increase or decrease upon purchasing the policy?

Hint: Does buying the policy add to the coverage still available or use up coverage already purchased?

Explanation: Prepaid Insurance increases because the company acquired unused coverage. Cash decreases, but a different asset increases.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does buying the policy add to the coverage still available or use up coverage already purchased?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does buying the policy add to the coverage still available or use up coverage already purchased?

### debit_credit

How is an increase to an Asset (Prepaid Insurance) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

### counter_account

Which other account changes when the policy is paid for today?

Hint: How much money actually left today, apart from a higher price that was never owed?

Explanation: Cash decreases by the stated agreed premium. No payable or separate refund exists: the insurer was paid in full at the reduced purchase price.

- a: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: How much money actually left today, apart from a higher price that was never owed?
- c: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- d: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: How much money actually left today, apart from a higher price that was never owed?

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: Unused protection was fully purchased at the agreed price. What benefit arrived and what money left?

Explanation: Debit Prepaid Insurance $100 and Credit Cash $100. Record the actual premium, not an unaccepted advertised price or separate discount income. None of the purchased coverage is an expense yet.

- a: Debit Prepaid Insurance $100 / Credit Accounts Payable $100 (id=opt_dr_prepaid_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Debit Prepaid Insurance $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=true; tag=)
  Stored hint: Unused protection was fully purchased at the agreed price. What benefit arrived and what money left?
- d: Debit Cash $100 / Credit Prepaid Insurance $100 (id=opt_dr_cash_cr_prepaid; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Does a never-recorded advertised price represent an additional resource or obligation in this purchase?

Explanation: Prepaid Insurance increases by $100 and Cash decreases by $100. Total assets, liabilities, and equity are unchanged. The upfront price reduction sets cost; it is not a separate revenue or debt-settlement event.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Does a never-recorded advertised price represent an additional resource or obligation in this purchase?
- c: Asset exchange: Prepaid Insurance increases (+$100) and Cash decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does a never-recorded advertised price represent an additional resource or obligation in this purchase?

## prepaid_consumption_closed_month

At month-end, a company records $100 of insurance protection that expired during this month from a policy paid for and recorded as a prepaid asset earlier. The business was temporarily closed for the month and served no customers, but its policy remained in force throughout the month. The stated amount is this month's used protection and has not previously been expensed. Additional coverage remains for future months. No cash is paid today. Record only the expired portion.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: Did the policy protect the business during the month even while no customers were served?

Explanation: Insurance Expense records protection consumed during the elapsed month. Being temporarily closed does not preserve that month's protection as a future asset or make its cost depend on customer revenue.

- a: Insurance Expense (id=opt_insurance_expense; correct=true; tag=)
  Stored hint: Did the policy protect the business during the month even while no customers were served?
- b: Rent Expense (id=opt_rent_expense; correct=false; tag=wrong_account)
  Stored hint: Did the policy protect the business during the month even while no customers were served?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=false; tag=)
  Stored hint: Did the policy protect the business during the month even while no customers were served?
- d: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?

### account_category

What category of account is Insurance Expense?

Hint: Does the account describe protection still available or protection already used during the month?

Explanation: Insurance Expense is an Expense: the cost of protection already used. Prepaid Insurance is the Asset for protection still available; these accounts describe different parts of the same policy.

- a: Revenue (id=opt_revenue; correct=false; tag=)
  Stored hint: Does the account describe protection still available or protection already used during the month?
- b: Asset (id=opt_asset; correct=false; tag=)
  Stored hint: Does the account describe protection still available or protection already used during the month?
- c: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Does the account describe protection still available or protection already used during the month?
- d: Expense (id=opt_expense; correct=true; tag=)
  Stored hint: Does the account describe protection still available or protection already used during the month?

### direction

Does Insurance Expense increase or decrease as coverage expires during the month?

Hint: Does recognizing the coverage used this month add to this period's costs or remove a previously recorded cost?

Explanation: Insurance Expense increases as consumed coverage is recorded. The prepaid asset decreases; the cost account increases.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does recognizing the coverage used this month add to this period's costs or remove a previously recorded cost?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does recognizing the coverage used this month add to this period's costs or remove a previously recorded cost?

### debit_credit

How is an increase in an Expense (Insurance Expense) recorded?

Hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

Explanation: Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an expense's normal balance?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

### counter_account

Which other account changes when already-paid insurance coverage is used?

Hint: Can a future-benefit balance still include protection for the month that has already ended?

Explanation: Prepaid Insurance decreases by the expired portion. Future months remain prepaid. Cash was paid earlier and must not be credited again merely because the business was closed.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: Can a future-benefit balance still include protection for the month that has already ended?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Can a future-benefit balance still include protection for the month that has already ended?
- d: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Can a future-benefit balance still include protection for the month that has already ended?

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: The month's protection expired while the business was closed. What cost arose and what prepaid portion was used?

Explanation: Debit Insurance Expense $100 and Credit Prepaid Insurance $100. No entry would keep expired protection as an asset. Do not write off the entire policy or repeat its earlier cash payment; future coverage remains.

- a: Debit Insurance Expense $100 / Credit Prepaid Insurance $100 (id=opt_dr_expense_cr_prepaid; correct=true; tag=)
  Stored hint: The month's protection expired while the business was closed. What cost arose and what prepaid portion was used?
- b: Debit Prepaid Insurance $100 / Credit Insurance Expense $100 (id=opt_dr_prepaid_cr_expense; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: No entry; the full policy remains prepaid. (id=opt_no_entry; correct=false; tag=prepaid_not_expensed_on_consumption)
  Stored hint: Can coverage already used this month still be reported as protection available for the future?
- d: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Does having no customer revenue keep the elapsed month's insurance protection available for future use?

Explanation: Assets decrease by $100 (Prepaid Insurance) and equity decreases through Insurance Expense. Cash and liabilities are unchanged. Only the expired portion is consumed; remaining future protection stays an asset.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Assets decrease by $100 (-Prepaid Insurance); Equity decreases by $100 (-Insurance Expense); Cash is unaffected. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does having no customer revenue keep the elapsed month's insurance protection available for future use?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does having no customer revenue keep the elapsed month's insurance protection available for future use?

## equipment_cash_secondhand

A company pays $100 cash today to buy and take delivery and ownership of a secondhand machine from a dealer. The machine was used by its previous owner, but it is expected to serve this company for several future years. It is bought for business use, not resale, and is not a repair service. No amount was previously recorded as owed, no trade-in occurs, and the dealer is paid in full. Record only this purchase, with no depreciation or additional costs.

### identify_account

Which account records the newly acquired machinery?

Hint: Does another owner's earlier use eliminate the machine's several years of future usefulness to this company?

Explanation: Equipment records the acquired machine at its purchase cost. Secondhand does not mean its benefit is already consumed by this buyer. The machine is held for use rather than resale, and no repair service is purchased.

- a: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Does another owner's earlier use eliminate the machine's several years of future usefulness to this company?
- b: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- c: Equipment (id=opt_equipment; correct=true; tag=)
  Stored hint: Does another owner's earlier use eliminate the machine's several years of future usefulness to this company?
- d: An expense for the cost of the machinery (id=opt_expense_misconception; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?

### account_category

What category of account is Equipment?

Hint: Is machinery that can be used in future years a resource still owned or a cost already consumed?

Explanation: Equipment is an Asset: a productive resource available over future years. Paying cash does not make the entire purchase an immediate operating expense.

- a: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is machinery that can be used in future years a resource still owned or a cost already consumed?
- b: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is machinery that can be used in future years a resource still owned or a cost already consumed?
- c: Expense (id=opt_expense; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?
- d: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is machinery that can be used in future years a resource still owned or a cost already consumed?

### direction

Does the Equipment account balance increase or decrease upon acquisition?

Hint: Did this company acquire a machine it did not previously own?

Explanation: Equipment increases by the purchase cost. The dealer's or prior owner's earlier records are not this company's equipment balance. No company machine is sold or traded in.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Did this company acquire a machine it did not previously own?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Did this company acquire a machine it did not previously own?

### debit_credit

How is an increase in an Asset (Equipment) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

### counter_account

Which other account changes when the machinery is paid for today?

Hint: Did this company pay the dealer now or leave any purchase price unpaid?

Explanation: Cash decreases because the dealer was paid in full today. No trade payable, new loan, or share issuance arises from this paid purchase.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Did this company pay the dealer now or leave any purchase price unpaid?
- b: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Did this company pay the dealer now or leave any purchase price unpaid?
- c: Common Stock (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: Did this company pay the dealer now or leave any purchase price unpaid?
- d: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

Hint: The bought machine still has years of future use for this company. What resource arrives and what resource leaves?

Explanation: Debit Equipment $100 and Credit Cash $100. Buying secondhand equipment is not immediate consumption of its entire future benefit. Use this company's purchase cost rather than the prior owner's cost or history.

- a: Debit an expense for the machinery $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?
- b: Debit Cash $100 / Credit Equipment $100 (id=opt_dr_cash_cr_equip; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Equipment $100 / Credit Cash $100 (id=opt_dr_equip_cr_cash; correct=true; tag=)
  Stored hint: The bought machine still has years of future use for this company. What resource arrives and what resource leaves?

### equation_effect

How does purchasing equipment for cash affect the accounting equation?

Hint: After paying the dealer, does this company own a resource useful in future years?

Explanation: Equipment increases by $100 and Cash decreases by $100. Total assets, liabilities, and equity are unchanged. Secondhand condition does not turn the whole purchase into a current expense; no depreciation is recorded here.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?
- b: Asset exchange: Equipment increases (+$100) and Cash decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: After paying the dealer, does this company own a resource useful in future years?
- c: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: After paying the dealer, does this company own a resource useful in future years?

## repay_principal_early_partial

A company pays $100 cash today to its bank as a voluntary early repayment of part of the principal on a previously recorded promissory note. The bank accepts the payment before its scheduled due date. The principal balance before payment is larger than today's payment, so some principal remains owed. This payment contains no interest, fee, or penalty, and no new borrowing occurs; ignore interest. Record only today's principal payment.

### identify_account

Which account records what the company owed the bank before today's principal payment?

Hint: Does paying before the scheduled due date still settle part of an existing repayment obligation?

Explanation: Notes Payable records the principal obligation being reduced. The payment date does not make principal a current expense or postpone recording an actual accepted repayment.

- a: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Does paying before the scheduled due date still settle part of an existing repayment obligation?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Does paying before the scheduled due date still settle part of an existing repayment obligation?
- c: An expense for the loan payment (id=opt_expense_misconception; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- d: Notes Payable (id=opt_notes_payable; correct=true; tag=)
  Stored hint: Does paying before the scheduled due date still settle part of an existing repayment obligation?

### account_category

What category of account is Notes Payable?

Hint: Does a signed promise to repay the bank describe a resource owned or an obligation owed?

Explanation: Notes Payable is a Liability: principal the company must repay. An expense records a cost incurred; the loan principal is previously borrowed money, not a new expense.

- a: Liability (id=opt_liability; correct=true; tag=)
  Stored hint: Does a signed promise to repay the bank describe a resource owned or an obligation owed?
- b: Expense (id=opt_expense; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- c: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Does a signed promise to repay the bank describe a resource owned or an obligation owed?
- d: Asset (id=opt_asset; correct=false; tag=)
  Stored hint: Does a signed promise to repay the bank describe a resource owned or an obligation owed?

### direction

Does Notes Payable increase or decrease when principal is repaid?

Hint: After the accepted partial payment, is less principal owed even though the note is not fully settled?

Explanation: Notes Payable decreases by the stated payment. The unpaid remainder stays a liability; neither the entire original balance nor zero is the new balance.

- a: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: After the accepted partial payment, is less principal owed even though the note is not fully settled?
- b: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: After the accepted partial payment, is less principal owed even though the note is not fully settled?

### debit_credit

How is a decrease in a Liability (Notes Payable) recorded?

Hint: Which side is a liability's normal balance, and does a decrease use that side or the opposite side?

Explanation: Liabilities have a normal credit balance, so a decrease is a Debit (left side). The debit reduces debt; it does not mean cash increased.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: Which side is a liability's normal balance, and does a decrease use that side or the opposite side?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: Which side is a liability's normal balance, and does a decrease use that side or the opposite side?

### counter_account

Which other account changes when the company pays principal today?

Hint: What resource left when the bank accepted the early payment?

Explanation: Cash decreases today. A later scheduled due date does not leave money already transferred in Cash. The payment creates no new note or trade payable.

- a: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: What resource left when the bank accepted the early payment?
- b: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource left when the bank accepted the early payment?
- c: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: What resource left when the bank accepted the early payment?

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

Hint: The bank accepted only part of the existing principal early. What obligation became smaller and what resource left?

Explanation: Debit Notes Payable $100 and Credit Cash $100. Record the actual payment now, not at the later due date. Do not debit an expense or clear the full note; unpaid principal remains.

- a: Debit Notes Payable $100 / Credit Cash $100 (id=opt_dr_notes_cr_cash; correct=true; tag=)
  Stored hint: The bank accepted only part of the existing principal early. What obligation became smaller and what resource left?
- b: Debit an expense for the loan payment $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- c: Debit Cash $100 / Credit Notes Payable $100 (id=opt_dr_cash_cr_notes; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does repaying loan principal affect the accounting equation?

Hint: Does an early principal payment create a cost, or remove part of a recorded obligation?

Explanation: Assets decrease by $100 (Cash) and liabilities decrease by $100 (Notes Payable). Equity is unchanged. No interest, fee, or penalty affects income, and the remaining principal stays owed.

- a: Assets decrease by $100 (-Cash); Liabilities decrease by $100 (-Notes Payable); Equity is unchanged. (id=opt_assets_down_liab_down; correct=true; tag=)
  Stored hint: Does an early principal payment create a cost, or remove part of a recorded obligation?
- b: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does an early principal payment create a cost, or remove part of a recorded obligation?
- c: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?

## dividend_cash_total_distribution

A corporation's board declares and pays a total cash dividend of $100 today, divided among several stockholders according to their share ownership. The stated amount is the corporation's entire distribution to all recipients together, not an amount for each stockholder. No dividend was declared or recorded as owed earlier. No shares are repurchased, no loan is repaid, and no payment is for work. Record only this total declaration and payment.

### identify_account

Which account records cash distributions paid directly to stockholders?

Hint: Are these payments for the owners' work or distributions because they hold shares?

Explanation: Dividends records the total newly declared ownership distribution. Splitting it among several stockholders does not create an operating expense, buy back shares, or repay a note.

- a: An expense for the payment to stockholders (id=opt_expense_misconception; correct=false; tag=expense_recorded_on_dividend)
  Stored hint: Were the stockholders paid for providing a service, or because they own shares?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Are these payments for the owners' work or distributions because they hold shares?
- c: Dividends (id=opt_dividends; correct=true; tag=)
  Stored hint: Are these payments for the owners' work or distributions because they hold shares?
- d: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Are these payments for the owners' work or distributions because they hold shares?

### account_category

What category of account is Dividends?

Hint: Does this account track operating costs, a resource owned, or distributions to owners?

Explanation: Dividends is the Dividends category, which reduces equity. It is not an Expense and does not reduce net income; the account is closed to Retained Earnings at period end.

- a: Dividends (id=opt_dividends_cat; correct=true; tag=)
  Stored hint: Does this account track operating costs, a resource owned, or distributions to owners?
- b: Asset (id=opt_asset; correct=false; tag=)
  Stored hint: Does this account track operating costs, a resource owned, or distributions to owners?
- c: Expense (id=opt_expense; correct=false; tag=expense_recorded_on_dividend)
  Stored hint: Were the stockholders paid for providing a service, or because they own shares?
- d: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Does this account track operating costs, a resource owned, or distributions to owners?

### direction

Does the Dividends account balance increase or decrease when recording this dividend distribution?

Hint: Does dividing the distribution among owners remove the distribution already declared today?

Explanation: Dividends increases by the total stated amount. The number of recipients does not multiply that total. This account records distributions rather than an increase in ownership capital.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does dividing the distribution among owners remove the distribution already declared today?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does dividing the distribution among owners remove the distribution already declared today?

### debit_credit

How is an increase in Dividends recorded?

Hint: Which side is the Dividends account's normal balance, and where does an increase belong?

Explanation: Dividends has a normal debit balance, so an increase is a Debit (left side). Crediting Dividends would reduce recorded distributions; an equity reduction does not mean the Dividends account decreases.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: Which side is the Dividends account's normal balance, and where does an increase belong?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: Which side is the Dividends account's normal balance, and where does an increase belong?

### counter_account

Which other account changes when this newly declared distribution is paid today?

Hint: How much company money leaves in total, regardless of how it is split among recipients?

Explanation: Cash decreases by the entire stated distribution once. It is not that amount per stockholder. No unpaid dividend remains after today's payment.

- a: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: How much company money leaves in total, regardless of how it is split among recipients?
- b: Dividends Payable (id=opt_dividends_payable; correct=false; tag=wrong_account)
  Stored hint: How much company money leaves in total, regardless of how it is split among recipients?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: How much company money leaves in total, regardless of how it is split among recipients?

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

Hint: Record the company's total distribution, not a per-owner amount. What distribution and payment arose?

Explanation: Debit Dividends $100 and Credit Cash $100. The sum paid to all stockholders is the stated amount. Do not multiply it by the number of owners or debit Common Stock: no shares were repurchased.

- a: Debit Dividends $100 / Credit Cash $100 (id=opt_dr_div_cr_cash; correct=true; tag=)
  Stored hint: Record the company's total distribution, not a per-owner amount. What distribution and payment arose?
- b: Debit an expense for the dividend $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_dividend)
  Stored hint: Were the stockholders paid for providing a service, or because they own shares?
- c: Debit Cash $100 / Credit Dividends $100 (id=opt_dr_cash_cr_div; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does paying cash dividends affect the accounting equation?

Hint: Does splitting a fixed total among several recipients change the corporation's total outflow?

Explanation: Assets decrease by $100 (Cash) and equity decreases by $100 through Dividends. Liabilities and net income are unchanged. Record the total once; owner allocations do not increase the company's distribution.

- a: Assets decrease by $100; Liabilities decrease by $100. (id=opt_assets_down_liab_down; correct=false; tag=equation_effect_missed)
  Stored hint: Does splitting a fixed total among several recipients change the corporation's total outflow?
- b: No change in equity; dividends are an asset. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does splitting a fixed total among several recipients change the corporation's total outflow?
- c: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Dividends); Liabilities are unchanged. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does splitting a fixed total among several recipients change the corporation's total outflow?

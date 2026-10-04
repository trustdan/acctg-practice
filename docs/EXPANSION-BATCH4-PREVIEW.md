# Expansion batch 4: complete teaching preview

Version-1 snapshots at $50, seed 100. All allowed amounts and scaffold levels are regression-tested. Shared stages and scenario overrides appear with every option, tag, and stored mistake hint. Review: [EXPANSION-BATCH4-REVIEW.md](EXPANSION-BATCH4-REVIEW.md).

## prepaid_purchase_renewal v1

A repair corporation pays $50 cash today to renew its insurance for a period beginning next month. Its existing policy still covers the current month. None of the new renewal coverage has begun or been used today, and no amount for the renewal was previously recorded as owed. Record only the payment for the future renewal, not any consumption of the existing policy.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Is today's payment buying protection already used or protection beginning after the current policy period?

Explanation: Prepaid Insurance records the unused renewal coverage. The existing policy protecting the business today does not make the newly purchased future coverage an expense. Cash records how the renewal was paid for.

- opt_cash: Cash (tag: ; hint: )
- opt_insurance_expense: Insurance Expense (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )

### account_category

What category of account is Prepaid Insurance?

Hint: Is unused coverage a future benefit the company controls or a cost already used up?

Explanation: Prepaid Insurance is an Asset: the right to future coverage. It becomes an expense as coverage is used; paying today does not mean the entire benefit was consumed today.

- opt_liability: Liability (tag: ; hint: )
- opt_expense: Expense (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_asset: Asset (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )

### direction

Does the Prepaid Insurance account balance increase or decrease upon purchasing the policy?

Hint: Does buying the policy add to the coverage still available or use up coverage already purchased?

Explanation: Prepaid Insurance increases because the company acquired unused coverage. Cash decreases, but a different asset increases.

- opt_decrease: Decrease (tag: ; hint: )
- opt_increase: Increase (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Prepaid Insurance) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

Which other account changes when the policy is paid for today?

Hint: Did the corporation pay the insurer now or leave the renewal price unpaid?

Explanation: Cash decreases because the insurer was paid today, so Cash is credited. A payable would report an unpaid renewal price, but no amount for this renewal remains owed.

- opt_insurance_expense: Insurance Expense (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Did the corporation pay the insurer now or leave the renewal price unpaid?)
- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: Set aside the existing policy. What future benefit did this renewal payment acquire, and what resource left today?

Explanation: Debit Prepaid Insurance $50 and Credit Cash $50. The new coverage starts next month. Expensing this renewal now would confuse the current policy's protection with the unused protection just purchased.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_dr_prepaid_cr_ap: Debit Prepaid Insurance $50 / Credit Accounts Payable $50 (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_dr_cash_cr_prepaid: Debit Cash $50 / Credit Prepaid Insurance $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_prepaid_cr_cash: Debit Prepaid Insurance $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Does paying for next month's protection leave a future resource after the cash payment?

Explanation: Prepaid Insurance increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged. Consumption of the separate current policy is outside this event.

- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Does paying for next month's protection leave a future resource after the cash payment?)
- opt_asset_swap: Asset exchange: Prepaid Insurance increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)

## prepaid_consumption_partial_policy v1

At month-end, a repair corporation records $50 of insurance coverage used during this month from a policy paid for and recorded as a prepaid asset earlier. The stated amount is only this month's consumed portion, not the full policy price. Additional unused coverage remains for later months. No cash is paid today, and this month's consumption has not previously been recorded.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: Can this month's used protection remain a future benefit just because other months are still covered?

Explanation: Insurance Expense records only the coverage used this month. The unused portion remains Prepaid Insurance. An unexpired policy can still contain a consumed portion that must be recognized as a cost.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_insurance_expense: Insurance Expense (tag: ; hint: )
- opt_rent_expense: Rent Expense (tag: wrong_account; hint: Can this month's used protection remain a future benefit just because other months are still covered?)

### account_category

What category of account is Insurance Expense?

Hint: Does the account describe protection still available or protection already used during the month?

Explanation: Insurance Expense is an Expense: the cost of protection already used. Prepaid Insurance is the Asset for protection still available; these accounts describe different parts of the same policy.

- opt_expense: Expense (tag: ; hint: )
- opt_revenue: Revenue (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )

### direction

Does Insurance Expense increase or decrease as coverage expires during the month?

Hint: Does recognizing this month's consumed portion add to or reverse insurance costs recorded this period?

Explanation: Insurance Expense increases by the consumed portion. The remaining policy coverage is still an asset, but it does not cancel the cost of protection already used.

- opt_decrease: Decrease (tag: ; hint: )
- opt_increase: Increase (tag: ; hint: )

### debit_credit

How is an increase in an Expense (Insurance Expense) recorded?

Hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

Explanation: Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

Which other account changes when already-paid insurance coverage is used?

Hint: Which earlier balance contains the coverage that was partly used this month?

Explanation: Prepaid Insurance decreases by the consumed portion, so it is credited. Only this month's used coverage is removed; unused later-month coverage remains. Crediting Cash would repeat a payment made earlier.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Which earlier balance contains the coverage that was partly used this month?)
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Which earlier balance contains the coverage that was partly used this month?)
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: Only the stated portion has been used and the policy was paid for earlier. What cost arose and which future benefit became smaller?

Explanation: Debit Insurance Expense $50 and Credit Prepaid Insurance $50. This records this month's portion rather than the whole policy. No entry would leave used coverage in the asset, while a Cash credit would record another payment.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_dr_prepaid_cr_expense: Debit Prepaid Insurance $50 / Credit Insurance Expense $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_prepaid: Debit Insurance Expense $50 / Credit Prepaid Insurance $50 (tag: ; hint: )
- opt_no_entry: No entry; the full policy remains prepaid. (tag: prepaid_not_expensed_on_consumption; hint: Can coverage already used this month still be reported as protection available for the future?)

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Does having some future coverage left make this month's used coverage a resource still available?

Explanation: Assets decrease by $50 (Prepaid Insurance) and equity decreases by $50 through Insurance Expense. Cash and liabilities are unchanged. The unused remainder stays an asset; only this month's consumption affects this entry.

- opt_assets_down_eq_down: Assets decrease by $50 (-Prepaid Insurance); Equity decreases by $50 (-Insurance Expense); Cash is unaffected. (tag: ; hint: )
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does having some future coverage left make this month's used coverage a resource still available?)
- opt_assets_down_cash: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

## repay_principal_equipment_loan v1

A repair corporation pays $50 cash today to its bank toward principal on a previously recorded promissory note. The loan originally helped buy a machine, whose purchase was recorded earlier. The machine remains owned and in use; no sale or disposal occurs today. This payment contains principal only; ignore interest. Record only this repayment, with no depreciation or new purchase.

### identify_account

Which account records what the company owed the bank before today's principal payment?

Hint: What existing obligation is reduced when the bank receives previously borrowed principal?

Explanation: Notes Payable records the principal owed to the bank. Repayment reduces that obligation. The machine purchase was already recorded and the machine is kept, so today's loan payment neither acquires nor disposes of equipment.

- opt_expense_misconception: An expense for the loan payment (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What existing obligation is reduced when the bank receives previously borrowed principal?)
- opt_notes_payable: Notes Payable (tag: ; hint: )
- opt_cash: Cash (tag: ; hint: )

### account_category

What category of account is Notes Payable?

Hint: Does a signed promise to repay the bank describe a resource owned or an obligation owed?

Explanation: Notes Payable is a Liability: principal the company must repay. An expense records a cost incurred; the loan principal is previously borrowed money, not a new expense.

- opt_asset: Asset (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )
- opt_expense: Expense (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_liability: Liability (tag: ; hint: )

### direction

Does Notes Payable increase or decrease when principal is repaid?

Hint: After repaying part of the principal, does the company owe the bank more or less?

Explanation: Notes Payable decreases as principal is settled. Borrowing would increase the debt, but today's event is repayment.

- opt_decrease: Decrease (tag: ; hint: )
- opt_increase: Increase (tag: ; hint: )

### debit_credit

How is a decrease in a Liability (Notes Payable) recorded?

Hint: Which side is a liability's normal balance, and does a decrease use that side or the opposite side?

Explanation: Liabilities have a normal credit balance, so a decrease is a Debit (left side). The debit reduces debt; it does not mean cash increased.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

Which other account changes when the company pays principal today?

Hint: What resource did the corporation transfer to the bank today?

Explanation: Cash decreases because the bank was paid today, so Cash is credited. The machine remains owned; its use as the original purpose of the loan does not make it the resource given up in this repayment.

- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: What resource did the corporation transfer to the bank today?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What resource did the corporation transfer to the bank today?)

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

Hint: The machine purchase is already recorded and only loan principal is settled today. What obligation shrinks and what resource leaves?

Explanation: Debit Notes Payable $50 and Credit Cash $50. Do not expense the principal or record the machine purchase again. No equipment is disposed of, and no interest cost is included.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_expense_cr_cash: Debit an expense for the loan payment $50 / Credit Cash $50 (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)

### equation_effect

How does repaying loan principal affect the accounting equation?

Hint: Does repaying financing for an earlier purchase create a new operating cost today?

Explanation: Assets decrease by $50 (Cash) and liabilities decrease by $50 (Notes Payable). Equity and Equipment are unchanged. Repaying principal does not reduce net income or remove the machine.

- opt_assets_down_liab_down: Assets decrease by $50 (-Cash); Liabilities decrease by $50 (-Notes Payable); Equity is unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does repaying financing for an earlier purchase create a new operating cost today?)

## dividend_cash_working_owners v1

A repair corporation's board declares and pays a $50 cash distribution to its stockholders today in proportion to their share ownership. Some stockholders also work for the company, but this payment is solely an ownership distribution, not pay for their work. No distribution was declared or recorded as owed earlier, and the payment does not repay any loan.

### identify_account

Which account records cash distributions paid directly to stockholders?

Hint: Is the payment based on work performed or on the recipients' share ownership?

Explanation: Dividends records the newly declared ownership distribution. Stockholders can also be employees, but this payment is based on shares rather than work. Treating it as an expense would misstate income.

- opt_dividends: Dividends (tag: ; hint: )
- opt_cash: Cash (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Is the payment based on work performed or on the recipients' share ownership?)
- opt_expense_misconception: An expense for the payment to stockholders (tag: expense_recorded_on_dividend; hint: Were the stockholders paid for providing a service, or because they own shares?)

### account_category

What category of account is Dividends?

Hint: Does this account track operating costs, a resource owned, or distributions to owners?

Explanation: Dividends is the Dividends category, which reduces equity. It is not an Expense and does not reduce net income; the account is closed to Retained Earnings at period end.

- opt_expense: Expense (tag: expense_recorded_on_dividend; hint: Were the stockholders paid for providing a service, or because they own shares?)
- opt_liability: Liability (tag: ; hint: )
- opt_dividends_cat: Dividends (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )

### direction

Does the Dividends account balance increase or decrease when recording this dividend distribution?

Hint: Does declaring this distribution add to the dividends accumulated this period or reverse a previous distribution?

Explanation: Dividends increases as a new distribution is recorded, even though equity decreases. The distribution account and total equity move in opposite directions.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase in Dividends recorded?

Hint: Which side is the Dividends account's normal balance, and where does an increase belong?

Explanation: Dividends has a normal debit balance, so an increase is a Debit (left side). Crediting Dividends would reduce recorded distributions; an equity reduction does not mean the Dividends account decreases.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

Which other account changes when this newly declared distribution is paid today?

Hint: Was this distribution merely promised for later or actually paid today?

Explanation: Cash decreases because declaration and payment both occur today. Dividends Payable would describe a distribution still owed, but no amount for this distribution remains unpaid.

- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Was this distribution merely promised for later or actually paid today?)
- opt_cash: Cash (tag: ; hint: )
- opt_dividends_payable: Dividends Payable (tag: wrong_account; hint: Was this distribution merely promised for later or actually paid today?)

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

Hint: No earlier distribution payable exists and this payment is for ownership rather than work. What distribution arose and what resource left?

Explanation: Debit Dividends $50 and Credit Cash $50. Employee status does not turn an ownership distribution into wages. Neither a loan nor an earlier dividend payable is being settled.

- opt_dr_cash_cr_div: Debit Cash $50 / Credit Dividends $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the dividend $50 / Credit Cash $50 (tag: expense_recorded_on_dividend; hint: Were the stockholders paid for providing a service, or because they own shares?)
- opt_dr_div_cr_cash: Debit Dividends $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does paying cash dividends affect the accounting equation?

Hint: Does a payment based on share ownership reduce operating income just because some owners also work here?

Explanation: Assets decrease by $50 (Cash) and equity decreases by $50 through Dividends. Liabilities and net income are unchanged. The ownership distribution reduces equity directly rather than creating an expense.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Dividends); Liabilities are unchanged. (tag: ; hint: )
- opt_assets_down_liab_down: Assets decrease by $50; Liabilities decrease by $50. (tag: equation_effect_missed; hint: Does a payment based on share ownership reduce operating income just because some owners also work here?)
- opt_no_net_change: No change in equity; dividends are an asset. (tag: equation_effect_missed; hint: Does a payment based on share ownership reduce operating income just because some owners also work here?)

## borrow_cash_stockholder_note v1

A repair corporation receives $50 cash today from one of its stockholders after signing a promissory note requiring repayment. The stockholder is acting as a lender for this transaction and receives no new shares. This is not payment for any customer work. Record only today's loan receipt; ignore interest.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource actually arrived today regardless of the lender's existing ownership?

Explanation: Cash records the money received today. The stockholder's existing shares do not change the resource received or make the loan receipt revenue.

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

Hint: Does the signed agreement require repayment or grant new ownership for this payment?

Explanation: Notes Payable records the obligation under the promissory note. Common Stock would mean new shares were issued, but none were. An existing owner can lend money without making another ownership contribution.

- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_notes_payable: Notes Payable (Liability) (tag: ; hint: )
- opt_common_stock: Common Stock (Equity) (tag: wrong_account; hint: Does the signed agreement require repayment or grant new ownership for this payment?)

### balanced_entry

What is the complete balanced journal entry for borrowing cash on a note?

Hint: The corporation received repayable money and issued no new shares. What resource and obligation arose?

Explanation: Debit Cash $50 and Credit Notes Payable $50. The source being a stockholder does not make this contributed capital. Crediting revenue would call a loan earning; reversing the entry would record repayment rather than receipt.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: ; hint: )
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)

### equation_effect

How does borrowing cash affect the accounting equation?

Hint: Does the lender already holding shares remove the corporation's promise to repay this new money?

Explanation: Assets increase by $50 (Cash) and liabilities increase by $50 (Notes Payable). Equity and net income are unchanged. The lender's existing ownership is separate from this loan obligation.

- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue); Liabilities unchanged. (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_no_net_change: No net change in total assets. (tag: equation_effect_missed; hint: Does the lender already holding shares remove the corporation's promise to repay this new money?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Notes Payable); Equity is unchanged. (tag: ; hint: )

# Expansion batch 6: complete teaching preview

Version-1 snapshots at $50, seed 100. All allowed amounts and scaffolds are regression-tested. Shared stages and overrides include all options, tags, and stored mistake hints. Review: [EXPANSION-BATCH6-REVIEW.md](EXPANSION-BATCH6-REVIEW.md).

## equipment_cash_existing_note v1

A company pays $50 from its existing cash today to buy and take ownership of a machine for several years of use. The company already has an unrelated bank loan recorded from last month, but it neither borrows nor repays any loan today. The dealer is paid in full and no amount for this machine was previously recorded as owed. Record only the delivered machine purchase, with no depreciation.

### identify_account

Which account records the newly acquired machinery?

Hint: What resource did the company acquire today, apart from its earlier bank loan?

Explanation: Equipment records the new machine owned for future use. The old loan is a separate obligation and does not make the machine an immediate expense or an unpaid purchase.

- opt_equipment: Equipment (tag: ; hint: )
- opt_expense_misconception: An expense for the cost of the machinery (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)
- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### account_category

What category of account is Equipment?

Hint: Is machinery that can be used in future years a resource still owned or a cost already consumed?

Explanation: Equipment is an Asset: a productive resource available over future years. Paying cash does not make the entire purchase an immediate operating expense.

- opt_liability: Liability (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )
- opt_expense: Expense (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)

### direction

Does the Equipment account balance increase or decrease upon acquisition?

Hint: After this purchase, does the company own more machinery or less?

Explanation: Equipment increases because new machinery was acquired. Cash decreases, but the Equipment account tracks machinery, not money.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase in an Asset (Equipment) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

Which other account changes when the machinery is paid for today?

Hint: Was the dealer paid from money already held or by a new promise to repay a lender?

Explanation: Cash decreases because existing money paid the dealer today. Notes Payable already has a balance, but no new borrowing or repayment occurs. Accounts Payable would report a dealer debt that this paid purchase does not create.

- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: Was the dealer paid from money already held or by a new promise to repay a lender?)
- opt_common_stock: Common Stock (tag: wrong_account; hint: Was the dealer paid from money already held or by a new promise to repay a lender?)
- opt_cash: Cash (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

Hint: Record only the machine bought with existing cash. What resource arrived and what resource left?

Explanation: Debit Equipment $50 and Credit Cash $50. The unrelated note is unchanged. A note or payable credit would create debt not incurred in this purchase; an expense debit would consume the entire multi-year benefit immediately.

- opt_dr_equip_cr_cash: Debit Equipment $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_cash_cr_equip: Debit Cash $50 / Credit Equipment $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the machinery $50 / Credit Cash $50 (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)

### equation_effect

How does purchasing equipment for cash affect the accounting equation?

Hint: Does merely having an old loan turn this fully paid purchase into another borrowing?

Explanation: Equipment increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged. The previously recorded bank loan remains unchanged.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Does merely having an old loan turn this fully paid purchase into another borrowing?)
- opt_asset_swap: Asset exchange: Equipment increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )

## prepaid_purchase_machine_cover v1

A company pays $50 cash today for an insurance policy protecting machinery it already owns. All the purchased coverage begins next month and none has been used today. No machine is bought, repaired, or improved by this payment. No amount for the policy was previously recorded as owed. Record only the insurance purchase.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Did the payment acquire another machine or a right to insurance protection not yet used?

Explanation: Prepaid Insurance records the future protection purchased. Equipment records the machinery already owned, which this policy payment does not increase. Insurance Expense would recognize protection before any purchased coverage has been used.

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

Hint: Did the company pay the insurer today or leave the policy price outstanding?

Explanation: Cash decreases because the insurer was paid today, so Cash is credited. Accounts Payable would imply an unpaid amount, while an expense counter-account would confuse payment with consumption.

- opt_insurance_expense: Insurance Expense (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Did the company pay the insurer today or leave the policy price outstanding?)
- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: Future protection was purchased with cash for a machine already owned. What benefit arrived and what resource left?

Explanation: Debit Prepaid Insurance $50 and Credit Cash $50. The object insured is machinery, but the resource purchased is unused protection. Do not increase Equipment or expense coverage before it is used.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_dr_prepaid_cr_ap: Debit Prepaid Insurance $50 / Credit Accounts Payable $50 (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_dr_cash_cr_prepaid: Debit Cash $50 / Credit Prepaid Insurance $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_prepaid_cr_cash: Debit Prepaid Insurance $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Does buying future protection leave a resource after paying the insurer?

Explanation: Prepaid Insurance increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged. Equipment is unchanged because no machine or improvement was acquired.

- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Does buying future protection leave a resource after paying the insurer?)
- opt_asset_swap: Asset exchange: Prepaid Insurance increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)

## prepaid_consumption_no_claim v1

At month-end, a company records $50 of insurance protection used during this month from a policy paid for and recorded as a prepaid asset last month. No insured accident occurred and no claim was filed, but this month's protection has expired. The stated amount is the cost of this month's used coverage and has not previously been expensed. No cash is paid today. Record only this consumption of prepaid coverage.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: Was protection available and used during the month even though no accident occurred?

Explanation: Insurance Expense records the cost of protection consumed with time. The absence of a claim does not mean the company received no protection or that expired coverage remains an asset.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_insurance_expense: Insurance Expense (tag: ; hint: )
- opt_rent_expense: Rent Expense (tag: wrong_account; hint: Was protection available and used during the month even though no accident occurred?)

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

Hint: Does recognizing the coverage used this month add to this period's costs or remove a previously recorded cost?

Explanation: Insurance Expense increases as consumed coverage is recorded. The prepaid asset decreases; the cost account increases.

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

Hint: Can last month's prepaid balance still include protection for a month that has ended?

Explanation: Prepaid Insurance decreases by the expired portion and is credited. The policy was paid for earlier, so Cash must not be credited again. Filing a claim is not required to recognize consumed coverage.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Can last month's prepaid balance still include protection for a month that has ended?)
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Can last month's prepaid balance still include protection for a month that has ended?)
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: The paid-for protection expired without a new payment. What cost arose and what future benefit became smaller?

Explanation: Debit Insurance Expense $50 and Credit Prepaid Insurance $50. No entry would leave expired protection in the asset. A Cash credit would repeat the earlier payment; no accident or claim is needed for this expense.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_dr_prepaid_cr_expense: Debit Prepaid Insurance $50 / Credit Insurance Expense $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_prepaid: Debit Insurance Expense $50 / Credit Prepaid Insurance $50 (tag: ; hint: )
- opt_no_entry: No entry; the full policy remains prepaid. (tag: prepaid_not_expensed_on_consumption; hint: Can coverage already used this month still be reported as protection available for the future?)

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Does having no accident keep this month's expired protection available for the future?

Explanation: Assets decrease by $50 (Prepaid Insurance) and equity decreases by $50 through Insurance Expense. Cash and liabilities are unchanged. The cost reflects protection used, not a payment for an accident.

- opt_assets_down_eq_down: Assets decrease by $50 (-Prepaid Insurance); Equity decreases by $50 (-Insurance Expense); Cash is unaffected. (tag: ; hint: )
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does having no accident keep this month's expired protection available for the future?)
- opt_assets_down_cash: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

## repay_principal_final_payment v1

A company pays $50 cash today to its bank to repay all remaining principal on a previously recorded promissory note. The stated amount equals the entire principal balance immediately before payment. No new loan is received today. This payment includes no interest or fees; ignore interest. Record only this final principal repayment.

### identify_account

Which account records what the company owed the bank before today's principal payment?

Hint: What previously recorded balance is completely settled by today's final principal payment?

Explanation: Notes Payable records the remaining principal owed to the bank. Paying it off reduces the liability to zero. An expense would treat repayment of old borrowed money as a new cost.

- opt_expense_misconception: An expense for the loan payment (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What previously recorded balance is completely settled by today's final principal payment?)
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

Hint: After the bank receives the entire remaining principal, is any principal still owed on this note?

Explanation: Notes Payable decreases by its entire remaining principal balance and reaches zero. It does not become negative or increase; the company received no new borrowing.

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

Hint: What resource left the company to settle the note today?

Explanation: Cash decreases because the bank received the payment today, so Cash is credited. The payment settles a bank note rather than creating or settling an unrelated trade payable.

- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: What resource left the company to settle the note today?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What resource left the company to settle the note today?)

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

Hint: Only the remaining principal is paid, with no interest or fee. What obligation ends and what resource leaves?

Explanation: Debit Notes Payable $50 and Credit Cash $50. The note's principal becomes zero. An expense debit would incorrectly reduce income, while reversing the entry would record borrowing instead of final repayment.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_expense_cr_cash: Debit an expense for the loan payment $50 / Credit Cash $50 (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)

### equation_effect

How does repaying loan principal affect the accounting equation?

Hint: Does eliminating the remaining principal create a new cost or remove an existing obligation?

Explanation: Assets decrease by $50 (Cash) and liabilities decrease by $50 (Notes Payable). Equity is unchanged. This note is fully settled; no interest or fee affects income.

- opt_assets_down_liab_down: Assets decrease by $50 (-Cash); Liabilities decrease by $50 (-Notes Payable); Equity is unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Does eliminating the remaining principal create a new cost or remove an existing obligation?)

## dividend_cash_prior_earnings v1

A corporation's board declares and pays a $50 cash distribution to stockholders today from profits accumulated and recorded in earlier periods. No dividend was declared or recorded as owed earlier. The company performs no customer work and earns no new revenue today. This distribution is solely for share ownership, not payment for work or repayment of a loan. Record only today's declaration and payment.

### identify_account

Which account records cash distributions paid directly to stockholders?

Hint: Is today's event new earning or a distribution of amounts already earned in earlier periods?

Explanation: Dividends records the newly declared owner distribution. Earlier profits were already recorded; distributing them creates neither new revenue nor an expense. Cash records the separate payment.

- opt_dividends: Dividends (tag: ; hint: )
- opt_cash: Cash (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Is today's event new earning or a distribution of amounts already earned in earlier periods?)
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

Hint: Did the corporation distribute money today or leave the newly declared amount unpaid?

Explanation: Cash decreases because the distribution is paid today. Dividends Payable would describe an unpaid declaration, but none remains owed after this payment.

- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Did the corporation distribute money today or leave the newly declared amount unpaid?)
- opt_cash: Cash (tag: ; hint: )
- opt_dividends_payable: Dividends Payable (tag: wrong_account; hint: Did the corporation distribute money today or leave the newly declared amount unpaid?)

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

Hint: Earlier profits are already recorded and the distribution is declared and paid today. What owner distribution arose and what resource left?

Explanation: Debit Dividends $50 and Credit Cash $50. Do not record earlier revenue again or treat the distribution as an operating expense. Debiting a dividend payable would assume an earlier declaration that the scenario excludes.

- opt_dr_cash_cr_div: Debit Cash $50 / Credit Dividends $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the dividend $50 / Credit Cash $50 (tag: expense_recorded_on_dividend; hint: Were the stockholders paid for providing a service, or because they own shares?)
- opt_dr_div_cr_cash: Debit Dividends $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does paying cash dividends affect the accounting equation?

Hint: Does distributing earlier profits change today's net income or the owners' equity directly?

Explanation: Assets decrease by $50 (Cash) and equity decreases by $50 through Dividends. Liabilities and net income are unchanged. Earlier earnings stay recorded; this event is a distribution rather than new earning.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Dividends); Liabilities are unchanged. (tag: ; hint: )
- opt_assets_down_liab_down: Assets decrease by $50; Liabilities decrease by $50. (tag: equation_effect_missed; hint: Does distributing earlier profits change today's net income or the owners' equity directly?)
- opt_no_net_change: No change in equity; dividends are an asset. (tag: equation_effect_missed; hint: Does distributing earlier profits change today's net income or the owners' equity directly?)

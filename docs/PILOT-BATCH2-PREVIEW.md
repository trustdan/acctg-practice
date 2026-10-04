# Pilot batch 2: complete teaching preview

Reviewed version-1 snapshots at $50 with seed 100. Allowed amounts: $50/$100/$200. All seven stages inherit reviewed family teaching and supply scenario overrides at identity, counter-account, entry, and equation stages; the ownership-investment scenario also customizes direction. This artifact records options, tags and snapshotted mistake hints. Provenance and contrast partners: [PILOT-BATCH2-REVIEW.md](PILOT-BATCH2-REVIEW.md). Approved package: [PILOT-BATCH2.json](PILOT-BATCH2.json).

## prepaid_purchase_clinic v1

A clinic pays $50 cash today for a one-month insurance policy covering next month. None of this policy's coverage has begun or been used today. No amount for this policy was previously recorded as owed, and the payment buys coverage rather than medical equipment.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Has the clinic used any of the protection it paid for, or is all of it still available for next month?

Explanation: Prepaid Insurance records the clinic's unused future coverage. Insurance Expense would mean protection was already consumed. Cash records the payment, not the coverage that remains available.

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

Hint: Did the clinic pay for the policy today, or only agree to pay an insurer later?

Explanation: Cash decreased because the clinic paid today, so Cash is credited. A payable would report an amount still owed even though this policy has been paid for.

- opt_insurance_expense: Insurance Expense (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Did the clinic pay for the policy today, or only agree to pay an insurer later?)
- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: All the purchased protection remains for next month and the clinic paid today. What benefit was acquired, and what resource left?

Explanation: Debit Prepaid Insurance $50 and Credit Cash $50. Insurance Expense would recognize coverage not yet used. A payable would ignore today's cash payment.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)
- opt_dr_prepaid_cr_ap: Debit Prepaid Insurance $50 / Credit Accounts Payable $50 (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_dr_cash_cr_prepaid: Debit Cash $50 / Credit Prepaid Insurance $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_prepaid_cr_cash: Debit Prepaid Insurance $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Did the payment leave the clinic with protection still available, or use that protection up today?

Explanation: Prepaid Insurance increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged; cash was exchanged for unused coverage.

- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Did the payment leave the clinic with protection still available, or use that protection up today?)
- opt_asset_swap: Asset exchange: Prepaid Insurance increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: expense_recorded_on_prepaid_purchase; hint: Has the purchased coverage already been consumed, or is its protection still available for future months?)

## prepaid_consumption_clinic v1

At the end of this month, a clinic records the $50 cost of a one-month insurance policy whose entire coverage has now expired. It paid for and recorded the policy last month before coverage began. No cash is paid today, and no coverage from this policy remains for future months.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: What happened to the protection the clinic had paid for before this month began?

Explanation: Insurance Expense records the protection used this month. Prepaid Insurance records protection still available, but none from this policy remains. Cash was paid last month, not today.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_insurance_expense: Insurance Expense (tag: ; hint: )
- opt_rent_expense: Rent Expense (tag: wrong_account; hint: What happened to the protection the clinic had paid for before this month began?)

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

Hint: Which earlier balance held the protection that has now been entirely used up?

Explanation: Prepaid Insurance is credited to remove the used-up coverage. Crediting Cash would record a second payment, and recording no entry would leave expired protection reported as a resource.

- opt_prepaid_insurance: Prepaid Insurance (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Which earlier balance held the protection that has now been entirely used up?)
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Which earlier balance held the protection that has now been entirely used up?)
- opt_cash: Cash (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: The clinic used the entire policy this month after paying last month. What cost arose, and what previously recorded resource is gone?

Explanation: Debit Insurance Expense $50 and Credit Prepaid Insurance $50. The full policy is consumed; no cash changes today. Leaving the prepaid balance unchanged would overstate future coverage.

- opt_dr_expense_cr_cash: Debit Insurance Expense $50 / Credit Cash $50 (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)
- opt_dr_prepaid_cr_expense: Debit Prepaid Insurance $50 / Credit Insurance Expense $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_prepaid: Debit Insurance Expense $50 / Credit Prepaid Insurance $50 (tag: ; hint: )
- opt_no_entry: No entry; the full policy remains prepaid. (tag: prepaid_not_expensed_on_consumption; hint: Can coverage already used this month still be reported as protection available for the future?)

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Can expired protection still be a future resource even though no cash moved today?

Explanation: Assets decrease by $50 (Prepaid Insurance) and equity decreases by $50 through Insurance Expense. Cash and liabilities are unchanged; using a prepaid resource still incurs a cost.

- opt_assets_down_eq_down: Assets decrease by $50 (-Prepaid Insurance); Equity decreases by $50 (-Insurance Expense); Cash is unaffected. (tag: ; hint: )
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Can expired protection still be a future resource even though no cash moved today?)
- opt_assets_down_cash: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Insurance Expense). (tag: cash_recorded_on_prepaid_expiration; hint: Was the policy paid for today, or was today's used coverage paid for earlier?)

## cash_rent_clinic v1

A clinic pays $50 cash today for using its treatment rooms during the current month. This payment covers no future months, and none of this month's room cost was previously recorded as owed or paid in advance. The clinic does not buy any building or equipment.

### identify_account

Which account records the cost of using the premises this month?

Hint: Does this payment buy future occupancy or pay for the rooms used during the current month?

Explanation: Rent Expense records this month's use of the rooms. Prepaid Rent would mean future occupancy, but none is purchased. Accounts Payable would mean an unpaid amount, but the clinic paid today.

- opt_prepaid_rent: Prepaid Rent (tag: wrong_account; hint: Does this payment buy future occupancy or pay for the rooms used during the current month?)
- opt_rent_expense: Rent Expense (tag: ; hint: )
- opt_ap: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### account_category

What category of account is Rent Expense?

Hint: Does this account track a future resource, an obligation, or a cost incurred in operating the business?

Explanation: Rent Expense is an Expense: a cost of using the premises this month. It reduces equity through income; it is not an asset providing future use.

- opt_expense: Expense (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_asset: Asset (tag: ; hint: )

### direction

Does the expense account balance increase or decrease when recording this rent payment?

Hint: Does recognizing this month's rent add to the costs accumulated this period or remove a previously recorded cost?

Explanation: Rent Expense increases as this month's cost is recorded. Cash decreases, but that does not make the expense account decrease.

- opt_decrease: Decrease (tag: ; hint: )
- opt_increase: Increase (tag: ; hint: )

### debit_credit

How is an increase in an Expense (Rent Expense) recorded?

Hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

Explanation: Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.

- opt_debit: Debit (Left side) (tag: ; hint: )
- opt_credit: Credit (Right side) (tag: ; hint: )

### counter_account

Which other account changes when the company pays for this month's use of the premises?

Hint: Did the clinic hand over money now, or leave this month's room cost unpaid?

Explanation: Cash is credited because the clinic paid today. Accounts Payable would imply a remaining unpaid amount, but no amount for this payment remains owed.

- opt_service_rev: Service Revenue (tag: wrong_account; hint: Did the clinic hand over money now, or leave this month's room cost unpaid?)
- opt_ap: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_cash: Cash (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for paying cash rent?

Hint: Only this month's room use is paid for today. What current cost arose, and what resource left?

Explanation: Debit Rent Expense $50 and Credit Cash $50. A prepaid debit would defer a current-month cost into the future. A payable credit would ignore the money already paid.

- opt_dr_prepaid_cr_cash: Debit Prepaid Rent $50 / Credit Cash $50 (tag: wrong_account; hint: Only this month's room use is paid for today. What current cost arose, and what resource left?)
- opt_dr_cash_cr_rent: Debit Cash $50 / Credit Rent Expense $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_rent_cr_cash: Debit Rent Expense $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_rent_cr_ap: Debit Rent Expense $50 / Credit Accounts Payable $50 (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### equation_effect

How does paying cash rent affect the accounting equation?

Hint: Unlike unused insurance, does this room payment leave a benefit for later months?

Explanation: Assets decrease by $50 (Cash) and equity decreases by $50 through Rent Expense. Liabilities are unchanged. This payment covers current use, not a new future resource.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Rent Expense); Liabilities unchanged. (tag: ; hint: )
- opt_no_net_change: No change; expenses do not affect balance sheet. (tag: equation_effect_missed; hint: Unlike unused insurance, does this room payment leave a benefit for later months?)
- opt_assets_down_liab_down: Assets decrease by $50; Liabilities decrease by $50; Equity unchanged. (tag: equation_effect_missed; hint: Unlike unused insurance, does this room payment leave a benefit for later months?)

## issue_shares_studio v1

A design studio incorporated as a corporation receives $50 cash today from investors in exchange for newly issued common shares. The investors do not lend the money or pay for any client work. The studio plans to use the funds for employee pay next month, but spends none today.

### identify_account

What did the company receive from the investors today?

Hint: What did the investors hand over today, before any employee payment occurs?

Explanation: Cash records the investors' payment to the studio. The ownership contribution is recorded separately; a plan to pay employees later is not a payment or expense today.

- opt_service_rev: Service Revenue (tag: revenue_recorded_on_share_issue; hint: Set aside how the payment was earned or financed. What resource actually arrived today?)
- opt_cash: Cash (tag: ; hint: )
- opt_common_stock: Common Stock (tag: wrong_account; hint: What arrived today, and is the company still waiting to collect that money?)

### account_category

What category of account is Cash?

Hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

Explanation: Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.

- opt_asset: Asset (tag: ; hint: )
- opt_liability: Liability (tag: ; hint: )
- opt_equity: Equity (tag: ; hint: )

### direction

Does Cash increase or decrease upon issuing stock for cash?

Hint: Does receiving investor money today add cash or remove cash, even if the studio plans to spend it later?

Explanation: Cash increases because the investors paid today. The intended future payroll does not decrease cash now; no spending occurred today.

- opt_increase: Increase (tag: ; hint: )
- opt_decrease: Decrease (tag: ; hint: )

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- opt_credit: Credit (Right side) (tag: ; hint: )
- opt_debit: Debit (Left side) (tag: ; hint: )

### counter_account

Which account balances the receipt from investors who bought shares?

Hint: Did the investors buy shares, expect a loan repayment, or pay for completed client work?

Explanation: Common Stock records capital contributed for newly issued shares. Notes Payable would describe a loan, and Service Revenue would describe completed work. Neither fits this ownership investment.

- opt_common_stock: Common Stock (Equity) (tag: ; hint: )
- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recorded_on_share_issue; hint: Did the investors pay for a service, or buy an ownership stake?)
- opt_notes_payable: Notes Payable (Liability) (tag: wrong_account; hint: Did the investors buy shares, expect a loan repayment, or pay for completed client work?)

### balanced_entry

What is the complete balanced journal entry for issuing common stock for cash?

Hint: The studio received money for new shares and spent none today. What resource arrived, and what ownership contribution was made?

Explanation: Debit Cash $50 and Credit Common Stock $50. Crediting revenue would call investment earning; crediting a note would turn owners into lenders. Future payroll is a separate event.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recorded_on_share_issue; hint: Did the investors pay for a service, or buy an ownership stake?)
- opt_dr_stock_cr_cash: Debit Common Stock $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_stock: Debit Cash $50 / Credit Common Stock $50 (tag: ; hint: )

### equation_effect

How does issuing common shares affect the accounting equation?

Hint: Does the future payroll plan change whether today's investors acquired ownership or became lenders?

Explanation: Assets increase by $50 (Cash) and equity increases by $50 (Common Stock). Liabilities are unchanged. The increase in equity is contributed capital, not revenue.

- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Common Stock); Liabilities unchanged. (tag: ; hint: )
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50; Equity unchanged. (tag: equation_effect_missed; hint: Does the future payroll plan change whether today's investors acquired ownership or became lenders?)
- opt_no_net_change: No net change in total assets. (tag: equation_effect_missed; hint: Does the future payroll plan change whether today's investors acquired ownership or became lenders?)

## repay_principal_studio v1

A corporate design studio pays $50 cash today to its bank to settle part of the principal on a previously recorded loan. The payment is principal only; ignore interest. It is not a payment for employee work or a distribution to stockholders, and the bank provides no new loan today.

### identify_account

Which account records what the company owed the bank before today's principal payment?

Hint: What previously recorded amount owed to the bank is today's principal payment settling?

Explanation: Notes Payable records the principal the studio owes the bank. The payment reduces that debt; an expense would treat principal as a new cost, and Dividends would describe a distribution to owners rather than a lender.

- opt_expense_misconception: An expense for the loan payment (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What previously recorded amount owed to the bank is today's principal payment settling?)
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

Hint: What resource did the studio hand over to settle part of the loan today?

Explanation: Cash decreases because the studio paid the bank today, so it is credited. Accounts Payable would refer to another obligation, not the money actually handed over.

- opt_cash: Cash (tag: ; hint: )
- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: What resource did the studio hand over to settle part of the loan today?)
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: What resource did the studio hand over to settle part of the loan today?)

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

Hint: This payment settles old loan principal and includes no interest. What obligation shrinks, and what resource leaves?

Explanation: Debit Notes Payable $50 and Credit Cash $50. An expense would reduce income for a principal repayment; Dividends would incorrectly classify the lender as an owner receiving a distribution.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_expense_cr_cash: Debit an expense for the loan payment $50 / Credit Cash $50 (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)

### equation_effect

How does repaying loan principal affect the accounting equation?

Hint: Is the bank being paid a new operating cost, or receiving money it previously lent?

Explanation: Assets decrease by $50 (Cash) and liabilities decrease by $50 (Notes Payable). Equity is unchanged because only old principal is repaid, with no interest or other expense.

- opt_assets_down_liab_down: Assets decrease by $50 (-Cash); Liabilities decrease by $50 (-Notes Payable); Equity is unchanged. (tag: ; hint: )
- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_loan_repayment; hint: Is this payment a new cost, or repayment of principal the company previously borrowed?)
- opt_no_net_change: No net change in total assets; asset swap only. (tag: equation_effect_missed; hint: Is the bank being paid a new operating cost, or receiving money it previously lent?)

## dividend_cash_studio v1

A corporate design studio's board declares and pays a $50 cash distribution to its stockholders today because they own shares. No distribution was declared earlier, and none was previously recorded as owed. The owners are not being paid for work, and this payment does not settle a bank loan.

### identify_account

Which account records cash distributions paid directly to stockholders?

Hint: Are the recipients being paid because they provided work, lent money, or own shares?

Explanation: Dividends records the studio's newly declared distribution to owners. An expense would imply payment for a business cost; Notes Payable would imply repayment to a lender. Neither happened.

- opt_dividends: Dividends (tag: ; hint: )
- opt_cash: Cash (tag: ; hint: )
- opt_service_revenue: Service Revenue (tag: wrong_account; hint: Are the recipients being paid because they provided work, lent money, or own shares?)
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

Hint: Was the distribution only announced for later payment, or did money leave the studio today?

Explanation: Cash is credited because the studio paid today. Dividends Payable would be used if the newly declared distribution remained unpaid, but declaration and payment both happened today.

- opt_accounts_payable: Accounts Payable (tag: wrong_account; hint: Was the distribution only announced for later payment, or did money leave the studio today?)
- opt_cash: Cash (tag: ; hint: )
- opt_dividends_payable: Dividends Payable (tag: wrong_account; hint: Was the distribution only announced for later payment, or did money leave the studio today?)

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

Hint: The studio declared and paid this distribution today without an earlier payable. What owner distribution arose, and what resource left?

Explanation: Debit Dividends $50 and Credit Cash $50. This is not an expense or a loan repayment. Debiting Dividends Payable would settle an earlier declaration that the scenario rules out.

- opt_dr_cash_cr_div: Debit Cash $50 / Credit Dividends $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the dividend $50 / Credit Cash $50 (tag: expense_recorded_on_dividend; hint: Were the stockholders paid for providing a service, or because they own shares?)
- opt_dr_div_cr_cash: Debit Dividends $50 / Credit Cash $50 (tag: ; hint: )

### equation_effect

How does paying cash dividends affect the accounting equation?

Hint: Does paying owners because they hold shares reduce operating income or their ownership claim directly?

Explanation: Assets decrease by $50 (Cash) and equity decreases by $50 through Dividends. Liabilities and net income are unchanged; this distribution is not a business expense.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Dividends); Liabilities are unchanged. (tag: ; hint: )
- opt_assets_down_liab_down: Assets decrease by $50; Liabilities decrease by $50. (tag: equation_effect_missed; hint: Does paying owners because they hold shares reduce operating income or their ownership claim directly?)
- opt_no_net_change: No change in equity; dividends are an asset. (tag: equation_effect_missed; hint: Does paying owners because they hold shares reduce operating income or their ownership claim directly?)


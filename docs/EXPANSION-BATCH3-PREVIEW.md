# Expansion batch 3: complete teaching preview

Version-1 snapshots at $50, seed 100. Every allowed amount and scaffold is regression-tested. Shared family teaching is included below alongside custom teaching, all options, tags, and stored mistake hints. Review: [EXPANSION-BATCH3-REVIEW.md](EXPANSION-BATCH3-REVIEW.md).

## equipment_purchase_workshop v1

A repair workshop pays $50 cash today to buy and take ownership of a machine expected to serve it for several years. The machine is delivered today, but production will begin next month. The full purchase price is paid now; no amount was previously recorded as owed. Record only this purchase, with no depreciation or other costs.

### identify_account

Which account records the newly acquired machinery?

Hint: Does waiting until next month to use the delivered machine change whether the workshop owns it today?

Explanation: Equipment records the machine now owned for future use. Waiting to begin production does not postpone this purchase. An expense would treat the entire multi-year benefit as already used; Cash records the payment.

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

Hint: Did money leave today, even though production will start later?

Explanation: Cash decreased when the workshop paid today, so Cash is credited. Neither Notes Payable nor Accounts Payable arises from this fully paid purchase; no shares were issued.

- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: Did money leave today, even though production will start later?)
- opt_common_stock: Common Stock (tag: wrong_account; hint: Did money leave today, even though production will start later?)
- opt_cash: Cash (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

Hint: The workshop owns the delivered machine and has paid for it. What resource arrived, and what resource left?

Explanation: Debit Equipment $50 and Credit Cash $50. The machine is acquired now, before production starts. An immediate expense would consume its full future benefit; reversing the entry would describe giving up equipment for cash.

- opt_dr_equip_cr_cash: Debit Equipment $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_cash_cr_equip: Debit Cash $50 / Credit Equipment $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the machinery $50 / Credit Cash $50 (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)

### equation_effect

How does purchasing equipment for cash affect the accounting equation?

Hint: Is cash the only resource the workshop has after buying the delivered machine?

Explanation: Equipment increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged. Starting production later does not remove the machine from today's resources.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Is cash the only resource the workshop has after buying the delivered machine?)
- opt_asset_swap: Asset exchange: Equipment increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )

## equipment_purchase_replacement v1

A repair workshop pays $50 cash today to buy a replacement machine from a dealer and takes delivery and ownership immediately. The new machine will be used for several years. This is a purchase of a whole machine, not a repair service. The old machine is kept; no sale, trade-in, or disposal occurs today. No amount for the new machine was previously recorded as owed. Record only the new purchase, with no depreciation.

### identify_account

Which account records the newly acquired machinery?

Hint: Did the workshop acquire another whole machine or merely pay to fix one it already owned?

Explanation: Equipment records the newly acquired machine. Calling it a replacement does not make buying a whole multi-year resource a repair cost. The old machine is kept, so this event records no disposal.

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

Hint: With the old machine kept and a new one acquired, is more or less machinery recorded?

Explanation: Equipment increases by the cost of the new machine. Its purpose as a replacement does not remove the old machine in this event. Cash decreases in a separate account.

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

Hint: How did the workshop settle the dealer's price today?

Explanation: Cash decreases because the dealer was paid in full today. A payable would report a remaining debt, while Common Stock would imply shares issued to the dealer; neither happened.

- opt_accounts_payable: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_notes_payable: Notes Payable (tag: wrong_account; hint: How did the workshop settle the dealer's price today?)
- opt_common_stock: Common Stock (tag: wrong_account; hint: How did the workshop settle the dealer's price today?)
- opt_cash: Cash (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

Hint: Only the new machine purchase is recorded, with no old-machine disposal. What resource was acquired and how was it paid for?

Explanation: Debit Equipment $50 and Credit Cash $50. Buying the whole machine creates a multi-year resource rather than an immediate repair expense. Crediting Equipment would remove machinery even though none was disposed of.

- opt_dr_equip_cr_cash: Debit Equipment $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_cash_cr_equip: Debit Cash $50 / Credit Equipment $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_expense_cr_cash: Debit an expense for the machinery $50 / Credit Cash $50 (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)

### equation_effect

How does purchasing equipment for cash affect the accounting equation?

Hint: Does the label replacement turn the new machine's entire future usefulness into a cost consumed today?

Explanation: Equipment increases by $50 and Cash decreases by $50. Total assets, liabilities, and equity are unchanged. No disposal or depreciation is part of this event.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Expense). (tag: expense_recorded_on_equipment_purchase; hint: Was the machinery's entire benefit used today, or can it provide use over future years?)
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50. (tag: equation_effect_missed; hint: Does the label replacement turn the new machine's entire future usefulness into a cost consumed today?)
- opt_asset_swap: Asset exchange: Equipment increases (+$50) and Cash decreases (-$50); Total Assets and Equity unchanged. (tag: ; hint: )

## cash_rent_workshop v1

A repair workshop pays $50 cash today for using rented workspace during the current month. It acquires no ownership of the building or any machinery, and this payment covers no future months. No amount for this month was previously recorded as owed or prepaid.

### identify_account

Which account records the cost of using the premises this month?

Hint: Did the payment acquire an owned resource for future years or pay for this month's use of space?

Explanation: Rent Expense records this month's workspace use. No equipment or building was purchased. Prepaid Rent would mean future occupancy, but this payment covers only the current month.

- opt_prepaid_rent: Prepaid Rent (tag: wrong_account; hint: Did the payment acquire an owned resource for future years or pay for this month's use of space?)
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

Hint: Was this month's workspace paid for today or left unpaid?

Explanation: Cash is credited because the workshop paid today. Accounts Payable would describe an unpaid obligation, which this event rules out.

- opt_service_rev: Service Revenue (tag: wrong_account; hint: Was this month's workspace paid for today or left unpaid?)
- opt_ap: Accounts Payable (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)
- opt_cash: Cash (tag: ; hint: )

### balanced_entry

What is the complete balanced journal entry for paying cash rent?

Hint: The space was used this month and paid for now. What current cost arose and what resource left?

Explanation: Debit Rent Expense $50 and Credit Cash $50. Unlike buying a machine, this payment acquires no owned multi-year resource. A prepaid debit would defer a cost for current use.

- opt_dr_prepaid_cr_cash: Debit Prepaid Rent $50 / Credit Cash $50 (tag: wrong_account; hint: The space was used this month and paid for now. What current cost arose and what resource left?)
- opt_dr_cash_cr_rent: Debit Cash $50 / Credit Rent Expense $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_rent_cr_cash: Debit Rent Expense $50 / Credit Cash $50 (tag: ; hint: )
- opt_dr_rent_cr_ap: Debit Rent Expense $50 / Credit Accounts Payable $50 (tag: payable_recorded_for_cash_payment; hint: Was the purchase paid for today, or is there still an unpaid amount owed?)

### equation_effect

How does paying cash rent affect the accounting equation?

Hint: After paying for current workspace use, does a purchased resource remain for future periods?

Explanation: Assets decrease by $50 (Cash) and equity decreases by $50 through Rent Expense. Liabilities are unchanged. This is current use of rented space, rather than an exchange of cash for owned machinery.

- opt_assets_down_eq_down: Assets decrease by $50 (-Cash); Equity decreases by $50 (-Rent Expense); Liabilities unchanged. (tag: ; hint: )
- opt_no_net_change: No change; expenses do not affect balance sheet. (tag: equation_effect_missed; hint: After paying for current workspace use, does a purchased resource remain for future periods?)
- opt_assets_down_liab_down: Assets decrease by $50; Liabilities decrease by $50; Equity unchanged. (tag: equation_effect_missed; hint: After paying for current workspace use, does a purchased resource remain for future periods?)

## borrow_cash_future_equipment v1

An incorporated repair workshop receives $50 cash today from a bank after signing a promissory note requiring repayment. It plans to buy a machine next month, but no machine has been ordered, delivered, or paid for today. The bank receives no shares, and this receipt is not payment for customer work. Record only today's borrowing; ignore interest.

### identify_account

What did the company receive today, and which account records it?

Hint: What did the workshop actually receive today before any machine purchase occurred?

Explanation: Cash records the money received from the bank. A plan to buy machinery does not create Equipment today. No customer service was provided in exchange for this loan.

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

Hint: Must the workshop repay this money, or did the bank purchase ownership?

Explanation: Notes Payable records the obligation under the promissory note. Common Stock would mean an ownership contribution, and Service Revenue would mean earning from work. Neither describes the bank loan.

- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_notes_payable: Notes Payable (Liability) (tag: ; hint: )
- opt_common_stock: Common Stock (Equity) (tag: wrong_account; hint: Must the workshop repay this money, or did the bank purchase ownership?)

### balanced_entry

What is the complete balanced journal entry for borrowing cash on a note?

Hint: Money arrived under a repayment promise, while the machine purchase is still only a plan. What resource and obligation arose today?

Explanation: Debit Cash $50 and Credit Notes Payable $50. Do not record Equipment before a purchase occurs. Crediting revenue would treat repayable financing as earning, and crediting stock would turn the lender into an owner.

- opt_dr_cash_cr_notes: Debit Cash $50 / Credit Notes Payable $50 (tag: ; hint: )
- opt_dr_notes_cr_cash: Debit Notes Payable $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)

### equation_effect

How does borrowing cash affect the accounting equation?

Hint: Does intending to buy equipment remove the obligation to repay the bank?

Explanation: Assets increase by $50 (Cash) and liabilities increase by $50 (Notes Payable). Equity is unchanged. The later equipment purchase is a separate event, so no equipment or expense is recorded today.

- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Revenue); Liabilities unchanged. (tag: revenue_recorded_on_borrowing; hint: Did the lender pay for completed work or supply money the company must repay?)
- opt_no_net_change: No net change in total assets. (tag: equation_effect_missed; hint: Does intending to buy equipment remove the obligation to repay the bank?)
- opt_assets_up_liab_up: Assets increase by $50 (+Cash); Liabilities increase by $50 (+Notes Payable); Equity is unchanged. (tag: ; hint: )

## issue_shares_future_equipment v1

An incorporated repair workshop receives $50 cash today from investors in exchange for newly issued common shares. It plans to buy a machine next month, but no machine has been ordered, delivered, or paid for today. The investors are not lenders and are not paying for customer work. Record only today's share issuance.

### identify_account

What did the company receive from the investors today?

Hint: What resource arrived today while the machine purchase remains a future plan?

Explanation: Cash records the investor money received today. Equipment has not been purchased. The investors paid for ownership rather than completed work, so the receipt is not revenue.

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

Hint: Does receiving the investment add cash today even though it may be spent next month?

Explanation: Cash increases when investors pay today. Planning a later machine purchase does not decrease cash or create Equipment now.

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

Hint: Did these investors acquire shares or a right to repayment under a loan?

Explanation: Common Stock records the contribution for newly issued shares. Notes Payable would report a loan obligation, but the investors bought ownership. Service Revenue would describe earned customer work, which did not occur.

- opt_common_stock: Common Stock (Equity) (tag: ; hint: )
- opt_service_rev: Service Revenue (Revenue) (tag: revenue_recorded_on_share_issue; hint: Did the investors pay for a service, or buy an ownership stake?)
- opt_notes_payable: Notes Payable (Liability) (tag: wrong_account; hint: Did these investors acquire shares or a right to repayment under a loan?)

### balanced_entry

What is the complete balanced journal entry for issuing common stock for cash?

Hint: Only the ownership investment happens today. What resource arrived and what capital contribution arose?

Explanation: Debit Cash $50 and Credit Common Stock $50. This is contributed capital rather than revenue or borrowing. The planned equipment purchase must be recorded separately when it occurs.

- opt_dr_cash_cr_rev: Debit Cash $50 / Credit Service Revenue $50 (tag: revenue_recorded_on_share_issue; hint: Did the investors pay for a service, or buy an ownership stake?)
- opt_dr_stock_cr_cash: Debit Common Stock $50 / Credit Cash $50 (tag: reversed_sides; hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?)
- opt_dr_cash_cr_stock: Debit Cash $50 / Credit Common Stock $50 (tag: ; hint: )

### equation_effect

How does issuing common shares affect the accounting equation?

Hint: Does buying ownership create a loan obligation just because the money may fund a future machine?

Explanation: Assets increase by $50 (Cash) and equity increases by $50 (Common Stock). Liabilities are unchanged. Equity rises through contributed capital, not income; no equipment purchase has occurred.

- opt_assets_up_eq_up: Assets increase by $50 (+Cash); Equity increases by $50 (+Common Stock); Liabilities unchanged. (tag: ; hint: )
- opt_assets_up_liab_up: Assets increase by $50; Liabilities increase by $50; Equity unchanged. (tag: equation_effect_missed; hint: Does buying ownership create a loan obligation just because the money may fund a future machine?)
- opt_no_net_change: No net change in total assets. (tag: equation_effect_missed; hint: Does buying ownership create a loan obligation just because the money may fund a future machine?)

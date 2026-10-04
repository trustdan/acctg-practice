# Expansion batch 11 rendered review

Draft preview at $100, seed 101. All seven stages include custom/inherited teaching and stored option hints.

## cash_rent_separate_insurance

A company pays $100 cash today for workspace used during the current month. It also has a separately paid and recorded insurance policy for next month; none of today's payment buys or consumes that insurance. This rent covers no future months and was not previously recorded as owed or prepaid. No building or equipment is purchased. Record only today's current-month rent payment.

### identify_account

Which account records the cost of using the premises this month?

Hint: What benefit does this payment cover: occupancy already used or the separate future policy?

Explanation: Rent Expense records this month's used workspace. The separately recorded future insurance remains an asset and is not part of this payment.

- a: Rent Expense (id=opt_rent_expense; correct=true; tag=)
  Stored hint: What benefit does this payment cover: occupancy already used or the separate future policy?
- b: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- c: Prepaid Rent (id=opt_prepaid_rent; correct=false; tag=wrong_account)
  Stored hint: What benefit does this payment cover: occupancy already used or the separate future policy?

### account_category

What category of account is Rent Expense?

Hint: Does this account track a future resource, an obligation, or a cost incurred in operating the business?

Explanation: Rent Expense is an Expense: a cost of using the premises this month. It reduces equity through income; it is not an asset providing future use.

- a: Asset (id=opt_asset; correct=false; tag=)
  Stored hint: Does this account track a future resource, an obligation, or a cost incurred in operating the business?
- b: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Does this account track a future resource, an obligation, or a cost incurred in operating the business?
- c: Expense (id=opt_expense; correct=true; tag=)
  Stored hint: Does this account track a future resource, an obligation, or a cost incurred in operating the business?

### direction

Does the expense account balance increase or decrease when recording this rent payment?

Hint: Does recognizing this month's rent add to the costs accumulated this period or remove a previously recorded cost?

Explanation: Rent Expense increases as this month's cost is recorded. Cash decreases, but that does not make the expense account decrease.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does recognizing this month's rent add to the costs accumulated this period or remove a previously recorded cost?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does recognizing this month's rent add to the costs accumulated this period or remove a previously recorded cost?

### debit_credit

How is an increase in an Expense (Rent Expense) recorded?

Hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

Explanation: Expenses have a normal debit balance, so an increase is a Debit (left side). Expenses reduce equity, but the expense account itself increases; credit would decrease it.

- a: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an expense's normal balance?
- b: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an expense's normal balance?

### counter_account

Which other account changes when the company pays for this month's use of the premises?

Hint: Does the current rent payment draw on cash or use up the separate insurance coverage?

Explanation: Cash decreases for rent paid today. Prepaid Insurance stays unchanged; it is a different benefit.

- a: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Does the current rent payment draw on cash or use up the separate insurance coverage?
- c: Service Revenue (id=opt_service_rev; correct=false; tag=wrong_account)
  Stored hint: Does the current rent payment draw on cash or use up the separate insurance coverage?

### balanced_entry

What is the complete balanced journal entry for paying cash rent?

Hint: Which current cost and payment occur while the future policy remains unused?

Explanation: Debit Rent Expense $100 and Credit Cash $100. Do not consume or purchase the separate recorded insurance again.

- a: Debit Rent Expense $100 / Credit Accounts Payable $100 (id=opt_dr_rent_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Rent Expense $100 / Credit Cash $100 (id=opt_dr_rent_cr_cash; correct=true; tag=)
  Stored hint: Which current cost and payment occur while the future policy remains unused?
- c: Debit Prepaid Rent $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=false; tag=wrong_account)
  Stored hint: Which current cost and payment occur while the future policy remains unused?
- d: Debit Cash $100 / Credit Rent Expense $100 (id=opt_dr_cash_cr_rent; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does paying cash rent affect the accounting equation?

Hint: Does an unrelated unused policy offset the cost of current occupancy?

Explanation: Assets decrease through Cash and equity decreases through Rent Expense by $100. Liabilities and the separate insurance asset are unchanged.

- a: Assets decrease by $100; Liabilities decrease by $100; Equity unchanged. (id=opt_assets_down_liab_down; correct=false; tag=equation_effect_missed)
  Stored hint: Does an unrelated unused policy offset the cost of current occupancy?
- b: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Rent Expense); Liabilities unchanged. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does an unrelated unused policy offset the cost of current occupancy?
- c: No change; expenses do not affect balance sheet. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does an unrelated unused policy offset the cost of current occupancy?

## borrow_cash_collateral

A company receives $100 cash today from a bank and signs a promissory note requiring repayment. A machine already owned and recorded by the company is pledged as collateral, but remains owned, held, and used by the company. No machine is sold, delivered to the bank, or purchased today, and no shares or customer services are involved. Record only today's new loan receipt; ignore interest and assume no fees.

### identify_account

What did the company receive today, and which account records it?

Hint: What new resource arrives while the already owned machine stays in use?

Explanation: Cash records the company's receipt of the bank loan. Pledging an existing machine does not acquire or dispose of it in this event.

- a: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What new resource arrives while the already owned machine stays in use?

### account_category

What category of account is Cash?

Hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

Explanation: Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.

- a: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?
- b: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?
- c: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

### direction

Does Cash increase or decrease upon borrowing money?

Hint: Compare the money held before and after today's receipt. Is there more or less?

Explanation: Cash increases because money arrived today. A future repayment obligation does not cancel today's receipt.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Compare the money held before and after today's receipt. Is there more or less?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Compare the money held before and after today's receipt. Is there more or less?

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

### counter_account

Which account balances the receipt from the lender?

Hint: Must the company repay the bank even though it offered security?

Explanation: Notes Payable records the repayment obligation. Collateral does not turn the bank into an owner or customer.

- a: Notes Payable (Liability) (id=opt_notes_payable; correct=true; tag=)
  Stored hint: Must the company repay the bank even though it offered security?
- b: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: Must the company repay the bank even though it offered security?
- c: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?

### balanced_entry

What is the complete balanced journal entry for borrowing cash on a note?

Hint: What money arrives and repayment promise arises without any machine changing ownership?

Explanation: Debit Cash $100 and Credit Notes Payable $100. Leave the already recorded Equipment balance unchanged; collateral is not a sale or new purchase.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?
- b: Debit Notes Payable $100 / Credit Cash $100 (id=opt_dr_notes_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Notes Payable $100 (id=opt_dr_cash_cr_notes; correct=true; tag=)
  Stored hint: What money arrives and repayment promise arises without any machine changing ownership?

### equation_effect

How does borrowing cash affect the accounting equation?

Hint: Does the company still own its machine after securing this new debt?

Explanation: Assets increase through Cash and liabilities increase through Notes Payable by $100. Equity and existing Equipment are unchanged.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Notes Payable); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: Does the company still own its machine after securing this new debt?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?
- c: No net change in total assets. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does the company still own its machine after securing this new debt?

## issue_shares_existing_lender

A corporation receives $100 cash today from an investor in exchange for newly issued common shares. The investor also holds a separate loan to the corporation, already recorded last month, but that debt is neither repaid, converted to shares, nor increased today. Today's cash buys ownership only and requires no repayment. No customer work is involved. Record only the new share issuance with the full receipt as share capital; ignore par-value splits and fees.

### identify_account

What did the company receive from the investors today?

Hint: What resource arrives for the newly issued shares today?

Explanation: Cash records the new receipt. The investor's separate previously recorded loan is not today's payment.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrives for the newly issued shares today?
- b: Common Stock (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- c: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?

### account_category

What category of account is Cash?

Hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

Explanation: Cash is an Asset because the company owns and can use it. Cash is not Revenue: earning, borrowing, and owner investment can all bring in cash.

- a: Equity (id=opt_equity; correct=false; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?
- b: Asset (id=opt_asset; correct=true; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?
- c: Liability (id=opt_liability; correct=false; tag=)
  Stored hint: Is cash a resource the company owns, an amount it owes, or an ownership claim?

### direction

Does Cash increase or decrease upon issuing stock for cash?

Hint: Does today's share payment add money even while the older loan stays separately owed?

Explanation: Cash increases by today's receipt for new shares. The earlier loan remains owed, but today's ownership payment adds no repayment obligation.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does today's share payment add money even while the older loan stays separately owed?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does today's share payment add money even while the older loan stays separately owed?

### debit_credit

How is an increase to an Asset (Cash) recorded?

Hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

Explanation: Assets have a normal debit balance, so an increase is a Debit (left side). A credit would decrease the asset; debit does not mean money came in.

- a: Debit (Left side) (id=opt_debit; correct=true; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?
- b: Credit (Right side) (id=opt_credit; correct=false; tag=)
  Stored hint: An account increases on its normal-balance side. Which side is an asset's normal balance?

### counter_account

Which account balances the receipt from investors who bought shares?

Hint: Does this new payment buy ownership or add to the investor's separate loan?

Explanation: Common Stock records the new ownership capital. Notes Payable stays unchanged because the earlier loan is neither settled nor increased.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- b: Notes Payable (Liability) (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Does this new payment buy ownership or add to the investor's separate loan?
- c: Common Stock (Equity) (id=opt_common_stock; correct=true; tag=)
  Stored hint: Does this new payment buy ownership or add to the investor's separate loan?

### balanced_entry

What is the complete balanced journal entry for issuing common stock for cash?

Hint: What new capital arrives while the investor's older debt claim remains intact?

Explanation: Debit Cash $100 and Credit Common Stock $100. Do not convert, repay, or repeat the separate recorded loan.

- a: Debit Cash $100 / Credit Common Stock $100 (id=opt_dr_cash_cr_stock; correct=true; tag=)
  Stored hint: What new capital arrives while the investor's older debt claim remains intact?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- c: Debit Common Stock $100 / Credit Cash $100 (id=opt_dr_stock_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does issuing common shares affect the accounting equation?

Hint: Does an investor's older lender role make this ownership payment repayable?

Explanation: Assets and equity increase by $100; liabilities are unchanged. Existing debt remains owed separately from the newly issued shares.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Common Stock); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Does an investor's older lender role make this ownership payment repayable?
- b: No net change in total assets. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does an investor's older lender role make this ownership payment repayable?
- c: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Does an investor's older lender role make this ownership payment repayable?

## prepaid_purchase_budget_label

A company pays $100 cash today for an insurance policy covering next month only. Its internal cash budget labels the payment insurance expense, but none of this policy's coverage has begun or been used today. No amount was previously recorded as owed or paid for this policy. The insurer is paid in full. Record the purchase based on the unused future protection, not the budget label.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Has the policy supplied protection yet, regardless of the cash-budget label?

Explanation: Prepaid Insurance records the unused future protection. A spending label in a cash budget does not establish that a benefit was consumed.

- a: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Has the policy supplied protection yet, regardless of the cash-budget label?
- d: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Has the policy supplied protection yet, regardless of the cash-budget label?

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

Hint: What resource left when the insurer was paid in full?

Explanation: Cash decreases by the premium paid today. There is no remaining payable and no consumption entry yet.

- a: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: What resource left when the insurer was paid in full?
- c: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- d: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource left when the insurer was paid in full?

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: Does paying for next month's protection use it up today?

Explanation: Debit Prepaid Insurance $100 and Credit Cash $100. Calling the budget line an expense does not change the unused benefit into current Insurance Expense.

- a: Debit Prepaid Insurance $100 / Credit Accounts Payable $100 (id=opt_dr_prepaid_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Debit Prepaid Insurance $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=true; tag=)
  Stored hint: Does paying for next month's protection use it up today?
- d: Debit Cash $100 / Credit Prepaid Insurance $100 (id=opt_dr_cash_cr_prepaid; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Is the payment only a loss of resources, or does unused protection replace cash?

Explanation: Prepaid Insurance increases and Cash decreases by $100. Total assets, liabilities, and equity are unchanged. The budget label creates no separate posting.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Is the payment only a loss of resources, or does unused protection replace cash?
- c: Asset exchange: Prepaid Insurance increases (+$100) and Cash decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Is the payment only a loss of resources, or does unused protection replace cash?

## prepaid_consumption_recent_payment

At month-end, a company records $100 of insurance protection that expired during this month. It bought the policy at the start of this same month, before any coverage expired, and already recorded the full payment as a prepaid asset then. The stated amount is only this month's consumed portion; additional coverage remains for later months. No cash is paid today and this consumption has not yet been recorded. Record only month-end consumption, even though the purchase and consumption occur in the same month.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: What cost has been used by month-end after the earlier purchase was recorded?

Explanation: Insurance Expense records this month's expired protection. Purchase and consumption remain different events even within one month.

- a: Insurance Expense (id=opt_insurance_expense; correct=true; tag=)
  Stored hint: What cost has been used by month-end after the earlier purchase was recorded?
- b: Rent Expense (id=opt_rent_expense; correct=false; tag=wrong_account)
  Stored hint: What cost has been used by month-end after the earlier purchase was recorded?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=false; tag=)
  Stored hint: What cost has been used by month-end after the earlier purchase was recorded?
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

Hint: Which recorded benefit has been partly consumed without another payment?

Explanation: Prepaid Insurance decreases for the expired portion. Cash was already credited at purchase and future coverage remains an asset.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: Which recorded benefit has been partly consumed without another payment?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Which recorded benefit has been partly consumed without another payment?
- d: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Which recorded benefit has been partly consumed without another payment?

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: What part of the recorded policy expired after its earlier same-month purchase?

Explanation: Debit Insurance Expense $100 and Credit Prepaid Insurance $100. Do not repeat the cash payment or expense the still-unused remainder.

- a: Debit Insurance Expense $100 / Credit Prepaid Insurance $100 (id=opt_dr_expense_cr_prepaid; correct=true; tag=)
  Stored hint: What part of the recorded policy expired after its earlier same-month purchase?
- b: Debit Prepaid Insurance $100 / Credit Insurance Expense $100 (id=opt_dr_prepaid_cr_expense; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: No entry; the full policy remains prepaid. (id=opt_no_entry; correct=false; tag=prepaid_not_expensed_on_consumption)
  Stored hint: Can coverage already used this month still be reported as protection available for the future?
- d: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Does recognizing expired protection require paying the premium a second time?

Explanation: Assets decrease through Prepaid Insurance and equity decreases through Insurance Expense by $100. Cash and liabilities are unchanged; future coverage remains.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Assets decrease by $100 (-Prepaid Insurance); Equity decreases by $100 (-Insurance Expense); Cash is unaffected. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does recognizing expired protection require paying the premium a second time?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does recognizing expired protection require paying the premium a second time?

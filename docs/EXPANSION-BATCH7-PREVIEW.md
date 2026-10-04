# Expansion batch 7 rendered review

Draft preview at $100, seed 101. All seven stages include inherited teaching and stored option hints.

## cash_rent_before_opening

A company pays $100 cash today for rented workspace used during the current month to prepare for opening. It has served no customers yet, but this month's occupancy has been used. The payment covers no future months and buys no building or equipment. No amount was previously recorded as owed or prepaid. Record only this current-month rent payment.

### identify_account

Which account records the cost of using the premises this month?

Hint: Was this month's use of space consumed even though customer work has not started?

Explanation: Rent Expense records current workspace use. Having no revenue yet does not preserve occupancy already used as a future asset. This event buys neither equipment nor future occupancy.

- a: Rent Expense (id=opt_rent_expense; correct=true; tag=)
  Stored hint: Was this month's use of space consumed even though customer work has not started?
- b: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- c: Prepaid Rent (id=opt_prepaid_rent; correct=false; tag=wrong_account)
  Stored hint: Was this month's use of space consumed even though customer work has not started?

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

Hint: Was the landlord paid today or left waiting for payment?

Explanation: Cash decreases because the landlord was paid today. No unpaid obligation remains and no prepaid balance is being used.

- a: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Was the landlord paid today or left waiting for payment?
- c: Service Revenue (id=opt_service_rev; correct=false; tag=wrong_account)
  Stored hint: Was the landlord paid today or left waiting for payment?

### balanced_entry

What is the complete balanced journal entry for paying cash rent?

Hint: Space was used this month and paid for today. What cost arose before the company opened?

Explanation: Debit Rent Expense $100 and Credit Cash $100. Current occupancy is an expense even before the first customer sale. An asset debit would retain a benefit already consumed.

- a: Debit Rent Expense $100 / Credit Accounts Payable $100 (id=opt_dr_rent_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Rent Expense $100 / Credit Cash $100 (id=opt_dr_rent_cr_cash; correct=true; tag=)
  Stored hint: Space was used this month and paid for today. What cost arose before the company opened?
- c: Debit Prepaid Rent $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=false; tag=wrong_account)
  Stored hint: Space was used this month and paid for today. What cost arose before the company opened?
- d: Debit Cash $100 / Credit Rent Expense $100 (id=opt_dr_cash_cr_rent; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does paying cash rent affect the accounting equation?

Hint: Does a lack of customer revenue prevent current occupancy from reducing the owners' equity?

Explanation: Assets decrease by $100 (Cash) and equity decreases by $100 through Rent Expense. Liabilities are unchanged. The current cost affects income even with no revenue yet.

- a: Assets decrease by $100; Liabilities decrease by $100; Equity unchanged. (id=opt_assets_down_liab_down; correct=false; tag=equation_effect_missed)
  Stored hint: Does a lack of customer revenue prevent current occupancy from reducing the owners' equity?
- b: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Rent Expense); Liabilities unchanged. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does a lack of customer revenue prevent current occupancy from reducing the owners' equity?
- c: No change; expenses do not affect balance sheet. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does a lack of customer revenue prevent current occupancy from reducing the owners' equity?

## cash_rent_stockholder_landlord

A corporation pays $100 cash today to a landlord who also owns shares in the corporation. The payment is solely the agreed rent for workspace used during the current month, not a distribution for share ownership or a loan repayment. It covers no future months and acquires no building. No rent was previously recorded as owed or prepaid. Record only this rent payment from the corporation's perspective.

### identify_account

Which account records the cost of using the premises this month?

Hint: Was this payment made for using space or solely because the recipient owns shares?

Explanation: Rent Expense records current workspace use. The landlord's ownership of shares does not turn payment for occupancy into Dividends. Classify the economic purpose, not just the recipient.

- a: Rent Expense (id=opt_rent_expense; correct=true; tag=)
  Stored hint: Was this payment made for using space or solely because the recipient owns shares?
- b: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- c: Prepaid Rent (id=opt_prepaid_rent; correct=false; tag=wrong_account)
  Stored hint: Was this payment made for using space or solely because the recipient owns shares?

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

Hint: What company resource was transferred to pay the rent today?

Explanation: Cash decreases and is credited. The landlord received payment now; no payable remains and no loan principal was settled.

- a: Accounts Payable (id=opt_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What company resource was transferred to pay the rent today?
- c: Service Revenue (id=opt_service_rev; correct=false; tag=wrong_account)
  Stored hint: What company resource was transferred to pay the rent today?

### balanced_entry

What is the complete balanced journal entry for paying cash rent?

Hint: The shareholder is being paid as landlord. What current cost and payment should the corporation record?

Explanation: Debit Rent Expense $100 and Credit Cash $100. Dividends would describe a distribution for ownership rather than payment for workspace use. A note debit would incorrectly settle a loan.

- a: Debit Rent Expense $100 / Credit Accounts Payable $100 (id=opt_dr_rent_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Rent Expense $100 / Credit Cash $100 (id=opt_dr_rent_cr_cash; correct=true; tag=)
  Stored hint: The shareholder is being paid as landlord. What current cost and payment should the corporation record?
- c: Debit Prepaid Rent $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=false; tag=wrong_account)
  Stored hint: The shareholder is being paid as landlord. What current cost and payment should the corporation record?
- d: Debit Cash $100 / Credit Rent Expense $100 (id=opt_dr_cash_cr_rent; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does paying cash rent affect the accounting equation?

Hint: Does paying a shareholder for current occupancy affect income or only record an owner distribution?

Explanation: Assets decrease by $100 (Cash) and equity decreases through Rent Expense. Liabilities are unchanged. This expense reduces income; a dividend has the same total equation decrease but does not reduce income.

- a: Assets decrease by $100; Liabilities decrease by $100; Equity unchanged. (id=opt_assets_down_liab_down; correct=false; tag=equation_effect_missed)
  Stored hint: Does paying a shareholder for current occupancy affect income or only record an owner distribution?
- b: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Rent Expense); Liabilities unchanged. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does paying a shareholder for current occupancy affect income or only record an owner distribution?
- c: No change; expenses do not affect balance sheet. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does paying a shareholder for current occupancy affect income or only record an owner distribution?

## issue_shares_customer_investor

A corporation receives $100 cash today from an existing customer in exchange for newly issued common shares. This payment buys ownership only; it is not payment for completed or future services, and it settles no customer invoice. No repayment is promised. The full receipt is recorded as share capital; ignore par-value splits and fees. Record only today's share issuance.

### identify_account

What did the company receive from the investors today?

Hint: What resource did the corporation receive from the customer acting as an investor?

Explanation: Cash records the money received today. Being an existing customer does not make every receipt payment for work or collection of an invoice.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource did the corporation receive from the customer acting as an investor?
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

Hint: Did the corporation receive additional money or pay any out in this issuance?

Explanation: Cash increases by the investor receipt. No customer invoice is settled and no company payment occurs.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Did the corporation receive additional money or pay any out in this issuance?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Did the corporation receive additional money or pay any out in this issuance?

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

Hint: What did the payer receive in exchange: services, repayment rights, or ownership?

Explanation: Common Stock records the newly issued ownership shares. Service Revenue would imply earned work; Unearned Revenue would imply future work owed; Notes Payable would imply a repayment obligation. None describes this ownership purchase.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- b: Notes Payable (Liability) (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: What did the payer receive in exchange: services, repayment rights, or ownership?
- c: Common Stock (Equity) (id=opt_common_stock; correct=true; tag=)
  Stored hint: What did the payer receive in exchange: services, repayment rights, or ownership?

### balanced_entry

What is the complete balanced journal entry for issuing common stock for cash?

Hint: The customer bought new ownership shares rather than services. What resource and capital contribution arose?

Explanation: Debit Cash $100 and Credit Common Stock $100. Customer identity does not justify revenue or an advance liability. No receivable was settled and no loan was received.

- a: Debit Cash $100 / Credit Common Stock $100 (id=opt_dr_cash_cr_stock; correct=true; tag=)
  Stored hint: The customer bought new ownership shares rather than services. What resource and capital contribution arose?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- c: Debit Common Stock $100 / Credit Cash $100 (id=opt_dr_stock_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does issuing common shares affect the accounting equation?

Hint: Does selling ownership to a customer create income from services?

Explanation: Assets increase by $100 (Cash) and equity increases by $100 through Common Stock. Liabilities and net income are unchanged. This is contributed capital, not earnings.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Common Stock); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Does selling ownership to a customer create income from services?
- b: No net change in total assets. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does selling ownership to a customer create income from services?
- c: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Does selling ownership to a customer create income from services?

## issue_shares_separate_private_sale

A corporation issues new common shares today and receives $100 cash directly from an investor for those new shares. Separately, two existing stockholders trade old shares between themselves using their own money; the corporation receives or pays nothing in that private trade. Neither event involves customer services or a loan. Record only the corporation's receipt for newly issued shares, with the full amount as share capital; ignore par-value splits and fees.

### identify_account

What did the company receive from the investors today?

Hint: Which share transaction actually brings money into the corporation?

Explanation: Cash records the corporation's receipt for its new shares. Money exchanged privately between stockholders belongs to them and does not change the corporation's cash.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Which share transaction actually brings money into the corporation?
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

Hint: Does the private trade cancel the money the corporation received for new shares?

Explanation: Cash increases by the stated new-share receipt. The separate private trade creates no company cash movement and does not reverse the issuance.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does the private trade cancel the money the corporation received for new shares?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does the private trade cancel the money the corporation received for new shares?

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

Hint: Did the corporation issue new ownership interests for the money it received?

Explanation: Common Stock increases for the corporation's new issuance. A private resale does not create additional company capital. The receipt is neither customer revenue nor repayable financing.

- a: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- b: Notes Payable (Liability) (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: Did the corporation issue new ownership interests for the money it received?
- c: Common Stock (Equity) (id=opt_common_stock; correct=true; tag=)
  Stored hint: Did the corporation issue new ownership interests for the money it received?

### balanced_entry

What is the complete balanced journal entry for issuing common stock for cash?

Hint: Record the company's own issuance only. What resource and ownership contribution arose?

Explanation: Debit Cash $100 and Credit Common Stock $100. No entry for the private trade belongs in the corporation's books. Recording nothing would omit the actual new issuance; revenue would misclassify contributed capital.

- a: Debit Cash $100 / Credit Common Stock $100 (id=opt_dr_cash_cr_stock; correct=true; tag=)
  Stored hint: Record the company's own issuance only. What resource and ownership contribution arose?
- b: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recorded_on_share_issue)
  Stored hint: Did the investors pay for a service, or buy an ownership stake?
- c: Debit Common Stock $100 / Credit Cash $100 (id=opt_dr_stock_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does issuing common shares affect the accounting equation?

Hint: Which transaction changes the corporation's resources and contributed capital?

Explanation: Assets increase by $100 (Cash) and equity increases by $100 (Common Stock) for the new issuance. Liabilities and income are unchanged. The private stockholder trade has no company equation effect.

- a: Assets increase by $100 (+Cash); Equity increases by $100 (+Common Stock); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=true; tag=)
  Stored hint: Which transaction changes the corporation's resources and contributed capital?
- b: No net change in total assets. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Which transaction changes the corporation's resources and contributed capital?
- c: Assets increase by $100; Liabilities increase by $100; Equity unchanged. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Which transaction changes the corporation's resources and contributed capital?

## borrow_cash_prior_investment

A corporation receives $100 cash today from a stockholder under a newly signed promissory note requiring repayment. The stockholder also bought newly issued shares last month, and that separate investment was already recorded then. No shares are issued today and no customer services are involved. Record only today's new loan; ignore interest.

### identify_account

What did the company receive today, and which account records it?

Hint: What resource arrived today separately from the investment already recorded last month?

Explanation: Cash records today's new loan proceeds. The prior investment remains recorded and is not recorded a second time.

- a: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: What arrived today, and is the company still waiting to collect that money?
- b: Service Revenue (id=opt_service_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Set aside how the payment was earned or financed. What resource actually arrived today?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource arrived today separately from the investment already recorded last month?

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

Hint: Does today's agreement require repayment or issue more ownership?

Explanation: Notes Payable records the new repayment obligation. Common Stock relates to last month's already recorded share issuance; being a stockholder does not make this new loan another capital contribution.

- a: Notes Payable (Liability) (id=opt_notes_payable; correct=true; tag=)
  Stored hint: Does today's agreement require repayment or issue more ownership?
- b: Common Stock (Equity) (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: Does today's agreement require repayment or issue more ownership?
- c: Service Revenue (Revenue) (id=opt_service_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?

### balanced_entry

What is the complete balanced journal entry for borrowing cash on a note?

Hint: Today's receipt requires repayment and issues no shares. What resource and obligation arose?

Explanation: Debit Cash $100 and Credit Notes Payable $100. Do not repeat last month's share contribution. Revenue would treat borrowing as earning; a stock credit would ignore the new note.

- a: Debit Cash $100 / Credit Service Revenue $100 (id=opt_dr_cash_cr_rev; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?
- b: Debit Notes Payable $100 / Credit Cash $100 (id=opt_dr_notes_cr_cash; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Cash $100 / Credit Notes Payable $100 (id=opt_dr_cash_cr_notes; correct=true; tag=)
  Stored hint: Today's receipt requires repayment and issues no shares. What resource and obligation arose?

### equation_effect

How does borrowing cash affect the accounting equation?

Hint: Does an earlier investment remove the repayment obligation under today's note?

Explanation: Assets increase by $100 (Cash) and liabilities increase by $100 (Notes Payable). Equity and income are unchanged. The previously recorded share capital remains unchanged.

- a: Assets increase by $100 (+Cash); Liabilities increase by $100 (+Notes Payable); Equity is unchanged. (id=opt_assets_up_liab_up; correct=true; tag=)
  Stored hint: Does an earlier investment remove the repayment obligation under today's note?
- b: Assets increase by $100 (+Cash); Equity increases by $100 (+Revenue); Liabilities unchanged. (id=opt_assets_up_eq_up; correct=false; tag=revenue_recorded_on_borrowing)
  Stored hint: Did the lender pay for completed work or supply money the company must repay?
- c: No net change in total assets. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does an earlier investment remove the repayment obligation under today's note?

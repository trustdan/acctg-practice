# Expansion batch 12 rendered review

Draft preview at $100, seed 101. All seven stages include custom/inherited teaching and stored option hints.

## equipment_cash_additional_machine

A company pays $100 cash today to buy and take delivery and ownership of an additional machine for several future years of business use. It already owns another machine whose cost was recorded earlier; that older machine stays owned and in use. Today's stated amount is only the new machine's full price, not the combined cost of both machines. No trade-in, repair, disposal, earlier payable, or loan transaction occurs. Record only the new purchase, with no depreciation or other costs.

### identify_account

Which account records the newly acquired machinery?

Hint: Does owning an older machine prevent the new delivered machine from adding a future resource?

Explanation: Equipment records the additional machine. The older recorded machine stays in the asset balance; the new purchase is not a cost of repairing it.

- a: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Does owning an older machine prevent the new delivered machine from adding a future resource?
- b: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- c: Equipment (id=opt_equipment; correct=true; tag=)
  Stored hint: Does owning an older machine prevent the new delivered machine from adding a future resource?
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

Hint: Is a second useful machine added or does it replace an entry already recorded?

Explanation: Equipment increases by the new machine's price. Its balance already includes the older machine, whose cost is not recorded again.

- a: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Is a second useful machine added or does it replace an entry already recorded?
- b: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Is a second useful machine added or does it replace an entry already recorded?

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

Hint: How was only the additional machine paid for today?

Explanation: Cash decreases by the new machine's full price. There is no new borrowing, unpaid balance, or payment for the older machine.

- a: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: How was only the additional machine paid for today?
- b: Notes Payable (id=opt_notes_payable; correct=false; tag=wrong_account)
  Stored hint: How was only the additional machine paid for today?
- c: Common Stock (id=opt_common_stock; correct=false; tag=wrong_account)
  Stored hint: How was only the additional machine paid for today?
- d: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?

### balanced_entry

What is the complete balanced journal entry for this cash equipment purchase?

Hint: What additional owned resource arrives and what payment occurs today?

Explanation: Debit Equipment $100 and Credit Cash $100. Leave the older machine's recorded cost intact; do not repost the combined cost or expense the new multi-year resource.

- a: Debit an expense for the machinery $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?
- b: Debit Cash $100 / Credit Equipment $100 (id=opt_dr_cash_cr_equip; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: Debit Equipment $100 / Credit Cash $100 (id=opt_dr_equip_cr_cash; correct=true; tag=)
  Stored hint: What additional owned resource arrives and what payment occurs today?

### equation_effect

How does purchasing equipment for cash affect the accounting equation?

Hint: Does adding a machine for cash change total resources or their composition?

Explanation: Equipment increases and Cash decreases by $100. Total assets, liabilities, and equity are unchanged; the older machine remains recorded.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_equipment_purchase)
  Stored hint: Was the machinery's entire benefit used today, or can it provide use over future years?
- b: Asset exchange: Equipment increases (+$100) and Cash decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does adding a machine for cash change total resources or their composition?
- c: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Does adding a machine for cash change total resources or their composition?

## repay_principal_bank_landlord

A corporation pays $100 cash today to a bank solely to repay part of the principal on a previously recorded promissory note. The bank also rents workspace to the corporation, but this month's rent was separately paid and recorded earlier; none of today's payment is rent. Some principal remains owed. No new borrowing occurs and today's payment contains no interest or fees; ignore interest. Record only the principal repayment.

### identify_account

Which account records what the company owed the bank before today's principal payment?

Hint: Which recorded obligation does today's principal payment reduce, despite the bank's other role?

Explanation: Notes Payable records the earlier loan. Today's payment reduces principal rather than recording rent already paid separately.

- a: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Which recorded obligation does today's principal payment reduce, despite the bank's other role?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Which recorded obligation does today's principal payment reduce, despite the bank's other role?
- c: An expense for the loan payment (id=opt_expense_misconception; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- d: Notes Payable (id=opt_notes_payable; correct=true; tag=)
  Stored hint: Which recorded obligation does today's principal payment reduce, despite the bank's other role?

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

Hint: After repaying part of the loan, is the principal owed higher or lower?

Explanation: Notes Payable decreases by the principal paid. The unpaid remainder stays owed; the bank's landlord role does not change the loan reduction.

- a: Increase (id=opt_increase; correct=false; tag=)
  Stored hint: After repaying part of the loan, is the principal owed higher or lower?
- b: Decrease (id=opt_decrease; correct=true; tag=)
  Stored hint: After repaying part of the loan, is the principal owed higher or lower?

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

Hint: What resource leaves to settle this part of the loan?

Explanation: Cash decreases by today's principal payment. Earlier rent is not paid or expensed again.

- a: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: What resource leaves to settle this part of the loan?
- b: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource leaves to settle this part of the loan?
- c: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: What resource leaves to settle this part of the loan?

### balanced_entry

What is the complete balanced journal entry for repaying note principal with cash?

Hint: What recorded loan balance decreases when cash pays principal only?

Explanation: Debit Notes Payable $100 and Credit Cash $100. Do not record Rent Expense or new borrowing; the separately recorded rent stays unchanged.

- a: Debit Notes Payable $100 / Credit Cash $100 (id=opt_dr_notes_cr_cash; correct=true; tag=)
  Stored hint: What recorded loan balance decreases when cash pays principal only?
- b: Debit an expense for the loan payment $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?
- c: Debit Cash $100 / Credit Notes Payable $100 (id=opt_dr_cash_cr_notes; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does repaying loan principal affect the accounting equation?

Hint: Does paying principal create another occupancy cost or reduce a recorded debt?

Explanation: Assets and liabilities decrease by $100. Equity and previously recorded Rent Expense are unchanged.

- a: Assets decrease by $100 (-Cash); Liabilities decrease by $100 (-Notes Payable); Equity is unchanged. (id=opt_assets_down_liab_down; correct=true; tag=)
  Stored hint: Does paying principal create another occupancy cost or reduce a recorded debt?
- b: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does paying principal create another occupancy cost or reduce a recorded debt?
- c: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_loan_repayment)
  Stored hint: Is this payment a new cost, or repayment of principal the company previously borrowed?

## dividend_cash_service_provider

A corporation's board declares and pays $100 cash today to stockholders solely in proportion to share ownership. One stockholder also supplied consulting services last month; those services and their separate payment were fully recorded then. None of today's distribution pays for services or settles a loan. No dividend was declared or recorded as owed earlier, and no shares are repurchased. Record only today's ownership distribution.

### identify_account

Which account records cash distributions paid directly to stockholders?

Hint: Does today's payment reward share ownership or buy the already settled consulting work?

Explanation: Dividends records the ownership distribution. A stockholder's separate service-provider role does not make this payment an operating expense.

- a: An expense for the payment to stockholders (id=opt_expense_misconception; correct=false; tag=expense_recorded_on_dividend)
  Stored hint: Were the stockholders paid for providing a service, or because they own shares?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Does today's payment reward share ownership or buy the already settled consulting work?
- c: Dividends (id=opt_dividends; correct=true; tag=)
  Stored hint: Does today's payment reward share ownership or buy the already settled consulting work?
- d: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Does today's payment reward share ownership or buy the already settled consulting work?

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

Hint: Does declaring this distribution add to the dividends accumulated this period or reverse a previous distribution?

Explanation: Dividends increases as a new distribution is recorded, even though equity decreases. The distribution account and total equity move in opposite directions.

- a: Decrease (id=opt_decrease; correct=false; tag=)
  Stored hint: Does declaring this distribution add to the dividends accumulated this period or reverse a previous distribution?
- b: Increase (id=opt_increase; correct=true; tag=)
  Stored hint: Does declaring this distribution add to the dividends accumulated this period or reverse a previous distribution?

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

Hint: What resource leaves for the newly declared and paid ownership distribution?

Explanation: Cash decreases by today's distribution. The earlier consulting payment is not repeated, and there is no earlier dividend payable to settle.

- a: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: What resource leaves for the newly declared and paid ownership distribution?
- b: Dividends Payable (id=opt_dividends_payable; correct=false; tag=wrong_account)
  Stored hint: What resource leaves for the newly declared and paid ownership distribution?
- c: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: What resource leaves for the newly declared and paid ownership distribution?

### balanced_entry

What is the complete balanced journal entry for paying cash dividends?

Hint: What ownership distribution and cash payment occur independently of last month's services?

Explanation: Debit Dividends $100 and Credit Cash $100. Do not record a service expense, loan repayment, or the already settled consulting payment.

- a: Debit Dividends $100 / Credit Cash $100 (id=opt_dr_div_cr_cash; correct=true; tag=)
  Stored hint: What ownership distribution and cash payment occur independently of last month's services?
- b: Debit an expense for the dividend $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_dividend)
  Stored hint: Were the stockholders paid for providing a service, or because they own shares?
- c: Debit Cash $100 / Credit Dividends $100 (id=opt_dr_cash_cr_div; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does paying cash dividends affect the accounting equation?

Hint: Does distributing cash for ownership buy another service or reduce the owners' stake?

Explanation: Assets and equity decrease by $100. Liabilities are unchanged. Dividends reduces equity directly and does not reduce net income as an expense.

- a: Assets decrease by $100; Liabilities decrease by $100. (id=opt_assets_down_liab_down; correct=false; tag=equation_effect_missed)
  Stored hint: Does distributing cash for ownership buy another service or reduce the owners' stake?
- b: No change in equity; dividends are an asset. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does distributing cash for ownership buy another service or reduce the owners' stake?
- c: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Dividends); Liabilities are unchanged. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does distributing cash for ownership buy another service or reduce the owners' stake?

## prepaid_purchase_broker_own_policy

An insurance brokerage company pays $100 cash today to an insurer for a policy protecting the brokerage's own office next month. It buys this protection for itself, not on behalf of a client, and receives no client payment or commission in this transaction. None of the coverage has begun or been used. The insurer is paid in full and no amount was previously recorded as owed. Record only the brokerage's purchase of its own future coverage.

### identify_account

Which account records the insurance coverage acquired today?

Hint: Is the brokerage buying protection for itself or earning a fee from a client?

Explanation: Prepaid Insurance records the company's own unused protection. Selling insurance services in other transactions does not make this purchase revenue or client money.

- a: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Is the brokerage buying protection for itself or earning a fee from a client?
- d: Cash (id=opt_cash; correct=false; tag=)
  Stored hint: Is the brokerage buying protection for itself or earning a fee from a client?

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

Hint: Does money arrive from a client or leave to buy the company's own coverage?

Explanation: Cash decreases by the premium paid to the insurer. No client receipt, commission, or unpaid obligation occurs here.

- a: Insurance Expense (id=opt_insurance_expense; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Does money arrive from a client or leave to buy the company's own coverage?
- c: Accounts Payable (id=opt_accounts_payable; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- d: Cash (id=opt_cash; correct=true; tag=)
  Stored hint: Does money arrive from a client or leave to buy the company's own coverage?

### balanced_entry

What is the complete balanced journal entry for purchasing this prepaid insurance policy?

Hint: What future benefit does the company itself acquire when it pays the insurer?

Explanation: Debit Prepaid Insurance $100 and Credit Cash $100. Do not record commission revenue or current Insurance Expense; the company's own protection is unused.

- a: Debit Prepaid Insurance $100 / Credit Accounts Payable $100 (id=opt_dr_prepaid_cr_ap; correct=false; tag=payable_recorded_for_cash_payment)
  Stored hint: Was the purchase paid for today, or is there still an unpaid amount owed?
- b: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- c: Debit Prepaid Insurance $100 / Credit Cash $100 (id=opt_dr_prepaid_cr_cash; correct=true; tag=)
  Stored hint: What future benefit does the company itself acquire when it pays the insurer?
- d: Debit Cash $100 / Credit Prepaid Insurance $100 (id=opt_dr_cash_cr_prepaid; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?

### equation_effect

How does purchasing prepaid insurance affect the accounting equation?

Hint: Does buying its own unused protection change total resources merely because the company is a brokerage?

Explanation: Prepaid Insurance increases and Cash decreases by $100. Total assets, liabilities, and equity are unchanged. Its line of business does not change this purchase.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_eq_down; correct=false; tag=expense_recorded_on_prepaid_purchase)
  Stored hint: Has the purchased coverage already been consumed, or is its protection still available for future months?
- b: Assets increase by $100; Liabilities increase by $100. (id=opt_assets_up_liab_up; correct=false; tag=equation_effect_missed)
  Stored hint: Does buying its own unused protection change total resources merely because the company is a brokerage?
- c: Asset exchange: Prepaid Insurance increases (+$100) and Cash decreases (-$100); Total Assets and Equity unchanged. (id=opt_asset_swap; correct=true; tag=)
  Stored hint: Does buying its own unused protection change total resources merely because the company is a brokerage?

## prepaid_consumption_final_remaining

At month-end, a company records $100 of insurance protection used during the final month of a multi-month policy. The policy was paid for and recorded as a prepaid asset earlier, and all earlier months' consumption was already recorded. The stated amount equals the entire remaining prepaid balance immediately before today's entry, not the original full premium. No future coverage remains, no cash is paid today, and this final portion has not previously been expensed. Record only the last remaining portion.

### identify_account

Which account records the cost of insurance coverage used this month?

Hint: What cost was used in the final month after earlier months were already recorded?

Explanation: Insurance Expense records only the final month's protection. Earlier consumption is not recorded again.

- a: Insurance Expense (id=opt_insurance_expense; correct=true; tag=)
  Stored hint: What cost was used in the final month after earlier months were already recorded?
- b: Rent Expense (id=opt_rent_expense; correct=false; tag=wrong_account)
  Stored hint: What cost was used in the final month after earlier months were already recorded?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=false; tag=)
  Stored hint: What cost was used in the final month after earlier months were already recorded?
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

Hint: Which remaining recorded benefit ends when the policy's final coverage expires?

Explanation: Prepaid Insurance decreases by the remaining balance to zero. The original premium and earlier consumed portions are not credited again; no cash moves.

- a: Cash (id=opt_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Accounts Payable (id=opt_accounts_payable; correct=false; tag=wrong_account)
  Stored hint: Which remaining recorded benefit ends when the policy's final coverage expires?
- c: Prepaid Insurance (id=opt_prepaid_insurance; correct=true; tag=)
  Stored hint: Which remaining recorded benefit ends when the policy's final coverage expires?
- d: Service Revenue (id=opt_service_revenue; correct=false; tag=wrong_account)
  Stored hint: Which remaining recorded benefit ends when the policy's final coverage expires?

### balanced_entry

What is the complete balanced journal entry for recording this expired insurance?

Hint: How much recorded protection remains to be consumed in this final entry?

Explanation: Debit Insurance Expense $100 and Credit Prepaid Insurance $100. Use only the remaining balance, not the historical full premium, and do not repeat payment.

- a: Debit Insurance Expense $100 / Credit Prepaid Insurance $100 (id=opt_dr_expense_cr_prepaid; correct=true; tag=)
  Stored hint: How much recorded protection remains to be consumed in this final entry?
- b: Debit Prepaid Insurance $100 / Credit Insurance Expense $100 (id=opt_dr_prepaid_cr_expense; correct=false; tag=reversed_sides)
  Stored hint: For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?
- c: No entry; the full policy remains prepaid. (id=opt_no_entry; correct=false; tag=prepaid_not_expensed_on_consumption)
  Stored hint: Can coverage already used this month still be reported as protection available for the future?
- d: Debit Insurance Expense $100 / Credit Cash $100 (id=opt_dr_expense_cr_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?

### equation_effect

How does recognizing expired insurance affect the accounting equation?

Hint: Does final expiry use the remaining benefit or require paying the original premium again?

Explanation: Assets decrease through Prepaid Insurance and equity decreases through Insurance Expense by $100. Cash and liabilities are unchanged; this policy's prepaid balance reaches zero.

- a: Assets decrease by $100 (-Cash); Equity decreases by $100 (-Insurance Expense). (id=opt_assets_down_cash; correct=false; tag=cash_recorded_on_prepaid_expiration)
  Stored hint: Was the policy paid for today, or was today's used coverage paid for earlier?
- b: Assets decrease by $100 (-Prepaid Insurance); Equity decreases by $100 (-Insurance Expense); Cash is unaffected. (id=opt_assets_down_eq_down; correct=true; tag=)
  Stored hint: Does final expiry use the remaining benefit or require paying the original premium again?
- c: No net change in total assets; asset swap only. (id=opt_no_net_change; correct=false; tag=equation_effect_missed)
  Stored hint: Does final expiry use the remaining benefit or require paying the original premium again?

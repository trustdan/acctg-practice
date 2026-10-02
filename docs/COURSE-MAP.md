# Course evidence map

## Course materials status

No local syllabus, slide decks, or problem-set files were located in the repository root or subdirectories during Stage 00 orientation. Per `docs/AGENT-CONTRACT.md` and `PLAN.md`, course conventions must not be invented; standard US GAAP financial accounting conventions govern.

## Convention status & evidence register

| Topic/convention | Source path + page/slide | Verified scope/week | Product consequence | Status |
|---|---|---|---|---|
| Basic double-entry & normal balances | Generic US GAAP baseline | Intro / Week 1 | Asset/Exp/Div = Dr normal; Liab/Eq/Rev = Cr normal | Generic baseline (provisional) |
| Journal entries | Generic 2-line & multi-line format | Week 1–2 | Standard Dr on top, Cr indented below | Generic baseline (provisional) |
| T-accounts | Generic debit (left) / credit (right) | Week 1–2 | Left = Dr, Right = Cr | Generic baseline (provisional) |
| Two-row three-column grid | Unverified (student report) | Week 1 | Assets/Liabilities/Equity across cols; Dr/Cr rows | Provisional generic layout (unverified) |
| Spreadsheet exercise format | Pending | Week 2+ | Multi-event horizontal transaction worksheet | Unverified (deferred to Stage 14) |
| Accounting standards framework | Generic US GAAP | Week 1–10 | US GAAP rules applied; IFRS flagged if course specifies | Generic baseline (US GAAP) |
| Assessment format / allowed aids | Pending | Exams | Exam mode suppresses hints and reference grid | Unverified (deferred to Stage 10+) |

## Vocabulary & confirmed aliases (Provisional)

- `Cash`: Asset, normal debit.
- `Accounts Receivable`: Asset, normal debit. (Alias: Trade Receivables).
- `Prepaid Expenses` (e.g. `Prepaid Insurance`): Asset, normal debit.
- `Equipment`: Asset, normal debit.
- `Accounts Payable`: Liability, normal credit. (Alias: Trade Payables).
- `Unearned Revenue`: Liability, normal credit. (Alias: Deferred Revenue; NOT to be called 'unearned income').
- `Notes Payable`: Liability, normal credit.
- `Common Stock`: Equity, normal credit. (Alias: Contributed Capital / Share Capital).
- `Retained Earnings`: Equity, normal credit.
- `Service Revenue`: Revenue, normal credit (increases Equity via Retained Earnings).
- `Rent Expense`, `Insurance Expense`: Expense, normal debit (decreases Equity).
- `Dividends`: Dividends, normal debit (contra-equity / direct reduction to Retained Earnings; not an expense).

## Scope boundaries

### Initial topic inclusion list (Stages 02–07)
1. **Cash service** (`cash_service`): Earned revenue received in cash immediately (Dr Cash, Cr Service Revenue).
2. **Service on credit** (`service_on_credit`): Earned revenue billed on account (Dr Accounts Receivable, Cr Service Revenue).
3. **Customer advance** (`customer_advance`): Unearned cash received for future performance (Dr Cash, Cr Unearned Revenue).
4. **Receivable collection** (`receivable_collection`): Cash collected for previously recognized revenue (Dr Cash, Cr Accounts Receivable). Contrast with revenue recognition.
5. **Earning an advance** (`earn_advance`): Performance of services previously paid for (Dr Unearned Revenue, Cr Service Revenue; no cash movement).
6. **Operating cash expense** (`cash_expense`): Current period expense paid in cash (Dr Expense, Cr Cash).
7. **Prepaid asset purchase** (`prepaid_purchase`): Future benefit paid in cash (Dr Prepaid Asset, Cr Cash).
8. **Prepaid consumption** (`prepaid_consumption`): Using up prepaid assets (Dr Expense, Cr Prepaid Asset).
9. **Equipment purchase for cash** (`equipment_purchase_cash`): Capital asset acquisition (Dr Equipment, Cr Cash).
10. **Borrowing via note** (`borrowing_note`): Financing receipt (Dr Cash, Cr Notes Payable; no revenue).
11. **Principal repayment** (`repay_note_principal`): Loan payoff (Dr Notes Payable, Cr Cash; no expense).
12. **Owner capital contribution** (`issue_stock_cash`): Equity financing (Dr Cash, Cr Common Stock; no revenue).
13. **Cash dividends** (`dividend_cash`): Distribution to owners (Dr Dividends, Cr Cash; not an expense).

### Initial topic exclusion list (Deferred until later stages or explicit course coverage)
- Accruals with multi-period interest compounding.
- Bad debt expense and Allowance for Doubtful Accounts (contra-asset).
- Inventory valuation methods (FIFO, LIFO, Weighted Average, Periodic vs Perpetual).
- Fixed asset depreciation methods (Straight-line, MACRS, Accumulated Depreciation contra-asset).
- Intangibles, amortization, and goodwill impairment.
- Long-term bonds (discounts, premiums, effective interest rate method).
- Direct vs indirect Statement of Cash Flows operating cash reconciliation.
- Multi-currency, tax accounting, and complex leases.

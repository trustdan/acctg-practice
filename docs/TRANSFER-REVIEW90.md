# Transfer and distractor review at 90

Current status (October 3, 2026): the focused quality work at 90 is complete under pedagogy policy 2. [Quality closeout](QUALITY-REVIEW90.md) records ten amount profiles, contextual account choices and the remaining external validation gates. Earlier sections retain the findings and next steps at their original checkpoints.

Reviewed 2026-10-04T02:05:01Z (October 3 local) by Codex under trustdan's continuing delegated review. This records assistant review, not personal human reading. No actual syllabus or slides were available. Bank remains 90 active, one retired, 13 families; no new topics, accounts or accounting rules.

## Policy and replay

`curriculum/pedagogy.json` policy 1 binds all 90 active templates at their existing versions. Question version plus pedagogy policy version identifies the new generated snapshot; instance IDs include the policy version. Strict loading rejects invalid provenance, unknown fields, duplicate references, unavailable/inactive partners and incompatible profiles. Unknown or changed templates receive no inferred setting credit. Older question snapshots, options, answers and attempts remain immutable.

Migration 9 adds snapshotted pedagogy and attempt setting/policy fields. Historical defaults are empty/zero, with no retrospective classification. Existing migration backup handling applies. Evidence projection version 2 preserves historical independent grades but can restore full guidance and the default three-day half-life until reviewed transfer evidence exists. No personal database was opened for this work.

## Selection and evidence

After a wrong account, counter-account, balanced-entry or equation decision, the next available declared partner can replace the next ordinary session question. Selection is bounded to one queue, preserves the amount where permitted, verifies target version/active status, and never adds session questions. Missing partners fall back to ordinary selection. Guided partners use full scaffolding, carry contrast assistance through every stage, and cannot queue another contrast.

Ordinary weighting retains its exact-ID repeat penalty and additionally penalizes the last three reviewed reasoning groups. Business names and changed amounts alone do not create a distinct setting. The grouped baseline intentionally combines cosmetic rewrites; other groups identify event isolation, counterparty role, partial/final amounts, paperwork/deadline cues, entity boundary, labels, performance timing and benefit boundary. These are reviewed reasoning distinctions, not validated psychometric categories.

Intermediate guidance requires score >=0.60 and two independent successes; faded guidance requires score >=0.80 and three. Both additionally require two reviewed successful groups and at least one delayed success changing the group from the most recent reviewed independent success. Delay is at least ten minutes since both the previous independent response and the most recent exposure, including assisted practice. Guided answers never improve independent evidence or transfer counts. Half-life remains three days until two delayed transfers support bounded adjustment. Duplicate attempt IDs and repeated concept substeps cannot multiply independent credit.

## Matched pairs

Each pair below is symmetric. Review checked the deciding event fact and balanced wrong-event rejection at all source amounts. Declared order determines preference when a source has multiple partners.

- `cash_service_unpaid_booking` / `customer_advance_training`
- `prepaid_consumption_clinic` / `prepaid_purchase_clinic`
- `equipment_purchase_workshop` / `prepaid_purchase_machine_cover`
- `cash_rent_before_opening` / `equipment_purchase_workshop`
- `borrow_cash_future_equipment` / `issue_shares_future_equipment`
- `borrow_cash_future_equipment` / `repay_principal_equipment_loan`
- `cash_rent_stockholder_landlord` / `dividend_cash_working_owners`
- `earn_advance_separate_unpaid_job` / `service_on_credit_new_client`
- `cash_service_stockholder_customer` / `issue_shares_customer_investor`
- `customer_advance_same_day` / `earn_advance_same_day`
- `cash_service_third_party` / `collect_receivable_third_party`
- `collect_receivable_before_due` / `customer_advance_full_price`
- `cash_service_undeposited_notes` / `service_on_credit_customer_withdrawal`

These comparisons distinguish receipt from earning, completed service from old-invoice collection, customer payment from ownership, company cash from customer cash, collection from advance, unused benefit from consumption, loan from capital, rent from owner distribution, equipment from future coverage, and principal repayment from borrowing. The same-day pair separates morning receipt from evening performance; the mixed-invoice pair isolates the prepaid job from another unpaid job.

## Revised distractors

Six scenario profiles add missing service-versus-capital, held-cash-versus-no-entry, collection-versus-advance, seller-books-versus-no-entry and earned-advance-versus-receivable choices. Correct keys and canonical entries are unchanged. Revised choices have stored causal mistake hints and deterministic ordering, unique IDs and no more than four choices. [Complete rendered preview](TRANSFER-DISTRACTOR-PREVIEW90.md) was read across all six profiles before approval. Tests cover all allowed amounts and scaffold levels, hints, retries, grading and contrast assistance.

- `cash_service_stockholder_customer`: `service_vs_capital`
- `collect_receivable_future_booking`: `collection_vs_advance`
- `cash_service_undeposited_notes`: `held_cash_vs_no_entry`
- `collect_receivable_before_due`: `collection_vs_advance`
- `customer_advance_customer_books`: `seller_receipt_vs_no_entry`
- `earn_advance_separate_unpaid_job`: `earned_advance_vs_receivable`

## Reviewed setting inventory

This explicit inventory governs credit; names or suffixes are not used to infer groups at runtime.

| Question | Template version | Reviewed group |
|---|---:|---|
| cash_service_basic | 4 | cash_service.baseline |
| customer_advance_basic | 3 | customer_advance.baseline |
| service_on_credit_basic | 3 | service_on_credit.baseline |
| collect_receivable_basic | 3 | collect_receivable.baseline |
| earn_advance_basic | 4 | earn_advance.baseline |
| cash_rent_basic | 3 | cash_rent.baseline |
| borrow_cash_basic | 2 | borrow_cash.baseline |
| issue_shares_basic | 2 | issue_shares.baseline |
| prepaid_insurance_retail | 2 | prepaid_purchase.baseline |
| prepaid_consumption_retail | 2 | prepaid_consumption.baseline |
| equipment_cash_rental | 2 | equipment_purchase_cash.baseline |
| repay_principal_logistics | 2 | repay_note_principal.baseline |
| dividend_cash_retail | 3 | dividend_cash.baseline |
| cash_service_webdev | 3 | cash_service.baseline |
| service_credit_consulting | 3 | service_on_credit.baseline |
| collect_receivable_consulting | 3 | collect_receivable.baseline |
| customer_advance_logistics | 3 | customer_advance.baseline |
| earn_advance_logistics | 4 | earn_advance.baseline |
| cash_service_training | 1 | cash_service.baseline |
| customer_advance_training | 1 | customer_advance.baseline |
| earn_advance_training | 1 | earn_advance.baseline |
| service_on_credit_repairs | 1 | service_on_credit.baseline |
| collect_receivable_repairs | 1 | collect_receivable.baseline |
| borrow_cash_future_payroll | 1 | borrow_cash.benefit_boundary |
| prepaid_purchase_clinic | 1 | prepaid_purchase.baseline |
| prepaid_consumption_clinic | 1 | prepaid_consumption.baseline |
| cash_rent_clinic | 1 | cash_rent.baseline |
| issue_shares_studio | 1 | issue_shares.baseline |
| repay_principal_studio | 1 | repay_note_principal.baseline |
| dividend_cash_studio | 1 | dividend_cash.baseline |
| equipment_purchase_workshop | 1 | equipment_purchase_cash.baseline |
| equipment_purchase_replacement | 1 | equipment_purchase_cash.benefit_boundary |
| cash_rent_workshop | 1 | cash_rent.baseline |
| borrow_cash_future_equipment | 1 | borrow_cash.benefit_boundary |
| issue_shares_future_equipment | 1 | issue_shares.benefit_boundary |
| prepaid_purchase_renewal | 1 | prepaid_purchase.benefit_boundary |
| prepaid_consumption_partial_policy | 1 | prepaid_consumption.partial_final |
| repay_principal_equipment_loan | 1 | repay_note_principal.benefit_boundary |
| dividend_cash_working_owners | 1 | dividend_cash.counterparty_role |
| borrow_cash_stockholder_note | 1 | borrow_cash.counterparty_role |
| cash_service_separate_job | 1 | cash_service.separate_event |
| service_on_credit_new_client | 1 | service_on_credit.separate_event |
| collect_receivable_installment | 1 | collect_receivable.partial_final |
| customer_advance_additional_receipt | 1 | customer_advance.partial_final |
| earn_advance_completed_session | 1 | earn_advance.partial_final |
| equipment_cash_existing_note | 1 | equipment_purchase_cash.separate_event |
| prepaid_purchase_machine_cover | 1 | prepaid_purchase.benefit_boundary |
| prepaid_consumption_no_claim | 1 | prepaid_consumption.benefit_boundary |
| repay_principal_final_payment | 1 | repay_note_principal.partial_final |
| dividend_cash_prior_earnings | 1 | dividend_cash.benefit_boundary |
| cash_rent_before_opening | 1 | cash_rent.performance_timing |
| cash_rent_stockholder_landlord | 1 | cash_rent.counterparty_role |
| issue_shares_customer_investor | 1 | issue_shares.counterparty_role |
| issue_shares_separate_private_sale | 1 | issue_shares.separate_event |
| borrow_cash_prior_investment | 1 | borrow_cash.separate_event |
| cash_service_unpaid_booking | 1 | cash_service.performance_timing |
| service_on_credit_due_today | 1 | service_on_credit.paperwork_deadline |
| collect_receivable_overdue | 1 | collect_receivable.paperwork_deadline |
| customer_advance_same_day | 1 | customer_advance.performance_timing |
| earn_advance_same_day | 1 | earn_advance.performance_timing |
| prepaid_purchase_discounted_policy | 1 | prepaid_purchase.benefit_boundary |
| prepaid_consumption_closed_month | 1 | prepaid_consumption.benefit_boundary |
| equipment_cash_secondhand | 1 | equipment_purchase_cash.benefit_boundary |
| repay_principal_early_partial | 1 | repay_note_principal.partial_final |
| dividend_cash_total_distribution | 1 | dividend_cash.partial_final |
| cash_service_third_party | 1 | cash_service.counterparty_role |
| service_on_credit_invoice_later | 1 | service_on_credit.paperwork_deadline |
| collect_receivable_third_party | 1 | collect_receivable.counterparty_role |
| customer_advance_full_price | 1 | customer_advance.label |
| earn_advance_paperwork_later | 1 | earn_advance.paperwork_deadline |
| cash_rent_separate_insurance | 1 | cash_rent.separate_event |
| borrow_cash_collateral | 1 | borrow_cash.benefit_boundary |
| issue_shares_existing_lender | 1 | issue_shares.separate_event |
| prepaid_purchase_budget_label | 1 | prepaid_purchase.label |
| prepaid_consumption_recent_payment | 1 | prepaid_consumption.performance_timing |
| equipment_cash_additional_machine | 1 | equipment_purchase_cash.partial_final |
| repay_principal_bank_landlord | 1 | repay_note_principal.counterparty_role |
| dividend_cash_service_provider | 1 | dividend_cash.separate_event |
| prepaid_purchase_broker_own_policy | 1 | prepaid_purchase.entity_boundary |
| prepaid_consumption_final_remaining | 1 | prepaid_consumption.partial_final |
| cash_service_stockholder_customer | 1 | cash_service.counterparty_role |
| service_on_credit_customer_withdrawal | 1 | service_on_credit.entity_boundary |
| collect_receivable_future_booking | 1 | collect_receivable.separate_event |
| customer_advance_internal_setup | 1 | customer_advance.performance_timing |
| earn_advance_final_remaining | 1 | earn_advance.partial_final |
| cash_service_undeposited_notes | 1 | cash_service.entity_boundary |
| service_on_credit_cash_sale_label | 1 | service_on_credit.label |
| collect_receivable_before_due | 1 | collect_receivable.paperwork_deadline |
| customer_advance_customer_books | 1 | customer_advance.entity_boundary |
| earn_advance_separate_unpaid_job | 1 | earn_advance.separate_event |

## Validation and limits

Representative Windows offline PTY check: faded stockholder service rejected the new Common Stock choice, displayed the ownership-versus-price hint, accepted retry and equation, then selected the share-issuance partner at the same $100 with full guidance and the assisted-comparison banner. Quit cleanly. In-memory SQLite, isolated ignored auth/cache and explicit OfflineTutor; no provider calls. Full guided persistence, unavailable partners, non-chaining, migration/restart, delayed transfer and assisted-spacing are automated checks, not claimed manual full sessions.

This implements the reviewed contrast/transfer gate and six priority distractor repairs. Other weak generic fillers and amount-interpretation choices remain for focused review. Passing deterministic checks does not establish learning effectiveness. The 180 target, real course-source verification and native-platform/live-provider release checks remain open.

Final verification: `go test ./...`, `go vet ./...`, `gofmt -l .` and `git diff --check` passed; rebuilt `acctg.exe` passed 90/90 in-memory reconciliations. Build emitted the existing nonfatal module stat-cache permission warning.

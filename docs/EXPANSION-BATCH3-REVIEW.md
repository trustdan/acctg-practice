# Expansion batch 3: equipment and financing review

Five version-1 scenarios reviewed under the continuing delegation recorded in HANDOFF.md. Reviewer: **Codex (review delegated by trustdan)**, approved at **2026-10-03T22:47:07Z**. This records assistant review rather than personal human reading of each string. Complete approved package: [EXPANSION-BATCH3.json](EXPANSION-BATCH3.json). Rendered seven-stage wording, options, tags, and stored mistake hints: [EXPANSION-BATCH3-PREVIEW.md](EXPANSION-BATCH3-PREVIEW.md).

No syllabus or slides were found in the repository inventory. These scenarios use existing supported families and rule version 1; they establish no additional course coverage. Allowed amounts remain $50/$100/$200. Small amounts are practice simplifications within the existing equipment rules, not guidance on real-world capitalization thresholds.

| Scenario | Meaningful reasoning cue | Canonical entry | Contrast |
|---|---|---|---|
| equipment_purchase_workshop | Delivery and ownership occur now; production begins later. Waiting to use the machine does not postpone the purchase. | Dr Equipment / Cr Cash | cash_rent_workshop; borrowing and share issuance below |
| equipment_purchase_replacement | A whole replacement machine is acquired; the old machine is kept. Replacement is neither a repair service nor an implicit disposal. | Dr Equipment / Cr Cash | cash_rent_workshop; baseline equipment acquisition |
| cash_rent_workshop | Current use of rented workspace confers no ownership or future occupancy. | Dr Rent Expense / Cr Cash | Both equipment purchases |
| borrow_cash_future_equipment | A note finances a planned machine purchase; no purchase has occurred. Cash purpose does not change its source or create equipment. | Dr Cash / Cr Notes Payable | issue_shares_future_equipment; equipment_purchase_workshop |
| issue_shares_future_equipment | Same future purchase intention, but investors acquire ownership rather than a repayment claim. | Dr Cash / Cr Common Stock | borrow_cash_future_equipment; equipment_purchase_workshop |

The rent scenario transfers the asset-versus-current-cost distinction to the same workshop setting. The two financing scenarios transfer the previously reviewed future-payroll cue to a prospective asset acquisition: neither an intention nor a funding receipt records a purchase. They deliberately share future-purchase wording, changing the financing source. The equipment pair adds recognition before use and replacement without disposal; the two purchases are not asserted to have different canonical entries.

Reviewed every stage, including inherited classification and side teaching. Custom identity, counter-account, entry, and equation hints ask causal questions; replacement and share issuance also customize direction. Explanations distinguish acquired assets, current costs, actual cash movements, debt, and contributed capital. Checked each option and stored first-error hint for consistency. Existing expense, payable, reversed-side, revenue-on-financing, and equation-error tags remain appropriate. Some generic wrong-account hints and family filler distractors remain the previously documented low-priority limitations; this batch does not implement new option sets.

Timing is explicit, no prior purchase payable is settled, and no depreciation, disposal, interest, repair-expense family, equipment rental family, or multi-event journal is introduced. Replacing a machine does not reduce an old asset without a disposal. Production timing does not alter delivery/ownership. Equity increases from share issuance without increasing net income. Every scenario has complete seven-stage teaching through reviewed shared stages plus its overrides.

Before activation, `go test ./internal/drill -run TestExpansionBatch3 -count=1` passed against the draft package. Shared verification covers every allowed amount, full/intermediate/faded scaffolds, canonical postings, balanced equation effects, rendered teaching, wrong-answer/retry sessions, and rejection of opposite contrast entries that balance. The helper now checks the expected package size instead of assuming all batches contain six items; earlier six-item pilot gates remain intact. Publication equality regression ensures active copies exactly match this approved package.

Windows terminal checks before activation used an ignored harness with an in-memory database, isolated auth/model-cache paths, and an explicit OfflineTutor:

- Full equipment purchase: selected immediate expense, saw the future-use misconception hint, retried Equipment correctly, read the tailored feedback, requested offline explanation, scrolled through it, and exited cleanly.
- Faded replacement purchase: rejected immediate expense, accepted Equipment/Cash on retry, accepted the asset-exchange equation effect, rendered the balanced T-account recap, and quit.
- Faded share issuance: rejected Cash/Service Revenue, saw the ownership hint, accepted Cash/Common Stock, accepted the equity equation effect, rendered the balanced recap, and quit.

No live provider calls or personal learner history were used. Harness source was removed after use; compiled scratch artifacts stay ignored. Automated checks cover the remaining scenarios, rather than claiming individual manual walkthroughs for all five.

After activation, `go test ./...`, `go vet ./...`, `gofmt -l .`, and `git diff --check` passed. Root `acctg.exe` rebuilt successfully, with the existing nonfatal module stat-cache permission warning. Its `--reconcile-all --db :memory:` command verified **35/35** journal/T-account/equation reconciliations. CLI reconciliation uses its own fixed amount; regressions cover the scenario parameter amounts.

Gate passed. Current bank: **35 active, one retired, 13 families**; equipment, rent, borrowing, and shares now each have three scenarios. Current coverage is 105 scenario/amount combinations and 245 full-scaffold scenario-stage combinations. Five more complete scenarios advance toward 90; neither 90 nor 180 is complete. Contrast links are documentation and regression expectations, not automatic matched scheduling. No commit, release publication, or version change is implied.

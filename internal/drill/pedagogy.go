package drill

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// Profiles are named in the reviewed, versioned pedagogy policy. They only change distractors.
func applyReviewedDistractors(inst *domain.QuestionInstance, amount domain.Money, seed int64, profile string) {
	if inst.Pedagogy.PolicyVersion >= 2 {
		replacePayableFillers(inst, profile)
	}
	if profile == "" {
		return
	}
	entry := inst.StageAnswers[domain.StageBalancedEntry]
	text := amount.FormatDollars()
	var replacement domain.AnswerOption
	replaceID := ""
	mistakeHint := ""
	switch profile {
	case "service_vs_capital":
		mistakeHint = "Does this payment buy new ownership or settle the price of completed work?"
		replacement = domain.AnswerOption{ID: "opt_dr_cash_cr_stock", Text: "Debit Cash " + text + " / Credit Common Stock " + text, ErrorTag: engine.TagWrongAccount}
		replaceID = "opt_dr_rev_cr_cash"
		counter := inst.StageAnswers[domain.StageCounterAccount]
		for i, opt := range counter.Options {
			if opt.ID == "opt_ap" {
				counter.Options[i] = domain.AnswerOption{ID: "opt_common_stock", Label: opt.Label, Text: "Common Stock", ErrorTag: engine.TagWrongAccount}
			}
		}
		inst.StageAnswers[domain.StageCounterAccount] = counter
	case "held_cash_vs_no_entry":
		mistakeHint = "Does the company already control this currency before it reaches the bank?"
		replacement = domain.AnswerOption{ID: "opt_no_entry", Text: "No entry until the cash is deposited in the bank", ErrorTag: engine.TagWrongAccount}
		replaceID = "opt_dr_rev_cr_cash"
	case "collection_vs_advance":
		mistakeHint = "What work does this payment cover: the already recorded invoice or an undelivered service?"
		replacement = domain.AnswerOption{ID: "opt_dr_cash_cr_unearned", Text: "Debit Cash " + text + " / Credit Unearned Revenue " + text, ErrorTag: engine.TagWrongAccount}
	case "seller_receipt_vs_no_entry":
		mistakeHint = "Does the customer's bookkeeping record the seller's money or finish the seller's promised lesson?"
		replacement = domain.AnswerOption{ID: "opt_no_entry", Text: "No seller entry because the customer recorded the payment", ErrorTag: engine.TagWrongAccount}
		replaceID = "opt_dr_unearned_cr_cash"
	case "earned_advance_vs_receivable":
		mistakeHint = "For this completed prepaid job, is its price still owed by the customer?"
		replacement = domain.AnswerOption{ID: "opt_dr_ar_cr_rev", Text: "Debit Accounts Receivable " + text + " / Credit Service Revenue " + text, ErrorTag: engine.TagAdvanceConfusedWithReceivable}
	}
	if amountProfile, ok := bank.ReviewedAmountDistractor(profile); ok {
		replacement = domain.AnswerOption{ID: "opt_wrong_amount_scope", Text: "Debit " + amountProfile.Debit + " / Credit " + amountProfile.Credit + ", both for " + amountProfile.Scope, ErrorTag: engine.TagWrongAmount}
		mistakeHint = amountProfile.Hint
		// Prefer replacing side reversal so the event/timing misconception remains selectable.
		if len(entry.Options) >= 4 {
			for _, opt := range entry.Options {
				if opt.ErrorTag == engine.TagReversedSides {
					replaceID = opt.ID
					break
				}
			}
		}

	}
	options := append([]domain.AnswerOption(nil), entry.Options...)
	if replaceID != "" {
		for i, opt := range options {
			if opt.ID == replaceID {
				options[i] = replacement
			}
		}
	} else {
		options = append(options, replacement)
	}
	// Stable option IDs remain the grading identity; reshuffle the revised set deterministically.
	entry.Options = shuffleOptions(options, seed, 9001)
	if entry.MistakeHints == nil {
		entry.MistakeHints = make(map[string]string)
	}
	entry.MistakeHints[replacement.ID] = mistakeHint
	inst.StageAnswers[domain.StageBalancedEntry] = entry
}

// Payables remain useful for cash expense misconceptions; replace only generic revenue-family fillers.
func replacePayableFillers(inst *domain.QuestionInstance, profile string) {
	for _, key := range []domain.DrillStage{domain.StageIdentifyAccount, domain.StageCounterAccount} {
		stage := inst.StageAnswers[key]
		var option domain.AnswerOption
		hint := ""
		switch inst.FamilyID {
		case bank.FamilyCashService:
			if key == domain.StageIdentifyAccount {
				option = domain.AnswerOption{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount}
				hint = "Which resource arrived today, rather than which ownership account might explain its source?"
			} else {
				option = domain.AnswerOption{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagWrongAccount}
				hint = "Was this completed work paid for today, or is its price still owed?"
			}
			if key == domain.StageCounterAccount && profile == "service_vs_capital" {
				option = domain.AnswerOption{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount}
				hint = "Does this payment buy new ownership or settle the price of completed work?"
			}
		case bank.FamilyCustomerAdvance:
			option = domain.AnswerOption{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount}
			hint = "What resource arrived from the customer, and is this a loan or payment for future work?"
		case bank.FamilyServiceOnCredit:
			if key == domain.StageIdentifyAccount {
				option = domain.AnswerOption{ID: "opt_service_rev", Text: "Service Revenue", ErrorTag: engine.TagWrongAccount}
				hint = "Which resource represents the unpaid customer's obligation, rather than the earning itself?"
			} else {
				option = domain.AnswerOption{ID: "opt_common_stock", Text: "Common Stock", ErrorTag: engine.TagWrongAccount}
				hint = "Did the customer buy ownership, or receive completed work?"
			}
		case bank.FamilyCollectReceivable:
			if key == domain.StageIdentifyAccount {
				option = domain.AnswerOption{ID: "opt_ar", Text: "Accounts Receivable", ErrorTag: engine.TagWrongAccount}
				hint = "Which resource arrived today, rather than the claim that the payment settles?"
			} else {
				option = domain.AnswerOption{ID: "opt_notes_payable", Text: "Notes Payable", ErrorTag: engine.TagWrongAccount}
				hint = "Does the customer owe the company, or is the company repaying a lender?"
			}
		default:
			continue
		}
		for i, old := range stage.Options {
			if old.ID == "opt_ap" {
				option.Label = old.Label
				stage.Options[i] = option
				if stage.MistakeHints == nil {
					stage.MistakeHints = map[string]string{}
				}
				stage.MistakeHints[option.ID] = hint
			}
		}
		if key == domain.StageCounterAccount {
			if stage.MistakeHints == nil {
				stage.MistakeHints = map[string]string{}
			}
			for _, opt := range stage.Options {
				switch {
				case inst.FamilyID == bank.FamilyCashService && opt.ID == "opt_notes_payable":
					stage.MistakeHints[opt.ID] = "Does this receipt pay for completed work, or create money the company must repay to a lender?"
				case inst.FamilyID == bank.FamilyCustomerAdvance && opt.ID == "opt_common_stock":
					stage.MistakeHints[opt.ID] = "Did this customer buy new ownership, or pay for work the company still owes?"
				case inst.FamilyID == bank.FamilyCollectReceivable && opt.ID == "opt_unearned_rev":
					stage.MistakeHints[opt.ID] = "Does this payment settle previously recorded work, or pay for a separate service still owed?"
				}
			}
		}
		inst.StageAnswers[key] = stage
	}
}

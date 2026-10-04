package drill

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/engine"
)

// Snapshot reviewed misconception hints at generation time. Session replay
// consumes these strings, never a potentially newer routing table.
func snapshotMistakeHints(inst *domain.QuestionInstance) {
	for key, stage := range inst.StageAnswers {
		for _, option := range stage.Options {
			if option.ID == stage.CorrectOptionID || option.ErrorTag == "" || stage.MistakeHints[option.ID] != "" {
				continue
			}
			if hint := mistakeHint(inst.FamilyID, stage, option); hint != "" {
				if stage.MistakeHints == nil {
					stage.MistakeHints = make(map[string]string)
				}
				stage.MistakeHints[option.ID] = hint
			}
		}
		inst.StageAnswers[key] = stage
	}
}

func mistakeHint(family string, stage domain.StageAnswer, option domain.AnswerOption) string {
	// Account identity precedes entry assembly: do not steer the first step
	// toward the counter-account when the learner must identify cash received.
	if stage.Stage == domain.StageIdentifyAccount {
		switch family {
		case bank.FamilyCustomerAdvance, bank.FamilyCashService, bank.FamilyBorrowCash, bank.FamilyIssueShares, bank.FamilyCollectReceivable:
			if option.ErrorTag == engine.TagWrongAccount {
				return "What arrived today, and is the company still waiting to collect that money?"
			}
			return "Set aside how the payment was earned or financed. What resource actually arrived today?"
		case bank.FamilyServiceOnCredit:
			if option.ErrorTag == engine.TagCashRecordedWhenUncollected {
				return "Did the customer pay today, or does the company still have a right to collect later?"
			}
		case bank.FamilyEarnAdvance:
			if option.ErrorTag == engine.TagAdvanceConfusedWithReceivable {
				return "For this prepaid work, was the company waiting for money or still owing the promised service?"
			}
		}
	}
	switch option.ErrorTag {
	case engine.TagDuplicateRevenueOnCollection:
		return "Was this work already recorded as revenue when it was completed, or is new work being earned today?"
	case engine.TagRevenueRecognizedPrematurely:
		return "Has the company performed the promised work yet, or does it still owe the customer that work?"
	case engine.TagRevenueDeferredWhenEarned:
		return "After completing today's work, does the company still owe that work to the customer?"
	case engine.TagCashRecordedOnEarningAdvance:
		return "Does completing the prepaid work bring another payment, or was its receipt already recorded?"
	case engine.TagExpenseRecordedOnPrepaidPurchase:
		return "Has the purchased coverage already been consumed, or is its protection still available for future months?"
	case engine.TagPrepaidNotExpensedOnConsumption:
		return "Can coverage already used this month still be reported as protection available for the future?"
	case engine.TagCashRecordedOnPrepaidExpiration:
		return "Was the policy paid for today, or was today's used coverage paid for earlier?"
	case engine.TagExpenseRecordedOnEquipmentPurchase:
		return "Was the machinery's entire benefit used today, or can it provide use over future years?"
	case engine.TagExpenseRecordedOnLoanRepayment:
		return "Is this payment a new cost, or repayment of principal the company previously borrowed?"
	case engine.TagExpenseRecordedOnDividend:
		return "Were the stockholders paid for providing a service, or because they own shares?"
	case engine.TagRevenueRecordedOnBorrowing:
		return "Did the lender pay for completed work or supply money the company must repay?"
	case engine.TagRevenueRecordedOnShareIssue:
		return "Did the investors pay for a service, or buy an ownership stake?"
	case engine.TagCashRecordedWhenUncollected:
		return "Did any payment arrive today, or does the customer still owe money for completed work?"
	case engine.TagAdvanceConfusedWithReceivable:
		return "The customer has already paid. Is there still money to collect for this work?"
	case engine.TagPayableRecordedForCashPayment:
		return "Was the purchase paid for today, or is there still an unpaid amount owed?"
	case engine.TagReversedSides:
		return "For each account, does its balance increase or decrease, and does that change belong on its normal side or the opposite side?"
	case engine.TagEquationEffectMissed:
		return stage.CausalHint
	case engine.TagWrongAccount:
		return stage.CausalHint
	default:
		return ""
	}
}

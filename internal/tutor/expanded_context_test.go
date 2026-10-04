package tutor

import (
	"strings"
	"testing"
)

func TestReviewedExplanationOnlySentForExplanationMode(t *testing.T) {
	req := Request{ProblemPrompt: "Seller receives money for next month's lesson.", Explanation: "Seller owes the lesson; customer books do not satisfy its promise."}
	if !strings.Contains(FormatUserPrompt(req, false), req.Explanation) {
		t.Fatal("reviewed facts missing from provider explanation")
	}
	if strings.Contains(FormatUserPrompt(req, true), req.Explanation) {
		t.Fatal("hint request leaked reviewed answer")
	}
}

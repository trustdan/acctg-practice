package tui

import (
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
)

func TestTutorHintUsesSubmittedSnapshotAndDoesNotCarryErrorForward(t *testing.T) {
	first := domain.StageAnswer{Stage: domain.StageIdentifyAccount, CorrectOptionID: "cash", CausalHint: "First stage hint?", Options: []domain.AnswerOption{{ID: "cash", Text: "Cash"}, {ID: "rev", Text: "Service Revenue", ErrorTag: engine.TagRevenueRecognizedPrematurely}}, MistakeHints: map[string]string{"rev": "Stored misconception hint?"}}
	next := domain.StageAnswer{Stage: domain.StageAccountCategory, CorrectOptionID: "asset", CausalHint: "Next stage hint?", Options: []domain.AnswerOption{{ID: "asset", Text: "Asset"}}}
	inst := &domain.QuestionInstance{InstanceID: "hint-context", StageAnswers: map[domain.DrillStage]domain.StageAnswer{first.Stage: first, next.Stage: next}}
	m := &Model{CurrentInstance: inst, Session: drill.NewSession("hint", inst)}
	feedback, err := m.Session.SubmitOption("rev", time.Unix(0, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	m.LastFeedback = feedback
	// The selection moved to Cash, but tutor context must retain the submitted error.
	req := m.buildTutorRequest(&first)
	if req.MistakeHint != first.MistakeHints["rev"] || req.CausalHint != req.MistakeHint || req.SelectedOption.ID != "rev" {
		t.Fatalf("wrong submitted-option context: %+v", req)
	}
	feedback, err = m.Session.SubmitOption("rev", time.Unix(1, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	m.LastFeedback = feedback
	req = m.buildTutorRequest(&next)
	if req.ErrorTag != "" || req.MistakeHint != "" || req.CausalHint != next.CausalHint {
		t.Fatalf("stale error carried into next stage: %+v", req)
	}
}

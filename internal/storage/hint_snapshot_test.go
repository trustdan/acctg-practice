package storage_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/storage"
)

func TestMistakeHintSnapshotSurvivesRestartAndDuplicateSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hints.db")
	db, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	inst := setupTestQuestionInstance(t)
	stage := inst.StageAnswers[domain.StageCounterAccount]
	stage.MistakeHints = map[string]string{"opt_service_rev": "Original stored hint?"}
	inst.StageAnswers[stage.Stage] = stage
	if err := db.SaveSession(storage.SessionRecord{ID: "hint-restart", Mode: "progressive_drill", StartedAt: time.Unix(0, 0).UTC()}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveQuestionInstance(inst, "hint-restart"); err != nil {
		t.Fatal(err)
	}
	stage.MistakeHints["opt_service_rev"] = "Newer wording?"
	if err := db.SaveQuestionInstance(inst, "hint-restart"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	saved, err := db.GetQuestionInstance(inst.InstanceID)
	if err != nil {
		t.Fatal(err)
	}
	if hint := saved.StageAnswers[stage.Stage].HintForOption("opt_service_rev"); hint != "Original stored hint?" {
		t.Fatalf("snapshot rewritten: %q", hint)
	}
}

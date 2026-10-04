package storage_test

import (
	"bytes"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/storage"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestQualityPolicySnapshotsSurviveRestartAndDuplicateSave(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	seeds, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	gen := drill.NewGenerator(cat, engine.NewEngine(cat))
	path := filepath.Join(t.TempDir(), "quality.sqlite")
	db, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if db != nil {
			db.Close()
		}
	})
	if err = db.SaveSession(storage.SessionRecord{ID: "quality", Mode: "drill", StartedAt: time.Unix(1, 0)}); err != nil {
		t.Fatal(err)
	}
	for _, q := range seeds.Questions {
		if !bank.IsActiveForPractice(q.Status) {
			continue
		}
		inst, err := gen.GenerateInstance(q, 101, map[string]int64{"amount_minor_units": q.Parameters["amount_minor_units"][0]})
		if err != nil {
			t.Fatal(err)
		}
		if err = db.SaveQuestionInstance(inst, "quality"); err != nil {
			t.Fatal(err)
		}
		original := inst.PromptText
		inst.PromptText = "replacement must not overwrite the snapshot"
		if err = db.SaveQuestionInstance(inst, "quality"); err != nil {
			t.Fatal(err)
		}
		inst.PromptText = original
		if err = db.Close(); err != nil {
			t.Fatal(err)
		}
		db, err = storage.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		// Concepts are stored with attempts, not the question_instances row.
		inst.Concepts = nil
		restored, err := db.GetQuestionInstance(inst.InstanceID)
		if err != nil || !reflect.DeepEqual(restored, inst) {
			t.Fatalf("snapshot changed after restart: %s %v", q.ID, err)
		}
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
}

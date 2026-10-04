package storage_test

import (
	"bytes"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/storage"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestExamSnapshotsPreserveOrderAndTeachingAfterRestart(t *testing.T) {
	cat, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), cat)
	if err != nil {
		t.Fatal(err)
	}
	eng := engine.NewEngine(cat)
	gen := drill.NewGenerator(cat, eng)
	runner, err := exam.NewExamRunner(exam.ExamConfig{SessionID: "snapshot-restart", Questions: pack.Questions, Catalog: cat, Engine: eng, Generator: gen, TotalQuestions: 90, Seed: 101})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "exam.sqlite")
	db, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.SaveExamSession(runner.ToSessionRecord()); err != nil {
		t.Fatal(err)
	}
	var original []*domain.QuestionInstance
	for _, eq := range runner.Questions {
		original = append(original, eq.Instance)
		if err = db.SaveQuestionInstance(eq.Instance, runner.SessionID); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 2; i++ {
		_, st, err := runner.CurrentQuestion()
		if err != nil {
			t.Fatal(err)
		}
		if err = runner.SubmitOption(st.CorrectOptionID, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	record := runner.Interrupt()
	if err = db.SaveExamSession(record); err != nil {
		t.Fatal(err)
	}
	for _, att := range runner.Attempts {
		if err = db.RecordExamAttempt(att); err != nil {
			t.Fatal(err)
		}
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	snapshots, err := db.GetExamQuestionSnapshots(runner.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	for i, inst := range snapshots {
		if inst.InstanceID != original[i].InstanceID || inst.PromptText != original[i].PromptText || !reflect.DeepEqual(inst.StageAnswers, original[i].StageAnswers) || !inst.Entry.EqualNormalized(original[i].Entry) {
			t.Fatalf("question %d order or teaching changed", i)
		}
	}
	attempts, err := db.GetExamAttempts(runner.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	resumed, err := exam.ResumeExamSnapshots(record, attempts, snapshots)
	if err != nil {
		t.Fatal(err)
	}
	if resumed.CurrentQIndex != 1 {
		t.Fatal("exam restart lost position")
	}
	evidence, err := db.GetAllAttempts()
	if err != nil || len(evidence) != 0 {
		t.Fatal("exam evidence contaminated learning history")
	}
}

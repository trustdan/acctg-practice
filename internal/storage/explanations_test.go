package storage_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/storage"
)

func TestSavedExplanationSurvivesRestartAndVersion7Upgrade(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notes.db")
	db, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = raw.Exec("DROP TABLE saved_explanations; DELETE FROM schema_migrations WHERE version=8"); err != nil {
		t.Fatal(err)
	}
	raw.Close()
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	e := storage.SavedExplanation{ID: "note-1", InstanceID: "old-instance", QuestionID: "question", QuestionVersion: 2, Stage: "debit_credit", Scenario: "Customer pays an existing receivable.", StagePrompt: "Which side?", Explanation: "Cash increases; the receivable decreases.", Provider: "test-provider", GeneratedAt: now, SavedAt: now.Add(time.Minute)}
	if err = db.SaveExplanation(e); err != nil {
		t.Fatal(err)
	}
	if err = db.SaveExplanation(e); err != nil {
		t.Fatal(err)
	}
	db.Close()
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.ListSavedExplanations()
	if err != nil || len(rows) != 1 || rows[0] != e {
		t.Fatalf("lost/duplicate note after restart: %+v %v", rows, err)
	}
	attempts, err := db.GetAllAttempts()
	if err != nil || len(attempts) != 0 {
		t.Fatal("note became learner evidence")
	}
}

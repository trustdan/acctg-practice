package storage_test

import (
	"database/sql"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/storage"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestVersion8UpgradePreservesLegacyEvidenceAndTransferSnapshots(t *testing.T) {
	path := filepath.Join(t.TempDir(), "transfer.sqlite")
	db, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
	sess := storage.SessionRecord{ID: "transfer", Mode: "drill", StartedAt: now}
	if err = db.SaveSession(sess); err != nil {
		t.Fatal(err)
	}
	inst := &domain.QuestionInstance{InstanceID: "legacy", QuestionID: "seed", Version: 1, FamilyID: "cash_service", RuleVersion: 1, PromptText: "Original prompt", Parameters: map[string]int64{"amount_minor_units": 10000}, StageAnswers: map[domain.DrillStage]domain.StageAnswer{}, Entry: domain.Entry{}}
	if err = db.SaveQuestionInstance(inst, sess.ID); err != nil {
		t.Fatal(err)
	}
	att := domain.Attempt{AttemptID: "original", SessionID: sess.ID, InstanceID: inst.InstanceID, QuestionID: inst.QuestionID, QuestionVersion: 1, Stage: domain.StageIdentifyAccount, ConceptID: "cash_vs_revenue", SelectedOptionID: "cash", IsCorrect: true, Assistance: domain.AssistanceNone, GradingVersion: 1, AnsweredAt: now}
	if err = db.RecordAttempt(att); err != nil {
		t.Fatal(err)
	}
	db.Close()
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = raw.Exec("ALTER TABLE question_instances DROP COLUMN pedagogy_json; ALTER TABLE attempts DROP COLUMN setting_group; ALTER TABLE attempts DROP COLUMN pedagogy_version; DELETE FROM schema_migrations WHERE version=9")
	raw.Close()
	if err != nil {
		t.Fatal(err)
	}
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	legacy, err := db.GetQuestionInstance("legacy")
	if err != nil || legacy.PromptText != "Original prompt" || legacy.Pedagogy.PolicyVersion != 0 {
		t.Fatal("legacy snapshot rewritten")
	}
	all, err := db.GetAllAttempts()
	if err != nil || len(all) != 1 || all[0].PedagogyVersion != 0 || all[0].SettingGroup != "" || !all[0].IsCorrect {
		t.Fatal("legacy attempts retagged")
	}
	stats := mastery.RebuildProjections(all)["cash_vs_revenue"]
	if stats.Alpha != 2 || stats.SuccessfulSettings != 0 || stats.EvidenceVersion != 2 {
		t.Fatal("historical grade or transfer projection wrong")
	}
	inst.InstanceID = "guided"
	inst.Pedagogy = domain.PedagogySnapshot{PolicyVersion: 1, SettingGroup: "reviewed", Remediation: true, Contrasts: []domain.QuestionRef{{QuestionID: "partner", QuestionVersion: 1}}}
	if err = db.SaveQuestionInstance(inst, sess.ID); err != nil {
		t.Fatal(err)
	}
	att.AttemptID = "guided"
	att.InstanceID = "guided"
	att.SettingGroup = "reviewed"
	att.PedagogyVersion = 1
	att.Assistance = domain.AssistanceContrast
	att.AnsweredAt = now.Add(time.Hour)
	if err = db.RecordAttemptsBatch([]domain.Attempt{att, att}); err != nil {
		t.Fatal(err)
	}
	db.Close()
	db, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	restored, err := db.GetQuestionInstance("guided")
	if err != nil || !reflect.DeepEqual(restored.Pedagogy, inst.Pedagogy) {
		t.Fatal("pedagogy snapshot lost after restart")
	}
	all, err = db.GetAllAttempts()
	if err != nil || len(all) != 2 {
		t.Fatal("duplicate/failed evidence persistence")
	}
	stats = mastery.RebuildProjections(all)["cash_vs_revenue"]
	if stats.Alpha != 2 || stats.AssistedAttempts != 1 || stats.SuccessfulSettings != 0 {
		t.Fatal("guided evidence inflated historic success")
	}
}

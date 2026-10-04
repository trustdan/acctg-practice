package exam

import (
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestExpandedBankExamCoverageAndSnapshotReplay(t *testing.T) {
	runner, cat, eng, gen, all := setupTestExam(t, 101, 91, 0)
	if len(runner.Questions) != 90 {
		t.Fatalf("expected only 90 active exam questions, got %d", len(runner.Questions))
	}
	var snapshots []*domain.QuestionInstance
	for _, eq := range runner.Questions {
		snapshots = append(snapshots, eq.Instance)
		if eq.Instance.QuestionID == "legacy_unclear_advance_v0" {
			t.Fatal("retired content in exam")
		}
	}
	for i := 0; i < 2; i++ {
		_, st, err := runner.CurrentQuestion()
		if err != nil {
			t.Fatal(err)
		}
		if err = runner.SubmitOption(st.CorrectOptionID, time.Unix(1, 0)); err != nil {
			t.Fatal(err)
		}
	}
	record := runner.Interrupt()
	// The current bank can be empty, reordered, or retired: replay only uses saved snapshots.
	for i := range all {
		all[i].Status = bank.StatusRetired
	}
	resumed, err := ResumeExamSnapshots(record, runner.Attempts, snapshots)
	if err != nil {
		t.Fatal(err)
	}
	if resumed.CurrentQIndex != 1 || !reflect.DeepEqual(resumed.Questions[0].Instance, runner.Questions[0].Instance) {
		t.Fatal("snapshot replay changed original question or position")
	}
	for !resumed.IsCompleted() {
		_, st, err := resumed.CurrentQuestion()
		if err != nil {
			t.Fatal(err)
		}
		if err = resumed.SubmitOption(st.CorrectOptionID, time.Unix(2, 0)); err != nil {
			t.Fatal(err)
		}
	}
	report, err := resumed.Finish(time.Unix(3, 0))
	if err != nil {
		t.Fatal(err)
	}
	if report.CorrectItems != 180 || len(report.Reviews) != 180 {
		t.Fatalf("exam review coverage %+v", report)
	}
	for _, review := range report.Reviews {
		if review.Explanation == "" || len(review.Postings) != 2 {
			t.Fatal("post-exam teaching/entry missing")
		}
	}
	if _, err = NewExamRunner(ExamConfig{Questions: all, Catalog: cat, Engine: eng, Generator: gen}); err == nil {
		t.Fatal("inactive-only exam should fail")
	}
	if _, err = ResumeExamSnapshots(record, runner.Attempts, snapshots[:1]); err == nil {
		t.Fatal("incomplete snapshots should fail")
	}
}

package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/exam"
)

// SaveExamSession creates or updates an exam session in the dedicated exam_sessions table.
// It also ensures a parent row exists in sessions with mode='exam' so question_instances foreign keys hold.
func (d *DB) SaveExamSession(sess exam.ExamSessionRecord) error {
	var compVal any = nil
	if sess.CompletedAt != nil {
		compVal = sess.CompletedAt.UTC()
	}

	// 1. Ensure sessions parent row exists
	sessQuery := `
INSERT INTO sessions (id, mode, started_at, completed_at)
VALUES (?, 'exam', ?, ?)
ON CONFLICT(id) DO UPDATE SET
    completed_at = excluded.completed_at;`

	if _, err := d.db.Exec(sessQuery, sess.ID, sess.StartedAt.UTC(), compVal); err != nil {
		return fmt.Errorf("failed ensuring parent session for exam: %w", err)
	}

	// 2. Insert or update in dedicated exam_sessions table
	query := `
INSERT INTO exam_sessions (
    id, total_questions, time_limit_seconds, elapsed_seconds,
    status, started_at, completed_at, score, correct_count,
    total_attempts, seed
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    elapsed_seconds = excluded.elapsed_seconds,
    status = excluded.status,
    completed_at = excluded.completed_at,
    score = excluded.score,
    correct_count = excluded.correct_count,
    total_attempts = excluded.total_attempts;`

	_, err := d.db.Exec(query,
		sess.ID, sess.TotalQuestions, sess.TimeLimitSeconds, sess.ElapsedSeconds,
		string(sess.Status), sess.StartedAt.UTC(), compVal, sess.Score,
		sess.CorrectCount, sess.TotalAttempts, sess.Seed,
	)
	if err != nil {
		return fmt.Errorf("failed to save exam session: %w", err)
	}
	return nil
}

// GetExamSession retrieves an exam session by ID.
func (d *DB) GetExamSession(sessionID string) (*exam.ExamSessionRecord, error) {
	query := `
SELECT id, total_questions, time_limit_seconds, elapsed_seconds,
       status, started_at, completed_at, score, correct_count,
       total_attempts, seed
FROM exam_sessions WHERE id = ?;`

	row := d.db.QueryRow(query, sessionID)
	var s exam.ExamSessionRecord
	var statusStr string
	var compTime sql.NullTime

	err := row.Scan(
		&s.ID, &s.TotalQuestions, &s.TimeLimitSeconds, &s.ElapsedSeconds,
		&statusStr, &s.StartedAt, &compTime, &s.Score, &s.CorrectCount,
		&s.TotalAttempts, &s.Seed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exam session %s not found", sessionID)
		}
		return nil, fmt.Errorf("failed querying exam session: %w", err)
	}

	s.Status = exam.ExamStatus(statusStr)
	if compTime.Valid {
		t := compTime.Time
		s.CompletedAt = &t
	}
	return &s, nil
}

// ListExamSessions returns all recorded exam sessions ordered from newest to oldest.
func (d *DB) ListExamSessions() ([]exam.ExamSessionRecord, error) {
	query := `
SELECT id, total_questions, time_limit_seconds, elapsed_seconds,
       status, started_at, completed_at, score, correct_count,
       total_attempts, seed
FROM exam_sessions
ORDER BY started_at DESC;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed querying exam sessions: %w", err)
	}
	defer rows.Close()

	var sessions []exam.ExamSessionRecord
	for rows.Next() {
		var s exam.ExamSessionRecord
		var statusStr string
		var compTime sql.NullTime

		err := rows.Scan(
			&s.ID, &s.TotalQuestions, &s.TimeLimitSeconds, &s.ElapsedSeconds,
			&statusStr, &s.StartedAt, &compTime, &s.Score, &s.CorrectCount,
			&s.TotalAttempts, &s.Seed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning exam session: %w", err)
		}
		s.Status = exam.ExamStatus(statusStr)
		if compTime.Valid {
			t := compTime.Time
			s.CompletedAt = &t
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// GetLatestInterruptedExamSession returns the most recent exam session with status "interrupted" or "in_progress".
func (d *DB) GetLatestInterruptedExamSession() (*exam.ExamSessionRecord, error) {
	query := `
SELECT id, total_questions, time_limit_seconds, elapsed_seconds,
       status, started_at, completed_at, score, correct_count,
       total_attempts, seed
FROM exam_sessions
WHERE status IN ('interrupted', 'in_progress')
ORDER BY started_at DESC
LIMIT 1;`

	row := d.db.QueryRow(query)
	var s exam.ExamSessionRecord
	var statusStr string
	var compTime sql.NullTime

	err := row.Scan(
		&s.ID, &s.TotalQuestions, &s.TimeLimitSeconds, &s.ElapsedSeconds,
		&statusStr, &s.StartedAt, &compTime, &s.Score, &s.CorrectCount,
		&s.TotalAttempts, &s.Seed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No interrupted session
		}
		return nil, fmt.Errorf("failed checking interrupted exam session: %w", err)
	}

	s.Status = exam.ExamStatus(statusStr)
	if compTime.Valid {
		t := compTime.Time
		s.CompletedAt = &t
	}
	return &s, nil
}

// CompleteExamSession marks an exam session as completed with final scores and stats.
func (d *DB) CompleteExamSession(sessionID string, completedAt time.Time, score float64, correctCount, totalAttempts, elapsedSec int) error {
	query := `
UPDATE exam_sessions SET
    status = 'completed',
    completed_at = ?,
    score = ?,
    correct_count = ?,
    total_attempts = ?,
    elapsed_seconds = ?
WHERE id = ?;`

	res, err := d.db.Exec(query, completedAt.UTC(), score, correctCount, totalAttempts, elapsedSec, sessionID)
	if err != nil {
		return fmt.Errorf("failed completing exam session: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exam session %s not found", sessionID)
	}
	return nil
}

// InterruptExamSession sets the status of an exam session to "interrupted".
func (d *DB) InterruptExamSession(sessionID string, elapsedSec int) error {
	query := `UPDATE exam_sessions SET status = 'interrupted', elapsed_seconds = ? WHERE id = ?;`
	_, err := d.db.Exec(query, elapsedSec, sessionID)
	if err != nil {
		return fmt.Errorf("failed interrupting exam session: %w", err)
	}
	return nil
}

// AbandonExamSession marks an interrupted exam session as abandoned.
func (d *DB) AbandonExamSession(sessionID string) error {
	query := `UPDATE exam_sessions SET status = 'abandoned' WHERE id = ?;`
	res, err := d.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed abandoning exam session: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exam session %s not found", sessionID)
	}
	return nil
}

// RecordExamAttempt saves an individual exam attempt into the separate exam_attempts table.
func (d *DB) RecordExamAttempt(att exam.ExamAttemptRecord) error {
	isCorrInt := 0
	if att.IsCorrect {
		isCorrInt = 1
	}

	query := `
INSERT INTO exam_attempts (
    id, exam_session_id, instance_id, question_id, question_index,
    stage, concept_id, selected_option_id, correct_option_id,
    is_correct, error_tag, grading_version, answered_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;`

	_, err := d.db.Exec(query,
		att.ID, att.ExamSessionID, att.InstanceID, att.QuestionID, att.QuestionIndex,
		string(att.Stage), att.ConceptID, att.SelectedOptionID, att.CorrectOptionID,
		isCorrInt, att.ErrorTag, att.GradingVersion, att.AnsweredAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed recording exam attempt: %w", err)
	}
	return nil
}

// RecordExamAttemptsBatch saves multiple exam attempts inside a single atomic transaction.
func (d *DB) RecordExamAttemptsBatch(attempts []exam.ExamAttemptRecord) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO exam_attempts (
    id, exam_session_id, instance_id, question_id, question_index,
    stage, concept_id, selected_option_id, correct_option_id,
    is_correct, error_tag, grading_version, answered_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;`)
	if err != nil {
		return fmt.Errorf("failed preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, att := range attempts {
		isCorrInt := 0
		if att.IsCorrect {
			isCorrInt = 1
		}
		_, err := stmt.Exec(
			att.ID, att.ExamSessionID, att.InstanceID, att.QuestionID, att.QuestionIndex,
			string(att.Stage), att.ConceptID, att.SelectedOptionID, att.CorrectOptionID,
			isCorrInt, att.ErrorTag, att.GradingVersion, att.AnsweredAt.UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed batch recording exam attempt %s: %w", att.ID, err)
		}
	}

	return tx.Commit()
}

// GetExamAttempts retrieves all recorded attempts for an exam session in chronological order.
func (d *DB) GetExamAttempts(sessionID string) ([]exam.ExamAttemptRecord, error) {
	query := `
SELECT id, exam_session_id, instance_id, question_id, question_index,
       stage, concept_id, selected_option_id, correct_option_id,
       is_correct, error_tag, grading_version, answered_at
FROM exam_attempts
WHERE exam_session_id = ?
ORDER BY answered_at ASC;`

	rows, err := d.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed querying exam attempts: %w", err)
	}
	defer rows.Close()

	var attempts []exam.ExamAttemptRecord
	for rows.Next() {
		var att exam.ExamAttemptRecord
		var stageStr string
		var isCorrInt int
		var errTag sql.NullString

		err := rows.Scan(
			&att.ID, &att.ExamSessionID, &att.InstanceID, &att.QuestionID, &att.QuestionIndex,
			&stageStr, &att.ConceptID, &att.SelectedOptionID, &att.CorrectOptionID,
			&isCorrInt, &errTag, &att.GradingVersion, &att.AnsweredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning exam attempt: %w", err)
		}
		att.Stage = domain.DrillStage(stageStr)
		att.IsCorrect = (isCorrInt == 1)
		if errTag.Valid {
			att.ErrorTag = errTag.String
		}
		attempts = append(attempts, att)
	}
	return attempts, nil
}

// GetExamQuestionSnapshots returns original questions in insertion (exam) order.
func (d *DB) GetExamQuestionSnapshots(sessionID string) ([]*domain.QuestionInstance, error) {
	rows, err := d.db.Query("SELECT id FROM question_instances WHERE session_id = ? ORDER BY rowid", sessionID)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	var instances []*domain.QuestionInstance
	for _, id := range ids {
		inst, err := d.GetQuestionInstance(id)
		if err != nil {
			return nil, err
		}
		instances = append(instances, inst)
	}
	return instances, nil
}

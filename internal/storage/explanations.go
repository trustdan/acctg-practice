package storage

import (
	"fmt"
	"time"
)

// SavedExplanation is advisory prose selected by the learner, never bank content.
type SavedExplanation struct {
	ID              string
	InstanceID      string
	QuestionID      string
	QuestionVersion int
	Stage           string
	Scenario        string
	StagePrompt     string
	Explanation     string
	Provider        string
	GeneratedAt     time.Time
	SavedAt         time.Time
}

func (d *DB) SaveExplanation(e SavedExplanation) error {
	if e.ID == "" || e.Explanation == "" {
		return fmt.Errorf("explanation ID and text are required")
	}
	_, err := d.db.Exec(`INSERT INTO saved_explanations
(id, instance_id, question_id, question_version, stage, scenario, stage_prompt, explanation, provider, generated_at, saved_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO NOTHING`,
		e.ID, e.InstanceID, e.QuestionID, e.QuestionVersion, e.Stage, e.Scenario, e.StagePrompt,
		e.Explanation, e.Provider, e.GeneratedAt.UTC().Format(time.RFC3339Nano), e.SavedAt.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("save explanation: %w", err)
	}
	return nil
}

func (d *DB) ListSavedExplanations() ([]SavedExplanation, error) {
	rows, err := d.db.Query(`SELECT id, instance_id, question_id, question_version, stage, scenario, stage_prompt,
explanation, provider, generated_at, saved_at FROM saved_explanations ORDER BY saved_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list saved explanations: %w", err)
	}
	defer rows.Close()
	var results []SavedExplanation
	for rows.Next() {
		var e SavedExplanation
		var generated, saved string
		if err := rows.Scan(&e.ID, &e.InstanceID, &e.QuestionID, &e.QuestionVersion, &e.Stage, &e.Scenario,
			&e.StagePrompt, &e.Explanation, &e.Provider, &generated, &saved); err != nil {
			return nil, err
		}
		if e.GeneratedAt, err = time.Parse(time.RFC3339Nano, generated); err != nil {
			return nil, err
		}
		if e.SavedAt, err = time.Parse(time.RFC3339Nano, saved); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, rows.Err()
}

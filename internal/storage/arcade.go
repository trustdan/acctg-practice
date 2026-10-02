package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// ArcadeHighScore represents a high-score achievement in the startup spaceship arcade game.
type ArcadeHighScore struct {
	ID              string    `json:"id"`
	Initials        string    `json:"initials"`
	Score           int       `json:"score"`
	BlastedCount    int       `json:"blasted_count"`
	SurvivalSeconds int       `json:"survival_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}

// SaveArcadeHighScore inserts a new high score entry into the database.
func (d *DB) SaveArcadeHighScore(entry ArcadeHighScore) error {
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("score-%d", time.Now().UnixNano())
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.Initials == "" {
		entry.Initials = "AAA"
	}
	if len(entry.Initials) > 3 {
		entry.Initials = entry.Initials[:3]
	}

	query := `
INSERT INTO arcade_high_scores (id, initials, score, blasted_count, survival_seconds, created_at)
VALUES (?, ?, ?, ?, ?, ?);`

	_, err := d.db.Exec(query,
		entry.ID,
		entry.Initials,
		entry.Score,
		entry.BlastedCount,
		entry.SurvivalSeconds,
		entry.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to save arcade high score: %w", err)
	}
	return nil
}

// GetTopArcadeHighScore returns the highest score recorded in the arcade game, or nil if none exists.
func (d *DB) GetTopArcadeHighScore() (*ArcadeHighScore, error) {
	query := `
SELECT id, initials, score, blasted_count, survival_seconds, created_at
FROM arcade_high_scores
ORDER BY score DESC, created_at ASC
LIMIT 1;`

	var h ArcadeHighScore
	var createdStr string
	err := d.db.QueryRow(query).Scan(
		&h.ID,
		&h.Initials,
		&h.Score,
		&h.BlastedCount,
		&h.SurvivalSeconds,
		&createdStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query top arcade high score: %w", err)
	}

	if t, parseErr := time.Parse(time.RFC3339, createdStr); parseErr == nil {
		h.CreatedAt = t
	} else if t, parseErr := time.Parse("2006-01-02 15:04:05.999999999-07:00", createdStr); parseErr == nil {
		h.CreatedAt = t
	} else if t, parseErr := time.Parse("2006-01-02 15:04:05", createdStr); parseErr == nil {
		h.CreatedAt = t
	}

	return &h, nil
}

// ListTopArcadeHighScores returns the top N high scores ordered by score descending.
func (d *DB) ListTopArcadeHighScores(limit int) ([]ArcadeHighScore, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
SELECT id, initials, score, blasted_count, survival_seconds, created_at
FROM arcade_high_scores
ORDER BY score DESC, created_at ASC
LIMIT ?;`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list arcade high scores: %w", err)
	}
	defer rows.Close()

	var results []ArcadeHighScore
	for rows.Next() {
		var h ArcadeHighScore
		var createdStr string
		if err := rows.Scan(
			&h.ID,
			&h.Initials,
			&h.Score,
			&h.BlastedCount,
			&h.SurvivalSeconds,
			&createdStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan arcade high score row: %w", err)
		}

		if t, parseErr := time.Parse(time.RFC3339, createdStr); parseErr == nil {
			h.CreatedAt = t
		} else if t, parseErr := time.Parse("2006-01-02 15:04:05.999999999-07:00", createdStr); parseErr == nil {
			h.CreatedAt = t
		} else if t, parseErr := time.Parse("2006-01-02 15:04:05", createdStr); parseErr == nil {
			h.CreatedAt = t
		}

		results = append(results, h)
	}

	return results, nil
}

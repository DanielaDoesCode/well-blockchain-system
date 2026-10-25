package db

import (
	"context"
	"fmt"
	"time"
)

type Contribution struct {
	ID             string
	Commitment     string
	Nullifier      string
	MerkleRoot     string
	PA0            string
	PA1            string
	PB00           string
	PB01           string
	PB10           string
	PB11           string
	PC0            string
	PC1            string
	ConditionGroup string
	Status         string
	CreatedAt      time.Time
	SubmittedAt    *time.Time
}

type Counter struct {
	ConditionGroup string `json:"condition_group"`
	TotalCount     int64  `json:"total_count"`
	PublishedCount int64  `json:"published_count"`
}

func (db *DB) InsertPendingContribution(ctx context.Context, c *Contribution) error {
	query := `
		INSERT INTO pending_contributions (
			commitment, nullifier, merkle_root,
			pA_0, pA_1,
			pB_0_0, pB_0_1, pB_1_0, pB_1_1,
			pC_0, pC_1,
			condition_group, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`
	_, err := db.Pool.Exec(ctx, query,
		c.Commitment, c.Nullifier, c.MerkleRoot,
		c.PA0, c.PA1,
		c.PB00, c.PB01, c.PB10, c.PB11,
		c.PC0, c.PC1,
		c.ConditionGroup, "pending",
	)
	if err != nil {
		return fmt.Errorf("insert pending contribution: %w", err)
	}
	return nil
}

func (db *DB) PendingCountForGroup(ctx context.Context, group string) (int64, error) {
	var count int64
	err := db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM pending_contributions WHERE condition_group = $1 AND status = 'pending'",
		group,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count pending: %w", err)
	}
	return count, nil
}

func (db *DB) GetPendingContributions(ctx context.Context, group string) ([]*Contribution, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, commitment, nullifier, merkle_root,
		        pA_0, pA_1,
		        pB_0_0, pB_0_1, pB_1_0, pB_1_1,
		        pC_0, pC_1,
		        condition_group, status, created_at, submitted_at
		 FROM pending_contributions
		 WHERE condition_group = $1 AND status = 'pending'
		 ORDER BY created_at ASC`,
		group,
	)
	if err != nil {
		return nil, fmt.Errorf("query pending: %w", err)
	}
	defer rows.Close()

	var out []*Contribution
	for rows.Next() {
		c := &Contribution{}
		err := rows.Scan(
			&c.ID, &c.Commitment, &c.Nullifier, &c.MerkleRoot,
			&c.PA0, &c.PA1,
			&c.PB00, &c.PB01, &c.PB10, &c.PB11,
			&c.PC0, &c.PC1,
			&c.ConditionGroup, &c.Status, &c.CreatedAt, &c.SubmittedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan pending: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}

func (db *DB) MarkContributionSubmitted(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := db.Pool.Exec(ctx,
		"UPDATE pending_contributions SET status = 'submitted', submitted_at = $1 WHERE id = $2",
		now, id,
	)
	if err != nil {
		return fmt.Errorf("mark submitted: %w", err)
	}
	return nil
}

func (db *DB) IncrementCounter(ctx context.Context, group string) error {
	_, err := db.Pool.Exec(ctx,
		`INSERT INTO research_counters (condition_group, total_count, published_count, updated_at)
		 VALUES ($1, 1, 0, NOW())
		 ON CONFLICT (condition_group)
		 DO UPDATE SET total_count = research_counters.total_count + 1, updated_at = NOW()`,
		group,
	)
	if err != nil {
		return fmt.Errorf("increment counter: %w", err)
	}
	return nil
}

func (db *DB) IncrementPublishedCount(ctx context.Context, group string, n int64) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE research_counters
		 SET published_count = published_count + $1, updated_at = NOW()
		 WHERE condition_group = $2`,
		n, group,
	)
	if err != nil {
		return fmt.Errorf("increment published: %w", err)
	}
	return nil
}

func (db *DB) GetCounters(ctx context.Context) ([]*Counter, error) {
	rows, err := db.Pool.Query(ctx,
		"SELECT condition_group, total_count, published_count FROM research_counters ORDER BY condition_group",
	)
	if err != nil {
		return nil, fmt.Errorf("query counters: %w", err)
	}
	defer rows.Close()

	var out []*Counter
	for rows.Next() {
		c := &Counter{}
		if err := rows.Scan(&c.ConditionGroup, &c.TotalCount, &c.PublishedCount); err != nil {
			return nil, fmt.Errorf("scan counter: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}

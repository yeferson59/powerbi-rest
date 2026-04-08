package metrics

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) List(ctx context.Context) ([]Record, error) {
	rows, err := s.db.Query(ctx, "SELECT id, request_id, route, method, complexity, n_param, response_ms, status_code, error_message, created_at FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("query metrics: %w", err)
	}
	defer rows.Close()

	records := make([]Record, 0)
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.RequestID, &record.Route, &record.Method, &record.Complexity, &record.NParam, &record.ResponseMs, &record.StatusCode, &record.ErrorMessage, &record.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan metric row: %w", err)
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metric rows: %w", err)
	}

	return records, nil
}

func (s *PostgresStore) Create(ctx context.Context, metric CreateInput) error {
	_, err := s.db.Exec(ctx, "INSERT INTO metrics (route, method, complexity, n_param, response_ms, status_code, error_message, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", metric.Route, metric.Method, metric.Complexity, metric.NParam, metric.ResponseMs, metric.StatusCode, metric.ErrorMsg, metric.RequestedAt)
	if err != nil {
		return fmt.Errorf("insert metric: %w", err)
	}

	return nil
}

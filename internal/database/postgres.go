package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	databaseURL string
}

func NewPostgresDB(databaseURL string) *PostgresDB {
	return &PostgresDB{
		databaseURL: databaseURL,
	}
}

func (db *PostgresDB) Connect(ctx context.Context) *pgxpool.Pool {
	conn, err := pgxpool.New(ctx, db.databaseURL)
	if err != nil {
		log.Fatal("failed to connect to database: %w", err)
	}

	return conn
}

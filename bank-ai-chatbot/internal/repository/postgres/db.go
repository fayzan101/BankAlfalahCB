package postgres

import (
	"context"
	"fmt"

	"bank-ai-chatbot/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	pool, err := database.NewPostgresPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("init postgres db: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	if d == nil || d.Pool == nil {
		return
	}
	d.Pool.Close()
}

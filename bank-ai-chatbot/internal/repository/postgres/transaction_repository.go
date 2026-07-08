package postgres

import (
	"context"
	"errors"
	"fmt"

	"bank-ai-chatbot/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TransactionRepository struct {
	db *DB
}

func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, amount, type, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, tx.ID, tx.UserID, tx.Amount, tx.Type, tx.Description).
		Scan(&tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, description, created_at
		FROM transactions
		WHERE id = $1
	`
	var tx models.Transaction
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Description, &tx.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get transaction by id: %w", err)
	}
	return &tx, nil
}

func (r *TransactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
		SELECT id, user_id, amount, type, description, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Description, &tx.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'credit' THEN amount
				WHEN type = 'debit' THEN -amount
				ELSE 0
			END
		), 0)
		FROM transactions
		WHERE user_id = $1
	`
	var balance float64
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}
	return nil
}

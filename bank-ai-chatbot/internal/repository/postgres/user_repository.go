package postgres

import (
	"context"
	"errors"
	"fmt"

	"bank-ai-chatbot/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, full_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, user.ID, user.FullName, user.Email, user.PasswordHash).
		Scan(&user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

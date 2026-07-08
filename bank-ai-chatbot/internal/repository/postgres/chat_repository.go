package postgres

import (
	"context"
	"errors"
	"fmt"

	"bank-ai-chatbot/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChatRepository struct {
	db *DB
}

func NewChatRepository(db *DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	query := `
		INSERT INTO chats (id, user_id, title)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, chat.ID, chat.UserID, chat.Title).Scan(&chat.CreatedAt)
	if err != nil {
		return fmt.Errorf("create chat: %w", err)
	}
	return nil
}

func (r *ChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	query := `
		SELECT id, user_id, title, created_at
		FROM chats
		WHERE id = $1
	`
	var chat models.Chat
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&chat.ID, &chat.UserID, &chat.Title, &chat.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get chat by id: %w", err)
	}
	return &chat, nil
}

func (r *ChatRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Chat, error) {
	query := `
		SELECT id, user_id, title, created_at
		FROM chats
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list chats: %w", err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.UserID, &chat.Title, &chat.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan chat: %w", err)
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chats: %w", err)
	}
	return chats, nil
}

func (r *ChatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM chats WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}
	return nil
}

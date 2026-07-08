package postgres

import (
	"context"
	"errors"
	"fmt"

	"bank-ai-chatbot/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	db *DB
}

func NewMessageRepository(db *DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) error {
	query := `
		INSERT INTO messages (id, chat_id, sender_type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, msg.ID, msg.ChatID, msg.SenderType, msg.Content).
		Scan(&msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, chat_id, sender_type, content, created_at
		FROM messages
		WHERE id = $1
	`
	var msg models.Message
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.ChatID, &msg.SenderType, &msg.Content, &msg.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get message by id: %w", err)
	}
	return &msg, nil
}

func (r *MessageRepository) ListByChatID(ctx context.Context, chatID uuid.UUID) ([]models.Message, error) {
	query := `
		SELECT id, chat_id, sender_type, content, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Pool.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderType, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}
	return messages, nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM messages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	return nil
}

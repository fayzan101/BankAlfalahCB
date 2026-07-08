package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SenderUser      = "user"
	SenderAssistant = "assistant"
	SenderSystem    = "system"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	ChatID     uuid.UUID `json:"chat_id"`
	SenderType string    `json:"sender_type"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

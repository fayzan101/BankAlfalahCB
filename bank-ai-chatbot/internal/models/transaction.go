package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TransactionCredit = "credit"
	TransactionDebit  = "debit"
)

type Transaction struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

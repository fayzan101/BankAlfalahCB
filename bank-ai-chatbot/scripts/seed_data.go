package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"bank-ai-chatbot/internal/security"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	userID := uuid.New()
	email := fmt.Sprintf("demo+%s@bank.local", userID.String()[:8])

	hash, err := security.HashPassword("DemoPass123")
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, full_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
	`, userID, "Demo Customer", email, hash)
	if err != nil {
		log.Fatalf("seed users failed: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO transactions (user_id, amount, type, description)
		VALUES
		($1, 5000.00, 'credit', 'Initial salary deposit'),
		($1, 1200.50, 'debit', 'ATM withdrawal'),
		($1, 230.75, 'debit', 'Utility bill payment')
	`, userID)
	if err != nil {
		log.Fatalf("seed transactions failed: %v", err)
	}

	log.Printf("seed complete: user_id=%s email=%s password=DemoPass123", userID.String(), email)
}

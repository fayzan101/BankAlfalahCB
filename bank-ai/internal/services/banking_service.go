package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bank-ai-chatbot/internal/dto"
	"bank-ai-chatbot/internal/models"
	"bank-ai-chatbot/internal/repository/postgres"
	apperrors "bank-ai-chatbot/pkg/errors"
	"github.com/google/uuid"
)

const defaultCurrency = "PKR"

type BankingService struct {
	transactions *postgres.TransactionRepository
}

func NewBankingService(transactions *postgres.TransactionRepository) *BankingService {
	return &BankingService{transactions: transactions}
}

func (s *BankingService) GetBalance(ctx context.Context, userID uuid.UUID) (*dto.BalanceResponse, error) {
	balance, err := s.transactions.GetBalance(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("failed to fetch balance", err)
	}

	return &dto.BalanceResponse{
		Balance:  balance,
		Currency: defaultCurrency,
	}, nil
}

func (s *BankingService) GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) (*dto.TransactionsResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	transactions, err := s.transactions.ListByUserID(ctx, userID, limit)
	if err != nil {
		return nil, apperrors.Internal("failed to fetch transactions", err)
	}

	items := make([]dto.TransactionItem, 0, len(transactions))
	for _, tx := range transactions {
		items = append(items, toTransactionItem(tx))
	}

	return &dto.TransactionsResponse{
		Transactions: items,
		Count:        len(items),
	}, nil
}

func (s *BankingService) ReplyForIntent(ctx context.Context, userID uuid.UUID, intent Intent) (string, error) {
	switch intent {
	case IntentBalance:
		result, err := s.GetBalance(ctx, userID)
		if err != nil {
			return "", err
		}
		return FormatBalanceReply(result.Balance, result.Currency), nil
	case IntentTransactions:
		result, err := s.GetRecentTransactions(ctx, userID, 10)
		if err != nil {
			return "", err
		}
		return FormatTransactionsReply(result.Transactions, defaultCurrency), nil
	default:
		return "", apperrors.BadRequest("unsupported banking intent")
	}
}

func FormatBalanceReply(balance float64, currency string) string {
	return fmt.Sprintf("Your current account balance is %s %s.", currency, formatAmount(balance))
}

func FormatTransactionsReply(transactions []dto.TransactionItem, currency string) string {
	if len(transactions) == 0 {
		return "You have no recent transactions on your account."
	}

	var builder strings.Builder
	builder.WriteString("Here are your recent transactions:\n")
	for i, tx := range transactions {
		sign := "+"
		if tx.Type == models.TransactionDebit {
			sign = "-"
		}
		builder.WriteString(fmt.Sprintf(
			"%d. %s %s%s — %s (%s)\n",
			i+1,
			tx.CreatedAt,
			sign,
			formatAmount(tx.Amount),
			tx.Description,
			currency,
		))
	}
	return strings.TrimSpace(builder.String())
}

func toTransactionItem(tx models.Transaction) dto.TransactionItem {
	return dto.TransactionItem{
		ID:          tx.ID.String(),
		Amount:      tx.Amount,
		Type:        tx.Type,
		Description: tx.Description,
		CreatedAt:   tx.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

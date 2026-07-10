package audit

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return &Logger{logger: logger}
}

func (l *Logger) AuthRegister(userID uuid.UUID, email, requestID, ip string, success bool) {
	l.log("auth.register", userID, email, requestID, ip, success, nil)
}

func (l *Logger) AuthLogin(userID uuid.UUID, email, requestID, ip string, success bool) {
	l.log("auth.login", userID, email, requestID, ip, success, nil)
}

func (l *Logger) BankingBalance(userID uuid.UUID, requestID, ip string) {
	l.log("banking.balance", userID, "", requestID, ip, true, nil)
}

func (l *Logger) BankingTransactions(userID uuid.UUID, requestID, ip string, limit int) {
	l.log("banking.transactions", userID, "", requestID, ip, true, map[string]any{
		"limit": limit,
	})
}

func (l *Logger) log(action string, userID uuid.UUID, email, requestID, ip string, success bool, extra map[string]any) {
	attrs := []any{
		"action", action,
		"success", success,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	}
	if userID != uuid.Nil {
		attrs = append(attrs, "user_id", userID.String())
	}
	if email != "" {
		attrs = append(attrs, "email", email)
	}
	if requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	if ip != "" {
		attrs = append(attrs, "ip", ip)
	}
	for k, v := range extra {
		attrs = append(attrs, k, v)
	}
	l.logger.Info("audit", attrs...)
}

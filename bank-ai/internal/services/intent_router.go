package services

import "strings"

type Intent string

const (
	IntentBalance      Intent = "balance"
	IntentTransactions Intent = "transactions"
	IntentGeneral      Intent = "general"
)

func DetectIntent(message string) Intent {
	normalized := strings.ToLower(strings.TrimSpace(message))

	balanceKeywords := []string{
		"balance",
		"account balance",
		"how much do i have",
		"how much money",
		"what is my balance",
		"what's my balance",
		"check my balance",
		"show my balance",
	}
	for _, keyword := range balanceKeywords {
		if strings.Contains(normalized, keyword) {
			return IntentBalance
		}
	}

	transactionKeywords := []string{
		"transaction",
		"recent payment",
		"recent transfer",
		"recent activity",
		"show recent",
		"latest transaction",
		"recent transaction",
		"payment history",
		"transaction history",
		"show my transactions",
		"list my transactions",
	}
	for _, keyword := range transactionKeywords {
		if strings.Contains(normalized, keyword) {
			return IntentTransactions
		}
	}

	return IntentGeneral
}

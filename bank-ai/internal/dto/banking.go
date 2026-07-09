package dto

type BalanceResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type TransactionItem struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

type TransactionsResponse struct {
	Transactions []TransactionItem `json:"transactions"`
	Count        int               `json:"count"`
}

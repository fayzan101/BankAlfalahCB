package handlers

import (
	"net/http"
	"strconv"

	"bank-ai-chatbot/internal/services"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
)

type BankingHandler struct {
	banking *services.BankingService
}

func NewBankingHandler(banking *services.BankingService) *BankingHandler {
	return &BankingHandler{banking: banking}
}

func (h *BankingHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	result, err := h.banking.GetBalance(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *BankingHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			response.Error(w, apperrors.BadRequest("limit must be a positive integer"))
			return
		}
		limit = parsed
	}

	result, err := h.banking.GetRecentTransactions(r.Context(), userID, limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

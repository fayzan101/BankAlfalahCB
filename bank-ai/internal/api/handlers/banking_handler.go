package handlers

import (
	"net/http"
	"strconv"

	"bank-ai-chatbot/internal/audit"
	"bank-ai-chatbot/internal/services"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type BankingHandler struct {
	banking *services.BankingService
	audit   *audit.Logger
}

func NewBankingHandler(banking *services.BankingService, auditLogger *audit.Logger) *BankingHandler {
	return &BankingHandler{banking: banking, audit: auditLogger}
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

	h.audit.BankingBalance(userID, chimiddleware.GetReqID(r.Context()), r.RemoteAddr)
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

	h.audit.BankingTransactions(userID, chimiddleware.GetReqID(r.Context()), r.RemoteAddr, limit)
	response.JSON(w, http.StatusOK, result)
}

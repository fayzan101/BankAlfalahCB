package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"bank-ai-chatbot/internal/audit"
	"bank-ai-chatbot/internal/dto"
	"bank-ai-chatbot/internal/services"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type AuthHandler struct {
	auth  *services.AuthService
	audit *audit.Logger
}

func NewAuthHandler(auth *services.AuthService, auditLogger *audit.Logger) *AuthHandler {
	return &AuthHandler{auth: auth, audit: auditLogger}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.BadRequest("invalid request body"))
		return
	}

	result, err := h.auth.Register(r.Context(), req)
	requestID := chimiddleware.GetReqID(r.Context())
	if err != nil {
		h.audit.AuthRegister(uuid.Nil, req.Email, requestID, r.RemoteAddr, false)
		response.Error(w, err)
		return
	}

	userID, _ := uuid.Parse(result.User.ID)
	h.audit.AuthRegister(userID, result.User.Email, requestID, r.RemoteAddr, true)
	response.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.BadRequest("invalid request body"))
		return
	}

	result, err := h.auth.Login(r.Context(), req)
	requestID := chimiddleware.GetReqID(r.Context())
	if err != nil {
		h.audit.AuthLogin(uuid.Nil, req.Email, requestID, r.RemoteAddr, false)
		response.Error(w, err)
		return
	}

	userID, _ := uuid.Parse(result.User.ID)
	h.audit.AuthLogin(userID, result.User.Email, requestID, r.RemoteAddr, true)
	response.JSON(w, http.StatusOK, result)
}

type MeHandler struct {
	auth *services.AuthService
}

func NewMeHandler(auth *services.AuthService) *MeHandler {
	return &MeHandler{auth: auth}
}

func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	user, err := h.auth.GetUser(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok && userID != uuid.Nil
}

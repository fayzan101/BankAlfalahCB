package handlers

import (
	"encoding/json"
	"net/http"

	"bank-ai-chatbot/internal/dto"
	"bank-ai-chatbot/internal/services"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chat *services.ChatService
}

func NewChatHandler(chat *services.ChatService) *ChatHandler {
	return &ChatHandler{chat: chat}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	var req dto.CreateChatRequest
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperrors.BadRequest("invalid request body"))
			return
		}
	}

	result, err := h.chat.CreateChat(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	chatID, err := parseChatID(chi.URLParam(r, "chat_id"))
	if err != nil {
		response.Error(w, err)
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.BadRequest("invalid request body"))
		return
	}

	result, err := h.chat.SendMessage(r.Context(), userID, chatID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.Unauthorized("authentication required"))
		return
	}

	chatID, err := parseChatID(chi.URLParam(r, "chat_id"))
	if err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.chat.GetHistory(r.Context(), userID, chatID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func parseChatID(raw string) (uuid.UUID, error) {
	if raw == "" {
		return uuid.Nil, apperrors.BadRequest("chat id is required")
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperrors.BadRequest("invalid chat id")
	}
	return id, nil
}

package services

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"bank-ai-chatbot/internal/dto"
	"bank-ai-chatbot/internal/models"
	"bank-ai-chatbot/internal/repository/postgres"
	apperrors "bank-ai-chatbot/pkg/errors"
	"github.com/google/uuid"
)

const (
	defaultChatTitle   = "New Chat"
	maxMessageLength   = 4000
	maxChatTitleLength = 100
)

type ChatService struct {
	chats     *postgres.ChatRepository
	messages  *postgres.MessageRepository
	llm       *LLMService
	llmEnabled bool
}

func NewChatService(chats *postgres.ChatRepository, messages *postgres.MessageRepository, llm *LLMService, llmEnabled bool) *ChatService {
	return &ChatService{
		chats:      chats,
		messages:   messages,
		llm:        llm,
		llmEnabled: llmEnabled,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, userID uuid.UUID, req dto.CreateChatRequest) (*dto.CreateChatResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = defaultChatTitle
	}
	if utf8.RuneCountInString(title) > maxChatTitleLength {
		return nil, apperrors.BadRequest("chat title must be at most 100 characters")
	}

	chat := &models.Chat{
		ID:     uuid.New(),
		UserID: userID,
		Title:  title,
	}

	if err := s.chats.Create(ctx, chat); err != nil {
		return nil, apperrors.Internal("failed to create chat", err)
	}

	return &dto.CreateChatResponse{
		Chat: toChatSummary(*chat),
	}, nil
}

func (s *ChatService) SendMessage(ctx context.Context, userID, chatID uuid.UUID, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, apperrors.BadRequest("message content is required")
	}
	if utf8.RuneCountInString(content) > maxMessageLength {
		return nil, apperrors.BadRequest("message content must be at most 4000 characters")
	}

	if _, err := s.getOwnedChat(ctx, userID, chatID); err != nil {
		return nil, err
	}

	userMessage := &models.Message{
		ID:         uuid.New(),
		ChatID:     chatID,
		SenderType: models.SenderUser,
		Content:    content,
	}
	if err := s.messages.Create(ctx, userMessage); err != nil {
		return nil, apperrors.Internal("failed to save user message", err)
	}

	if !s.llmEnabled || s.llm == nil {
		return nil, apperrors.ServiceUnavailable("assistant is temporarily unavailable")
	}

	history, err := s.messages.ListByChatID(ctx, chatID)
	if err != nil {
		return nil, apperrors.Internal("failed to load chat history", err)
	}

	reply, err := s.llm.GenerateReply(ctx, history)
	if err != nil {
		return nil, apperrors.ServiceUnavailable("assistant is temporarily unavailable")
	}

	assistantMessage := &models.Message{
		ID:         uuid.New(),
		ChatID:     chatID,
		SenderType: models.SenderAssistant,
		Content:    reply,
	}
	if err := s.messages.Create(ctx, assistantMessage); err != nil {
		return nil, apperrors.Internal("failed to save assistant message", err)
	}

	return &dto.SendMessageResponse{
		UserMessage:      toMessageSummary(*userMessage),
		AssistantMessage: toMessageSummary(*assistantMessage),
	}, nil
}

func (s *ChatService) GetHistory(ctx context.Context, userID, chatID uuid.UUID) (*dto.ChatHistoryResponse, error) {
	chat, err := s.getOwnedChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	messages, err := s.messages.ListByChatID(ctx, chatID)
	if err != nil {
		return nil, apperrors.Internal("failed to load chat history", err)
	}

	summaries := make([]dto.MessageSummary, 0, len(messages))
	for _, msg := range messages {
		summaries = append(summaries, toMessageSummary(msg))
	}

	return &dto.ChatHistoryResponse{
		Chat:     toChatSummary(*chat),
		Messages: summaries,
	}, nil
}

func (s *ChatService) getOwnedChat(ctx context.Context, userID, chatID uuid.UUID) (*models.Chat, error) {
	chat, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		return nil, apperrors.Internal("failed to lookup chat", err)
	}
	if chat == nil {
		return nil, apperrors.NotFound("chat not found")
	}
	if chat.UserID != userID {
		return nil, apperrors.Forbidden("you do not have access to this chat")
	}
	return chat, nil
}

func toChatSummary(chat models.Chat) dto.ChatSummary {
	return dto.ChatSummary{
		ID:        chat.ID.String(),
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toMessageSummary(msg models.Message) dto.MessageSummary {
	return dto.MessageSummary{
		ID:         msg.ID.String(),
		SenderType: msg.SenderType,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt.UTC().Format(time.RFC3339),
	}
}

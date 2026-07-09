package services

import (
	"context"

	"bank-ai-chatbot/internal/models"
)

type LLMGenerator interface {
	GenerateReply(ctx context.Context, history []models.Message, maxHistory int) (string, error)
}

type LLMService struct {
	client     LLMGenerator
	maxHistory int
}

func NewLLMService(client LLMGenerator, maxHistory int) *LLMService {
	return &LLMService{
		client:     client,
		maxHistory: maxHistory,
	}
}

func (s *LLMService) GenerateReply(ctx context.Context, history []models.Message) (string, error) {
	return s.client.GenerateReply(ctx, history, s.maxHistory)
}

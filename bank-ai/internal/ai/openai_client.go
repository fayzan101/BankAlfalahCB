package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bank-ai-chatbot/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client    *openai.Client
	model     string
	maxTokens int
}

func NewClient(apiKey, model string, maxTokens int, timeout time.Duration) *Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.HTTPClient = &http.Client{Timeout: timeout}
	return &Client{
		client:    openai.NewClientWithConfig(cfg),
		model:     model,
		maxTokens: maxTokens,
	}
}

func (c *Client) GenerateReply(ctx context.Context, history []models.Message, maxHistory int) (string, error) {
	messages := BuildChatMessages(history, maxHistory)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: c.maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("openai returned empty content")
	}

	return content, nil
}

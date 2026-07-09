package ai

import (
	"bank-ai-chatbot/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

func BuildChatMessages(history []models.Message, maxHistory int) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: BankingSystemPrompt,
		},
	}

	start := 0
	if maxHistory > 0 && len(history) > maxHistory {
		start = len(history) - maxHistory
	}

	for _, msg := range history[start:] {
		role := openai.ChatMessageRoleUser
		switch msg.SenderType {
		case models.SenderAssistant:
			role = openai.ChatMessageRoleAssistant
		case models.SenderSystem:
			role = openai.ChatMessageRoleSystem
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	return messages
}

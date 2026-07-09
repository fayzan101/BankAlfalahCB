package dto

type CreateChatRequest struct {
	Title string `json:"title"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type ChatSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type MessageSummary struct {
	ID         string `json:"id"`
	SenderType string `json:"sender_type"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

type CreateChatResponse struct {
	Chat ChatSummary `json:"chat"`
}

type SendMessageResponse struct {
	UserMessage      MessageSummary `json:"user_message"`
	AssistantMessage MessageSummary `json:"assistant_message"`
}

type ChatHistoryResponse struct {
	Chat     ChatSummary      `json:"chat"`
	Messages []MessageSummary `json:"messages"`
}

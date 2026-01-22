package httpapi

import (
	"context"

	"hirifyGOTest/internal/domain/chat"
)

type ChatService interface {
	CreateChat(ctx context.Context, title string) (chat.Chat, error)
	CreateMessage(ctx context.Context, chatID int64, text string) (chat.Message, error)
	GetChatWithLastMessages(ctx context.Context, chatID int64, limit int) (chat.Chat, []chat.Message, error)
	DeleteChat(ctx context.Context, chatID int64) error
}

type createChatRequest struct {
	Title string `json:"title"`
}

type createMessageRequest struct {
	Text string `json:"text"`
}

type getChatResponse struct {
	Chat     chat.Chat      `json:"chat"`
	Messages []chat.Message `json:"messages"`
}

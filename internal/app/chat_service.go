package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"hirifyGOTest/internal/domain/chat"
)

type ChatService struct {
	chats    chat.ChatRepository
	messages chat.MessageRepository
}

func NewChatService(chats chat.ChatRepository, messages chat.MessageRepository) *ChatService {
	return &ChatService{chats: chats, messages: messages}
}

func (s *ChatService) CreateChat(ctx context.Context, title string) (chat.Chat, error) {
	title = strings.TrimSpace(title)
	if len(title) < 1 || len(title) > 200 {
		return chat.Chat{}, fmt.Errorf("%w: title must be 1..200 chars", chat.ErrValidation)
	}
	return s.chats.Create(ctx, title)
}

func (s *ChatService) CreateMessage(ctx context.Context, chatID int64, text string) (chat.Message, error) {
	text = strings.TrimSpace(text)
	if len(text) < 1 || len(text) > 5000 {
		return chat.Message{}, fmt.Errorf("%w: text must be 1..5000 chars", chat.ErrValidation)
	}

	if _, err := s.chats.GetByID(ctx, chatID); err != nil {
		if errors.Is(err, chat.ErrNotFound) {
			return chat.Message{}, err
		}
		return chat.Message{}, fmt.Errorf("get chat: %w", err)
	}

	return s.messages.Create(ctx, chatID, text)
}

func (s *ChatService) GetChatWithLastMessages(ctx context.Context, chatID int64, limit int) (chat.Chat, []chat.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	c, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		return chat.Chat{}, nil, err
	}

	msgs, err := s.messages.ListLastByChatID(ctx, chatID, limit)
	if err != nil {
		return chat.Chat{}, nil, err
	}

	return c, msgs, nil
}

func (s *ChatService) DeleteChat(ctx context.Context, chatID int64) error {
	deleted, err := s.chats.DeleteByID(ctx, chatID)
	if err != nil {
		return err
	}
	if !deleted {
		return chat.ErrNotFound
	}
	return nil
}

package postgres

import (
	"context"
	"fmt"

	"hirifyGOTest/internal/domain/chat"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, chatID int64, text string) (chat.Message, error) {
	m := MessageModel{ChatID: chatID, Text: text}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return chat.Message{}, fmt.Errorf("create message: %w", err)
	}
	return chat.Message{ID: m.ID, ChatID: m.ChatID, Text: m.Text, CreatedAt: m.CreatedAt}, nil
}

func (r *MessageRepository) ListLastByChatID(ctx context.Context, chatID int64, limit int) ([]chat.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// селект
	// SELECT * FROM (
	//   SELECT * FROM messages WHERE chat_id=? ORDER BY created_at DESC, id DESC LIMIT ?
	// ) t ORDER BY created_at ASC, id ASC;
	sub := r.db.WithContext(ctx).
		Model(&MessageModel{}).
		Where("chat_id = ?", chatID).
		Order("created_at desc, id desc").
		Limit(limit)

	var models []MessageModel
	err := r.db.WithContext(ctx).
		Table("(?) as t", sub).
		Order("created_at asc, id asc").
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}

	out := make([]chat.Message, 0, len(models))
	for _, m := range models {
		out = append(out, chat.Message{ID: m.ID, ChatID: m.ChatID, Text: m.Text, CreatedAt: m.CreatedAt})
	}
	return out, nil
}

var _ chat.MessageRepository = (*MessageRepository)(nil)

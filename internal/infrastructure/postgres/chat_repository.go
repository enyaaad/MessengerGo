package postgres

import (
	"context"
	"errors"
	"fmt"

	"hirifyGOTest/internal/domain/chat"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, title string) (chat.Chat, error) {
	m := ChatModel{Title: title}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return chat.Chat{}, fmt.Errorf("create chat: %w", err)
	}
	return chat.Chat{ID: m.ID, Title: m.Title, CreatedAt: m.CreatedAt}, nil
}

func (r *ChatRepository) GetByID(ctx context.Context, id int64) (chat.Chat, error) {
	var m ChatModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chat.Chat{}, chat.ErrNotFound
		}
		return chat.Chat{}, fmt.Errorf("get chat: %w", err)
	}
	return chat.Chat{ID: m.ID, Title: m.Title, CreatedAt: m.CreatedAt}, nil
}

func (r *ChatRepository) DeleteByID(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Delete(&ChatModel{}, "id = ?", id)
	if res.Error != nil {
		return false, fmt.Errorf("delete chat: %w", res.Error)
	}
	return res.RowsAffected > 0, nil
}

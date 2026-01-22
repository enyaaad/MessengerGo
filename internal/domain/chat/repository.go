package chat

import "context"

type ChatRepository interface {
	Create(ctx context.Context, title string) (Chat, error)
	GetByID(ctx context.Context, id int64) (Chat, error)
	DeleteByID(ctx context.Context, id int64) (deleted bool, err error)
}

type MessageRepository interface {
	Create(ctx context.Context, chatID int64, text string) (Message, error)
	ListLastByChatID(ctx context.Context, chatID int64, limit int) ([]Message, error)
}

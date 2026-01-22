package postgres

import "time"

type ChatModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (ChatModel) TableName() string { return "chats" }

type MessageModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	ChatID    int64     `gorm:"not null;index:idx_messages_chat_created_at,priority:1"`
	Text      string    `gorm:"type:varchar(5000);not null"`
	CreatedAt time.Time `gorm:"not null;index:idx_messages_chat_created_at,priority:2,sort:desc"`
}

func (MessageModel) TableName() string { return "messages" }

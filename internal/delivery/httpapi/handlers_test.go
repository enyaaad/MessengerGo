package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"hirifyGOTest/internal/app"
	"hirifyGOTest/internal/domain/chat"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type memChatRepo struct {
	mu    sync.Mutex
	next  int64
	chats map[int64]chat.Chat
}

func newMemChatRepo() *memChatRepo {
	return &memChatRepo{next: 1, chats: make(map[int64]chat.Chat)}
}

func (r *memChatRepo) Create(_ context.Context, title string) (chat.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.next
	r.next++
	c := chat.Chat{ID: id, Title: title, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	r.chats[id] = c
	return c, nil
}

func (r *memChatRepo) GetByID(_ context.Context, id int64) (chat.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.chats[id]
	if !ok {
		return chat.Chat{}, chat.ErrNotFound
	}
	return c, nil
}

func (r *memChatRepo) DeleteByID(_ context.Context, id int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.chats[id]; !ok {
		return false, nil
	}
	delete(r.chats, id)
	return true, nil
}

type memMsgRepo struct {
	mu       sync.Mutex
	next     int64
	messages map[int64][]chat.Message
}

func newMemMsgRepo() *memMsgRepo {
	return &memMsgRepo{next: 1, messages: make(map[int64][]chat.Message)}
}

func (r *memMsgRepo) Create(_ context.Context, chatID int64, text string) (chat.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.next
	r.next++
	m := chat.Message{ID: id, ChatID: chatID, Text: text, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	r.messages[chatID] = append(r.messages[chatID], m)
	return m, nil
}

func (r *memMsgRepo) ListLastByChatID(_ context.Context, chatID int64, limit int) ([]chat.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	all := r.messages[chatID]
	if limit <= 0 || limit > len(all) {
		limit = len(all)
	}

	out := make([]chat.Message, 0, limit)
	for i := len(all) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, all[i])
	}
	return out, nil
}

func TestCreateChat_TrimsTitle(t *testing.T) {
	logger := zerolog.New(bytes.NewBuffer(nil))

	chatRepo := newMemChatRepo()
	msgRepo := newMemMsgRepo()
	svc := app.NewChatService(chatRepo, msgRepo)
	router := NewRouter(svc, logger)

	body := []byte(`{"title":"  General  "}`)
	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var got chat.Chat
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Equal(t, int64(1), got.ID)
	require.Equal(t, "General", got.Title)
}

func TestCreateMessage_ChatNotFound(t *testing.T) {
	logger := zerolog.New(bytes.NewBuffer(nil))

	chatRepo := newMemChatRepo()
	msgRepo := newMemMsgRepo()
	svc := app.NewChatService(chatRepo, msgRepo)
	router := NewRouter(svc, logger)

	body := []byte(`{"text":"hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/chats/999/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

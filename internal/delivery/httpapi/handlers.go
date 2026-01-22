package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"hirifyGOTest/internal/domain/chat"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handlers struct {
	chatSvc ChatService
	logger  zerolog.Logger
}

func NewHandlers(chatSvc ChatService, logger zerolog.Logger) *Handlers {
	return &Handlers{chatSvc: chatSvc, logger: logger}
}

func (h *Handlers) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req createChatRequest
	if err := readJSON(w, r, &req); err != nil {
		return
	}

	created, err := h.chatSvc.CreateChat(r.Context(), req.Title)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handlers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	chatID, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var req createMessageRequest
	if err := readJSON(w, r, &req); err != nil {
		return
	}

	created, err := h.chatSvc.CreateMessage(r.Context(), chatID, req.Text)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handlers) GetChat(w http.ResponseWriter, r *http.Request) {
	chatID, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			writeError(w, http.StatusBadRequest, "limit must be positive integer")
			return
		}
		if n > 100 {
			n = 100
		}
		limit = n
	}

	c, msgs, err := h.chatSvc.GetChatWithLastMessages(r.Context(), chatID, limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, getChatResponse{Chat: c, Messages: msgs})
}

func (h *Handlers) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	if err := h.chatSvc.DeleteChat(r.Context(), chatID); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	if raw == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, chat.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, chat.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

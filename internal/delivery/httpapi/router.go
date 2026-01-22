package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func NewRouter(chatSvc ChatService, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	h := NewHandlers(chatSvc, logger)

	r.Get("/openapi.yaml", h.OpenAPI)
	r.Get("/swagger", h.SwaggerUI)

	r.Post("/chats", h.CreateChat)
	r.Post("/chats/{id}/messages", h.CreateMessage)
	r.Get("/chats/{id}", h.GetChat)
	r.Delete("/chats/{id}", h.DeleteChat)

	return r
}

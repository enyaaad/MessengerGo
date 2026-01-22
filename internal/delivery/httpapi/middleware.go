package httpapi

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		event := logger.Info()
		if sw.status >= 500 {
			event = logger.Error()
		} else if sw.status >= 400 {
			event = logger.Warn()
		}

		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", sw.status).
			Dur("duration", time.Since(start)).
			Msg("http")
	})
}

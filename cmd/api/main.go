package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messengerTest/internal/app"
	"messengerTest/internal/config"
	httpapi "messengerTest/internal/delivery/httpapi"
	"messengerTest/internal/infrastructure/postgres"
	"messengerTest/pkg/logging"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.FromEnv()
	logger := logging.New(cfg.ServiceName)

	db, err := postgres.OpenGorm(cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("open database")
	}

	chatRepo := postgres.NewChatRepository(db)
	msgRepo := postgres.NewMessageRepository(db)

	chatSvc := app.NewChatService(chatRepo, msgRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", promhttp.Handler())

	apiRouter := httpapi.NewRouter(chatSvc, logger)
	mux.Handle("/", apiRouter)

	handler := httpapi.LoggingMiddleware(logger, httpapi.MetricsMiddleware(mux))

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info().Str("addr", cfg.HTTPAddr).Msg("http server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("http server failed")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Info().Msg("http server stopped")
}

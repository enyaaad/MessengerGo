package main

import (
	"database/sql"
	"fmt"
	"os"

	"messengerTest/internal/config"
	"messengerTest/pkg/logging"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.FromEnv()
	logger := logging.New(cfg.ServiceName + "-migrate")

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("sql open")
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal().Err(err).Msg("goose set dialect")
	}

	if err := goose.Up(db, cfg.MigrationsDir); err != nil {
		logger.Fatal().Err(err).Msg("goose up")
	}

	_, _ = fmt.Fprintln(os.Stdout, "migrations applied")
}

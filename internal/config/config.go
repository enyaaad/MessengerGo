package config

import (
	"os"
)

type Config struct {
	ServiceName   string
	HTTPAddr      string
	DatabaseDSN   string
	MigrationsDir string
}

func FromEnv() Config {
	return Config{
		ServiceName:   getEnv("SERVICE_NAME", "messenger-test"),
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		DatabaseDSN:   getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/hirify?sslmode=disable"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "migrations"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

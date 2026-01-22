package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func New(serviceName string) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	return zerolog.New(output).With().Timestamp().Str("service", serviceName).Logger()
}

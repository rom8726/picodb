package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"picodb/internal/config"

	"github.com/rs/zerolog"
)

const (
	loggerTimestampFormat = "2006-01-02 15:04:05"
)

func main() {
	console := consoleLogger()
	if err := run(console); err != nil {
		console.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func run(console *zerolog.Logger) error {
	console.Info().Msg("init config...")
	_, err := config.Init()
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()

	return nil
}

func consoleLogger() *zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: loggerTimestampFormat}
	logger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	return &logger
}

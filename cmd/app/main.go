package main

import (
	"context"
	"log/slog"
	"os"
	"web_chat/config"
	"web_chat/internal/infrastructure/postgres"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	dbPool, closePool, err := postgres.NewPool(context.Background(), cfg)
	if err != nil {
		logger.Error("failed to init db connection", slog.Any("error", err))
		return err
	}
	defer closePool()

	if err := NewApp(logger, dbPool).Run(); err != nil {
		logger.Error("server exited with error", slog.Any("error", err))
		return err
	}

	return nil
}

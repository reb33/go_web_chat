package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	if err := NewApp(logger).Run(); err != nil {
		logger.Error("server exited with error", slog.Any("error", err))
		return err
	}

	return nil
}

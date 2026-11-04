package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_chat/internal/infrastructure/postgres"
	api "web_chat/internal/transport/httpserver"
	"web_chat/internal/usecase/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	GracefulShutdownTimeout = 5 * time.Second
	appServerCount          = 3
	debugPort               = 9001
)

type App struct {
	router *echo.Echo
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, dbPool *pgxpool.Pool) *App {
	authRepo := auth.Repository(postgres.NewAuthRepo(dbPool))
	authUseCase := auth.New(authRepo)
	router := api.NewRouter(logger, authUseCase)
	return &App{router: router, logger: logger}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return a.run(ctx)
}

func (a *App) run(ctx context.Context) error {
	port := flag.Int("port", 9000, "port")
	flag.Parse()
	serverErrors := make(chan error, appServerCount)
	go func() {
		addr := fmt.Sprintf(":%d", *port)
		a.logger.InfoContext(ctx, "Starting server", slog.Any("address", addr))
		if err := a.router.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("API server: %w", err)
		}
	}()

	var debugServer *http.Server
	if debugPort != 0 {
		debugServer = newDebugServer(debugPort)
		go func() {
			a.logger.InfoContext(ctx, "Starting debug server", slog.Any("address", debugServer.Addr))
			if err := debugServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				serverErrors <- fmt.Errorf("debug server: %w", err)
			}
		}()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-serverErrors:
	}
	a.logger.InfoContext(ctx, "Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if debugServer != nil {
		if err := debugServer.Shutdown(shutdownCtx); err != nil {
			a.logger.ErrorContext(ctx, "failed to shutdown debug server", slog.Any("error", err))
		}
	}

	apiErr := a.router.Shutdown(shutdownCtx)
	return errors.Join(runErr, apiErr)
}

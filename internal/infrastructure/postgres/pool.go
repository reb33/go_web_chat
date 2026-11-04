package postgres

import (
	"context"
	"fmt"
	"time"
	"web_chat/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// трасер таймаутов
type timeoutTracer struct {
	timeout time.Duration
}

type commandCancelKey struct{}

func (t *timeoutTracer) TraceQueryStart(
	ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData,
) context.Context {
	queryCtx, cancel := context.WithTimeout(ctx, t.timeout)
	return context.WithValue(queryCtx, commandCancelKey{}, cancel)
}

func (t *timeoutTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	if cancel, ok := ctx.Value(commandCancelKey{}).(context.CancelFunc); ok {
		cancel()
	}
}

// Функция инициализации пула с применением трасера
func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
	// 1. Парсим стандартный конфиг из DSN строки
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, nil, err
	}

	// 2. Внедряем наш кастомный трасер в конфигурацию соединения pgx
	poolConfig.ConnConfig.Tracer = &timeoutTracer{
		timeout: time.Duration(cfg.DB.CommandTimeout) * time.Second,
	}

	// 3. Создаем пул с готовым конфигом
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("postgres: ping: %w", err)
	}

	return pool, pool.Close, nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"
	"web_chat/internal/domain/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, created_at FROM users.users WHERE username = '$1'`

	var user entities.User
	err := r.pool.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		// User отсутствует
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("auth_repo: get user: %w", err)
	}
	return &user, nil
}

func (r *AuthRepo) CreateUser(ctx context.Context, username string) (*entities.User, error) {
	query := `INSERT INTO users.users (username) VALUES ('$1')`
	_, err := r.pool.Exec(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("auth_repo: create user: %w", err)
	}
	return r.GetUserByUsername(ctx, username)
}

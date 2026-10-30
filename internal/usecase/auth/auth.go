package auth

import (
	"context"
	"web_chat/internal/domain/entities"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	CreateUser(ctx context.Context, username string) (*entities.User, error)
}


type UseCase struct {
	repo    Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}


func (u *UseCase)Execute(ctx context.Context, req Request)(*Response, error) {
	user, err := u.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &Response{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt}, nil
	}
	user, err = u.repo.CreateUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &Response{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt}, nil
}

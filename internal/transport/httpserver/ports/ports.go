package ports

import (
	"context"
	"web_chat/internal/usecase/auth"
)

type AuthUseCase interface {
	Execute(ctx context.Context, req auth.Request) (*auth.Response, error)
}

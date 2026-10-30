package auth

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Username string
}

type Response struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
}

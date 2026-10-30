package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
}

type Group struct {
	ID        uuid.UUID
	Name      string
	CreatedBy uuid.UUID
	CreatedAt time.Time
}

type Message struct {
	ID        uuid.UUID
	GroupID   uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

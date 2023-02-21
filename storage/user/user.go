package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         string
	TelegramID string
	CreatedAt  time.Time
}

func NewUser(telegramID string) *User {
	return &User{ID: uuid.New().String(), TelegramID: telegramID, CreatedAt: time.Now()}
}

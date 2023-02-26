package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         string
	Username   string
	FullName   string
	TelegramID string
	Age        int
	Sex        string
	CreatedAt  time.Time
}

func NewUser(telegramID string, username string, name string) *User {
	return &User{
		ID:         uuid.New().String(),
		TelegramID: telegramID,
		Username:   username,
		FullName:   name,
		CreatedAt:  time.Now(),
	}
}

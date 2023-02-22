package handler

import (
	"context"

	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRequest struct {
	FromID int64

	Update    tgbotapi.Update
	User      user.User
	IsNewUser bool

	StopPropagation bool
}

type Handler interface {
	Support(*BotRequest) bool
	Handle(context.Context, *BotRequest) error
}

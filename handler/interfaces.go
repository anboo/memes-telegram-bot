package handler

import (
	"context"

	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotContext struct {
	Update    tgbotapi.Update
	User      user.User
	IsNewUser bool
}

type Handler interface {
	Support(BotContext) bool
	Handle(context.Context, BotContext) error
}

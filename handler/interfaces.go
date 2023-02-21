package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotContext struct {
	Update    tgbotapi.Update
	IsNewUser bool
}

type Handler interface {
	Support(BotContext) bool
	Handle(BotContext) error
}

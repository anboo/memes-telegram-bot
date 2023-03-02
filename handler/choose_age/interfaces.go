package choose_age

import (
	"context"

	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

type UserRepository interface {
	Upsert(ctx context.Context, u user.User) (user.User, bool, error)
}

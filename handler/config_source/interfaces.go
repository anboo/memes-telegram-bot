package config_source

import (
	"context"

	"memes-bot/storage/user_source"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSourceRepository interface {
	Create(ctx context.Context, us user_source.UserSource) error
	ByUserId(ctx context.Context, userId string) (res []user_source.UserSource, err error)
}

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

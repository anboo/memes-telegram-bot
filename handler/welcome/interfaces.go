package welcome

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

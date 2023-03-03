package send_mem

import (
	"context"

	"memes-bot/storage/mem"
	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MemRepository interface {
	FindRelevantMemForUser(ctx context.Context, u user.User) (mem.Mem, error)
	ReserveNewMem(ctx context.Context, user user.User, mem mem.Mem) error
}

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

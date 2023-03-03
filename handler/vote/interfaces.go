package vote

import (
	"context"

	"memes-bot/storage/mem"
	"memes-bot/storage/vote"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MemRepository interface {
	Find(ctx context.Context, id string) (mem.Mem, error)
	UpdateRating(ctx context.Context, memId string, diff int) error
}

type VoteRepository interface {
	Save(ctx context.Context, v *vote.Vote) error
}

type TelegramAPI interface {
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

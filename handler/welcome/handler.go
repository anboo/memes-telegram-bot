package welcome

import (
	"context"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Handler struct {
	bot *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{bot: bot}
}

func (h Handler) Support(r *handler.BotRequest) bool {
	return r.IsNewUser
}

func (h Handler) Handle(ctx context.Context, r *handler.BotRequest) error {
	msg := tgbotapi.NewMessage(
		r.FromID,
		"Вам нужно оценить минимум 20 мемов, прежде чем рекомендации станут более менее релевантными.",
	)
	_, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send welcome message")
	}
	return nil
}

package welcome

import (
	"context"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Handler struct {
	bot TelegramAPI
}

func NewHandler(bot TelegramAPI) *Handler {
	return &Handler{bot: bot}
}

func (h *Handler) String() string {
	return "welcome"
}

func (h Handler) Support(r *handler.BotRequest) bool {
	return r.IsNewUser
}

func (h Handler) Handle(ctx context.Context, r *handler.BotRequest) error {
	if !r.IsNewUser {
		return errors.New("is not new user")
	}

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

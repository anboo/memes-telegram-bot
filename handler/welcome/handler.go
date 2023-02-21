package welcome

import (
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

func (h Handler) Support(c handler.BotContext) bool {
	return c.IsNewUser
}

func (h Handler) Handle(c handler.BotContext) error {
	msg := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "Вам нужно оценить минимум 20 мемов, прежде чем рекомендации станут более менее релевантными.")
	msg.ReplyToMessageID = c.Update.Message.MessageID

	_, err := h.bot.Send(msg)

	return errors.Wrap(err, "send welcome message")
}

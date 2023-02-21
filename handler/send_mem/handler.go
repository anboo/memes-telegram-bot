package send_mem

import (
	"context"

	"memes-bot/handler"
	"memes-bot/storage/mem"
	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Handler struct {
	bot           *tgbotapi.BotAPI
	memRepository *mem.Repository
	log           *zerolog.Logger
}

func NewHandler(bot *tgbotapi.BotAPI, memRepository *mem.Repository, log *zerolog.Logger) *Handler {
	return &Handler{
		bot:           bot,
		memRepository: memRepository,
		log:           log,
	}
}

func (h Handler) Support(c handler.BotContext) bool {
	return true
}

func (h Handler) Handle(ctx context.Context, c handler.BotContext) error {
	mem, err := h.memRepository.FindRelevantMemForUser(ctx, user.User{})
	if err != nil {
		return errors.Wrap(err, "send mem handler")
	}

	msg := tgbotapi.NewMessage(c.Update.FromChat().ID, mem.Img)
	_, err = h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send mem handler bot send")
	}

	p := tgbotapi.NewPhoto(c.Update.FromChat().ID, tgbotapi.FileURL(mem.Img))
	p.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("👍", "up_"+mem.ID),
			tgbotapi.NewInlineKeyboardButtonData("👎", "down_"+mem.ID),
			tgbotapi.NewInlineKeyboardButtonData("🆘", "sos_"+mem.ID),
		),
	)
	h.bot.Send(p)

	return nil
}

package send_mem

import (
	"context"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Handler struct {
	bot           TelegramAPI
	memRepository MemRepository
	log           *zerolog.Logger
}

func NewHandler(bot TelegramAPI, memRepository MemRepository, log *zerolog.Logger) *Handler {
	return &Handler{
		bot:           bot,
		memRepository: memRepository,
		log:           log,
	}
}

func (h Handler) Support(r *handler.BotRequest) bool {
	return true
}

func (h Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	m, err := h.memRepository.FindRelevantMemForUser(ctx, request.User)
	if err != nil {
		return errors.Wrap(err, "send m handler")
	}

	p := tgbotapi.NewPhoto(request.FromID, tgbotapi.FileURL(m.Img))
	p.Caption = m.Text
	p.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👍", "up_"+m.ID),
			tgbotapi.NewInlineKeyboardButtonData("👎", "down_"+m.ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🆘", "sos_"+m.ID),
		),
	)

	_, err = h.bot.Send(p)
	if err != nil {
		return errors.Wrap(err, "try send send m message")
	}

	err = h.memRepository.ReserveNewMem(ctx, request.User, m)
	if err != nil {
		return errors.Wrap(err, "reserve m")
	}

	return nil
}

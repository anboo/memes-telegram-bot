package vote

import (
	"context"
	"strings"

	"memes-bot/handler"
	"memes-bot/storage/mem"
	"memes-bot/storage/vote"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	UpPrefix   = "up_"
	DownPrefix = "down_"
	SosPrefix  = "sos_"
)

type Handler struct {
	bot            *tgbotapi.BotAPI
	memRepository  *mem.Repository
	voteRepository *vote.Repository
	logger         *zerolog.Logger
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	memRepository *mem.Repository,
	voteRepository *vote.Repository,
	logger *zerolog.Logger,
) *Handler {
	return &Handler{
		bot:            bot,
		memRepository:  memRepository,
		voteRepository: voteRepository,
		logger:         logger,
	}
}

func (h Handler) Support(request *handler.BotRequest) bool {
	return request.Update.CallbackQuery != nil && (strings.HasPrefix(request.Update.CallbackQuery.Data, UpPrefix) ||
		strings.HasPrefix(request.Update.CallbackQuery.Data, DownPrefix) ||
		strings.HasPrefix(request.Update.CallbackQuery.Data, SosPrefix))
}

func (h Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	data := strings.Split(request.Update.CallbackQuery.Data, "_")
	if len(data) < 2 {
		return errors.New("incorrect data")
	}

	prefix, memId := data[0], data[1]

	rating := 0
	switch prefix {
	case UpPrefix:
		rating = 1
	case DownPrefix:
		rating = -1
	case SosPrefix:
		rating = -2
	}

	meme, err := h.memRepository.Find(ctx, memId)
	if err != nil {
		return errors.Wrap(err, "vote handler find meme")
	}

	err = h.memRepository.UpdateRating(ctx, meme.ID, rating)
	if err != nil {
		return errors.Wrap(err, "vote handler try update mem rating")
	}

	err = h.voteRepository.Save(ctx, *vote.NewVote(memId, request.User.ID, rating))
	if err != nil {
		return errors.Wrap(err, "vote handler save vote")
	}

	_, err = h.bot.Request(tgbotapi.NewCallback(request.Update.CallbackQuery.ID, "Спасибо за вашу оценку!"))
	if err != nil {
		return errors.Wrap(err, "vote handler reply button callback")
	}

	return nil
}

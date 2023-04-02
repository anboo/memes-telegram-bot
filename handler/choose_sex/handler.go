package choose_sex

import (
	"context"
	"strings"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var (
	ChooseSexPrefix = "choose_sex_"

	SexMen  = "m"
	SexGirl = "g"
	SexFish = "f"
)

type Handler struct {
	bot            TelegramAPI
	userRepository UserRepository
}

func NewHandler(bot TelegramAPI, userRepository UserRepository) *Handler {
	return &Handler{
		bot:            bot,
		userRepository: userRepository,
	}
}

func (h *Handler) String() string {
	return "choose_sex"
}

func (h *Handler) Support(request *handler.BotRequest) bool {
	return request.User.Sex == "" ||
		(request.Update.CallbackQuery != nil && strings.HasPrefix(request.Update.CallbackQuery.Data, ChooseSexPrefix))
}

func (h *Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	switch {
	case request.Update.CallbackQuery != nil && strings.HasPrefix(request.Update.CallbackQuery.Data, ChooseSexPrefix):
		return h.handleClick(ctx, request)
	case request.User.Sex == "":
		return h.handleShowMenu(ctx, request)
	default:
		return errors.New("incorrect handle choose sex")
	}
}

func (h *Handler) handleClick(ctx context.Context, request *handler.BotRequest) error {
	sex := strings.TrimLeft(request.Update.CallbackQuery.Data, ChooseSexPrefix)
	if sex != SexMen && sex != SexGirl && sex != SexFish {
		return errors.New("incorrect sex")
	}

	request.User.Sex = sex
	_, _, err := h.userRepository.Upsert(ctx, request.User)
	if err != nil {
		return errors.Wrap(err, "try set new sex")
	}

	_, err = h.bot.Request(tgbotapi.NewCallback(request.Update.CallbackQuery.ID, "Спасибо"))
	if err != nil {
		return errors.Wrap(err, "try response callback choose sex")
	}

	return nil
}

func (h *Handler) handleShowMenu(ctx context.Context, request *handler.BotRequest) error {
	msg := tgbotapi.NewMessage(request.FromID, "Выберите ваш пол")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👱🏼‍♂️Мужчина", ChooseSexPrefix+SexMen),
			tgbotapi.NewInlineKeyboardButtonData("👱🏼‍♀️Девушка", ChooseSexPrefix+SexGirl),
			tgbotapi.NewInlineKeyboardButtonData("🐟 My son?? Where is my son??", ChooseSexPrefix+SexFish),
		),
	)
	_, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send chose sex message")
	}

	request.StopPropagation = true

	return nil
}

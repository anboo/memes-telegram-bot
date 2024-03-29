package choose_age

import (
	"context"
	"strconv"
	"strings"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var (
	ChooseAgePrefix = "choose_age_"

	AgeLessThen18  = "1"
	AgeBetween1825 = "2"
	AgeGreatThen25 = "3"
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
	return "choose_age"
}

func (h *Handler) Support(request *handler.BotRequest) bool {
	return request.User.Age == 0 ||
		(request.Update.CallbackQuery != nil && strings.HasPrefix(request.Update.CallbackQuery.Data, ChooseAgePrefix))
}

func (h *Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	switch {
	case request.Update.CallbackQuery != nil && strings.HasPrefix(request.Update.CallbackQuery.Data, ChooseAgePrefix):
		return h.handleClick(ctx, request)
	case request.User.Age == 0:
		return h.showMenu(ctx, request)
	default:
		return errors.New("incorrect handle choose age")
	}
}

func (h *Handler) showMenu(ctx context.Context, request *handler.BotRequest) error {
	msg := tgbotapi.NewMessage(request.FromID, "Выберите ваш возраст")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("До 18", ChooseAgePrefix+AgeLessThen18),
			tgbotapi.NewInlineKeyboardButtonData("От 18 до 25", ChooseAgePrefix+AgeBetween1825),
			tgbotapi.NewInlineKeyboardButtonData("Больше 25", ChooseAgePrefix+AgeGreatThen25),
		),
	)
	_, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send chose age message")
	}

	request.StopPropagation = true

	return nil
}

func (h *Handler) handleClick(ctx context.Context, request *handler.BotRequest) error {
	age := strings.TrimLeft(request.Update.CallbackQuery.Data, ChooseAgePrefix)
	if age != AgeLessThen18 && age != AgeBetween1825 && age != AgeGreatThen25 {
		return errors.New("incorrect age")
	}

	userAge, _ := strconv.Atoi(age)
	request.User.Age = userAge
	_, _, err := h.userRepository.Upsert(ctx, request.User)
	if err != nil {
		return errors.Wrap(err, "try set new age")
	}

	_, err = h.bot.Request(tgbotapi.NewCallback(request.Update.CallbackQuery.ID, "Спасибо"))
	if err != nil {
		return errors.Wrap(err, "try response callback choose age")
	}

	return nil
}

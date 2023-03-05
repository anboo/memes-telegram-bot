package config_source

import (
	"context"
	"fmt"
	"strings"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var (
	callbackKeyPrefix = "choose_config_"
)

type Handler struct {
	bot                  TelegramAPI
	userSourceRepository UserSourceRepository
	vkGroups             []string
}

func NewHandler(bot TelegramAPI, userSourceRepository UserSourceRepository, vkGroups string) *Handler {
	return &Handler{
		bot:                  bot,
		userSourceRepository: userSourceRepository,
		vkGroups:             strings.Split(vkGroups, ","),
	}
}

func (h *Handler) Support(request *handler.BotRequest) bool {
	return (request.Update.CallbackQuery != nil &&
		strings.HasPrefix(request.Update.CallbackQuery.Data, callbackKeyPrefix)) ||
		request.Update.Message.Command() == "/settings"
}

func (h *Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	userSources, err := h.userSourceRepository.ByUserId(ctx, request.User.ID)
	if err != nil {
		return errors.Wrap(err, "try fetch user sources")
	}

	buttons := make([]tgbotapi.InlineKeyboardButton, len(h.vkGroups))
	for j, group := range h.vkGroups {
		keyboardKey := "disable_" + group
		keyboardValue := fmt.Sprintf("❌ %s", group)

		for _, userSource := range userSources {
			if userSource.Source == keyboardKey && !userSource.Enabled {
				keyboardKey = "enable_" + group
				keyboardValue = fmt.Sprintf("✅ %s", group)
			}
		}

		buttons[j] = tgbotapi.NewInlineKeyboardButtonData(keyboardValue, keyboardKey)
	}

	keyboard := tgbotapi.NewInlineKeyboardRow(buttons...)
	msg := tgbotapi.NewMessage(request.FromID, "")
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "try send config source message")
	}

	return nil
}

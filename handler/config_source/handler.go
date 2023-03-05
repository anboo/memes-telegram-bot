package config_source

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"memes-bot/handler"
	"memes-bot/storage/user_source"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var (
	callbackKeyPrefix           = "choose_config_"
	callbackDisableSourcePrefix = "disable_"
	callbackEnableSourcePrefix  = "enable_"
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
		(request.Update.Message != nil && request.Update.Message.Command() == "settings")
}

func (h *Handler) Handle(ctx context.Context, request *handler.BotRequest) error {
	switch {
	case request.Update.Message != nil && request.Update.Message.Command() == "settings":
		return h.showSettings(ctx, request)
	case request.Update.CallbackQuery != nil && strings.HasPrefix(request.Update.CallbackQuery.Data, callbackKeyPrefix):
		return h.handleButton(ctx, request)
	}

	return nil
}

func (h *Handler) handleButton(ctx context.Context, request *handler.BotRequest) error {
	request.StopPropagation = true

	sourceData := strings.Split(strings.TrimPrefix(request.Update.CallbackQuery.Data, callbackKeyPrefix), "_")
	if len(sourceData) < 2 {
		return errors.New("incorrect data")
	}
	messageID, action, source := sourceData[0], sourceData[1], sourceData[2]

	found := false
	for _, s := range h.vkGroups {
		if source == s {
			found = true
		}
	}

	if !found {
		_, err := h.bot.Request(tgbotapi.NewCallback(request.Update.CallbackQuery.ID, "Ошибка"))
		if err != nil {
			return errors.Wrap(err, "try send source not found callback response")
		}
		return errors.New("not found source")
	}

	userSource := user_source.New(request.User.ID, source)
	if action == "disable" {
		userSource.Enabled = false
	}

	err := h.userSourceRepository.Create(ctx, userSource)
	if err != nil {
		return errors.Wrap(err, "try create user_source failed")
	}

	_, err = h.bot.Request(tgbotapi.NewCallback(request.Update.CallbackQuery.ID, "Сохранено"))
	if err != nil {
		return errors.Wrap(err, "try send callback ok")
	}

	telegramMessageId, err := strconv.Atoi(messageID)
	if err != nil {
		return errors.Wrap(err, "try parse message id")
	}

	keyboard, err := h.generateSettingsKeyboard(ctx, request, telegramMessageId)
	if err != nil {
		return errors.Wrap(err, "try generate keyboard")
	}

	_, err = h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(request.FromID, telegramMessageId, keyboard))
	if err != nil {
		return errors.Wrap(err, "try send new keyboard")
	}

	return nil
}

func (h *Handler) showSettings(ctx context.Context, request *handler.BotRequest) error {
	request.StopPropagation = true

	msg, err := h.bot.Send(tgbotapi.NewMessage(request.FromID, "Выберите источники откуда получать мемы:"))
	if err != nil {
		return errors.Wrap(err, "try send config source message")
	}

	keyboard, err := h.generateSettingsKeyboard(ctx, request, msg.MessageID)
	if err != nil {
		return errors.Wrap(err, "try generate keyboard")
	}

	_, err = h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(request.FromID, msg.MessageID, keyboard))
	if err != nil {
		return errors.Wrap(err, "try update keyboard")
	}

	return nil
}

func (h *Handler) generateSettingsKeyboard(ctx context.Context, request *handler.BotRequest, messageID int) (tgbotapi.InlineKeyboardMarkup, error) {
	userSources, err := h.userSourceRepository.ByUserId(ctx, request.User.ID)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, errors.Wrap(err, "try fetch user sources")
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, len(h.vkGroups))
	for j, group := range h.vkGroups {
		keyboardKey := callbackKeyPrefix + strconv.Itoa(messageID) + "_" + callbackDisableSourcePrefix + group
		keyboardValue := fmt.Sprintf("❌ %s", group)

		for _, userSource := range userSources {
			if userSource.Source == group && !userSource.Enabled {
				keyboardKey = callbackKeyPrefix + strconv.Itoa(messageID) + "_" + callbackEnableSourcePrefix + group
				keyboardValue = fmt.Sprintf("✅ %s", group)
			}
		}

		buttons[j] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(keyboardValue, keyboardKey))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...), nil
}

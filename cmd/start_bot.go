package cmd

import (
	"context"
	"strconv"

	"memes-bot/handler"
	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type StartBotCmd struct {
	bot            *tgbotapi.BotAPI
	userRepository *user.Repository
	router         *handler.Router
	l              *zerolog.Logger
}

func NewStartBotCmd(bot *tgbotapi.BotAPI, userRepository *user.Repository, router *handler.Router, l *zerolog.Logger) *StartBotCmd {
	return &StartBotCmd{
		bot:            bot,
		userRepository: userRepository,
		router:         router,
		l:              l,
	}
}

func (c *StartBotCmd) Execute(ctx context.Context) error {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60
	updates := c.bot.GetUpdatesChan(config)

	for update := range updates {
		select {
		default:
			var (
				telegramID string
				fromID     int64
			)

			switch {
			case update.Message != nil:
				fromID = update.Message.From.ID
				telegramID = strconv.Itoa(int(fromID))
			case update.CallbackQuery != nil:
				fromID = update.CallbackQuery.From.ID
				telegramID = strconv.Itoa(int(fromID))
			default:
				continue
			}

			isUpdated := true
			u, err := c.userRepository.ByTelegramID(ctx, telegramID)
			if err != nil {
				u, isUpdated, err = c.userRepository.Upsert(ctx, *user.NewUser(telegramID))
				if err != nil {
					c.l.Err(err).Interface("update", update).Msg("try upsert update")
				}
			}

			botContext := handler.BotRequest{FromID: fromID, Update: update, IsNewUser: !isUpdated, User: u}
			err = c.router.Handle(ctx, botContext)
			if err != nil {
				c.l.Err(err).Msg("error router")
			}
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

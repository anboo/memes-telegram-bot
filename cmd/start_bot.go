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

			var from *tgbotapi.User
			switch {
			case update.Message != nil:
				from = update.Message.From
			case update.CallbackQuery != nil:
				from = update.CallbackQuery.From
			default:
				continue
			}

			fromID = from.ID
			telegramID = strconv.Itoa(int(fromID))

			isUpdated := true
			u, err := c.userRepository.ByTelegramID(ctx, telegramID)
			if err != nil {
				u, isUpdated, err = c.userRepository.Upsert(ctx, *user.NewUser(telegramID, from.UserName, from.FirstName+" "+from.LastName))
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

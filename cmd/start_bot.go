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
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

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

			user, isUpdated, err := c.userRepository.Upsert(ctx, *user.NewUser(telegramID))
			if err != nil {
				c.l.Err(err).Interface("update", update).Msg("try upsert user")
			}

			botContext := handler.BotContext{FromID: fromID, Update: update, IsNewUser: !isUpdated, User: user}
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

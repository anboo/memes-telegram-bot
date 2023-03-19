package handler

import (
	"context"
	"time"

	"memes-bot/metrics"
	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BotRequest struct {
	FromID int64

	Update    tgbotapi.Update
	User      user.User
	IsNewUser bool

	StopPropagation bool
}

type Router struct {
	handlers []Handler
	logger   *zerolog.Logger
}

func NewRouter(logger *zerolog.Logger, handlers ...Handler) *Router {
	return &Router{
		logger:   logger,
		handlers: handlers,
	}
}

func (r *Router) Handle(ctx context.Context, request BotRequest) error {
	var found bool

	for _, h := range r.handlers {
		if h.Support(&request) {
			var err error

			startedAt := time.Now()
			found = true

			err = h.Handle(ctx, &request)
			if err != nil {
				metrics.TelegramHandlerDuration.WithLabelValues(h.String(), "error").Observe(time.Since(startedAt).Seconds())
				r.logger.Err(err).Str("userId", request.User.ID).Int64("telegramId", request.FromID).Msg("handler error")
				return nil
			} else {
				metrics.TelegramHandlerDuration.WithLabelValues(h.String(), "success").Observe(time.Since(startedAt).Seconds())
			}

			if request.StopPropagation {
				return nil
			}
		}
	}

	if !found {
		return errors.New("not found handler")
	}
	return nil
}

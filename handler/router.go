package handler

import (
	"context"

	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
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
}

func NewRouter(handlers ...Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Handle(ctx context.Context, request BotRequest) error {
	var found bool

	for _, h := range r.handlers {
		if h.Support(&request) {
			found = true
			err := h.Handle(ctx, &request)
			if err != nil {
				return errors.Wrap(err, "handler")
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

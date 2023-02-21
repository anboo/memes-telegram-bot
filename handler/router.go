package handler

import (
	"context"

	"github.com/pkg/errors"
)

type Router struct {
	handlers []Handler
}

func NewRouter(handlers ...Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Handle(ctx context.Context, botContext BotContext) error {
	for _, h := range r.handlers {
		if h.Support(botContext) {
			return h.Handle(ctx, botContext)
		}
	}
	return errors.New("not found handler")
}

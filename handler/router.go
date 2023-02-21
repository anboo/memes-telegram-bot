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
	var found bool

	for _, h := range r.handlers {
		if h.Support(&botContext) {
			found = true
			err := h.Handle(ctx, &botContext)
			if err != nil {
				return errors.Wrap(err, "handler")
			}
			if botContext.StopPropagation {
				return nil
			}
		}
	}

	if !found {
		return errors.New("not found handler")
	}
	return nil
}

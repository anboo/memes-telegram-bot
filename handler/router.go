package handler

import (
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

func (r *Router) Handle(ctx BotContext) error {
	for _, h := range r.handlers {
		if h.Support(ctx) {
			return h.Handle(ctx)
		}
	}
	return errors.New("not found handler")
}

package handler

import (
	"context"
)

type Handler interface {
	Support(*BotRequest) bool
	Handle(context.Context, *BotRequest) error
}

package handler

import (
	"context"
)

type Handler interface {
	String() string
	Support(*BotRequest) bool
	Handle(context.Context, *BotRequest) error
}

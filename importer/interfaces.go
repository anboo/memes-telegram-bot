package importer

import (
	"context"

	"memes-bot/storage/mem"
)

type Importer interface {
	Import(ctx context.Context) (chan mem.Mem, chan struct{})
}

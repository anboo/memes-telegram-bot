package importer

import (
	"context"
	"strconv"

	"memes-bot/storage/mem"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/rs/zerolog"
)

var Groups = []string{
	"borsch",
	"agil_vk",
	"fuck_humor",
	"leprazo",
	"s_arcazm",
	"in.humour",
	"dzenpub",
}

type VkImporter struct {
	l *zerolog.Logger
}

func NewVkImporter(l *zerolog.Logger) *VkImporter {
	return &VkImporter{
		l: l,
	}
}

func (i *VkImporter) Import(ctx context.Context) (chan mem.Mem, chan struct{}) {
	var (
		ch   = make(chan mem.Mem)
		stop = make(chan struct{})
	)

	go i.startParsing(ctx, ch, stop)

	return ch, stop
}

func (i *VkImporter) startParsing(ctx context.Context, ch chan mem.Mem, stop chan struct{}) {
	vk := api.NewVK("4fb4ae9e4fb4ae9e4fb4ae9e144fcc4e8544fb44fb4ae9e2ebb0c3da28feb0386eb89b5")

	for _, g := range Groups {
		select {
		default:
			params := map[string]interface{}{
				"domain":       g,
				"access_token": "4fb4ae9e4fb4ae9e4fb4ae9e144fcc4e8544fb44fb4ae9e2ebb0c3da28feb0386eb89b5",
				"v":            "5.131",
				"count":        100,
				"offset":       0,
			}

			res, err := vk.WallGet(params)
			if err != nil {
				i.l.Err(err).Interface("params", params).Msg("import from vk")
				continue
			}

			for _, item := range res.Items {
				if len(item.Attachments) == 0 ||
					len(item.Attachments[0].Photo.Sizes) == 0 {
					continue
				}

				if BlackListed(item.Text) {
					continue
				}

				sizes := item.Attachments[0].Photo.Sizes
				ch <- *mem.NewMem(strconv.Itoa(item.ID), item.Text, "vk", g, sizes[len(sizes)-1].URL)
			}
		case <-ctx.Done():
			return
		}
	}

	stop <- struct{}{}
}

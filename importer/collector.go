package importer

import (
	"context"

	"memes-bot/storage/mem"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Collector struct {
	importers []Importer
	rep       *mem.Repository
}

func NewCollector(rep *mem.Repository, importers []Importer) *Collector {
	return &Collector{
		importers: importers,
		rep:       rep,
	}
}

func (c *Collector) Import(ctx context.Context) error {
	errg, _ := errgroup.WithContext(ctx)

	for _, i := range c.importers {
		errg.Go(func() error {
			ch, stop := i.Import(ctx)

			errg.Go(func() error {
				return c.save(ctx, ch, stop)
			})

			return nil
		})
	}

	err := errg.Wait()
	if err != nil {
		return errors.Wrapf(err, "import wait")
	}

	return nil
}

func (c *Collector) save(ctx context.Context, ch chan mem.Mem, stop chan struct{}) error {
	for {
		select {
		case m := <-ch:
			_, err := c.rep.UpsertMem(ctx, m)
			if err != nil {
				return errors.Wrap(err, "save imported mem")
			}
		case <-stop:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}

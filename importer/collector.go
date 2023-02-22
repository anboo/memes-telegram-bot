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
	var (
		memes []mem.Mem
		j     = 0
	)

	for {
		select {
		case m := <-ch:
			memes = append(memes, m)
			j++
			if j >= 1000 {
				memes = []mem.Mem{}
				err := c.flush(ctx, memes)
				if err != nil {
					return errors.Wrap(err, "save imported memes")
				}
				j = 0
			}
		case <-stop:
			err := c.flush(ctx, memes)
			if err != nil {
				return errors.Wrap(err, "save imported memes")
			}
			return nil
		case <-ctx.Done():
			err := c.flush(ctx, memes)
			if err != nil {
				return errors.Wrap(err, "save imported memes")
			}
			return nil
		}
	}
}

func (c *Collector) flush(ctx context.Context, memes []mem.Mem) error {
	err := c.rep.BatchCreate(ctx, memes)
	if err != nil {
		return errors.Wrap(err, "save imported memes")
	}
	return nil
}

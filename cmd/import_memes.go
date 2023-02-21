package cmd

import (
	"context"

	"memes-bot/importer"
)

type ImportMemesCmd struct {
	collector *importer.Collector
}

func NewImportMemesCmd(collector *importer.Collector) *ImportMemesCmd {
	return &ImportMemesCmd{
		collector: collector,
	}
}

func (c *ImportMemesCmd) Execute(ctx context.Context) error {
	err := c.collector.Import(ctx)
	if err != nil {
		return err
	}
	return nil
}

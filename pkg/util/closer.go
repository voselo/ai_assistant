package closer

import (
	"context"
)

//nolint:gochecknoglobals
var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func AddWithCtx(f ...func(ctx context.Context) error) {
	globalCloser.AddWithCtx(f...)
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	funcs        []func() error
	funcsWithCtx []func(ctx context.Context) error
}

func New() *Closer {
	return &Closer{
		funcs:        make([]func() error, 0),
		funcsWithCtx: make([]func(ctx context.Context) error, 0),
	}
}

func (c *Closer) Add(f ...func() error) {
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) AddWithCtx(f ...func(ctx context.Context) error) {
	c.funcsWithCtx = append(c.funcsWithCtx, f...)
}

func (c *Closer) CloseAll() {
	for _, f := range c.funcs {
		if err := f(); err != nil {
			// .Errorf("err close %v", err)
		}
	}
	// TODO: flush all loggers
	// if err := logger.Zap.Sync(); err != nil {
	// logger.SugarZap.Errorf("err sync logger %v", err)
	// }
}

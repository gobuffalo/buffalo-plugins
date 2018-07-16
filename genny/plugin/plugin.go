package plugin

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	err := NormalizeOptions(opts)
	if err != nil {
		return g, errors.WithStack(err)
	}

	g.Box(packr.NewBox("./templates"))
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Replace("-dot-", "."))

	return g, nil
}

package with

import (
	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func GenerateCmd(opts *plugin.Options) (*genny.Generator, error) {
	g := genny.New()
	box := packr.NewBox("./generate/templates")
	if err := g.Box(box); err != nil {
		return g, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))

	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Dot())

	return g, nil
}

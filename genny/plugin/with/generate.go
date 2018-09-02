package with

import (
	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/genny/new"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func GenerateCmd(opts *plugin.Options) (*genny.Group, error) {
	gg := &genny.Group{}
	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	box := packr.NewBox("./generate/templates")
	if err := g.Box(box); err != nil {
		return gg, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))

	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Dot())
	gg.Add(g)

	g, err := new.New(&new.Options{
		Name:   opts.ShortName,
		Prefix: "genny",
	})
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(g)

	return gg, nil
}

package plugin

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/gobuffalo/genny/movinglater/licenser"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	g.Box(packr.NewBox("../plugin/templates"))
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Dot())
	gg.Add(g)

	lopts := &licenser.Options{
		Author: opts.Author,
		Name:   opts.License,
	}

	g, err := licenser.New(lopts)
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(g)

	gm, err := gomods.New(opts.PluginPkg, opts.Root)
	if err != nil && errors.Cause(err) != gomods.ErrModsOff {
		return gg, errors.WithStack(err)
	}
	gg.Merge(gm)
	return gg, nil
}

package with

import (
	"strings"

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

	t := genny.NewTransformer("*", func(f genny.File) (genny.File, error) {
		name := f.Name()
		name = strings.Replace(name, "-shortName-", opts.ShortName, -1)
		name = strings.Replace(name, "-dot-", ".", -1)
		return genny.NewFile(name, f), nil
	})
	g.Transformer(t)

	return g, nil
}

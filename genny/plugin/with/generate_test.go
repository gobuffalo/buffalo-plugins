package with

import (
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_GenerateCmd(t *testing.T) {
	r := require.New(t)

	opts := &plugin.Options{
		PluginPkg: "github.com/foo/buffalo-bar",
		Year:      1999,
		Author:    "Homer Simpson",
		ShortName: "bar",
	}

	run := genny.DryRunner(context.Background())

	g, err := GenerateCmd(opts)
	r.NoError(err)
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 5)

	f := res.Files[0]
	r.Equal("cmd/available.go", f.Name())
	r.Contains(f.String(), `{Name: generateCmd.Use, BuffaloCommand: "generate", Description: generateCmd.Short, Aliases: generateCmd.Aliases}`)

	f = res.Files[1]
	r.Equal("cmd/generate.go", f.Name())
	r.Contains(f.String(), opts.PluginPkg+"/genny/")

	f = res.Files[2]
	r.Equal("genny/bar/bar.go", f.Name())
	r.Contains(f.String(), "package bar")
	r.Contains(f.String(), "func New(opts *Options) (*genny.Generator, error)")

	f = res.Files[3]
	r.Equal("genny/bar/options.go", f.Name())
	r.Contains(f.String(), "package bar")
	r.Contains(f.String(), "type Options struct {")

	f = res.Files[4]
	r.Equal("genny/bar/templates/example.txt", f.Name())

}

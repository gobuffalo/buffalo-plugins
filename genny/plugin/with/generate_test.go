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

	run.FileFn = func(f genny.File) error {
		return nil
	}

	for _, f := range run.Results().Files {
		body := f.String()
		switch f.Name() {
		case "cmd/generate.go":
			r.Contains(body, opts.PluginPkg+"/genny/")
		case "genny/bar/bar.go":
			r.Contains(body, "package bar")
			r.Contains(body, "func New(opts *Options) (*genny.Generator, error)")
		case "genny/bar/options.go":
			r.Contains(body, "package bar")
			r.Contains(body, "type Options struct {")
		case "cmd/available.go":
			r.Contains(body, `{Name: generateCmd.Use, BuffaloCommand: "generate", Description: generateCmd.Short, Aliases: generateCmd.Aliases}`)
		}
	}

}

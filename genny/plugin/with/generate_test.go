package with

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/sirupsen/logrus"
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
	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb
	run.Logger = l

	run.FileFn = func(f genny.File) error {
		run.Logger.Infof(f.Name())
		b, err := ioutil.ReadAll(f)
		r.NoError(err)
		body := string(b)
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
			r.Contains(body, `{Name: "generate", BuffaloCommand: "generate", Description: generateCmd.Short, Aliases: generateCmd.Aliases}`)
		}
		return nil
	}

	g, err := GenerateCmd(opts)
	r.NoError(err)

	run.With(g)

	r.NoError(run.Run())
	out := bb.String()
	r.Contains(out, "cmd/generate.go")
	r.Contains(out, "genny/bar/bar.go")
	r.Contains(out, "cmd/available.go")

}

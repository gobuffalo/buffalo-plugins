package plugin

import (
	"context"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	r := require.New(t)
	opts := &Options{
		PluginPkg: "github.com/foo/buffalo-bar",
		Year:      1999,
		Author:    "Homer Simpson",
		ShortName: "bar",
	}
	run := genny.DryRunner(context.Background())

	run.FileFn = func(f genny.File) error {
		b, err := ioutil.ReadAll(f)
		r.NoError(err)
		body := string(b)
		switch f.Name() {
		case "LICENSE":
			r.Contains(body, strconv.Itoa(opts.Year))
			r.Contains(body, opts.Author)
		case "README.md":
			r.Contains(body, opts.PluginPkg)
		case "main.go":
			r.Contains(body, opts.PluginPkg+"/cmd")
		}
		return nil
	}

	g, err := New(opts)
	r.NoError(err)
	run.With(g)

	r.NoError(run.Run())
}

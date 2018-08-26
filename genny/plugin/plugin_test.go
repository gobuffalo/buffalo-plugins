package plugin

import (
	"context"
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

	g, err := New(opts)
	r.NoError(err)
	run.With(g)

	r.NoError(run.Run())
	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 9)

	f := res.Files[0]
	r.Equal(".travis.yml", f.Name())

	f = res.Files[1]
	r.Equal("Makefile", f.Name())

	f = res.Files[2]
	r.Equal("README.md", f.Name())
	r.Contains(f.String(), opts.PluginPkg)

	f = res.Files[3]
	r.Equal("bar/version.go", f.Name())
	r.Contains(f.String(), opts.ShortName)
	r.Contains(f.String(), "v0.0.0")

	f = res.Files[4]
	r.Equal("cmd/available.go", f.Name())

	f = res.Files[5]
	r.Equal("cmd/bar.go", f.Name())

	f = res.Files[6]
	r.Equal("cmd/root.go", f.Name())

	f = res.Files[7]
	r.Equal("cmd/version.go", f.Name())
	r.Contains(f.String(), opts.PluginPkg+"/"+opts.ShortName)
	r.Contains(f.String(), opts.ShortName+".Version")

	f = res.Files[8]
	r.Equal("main.go", f.Name())
	r.Contains(f.String(), "github.com/foo/buffalo-bar/cmd")

}

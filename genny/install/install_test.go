package install

import (
	"bytes"
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		App: meta.New("."),
		Plugins: []plugdeps.Plugin{
			{Binary: "buffalo-pop", GoGet: "github.com/gobuffalo/buffalo-pop"},
		},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	c := build.Default
	run.Disk.Add(genny.NewFile(filepath.Join(c.GOPATH, "bin", "buffalo-pop"), &bytes.Buffer{}))
	run.FileFn = func(f genny.File) (genny.File, error) {
		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, f); err != nil {
			return f, errors.WithStack(err)
		}
		return genny.NewFile(f.Name(), bb), nil
	}

	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	ecmds := []string{"go get github.com/gobuffalo/buffalo-pop"}
	for i, c := range res.Commands {
		r.Equal(ecmds[i], strings.Join(c.Args, " "))
	}
	r.Len(res.Commands, len(ecmds))

	efiles := []string{"bin/buffalo-pop", "config/buffalo-plugins.toml"}
	for i, f := range res.Files {
		r.True(strings.HasSuffix(f.Name(), efiles[i]))
	}
	r.Len(res.Files, len(efiles))
}

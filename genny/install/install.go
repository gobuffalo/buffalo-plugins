package install

import (
	"bytes"
	"go/build"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	proot := filepath.Join(opts.App.Root, "plugins")
	for _, p := range opts.Plugins {
		g.RunFn(pRun(proot, p))
	}

	bb := &bytes.Buffer{}
	if err := toml.NewEncoder(bb).Encode(plugdeps.Plugins{Plugins: opts.Plugins}); err != nil {
		return g, errors.WithStack(err)
	}

	cpath := filepath.Join(opts.App.Root, "config", "buffalo-plugins.toml")
	g.File(genny.NewFile(cpath, bb))

	return g, nil
}

func pRun(proot string, p plugdeps.Plugin) genny.RunFn {
	return func(r *genny.Runner) error {
		if err := gotools.Install(p.GoGet)(r); err != nil {
			return errors.WithStack(err)
		}

		c := build.Default
		if c.GOOS == "windows" {
			p.Binary += ".exe"
		}

		bp := filepath.Join(c.GOPATH, "bin", p.Binary)
		sf, err := r.FindFile(bp)
		if err != nil {
			return errors.WithStack(err)
		}

		pbp := filepath.Join(proot, p.Binary)
		r.Disk.Delete(pbp)

		df := genny.NewFile(pbp, sf)
		if err := r.File(df); err != nil {
			return errors.WithStack(err)
		}

		os.Chmod(pbp, 0555)

		return nil
	}
}

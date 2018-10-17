package install

import (
	"bytes"
	"go/build"
	"os"
	"path/filepath"

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
		if len(p.GoGet) == 0 {
			continue
		}
		g.RunFn(gotools.Install(p.GoGet))
		if opts.Vendor {
			g.RunFn(pRun(proot, p))
		}
	}

	bb := &bytes.Buffer{}
	plugs := plugdeps.New()
	plugs.Add(opts.Plugins...)
	if err := plugs.Encode(bb); err != nil {
		return g, errors.WithStack(err)
	}

	cpath := filepath.Join(opts.App.Root, "config", "buffalo-plugins.toml")
	g.File(genny.NewFile(cpath, bb))

	return g, nil
}

func pRun(proot string, p plugdeps.Plugin) genny.RunFn {
	return func(r *genny.Runner) error {
		c := build.Default
		if c.GOOS == "windows" {
			return errors.New("vendoring of plugins is currently not supported on windows. PRs are VERY welcome! :)")
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

package install

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-plugins/genny/add"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	aopts := &add.Options{
		App:     opts.App,
		Plugins: opts.Plugins,
	}

	if err := aopts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	proot := filepath.Join(opts.App.Root, "plugins")
	for _, p := range opts.Plugins {
		if len(p.GoGet) == 0 {
			continue
		}

		var args []string
		if len(p.Tags) > 0 {
			args = append(args, "-tags", p.Tags.String())
		}
		g.Command(gogen.Get(p.GoGet, args...))
		if opts.Vendor {
			g.RunFn(pRun(proot, p))
		}
	}
	gg.Add(g)

	g, err := add.New(aopts)
	if err != nil {
		return gg, errors.WithStack(err)
	}

	gg.Add(g)

	return gg, nil
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

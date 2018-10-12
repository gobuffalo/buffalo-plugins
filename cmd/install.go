package cmd

import (
	"bytes"
	"context"
	"io"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/buffalo-plugins/genny/install"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var installOptions = struct {
	dryRun bool
}{}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "installs plugins listed in config/buffalo-plugins.toml to ./plugins/",
	RunE: func(cmd *cobra.Command, args []string) error {
		run := genny.WetRunner(context.Background())
		if installOptions.dryRun {
			run = genny.DryRunner(context.Background())
			run.FileFn = func(f genny.File) (genny.File, error) {
				bb := &bytes.Buffer{}
				if _, err := io.Copy(bb, f); err != nil {
					return f, errors.WithStack(err)
				}
				return genny.NewFile(f.Name(), bb), nil
			}
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil {
			return errors.WithStack(err)
		}

		err = run.WithNew(install.New(&install.Options{
			App:     app,
			Plugins: plugs.Plugins,
		}))
		if err != nil {
			return errors.WithStack(err)
		}

		run.WithRun(func(r *genny.Runner) error {
			bb := &bytes.Buffer{}
			if err := toml.NewEncoder(bb).Encode(plugs); err != nil {
				return errors.WithStack(err)
			}
			cpath := filepath.Join(app.Root, "config", "buffalo-plugins.toml")
			return r.File(genny.NewFile(cpath, bb))
		})

		return run.Run()
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installOptions.dryRun, "dry-run", "d", false, "dry run")
	pluginsCmd.AddCommand(installCmd)
}

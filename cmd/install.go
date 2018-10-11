package cmd

import (
	"bytes"
	"context"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
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
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil {
			return errors.WithStack(err)
		}

		proot := filepath.Join(app.Root, "plugins")
		if err := os.MkdirAll(proot, 0755); err != nil {
			return errors.WithStack(err)
		}

		for _, p := range plugs.Plugins {
			err := func(p plugdeps.Plugin) error {
				p.GoGet = strings.TrimSpace(p.GoGet)
				if len(p.GoGet) == 0 {
					return errors.Errorf("go get instructions missing for %s", p.Binary)
				}
				run.WithRun(func(r *genny.Runner) error {
					if err := gotools.Install(p.GoGet)(r); err != nil {
						return errors.WithStack(err)
					}

					c := build.Default

					sf, err := os.Open(filepath.Join(c.GOPATH, "bin", p.Binary))
					if err != nil {
						return errors.WithStack(err)
					}
					defer sf.Close()

					bpath := filepath.Join(proot, p.Binary)
					os.Remove(bpath)
					df, err := os.OpenFile(bpath, os.O_RDWR|os.O_CREATE, 0555)
					if err != nil {
						return errors.WithStack(err)
					}
					defer df.Close()

					_, err = io.Copy(df, sf)
					if err != nil {
						return errors.WithStack(err)
					}
					return nil
				})
				return nil
			}(p)

			if err != nil {
				return errors.WithStack(err)
			}
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

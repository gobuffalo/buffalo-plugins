package cmd

import (
	"bytes"
	"context"
	"path"
	"strings"

	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var removeOptions = struct {
	dryRun bool
	vendor bool
}{}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes plugin from config/buffalo-plugins.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify at least one package")
		}
		run := genny.WetRunner(context.Background())
		if removeOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil && (errors.Cause(err) != plugdeps.ErrMissingConfig) {
			return errors.WithStack(err)
		}

		for _, a := range args {
			a = strings.TrimSpace(a)
			bin := path.Base(a)
			plugs.Remove(plugdeps.Plugin{
				Binary: bin,
				GoGet:  a,
			})
		}

		run.WithRun(func(r *genny.Runner) error {
			p := plugdeps.ConfigPath(app)
			bb := &bytes.Buffer{}
			if err := plugs.Encode(bb); err != nil {
				return errors.WithStack(err)
			}
			return r.File(genny.NewFile(p, bb))
		})
		if err != nil {
			return errors.WithStack(err)
		}

		return run.Run()
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeOptions.dryRun, "dry-run", "d", false, "dry run")
	removeCmd.Flags().BoolVar(&removeOptions.vendor, "vendor", false, "will install plugin binaries into ./plugins [WINDOWS not currently supported]")
}

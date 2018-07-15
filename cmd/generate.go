package cmd

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateOptions = struct {
	*plugin.Options
	dryRun bool
}{
	Options: &plugin.Options{},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "plugin",
	Short: "generates a new buffalo plugin",
	Long:  "buffalo generate plugin github.com/foo/buffalo-bar",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			generateOptions.PluginPkg = args[0]
		}
		if err := plugin.NormalizeOptions(generateOptions.Options); err != nil {
			return errors.WithStack(err)
		}

		r := genny.WetRunner(context.Background())
		if generateOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}
		r.Root = filepath.Join(envy.GoPath(), "src")
		r.Root = filepath.Join(r.Root, generateOptions.Options.PluginPkg)

		g, err := plugin.New(generateOptions.Options)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gotools.GoFmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	generateCmd.Flags().StringVarP(&generateOptions.Author, "author", "a", "", "author's name")
	generateCmd.Flags().StringVarP(&generateOptions.ShortName, "short-name", "s", "", "a 'short' name for the package")
	rootCmd.AddCommand(generateCmd)
}

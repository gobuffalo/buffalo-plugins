package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-plugins/genny/plugin"
	"github.com/gobuffalo/buffalo-plugins/genny/plugin/with"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/gobuffalo/genny/movinglater/licenser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "plugin",
	Short: "generates a new buffalo plugin",
	Long:  "buffalo generate plugin github.com/foo/buffalo-bar",
	RunE: func(cmd *cobra.Command, args []string) error {
		popts := &plugin.Options{
			Author:    viper.GetString("author"),
			ShortName: viper.GetString("short-name"),
		}
		if len(args) > 0 {
			popts.PluginPkg = args[0]
		}
		if err := plugin.NormalizeOptions(popts); err != nil {
			return errors.WithStack(err)
		}

		r := genny.WetRunner(context.Background())
		if viper.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}
		r.Root = filepath.Join(envy.GoPath(), "src")
		r.Root = filepath.Join(r.Root, popts.PluginPkg)
		r.WithRun(genny.Force(r.Root, viper.GetBool("force")))

		g, err := plugin.New(popts)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		if viper.GetBool("with-gen") {
			g, err = with.GenerateCmd(popts)
			if err != nil {
				return errors.WithStack(err)
			}
			r.With(g)
		}

		lopts := &licenser.Options{
			Author: viper.GetString("author"),
			Name:   viper.GetString("license"),
		}

		g, err = licenser.New(lopts)
		if err := licenser.NormalizeOptions(lopts); err != nil {
			return errors.WithStack(err)
		}

		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gotools.GoFmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		gg, err := gomods.New(popts.PluginPkg, r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		gg.With(r)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolP("dry-run", "d", false, "run the generator without creating files or running commands")
	generateCmd.Flags().Bool("with-gen", false, "creates a generator plugin")
	generateCmd.Flags().BoolP("force", "f", false, "will delete the target directory if it exists")
	generateCmd.Flags().StringP("author", "a", "", "author's name")
	generateCmd.Flags().StringP("license", "l", "mit", fmt.Sprintf("choose a license from: [%s]", strings.Join(licenser.Available, ", ")))
	generateCmd.Flags().StringP("short-name", "s", "", "a 'short' name for the package")
	viper.BindPFlags(generateCmd.Flags())
	rootCmd.AddCommand(generateCmd)
}

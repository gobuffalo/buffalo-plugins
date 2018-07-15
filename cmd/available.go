package cmd

import (
	"encoding/json"
	"os"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "tools for working with buffalo plugins",
}

// availableCmd represents the available command
var availableCmd = &cobra.Command{
	Use:   "available",
	Short: "a list of available buffalo plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		plugs := plugins.Commands{
			{Name: generateCmd.Use, BuffaloCommand: "generate", Description: generateCmd.Short, Aliases: generateCmd.Aliases},
			{Name: pluginsCmd.Use, BuffaloCommand: "root", Description: pluginsCmd.Short, Aliases: pluginsCmd.Aliases},
		}
		return json.NewEncoder(os.Stdout).Encode(plugs)
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
	rootCmd.AddCommand(availableCmd)
}

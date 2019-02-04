package cache

import (
	"os"

	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "rebuilds the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		os.RemoveAll(plugins.CachePath)
		_, err := plugins.Available()
		return err
	},
}

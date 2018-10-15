package cmd

import (
	"encoding/json"

	"github.com/gobuffalo/events"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listens to github.com/gobuffalo/events",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("must pass a payload")
		}

		e := events.Event{}
		err := json.Unmarshal([]byte(args[0]), &e)
		if err != nil {
			return errors.WithStack(err)
		}

		if e.Kind != "buffalo:setup:started" {
			return nil
		}

		return installCmd.RunE(cmd, []string{})
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

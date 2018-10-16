package cmd

import (
	"context"
	"encoding/json"

	"github.com/gobuffalo/buffalo-plugins/genny/install"
	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny"
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

		run := genny.WetRunner(context.Background())

		opts := &install.Options{}
		err = run.WithNew(install.New(opts))
		if err != nil {
			return errors.WithStack(err)
		}
		payload := e.Payload
		payload["plugins"] = opts.Plugins
		events.EmitPayload(plugins.EvtSetupStarted, payload)
		if err := run.Run(); err != nil {
			events.EmitError(plugins.EvtSetupErr, err, payload)
			return errors.WithStack(err)
		}
		events.EmitPayload(plugins.EvtSetupFinished, payload)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

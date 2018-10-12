package install

import (
	"os"

	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

type Options struct {
	App     meta.App
	Plugins []plugdeps.Plugin
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		opts.App = meta.New(pwd)
	}
	return nil
}

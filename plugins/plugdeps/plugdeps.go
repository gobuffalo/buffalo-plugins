package plugdeps

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

// ErrMissingConfig is if config/buffalo-plugins.toml file is not found. Use plugdeps#On(app) to test if plugdeps are being used
var ErrMissingConfig = errors.Errorf("could not find a buffalo-plugins config file at %s", ConfigPath(meta.New(".")))

// List all of the plugins the application depeneds on. Will return ErrMissingConfig
// if the app is not using config/buffalo-plugins.toml to manage their plugins.
// Use plugdeps#On(app) to test if plugdeps are being used.
func List(app meta.App) (*Plugins, error) {
	plugs := New()
	if app.WithPop {
		plugs.Add(pop)
	}

	if !On(app) {
		return plugs, ErrMissingConfig
	}

	p := ConfigPath(app)
	tf, err := os.Open(p)
	if err != nil {
		return plugs, errors.WithStack(err)
	}
	if err := plugs.Decode(tf); err != nil {
		return plugs, errors.WithStack(err)
	}

	return plugs, nil
}

// ConfigPath returns the path to the config/buffalo-plugins.toml file
// relative to the app
func ConfigPath(app meta.App) string {
	return filepath.Join(app.Root, "config", "buffalo-plugins.toml")
}

// On checks for the existence of config/buffalo-plugins.toml if this
// file exists its contents will be used to list plugins. If the file is not
// found, then the BUFFALO_PLUGIN_PATH and ./plugins folders are consulted.
func On(app meta.App) bool {
	_, err := os.Stat(ConfigPath(app))
	return err == nil
}

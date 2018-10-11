package plugdeps

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/meta"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type Plugin struct {
	Binary string `toml:"binary"`
	GoGet  string `toml:"go_get"`
}

type Plugins struct {
	Plugins []Plugin `toml:"plugin"`
}

func List(app meta.App) (Plugins, error) {
	plugs := Plugins{}
	tfp := filepath.Join(app.Root, "config", "buffalo-plugins.toml")
	tf, err := os.Open(tfp)
	if err != nil {
		if app.WithPop {
			plugs.Plugins = append(plugs.Plugins, Plugin{
				Binary: "buffalo-pop",
				GoGet:  "github.com/gobuffalo/buffalo-pop",
			})
		}
		return plugs, nil
	}
	if err := toml.NewDecoder(tf).Decode(&plugs); err != nil {
		return plugs, errors.WithStack(err)
	}
	return plugs, nil
}

package plugdeps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_Plugin(t *testing.T) {
	r := require.New(t)

	plugs := Plugins{
		Plugins: []Plugin{
			{
				Binary: "buffalo-pop",
				GoGet:  "github.com/gobuffalo/buffalo-pop",
			},
			{
				Binary: "buffalo-heroku",
				GoGet:  "github.com/gobuffalo/buffalo-heroku",
			},
		},
	}

	err := toml.NewEncoder(os.Stdout).Encode(plugs)
	r.NoError(err)
}

func Test_List(t *testing.T) {
	r := require.New(t)

	app := meta.App{}
	plugs, err := List(app)
	r.NoError(err)
	r.Len(plugs.Plugins, 0)

	app.WithPop = true

	plugs, err = List(app)
	r.NoError(err)
	r.Len(plugs.Plugins, 1)

	p := filepath.Join(os.TempDir(), "config")
	r.NoError(os.MkdirAll(p, 0755))
	f, err := os.Create(filepath.Join(p, "buffalo-plugins.toml"))
	r.NoError(err)
	f.WriteString(eToml)
	r.NoError(f.Close())

	app = meta.New(os.TempDir())
	plugs, err = List(app)
	r.NoError(err)
	r.Len(plugs.Plugins, 2)
}

const eToml = `[[plugin]]
  binary = "buffalo-pop"
  go_get = "github.com/gobuffalo/buffalo-pop"

[[plugin]]
  binary = "buffalo-heroku"
  go_get = "github.com/gobuffalo/buffalo-heroku"
`

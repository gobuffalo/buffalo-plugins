package plugdeps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var heroku = Plugin{
	Binary: "buffalo-heroku",
	GoGet:  "github.com/gobuffalo/buffalo-heroku",
	Commands: []Command{
		{Name: "deploy", Flags: []string{"-v"}},
	},
	Tags: []string{"foo", "bar"},
}

var local = Plugin{
	Binary: "buffalo-hello.rb",
	Local:  "./plugins/buffalo-hello.rb",
}

func Test_ConfigPath(t *testing.T) {
	r := require.New(t)

	x := ConfigPath(meta.App{Root: "foo"})
	r.Equal(x, filepath.Join("foo", "config", "buffalo-plugins.toml"))
}

func Test_List_Off(t *testing.T) {
	r := require.New(t)

	app := meta.App{}
	plugs, err := List(app)
	r.Error(err)
	r.Equal(errors.Cause(err), ErrMissingConfig)
	r.Len(plugs.List(), 1)
}

func Test_List_On(t *testing.T) {
	r := require.New(t)

	app := meta.New(os.TempDir())

	p := ConfigPath(app)
	r.NoError(os.MkdirAll(filepath.Dir(p), 0755))
	f, err := os.Create(p)
	r.NoError(err)
	f.WriteString(eToml)
	r.NoError(f.Close())

	plugs, err := List(app)
	r.NoError(err)
	r.Len(plugs.List(), 4)
}

const eToml = `[[plugin]]
  binary = "buffalo-hello.rb"
  local = "./plugins/buffalo-hello.rb"

[[plugin]]
  binary = "buffalo-heroku"
  go_get = "github.com/gobuffalo/buffalo-heroku"
  tags = ["foo", "bar"]

  [[plugin.command]]
    name = "deploy"
    flags = ["-v"]

[[plugin]]
  binary = "buffalo-plugins"
  go_get = "github.com/gobuffalo/buffalo-plugins"

[[plugin]]
  binary = "buffalo-pop"
  go_get = "github.com/gobuffalo/buffalo-pop"
`

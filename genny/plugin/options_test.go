package plugin

import (
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_normalizeOptions(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := NormalizeOptions(opts)
	r.Error(err)

	opts.PluginPkg = "github.com/foo/bar"

	err = NormalizeOptions(opts)
	r.NoError(err)
	r.Equal("github.com/foo/buffalo-bar", opts.PluginPkg)

	year := time.Now().Year()
	r.Equal(opts.Year, year)

	u, err := user.Current()
	r.NoError(err)
	r.Equal(u.Name, opts.Author)
}

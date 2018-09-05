package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	r := require.New(t)
	opts := &Options{
		PluginPkg: "github.com/foo/buffalo-bar",
		Year:      1999,
		Author:    "Homer Simpson",
		ShortName: "bar",
	}
	run := genny.DryRunner(context.Background())
	run.Root = os.TempDir()

	gg, err := New(opts)
	r.NoError(err)
	run.WithGroup(gg)

	r.NoError(run.Run())
	res := run.Results()
	if gomods.On() {
		r.Len(res.Commands, 2)

		c := res.Commands[0]
		r.Equal("go mod init github.com/foo/buffalo-bar", strings.Join(c.Args, " "))

		c = res.Commands[1]
		r.Equal("go mod tidy", strings.Join(c.Args, " "))
	} else {
		r.Len(res.Commands, 0)
	}

	r.Len(res.Files, 11)

	for i, f := range res.Files {
		if i == 0 {
			fmt.Printf("f := res.Files[%d]\n", i)
		} else {
			fmt.Printf("f = res.Files[%d]\n", i)
		}
		fmt.Printf(`r.Equal("%s", f.Name())`+"\n\n", f.Name())
	}
	f := res.Files[0]
	r.Equal(".goreleaser.yml", f.Name())

	f = res.Files[1]
	r.Equal(".travis.yml", f.Name())

	f = res.Files[2]
	r.Equal("LICENSE", f.Name())

	f = res.Files[3]
	r.Equal("Makefile", f.Name())

	f = res.Files[4]
	r.Equal("README.md", f.Name())
	r.Contains(f.String(), opts.PluginPkg)

	f = res.Files[5]
	r.Equal("bar/version.go", f.Name())

	f = res.Files[6]
	r.Equal("cmd/available.go", f.Name())

	f = res.Files[7]
	r.Equal("cmd/bar.go", f.Name())

	f = res.Files[8]
	r.Equal("cmd/root.go", f.Name())

	f = res.Files[9]
	r.Equal("cmd/version.go", f.Name())
	r.Contains(f.String(), opts.PluginPkg+"/"+opts.ShortName)
	r.Contains(f.String(), opts.ShortName+".Version")

	f = res.Files[10]
	r.Equal("main.go", f.Name())
	r.Contains(f.String(), "github.com/foo/buffalo-bar/cmd")

}

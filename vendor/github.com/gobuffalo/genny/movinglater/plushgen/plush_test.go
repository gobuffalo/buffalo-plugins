package plushgen

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_Transformer(t *testing.T) {
	r := require.New(t)

	ctx := plush.NewContext()
	ctx.Set("name", "mark")
	f := genny.NewFile("foo.plush.txt", strings.NewReader("Hello <%= name %>"))

	tr := Transformer(ctx)
	f, err := tr.Transform(f)
	r.NoError(err)
	r.Equal("foo.txt", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("Hello mark", string(b))
}

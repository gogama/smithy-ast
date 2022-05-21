package prelude

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/gogama/smithy-ast/ast"
	"github.com/stretchr/testify/require"
)

// TestMinifyAndGZIP is both a test AND a crucial part of the build
// process, since it generates the GZIPped minified prelude model JSON
// which is embedded in prelude.go.
func TestMinifyAndGZIP(t *testing.T) {
	raw, err := os.Open("prelude.json")
	require.NoError(t, err)
	defer func() {
		_ = raw.Close()
	}()

	m, err := ast.ReadModel(raw)
	require.NoError(t, err)

	min, err := os.OpenFile("prelude_min.json", os.O_WRONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	defer func() {
		_ = min.Close()
	}()

	err = ast.WriteModel(m, min)
	require.NoError(t, err)

	gzf, err := os.OpenFile("prelude_min.json.gz", os.O_WRONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	defer func() {
		_ = gzf.Close()
	}()

	gz := gzip.NewWriter(gzf)
	defer func() {
		_ = gz.Close()
	}()
	err = ast.WriteModel(m, gz)
	require.NoError(t, err)
}

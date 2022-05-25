package prelude

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestNewReader(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		r := NewReader()
		require.NotNil(t, r)

		data, err := io.ReadAll(r)
		require.NoError(t, err)

		v := json.Valid(data)
		assert.True(t, v)
	})

	t.Run("Parseable AST", func(t *testing.T) {
		r := NewReader()
		require.NotNil(t, r)

		m, err := ast.ReadModel(r)
		require.NoError(t, err)

		assert.Equal(t, "1.0", m.Version.Value)
	})
}

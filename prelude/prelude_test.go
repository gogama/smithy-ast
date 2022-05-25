package prelude

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogama/smithy-ast/ast"
	"github.com/stretchr/testify/require"
)

// TestMinifyAndGZIP is both a test AND a crucial part of the build
// process, since it generates the GZIPped minified prelude model JSON
// which is embedded in prelude.go.
//
// This "test" is separated into a "Pipeline" of sub-steps to ensure
// that the output of each sub-step is valid input to the next sub-step.

func TestMinifyAndGZIP(t *testing.T) {
	var m1, m2 ast.Model

	t.Run("Minify", func(t *testing.T) {
		raw, err := os.Open("prelude.json")
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = raw.Close()
		})

		m1, err = ast.ReadModel(raw)
		require.NoError(t, err)

		min, err := os.OpenFile("prelude_min.json", os.O_WRONLY|os.O_CREATE, 0644)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = min.Close()
		})

		err = ast.WriteModel(m1, min)
		require.NoError(t, err)
	})

	t.Run("GZIP", func(t *testing.T) {
		min, err := os.Open("prelude_min.json")
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = min.Close()
		})

		m2, err = ast.ReadModel(min)
		require.NoError(t, err)
		require.True(t, reflect.DeepEqual(m1, m2), "m1 and m2 must be deeply equal")

		gzf, err := os.OpenFile("prelude_min.json.gz", os.O_WRONLY|os.O_CREATE, 0644)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = gzf.Close()
		})

		gz := gzip.NewWriter(gzf)
		t.Cleanup(func() {
			_ = gz.Close()
		})
		err = ast.WriteModel(m2, gz)
		require.NoError(t, err)
	})

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

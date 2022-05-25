package prelude

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
)

//go:embed prelude_min.json.gz
var gzipJSON []byte

// NewReader returns a new Reader reading the prelude model AST JSON.
func NewReader() io.Reader {
	br := bytes.NewReader(gzipJSON)
	gzr, err := gzip.NewReader(br)
	if err != nil {
		panic(err)
	}
	return gzr
}

package prelude

import (
	_ "embed"
)

//go:embed prelude_min.json.gz
var gzipJSON []byte

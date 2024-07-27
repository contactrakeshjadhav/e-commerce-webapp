package swagger

import (
	_ "embed"
)

// We embed the json file to be compile in build time
//
//go:embed docs/swagger.json
var JsonFile []byte

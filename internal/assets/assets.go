package assets

import _ "embed"

//go:embed schema.json
var SchemaJSON []byte

//go:embed dict.yaml
var DictYAML []byte

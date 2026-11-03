// Package api embeds the reviewed API contract served by the HTTP adapter.
package api

import _ "embed"

//go:embed openapi.yaml
var Document []byte

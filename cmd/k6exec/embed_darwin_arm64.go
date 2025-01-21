//go:build darwin && arm64
// +build darwin,arm64

package main

import "embed"

//go:embed k6/darwin/arm64
var binFS embed.FS

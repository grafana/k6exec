//go:build windows && arm64
// +build windows,arm64

package main

import "embed"

//go:embed k6/windows/arm64
var binFS embed.FS

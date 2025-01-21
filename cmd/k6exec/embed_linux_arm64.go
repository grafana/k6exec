//go:build linux && arm64
// +build linux,arm64

package main

import "embed"

//go:embed k6/linux/arm64
var binFS embed.FS

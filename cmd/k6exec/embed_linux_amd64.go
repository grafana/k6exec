//go:build linux && amd64
// +build linux,amd64

package main

import "embed"

//go:embed k6/linux/amd64
var binFS embed.FS

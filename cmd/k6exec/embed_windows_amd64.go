//go:build windows && amd64
// +build windows,amd64

package main

import "embed"

//go:embed k6/windows/amd64
var binFS embed.FS

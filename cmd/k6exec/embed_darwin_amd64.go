//go:build darwin && amd64
// +build darwin,amd64

package main

import "embed"

//go:embed k6/darwin/amd64
var binFS embed.FS

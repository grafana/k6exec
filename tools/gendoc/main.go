// Package main contains CLI documentation generator tool.
package main

import (
	"github.com/grafana/clireadme"
	"github.com/grafana/k6exec/cmd"
)

func main() {
	root := cmd.New(nil)
	clireadme.Main(root, 1)
}

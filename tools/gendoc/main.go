// Package main contains CLI documentation generator tool.
package main

import (
	_ "embed"
	"strings"

	"github.com/grafana/clireadme"
	"github.com/grafana/k6exec/cmd"
)

func main() {
	root := cmd.New(nil)
	root.Use = strings.ReplaceAll(root.Use, "exec", "k6exec")
	clireadme.Main(root, 1)
}

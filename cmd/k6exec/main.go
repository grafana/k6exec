// Package main contains the main function for k6exec CLI tool.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/grafana/k6exec/cmd"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	appname = "k6exec"
	version = "dev"
)

func main() {
	runCmd(newCmd(os.Args[1:])) //nolint:forbidigo
}

func newCmd(args []string) *cobra.Command {
	cmd := cmd.New()
	cmd.Use = strings.Replace(cmd.Use, cmd.Name(), appname, 1)
	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	log.SetFlags(0)
	log.Writer()

	if err := cmd.Execute(); err != nil {
		log.Fatal(formatError(err))
	}
}

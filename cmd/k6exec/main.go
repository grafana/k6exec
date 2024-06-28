// Package main contains the main function for k6exec CLI tool.
package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grafana/k6exec/cmd"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	appname = "k6exec"
	version = "dev"
)

func initLogging() *slog.LevelVar {
	var levelVar = new(slog.LevelVar)

	logrus.SetLevel(logrus.DebugLevel)

	logger := slog.New(sloglogrus.Option{Level: levelVar}.NewLogrusHandler())
	logger = logger.With("app", appname)

	slog.SetDefault(logger)

	return levelVar
}

func main() {
	levelVar := initLogging()
	runCmd(newCmd(os.Args[1:], levelVar)) //nolint:forbidigo
}

func newCmd(args []string, levelVar *slog.LevelVar) *cobra.Command {
	cmd := cmd.New(levelVar)
	cmd.Use = strings.Replace(cmd.Use, cmd.Name(), appname, 1)
	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		slog.Error(formatError(err))
		os.Exit(1) //nolint:forbidigo
	}
}

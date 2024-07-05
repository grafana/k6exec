// Package main contains the main function for k6exec.
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

func initLogging(app string) *slog.LevelVar {
	levelVar := new(slog.LevelVar)

	logrus.SetLevel(logrus.DebugLevel)

	logger := slog.New(sloglogrus.Option{Level: levelVar}.NewLogrusHandler())
	logger = logger.With("app", app)

	slog.SetDefault(logger)

	return levelVar
}

func main() {
	runCmd(newCmd(os.Args[1:], initLogging(appname))) //nolint:forbidigo
}

func newCmd(args []string, levelVar *slog.LevelVar) *cobra.Command {
	cmd := cmd.New(levelVar)
	cmd.Use = strings.ReplaceAll(cmd.Use, "exec", appname)
	cmd.Version = version

	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		args[0] = "help"
	}

	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		slog.Error(formatError(err))
		os.Exit(1) //nolint:forbidigo
	}
}

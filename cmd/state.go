package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/exec"

	"github.com/grafana/k6exec"
	"github.com/spf13/cobra"
)

const (
	defaultBuildServiceURL = "https://ingest.k6.io/builder/api/v1"
)

type state struct {
	k6exec.Options
	buildServiceURL string
	verbose         bool
	quiet           bool
	nocolor         bool
	version         bool
	usage           bool
	levelVar        *slog.LevelVar
	cmd             *exec.Cmd
	cleanup         func() error
	configFile      string
}

func newState(levelVar *slog.LevelVar) *state {
	s := new(state)

	s.levelVar = levelVar

	return s
}

func (s *state) persistentPreRunE(cmd *cobra.Command, _ []string) error {
	var err error

	// get URL to build service: first provided from flag, then from environment variable, then default
	buildServiceURL := s.buildServiceURL

	if len(buildServiceURL) == 0 {
		buildServiceURL = os.Getenv("K6_BUILD_SERVICE_URL") //nolint:forbidigo
	}
	if len(buildServiceURL) == 0 {
		buildServiceURL = defaultBuildServiceURL
	}

	s.Options.BuildServiceURL = buildServiceURL

	// get authorization token for the build service
	auth := os.Getenv("K6_CLOUD_TOKEN") //nolint:forbidigo

	if len(auth) == 0 {
		// allow overriding the config file for testing
		configFile := s.configFile
		if configFile == "" {
			// check if the command has a 'config' flag and get the value
			configFile, err = getFlagValue(cmd, "--config", "-c")
			if err != nil {
				return err
			}
		}

		config, err := loadConfig(configFile)
		if err != nil {
			return err
		}

		auth = config.Collectors.Cloud.Token
	}

	s.Options.BuildServiceToken = auth

	if s.verbose && s.levelVar != nil {
		s.levelVar.Set(slog.LevelDebug)
	}

	return nil
}

func (s *state) preRunE(sub *cobra.Command, args []string) error {
	cmdargs := make([]string, 0, len(args))

	if sub.Name() != s.Options.AppName {
		cmdargs = append(cmdargs, sub.Name())
	}

	if s.version {
		cmdargs = append(cmdargs, "--version")
	}

	if s.verbose {
		cmdargs = append(cmdargs, "-v")
	}

	if s.quiet {
		cmdargs = append(cmdargs, "-q")
	}

	if s.nocolor {
		cmdargs = append(cmdargs, "--no-color")
	}

	if subargs := getArgs(sub); subargs != nil {
		cmdargs = append(cmdargs, subargs...)
	} else {
		cmdargs = append(cmdargs, args...)
	}

	ctx := sub.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cmd, cleanup, err := k6exec.Command(ctx, cmdargs, &s.Options)
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stderr //nolint:forbidigo
	cmd.Stdout = os.Stdout //nolint:forbidigo
	cmd.Stdin = os.Stdin   //nolint:forbidigo

	s.cmd = cmd
	s.cleanup = cleanup

	return nil
}

func (s *state) runE(_ *cobra.Command, _ []string) error {
	var err error

	// FIXME: I think this code is not setting the error to the cleanup function (pablochacin)
	defer func() {
		e := s.cleanup()
		if err == nil {
			err = e
		}
	}()

	slog.Debug("running", "k6 binary", s.cmd.Path, "args", s.cmd.Args[1:])
	err = s.cmd.Run()

	return err
}

func (s *state) helpFunc(cmd *cobra.Command, args []string) {
	err := s.preRunE(cmd, append(args, "-h"))
	if err != nil {
		cmd.PrintErr(err)
		// FIXME: added this return because in case of error provisioning the binary,
		// it doesn't make sense to continue (pablochacin)
		return
	}

	err = s.runE(cmd, args)
	if err != nil {
		cmd.PrintErr(err)
	}
}

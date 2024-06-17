// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"os"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6exec"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

// New creates new cobra command for exec command.
func New() *cobra.Command {
	root := &cobra.Command{
		Use:               "exec",
		Short:             "Launching k6 with extensions",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	for _, cmd := range subcommands() {
		cmd := cmd
		cmd.RunE = exec
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		cmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}
		cmd.SetHelpFunc(usage)

		root.AddCommand(&cmd)
	}

	return root
}

func usage(cmd *cobra.Command, args []string) {
	err := exec(cmd, append(args, "-h"))
	if err != nil {
		cmd.PrintErr(err)
	}
}

func exec(sub *cobra.Command, args []string) error {
	var (
		deps k6deps.Dependencies
		err  error
	)

	if scriptname, hasScript := scriptArg(sub, args); hasScript {
		deps, err = k6deps.Analyze(&k6deps.Options{
			Script: k6deps.Source{
				Name: scriptname,
			},
		})
		if err != nil {
			return err
		}
	}

	args = append([]string{sub.Name()}, args...)

	cmd, err := k6exec.Command(context.TODO(), args, deps, nil)
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stderr //nolint:forbidigo
	cmd.Stdout = os.Stdout //nolint:forbidigo
	cmd.Stdin = os.Stdin   //nolint:forbidigo

	return cmd.Run()
}

func scriptArg(cmd *cobra.Command, args []string) (string, bool) {
	if len(cmd.Annotations) == 0 {
		return "", false
	}

	if _, use := cmd.Annotations[useExtensions]; !use {
		return "", false
	}

	if len(args) == 0 {
		return "", false
	}

	last := args[len(args)-1]
	if len(last) == 0 || last[0] == '-' {
		return "", false
	}

	return last, true
}

func subcommands() []cobra.Command {
	annext := map[string]string{useExtensions: "true"}

	return []cobra.Command{
		{Use: "help", Short: "Help about any command"},
		{Use: "resume", Short: "Resume a paused test"},
		{Use: "scale", Short: "Scale a running test"},
		{Use: "cloud", Short: "Run a test on the cloud"},
		{Use: "completion", Short: "Generate the autocompletion script for the specified shell"},
		{Use: "inspect", Short: "Inspect a script or archive"},
		{Use: "pause", Short: "Pause a running test"},
		{Use: "status", Short: "Show test status"},
		{Use: "login", Short: "Authenticate with a service"},
		{Use: "stats", Short: "Show test metrics"},
		{Use: "version", Short: "Show application version"},
		{Use: "new", Short: "Create and initialize a new k6 script"},
		{Use: "run", Short: "Start a test", Annotations: annext},
		{Use: "archive", Short: "Create an archive", Annotations: annext},
	}
}

const useExtensions = "useExtensions"

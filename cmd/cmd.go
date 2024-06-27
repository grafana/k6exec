// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"net/url"
	"os"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6exec"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	k6exec.Options
	buildServiceURL     string
	extensionCatalogURL string
}

func (o *options) postProcess() error {
	if len(o.buildServiceURL) > 0 {
		val, err := url.Parse(o.buildServiceURL)
		if err != nil {
			return err
		}

		o.BuildServiceURL = val
	}

	if len(o.extensionCatalogURL) > 0 {
		val, err := url.Parse(o.extensionCatalogURL)
		if err != nil {
			return err
		}

		o.ExtensionCatalogURL = val
	}

	return nil
}

//nolint:forbidigo
func (o *options) init() {
	if value, found := os.LookupEnv("K6_BUILD_SERVICE_URL"); found {
		o.buildServiceURL = value
	}

	if value, found := os.LookupEnv("K6_EXTENSION_CATALOG_URL"); found {
		o.extensionCatalogURL = value
	}
}

// New creates new cobra command for exec command.
func New() *cobra.Command {
	opts := new(options)

	opts.init()

	root := &cobra.Command{
		Use:               "exec [flags] [command]",
		Short:             "Lanch k6 with extensions",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE:              func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return opts.postProcess() },
	}

	root.AddCommand(subcommands(&opts.Options)...)

	flags := root.PersistentFlags()

	flags.StringVar(
		&opts.extensionCatalogURL,
		"extension-catalog-url",
		opts.extensionCatalogURL,
		"URL of the k6 extension catalog to be used",
	)
	flags.StringVar(
		&opts.buildServiceURL,
		"build-service-url",
		opts.buildServiceURL,
		"URL of the k6 build service to be used",
	)
	flags.BoolVarP(
		&opts.Verbose,
		"verbose",
		"v",
		false,
		"enable verbose logging",
	)

	root.MarkFlagsMutuallyExclusive("extension-catalog-url", "build-service-url")

	return root
}

func usage(cmd *cobra.Command, args []string) {
	err := exec(cmd, append(args, "-h"), nil)
	if err != nil {
		cmd.PrintErr(err)
	}
}

func exec(sub *cobra.Command, args []string, opts *k6exec.Options) error {
	var (
		deps  k6deps.Dependencies
		err   error
		dopts k6deps.Options
	)

	if scriptname, hasScript := scriptArg(sub, args); hasScript {
		dopts.Script.Name = scriptname
	}

	dopts.Manifest.Name = "package.json"

	deps, err = k6deps.Analyze(&dopts)
	if err != nil {
		return err
	}

	args = append([]string{sub.Name()}, args...)

	cmd, err := k6exec.Command(context.Background(), args, deps, opts)
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stderr //nolint:forbidigo
	cmd.Stdout = os.Stdout //nolint:forbidigo
	cmd.Stdin = os.Stdin   //nolint:forbidigo

	defer k6exec.CleanupState(opts) //nolint:errcheck

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

func subcommands(opts *k6exec.Options) []*cobra.Command {
	annext := map[string]string{useExtensions: "true"}

	all := make([]*cobra.Command, 0, len(commands))

	for _, name := range commands {
		cmd := &cobra.Command{
			Use:                name,
			RunE:               func(cmd *cobra.Command, args []string) error { return exec(cmd, args, opts) },
			SilenceErrors:      true,
			SilenceUsage:       true,
			FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
			Hidden:             true,
		}
		cmd.SetHelpFunc(usage)

		if name == "run" || name == "archive" {
			cmd.Annotations = annext
		}

		all = append(all, cmd)
	}

	return all
}

const useExtensions = "useExtensions"

var commands = []string{ //nolint:gochecknoglobals
	"help",
	"resume",
	"scale",
	"cloud",
	"completion",
	"inspect",
	"pause",
	"status",
	"login",
	"stats",
	"version",
	"new",
	"run",
	"archive",
}

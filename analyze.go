package k6exec

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grafana/k6deps"
)

func analyze(args []string, opts *Options) (k6deps.Dependencies, error) {
	depsOpts := newDepsOptions(args, opts)

	// we call Analyze before logging because it will return the name of the manifest, in any
	deps, err := k6deps.Analyze(depsOpts)

	slog.Debug("analyzing sources", depsOptsAttrs(depsOpts)...)

	if err == nil && len(deps) > 0 {
		slog.Debug("found dependencies", "deps", deps.String())
	}

	return deps, err
}

func newDepsOptions(args []string, opts *Options) *k6deps.Options {
	dopts := &k6deps.Options{
		Env:          opts.Env,
		Manifest:     opts.Manifest,
		LookupEnv:    opts.LookupEnv,
		FindManifest: opts.FindManifest,
	}

	scriptname, hasScript := scriptArg(args)
	if !hasScript {
		return dopts
	}

	if _, err := os.Stat(scriptname); err != nil { //nolint:forbidigo
		return dopts
	}

	if strings.HasSuffix(scriptname, ".tar") {
		dopts.Archive.Name = scriptname
	} else {
		dopts.Script.Name = scriptname
	}

	return dopts
}

func scriptArg(args []string) (string, bool) {
	if len(args) == 0 {
		return "", false
	}

	cmd := args[0]
	if cmd != "run" && cmd != "archive" && cmd != "inspect" && cmd != "cloud" {
		return "", false
	}

	if len(args) == 1 {
		return "", false
	}

	last := args[len(args)-1]
	if last[0] == '-' {
		return "", false
	}

	return last, true
}

func depsOptsAttrs(opts *k6deps.Options) []any {
	attrs := []any{}

	if opts.Manifest.Name != "" {
		attrs = append(attrs, "Manifest", opts.Manifest.Name)
	}

	if opts.Archive.Name != "" {
		attrs = append(attrs, "Archive", opts.Archive.Name)
	}

	// ignore script if archive is present
	if opts.Archive.Name == "" && opts.Script.Name != "" {
		attrs = append(attrs, "Script", opts.Script.Name)
	}

	if opts.Env.Name != "" {
		attrs = append(attrs, "Env", opts.Env.Name)
	}

	return attrs
}

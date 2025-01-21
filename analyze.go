package k6exec

import (
	"os"
	"strings"

	"github.com/grafana/k6deps"
)

func analyze(args []string, opts *Options) (k6deps.Dependencies, error) {
	return k6deps.Analyze(newDepsOptions(args, opts))
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
	if cmd != "run" && cmd != "archive" && cmd != "inspect" {
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

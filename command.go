package k6exec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/grafana/k6provision"
)

// Command returns the exec.Cmd struct to execute k6 with the given arguments.
// If the given subcommand has a script argument, it analyzes the dependencies
// in the script and provisions a k6 executable based on them.
// In Options, you can also specify environment variable and manifest file as dependency sources.
// If no errors occur, the provisioned k6 executable will be placed in a temporary directory.
// The second return value is a cleanup function that is used to delete this temporary directory.
func Command(ctx context.Context, args []string, opts *Options) (*exec.Cmd, func() error, error) {
	deps, err := analyze(args, opts)
	if err != nil {
		return nil, nil, err
	}

	dir, err := os.MkdirTemp("", "k6exec-*") //nolint:forbidigo
	if err != nil {
		return nil, nil, err
	}

	exe := filepath.Join(dir, k6provision.ExeName)

	if err := provision(ctx, deps, exe, opts); err != nil {
		return nil, nil, err
	}

	cmd := exec.CommandContext(ctx, exe, args...) //nolint:gosec

	return cmd, func() error { return os.RemoveAll(dir) }, nil //nolint:forbidigo
}

/*
	sum := sha256.Sum256([]byte(deps.String()))
	dir := filepath.Join(os.TempDir(), "k6exec-"+hex.EncodeToString(sum[:]))

	if err := os.MkdirAll(dir, 0o700); err != nil { //nolint:forbidigo
		return nil, nil, err
	}

	exe := filepath.Join(dir, k6provision.ExeName)

*/

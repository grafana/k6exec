package k6exec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/grafana/k6deps"
)

// Command returns the exec.Cmd struct to execute k6 with the given dependencies and arguments.
func Command(ctx context.Context, args []string, deps k6deps.Dependencies, opts *Options) (*exec.Cmd, error) {
	dir, err := opts.stateSubdir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrState, err.Error())
	}

	exe := filepath.Join(dir, k6binary)

	loc, err := build(ctx, deps, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBuild, err.Error())
	}

	client, err := opts.client()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDownload, err.Error())
	}

	if err = download(ctx, loc, exe, client); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDownload, err.Error())
	}

	cmd := exec.CommandContext(ctx, exe, args...) //nolint:gosec

	return cmd, nil
}

var (
	// ErrDownload is returned if an error occurs during download.
	ErrDownload = errors.New("download error")
	// ErrBuild is returned if an error occurs during build.
	ErrBuild = errors.New("build error")
	// ErrCache is returned if an error occurs during cache handling.
	ErrCache = errors.New("cache error")
	// ErrState is returned if an error occurs during state handling.
	ErrState = errors.New("state error")
)

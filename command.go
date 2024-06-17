package k6exec

import (
	"context"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/grafana/k6deps"
)

// Command returns the exec.Cmd struct to execute k6 with the given dependencies and arguments.
func Command(ctx context.Context, args []string, deps k6deps.Dependencies, opts *Options) (*exec.Cmd, error) {
	cachedir, err := xdg.CacheFile(opts.appname())
	if err != nil {
		return nil, err
	}

	exe := filepath.Join(cachedir, k6binary)

	mods, err := unmarshalVersionOutput(ctx, exe)
	if err != nil {
		return nil, err
	}

	if !mods.fulfill(deps) {
		demands := mods.merge(deps)

		loc, err := build(demands.Sorted())
		if err != nil {
			return nil, err
		}

		if err = download(ctx, loc, cachedir, opts.client()); err != nil {
			return nil, err
		}
	}

	return exec.Command(exe, args...), nil //nolint:gosec
}

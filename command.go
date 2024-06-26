package k6exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/adrg/xdg"
	"github.com/grafana/k6deps"
)

//nolint:forbidigo
func exists(file string) bool {
	_, err := os.Stat(file)

	return err == nil || !errors.Is(err, os.ErrNotExist)
}

// Command returns the exec.Cmd struct to execute k6 with the given dependencies and arguments.
func Command(ctx context.Context, args []string, deps k6deps.Dependencies, opts *Options) (*exec.Cmd, error) {
	cachedir, err := xdg.CacheFile(opts.appname())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCache, err.Error())
	}

	err = os.MkdirAll(cachedir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCache, err.Error())
	}

	exe := filepath.Join(cachedir, k6binary)

	var mods modules

	if !opts.forceUpdate() && exists(exe) {
		mods, err = unmarshalVersionOutput(ctx, exe)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrCache, err.Error())
		}
	}

	if opts.forceUpdate() || !mods.fulfill(deps) {
		demands := mods.merge(deps)

		loc, err := build(ctx, demands.Sorted(), opts)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrBuild, err.Error())
		}

		if err = download(ctx, loc, exe, opts.client()); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrDownload, err.Error())
		}
	}

	return exec.Command(exe, args...), nil //nolint:gosec
}

var (
	// ErrDownload is returned if an error occurs during download.
	ErrDownload = errors.New("download error")
	// ErrBuild is returned if an error occurs during build.
	ErrBuild = errors.New("build error")
	// ErrCache is returned if an error occurs during cache handling.
	ErrCache = errors.New("cache error")
)

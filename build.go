package k6exec

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/grafana/k6build"
	"github.com/grafana/k6deps"
)

const (
	platform = runtime.GOOS + "/" + runtime.GOARCH

	k6module = "k6"
)

func depsConvert(deps k6deps.Dependencies) (string, []k6build.Dependency) {
	bdeps := make([]k6build.Dependency, 0, len(deps))
	k6constraint := "*"

	for _, dep := range deps {
		if dep.Name == k6module {
			k6constraint = dep.GetConstraints().String()
			continue
		}

		bdeps = append(bdeps, k6build.Dependency{Name: dep.Name, Constraints: dep.GetConstraints().String()})
	}

	return k6constraint, bdeps
}

func build(ctx context.Context, deps k6deps.Dependencies, opts *Options) (*url.URL, error) {
	svc, err := newBuildService(ctx, opts)
	if err != nil {
		return nil, err
	}

	k6constraints, bdeps := depsConvert(deps)

	artifact, err := svc.Build(ctx, platform, k6constraints, bdeps)
	if err != nil {
		return nil, err
	}

	return url.Parse(artifact.URL)
}

func newBuildService(ctx context.Context, opts *Options) (k6build.BuildService, error) {
	if opts.BuildServiceURL != nil {
		return newBuildServiceClient(opts)
	}

	return newLocalBuildService(ctx, opts)
}

func newLocalBuildService(ctx context.Context, opts *Options) (k6build.BuildService, error) {
	statedir, err := opts.stateSubdir()
	if err != nil {
		return nil, err
	}

	catfile := filepath.Join(statedir, "catalog.json")

	client, err := opts.client()
	if err != nil {
		return nil, err
	}

	if err := download(ctx, opts.catalogURL(), catfile, client); err != nil {
		return nil, err
	}

	cachedir, err := opts.cacheDir()
	if err != nil {
		return nil, err
	}

	conf := k6build.LocalBuildServiceConfig{
		BuildEnv:  map[string]string{"GOWORK": "off"},
		Catalog:   catfile,
		CopyGoEnv: true,
		CacheDir:  filepath.Join(cachedir, "build"),
		Verbose:   opts.verbose(),
	}

	return k6build.NewLocalBuildService(ctx, conf)
}

func newBuildServiceClient(opts *Options) (k6build.BuildService, error) {
	return k6build.NewBuildServiceClient(k6build.BuildServiceClientConfig{
		URL: opts.BuildServiceURL.String(),
	})
}

//nolint:forbidigo
func download(ctx context.Context, from *url.URL, dest string, client *http.Client) error {
	tmp, err := os.CreateTemp(filepath.Dir(dest), filepath.Base(dest)+"-*")
	if err != nil {
		return err
	}

	if from.Scheme == "file" {
		err = fileDownload(from, tmp)
	} else {
		err = httpDownload(ctx, from, tmp, client)
	}

	if err != nil {
		_ = os.Remove(tmp.Name())

		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	err = os.Chmod(tmp.Name(), syscall.S_IRUSR|syscall.S_IXUSR)
	if err != nil {
		return err
	}

	return os.Rename(tmp.Name(), dest)
}

//nolint:forbidigo
func fileDownload(from *url.URL, dest *os.File) error {
	src, err := os.Open(from.Path)
	if err != nil {
		return err
	}

	defer src.Close() //nolint:errcheck

	_, err = io.Copy(dest, src)

	return err
}

//nolint:forbidigo
func httpDownload(ctx context.Context, from *url.URL, dest *os.File, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, from.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", os.ErrNotExist, resp.Status)
	}

	defer resp.Body.Close() //nolint:errcheck

	_, err = io.Copy(dest, resp.Body)

	return err
}

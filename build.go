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

const platform = runtime.GOOS + "/" + runtime.GOARCH

func depsConvert(deps []*k6deps.Dependency) (string, []k6build.Dependency) {
	bdeps := make([]k6build.Dependency, len(deps)-1)

	for idx, dep := range deps[1:] {
		bdeps[idx] = k6build.Dependency{Name: dep.Name, Constraints: dep.GetConstraints().String()}
	}

	return deps[0].GetConstraints().String(), bdeps
}

func build(ctx context.Context, deps []*k6deps.Dependency, _ *Options) (*url.URL, error) {
	svc, err := k6build.DefaultLocalBuildService()
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

//nolint:forbidigo
func download(ctx context.Context, from *url.URL, dest string, client *http.Client) error {
	tmp, err := os.CreateTemp(filepath.Dir(dest), k6temp)
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
	if err != nil {
		return err
	}

	err = os.Chmod(dest.Name(), syscall.S_IRUSR|syscall.S_IXUSR)
	if err != nil {
		return err
	}

	return nil
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
	if err != nil {
		return err
	}

	err = os.Chmod(dest.Name(), syscall.S_IRUSR|syscall.S_IXUSR)
	if err != nil {
		return err
	}

	return nil
}

package k6exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/grafana/k6deps"
)

func build(_ []*k6deps.Dependency) (*url.URL, error) {
	loc, err := exec.LookPath(k6binary)
	if err != nil {
		return nil, err
	}

	return url.Parse("file://" + filepath.ToSlash(loc))
}

func download(ctx context.Context, from *url.URL, dir string, client *http.Client) error {
	dest, err := os.CreateTemp(dir, k6temp) //nolint:forbidigo
	if err != nil {
		return err
	}

	if from.Scheme == "file" {
		err = fileDownload(from, dest)
	} else {
		err = httpDownload(ctx, from, dest, client)
	}

	if err != nil {
		_ = os.Remove(dest.Name()) //nolint:forbidigo

		return err
	}

	if err = dest.Close(); err != nil {
		return err
	}

	return os.Rename(dest.Name(), filepath.Join(dir, k6binary)) //nolint:forbidigo
}

func fileDownload(from *url.URL, dest *os.File) error { //nolint:forbidigo
	src, err := os.Open(from.Path) //nolint:forbidigo
	if err != nil {
		return err
	}

	defer src.Close() //nolint:errcheck

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return os.Chmod(dest.Name(), syscall.S_IRUSR|syscall.S_IXUSR) //nolint:forbidigo
}

func httpDownload(ctx context.Context, from *url.URL, dest *os.File, client *http.Client) error { //nolint:forbidigo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, from.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", errDownload, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", errDownload, resp.Status)
	}

	defer resp.Body.Close() //nolint:errcheck

	_, err = io.Copy(dest, resp.Body)

	return err
}

var errDownload = errors.New("download error")

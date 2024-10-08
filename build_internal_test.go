package k6exec

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/grafana/k6build/pkg/client"

	"github.com/grafana/k6deps"
	"github.com/stretchr/testify/require"
)

func Test_depsConvert(t *testing.T) {
	t.Parallel()

	src := make(k6deps.Dependencies)

	err := src.UnmarshalText([]byte("k6>0.50;k6/x/faker>0.2.0"))

	require.NoError(t, err)

	k6Constraints, deps := depsConvert(src)

	require.Equal(t, ">0.50", k6Constraints)

	require.Equal(t, "k6/x/faker", deps[0].Name)
	require.Equal(t, src["k6/x/faker"].Constraints.String(), deps[0].Constraints)

	err = src.UnmarshalText([]byte("k6/x/faker*"))

	require.NoError(t, err)

	k6Constraints, deps = depsConvert(src)

	require.Equal(t, "*", k6Constraints)

	require.Equal(t, "k6/x/faker", deps[0].Name)
	require.Equal(t, "*", deps[0].Constraints)
}

func Test_newBuildService(t *testing.T) {
	t.Parallel()

	opts := &Options{
		StateDir: t.TempDir(),
		CacheDir: t.TempDir(),
	}

	opts.BuildServiceURL, _ = url.Parse("http://localhost:8000")

	svc, err := newBuildService(context.Background(), opts)

	require.NoError(t, err)
	require.IsType(t, new(client.BuildClient), svc)
	require.Equal(t, "*client.BuildClient", fmt.Sprintf("%T", svc))

	srv := testWebServer(t)
	defer srv.Close()

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/empty-catalog.json")
	require.NoError(t, err)
	opts.BuildServiceURL = nil

	svc, err = newBuildService(context.Background(), opts)

	require.NoError(t, err)

	require.NotEqual(t, "*k6build.BuildClient", fmt.Sprintf("%T", svc))
}

func Test_newLocalBuildService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	abs, err := filepath.Abs(filepath.Join("testdata", "minimal-catalog.json"))
	require.NoError(t, err)

	u, err := url.Parse("file://" + filepath.ToSlash(abs))
	require.NoError(t, err)

	opts := &Options{ExtensionCatalogURL: u}

	_, err = newLocalBuildService(ctx, opts)
	require.NoError(t, err)

	opts.StateDir = ""
	opts.AppName = invalidAppName(t)
	_, err = newLocalBuildService(ctx, opts)
	require.Error(t, err)
}

//nolint:forbidigo
func Test_httpDownload(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	ctx := context.Background()
	from, err := url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	dest, err := os.Create(filepath.Clean(filepath.Join(tmp, "catalog.json")))

	require.NoError(t, err)

	err = httpDownload(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	require.NoError(t, dest.Close())

	contents, err := os.ReadFile(dest.Name())

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))
}

//nolint:forbidigo
func Test_fileDownload(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	abs, err := filepath.Abs(filepath.Join("testdata", "empty-catalog.json"))

	require.NoError(t, err)

	from, err := url.Parse("file://" + filepath.ToSlash(abs))

	require.NoError(t, err)

	dest, err := os.Create(filepath.Clean(filepath.Join(tmp, "catalog.json")))

	require.NoError(t, err)

	err = fileDownload(from, dest)
	require.NoError(t, err)

	require.NoError(t, dest.Close())

	contents, err := os.ReadFile(dest.Name())

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))

	from, err = url.Parse("file:///" + tmp + "/no_such_file")

	require.NoError(t, err)

	err = fileDownload(from, dest)

	require.Error(t, err)
}

//nolint:forbidigo
func Test_download(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	ctx := context.Background()
	from, err := url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	dest := filepath.Clean(filepath.Join(tmp, "catalog.json"))

	require.NoError(t, err)

	err = download(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	contents, err := os.ReadFile(dest)

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))

	abs, err := filepath.Abs(filepath.Join("testdata", "empty-catalog.json"))

	require.NoError(t, err)

	from, err = url.Parse("file://" + filepath.ToSlash(abs))

	require.NoError(t, err)

	err = download(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	contents, err = os.ReadFile(dest)

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))
}

func testWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.FileServer(http.Dir("testdata")))
}

func Test_build(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" { // TODO - Re-enable as soon as k6build supports Windows!
		t.Skip("Skip because k6build doesn't work on Windows yet!")
	}

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/minimal-catalog.json")

	require.NoError(t, err)

	ctx := context.Background()

	opts := &Options{CacheDir: t.TempDir(), StateDir: t.TempDir(), ExtensionCatalogURL: u}

	loc, err := build(ctx, make(k6deps.Dependencies), opts)

	require.NoError(t, err)

	tmp := t.TempDir()

	dest := filepath.Join(tmp, k6binary)

	err = download(ctx, loc, dest, nil)

	require.NoError(t, err)

	cmd := exec.Command(filepath.Clean(dest), "version") //nolint:gosec

	out, err := cmd.Output()

	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(out), "k6 "))

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	_, err = build(ctx, make(k6deps.Dependencies), opts)

	require.Error(t, err)

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/missing-catalog.json")

	require.NoError(t, err)

	_, err = build(ctx, make(k6deps.Dependencies), opts)

	require.Error(t, err)
}

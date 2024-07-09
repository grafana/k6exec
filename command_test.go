package k6exec_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/grafana/k6exec"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" { // TODO - Re-enable as soon as k6build supports Windows!
		t.Skip("Skip because k6build doesn't work on Windows yet!")
	}

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/minimal-catalog.json")
	require.NoError(t, err)

	ctx := context.Background()

	opts := &k6exec.Options{StateDir: t.TempDir(), CacheDir: t.TempDir(), ExtensionCatalogURL: u}

	cmd, err := k6exec.Command(ctx, []string{"version"}, nil, opts)

	require.NoError(t, err)

	out, err := cmd.Output()

	require.NoError(t, err)

	require.True(t, strings.HasPrefix(string(out), "k6 "))
}

func TestCommand_errors(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/missing-catalog.json")
	require.NoError(t, err)

	ctx := context.Background()

	_, err = k6exec.Command(ctx, nil, nil, &k6exec.Options{AppName: invalidAppName(t)})
	require.Error(t, err)
	require.ErrorIs(t, err, k6exec.ErrState)

	_, err = k6exec.Command(ctx, nil, nil, &k6exec.Options{ExtensionCatalogURL: u})
	require.Error(t, err)
	require.ErrorIs(t, err, k6exec.ErrBuild)
}

func testWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.FileServer(http.Dir("testdata")))
}

func invalidAppName(t *testing.T) string {
	t.Helper()

	return strings.Repeat("too long", 2048)
}

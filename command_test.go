package k6exec_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/grafana/k6exec"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/minimal-catalog.json")
	require.NoError(t, err)

	ctx := context.Background()

	cmd, err := k6exec.Command(ctx, []string{"version"}, nil, &k6exec.Options{ExtensionCatalogURL: u})

	require.NoError(t, err)

	out, err := cmd.Output()

	require.NoError(t, err)

	require.True(t, strings.HasPrefix(string(out), "k6 "))
}

func testWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.FileServer(http.Dir("testdata")))
}

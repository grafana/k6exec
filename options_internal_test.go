package k6exec

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/stretchr/testify/require"
)

//nolint:forbidigo
func Test_Options_appname(t *testing.T) { //nolint:paralleltest
	const defaultAppName = "foo"
	saved := os.Args[0]

	os.Args[0] = defaultAppName
	defer func() { os.Args[0] = saved }()

	require.Equal(t, defaultAppName, (*Options)(nil).appname())
	require.Equal(t, defaultAppName, new(Options).appname())
	require.Equal(t, "bar", (&Options{AppName: "bar"}).appname())
}

func Test_Options_extensionCatalogURL(t *testing.T) {
	t.Parallel()

	defaultValue, err := url.Parse(DefaultExtensionCatalogURL)

	require.NoError(t, err)
	require.Equal(t, defaultValue, (*Options)(nil).extensionCatalogURL())
	require.Equal(t, defaultValue, new(Options).extensionCatalogURL())

	value, err := url.Parse("https://example.com/catalog.json")
	require.NoError(t, err)

	require.Equal(t, value, (&Options{ExtensionCatalogURL: value}).extensionCatalogURL())
}

func Test_Options_client(t *testing.T) {
	t.Parallel()

	client, err := (*Options)(nil).client()

	require.NoError(t, err)
	require.NotNil(t, client)
	require.IsType(t, new(httpcache.Transport), client.Transport)
	require.IsType(t, new(diskcache.Cache), client.Transport.(*httpcache.Transport).Cache)

	client, err = new(Options).client()

	require.NoError(t, err)
	require.NotNil(t, client)
	require.IsType(t, new(httpcache.Transport), client.Transport)
	require.IsType(t, new(diskcache.Cache), client.Transport.(*httpcache.Transport).Cache)

	client, err = (&Options{Client: http.DefaultClient}).client()

	require.NoError(t, err)
	require.Same(t, http.DefaultClient, client)

	_, err = (&Options{AppName: invalidAppName(t)}).client()

	require.Error(t, err)
	require.ErrorIs(t, err, ErrCache)
}

func Test_Options_cacheDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	xdg.Reload()

	expected := filepath.Join(tmp, (*Options)(nil).appname())

	dir, err := (*Options)(nil).cacheDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	dir, err = new(Options).cacheDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	t.Setenv("XDG_CACHE_HOME", "")

	xdg.Reload()

	tmp = t.TempDir()
	expected = filepath.Join(tmp, (*Options)(nil).appname())

	dir, err = (&Options{CacheDir: expected}).cacheDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	_, err = (&Options{AppName: invalidAppName(t)}).cacheDir()

	require.Error(t, err)
	require.ErrorIs(t, err, ErrCache)
}

func Test_Options_stateDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_STATE_HOME", tmp)

	xdg.Reload()

	expected := filepath.Join(tmp, (*Options)(nil).appname())

	dir, err := (*Options)(nil).stateDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	dir, err = new(Options).stateDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	t.Setenv("XDG_STATE_HOME", "")

	xdg.Reload()

	tmp = t.TempDir()
	expected = filepath.Join(tmp, (*Options)(nil).appname())

	dir, err = (&Options{StateDir: expected}).stateDir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	_, err = (&Options{AppName: invalidAppName(t)}).stateDir()

	require.Error(t, err)
	require.ErrorIs(t, err, ErrState)
}

func Test_Options_stateSubDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_STATE_HOME", tmp)

	xdg.Reload()

	expected := filepath.Join(tmp, (*Options)(nil).appname(), strconv.Itoa(os.Getpid())) //nolint:forbidigo

	dir, err := (*Options)(nil).stateSubdir()

	require.NoError(t, err)
	require.Equal(t, expected, dir)

	_, err = (&Options{AppName: invalidAppName(t)}).stateSubdir()

	require.Error(t, err)
	require.ErrorIs(t, err, ErrState)
}

func invalidAppName(t *testing.T) string {
	t.Helper()

	return strings.Repeat("too long", 2048)
}

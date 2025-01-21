package cmd

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/k6exec"
	"github.com/stretchr/testify/require"
)

func Test_newState(t *testing.T) {
	lvar := new(slog.LevelVar)

	t.Setenv("K6_BUILD_SERVICE_URL", "")
	t.Setenv("K6_EXTENSION_CATALOG_URL", "")

	st := newState(lvar, nil)

	require.Same(t, lvar, st.levelVar)
	require.Empty(t, st.buildServiceURL)
	require.Empty(t, st.extensionCatalogURL)

	t.Setenv("K6_BUILD_SERVICE_URL", "foo")
	t.Setenv("K6_EXTENSION_CATALOG_URL", "bar")

	st = newState(lvar, nil)

	require.Equal(t, "foo", st.buildServiceURL)
	require.Equal(t, "bar", st.extensionCatalogURL)
}

func Test_persistentPreRunE(t *testing.T) {
	t.Parallel()

	st := &state{levelVar: new(slog.LevelVar)}

	require.NoError(t, st.persistentPreRunE(nil, nil))
	require.Nil(t, st.BuildServiceURL)
	require.Nil(t, st.ExtensionCatalogURL)
	require.Equal(t, slog.LevelInfo, st.levelVar.Level())

	st.buildServiceURL = "http://example.com"
	st.extensionCatalogURL = "http://example.net"

	require.NoError(t, st.persistentPreRunE(nil, nil))
	require.Equal(t, "http://example.com", st.BuildServiceURL.String())
	require.Equal(t, "http://example.net", st.ExtensionCatalogURL.String())

	st.buildServiceURL = "http://example.com/%"
	require.Error(t, st.persistentPreRunE(nil, nil))

	st.buildServiceURL = "http://example.com"
	st.extensionCatalogURL = "http://example.net/%"
	require.Error(t, st.persistentPreRunE(nil, nil))

	st.buildServiceURL = "http://example.com"
	st.extensionCatalogURL = "http://example.net"
	st.verbose = true

	require.NoError(t, st.persistentPreRunE(nil, nil))
	require.Equal(t, slog.LevelDebug, st.levelVar.Level())

	st.levelVar = nil
	require.NoError(t, st.persistentPreRunE(nil, nil))
}

func Test_preRunE(t *testing.T) {
	t.Parallel()

	st := &state{
		levelVar: new(slog.LevelVar),
		Options:  k6exec.Options{CacheDir: t.TempDir()},
	}

	sub := newSubcommand("version", st)

	require.NoError(t, st.preRunE(sub, nil))
	require.NotContains(t, st.cmd.Args, "-v")
	require.NotContains(t, st.cmd.Args, "-q")
	require.NotContains(t, st.cmd.Args, "--no-color")

	st.verbose = true
	st.nocolor = true
	require.NoError(t, st.preRunE(sub, nil))
	require.Contains(t, st.cmd.Args, "-v")
	require.NotContains(t, st.cmd.Args, "-q")
	require.Contains(t, st.cmd.Args, "--no-color")

	st.verbose = false
	st.nocolor = false
	st.quiet = true
	require.NoError(t, st.preRunE(sub, nil))
	require.NotContains(t, st.cmd.Args, "-v")
	require.Contains(t, st.cmd.Args, "-q")
	require.NotContains(t, st.cmd.Args, "--no-color")

	st.quiet = false

	sub = newSubcommand("run", st)

	arg := filepath.Join("testdata", "script.js")
	require.NoError(t, st.preRunE(sub, []string{arg}))

	arg = filepath.Join("testdata", "archive.tar")
	require.NoError(t, st.preRunE(sub, []string{arg}))

	arg = filepath.Join("testdata", "invalid_constraint.js")
	require.Error(t, st.preRunE(sub, []string{arg}))

	arg = filepath.Join("testdata", "no_such_version.js")
	require.Error(t, st.preRunE(sub, []string{arg}))
}

func Test_runE(t *testing.T) {
	t.Parallel()

	st := &state{
		levelVar: new(slog.LevelVar),
		Options:  k6exec.Options{CacheDir: t.TempDir()},
	}

	err := st.preRunE(newSubcommand("version", st), nil)

	require.NoError(t, err)

	require.True(t, exists(t, st.cmd.Path))

	err = st.runE(nil, nil)

	require.NoError(t, err)

	require.False(t, exists(t, st.cmd.Path))
}

func Test_helpFunc(t *testing.T) { //nolint:paralleltest
	st := &state{
		levelVar: new(slog.LevelVar),
		Options:  k6exec.Options{CacheDir: t.TempDir()},
	}

	out := captureStderr(t, func() { st.helpFunc(newSubcommand("version", st), nil) })

	require.Empty(t, out)
}

func exists(t *testing.T, filename string) bool {
	t.Helper()

	_, err := os.Stat(filename) //nolint:forbidigo

	return err == nil
}

//nolint:forbidigo
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stderr
	defer func() { os.Stderr = orig }()

	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	require.NoError(t, w.Close())

	out, err := io.ReadAll(r)

	require.NoError(t, err)

	return string(out)
}

package cmd_test

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/grafana/k6exec/cmd"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) { //nolint:paralleltest
	lvar := new(slog.LevelVar)

	c := cmd.New(lvar, nil)

	require.Contains(t, c.Use, "k6exec")
	require.Contains(t, c.Long, "k6exec")

	require.NotNil(t, c.PreRunE)
	require.NotNil(t, c.RunE)
	require.NotNil(t, c.PersistentPreRunE)

	flags := c.PersistentFlags()

	require.NotNil(t, flags.Lookup("extension-catalog-url"))
	require.NotNil(t, flags.Lookup("build-service-url"))
	require.NotNil(t, flags.Lookup("verbose"))
	require.NotNil(t, flags.Lookup("quiet"))
	require.NotNil(t, flags.Lookup("no-color"))

	out := captureStdout(t, func() { require.NoError(t, c.Execute()) })

	require.True(t, strings.Contains(out, "  k6"))
}

//nolint:forbidigo
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	defer func() { os.Stdout = orig }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	require.NoError(t, w.Close())

	out, err := io.ReadAll(r)

	require.NoError(t, err)

	return string(out)
}

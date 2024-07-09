package cmd_test

import (
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/grafana/k6exec/cmd"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) { //nolint:paralleltest
	if runtime.GOOS == "windows" { // TODO - Re-enable as soon as k6build supports Windows!
		t.Skip("Skip because k6build doesn't work on Windows yet!")
	}

	lvar := new(slog.LevelVar)

	c := cmd.New(lvar)

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
	require.NotNil(t, flags.Lookup("usage"))

	c.SetArgs([]string{"--usage"})

	out := captureStdout(t, func() { require.NoError(t, c.Execute()) })

	require.Equal(t, c.Long+"\n"+c.UsageString(), out)

	c = cmd.New(lvar)

	c.SetArgs([]string{})

	out = captureStdout(t, func() { require.NoError(t, c.Execute()) })

	require.True(t, strings.Contains(out, "  k6 [command]"))
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

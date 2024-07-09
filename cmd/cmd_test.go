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

//nolint:forbidigo
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

	orig := os.Stdout
	defer func() {
		os.Stdout = orig
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	c.SetArgs([]string{"--usage"})

	require.NoError(t, c.Execute())

	require.NoError(t, w.Close())

	out, _ := io.ReadAll(r)

	require.Equal(t, c.Long+"\n"+c.UsageString(), string(out))

	r, w, _ = os.Pipe()
	os.Stdout = w

	c = cmd.New(lvar)

	c.SetArgs([]string{})

	require.NoError(t, c.Execute())

	require.NoError(t, w.Close())

	out, _ = io.ReadAll(r)

	require.True(t, strings.Contains(string(out), "  k6 [command]"))
}

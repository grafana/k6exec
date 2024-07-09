package main

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newCmd(t *testing.T) { //nolint:paralleltest
	if runtime.GOOS == "windows" { // TODO - Re-enable as soon as k6build supports Windows!
		t.Skip("Skip because k6build doesn't work on windows yet!")
	}

	abs, err := filepath.Abs(filepath.Join("..", "..", "examples", "combined.js"))

	require.NoError(t, err)

	lvar := new(slog.LevelVar)

	cmd := newCmd([]string{"run", abs}, lvar)

	require.NoError(t, cmd.Execute())
}

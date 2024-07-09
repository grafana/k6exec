package main

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newCmd(t *testing.T) { //nolint:paralleltest
	abs, err := filepath.Abs(filepath.Join("..", "..", "examples", "combined.js"))

	require.NoError(t, err)

	lvar := new(slog.LevelVar)

	cmd := newCmd([]string{"run", abs}, lvar)

	require.NoError(t, cmd.Execute())
}

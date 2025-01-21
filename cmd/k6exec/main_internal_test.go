package main

import (
	"log/slog"
	"path/filepath"
	"testing"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_newCmd(t *testing.T) { //nolint:paralleltest
	abs, err := filepath.Abs(filepath.Join("..", "..", "examples", "combined.js"))

	require.NoError(t, err)

	lvar := new(slog.LevelVar)

	cmd := newCmd([]string{"run", abs}, lvar)

	require.NoError(t, cmd.Execute())
}

func Test_initLogging(t *testing.T) { //nolint:paralleltest
	lvar := initLogging(appname)

	require.NotNil(t, lvar)
	require.Equal(t, logrus.DebugLevel, logrus.GetLevel())
	require.IsType(t, new(sloglogrus.LogrusHandler), slog.Default().Handler())
}

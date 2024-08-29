package cmd

import (
	"path/filepath"
	"testing"

	"github.com/grafana/k6deps"
	"github.com/stretchr/testify/require"
)

func Test_analyzeArchive(t *testing.T) {
	t.Parallel()

	actual, err := analyzeArchive(filepath.Join("testdata", "archive.tar"))

	require.NoError(t, err)

	opts := &k6deps.Options{
		Script:   k6deps.Source{Name: filepath.Join("testdata", "combined.js")},
		Manifest: k6deps.Source{Ignore: true},
		Env:      k6deps.Source{Ignore: true},
	}

	expected, err := k6deps.Analyze(opts)

	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

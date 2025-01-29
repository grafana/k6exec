package k6exec_test

import (
	"context"
	"strings"
	"testing"

	"github.com/grafana/k6build/pkg/testutils"
	"github.com/grafana/k6deps"
	"github.com/grafana/k6exec"
	"github.com/grafana/k6provider"

	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	env, err := testutils.NewTestEnv(testutils.TestEnvConfig{
		WorkDir:    t.TempDir(),
		CatalogURL: "testdata/minimal-catalog.json",
	})
	require.NoError(t, err)

	t.Cleanup(env.Cleanup)

	opts := &k6exec.Options{
		Env:             k6deps.Source{Ignore: true},
		Manifest:        k6deps.Source{Ignore: true},
		BuildServiceURL: env.BuildServiceURL(),
	}

	cmd, cleanup, err := k6exec.Command(context.TODO(), []string{"version"}, opts)
	defer func() { require.NoError(t, cleanup()) }()

	require.NoError(t, err)

	out, err := cmd.Output()

	require.NoError(t, err)

	require.True(t, strings.HasPrefix(string(out), "k6"))
}

func TestCommand_errors(t *testing.T) {
	t.Parallel()

	env, err := testutils.NewTestEnv(testutils.TestEnvConfig{
		WorkDir:    t.TempDir(),
		CatalogURL: "testdata/empty-catalog.json",
	})
	require.NoError(t, err)

	t.Cleanup(env.Cleanup)

	opts := &k6exec.Options{
		Env:             k6deps.Source{Ignore: true},
		Manifest:        k6deps.Source{Ignore: true},
		BuildServiceURL: env.BuildServiceURL(),
	}

	_, _, err = k6exec.Command(context.TODO(), nil, opts)
	require.Error(t, err)
	require.ErrorIs(t, err, k6provider.ErrInvalidParameters)
}

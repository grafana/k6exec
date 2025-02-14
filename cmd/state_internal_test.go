package cmd

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/k6build/pkg/testutils"
	"github.com/grafana/k6exec"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_interal_state(t *testing.T) {
	t.Setenv("K6_BUILD_SERVICE_URL", "")
	t.Setenv("K6_CLOUD_TOKEN", "")

	env, err := testutils.NewTestEnv(testutils.TestEnvConfig{
		WorkDir: t.TempDir(),
	})
	require.NoError(t, err)

	t.Cleanup(env.Cleanup)

	t.Run("Test_newState", func(t *testing.T) { //nolint:paralleltest
		lvar := new(slog.LevelVar)

		st := newState(lvar)

		require.Same(t, lvar, st.levelVar)
	})

	t.Run("Test_persistentPreRunE", func(t *testing.T) { //nolint:paralleltest
		cmd := &cobra.Command{}
		st := &state{levelVar: new(slog.LevelVar)}

		require.NoError(t, st.persistentPreRunE(cmd, nil))
		require.Equal(t, defaultBuildServiceURL, st.BuildServiceURL)
		require.Equal(t, slog.LevelInfo, st.levelVar.Level())

		st.buildServiceURL = "http://example.com"

		require.NoError(t, st.persistentPreRunE(cmd, nil))
		require.Equal(t, "http://example.com", st.BuildServiceURL)

		st.buildServiceURL = "http://example.com"
		st.verbose = true

		require.NoError(t, st.persistentPreRunE(cmd, nil))
		require.Equal(t, slog.LevelDebug, st.levelVar.Level())

		st.levelVar = nil
		require.NoError(t, st.persistentPreRunE(cmd, nil))
	})

	t.Run("Test_loadConfig", func(t *testing.T) { //nolint:paralleltest
		st := &state{
			levelVar:   new(slog.LevelVar),
			configFile: filepath.Join("testdata", "config", "valid.json"),
		}

		require.NoError(t, st.persistentPreRunE(&cobra.Command{}, nil))
		require.Equal(t, "token", st.Options.BuildServiceToken)

		st = &state{
			levelVar:   new(slog.LevelVar),
			configFile: filepath.Join("testdata", "config", "empty.json"),
		}

		require.NoError(t, st.persistentPreRunE(&cobra.Command{}, nil))
		require.Empty(t, st.Options.BuildServiceToken)

		// test config override from flag
		cmd := &cobra.Command{Use: "test"}
		cmd.SetContext(context.WithValue(context.Background(), argsKey{}, []string{"test", "--config", "no_such_file.json"}))
		st = &state{
			levelVar: new(slog.LevelVar),
		}
		require.Error(t, st.persistentPreRunE(cmd, nil))
	})

	t.Run("Test_preRunE", func(t *testing.T) { //nolint:paralleltest
		st := &state{
			levelVar: new(slog.LevelVar),
			Options:  k6exec.Options{BuildServiceURL: env.BuildServiceURL()},
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
	})

	t.Run("Test_runE", func(t *testing.T) { //nolint:paralleltest
		st := &state{
			levelVar: new(slog.LevelVar),
			Options:  k6exec.Options{BuildServiceURL: env.BuildServiceURL()},
		}

		err = st.preRunE(newSubcommand("version", st), nil)

		require.NoError(t, err)

		require.True(t, exists(t, st.cmd.Path))

		err = st.runE(nil, nil)

		require.NoError(t, err)
	})

	t.Run("Test_helpFunc", func(t *testing.T) { //nolint:paralleltest
		st := &state{
			levelVar: new(slog.LevelVar),
			Options:  k6exec.Options{BuildServiceURL: env.BuildServiceURL()},
		}

		out := captureStderr(t, func() { st.helpFunc(newSubcommand("version", st), nil) })

		require.Empty(t, out)
	})
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

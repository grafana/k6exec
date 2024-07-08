package k6exec_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/grafana/k6exec"
	"github.com/stretchr/testify/require"
)

func exists(t *testing.T, filename string) bool {
	t.Helper()

	_, err := os.Stat(filename)

	return err == nil
}

//nolint:forbidigo
func TestCleanupState(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	opts := &k6exec.Options{StateDir: tmp}

	subdir := filepath.Join(tmp, strconv.Itoa(os.Getpid()))

	require.NoError(t, os.MkdirAll(subdir, 0700))

	name := filepath.Join(subdir, "hello.txt")

	err := os.WriteFile(name, []byte("Hello, World!\n"), 0o600)

	require.NoError(t, err)

	require.True(t, exists(t, name))

	err = k6exec.CleanupState(opts)
	require.NoError(t, err)

	require.False(t, exists(t, subdir))
	require.True(t, exists(t, tmp))

	opts = &k6exec.Options{AppName: strings.Repeat("too long", 2048)}
	err = k6exec.CleanupState(opts)
	require.Error(t, err)
}

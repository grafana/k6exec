package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_scriptArg(t *testing.T) {
	t.Parallel()

	cmd := new(cobra.Command)

	sarg, has := scriptArg(cmd, nil)

	require.False(t, has)
	require.Empty(t, sarg)

	cmd.Annotations = map[string]string{"foo": "bar"}

	sarg, has = scriptArg(cmd, nil)

	require.False(t, has)
	require.Empty(t, sarg)

	cmd.Annotations = map[string]string{useExtensions: "true"}

	sarg, has = scriptArg(cmd, nil)

	require.False(t, has)
	require.Empty(t, sarg)

	sarg, has = scriptArg(cmd, []string{"-"})

	require.False(t, has)
	require.Empty(t, sarg)

	sarg, has = scriptArg(cmd, []string{"script.js"})

	require.True(t, has)
	require.Equal(t, "script.js", sarg)
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

type argsKey struct{}

// SetArgs set arguments for cmd command and store the args in the context.
func SetArgs(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cmd.SetContext(context.WithValue(ctx, argsKey{}, args))

	cmd.SetArgs(args)
}

func getArgs(cmd *cobra.Command) []string {
	ctx := cmd.Context()
	if ctx == nil {
		return nil
	}

	value := ctx.Value(argsKey{})
	if value == nil {
		return nil
	}

	args, ok := value.([]string)
	if !ok {
		return nil
	}

	for idx := range args {
		if args[idx] == cmd.Name() {
			return args[idx+1:]
		}
	}

	return nil
}

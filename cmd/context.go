package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type argsKey struct{}

// SetArgs set arguments for cmd command and store the args in the command context.
func SetArgs(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cmd.SetContext(context.WithValue(ctx, argsKey{}, args))

	cmd.SetArgs(args)
}

// getArgs returns the arguments from the command context, which were previously stored by SetArgs.
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

// getFlagValue returns the value of the flag from the command arguments given its name and (optional) shortName.
// If the flag is not found it returns an empty string.
// if it is found but the next element is not its value, it returns an error.
func getFlagValue(cmd *cobra.Command, fullName string, shortName string) (string, error) {
	args := getArgs(cmd)

	for i, arg := range args {
		if arg == fullName || (shortName != "" && arg == shortName) {
			// if this is the last argument or the next element is a flag, return an empty string
			// this is an error in thc CLI arguments
			if i+1 == len(args) || strings.HasPrefix(args[i+1], "-") {
				return "", fmt.Errorf("flag %s is missing a value", arg)
			}

			return args[i+1], nil
		}
	}

	return "", nil
}

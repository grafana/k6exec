package k6exec

import (
	"fmt"
	"os"
)

// CleanupState deletes the state directory belonging to the current process.
func CleanupState(opts *Options) error {
	dir, err := opts.stateSubdir()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrState, err.Error())
	}

	if err = os.RemoveAll(dir); err != nil { //nolint:forbidigo
		return fmt.Errorf("%w: %s", ErrState, err.Error())
	}

	return nil
}

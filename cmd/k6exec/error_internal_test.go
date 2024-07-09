package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testError struct{}

func (*testError) Error() string {
	return "test error"
}

func (*testError) Format(_ int, _ bool) string {
	return "formatted test error"
}

func Test_formatError(t *testing.T) {
	t.Parallel()

	require.Equal(t, errors.ErrUnsupported.Error(), formatError(errors.ErrUnsupported))
	require.Equal(t, "formatted test error", formatError(new(testError)))
}

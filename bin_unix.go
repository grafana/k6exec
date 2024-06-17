//go:build !windows

package k6exec

const (
	k6binary = "k6"
	k6temp   = "k6-*"
)

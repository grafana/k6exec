package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/grafana/k6deps"
	"github.com/grafana/k6exec"
	"github.com/grafana/k6provision"
)

const (
	dirMode     = syscall.S_IXUSR | syscall.S_IRUSR | syscall.S_IWUSR
	exeMode     = syscall.S_IXUSR | syscall.S_IRUSR | syscall.S_IWUSR
	dataMode    = syscall.S_IRUSR | syscall.S_IWUSR
	versionFile = "version.txt"
	hashFile    = "sha256.txt"
)

var binBaseDir = filepath.Join("k6", runtime.GOOS, runtime.GOARCH) //nolint:gochecknoglobals

func tryEmbedded(ctx context.Context, deps k6deps.Dependencies, exe string, next k6exec.ProvisionerFunc) error {
	if len(deps) > 1 { // at least one non-k6 dependency exists
		return next(ctx, deps, exe, next)
	}

	ek6ver := embeddedVersion()
	k6dep := &k6deps.Dependency{Name: "k6"}

	if len(deps) == 1 {
		dep, found := deps[k6dep.Name]
		if !found { // non-k6 dependency exists
			return next(ctx, deps, exe, next)
		}

		k6dep = dep
	}

	if !k6dep.GetConstraints().Check(ek6ver) {
		return next(ctx, deps, exe, next)
	}

	slog.Debug("Using the embedded k6", "version", ek6ver.String(), "constraints", k6dep.GetConstraints())

	if alreadyExtracted(exe) {
		return nil
	}

	return extractExe(exe)
}

func alreadyExtracted(exe string) bool {
	_, err := exec.LookPath(exe)

	if err != nil && !errors.Is(err, exec.ErrDot) {
		return false
	}

	file, err := os.Open(filepath.Clean(exe)) //nolint:forbidigo
	if err != nil {
		return false
	}

	defer file.Close() //nolint:errcheck

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	ehash, err := embeddedHash()
	if err != nil {
		return false
	}

	return bytes.Equal(ehash, hash.Sum(nil))
}

func embeddedHash() ([]byte, error) {
	strhash, err := binFS.ReadFile(filepath.Join(binBaseDir, hashFile))
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(string(strhash[0:64]))
}

func embeddedVersion() *semver.Version {
	binver, err := binFS.ReadFile(filepath.Join(binBaseDir, versionFile))
	if err == nil {
		ver, err := semver.NewVersion(string(binver))
		if err == nil {
			return ver
		}
	}

	return semver.New(0, 0, 0, "", "")
}

//nolint:forbidigo
func extractExe(exe string) error {
	slog.Debug("extract exe")

	src, err := binFS.Open(filepath.Join(binBaseDir, k6provision.ExeName+".gz"))
	if err != nil {
		return err
	}

	defer src.Close() //nolint:errcheck

	reader, err := gzip.NewReader(src)
	if err != nil {
		return err
	}

	defer reader.Close() //nolint:errcheck

	dst, err := os.Create(filepath.Clean(exe))
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, reader) //nolint:gosec
	if err != nil {
		return err
	}

	if err := dst.Close(); err != nil {
		return err
	}

	return os.Chmod(exe, exeMode)
}

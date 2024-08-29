package cmd

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6pack"
)

//nolint:forbidigo
func analyzeArchive(filename string) (k6deps.Dependencies, error) {
	dir, err := os.MkdirTemp("", "k6-archive-*")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dir) //nolint:errcheck

	err = extractArchive(dir, filename)
	if err != nil {
		return nil, err
	}

	opts, err := loadMetadata(dir)
	if err != nil {
		return nil, err
	}

	return k6deps.Analyze(opts)
}

//nolint:forbidigo
func loadMetadata(dir string) (*k6deps.Options, error) {
	var meta archiveMetadata

	data, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "metadata.json"))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	opts := new(k6deps.Options)

	opts.Manifest.Ignore = true // no manifest (yet) in archive

	opts.Script.Name = filepath.Join(
		dir,
		"file",
		filepath.FromSlash(strings.TrimPrefix(meta.Filename, "file:///")),
	)

	if value, found := meta.Env[k6deps.EnvDependencies]; found {
		opts.Env.Name = k6deps.EnvDependencies
		opts.Env.Contents = []byte(value)
	} else {
		opts.Env.Ignore = true
	}

	contents, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "data"))
	if err != nil {
		return nil, err
	}

	script, _, err := k6pack.Pack(string(contents), &k6pack.Options{Filename: opts.Script.Name})
	if err != nil {
		return nil, err
	}

	opts.Script.Contents = script

	return opts, nil
}

type archiveMetadata struct {
	Filename string            `json:"filename"`
	Env      map[string]string `json:"env"`
}

//nolint:forbidigo
func extractArchive(dir string, filename string) error {
	input, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer input.Close() //nolint:errcheck

	reader := tar.NewReader(input)

	const maxFileSize = 1024 * 1024 * 10 // 10M

	for {
		header, err := reader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dir, filepath.Clean(filepath.FromSlash(header.Name)))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o750); err != nil {
				return err
			}

		case tar.TypeReg:
			if ext := filepath.Ext(target); ext == ".csv" || (ext == ".json" && filepath.Base(target) != "metadata.json") {
				continue
			}

			file, err := os.OpenFile(filepath.Clean(target), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.CopyN(file, reader, maxFileSize); err != nil && !errors.Is(err, io.EOF) {
				return err
			}

			if err = file.Close(); err != nil {
				return err
			}
		}
	}
}

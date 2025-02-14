package k6exec

import (
	"context"
	"log/slog"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6provider"
)

func provision(ctx context.Context, deps k6deps.Dependencies, opts *Options) (string, error) {
	config := k6provider.Config{}

	if opts != nil {
		config.BuildServiceURL = opts.BuildServiceURL
		config.BuildServiceAuth = opts.BuildServiceToken
	}

	provider, err := k6provider.NewProvider(config)
	if err != nil {
		return "", err
	}

	slog.Debug("fetching binary", "build service URL: ", opts.BuildServiceURL)

	binary, err := provider.GetBinary(ctx, deps)
	if err != nil {
		return "", err
	}

	slog.Debug("binary fetched",
		"Path: ", binary.Path,
		"dependencies", deps.String(),
		"checksum", binary.Checksum,
		"chached", binary.Cached,
		"download URL", binary.DownloadURL,
	)

	return binary.Path, nil
}

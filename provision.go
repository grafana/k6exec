package k6exec

import (
	"context"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6provision"
)

func provision(ctx context.Context, deps k6deps.Dependencies, exe string, opts *Options) error {
	popts := new(k6provision.Options)

	if opts != nil {
		popts.AppName = opts.AppName
		popts.CacheDir = opts.CacheDir
		popts.Client = opts.Client
		popts.BuildServiceURL = opts.BuildServiceURL
		popts.ExtensionCatalogURL = opts.ExtensionCatalogURL
	}

	return k6provision.Provision(ctx, deps, exe, popts)
}

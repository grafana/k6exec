package k6exec

import (
	"net/http"
	"os"
	"path/filepath"
)

// Options contains the optional parameters of the Command function.
type Options struct {
	// AppName contains the name of the application. It is used to define the default value of CacheDir.
	// If empty, it defaults to os.Args[0].
	AppName string
	// CacheDir specifies the name of the directory where the k6 binary can be cached.
	// Its default is determined based on the XDG Base Directory Specification.
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	CacheDir string
	// NoCache can be used to disable the k6 binary cache.
	NoCache bool
	// Client is used during HTTP communication with the build service.
	// If absent, http.DefaultClient will be used.
	Client *http.Client
}

func (o *Options) appname() string {
	if o != nil && len(o.AppName) > 0 {
		return o.AppName
	}

	return filepath.Base(os.Args[0]) //nolint:forbidigo
}

func (o *Options) client() *http.Client {
	if o != nil && o.Client != nil {
		return o.Client
	}

	return http.DefaultClient
}

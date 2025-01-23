package k6exec

import (
	"net/http"
	"net/url"

	"github.com/grafana/k6deps"
)

// Options contains the optional parameters of the Command function.
type Options struct {
	// Manifest contains the properties of the manifest file to be analyzed.
	// If the Ignore property is not set and no manifest file is specified,
	// the package.json file closest to the script is searched for.
	Manifest k6deps.Source
	// Env contains the properties of the environment variable to be analyzed.
	// If the Ignore property is not set and no variable is specified,
	// the value of the variable named K6_DEPENDENCIES is read.
	Env k6deps.Source
	// LookupEnv function is used to query the value of the environment variable
	// specified in the Env option Name if the Contents of the Env option is empty.
	// If empty, os.LookupEnv will be used.
	LookupEnv func(key string) (value string, ok bool)
	// FindManifest function is used to find manifest file for the given scriptfile
	// if the Contents of Manifest option is empty.
	// If the scriptfile parameter is empty, FindManifest starts searching
	// for the manifest file from the current directory
	// If missing, the closest manifest file will be used.
	FindManifest func(scriptfile string) (contents []byte, filename string, ok bool, err error)
	// AppName contains the name of the application. It is used to define the default value of CacheDir.
	// If empty, it defaults to os.Args[0].
	AppName string
	// CacheDir specifies the name of the directory where the cacheable files can be cached.
	// Its default is determined based on the XDG Base Directory Specification.
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	CacheDir string
	// Client is used during HTTP communication with the build service.
	// If absent, http.DefaultClient will be used.
	Client *http.Client
	// ExtensionCatalogURL contains the URL of the k6 extension catalog to be used.
	// If absent, DefaultExtensionCatalogURL will be used.
	ExtensionCatalogURL *url.URL
	// BuildServiceURL contains the URL of the k6 build service to be used.
	// If the value is not nil, the k6 binary is built using the build service instead of the local build.
	BuildServiceURL *url.URL
}

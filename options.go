package k6exec

import (
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
	FindManifest func(scriptfile string) (filename string, ok bool, err error)
	// AppName contains the name of the application. It is used to define the default value of CacheDir.
	// If empty, it defaults to os.Args[0].
	AppName string
	// BuildServiceURL contains the URL of the k6 build service to be used.
	// If the value is not nil, the k6 binary is built using the build service instead of the local build.
	BuildServiceURL string
	// BuildServiceToken contains the token to be used to authenticate with the build service.
	// Defaults to K6_CLOUD_TOKEN environment variable is set, or the value stored in the k6 config file.
	BuildServiceToken string
}

package k6exec

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/adrg/xdg"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// Options contains the optional parameters of the Command function.
type Options struct {
	// AppName contains the name of the application. It is used to define the default value of CacheDir.
	// If empty, it defaults to os.Args[0].
	AppName string
	// CacheDir specifies the name of the directory where the cacheable files can be cached.
	// Its default is determined based on the XDG Base Directory Specification.
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	CacheDir string
	// StateDir specifies the name of the directory where the k6 running state is stored,
	// including the k6 binary and extension catalog. Each execution has a sub-directory,
	// which is deleted after successful execution.
	// Its default is determined based on the XDG Base Directory Specification.
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	StateDir string
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

// DefaultExtensionCatalogURL contains the address of the default k6 extension catalog.
const DefaultExtensionCatalogURL = "https://grafana.github.io/k6-extension-registry/catalog-registered.json"

func (o *Options) appname() string {
	if o != nil && len(o.AppName) > 0 {
		return o.AppName
	}

	return filepath.Base(os.Args[0]) //nolint:forbidigo
}

func (o *Options) client() (*http.Client, error) {
	if o != nil && o.Client != nil {
		return o.Client, nil
	}

	cachedir, err := o.cacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cachedir, "http")

	err = os.MkdirAll(dir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCache, err.Error())
	}

	transport := httpcache.NewTransport(diskcache.New(dir))

	return &http.Client{Transport: transport}, nil
}

func (o *Options) extensionCatalogURL() *url.URL {
	if o != nil && o.ExtensionCatalogURL != nil {
		return o.ExtensionCatalogURL
	}

	loc, _ := url.Parse(DefaultExtensionCatalogURL)

	return loc
}

func (o *Options) xdgDir(option string, xdgfunc func(string) (string, error), e error) (string, error) {
	var xdgdir string

	if o != nil && len(option) != 0 {
		xdgdir = option
	} else {
		dir, err := xdgfunc(o.appname())
		if err != nil {
			return "", fmt.Errorf("%w: %s", e, err.Error())
		}

		xdgdir = dir
	}

	err := os.MkdirAll(xdgdir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return "", fmt.Errorf("%w: %s", e, err.Error())
	}

	return xdgdir, nil
}

func (o *Options) cacheDir() (string, error) {
	var option string
	if o != nil {
		option = o.CacheDir
	}

	return o.xdgDir(option, xdg.CacheFile, ErrCache)
}

func (o *Options) stateDir() (string, error) {
	var option string
	if o != nil {
		option = o.StateDir
	}

	return o.xdgDir(option, xdg.StateFile, ErrState)
}

func (o *Options) stateSubdir() (string, error) {
	dir, err := o.stateDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, strconv.Itoa(os.Getpid())) //nolint:forbidigo

	err = os.MkdirAll(dir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return "", err
	}

	return dir, nil
}

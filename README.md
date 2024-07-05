[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6exec)](https://github.com/grafana/k6exec/releases/)
[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6exec.svg)](https://pkg.go.dev/github.com/grafana/k6exec)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6exec)](https://goreportcard.com/report/github.com/grafana/k6exec)
[![GitHub Actions](https://github.com/grafana/k6exec/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6exec/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6exec/graph/badge.svg?token=6MP3G02V9C)](https://codecov.io/gh/grafana/k6exec)
![GitHub Downloads](https://img.shields.io/github/downloads/grafana/k6exec/total)

<h1 name="title">k6exec</h1>

**Launcher for k6 with seamless use of extensions.**

k6exec is a launcher library for k6 with seamless use of extensions. The launcher will always run the k6 test script with the appropriate k6 binary, which contains the extensions used by the script. Extensions can also be recognized from the environment variable (default `K6_DEPENDENCIES`) or from the `dependencies` property of the manifest file.

k6exec is primarily used as a go library. In addition, it also contains a command-line tool, which is suitable for listing the dependencies of k6 test scripts.

The command line tool can be integrated into other command line tools as a subcommand. For this purpose, the library also contains the functionality of the command line tool as a factrory function that returns [cobra.Command](https://pkg.go.dev/github.com/spf13/cobra#Command).

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6exec/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6exec/cmd/k6exec@latest
```

Docker images are available on the GitHub [Packages](https://github.com/grafana/k6exec/pkgs/container/k6exec):

```
docker pull ghcr.io/grafana/k6exec:latest
```

## Usage

<!-- #region cli -->
## k6exec

Run k6 with extensions

### Synopsis

Run k6 with a seamless extension user experience.

`k6exec` is a [k6] launcher that automatically provides [k6] with the [extensions] used by the test. In order to do this, it analyzes the script arguments of the `run` and `archive` subcommands, detects the extensions to be used and their version constraints.

The launcher acts as a drop-in replacement for the `k6` command. For more convenient use, it is advisable to create an alias or shell script called `k6` for the launcher. The alias can be used in exactly the same way as the `k6` command, with the difference that it generates the real `k6` on the fly based on the extensions you want to use.

Any k6 command can be used. Use the `help` command to list the available k6 commands.

Since k6exec tries to emulate the `k6` command line, the `help` command or the `--help` flag cannot be used to display help from `k6exec` command itself. The `k6exec` help can be displayed using the `--usage` flag:

    k6exec --usage

### Prerequisites

k6exec tries to provide the appropriate k6 executable after detecting the extension dependencies. This can be done using a build service or a native builder.

#### Build Service

No additional installation is required to use the build service, just provide the build service URL.

The build service URL can be specified in the `K6_BUILD_SERVICE_URL` environment variable or by using the `--build-service-url` flag.

There is no default URL for the build service, otherwise k6exec will automatically provide k6 with the native builder.

#### Native Builder

To use the native builder, you only need to install the [Go language toolkit](https://go.dev/doc/install).

The native builder uses a k6 extension catalog to resolve extension URLs and versions. The extension catalog URL has a default value. A different extension catalog URL can be specified in the `K6_EXTENSION_CATALOG_URL` environment variable or by using the `--extension-catalog-url` flag.

### Pragma

Version constraints can be specified using the JavaScript `"use ..."` pragma syntax for k6 and extensions. Put the following lines at the beginning of the test script:

```js
"use k6 >= v0.52";
"use k6 with k6/x/faker > 0.2";
```

Any number of `"use k6"` pragmas can be used.

> **Note**
> The use of pragmas is completely optional for JavaScript type extensions, it is only necessary if you want to specify version constraints.

The pragma syntax can also be used to specify an extension dependency that is not referenced in an import expression. A typical example of this is the Output type extension such as [xk6-top]:

```js
"use k6 with top >= 0.1";
```

Read the version constraints syntax in the [Version Constraints](#version-constraints) section

### Environment

The extensions to be used and optionally their version constraints can also be specified in the `K6_DEPENDENCIES` environment variable. The value of the environment variable K6_DEPENDENCIES is a list of elements separated by semicolons. Each element specifies an extension (or k6 itself) and optionally its version constraint.

```
k6>=0.52;k6/x/faker>=0.3;k6/x/sql>=0.4
```

### Manifest

The manifest file is a JSON file, the `dependencies` property of which can specify extension dependencies and version constraints. The value of the `dependencies` property is a JSON object. The property names of this object are the extension names (or k6) and the values ​​are the version constraints.

```json
{
  "dependencies": {
    "k6": ">=0.52",
    "k6/x/faker": ">=0.3",
    "k6/x/sql": ">=0.4"
  }
}
```

The manifest file is a file named `package.json`, which is located closest to the k6 test script or the current directory, depending on whether the given subcommand has a test script argument (e.g. run, archive) or not (e.g. version). The `package.json` file is searched for up to the root of the directory hierarchy.

### Limitations

Version constraints can be specified in several sources ([pragma](#pragma), [environment](#environment), [manifest](#manifest)) but cannot be overwritten. That is, for a given extension, the version constraints from different sources must either be equal, or only one source can contain a version constraint.

### Version Constraints

*This section is based on the [Masterminds/semver] documentation.*

#### Basic Comparisons

There are two elements to the comparisons. First, a comparison string is a list
of space or comma separated AND comparisons. These are then separated by || (OR)
comparisons. For example, `">= 1.2 < 3.0.0 || >= 4.2.3"` is looking for a
comparison that's greater than or equal to 1.2 and less than 3.0.0 or is
greater than or equal to 4.2.3.

The basic comparisons are:

* `=`: equal (aliased to no operator)
* `!=`: not equal
* `>`: greater than
* `<`: less than
* `>=`: greater than or equal to
* `<=`: less than or equal to

#### Hyphen Range Comparisons

There are multiple methods to handle ranges and the first is hyphens ranges.
These look like:

* `1.2 - 1.4.5` which is equivalent to `>= 1.2 <= 1.4.5`
* `2.3.4 - 4.5` which is equivalent to `>= 2.3.4 <= 4.5`

#### Wildcards In Comparisons

The `x`, `X`, and `*` characters can be used as a wildcard character. This works
for all comparison operators. When used on the `=` operator it falls
back to the patch level comparison (see tilde below). For example,

* `1.2.x` is equivalent to `>= 1.2.0, < 1.3.0`
* `>= 1.2.x` is equivalent to `>= 1.2.0`
* `<= 2.x` is equivalent to `< 3`
* `*` is equivalent to `>= 0.0.0`

#### Tilde Range Comparisons (Patch)

The tilde (`~`) comparison operator is for patch level ranges when a minor
version is specified and major level changes when the minor number is missing.
For example,

* `~1.2.3` is equivalent to `>= 1.2.3, < 1.3.0`
* `~1` is equivalent to `>= 1, < 2`
* `~2.3` is equivalent to `>= 2.3, < 2.4`
* `~1.2.x` is equivalent to `>= 1.2.0, < 1.3.0`
* `~1.x` is equivalent to `>= 1, < 2`

#### Caret Range Comparisons (Major)

The caret (`^`) comparison operator is for major level changes once a stable
(1.0.0) release has occurred. Prior to a 1.0.0 release the minor versions acts
as the API stability level. This is useful when comparisons of API versions as a
major change is API breaking. For example,

* `^1.2.3` is equivalent to `>= 1.2.3, < 2.0.0`
* `^1.2.x` is equivalent to `>= 1.2.0, < 2.0.0`
* `^2.3` is equivalent to `>= 2.3, < 3`
* `^2.x` is equivalent to `>= 2.0.0, < 3`
* `^0.2.3` is equivalent to `>=0.2.3 <0.3.0`
* `^0.2` is equivalent to `>=0.2.0 <0.3.0`
* `^0.0.3` is equivalent to `>=0.0.3 <0.0.4`
* `^0.0` is equivalent to `>=0.0.0 <0.1.0`
* `^0` is equivalent to `>=0.0.0 <1.0.0`

[k6]: https://k6.io
[extensions]: https://grafana.com/docs/k6/latest/extensions/
[xk6-top]: https://github.com/szkiba/xk6-top
[Masterminds/semver]: https://github.com/Masterminds/semver


```
k6exec [flags] [command]
```

### Flags

```
      --build-service-url string       URL of the k6 build service to be used
      --extension-catalog-url string   URL of the k6 extension catalog to be used
  -h, --help                           help for k6
      --no-color                       disable colored output
  -q, --quiet                          disable progress updates
      --usage                          print launcher usage
  -v, --verbose                        enable verbose logging
```

<!-- #endregion cli -->

## Development

### Tasks

This section contains a description of the tasks performed during development. Commands must be issued in the repository base directory. If you have the [xc](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command in the repository base directory.

<details><summary>Click to expand</summary>

#### readme

Update documentation in README.md.

```sh
go run ./tools/gendoc README.md
```

#### lint

Run the static analyzer.

```
golangci-lint run
```

#### test

Run the tests.

```
go test -count 1 -race -coverprofile=coverage.txt ./...
```

#### coverage

View the test coverage report.

```
go tool cover -html=coverage.txt
```

#### build

Build the executable binary.

This is the easiest way to create an executable binary (although the release process uses the goreleaser tool to create release versions).

```sh
go build -ldflags="-w -s" -o k6exec ./cmd/k6exec
```

#### snapshot

Creating an executable binary with a snapshot version.

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```sh
goreleaser build --snapshot --clean --single-target -o k6exec
```

#### docker

Building a Docker image. Before building the image, it is advisable to perform a snapshot build using goreleaser. To build the image, it is advisable to use the same `Docker.goreleaser` file that `goreleaser` uses during release.

Requires: snapshot

```sh
docker build -t k6exec -f Dockerfile.goreleaser .
```

#### examples

Run all scripts in the `examples` directory with a fresh build.

Requires: clean, snapshot

```sh
find  examples -type f | xargs -n 1 ./k6exec run
```

#### clean

Delete the build directory.

```
rm -rf build
```

</details>

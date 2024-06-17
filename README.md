[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6exec.svg)](https://pkg.go.dev/github.com/grafana/k6exec)
[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6exec)](https://github.com/grafana/k6exec/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6exec)](https://goreportcard.com/report/github.com/grafana/k6exec)
[![GitHub Actions](https://github.com/grafana/k6exec/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6exec/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6exec/graph/badge.svg?token=6MP3G02V9C)](https://codecov.io/gh/grafana/k6exec)
![GitHub Downloads](https://img.shields.io/github/downloads/grafana/k6exec/total)

<h1 name="title">k6exec</h1>

**Launching k6 with extensions**

The purpose of k6exec is to launch the k6 binary containing the extensions used by the k6 test script. For this purpose, k6exec analyzes the k6 test scripts.

k6exec is primarily used as a go library for [k6](https://github.com/grafana/k6) and [xk6](https://github.com/grafana/xk6). In addition, it also contains a command-line tool, which is suitable for launching k6 with extensions based on the dependencies of k6 test scripts.

The command line tool can be integrated into other command line tools as a subcommand. For this purpose, the library also contains the functionality of the command line tool as a factrory function that returns [cobra.Command](https://pkg.go.dev/github.com/spf13/cobra#Command).

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6exec/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6exec/cmd/k6exec@latest
```

## Usage

<!-- #region cli -->
## k6exec

**Launching k6 with extensions**

Launching k6 containing the extensions used by the test script.


### Commands

* [k6exec archive](#k6exec-archive)	 - Create an archive
* [k6exec cloud](#k6exec-cloud)	 - Run a test on the cloud
* [k6exec completion](#k6exec-completion)	 - Generate the autocompletion script for the specified shell
* [k6exec help](#k6exec-help)	 - Help about any command
* [k6exec inspect](#k6exec-inspect)	 - Inspect a script or archive
* [k6exec login](#k6exec-login)	 - Authenticate with a service
* [k6exec new](#k6exec-new)	 - Create and initialize a new k6 script
* [k6exec pause](#k6exec-pause)	 - Pause a running test
* [k6exec resume](#k6exec-resume)	 - Resume a paused test
* [k6exec run](#k6exec-run)	 - Start a test
* [k6exec scale](#k6exec-scale)	 - Scale a running test
* [k6exec stats](#k6exec-stats)	 - Show test metrics
* [k6exec status](#k6exec-status)	 - Show test status
* [k6exec version](#k6exec-version)	 - Show application version

---
## k6exec archive

Create an archive

```
k6exec archive [flags]
```

### Flags

```
  -h, --help   help for archive
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec cloud

Run a test on the cloud

```
k6exec cloud [flags]
```

### Flags

```
  -h, --help   help for cloud
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec completion

Generate the autocompletion script for the specified shell

```
k6exec completion [flags]
```

### Flags

```
  -h, --help   help for completion
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec help

Help about any command

```
k6exec help [flags]
```

### Flags

```
  -h, --help   help for help
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec inspect

Inspect a script or archive

```
k6exec inspect [flags]
```

### Flags

```
  -h, --help   help for inspect
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec login

Authenticate with a service

```
k6exec login [flags]
```

### Flags

```
  -h, --help   help for login
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec new

Create and initialize a new k6 script

```
k6exec new [flags]
```

### Flags

```
  -h, --help   help for new
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec pause

Pause a running test

```
k6exec pause [flags]
```

### Flags

```
  -h, --help   help for pause
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec resume

Resume a paused test

```
k6exec resume [flags]
```

### Flags

```
  -h, --help   help for resume
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec run

Start a test

```
k6exec run [flags]
```

### Flags

```
  -h, --help   help for run
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec scale

Scale a running test

```
k6exec scale [flags]
```

### Flags

```
  -h, --help   help for scale
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec stats

Show test metrics

```
k6exec stats [flags]
```

### Flags

```
  -h, --help   help for stats
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec status

Show test status

```
k6exec status [flags]
```

### Flags

```
  -h, --help   help for status
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

---
## k6exec version

Show application version

```
k6exec version [flags]
```

### Flags

```
  -h, --help   help for version
```

### SEE ALSO

* [k6exec](#k6exec)	 - Launching k6 with extensions

<!-- #endregion cli -->

## Development

### Tasks

This section contains a description of the tasks performed during development. If you have the [xc (Markdown defined task runner)](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command.

<details><summary>Click to expand</summary>

#### readme

Update documentation in README.md.

```
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
go test -count 1 -race -coverprofile=build/coverage.txt ./...
```

#### coverage

View the test coverage report.

```
go tool cover -html=build/coverage.txt
```

#### build

Build the executable binary.

This is the easiest way to create an executable binary (although the release process uses the goreleaser tool to create release versions).

```
go build -ldflags="-w -s" -o build/k6exec ./cmd/k6exec
```

#### snapshot

Creating an executable binary with a snapshot version.

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```
goreleaser build --snapshot --clean --single-target -o build/k6exec
```

#### clean

Delete the build directory.

```
rm -rf build
```

</details>

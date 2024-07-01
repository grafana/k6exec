[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6exec.svg)](https://pkg.go.dev/github.com/grafana/k6exec)
[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6exec)](https://github.com/grafana/k6exec/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6exec)](https://goreportcard.com/report/github.com/grafana/k6exec)
[![GitHub Actions](https://github.com/grafana/k6exec/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6exec/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6exec/graph/badge.svg?token=6MP3G02V9C)](https://codecov.io/gh/grafana/k6exec)

<h1 name="title">k6exec</h1>

**Launcher for k6 with seamless use of extensions.**

k6exec is a launcher library for k6 with seamless use of extensions. The launcher will always run the k6 test script with the appropriate k6 binary, which contains the extensions used by the script. Extensions can also be recognized from the environment variable (default `K6_DEPENDENCIES`) or from the `dependencies` property of the manifest file.

## Development

### Tasks

This section contains a description of the tasks performed during development. Commands must be issued in the repository base directory. If you have the [xc](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command in the repository base directory.

<details><summary>Click to expand</summary>

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

#### clean

Delete the build directory.

```
rm -rf build
```

</details>

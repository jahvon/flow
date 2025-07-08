# flow Development

[![Go Report Card](https://goreportcard.com/badge/github.com/flowexec/flow)](https://goreportcard.com/report/github.com/flowexec/flow)
[![Go Reference](https://pkg.go.dev/badge/github.com/flowexec/flow.svg)](https://pkg.go.dev/github.com/flowexec/flow)
[![GitHub branch status](https://img.shields.io/github/checks-status/flowexec/flow/main)](https://github.com/flowexec/flow/actions?query=branch%3Amain)
[![Codecov](https://img.shields.io/codecov/c/github/flowexec/flow)](https://app.codecov.io/gh/flowexec/flow)

Before getting started, please read the [Code of Conduct](../.github/CODE_OF_CONDUCT.md) and [Contributing Guidelines](../.github/CONTRIBUTING.md).

flow is written in [Go](https://golang.org/). See the [go.mod](../go.mod) file for the current Go version used in 
building the project.

## Getting Started

Before developing on this project, you will need to make sure you have the latest `flow` version installed.
Refer to the [Installation](installation.md) section for more information.

After cloning the repository, you can start using the below commands after registering the repo workspace:

```sh
flow workspace create flow <repo-path>
```

### Development Executables

The `flow` project contains a few development executables that can be run locally. After registering the repo
workspace, you can run the following commands:

```sh
# Install Go tool dependencies
flow install tools

# Build the CLI binary
flow build binary <output-path>

# Validate code changes (runs tests, linters, codegen, etc)
flow validate

# Only generate code
flow generate

# Only run tests
flow test all
```

### Working with generated types

The `flow` project uses [go-jsonschema](github.com/atombender/go-jsonschema) with [go generate](https://blog.golang.org/generate) 
to generate Go types from JSON schema files (defined in YAML). If you need to make changes to the generated types 
(found in the `types` package), you should update the associated `*schema.yaml` file and run the flow `run generate` executable
or go generate directly.

Note that go generate alone does not update generated documentation. 
Be sure to regenerate the JSON schema files and markdown documentation before submitting a PR.

### Working with tuikit

The `flow` project uses the [tuikit](tuikit.md) framework for building the terminal UI.
Contributions to the components and helpers in `tuikit` are welcome.

_You should test all tuikit changes with a local flow build before submitting a PR._
    
```sh
  go mod edit -replace github.com/flowexec/tuikit=../tuikit
```

## Development Tools

Required tools for development:

- [mockgen](https://github.com/uber-go/mock) for generating test mocks
- [golangci-lint](https://golangci-lint.run/) for linting

Other tools used in the project:
- [goreleaser](https://goreleaser.com/) for releasing the project
- [ginkgo](https://onsi.github.io/ginkgo/) and [gomega](https://onsi.github.io/gomega/) for testing

# flow Development

Before getting started, please read the [Code of Conduct](.github/CODE_OF_CONDUCT.md) and [Contributing Guidelines](.github/CONTRIBUTING.md).

flow is written in [Go](https://golang.org/). See the [go.mod](go.mod) file for the current Go version used in 
building the project.

## Getting Started

Before developing on this project, you will need to make sure you have the latest `flow` version installed.
Refer to the [Installation](README.md#installation) section for more information.

After cloning the repository, you can start using the below commands after registering the repo workspace:

```sh
flow init workspace flow <repo-path>
```

### Development Executables

The `flow` project contains a few development executables that can be run locally. After registering the repo
workspace, you can run the following commands:

```sh
# Install local dependencies
flow install deps

# Build the project
flow run build <output-path>

# Run tests
flow run tests

# Validate code changes
flow run validate

# Install the flow binary in your $GOPATH
flow install gopath
```

### Working with generated types

The `flow` project uses [go-jsonschema](github.com/atombender/go-jsonschema) with [go generate](https://blog.golang.org/generate) 
to generate Go types from JSON schema files (defined in YAML). If you need to make changes to the generated types 
(found in the `types` package), you should update the associated `*schema.yaml` file and run the flow `run generate` executable
or go generate directly.

Note that go generate alone does not update generated documentation. 
Be sure to regenerate the JSON schema files and markdown documentation before submitting a PR.

### Working with tuikit

The `flow` project uses the [tuikit](github.com/jahvon/tuikit) framework for building the terminal UI.
Contributions to the components and helpers in `tuikit` are welcome.

_You should test all tuikit changes with a local flow build before submitting a PR._
    
```sh
  go mod edit -replace github.com/jahvon/tuikit=../tuikit
```

## Development Tools

The following tools are required for development:

- [mockgen](https://github.com/uber-go/mock) for generating test mocks
- [golangci-lint](https://golangci-lint.run/) for linting

Additionally, we recommend using the [ginkgo CLI](https://onsi.github.io/ginkgo/#ginkgo-cli-overview) for setting up and running tests.


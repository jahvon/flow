# flow Development

Before getting started, please read the [Code of Conduct](.github/CODE_OF_CONDUCT.md) and [Contributing Guidelines](.github/CONTRIBUTING.md).

flow is written in [Go](https://golang.org/). See the [go.mod](go.mod) file for the current Go version used in 
building the project.

## Development Tools

The following tools are required for development:

- [mockgen](https://github.com/uber-go/mock) for generating test mocks
- [golangci-lint](https://golangci-lint.run/) for linting

Additionally, we recommend using the [ginkgo CLI](https://onsi.github.io/ginkgo/#ginkgo-cli-overview) for setting up and running tests.

## Running a local build

## Make Commands

The following `make` commands are available for development:

|                                | Make Command      |
|--------------------------------|-------------------|
| **Install Local Dependencies** | `make local/deps` |
| **Build**                      | `make go/build`   |
| **Test**                       | `make go/test`    |
| **Pre-commit**                 | `make pre-commit` |

## Installing via Source

```bash
$ git clone github.com/jahvon/flow
$ cd flow
$ go generate ./...
$ go install
```

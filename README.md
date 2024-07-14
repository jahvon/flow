# flow 

[![Go Report Card](https://goreportcard.com/badge/github.com/jahvon/flow)](https://goreportcard.com/report/github.com/jahvon/flow)
[![Go Reference](https://pkg.go.dev/badge/github.com/jahvon/flow.svg)](https://pkg.go.dev/github.com/jahvon/flow)
[![GitHub release](https://img.shields.io/github/v/release/jahvon/flow)](https://github.com/jahvon/flow/releases)

flow is a versatile Command Line Interface (CLI) tool designed to streamline and enhance your local development and operations
workflows. Whether you're running scripts, transforming API responses, or opening applications, flow has you covered. 
It's driven by "executable" YAML configurations organized across workspaces and namespaces defined in a workspace, easily
discoverable from anywhere on your system.

## Key Features

- **Workflow Runner**: Easily define, manage, and run your workflows from the command line. Example use cases includes running a set of scripts, opening an application after running setup tasks,
  and rendering a markdown document that is dynamically generated from data in an external system.
- **Secret Management**: Safely store sensitive secrets in an encrypted local vault, easily referencable in your executable configurations.
- **Input Handling**: Pass values into executables via environment variables defined by configuration, flags, or interactive prompts, ensuring your workflows are flexible and easy to use.
- **Executable Organization**: Group your executables into workspace and namespace, making them easily discoverable and accessible from anywhere on your system. Tag your executables and workspaces for easy filtering and searching.
- **Customizable TUI**: Enjoy a seamless TUI experience with log formatting, log archiving, and execution notifications, enhancing your productivity.
- **Generate w/ Templates and Comments**: Automatically generate executable configurations from flow templates or comments in your script files, making it easy to onboard new executables.
- **Comprehensive Commands**:
    - `get` and `list`: Retrieve and display executable details, formatted as markdown documentation, YAML, or JSON.
    - `library`: Access a fully searchable and interactive TUI for running and exploring executables.

## Getting Started

### System Requirements

At this time, flow is only supported on Linux and MacOS systems. 
While it may work on Windows, it has not been tested and is not officially supported.

### Installation

You can install the pre-compiled binary or compile from source. To install the latest pre-compiled binary,
run the following command:

```bash
curl -sSL https://raw.githubusercontent.com/jahvon/flow/main/scripts/install.sh | bash
```

Alternatively, you can install the binary via Go using the following command:

```bash
go install github.com/jahvon/flow@latest
```

### Basic Usage

#### Setting up a Workspace

A Workspace is a directory that contains workflows and configuration files for flow to manage.
A workspace can be created anywhere on your system but must be registered in the user config file in order to
have executables discovered by flow.

To create a new workspace:

```bash
flow init workspace <name> <path>
```

This command will register the Workspace and create the root config file for you.
For more information on Workspaces and it's config, see [Workspaces](docs/types/workspace.md).

#### Defining Executables

Executables are the core of flow. Each executable is driven by its definition within a flow file (`*.flow`).
There are several types of executables that can be defined. For more information on Executables and the flow file, see [FlowFile.md](docs/types/flowfile.md).

There is also a JSON Schema that can be used in IDEs with the Language Server Protocol (LSP) to perform intelligent 
suggestions. You can add the following comment to the top of your flow files to enable this:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/jahvon/flow/HEAD/schemas/flowfile_schema.json
```

See the [schemas](schemas/) directory all available schemas.

#### Running and managing workflows

The main command for running workflows is `flow exec`. This command will execute the workflow with the provided
executable ID. `exec` can be replaced with any verb known to flow but should match the verb defined in the executable's 
definition or an alias of the verb.

As you make changes to executables on your system, you can run `flow sync` to trigger a re-index of your executables.

#### Autocompletion

flow supports autocompletion for bash, zsh, and fish shells. Setup depends on the shell you are using. For instance, if
you are using zsh with oh-my-zsh, you can run the following command to generate the autocompletion script:

```bash
flow completion zsh > ~/.oh-my-zsh/completions/_flow
```

## Documentation

See [flow CLI documentation](docs/cli/flow.md) for more information on all available commands.
See [flow Configuration documentation](docs/types/index.md) for more information on the configuration file options.

### Examples

Check out some example configurations and workflows in the [examples](examples/) directory. You can clone this repository
and follow the workspace setup instructions above to run these examples on your system.

## Contributing

We welcome contributions from the community! Please see our [Contributing Guide](.github/CONTRIBUTING.md) for more details on how to get involved.





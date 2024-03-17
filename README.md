# flow 

[![Go Report Card](https://goreportcard.com/badge/github.com/jahvon/flow)](https://goreportcard.com/report/github.com/jahvon/flow)
[![Go Reference](https://pkg.go.dev/badge/github.com/jahvon/flow.svg)](https://pkg.go.dev/github.com/jahvon/flow)
[![GitHub release](https://img.shields.io/github/v/release/jahvon/flow)](https://github.com/jahvon/flow/releases)

flow is a command line interface designed to make managing and running development workflows easier. It's driven by
"executables" organized across workspaces and namespaces defined in a workspace.

Some common use cases includes running a set of scripts, opening an application after running setup tasks,
and rendering a markdown document that is dynamically generated from data in an external system.
That's just the start; all aspects of flow are meant to be easily configurable, easily discoverable, and highly extensible.

## Getting Started

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

See the [Development](DEVELOPMENT.md) guide for more information on how to build from source.

### Setting up a Workspace

A Workspace is a directory that contains workflows and configuration files for `flow` to manage.
A workspace can be created anywhere on your system but must be registered in the User Config file in order to
have executables discovered by `flow`.

To create a new workspace:
    
```bash
flow init workspace <name> <path>
```

This command will register the Workspace and create the root config file for you.
For more information on Workspaces and it's config, see [Workspaces](docs/config/workspace_config.md).

### Defining Executables

Executables are the core of `flow`. They are the workflows that `flow` will execute when running a workflow.
Each executable is driven by its definition within an executable definition file (`*.flow` file). There are
several types of executables that can be defined. 
For more information on Executables and it's config, see [Executables](docs/config/executables.md).


### Running and managing workflows

The main command for running workflows is `flow exec`. This command will execute the workflow with the provided
executable ID. `exec` can be replaced with any verb but should match the verb defined in the executable's definition or
an alias of the verb.

As you make changes to executables on your system, you can run `flow sync` to trigger a re-index of your executables.

See [flow CLI documentation](docs/cli/flow.md) for more information on all available commands.

**Autocompletion**

Example autocompletion setup script: `flow completion zsh > ~/.oh-my-zsh/completions/_flow`

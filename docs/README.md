<p align="center"><img src="_media/logo.png" alt="flow" width="300"/></p>

<p align="center">
  <a href="https://img.shields.io/github/v/release/jahvon/flow"><img src="https://img.shields.io/github/v/release/jahvon/flow" alt="GitHub release"></a>
  <a href="https://pkg.go.dev/github.com/jahvon/flow"><img src="https://pkg.go.dev/badge/github.com/jahvon/flow.svg" alt="Go Reference"></a>
</p>

flow is a customizable and interactive CLI tool designed to streamline how you manage and run local development and 
operations workflows.

**This project is currently in beta and documentation is a work in progress.** We welcome contributions and feedback.

## Features

- **Task Runner**: Easily define, manage, and run your tasks (called [executables](guide/executable.md)) from the command line.
- **Secret Management**: Safely store sensitive secrets in a secure local [vault](guide/vault.md), easily referencable in your executable configurations.
- **Input Handling**: Pass values into executables via environment variables defined by configuration, flags, or interactive prompts, ensuring your workflows are flexible and easy to use.
- **Executable Organization**: Group your executables into workspace and namespace, making them easily discoverable and accessible from anywhere on your system. Tag your executables and workspaces for easy filtering and searching.
- **Customizable TUI**: Enjoy a seamless [TUI experience](guide/interactive.md) with log formatting, log archiving, and execution notifications, enhancing your productivity.
- **Generate w/ Templates and Comments**: Automatically generate executable configurations and workspace scaffolding from [flow file templates](guide/templating.md) or comments in your script files, making it easy to onboard new executables.
- **Comprehensive Commands**:
    - `get` and `list`: Retrieve and display executable details, formatted as markdown documentation, YAML, or JSON.
    - `library`: Access a fully searchable and interactive TUI for running and exploring executables.

## Installation and setup

### System Requirements

At this time, flow is only supported on Linux and MacOS systems.
While it may work on Windows, it has not been tested and is not officially supported.

### Installation
`
You can install the pre-compiled binary or compile from source. To install the latest pre-compiled binary,
run the following command:

```bash
curl -sSL https://raw.githubusercontent.com/jahvon/flow/main/scripts/install.sh | bash
```

Alternatively, you can install the binary via Go using the following command:

```bash
go install github.com/jahvon/flow@latest
```

### Autocompletion

flow supports autocompletion for bash, zsh, and fish shells. Setup depends on the shell you are using. For instance, if
you are using zsh with oh-my-zsh, you can run the following command to generate the autocompletion script:

```bash
flow completion zsh > ~/.oh-my-zsh/completions/_flow
```

## Getting started

### Setting up a workspace

A Workspace is a directory that contains workflows and configuration files for flow to manage.
A workspace can be created anywhere on your system but must be registered in the user config file in order to
have executables discovered by flow.

To create a new workspace:

```bash
flow init workspace <name> <path>
```

This command will register the Workspace and create the root config file for you.
For more information on Workspaces and it's config, see [Workspaces](/types/workspace.md).

### Defining executables

Executables are the core of flow. Each executable is driven by its definition within a flow file (`*.flow`).
There are several types of executables that can be defined. For more information on Executables and the flow file, see [FlowFile.md](types/flowfile.md).

There is also a JSON Schema that can be used in IDEs with the Language Server Protocol (LSP) to perform intelligent
suggestions. You can add the following comment to the top of your flow files to enable this:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/jahvon/flow/HEAD/schemas/flowfile_schema.json
```

See the [schemas](../schemas/) directory all available schemas.

### Running and managing executables

The main command for running executables is `flow exec`. This command will execute the workflow with the provided
executable ID. `exec` can be replaced with any verb known to flow but should match the verb defined in the flow file
configurations or an alias of that verb.

As you make changes to executables on your system, you can run `flow sync` to trigger a re-index of your executables.

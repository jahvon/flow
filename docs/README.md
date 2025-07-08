<p align="center"><img src="_media/logo.png" alt="flow" width="300"/></p>

<p align="center">
  <a href="https://img.shields.io/github/v/release/flowexec/flow"><img src="https://img.shields.io/github/v/release/flowexec/flow" alt="GitHub release"></a>
  <a href="https://pkg.go.dev/github.com/flowexec/flow"><img src="https://pkg.go.dev/badge/github.com/flowexec/flow.svg" alt="Go Reference"></a>
</p>

flow is a customizable and interactive CLI tool designed to streamline how you manage and run local development and 
operations workflows.

**This project is currently in beta and documentation is a work in progress.** Contributions and feedback are welcome.

#### _Features_ <!-- {docsify-ignore} -->

- **Task Runner**: Easily define, manage, and run your tasks (called [executables](guide/executable.md)) from the command line.
- **Secret Vault**: Store sensitive secrets in a secure local [vault](guide/secret.md#vault-setup).
- **Template Generator**: Generate executables and workspace scaffolding with [flow file templates](guide/templating.md).
- **TUI Library**: Explore and run executables from the interactive and searchable TUI [library](cli/flow_browse.md).
- **Executable Organizer**: Group, reference, and search for executables by workspace, namespace, verbs, and tags.
- **Input Handler**: Pass values into executables with environment variables defined by secrets, command-line args, or interactive prompts.
- **Customizable TUI**: Personalize your [TUI experience](guide/interactive.md) with settings for log formatting, log archiving, and execution notifications, and more.

---

<p align="center"><img src="_media/demo.gif" width="1600"></p>


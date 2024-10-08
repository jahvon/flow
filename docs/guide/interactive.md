## Interactive Configurations

The interactive TUI can be customized in the [flow config file](../types/config.md). Additionally,
there are several [flow config commands](../cli/flow_config.md) that can be used to change the TUI settings.

### Changing the log mode

There are 4 log modes available in the TUI:
- `logfmt`: Includes the log level, timestamp, and log message.
- `text`: Includes just the log message.
- `json`: Includes the log level, timestamp, and log message in JSON format.
- `hidden`: Hides the log messages.

The default log mode is `logfmt`. Use the following command to change the log mode:

```shell
flow config set log-mode (logfmt|text|json|hidden)
```

`exec` executables can be also configured to use a specific log mode. See the [flowfile configuration](../types/flowfile.md#executableexecexecutabletype) for more information.

```yaml
executables:
  - name: my-task
    exec:
      logMode: text
      cmd: echo "Hello, world!"
```

Note: the [flow logs](../cli/flow_logs.md) command will always display logs in `json` mode.

### Changing the workspace mode

There are 2 workspace modes available in the TUI:
- `fixed`: The current workspace is fixed to the one you've set with [flow workspace set](../cli/flow_workspace_set.md).
- `dynamic`: The current workspace is determined by your current working directory. If you're in a workspace directory, the TUI will automatically switch to that workspace. Otherwise, the TUI will use the workspace you've set with [flow workspace set](../cli/flow_workspace_set.md).

See the [workspace guide](workspace.md) for more information on workspaces.

### Disable the TUI

In some cases, you may want to disable the interactive TUI (in CI/CD pipelines and containers, for example). 
Use the following command will switch all TUI commands to their non-interactive counterparts:

```shell
flow config set tui false
```

Alternatively, you can set the `DISABLE_FLOW_INTERACTIVE` environment variable to `true` to disable the TUI.

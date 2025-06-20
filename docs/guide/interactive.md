## Interactive Configurations

The interactive TUI can be customized in the [flow config file](../types/config.md). Additionally,
there are several [flow config commands](../cli/flow_config.md) that can be used to change the TUI settings.

> [!TIP]
> You can view your current settings with the config view command:
> ```shell
> flow config get
> ```

### Changing the TUI theme

There are several themes available in the TUI:
- `default` (everforest)
- `light`
- `dark`
- `dracula`
- `tokyo-night`

Use the following command to change the theme:

```shell
flow config set theme (default|light|dark|dracula|tokyo-night)
```

**Overriding the theme's colors**

Additionally, you can override the theme colors by setting the `colorOverride` field in the config file. Any color not 
set in the `colorOverride` field will use the default color for the set theme.
See the [config file reference](../types/config.md#ColorPalette) for more information.

### Changing desktop notification settings

Desktop notifications can be sent when executables are completed. Use the following command to enable or disable desktop notifications:

```shell
flow config set notifications (true|false) # --sound
```

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

### Changing the default executable timeout

The global default executable timeout is 30 minutes. Use the following command to change the default executable timeout:

```shell
flow config set timeout <duration>
```

### Disable the TUI

In some cases, you may want to disable the interactive TUI (in CI/CD pipelines and containers, for example). 
Use the following command will switch all TUI commands to their non-interactive counterparts:

```shell
flow config set tui false
```

Alternatively, you can set the `DISABLE_FLOW_INTERACTIVE` environment variable to `true` to disable the TUI.

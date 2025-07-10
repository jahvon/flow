Flow is built around [executables](executable.md), which are defined in [flowfiles](../types/flowfile.md) and organized into 
directory structures called **workspaces**. Executables can also be grouped across **namespaces** nested within the flowfile trees of these workspaces.

## Workspace Registration

You can register multiple workspaces with Flow. To do so, use the [flow workspace add](../cli/flow_workspace_add.md) command:

```shell
# Create a workspace named "my-workspace" in the current directory
flow workspace add my-workspace .
# Create a workspace named "my-workspace" in a specific directory and set it as the current workspace
flow workspace add my-workspace /path/to/directory --set
```

When a workspace is created, a [configuration](#workspace-configuration) file is added to the root of the workspace if one does not already exist.

### Workspace Configuration

The workspace configuration file is a YAML file that contains the configuration options for the workspace. This file is located in the root directory of the workspace and is named `flow.yaml`.

```yaml
# Example workspace configuration
displayName: "My Project"
description: "A sample Flow workspace"
descriptionFile: README.md
tags: ["development", "web"]
verbAliases:
  run: ["exec", "start"]
  build: ["compile"]
executables:
  included: ["scripts/", "tools/"]
  excluded: ["node_modules/", ".git/"]
```

**Key Configuration Options:**

- **verbAliases**: Customize which verb aliases are available when running executables. Set to `{}` to disable all aliases. See [custom verb aliases](executable.md#custom-verb-aliases) for more details.
- **executables**: Configure which directories to include or exclude when discovering executables.
- **displayName**,  **description**, and **descriptionFile**: Used in the interactive UI and library views.
- **tags**: Used for filtering workspaces.

For more details about workspace configuration options, see [Workspace](../types/workspace.md).


## Changing the Current Workspace

To change the current workspace, use the [flow workspace switch](../cli/flow_workspace_switch.md) command:

```shell
# Set the current workspace to "my-workspace"
flow workspace switch my-workspace
```

Also see the [workspace mode](interactive.md#changing-the-workspace-mode) documentation for more information on workspace modes.

## Deleting a Workspace

To delete a workspace, use the [flow workspace remove](../cli/flow_workspace_remove.md) command:

```shell
# Delete the workspace named "my-workspace"
flow workspace remove my-workspace
```

## Workspace Viewer

The [flow workspace](../cli/flow_workspace.md) command provides various viewing and management options for workspaces:

```shell
# View the current workspace
flow workspace get
# View a specific workspace
flow workspace get my-workspace
# List all registered workspaces
flow workspace list
# List all workspaces with a specific tag
flow workspace list --tag my-tag
```

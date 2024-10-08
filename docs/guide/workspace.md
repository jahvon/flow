Flow is built around [executables](executable.md), which are defined in [flowfiles](../types/flowfile.md) and organized into 
directory structures called **workspaces**. Executables can also be grouped across **namespaces** nested within the flowfile trees of these workspaces.

## Workspace Registration

You can register multiple workspaces with Flow. To do so, use the [flow workspace create](../cli/flow_workspace_create.md) command:

```shell
# Create a workspace named "my-workspace" in the current directory
flow workspace create my-workspace .
# Create a workspace named "my-workspace" in a specific directory and set it as the current workspace
flow workspace create my-workspace /path/to/directory --set
```

When a workspace is created, a [configuration](#workspace-configuration) file is added to the root of the workspace if one does not already exist.

### Workspace Configuration

The workspace configuration file is a YAML file that contains the configuration options for the workspace. This file is located in the root directory of the workspace and is named `flow.yaml`.

For more details about workspace configuration options, see [Workspace](../types/workspace.md).


## Changing the Current Workspace

To change the current workspace, use the [flow workspace set](../cli/flow_workspace_set.md) command:

```shell
# Set the current workspace to "my-workspace"
flow workspace set my-workspace
```

Also see the [workspace mode](interactive.md#changing-the-workspace-mode) documentation for more information on workspace modes.

## Deleting a Workspace

To delete a workspace, use the [flow workspace delete](../cli/flow_workspace_delete.md) command:

```shell
# Delete the workspace named "my-workspace"
flow workspace delete my-workspace
```

## Workspace Viewer

The [flow workspace](../cli/flow_workspace.md) command provides various viewing and management options for workspaces:

```shell
# View the current workspace
flow workspace view
# View a specific workspace
flow workspace view my-workspace
# List all registered workspaces
flow workspace list
# List all workspaces with a specific tag
flow workspace list --tag my-tag
```

## flow workspace delete

Remove an existing workspace from the global configuration's workspaces list.

### Synopsis

Remove an existing workspace. File contents will remain in the corresponding directory but the workspace will be unlinked from the flow global configurations.
Note: You cannot remove the current workspace.

```
flow workspace delete NAME [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow workspace](flow_workspace.md)	 - Manage flow workspaces.


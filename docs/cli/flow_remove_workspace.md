## flow remove workspace

Remove an existing workspace from the list of known workspaces.

### Synopsis

Remove an existing workspace. File contents will remain in the corresponding directory but the workspace will be unlinked from the flow global configurations.
Note: You cannot remove the current workspace.

```
flow remove workspace <name> [flags]
```

### Options

```
  -h, --help   help for workspace
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (from 0 to 4 where 4 is most verbose)
```

### SEE ALSO

* [flow remove](flow_remove.md)	 - Remove a flow object.


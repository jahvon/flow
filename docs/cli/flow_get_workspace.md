## flow get workspace

Print a workspaces configuration. If the name is omitted, the current workspace is used.

```
flow get workspace [NAME] [flags]
```

### Options

```
  -h, --help            help for workspace
  -o, --output string   Output format. One of: yaml, json, doc, or list.
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow get](flow_get.md)	 - Print a flow entity.


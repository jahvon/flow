## flow list workspaces

Print a list of the registered flow workspaces.

```
flow list workspaces [flags]
```

### Options

```
  -h, --help              help for workspaces
  -o, --output string     Output format. One of: yaml, json, doc, or list.
  -t, --tag stringArray   Filter by tags.
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow list](flow_list.md)	 - Print a list of flow entities.


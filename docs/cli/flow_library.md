## flow library

View and manage your library of workspaces and executables.

```
flow library [flags]
```

### Options

```
  -h, --help               help for library
  -n, --namespace string   Filter executables by namespace.
  -t, --tag stringArray    Filter by tags.
  -v, --verb string        Filter executables by verb. One of: [add apply build configure delete deploy destroy edit exec generate install launch manage monitor new open refresh release reload remove render run setup show start transform undeploy uninstall update upgrade view]
  -w, --workspace string   Filter executables by workspace.
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow](flow.md)	 - flow is a command line interface designed to make managing and running development workflows easier.


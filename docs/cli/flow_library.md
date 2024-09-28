## flow library

View and manage your library of workspaces and executables.

```
flow library [flags]
```

### Options

```
  -a, --all                List from all namespaces.
  -f, --filter string      Filter executable by reference substring.
  -h, --help               help for library
  -n, --namespace string   Filter executables by namespace.
  -t, --tag stringArray    Filter by tags.
  -v, --verb string        Filter executables by verb. One of: [activate add analyze apply audit build bundle check clean clear configure create deactivate delete deploy destroy disable edit enable exec generate init install kill launch lint manage new open package patch publish purge push reboot refresh release reload remove render restart run scan send set setup show start stop teardown test transform undeploy uninstall unset update upgrade validate verify view watch]
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
* [flow library glance](flow_library_glance.md)	 - View a list of just executables.
* [flow library view](flow_library_view.md)	 - View an executable's documentation. The executable is found by reference.


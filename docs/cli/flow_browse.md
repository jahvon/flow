## flow browse

Discover and explore available executables.

### Synopsis

Browse executables across workspaces.

  flow browse                # Interactive multi-pane executable browser
  flow browse --list         # Simple list view of executables
  flow browse VERB [ID]      # Detailed view of specific executable

See https://flowexec.io/#/types/flowfile#executableverb for more information on executable verbs and https://flowexec.io/#/types/flowfile#executableref for more information on executable references.

```
flow browse [EXECUTABLE-REFERENCE] [flags]
```

### Options

```
  -a, --all                List from all namespaces.
  -f, --filter string      Filter executable by reference substring.
  -h, --help               help for browse
  -l, --list               Show a simple list view of executables instead of interactive discovery.
  -n, --namespace string   Filter executables by namespace.
  -o, --output string      Output format. One of: yaml, json, or tui. (default "tui")
  -t, --tag stringArray    Filter by tags.
  -v, --verb string        Filter executables by verb. One of: [abort activate add analyze apply build bundle check clean clear compile create deactivate debug delete deploy destroy disable edit enable erase exec execute fetch fix generate get init inspect install kill launch lint modify monitor new open package patch pause profile publish purge push reboot refresh release reload remove request reset restart retrieve run scan send set setup show start stop teardown terminate test tidy trace track transform trigger undeploy uninstall unset update upgrade validate verify view watch]
  -w, --workspace string   Filter executables by workspace.
```

### Options inherited from parent commands

```
  -L, --log-level string   Log verbosity level (debug, info, fatal) (default "info")
      --sync               Sync flow cache and workspaces
```

### SEE ALSO

* [flow](flow.md)	 - flow is a command line interface designed to make managing and running development workflows easier.


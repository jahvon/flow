## flow template view

View a flowfile template's documentation. Either it's registered name or file path can be used.

```
flow template view [flags]
```

### Options

```
  -f, --file string                  Path to the template file. It must be a valid flow file template.
  -h, --help                         help for view
  -o, --output string                Output format. One of: yaml, json, doc, or list.
  -t, --template flow set template   Registered template name. Templates can be registered in the flow configuration file or with flow set template.
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow template](flow_template.md)	 - Manage flowfile templates.


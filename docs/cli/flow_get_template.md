## flow get template

Print a flowfile template using it's registered name or file path.

```
flow get template [flags]
```

### Options

```
  -f, --file string                  Path to the template file. It must be a valid flow file template.
  -h, --help                         help for template
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

* [flow get](flow_get.md)	 - Print a flow entity.


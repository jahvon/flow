## flow list templates

Print a list of registered flowfile templates.

```
flow list templates [flags]
```

### Options

```
  -f, --file string                  Path to the template file. It must be a valid flow file template.
  -h, --help                         help for templates
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

* [flow list](flow_list.md)	 - Print a list of flow entities.


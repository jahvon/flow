## flow exec

Execute a flow by ID.

### Synopsis

Execute a flow where <executable-id> is the target executable's ID in the form of 'ws/ns:name'.
The flow subcommand used should match the target executable's verb or one of its aliases.

See https://github.com/jahvon/flow/blob/main/docs/config/executables.md#Verbfor more information on executable verbs.See https://github.com/jahvon/flow/blob/main/docs/config/executables.md#Reffor more information on executable IDs.

```
flow exec <executable-id> [flags]
```

### Options

```
  -h, --help   help for exec
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow](flow.md)	 - flow is a command line interface designed to make managing and running development workflows easier.


## flow get executable

Print an executable flow by reference.

### Synopsis

Print an executable by the executable's verb and ID.
The target executable's ID should be in the  form of 'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.

Seehttps://github.com/jahvon/flow/blob/main/docs/config/executables.md#Verbfor more information on executable verbs.Seehttps://github.com/jahvon/flow/blob/main/docs/config/executable.md#Reffor more information on executable IDs.

```
flow get executable VERB ID [flags]
```

### Options

```
  -h, --help   help for executable
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow get](flow_get.md)	 - Print a flow entity.


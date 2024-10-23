## flow store set

Set a key-value pair in the data store.

### Synopsis

The data store is a key-value store that can be used to persist data across executions. Values that are set outside of an executable will persist across all executions until they are cleared. When set within an executable, the data will only persist across serial or parallel sub-executables but all values will be cleared when the parent executable completes.

This will overwrite any existing value for the key.

```
flow store set [flags]
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow store](flow_store.md)	 - Manage the data store.


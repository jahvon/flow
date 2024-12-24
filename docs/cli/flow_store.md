## flow store

Manage the data store for persisting key-value data.

### Synopsis

Manage the flow data store - a key-value store that persists data within and across executable runs. Values set outside executables persist globally, while values set within executables persist only for that execution scope.

### Options

```
  -h, --help   help for store
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow](flow.md)	 - flow is a command line interface designed to make managing and running development workflows easier.
* [flow store clear](flow_store_clear.md)	 - Clear data from the store. Use --full to remove all stored data.
* [flow store get](flow_store_get.md)	 - Get a value from the store by its key.
* [flow store set](flow_store_set.md)	 - Set a key-value pair in the store.


## flow vault migrate

Migrate the legacy vault to a newer vault.

### Synopsis

Migrate the legacy vault to a newer vault type. The target vault must exist and the encryption key must be set for the legacy vault. Note: This will not remove the legacy vault, but will copy its contents to the target vault.

```
flow vault migrate TARGET [flags]
```

### Options

```
  -h, --help   help for migrate
```

### Options inherited from parent commands

```
  -L, --log-level string   Log verbosity level (debug, info, fatal) (default "info")
      --sync               Sync flow cache and workspaces
```

### SEE ALSO

* [flow vault](flow_vault.md)	 - Manage sensitive secret stores.


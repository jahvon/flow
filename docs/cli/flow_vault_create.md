## flow vault create

Create a new vault.

```
flow vault create NAME [flags]
```

### Options

```
  -h, --help                   help for create
      --identity-env string    Environment variable name for the Age vault identity. Only used for Age vaults.
      --identity-file string   File path for the Age vault identity. An absolute path is recommended. Only used for Age vaults.
      --key-env string         Environment variable name for the vault encryption key. Only used for AES256 vaults.
      --key-file string        File path for the vault encryption key. An absolute path is recommended. Only used for AES256 vaults.
  -p, --path string            Directory that the vault will use to store its data. If not set, the vault will be stored in the flow cache directory.
      --recipients string      Comma-separated list of recipient keys for the vault. Only used for Age vaults.
  -s, --set                    Set the newly created vault as the current vault
  -t, --type string            Vault type. Either age or aes256 (default "aes256")
```

### Options inherited from parent commands

```
  -L, --log-level string   Log verbosity level (debug, info, fatal) (default "info")
      --sync               Sync flow cache and workspaces
```

### SEE ALSO

* [flow vault](flow_vault.md)	 - Manage sensitive secret stores.


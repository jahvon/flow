## flow get secret

Print the value of a secret in the flow secret vault.

```
flow get secret REFERENCE [flags]
```

### Options

```
  -h, --help        help for secret
  -p, --plainText   Output the secret value as plain text instead of an obfuscated string
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow get](flow_get.md)	 - Print a flow entity.


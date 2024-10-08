## flow secret view

Show the value of a secret in the secret vault.

```
flow secret view REFERENCE [flags]
```

### Options

```
      --copy        Copy the secret value to the clipboard
  -h, --help        help for view
  -p, --plainText   Output the secret value as plain text instead of an obfuscated string
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow secret](flow_secret.md)	 - Manage flow secrets.


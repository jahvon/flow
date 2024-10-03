## Vault Setup

The flow CLI has an integrated vault that can be used to store secrets. The vault is encrypted using a password that you provide. 
Secrets are stored locally and encrypted with a generated key. To setup a new vault, run the following command:

```shell
flow secret vault create
```

The output will include the randomly generated key. You will need to store this key in a safe place. If you lose this key, 
you will not be able to access your secrets.

Whenever you invoke an executable that requires access to the vault or modify vault data, you will be prompted to enter 
the generated key. The `FLOW_VAULT_KEY` environment variable can be used to set the key. You could include this in your 
shell profile to avoid having to enter the key each time.

> [!TIP]
> You can create multiple vaults by repeating the above command. Switch your encryption key to switch between vaults.

## Adding Secrets

To add a secret to the vault, run the following command:

```shell
flow secret set KEY VALUE
# Alternatively, you can just include the key and the CLI will prompt you for the value
flow secret set KEY
```

## Retrieving Secrets

See the [executable guide](executable.md#environment-variables) for information on how to include secrets as executable
environment variables. The `secretRef` provided is equivalent to the key you used when adding the secret to the vault.

Additionally, you can view secrets in the vault by running the following commands:

```shell
flow secret list # List all secrets in the vault
flow secret view KEY # View the value of a specific secret
```

By default, those commands will not display the secret values. You will need to provide the `--plainText` flag to view 
the values.

## Removing Secrets

To remove a secret from the vault, run the following command:

```shell
flow secret delete KEY
```

You can also delete secrets in the interactive views when retrieving secrets.

## Backup and Restore

The vault data is stored in flow cache directory. On Linux, this is typically `~/.cache/flow/vault` or `$XDG_CACHE_HOME/flow/vault`.
On MacOS, this is typically `~/Library/Caches/flow/vault`. There is a directory for each vault you create. 

You can back up the vault by copying the directory to a safe location. To restore the vault, copy the directory back to the cache location.

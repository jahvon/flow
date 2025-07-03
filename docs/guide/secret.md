## Vault Setup

The flow CLI has an integrated vault system that can be used to store secrets securely. 
Vaults support multiple encryption backends and can be easily switched between different environments or teams.

#### Creating a New Vault

To create a new vault, use the `flow vault create` command with a name:

```shell
flow vault create development
```

This creates an AES256-encrypted vault with a randomly generated key that will be displayed in the output. 
Store this key securely - if you lose it, you won't be able to access your secrets.

#### Vault Types

flow supports multiple vault types:

<!-- tabs:start -->

#### **AES256 (Default)**

Symmetric encryption using a generated key. This is the simplest vault type - flow generates a random encryption key for you.

```shell
# Create an AES256 vault (default type)
flow vault create myapp
# or explicitly specify the type
flow vault create myapp --type aes256
```

**Key Management Options:**
```shell
# Store key in environment variable
flow vault create myapp --key-env MY_VAULT_KEY

# Store key in file
flow vault create myapp --key-file ~/mykeys/myapp.key
```

**Key Sharing:**
If you specify a `--key-env` and that environment variable already contains a valid encryption key, the vault will use that existing key instead of generating a new one:

```shell
# Set a shared key in environment
export SHARED_VAULT_KEY="your-existing-valid-key"

# Create multiple vaults using the same key
flow vault create dev --key-env SHARED_VAULT_KEY
flow vault create staging --key-env SHARED_VAULT_KEY
```

> [!NOTE]
> **Valid Key Format**: The existing key must be a base64-encoded 32-byte (256-bit) encryption key. You can generate a compatible key using `flow vault create` and copying the output, or by generating 32 random bytes and base64-encoding them. If the environment variable contains an invalid key format, vault creation will fail.

#### **Age**

Asymmetric encryption using recipient keys. This is ideal for team vaults where multiple people need access.

**Prerequisites:**
Install and use the [age-keygen](https://github.com/FiloSottile/age) tool to generate keys:

```shell
# Generate an age identity (private key)
age-keygen -o ~/.age/identity.txt

# Extract the public key (recipient) from the identity
age-keygen -y ~/.age/identity.txt
```

The public key output is what you use as recipients, and the identity file contains your private key for decryption.

**Creating Age Vaults:**
```shell
# Create vault with recipient keys
flow vault create team --type age --recipients key1,key2,key3 --identity-file ~/.age/identity.txt

# With identity environment variable
flow vault create team --type age --recipients key1,key2,key3 --identity-env MY_IDENTITY
```

<!-- tabs:end -->

See the [vault command reference](../cli/flow_vault.md) for more details on managing vaults, including listing existing vaults, switching between them, and deleting vaults.

#### Authentication

The environment variable or file that you provide at setup will be used to resolve the encryption key when accessing the vault. 
If you did not provide a key or file, these default environment variables will be used:

- For AES256 vaults: `FLOW_VAULT_KEY` environment variable
- For Age vaults: `FLOW_VAULT_IDENTITY` environment variable

At least one of the key or file will be used. You can configure key storage during vault creation:

```shell
# Expect to store the key in a specific environment variable
flow vault create myapp --key-env MY_VAULT_KEY

# Store key in file (file is created with the key if it doesn't exist)
flow vault create myapp --key-file ~/mykeys/myapp.key

# Age vault with identity file
flow vault create team --type age --identity-file ~/identities/identity.txt --identity-env MY_IDENTITY
```

#### Custom Vault Storage Location

You can specify a custom storage location for the encrypted data when creating a vault:

```shell
flow vault create myapp --path /storage/myapp
```

This data is encrypted, so you can safely store it as-is without worrying about plaintext secrets being exposed.


#### Pre-v1 Migration

If you have a (pre-v1.0) legacy vault, you can migrate it to a v1 vault:

```shell
flow vault create new-vault --key-env MY_NEW_VAULT_KEY
# Be sure to set the new and old vault keys if needed
flow vault migrate new-vault
```

This migrates secrets from the old vault format to the new named vault system. Note that this requires the old vault to 
be accessible with its key set in the `FLOW_VAULT_KEY` environment variable.

## Managing Secrets

#### Adding Secrets

To add a secret to the current vault, use the `flow secret set` command:

```shell
# Set a secret with a value
flow secret set api-key "my-secret-value"

# Set a secret interactively (you'll be prompted for the value)
flow secret set api-key
```

Secrets are stored in the currently active vault. Use `flow vault switch` to change which vault receives new secrets.

#### Retrieving Secrets

**List all secrets in the current vault:**
```shell
flow secret list
```

**Get a specific secret value:**
```shell
flow secret get api-key --plaintext
```

By default, secret commands don't display actual secret values for security. Use the `--plaintext` flag to view the actual values.

**Using secrets in executables:**
See the [executable guide](executable.md#environment-variables) for information on how to include secrets as executable environment variables. The `secretRef` provided is equivalent to the key you used when adding the secret to the vault.

#### Removing Secrets

To remove a secret from the current vault:

```shell
flow secret remove api-key
```

You can also delete secrets interactively when using the TUI views for listing secrets.

### Working with Multiple Vaults

When working with multiple vaults, secrets are isolated per vault but the vault's name can be used to reference secrets across vaults.
You can retrieve secrets from a specific vault without switching to it by using the vault name as a prefix:

```shell
# Retrieve secrets from different vaults without switching
flow secret get production/db-password
flow secret get development/api-key
```

### Backup and Restore

By default, vault data is stored in the flow cache directory, with each vault having its own directory:

- **Linux**: `~/.cache/flow/vaults/` or `$XDG_CACHE_HOME/flow/vaults/`
- **MacOS**: `~/Library/Caches/flow/vaults/`

Each vault you create gets its own configuration file and data file. 
You can back up these directories to ensure you have a copy of your vaults. 
Note that if you are using a custom storage path, you should include that in your backup strategy.

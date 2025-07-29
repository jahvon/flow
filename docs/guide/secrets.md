# Working with Secrets

flow's built-in vault keeps your sensitive data secure while making it easy to use in your workflows. 
Whether you're managing API keys, database passwords, or deployment tokens, the vault has you covered.

## Quick Start

Create your first vault and add a secret:

```shell
# Create a vault and set it as current (generates a key and shows it in output)
flow vault create my-vault --set

# Set the generated key in the default environment variable
export FLOW_VAULT_KEY="<key-from-output>"

# Add a secret (you'll be prompted for the value)
flow secret set database-password
```

```yaml
# Use it in an executable
executables:
  - verb: backup
    name: database
    exec:
      params:
        - secretRef: database-password
          envKey: DB_PASSWORD
      cmd: pg_dump -h localhost -U admin mydb
```

## Vault Types

flow supports multiple vault backends for different security needs:

<!-- tabs:start -->

#### **AES256 (Default)**

Symmetric encryption using a generated key. This is the simplest vault type - flow generates a random encryption key for you.

```shell
# Create an AES256 vault (default type)
flow vault create myapp
# or explicitly specify the type
flow vault create myapp --type aes256
```

This creates an AES256-encrypted vault with a randomly generated key that will be displayed in the output. 
Store this key securely - if you lose it, you won't be able to access your secrets.

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
# Create first vault and set the generated key
flow vault create dev --key-env SHARED_VAULT_KEY
export SHARED_VAULT_KEY="<key-from-output>"

# Create additional vaults using the same key
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

#### Pre-v1 Migration

If you have a (pre-v1.0) legacy vault, you can migrate it to a v1 vault:

```shell
flow vault create new-vault --key-env MY_NEW_VAULT_KEY
# Be sure to set the new and old vault keys if needed
flow vault migrate new-vault
```

This migrates secrets from the old vault format to the new named vault system. Note that this requires the old vault to be accessible with its key set in the `FLOW_VAULT_KEY` environment variable.

## Using Secrets in Workflows

### Basic Usage <!-- {docsify-ignore} -->

```yaml
executables:
  - verb: deploy
    name: app
    exec:
      params:
        - secretRef: api-key
          envKey: API_KEY
        - secretRef: database-url
          envKey: DATABASE_URL
      cmd: ./deploy.sh
```

### Cross-Vault References <!-- {docsify-ignore} -->

Reference secrets from different vaults:

```yaml
executables:
  - verb: sync
    name: environments
    exec:
      params:
        - secretRef: production/api-key
          envKey: PROD_API_KEY
        - secretRef: staging/api-key
          envKey: STAGING_API_KEY
      cmd: ./sync-environments.sh
```

## Secret Management

### Adding Secrets <!-- {docsify-ignore} -->

```shell
# Interactive prompt (recommended)
flow secret set my-secret

# From command line (less secure)
flow secret set my-secret "secret-value"

# From file
cat secret.txt | flow secret set my-secret
# OR
flow secret set my-secret --file secret.txt
```

### Viewing Secrets <!-- {docsify-ignore} -->

```shell
# List all secrets (values hidden)
flow secret list

# Get specific secret (obfuscated)
flow secret get my-secret

# Get plaintext value
flow secret get my-secret --plaintext

# Copy to clipboard
flow secret get my-secret --copy
```

### Updating and Removing <!-- {docsify-ignore} -->

```shell
# Update a secret (prompts for new value)
flow secret set existing-secret

# Remove a secret
flow secret remove old-secret
```

### Working with Multiple Vaults  <!-- {docsify-ignore} -->

When working with multiple vaults, secrets are isolated per vault but the vault's name can be used to reference secrets across vaults.
You can retrieve secrets from a specific vault without switching to it by using the vault name as a prefix:

```shell
# Retrieve secrets from different vaults without switching
flow secret get production/db-password
flow secret get development/api-key
```

## Vault Management

See the [vault command reference](../cli/flow_vault.md) for detailed commands and options.

### Vault Configuration <!-- {docsify-ignore} -->

```shell
# View the current vault
flow vault get

# View specific vault details
flow vault get my-vault

# Edit vault settings
flow vault edit my-vault --key-env NEW_KEY_VAR

# Remove vault (data remains, just unlinks)
flow vault remove old-vault
```

#### Custom Vault Storage Location

You can specify a custom storage location for the encrypted data when creating a vault:

```shell
flow vault create myapp --path /storage/myapp
```

This data is encrypted, so you can safely store it as-is without worrying about plaintext secrets being exposed.

### Managing Multiple Vaults <!-- {docsify-ignore} -->

Switch between vaults for different projects or environments:

```shell
# List all vaults
# Authentication for the created vaults must be resolvable by the environment variable or file you
# specified during vault creation in order to list them.
flow vault list

# Switch to a different vault
flow vault switch production

# Work with secrets in current vault
flow secret set api-key
flow secret list
```

### Backup and Recovery <!-- {docsify-ignore} -->

Vault data is stored in your flow config directory:

```shell
# Find your vaults
ls ~/.config/flow/vaults/  # Linux
ls ~/Library/Caches/flow/vaults/  # macOS

# Backup (encrypted data is safe to copy)
cp -r ~/.config/flow/vaults/ ~/backups/
```

Each vault you create gets its own configuration file and data file.
You can back up these directories to ensure you have a copy of your vaults.
Note that if you are using a custom storage path, you should include that in your backup strategy.

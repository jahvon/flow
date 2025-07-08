> **System Requirements:** At this time, flow is only supported on Linux and MacOS systems. On Linux, you will need `xclip` installed to use the clipboard feature.

### Pre-compiled binary

Run the following command to install the latest version of flow:

```shell
curl -sSL https://raw.githubusercontent.com/flowexec/flow/main/scripts/install.sh | bash
```

Alternatively, you can download the latest release from the [releases page](https://github.com/flowexec/flow/releases) and 
add the binary to your `$PATH`. A checksum is provided for each release to verify the download.

### Homebrew

```shell
brew install jahvon/tap/flow
```

### Go

```bash
go install github.com/flowexec/flow@latest
```

For CI/CD integrations and containerized environments, see the [integrations guide](guide/integrations.md).

### Autocompletion

flow supports autocompletion for bash, zsh, and fish shells. Setup depends on the shell you are using. For instance, if
you are using zsh with oh-my-zsh, you can run the following command to generate the autocompletion script:

```bash
flow completion zsh > ~/.oh-my-zsh/completions/_flow
```

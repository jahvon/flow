> **System Requirements:** At this time, flow is only supported on Linux and MacOS systems. On Linux, you will need `xclip` installed to use the clipboard feature.

### Pre-compiled binary

Run the following command to install the latest version of flow:

```shell
curl -sSL https://raw.githubusercontent.com/jahvon/flow/main/scripts/install.sh | bash
```

Alternatively, you can download the latest release from the [releases page](https://github.com/jahvon/flow/releases) and 
add the binary to your `$PATH`. A checksum is provided for each release to verify the download.

### Homebrew

```shell
brew install jahvon/tap/flow
```

### Go

```bash
go install github.com/jahvon/flow@latest
```

### Docker (experimental)

You can also run flow in a Docker container. This is useful if you want to run flow in a CI/CD pipeline or in a containerized environment.
The `GIT_REPO`, `BRANCH`, and `WORKSPACE` environment variables are optional and can be used to clone a specific flow workspace.

```shell
docker run -it --rm -t ghcr.io/jahvon/flow -e GIT_REPO=$GIT_REPO -e BRANCH=$BRANCH -e WORKSPACE=$WORKSPACE
```

### Autocompletion

flow supports autocompletion for bash, zsh, and fish shells. Setup depends on the shell you are using. For instance, if
you are using zsh with oh-my-zsh, you can run the following command to generate the autocompletion script:

```bash
flow completion zsh > ~/.oh-my-zsh/completions/_flow
```

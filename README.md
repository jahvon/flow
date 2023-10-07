# flow - Local, CLI Workflow Manager

## Configuration and Workflow Definition Files

`flow` uses a configuration file to define the workspaces and workflows that it manages. 
The location for this flow configuration file is `~/.flow/config.yaml`. Below is an example of a flow configuration file.

```yaml
currentWorkspace: workspace1
workspaces:
  workspace1: /path/to/workspace1
  workspace2: /path/to/workspace2
backends:
  secret:
    backend: envFile
  auth:
    backend: keyring
    preferredMode: password
    rememberMe: true
    rememberDuration: 24h
```

Workflows can be defined anywhere within a workspace's directory with the `.flow` file extension.
Below is an example of a workflow definition file.

```yaml
namespace: ns
tags:
  - example
executables:
  - type: open
    name: config
    tags:
      - config
    description: open flow config in vscode
    spec:
      args:
        - ~/.flow
      application: code
```

Running `flow open ns:config` will run the above workflow.

## CLI Usage

```bash
flow [command]
``` 

### Available Commands

**Updating configs and data used by `flow`**

- `flow get` - Get current value of various configuration options and data
- `flow set` - Update various configuration options and data
- `flow create` - Create ...
- `flow delete` - Delete ...
- `flow login` - Login to auth backend. This is needed only when an auth backend is set.

**Executing workflows**

- `flow run` - Run a workflow
- `flow open` - Open a workflow in the browser

**Autocompletion**

- `flow completion` - Generate shell autocompletion for `flow`

Example autocompletion setup script: `flow completion zsh > ~/.oh-my-zsh/completions/_flow`

## Install

You can install the pre-compiled binary (in several ways), use Docker or compile from source (when on OSS).

Below you can find the steps for each of them.

<details>
  <summary><h3>homebrew tap</h3></summary>

```bash
brew install jahvon/tap/flow
```

</details>

<details>
  <summary><h3>apt</h3></summary>

```bash
echo 'deb [trusted=yes] https://apt.fury.io/jahvon/ /' | sudo tee /etc/apt/sources.list.d/jahvon.list
sudo apt update
sudo apt install flow
```

</details>

<details>
  <summary><h3>yum</h3></summary>

```bash
echo '[jahvon]
name=Gemfury jahvon repository
baseurl=https://yum.fury.io/jahvon/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/jahvon.repo
sudo yum install goreleaser
```

</details>

<details>
  <summary><h3>deb, rpm and apk packages</h3></summary>
Download the .deb, .rpm or .apk packages from the [release page](https://github.com/jahvon/flow/releases) and install them with the appropriate tools.
</details>

<details>
  <summary><h3>go install</h3></summary>

```bash
go install github.com/jahvon/flow@latest
```

</details>

<details>
  <summary><h3>from the GitHub releases</h3></summary>

Download the pre-compiled binaries from the [release page](https://github.com/jahvon/flow/releases) page and copy them to the desired location.

```bash
$ VERSION=v1.0.0
$ OS=Linux
$ ARCH=x86_64
$ TAR_FILE=flow_${OS}_${ARCH}.tar.gz
$ wget https://github.com/jahvon/flow/releases/download/${VERSION}/${TAR_FILE}
$ sudo tar xvf ${TAR_FILE} flow -C /usr/local/bin
$ rm -f ${TAR_FILE}
```

</details>

<details>
  <summary><h3>manually</h3></summary>

```bash
$ git clone github.com/jahvon/flow
$ cd flow
$ go generate ./...
$ go install
```

</details>

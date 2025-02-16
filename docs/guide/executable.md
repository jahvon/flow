## Finding Executables

Executables are customizable actions defined in a YAML [flowfile](#flowfile). There are a few [flow library](../cli/flow_library.md)
command that can be used to find executables:
    
```shell
flow library # Multi-pane view for browsing executables
flow library glance # Single-pane view for browsing executables
flow library view # View documentation for a single executable
```

The `flow library` and `flow library glance` commands accept optional command-line flags to filter the list of 
executables by workspace, namespace, verb, or tag:

```shell
flow library --workspace ws --namespace ns --verb exec --tag my-tag 
# additionally, the --all flag can be used to show executables from all namespaces and the 
# --filter flag can be used to search the executable names and descriptions
flow library --all --filter "search string"
```

## Running Executables

The [flow exec](../cli/flow_exec.md) command is used 
to run executables, which can be individual tasks or workflows ([serial](#serial) & [parallel](#parallel)). 

```shell
flow VERB EXECUTABLE_ID [flags] [args]
# by default, the verb is 'exec'
flow exec EXECUTABLE_ID [flags] [args]
```

**Verbs**

The "verb" is a single word that describes the operation being executed. It can be configured in the flowfile for 
each executable.

When running in the CLI, the configured verb can be replaced with any synonym/alias that describes the operation.

For instance, `flow test my-app` is equivalent to `flow validate my-app`. This allows for a more natural language-like 
interaction with the CLI, making it easier to remember and use. 
*See the [verb reference](../types/flowfile.md#verb-groups) for a list all verbs and their synonyms.*

> [!TIP]
> Create shell aliases for commonly used verbs to make running executables easier. For example:
> ```shell
> alias build="flow build"
> ```
> This allows you to run `build my-app` instead of `flow build my-app` or the synonym `flow package my-app`.

**Executable IDs**

Executables are identified by their unique ID, which is a combination of the workspace, namespace, and name - using the 
format `workspace/namespace:name`. If the workspace and namespace are omitted (`ws/name`, `ns:name`, or just `name`), 
the current workspace and namespace are assumed. See the [workspace guide](workspace.md) for more information on workspaces and namespaces.

The name of an executable can also be replaced with an alias if one is defined in the flowfile.

**Nameless Executables**

If an executable is defined without a namespace and name, it is considered "nameless" and can be run by its verb alone.
Nameless executables cannot exist within a namespace.

```shell
# run the current workspaces nameless executable with the verb 'validate'
flow validate
# run the (non-current) 'my-project' workspaces nameless executable with the verb 'build'
flow build my-project/
```

## Flowfile

The flowfile is the primary configuration file that defines what an executable should do. The file is written in YAML but
should have a `.flow` extension.

The [flow sync](../cli/flow_sync.md) command is used to trigger a discovery of executables within workspaces. This 
command should be run whenever an executable is created, moved, or deleted from a flowfile.

```shell
flow sync
# alternatively, use the --sync flag with any other command to sync before running
flow exec my-task --sync
```

On sync, the CLI reads all flowfiles in your workspaces and updates the index of executables. You can configure where
the CLI should look for flowfiles in your [workspace configuration](workspace.md).

**Example Structure**

Below is an example of a flowfile. It contains a single executable named `my-task` that prints a message to the console.
_See the [flowfile reference](../types/flowfile.md) for a more detailed explanation of the flowfile schema._

```yaml
visibility: internal
namespace: example
executables:
  - name: my-task
    exec:
      cmd: echo "Hello, world!"
```

### Executables

The only required field in an executable's configuration is the `name` (if the `namespace` is unset).

Additionally, you can define the following fields:

- **visibility**: The visibility of the executable. Can be `public`, `private`, `internal`, or `hidden`.
  - `public`: viewable in the library and can be run from any workspace.
  - `private`: viewable in the library but can only be run if your current workspace is the same as the executable's workspace.
  - `internal`: not viewable in the library and can only be run from the same workspace.
  - `hidden` not viewable in the library and cannot be run from the CLI.
- **description**: Markdown description of the executable to display in the library.
- **tags**: A list of tags to categorize the executable.
- **aliases**: A list of alternative names for the executable.
- **timeout**: The maximum time the executable is allowed to run before being terminated.

One of the following executable types must be defined:

- [exec](#exec): Execute a command directly in the shell.
- [serial](#serial): Run a list of executables sequentially.
- [parallel](#parallel): Run a list of executables concurrently.
- [launch](#launch): Open a service or application.
- [request](#request): Make HTTP requests to APIs.
- [render](#render): Generate and view markdown created dynamically with templates or configurations.

#### Environment variables

Environment variables are used to customize executable behavior without modifying the flowfile itself. 
In addition to inheriting environment variables from the shell, you can define custom environment variables by 
setting `params` or `args` in an executable's configuration. Here's an overview of the different options:

**Params**

```yaml
executables:
  - verb: "deploy"
    name: "devbox"
    exec:
      file: "dev-deploy.sh"
      params:
        # secret param
        - secretRef: "dev-api-token"
          envKey: "API_TOKEN"
        # prompt / form param
        - prompt: "What application are you deploying?"
          envKey: "APP_NAME"
        # static param
        - text: "false"
          envKey: "DRY_RUN"
```

In the example above, the `devbox` executable has three parameters: `API_TOKEN`, `APP_NAME`, and `PORT`.
The `secretRef` parameter type is used to reference a secret stored in the vault (see the [secret vault](secret.md) guide 
for more information). The `prompt` parameter type prompts the user for input when the executable is run. The `text`
parameter type sets a static value for the environment variable.

_This example used the `exec` type, but the `params` field can be used with any executable type._

**Args**


```yaml
executables:
  - verb: "build"
    name: "container"
    exec:
      file: "build-image.sh"
      args:
        # positional argument
        - pos: 1
          envKey: "TAG"
          required: true
        # flag arguments
        - flag: "publish"
          envKey: "PUBLISH"
          type: "bool"
        - flag: "builder"
          envKey: "BUILDER"
          default: "podman"
```

In the example above, the `container` executable has three arguments: `TAG`, `PUBLISH`, and `BUILDER`.
The `pos` argument type is a positional CLI argument that must be provided when running the executable. The `flag` argument
type is a "flag" that can be set when running the executable. The `type` field can be used to specify the type to validate
the argument against. 

Here is an example of how this would be run:

```shell
flow build container 1.0.0 publish=true builder=docker
```

_This example used the `exec` type, but the `args` field can be used with any executable type._

#### Changing directories

You can use the `dir` field in the executable configuration to specify the working directory for the executable. By default,
the working directory is the directory where the flowfile is located. The path can include environment variables in the 
form of `$VAR`. Additionally, the following prefixes can be used to reference specific directories:

- `//`: workspace root directory
- `~/`: user home directory
- `./`: current working directory

Other values are assumed to be relative to the flowfile's directory.

Specifying the value `f:tmp` will create a temporary directory that is automatically cleaned up after the executable finishes.

```yaml
executables:
  - verb: "clean"
    name: "downloads"
    exec:
      file: "cleanup.sh"
      dir: "$HOME/downloads"
  - verb: "build"
    name: "app"
    exec:
      file: "build.sh"
      dir: "//"
  - verb: "test"
    name: "unit"
    exec:
      cmd: "cp $HOME/unit-tests.sh . && ./unit-tests.sh"
      dir: "f:tmp"
```

_This example used the `exec` type, but the `dir` field can be used with the `serial` and `parallel` types as well._

### Executable Type Examples

> [!TIP]
> Check out the [flowfile examples](https://github.com/jahvon/flow/tree/main/examples) found on GitHub for more examples 
> of executable configurations.

##### exec

The `exec` type is used to run a command directly in the shell. The command can be a single command or a script file.

```yaml
fromFile:
  - "generated.sh"
executables:
  - verb: "init"
    name: "chezmoi"
    exec:
      file: "init-chezmoi.sh"
  - verb: "apply"
    name: "dotfiles"
    exec:
      cmd: "chezmoi apply"
```

**Generated Executable**

Executables can also be generated from comments in a `.sh` file. Include the file name in the `fromFile` field of the flowfile
and add flow comments to the file. The comments should be in the format `f:key=value` to define the executable properties.

The following keys are supported:

- `name`
- `verb`
- `description`
- `tags`
- `aliases`
- `visibility`
- `timeout`

Here is an example of valid script with flow comments. Note that any comments after the first non-comment line will be
ignored.

```shell
#!/bin/sh

# f:name=generated f:verb=run
# f:description="start of the description"
# f:tags=example,generated
# I'm ignoring this comment
#
# <f|description>
# continued description
# <f|description>

echo "Hello from a generated executable!"
````

##### serial

The `serial` type is used to run a list of executables sequentially. For each `exec` in the list, you must define
either a `ref` to another executable or a `cmd` to run directly in the shell.

The [executable environment variables](#environment-variables) and [executable directory](#changing-directories)
of the parent executable are inherited by the child executables.

```yaml
executables:
  - verb: "setup"
    name: "flow-system"
    serial:
      failFast: true # Setting `failFast` will stop the execution if any of the executables fail
      execs:
        - ref: "upgrade flow:cli"
        - ref: "install flow:workspaces"
          args: ["all"] # When referring to another executable that requires arguments, you can pass them in the `args` field
          retries: 2 # retry the executable up to 2 times if it fails
        - ref: "apply flow:config"
          reviewRequired: true # serial execution will pause here for review (user input)
        - cmd: "flow sync"
```

##### parallel

The `parallel` type is used to run a list of executables concurrently. For each `exec` in the list, you must define
either a `ref` to another executable or a `cmd` to run directly in the shell.

The [executable environment variables](#environment-variables) and [executable directory](#changing-directories)
of the parent executable are inherited by the child executables.

```yaml
executables:
  - verb: "deploy"
    name: "apps"
    parallel:
      failFast: true # Setting `failFast` will stop the execution if any of the executables fail
      maxThreads: 2 # Setting `maxThreads` will limit the number of concurrent threads that can be run at once
      execs:
        - ref: "deploy helm-app"
          args: ["tailscale"] # When referring to another executable that requires arguments, you can pass them in the `args` field
        - ref: "deploy helm-app"
          args: ["metallb"]
        - ref: "setup gloo-gateway"
          retries: 1 # retry the executable up to 1 time if it fails
        - cmd: "kubectl apply -f external-services.yaml"
```

##### launch

The `launch` type is used to open a service or application. The `uri` field is required and can include environment variables
(including those resolved from params and args)

```yaml
executables:
  - verb: "open"
    name: "workspace"
    launch:
      uri: "$FLOW_WORKSPACE_PATH"
      app: "Visual Studio Code" # optional application to open the URI with
      wait: true # wait for the application to close before continuing
```

##### request

The `request` type is used to make HTTP requests to APIs. The `url` field is required, and the `method` field defaults 
to `GET`.

Additionally, you can define the `body` field to include a request body and the `headers` field to include custom headers.

```yaml
executables:
  - verb: "pause"
    name: "pihole"
    request:
      method: "POST"
      url: "http://pi.hole/admin/api.php?disable=$DURATION&auth=$PWHASH"
      logResponse: true # log the response body
      validStatusCodes: [200] # only consider the execution successful if the status code is 200
      # transform the response body with a Expr expression
      transformResponse: |
        "paused: " + string(fromJSON(body)["status"] == "disabled")
```

##### render

The `render` type is used to generate and view markdown created dynamically with templates or configurations. 
The `templateFile` field is required and can include environment variables (including those resolved from params and args).

The markdown template can include [Go template](https://pkg.go.dev/text/template) syntax to dynamically generate content.
[Sprig functions](https://masterminds.github.io/sprig/) are also available for use in the template.

```yaml
executables:
  - verb: "show"
    name: "cluster-summary"
    render:
      templateFile: "cluster-template.md"
      # Optionally, you can define a data file to use with the template
      # It can be a JSON or YAML file.
      templateDataFile: "kubectl-out.json"
```

[comment]: # (Documentation autogenerated by docsgen. Do not edit directly.)

# FlowFile

Configuration for a group of Flow CLI executables. The file must have the extension `.flow`, `.flow.yaml`, or `.flow.yml` 
in order to be discovered by the CLI. It's configuration is used to define a group of executables with shared metadata 
(namespace, tags, etc). A workspace can have multiple flow files located anywhere in the workspace directory


## Properties


**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `description` | A description of the executables defined within the flow file. This description will used as a shared description for all executables in the flow file.  | `string` |  |  |
| `descriptionFile` | A path to a markdown file that contains the description of the executables defined within the flow file. | `string` |  |  |
| `executables` |  | `array` ([Executable](#Executable)) | [] |  |
| `fromFile` | DEPRECATED: Use `imports` instead | [FromFile](#FromFile) | [] |  |
| `imports` |  | [FromFile](#FromFile) | [] |  |
| `namespace` | The namespace to be given to all executables in the flow file. If not set, the executables in the file will be grouped into the root (*) namespace.  Namespaces can be reused across multiple flow files.  Namespaces are used to reference executables in the CLI using the format `workspace:namespace/name`.  | `string` |  |  |
| `tags` | Tags to be applied to all executables defined within the flow file. | `array` (`string`) | [] |  |
| `visibility` |  | [CommonVisibility](#CommonVisibility) | <no value> |  |


## Definitions

### CommonAliases

Alternate names that can be used to reference the executable in the CLI.

**Type:** `array` (`string`)




### CommonTags

A list of tags.
Tags can be used with list commands to filter returned data.


**Type:** `array` (`string`)




### CommonVisibility

The visibility of the executables to Flow.
If not set, the visibility will default to `public`.

`public` executables can be executed and listed from anywhere.
`private` executables can be executed and listed only within their own workspace.
`internal` executables can be executed within their own workspace but are not listed.
`hidden` executables cannot be executed or listed.


**Type:** `string`
**Default:** `public`
**Valid values:**
- `public`
- `private`
- `internal`
- `hidden`



### Executable

The executable schema defines the structure of an executable in the Flow CLI.
Executables are the building blocks of workflows and are used to define the actions that can be performed in a workspace.


**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `aliases` |  | [CommonAliases](#CommonAliases) | [] |  |
| `description` | A description of the executable. This description is rendered as markdown in the interactive UI.  | `string` |  |  |
| `exec` |  | [ExecutableExecExecutableType](#ExecutableExecExecutableType) | <no value> |  |
| `launch` |  | [ExecutableLaunchExecutableType](#ExecutableLaunchExecutableType) | <no value> |  |
| `name` | An optional name for the executable.  Name is used to reference the executable in the CLI using the format `workspace/namespace:name`. [Verb group + Name] must be unique within the namespace of the workspace.  | `string` |  |  |
| `parallel` |  | [ExecutableParallelExecutableType](#ExecutableParallelExecutableType) | <no value> |  |
| `render` |  | [ExecutableRenderExecutableType](#ExecutableRenderExecutableType) | <no value> |  |
| `request` |  | [ExecutableRequestExecutableType](#ExecutableRequestExecutableType) | <no value> |  |
| `serial` |  | [ExecutableSerialExecutableType](#ExecutableSerialExecutableType) | <no value> |  |
| `tags` |  | [CommonTags](#CommonTags) | [] |  |
| `timeout` | The maximum amount of time the executable is allowed to run before being terminated. The timeout is specified in Go duration format (e.g. 30s, 5m, 1h).  | `string` | <no value> |  |
| `verb` |  | [ExecutableVerb](#ExecutableVerb) | exec | ✘ |
| `verbAliases` | A list of aliases for the verb. This allows the executable to be referenced with multiple verbs. | `array` ([Verb](#Verb)) | [] |  |
| `visibility` |  | [CommonVisibility](#CommonVisibility) | <no value> |  |

### ExecutableArgument



**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `default` | The default value to use if the argument is not provided. If the argument is required and no default is provided, the executable will fail.  | `string` |  |  |
| `envKey` | The name of the environment variable that will be assigned the value. | `string` |  |  |
| `flag` | The flag to use when setting the argument from the command line. Either `flag` or `pos` must be set, but not both.  | `string` |  |  |
| `outputFile` | A path where the argument value will be temporarily written to disk. The file will be created before execution and cleaned up afterwards.  | `string` |  |  |
| `pos` | The position of the argument in the command line ArgumentList. Values start at 1. Either `flag` or `pos` must be set, but not both.  | `integer` | <no value> |  |
| `required` | If the argument is required, the executable will fail if the argument is not provided. If the argument is not required, the default value will be used if the argument is not provided.  | `boolean` | false |  |
| `type` | The type of the argument. This is used to determine how to parse the value of the argument. | `string` | string |  |

### ExecutableArgumentList



**Type:** `array` ([ExecutableArgument](#ExecutableArgument))




### ExecutableDirectory

The directory to execute the command in.
If unset, the directory of the flow file will be used.
If set to `f:tmp`, a temporary directory will be created for the process.
If prefixed with `./`, the path will be relative to the current working directory.
If prefixed with `//`, the path will be relative to the workspace root.
Environment variables in the path will be expended at runtime.


**Type:** `string`




### ExecutableExecExecutableType

Standard executable type. Runs a command/file in a subprocess.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `cmd` | The command to execute. Only one of `cmd` or `file` must be set.  | `string` |  |  |
| `dir` |  | [ExecutableDirectory](#ExecutableDirectory) |  |  |
| `file` | The file to execute. Only one of `cmd` or `file` must be set.  | `string` |  |  |
| `logMode` | The log mode to use when running the executable. This can either be `hidden`, `json`, `logfmt` or `text`  | `string` | logfmt |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |

### ExecutableLaunchExecutableType

Launches an application or opens a URI.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `app` | The application to launch the URI with. | `string` |  |  |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |
| `uri` | The URI to launch. This can be a file path or a web URL. | `string` |  | ✘ |

### ExecutableParallelExecutableType



**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `dir` |  | [ExecutableDirectory](#ExecutableDirectory) |  |  |
| `execs` | A list of executables to run in parallel. Each executable can be a command or a reference to another executable.  | [ExecutableParallelRefConfigList](#ExecutableParallelRefConfigList) | <no value> | ✘ |
| `failFast` | End the parallel execution as soon as an exec exits with a non-zero status. This is the default behavior. When set to false, all execs will be run regardless of the exit status of parallel execs.  | `boolean` | <no value> |  |
| `maxThreads` | The maximum number of threads to use when executing the parallel executables. | `integer` | 5 |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |

### ExecutableParallelRefConfig

Configuration for a parallel executable.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` | Arguments to pass to the executable. | `array` (`string`) | [] |  |
| `cmd` | The command to execute. One of `cmd` or `ref` must be set.  | `string` |  |  |
| `if` | An expression that determines whether the executable should run, using the Expr language syntax. The expression is evaluated at runtime and must resolve to a boolean value.  The expression has access to OS/architecture information (os, arch), environment variables (env), stored data (store), and context information (ctx) like workspace and paths.  For example, `os == "darwin"` will only run on macOS, `len(store["feature"]) > 0` will run if a value exists in the store, and `env["CI"] == "true"` will run in CI environments. See the [Expr documentation](https://expr-lang.org/docs/language-definition) for more information.  | `string` |  |  |
| `ref` | A reference to another executable to run in serial. One of `cmd` or `ref` must be set.  | [ExecutableRef](#ExecutableRef) |  |  |
| `retries` | The number of times to retry the executable if it fails. | `integer` | 0 |  |

### ExecutableParallelRefConfigList

A list of executables to run in parallel. The executables can be defined by it's exec `cmd` or `ref`.


**Type:** `array` ([ExecutableParallelRefConfig](#ExecutableParallelRefConfig))




### ExecutableParameter

A parameter is a value that can be passed to an executable and all of its sub-executables.
Only one of `text`, `secretRef`, `prompt`, or `file` must be set. Specifying more than one will result in an error.


**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `envKey` | The name of the environment variable that will be assigned the value. | `string` |  |  |
| `outputFile` | A path where the parameter value will be temporarily written to disk. The file will be created before execution and cleaned up afterwards.  | `string` |  |  |
| `prompt` | A prompt to be displayed to the user when collecting an input value. | `string` |  |  |
| `secretRef` | A reference to a secret to be passed to the executable. | `string` |  |  |
| `text` | A static value to be passed to the executable. | `string` |  |  |

### ExecutableParameterList



**Type:** `array` ([ExecutableParameter](#ExecutableParameter))




### ExecutableRef

A reference to an executable.
The format is `<verb> <workspace>/<namespace>:<executable name>`.
For example, `exec ws/ns:my-workflow`.

- If the workspace is not specified, the current workspace will be used.
- If the namespace is not specified, the current namespace will be used.
- Excluding the name will reference the executable with a matching verb but an unspecified name and namespace (e.g. `exec ws` or simply `exec`).


**Type:** `string`




### ExecutableRenderExecutableType

Renders a markdown template file with data.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `dir` |  | [ExecutableDirectory](#ExecutableDirectory) |  |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |
| `templateDataFile` | The path to the JSON or YAML file containing the template data. | `string` |  |  |
| `templateFile` | The path to the markdown template file to render. | `string` |  |  |

### ExecutableRequestExecutableType

Makes an HTTP request.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `body` | The body of the request. | `string` |  |  |
| `headers` | A map of headers to include in the request. | `map` (`string` -> `string`) | map[] |  |
| `logResponse` | If set to true, the response will be logged as program output. | `boolean` | false |  |
| `method` | The HTTP method to use when making the request. | `string` | GET |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |
| `responseFile` |  | [ExecutableRequestResponseFile](#ExecutableRequestResponseFile) | <no value> |  |
| `timeout` | The timeout for the request in Go duration format (e.g. 30s, 5m, 1h). | `string` | 30m0s |  |
| `transformResponse` | [Expr](https://expr-lang.org/docs/language-definition) expression used to transform the response before saving it to a file or outputting it.  The following variables are available in the expression:   - `status`: The response status string.   - `code`: The response status code.   - `body`: The response body.   - `headers`: The response headers.  For example, to capitalize a JSON body field's value, you can use `upper(fromJSON(body)["field"])`.  | `string` |  |  |
| `url` | The URL to make the request to. | `string` |  | ✘ |
| `validStatusCodes` | A list of valid status codes. If the response status code is not in this list, the executable will fail. If not set, the response status code will not be checked.  | `array` (`integer`) | [] |  |

### ExecutableRequestResponseFile

Configuration for saving the response of a request to a file.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `dir` |  | [ExecutableDirectory](#ExecutableDirectory) |  |  |
| `filename` | The name of the file to save the response to. | `string` |  | ✘ |
| `saveAs` | The format to save the response as. | `string` | raw |  |

### ExecutableSerialExecutableType

Executes a list of executables in serial.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` |  | [ExecutableArgumentList](#ExecutableArgumentList) | <no value> |  |
| `dir` |  | [ExecutableDirectory](#ExecutableDirectory) |  |  |
| `execs` | A list of executables to run in serial. Each executable can be a command or a reference to another executable.  | [ExecutableSerialRefConfigList](#ExecutableSerialRefConfigList) | <no value> | ✘ |
| `failFast` | End the serial execution as soon as an exec exits with a non-zero status. This is the default behavior. When set to false, all execs will be run regardless of the exit status of the previous exec.  | `boolean` | <no value> |  |
| `params` |  | [ExecutableParameterList](#ExecutableParameterList) | <no value> |  |

### ExecutableSerialRefConfig

Configuration for a serial executable.

**Type:** `object`



**Properties:**

| Field | Description | Type | Default | Required |
| ----- | ----------- | ---- | ------- | :--------: |
| `args` | Arguments to pass to the executable. | `array` (`string`) | [] |  |
| `cmd` | The command to execute. One of `cmd` or `ref` must be set.  | `string` |  |  |
| `if` | An expression that determines whether the executable should run, using the Expr language syntax. The expression is evaluated at runtime and must resolve to a boolean value.  The expression has access to OS/architecture information (os, arch), environment variables (env), stored data (store), and context information (ctx) like workspace and paths.  For example, `os == "darwin"` will only run on macOS, `len(store["feature"]) > 0` will run if a value exists in the store, and `env["CI"] == "true"` will run in CI environments. See the [Expr documentation](https://expr-lang.org/docs/language-definition) for more information.  | `string` |  |  |
| `ref` | A reference to another executable to run in serial. One of `cmd` or `ref` must be set.  | [ExecutableRef](#ExecutableRef) |  |  |
| `retries` | The number of times to retry the executable if it fails. | `integer` | 0 |  |
| `reviewRequired` | If set to true, the user will be prompted to review the output of the executable before continuing. | `boolean` | false |  |

### ExecutableSerialRefConfigList

A list of executables to run in serial. The executables can be defined by it's exec `cmd` or `ref`.


**Type:** `array` ([ExecutableSerialRefConfig](#ExecutableSerialRefConfig))




### ExecutableVerb

Keywords that describe the action an executable performs. Executables are configured with a single verb,
but core verbs have aliases that can be used interchangeably when referencing executables. This allows users 
to use the verb that best describes the action they are performing.

### Default Verb Aliases

- **Execution Group**: `exec`, `run`, `execute`
- **Retrieval Group**: `get`, `fetch`, `retrieve`
- **Display Group**: `show`, `view`, `list`
- **Configuration Group**: `configure`, `setup`
- **Update Group**: `update`, `upgrade`

### Usage Notes

1. [Verb + Name] must be unique within the namespace of the workspace.
2. When referencing an executable, users can use any verb from the default or configured alias group.
3. All other verbs are standalone and self-descriptive.

### Examples

- An executable configured with the `exec` verb can also be referenced using "run" or "execute".
- An executable configured with `get` can also be called with "list", "show", or "view".
- Operations like `backup`, `migrate`, `flush` are standalone verbs without aliases.
- Use domain-specific verbs like `deploy`, `scale`, `tunnel` for clear operational intent.

By providing minimal aliasing with comprehensive verb coverage, flow enables natural language operations
while maintaining simplicity and flexibility for diverse development and operations workflows.


**Type:** `string`
**Default:** `exec`
**Valid values:**
- `abort`
- `activate`
- `add`
- `analyze`
- `apply`
- `archive`
- `audit`
- `backup`
- `benchmark`
- `build`
- `bundle`
- `check`
- `clean`
- `clear`
- `commit`
- `compile`
- `compress`
- `configure`
- `connect`
- `copy`
- `create`
- `deactivate`
- `debug`
- `decompress`
- `decrypt`
- `delete`
- `deploy`
- `destroy`
- `disable`
- `disconnect`
- `edit`
- `enable`
- `encrypt`
- `erase`
- `exec`
- `execute`
- `export`
- `expose`
- `fetch`
- `fix`
- `flush`
- `format`
- `generate`
- `get`
- `import`
- `index`
- `init`
- `inspect`
- `install`
- `join`
- `kill`
- `launch`
- `lint`
- `list`
- `load`
- `lock`
- `login`
- `logout`
- `manage`
- `merge`
- `migrate`
- `modify`
- `monitor`
- `mount`
- `new`
- `notify`
- `open`
- `package`
- `partition`
- `patch`
- `pause`
- `ping`
- `preload`
- `prefetch`
- `profile`
- `provision`
- `publish`
- `purge`
- `push`
- `queue`
- `reboot`
- `recover`
- `refresh`
- `release`
- `reload`
- `remove`
- `request`
- `reset`
- `restart`
- `restore`
- `retrieve`
- `rollback`
- `run`
- `save`
- `scale`
- `scan`
- `schedule`
- `seed`
- `send`
- `serve`
- `set`
- `setup`
- `show`
- `snapshot`
- `start`
- `stash`
- `stop`
- `tag`
- `teardown`
- `terminate`
- `test`
- `tidy`
- `trace`
- `transform`
- `trigger`
- `tunnel`
- `undeploy`
- `uninstall`
- `unmount`
- `unset`
- `update`
- `upgrade`
- `validate`
- `verify`
- `view`
- `watch`



### FromFile

A list of `.sh` files to convert into generated executables in the file's executable group.

**Type:** `array` (`string`)




### Ref








### Verb










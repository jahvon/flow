Templates are a powerful feature in flow that allow you to define reusable flowfile and workspace structures. 
This guide will walk you through the process of creating and using templates in flow. It will use a simple example
of create a set of executables for managing a Kubernetes deployment.

## Registering templates

Templates are registered with flow using the [flow template register](../cli/flow_template_register.md) command. This command accepts a
the name to be given to the template and the path to the template file.

```shell
# Register the k8s-deployment template
flow template register k8s-deployment --file /path/to/k8s-deployment.flow.tmpl
```

## Generating scaffolding from a template

Templates are rendered using the [flow template generate](../cli/flow_template_generate.md) command. This command accepts a
the name to be given to the flowfile (if applicable) when rendering its template and several flags to control the rendering process.

```shell
# Run the kes-deployment template generation in the mealie directory of the homelab workspace. 
# The rendered flowfile will be given the name mealie.flow
flow template generate mealie --output mealie --template k8s-deployment --workspace homelab
```

Alternatively, you can reference a flowfile template directly from a file using the `--file` flag.

```shell
flow template generate mealie --output mealie --file /path/to/k8s-deployment.flow.tmpl --workspace homelab
```

## Viewing templates

Individual templates can be viewed using the [flow template view](../cli/flow_template_view.md) command. This command 
accepts either the registered name of the template or the path to the template file.

```shell
# View the k8s-deployment template by its registered name
flow template view --template k8s-deployment
# View the k8s-deployment template by its file path
flow template view --file /path/to/k8s-deployment.flow.tmpl
```

You can also list all registered templates using the [flow template list](../cli/flow_template_list.md) command.

```shell
flow template list
```

## Creating a template file

Templates are defined as YAML files that contain a string template of a [flowfile](executable.md#flowfile), artifacts to be copied,
executables to be run during the templating process, and form fields that can be used throughout the template's fields.

Check out the [template configurations](../types/template.md) for more details on the structure the template file.

### Form inputs

The form section defines the fields that will be prompted to the user when the template is rendered. Each field has a key,
prompt, and optional default value, type, and validation. It's the first step in the template rendering process and is what
provides the data for the [Go text templating](https://pkg.go.dev/text/template) that is used throughout the template file.

For instance, the following form section will prompt the user for the namespace, image, replicas, and 
whether the app should be deployed immediately.

```yaml
form:
  - key: "Namespace"
    prompt: "What namespace should the deployment be created in?"
    default: "apps"
  - key: "Deploy"
    prompt: "Should the app be deployed immediately?"
    type: "confirm" # This will prompt the user with a yes/no question
  - key: "Image"
    prompt: "What image should be used for the deployment?"
    required: true # The template will not render if this field is not provided
  - key: "Replicas"
    prompt: "How many replicas should be created?"
    default: "1"
    validate: "^[0-9]+$" # This will validate that the input is a number
  - key: "Type"
    prompt: "Should the deployment be a Helm chart or a Kubernetes manifest?"
    default: "Helm"
    validate: "^(Helm|K8s)$" # This will validate that the input is either "Helm" or "K8s"
```

### Artifacts

The artifacts section defines the files that will be copied to the output directory when the template is rendered. 
Each artifact has a source, destination, and optional template flag that will render the file as a Go text template.

In the following example, the template will copy the `helm-deploy.sh`, `deploy.sh`, `values.yaml.tmpl`, and `resources.yaml.tmpl` files.
The `values.yaml.tmpl` and `resources.yaml.tmpl` files will be rendered as Go text templates (using the data provided through the form).

The `if` field can be used to conditionally copy the file based on the value of a form field. Below the `helm-deploy.sh` and `values.yaml.tmpl` 
files will only be copied if the `Type` field is set to `Helm`. The `deploy.sh` and `resources.yaml.tmpl` files will only be copied if the 
`Type` field is set to `K8s`.

```yaml
artifacts:
  - srcName: "helm-deploy.sh"
    srcDir: "scripts" # By default, the file will be copied from the template directory. This field can be used to specify a different directory.
    if: "{{ eq .Type 'Helm' }}" # This will only copy the file if the Helm field is true
  - srcName: "deploy.sh"
    srcDir: "scripts"
    if: "{{ eq .Type 'K8s' }}" # This will only copy the file if the K8s field is true
  - srcName: "values.yaml.tmpl"
    asTemplate: true
    dstName: "values.yaml"
    if: "{{ eq .Type 'Helm' }}" # This will only copy the file if the Helm field is true
  - srcName: "resources.yaml.tmpl"
    asTemplate: true
    dstName: "resources.yaml"
    if: "{{ eq .Type 'K8s' }}" # This will only copy the file if the K8s field is true
```

### flowfile template string

The template section defines the string template of the flowfile that will be rendered. The template can be as simple or 
complex as needed and can include Go text templating to reference the form fields provided by the user.

In the following example, the template will create a set flowfile with executables for deploying, restarting, and opening the app.

```yaml
template: |
  tags: [k8s]
  executables:
    - verb: deploy
      name: "{{ .FlowFileName }}"
      exec:
        file: "{{ if eq .Type 'Helm' }}helm-deploy.sh{{ else }}deploy.sh{{ end }}"
        params:
          - envKey: "NAMESPACE"
            text: "{{ .Namespace }}"
          - envKey: "APP_NAME"
            text: "{{ .FlowFileName }}"
    - verb: restart
      name: "{{ .FlowFileName }}"
      exec:
        cmd: "kubectl rollout restart deployment/{{ .FlowFileName }} -n {{ .Namespace }}"
    - verb: open
      name: "{{ .FlowFileName }}"
      launch:
        uri: "https://{{ .FlowFileName }}.my.haus"
```

### Pre- and post- run executables

The preRun and postRun sections define the executables that will be run before and after the template is rendered.
These executables can be used to extend the templating process by running additional commands.

In the following example, the template will run a validation executable before copying artifacts and rendering the flowfile.
Before exiting, it will also run a simple command and either open the flowfile in vscode or deploy the app based on the user's input.

```yaml
preRun:
  - ref: "validate k8s/validation:context" # You can reference other executables that you have on your system
    args: ["homelab"]
    if: "{{ .Deploy }}"
postRun:
  - cmd: "echo 'Rendered {{ if .Helm }}Helm values{{ else }}k8s manifest{{end}}'; ls -al"
  - ref: "edit vscode"
    args: ["{{ .FlowFilePath }}"]
    if: "{{ not .Deploy }}"
  - ref: "deploy {{ .FlowFileName }}"
    if: "{{ .Deploy }}"
```

**Note**: preRun executables are run from the template directory, while postRun executables are run from the output directory.

### Full template example

Bringing it all together, the following is a full example of the k8s deployment template. It's not required to have all the sections
in a template, but it's a good starting point for creating your own templates.

```yaml
form:
  - key: "Namespace"
    prompt: "What namespace should the deployment be created in?"
    default: "apps"
  - key: "Deploy"
    prompt: "Should the app be deployed immediately?"
    type: "confirm"
  - key: "Image"
    prompt: "What image should be used for the deployment?"
  - key: "Replicas"
    prompt: "How many replicas should be created?"
    default: "1"
    validate: "^[0-9]+$"
  - key: "Type"
    prompt: "Should the deployment be a Helm chart or a Kubernetes manifest?"
    default: "Helm"
    validate: "^(Helm|K8s)$"
artifacts:
  - srcName: "helm-deploy.sh"
    srcDir: "scripts"
    if: "{{ eq .Type 'Helm' }}"
  - srcName: "deploy.sh"
    srcDir: "scripts"
    if: "{{ eq .Type 'K8s' }}"
  - srcName: "values.yaml.tmpl"
    asTemplate: true
    dstName: "values.yaml"
    if: "{{ eq .Type 'Helm' }}"
  - srcName: "resources.yaml.tmpl"
    asTemplate: true
    dstName: "resources.yaml"
    if: "{{ eq .Type 'K8s' }}"
preRun:
  - ref: "validate k8s/validation:context"
    args: ["homelab"]
    if: "{{ .Deploy }}"
postRun:
  - cmd: "echo 'Rendered {{ if .Helm }}Helm values{{ else }}k8s manifest{{end}}'; ls -al"
  - ref: "edit vscode"
    args: ["{{ .FlowFilePath }}"]
    if: "{{ not .Deploy }}"
  - ref: "deploy {{ .FlowFileName }}"
    if: "{{ .Deploy }}"
template: |
  tags: [k8s]
  executables:
    - verb: deploy
      name: "{{ .FlowFileName }}"
      exec:
        file: "{{ if eq .Type 'Helm' }}helm-deploy.sh{{ else }}deploy.sh{{ end }}"
        params:
          - envKey: "NAMESPACE"
            text: "{{ .Namespace }}"
          - envKey: "APP_NAME"
            text: "{{ .FlowFileName }}"
    - verb: restart
      name: "{{ .FlowFileName }}"
      exec:
        cmd: "kubectl rollout restart deployment/{{ .FlowFileName }} -n {{ .Namespace }}"
    - verb: open
      name: "{{ .FlowFileName }}"
      launch:
        uri: "https://{{ .FlowFileName }}.my.haus"
``` 

### Template helpers

Templates can use the [Sprig functions](https://masterminds.github.io/sprig/) to manipulate data during the rendering process.
Additionally, the following keys are available to the template:

- `FlowFileName`: The name of the flowfile being rendered. This is the argument passed into the `generate` command.
- `FlowFilePath`: The output path to the flowfile being rendered.
- `FlowWorkspace`: The name of the workspace that the template is rendering into.
- `FlowWorkspacePath`: The path to the workspace root directory.

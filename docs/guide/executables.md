# Executables

Executables are the building blocks of flow automation. They can be simple commands, complex multi-step workflows, HTTP requests, or even GUI applications. 
This guide covers all executable types and configuration options.

## Finding Executables

Use the `flow browse` command to discover executables across your workspaces:

```shell
flow browse            # Interactive multi-pane browser
flow browse --list     # Simple list view
flow browse VERB ID    # View specific executable details
```

Filter executables by workspace, namespace, verb, or tag:

```shell
flow browse --workspace api --namespace v1 --verb deploy --tag production
flow browse --all --filter "database"  # Search names and descriptions
```

## Executable Configuration

### Basic Structure <!-- {docsify-ignore} -->

Every executable needs a verb and optionally a name:

```yaml
executables:
  - verb: run
    name: my-task
    description: "Does something useful"
    tags: [development, automation]
    aliases: [task, job]
    timeout: 5m
    visibility: public
    exec:
      cmd: echo "Hello, world!"
```

### Common Fields <!-- {docsify-ignore} -->

- **verb**: Action type (run, build, test, deploy, etc.)
- **name**: Unique identifier within the namespace
- **description**: Markdown documentation for the executable
- **tags**: Labels for categorization and filtering
- **aliases**: Alternative names for the executable
- **timeout**: Maximum execution time (e.g., 30s, 5m, 1h)
- **visibility**: Access control (public, private, internal, hidden)

### Visibility Levels <!-- {docsify-ignore} -->

- **public**: Available from any workspace
- **private**: Only available within the same workspace but shown in browse lists (default)
- **internal**: Available within workspace but hidden from browse lists
- **hidden**: Cannot be run or listed

## Environment Variables

Customize executable behavior with environment variables using `params` or `args`.

### Parameters (`params`) <!-- {docsify-ignore} -->

Set environment variables from various sources:

```yaml
executables:
  - verb: deploy
    name: app
    exec:
      file: deploy.sh
      params:
        # From secrets
        - secretRef: api-token
          envKey: API_TOKEN
        - secretRef: production/database-url
          envKey: DATABASE_URL
        
        # Interactive prompts
        - prompt: "Which environment?"
          envKey: ENVIRONMENT
        
        # Static values
        - text: "production"
          envKey: DEPLOY_ENV
```

**Parameter types:**
- `secretRef`: Reference to vault secret
- `prompt`: Interactive user input
- `text`: Static value

### Arguments (`args`) <!-- {docsify-ignore} -->

Handle command-line arguments:

```yaml
executables:
  - verb: build
    name: container
    exec:
      file: build.sh
      args:
        # Positional argument
        - pos: 1
          envKey: IMAGE_TAG
          required: true
        
        # Flag arguments
        - flag: publish
          envKey: PUBLISH
          type: bool
          default: false
        
        - flag: registry
          envKey: REGISTRY
          default: "docker.io"
```

**Run with arguments:**
```shell
flow build container v1.2.3 publish=true registry=my-registry.com
```

**Argument types:**
- `pos`: Positional argument (by position number, starting from 1)
- `flag`: Named flag argument
- `type`: Validation type (string, int, float, bool)

### Command-Line Overrides <!-- {docsify-ignore} -->

Override any environment variable with `--param`:

```shell
flow deploy app --param API_TOKEN=override --param ENVIRONMENT=staging
```

## Working Directories

Control where executables run with the `dir` field:

```yaml
executables:
  - verb: build
    name: frontend
    exec:
      cmd: npm run build
      dir: "./frontend"  # Relative to flowfile
  
  - verb: clean
    name: downloads
    exec:
      cmd: rm -rf downloads/*
      dir: "~/Downloads"  # User home directory
  
  - verb: deploy
    name: from-root
    exec:
      cmd: kubectl apply -f k8s/
      dir: "//"  # Workspace root
  
  - verb: test
    name: isolated
    exec:
      cmd: |
        echo "Running in temporary directory"
        ls -la
      dir: "f:tmp"  # Temporary directory (auto-cleaned)
```

**Directory prefixes:**
- `//`: Workspace root directory
- `~/`: User home directory
- `./`: Current working directory
- `f:tmp`: Temporary directory (auto-cleaned)
- `$VAR`: Environment variable expansion

## Executable Types <!-- {docsify-ignore} -->

### exec - Shell Commands

Run commands or scripts directly:

```yaml
executables:
  - verb: build
    name: app
    exec:
      cmd: npm run build && npm test
  
  - verb: deploy
    name: app
    exec:
      file: deploy.sh
      logMode: json  # text, logfmt, json, or hidden
```

**Options:**
- `cmd`: Inline command to run
- `file`: Script file to execute
- `logMode`: How to format command output

### serial - Sequential Execution

Run multiple steps in order:

```yaml
executables:
  - verb: deploy
    name: full-stack
    serial:
      failFast: true  # Stop on first failure
      execs:
        - cmd: docker build -t api .
        - cmd: docker build -t web ./frontend
        - ref: test api
        - cmd: kubectl apply -f k8s/
          retries: 3
        - cmd: kubectl rollout status deployment/api
          reviewRequired: true  # Pause for user confirmation
```

The [executable environment variables](#environment-variables) and [executable directory](#working-directories)
of the parent executable are inherited by the child executables.

**Options:**
- `failFast`: Stop execution on first failure (default: true)
- `retries`: Number of times to retry failed steps
- `reviewRequired`: Pause for user confirmation

### parallel - Concurrent Execution

Run multiple steps simultaneously:

```yaml
executables:
  - verb: test
    name: all-suites
    parallel:
      maxThreads: 4  # Limit concurrent operations
      failFast: false  # Run all tests even if some fail
      execs:
        - cmd: npm run test:unit
        - cmd: npm run test:integration
        - cmd: npm run test:e2e
        - ref: lint code
          retries: 1
```

The [executable environment variables](#environment-variables) and [executable directory](#working-directories)
of the parent executable are inherited by the child executables.

**Options:**
- `maxThreads`: Maximum concurrent operations (default: 5)
- `failFast`: Stop all operations on first failure (default: true)
- `retries`: Number of times to retry failed operations

### launch - Open Applications

Open files, URLs, or applications:

```yaml
executables:
  - verb: open
    name: workspace
    launch:
      uri: "$FLOW_WORKSPACE_PATH"
      app: "Visual Studio Code"
  
  - verb: open
    name: docs
    launch:
      uri: "https://flowexec.io"
  
  - verb: open
    name: note
    launch:
      uri: "./note.md"
      app: "Obsidian"
```

**Options:**
- `uri`: File path or URL to open (required)
- `app`: Specific application to use

### request - HTTP Requests

Make HTTP requests to APIs:

```yaml
executables:
  - verb: deploy
    name: webhook
    request:
      method: POST
      url: "https://api.example.com/deploy"
      headers:
        Authorization: "Bearer $API_TOKEN"
        Content-Type: "application/json"
      body: |
        {
          "environment": "$ENVIRONMENT",
          "version": "$VERSION"
        }
      timeout: 30s
      validStatusCodes: [200, 201]
      logResponse: true
      transformResponse: |
        "Deployment " + fromJSON(data)["status"]
      responseFile:
        filename: "deploy-response.json"
```

**Options:**
- `method`: HTTP method (GET, POST, PUT, PATCH, DELETE)
- `url`: Request URL (required)
- `headers`: Custom headers
- `body`: Request body
- `timeout`: Request timeout
- `validStatusCodes`: Acceptable status codes
- `logResponse`: Log response body
- `transformResponse`: Transform response with Expr
- `responseFile`: Save response to file

### render - Dynamic Documentation

Generate and display markdown with templates:

```yaml
executables:
  - verb: show
    name: status
    render:
      templateFile: "status-template.md"
      templateDataFile: "status-data.json"
```

**Template file example:**
```markdown
# System Status

Current time: {{ .timestamp }}

## Services
{{- range .services }}
- **{{ .name }}**: {{ .status }}
{{- end }}

## Metrics
- CPU: {{ .cpu }}%
- Memory: {{ .memory }}%
```

**Options:**
- `templateFile`: Markdown template file (required)
- `templateDataFile`: JSON/YAML data file

## Generated Executables

Generate executables from shell scripts with special comments:

```yaml
# In flowfile
fromFile:
  - "scripts/deploy.sh"
  - "scripts/backup.sh"
```

```bash
#!/bin/bash
# scripts/deploy.sh

# f:name=production f:verb=deploy
# f:description="Deploy to production environment"
# f:tags=production,critical
# f:aliases=prod-deploy
# f:visibility=internal
# f:timeout=10m

echo "Deploying to production..."
kubectl apply -f k8s/
```

**Supported comment keys:**
- `name`, `verb`, `description`, `tags`, `aliases`, `visibility`, `timeout`

**Multi-line descriptions:**
```bash
# f:name=backup f:verb=run
# <f|description>
# Creates a backup of the database
# and uploads it to S3 storage
# <f|description>
```

## Executable References

Reference other executables to build modular workflows:

```yaml
executables:
  # Reusable components
  - verb: build
    name: api
    exec:
      cmd: docker build -t api .
  
  - verb: test
    name: api
    exec:
      cmd: npm test
  
  # Composite workflows
  - verb: deploy
    name: full
    serial:
      execs:
        - ref: build api
        - ref: test api
        - cmd: kubectl apply -f api.yaml
  
  # Cross-workspace references (requires public visibility)
  - verb: deploy
    name: with-monitoring
    serial:
      execs:
        - ref: deploy full
        - ref: trigger monitoring/slack:deployment-complete
```

**Reference formats:**
- `ref: build api` - Current workspace/namespace
- `ref: build workspace/namespace:api` - Full reference
- `ref: build workspace/api` - Specific workspace
- `ref: build namespace:api` - Specific namespace

**Cross-workspace requirements:**
- Referenced executables must have `visibility: public`
- Private, internal, and hidden executables cannot be cross-referenced

## What's Next? <!-- {docsify-ignore} -->

Now that you understand all executable types and options:

- **Build complex workflows** → [Advanced workflows](advanced.md)
- **Secure your automation** → [Working with secrets](secrets.md)
- **Generate project templates** → [Templates & code generation](templating.md)
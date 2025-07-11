# Integrations

flow integrates with popular CI/CD platforms and containerized environments to bring your automation anywhere.

## GitHub Actions

Execute flow workflows directly in your GitHub Actions pipelines with the official action.

```yaml
name: Build and Deploy
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: flowexec/action@v1
        with:
          executable: 'build app'
```

> **Complete documentation**: Visit the [Flow Execute Action](https://github.com/marketplace/actions/flow-execute) on GitHub Marketplace.

## Docker

Run flow in containerized environments for CI/CD pipelines or isolated execution.

### Basic Usage <!-- {docsify-ignore} -->

```shell
# Run with default workspace
docker run -it --rm ghcr.io/flowexec/flow

# Execute specific executable
docker run -it --rm ghcr.io/flowexec/flow validate
```

**Environment Variables**
- `REPO`: Repository URL to clone (defaults to flow's repo)
- `BRANCH`: Git branch to checkout (optional)
- `WORKSPACE`: Workspace name to use (defaults to "flow")


### Workspace from Git <!-- {docsify-ignore} -->

Automatically clone and configure a workspace:

```shell
docker run -it --rm \
  -e REPO=https://github.com/your-org/your-workspace \
  -e BRANCH=main \
  -e WORKSPACE=my-workspace \
  ghcr.io/flowexec/flow exec "deploy app"
```

### Local Workspace <!-- {docsify-ignore} -->

Mount your local workspace:

```shell
docker run -it --rm \
  -v $(pwd):/workspaces/my-workspace \
  -w /workspaces/my-workspace \
  -e WORKSPACE=my-workspace \
  ghcr.io/flowexec/flow exec "build app"
```

### In CI/CD Pipelines <!-- {docsify-ignore} -->

Any CI/CD platform that supports Docker can run flow. The key is:

1. **Use the Docker image**: `ghcr.io/flowexec/flow`
2. **Set environment variables**: `REPO`, `WORKSPACE`, `BRANCH` as needed
3. **Execute your flow commands**: `flow exec "your-executable"`

> **Note**: While this should work, the Docker integration hasn't been extensively tested. If you try flow with other CI/CD platforms, we'd love to hear about your experience!

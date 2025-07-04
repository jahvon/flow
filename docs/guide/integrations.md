# Integrations

> **Note:** Integrations are still considered experimental and may change in future releases. Please provide feedback on your experience.

## GitHub Actions

flow provides seamless integration with GitHub Actions through the [flow-execute](https://github.com/marketplace/actions/flow-execute) marketplace action. This enables you to execute 
your flow executables directly within your CI/CD pipelines.

#### Quick Start

```yaml
- uses: jahvon/flow-action@v1.0.0-beta1
  with:
    executable: 'build app'
```

For complete documentation, configuration options, and examples, visit the [Flow Execute Action](https://github.com/marketplace/actions/flow-execute) on the GitHub Marketplace.

## Docker

flow can run in containerized environments, making it useful for CI/CD pipelines or isolated execution environments.

#### Basic Usage

```shell
docker run -it --rm ghcr.io/jahvon/flow
```

This runs the container with the default flow workspace and shows the version.

#### With Workspace Configuration

You can automatically clone and configure a flow workspace by setting environment variables:

```shell
docker run -it --rm \
  -e REPO=https://github.com/your-org/your-workspace \
  -e BRANCH=main \
  -e WORKSPACE=my-workspace \
  ghcr.io/jahvon/flow exec "your-executable"
```

#### Environment Variables

- `REPO`: Repository URL to clone (defaults to flow's own repo)
- `BRANCH`: Git branch to checkout (optional, defaults to default branch)
- `WORKSPACE`: Workspace name to use (defaults to "flow")

#### Volume Mounting

For persistent data or to use local workspaces:

```shell
docker run -it --rm \
  -v $(pwd):/workspaces/my-workspace \
  -w /workspaces/my-workspace \
  -e WORKSPACE=my-workspace \
  ghcr.io/jahvon/flow "exec your-executable"
```

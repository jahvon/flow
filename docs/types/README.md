# flow YAML Configurations

- [FlowFile](flowfile.md)
- [Template](template.md)
- [Workspace](workspace.md)
- [Config](config.md)

## Schemas

All flow configuration files also have YAML schemas available for use in IDEs with the Language Server Protocol (LSP) to perform intelligent suggestions.
You can add the following comment to the top of your flow files to enable this:

```yaml
# yaml-language-server: $schema=https://flowexec.io/schemas/flowfile_schema.json
```

See the [schemas directory on GitHub](https://github.com/jahvon/flow/tree/main/docs/schemas) for all available schemas.

Note: If using the flow file schema, you will need to make sure your IDE is configured to treat `*.flow`, `*.flow.yaml`, and `*.flow.yml` files as YAML files.

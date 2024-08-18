## flow init executables

Add rendered executables from an executable definition template to a workspace

### Synopsis

Add rendered executables from an executable definition template to a workspace.

The WORKSPACE_NAME is the name of the workspace to initialize the executables in.
The DEFINITION_NAME is the name of the definition to use in rendering the template.
This name will become the name of the file containing the copied executable definition.

One one of -f or -t must be provided and must point to a valid executable definition template.
The -p flag can be used to specify a sub-path within the workspace to create the executable definition and its artifacts.

```
flow init executables FLOWFILE_NAME [-w WORKSPACE ] [-o OUTPUT_DIR] [-f FILE | -t TEMPLATE] [flags]
```

### Options

```
  -f, --file string                  Path to the template file. It must be a valid flow file template.
  -h, --help                         help for executables
  -o, --output string                Output directory (within the workspace) to create the flow file and its artifacts. If the directory does not exist, it will be created.
  -t, --template flow set template   Registered template name. Templates can be registered in the flow configuration file or with flow set template.
  -w, --workspace string             Workspace to create the flow file and its artifacts. Defaults to the current workspace.
```

### Options inherited from parent commands

```
  -x, --non-interactive   Disable displaying flow output via terminal UI rendering. This is only needed if the interactive output is enabled by default in flow's configuration.
      --sync              Sync flow cache and workspaces
      --verbosity int     Log verbosity level (-1 to 1)
```

### SEE ALSO

* [flow init](flow_init.md)	 - Initialize or restore the flow application state.


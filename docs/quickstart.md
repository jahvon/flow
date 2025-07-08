> [!NOTE]
> Before getting started, install the latest `flow` version using one of the methods described in the 
> [installation guide](installation.md).

This guide will walk you through the basic steps of creating and running an executable with `flow`.

**Create a new workspace**

A workspace can be created anywhere on your system but must be registered in order to have executables discovered by flow.

To create a new workspace, run the following command in the directory where you want the workspace to be created:

```shell
flow workspace add my-workspace . --set
```

You can replace `my-workspace` with any name you want to give your workspace and `.` with the path to the root directory 
of the new workspace.

This command will register the workspace and create a `flow.yaml` file in the root directory. This file contains the new 
workspace's configurations. For more information on workspaces see the [workspace guide](guide/workspace.md).

**Create an executable**

Executables are the core of flow. Each executable is driven by its definition within a flow file (`*.flow`, `*.flow.yaml`, or `*.flow.yml`).
There are several types of executables that can be defined. For more information on executables and the flow file, see the [executable guide](guide/executable.md).

To get started, create a new flow file in the workspace directory.
    
```shell
touch executables.flow
```

Open the file in your favorite text editor and add the following content:

```yaml
executables:
  - verb: run
    name: my-task
    exec:
      params:
      - prompt: What is your favorite color?
        envKey: COLOR
      cmd: echo "Your favorite color is $COLOR"
```

This flow file defines a single executable named `my-task`. When run, it will prompt the user for their favorite color 
and then echo that color back to the console.

**Add another workspace**

If you're new to flow, try adding the `flow` workspace and explore the executables it contains.

```shell
git clone github.com/jahvon/flow
flow workspace add flow flow
```

When you have multiple workspaces, you can switch between them using a command like the following:

```shell
flow workspace set flow
```

**Running an executable**

Whenever you create, move, or delete executables and flow files, you will need to update the index of executables before running them.

```shell
flow sync
```

The main command for running executables is `flow exec`. This command will execute the workflow with the provided
executable ID. `exec` can be replaced with any verb known to flow but should match the verb defined in the flow file
configurations or an alias of that verb.

In our case, we will use the `run` verb:

```shell
flow run my-task
```

> [!TIP]
> You can also run the executable by its full name, `run my-workspace/my-task` or an alias if one is defined.

Try adding more executables to the workspace! You can create multiple flow files anywhere in the workspace. As you add more
executables, try viewing them from the interactive UI:

```shell
flow browse
```

When in the library, you can press the <kbd>R</kbd> key on a selected executable to run it.

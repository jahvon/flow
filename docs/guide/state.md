## State Management

flow provides several mechanisms for managing state across [executable](executable.md) runs. This guide explores how to use the data store
and temporary directories to maintain state, as well as how to make execution decisions based on that state.

### Data Store

The data store is a key-value store that persists data across executions. It provides a simple way to share information
between executables and maintain state between runs.

#### Store Persistence Scope

Values in the data store have different persistence scopes depending on where they are set:

- Values set outside an executable (using the CLI directly) persist across all executions until explicitly cleared
- Values set within an executable persist only across that executable's sub-executables (both serial and parallel)
- All values set within an executable are automatically cleared when the parent executable completes

#### Managing Store Data

The data store can be managed at a global level in the CLI and within an executable's script. Here are the key operations:

**Setting Values**

```shell
# Direct CLI usage
flow store set KEY VALUE

# for example:
flow store set my-key "my value"
# or pipe a value from a command
echo "my value" | flow store set my-key
```

**Getting Values**

```shell
# Direct CLI usage
flow store get KEY

# for example:
value=$(flow store get my-key)
```

**Clearing Values**

```shell
# Clear all values
flow store clear
# Clear with full flag to remove all stored data
flow store clear --full
```

#### Using the Store in Executables

Here's an example of how to use the data store within an executable:

```yaml
executables:
  - verb: run
    name: data-store-demo
    serial:
      params:
      execs:
        # Set some values in the store
        - cmd: |
            flow store set user-preference dark-mode
            flow store set last-run "$(date)"
        # Use those values in a subsequent step
        - cmd: |
            preference=$(flow store get user-preference)
            echo "User preference is: $preference"
            echo "Last run: $(flow store get last-run)"
```

#### Store-Based Conditional Execution

The data store's contents can be accessed in executable `if` conditions using the `data` context variable. This allows for
dynamic execution paths based on stored values:

```yaml
executables:
  - verb: run
    name: conditional-demo
    serial:
      execs:
        - cmd: flow store set feature-enabled true
        # This will execute because feature-enabled is set to "true"
        - if: data["feature-enabled"] == "true"
          cmd: echo "Feature is enabled"
        # This will not execute because test-key is not set
        - if: len(data["test-key"]) > 0
          cmd: echo "Test key exists"
```

See the [Conditional Execution](conditional.md) guide for more examples of using conditions in Flow.

### Temporary Directories

Flow provides a special directory reference `f:tmp` that creates an isolated temporary directory for an executable. This
directory is automatically cleaned up when the executable completes.

#### Using Temporary Directories

To use a temporary directory, set the `dir` field in your executable configuration:

```yaml
executables:
  - verb: build
    name: temp-workspace
    exec:
      dir: f:tmp  # Creates and uses a temporary directory
      cmd: |
        # All commands run in an isolated temp directory
        git clone https://github.com/user/repo .
        make build
```

#### Sharing Temporary Files

While temporary directories are isolated, you can share files between steps in a serial or parallel executable by using
the same temporary directory:

```yaml
executables:
  - verb: process
    name: shared-temp
    serial:
      dir: f:tmp  # All sub-executables share this temp directory
      execs:
        - cmd: echo "Step 1" > output.txt
        - cmd: cat output.txt && echo "Step 2" >> output.txt
        - cmd: cat output.txt
```

### Combining State Management Approaches

The data store and temporary directories can be used together for more complex state management:

```yaml
executables:
  - verb: build
    name: complex-state
    serial:
      dir: f:tmp
      execs:
        # Generate a build ID and store it
        - cmd: |
            build_id=$(date +%Y%m%d_%H%M%S)
            flow store set current-build $build_id
        # Use the stored build ID for conditional execution
        - if: len(data["current-build"]) > 0
          cmd: |
            echo "Building artifacts for ${build_id}"
            make build
        # Clean up based on stored state
        - if: data["cleanup-enabled"] == "true"
          cmd: make clean
```

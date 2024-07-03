package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jahvon/tuikit/components"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/config"
	argUtils "github.com/jahvon/flow/config/args"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/exec"
	"github.com/jahvon/flow/internal/runner/launch"
	"github.com/jahvon/flow/internal/runner/parallel"
	"github.com/jahvon/flow/internal/runner/render"
	"github.com/jahvon/flow/internal/runner/request"
	"github.com/jahvon/flow/internal/runner/serial"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterExecCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "exec EXECUTABLE_ID [args...]",
		Aliases: config.SortedValidVerbs(),
		Short:   "Execute a flow by ID.",
		Long: execDocumentation + "\n\n" + execExamples + "\n\n" +
			"See " + io.ConfigDocsURL("executables", "Verb") + "for more information on executable verbs." +
			"See " + io.ConfigDocsURL("executables", "Ref") + "for more information on executable IDs.",
		Args: cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			execList, err := ctx.ExecutableCache.GetExecutableList(ctx.Logger)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			execIDs := make([]string, 0, len(execList))
			for _, e := range execList {
				execIDs = append(execIDs, e.ID())
			}
			return execIDs, cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			execPreRun(ctx, cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			verbStr := cmd.CalledAs()
			verb := config.Verb(verbStr)
			execFunc(ctx, cmd, verb, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func execPreRun(ctx *context.Context, cmd *cobra.Command, _ []string) {
	runner.RegisterRunner(exec.NewRunner())
	runner.RegisterRunner(launch.NewRunner())
	runner.RegisterRunner(request.NewRunner())
	runner.RegisterRunner(render.NewRunner())
	runner.RegisterRunner(serial.NewRunner())
	runner.RegisterRunner(parallel.NewRunner())
	interactive.InitInteractiveCommand(ctx, cmd)
}

//nolint:gocognit
func execFunc(ctx *context.Context, cmd *cobra.Command, verb config.Verb, args []string) {
	logger := ctx.Logger
	if err := verb.Validate(); err != nil {
		logger.FatalErr(err)
	}

	idArg := args[0]
	ref := context.ExpandRef(ctx, config.NewRef(idArg, verb))
	executable, err := ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
		logger.Debugf("Executable %s not found in cache, syncing cache", ref)
		if err := ctx.ExecutableCache.Update(logger); err != nil {
			logger.FatalErr(err)
		}
		executable, err = ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	}
	if err != nil {
		logger.FatalErr(err)
	}

	if err := executable.Validate(); err != nil {
		logger.FatalErr(err)
	}

	if !executable.IsExecutableFromWorkspace(ctx.CurrentWorkspace.AssignedName()) {
		logger.FatalErr(fmt.Errorf(
			"executable '%s' cannot be executed from workspace %s",
			ref,
			ctx.UserConfig.CurrentWorkspace,
		))
	}

	execArgs := args[1:]
	envMap, err := argUtils.ProcessArgs(executable, execArgs)
	if err != nil {
		logger.FatalErr(err)
	}
	if envMap == nil {
		envMap = make(map[string]string)
	}

	setAuthEnv(ctx, executable)
	textInputs := pendingTextInputs(ctx, executable)
	if len(textInputs) > 0 {
		inputs, err := components.ProcessInputs(io.Theme(), textInputs...)
		if err != nil {
			logger.FatalErr(err)
		}
		for _, input := range inputs {
			envMap[input.Key] = input.Value()
		}
	}
	startTime := time.Now()
	if err := runner.Exec(ctx, executable, envMap); err != nil {
		logger.FatalErr(err)
	}
	dur := time.Since(startTime)
	logger.Infox(fmt.Sprintf("%s flow completed", ref), "Elapsed", dur.Round(time.Millisecond))
	if interactive.UIEnabled(ctx, cmd) {
		if dur > 1*time.Minute && ctx.UserConfig.Interactive.SoundOnCompletion {
			_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		}
		if dur > 1*time.Minute && ctx.UserConfig.Interactive.NotifyOnCompletion {
			_ = beeep.Notify("Flow", "Flow completed", "")
		}
	}
}

func runByRef(ctx *context.Context, cmd *cobra.Command, argsStr string) error {
	s := strings.Split(argsStr, " ")
	if len(s) != 2 {
		return fmt.Errorf("invalid reference string %s", argsStr)
	}
	verbStr := s[0]
	verb := config.Verb(verbStr)
	id := s[1]

	cmds := cmd.Root().Commands()
	var execCmd *cobra.Command
	for _, c := range cmds {
		if c.Name() == "exec" {
			execCmd = c
			break
		}
	}

	if execCmd == nil {
		return errors.New("exec command not found")
	}
	execCmd.SetArgs([]string{verbStr, id})
	execCmd.SetOut(ctx.StdOut())
	execCmd.SetErr(ctx.StdOut())
	execCmd.SetIn(ctx.StdIn())
	execPreRun(ctx, execCmd, []string{id})
	execFunc(ctx, execCmd, verb, []string{id})
	ctx.CancelFunc()
	return nil
}

func setAuthEnv(ctx *context.Context, executable *config.Executable) {
	if authRequired(ctx, executable) {
		resp, err := components.ProcessInputs(
			io.Theme(),
			&components.TextInput{
				Key:    vault.EncryptionKeyEnvVar,
				Prompt: "Enter vault encryption key",
				Hidden: true,
			})
		if err != nil {
			ctx.Logger.FatalErr(err)
		}
		val := resp.ValueMap()[vault.EncryptionKeyEnvVar]
		if val == "" {
			ctx.Logger.FatalErr(fmt.Errorf("vault encryption key required"))
		}
		if err := os.Setenv(vault.EncryptionKeyEnvVar, val); err != nil {
			ctx.Logger.FatalErr(fmt.Errorf("failed to set vault encryption key\n%w", err))
		}
	}
}

//nolint:gocognit
func authRequired(ctx *context.Context, rootExec *config.Executable) bool {
	if rootExec.Type == nil || os.Getenv(vault.EncryptionKeyEnvVar) != "" {
		return false
	}
	switch {
	case rootExec.Type.Exec != nil:
		for _, param := range rootExec.Type.Exec.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Type.Launch != nil:
		for _, param := range rootExec.Type.Launch.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Type.Request != nil:
		for _, param := range rootExec.Type.Request.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Type.Render != nil:
		for _, param := range rootExec.Type.Render.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Type.Serial != nil:
		for _, param := range rootExec.Type.Serial.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
		for _, child := range rootExec.Type.Serial.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			if authRequired(ctx, childExec) {
				return true
			}
		}
	case rootExec.Type.Parallel != nil:
		for _, param := range rootExec.Type.Parallel.Parameters {
			if param.SecretRef != "" {
				return true
			}
		}
		for _, child := range rootExec.Type.Parallel.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			if authRequired(ctx, childExec) {
				return true
			}
		}
	}
	return false
}

//nolint:gocognit
func pendingTextInputs(ctx *context.Context, rootExec *config.Executable) []*components.TextInput {
	pending := make([]*components.TextInput, 0)
	if rootExec.Type == nil {
		return nil
	}
	switch {
	case rootExec.Type.Exec != nil:
		for _, param := range rootExec.Type.Exec.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Launch != nil:
		for _, param := range rootExec.Type.Launch.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Request != nil:
		for _, param := range rootExec.Type.Request.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Render != nil:
		for _, param := range rootExec.Type.Render.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Serial != nil:
		for _, param := range rootExec.Type.Serial.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
		for _, child := range rootExec.Type.Serial.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			childPending := pendingTextInputs(ctx, childExec)
			pending = append(pending, childPending...)
		}
	case rootExec.Type.Parallel != nil:
		for _, param := range rootExec.Type.Parallel.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &components.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
		for _, child := range rootExec.Type.Parallel.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			childPending := pendingTextInputs(ctx, childExec)
			pending = append(pending, childPending...)
		}
	}
	return pending
}

var (
	//nolint:lll
	execDocumentation = `
Execute a flow where <executable-id> is the target executable's ID in the form of 'ws/ns:name'.
The flow subcommand used should match the target executable's verb or one of its aliases.

If the target executable accept arguments, they can be passed in the form of flag or positional arguments.
Flag arguments are specified with the format 'flag=value' and positional arguments are specified as values without any prefix.
`
	execExamples = `
# Execute the 'build' flow in the current workspace and namespace
flow exec build
flow run build # Equivalent to the above since 'run' is an alias for the 'exec' verb

# Execute the 'docs' flow with the 'show' verb in the current workspace and namespace
flow show docs

# Execute the 'build' flow in the 'ws' workspace and 'ns' namespace
flow exec ws/ns:build

# Execute the 'build' flow in the 'ws' workspace and 'ns' namespace with flag and positional arguments
flow exec ws/ns:build flag1=value1 flag2=value2 value3 value4
`
)

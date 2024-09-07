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
	argUtils "github.com/jahvon/flow/internal/utils/args"
	"github.com/jahvon/flow/internal/vault"
	"github.com/jahvon/flow/types/executable"
)

func RegisterExecCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "exec EXECUTABLE_ID [args...]",
		Aliases: executable.SortedValidVerbs(),
		Short:   "Execute a flow by ID.",
		Long: execDocumentation + fmt.Sprintf(
			"\n\n%s\n\nSee %s for more information on executable verbs."+
				"See %s for more information on executable IDs.",
			execExamples, io.TypesDocsURL("flowfile", "ExecutableVerb"),
			io.TypesDocsURL("flowfile", "ExecutableRef"),
		),
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
			verb := executable.Verb(verbStr)
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
	SetLoadingView(ctx, cmd)
}

//nolint:gocognit
func execFunc(ctx *context.Context, cmd *cobra.Command, verb executable.Verb, args []string) {
	logger := ctx.Logger
	if err := verb.Validate(); err != nil {
		logger.FatalErr(err)
	}

	idArg := args[0]
	ref := context.ExpandRef(ctx, executable.NewRef(idArg, verb))
	e, err := ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
		logger.Debugf("Executable %s not found in cache, syncing cache", ref)
		if err := ctx.ExecutableCache.Update(logger); err != nil {
			logger.FatalErr(err)
		}
		e, err = ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	}
	if err != nil {
		logger.FatalErr(err)
	}

	if err := e.Validate(); err != nil {
		logger.FatalErr(err)
	}

	if !e.IsExecutableFromWorkspace(ctx.CurrentWorkspace.AssignedName()) {
		logger.FatalErr(fmt.Errorf(
			"e '%s' cannot be executed from workspace %s",
			ref,
			ctx.Config.CurrentWorkspace,
		))
	}

	execArgs := args[1:]
	envMap, err := argUtils.ProcessArgs(e, execArgs)
	if err != nil {
		logger.FatalErr(err)
	}
	if envMap == nil {
		envMap = make(map[string]string)
	}

	setAuthEnv(ctx, e)
	textInputs := pendingFormFields(ctx, e)
	if len(textInputs) > 0 {
		form, err := components.NewForm(io.Theme(), ctx.StdIn(), ctx.StdOut(), textInputs...)
		if err != nil {
			logger.FatalErr(err)
		}
		if err := ctx.SetView(form); err != nil {
			logger.FatalErr(err)
		}
		for key, val := range form.ValueMap() {
			envMap[key] = fmt.Sprintf("%v", val)
		}
	}
	startTime := time.Now()
	if err := runner.Exec(ctx, e, envMap); err != nil {
		logger.FatalErr(err)
	}
	dur := time.Since(startTime)
	logger.Infox(fmt.Sprintf("%s flow completed", ref), "Elapsed", dur.Round(time.Millisecond))
	if UIEnabled(ctx, cmd) {
		if dur > 1*time.Minute && ctx.Config.SendSoundNotification() {
			_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		}
		if dur > 1*time.Minute && ctx.Config.SendTextNotification() {
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
	verb := executable.Verb(verbStr)
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

func setAuthEnv(ctx *context.Context, executable *executable.Executable) {
	if authRequired(ctx, executable) {
		form, err := components.NewForm(
			io.Theme(),
			ctx.StdIn(),
			ctx.StdOut(),
			&components.FormField{
				Key:   vault.EncryptionKeyEnvVar,
				Title: "Enter vault encryption key",
				Type:  components.PromptTypeMasked,
			})
		if err != nil {
			ctx.Logger.FatalErr(err)
		}
		if err := ctx.SetView(form); err != nil {
			ctx.Logger.FatalErr(err)
		}
		val := form.FindByKey(vault.EncryptionKeyEnvVar).Value()
		if val == "" {
			ctx.Logger.FatalErr(fmt.Errorf("vault encryption key required"))
		}
		if err := os.Setenv(vault.EncryptionKeyEnvVar, val); err != nil {
			ctx.Logger.FatalErr(fmt.Errorf("failed to set vault encryption key\n%w", err))
		}
	}
}

//nolint:gocognit
func authRequired(ctx *context.Context, rootExec *executable.Executable) bool {
	if os.Getenv(vault.EncryptionKeyEnvVar) != "" {
		return false
	}
	switch {
	case rootExec.Exec != nil:
		for _, param := range rootExec.Exec.Params {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Launch != nil:
		for _, param := range rootExec.Launch.Params {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Request != nil:
		for _, param := range rootExec.Request.Params {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Render != nil:
		for _, param := range rootExec.Render.Params {
			if param.SecretRef != "" {
				return true
			}
		}
	case rootExec.Serial != nil:
		for _, param := range rootExec.Serial.Params {
			if param.SecretRef != "" {
				return true
			}
		}
		for _, child := range rootExec.Serial.Refs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			if authRequired(ctx, childExec) {
				return true
			}
		}
		for _, e := range rootExec.Serial.Execs {
			if e.Ref != "" {
				childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, e.Ref)
				if err != nil {
					continue
				}
				if authRequired(ctx, childExec) {
					return true
				}
			}
		}
	case rootExec.Parallel != nil:
		for _, param := range rootExec.Parallel.Params {
			if param.SecretRef != "" {
				return true
			}
		}
		for _, child := range rootExec.Parallel.Refs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			if authRequired(ctx, childExec) {
				return true
			}
		}
		for _, e := range rootExec.Parallel.Execs {
			if e.Ref != "" {
				childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, e.Ref)
				if err != nil {
					continue
				}
				if authRequired(ctx, childExec) {
					return true
				}
			}
		}
	}
	return false
}

//nolint:gocognit
func pendingFormFields(ctx *context.Context, rootExec *executable.Executable) []*components.FormField {
	pending := make([]*components.FormField, 0)
	switch {
	case rootExec.Exec != nil:
		for _, param := range rootExec.Exec.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
	case rootExec.Launch != nil:
		for _, param := range rootExec.Launch.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
	case rootExec.Request != nil:
		for _, param := range rootExec.Request.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
	case rootExec.Render != nil:
		for _, param := range rootExec.Render.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
	case rootExec.Serial != nil:
		for _, param := range rootExec.Serial.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
		for _, child := range rootExec.Serial.Refs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			childPending := pendingFormFields(ctx, childExec)
			pending = append(pending, childPending...)
		}
		for _, child := range rootExec.Serial.Execs {
			if child.Ref != "" {
				childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child.Ref)
				if err != nil {
					continue
				}
				childPending := pendingFormFields(ctx, childExec)
				pending = append(pending, childPending...)
			}
		}
	case rootExec.Parallel != nil:
		for _, param := range rootExec.Parallel.Params {
			if param.Prompt != "" {
				pending = append(pending, &components.FormField{Key: param.EnvKey, Title: param.Prompt})
			}
		}
		for _, child := range rootExec.Parallel.Refs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child)
			if err != nil {
				continue
			}
			childPending := pendingFormFields(ctx, childExec)
			pending = append(pending, childPending...)
		}
		for _, child := range rootExec.Parallel.Execs {
			if child.Ref != "" {
				childExec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, child.Ref)
				if err != nil {
					continue
				}
				childPending := pendingFormFields(ctx, childExec)
				pending = append(pending, childPending...)
			}
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

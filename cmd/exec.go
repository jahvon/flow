package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/ui/types"
	"github.com/jahvon/flow/internal/io/ui/views"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/exec"
	"github.com/jahvon/flow/internal/runner/launch"
	"github.com/jahvon/flow/internal/runner/parallel"
	"github.com/jahvon/flow/internal/runner/render"
	"github.com/jahvon/flow/internal/runner/request"
	"github.com/jahvon/flow/internal/runner/serial"
	"github.com/jahvon/flow/internal/vault"
)

var execCmd = &cobra.Command{
	Use:     "exec <executable-id>",
	Aliases: config.SortedValidVerbs(),
	Short:   "Execute a flow by ID.",
	Long: "Execute a flow where <executable-id> is the target executable's ID in the form of 'ws/ns:name'.\n" +
		"The flow subcommand used should match the target executable's verb or one of its aliases.\n\n" +
		"See " + io.DocsURL("executable-verbs") + "for more information on executable verbs." +
		"See " + io.DocsURL("executable-ids") + "for more information on executable IDs.",
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		startApp(cmd, args)
		runner.RegisterRunner(exec.NewRunner())
		runner.RegisterRunner(launch.NewRunner())
		runner.RegisterRunner(request.NewRunner())
		runner.RegisterRunner(render.NewRunner())
		runner.RegisterRunner(serial.NewRunner())
		runner.RegisterRunner(parallel.NewRunner())
		setTermView(cmd, args)
	},
	PostRun: waitForExit,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		verbStr := cmd.CalledAs()
		verb := config.Verb(verbStr)
		if err := verb.Validate(); err != nil {
			logger.FatalErr(err)
		}

		idArg := args[0]
		ref := context.ExpandRef(curCtx, config.NewRef(idArg, verb))
		executable, err := curCtx.ExecutableCache.GetExecutableByRef(ref)
		if err != nil {
			handleError(err)
		} else if executable == nil {
			handleError(fmt.Errorf("executable %s not found", ref))
		}

		if err := executable.Validate(); err != nil {
			handleError(err)
		}

		if !executable.IsExecutableFromWorkspace(curCtx.UserConfig.CurrentWorkspace) {
			handleError(fmt.Errorf(
				"executable '%s' cannot be executed from workspace %s",
				ref,
				curCtx.UserConfig.CurrentWorkspace,
			))
		}

		setAuthEnv(executable)
		if interactiveUIEnabled() {
			curCtx.App.SetNotice("... processing ...", types.NoticeLevelInfo)
		}
		textInputs := pendingTextInputs(curCtx, executable)
		var envMap map[string]string
		if len(textInputs) > 0 {
			envMap = processUserInput(textInputs...)
		}
		startTime := time.Now()
		if err := runner.Exec(curCtx, executable, envMap); err != nil {
			handleError(err)
		}
		dur := time.Since(startTime)
		logger.PlainTextSuccess(fmt.Sprintf("%s flow completed", ref))
		if interactiveUIEnabled() {
			curCtx.App.SetNotice(fmt.Sprintf("Elapsed: %s", dur.Round(time.Millisecond)), types.NoticeLevelInfo)
			if dur > 1*time.Minute && curCtx.UserConfig.Interactive.SoundOnCompletion {
				_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
			}
			if dur > 1*time.Minute && curCtx.UserConfig.Interactive.NotifyOnCompletion {
				_ = beeep.Notify("Flow", "Flow completed", "")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func setAuthEnv(executable *config.Executable) {
	if authRequired(curCtx, executable) {
		if interactiveUIEnabled() {
			curCtx.App.SetNotice("... authenticating ...", types.NoticeLevelInfo)
		}
		resp := processUserInput(&views.TextInput{
			Key:    vault.EncryptionKeyEnvVar,
			Prompt: "Enter vault encryption key",
			Hidden: true,
		})
		val, ok := resp[vault.EncryptionKeyEnvVar]
		if !ok || val == "" {
			handleError(fmt.Errorf("vault encryption key required"))
		}
		if err := os.Setenv(vault.EncryptionKeyEnvVar, val); err != nil {
			handleError(fmt.Errorf("failed to set vault encryption key\n%w", err))
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
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(child)
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
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(child)
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
func pendingTextInputs(ctx *context.Context, rootExec *config.Executable) []*views.TextInput {
	pending := make([]*views.TextInput, 0)
	if rootExec.Type == nil {
		return nil
	}
	switch {
	case rootExec.Type.Exec != nil:
		for _, param := range rootExec.Type.Exec.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Launch != nil:
		for _, param := range rootExec.Type.Launch.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Request != nil:
		for _, param := range rootExec.Type.Request.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Render != nil:
		for _, param := range rootExec.Type.Render.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
	case rootExec.Type.Serial != nil:
		for _, param := range rootExec.Type.Serial.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
		for _, child := range rootExec.Type.Serial.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(child)
			if err != nil {
				continue
			}
			childPending := pendingTextInputs(ctx, childExec)
			pending = append(pending, childPending...)
		}
	case rootExec.Type.Parallel != nil:
		for _, param := range rootExec.Type.Parallel.Parameters {
			if param.Prompt != "" {
				pending = append(pending, &views.TextInput{Key: param.EnvKey, Prompt: param.Prompt})
			}
		}
		for _, child := range rootExec.Type.Parallel.ExecutableRefs {
			childExec, err := ctx.ExecutableCache.GetExecutableByRef(child)
			if err != nil {
				continue
			}
			childPending := pendingTextInputs(ctx, childExec)
			pending = append(pending, childPending...)
		}
	}
	return pending
}

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
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

var execCmd = &cobra.Command{
	Use:     "exec <executable-id>",
	Aliases: config.SortedValidVerbs(),
	Short:   "Execute a flow by ID.",
	Long: "Execute a flow where <executable-id> is the target executable's ID in the form of 'ws/ns:name'.\n" +
		"The flow subcommand used should match the target executable's verb or one of its aliases.\n\n" +
		"See " + io.ConfigDocsURL("executables", "Verb") + "for more information on executable verbs." +
		"See " + io.ConfigDocsURL("executables", "Ref") + "for more information on executable IDs.",
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		runner.RegisterRunner(exec.NewRunner())
		runner.RegisterRunner(launch.NewRunner())
		runner.RegisterRunner(request.NewRunner())
		runner.RegisterRunner(render.NewRunner())
		runner.RegisterRunner(serial.NewRunner())
		runner.RegisterRunner(parallel.NewRunner())
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		verbStr := cmd.CalledAs()
		verb := config.Verb(verbStr)
		if err := verb.Validate(); err != nil {
			logger.FatalErr(err)
		}

		idArg := args[0]
		ref := context.ExpandRef(curCtx, config.NewRef(idArg, verb))
		executable, err := curCtx.ExecutableCache.GetExecutableByRef(logger, ref)
		if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
			logger.Debugf("Executable %s not found in cache, syncing cache", ref)
			if err := curCtx.ExecutableCache.Update(logger); err != nil {
				logger.FatalErr(err)
			}
			executable, err = curCtx.ExecutableCache.GetExecutableByRef(logger, ref)
		}
		if err != nil {
			logger.FatalErr(err)
		}

		if err := executable.Validate(); err != nil {
			logger.FatalErr(err)
		}

		if !executable.IsExecutableFromWorkspace(curCtx.UserConfig.CurrentWorkspace) {
			logger.FatalErr(fmt.Errorf(
				"executable '%s' cannot be executed from workspace %s",
				ref,
				curCtx.UserConfig.CurrentWorkspace,
			))
		}

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Notice = fmt.Sprintf("Executing %s", ref)
			header.NoticeLevel = components.NoticeLevelInfo
			header.Print()
		}
		setAuthEnv(executable)
		textInputs := pendingTextInputs(curCtx, executable)
		var envMap map[string]string
		if len(textInputs) > 0 {
			inputs, err := components.ProcessInputs(io.Styles(), textInputs...)
			if err != nil {
				logger.FatalErr(err)
			}
			envMap = inputs.ValueMap()
		}
		startTime := time.Now()
		if err := runner.Exec(curCtx, executable, envMap); err != nil {
			logger.FatalErr(err)
		}
		dur := time.Since(startTime)
		logger.Infox(fmt.Sprintf("%s flow completed", ref), "Elapsed", dur.Round(time.Millisecond))
		if interactiveUIEnabled() {
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
		resp, err := components.ProcessInputs(
			io.Styles(),
			&components.TextInput{
				Key:    vault.EncryptionKeyEnvVar,
				Prompt: "Enter vault encryption key",
				Hidden: true,
			})
		if err != nil {
			curCtx.Logger.FatalErr(err)
		}
		val := resp.ValueMap()[vault.EncryptionKeyEnvVar]
		if val == "" {
			curCtx.Logger.FatalErr(fmt.Errorf("vault encryption key required"))
		}
		if err := os.Setenv(vault.EncryptionKeyEnvVar, val); err != nil {
			curCtx.Logger.FatalErr(fmt.Errorf("failed to set vault encryption key\n%w", err))
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

package serial

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/engine"
	argUtils "github.com/jahvon/flow/internal/utils/args"
	execUtils "github.com/jahvon/flow/internal/utils/executables"
	"github.com/jahvon/flow/types/executable"
)

type serialRunner struct{}

func NewRunner() runner.Runner {
	return &serialRunner{}
}

func (r *serialRunner) Name() string {
	return "serial"
}

func (r *serialRunner) IsCompatible(executable *executable.Executable) bool {
	if executable == nil || executable.Serial == nil {
		return false
	}
	return true
}

func (r *serialRunner) Exec(
	ctx *context.Context,
	e *executable.Executable,
	eng engine.Engine,
	inputEnv map[string]string,
) error {
	serialSpec := e.Serial
	if err := runner.SetEnv(ctx.Logger, e.Env(), inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if len(serialSpec.Execs) > 0 {
		return handleExec(ctx, e, eng, serialSpec, inputEnv)
	}
	return fmt.Errorf("no serial executables to run")
}

func handleExec(
	ctx *context.Context,
	parent *executable.Executable,
	eng engine.Engine,
	serialSpec *executable.SerialExecutableType,
	promptedEnv map[string]string,
) error {
	var execs []engine.Exec
	for i, refConfig := range serialSpec.Execs {
		var exec *executable.Executable
		switch {
		case refConfig.Ref != "":
			var err error
			exec, err = execUtils.ExecutableForRef(ctx, refConfig.Ref)
			if err != nil {
				return err
			}
		case refConfig.Cmd != "":
			exec = execUtils.ExecutableForCmd(parent, refConfig.Cmd, i)
		default:
			return errors.New("serial executable must have a ref or cmd")
		}
		ctx.Logger.Debugf("executing %s (%d/%d)", exec.Ref(), i+1, len(serialSpec.Execs))

		execPromptedEnv := make(map[string]string)
		maps.Copy(promptedEnv, execPromptedEnv)
		if len(refConfig.Args) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Args, execPromptedEnv)
			if err != nil {
				ctx.Logger.Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}
		fields := map[string]interface{}{"step": exec.ID()}
		exec.Exec.SetLogFields(fields)

		runExec := func() error {
			err := runner.Exec(ctx, exec, eng, execPromptedEnv)
			if err != nil {
				return err
			}
			if i < len(serialSpec.Execs) && refConfig.ReviewRequired {
				ctx.Logger.Println("Do you want to proceed with the next execution? (y/n)")
				if !inputConfirmed(os.Stdin) {
					return fmt.Errorf("stopping runner early (%d/%d)", i+1, len(serialSpec.Execs))
				}
			}
			return nil
		}

		execs = append(execs, engine.Exec{ID: exec.Ref().String(), Function: runExec, MaxRetries: refConfig.Retries})
	}
	results := eng.Execute(ctx.Ctx, execs, engine.WithMode(engine.Serial), engine.WithFailFast(parent.Serial.FailFast))
	if results.HasErrors() {
		return errors.New(results.String())
	}
	return nil
}

func inputConfirmed(in *os.File) bool {
	var response string
	_, _ = fmt.Fscanf(in, response)
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

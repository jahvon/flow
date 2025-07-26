package serial

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/logger"
	"github.com/flowexec/flow/internal/runner"
	"github.com/flowexec/flow/internal/runner/engine"
	"github.com/flowexec/flow/internal/services/expr"
	"github.com/flowexec/flow/internal/services/store"
	argUtils "github.com/flowexec/flow/internal/utils/args"
	execUtils "github.com/flowexec/flow/internal/utils/executables"
	"github.com/flowexec/flow/types/executable"
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
	if err := runner.SetEnv(ctx.Config.CurrentVaultName(), e.Env(), inputEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if len(serialSpec.Execs) > 0 {
		str, err := store.NewStore(store.Path())
		if err != nil {
			return err
		}
		if err := str.CreateBucket(store.EnvironmentBucket()); err != nil {
			return err
		}
		cacheData, err := str.GetAll()
		if err != nil {
			return err
		}
		if err := str.Close(); err != nil {
			logger.Log().Error(err, "unable to close store")
		}

		return handleExec(ctx, e, eng, serialSpec, inputEnv, cacheData)
	}
	return fmt.Errorf("no serial executables to run")
}

//nolint:gocognit
func handleExec(
	ctx *context.Context,
	parent *executable.Executable,
	eng engine.Engine,
	serialSpec *executable.SerialExecutableType,
	promptedEnv map[string]string,
	cacheData map[string]string,
) error {
	dataMap := expr.ExpressionEnv(ctx, parent, cacheData, promptedEnv)

	var execs []engine.Exec
	for i, refConfig := range serialSpec.Execs {
		if refConfig.If != "" {
			truthy, err := expr.IsTruthy(refConfig.If, &dataMap)
			if err != nil {
				return err
			}
			if !truthy {
				logger.Log().Debugf("skipping execution %d/%d", i+1, len(serialSpec.Execs))
				continue
			}
			logger.Log().Debugf("condition %s is true", refConfig.If)
		}
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
		logger.Log().Debugf("executing %s (%d/%d)", exec.Ref(), i+1, len(serialSpec.Execs))

		execPromptedEnv := make(map[string]string)
		maps.Copy(promptedEnv, execPromptedEnv)
		if len(refConfig.Args) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Args, execPromptedEnv)
			if err != nil {
				logger.Log().Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}

		switch {
		case exec.Exec != nil:
			fields := map[string]interface{}{"step": exec.Ref().String()}
			exec.Exec.SetLogFields(fields)
			if serialSpec.Dir != "" && exec.Exec.Dir == "" {
				exec.Exec.Dir = serialSpec.Dir
			}
		case exec.Parallel != nil:
			if serialSpec.Dir != "" && exec.Parallel.Dir == "" {
				exec.Parallel.Dir = serialSpec.Dir
			}
		case exec.Serial != nil:
			if serialSpec.Dir != "" && exec.Serial.Dir == "" {
				exec.Serial.Dir = serialSpec.Dir
			}
		}

		runExec := func() error {
			return runSerialExecFunc(ctx, i, refConfig, exec, eng, execPromptedEnv, serialSpec)
		}

		execs = append(execs, engine.Exec{ID: exec.Ref().String(), Function: runExec, MaxRetries: refConfig.Retries})
	}
	results := eng.Execute(ctx.Ctx, execs, engine.WithMode(engine.Serial), engine.WithFailFast(parent.Serial.FailFast))
	if results.HasErrors() {
		return errors.New(results.String())
	}
	return nil
}

func runSerialExecFunc(
	ctx *context.Context,
	step int,
	refConfig executable.SerialRefConfig,
	exec *executable.Executable,
	eng engine.Engine,
	execPromptedEnv map[string]string,
	serialSpec *executable.SerialExecutableType,
) error {
	err := runner.Exec(ctx, exec, eng, execPromptedEnv)
	if err != nil {
		return err
	}
	if step < len(serialSpec.Execs) && refConfig.ReviewRequired {
		logger.Log().Println("Do you want to proceed with the next execution? (y/n)")
		if !inputConfirmed(os.Stdin) {
			return fmt.Errorf("stopping runner early (%d/%d)", step+1, len(serialSpec.Execs))
		}
	}
	return nil
}

func inputConfirmed(in *os.File) bool {
	var response string
	_, _ = fmt.Fscanf(in, response)
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

package serial

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/utils"
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

func (r *serialRunner) Exec(ctx *context.Context, e *executable.Executable, promptedEnv map[string]string) error {
	serialSpec := e.Serial
	if err := runner.SetEnv(ctx.Logger, e.Env(), promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	if err := utils.ValidateOneOf("executable list", serialSpec.Refs, serialSpec.Execs); err != nil {
		return err
	}

	if len(serialSpec.Refs) > 0 {
		return handleExecRef(ctx, serialSpec, promptedEnv)
	} else if len(serialSpec.Execs) > 0 {
		return handleExec(ctx, e, serialSpec, promptedEnv)
	}
	return fmt.Errorf("no serial executables to run")
}

func handleExecRef(ctx *context.Context, serialSpec *executable.SerialExecutableType, promptedEnv map[string]string) error {
	order := serialSpec.Refs

	var errs []error
	for i, executableRef := range order {
		ctx.Logger.Debugf("executing %s (%d/%d)", executableRef, i+1, len(order))
		exec, err := execUtils.ExecutableForRef(ctx, executableRef)
		if err != nil {
			return err
		} else if exec == nil {
			return fmt.Errorf("unable to find e with reference %s", executableRef)
		}

		if exec.Exec != nil {
			fields := map[string]interface{}{
				"e": exec.ID(),
			}
			exec.Exec.SetLogFields(fields)
		}

		if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
			if serialSpec.FailFast {
				return errors.Wrapf(err, "execution error ref='%s'", executableRef)
			}
			errs = append(errs, err)
			ctx.Logger.Error(err, fmt.Sprintf("execution error for %s", executableRef))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d execution errors - %v", len(errs), errs)
	}
	return nil
}

//nolint:gocognit
func handleExec(
	ctx *context.Context,
	parent *executable.Executable,
	serialSpec *executable.SerialExecutableType,
	promptedEnv map[string]string,
) error {
	var errs []error
	for i, refConfig := range serialSpec.Execs {
		var exec *executable.Executable
		switch {
		case len(refConfig.Ref) > 0:
			var err error
			exec, err = execUtils.ExecutableForRef(ctx, refConfig.Ref)
			if err != nil {
				return err
			}
		case refConfig.Cmd != "":
			exec = execUtils.ExecutableForCmd(parent, refConfig.Cmd, i)
		}
		ctx.Logger.Debugf("executing %s (%d/%d)", exec.Ref(), i+1, len(serialSpec.Execs))

		execPromptedEnv := make(map[string]string)
		maps.Copy(promptedEnv, execPromptedEnv)
		if len(refConfig.Args) > 0 {
			a, err := argUtils.ProcessArgs(exec, refConfig.Args)
			if err != nil {
				ctx.Logger.Error(err, "unable to process arguments")
			}
			maps.Copy(execPromptedEnv, a)
		}

		for {
			if err := runner.Exec(ctx, exec, execPromptedEnv); err != nil {
				switch {
				case refConfig.Retries == 0 && serialSpec.FailFast:
					return errors.Wrapf(err, "execution error ref='%s'", exec.Ref())
				case refConfig.Retries == 0 && !serialSpec.FailFast:
					errs = append(errs, err)
					ctx.Logger.Error(err, fmt.Sprintf("execution error ref='%s'", exec.Ref()))
				case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && serialSpec.FailFast:
					return fmt.Errorf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries)
				case refConfig.Retries != 0 && refConfig.AttemptedMaxTimes() && !serialSpec.FailFast:
					errs = append(errs, err)
					ctx.Logger.Error(err, fmt.Sprintf("retries exceeded ref='%s' max=%d", exec.Ref(), refConfig.Retries))
				case refConfig.Retries != 0 && !refConfig.AttemptedMaxTimes():
					refConfig.RecordAttempt()
					ctx.Logger.Warnf("retrying ref='%s'", exec.Ref())
				default:
					return errors.Wrapf(err, "unexpected error handling ref='%s'", exec.Ref())
				}
			}
			break
		}

		if i < len(serialSpec.Execs) && refConfig.ReviewRequired {
			ctx.Logger.Println("Do you want to proceed with the next execution? (y/n)")
			if !inputConfirmed(os.Stdin) {
				return fmt.Errorf("stopping runner early (%d/%d)", i+1, len(serialSpec.Execs))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d execution errors - %v", len(errs), errs)
	}
	return nil
}

func inputConfirmed(in *os.File) bool {
	var response string
	_, _ = fmt.Fscanf(in, response)
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

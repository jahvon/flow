package serial

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
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

	order := serialSpec.Refs
	var errs []error
	for i, executableRef := range order {
		ctx.Logger.Debugf("executing %s (%d/%d)", executableRef, i+1, len(order))
		executableRef = context.ExpandRef(ctx, executableRef)
		exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, executableRef)
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
				return errors.Wrapf(err, "execution error for %s", executableRef)
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

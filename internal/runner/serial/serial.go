package serial

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
)

type serialRunner struct{}

func NewRunner() runner.Runner {
	return &serialRunner{}
}

func (r *serialRunner) Name() string {
	return "serial"
}

func (r *serialRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Serial == nil {
		return false
	}
	return true
}

func (r *serialRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	serialSpec := executable.Type.Serial
	if err := runner.SetEnv(ctx.Logger, &serialSpec.ExecutableEnvironment, promptedEnv); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}

	order := serialSpec.ExecutableRefs
	for i, executableRef := range order {
		ctx.Logger.Debugf("executing %s (%d/%d)", executableRef, i+1, len(order))
		executableRef = context.ExpandRef(ctx, executableRef)
		exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, executableRef)
		if err != nil {
			return err
		} else if exec == nil {
			return fmt.Errorf("unable to find executable with reference %s", executableRef)
		}

		if exec.Type.Exec != nil {
			fields := map[string]interface{}{
				"executable": exec.ID(),
			}
			exec.Type.Exec.SetLogFields(fields)
		}

		if err := runner.Exec(ctx, exec, promptedEnv); err != nil {
			if serialSpec.FailFast {
				return errors.Wrapf(err, "execution error for %s", executableRef)
			}
			ctx.Logger.Error(err, fmt.Sprintf("execution error for %s", executableRef))
		}
	}

	return nil
}

package serial

import (
	"fmt"

	"github.com/rs/zerolog/log"

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

func (r *serialRunner) Exec(ctx *context.Context, executable *config.Executable) error {
	serialSpec := executable.Type.Serial
	if err := runner.SetEnv(&serialSpec.ParameterizedExecutable); err != nil {
		return fmt.Errorf("unable to set parameters to env - %w", err)
	}

	order := serialSpec.ExecutableRefs
	for i, executableRef := range order {
		log.Debug().Msgf("executing %s (%d/%d)", executableRef, i+1, len(order))
		exec, err := ctx.ExecutableCache.GetExecutableByRef(executableRef)
		if err != nil {
			return fmt.Errorf("unable to get executable by ref %s - %w", executableRef, err)
		} else if exec == nil {
			return fmt.Errorf("unable to find executable with reference %s", executableRef)
		}

		if err := runner.Exec(ctx, exec); err != nil {
			return fmt.Errorf("execution error for %s - %w", executableRef, err)
		}
	}

	return nil
}

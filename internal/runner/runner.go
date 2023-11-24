package runner

import (
	"fmt"
	"time"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log().With().Str("scope", "runner").Logger()

type Runner interface {
	Name() string
	Exec(ctx *context.Context, executable *config.Executable) error
	IsCompatible(executable *config.Executable) bool
}

var registeredRunners []Runner

func init() {
	registeredRunners = make([]Runner, 0)
}

func RegisterRunner(runner Runner) {
	registeredRunners = append(registeredRunners, runner)
}

func Exec(ctx *context.Context, executable *config.Executable) error {
	var assignedRunner Runner
	for _, runner := range registeredRunners {
		if runner.IsCompatible(executable) {
			assignedRunner = runner
			log.Trace().Msgf("assigned %s runner", assignedRunner.Name())
			break
		}
	}
	if assignedRunner == nil {
		return fmt.Errorf("comptatible runner not found for executable %s", executable.ID())
	}

	if executable.Timeout == 0 {
		return assignedRunner.Exec(ctx, executable)
	}

	done := make(chan error, 1)
	go func() {
		done <- assignedRunner.Exec(ctx, executable)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(executable.Timeout):
		return fmt.Errorf("timeout after %v", executable.Timeout)
	}
}

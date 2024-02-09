package runner

import (
	"fmt"
	"time"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

type Runner interface {
	Name() string
	Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error
	IsCompatible(executable *config.Executable) bool
}

var registeredRunners []Runner

func init() {
	registeredRunners = make([]Runner, 0)
}

func RegisterRunner(runner Runner) {
	registeredRunners = append(registeredRunners, runner)
}

func Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	var assignedRunner Runner
	for _, runner := range registeredRunners {
		if runner.IsCompatible(executable) {
			assignedRunner = runner
			break
		}
	}
	if assignedRunner == nil {
		return fmt.Errorf("comptatible runner not found for executable %s", executable.ID())
	}

	if executable.Timeout == 0 {
		return assignedRunner.Exec(ctx, executable, promptedEnv)
	}

	done := make(chan error, 1)
	go func() {
		done <- assignedRunner.Exec(ctx, executable, promptedEnv)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(executable.Timeout):
		return fmt.Errorf("timeout after %v", executable.Timeout)
	}
}

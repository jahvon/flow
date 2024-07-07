package launch

import (
	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/services/open"
	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/executable"
)

type launchRunner struct{}

func NewRunner() runner.Runner {
	return &launchRunner{}
}

func (r *launchRunner) Name() string {
	return "launch"
}

func (r *launchRunner) IsCompatible(executable *executable.Executable) bool {
	if executable == nil || executable.Launch == nil {
		return false
	}
	return true
}

func (r *launchRunner) Exec(ctx *context.Context, e *executable.Executable, inputEnv map[string]string) error {
	launchSpec := e.Launch
	envMap, err := runner.BuildEnvMap(
		ctx.Logger,
		e.Env(),
		inputEnv,
		runner.DefaultEnv(ctx, e),
	)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	if err := runner.SetEnv(ctx.Logger, e.Env(), envMap); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	targetURI := utils.ExpandDirectory(
		ctx.Logger,
		launchSpec.URI,
		e.WorkspacePath(),
		e.ConfigPath(),
		envMap,
	)

	if launchSpec.App != "" {
		return open.OpenWith(launchSpec.App, targetURI, launchSpec.Wait)
	}
	return open.Open(targetURI, launchSpec.Wait)
}

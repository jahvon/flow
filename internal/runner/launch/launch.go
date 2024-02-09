package launch

import (
	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/services/open"
	"github.com/jahvon/flow/internal/utils"
)

type launchRunner struct{}

func NewRunner() runner.Runner {
	return &launchRunner{}
}

func (r *launchRunner) Name() string {
	return "launch"
}

func (r *launchRunner) IsCompatible(executable *config.Executable) bool {
	if executable == nil || executable.Type == nil || executable.Type.Launch == nil {
		return false
	}
	return true
}

func (r *launchRunner) Exec(ctx *context.Context, executable *config.Executable, promptedEnv map[string]string) error {
	launchSpec := executable.Type.Launch
	envMap, err := runner.ParametersToEnvMap(ctx.Logger, &launchSpec.ParameterizedExecutable, promptedEnv)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	if err := runner.SetEnv(ctx.Logger, &launchSpec.ParameterizedExecutable, envMap); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	targetURI := utils.ExpandDirectory(
		ctx.Logger,
		launchSpec.URI,
		executable.WorkspacePath(),
		executable.DefinitionPath(),
		envMap,
	)

	if launchSpec.App != "" {
		return open.OpenWith(launchSpec.App, targetURI, launchSpec.Wait)
	}
	return open.Open(targetURI, launchSpec.Wait)
}

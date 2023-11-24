package launch

import (
	"fmt"

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

func (r *launchRunner) Exec(_ *context.Context, executable *config.Executable) error {
	launchSpec := executable.Type.Launch
	envMap, err := runner.ParametersToEnvMap(&launchSpec.ParameterizedExecutable)
	if err != nil {
		return fmt.Errorf("env setup failed  - %w", err)
	}
	if err := runner.SetEnv(&launchSpec.ParameterizedExecutable); err != nil {
		return fmt.Errorf("unable to set parameters to env - %w", err)
	}
	targetURI := utils.ExpandDirectory(launchSpec.URI, executable.WorkspacePath(), executable.DefinitionPath(), envMap)

	if launchSpec.App != "" {
		return open.OpenWith(launchSpec.App, targetURI, launchSpec.Wait)
	}
	return open.Open(targetURI, launchSpec.Wait)
}

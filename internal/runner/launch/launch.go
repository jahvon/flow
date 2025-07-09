package launch

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/runner"
	"github.com/flowexec/flow/internal/runner/engine"
	"github.com/flowexec/flow/internal/services/open"
	"github.com/flowexec/flow/internal/utils"
	"github.com/flowexec/flow/types/executable"
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

func (r *launchRunner) Exec(
	ctx *context.Context,
	e *executable.Executable,
	_ engine.Engine,
	inputEnv map[string]string,
) error {
	launchSpec := e.Launch
	envMap, err := runner.BuildEnvMap(
		ctx.Logger,
		ctx.Config.CurrentVaultName(),
		e.Env(),
		inputEnv,
		runner.DefaultEnv(ctx, e),
	)
	if err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	if err := runner.SetEnv(ctx.Logger, ctx.Config.CurrentVaultName(), e.Env(), envMap); err != nil {
		return errors.Wrap(err, "unable to set parameters to env")
	}
	launchSpec.URI = os.ExpandEnv(launchSpec.URI)
	targetURI := launchSpec.URI
	if !strings.HasPrefix(targetURI, "http") {
		targetURI = utils.ExpandDirectory(
			ctx.Logger,
			launchSpec.URI,
			e.WorkspacePath(),
			e.FlowFilePath(),
			envMap,
		)
	}

	if launchSpec.App != "" {
		return open.OpenWith(launchSpec.App, targetURI)
	}
	return open.Open(targetURI)
}

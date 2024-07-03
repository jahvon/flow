package builder

import (
	"time"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

func ExampleExecutableDefinition(ctx *context.Context, path string) *config.ExecutableDefinition {
	return &config.ExecutableDefinition{
		Namespace:  "examples",
		Visibility: config.VisibilityInternal,
		Executables: []*config.Executable{
			SimpleExec(ctx, "simple", path),
			ExecWithPauses(ctx, "with-pauses", path),
			ExecWithExitCode(ctx, "with-exit-code", path, 1),
			ExecWithTmpDir(ctx, "with-tmp-dir", path),
			ExecWithArgs(ctx, "with-args", path, config.ArgumentList{{EnvKey: "ARG1", Pos: 0}}),
			ExecWithParams(ctx, "with-params", path, config.ParameterList{{EnvKey: "PARAM1", Text: "value1"}}),
			ExecWithLogMode(ctx, "with-plaintext", path, tuikitIO.Text),
			ExecWithTimeout(ctx, "with-timeout", path, 3*time.Second),
		},
	}
}

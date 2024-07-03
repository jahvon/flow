package builder

import (
	"fmt"
	"strings"
	"time"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

func SimpleExec(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Simple executable.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: fmt.Sprintf("echo 'hello from %s'", name),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithPauses(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable with pauses.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: fmt.Sprintf(
					"echo 'hello from %[1]s'; sleep 1; echo 'hello from %[1]s'; sleep 1; echo 'hello from %[1]s'",
					name,
				),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithTmpDir(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable with a temporary directory specified with the tmp dir label.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
				Command:             fmt.Sprintf("echo 'hello from %[1]s';mkdir %[1]s; cd %[1]s; pwd", name),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithArgs(ctx *context.Context, name, definitionPath string, args config.ArgumentList) *config.Executable {
	var argStrs []string
	for _, arg := range args {
		if arg.Pos >= 0 {
			argStrs = append(argStrs, fmt.Sprintf("echo 'pos=%d, key=%s'", arg.Pos, arg.EnvKey))
		} else if arg.Flag != "" {
			argStrs = append(argStrs, fmt.Sprintf("echo 'flag=%s, key=%s'", arg.Flag, arg.EnvKey))
		}
	}
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable with arguments.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				ExecutableEnvironment: config.ExecutableEnvironment{Args: args},
				Command: fmt.Sprintf(
					"echo 'hello from %s'; %s",
					name, strings.Join(argStrs, "; "),
				),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithParams(ctx *context.Context, name, definitionPath string, params config.ParameterList) *config.Executable {
	var paramStrs []string
	for _, param := range params {
		if param.Text != "" {
			paramStrs = append(paramStrs, fmt.Sprintf("echo 'key=%s, value=%s'", param.EnvKey, param.Text))
		} else if param.SecretRef != "" {
			paramStrs = append(paramStrs, fmt.Sprintf("echo 'key=%s, secret=%s'", param.EnvKey, param.SecretRef))
		} else if param.Prompt != "" {
			paramStrs = append(
				paramStrs,
				fmt.Sprintf("echo 'key=%s, prompt=%s', value=$%[1]s", param.EnvKey, param.Prompt),
			)
		}
	}
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable with parameters.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				ExecutableEnvironment: config.ExecutableEnvironment{Parameters: params},
				Command: fmt.Sprintf(
					"echo 'hello from %s'; %s",
					name, strings.Join(paramStrs, "; "),
				),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e

}

func ExecWithLogMode(ctx *context.Context, name, definitionPath string, logMode tuikitIO.LogMode) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable with a log mode specified.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				LogMode: logMode,
				Command: fmt.Sprintf(
					"echo 'hello from %s'; echo 'line 2'; echo 'line 3'; echo 'line 4'",
					name,
				),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithExitCode(ctx *context.Context, name, definitionPath string, exitCode int) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable that will exit with the provided exit code.",
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: fmt.Sprintf("exit %d", exitCode),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func ExecWithTimeout(ctx *context.Context, name, definitionPath string, timeout time.Duration) *config.Executable {
	e := &config.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  config.VisibilityInternal.NewPointer(),
		Description: "Executable that will sleep for the provided duration.",
		Timeout:     timeout,
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: fmt.Sprintf("sleep %d", int(timeout.Seconds()+10)),
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

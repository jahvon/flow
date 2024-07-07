package builder

import (
	"fmt"
	"strings"
	"time"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/types/executable"
)

func SimpleExec(opts ...Option) *executable.Executable {
	name := "simple-print"
	e := &executable.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: privateExecVisibility(),
		Exec: &executable.ExecExecutableType{
			Cmd: fmt.Sprintf("echo 'hello from %s'", name),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithPauses(opts ...Option) *executable.Executable {
	name := "with-pauses"
	e := &executable.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: privateExecVisibility(),
		Exec: &executable.ExecExecutableType{
			Cmd: fmt.Sprintf(
				"echo 'hello from %[1]s'; sleep 1; echo 'hello from %[1]s'; sleep 1; echo 'hello from %[1]s'",
				name,
			),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithExit(opts ...Option) *executable.Executable {
	name := "with-exit"
	exitCode := 1
	e := &executable.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: privateExecVisibility(),
		Exec: &executable.ExecExecutableType{
			Cmd: fmt.Sprintf("exit %d", exitCode),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithTimeout(opts ...Option) *executable.Executable {
	name := "with-timeout"
	timeout := 3 * time.Second
	docstring := "The `timeout` field can be set to limit the amount of time the executable will run.\n" +
		"If the executable runs longer than the timeout, it will be killed and the execution will fail."
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Timeout:     timeout,
		Exec: &executable.ExecExecutableType{
			Cmd: fmt.Sprintf("sleep %d", int(timeout.Seconds()+10)),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithTmpDir(opts ...Option) *executable.Executable {
	name := "with-tmp-dir"
	docstring := fmt.Sprintf(
		"Executables will be run from a new temporary direction when the `dir` field is set to `%s`.\n"+
			"If the executable is a parallel or serial executable, all sub-executables will run from the same temporary directory.\n"+
			"Any files created during the execution will be deleted after the executable completes.",
		executable.TmpDirLabel,
	)
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Exec: &executable.ExecExecutableType{
			Dir: executable.Directory(executable.TmpDirLabel),
			Cmd: fmt.Sprintf("echo 'hello from %[1]s';mkdir %[1]s; cd %[1]s; pwd", name),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithArgs(opts ...Option) *executable.Executable {
	name := "with-args"
	args := executable.ArgumentList{{EnvKey: "ARG1", Pos: 1}, {EnvKey: "ARG2", Flag: "x", Default: "yz"}}
	var argCmds []string
	for _, arg := range args {
		if arg.Pos > 0 {
			argCmds = append(argCmds, fmt.Sprintf("echo 'pos=%d, key=%s'", arg.Pos, arg.EnvKey))
		} else if arg.Flag != "" {
			argCmds = append(argCmds, fmt.Sprintf("echo 'flag=%s, key=%s'", arg.Flag, arg.EnvKey))
		}
	}
	docstring := "Command line arguments can be passed to the executable using the `args` field.\n" +
		"Arguments can be positional or flags, and can have default values.\n" +
		"**You must specify the `envKey` field for each argument and one of `pos` or `flag`** " +
		"The value of the argument will be available in the environment variable specified by `envKey`.\n" +
		"The first positional argument is at position 1 and following arguments are at increasing positions. " +
		"Flags are specified with the defined flag value and is followed by `=` and it's value (no spaces).\n" +
		"If a default value is provided, it will be used if the argument is not provided. The executable will " +
		"fail if a required argument is not provided."
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Exec: &executable.ExecExecutableType{
			Args: args,
			Cmd:  fmt.Sprintf("echo 'hello from %s'; %s", name, strings.Join(argCmds, "; ")),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ExecWithParams(opts ...Option) *executable.Executable {
	name := "with-params"
	docstring := "Parameters can be passed to the executable using the `params` field.\n" +
		"Parameters can be text, secrets, or prompts. Text parameters will be available in the environment variable " +
		"specified by `envKey`. Secret parameters will be resolved from the secret store and will be available in the " +
		"environment variable specified by `envKey`. Prompt parameters will prompt the user for a value and will be " +
		"available in the environment variable specified by `envKey`."
	params := executable.ParameterList{
		{EnvKey: "PARAM1", Text: "value1"},
		{EnvKey: "PARAM2", SecretRef: "flow-example-secret"},
		{EnvKey: "PARAM3", Prompt: "Enter a value"},
	}
	var paramCmds []string
	for _, param := range params {
		if param.Text != "" {
			paramCmds = append(paramCmds, fmt.Sprintf("echo 'key=%s, value=%s'", param.EnvKey, param.Text))
		} else if param.SecretRef != "" {
			paramCmds = append(paramCmds, fmt.Sprintf("echo 'key=%s, secret=%s'", param.EnvKey, param.SecretRef))
		} else if param.Prompt != "" {
			paramCmds = append(
				paramCmds,
				fmt.Sprintf("echo 'key=%s, prompt=%s', value=$%[1]s", param.EnvKey, param.Prompt),
			)
		}
	}
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Exec: &executable.ExecExecutableType{
			Params: params,
			Cmd:    fmt.Sprintf("echo 'hello from %s'; %s", name, strings.Join(paramCmds, "; ")),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e

}

func ExecWithLogMode(opts ...Option) *executable.Executable {
	name := "with-plaintext"
	docstring := "The `logMode` field can be set to change the formatting of the executable's output logs.\n" +
		"Valid values are `logfmt`, `text`, `json`, and `hidden`. The default value is determined by the user's configuration."
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Exec: &executable.ExecExecutableType{
			LogMode: tuikitIO.Text,
			Cmd: fmt.Sprintf(
				"echo 'hello from %s'; echo 'line 2'; echo 'line 3'; echo 'line 4'",
				name,
			),
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

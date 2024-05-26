package examples_test

import (
	"time"

	"github.com/jahvon/flow/config"
)

var (
	serialBaseDesc = "Serial flows are executed in the order they are defined, one after the other."
)

var SerialExecRoot = &config.Executable{
	Verb:       "start",
	Name:       "serial",
	Visibility: config.VisibilityPrivate.NewPointer(),
	Description: serialBaseDesc +
		"\n\n- The environment defined on the root executable is inherited by all the child executables." +
		"\n- Setting `f:tmp` as the directory on child executables will create a shared temporary directory for execution.",
	Timeout: 10 * time.Second,
	Type: &config.ExecutableTypeSpec{
		Serial: &config.SerialExecutableType{
			ExecutableEnvironment: config.ExecutableEnvironment{
				Parameters: []config.Parameter{
					{
						EnvKey: "PARAM1",
						Text:   "value1",
					},
				},
			},
			ExecutableRefs: []config.Ref{
				"run examples:serial-exec1",
				"run examples:serial-exec2",
				"run examples:serial-exec3",
			},
		},
	},
}

var SerialWithExitRoot = &config.Executable{
	Verb:       "start",
	Name:       "serial-with-exit",
	Aliases:    []string{"serial-exit"},
	Visibility: config.VisibilityPrivate.NewPointer(),
	Description: serialBaseDesc +
		"\n\n The `failFast` option can be set to `true` to stop the flow if a child executable fails.",
	Timeout: 10 * time.Second,
	Type: &config.ExecutableTypeSpec{
		Serial: &config.SerialExecutableType{
			FailFast: true,
			ExecutableRefs: []config.Ref{
				"run examples:serial-exec1",
				"run examples:serial-exec2",
				"run examples:exit",
				"run examples:serial-exec3",
			},
		},
	},
}

var SerialExec1 = &config.Executable{
	Verb:        "run",
	Name:        "serial-exec1",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "First serial executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'hello from 1';mkdir one;echo $PARAM1 > one/param1.txt",
		},
	},
}

var SerialExec2 = &config.Executable{
	Verb:        "run",
	Name:        "serial-exec2",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "Second serial executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'hello from 2';mkdir two;echo $PARAM1 > two/param1.txt",
		},
	},
}

var SerialExec3 = &config.Executable{
	Verb:        "run",
	Name:        "serial-exec3",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "Third serial executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'hello from 3';mkdir three;echo $PARAM1 > three/param1.txt;ls -1 one two three",
		},
	},
}

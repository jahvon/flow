package examples_test

import (
	"time"

	"github.com/jahvon/flow/config"
)

var ExecWithExit = &config.Executable{
	Verb:       "run",
	Name:       "exit",
	Visibility: config.VisibilityPrivate,
	Description: "This executable will exit with the provided exit code." +
		"\n\n- The exit code can be set using the `exitCode` parameter.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableEnvironment: config.ExecutableEnvironment{
				Args: config.ArgumentList{
					{
						EnvKey:  "exitCode",
						Default: "1",
						Type:    "int",
						Pos:     0,
					},
				},
			},
			Command: "exit $exitCode",
		},
	},
}

var ExecWithTimeout = &config.Executable{
	Verb:       "run",
	Name:       "timeout",
	Visibility: config.VisibilityPrivate,
	Description: "This executable will sleep for the provided duration." +
		"\n\n- The sleep duration can be set using the `duration` parameter.",
	Timeout: 250 * time.Millisecond,
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableEnvironment: config.ExecutableEnvironment{
				Args: config.ArgumentList{
					{
						EnvKey:  "duration",
						Default: "1",
						Type:    "int",
						Pos:     0,
					},
				},
			},
			Command: "sleep $duration",
		},
	},
}

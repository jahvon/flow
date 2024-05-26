package examples_test

import (
	"time"

	"github.com/jahvon/flow/config"
)

var (
	parallelBaseDesc = "Parallel flows are executed concurrently."
)

var ParallelExecRoot = &config.Executable{
	Verb:       "start",
	Name:       "parallel",
	Visibility: config.VisibilityPrivate.NewPointer(),
	Description: parallelBaseDesc +
		"\n\n- The environment defined on the root executable is inherited by all the child executables." +
		"\n- Setting `f:tmp` as the directory on child executables will create a shared temporary directory for execution.",
	Timeout: 10 * time.Second,
	Type: &config.ExecutableTypeSpec{
		Parallel: &config.ParallelExecutableType{
			ExecutableRefs: []config.Ref{
				"run examples:parallel-exec1",
				"run examples:parallel-exec2",
				"run examples:parallel-exec3",
			},
		},
	},
}

var ParallelExecRootWithExit = &config.Executable{
	Verb:       "start",
	Name:       "parallel-exit",
	Aliases:    []string{"parallel-with-exit"},
	Visibility: config.VisibilityPrivate.NewPointer(),
	Description: parallelBaseDesc +
		"\n\nThe `failFast` option can be set to `true` to stop the flow if a child executable fails.",
	Timeout: 10 * time.Second,
	Type: &config.ExecutableTypeSpec{
		Parallel: &config.ParallelExecutableType{
			FailFast: true,
			ExecutableRefs: []config.Ref{
				"run examples:parallel-exec1",
				"run examples:parallel-exec2",
				"run examples:exit",
				"run examples:parallel-exec3",
			},
		},
	},
}

var ParallelExecRootWithMaxThreads = &config.Executable{
	Verb:       "start",
	Name:       "parallel-max",
	Aliases:    []string{"parallel-with-max-threads"},
	Visibility: config.VisibilityPrivate.NewPointer(),
	Description: parallelBaseDesc +
		"\n\nThe `maxThreads` option can be set to limit the number of concurrent executions.",
	Timeout: 15 * time.Second,
	Type: &config.ExecutableTypeSpec{
		Parallel: &config.ParallelExecutableType{
			MaxThreads: 1,
			ExecutableRefs: []config.Ref{
				"run examples:parallel-exec1",
				"run examples:parallel-exec2",
				"run examples:parallel-exec3",
			},
		},
	},
}

var ParallelExec1 = &config.Executable{
	Verb:        "run",
	Name:        "parallel-exec1",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "First parallel executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'mkdir one;sleep 1;hello from 1;sleep 1;hello from 1'",
		},
	},
}

var ParallelExec2 = &config.Executable{
	Verb:        "run",
	Name:        "parallel-exec2",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "Second parallel executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'mkdir two;hello from 2;sleep 1;hello from 2'",
		},
	},
}

var ParallelExec3 = &config.Executable{
	Verb:        "run",
	Name:        "parallel-exec3",
	Visibility:  config.VisibilityInternal.NewPointer(),
	Description: "Third parallel executable.",
	Type: &config.ExecutableTypeSpec{
		Exec: &config.ExecExecutableType{
			ExecutableDirectory: config.ExecutableDirectory{Directory: config.TmpDirLabel},
			Command:             "echo 'hello from 3;sleep 2;hello from 3;ls -1'",
		},
	},
}

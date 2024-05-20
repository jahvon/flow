package examples_test

import "github.com/jahvon/flow/config"

// TODO: convert testdata into generated documentation
var TestExecutableDefinition = config.ExecutableDefinition{
	Namespace:  "examples",
	Visibility: config.VisibilityInternal,
	Executables: []*config.Executable{
		ParallelExecRoot,
		ParallelExecRootWithExit,
		ParallelExecRootWithMaxThreads,
		ParallelExec1,
		ParallelExec2,
		ParallelExec3,
		SerialExecRoot,
		SerialWithExitRoot,
		SerialExec1,
		SerialExec2,
		SerialExec3,
		ExecWithExit,
	},
}

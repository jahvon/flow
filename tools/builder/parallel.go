//nolint:lll
package builder

import (
	"github.com/jahvon/flow/types/executable"
)

const (
	parallelBaseDesc = "Multiple executables can be run concurrently using a parallel executable."
)

func ParallelExecByRef(opts ...Option) *executable.Executable {
	name := "parallel"
	docstring := parallelBaseDesc +
		"The `refs` field is required and must be a valid reference to an executable.\n" +
		"The environment derived from parameters and arguments on the root executable is inherited by all the sub-executables.\n"
	e1 := ExecWithPauses(opts...)
	e2 := ExecWithPauses(opts...)
	e3 := ExecWithPauses(opts...)
	e := &executable.Executable{
		Verb:        "start",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Parallel: &executable.ParallelExecutableType{
			Refs: []executable.Ref{e1.Ref(), e2.Ref(), e3.Ref()},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ParallelExecByRefConfig(opts ...Option) *executable.Executable {
	name := "parallel-config"
	e1 := SimpleExec(opts...)
	e2 := ExecWithArgs(opts...)
	e := &executable.Executable{
		Verb:       "start",
		Name:       name,
		Visibility: privateExecVisibility(),
		Description: parallelBaseDesc +
			"The `execs` field can be used to define the child executables with more options. " +
			"This includes defining an executable inline, retries, arguments, and more.",
		Parallel: &executable.ParallelExecutableType{
			Execs: []executable.ParallelRefConfig{
				{Ref: e1.Ref()},
				{Ref: e2.Ref(), Args: []string{"hello", "x=123"}},
				{Cmd: "echo 'hello from serial command'"},
			},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ParallelExecWithExit(opts ...Option) *executable.Executable {
	name := "parallel-with-failure"
	e1 := ExecWithPauses(opts...)
	e2 := ExecWithExit(opts...)
	e3 := ExecWithPauses(opts...)
	e := &executable.Executable{
		Verb:       "start",
		Name:       name,
		Aliases:    []string{"parallel-exit"},
		Visibility: privateExecVisibility(),
		Description: parallelBaseDesc +
			"The `failFast` option can be set to `true` to stop the flow if a sub-executable fails.",
		Parallel: &executable.ParallelExecutableType{
			FailFast: true,
			Refs:     executable.RefList{e1.Ref(), e2.Ref(), e3.Ref()},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func ParallelExecWithMaxThreads(opts ...Option) *executable.Executable {
	e := ParallelExecByRef(opts...)
	e.Description = parallelBaseDesc +
		"\n\nThe `maxThreads` option can be set to limit the number of concurrent executions."
	e.Parallel.MaxThreads = 1
	return e
}

package builder

import (
	"github.com/jahvon/flow/types/executable"
)

const (
	serialBaseDesc = "Multiple executables can be run in sequence using a serial executable.\n"
)

func SerialExecByRef(opts ...Option) *executable.Executable {
	name := "serial"
	docstring := serialBaseDesc +
		"The `refs` field is required and must be a valid reference to an executable.\n" +
		"The environment derived from parameters and arguments on the root executable is inherited by all the sub-executables."
	e1 := SimpleExec(opts...)
	e2 := SimpleExec(opts...)
	e3 := SimpleExec(opts...)
	e := &executable.Executable{
		Verb:        "start",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Serial: &executable.SerialExecutableType{
			Refs: executable.RefList{e1.Ref(), e2.Ref(), e3.Ref()},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func SerialExecWithExit(opts ...Option) *executable.Executable {
	name := "serial-with-failure"
	e1 := SimpleExec(opts...)
	e2 := ExecWithExit(opts...)
	e3 := SimpleExec(opts...)
	e := &executable.Executable{
		Verb:       "start",
		Name:       name,
		Aliases:    []string{"serial-exit"},
		Visibility: privateExecVisibility(),
		Description: serialBaseDesc +
			"The `failFast` option can be set to `true` to stop the executable if a sub-executable fails.",
		Serial: &executable.SerialExecutableType{
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

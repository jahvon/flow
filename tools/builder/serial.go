package main

import (
	"github.com/jahvon/flow/types/executable"
)

const (
	serialBaseDesc = "Multiple executables can be run in sequence using a serial executable.\n"
)

func SerialExecByRefConfig(opts ...Option) *executable.Executable {
	name := "serial-config"
	e1 := SimpleExec(opts...)
	e2 := ExecWithArgs(opts...)
	e := &executable.Executable{
		Verb:       "start",
		Name:       name,
		Visibility: privateExecVisibility(),
		Description: "The `execs` field can be used to define the child executables with more options. " +
			"This includes defining an executable inline, retries, arguments, and more.",
		Serial: &executable.SerialExecutableType{
			Execs: []executable.SerialRefConfig{
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

func SerialExecWithExit(opts ...Option) *executable.Executable {
	name := "serial-with-failure"
	e1 := SimpleExec(opts...)
	e2 := ExecWithExit(opts...)
	e3 := SimpleExec(opts...)
	ff := true
	e := &executable.Executable{
		Verb:        "start",
		Name:        name,
		Aliases:     []string{"serial-exit"},
		Visibility:  privateExecVisibility(),
		Description: "The `failFast` option can be set to `true` to stop the executable if a sub-executable fails.",
		Serial: &executable.SerialExecutableType{
			FailFast: &ff,
			Execs:    []executable.SerialRefConfig{{Ref: e1.Ref()}, {Ref: e2.Ref()}, {Ref: e3.Ref()}},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

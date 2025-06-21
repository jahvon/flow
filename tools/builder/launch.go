package main

import (
	"github.com/jahvon/flow/types/executable"
)

func LaunchGitHubExample(opts ...Option) *executable.Executable {
	name := "github"
	e := &executable.Executable{
		Verb:        "open",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: "Example of launching a web URL.",
		Launch: &executable.LaunchExecutableType{
			URI: "https://www.github.com",
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func LaunchMacSettingsExample(opts ...Option) *executable.Executable {
	name := "mac-settings"
	e := &executable.Executable{
		Verb:        "open",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: "Example of launching a macOS application.",
		Launch: &executable.LaunchExecutableType{
			App: "System Preferences",
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func LaunchWorkspaceExample(opts ...Option) *executable.Executable {
	name := "ws-config"
	e := &executable.Executable{
		Verb:        "open",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: "Example of launching a workspace path with wait option.",
		Launch: &executable.LaunchExecutableType{
			URI:  "$FLOW_WORKSPACE_PATH",
			Wait: true,
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

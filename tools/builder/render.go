package main

import (
	"github.com/jahvon/flow/types/executable"
)

func RenderMarkdownExample(opts ...Option) *executable.Executable {
	name := "markdown"
	e := &executable.Executable{
		Verb:        "render",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: "Example of rendering a markdown template with data and parameters.",
		Render: &executable.RenderExecutableType{
			TemplateFile:     "template.md",
			TemplateDataFile: "template-data.yaml",
			Params: executable.ParameterList{
				{Prompt: "What is your name?", EnvKey: "NAME"},
				{Text: "Hi", EnvKey: "GREETING"},
			},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

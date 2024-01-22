package config

import (
	"github.com/samber/lo"
)

type ExecutableDefinition struct {
	// +docsgen:namespace
	// The namespace of the executable definition. This is used to group executables together.
	// If not set, the executables in the definition will be grouped into the root (*) namespace.
	// Namespaces can be reused across multiple definitions.
	Namespace string `yaml:"namespace"`
	Tags      Tags   `yaml:"tags"`
	// +docsgen:visibility
	// The visibility of the executables to Flow.
	// If not set, the visibility will default to `public`.
	//
	// `public` executables can be executed and listed from anywhere.
	// `private` executables can be executed and listed only within their own workspace.
	// `internal` executables can be executed within their own workspace but are not listed.
	// `hidden` executables cannot be executed or listed.
	Visibility VisibilityType `yaml:"visibility"`
	// +docsgen:executables
	// A list of executables to be defined in the executable definition.
	Executables ExecutableList `yaml:"executables"`

	workspaceName, workspacePath, definitionPath string
}

type ExecutableDefinitionList []*ExecutableDefinition

func (d *ExecutableDefinition) SetContext(workspaceName, workspacePath, definitionPath string) {
	d.workspaceName = workspaceName
	d.workspacePath = workspacePath
	d.definitionPath = definitionPath
	for _, exec := range d.Executables {
		exec.SetContext(workspaceName, workspacePath, d.Namespace, definitionPath)
		exec.SetDefaults()
		exec.MergeTags(d.Tags)
		exec.MergeVisibility(d.Visibility)
	}
}

func (d *ExecutableDefinition) SetDefaults() {
	if d.Visibility == "" {
		d.Visibility = VisibilityPrivate
	}
}

func (d *ExecutableDefinition) WorkspacePath() string {
	return d.workspacePath
}

func (d *ExecutableDefinition) DefinitionPath() string {
	return d.definitionPath
}

func (l *ExecutableDefinitionList) FilterByNamespace(namespace string) ExecutableDefinitionList {
	definitions := lo.Filter(*l, func(definition *ExecutableDefinition, _ int) bool {
		return definition.Namespace == namespace
	})
	log.Trace().Int("definitions", len(definitions)).Msgf("filtered definitions by namespace %s", namespace)
	return definitions
}

func (l *ExecutableDefinitionList) FilterByTag(tag string) ExecutableDefinitionList {
	definitions := lo.Filter(*l, func(definition *ExecutableDefinition, _ int) bool {
		return definition.Tags.HasTag(tag)
	})
	log.Trace().Int("definitions", len(definitions)).Msgf("filtered definitions by tag %s", tag)
	return definitions
}

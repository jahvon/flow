package config

import (
	"github.com/samber/lo"
)

type ExecutableDefinition struct {
	Namespace   string         `yaml:"namespace"`
	Tags        Tags           `yaml:"tags"`
	Visibility  VisibilityType `yaml:"visibility"`
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

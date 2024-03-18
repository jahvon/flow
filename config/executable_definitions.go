package config

import (
	"errors"
	"fmt"

	"github.com/samber/lo"
)

type ExecutableDefinition struct {
	// +docsgen:namespace
	// The namespace of the executable definition. This is used to group executables together.
	// If not set, the executables in the definition will be grouped into the root (*) namespace.
	// Namespaces can be reused across multiple definitions.
	Namespace string `yaml:"namespace,omitempty"`
	Tags      Tags   `yaml:"tags,omitempty"`
	// +docsgen:visibility
	// The visibility of the executables to Flow.
	// If not set, the visibility will default to `public`.
	//
	// `public` executables can be executed and listed from anywhere.
	// `private` executables can be executed and listed only within their own workspace.
	// `internal` executables can be executed within their own workspace but are not listed.
	// `hidden` executables cannot be executed or listed.
	Visibility VisibilityType `yaml:"visibility,omitempty"`
	// +docsgen:executables
	// A list of executables to be defined in the executable definition.
	Executables ExecutableList `yaml:"executables,omitempty"`

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
	return definitions
}

func (l *ExecutableDefinitionList) FilterByTag(tag string) ExecutableDefinitionList {
	definitions := lo.Filter(*l, func(definition *ExecutableDefinition, _ int) bool {
		return definition.Tags.HasTag(tag)
	})
	return definitions
}

type TemplateDataEntry struct {
	// +docsgen:key
	// The key to associate the data with. This is used as the key in the template data map.
	Key string `yaml:"key"`
	// +docsgen:prompt
	// A prompt to be displayed to the user when collecting an input value.
	Prompt string `yaml:"prompt"`
	// +docsgen:default
	// The default value to use if the template data is not set.
	Default string `yaml:"default"`
	// +docsgen:required
	// If true, the template data must be set. If false, the default value will be used if the template data is not set.
	Required bool `yaml:"required"`

	value string
}

func (t *TemplateDataEntry) Set(value string) {
	t.value = value
}

func (t *TemplateDataEntry) Value() string {
	if t.value == "" {
		return t.Default
	}
	return t.value
}

func (t *TemplateDataEntry) Validate() error {
	if t.Prompt == "" {
		return errors.New("must specify prompt for template data")
	}
	if t.Key == "" {
		return errors.New("must specify key for template data")
	}
	return nil
}

func (t *TemplateDataEntry) ValidateValue() error {
	if t.value == "" && t.Required {
		return fmt.Errorf("required template data not set")
	}
	return nil
}

type TemplateData []TemplateDataEntry

func (t *TemplateData) Set(key, value string) {
	for i, entry := range *t {
		if entry.Key == key {
			(*t)[i].Set(value)
			return
		}
	}
}

func (t *TemplateData) MapInterface() map[string]interface{} {
	data := map[string]interface{}{}
	for _, entry := range *t {
		data[entry.Key] = entry.Value()
	}
	return data
}

func (t *TemplateData) Validate() error {
	for _, entry := range *t {
		if err := entry.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplateData) ValidateValues() error {
	for _, entry := range *t {
		if err := entry.ValidateValue(); err != nil {
			return err
		}
	}
	return nil
}

type ExecutableDefinitionTemplate struct {
	// +docsgen:data
	// A list of template data to be used when rendering the executable definition.
	Data TemplateData `yaml:"data"`
	// +docsgen:artifacts
	// A list of files to include when copying the template in a new location. The files are copied as-is.
	Artifacts []string `yaml:"artifacts,omitempty"`

	*ExecutableDefinition `yaml:",inline"`

	location string
}

func (t *ExecutableDefinitionTemplate) SetContext(location string) {
	t.location = location
}

func (t *ExecutableDefinitionTemplate) Location() string {
	return t.location
}

func (t *ExecutableDefinitionTemplate) Validate() error {
	return t.Data.Validate()
}

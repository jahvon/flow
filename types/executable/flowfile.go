package executable

import (
	"errors"
	"fmt"

	"github.com/jahvon/flow/types/common"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p executable -o flowfile.gen.go flowfile_schema.yaml

type FlowFileList []*FlowFile

func (f *FlowFile) SetContext(workspaceName, workspacePath, configPath string) {
	f.workspace = workspaceName
	f.workspacePath = workspacePath
	f.configPath = configPath
	for _, exec := range f.Executables {
		exec.SetContext(workspaceName, workspacePath, f.Namespace, configPath)
		if exec.Visibility == nil && f.Visibility != nil {
			v := ExecutableVisibility(*f.Visibility)
			exec.Visibility = &v
		}
		exec.SetDefaults()
		exec.SetInheritedFields(f)
	}
}

func (f *FlowFile) SetDefaults() {
	if f.Visibility == nil || *f.Visibility == "" {
		v := FlowFileVisibility(common.VisibilityPrivate)
		f.Visibility = &v
	}
}

func (f *FlowFile) WorkspacePath() string {
	return f.workspacePath
}

func (f *FlowFile) ConfigPath() string {
	return f.configPath
}

func (l *FlowFileList) FilterByNamespace(namespace string) FlowFileList {
	filteredCfgs := make(FlowFileList, 0)
	for _, cfg := range *l {
		if cfg.Namespace == namespace {
			filteredCfgs = append(filteredCfgs, cfg)
		}
	}
	return filteredCfgs
}

func (l *FlowFileList) FilterByTag(tag string) FlowFileList {
	filteredCfgs := make(FlowFileList, 0)
	for _, cfg := range *l {
		t := common.Tags(cfg.Tags)
		if t.HasTag(tag) {
			filteredCfgs = append(filteredCfgs, cfg)
		}
	}
	return filteredCfgs
}

type TemplateDataEntry struct {
	// The key to associate the data with. This is used as the key in the template data map.
	Key string `yaml:"key"`
	// A prompt to be displayed to the user when collecting an input value.
	Prompt string `yaml:"prompt"`
	// The default value to use if a value is not set.
	Default string `yaml:"default"`
	// If true, a value must be set. If false, the default value will be used if a value is not set.
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

type FlowFileTemplate struct {
	// A list of template data to be used when rendering the flow executable config file.
	Data TemplateData `yaml:"data"`
	// A list of files to include when copying the template in a new location. The files are copied as-is.
	Artifacts []string `yaml:"artifacts,omitempty"`

	*FlowFile `yaml:",inline"`

	location string
}

func (t *FlowFileTemplate) SetContext(location string) {
	t.location = location
}

func (t *FlowFileTemplate) Location() string {
	return t.location
}

func (t *FlowFileTemplate) Validate() error {
	return t.Data.Validate()
}

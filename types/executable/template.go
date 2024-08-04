package executable

import (
	"errors"
	"fmt"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p executable -o template.gen.go template_schema.yaml

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

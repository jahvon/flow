package executable

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jahvon/tuikit/types"
	"gopkg.in/yaml.v3"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p executable -o template.gen.go template_schema.yaml

type TemplateList []*Template

func (f *Field) Set(value string) {
	f.value = &value
}

func (f *Field) Value() string {
	if f.value == nil {
		return f.Default
	}
	return *f.value
}

func (f *Field) ValidateConfig() error {
	if f.Key == "" {
		return errors.New("field is missing a key")
	}
	if f.Prompt == "" {
		return fmt.Errorf("field %s is missing a prompt", f.Key)
	}
	return nil
}

func (f *Field) ValidateValue() error {
	if f.Value() == "" && f.Required {
		return fmt.Errorf("required field with key %s not set", f.Key)
	}

	if f.Validate != "" {
		r, err := regexp.Compile(f.Validate)
		if err != nil {
			return fmt.Errorf("unable to compile validation regex for field with key %s: %w", f.Key, err)
		}
		if !r.MatchString(f.Value()) {
			return fmt.Errorf("validation (%s) failed for field with key %s", f.Validate, f.Key)
		}
	}
	return nil
}

type FormFields []*Field

func (f FormFields) Set(key, value string) {
	for i, entry := range f {
		if entry.Key == key {
			f[i].Set(value)
			return
		}
	}
}

func (f FormFields) MapInterface() map[string]interface{} {
	data := map[string]interface{}{}
	for _, entry := range f {
		data[entry.Key] = entry.Value()
	}
	return data
}

func (f FormFields) ValidateConfig() error {
	for _, field := range f {
		if err := field.ValidateConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (f FormFields) ValidateValues() error {
	for _, field := range f {
		if err := field.ValidateValue(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Template) SetContext(name, location string) {
	if t == nil {
		return
	}
	t.location = &location
	t.assignedName = &name
	if name == "" {
		fn := filepath.Base(location)
		fn = strings.TrimSuffix(fn, filepath.Ext(fn))
		if filepath.Ext(fn) != "" {
			// remove second extension part
			fn = strings.TrimSuffix(fn, filepath.Ext(fn))
		}
		t.assignedName = &fn
	}
}

func (t *Template) Location() string {
	if t.location == nil {
		return ""
	}
	return *t.location
}

func (t *Template) Name() string {
	return *t.assignedName
}

func (t *Template) ValidateFormConfig() error {
	return t.Form.ValidateConfig()
}

func (t *Template) ValidateFormValues() error {
	return t.Form.ValidateValues()
}

func (t *Template) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template - %w", err)
	}
	return string(yamlBytes), nil
}

func (t *Template) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal template - %w", err)
	}
	return string(jsonBytes), nil
}

func (t *Template) Markdown() string {
	return templateMarkdown(t)
}

func (t TemplateList) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template list - %w", err)
	}
	return string(yamlBytes), nil
}

func (t TemplateList) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal template list - %w", err)
	}
	return string(jsonBytes), nil
}

func (t TemplateList) Singular() string {
	return "template"
}

func (t TemplateList) Plural() string {
	return "templates"
}

func (t TemplateList) Items() []*types.CollectionItem {
	items := make([]*types.CollectionItem, len(t))
	for i, template := range t {
		items[i] = &types.CollectionItem{
			ID:     template.Name(),
			Header: template.Name(),
		}
	}
	return items
}

func (t TemplateList) Find(name string) *Template {
	for _, template := range t {
		if template.Name() == name {
			return template
		}
	}
	return nil
}

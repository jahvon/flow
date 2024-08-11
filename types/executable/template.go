package executable

import (
	"errors"
	"fmt"
	"regexp"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p executable -o template.gen.go template_schema.yaml

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
			return fmt.Errorf("unable to compile validation regex for field with key %s: %v", f.Key, err)
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

func (t *Template) SetContext(location string) {
	*t.location = location
}

func (t *Template) Location() string {
	if t.location == nil {
		return ""
	}
	return *t.location
}

func (t *Template) ValidateFormConfig() error {
	return t.Form.ValidateConfig()
}

func (t *Template) ValidateFormValues() error {
	return t.Form.ValidateValues()
}

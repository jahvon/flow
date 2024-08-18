package templates

import (
	"fmt"

	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/types/executable"
)

func showForm(fields executable.FormFields) error {
	if len(fields) == 0 {
		return nil
	}

	if err := fields.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid form fields: %w", err)
	}
	var inputs []*components.TextInput
	for _, f := range fields {
		inputs = append(inputs, &components.TextInput{
			Key:         f.Key,
			Prompt:      f.Prompt,
			Placeholder: f.Default,
		})
	}
	inputs, err := components.ProcessInputs(io.Theme(), inputs...)
	if err != nil {
		return fmt.Errorf("unable to show form: %w", err)
	}
	for _, input := range inputs {
		fields.Set(input.Key, input.Value())
	}
	if err := fields.ValidateValues(); err != nil {
		return err
	}
	return nil
}

package templates

import (
	"fmt"

	"github.com/jahvon/tuikit/views"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/types/executable"
)

func showForm(ctx *context.Context, fields executable.FormFields) error {
	if len(fields) == 0 {
		return nil
	}
	in := ctx.StdIn()
	out := ctx.StdOut()

	if err := fields.Validate(); err != nil {
		return fmt.Errorf("invalid form fields: %w", err)
	}
	var ff []*views.FormField
	for _, f := range fields {
		ff = append(ff, &views.FormField{
			Key:            f.Key,
			Group:          uint(f.Group),
			Description:    f.Description,
			Default:        f.Default,
			Title:          f.Prompt,
			Placeholder:    f.Default,
			Required:       f.Required,
			ValidationExpr: f.Validate,
		})
	}
	form, err := views.NewForm(io.Theme(ctx.Config.Theme.String()), in, out, ff...)
	if err != nil {
		return fmt.Errorf("encountered form init error: %w", err)
	}
	if err = form.Run(ctx.Ctx); err != nil {
		return fmt.Errorf("encountered form run error: %w", err)
	}
	for _, f := range fields {
		v, ok := form.ValueMap()[f.Key]
		if !ok {
			continue
		}
		f.Set(fmt.Sprintf("%v", v))
	}
	return nil
}

package views

import (
	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io/ui/types"
)

func NewUserConfigView(parent types.ParentView, cfg config.UserConfig, format config.OutputFormat) types.ViewBuilder {
	return NewEntityView(parent, &cfg, format)
}

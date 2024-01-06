package ui

import (
	"github.com/jahvon/flow/config"
)

func NewUserConfigView(app *Application, cfg config.UserConfig, format config.OutputFormat) ViewBuilder {
	return NewEntityView(app, &cfg, format)
}

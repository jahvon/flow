package config

import (
	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

func NewUserConfigView(
	container *components.ContainerView,
	cfg config.UserConfig,
	format components.Format,
) components.TeaModel {
	state := &components.TerminalState{
		Theme:  io.Styles(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewEntityView(state, &cfg, format)
}

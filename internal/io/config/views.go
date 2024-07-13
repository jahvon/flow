package config

import (
	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/types/config"
)

func NewUserConfigView(
	container *components.ContainerView,
	cfg config.Config,
	format components.Format,
) components.TeaModel {
	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewEntityView(state, &cfg, format)
}

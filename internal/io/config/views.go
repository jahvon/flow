package config

import (
	"github.com/jahvon/tuikit"
	"github.com/jahvon/tuikit/types"
	"github.com/jahvon/tuikit/views"

	"github.com/jahvon/flow/types/config"
)

func NewUserConfigView(
	container *tuikit.Container,
	cfg config.Config,
	format types.Format,
) tuikit.View {
	return views.NewEntityView(container.RenderState(), &cfg, format)
}

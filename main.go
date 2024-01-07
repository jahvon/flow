package main

import (
	stdCtx "context"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/ui"
)

var log = io.Log().With().Str("scope", "main").Logger()

func main() {
	var cancel stdCtx.CancelFunc
	ctx := context.NewContext(stdCtx.Background())
	ctx.Ctx, cancel = stdCtx.WithCancel(ctx.Ctx)
	if ctx == nil {
		log.Panic().Msg("failed to initialize context")
	}

	cfg := ctx.UserConfig
	if cfg == nil {
		log.Panic().Msg("failed to load user config")
	}
	if err := cfg.Validate(); err != nil {
		log.Panic().Err(err).Msg("user config validation error")
	}

	if cfg.InteractiveUI {
		app := ui.StartApplication(
			ctx.Ctx,
			cancel,
			ui.WithCurrentWorkspace(cfg.CurrentWorkspace),
			ui.WithCurrentNamespace(cfg.CurrentNamespace),
		)
		ctx.App = app
	}

	cmd.Execute(ctx)

	if cfg.InteractiveUI {
		// Keep the app running until the context is cancelled.
		for range ctx.Ctx.Done() {
			return
		}
	}
}

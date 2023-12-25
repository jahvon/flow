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
	ctx := context.NewContext(stdCtx.Background())
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

	if cfg.UIEnabled {
		app := ui.StartApplication(ctx.CancelFunc)
		ctx.App = app
	}

	cmd.Execute(ctx)

	if cfg.UIEnabled {
		// Keep the app running until the context is cancelled.
		for {
			select {
			case <-ctx.Ctx.Done():
				return
			}
		}
	}
}

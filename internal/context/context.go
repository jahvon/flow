package context

import (
	"context"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
)

type Context struct {
	Ctx              context.Context
	UserConfig       *config.UserConfig
	CurrentWorkspace *config.WorkspaceConfig
	WorkspacesCache  *cache.WorkspaceCache
	ExecutableCache  *cache.ExecutableCache
}

var log = io.Log().With().Str("scope", "context").Logger()

func NewContext(ctx context.Context, syncCache bool) *Context {
	userConfg := file.LoadUserConfig()
	if userConfg == nil {
		return nil
	}
	if err := userConfg.Validate(); err != nil {
		log.Err(err).Msg("user config validation error")
		return nil
	}

	wsConfg, err := file.LoadWorkspaceConfig(
		userConfg.CurrentWorkspace,
		userConfg.Workspaces[userConfg.CurrentWorkspace],
	)
	if err != nil {
		log.Err(err).Msg("workspace config load error")
		return nil
	} else if wsConfg == nil {
		log.Error().Msg("workspace config not found")
		return nil
	}

	workspaceCache := cache.NewWorkspaceCache()
	if workspaceCache == nil {
		log.Error().Msg("workspace cache initialization error")
		return nil
	}
	executableCache := cache.NewExecutableCache()
	if executableCache == nil {
		log.Error().Msg("executable cache initialization error")
		return nil
	}

	if syncCache {
		if err := workspaceCache.Update(); err != nil {
			log.Err(err).Msg("workspace cache update error")
			return nil
		}
		if err := executableCache.Update(); err != nil {
			log.Err(err).Msg("executable cache update error")
			return nil
		}
	}

	return &Context{
		Ctx:              ctx,
		UserConfig:       userConfg,
		CurrentWorkspace: wsConfg,
		WorkspacesCache:  cache.NewWorkspaceCache(),
		ExecutableCache:  cache.NewExecutableCache(),
	}
}

package context

import (
	"context"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/ui"
)

type Context struct {
	Ctx              context.Context
	CancelFunc       context.CancelFunc
	UserConfig       *config.UserConfig
	CurrentWorkspace *config.WorkspaceConfig
	WorkspacesCache  *cache.WorkspaceCache
	ExecutableCache  *cache.ExecutableCache
	App              *ui.Application
}

var log = io.Log().With().Str("scope", "context").Logger()

func NewContext(ctx context.Context) *Context {
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

	ctxx, cancel := context.WithCancel(ctx)
	return &Context{
		Ctx:              ctxx,
		CancelFunc:       cancel,
		UserConfig:       userConfg,
		CurrentWorkspace: wsConfg,
		WorkspacesCache:  cache.NewWorkspaceCache(),
		ExecutableCache:  cache.NewExecutableCache(),
	}
}

func ExpandRef(ctx *Context, ref config.Ref) config.Ref {
	id := ref.GetID()
	ws, ns, name := config.ParseExecutableID(id)
	if ws == "" {
		ws = ctx.UserConfig.CurrentWorkspace
	}
	if ns == "" {
		ns = ctx.UserConfig.CurrentNamespace
	}
	return config.NewRef(config.NewExecutableID(ws, ns, name), ref.GetVerb())
}

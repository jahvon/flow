package context

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/ui"
)

type Context struct {
	Ctx              context.Context
	CancelFunc       context.CancelFunc
	Logger           *io.Logger
	UserConfig       *config.UserConfig
	CurrentWorkspace *config.WorkspaceConfig
	WorkspacesCache  *cache.WorkspaceCache
	ExecutableCache  *cache.ExecutableCache
	App              *ui.Application

	// ProcessTmpDir is the temporary directory for the current process. If set, it will be
	// used to store temporary files all executable runs when the tmpDir value is specified.
	ProcessTmpDir string
}

func NewContext(ctx context.Context) *Context {
	userConfig := file.LoadUserConfig()
	if userConfig == nil {
		panic("user config load error")
	}
	if err := userConfig.Validate(); err != nil {
		panic(fmt.Sprintf("user config validation error\nconfig location: %s\n%s", file.UserConfigPath, err))
	}

	wsConfig, err := file.LoadWorkspaceConfig(
		userConfig.CurrentWorkspace,
		userConfig.Workspaces[userConfig.CurrentWorkspace],
	)
	if err != nil {
		panic(fmt.Sprintf("workspace config load error\n%s", err))
	} else if wsConfig == nil {
		panic(fmt.Sprintf("workspace config not found in current workspace (%s)", userConfig.CurrentWorkspace))
	}

	workspaceCache := cache.NewWorkspaceCache()
	if workspaceCache == nil {
		panic("workspace cache initialization error")
	}
	executableCache := cache.NewExecutableCache()
	if executableCache == nil {
		panic("executable cache initialization error")
	}

	var logger *io.Logger
	if userConfig.Interactive != nil && userConfig.Interactive.Enabled {
		logger = io.NewLogger(io.HumanReadable, true)
	} else {
		logger = io.NewLogger(io.HumanReadable, false)
	}

	ctxx, cancel := context.WithCancel(ctx)
	return &Context{
		Ctx:              ctxx,
		CancelFunc:       cancel,
		UserConfig:       userConfig,
		CurrentWorkspace: wsConfig,
		WorkspacesCache:  cache.NewWorkspaceCache(),
		ExecutableCache:  cache.NewExecutableCache(),
		Logger:           logger,
	}
}

func (ctx *Context) Finalize() {
	if ctx.ProcessTmpDir != "" {
		files, err := filepath.Glob(filepath.Join(ctx.ProcessTmpDir, "*"))
		if err != nil {
			ctx.Logger.Error(err, fmt.Sprintf("unable to list files in temp dir %s", ctx.ProcessTmpDir))
			return
		}
		for _, f := range files {
			err = os.RemoveAll(f)
			if err != nil {
				ctx.Logger.Error(err, fmt.Sprintf("unable to remove file %s", f))
			}
		}
		if err := os.Remove(ctx.ProcessTmpDir); err != nil {
			ctx.Logger.Error(err, fmt.Sprintf("unable to remove temp dir %s", ctx.ProcessTmpDir))
		}
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

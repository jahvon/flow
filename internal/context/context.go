package context

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/styles"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
)

type Context struct {
	Ctx                  context.Context
	CancelFunc           context.CancelFunc
	Logger               io.Logger
	UserConfig           *config.UserConfig
	CurrentWorkspace     *config.WorkspaceConfig
	WorkspacesCache      *cache.WorkspaceCache
	ExecutableCache      *cache.ExecutableCache
	InteractiveContainer *components.ContainerView

	// ProcessTmpDir is the temporary directory for the current process. If set, it will be
	// used to store temporary files all executable runs when the tmpDir value is specified.
	ProcessTmpDir string

	stdOut, stdIn *os.File
}

func NewContext(ctx context.Context, stdIn, stdOut *os.File) *Context {
	userConfig, err := file.LoadUserConfig()
	if err != nil {
		panic(errors.Wrap(err, "user config load error"))
	}

	wsConfig, err := currentWorkspace(userConfig)
	if err != nil {
		panic(errors.Wrap(err, "workspace config load error"))
	} else if wsConfig == nil {
		panic(fmt.Errorf("workspace config not found in current workspace (%s)", userConfig.CurrentWorkspace))
	}

	workspaceCache := cache.NewWorkspaceCache()
	if workspaceCache == nil {
		panic("workspace cache initialization error")
	}
	executableCache := cache.NewExecutableCache()
	if executableCache == nil {
		panic("executable cache initialization error")
	}

	ctxx, cancel := context.WithCancel(ctx)
	theme := styles.EverforestTheme()
	logMode := userConfig.DefaultLogMode
	return &Context{
		Ctx:              ctxx,
		CancelFunc:       cancel,
		UserConfig:       userConfig,
		CurrentWorkspace: wsConfig,
		WorkspacesCache:  cache.NewWorkspaceCache(),
		ExecutableCache:  cache.NewExecutableCache(),
		Logger:           io.NewLogger(stdOut, theme, logMode, file.LogsDir()),
		stdOut:           stdOut,
		stdIn:            stdIn,
	}
}

func (ctx *Context) StdOut() *os.File {
	return ctx.stdOut
}

func (ctx *Context) StdIn() *os.File {
	return ctx.stdIn
}

// SetIO sets the standard input and output for the context
// This function should NOT be used outside of tests! The standard input and output
// should be set when creating the context.
func (ctx *Context) SetIO(stdIn, stdOut *os.File) {
	ctx.stdIn = stdIn
	ctx.stdOut = stdOut
}

func (ctx *Context) Finalize() {
	_ = ctx.stdIn.Close()
	_ = ctx.stdOut.Close()

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
	if err := ctx.Logger.Flush(); err != nil {
		if errors.Is(err, os.ErrClosed) {
			return
		}
		panic(err)
	}
}

func ExpandRef(ctx *Context, ref config.Ref) config.Ref {
	id := ref.GetID()
	ws, ns, name := config.ParseExecutableID(id)
	if ws == "" {
		ws = ctx.CurrentWorkspace.AssignedName()
	}
	if ns == "" {
		ns = ctx.UserConfig.CurrentNamespace
	}
	return config.NewRef(config.NewExecutableID(ws, ns, name), ref.GetVerb())
}

func currentWorkspace(userConfig *config.UserConfig) (*config.WorkspaceConfig, error) {
	var ws, wsPath string
	mode := userConfig.WorkspaceMode
	if mode == "" {
		mode = config.WorkspaceModeDynamic
	}

	switch mode {
	case config.WorkspaceModeDynamic:
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		for wsName, path := range userConfig.Workspaces {
			rel, err := filepath.Rel(filepath.Clean(path), filepath.Clean(wd))
			if err != nil {
				return nil, err
			}
			if !strings.HasPrefix(rel, "..") {
				ws = wsName
				wsPath = path
				break
			}
		}
		fallthrough
	case config.WorkspaceModeFixed:
		if ws != "" && wsPath != "" {
			break
		}
		ws = userConfig.CurrentWorkspace
		wsPath = userConfig.Workspaces[ws]
	}
	if ws == "" || wsPath == "" {
		return nil, fmt.Errorf("current workspace not found")
	}

	return file.LoadWorkspaceConfig(ws, wsPath)
}

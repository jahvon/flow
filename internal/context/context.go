package context

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jahvon/tuikit"
	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/filesystem"
	flowIO "github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/types/config"
	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

const (
	AppName      = "flow"
	HeaderCtxKey = "ctx"
)

type Context struct {
	Ctx              context.Context
	CancelFunc       context.CancelFunc
	Logger           io.Logger
	Config           *config.Config
	CurrentWorkspace *workspace.Workspace
	TUIContainer     *tuikit.Container
	WorkspacesCache  cache.WorkspaceCache
	ExecutableCache  cache.ExecutableCache

	// ProcessTmpDir is the temporary directory for the current process. If set, it will be
	// used to store temporary files all executable runs when the tmpDir value is specified.
	ProcessTmpDir string

	stdOut, stdIn *os.File
}

func NewContext(ctx context.Context, stdIn, stdOut *os.File) *Context {
	cfg, err := filesystem.LoadConfig()
	if err != nil {
		panic(errors.Wrap(err, "user config load error"))
	}

	cfg.SetDefaults()
	wsConfig, err := currentWorkspace(cfg)
	if err != nil {
		panic(errors.Wrap(err, "workspace config load error"))
	} else if wsConfig == nil {
		panic(fmt.Errorf("workspace config not found in current workspace (%s)", cfg.CurrentWorkspace))
	}

	workspaceCache := cache.NewWorkspaceCache()
	if workspaceCache == nil {
		panic("workspace cache initialization error")
	}
	executableCache := cache.NewExecutableCache(workspaceCache)
	if executableCache == nil {
		panic("executable cache initialization error")
	}

	ctxx, cancel := context.WithCancel(ctx)
	logMode := cfg.DefaultLogMode
	return &Context{
		Ctx:              ctxx,
		CancelFunc:       cancel,
		Config:           cfg,
		CurrentWorkspace: wsConfig,
		WorkspacesCache:  workspaceCache,
		ExecutableCache:  executableCache,
		Logger:           io.NewLogger(stdOut, flowIO.Theme(), logMode, filesystem.LogsDir()),
		stdOut:           stdOut,
		stdIn:            stdIn,
	}
}

func (ctx *Context) String() string {
	ws := ctx.CurrentWorkspace.AssignedName()
	ns := ctx.Config.CurrentNamespace
	if ws == "" {
		ws = "unk"
	}
	if ns == "" {
		ns = "*"
	}
	return fmt.Sprintf("%s/%s", ws, ns)
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

func (ctx *Context) SetView(view tuikit.View) error {
	if ctx.TUIContainer == nil {
		ctx.TUIContainer = newContainer(ctx)
	}
	return ctx.TUIContainer.SetView(view)
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

func newContainer(ctx *Context) *tuikit.Container {
	return tuikit.NewContainer(ctx.Ctx, ctx.StdIn(), ctx.StdOut(), flowIO.Theme()).
		WithAppName(AppName).
		WithHeaderContext(HeaderCtxKey, ctx.String()).
		WithLoadingMsg("thinking...")
}

func ExpandRef(ctx *Context, ref executable.Ref) executable.Ref {
	id := ref.GetID()
	ws, ns, name := executable.ParseExecutableID(id)
	if ws == "" {
		ws = ctx.CurrentWorkspace.AssignedName()
	}
	if ns == "" {
		ns = ctx.Config.CurrentNamespace
	}
	return executable.NewRef(executable.NewExecutableID(ws, ns, name), ref.GetVerb())
}

func currentWorkspace(cfg *config.Config) (*workspace.Workspace, error) {
	var ws, wsPath string
	mode := cfg.WorkspaceMode

	switch mode {
	case config.ConfigWorkspaceModeDynamic:
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		for wsName, path := range cfg.Workspaces {
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
	case config.ConfigWorkspaceModeFixed:
		if ws != "" && wsPath != "" {
			break
		}
		ws = cfg.CurrentWorkspace
		wsPath = cfg.Workspaces[ws]
	}
	if ws == "" || wsPath == "" {
		return nil, fmt.Errorf("current workspace not found")
	}

	return filesystem.LoadWorkspaceConfig(ws, wsPath)
}

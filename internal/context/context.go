package context

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jahvon/tuikit"
	"github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/styles"
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
	if cfg.DefaultTimeout != 0 && os.Getenv(executable.TimeoutOverrideEnv) == "" {
		// HACK: Set the default timeout as an environment variable to be used by the exec runner
		// This is a temporary solution until the config handling is refactored a bit
		_ = os.Setenv(executable.TimeoutOverrideEnv, cfg.DefaultTimeout.String())
	}
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
	c := &Context{
		Ctx:              ctxx,
		CancelFunc:       cancel,
		Config:           cfg,
		CurrentWorkspace: wsConfig,
		WorkspacesCache:  workspaceCache,
		ExecutableCache:  executableCache,
		Logger:           io.NewLogger(stdOut, flowIO.Theme(cfg.Theme.String()), logMode, filesystem.LogsDir()),
		stdOut:           stdOut,
		stdIn:            stdIn,
	}

	app := tuikit.NewApplication(
		AppName,
		tuikit.WithState(HeaderCtxKey, c.String()),
		tuikit.WithLoadingMsg("thinking..."),
	)

	theme := flowIO.Theme(cfg.Theme.String())
	if cfg.ColorOverride != nil {
		theme = overrideThemeColor(theme, cfg.ColorOverride)
	}
	c.TUIContainer, err = tuikit.NewContainer(
		ctx, app,
		tuikit.WithInput(stdIn),
		tuikit.WithOutput(stdOut),
		tuikit.WithTheme(theme),
	)
	if err != nil {
		panic(errors.Wrap(err, "TUI container initialization error"))
	}
	return c
}

func (ctx *Context) String() string {
	ws := ctx.CurrentWorkspace.AssignedName()
	ns := ctx.Config.CurrentNamespace
	if ws == "" {
		ws = "unk"
	}
	if ns == "" {
		ns = executable.WildcardNamespace
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

func ExpandRef(ctx *Context, ref executable.Ref) executable.Ref {
	id := ref.ID()
	ws, ns, name := executable.MustParseExecutableID(id)
	if ws == "" || ws == executable.WildcardWorkspace {
		ws = ctx.CurrentWorkspace.AssignedName()
	}
	if ns == "" {
		ns = ctx.Config.CurrentNamespace
	}
	return executable.NewRef(executable.NewExecutableID(ws, ns, name), ref.Verb())
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
		if runtime.GOOS == "darwin" {
			// On macOS, paths that start with /tmp (and some other system directories)
			// are actually symbolic links to paths under /private. The OS may return
			// either form of the path - e.g., both "/tmp/file" and "/private/tmp/file"
			// refer to the same location. We strip the "/private" prefix for consistent
			// path comparison, while preserving the original paths for filesystem operations.
			wd = strings.TrimPrefix(wd, "/private")
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

func overrideThemeColor(theme styles.Theme, palette *config.ColorPalette) styles.Theme {
	if palette == nil {
		return theme
	}
	if palette.Primary != nil {
		theme.PrimaryColor = lipgloss.Color(*palette.Primary)
	}
	if palette.Secondary != nil {
		theme.SecondaryColor = lipgloss.Color(*palette.Secondary)
	}
	if palette.Tertiary != nil {
		theme.TertiaryColor = lipgloss.Color(*palette.Tertiary)
	}
	if palette.Success != nil {
		theme.SuccessColor = lipgloss.Color(*palette.Success)
	}
	if palette.Warning != nil {
		theme.WarningColor = lipgloss.Color(*palette.Warning)
	}
	if palette.Error != nil {
		theme.ErrorColor = lipgloss.Color(*palette.Error)
	}
	if palette.Info != nil {
		theme.InfoColor = lipgloss.Color(*palette.Info)
	}
	if palette.Body != nil {
		theme.BodyColor = lipgloss.Color(*palette.Body)
	}
	if palette.Emphasis != nil {
		theme.EmphasisColor = lipgloss.Color(*palette.Emphasis)
	}
	if palette.White != nil {
		theme.White = lipgloss.Color(*palette.White)
	}
	if palette.Black != nil {
		theme.Black = lipgloss.Color(*palette.Black)
	}
	if palette.Gray != nil {
		theme.Gray = lipgloss.Color(*palette.Gray)
	}
	if palette.CodeStyle != nil {
		theme.ChromaCodeStyle = *palette.CodeStyle
	}
	return theme
}

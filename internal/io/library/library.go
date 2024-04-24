package library

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

type Library struct {
	ctx                   *context.Context
	termWidth, termHeight int

	allWorkspaces  config.WorkspaceConfigList
	allExecutables config.ExecutableList
	curWsConfig    *config.WorkspaceConfig

	filter             Filter
	showNamespaces     bool
	visibleWorkspaces  []string
	visibleNamespaces  []string
	visibleExecutables config.ExecutableList

	headerModel   components.Header
	loadingScreen tea.Model
	infoText      string
	showHelp      bool

	paneZeroViewport, paneOneViewport, paneTwoViewport                 viewport.Model
	currentPane, currentWorkspace, currentNamespace, currentExecutable uint
}

type Filter struct {
	Workspace, Namespace string
	Verb                 config.Verb
	Tags                 config.Tags
}

func NewLibrary(
	ctx *context.Context,
	workspaces config.WorkspaceConfigList,
	execs config.ExecutableList,
	filter Filter,
) *Library {
	p1 := viewport.New(0, 0)
	p2 := viewport.New(0, 0)
	p3 := viewport.New(0, 0)
	headerModel := components.Header{
		Styles: styles,
		Name:   "flow",
		CtxKey: "ctx",
		CtxVal: ctxVal(ctx.UserConfig.CurrentWorkspace, ctx.UserConfig.CurrentNamespace),
	}
	loadingModel := components.NewLoadingView("Loading...", styles)

	return &Library{
		ctx:              ctx,
		allWorkspaces:    workspaces,
		allExecutables:   execs,
		filter:           filter,
		headerModel:      headerModel,
		loadingScreen:    loadingModel,
		paneZeroViewport: p1,
		paneOneViewport:  p2,
		paneTwoViewport:  p3,
	}
}

func ctxVal(ws, ns string) string {
	if ws == "" {
		ws = "unk"
	}
	if ns == "" {
		ns = "*"
	}
	return fmt.Sprintf("%s/%s", ws, ns)
}

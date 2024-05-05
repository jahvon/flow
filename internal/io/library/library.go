package library

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

const (
	appName = "flow library"
)

type Library struct {
	ctx                      *context.Context
	termWidth, termHeight    int
	noticeText               string
	showHelp, showNamespaces bool

	visibleWorkspaces  []string
	visibleNamespaces  []string
	visibleExecutables config.ExecutableList
	allWorkspaces      config.WorkspaceConfigList
	allExecutables     config.ExecutableList
	filter             Filter
	theme              styles.Theme
	selectedWsConfig   *config.WorkspaceConfig

	currentPane, currentWorkspace, currentNamespace, currentExecutable uint
	paneZeroViewport, paneOneViewport, paneTwoViewport                 viewport.Model
	loadingScreen                                                      tea.Model
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
	theme styles.Theme,
) *Library {
	p1 := viewport.New(0, 0)
	p2 := viewport.New(0, 0)
	p3 := viewport.New(0, 0)
	loadingModel := components.NewLoadingView("Loading...", theme)

	return &Library{
		ctx:              ctx,
		allWorkspaces:    workspaces,
		allExecutables:   execs,
		filter:           filter,
		loadingScreen:    loadingModel,
		paneZeroViewport: p1,
		paneOneViewport:  p2,
		paneTwoViewport:  p3,
		theme:            theme,
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

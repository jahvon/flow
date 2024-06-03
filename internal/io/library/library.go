package library

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

const (
	appName = "flow library"
)

type Library struct {
	ctx                                 *context.Context
	termWidth, termHeight               int
	noticeText                          string
	showHelp, showNamespaces, splitView bool

	visibleWorkspaces  []string
	visibleNamespaces  []string
	visibleExecutables config.ExecutableList
	allWorkspaces      config.WorkspaceConfigList
	allExecutables     config.ExecutableList
	filter             Filter
	theme              styles.Theme
	selectedWsConfig   *config.WorkspaceConfig

	currentPane, currentWorkspace, currentNamespace, currentExecutable uint
	currentFormat, currentHelpPage                                     uint
	paneZeroViewport, paneOneViewport, paneTwoViewport                 viewport.Model

	cmdRunFunc func(string) error
}

type Filter struct {
	Workspace, Namespace string
	Verb                 config.Verb
	Tags                 config.Tags
	Substring            string
}

func NewLibrary(
	ctx *context.Context,
	workspaces config.WorkspaceConfigList,
	execs config.ExecutableList,
	filter Filter,
	theme styles.Theme,
	runFunc func(string) error,
) *Library {
	p1 := viewport.New(0, 0)
	p2 := viewport.New(0, 0)
	p3 := viewport.New(0, 0)
	return &Library{
		ctx:              ctx,
		allWorkspaces:    workspaces,
		allExecutables:   execs,
		filter:           filter,
		paneZeroViewport: p1,
		paneOneViewport:  p2,
		paneTwoViewport:  p3,
		theme:            theme,
		cmdRunFunc:       runFunc,
	}
}

func NewLibraryView(
	ctx *context.Context,
	workspaces config.WorkspaceConfigList,
	execs config.ExecutableList,
	filter Filter,
	theme styles.Theme,
	runFunc func(string) error,
) components.TeaModel {
	l := NewLibrary(ctx, workspaces, execs, filter, theme, runFunc)
	return components.NewFrameView(l)
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

var formatOrder = []string{"markdown", "yaml", "json"}

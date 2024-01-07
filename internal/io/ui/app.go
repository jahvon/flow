package ui

import (
	"context"
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jahvon/flow/config"
)

const (
	syncAppMsg    = "sync-app"
	heightPadding = 5
)

type Application struct {
	ctx         context.Context
	pendingView ViewBuilder
	activeView  ViewBuilder
	lastView    ViewBuilder
	program     *tea.Program

	curWs, curNs, curNotice string
	curNoticeLvl            NoticeLevel
	curFilter               config.Tags

	width, height int
	ready         bool
}

func StartApplication(ctx context.Context, cancel context.CancelFunc, opts ...HeaderOption) *Application {
	activeView := NewLoadingView("")
	a := &Application{
		ctx:          ctx,
		curWs:        "unk",
		curNs:        "unk",
		curNotice:    "",
		curNoticeLvl: NoticeLevelInfo,
		activeView:   activeView,
	}
	a.SetHeader(opts...)
	prgm := tea.NewProgram(a, tea.WithContext(ctx))
	go func() {
		var err error
		if _, err = prgm.Run(); err != nil {
			log.Panic().Err(err).Msg("error running application")
		}
		log.Debug().Msg("application exited")
		cancel()
	}()
	a.program = prgm
	return a
}

func (a *Application) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	if a.activeView != nil {
		cmds = append(cmds, a.activeView.Init())
	}
	cmds = append(cmds, tea.SetWindowTitle("flow"))
	return tea.Batch(cmds...)
}

func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.activeView.Update(tea.Quit())
			return a, tea.Quit
		case "esc", "backspace":
			if a.lastView == nil || a.activeView == a.lastView {
				a.activeView.Update(tea.Quit())
				return a, tea.Quit
			} else {
				a.activeView = a.lastView
				a.lastView = nil
				return a, nil
			}
		}
	case string:
		if msg == syncAppMsg {
			log.Trace().Msg("sync triggered")
		}
	case tea.WindowSizeMsg:
		if !a.Ready() {
			a.width = int(math.Floor(float64(msg.Width) * 0.75))
			a.height = msg.Height - heightPadding
			if a.pendingView != nil {
				a.activeView = a.pendingView
				a.pendingView = nil
			}
			a.ready = true
		} else {
			a.width = int(math.Floor(float64(msg.Width) * 0.75))
			a.height = msg.Height - heightPadding
		}
	}

	_, cmd := a.activeView.Update(msg)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

func (a *Application) Ready() bool {
	return a.ready
}

func (a *Application) Height() int {
	return a.height
}

func (a *Application) Width() int {
	return a.width
}

func (a *Application) View() string {
	header := fmt.Sprintf(
		"%s%s%s",
		brandStyle("flow"),
		ContextStr(a.curWs, a.curNs),
		FilterStr(a.curFilter),
	)
	var help string
	if a.activeView.FooterEnabled() && a.activeView.HelpMsg() != "" {
		help = helpStyle(fmt.Sprintf("\n %s | esc: back • q: quit • ↑/↓: navigate", a.activeView.HelpMsg()))
	} else if a.activeView.FooterEnabled() {
		help = helpStyle("\n esc: back • q: quit • ↑/↓: navigate")
	}
	var notice string
	if a.curNotice != "" {
		notice = fmt.Sprintf("\n %s", NoticeStr(a.curNotice, a.curNoticeLvl))
	}

	if !a.Ready() && a.activeView.Type() != loadingViewType {
		a.activeView = NewLoadingView("")
	}

	return header + "\n" + a.activeView.View() + help + notice
}

func (a *Application) SetHeader(opts ...HeaderOption) {
	for _, opt := range opts {
		opt(a)
	}
}

func (a *Application) BuildAndSetView(viewBuilder ViewBuilder, opts ...HeaderOption) {
	a.SetHeader(opts...)
	if !a.Ready() {
		a.pendingView = viewBuilder
		return
	}
	if a.activeView != nil && a.activeView.Type() != loadingViewType && a.activeView.Type() != viewBuilder.Type() {
		a.lastView = a.activeView
	}
	a.activeView = viewBuilder
	a.Update(syncAppMsg)
}

func (a *Application) HandleUserError(err error) {
	if err == nil {
		return
	}
	a.curNoticeLvl = NoticeLevelWarning
	a.curNotice = err.Error()
	a.Update(syncAppMsg)
}

func (a *Application) HandleInternalError(err error) {
	if err == nil {
		return
	}
	a.curNoticeLvl = NoticeLevelError
	a.curNotice = err.Error()
	a.Update(syncAppMsg)
}

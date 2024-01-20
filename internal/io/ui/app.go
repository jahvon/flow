package ui

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
	"github.com/jahvon/flow/internal/io/ui/views"
)

const (
	heightPadding = 5
)

type Application struct {
	ctx         context.Context
	pendingView types.ViewBuilder
	activeView  types.ViewBuilder
	lastView    types.ViewBuilder
	// nextView    ViewBuilder
	program *tea.Program

	curWs, curNs, curNotice, curError string
	curNoticeLvl                      types.NoticeLevel

	width, height int
	ready         bool
}

func StartApplication(ctx context.Context, cancel context.CancelFunc) *Application {
	activeView := views.NewLoadingView("")
	a := &Application{
		ctx:          ctx,
		curWs:        "unk",
		curNs:        "unk",
		curNotice:    "",
		curNoticeLvl: types.NoticeLevelInfo,
		activeView:   activeView,
	}
	prgm := tea.NewProgram(a, tea.WithContext(ctx))
	go func() {
		var err error
		if _, err = prgm.Run(); err != nil {
			panic(fmt.Errorf("error running application: %w", err))
		}
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
	cmds = append(
		cmds,
		tea.SetWindowTitle("flow"),
		tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return types.TickMsg(t)
		}),
	)
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
			a.width = int(math.Floor(float64(msg.Width) * 0.90))
			a.height = msg.Height - heightPadding
		}
	case types.TickMsg:
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return types.TickMsg(t)
		}))
	case tea.Cmd:
		cmds = append(cmds, msg)
	}

	_, cmd := a.activeView.Update(msg)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

func (a *Application) Ready() bool {
	return a.ready && a.curWs != "unk" && a.curNs != "unk"
}

func (a *Application) Finalize() {
	a.curNotice = ""
}

func (a *Application) Height() int {
	return a.height
}

func (a *Application) Width() int {
	return a.width
}

func (a *Application) View() string {
	header := fmt.Sprintf("%s%s", styles.RenderBrand("flow"), ContextStr(a.curWs, a.curNs))
	header += NoticeStr(a.curNotice, a.curNoticeLvl)

	if a.curError != "" {
		a.ready = true
		padded := lipgloss.Style{}.MarginLeft(2).Render
		return header + "\n\n" + padded(styles.RenderError(a.curError)) + "\n\n"
	}

	var help string
	var lastViewHelp string
	if a.lastView != nil {
		lastViewHelp = "esc: back • "
	}
	if a.activeView.FooterEnabled() && a.activeView.HelpMsg() != "" {
		help = styles.RenderHelp(fmt.Sprintf("\n %s | %sq: quit • ↑/↓: navigate", a.activeView.HelpMsg(), lastViewHelp))
	} else if a.activeView.FooterEnabled() {
		help = styles.RenderHelp(fmt.Sprintf("\n %sq: quit • ↑/↓: navigate", lastViewHelp))
	}
	if !a.Ready() && a.activeView.Type() != views.LoadingViewType {
		a.activeView = views.NewLoadingView("")
	}

	return header + "\n" + a.activeView.View() + help
}

func (a *Application) SetContext(ws, ns string) {
	if ws != "" {
		a.curWs = ws
	}
	if ns != "" {
		a.curNs = ns
	} else {
		a.curNs = "*"
	}
}

func (a *Application) GetContext() (string, string) {
	return a.curWs, a.curNs
}

func (a *Application) SetNotice(notice string, lvl types.NoticeLevel) {
	a.curNotice = notice
	a.curNoticeLvl = lvl
}

func (a *Application) BuildAndSetView(viewBuilder types.ViewBuilder) {
	if !a.Ready() {
		a.pendingView = viewBuilder
		return
	}
	if a.activeView != nil && a.activeView.Type() != views.LoadingViewType && a.activeView.Type() != viewBuilder.Type() {
		a.lastView = a.activeView
	}
	a.activeView = viewBuilder
	cmd := a.activeView.Init()
	a.Update(cmd)
}

func (a *Application) HandleUserError(err error) {
	if err == nil {
		return
	}
	a.curError = fmt.Errorf("!! user error !!\n%w", err).Error()
	a.ready = false // not ready until rendered
}

func (a *Application) HandleInternalError(err error) {
	if err == nil {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if strings.Contains(file, "common") {
		maxDepth := 5
		curDepth := 2
		for {
			if strings.Contains(file, "common") {
				_, file, line, ok = runtime.Caller(curDepth)
				curDepth++
			} else if curDepth < maxDepth {
				break
			}
		}
	}
	if ok {
		a.curError = fmt.Errorf("!! encountered error !!\n%w\n%s:%d", err, file, line).Error()
	} else {
		a.curError = fmt.Errorf("!! encountered error !!\n%w", err).Error()
	}
	a.ready = false // not ready until rendered
}

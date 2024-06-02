//nolint:funlen,gocritic,gocognit,gocyclo,cyclop
package library

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jahvon/tuikit/styles"
	"golang.design/x/clipboard"

	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/internal/services/open"
)

func (l *Library) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case tea.KeyLeft.String():
			if l.currentPane == 0 {
				break
			}
			l.currentPane--

			// Reset the current executable when switching back to the workspaces pane
			if l.currentPane == 0 {
				l.currentExecutable = 0
				l.paneOneViewport.GotoTop()
			}
		case tea.KeyRight.String():
			if l.currentPane == 2 {
				break
			}
			l.currentPane++
		case "h":
			l.showHelp = !l.showHelp
		}
	}

	wsPane, wsCmd := l.updateWsPane(msg)
	l.paneZeroViewport = wsPane
	execPane, execCmd := l.updateExecPanes(msg)
	if l.currentPane == 1 {
		l.paneOneViewport = execPane
	} else if l.currentPane == 2 {
		l.paneTwoViewport = execPane
	}

	l.setVisibleWorkspaces()
	l.setVisibleNamespaces()
	l.setVisibleExecs()

	cmds = append(cmds, wsCmd, execCmd)
	return l, tea.Batch(cmds...)
}

func (l *Library) updateWsPane(msg tea.Msg) (viewport.Model, tea.Cmd) {
	if l.currentPane != 0 {
		return l.paneZeroViewport, nil
	}

	numWs := len(l.visibleWorkspaces)
	numNs := len(l.visibleNamespaces)
	if numWs == 0 {
		return l.paneZeroViewport, nil
	}

	curWs := l.visibleWorkspaces[l.currentWorkspace]
	curWsCfg := l.selectedWsConfig
	wsCanMoveUp := numWs > 1 && l.currentWorkspace >= 1 && l.currentWorkspace < uint(numWs)
	wsCanMoveDown := numWs > 1 && l.currentWorkspace < uint(numWs-1)

	var curNs string
	if len(l.visibleNamespaces) > 0 {
		curNs = l.visibleNamespaces[l.currentNamespace]
	}
	nsCanMoveUp := curNs != "" && numNs > 1 && l.currentNamespace >= 1 && l.currentNamespace < uint(numNs)
	nsCanMoveDown := curNs != "" && numNs > 1 && l.currentNamespace < uint(numNs-1)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case tea.KeyDown.String():
			if l.showNamespaces && nsCanMoveDown {
				l.currentNamespace++
			} else if !l.showNamespaces && wsCanMoveDown {
				l.currentWorkspace++
			}
		case tea.KeyUp.String():
			if l.showNamespaces && nsCanMoveUp {
				l.currentNamespace--
			} else if !l.showNamespaces && wsCanMoveUp {
				l.currentWorkspace--
			}
		case tea.KeySpace.String():
			if numNs > 0 {
				l.showNamespaces = !l.showNamespaces
				l.currentNamespace = 0
			}
		case "o":
			if curWsCfg == nil {
				l.SetNotice("no workspace selected", styles.NoticeLevelError)
				break
			}

			if err := open.Open(curWsCfg.Location(), false); err != nil {
				l.ctx.Logger.Error(err, "unable to open workspace")
				l.SetNotice("unable to open workspace", styles.NoticeLevelError)
			}
		case "e":
			if curWsCfg == nil {
				l.SetNotice("no workspace selected", styles.NoticeLevelError)
				break
			}

			if err := common.OpenInEditor(
				filepath.Join(curWsCfg.Location(), file.WorkspaceConfigFileName),
				l.ctx.StdIn(), l.ctx.StdOut(),
			); err != nil {
				l.ctx.Logger.Error(err, "unable to open workspace in editor")
				l.SetNotice("unable to open workspace in editor", styles.NoticeLevelError)
			}
		case "s":
			if curWsCfg == nil {
				l.SetNotice("no workspace selected", styles.NoticeLevelError)
				break
			}

			curCfg, err := file.LoadUserConfig()
			if err != nil {
				l.ctx.Logger.Error(err, "unable to load user config")
				l.SetNotice("unable to load user config", styles.NoticeLevelError)
				break
			}

			switch {
			case l.showNamespaces && curNs == withoutNamespaceLabel:
				curCfg.CurrentNamespace = ""
			case l.showNamespaces && curNs == allNamespacesLabel:
				l.SetNotice("no namespace selected", styles.NoticeLevelError)
			case l.showNamespaces && curNs != "":
				curCfg.CurrentNamespace = curNs
			case !l.showNamespaces && curWs == allWorkspacesLabel:
				l.SetNotice("no workspace selected", styles.NoticeLevelError)
			case !l.showNamespaces && curWs != "":
				if curWs != curWsCfg.AssignedName() {
					l.SetNotice("current workspace out of sync", styles.NoticeLevelError)
					break
				}
				curCfg.CurrentWorkspace = curWsCfg.AssignedName()
			}

			if err := file.WriteUserConfig(curCfg); err != nil {
				l.ctx.Logger.Error(err, "unable to write user config")
				l.SetNotice("unable to write user config", styles.NoticeLevelError)
				break
			}

			l.ctx.UserConfig.CurrentWorkspace = curCfg.CurrentWorkspace
			l.ctx.UserConfig.CurrentNamespace = curCfg.CurrentNamespace
			l.SetNotice("context updated", styles.NoticeLevelInfo)
		}
	}

	return l.paneZeroViewport.Update(msg)
}

func (l *Library) updateExecPanes(msg tea.Msg) (viewport.Model, tea.Cmd) {
	if l.currentPane != 1 && l.currentPane != 2 {
		return l.paneOneViewport, nil
	}

	var pane viewport.Model
	if l.currentPane == 1 {
		pane = l.paneOneViewport
	} else if l.currentPane == 2 {
		pane = l.paneTwoViewport
	}

	numExecs := len(l.visibleExecutables)
	if numExecs == 0 {
		return pane, nil
	}

	curExec := l.visibleExecutables[l.currentExecutable]
	canMoveUp := numExecs > 1 && l.currentExecutable >= 1 && l.currentExecutable < uint(numExecs)
	canMoveDown := numExecs > 1 && l.currentExecutable < uint(numExecs-1)

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case tea.KeyDown.String():
			if l.currentPane == 1 && canMoveDown {
				l.currentExecutable++
			}
		case tea.KeyUp.String():
			if l.currentPane == 1 && canMoveUp {
				l.currentExecutable--
			}
		case "e":
			if curExec == nil {
				l.SetNotice("no executable selected", styles.NoticeLevelError)
				break
			}

			if err := common.OpenInEditor(curExec.DefinitionPath(), l.ctx.StdIn(), l.ctx.StdOut()); err != nil {
				l.ctx.Logger.Error(err, "unable to open executable in editor")
				l.SetNotice("unable to open executable in editor", styles.NoticeLevelError)
			}
		case "c":
			if curExec == nil {
				l.SetNotice("no executable selected", styles.NoticeLevelError)
				break
			}

			if err := clipboard.Init(); err != nil {
				l.ctx.Logger.Error(err, "unable to initialize clipboard")
				l.SetNotice("unable to initialize clipboard", styles.NoticeLevelError)
				break
			}

			clipboard.Write(clipboard.FmtText, []byte(curExec.Ref().String()))
			l.SetNotice("copied reference to clipboard", styles.NoticeLevelInfo)
		case "r":
			if curExec == nil {
				l.SetNotice("no executable selected", styles.NoticeLevelError)
				break
			}

			go func() {
				l.ctx.InteractiveContainer.Shutdown()
				if err := l.cmdRunFunc(curExec.Ref().String()); err != nil {
					l.ctx.Logger.Fatalx("unable to execute command", "error", err)
				}
			}()
		}
	}

	return pane.Update(msg)
}

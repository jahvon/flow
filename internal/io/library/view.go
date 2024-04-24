package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/config"
)

const (
	// heightPadding is used when determining the height of the panes.
	// It is used to account for the header and footer
	heightPadding = 4
	// widthPadding is used when determining the width of the panes.
	// It is used to account for the padding on the left and right of the panes
	widthPadding = 4

	allWorkspacesLabel    = "all workspaces"
	withoutNamespaceLabel = "w/o namespace"
	allNamespacesLabel    = "all namespaces"

	footerPrefix        = "[ q/ctrl+c: quit] [ h: help ]"
	helpPrefix          = "[ ↑/↓: navigate pane ] [ ←/→: change pane ]"
	paneOneHelp         = "[ o: open ] [ e: edit ] [ s:set ] ● [ space: show namespaces ]"
	paneOneExpandedHelp = "[ o: open ] [ e: edit ] [ s:set ] ● [ space: hide namespaces ]"
	paneTwoAndThreeHelp = "[ e: edit ] [ c: copy ref ]  ● [ f: change format ]"
)

func (l *Library) View() string {
	if l.loadingScreen != nil {
		return l.loadingScreen.View()
	}

	l.paneZeroViewport.Style = paneStyle(l.currentPane == 0, 0)
	l.paneZeroViewport.SetContent(l.paneZeroContent())

	l.paneOneViewport.Style = paneStyle(l.currentPane == 1, 1)
	l.paneOneViewport.SetContent(l.paneOneContent())

	l.paneTwoViewport.Style = paneStyle(l.currentPane == 2, 2)
	l.paneTwoViewport.SetContent(l.paneTwoContent())
	header := renderHeader(ctxVal(l.ctx.UserConfig.CurrentWorkspace, l.ctx.UserConfig.CurrentNamespace))
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		l.paneZeroViewport.View(),
		l.paneOneViewport.View(),
		l.paneTwoViewport.View(),
	)
	footer := l.footerContent()

	return lipgloss.JoinVertical(lipgloss.Top, header, panes, footer)

}

func (l *Library) SetNotice(notice string, level components.NoticeLevel) {
	if level == "" {
		level = components.NoticeLevelInfo
	}
	l.headerModel.Notice = notice
	l.headerModel.NoticeLevel = level

}

func (l *Library) paneZeroContent() string {
	var sb strings.Builder
	workspaces := l.visibleWorkspaces
	namespaces := l.visibleNamespaces
	sb.WriteString(renderPaneTitle("Workspaces", len(workspaces), l.currentPane == 0))

	numWs := len(workspaces)
	numNs := len(namespaces)
	if numWs == 0 {
		sb.WriteString(styles.RenderError("No workspaces found"))
		return sb.String()
	}
	paneWidth, _, _ := calculateViewportWidths(l.termWidth)

	for i, ws := range workspaces {
		prefix := "◌ "
		if uint(i) == l.currentWorkspace && !l.showNamespaces {
			prefix = "● "
		} else if uint(i) == l.currentWorkspace && l.showNamespaces {
			prefix = "◉ "
		}

		if uint(i) == l.currentWorkspace {
			sb.WriteString(renderSelection(prefix + truncateText(ws, paneWidth)))
			sb.WriteString("\n")
			if numNs == 1 {
				sb.WriteString(renderDescription(fmt.Sprintf("  %d namespace", numNs)))
			} else {
				sb.WriteString(renderDescription(fmt.Sprintf("  %d namespaces", numNs)))
			}
			sb.WriteString("\n")

			if l.showNamespaces {
				for j, ns := range namespaces {
					if uint(j) == l.currentNamespace {
						sb.WriteString(renderSecondarySelection("  > " + truncateText(ns, paneWidth)))
					} else {
						sb.WriteString(renderInactive("    " + truncateText(ns, paneWidth)))
					}
					sb.WriteString("\n")
				}
			}
		} else {
			sb.WriteString(renderInactive(prefix + truncateText(ws, paneWidth)))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (l *Library) paneOneContent() string {
	var sb strings.Builder
	sb.WriteString(renderPaneTitle("Executables", len(l.visibleExecutables), l.currentPane == 1))
	if len(l.visibleExecutables) == 0 {
		sb.WriteString(styles.RenderError("No executables found"))
		return sb.String()
	}

	_, paneWidth, _ := calculateViewportWidths(l.termWidth)

	for i, ex := range l.visibleExecutables {
		if uint(i) == l.currentExecutable {
			sb.WriteString(renderSelection("* " + truncateText(ex.Ref().GetID(), paneWidth)))
		} else {
			sb.WriteString(renderInactive("  " + truncateText(ex.Ref().GetID(), paneWidth)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (l *Library) paneTwoContent() string {
	if len(l.visibleExecutables) == 0 {
		return ""
	}

	_, _, maxWidth := calculateViewportWidths(l.termWidth)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(styles.MarkdownStyleJSON)),
		glamour.WithWordWrap(maxWidth-2),
	)
	if err != nil {
		l.ctx.Logger.Error(err, "unable to create term renderer")
		return styles.RenderError("unable to render markdown")
	}

	ex := l.visibleExecutables[l.currentExecutable]
	content := ex.Markdown()
	viewStr, err := renderer.Render(content)
	if err != nil {
		l.ctx.Logger.Error(err, "unable to render markdown")
		return styles.RenderError("unable to render markdown")
	}

	return viewStr
}

func (l *Library) footerContent() string {
	help := l.showHelp
	switch l.currentPane {
	case 0:
		if help && l.showNamespaces {
			return renderFooter(fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, paneOneExpandedHelp))
		} else if help {
			return renderFooter(fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, paneOneHelp))
		} else if l.currentWorkspace < uint(len(l.visibleWorkspaces)) {
			ws := l.visibleWorkspaces[l.currentWorkspace]
			if ws == allWorkspacesLabel {
				break
			}
			var wsCfg *config.WorkspaceConfig
			for _, w := range l.allWorkspaces {
				if w.AssignedName() == ws {
					wsCfg = &w
				}
			}
			if wsCfg == nil {
				l.ctx.Logger.Errorf("unable to find workspace config for %s", ws)
				break
			}

			path, err := relativePathFromWd(wsCfg.Location())
			if err != nil {
				l.ctx.Logger.Error(err, "unable to get relative path from wd")
				break
			}
			var info string
			if len(wsCfg.Tags) > 0 {
				info = fmt.Sprintf("%s(%s) -> %s", wsCfg.DisplayName, wsCfg.Tags.PreviewString(), path)
			} else {
				info = fmt.Sprintf("%s -> %s", wsCfg.DisplayName, path)
			}
			return renderFooter(fmt.Sprintf("%s ● %s", footerPrefix, info))
		}
	case 1, 2:
		if help {
			return renderFooter(fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, paneTwoAndThreeHelp))
		} else if l.currentExecutable < uint(len(l.visibleExecutables)) {
			exec := l.visibleExecutables[l.currentExecutable]
			path, err := relativePathFromWd(exec.DefinitionPath())
			if err != nil {
				l.ctx.Logger.Error(err, "unable to get relative path from wd")
				break
			}
			return renderFooter(fmt.Sprintf("%s ● %s", footerPrefix, path))
		}
	}
	return renderFooter(footerPrefix)
}

func relativePathFromWd(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Rel(wd, path)
}

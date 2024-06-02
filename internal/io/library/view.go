//nolint:gocognit,gocritic,nestif
package library

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/jahvon/tuikit/styles"

	"github.com/jahvon/flow/config"
)

const (
	// widthPadding is used when determining the width of the panes.
	// It is used to account for the spacing on left/right spacing of the panes
	widthPadding = 4

	allWorkspacesLabel    = "all workspaces"
	withoutNamespaceLabel = "w/o namespace"
	allNamespacesLabel    = "all namespaces"

	footerPrefix        = "[ q/ctrl+c: quit] [ h: help ]"
	helpPrefix          = "[ ↑/↓: navigate pane ] [ ←/→: change pane ]"
	paneOneHelp         = "[ o: open ] [ e: edit ] [ s:set ] ● [ space: show namespaces ]"
	paneOneExpandedHelp = "[ o: open ] [ e: edit ] [ s:set ] ● [ space: hide namespaces ]"
	paneTwoHelp         = "[ r: run ] [ e: edit ] [ c: copy ref ]  ● [ f: change format ]"
	paneThreeHelp       = "[ r: run ] [ e: edit ] [ c: copy ref ]  ● [ f: change format ]"
)

var (
	// heightPadding is used when determining the height of the panes.
	// It is used to account for the header and footer
	heightPadding = styles.HeaderHeight + styles.FooterHeight
)

func (l *Library) View() string {
	l.paneZeroViewport.Style = paneStyle(0, l.theme)
	l.paneZeroViewport.SetContent(l.paneZeroContent())

	l.paneOneViewport.Style = paneStyle(1, l.theme)
	l.paneOneViewport.SetContent(l.paneOneContent())

	l.paneTwoViewport.Style = paneStyle(2, l.theme)
	l.paneTwoViewport.SetContent(l.paneTwoContent())
	v := ctxVal(l.ctx.UserConfig.CurrentWorkspace, l.ctx.UserConfig.CurrentNamespace)
	header := l.theme.RenderHeader(appName, "ctx", v, l.termWidth)
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		l.paneZeroViewport.View(),
		l.paneOneViewport.View(),
		l.paneTwoViewport.View(),
	)
	footer := l.footerContent()

	return lipgloss.JoinVertical(lipgloss.Top, header, panes, footer)
}

func (l *Library) SetNotice(notice string, level styles.NoticeLevel) {
	if level == "" {
		level = styles.NoticeLevelInfo
	}
	l.noticeText = l.theme.RenderNotice(notice, level)
}

func (l *Library) paneZeroContent() string {
	var sb strings.Builder
	workspaces := l.visibleWorkspaces
	namespaces := l.visibleNamespaces
	sb.WriteString(renderPaneTitle("Workspaces", len(workspaces), l.currentPane == 0, l.theme))

	numWs := len(workspaces)
	numNs := len(namespaces)
	if numWs == 0 {
		sb.WriteString(l.theme.RenderError("No workspaces found"))
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
			sb.WriteString(renderSelection(prefix+truncateText(ws, paneWidth), l.theme))
			sb.WriteString("\n")
			if numNs == 1 {
				sb.WriteString(renderDescription(fmt.Sprintf("  %d namespace", numNs), l.theme))
			} else {
				sb.WriteString(renderDescription(fmt.Sprintf("  %d namespaces", numNs), l.theme))
			}
			sb.WriteString("\n")

			if l.showNamespaces {
				for j, ns := range namespaces {
					if uint(j) == l.currentNamespace {
						sb.WriteString(renderSecondarySelection("  > "+truncateText(ns, paneWidth), l.theme))
					} else {
						sb.WriteString(renderInactive("    "+truncateText(ns, paneWidth), l.theme))
					}
					sb.WriteString("\n")
				}
			}
		} else {
			sb.WriteString(renderInactive(prefix+truncateText(ws, paneWidth), l.theme))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (l *Library) paneOneContent() string {
	var sb strings.Builder
	sb.WriteString(renderPaneTitle("Executables", len(l.visibleExecutables), l.currentPane == 1, l.theme))
	if len(l.visibleExecutables) == 0 {
		sb.WriteString(l.theme.RenderError("No executables found"))
		return sb.String()
	}

	_, paneWidth, _ := calculateViewportWidths(l.termWidth)

	for i, ex := range l.visibleExecutables {
		if uint(i) == l.currentExecutable {
			sb.WriteString(renderSelection("* "+truncateText(ex.Ref().GetID(), paneWidth), l.theme))
		} else {
			sb.WriteString(renderInactive("  "+truncateText(ex.Ref().GetID(), paneWidth), l.theme))
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
	paneTwoMaxWidth := math.Floor(float64(maxWidth) * 0.95)
	mdStyles, err := l.theme.MarkdownStyleJSON()
	if err != nil {
		return l.theme.RenderError(fmt.Sprintf("unable to render markdown: %s", err.Error()))
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(mdStyles)),
		glamour.WithPreservedNewLines(),
		glamour.WithWordWrap(int(paneTwoMaxWidth)),
	)
	if err != nil {
		return l.theme.RenderError(fmt.Sprintf("unable to render markdown: %s", err.Error()))
	}

	ex := l.visibleExecutables[l.currentExecutable]
	content := ex.Markdown()
	viewStr, err := renderer.Render(content)
	if err != nil {
		return l.theme.RenderError(fmt.Sprintf("unable to render markdown: %s", err.Error()))
	}

	return viewStr
}

func (l *Library) footerContent() string {
	help := l.showHelp
	switch l.currentPane {
	case 0:
		if help && l.showNamespaces {
			return l.theme.RenderFooter(
				fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, paneOneExpandedHelp), l.termWidth,
			)
		} else if help {
			return l.theme.RenderFooter(fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, paneOneHelp), l.termWidth)
		} else if l.currentWorkspace < uint(len(l.visibleWorkspaces)) {
			ws := l.visibleWorkspaces[l.currentWorkspace]
			if ws == allWorkspacesLabel {
				break
			}
			var wsCfg *config.WorkspaceConfig
			for i, w := range l.allWorkspaces {
				if w.AssignedName() == ws {
					wsCfg = &l.allWorkspaces[i]
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
			return l.theme.RenderFooter(fmt.Sprintf("%s ● %s", footerPrefix, info), l.termWidth)
		}
	case 1, 2:
		if help {
			helpStr := paneTwoHelp
			if l.currentPane == 3 {
				helpStr = paneThreeHelp
			}

			return l.theme.RenderFooter(
				fmt.Sprintf("%s ● %s ● %s", footerPrefix, helpPrefix, helpStr), l.termWidth,
			)
		} else if l.currentExecutable < uint(len(l.visibleExecutables)) {
			exec := l.visibleExecutables[l.currentExecutable]
			path, err := relativePathFromWd(exec.DefinitionPath())
			if err != nil {
				l.ctx.Logger.Error(err, "unable to get relative path from wd")
				break
			}
			return l.theme.RenderFooter(fmt.Sprintf("%s ● %s", footerPrefix, path), l.termWidth)
		}
	}
	return l.theme.RenderFooter(footerPrefix, l.termWidth)
}

func relativePathFromWd(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Rel(wd, path)
}

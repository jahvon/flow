//nolint:gocritic
package library

import (
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/config"
)

func (l *Library) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	cmds = append(
		cmds,
		tea.SetWindowTitle("flow library"),
		tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
			return components.TickMsg(t)
		}),
	)
	cmds = append(
		cmds,
		l.paneZeroViewport.Init(),
		l.paneOneViewport.Init(),
		l.paneTwoViewport.Init(),
	)

	if l.ctx.InteractiveContainer.Width() >= 150 {
		l.splitView = true
	}
	l.setSize()
	go func() {
		l.setVisibleWorkspaces()
		l.setVisibleNamespaces()
		l.setVisibleExecs()
	}()

	return tea.Batch(cmds...)
}

func (l *Library) setVisibleExecs() {
	if len(l.allExecutables) == 0 || len(l.visibleWorkspaces) == 0 {
		return
	}

	curWs := l.filter.Workspace
	if label := l.visibleWorkspaces[l.currentWorkspace]; label != "" && label != allWorkspacesLabel {
		curWs = label
	} else if curWs == allWorkspacesLabel {
		curWs = ""
	}

	curNs := l.filter.Namespace
	if l.showNamespaces && len(l.visibleNamespaces) > 0 { //nolint:nestif
		if label := l.visibleNamespaces[l.currentNamespace]; label != "" {
			if label == withoutNamespaceLabel {
				curNs = ""
			} else if label == allNamespacesLabel {
				curNs = "*"
			} else {
				curNs = label
			}
		}
	}

	filter := l.filter
	filteredExec := l.allExecutables
	filteredExec = filteredExec.
		FilterByWorkspace(curWs).
		FilterByNamespace(curNs).
		FilterByVerb(filter.Verb).
		FilterByTags(filter.Tags).
		FilterBySubstring(filter.Substring)

	slices.SortFunc(filteredExec, func(i, j *config.Executable) int {
		return strings.Compare(i.Ref().String(), j.Ref().String())
	})
	l.visibleExecutables = filteredExec
}

func (l *Library) setVisibleWorkspaces() {
	if l.allWorkspaces == nil {
		return
	}

	filter := l.filter
	filteredWs := l.allWorkspaces
	if filter.Workspace != "" {
		for _, ws := range l.allWorkspaces {
			if ws.AssignedName() == filter.Workspace {
				filteredWs = config.WorkspaceConfigList{ws}
				break
			}
		}
	}

	var labels, prepend []string
	if len(filteredWs) > 1 {
		prepend = []string{allWorkspacesLabel}
	}
	for _, ws := range filteredWs {
		labels = append(labels, ws.AssignedName())
	}
	slices.Sort(labels)
	l.visibleWorkspaces = append(prepend, labels...) //nolint:gocritic
}

func (l *Library) setVisibleNamespaces() {
	if l.allExecutables == nil || len(l.visibleWorkspaces) == 0 {
		return
	}

	var labels, prepend []string
	var someWithoutNs bool
	filter := l.filter
	filterWs := l.visibleWorkspaces[l.currentWorkspace]
	nsSet := make(map[string]struct{})
	for _, ex := range l.allExecutables {
		ns := ex.Ref().GetNamespace()
		ws := ex.Ref().GetWorkspace()
		if filter.Namespace != "*" && filter.Namespace != "" && ns != filter.Namespace {
			continue
		} else if filterWs != allWorkspacesLabel && filterWs != "" && ws != filterWs {
			continue
		} else if ns == "" || ns == "*" {
			someWithoutNs = true
			continue
		}

		if _, ok := nsSet[ns]; !ok {
			nsSet[ns] = struct{}{}
			labels = append(labels, ns)
		}
	}
	slices.Sort(labels)
	if len(labels) > 1 {
		prepend = append(prepend, allNamespacesLabel)
	}
	if someWithoutNs {
		prepend = append(prepend, withoutNamespaceLabel)
	}
	l.visibleNamespaces = append(prepend, labels...) //nolint:gocritic
}

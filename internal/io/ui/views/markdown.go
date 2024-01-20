package views

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

type GenericMarkdownView struct {
	parent   types.ParentView
	viewport viewport.Model
	content  string
}

func NewGenericMarkdownView(parent types.ParentView, content string) types.ViewBuilder {
	vp := viewport.New(parent.Width(), parent.Height())
	vp.Style = styles.EntityViewStyle().Width(parent.Width())
	return &GenericMarkdownView{
		parent:   parent,
		content:  content,
		viewport: vp,
	}
}

func (v *GenericMarkdownView) Init() tea.Cmd {
	return nil
}

func (v *GenericMarkdownView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

func (v *GenericMarkdownView) View() string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(styles.MarkdownStyleJSON)),
		glamour.WithWordWrap(v.parent.Width()),
	)
	if err != nil {
		v.parent.HandleInternalError(err)
	}

	viewStr, err := renderer.Render(v.content)
	if err != nil {
		v.parent.HandleInternalError(err)
	}
	v.viewport.SetContent(viewStr)
	return v.viewport.View()
}

func (v *GenericMarkdownView) HelpMsg() string {
	return ""
}

func (v *GenericMarkdownView) FooterEnabled() bool {
	return true
}

func (v *GenericMarkdownView) Type() string {
	return "markdown"
}

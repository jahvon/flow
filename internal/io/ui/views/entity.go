package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

type EntityView struct {
	parent    types.ParentView
	viewport  viewport.Model
	entity    config.Entity
	format    config.OutputFormat
	callbacks []types.KeyCallback
}

func NewEntityView(
	parent types.ParentView,
	entity config.Entity,
	format config.OutputFormat,
	keys ...types.KeyCallback,
) types.ViewBuilder {
	if format == "" {
		format = config.UNSET
	}
	vp := viewport.New(parent.Width(), parent.Height())
	vp.Style = styles.EntityViewStyle().Width(parent.Width())
	return &EntityView{
		parent:    parent,
		entity:    entity,
		format:    format,
		callbacks: keys,
		viewport:  vp,
	}
}

func (v *EntityView) Init() tea.Cmd {
	return nil
}

//nolint:gocognit
func (v *EntityView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if v.parent.Ready() {
			v.viewport.Width = v.parent.Width()
			v.viewport.Height = v.parent.Height()
			v.viewport.SetContent(v.renderedContent())
		}
	case tea.KeyMsg:
		if !v.parent.Ready() {
			return v, nil
		}

		switch msg.String() {
		case "", "-":
			if v.format == config.UNSET {
				return v, nil
			}
			v.format = config.UNSET
			v.viewport.GotoTop()
		case "y":
			if v.format == config.YAML {
				return v, nil
			}
			v.format = config.YAML
			v.viewport.GotoTop()
		case "j":
			if v.format == config.JSON {
				return v, nil
			}
			v.format = config.JSON
			v.viewport.GotoTop()
		case "f":
			if v.format == config.FormattedJSON {
				return v, nil
			}
			v.format = config.FormattedJSON
			v.viewport.GotoTop()
		case "up":
			v.viewport.LineUp(1)
		case "down":
			v.viewport.LineDown(1)
		default:
			for _, cb := range v.callbacks {
				if cb.Key == msg.String() {
					if err := cb.Callback(); err != nil {
						v.parent.HandleInternalError(err)
					}
				}
			}
		}
	}

	return v, cmd
}

func (v *EntityView) renderedContent() string {
	var content string
	var err error
	switch v.format {
	case config.YAML:
		content, err = v.entity.YAML()
		content = fmt.Sprintf("```yaml\n%s\n```", content)
	case config.JSON:
		content, err = v.entity.JSON(false)
		content = fmt.Sprintf("```json\n%s\n```", content)
	case config.FormattedJSON:
		content, err = v.entity.JSON(true)
		content = fmt.Sprintf("```json\n%s\n```", content)
	case config.UNSET:
		fallthrough
	default:
		content = v.entity.Markdown()
	}
	if err != nil {
		v.parent.HandleInternalError(err)
	}
	if content == "" {
		content = "no data"
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(styles.MarkdownStyleJSON)),
		glamour.WithWordWrap(v.parent.Width()),
	)
	if err != nil {
		v.parent.HandleInternalError(err)
	}

	viewStr, err := renderer.Render(content)
	if err != nil {
		v.parent.HandleInternalError(err)
	}
	return viewStr
}

func (v *EntityView) View() string {
	v.viewport.SetContent(v.renderedContent())
	return v.viewport.View()
}

func (v *EntityView) HelpMsg() string {
	msg := "y: yaml • j: json • f: formatted json"

	var extendedHelp string
	for i, cb := range v.callbacks {
		switch {
		case cb.Key == "" || cb.Label == "":
			continue
		case i == 0:
			extendedHelp += fmt.Sprintf("%s: %s", cb.Key, cb.Label)
		default:
			extendedHelp += fmt.Sprintf(" • %s: %s", cb.Key, cb.Label)
		}
	}
	if extendedHelp != "" {
		msg = fmt.Sprintf("%s | %s", extendedHelp, msg)
	}
	return msg
}

func (v *EntityView) FooterEnabled() bool {
	return true
}

func (v *EntityView) Type() string {
	return "entity"
}

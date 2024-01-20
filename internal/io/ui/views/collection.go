package views

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

type CollectionViewItem struct {
	title, desc, filterVal string
}

func (i *CollectionViewItem) Title() string       { return i.title }
func (i *CollectionViewItem) Description() string { return i.desc }
func (i *CollectionViewItem) FilterValue() string { return i.filterVal }

func listItemDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles = styles.ListItemStyles()
	return delegate
}

type CollectionView struct {
	app          types.ParentView
	collection   config.Collection
	model        *list.Model
	items        []list.Item
	format       config.OutputFormat
	callbacks    []types.KeyCallback
	selectedFunc func(filterVal string) error
}

func NewCollectionView(
	app types.ParentView,
	collection config.Collection,
	format config.OutputFormat,
	selectedFunc func(filterVal string) error,
	keys ...types.KeyCallback,
) types.ViewBuilder {
	if format == "" {
		format = config.UNSET
	}
	items := make([]list.Item, 0)
	for _, item := range collection.Items() {
		items = append(items, &CollectionViewItem{
			title:     collectionItemTitle(item),
			desc:      collectionItemDesc(item),
			filterVal: item.Header,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].FilterValue() < items[j].FilterValue()
	})
	model := list.New(items, listItemDelegate(), app.Width(), app.Height())
	model.SetShowTitle(false)
	model.SetShowHelp(false)
	model.SetShowPagination(false)
	model.SetStatusBarItemName(collection.Singular(), collection.Plural())
	model.Styles = styles.ListStyles()
	return &CollectionView{
		app:          app,
		collection:   collection,
		model:        &model,
		items:        items,
		format:       format,
		selectedFunc: selectedFunc,
		callbacks:    keys,
	}
}

func (v *CollectionView) Init() tea.Cmd {
	return nil
}

//nolint:gocognit
func (v *CollectionView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if v.app.Ready() {
			v.model.SetSize(v.app.Height(), v.app.Width())
		}
	case tea.KeyMsg:
		if !v.app.Ready() {
			return v, nil
		}

		switch msg.String() {
		case "", "-":
			if v.format == config.UNSET {
				return v, nil
			}
			v.format = config.UNSET
		case "y":
			if v.format == config.YAML {
				return v, nil
			}
			v.format = config.YAML
		case "j":
			if v.format == config.JSON {
				return v, nil
			}
			v.format = config.JSON
		case "f":
			if v.format == config.FormattedJSON {
				return v, nil
			}
			v.format = config.FormattedJSON
		case tea.KeyEnter.String():
			if v.selectedFunc == nil {
				return v, nil
			}
			selected := v.model.SelectedItem()
			if selected == nil {
				return v, nil
			}

			if err := v.selectedFunc(selected.FilterValue()); err != nil {
				v.app.HandleInternalError(err)
			}
			return v, nil
		default:
			for _, cb := range v.callbacks {
				if cb.Key == msg.String() {
					if err := cb.Callback(); err != nil {
						v.app.HandleInternalError(err)
					}
					return v, nil
				}
			}
		}
	}

	model, cmd := v.model.Update(msg)
	v.model = &model
	return v, cmd
}

func (v *CollectionView) UpdateItemsFromCollections() {
	items := make([]list.Item, 0)
	for _, item := range v.collection.Items() {
		items = append(items, &CollectionViewItem{
			title:     collectionItemTitle(item),
			desc:      collectionItemDesc(item),
			filterVal: item.Header,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].FilterValue() < items[j].FilterValue()
	})
	v.items = items
}

func (v *CollectionView) Items() []list.Item {
	return v.model.Items()
}

func (v *CollectionView) renderedContent() string {
	var content string
	var isMkdwn bool
	var err error
	switch v.format {
	case config.YAML:
		content, err = v.collection.YAML()
		content = fmt.Sprintf("```yaml\n%s\n```", content)
		isMkdwn = true
	case config.JSON:
		content, err = v.collection.JSON(false)
		content = fmt.Sprintf("```json\n%s\n```", content)
		isMkdwn = true
	case config.FormattedJSON:
		content, err = v.collection.JSON(true)
		content = fmt.Sprintf("```json\n%s\n```", content)
		isMkdwn = true
	case config.UNSET:
		fallthrough
	default:
		v.model.SetSize(v.app.Width(), v.app.Height())
		v.UpdateItemsFromCollections()
		style := styles.CollectionStyle().Width(v.app.Width())
		content = style.Render(v.model.View())
	}
	if err != nil {
		v.app.HandleInternalError(err)
	}
	if content == "" {
		content = "no data"
	}

	if !isMkdwn {
		return content
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(styles.MarkdownStyleJSON)),
		glamour.WithWordWrap(v.app.Width()),
	)
	if err != nil {
		v.app.HandleInternalError(err)
	}

	viewStr, err := renderer.Render(content)
	if err != nil {
		v.app.HandleInternalError(err)
	}
	return viewStr
}

func (v *CollectionView) View() string {
	return v.renderedContent()
}

func (v *CollectionView) HelpMsg() string {
	var selectHelp string
	if v.selectedFunc != nil {
		selectHelp = "enter: select • "
	}
	msg := fmt.Sprintf("%s/: filter | y: yaml • j: json • f: formatted json", selectHelp)

	var extendedHelp string
	for _, cb := range v.callbacks {
		if cb.Key == "" || cb.Label == "" {
			continue
		}
		extendedHelp += fmt.Sprintf(" • %s: %s", cb.Key, cb.Label)
	}
	if extendedHelp != "" {
		msg = fmt.Sprintf("%s | %s", extendedHelp, msg)
	}
	return msg
}

func (v *CollectionView) FooterEnabled() bool {
	return true
}

func (v *CollectionView) Type() string {
	return "collection"
}

func collectionItemTitle(item config.CollectionItem) string {
	title := item.Header
	if item.SubHeader != "" {
		title += fmt.Sprintf(" (%s)", item.SubHeader)
	}
	return title
}

func collectionItemDesc(item config.CollectionItem) string {
	desc := item.Description
	if len(item.Tags) > 0 {
		desc = fmt.Sprintf("[%s] ", item.Tags.PreviewString()) + desc
	}
	return desc
}

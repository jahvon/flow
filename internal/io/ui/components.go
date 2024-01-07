package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

type NoticeLevel string

const (
	NoticeLevelInfo    NoticeLevel = "info"
	NoticeLevelWarning NoticeLevel = "warning"
	NoticeLevelError   NoticeLevel = "error"

	loadingViewType = "loading"
)

var (
	log = io.Log().With().Str("scope", "io/ui").Logger()
)

type KeyCallback struct {
	Key      string
	Label    string
	Callback func() error
}

func ContextStr(ws, ns string) string {
	if ws == "" {
		ws = unknownStyle("unk")
	}
	if ns == "" {
		ns = "*"
	}
	return contextStyle(fmt.Sprintf("ctx: %s/%s ", ws, ns))
}

func NoticeStr(notice string, lvl NoticeLevel) string {
	if notice == "" {
		return ""
	}

	padded := lipgloss.Style{}.MarginLeft(2).Render
	switch lvl {
	case NoticeLevelInfo:
		return padded(infoStyle("INFO:" + notice))
	case NoticeLevelWarning:
		return padded(warningStyle("WARN: " + notice))
	case NoticeLevelError:
		return padded(errorStyle("ERR: " + notice))
	default:
		return padded(unknownStyle(notice))
	}
}

func FilterStr(tags config.Tags) string {
	if len(tags) == 0 {
		return filterStyle("filter: *")
	}
	return filterStyle(fmt.Sprintf("filter: %s ", tags.ContextString()))
}

type LoadingView struct {
	msg     string
	spinner spinner.Model
}

func (v *LoadingView) Init() tea.Cmd {
	return v.spinner.Tick
}

func (v *LoadingView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		v.msg = msg.Error()
	case string:
		v.msg = msg
	}
	v.spinner, cmd = v.spinner.Update(msg)
	return v, cmd
}

func (v *LoadingView) View() string {
	var txt string
	if v.msg == "" {
		txt = fmt.Sprintf("\n\n %s %s\n\n", v.spinner.View(), infoStyle("thinking..."))
	} else {
		txt = fmt.Sprintf("\n\n %s %s\n\n", v.spinner.View(), infoStyle(v.msg))
	}
	return txt
}

func (v *LoadingView) HelpMsg() string {
	return ""
}

func (v *LoadingView) FooterEnabled() bool {
	return false
}

func (v *LoadingView) Type() string {
	return loadingViewType
}

func NewLoadingView(msg string) ViewBuilder {
	spin := spinner.New()
	spin.Style = spinnerStyle
	spin.Spinner = spinner.Points
	return &LoadingView{
		msg:     msg,
		spinner: spin,
	}
}

type EntityView struct {
	app       *Application
	viewport  viewport.Model
	entity    config.Entity
	format    config.OutputFormat
	callbacks []KeyCallback
}

func NewEntityView(
	app *Application,
	entity config.Entity,
	format config.OutputFormat,
	keys ...KeyCallback,
) ViewBuilder {
	if format == "" {
		format = config.UNSET
	}
	vp := viewport.New(app.Width(), app.Height())
	vp.Style = entityViewStyle
	return &EntityView{
		app:       app,
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
		if v.app.Ready() {
			v.viewport.Width = v.app.Width()
			v.viewport.Height = v.app.Height()
			v.viewport.SetContent(v.renderedContent())
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
						v.app.HandleInternalError(err)
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
		v.app.HandleInternalError(err)
	}
	if content == "" {
		content = "no data"
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes([]byte(markdownStyleJSON)),
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

type CollectionViewItem struct {
	title, desc, filterVal string
}

func (i *CollectionViewItem) Title() string       { return i.title }
func (i *CollectionViewItem) Description() string { return i.desc }
func (i *CollectionViewItem) FilterValue() string { return i.filterVal }

func listItemDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles = listItemStyles()
	return delegate
}

type CollectionView struct {
	app          *Application
	collection   config.Collection
	model        *list.Model
	items        []list.Item
	format       config.OutputFormat
	callbacks    []KeyCallback
	selectedFunc func(filterVal string) error
}

func NewCollectionView(
	app *Application,
	collection config.Collection,
	format config.OutputFormat,
	selectedFunc func(filterVal string) error,
	keys ...KeyCallback,
) ViewBuilder {
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
	model.Styles = listStyles()
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
		case "enter":
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
		content = collectionStyle.Render(v.model.View())
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
		glamour.WithStylesFromJSONBytes([]byte(markdownStyleJSON)),
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

package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io/ui"
	"github.com/jahvon/flow/internal/services/open"
)

type WorkspaceView struct {
	curCtx  *context.Context
	subject *config.WorkspaceConfig
	format  config.OutputFormat
}

func NewWorkspaceView(ctx context.Context, ws *config.WorkspaceConfig, format config.OutputFormat) ViewBuilder {
	if format == "" {
		format = config.INTERACTIVE
	}
	if ws == nil {
		ws = ctx.CurrentWorkspace
	}
	return &WorkspaceView{
		curCtx:  &ctx,
		subject: ws,
		format:  format,
	}
}

func (v *WorkspaceView) Title() string {
	return fmt.Sprintf("workspace: %s", v.subject.AssignedName())
}

func (v *WorkspaceView) Help() string {
	help := `
		i: interactive
		y: yaml
		j: json
		n: indented json
		
		h: help
		esc: back
		ctrl+c: quit
	`
	return fmt.Sprintf(help)
}

func (v *WorkspaceView) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'i':
			if v.format == config.INTERACTIVE {
				break
			}
			v.format = config.INTERACTIVE
		case 'y':
			if v.format == config.YAML {
				break
			}
			v.format = config.YAML
		case 'j':
			if v.format == config.JSON {
				break
			}
			v.format = config.JSON
		case 'v':
			if v.format == config.JSONP {
				break
			}
			v.format = config.JSONP
		case 'h':
			v.curCtx.App.ShowModal(v.Help())
			break
		}

		view, err := v.View()
		if err != nil {
			v.curCtx.App.SetErrorPage(err)
		}
		v.curCtx.App.SetPage(v.Title(), view)

		return nil
	}
	return event
}

func (v *WorkspaceView) View() (tview.Primitive, error) {
	var wsView tview.Primitive
	switch v.format {
	case config.JSON:
		json, err := v.subject.JSON(false)
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(json)
	case config.JSONP:
		json, err := v.subject.JSON(true)
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(json)
	case config.YAML:
		yaml, err := v.subject.YAML()
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(yaml)
	case config.INTERACTIVE:
		buttons := []ui.ViewButton{v.launchWsButton(), v.editWsButton()}
		if v.subject.AssignedName() != v.curCtx.UserConfig.CurrentWorkspace {
			buttons = append(buttons, v.setWsButton())
		}
		wsView = ui.DetailsView(v.subject.DetailsString(), buttons...)
	default:
		return nil, fmt.Errorf("unknown output format: %s", v.format)
	}
	return wsView, nil
}

func (v *WorkspaceView) setWsButton() ui.ViewButton {
	return ui.ViewButton{
		Label: "Set",
		Func: func() {
			v.curCtx.UserConfig.CurrentWorkspace = v.subject.AssignedName()
			if err := file.WriteUserConfig(v.curCtx.UserConfig); err != nil {
				log.Err(err).Msg("failed to save user config")
			}
			v.curCtx.App.SetHeader(
				ui.WithCurrentNamespace(v.curCtx.UserConfig.CurrentNamespace),
				ui.WithCurrentWorkspace(v.curCtx.UserConfig.CurrentWorkspace),
			)
		},
	}
}

func (v *WorkspaceView) launchWsButton() ui.ViewButton {
	return ui.ViewButton{
		Label: "Open",
		Func: func() {
			if err := open.Open(v.subject.Location(), false); err != nil {
				log.Err(err).Msg("unable to open workspace")
			}
		},
	}
}

func (v *WorkspaceView) editWsButton() ui.ViewButton {
	return ui.ViewButton{
		Label: "Edit",
		Func: func() {
			app := v.curCtx.UserConfig.AppPreferences.Edit
			if app != "" {
				if err := open.OpenWith(app, v.subject.Location(), false); err != nil {
					log.Err(err).Msg("unable to open workspace")
				}
			} else {
				v.curCtx.App.Suspend(func() {
					cmd := exec.Command("vim", filepath.Join(v.subject.Location(), file.WorkspaceConfigFileName))
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					err := cmd.Run()
					if err != nil {
						log.Err(err).Msg("unable to open vim")
					}
				})
			}
		},
	}
}

type WorkspaceListView struct {
	curCtx     *context.Context
	format     config.OutputFormat
	workspaces config.WorkspaceConfigList
}

func NewWorkspaceListView(
	ctx context.Context,
	workspaces config.WorkspaceConfigList,
	format config.OutputFormat,
) ViewBuilder {
	if format == "" {
		format = config.INTERACTIVE
	}
	return &WorkspaceListView{
		curCtx:     &ctx,
		format:     format,
		workspaces: workspaces,
	}
}

func (v *WorkspaceListView) Title() string {
	return "workspaces"
}

func (v *WorkspaceListView) Help() string {
	help := `
		i: interactive
		y: yaml
		j: json
		n: indented json
		
		h: help
		esc: back
		ctrl+c: quit
	`
	return fmt.Sprintf(help)
}

func (v *WorkspaceListView) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'i':
			if v.format == config.INTERACTIVE {
				break
			}
			v.format = config.INTERACTIVE
		case 'y':
			if v.format == config.YAML {
				break
			}
			v.format = config.YAML
		case 'j':
			if v.format == config.JSON {
				break
			}
			v.format = config.JSON
		case 'n':
			if v.format == config.JSONP {
				break
			}
			v.format = config.JSONP
		case 'h':
			v.curCtx.App.ShowModal(v.Help())
			break
		}

		view, err := v.View()
		if err != nil {
			v.curCtx.App.SetErrorPage(err)
		}
		v.curCtx.App.SetPage(v.Title(), view)

		return nil
	}
	return event
}

func (v *WorkspaceListView) View() (tview.Primitive, error) {
	var wsView tview.Primitive
	switch v.format {
	case config.JSON:
		json, err := v.workspaces.JSON(false)
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(json)
	case config.JSONP:
		json, err := v.workspaces.JSON(true)
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(json)
	case config.YAML:
		yaml, err := v.workspaces.YAML()
		if err != nil {
			return nil, err
		}
		wsView = ui.TextView(yaml)
	case config.INTERACTIVE:
		header, rows := v.workspaces.TableData()
		wsView = ui.TableView(header, rows, v.selectRowFunc())
	default:
		return nil, fmt.Errorf("unknown output format: %s", v.format)
	}
	return wsView, nil
}

func (v *WorkspaceListView) selectRowFunc() func(rowData []string) {
	return func(rowData []string) {
		if len(rowData) == 0 {
			log.Err(fmt.Errorf("row data is empty"))
			return
		}
		wsName := rowData[0]
		ws, found := lo.Find(v.workspaces, func(s config.WorkspaceConfig) bool {
			if s.AssignedName() == wsName || s.DisplayName == wsName {
				return true
			}
			return false
		})
		if !found {
			log.Err(fmt.Errorf("workspace '%s' not found", wsName))
			return
		}
		wsView := NewWorkspaceView(*v.curCtx, &ws, v.format)
		err := BuildAndSetView(v.curCtx.App, wsView)
		if err != nil {
			v.curCtx.App.SetErrorPage(err)
		}
	}
}

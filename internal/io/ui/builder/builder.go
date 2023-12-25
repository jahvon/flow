package builder

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/jahvon/flow/internal/io/ui"
)

type ViewBuilder interface {
	Title() string
	View() (tview.Primitive, error)
	Help() string
	HandleInput(event *tcell.EventKey) *tcell.EventKey
}

func BuildAndSetView(app *ui.Application, builder ViewBuilder, headerOpts ...ui.HeaderOption) error {
	if app == nil {
		return nil
	}
	if builder == nil {
		return nil
	}
	if len(headerOpts) > 0 {
		app.SetHeader(headerOpts...)
	}
	view, err := builder.View()
	if err != nil {
		return err
	}
	app.SetPage(builder.Title(), view)
	app.RegisterKeyHandler(builder.Title(), builder.HandleInput)
	return nil
}

package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DefaultView = WorkspaceView

type View string

const (
	WorkspaceView  View = "workspace"
	NamespaceView  View = "namespace"
	ExecutableView View = "executable"
)

type TableData struct {
	Headers []string
	Rows    [][]string
}

func (t TableData) TViewTable() *tview.Table {
	table := tview.NewTable().
		InsertRow(len(t.Rows)).
		InsertColumn(len(t.Headers))

	for i, header := range t.Headers {
		cell := tview.NewTableCell(strings.ToUpper(strWithPadding(header))).
			SetBackgroundColor(tcell.ColorRed).
			SetStyle(
				tcell.StyleDefault.
					Bold(true).
					Background(tcell.ColorLightGray).
					Foreground(tcell.ColorBlack),
			).
			SetSelectable(false).
			SetExpansion(expansionMap(header))
		table.SetCell(0, i, cell)
	}

	for i, row := range t.Rows {
		for j, val := range row {
			cell := tview.NewTableCell(strWithPadding(val))
			if i%2 == 0 {
				cell = cell.SetTransparency(true).SetTextColor(tcell.ColorWhite)
			} else {
				cell = cell.SetBackgroundColor(tcell.ColorGray).SetTextColor(tcell.ColorWhite)
			}
			table.SetCell(i+1, j, cell).SetSelectable(false, false)
		}
	}
	table.SetBorderPadding(0, 1, 1, 1)
	table.SetBorder(false)
	return table
}

func strWithPadding(str string) string {
	return fmt.Sprintf(" %s ", str)
}

func expansionMap(str string) int {
	str = strings.ToLower(str)
	switch str {
	case "key", "name", "id":
		return 1
	case "value", "description":
		return 2
	default:
		return 0
	}
}

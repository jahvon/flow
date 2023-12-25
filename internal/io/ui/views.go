package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TextView(text string) tview.Primitive {
	textView := tview.NewTextView().
		SetText(text).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft)
	return textView
}

type ViewButton struct {
	Label string
	Func  func()
}

func DetailsView(text string, buttons ...ViewButton) tview.Primitive {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).SetBorderPadding(1, 1, 0, 0)
	if len(buttons) > 0 {
		buttonGroup := buttonGroupView(buttons...)
		flex.AddItem(buttonGroup, 2, 0, true)
	}
	textView := tview.NewTextView().
		SetText(text).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true).
		SetWordWrap(true)
	textView.SetBorderPadding(0, 0, 2, 3)
	flex.AddItem(textView, 0, 1, false)
	return flex
}

func buttonView(label string, action func()) *tview.Button {
	button := tview.NewButton(label).SetSelectedFunc(action)
	button.SetLabelColor(IdleState.PrimaryFGColor()).
		SetBackgroundColor(IdleState.PrimaryBGColor())
	button.SetLabelColorActivated(IdleState.PrimaryFGColor()).
		SetBackgroundColorActivated(IdleState.PrimaryBGColor())
	return button
}

func buttonGroupView(buttons ...ViewButton) tview.Primitive {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn).SetBorderPadding(0, 1, 0, 0)
	flex.AddItem(nil, 0, len(buttons)*3, false)
	for i, button := range buttons {
		flex.AddItem(buttonView(button.Label, button.Func), 0, 1, false)
		if i < len(buttons)-1 {
			flex.AddItem(nil, 2, 0, false)
		}
	}
	flex.AddItem(nil, 0, len(buttons)*3, false)
	return flex
}

func TableView(headers []string, rows [][]string, selectFunc func(rowData []string)) tview.Primitive {
	tableData := TableData{
		Headers: headers,
		Rows:    rows,
	}
	return tableData.TViewTable(selectFunc)
}

func LoadingView() tview.Primitive {
	textView := tview.NewTextView().
		SetText("\nLoading...").
		SetTextColor(ProgressingState.SecondaryFGColor()).
		SetTextAlign(tview.AlignCenter)
	return textView
}

func ErrorView(err error) tview.Primitive {
	textView := tview.NewTextView().
		SetText(fmt.Sprintf("\nEncountered error: %s", err.Error())).
		SetTextColor(ErrorState.SecondaryFGColor()).
		SetTextAlign(tview.AlignCenter)
	return textView
}

func HelpView(txt string) tview.Primitive {
	textView := tview.NewTextView().
		SetText(txt).
		SetTextColor(IdleState.SecondaryFGColor()).
		SetTextAlign(tview.AlignLeft)
	return textView
}

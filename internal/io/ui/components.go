package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

func brandTxt(state State) *tview.TextView {
	blink := state == ProgressingState
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextStyle(
			tcell.StyleDefault.
				Bold(true).
				Italic(true).
				Blink(blink).
				Background(state.PrimaryBGColor()).
				Foreground(state.PrimaryFGColor()),
		).
		SetScrollable(false).
		SetText("flow")
}

func contextTxt(workspace, namespace string, state State) *tview.TextView {
	if workspace == "" {
		workspace = "*"
	}
	if namespace == "" {
		namespace = "*"
	}
	txt := fmt.Sprintf("[ ctx: %s/%s ]", workspace, namespace)

	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextStyle(
			tcell.StyleDefault.
				Bold(true).
				Background(state.SecondaryBGColor()).
				Foreground(state.SecondaryFGColor()),
		).
		SetScrollable(false).
		SetText(txt)
}

func filterTxt(filter string, state State) *tview.TextView {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		filter = "*"
	}
	txt := fmt.Sprintf("[ filter: %s ]", filter)

	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextStyle(
			tcell.StyleDefault.
				Bold(true).
				Background(state.SecondaryBGColor()).
				Foreground(state.SecondaryFGColor()),
		).
		SetScrollable(false).
		SetWrap(false).
		SetText(txt)
}

func noticeTxt(notice string, state State) *tview.TextView {
	notice = strings.TrimSpace(notice)
	if notice == "" {
		return nil
	}

	return tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetTextStyle(
			tcell.StyleDefault.
				Italic(true).
				Background(state.SecondaryBGColor()).
				Foreground(tcell.ColorBlack),
		).
		SetText(notice + " ")
}

func contentBox(title string, opts FrameOptions) *tview.Flex {
	flex := tview.NewFlex()
	if opts.ObjectContent != nil {
		flex.SetDirection(tview.FlexRow).
			AddItem(opts.ObjectContent.TViewTable(), 0, 1, false)

	} else if opts.ObjectList != nil {
		flex.SetDirection(tview.FlexRow).
			AddItem(opts.ObjectList.TViewTable(), 0, 1, false)
	} else {
		flex.SetDirection(tview.FlexColumn).
			AddItem(tview.NewTextView().
				SetText("\ninternal error").
				SetTextColor(tcell.ColorRed).
				SetTextAlign(tview.AlignCenter),
				0, 1, false)
	}

	flex.SetBorder(true)
	flex.SetBorderColor(tcell.ColorWhite)
	flex.SetTitle(title)
	flex.SetTitleColor(opts.CurrentState.PrimaryFGColor())
	flex.SetTitleAlign(tview.AlignCenter)

	return flex
}

func emptyBox() *tview.Box {
	return tview.NewBox().SetBorder(false)
}

func textViewWidth(tv *tview.TextView) int {
	txt := tv.GetText(false)
	return runewidth.StringWidth(txt)
}
